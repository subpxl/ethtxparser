// Package monitor implements blockchain monitoring functionality for tracking
package monitor

import (
	"blockchain-parser/internal/logger"
	"blockchain-parser/internal/notification"
	"blockchain-parser/internal/parser"
	"blockchain-parser/internal/storage"
	"blockchain-parser/internal/utils"
	"fmt"
	"strings"
	"time"
)

// BlockMonitor watches for new blocks and processes their transactions
type BlockMonitor struct {
	parser    parser.Parser
	rpcClient *parser.RPCClient
	notifier  notification.NotificationService
}

// NewBlockMonitor creates a new block monitor instance
func NewBlockMonitor(p parser.Parser, rpc *parser.RPCClient, notifier notification.NotificationService) *BlockMonitor {
	return &BlockMonitor{
		parser:    p,
		rpcClient: rpc,
		notifier:  notifier,
	}
}

// StartMonitoring begins continuous monitoring of new blocks.
func (m *BlockMonitor) StartMonitoring() {
	fmt.Println("Starting block monitoring...")

	// Add rate limiting for public nodes
	rateLimiter := time.NewTicker(1 * time.Second)
	defer rateLimiter.Stop()

	for {
		select {
		case <-rateLimiter.C:
			m.processNewBlocks()
		}
	}
}

// processNewBlocks checks for and processes any new blocks
func (m *BlockMonitor) processNewBlocks() error {
	result, err := m.rpcClient.MakeCall("eth_blockNumber", nil)
	if err != nil {
		return NewMonitorError(ErrBlockNumberFetch, "Failed to fetch block number", err)
	}

	latestBlockHex, ok := result.Result.(string)
	if !ok {
		return NewMonitorError(ErrBlockNumberFetch, "Invalid block number response format", nil)
	}

	latestBlock, err := m.parseHexToInt64(strings.TrimPrefix(latestBlockHex, "0x"))
	if err != nil {
		return NewMonitorError(ErrBlockNumberParse, "Failed to parse block number", err)
	}

	currentBlock := m.parser.GetCurrentBlock()
	if latestBlock > currentBlock {
		logger.Info("Processing new block: %d (current: %d)", latestBlock, currentBlock)
		if err := m.processBlock(latestBlock); err != nil {
			return err
		}
	} else {
		logger.Debug("No new blocks to process. Current: %d, Latest: %d", currentBlock, latestBlock)
	}

	return nil
}

// processBlock processes a single block and its transactions
func (m *BlockMonitor) processBlock(blockNumber int64) error {
	logger.Debug("Fetching block details for block %d", blockNumber)

	blockResult, err := m.rpcClient.MakeCall("eth_getBlockByNumber",
		[]interface{}{fmt.Sprintf("0x%x", blockNumber), true})
	if err != nil {
		return NewMonitorError(ErrBlockFetch, fmt.Sprintf("Failed to fetch block %d", blockNumber), err)
	}

	block, ok := blockResult.Result.(map[string]interface{})
	if !ok {
		return NewMonitorError(ErrBlockFetch, "Invalid block response format", nil)
	}

	transactions, ok := block["transactions"].([]interface{})
	if !ok {
		logger.Warn("No transactions found in block %d", blockNumber)
		return nil
	}

	timestampHex, ok := block["timestamp"].(string)
	if !ok {
		return NewMonitorError(ErrTimestampParse, "Invalid timestamp format", nil)
	}

	timestamp, err := m.parseHexToInt64(strings.TrimPrefix(timestampHex, "0x"))
	if err != nil {
		return NewMonitorError(ErrTimestampParse, "Failed to parse block timestamp", err)
	}

	logger.Info("Processing %d transactions from block %d", len(transactions), blockNumber)

	for i, tx := range transactions {
		txMap, ok := tx.(map[string]interface{})
		if !ok {
			logger.Error("Invalid transaction format at index %d", i)
			continue
		}

		processedTx, err := m.parser.ProcessTransaction(txMap, timestamp)
		if err != nil {
			logger.Error("Failed to process transaction: %v", err)
			continue
		}

		if processedTx != nil {
			logger.Debug("Successfully processed transaction %s", processedTx.Hash)
			m.notifyTransaction(*processedTx)
		}
	}

	m.parser.UpdateCurrentBlock(blockNumber)
	logger.Info("Successfully processed block %d", blockNumber)
	return nil
}

// notifyTransaction sends notifications for relevant transactions
func (m *BlockMonitor) notifyTransaction(tx storage.Transaction) {
	isIncoming := m.parser.IsSubscribed(tx.ToAddress)
	direction := "Outgoing"
	address := tx.FromAddress

	if isIncoming {
		direction = "Incoming"
		address = tx.ToAddress
	}

	logger.Info("%s transaction detected for %s", direction, address)
	logger.Debug("Transaction details: Hash: %s, From: %s, To: %s, Value: %f ETH, Block: %d, Time: %s",
		tx.Hash,
		tx.FromAddress,
		tx.ToAddress,
		tx.Value,
		tx.BlockNumber,
		time.Unix(tx.Timestamp, 0).Format("2006-01-02 15:04:05"))

	m.notifier.Notify(notification.Notification{
		Type:        notification.TransactionReceived,
		Address:     address,
		Transaction: tx,
		Timestamp:   tx.Timestamp,
	})
}

// parseHexToInt64 converts a hex string to int64
func (m *BlockMonitor) parseHexToInt64(hex string) (int64, error) {
	value, err := utils.String2Int64(hex, 16)
	if err != nil {
		logger.Error("Failed to parse hex value %s: %v", hex, err)
		return 0, err
	}
	return value, nil
}


