package bitcoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/studyzy/runestone"
	"lukechampine.com/uint128"

	"path/filepath"

	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/icon-project/goloop/common/codec"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/icon-project/centralized-relay/relayer/chains/bitcoin/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"github.com/ethereum/go-ethereum/common"
)

var (
	BTCToken           = "0:0"
	MethodDeposit      = "Deposit"
	MethodWithdrawTo   = "WithdrawTo"
	MethodRefundTo     = "RefundTo"
	MethodRollback     = "Rollback"
	MasterMode         = "master"
	SlaveMode          = "slave"
	BtcDB              = "btc.db"
	WitnessSize        = 380
	NumberRequiredSigs = 3

	// Define DB key
	LastProcessedSequenceNumber = "LastSequenceNumber"
	LastSequenceNumber = "LastSequenceNumber"
	DefaultCreateTimeout   = time.Second * 10
	
)

var chainIdToName = map[uint8]string{
	1: "0x1.icon",
	2: "0x1.btc",
	3: "0x2.icon",
	4: "0x2.btc",
}

var (
	MessageStatusPending = "pending"
	MessageStatusSuccess = "success"
)

type MessageDecoded struct {
	Action       string
	TokenAddress string
	To           string
	Amount       []byte
}

type CSMessageResult struct {
	Sn      *big.Int
	Code    uint8
	Message []byte
}

type slaveRequestParams struct {
	MsgSn   string              `json:"msgSn"`
	FeeRate int                 `json:"feeRate"`
	Inputs  []slaveRequestInput `json:"inputs"`
}
type slaveRequestInput struct {
	TxHash string `json:"txHash"`
	Output int    `json:"output"`
	Amount int64  `json:"amount"`
}

type slaveRequestUpdateRelayMessageStatus struct {
	MsgSn  string `json:"msgSn"`
	Status string `json:"status"`
	TxHash string `json:"txHash"`
}

type slaveResponse struct {
	order int
	sigs  [][]byte
	err   error
}

type StoredMessageData struct {
	OriginalMessage  *relayTypes.Message
	TxHash           string
	OutputIndex      uint32
	Amount           uint64
	RecipientAddress string
	SenderAddress    string
	RuneId           string
	RuneAmount       uint64
	ActionMethod     string
	TokenAddress     string
}

type StoredRelayMessage struct {
	Message *relayTypes.Message
	Status  string
	TxHash  string
}

type Provider struct {
	logger              *zap.Logger
	cfg                 *Config
	client              IClient
	LastSavedHeightFunc func() uint64
	LastSerialNumFunc   func() *big.Int
	multisigAddrScript  []byte
	httpServer          chan struct{}
	db                  *leveldb.DB
	chainParam          *chaincfg.Params
	eth       			*ethclient.Client
	bitcoinState        *abi.BitcoinState
	runeFactory         *abi.Runefactory
}

type Config struct {
	provider.CommonConfig `json:",inline" yaml:",inline"`
	OpCode                int      `json:"op-code" yaml:"op-code"`
	RequestTimeout        int      `json:"request-timeout" yaml:"request-timeout"` // seconds
	UniSatURL             string   `json:"unisat-url" yaml:"unisat-url"`
	UniSatKey             string   `json:"unisat-key" yaml:"unisat-key"`
	UniSatWalletURL       string   `json:"unisat-wallet-url" yaml:"unisat-wallet-url"`
	MempoolURL            string   `json:"mempool-url" yaml:"mempool-url"`
	Type                  string   `json:"type" yaml:"type"`
	User                  string   `json:"rpc-user" yaml:"rpc-user"`
	Password              string   `json:"rpc-password" yaml:"rpc-password"`
	Mode                  string   `json:"mode" yaml:"mode"`
	SlaveServer1          string   `json:"slave-server-1" yaml:"slave-server-1"`
	SlaveServer2          string   `json:"slave-server-2" yaml:"slave-server-2"`
	Port                  string   `json:"port" yaml:"port"`
	ApiKey                string   `json:"api-key" yaml:"api-key"`
	MasterPubKey          string   `json:"masterPubKey" yaml:"masterPubKey"`
	Slave1PubKey          string   `json:"slave1PubKey" yaml:"slave1PubKey"`
	Slave2PubKey          string   `json:"slave2PubKey" yaml:"slave2PubKey"`
	RelayerPrivKey        string   `json:"relayerPrivKey" yaml:"relayerPrivKey"`
	RecoveryLockTime      int      `json:"recoveryLockTime" yaml:"recoveryLockTime"`
	Connections           []string `json:"connections" yaml:"connections"`
	EthRPC           	  string   `json:"eth-prc" yaml:"connections"`
	LastSequenceNumber    uint64   `json:"last-sequence-number" yaml:"last-sequence-number"`
	RuneFactory           string   `json:"rune-factory" yaml:"rune-factory"`
	BitcoinState          string   `json:"bitcoin-state" yaml:"bitcoin-state"`
}

