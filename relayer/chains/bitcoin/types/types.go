package types

import "encoding/json"

// RPCRequest represents a JSON-RPC request.
type RPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// RPCResponse represents a JSON-RPC response.
type RPCResponse struct {
	Jsonrpc string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Result  json.RawMessage   `json:"result"`
	Error   *RPCResponseError `json:"error"`
}

// RPCResponseError represents an error in a JSON-RPC response.
type RPCResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
