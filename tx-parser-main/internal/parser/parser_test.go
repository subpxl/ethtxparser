package parser

import (
	"blockchain-parser/internal/storage"
	"testing"
)

// MockStorage implements storage.StorageInterface for testing
type MockStorage struct {
	currentBlock int64
	subscribers  map[string]bool
	transactions []storage.Transaction
}

func newMockStorage() *MockStorage {
	return &MockStorage{
		subscribers: make(map[string]bool),
	}
}

func (m *MockStorage) GetCurrentBlock() int64         { return m.currentBlock }
func (m *MockStorage) UpdateCurrentBlock(block int64) { m.currentBlock = block }
func (m *MockStorage) AddSubscriber(address string) bool {
	if _, exists := m.subscribers[address]; exists {
		return false
	}
	m.subscribers[address] = true
	return true
}
func (m *MockStorage) IsSubscribed(address string) bool { return m.subscribers[address] }
func (m *MockStorage) GetSubscribers() []string {
	subs := make([]string, 0, len(m.subscribers))
	for addr := range m.subscribers {
		subs = append(subs, addr)
	}
	return subs
}
func (m *MockStorage) GetTransactions(address string) []storage.Transaction { return m.transactions }
func (m *MockStorage) StoreTransaction(tx storage.Transaction) {
	m.transactions = append(m.transactions, tx)
}

func TestNewParser(t *testing.T) {
	mockStorage := newMockStorage()
	mockRPC := &RPCClient{}
	parser := NewParser(mockStorage, mockRPC)

	if parser == nil {
		t.Error("Expected non-nil Parser")
	}
}

func TestSubscribeOperations(t *testing.T) {
	mockStorage := newMockStorage()
	parser := NewParser(mockStorage, nil)

	testCases := []struct {
		name      string
		address   string
		expectAdd bool
	}{
		{
			name:      "new subscriber",
			address:   "0x123",
			expectAdd: true,
		},
		{
			name:      "duplicate subscriber",
			address:   "0x123",
			expectAdd: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parser.Subscribe(tc.address)
			if result != tc.expectAdd {
				t.Errorf("Expected Subscribe result %v, got %v", tc.expectAdd, result)
			}

			if !parser.IsSubscribed(tc.address) {
				t.Error("Expected address to be subscribed")
			}
		})
	}
}

func TestProcessTransaction(t *testing.T) {
	mockStorage := newMockStorage()
	parser := NewParser(mockStorage, nil)

	// Add a subscriber
	parser.Subscribe("0x123")

	testCases := []struct {
		name          string
		tx            map[string]interface{}
		timestamp     int64
		expectSuccess bool
	}{
		{
			name: "valid transaction",
			tx: map[string]interface{}{
				"hash":        "0xabc",
				"from":        "0x123",
				"to":          "0x456",
				"value":       "0x1",
				"blockNumber": "0x1",
			},
			timestamp:     1000,
			expectSuccess: true,
		},
		{
			name: "invalid transaction data",
			tx: map[string]interface{}{
				"from": "0x123",
				// missing required fields
			},
			timestamp:     1000,
			expectSuccess: false,
		},
		{
			name: "invalid hex value",
			tx: map[string]interface{}{
				"hash":        "0xabc",
				"from":        "0x123",
				"to":          "0x456",
				"value":       "invalid",
				"blockNumber": "0x1",
			},
			timestamp:     1000,
			expectSuccess: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := parser.ProcessTransaction(tc.tx, tc.timestamp)
			if tc.expectSuccess {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tx == nil {
					t.Error("Expected non-nil transaction")
				}
			} else {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			}
		})
	}
}

func TestParseHexToInt64(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		{
			name:        "valid hex",
			input:       "a",
			expected:    10,
			expectError: false,
		},
		{
			name:        "invalid hex",
			input:       "invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseHexToInt64(tc.input)
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tc.expectError && result != tc.expected {
				t.Errorf("Expected %d but got %d", tc.expected, result)
			}
		})
	}
}