// NewProvider returns new Icon provider
func (c *Config) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := c.sanitize(); err != nil {
		return nil, err
	}

	// Create the database file path
	dbPath := filepath.Join(homepath+"/data", BtcDB)

	// Open the database, creating it if it doesn't exist
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create database: %v", err)
	}

	// init sequence number 
	lastSeqKey := AddPrefixChainName(c.NID, []byte(LastSequenceNumber))
	_, err = db.Get(lastSeqKey, nil)
	if err == leveldb.ErrNotFound {
		db.Put(AddPrefixChainName(c.NID, []byte(LastSequenceNumber)), []byte(fmt.Sprintf("%d", c.LastSequenceNumber)), nil)
	}

	client, err := newClient(ctx, c.RPCUrl, c.User, c.Password, true, false, log)
	if err != nil {
		db.Close() // Close the database if client creation fails
		return nil, fmt.Errorf("failed to create new client: %v", err)
	}
	chainParam := &chaincfg.TestNet3Params
	if c.NID == "0x1.btc" {
		chainParam = &chaincfg.MainNetParams
	}
	c.HomeDir = homepath

	msPubkey, err := client.DecodeAddress(c.Address)
	if err != nil {
		return nil, err
	}

	createCtx, cancel := context.WithTimeout(ctx, DefaultCreateTimeout)
	defer cancel()
	rpc, err := ethclient.DialContext(createCtx, c.EthRPC)
	if err != nil {
		return nil, err
	}

	// init rune contract instance
	runeFactory, err := abi.NewRunefactory(common.HexToAddress(c.RuneFactory), rpc)
	if err != nil {
		return nil, err
	}

	// init bitcoinState contract instance
	bitcoinState, err := abi.NewBitcoinState(common.HexToAddress(c.BitcoinState), rpc)
	if err != nil {
		return nil, err
	}

	p := &Provider{
		logger:             log.With(zap.Stringp("nid", &c.NID), zap.Stringp("name", &c.ChainName)),
		cfg:                c,
		client:             client,
		LastSerialNumFunc:  func() *big.Int { return big.NewInt(0) },
		httpServer:         make(chan struct{}),
		db:                 db, // Add the database to the Provider
		chainParam:         chainParam,
		multisigAddrScript: msPubkey,
		eth: rpc,
		runeFactory: runeFactory,
		bitcoinState: bitcoinState,
	}
	// Run an http server to help btc interact each others
	go func() {
		if c.Mode == MasterMode {
			startMaster(c)
		} else {
			startSlave(c, p)
		}
		close(p.httpServer)
	}()

	return p, nil
}

func (p *Provider) CallSlaves(slaveRequestData []byte, path string) [][][]byte {
	resultChan := make(chan [][][]byte)
	go func() {
		responses := make(chan slaveResponse, 2)
		var wg sync.WaitGroup
		wg.Add(2)

		go requestPartialSign(p.cfg.ApiKey, p.cfg.SlaveServer1+path, slaveRequestData, responses, 1, &wg)
		go requestPartialSign(p.cfg.ApiKey, p.cfg.SlaveServer2+path, slaveRequestData, responses, 2, &wg)

		go func() {
			wg.Wait()
			close(responses)
		}()

		results := make([][][]byte, 2)
		for res := range responses {
			if res.err != nil {
				p.logger.Error("failed to call slaves", zap.Error(res.err))
				continue
			}
			results[res.order-1] = res.sigs
		}
		resultChan <- results
	}()

	return <-resultChan
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx)
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*relayTypes.Receipt, error) {
	res, err := p.client.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return &relayTypes.Receipt{
		TxHash: res.Txid,
	}, nil
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Init(ctx context.Context, homePath string, kms kms.KMS) error {
	return nil
}

// Wallet returns the wallet of the provider
func (p *Provider) Wallet() (*multisig.MultisigWallet, error) {
	return p.buildMultisigWallet()
}

func (p *Provider) Type() string {
	return p.cfg.ChainName
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

func (p *Provider) Listener(ctx context.Context, lastProcessedTx relayTypes.LastProcessedTx, blockInfoChan chan *relayTypes.BlockInfo) error {
	// run http server to help btc interact each others
	
	// loop thru 
	return nil
}

func (p *Provider) GetBitcoinUTXOs(server, address string, amountRequired int64, addressPkScript []byte) ([]*multisig.Input, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*int64(p.cfg.RequestTimeout)))
	defer cancel()
	resp, err := GetBtcUtxo(ctx, server, p.cfg.UniSatKey, address, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to query bitcoin UTXOs from unisat: %v", err)
	}
	inputs := []*multisig.Input{}
	var totalAmount int64

	utxos := resp.Data.Utxo
	// sort utxos by amount in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Satoshi.Cmp(utxos[j].Satoshi) == 1
	})

	for _, utxo := range utxos {
		if totalAmount >= amountRequired {
			break
		}
		isSpent, _ := p.isSpentUTXO(utxo.TxId, utxo.Vout)
		if isSpent {
			continue
		}
		outputAmount := utxo.Satoshi.Int64()
		inputs = append(inputs, &multisig.Input{
			TxHash:       utxo.TxId,
			OutputIdx:    uint32(utxo.Vout),
			OutputAmount: outputAmount,
			PkScript:     addressPkScript,
		})
		totalAmount += outputAmount
	}

	return inputs, nil
}

