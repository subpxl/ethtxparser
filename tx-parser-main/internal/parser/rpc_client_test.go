package parser

import (
    "blockchain-parser/config"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestNewRPCClient(t *testing.T) {
    // Create test configuration
    cfg := &config.Config{
        RPCEndpoint: "http://localhost:8545",
        Network: config.NetworkConfig{
            RequestTimeout: 5 * time.Second,
            RetryAttempts: 3,
            RetryDelay:    time.Second,
            RateLimitDelay: 100 * time.Millisecond,
        },
    }

    client := NewRPCClient(cfg)
    
    if client == nil {
        t.Error("Expected non-nil RPCClient")
    }
    if client.url != cfg.RPCEndpoint {
        t.Errorf("Expected RPC URL %s, got %s", cfg.RPCEndpoint, client.url)
    }
    if client.config.RequestTimeout != cfg.Network.RequestTimeout {
        t.Errorf("Expected timeout %v, got %v", cfg.Network.RequestTimeout, client.config.RequestTimeout)
    }
}

func TestMakeCall(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        params         []interface{}
        serverResponse *JSONRPCResponse
        serverStatus   int
        expectError    bool
    }{
        {
            name:   "successful eth_blockNumber call",
            method: "eth_blockNumber",
            params: nil,
            serverResponse: &JSONRPCResponse{
                Result: "0x1234",
            },
            serverStatus: http.StatusOK,
            expectError:  false,
        },
        {
            name:   "successful eth_getBlockByNumber call",
            method: "eth_getBlockByNumber",
            params: []interface{}{"0x1", true},
            serverResponse: &JSONRPCResponse{
                Result: map[string]interface{}{
                    "number": "0x1",
                    "hash":   "0x1234",
                    "transactions": []interface{}{},
                    "timestamp":    "0x123456",
                },
            },
            serverStatus: http.StatusOK,
            expectError:  false,
        },
        {
            name:   "rpc error response",
            method: "eth_blockNumber",
            params: nil,
            serverResponse: &JSONRPCResponse{
                Error: map[string]interface{}{
                    "code":    -32600,
                    "message": "Invalid request",
                },
            },
            serverStatus: http.StatusOK,
            expectError:  true,
        },
        {
            name:           "server error",
            method:        "eth_blockNumber",
            params:        nil,
            serverStatus:  http.StatusInternalServerError,
            expectError:   true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create test server
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Verify request method and content type
                if r.Method != http.MethodPost {
                    t.Errorf("Expected POST request, got %s", r.Method)
                }
                if ct := r.Header.Get("Content-Type"); ct != "application/json" {
                    t.Errorf("Expected Content-Type application/json, got %s", ct)
                }

                // Verify request body
                var request JSONRPCRequest
                if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
                    t.Errorf("Failed to decode request: %v", err)
                    return
                }

                // Verify request fields
                if request.JsonRPC != "2.0" {
                    t.Errorf("Expected JSON-RPC 2.0, got %s", request.JsonRPC)
                }
                if request.Method != tc.method {
                    t.Errorf("Expected method %s, got %s", tc.method, request.Method)
                }

                // Set response status and body
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(tc.serverStatus)
                if tc.serverResponse != nil {
                    json.NewEncoder(w).Encode(tc.serverResponse)
                }
            }))
            defer server.Close()

            // Create client with test configuration
            cfg := &config.Config{
                RPCEndpoint: server.URL,
                Network: config.NetworkConfig{
                    RequestTimeout: time.Second,
                    RetryAttempts: 1,
                    RetryDelay:    time.Millisecond,
                    RateLimitDelay: time.Millisecond,
                },
            }
            client := NewRPCClient(cfg)

            // Make the RPC call
            response, err := client.MakeCall(tc.method, tc.params)

            // Check results
            if tc.expectError {
                if err == nil {
                    t.Error("Expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error: %v", err)
                }
                if response == nil {
                    t.Error("Expected non-nil response")
                    return
                }

                // Compare response with expected
                if tc.serverResponse != nil {
                    expectedJSON, _ := json.Marshal(tc.serverResponse)
                    actualJSON, _ := json.Marshal(response)
                    if string(actualJSON) != string(expectedJSON) {
                        t.Errorf("Expected response %s, got %s", expectedJSON, actualJSON)
                    }
                }
            }
        })
    }
}

func TestRetryLogic(t *testing.T) {
    attempts := 0
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attempts++
        if attempts < 3 {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        json.NewEncoder(w).Encode(&JSONRPCResponse{Result: "0x1234"})
    }))
    defer server.Close()

    cfg := &config.Config{
        RPCEndpoint: server.URL,
        Network: config.NetworkConfig{
            RequestTimeout: time.Second,
            RetryAttempts: 3,
            RetryDelay:    time.Millisecond,
            RateLimitDelay: time.Millisecond,
        },
    }
    client := NewRPCClient(cfg)

    response, err := client.MakeCall("eth_blockNumber", nil)
    if err != nil {
        t.Errorf("Unexpected error after retries: %v", err)
    }
    if attempts != 3 {
        t.Errorf("Expected 3 attempts, got %d", attempts)
    }
    if response == nil {
        t.Error("Expected non-nil response after successful retry")
    }
}