// Package notification provides functionality for notifying users about blockchain transactions.
package notification

import (
	"blockchain-parser/internal/logger"
	"blockchain-parser/internal/storage"
	"fmt"
)

// NotificationType defines the type of transaction notification
type NotificationType string

const (
	// TransactionReceived indicates an incoming transaction
	TransactionReceived NotificationType = "TRANSACTION_RECEIVED"
	
	// TransactionSent indicates an outgoing transaction
	TransactionSent     NotificationType = "TRANSACTION_SENT"
)

// Notification represents a transaction notification with all relevant details
type Notification struct {
	Type        NotificationType
	Address     string
	Transaction storage.Transaction
	Timestamp   int64
}

// NotificationService defines the interface for notification delivery
type NotificationService interface {
	Notify(notification Notification) error
}

// ConsoleNotificationService implements NotificationService
type ConsoleNotificationService struct{}

func NewConsoleNotificationService() NotificationService {
	return &ConsoleNotificationService{}
}

// NewConsoleNotificationService creates a new console notification service
func (s *ConsoleNotificationService) Notify(n Notification) error {
	var direction string
	if n.Type == TransactionReceived {
		direction = "Incoming"
	} else {
		direction = "Outgoing"
	}

	fmt.Printf("\n=== %s Transaction Notification ===\n", direction)
	fmt.Printf("Address: %s\n", n.Address)
	fmt.Printf("Transaction Hash: %s\n", n.Transaction.Hash)
	fmt.Printf("From: %s\n", n.Transaction.FromAddress)
	fmt.Printf("To: %s\n", n.Transaction.ToAddress)
	fmt.Printf("Value: %f ETH\n", n.Transaction.Value)
	fmt.Printf("Block Number: %d\n", n.Transaction.BlockNumber)
	fmt.Printf("================================\n\n")

	logger.Info("Notification sent for %s transaction to address %s", direction, n.Address)
	return nil
}
