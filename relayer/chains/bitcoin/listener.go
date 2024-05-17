package bitcoin

import (
	"fmt"
	"log"
	"sync"
)

func ListenerBitcoin(config *Config) {
	blockNumber := config.StartBlock

	// Get the block hash for the given block number.
	blockHash, err := GetBlockHash(config, blockNumber)
	if err != nil {
		log.Fatalf("Failed to get block hash: %v", err)
	}

	// Get the block for the given block hash.
	block, err := GetBlock(config, blockHash)
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	// Print the block details.
	fmt.Printf("Block number: %d\n", blockNumber)
	fmt.Printf("New block detected: %d\n", blockNumber)

	txs, _ := block["tx"].([]interface{})
	fmt.Printf("Length transaction: %d\n", len(txs))

	var wg sync.WaitGroup
	for _, tx := range txs {
		txIDStr, ok := tx.(string)
		if !ok {
			log.Println("Invalid transaction ID format")
			continue
		}
		wg.Add(1)
		go CheckTransaction(config, txIDStr, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}
