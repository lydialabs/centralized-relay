package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RPCClient struct {
	*http.Client
}

// NewHttpClient to get http client instance
func NewRPCClient() *RPCClient {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	return &RPCClient{
		httpClient,
	}
}

func buildRPCServerAddress(protocol string, host string, port string) string {
	url := host
	if protocol != "" {
		url = protocol + "://" + url
	}
	if port != "" {
		url = url + ":" + port
	}
	return url
}

func (client *RPCClient) RPCCall(
	rpcProtocol string,
	rpcHost string,
	rpcPortStr string,
	method string,
	params interface{},
	rpcResponse interface{},
) (err error) {
	rpcEndpoint := buildRPCServerAddress(rpcProtocol, rpcHost, rpcPortStr)

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(rpcEndpoint, "application/json", bytes.NewBuffer(payloadInBytes))
	if err != nil {
		return err
	}
	respBody := resp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, rpcResponse)
	if err != nil {
		return err
	}
	return nil
}

func (client *RPCClient) MakeRequest(
	endpoint string,
	method string,
	params interface{},
	isGetMethod bool,
	dataResponse interface{},
) (interface{}, error) {

	resp := &http.Response{}
	var err error

	if isGetMethod {
		url := endpoint
		resp, err = client.Get(url)
		if err != nil {
			return nil, err
		}
	} else {
		payloadInBytes, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		log.Println("payloadInBytes", string(payloadInBytes))

		resp, err = client.Post(endpoint, "application/json", bytes.NewBuffer(payloadInBytes))
		if err != nil {
			return nil, err
		}
	}

	respBody := resp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, err
	}

	finalResp := struct {
		Error  interface{}
		status bool
		Data   interface{}
	}{
		Data: dataResponse,
	}

	err = json.Unmarshal(body, &finalResp)
	if err != nil {
		fmt.Printf("Unmarshal resp error: %v\n", string(body))
		return nil, err
	}
	return finalResp.Data, err
}
