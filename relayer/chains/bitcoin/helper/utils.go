package helper

import (
	"golang.org/x/crypto/sha3"
	"crypto/ecdsa"
	"math/big"
	"errors"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type RPCError struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

type RPCBaseRes struct {
	Id       int       `json:"Id"`
	RPCError *RPCError `json:"Error"`
}

func Rawsha3(b []byte) []byte {
	hashF := sha3.NewLegacyKeccak256()
	hashF.Write(b)
	buf := hashF.Sum(nil)
	return buf
}

type ValidatorKeyInfo struct {
	Address string
	PrivKey *ecdsa.PrivateKey
	Big     *big.Int
}

// Implement the sort.Interface for the ValidatorKeyInfo type
type ValidatorKeyInfoSlice []ValidatorKeyInfo

func (s ValidatorKeyInfoSlice) Len() int {
	return len(s)
}

func (s ValidatorKeyInfoSlice) Less(i, j int) bool {
	return s[i].Big.Cmp(s[j].Big) == -1
}

func (s ValidatorKeyInfoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type EstimateRes struct {
	RPCBaseRes
	Result interface{} `json:"Result"`
}

func ToCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}


func EstimateGas(msg ethereum.CallMsg, fullnode string) (uint64, error) {
	var resp EstimateRes
	rpcClient := NewRPCClient()
	params := []interface{}{
		ToCallArg(msg),
		"latest",
	}
	err := rpcClient.RPCCall(
		"",
		fullnode,
		"",
		"eth_estimateGas",
		params,
		&resp,
	)
	if err != nil {
		return 0, err
	}

	if resp.RPCError != nil {
		return 0, errors.New(resp.RPCError.Message)
	}

	hexString := resp.Result.(string)

	n := new(big.Int)
	n.SetString(hexString[2:], 16)
	if n.Cmp(big.NewInt(0)) == 0 {
		return 0, errors.New("call estimate gas got error")
	}
	return n.Uint64(), nil
}