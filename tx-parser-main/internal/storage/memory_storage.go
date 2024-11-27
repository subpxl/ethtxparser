// Package storage provides interfaces and implementations for storing blockchain transaction data

package storage

import (
	"strings"
)

// Transaction represents a blockchain transaction with its key details
type Transaction struct {
	Hash        string
	FromAddress string
	ToAddress   string
	Value       float64
	BlockNumber int64
	Timestamp   int64
}

// StorageInterface defines the required methods for a storage implementation
type StorageInterface interface {
	StoreTransaction(transaction Transaction)
	GetTransactions(address string) []Transaction
	AddSubscriber(address string) bool
	IsSubscribed(address string) bool
	GetSubscribers() []string
	UpdateCurrentBlock(blockNumber int64)
	GetCurrentBlock() int64
}

// MemoryStorage implements StorageInterface using in-memory data structures
type MemoryStorage struct {
	transactions map[string][]Transaction
	subscribers  map[string]bool
	currentBlock int64
}

// NewMemoryStorage creates and initializes a new MemoryStorage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: make(map[string][]Transaction),
		subscribers:  make(map[string]bool),
		currentBlock: 0,
	}
}

// Store a transaction
func (ms *MemoryStorage) StoreTransaction(transaction Transaction) {
	address := strings.ToLower(transaction.FromAddress)
	if _, exists := ms.transactions[address]; !exists {
		ms.transactions[address] = []Transaction{}
	}
	ms.transactions[address] = append(ms.transactions[address], transaction)

	if transaction.ToAddress != "" {
		address = strings.ToLower(transaction.ToAddress)
		if _, exists := ms.transactions[address]; !exists {
			ms.transactions[address] = []Transaction{}
		}
		ms.transactions[address] = append(ms.transactions[address], transaction)
	}
}

func (ms *MemoryStorage) GetTransactions(address string) []Transaction {
	return ms.transactions[strings.ToLower(address)]
}

// AddSubscriber adds subscriber to memorystorage
func (ms *MemoryStorage) AddSubscriber(address string) bool {
	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		return false
	}
	ms.subscribers[strings.ToLower(address)] = true
	return true
}

// IsSubscribed for address is subscribed
func (ms *MemoryStorage) IsSubscribed(address string) bool {
	_, exists := ms.subscribers[strings.ToLower(address)]
	return exists
}

// GetSubscribers gets all subscribers
func (ms *MemoryStorage) GetSubscribers() []string {
	subscribers := make([]string, 0, len(ms.subscribers))
	for addr := range ms.subscribers {
		subscribers = append(subscribers, addr)
	}
	return subscribers
}

// UpdateCurrentBlock  Updates current block
func (ms *MemoryStorage) UpdateCurrentBlock(blockNumber int64) {
	ms.currentBlock = blockNumber
}

// GetCurrentBlock gets current block
func (ms *MemoryStorage) GetCurrentBlock() int64 {
	return ms.currentBlock
}
