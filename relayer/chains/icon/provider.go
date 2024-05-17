package icon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/module"
	"go.uber.org/zap"
)

type Config struct {
	ChainName     string                         `json:"-" yaml:"-"`
	RPCUrl        string                         `json:"rpc-url" yaml:"rpc-url"`
	Address       string                         `json:"address" yaml:"address"`
	StartHeight   uint64                         `json:"start-height" yaml:"start-height"` // would be of highest priority
	Contracts     relayerTypes.ContractConfigMap `json:"contracts" yaml:"contracts"`
	NetworkID     int64                          `json:"network-id" yaml:"network-id"`
	FinalityBlock uint64                         `json:"finality-block" yaml:"finality-block"`
	NID           string                         `json:"nid" yaml:"nid"`
	StepMin       int64                          `json:"step-min" yaml:"step-min"`
	StepLimit     int64                          `json:"step-limit" yaml:"step-limit"`
	StepBuffer    int                            `json:"step-buffer" yaml:"step-buffer"`
	HomeDir       string                         `json:"-" yaml:"-"`
}

// NewProvider returns new Icon provider
func (c *Config) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := c.sanitize(); err != nil {
		return nil, err
	}

	c.ChainName = chainName
	c.HomeDir = homepath

	return &Provider{
		log:       log.With(zap.Stringp("nid ", &c.NID), zap.Stringp("name", &c.ChainName)),
		client:    NewClient(ctx, c.RPCUrl, log),
		cfg:       c,
		contracts: c.eventMap(),
	}, nil
}

func (c *Config) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	if err := c.Contracts.Validate(); err != nil {
		return fmt.Errorf("contracts are not valid: %s", err)
	}

	// TODO: validation for keystore
	// TODO: contractaddress validation
	// TODO: account should have some balance no balance then use another accoutn

	return nil
}

func (c *Config) sanitize() error {
	if c.StepBuffer == 0 {
		c.StepBuffer = 10
	}
	return nil
}

func (p *Config) SetWallet(addr string) {
	p.Address = addr
}

func (p *Config) GetWallet() string {
	return p.Address
}

type Provider struct {
	log       *zap.Logger
	cfg       *Config
	wallet    module.Wallet
	client    *Client
	kms       kms.KMS
	contracts map[string]providerTypes.EventMap
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Init(ctx context.Context, homepath string, kms kms.KMS) error {
	p.kms = kms
	return nil
}

func (p *Provider) Type() string {
	return "icon"
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Wallet() (module.Wallet, error) {
	if p.wallet == nil {
		if err := p.RestoreKeystore(context.Background()); err != nil {
			return nil, err
		}
	}
	return p.wallet, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

// MessageReceived checks if the message is received
func (p *Provider) MessageReceived(ctx context.Context, messageKey *providerTypes.MessageKey) (bool, error) {
	callParam := p.prepareCallParams(MethodGetReceipts, p.cfg.Contracts[providerTypes.ConnectionContract], map[string]interface{}{
		"srcNetwork": messageKey.Src,
		"_connSn":    types.NewHexInt(int64(messageKey.Sn)),
	})

	var status types.HexInt
	if err := p.client.Call(callParam, &status); err != nil {
		return false, fmt.Errorf("MessageReceived: %v", err)
	}
	return status == types.NewHexInt(1), nil
}

// ReverseMessage reverts a message
func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	params := map[string]interface{}{"_sn": types.NewHexInt(sn.Int64())}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), params, MethodRevertMessage)
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}

// SetAdmin sets the admin address of the bridge contract
func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	callParam := map[string]interface{}{
		"_relayer": admin,
	}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), callParam, MethodSetAdmin)

	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return fmt.Errorf("SetAdmin: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.HexBytes(txHash)})
	if err != nil {
		return fmt.Errorf("SetAdmin: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("SetAdmin: failed to set admin: %s", txr.TxHash)
	}
	return nil
}

// GetFee
func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	callParam := p.prepareCallParams(MethodGetFee, p.cfg.Contracts[providerTypes.ConnectionContract], map[string]interface{}{
		"to":       networkID,
		"response": types.NewHexInt(1),
	})

	var status types.HexInt
	if err := p.client.Call(callParam, &status); err != nil {
		return 0, fmt.Errorf("GetFee: %v", err)
	}
	fee, err := status.BigInt()
	if err != nil {
		return 0, fmt.Errorf("GetFee: %v", err)
	}
	return fee.Uint64(), nil
}

// SetFees
func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	callParam := map[string]interface{}{
		"networkId":   networkID,
		"messageFee":  types.NewHexInt(int64(msgFee)),
		"responseFee": types.NewHexInt(int64(resFee)),
	}

	msg := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), callParam, MethodSetFee)
	txHash, err := p.SendTransaction(ctx, msg)
	if err != nil {
		return fmt.Errorf("SetFee: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		fmt.Println("SetFee: WaitForResults: %v", err)
		return fmt.Errorf("SetFee: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("SetFee: failed to claim fees: %s", txr.TxHash)
	}
	return nil
}

// ClaimFees
func (p *Provider) ClaimFee(ctx context.Context) error {
	msg := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), map[string]interface{}{}, MethodClaimFees)
	txHash, err := p.SendTransaction(ctx, msg)
	if err != nil {
		return fmt.Errorf("ClaimFees: %v", err)
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return fmt.Errorf("ClaimFees: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("ClaimFees: failed to claim fees: %s", txr.TxHash)
	}
	return nil
}

// ExecuteRollback
func (p *Provider) ExecuteRollback(ctx context.Context, sn uint64) error {
	params := map[string]interface{}{"_sn": types.NewHexInt(int64(sn))}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.XcallContract]), params, MethodExecuteRollback)
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}
