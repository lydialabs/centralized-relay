package wasm

import (
	"context"
	"fmt"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

const (
	ChainType string = "wasm"
)

type Provider struct {
	logger  *zap.Logger
	config  ProviderConfig
	client  client.IClient
	txMutex sync.Mutex
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
		TxHash: txHash,
		Height: uint64(res.TxResponse.Height),
		Status: abiTypes.CodeTypeOK == res.TxResponse.Code,
	}, nil
}

func (p *Provider) NID() string {
	return p.config.NID
}

func (p *Provider) ChainName() string {
	return p.config.ChainName
}

func (p *Provider) Init(ctx context.Context) error {
	return nil
}

func (p *Provider) Type() string {
	return ChainType
}

func (p *Provider) ProviderConfig() provider.ProviderConfig {
	return p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan relayTypes.BlockInfo) error {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		p.logger.Error("failed to get latest block height: ", zap.Error(err))
		return err
	}

	startHeight, err := p.getStartHeight(ctx, latestHeight, lastSavedHeight)
	if err != nil {
		p.logger.Error("failed to determine start height: ", zap.Error(err))
		return err
	}

	blockInterval, err := time.ParseDuration(p.config.BlockInterval)
	if err != nil {
		p.logger.Error("failed to parse block interval: ", zap.Error(err))
		return err
	}

	blockIntervalTicker := time.NewTicker(blockInterval)
	defer blockIntervalTicker.Stop()

	p.logger.Info("start querying from height", zap.Uint64("start-height", startHeight))

	for {
		select {
		case <-blockIntervalTicker.C:
			func() {
				done := make(chan interface{})
				defer close(done)

				heightStream := p.getHeightStream(done, startHeight, latestHeight)

				numOfPipelines := runtime.NumCPU() //Todo tune or configure this

				pipelines := make([]<-chan interface{}, numOfPipelines)

				for i := 0; i < numOfPipelines; i++ {
					pipelines[i] = p.getBlockInfoStream(done, heightStream)
				}

				for bn := range concurrency.FanIn(done, pipelines...) {
					block, ok := bn.(relayTypes.BlockInfo)
					if !ok || block.HasError() {
						if !block.HasError() {
							block.Error = fmt.Errorf("received invalid block type -> required: %T, got: %T", relayTypes.BlockInfo{}, bn)
						}
						p.logger.Error("error receiving block: ", zap.Error(block.Error))
						continue
					}
					blockInfo <- relayTypes.BlockInfo{
						Height: block.Height, Messages: block.Messages,
					}
				}
			}()
		}
	}

	return nil
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {

	//Build message
	msg := p.getMsgExecuteContract(message)

	res, err := p.client.SendTx(ctx, p.buildTxFactory(), []sdkTypes.Msg{&msg})
	if err != nil || res.Code != abiTypes.CodeTypeOK {
		if err == nil {
			err = fmt.Errorf("failed to send tx: %v", res.RawLog)
		}
		p.logger.Error("failed to route message: ", zap.Error(err))
		callback(message.MessageKey(), relayTypes.TxResponse{}, err)
		return err
	}

	callback(message.MessageKey(), relayTypes.TxResponse{
		Height:    res.Height,
		TxHash:    res.TxHash,
		Codespace: res.Codespace,
		Code:      relayTypes.ResponseCode(res.Code),
		Data:      res.Data,
	}, nil)

	return nil
}

func (p *Provider) MessageReceived(ctx context.Context, key relayTypes.MessageKey) (bool, error) {
	_, err := p.client.QuerySmartContract(ctx, p.config.ContractAddress, []byte("hello"))
	if err != nil {
		p.logger.Error("failed to check if message is received: ", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	coin, err := p.client.GetBalance(ctx, addr, "denomination")
	if err != nil {
		p.logger.Error("failed to query balance: ", zap.Error(err))
		return nil, err
	}
	return &relayTypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Uint64(),
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *relayTypes.MessageKeyWithMessageHeight) (*relayTypes.Message, error) {
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) getStartHeight(ctx context.Context, latestHeight, lastSavedHeight uint64) (uint64, error) {
	if lastSavedHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		return lastSavedHeight, nil
	}

	return latestHeight, nil
}

func (p *Provider) getHeightStream(done <-chan interface{}, fromHeight, toHeight uint64) <-chan uint64 {
	heightStream := make(chan uint64)
	go func() {
		defer close(heightStream)
		for i := fromHeight; i <= toHeight; i++ {
			select {
			case <-done:
				return
			case heightStream <- i:
			}
		}
	}()
	return heightStream
}

func (p *Provider) getBlockInfoStream(done <-chan interface{}, heightStream <-chan uint64) <-chan interface{} {
	blockInfoStream := make(chan interface{})
	go func() {
		defer close(blockInfoStream)
		for {
			select {
			case <-done:
				return
			case height := <-heightStream:
				searchParam := types.TxSearchParam{}
				messages, err := p.client.GetMessages(context.Background(), searchParam)
				blockInfoStream <- relayTypes.BlockInfo{
					Height:   height,
					Messages: messages,
					Error:    err,
				}
			}
		}
	}()
	return blockInfoStream
}

func (p *Provider) buildTxFactory() tx.Factory {
	return tx.Factory{}.
		WithKeybase(p.client.Context().Keyring).
		WithFeePayer(p.client.Context().FeePayer).
		WithChainID(p.config.ChainID).
		WithGasAdjustment(p.config.GasAdjustment)
}

func (p *Provider) getMsgExecuteContract(message *relayTypes.Message) wasmTypes.MsgExecuteContract {
	return wasmTypes.MsgExecuteContract{
		Sender:   p.client.Context().FromAddress.String(),
		Contract: p.config.ContractAddress,
		Msg:      []byte("msg here"),
	}
}
