// Package parser provides functionality for interacting with blockchain nodes via JSON-RPC.
package parser

import (
	"blockchain-parser/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RPCClient handles communication with a blockchain node's JSON-RPC API
type RPCClient struct {
	config config.NetworkConfig
	client *http.Client
	url    string
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}


// BlockResponse represents the Ethereum block data
type BlockResponse struct {
    Number          string   `json:"number"`
    Hash            string   `json:"hash"`
    ParentHash      string   `json:"parentHash"`
    Nonce           string   `json:"nonce"`
    Sha3Uncles      string   `json:"sha3Uncles"`
    LogsBloom       string   `json:"logsBloom"`
    TransactionsRoot string   `json:"transactionsRoot"`
    StateRoot       string   `json:"stateRoot"`
    ReceiptsRoot    string   `json:"receiptsRoot"`
    Miner           string   `json:"miner"`
    Difficulty      string   `json:"difficulty"`
    TotalDifficulty string   `json:"totalDifficulty"`
    ExtraData       string   `json:"extraData"`
    Size            string   `json:"size"`
    GasLimit        string   `json:"gasLimit"`
    GasUsed         string   `json:"gasUsed"`
    Timestamp       string   `json:"timestamp"`
    Transactions    []string `json:"transactions"`
    Uncles          []string `json:"uncles"`
}

// NewRPCClient creates a new RPC client instance for the given URL
func NewRPCClient(cfg *config.Config) *RPCClient {
	return &RPCClient{
		config: cfg.Network,
		url:    cfg.RPCEndpoint,
		client: &http.Client{
			Timeout: cfg.Network.RequestTimeout,
		},
	}
}

// MakeCall sends a JSON-RPC request to the blockchain node
func (rc *RPCClient) MakeCall(method string, params []interface{}) (*JSONRPCResponse, error) {
	payload := JSONRPCRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(rc.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("RPC error: %v", result.Error)
	}

	return &result, nil
}