func (p *Provider) GetRuneUTXOs(server, address, runeId string, amountRequired uint128.Uint128, addressPkScript []byte) ([]*multisig.Input, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*int64(p.cfg.RequestTimeout)))
	defer cancel()
	resp, err := GetRuneUtxo(ctx, server, p.cfg.UniSatKey, address, runeId, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to query rune UTXOs from unisat: %v", err)
	}

	utxos := resp.Data.Utxo
	// sort utxos by amount in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Satoshi.Cmp(utxos[j].Satoshi) == 1
	})

	inputs := []*multisig.Input{}
	var totalAmount uint128.Uint128
	for _, utxo := range utxos {
		if totalAmount.Cmp(amountRequired) >= 0 {
			break
		}
		if len(utxo.Runes) != 1 {
			continue
		}
		isSpent, _ := p.isSpentUTXO(utxo.TxId, utxo.Vout)
		if isSpent {
			continue
		}
		inputs = append(inputs, &multisig.Input{
			TxHash:       utxo.TxId,
			OutputIdx:    uint32(utxo.Vout),
			OutputAmount: utxo.Satoshi.Int64(),
			PkScript:     addressPkScript,
		})
		runeAmount, _ := uint128.FromString(utxo.Runes[0].Amount)
		totalAmount = totalAmount.Add(runeAmount)
	}

	return inputs, nil
}

func (p *Provider) GetUTXORuneBalance(server, txId string, index int) (*ResponseUtxoRuneBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*int64(p.cfg.RequestTimeout)))
	defer cancel()
	resp, err := GetUtxoRuneBalance(ctx, server, p.cfg.UniSatKey, txId, index)
	if err != nil {
		return nil, fmt.Errorf("failed to query rune UTXO balance from unisat: %v", err)
	}
	return &resp, nil
}

func (p *Provider) GetFastestFee() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*int64(p.cfg.RequestTimeout)))
	defer cancel()
	fastestFee, err := GetFastestFee(ctx, p.cfg.UniSatWalletURL)
	if err != nil {
		return 0, err
	}
	return fastestFee, nil
}

func (p *Provider) CreateBitcoinMultisigTx(
	outputData []*wire.TxOut,
	feeRate int64,
	decodedData *MessageDecoded,
	msWallet *multisig.MultisigWallet,
	reqInputs []slaveRequestInput,
) ([]*multisig.Input, *wire.MsgTx, error) {
	// build receiver Pk Script
	receiverAddr, err := btcutil.DecodeAddress(decodedData.To, p.chainParam)
	if err != nil {
		return nil, nil, err
	}
	receiverPkScript, err := txscript.PayToAddrScript(receiverAddr)
	if err != nil {
		return nil, nil, err
	}
	// ----- BUILD OUTPUTS -----
	var outputs []*wire.TxOut
	var bitcoinAmountRequired int64
	var runeAmountRequired uint128.Uint128

	rlMsAddress, err := multisig.AddressOnChain(p.chainParam, msWallet)
	if err != nil {
		return nil, nil, err
	}
	msAddressStr := rlMsAddress.String()

	// add withdraw output
	amount := new(big.Int).SetBytes(decodedData.Amount)
	if decodedData.Action == MethodWithdrawTo || decodedData.Action == MethodDeposit || decodedData.Action == MethodRollback {
		if decodedData.TokenAddress == BTCToken {
			// transfer btc
			bitcoinAmountRequired = amount.Int64()
			outputs = []*wire.TxOut{
				// bitcoin send to receiver
				{
					Value:    bitcoinAmountRequired,
					PkScript: receiverPkScript,
				},
			}
		} else {
			// transfer rune
			runeAmountRequired = uint128.FromBig(amount)
			runeRequired, err := runestone.RuneIdFromString(decodedData.TokenAddress)
			if err != nil {
				return nil, nil, err
			}
			changeOutputId := uint32(1)
			runeOutput := &runestone.Runestone{
				Edicts: []runestone.Edict{
					{
						ID:     *runeRequired,
						Amount: runeAmountRequired,
						Output: 0,
					},
				},
				Pointer: &changeOutputId,
			}
			runeScript, err := runeOutput.Encipher()
			if err != nil {
				return nil, nil, err
			}

			outputs = []*wire.TxOut{
				// rune send to receiver
				{
					Value:    multisig.DUST_UTXO_AMOUNT,
					PkScript: receiverPkScript,
				},
				// rune change output
				{
					Value:    multisig.DUST_UTXO_AMOUNT,
					PkScript: msWallet.PKScript,
				},
				// rune OP_RETURN
				{
					Value:    0,
					PkScript: runeScript,
				},
			}

			bitcoinAmountRequired = multisig.DUST_UTXO_AMOUNT * 2
		}
	} else if decodedData.Action == MethodRefundTo {
		if decodedData.TokenAddress != BTCToken {
			return nil, nil, fmt.Errorf("refund is only supported for btc token, current token: %s", decodedData.TokenAddress)
		}
		uintAmount := amount.Uint64()
		bitcoinAmountRequired = int64(uintAmount)
		outputs = []*wire.TxOut{
			// bitcoin send to receiver
			{
				Value:    bitcoinAmountRequired,
				PkScript: receiverPkScript,
			},
		}
	}

	outputs = append(outputs, outputData...)
	// ----- BUILD INPUTS -----
	inputs, estFee, err := p.computeTx(feeRate, bitcoinAmountRequired, runeAmountRequired, decodedData.TokenAddress, msAddressStr, outputs, msWallet, reqInputs)
	if err != nil {
		return nil, nil, err
	}

	if decodedData.Action == MethodRefundTo {
		if estFee >= bitcoinAmountRequired {
			return nil, nil, fmt.Errorf("estimated fee is greater than the amount")
		} else {
			for _, output := range outputs {
				if bytes.Equal(output.PkScript, receiverPkScript) {
					output.Value = bitcoinAmountRequired - estFee
					break
				}
			}
		}
	}
	// create raw tx
	msgTx, err := multisig.CreateTx(inputs, outputs, msWallet.PKScript, estFee, 0)

	return inputs, msgTx, err
}

