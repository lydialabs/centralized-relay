package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
)

func startMaster(c *Config, p *Provider) {
	http.HandleFunc("/execute", handleExecute)
	http.HandleFunc("/broadcast-request", func(w http.ResponseWriter, r *http.Request) {
		handleBroadcastNewRequest(w, r, p)
	})
	port := c.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	log.Printf("Master starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("x-api-key")
	if apiKey == "" {
		http.Error(w, "Missing API Key", http.StatusUnauthorized)
		return
	}
	apiKeyHeader := os.Getenv("API_KEY")
	if apiKey != apiKeyHeader {
		http.Error(w, "Invalid API Key", http.StatusForbidden)
		return
	}

	var msg string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "msg": msg}
	json.NewEncoder(w).Encode(response)
}

func handleBroadcastNewRequest(w http.ResponseWriter, r *http.Request, p *Provider) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("x-api-key")
	if apiKey == "" {
		http.Error(w, "Missing API Key", http.StatusUnauthorized)
		return
	}
	apiKeyHeader := os.Getenv("API_KEY")
	if apiKey != apiKeyHeader {
		http.Error(w, "Invalid API Key", http.StatusForbidden)
		return
	}

	// var msg string
	var msg slaveNewRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success", "msg": msg.RawTransaction}
	json.NewEncoder(w).Encode(response)

	// todo: validation raw tx
	// store the request to db 
	storeData, err := hex.DecodeString(msg.RawTransaction)
	if err != nil {
		p.logger.Error("Error decoding hex string", zap.Error(err))
		http.Error(w, "Error decoding hex string", http.StatusInternalServerError)
		return
	}
	err = p.StoreNewPendingRequest(storeData)
	if err != nil {
		http.Error(w, "Failed to store new request", http.StatusBadRequest)
		return
	}
	
	// route
	route := "/add-new-request"

	// fwd data to slavers
	sendRequestToSlaver(p.cfg.ApiKey, p.cfg.SlaveServer1+route, msg, p.logger)
	sendRequestToSlaver(p.cfg.ApiKey, p.cfg.SlaveServer2+route, msg, p.logger)
}

func sendRequestToSlaver(apiKey, url string, msg slaveNewRequest, logger *zap.Logger) error {
	client := &http.Client{}
	requestData, _ := json.Marshal(msg)
	payload := bytes.NewBuffer(requestData)

	maxRetry := 0
	for {
		req, err := http.NewRequest("POST", url, payload)
		if err != nil {
			logger.Error("sendRequestToSlaver: failed to create request: ", zap.Error(err))
			return err
		}

		req.Header.Add("x-api-key", apiKey)
		_, err = client.Do(req)
		if err != nil {
			logger.Error("sendRequestToSlaver: failed to send request: ", zap.Error(err))
			// retry 
			maxRetry++
			if maxRetry == 5 {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func requestPartialSign(apiKey string, url string, slaveRequestData []byte, responses chan<- slaveResponse, order int, wg *sync.WaitGroup) {
	defer wg.Done()
	response := slaveResponse{}
	client := &http.Client{}
	payload := bytes.NewBuffer(slaveRequestData)
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		response.err = fmt.Errorf("failed to create request: %v", err)
		responses <- response
		return
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)

	if err != nil {
		response.err = fmt.Errorf("failed to send request: %v", err)
		responses <- response
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.err = fmt.Errorf("error reading response: %v", err)
		responses <- response
		return
	}

	sigs := [][]byte{}
	err = json.Unmarshal(body, &sigs)
	if err != nil {
		response.err = fmt.Errorf("err Unmarshal: %v", err)
		responses <- response
		return
	}

	response.order = order
	response.sigs = sigs
	response.err = nil
	responses <- response
}
