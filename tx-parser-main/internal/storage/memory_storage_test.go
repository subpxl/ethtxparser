package storage

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	if storage == nil {
		t.Error("Expected non-nil MemoryStorage")
	}
	if storage.currentBlock != 0 {
		t.Errorf("Expected initial block 0, got %d", storage.currentBlock)
	}
	if len(storage.transactions) != 0 {
		t.Error("Expected empty transactions map")
	}
	if len(storage.subscribers) != 0 {
		t.Error("Expected empty subscribers map")
	}
}

func TestSubscriberOperations(t *testing.T) {
	storage := NewMemoryStorage()

	testCases := []struct {
		name      string
		address   string
		expectAdd bool
	}{
		{
			name:      "valid address",
			address:   "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			expectAdd: true,
		},
		{
			name:      "invalid address - no 0x prefix",
			address:   "742d35Cc6634C0532925a3b844Bc454e4438f44e",
			expectAdd: false,
		},
		{
			name:      "invalid address - wrong length",
			address:   "0x742d35Cc",
			expectAdd: false,
		},
		{
			name:      "duplicate address",
			address:   "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			expectAdd: true, // should still return true for existing address
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := storage.AddSubscriber(tc.address)
			if result != tc.expectAdd {
				t.Errorf("Expected AddSubscriber result %v, got %v", tc.expectAdd, result)
			}

			if tc.expectAdd {
				if !storage.IsSubscribed(tc.address) {
					t.Error("Expected address to be subscribed")
				}
			}
		})
	}

	// Test GetSubscribers
	subscribers := storage.GetSubscribers()
	if len(subscribers) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", len(subscribers))
	}
}

func TestTransactionOperations(t *testing.T) {
	storage := NewMemoryStorage()

	tx := Transaction{
		Hash:        "0x123",
		FromAddress: "0xabc",
		ToAddress:   "0xdef",
		Value:       1.0,
		BlockNumber: 100,
		Timestamp:   1000,
	}

	// Test storing transaction
	storage.StoreTransaction(tx)

	// Test retrieving transactions for sender
	fromTxs := storage.GetTransactions(tx.FromAddress)
	if len(fromTxs) != 1 {
		t.Errorf("Expected 1 transaction for sender, got %d", len(fromTxs))
	}
	if !reflect.DeepEqual(fromTxs[0], tx) {
		t.Error("Retrieved transaction doesn't match stored transaction")
	}

	// Test retrieving transactions for receiver
	toTxs := storage.GetTransactions(tx.ToAddress)
	if len(toTxs) != 1 {
		t.Errorf("Expected 1 transaction for receiver, got %d", len(toTxs))
	}
	if !reflect.DeepEqual(toTxs[0], tx) {
		t.Error("Retrieved transaction doesn't match stored transaction")
	}

	// Test case sensitivity
	upperFromTxs := storage.GetTransactions(strings.ToUpper(tx.FromAddress))
	if len(upperFromTxs) != 1 {
		t.Error("Case sensitivity affected transaction retrieval")
	}
}

func TestBlockOperations(t *testing.T) {
	storage := NewMemoryStorage()

	testBlocks := []int64{0, 1, 100, 1000}

	for _, block := range testBlocks {
		storage.UpdateCurrentBlock(block)
		if got := storage.GetCurrentBlock(); got != block {
			t.Errorf("Expected current block %d, got %d", block, got)
		}
	}
}

func TestEmptyStorage(t *testing.T) {
	storage := NewMemoryStorage()

	// Test getting transactions for non-existent address
	txs := storage.GetTransactions("0x123")
	if txs != nil {
		t.Error("Expected nil transactions for non-existent address")
	}

	// Test checking subscription for non-existent address
	if storage.IsSubscribed("0x123") {
		t.Error("Expected false for non-existent subscriber")
	}
}

func TestTransactionWithEmptyToAddress(t *testing.T) {
	storage := NewMemoryStorage()

	tx := Transaction{
		Hash:        "0x123",
		FromAddress: "0xabc",
		ToAddress:   "", // Empty to address
		Value:       1.0,
		BlockNumber: 100,
		Timestamp:   1000,
	}

	storage.StoreTransaction(tx)

	// Should only store for FromAddress
	fromTxs := storage.GetTransactions(tx.FromAddress)
	if len(fromTxs) != 1 {
		t.Errorf("Expected 1 transaction for sender, got %d", len(fromTxs))
	}
}