// calculateTxSize calculates the size of a transaction given the inputs, outputs, estimated fee, change address, chain parameters, and multisig wallet.
// It returns the size of the transaction in bytes and an error if any occurs during the process.
func (p *Provider) calculateTxSize(inputs []*multisig.Input, outputs []*wire.TxOut, estFee int64, msWallet *multisig.MultisigWallet) (int, error) {
	msgTx, err := multisig.CreateTx(inputs, outputs, msWallet.PKScript, estFee, 0)
	if err != nil {
		return 0, err
	}
	var rawTxBytes bytes.Buffer
	err = msgTx.Serialize(&rawTxBytes)
	if err != nil {
		return 0, err
	}
	baseSize := len(rawTxBytes.Bytes())
	totalSize := baseSize + len(inputs)*WitnessSize
	txSize := (baseSize*3 + totalSize) / 4
	return txSize, nil
}

func (p *Provider) computeTx(feeRate int64, satsToSend int64, runeToSend uint128.Uint128, runeId, address string, outputs []*wire.TxOut, msWallet *multisig.MultisigWallet, reqInputs []slaveRequestInput) ([]*multisig.Input, int64, error) {

	outputsCopy := make([]*wire.TxOut, len(outputs))
	copy(outputsCopy, outputs)

	inputs, err := p.selectUnspentUTXOs(satsToSend, runeToSend, runeId, address, msWallet.PKScript, reqInputs)
	sumSelectedInputs := multisig.SumInputsSat(inputs)
	if err != nil {
		return nil, 0, err
	}

	txSize, err := p.calculateTxSize(inputs, outputsCopy, 0, msWallet)
	if err != nil {
		return nil, 0, err
	}

	estFee := int64(txSize) * feeRate
	count := 0

	for sumSelectedInputs < satsToSend+estFee {
		// Create a fresh copy of outputs for each iteration
		iterationOutputs := make([]*wire.TxOut, len(outputs))
		copy(iterationOutputs, outputs)

		newSatsToSend := satsToSend + estFee
		var err error
		inputs, err = p.selectUnspentUTXOs(newSatsToSend, runeToSend, runeId, address, msWallet.PKScript, reqInputs)
		if err != nil {
			return nil, 0, err
		}

		sumSelectedInputs = multisig.SumInputsSat(inputs)

		txSize, err := p.calculateTxSize(inputs, iterationOutputs, estFee, msWallet)
		if err != nil {
			return nil, 0, err
		}

		estFee = feeRate * int64(txSize)

		count += 1
		if count > 500 {
			return nil, 0, fmt.Errorf("not enough sats for fee")
		}
	}

	return inputs, estFee, nil
}

