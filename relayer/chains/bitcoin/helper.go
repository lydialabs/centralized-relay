package bitcoin

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	bitcoinABI "github.com/icon-project/centralized-relay/relayer/chains/bitcoin/abi"
	"github.com/icon-project/centralized-relay/utils/multisig"
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

func ToXCallMessage(data interface{}, from, to string, sn uint, protocols []string, requester, token0, token1 common.Address) ([]byte, error) {
	var res, calldata []byte
	bitcoinStateAbi, _ := abi.JSON(strings.NewReader(bitcoinABI.BitcoinStateMetaData.ABI))
	nonfungibleABI, _ := abi.JSON(strings.NewReader(bitcoinABI.InonfungibleTokenMetaData.ABI))
	routerABI, _ := abi.JSON(strings.NewReader(bitcoinABI.IswaprouterMetaData.ABI))
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

	amount0, _ := big.NewInt(0).SetString("18999999999999999977305673", 10)

	switch data.(type) {
	case multisig.RadFiProvideLiquidityMsg:
		dataMint := data.(multisig.RadFiProvideLiquidityMsg)
		mintParams := bitcoinABI.INonfungiblePositionManagerMintParams{
			Token0:         token0,
			Token1:         token1,
			Fee:            big.NewInt(int64(dataMint.Detail.Fee) * 100),
			TickLower:      big.NewInt(int64(dataMint.Detail.LowerTick)),
			TickUpper:      big.NewInt(int64(dataMint.Detail.UpperTick)),
			Amount0Desired: amount0,
			Amount1Desired: big.NewInt(539580403982610478),
			Recipient:      common.HexToAddress(to),
			Deadline:       big.NewInt(1000000000),
		}

		mintParams.Amount0Min = mulDiv(mintParams.Amount0Desired, big.NewInt(int64(dataMint.Detail.Min0)), big.NewInt(1e4))
		mintParams.Amount1Min = mulDiv(mintParams.Amount1Desired, big.NewInt(int64(dataMint.Detail.Min1)), big.NewInt(1e4))

		var err error
		if dataMint.InitPrice == uint256.NewInt(0) {
			calldata, err = bitcoinStateAbi.Pack("initPool", mintParams, "btc", "rad", 1e0)
			if err != nil {
				return nil, err
			}
		} else {
			calldata, err = nonfungibleABI.Pack("mint", mintParams)
		}
		if err != nil {
			return nil, err
		}

	case multisig.RadFiWithdrawLiquidityMsg:
		withdrawLiquidityInfo := data.(multisig.RadFiWithdrawLiquidityMsg)
		
		decreaseLiquidityData := bitcoinABI.INonfungiblePositionManagerDecreaseLiquidityParams{
			TokenId: withdrawLiquidityInfo.NftId.ToBig(),
			Amount0Min: big.NewInt(0),
			Amount1Min: big.NewInt(0),
			Liquidity: withdrawLiquidityInfo.LiquidityValue.ToBig(),
			Deadline: big.NewInt(1000000000),
		}

		decreaseLiquidityCalldata, err := nonfungibleABI.Pack("decreaseLiquidity", decreaseLiquidityData)
		if err != nil {
			return nil, err
		}

		uint8Ty, _ := abi.NewType("uint8", "", nil)
		bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	
		withdrawLiquidityArgs := abi.Arguments{
			{
				Type: addressTy,
			},
			{
				Type: bytes,
			},
			{
				Type: uint8Ty,
			},
			{
				Type: bytes32Ty,
			},
		}

		calldata, _ = withdrawLiquidityArgs.Pack(decreaseLiquidityCalldata, withdrawLiquidityInfo.V, withdrawLiquidityInfo.R, withdrawLiquidityInfo.S)

	case multisig.RadFiIncreaseLiquidityMsg:
		increaseLiquidityInfo := data.(multisig.RadFiIncreaseLiquidityMsg)
		increaseLiquidityData := bitcoinABI.INonfungiblePositionManagerIncreaseLiquidityParams{
			TokenId: increaseLiquidityInfo.NftId.ToBig(),
			Amount0Desired: big.NewInt(0), //todo fill in	
			Amount1Desired: big.NewInt(0), //todo fill in
			Deadline: big.NewInt(1000000000),
			Amount0Min: big.NewInt(0),
			Amount1Min: big.NewInt(0),
		}

		var err error
		calldata, err = nonfungibleABI.Pack("increaseLiquidity", increaseLiquidityData)
		if err != nil {
			return nil, err
		}

	case multisig.RadFiSwapMsg:
		swapInfo := data.(multisig.RadFiSwapMsg)
		path, err := buildPath([]common.Address{token0, token1}, []int64{int64(swapInfo.TokenOutIndex)})
		if err != nil {
			return nil, err
		}

		if swapInfo.IsExactInOut {
			//exact in
			swapExactInData := bitcoinABI.ISwapRouterExactInputParams{
				Path: path,
				AmountIn: big.NewInt(0), // todo:
				AmountOutMinimum: big.NewInt(0), // todo:
				Recipient: common.HexToAddress(to),
				Deadline: big.NewInt(1000000000),
			}

			calldata, err = routerABI.Pack("exactInput", swapExactInData)
			if err != nil {
				return nil, err
			}
	
		} else {
			//exact out	
			swapExactOutData := bitcoinABI.ISwapRouterExactOutputParams{
				Path: path,
				AmountOut: amount0,
				Recipient: common.HexToAddress(to), // todo: review
				Deadline: big.NewInt(1000000000),
				AmountInMaximum: big.NewInt(0), // todo:
			}

			calldata, err = nonfungibleABI.Pack("exactOutput", swapExactOutData)
			if err != nil {
				return nil, err
			}
	
		}
	case multisig.RadFiCollectFeesMsg:
		// todo:

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

	return finalMessage, nil
}

func mulDiv(a, nNumerator, nDenominator *big.Int) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(a, nNumerator), nDenominator)
}

func buildPath(paths []common.Address, fees []int64) ([]byte, error) {
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