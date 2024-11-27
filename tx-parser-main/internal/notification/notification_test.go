package notification

import (
	"blockchain-parser/internal/storage"
	"testing"
	"time"
)

func TestNewConsoleNotificationService(t *testing.T) {
	service := NewConsoleNotificationService()
	if service == nil {
		t.Error("Expected non-nil NotificationService")
	}
}

func TestConsoleNotificationService_Notify(t *testing.T) {
	testCases := []struct {
		name         string
		notification Notification
	}{
		{
			name: "incoming transaction",
			notification: Notification{
				Type:    TransactionReceived,
				Address: "0x123",
				Transaction: storage.Transaction{
					Hash:        "0xabc",
					FromAddress: "0x456",
					ToAddress:   "0x123",
					Value:       1.0,
					BlockNumber: 100,
					Timestamp:   time.Now().Unix(),
				},
			},
		},
		{
			name: "outgoing transaction",
			notification: Notification{
				Type:    TransactionSent,
				Address: "0x456",
				Transaction: storage.Transaction{
					Hash:        "0xdef",
					FromAddress: "0x456",
					ToAddress:   "0x789",
					Value:       2.0,
					BlockNumber: 101,
					Timestamp:   time.Now().Unix(),
				},
			},
		},
	}

	service := NewConsoleNotificationService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.Notify(tc.notification)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestNotificationTypes(t *testing.T) {
	if TransactionReceived != "TRANSACTION_RECEIVED" {
		t.Errorf("Expected TransactionReceived to be 'TRANSACTION_RECEIVED', got %s", TransactionReceived)
	}
	if TransactionSent != "TRANSACTION_SENT" {
		t.Errorf("Expected TransactionSent to be 'TRANSACTION_SENT', got %s", TransactionSent)
	}
}