func (p *Provider) selectUnspentUTXOs(satToSend int64, runeToSend uint128.Uint128, runeId string, address string, addressPkScript []byte, reqInputs []slaveRequestInput) ([]*multisig.Input, error) {
	// add tx fee the the required bitcoin amount
	inputs := []*multisig.Input{}
	if len(reqInputs) > 0 {
		for _, input := range reqInputs {
			inputs = append(inputs, &multisig.Input{
				TxHash:       input.TxHash,
				OutputIdx:    uint32(input.Output),
				OutputAmount: input.Amount,
				PkScript:     addressPkScript,
			})
		}
		return inputs, nil
	}
	if !runeToSend.IsZero() {
		// query rune UTXOs from unisat
		runeInputs, err := p.GetRuneUTXOs(p.cfg.UniSatURL, address, runeId, runeToSend, addressPkScript)
		if err != nil {
			return nil, err
		}
		if len(runeInputs) == 0 {
			return nil, fmt.Errorf("no rune UTXOs found")
		}
		inputs = append(inputs, runeInputs...)
	}

	// query bitcoin UTXOs from unisat
	bitcoinInputs, err := p.GetBitcoinUTXOs(p.cfg.UniSatURL, address, satToSend, addressPkScript)
	if err != nil {
		return nil, err
	}
	if len(bitcoinInputs) == 0 {
		return nil, fmt.Errorf("no bitcoin UTXOs found")
	}
	inputs = append(inputs, bitcoinInputs...)

	return inputs, nil
}

func (p *Provider) buildMultisigWallet() (*multisig.MultisigWallet, error) {
	masterPubKey, _ := hex.DecodeString(p.cfg.MasterPubKey)
	slave1PubKey, _ := hex.DecodeString(p.cfg.Slave1PubKey)
	slave2PubKey, _ := hex.DecodeString(p.cfg.Slave2PubKey)
	relayersMultisigInfo := multisig.MultisigInfo{
		PubKeys:            [][]byte{masterPubKey, slave1PubKey, slave2PubKey},
		EcPubKeys:          nil,
		NumberRequiredSigs: NumberRequiredSigs,
		RecoveryPubKey:     masterPubKey,
		RecoveryLockTime:   uint64(p.cfg.RecoveryLockTime),
	}
	msWallet, err := multisig.BuildMultisigWallet(&relayersMultisigInfo)
	if err != nil {
		p.logger.Error("failed to build multisig wallet: %v", zap.Error(err))
		return nil, err
	}

	return msWallet, nil
}

func (p *Provider) HandleBitcoinMessageTx(message *relayTypes.Message, feeRate int, reqInputs []slaveRequestInput) ([]*multisig.Input, *multisig.MultisigWallet, *wire.MsgTx, [][]byte, int, error) {
	msWallet, err := p.buildMultisigWallet()
	if err != nil {
		return nil, nil, nil, nil, 0, err
	}
	if feeRate == 0 {
		feeRate, err = p.GetFastestFee()
		if err != nil {
			p.logger.Error("failed to get fastest fee from unisat", zap.Error(err))
			feeRate = 1
		}
	}
	inputs, msgTx, err := p.buildTxMessage(message, int64(feeRate), msWallet, reqInputs)
	if err != nil {
		p.logger.Error("failed to build tx message: %v", zap.Error(err))
		return nil, nil, nil, nil, 0, err
	}
	relayerSigs, err := multisig.SignTapMultisig(p.cfg.RelayerPrivKey, msgTx, inputs, msWallet, 0)
	if err != nil {
		p.logger.Error("failed to sign tx message: %v", zap.Error(err))
		return nil, nil, nil, nil, 0, err
	}
	return inputs, msWallet, msgTx, relayerSigs, feeRate, err
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	if (strings.HasSuffix(message.Src, "icon") || strings.HasSuffix(message.Src, "btc")) && strings.HasSuffix(message.Dst, "btc") {
		p.logger.Info("starting to route message", zap.Any("message", message))

		key := []byte(message.Sn.String())
		storedRelayerMessage := StoredRelayMessage{}
		data, _ := p.db.Get(key, nil)
		if len(data) > 0 {
			json.Unmarshal(data, &storedRelayerMessage)
			if storedRelayerMessage.Status == MessageStatusSuccess {
				p.logger.Info("Message already success", zap.Any("message", storedRelayerMessage))
				callback(message.MessageKey(), &relayTypes.TxResponse{Code: relayTypes.Success, TxHash: storedRelayerMessage.TxHash}, nil)
				return nil
			}
		}

		storedMessage := StoredRelayMessage{
			Message: message,
			Status:  MessageStatusPending,
		}
		value, _ := json.Marshal(storedMessage)
		err := p.db.Put(key, value, nil)
		if err != nil {
			p.logger.Error("failed to store message in LevelDB: %v", zap.Error(err))
			return err
		}
		p.logger.Info("Message stored in LevelDB", zap.String("key", string(key)), zap.Any("message", message))

		if p.cfg.Mode != MasterMode {
			return nil
		}
		txHash, err := p.sendTransaction(ctx, message)
		if err != nil {
			p.logger.Error("failed to send transaction: %v", zap.Error(err))
			return err
		}
		p.logger.Info("transaction sent", zap.String("txHash", txHash), zap.Any("message", message))

		// update message status to success and save
		storedMessage.Status = MessageStatusSuccess
		storedMessage.TxHash = txHash
		value, _ = json.Marshal(storedMessage)
		p.db.Put(key, value, nil)

		// call to slave to update message status
		rsi := slaveRequestUpdateRelayMessageStatus{
			MsgSn:  message.Sn.String(),
			Status: MessageStatusSuccess,
			TxHash: txHash,
		}
		slaveRequestData, _ := json.Marshal(rsi)
		p.CallSlaves(slaveRequestData, "/update-relayer-message-status")

		// callback to clear relayer message
		callback(message.MessageKey(), &relayTypes.TxResponse{Code: relayTypes.Success, TxHash: txHash}, nil)
		return nil
	}
	return nil
}

