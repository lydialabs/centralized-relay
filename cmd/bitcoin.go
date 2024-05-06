package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func btcCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bitcoin",
		Aliases: []string{"btc"},
		Short:   "Command line of bitcoin",
	}

	cmd.AddCommand(
		syncBlockEvent(a),
	)

	return cmd
}

func syncBlockEvent(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"sync"},
		Short:   "Sync block on-chain of bitcoin",
		Args:    withUsage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := "https://mempool.space/testnet/api/block-height/2633011"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error fetching URL: ", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response: ", err)
			}
			txsHashOfBlock := string(body)
			url = "https://mempool.space/testnet/api/block/" + txsHashOfBlock + "/txids"
			resp, err = http.Get(url)
			if err != nil {
				fmt.Println("Error fetching URL: ", err)
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response: ", err)
			}
			var txs []string
			if err := json.Unmarshal(body, &txs); err != nil {
				fmt.Println("Error decoding JSON:", err)
			}
			var wg sync.WaitGroup
			client := &http.Client{Timeout: 10 * time.Second} // Set a timeout
			for _, tx := range txs {
				wg.Add(1)
				go fetchTransactionData(tx, client, &wg)
			}
			wg.Wait()
			fmt.Println("Done")
			return nil
		},
	}
	return yamlFlag(a.viper, jsonFlag(a.viper, cmd))
}

func fetchTransactionData(txID string, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the wait group counter is decremented

	url := "https://mempool.space/testnet/api/tx/" + txID
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close() // Close the response body right after checking the error

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var transactionData Transaction
	if err := json.Unmarshal(body, &transactionData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(transactionData.Vout) >= 2 && transactionData.Vout[1].ScriptPubKeyAddress == "tb1qntuupajv5k9hr92h3ecx7uvy6gprjeltnwlp8g" {
		fmt.Println(transactionData.Vout[1])
	}
}
