package bitcoin

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/utils/multisig"

	radfiAbi "github.com/icon-project/centralized-relay/relayer/chains/bitcoin/abi"
	"github.com/icon-project/icon-bridge/common/codec"
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

func ToXCallMessage(data interface{}, from, to string, sn uint, protocols []string, requester string, runeFactory *radfiAbi.Runefactory) ([]byte, []byte, error) {
	var res, calldata []byte
	bitcoinStateAbi, _ := abi.JSON(strings.NewReader(radfiAbi.BitcoinstateMetaData.ABI))
	nonfungibleABI, _ := abi.JSON(strings.NewReader(radfiAbi.NonfungiblePositionManagerMetaData.ABI))
	routerABI, _ := abi.JSON(strings.NewReader(radfiAbi.IrouterMetaData.ABI))
	
	// addressTy, _ := abi.NewType("address", "", nil)
	// bytes, _ := abi.NewType("bytes", "", nil)

	// arguments := abi.Arguments{
	// 	{
	// 		Type: addressTy,
	// 	},
	// 	{
	// 		Type: bytes,
	// 	},
	// }

	switch v := data.(type) {
	case multisig.RadFiProvideLiquidityMsg:
		dataMint := v

		// get address token0 token1 from contract
		token0, err := runeFactory.ComputeTokenAddress(nil, dataMint.Token0Id.String())
		if err != nil {
			return nil, nil, err
		}

		token1, err := runeFactory.ComputeTokenAddress(nil, dataMint.Token1Id.String())
		if err != nil {
			return nil, nil, err
		}

		mintParams := radfiAbi.INonfungiblePositionManagerMintParams{
			Token0:         token0,
			Token1:         token1,
			Fee:            big.NewInt(int64(dataMint.Fee)),
			TickLower:      big.NewInt(int64(dataMint.Ticks.LowerTick)),
			TickUpper:      big.NewInt(int64(dataMint.Ticks.UpperTick)),
			Amount0Min: big.NewInt(0).Div(big.NewInt(0).Mul(dataMint.Amount0Desired.Big(), big.NewInt(99999)), big.NewInt(100000)),
			Amount1Min: big.NewInt(0).Div(big.NewInt(0).Mul(dataMint.Amount1Desired.Big(), big.NewInt(99999)), big.NewInt(100000)),
			Amount0Desired: dataMint.Amount0Desired.Big(),
			Amount1Desired: dataMint.Amount1Desired.Big(),
			Recipient:      common.HexToAddress(to),
			Deadline:       big.NewInt(10000000000),
		}

		if !dataMint.InitPrice.IsZero() {
			encodeInitPoolArgs, err := nonfungibleABI.Pack("initPoolHelper", mintParams, dataMint.Token0Id.String(), dataMint.Token0Decimal, dataMint.Token1Id.String(), dataMint.Token1Decimal, dataMint.InitPrice.Big())
			if err != nil {
				return nil, nil, err
			}

			calldata, err = bitcoinStateAbi.Pack("initPool", encodeInitPoolArgs[4:])
		} else {
			calldata, err = nonfungibleABI.Pack("mint", mintParams)
		}

		if err != nil {
			return nil, nil, err
		}

	case multisig.RadFiWithdrawLiquidityMsg:
		withdrawLiquidityInfo := v

		decreaseLiquidityData := radfiAbi.INonfungiblePositionManagerDecreaseLiquidityParams{
			TokenId:    withdrawLiquidityInfo.NftId.Big(),
			Amount0Min: withdrawLiquidityInfo.Amount0.Big(),
			Amount1Min: withdrawLiquidityInfo.Amount1.Big(),
			Liquidity:  withdrawLiquidityInfo.LiquidityValue.Big(),
			Deadline:   big.NewInt(10000000000),
		}

		decreaseLiquidityCalldata, err := nonfungibleABI.Pack("decreaseLiquidity", decreaseLiquidityData)
		if err != nil {
			return nil, nil, err
		}

		calldata, err = bitcoinStateAbi.Pack("removeLiquidity", decreaseLiquidityCalldata)
		if err != nil {
			return nil, nil, err
		}

	case multisig.RadFiIncreaseLiquidityMsg:
		increaseLiquidityInfo := v
		increaseLiquidityData := radfiAbi.INonfungiblePositionManagerIncreaseLiquidityParams{
			TokenId:        increaseLiquidityInfo.NftId.Big(),
			Amount0Desired: increaseLiquidityInfo.Amount0.Big(), //todo fill in
			Amount1Desired: increaseLiquidityInfo.Amount1.Big(), //todo fill in
			Deadline:       big.NewInt(10000000000),
			Amount0Min:     increaseLiquidityInfo.Amount0.Big(),
			Amount1Min:     increaseLiquidityInfo.Amount1.Big(),
		}

		var err error
		calldata, err = nonfungibleABI.Pack("increaseLiquidity", increaseLiquidityData)
		if err != nil {
			return nil, nil, err
		}

	case multisig.RadFiSwapMsg:
		swapInfo := v
		var err error

		// get token from contract to build path
		tokens := []common.Address{}
		for _, v := range tokens {
			tokenAddr, err := runeFactory.ComputeTokenAddress(nil, v.String())
			if err != nil {
				return nil, nil, err
			}

			tokens = append(tokens, tokenAddr)
		}

		path, err := BuildPath(tokens, swapInfo.Fees)
		if err != nil {
			return nil, nil, err
		}

		if swapInfo.IsExactIn {
			//exact in
			swapExactInData := radfiAbi.ISwapRouterExactInputParams{
				Path:             path,
				AmountIn:         swapInfo.AmountIn.Big(),  // todo:
				AmountOutMinimum: swapInfo.AmountOut.Big(), // todo:
				Recipient:        common.HexToAddress(to),
				Deadline:         big.NewInt(10000000000),
			}

			calldata, err = routerABI.Pack("exactInput", swapExactInData)
			if err != nil {
				return nil, nil, err
			}

		} else {
			//exact out
			swapExactOutData := radfiAbi.ISwapRouterExactOutputParams{
				Path:            path,
				Recipient:       common.HexToAddress(to), // todo: review
				Deadline:        big.NewInt(10000000000),
				AmountInMaximum: swapInfo.AmountIn.Big(), // todo:
				AmountOut:       swapInfo.AmountIn.Big(),
			}

			calldata, err = nonfungibleABI.Pack("exactOutput", swapExactOutData)
			if err != nil {
				return nil, nil, err
			}

		}
	case multisig.RadFiCollectFeesMsg:
		collectInfo := v
		collectParams := radfiAbi.INonfungiblePositionManagerCollectParams{
			TokenId:    collectInfo.NftId.Big(),
			Amount0Max: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)),
			Amount1Max: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)),
			Recipient:  common.HexToAddress(to), // todo:
		}

		collectCalldata, err := nonfungibleABI.Pack("collect", collectParams)
		if err != nil {
			return nil, nil, err
		}

		calldata, err = bitcoinStateAbi.Pack("removeLiquidity", collectCalldata)
		if err != nil {
			return nil, nil, err
		}

	default:
		return nil, nil, fmt.Errorf("not supported")
	}

	// todo: uncomment for batch process encode with requester
	// calldataWithRequester, err := arguments.Pack(requester, calldata)
	// if err != nil {
	// 	return nil, err
	// }

	// encode to xcall format
	res, err := XcallFormat(calldata, from, to, sn, protocols)
	if err != nil {
		return nil, nil, err
	}

	return res, calldata, nil
}

func XcallFormat(callData []byte, from, to string, sn uint, protocols []string) ([]byte, error) {
	//
	csV2 := CSMessageRequestV2{
		From:        from,
		To:          to,
		Sn:          big.NewInt(int64(sn)).Bytes(),
		MessageType: uint8(CALL_MESSAGE_TYPE),
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

	fmt.Println("fucker")
	fmt.Println(protocols)

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

func BuildPath(paths []common.Address, fees []uint32) ([]byte, error) {
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