func (p *Provider) cacheSpentUTXOs(inputs []*multisig.Input) {
	prefix := fmt.Sprintf("%s_utxo_spent", p.cfg.NID)
	for _, input := range inputs {
		key := fmt.Sprintf("%s_%s_%d", prefix, input.TxHash, input.OutputIdx)
		p.db.Put([]byte(key), []byte{1}, nil)
	}
}

func (p *Provider) removeCachedSpentUTXOs(inputs []*multisig.Input) {
	prefix := fmt.Sprintf("%s_utxo_spent", p.cfg.NID)
	for _, input := range inputs {
		key := fmt.Sprintf("%s_%s_%d", prefix, input.TxHash, input.OutputIdx)
		p.db.Delete([]byte(key), nil)
	}
}

func (p *Provider) isSpentUTXO(txHash string, outputIdx int) (bool, error) {
	prefix := fmt.Sprintf("%s_utxo_spent", p.cfg.NID)
	key := fmt.Sprintf("%s_%s_%d", prefix, txHash, outputIdx)
	data, err := p.db.Get([]byte(key), nil)
	if err != nil {
		return false, err
	}
	return len(data) > 0, nil
}

func (p *Provider) decodeMessage(message *relayTypes.Message) (CSMessageResult, error) {

	wrapperInfo := CSMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(message.Data, &wrapperInfo)
	if err != nil {
		p.logger.Error("failed to unmarshal message: %v", zap.Error(err))
		return CSMessageResult{}, err
	}

	messageDecoded := CSMessageResult{}
	_, err = codec.RLP.UnmarshalFromBytes(wrapperInfo.Payload, &messageDecoded)
	if err != nil {
		p.logger.Error("failed to unmarshal message: %v", zap.Error(err))
		return CSMessageResult{}, err
	}

	return messageDecoded, nil
}

func (p *Provider) decodeWithdrawToMessage(input []byte) (*MessageDecoded, []byte, error) {
	withdrawInfoWrapper := CSMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(input, &withdrawInfoWrapper)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal withdraw info wrapper: %v", err)
	}

	// withdraw info data
	withdrawInfoWrapperV2 := CSMessageRequestV2{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapper.Payload, &withdrawInfoWrapperV2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal withdraw info wrapper: %v", err)
	}
	// withdraw info
	withdrawInfo := &MessageDecoded{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapperV2.Data, &withdrawInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal withdraw info: %v", err)
	}

	return withdrawInfo, withdrawInfoWrapperV2.Data, nil
}

func (p *Provider) storedDataToMessageDecoded(sn string) (*MessageDecoded, []byte, error) {
	data, err := p.db.Get([]byte("RB"+sn), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve stored rollback message data: %v", err)
	}
	var storedData StoredMessageData
	err = json.Unmarshal(data, &storedData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal stored rollback message data: %v", err)
	}
	decodedData := &MessageDecoded{
		Action:       storedData.ActionMethod,
		To:           storedData.SenderAddress,
		TokenAddress: storedData.TokenAddress,
		Amount:       uint64ToBytes(storedData.Amount),
	}
	if storedData.RuneId != "" {
		decodedData.Amount = uint64ToBytes(storedData.RuneAmount)
	}
	decodeDataBuffer, _ := codec.RLP.MarshalToBytes(decodedData)
	return decodedData, decodeDataBuffer, nil
}

