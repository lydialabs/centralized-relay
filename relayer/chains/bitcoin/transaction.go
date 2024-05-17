package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/icon-project/centralized-relay/relayer/chains/bitcoin/types"
)

// GetBlockHash retrieves the block hash for a given block number.
func GetBlockHash(config *Config, blockNumber int) (string, error) {
	reqBody := types.RPCRequest{
		Jsonrpc: "1.0",
		ID:      1,
		Method:  "getblockhash",
		Params:  []interface{}{blockNumber},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(config.BitcoinRPC, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var rpcResp types.RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return "", err
	}

	if rpcResp.Error != nil {
		return "", fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	var blockHash string
	if err := json.Unmarshal(rpcResp.Result, &blockHash); err != nil {
		return "", err
	}

	return blockHash, nil
}

// GetBlock retrieves the block for a given block hash.
func GetBlock(config *Config, blockHash string) (map[string]interface{}, error) {
	reqBody := types.RPCRequest{
		Jsonrpc: "1.0",
		ID:      1,
		Method:  "getblock",
		Params:  []interface{}{blockHash},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(config.BitcoinRPC, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp types.RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	var block map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &block); err != nil {
		return nil, err
	}
	return block, nil
}

// GetTransaction retrieves the transaction for a given transaction ID.
func GetTransaction(config *Config, txID string) (map[string]interface{}, error) {
	reqBody := types.RPCRequest{
		Jsonrpc: "1.0",
		ID:      1,
		Method:  "getrawtransaction",
		Params:  []interface{}{txID, 1}, // 1 means get the decoded transaction
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(config.BitcoinRPC, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp types.RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	var tx map[string]interface{}
	if err := json.Unmarshal(rpcResp.Result, &tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// CheckTransaction checks if a transaction has a Vout with the target address.
func CheckTransaction(config *Config, txIDStr string, wg *sync.WaitGroup) {
	defer wg.Done()
	tx, err := GetTransaction(config, txIDStr)
	if err != nil {
		log.Printf("Error getting transaction %s: %v", txIDStr, err)
		return
	}
	vouts, ok := tx["vout"].([]interface{})
	if !ok {
		log.Println("Invalid Vout format")
		return
	}

	for _, v := range vouts {
		vout, ok := v.(map[string]interface{})
		if !ok {
			log.Println("Invalid Vout format")
			continue
		}
		scriptPubKey, ok := vout["scriptPubKey"].(map[string]interface{})
		if !ok {
			log.Println("Invalid scriptPubKey format")
			continue
		}
		address, ok := scriptPubKey["address"].(string)
		if ok {
			if address == config.AdminWallet {
				fmt.Printf("Transaction ID: %s\n", tx["txid"])
				fmt.Printf("Output: %v\n", vout)
			}
			return
		} else {
			fmt.Println("Not found")
			fmt.Println(vout)
			return
		}

	}
}
