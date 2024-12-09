package bitcoin

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"encoding/hex"
	"fmt"
	"io"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/utils/multisig"

	"github.com/icon-project/icon-bridge/common/codec"
	evmAbi "github.com/icon-project/centralized-relay/relayer/chains/evm/abi"
)

func GetRuneTxIndex(endpoint, method, bearToken, txId string, index int) (*RuneTxIndexResponse, error) {
	client := &http.Client{}
	endpoint = endpoint + "/runes/utxo/" + txId + "/" + strconv.FormatUint(uint64(index), 10) + "/balance"
	fmt.Println(endpoint)
	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", bearToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var resp *RuneTxIndexResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func ToXCallMessage(data interface{}, from, to string, sn uint, protocols []string, requester common.Address) ([]byte, error) {
	var res, calldata []byte
	bitcoinStateAbi, _ := abi.JSON(strings.NewReader(evmAbi.BitcoinStateMetaData.ABI))
	nonfungibleABI, _ := abi.JSON(strings.NewReader(evmAbi.NonfungiblePositionManagerMetaData.ABI))
	routerABI, _ := abi.JSON(strings.NewReader(evmAbi.IrouterMetaData.ABI))
	addressTy, _ := abi.NewType("address", "", nil)
	bytes, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: bytes,
		},
	}

	switch data.(type) {
	case multisig.RadFiProvideLiquidityMsg:
		dataMint := data.(multisig.RadFiProvideLiquidityMsg)
		mintParams := evmAbi.INonfungiblePositionManagerMintParams{
			Token0:         dataMint.Detail.Token0,
			Token1:         dataMint.Detail.Token1,
			Fee:            big.NewInt(int64(dataMint.Detail.Fee) * 100),
			TickLower:      big.NewInt(int64(dataMint.Detail.LowerTick)),
			TickUpper:      big.NewInt(int64(dataMint.Detail.UpperTick)),
			Amount0Desired: dataMint.Detail.Amount0Desired.Big(),
			Amount1Desired: dataMint.Detail.Amount1Desired.Big(),
			Recipient:      common.HexToAddress(to),
			Deadline:       big.NewInt(10000000000),
		}

		mintParams.Amount0Min = mulDiv(mintParams.Amount0Desired, big.NewInt(int64(dataMint.Detail.Min0)), big.NewInt(1e4))
		mintParams.Amount1Min = mulDiv(mintParams.Amount1Desired, big.NewInt(int64(dataMint.Detail.Min1)), big.NewInt(1e4))

		var err error
		if dataMint.InitPrice != nil && !dataMint.InitPrice.IsZero() {
			encodeInitPoolArgs, err := nonfungibleABI.Pack("initPoolHelper", mintParams, dataMint.Token0, dataMint.Token1, dataMint.InitPrice.Big())
			if err != nil {
				return nil, err
			}

			calldata, err = bitcoinStateAbi.Pack("initPool", encodeInitPoolArgs[4:])

		} else {
			calldata, err = nonfungibleABI.Pack("mint", mintParams)
		}

		if err != nil {
			return nil, err
		}

	case multisig.RadFiWithdrawLiquidityMsg:
		withdrawLiquidityInfo := data.(multisig.RadFiWithdrawLiquidityMsg)

		decreaseLiquidityData := evmAbi.INonfungiblePositionManagerDecreaseLiquidityParams{
			TokenId: withdrawLiquidityInfo.NftId.Big(),
			Amount0Min: withdrawLiquidityInfo.Amount0Min.Big(),
			Amount1Min: withdrawLiquidityInfo.Amount1Min.Big(),
			Liquidity: withdrawLiquidityInfo.LiquidityValue.Big(),
			Deadline: big.NewInt(10000000000),
		}

		decreaseLiquidityCalldata, err := nonfungibleABI.Pack("decreaseLiquidity", decreaseLiquidityData)
		if err != nil {
			return nil, err
		}

		calldata, err = bitcoinStateAbi.Pack("removeLiquidity", decreaseLiquidityCalldata)
		if err != nil {
			return nil, err
		}

	case multisig.RadFiIncreaseLiquidityMsg:
		increaseLiquidityInfo := data.(multisig.RadFiIncreaseLiquidityMsg)
		increaseLiquidityData := evmAbi.INonfungiblePositionManagerIncreaseLiquidityParams{
			TokenId: increaseLiquidityInfo.NftId.Big(),
			Amount0Desired: increaseLiquidityInfo.Amount0Desired.Big(), //todo fill in
			Amount1Desired: increaseLiquidityInfo.Amount1Desired.Big(), //todo fill in
			Deadline: big.NewInt(10000000000),
			Amount0Min: increaseLiquidityInfo.Amount0Min.Big(),
			Amount1Min: increaseLiquidityInfo.Amount1Min.Big(),
		}

		var err error
		calldata, err = nonfungibleABI.Pack("increaseLiquidity", increaseLiquidityData)
		if err != nil {
			return nil, err
		}

	case multisig.RadFiSwapMsg:
		swapInfo := data.(multisig.RadFiSwapMsg)
		var err error
		if swapInfo.IsExactInOut {
			//exact in
			swapExactInData := evmAbi.ISwapRouterExactInputParams{
				Path: swapInfo.Path,
				AmountIn: swapInfo.AmountIn.Big(), // todo:
				AmountOutMinimum: swapInfo.AmountOutMinimum.Big(), // todo:
				Recipient: common.HexToAddress(to),
				Deadline: big.NewInt(10000000000),
			}

			calldata, err = routerABI.Pack("exactInput", swapExactInData)
			if err != nil {
				return nil, err
			}

		} else {
			//exact out
			swapExactOutData := evmAbi.ISwapRouterExactOutputParams{
				Path: swapInfo.Path,
				Recipient: common.HexToAddress(to), // todo: review
				Deadline: big.NewInt(10000000000),
				AmountInMaximum: swapInfo.AmountInMaximum.Big(), // todo:
				AmountOut: swapInfo.AmountIn.Big(),
			}

			calldata, err = nonfungibleABI.Pack("exactOutput", swapExactOutData)
			if err != nil {
				return nil, err
			}

		}
	case multisig.RadFiCollectFeesMsg:
		collectInfo := data.(multisig.RadFiCollectFeesMsg)
		collectParams := evmAbi.INonfungiblePositionManagerCollectParams{
			TokenId: collectInfo.NftId.Big(),
			Amount0Max: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)),
			Amount1Max: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)),
			Recipient: common.HexToAddress(to), // todo:
		}

		collectCalldata, err := nonfungibleABI.Pack("collect", collectParams)
		if err != nil {
			return nil, err
		}

		calldata, err = bitcoinStateAbi.Pack("removeLiquidity", collectCalldata)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("not supported")
	}

	// encode with requester
	calldataWithRequester, err := arguments.Pack(requester, calldata)
	if err != nil {
		return nil, err
	}
	// encode to xcall format
	res, err = XcallFormat(calldataWithRequester, from, to, sn, protocols)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func XcallFormat(callData []byte, from, to string, sn uint, protocols []string, messType uint8) ([]byte, error) {
	//
	csV2 := CSMessageRequestV2{
		From:        from,
		To:          to,
		Sn:          big.NewInt(int64(sn)).Bytes(),
		MessageType: messType,
		Data:        callData,
		Protocols:   protocols,
	}

	//
	cvV2EncodeMsg, err := codec.RLP.MarshalToBytes(csV2)
	if err != nil {
		return nil, err
	}

	message := CSMessage{
		MsgType: big.NewInt(int64(CS_REQUEST)).Bytes(),
		Payload: cvV2EncodeMsg,
	}

	//
	finalMessage, err := codec.RLP.MarshalToBytes(message)
	if err != nil {
		return nil, err
	}

	return finalMessage, nil
}

func uint64ToBytes(amount uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, amount)
	return bytes
}

// Helper function to get minimum of two uint64 values
func min(a, b uint64) uint64 {
	if a <= b {
		return a
	}
	return b
}

func mulDiv(a, nNumerator, nDenominator *big.Int) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(a, nNumerator), nDenominator)
}

func BuildPath(paths []common.Address, fees []int64) ([]byte, error) {
	var temp []byte
	for i := 0; i < len(fees); i++ {
		temp = append(temp, paths[i].Bytes()...)
		temp1 := fmt.Sprintf("%06x", fees[i])
		fee, err := hex.DecodeString(temp1)
		if err != nil {
			return nil, err
		}
		temp = append(temp, fee...)
	}
	temp = append(temp, paths[len(paths)-1].Bytes()...)
	return temp, nil
}

func AddPrefixChainName(chainName string, key []byte) []byte {
	prefix := fmt.Sprintf("%s_", chainName)
	return append([]byte(prefix), key...)
}