// todo: back to this fucntion later
func (p *Provider) buildTxMessage(message *relayTypes.Message, feeRate int64, msWallet *multisig.MultisigWallet, reqInputs []slaveRequestInput) ([]*multisig.Input, *wire.MsgTx, error) {
	outputs := []*wire.TxOut{}
	decodedData := &MessageDecoded{}
	switch message.EventType {
	case events.EmitMessage:
		messageDecoded, err := p.decodeMessage(message)

		// transaction message from icon decoded successfully
		// withdraw message has no messageDecoded.Code
		if err == nil {
			// 1 is transaction already success, no need to rollback
			if messageDecoded.Code == 1 {
				return nil, nil, fmt.Errorf("transaction already success")
			}

			// message code 0 is need to rollback
			// Process RollbackMessage
			// p.logger.Info("Detected rollback message", zap.String("sn", messageDecoded.Sn.String()))
			// messageDecoded, decodeDataBuffer, err := p.storedDataToMessageDecoded(messageDecoded.Sn.String())
			// if err != nil {
			// 	p.logger.Error("failed to get decode data: %v", zap.Error(err))
			// 	return nil, nil, err
			// }
			// scripts, _ := multisig.EncodePayloadToScripts(decodeDataBuffer)
			// outputs = multisig.BuildBridgeScriptsOutputs(outputs, scripts)
			// decodedData = messageDecoded
		} else {
			// Perform WithdrawData
			// data, opBufferData, err := p.decodeWithdrawToMessage(message.Data)
			// decodedData = data
			// if err != nil {
			// 	p.logger.Error("failed to decode message: %v", zap.Error(err))
			// 	return nil, nil, err
			// }
			// scripts, _ := multisig.EncodePayloadToScripts(opBufferData)
			// outputs = multisig.BuildBridgeScriptsOutputs(outputs, scripts)
		}
	case events.RollbackMessage:
		// p.logger.Info("Detected refund message", zap.String("sn", message.Sn.String()))
		// messageDecoded, decodeDataBuffer, err := p.storedDataToMessageDecoded(message.Sn.String())
		// if err != nil {
		// 	p.logger.Error("failed to get decode data: %v", zap.Error(err))
		// 	return nil, nil, err
		// }
		// if messageDecoded.TokenAddress != BTCToken {
		// 	return nil, nil, fmt.Errorf("only support refund for BTC")
		// }
		// decodedData = messageDecoded
		// scripts, _ := multisig.EncodePayloadToScripts(decodeDataBuffer)
		// outputs = multisig.BuildBridgeScriptsOutputs(outputs, scripts)

	default:
		return nil, nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}

	inputs, msgTx, err := p.CreateBitcoinMultisigTx(outputs, feeRate, decodedData, msWallet, reqInputs)
	return inputs, msgTx, err
}

func (p *Provider) call(ctx context.Context, message *relayTypes.Message) (string, error) {
	return "", nil
}

func (p *Provider) sendTransaction(ctx context.Context, message *relayTypes.Message) (string, error) {
	inputs, msWallet, msgTx, relayerSigs, feeRate, err := p.HandleBitcoinMessageTx(message, 0, []slaveRequestInput{})
	if err != nil {
		p.logger.Error("failed to handle bitcoin message tx: %v", zap.Error(err))
		return "", err
	}
	totalSigs := [][][]byte{relayerSigs}

	// send message sn and inputs to 2 slave relayers to get sign
	slaveInputs := []slaveRequestInput{}
	for _, input := range inputs {
		slaveInputs = append(slaveInputs, slaveRequestInput{
			TxHash: input.TxHash,
			Output: int(input.OutputIdx),
			Amount: input.OutputAmount,
		})
	}
	rsi := slaveRequestParams{
		MsgSn:   message.Sn.String(),
		FeeRate: feeRate,
		Inputs:  slaveInputs,
	}

	slaveRequestData, _ := json.Marshal(rsi)
	slaveSigs := p.CallSlaves(slaveRequestData, "")

	p.logger.Info("Slave signatures", zap.Any("slave sigs", slaveSigs))
	if len(slaveSigs) < 2 || len(slaveSigs[0]) == 0 || len(slaveSigs[1]) == 0 {
		return "", fmt.Errorf("slave sigs is empty")
	}

	totalSigs = append(totalSigs, slaveSigs...)

	// combine sigs
	signedMsgTx, err := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, msWallet, 0)

	if err != nil {
		p.logger.Error("err combine tx: ", zap.Error(err))
		return "", err
	}
	p.logger.Info("signedMsgTx", zap.Any("transaction", signedMsgTx))
	var buf bytes.Buffer
	err = signedMsgTx.Serialize(&buf)

	if err != nil {
		p.logger.Error("err serialize tx: ", zap.Error(err))
		return "", err
	}

	txSize := len(buf.Bytes())
	p.logger.Info("--------------------txSize--------------------", zap.Int("size", txSize))
	signedMsgTxHex := hex.EncodeToString(buf.Bytes())
	p.logger.Info("signedMsgTxHex", zap.String("transaction_hex", signedMsgTxHex))

	p.cacheSpentUTXOs(inputs)

	// Broadcast transaction to bitcoin network
	txHash, err := p.client.SendRawTransaction(p.cfg.MempoolURL, []json.RawMessage{json.RawMessage(signedMsgTxHex)})
	if err != nil {
		p.removeCachedSpentUTXOs(inputs)
		p.logger.Error("failed to send raw transaction", zap.Error(err))
		return "", err
	}

	p.logger.Info("txHash", zap.String("transaction_hash", txHash))
	return txHash, nil
}

