// Package parser implements blockchain transaction parsing and processing functionality.
package parser

import (
	"blockchain-parser/internal/storage"
	"blockchain-parser/internal/utils"
	"fmt"
	"strconv"
	"strings"
)

// Parser defines the interface for blockchain transaction parsing operations
type Parser interface {

	// GetCurrentBlock returns the latest processed block number
	GetCurrentBlock() int64

	// UpdateCurrentBlock updates the latest processed block number
	UpdateCurrentBlock(blockNumber int64)

	// Subscribe adds an address to monitor, returns false if address is invalid
	Subscribe(address string) bool

	// IsSubscribed checks if an address is being monitored
	IsSubscribed(address string) bool

	// GetSubscribers returns all monitored addresses
	GetSubscribers() []string

	// GetTransactions returns all transactions for a given address
	GetTransactions(address string) []storage.Transaction

	// ProcessTransaction processes a raw transaction and stores it if relevant
	ProcessTransaction(tx map[string]interface{}, blockTimestamp int64) (*storage.Transaction, error)
}

// parserImpl implements the Parser interface
type parserImpl struct {
	storage   storage.StorageInterface
	rpcClient *RPCClient
}

// NewParser creates a new Parser instance with the given storage and RPC client
func NewParser(storage storage.StorageInterface, rpcClient *RPCClient) Parser {
	return &parserImpl{
		storage:   storage,
		rpcClient: rpcClient,
	}
}

// func (p *parserImpl) GetCurrentBlock() int64 {
// 	return p.storage.GetCurrentBlock()
// }

func (p *parserImpl) GetCurrentBlock() int64 {
	return p.storage.GetCurrentBlock()
}

func (p *parserImpl) UpdateCurrentBlock(blockNumber int64) {
	p.storage.UpdateCurrentBlock(blockNumber)
}

func (p *parserImpl) Subscribe(address string) bool {
	return p.storage.AddSubscriber(address)
}

func (p *parserImpl) IsSubscribed(address string) bool {
	return p.storage.IsSubscribed(address)
}

func (p *parserImpl) GetSubscribers() []string {
	return p.storage.GetSubscribers()
}

func (p *parserImpl) GetTransactions(address string) []storage.Transaction {
	return p.storage.GetTransactions(address)
}

func (p *parserImpl) ProcessTransaction(tx map[string]interface{}, blockTimestamp int64) (*storage.Transaction, error) {
	// Check if required fields exist and are not nil
	if tx == nil {
		return nil, fmt.Errorf("transaction data is nil")
	}

	// Safely get values with type checking
	from, ok := tx["from"].(string)
	if !ok || from == "" {
		return nil, fmt.Errorf("invalid or missing 'from' address")
	}

	value, ok := tx["value"].(string)
	if !ok || value == "" {
		return nil, fmt.Errorf("invalid or missing 'value'")
	}

	hash, ok := tx["hash"].(string)
	if !ok || hash == "" {
		return nil, fmt.Errorf("invalid or missing 'hash'")
	}

	blockNumber, ok := tx["blockNumber"].(string)
	if !ok || blockNumber == "" {
		return nil, fmt.Errorf("invalid or missing 'blockNumber'")
	}

	// Handle optional 'to' address (could be nil for contract creation)
	var toAddress string
	if to, ok := tx["to"].(string); ok && to != "" {
		toAddress = strings.ToLower(to)
	}

	// Parse values
	valueHex := strings.TrimPrefix(value, "0x")
	valueInt64, err := parseHexToInt64(valueHex)
	if err != nil {
		return nil, fmt.Errorf("error parsing transaction value: %v", err)
	}
	valueETH := float64(valueInt64) / 1e18

	blockNumberHex := strings.TrimPrefix(blockNumber, "0x")
	blockNum, err := parseHexToInt64(blockNumberHex)
	if err != nil {
		return nil, fmt.Errorf("error parsing block number: %v", err)
	}

	// Create transaction object
	transaction := storage.Transaction{
		Hash:        hash,
		FromAddress: strings.ToLower(from),
		ToAddress:   toAddress,
		Value:       valueETH,
		BlockNumber: blockNum,
		Timestamp:   blockTimestamp,
	}

	// Check if we should store this transaction
	if p.storage.IsSubscribed(transaction.FromAddress) ||
		(transaction.ToAddress != "" && p.storage.IsSubscribed(transaction.ToAddress)) {
		p.storage.StoreTransaction(transaction)
		return &transaction, nil
	}

	return nil, nil
}

func parseHexToInt64(hex string) (int64, error) {

	return utils.String2Int64(hex, 16) // Use the utility function
}

func String2Int64(s string, base int) (int64, error) {
	return strconv.ParseInt(s, base, 64)
}