func (p *Provider) handleSequence(ctx context.Context) error {
	return nil
}

func (p *Provider) waitForTxResult(ctx context.Context, mk *relayTypes.MessageKey, txHash string, callback relayTypes.TxResponseFunc) {

}

func (p *Provider) MessageReceived(ctx context.Context, key *relayTypes.MessageKey) (bool, error) {
	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	return nil, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*relayTypes.Message, error) {
	// blocks, err := p.fetchBlockMessages(ctx, &HeightRange{fromHeight, toHeight})
	// if err != nil {
	// 	return nil, err
	// }
	// var messages []*relayTypes.Message
	// for _, block := range blocks {
	// 	messages = append(messages, block.Messages...)
	// }

	// todo: update generate mess here
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	msg := &relayTypes.Message{
		Sn:        sn,
		EventType: events.RevertMessage,
	}
	_, err := p.call(ctx, msg)
	return err
}

// SetFee
func (p *Provider) SetFee(ctx context.Context, networkdID string, msgFee, resFee *big.Int) error {
	msg := &relayTypes.Message{
		Src:       networkdID,
		Sn:        msgFee,
		ReqID:     resFee,
		EventType: events.SetFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ClaimFee
func (p *Provider) ClaimFee(ctx context.Context) error {
	msg := &relayTypes.Message{
		EventType: events.ClaimFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	return 0, nil
}

func (p *Provider) SetAdmin(ctx context.Context, address string) error {
	msg := &relayTypes.Message{
		Src:       address,
		EventType: events.SetAdmin,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ExecuteRollback
func (p *Provider) ExecuteRollback(ctx context.Context, sn *big.Int) error {
	return nil
}

func (p *Provider) getStartHeight(latestHeight, lastSavedHeight uint64) (uint64, error) {
	startHeight := lastSavedHeight
	if p.cfg.StartHeight > 0 && p.cfg.StartHeight < latestHeight {
		return p.cfg.StartHeight, nil
	}

	if startHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if startHeight != 0 && startHeight < latestHeight {
		return startHeight, nil
	}

	return latestHeight, nil
}

func (p *Provider) getHeightStream(done <-chan bool, fromHeight, toHeight uint64) <-chan *HeightRange {
	heightChan := make(chan *HeightRange)
	go func(fromHeight, toHeight uint64, heightChan chan *HeightRange) {
		defer close(heightChan)
		for fromHeight < toHeight {
			select {
			case <-done:
				return
			case heightChan <- &HeightRange{
				Start: fromHeight,
				End:   fromHeight + min(2, toHeight-fromHeight),
			}:
				fromHeight += min(2, toHeight-fromHeight)
			}
		}
	}(fromHeight, toHeight, heightChan)
	return heightChan
}

func (p *Provider) extractOutputReceiver(tx *wire.MsgTx) []string {
	receiverAddresses := []string{}
	for _, out := range tx.TxOut {
		receiverAddresses = append(receiverAddresses, p.getAddressesFromTx(out, p.chainParam)...)
	}
	return receiverAddresses
}

func (p *Provider) getNumOfPipelines(diff int) int {
	if diff <= runtime.NumCPU() {
		return diff
	}
	return runtime.NumCPU() / 2
}

func (p *Provider) getAddressesFromTx(txOut *wire.TxOut, chainParams *chaincfg.Params) []string {
	receiverAddresses := []string{}

	scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, chainParams)
	if err != nil {
		p.logger.Error("Script: Unable to parse (possibly OP_RETURN)", zap.Error(err))
	} else {
		p.logger.Info("Script Class", zap.String("scriptClass", scriptClass.String()))
		if len(addresses) > 0 {
			p.logger.Info("Receiver Address", zap.String("address", addresses[0].String()))
			receiverAddresses = append(receiverAddresses, addresses[0].String())
		}
	}

	return receiverAddresses
}

// SubscribeMessageEvents subscribes to the message events
// Expermental: Allows to subscribe to the message events realtime without fully syncing the chain
// func (p *Provider) SubscribeMessageEvents(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, opts *types.SubscribeOpts, resetFunc func()) error {
// 	return nil
// }

// SetLastSavedHeightFunc sets the function to save the last saved height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}

// GetLastSavedHeight returns the last saved height
func (p *Provider) GetLastSavedHeight() uint64 {
	return p.LastSavedHeightFunc()
}

func (p *Provider) FetchTxMessages(ctx context.Context, txHash string) ([]*relayTypes.Message, error) {
	return nil, nil
}

func (p *Config) sanitize() error {
	// TODO:
	return nil
}

func (c *Config) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("bitcoin provider rpc endpoint is empty")
	}
	return nil
}
