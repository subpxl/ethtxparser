// Package main implements the entry point for the blockchain parser application.
// It initializes and coordinates all components including configuration,
// logging, storage, parsing, monitoring, and the HTTP API server.

package main

import (
	"blockchain-parser/config"
	"blockchain-parser/internal/api"
	"blockchain-parser/internal/logger"
	"blockchain-parser/internal/monitor"
	"blockchain-parser/internal/notification"
	"blockchain-parser/internal/parser"
	"blockchain-parser/internal/storage"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Environment variable constants with their default values

const ethUrl string = "https://ethereum-sepolia-rpc.publicnode.com"

const (
	defaultRPCEndpoint = "http://127.0.0.1:7545"
	defaultServerHost  = "127.0.0.1"
	defaultServerPort  = "8000"
	defaultLogPath     = "./logs/blockchain-parser.log"
	defaultEnv         = "development"
	defaultDelay       = 5
)

// getEnvOrDefault retrieves an environment variable value or returns
// the default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault retrieves an environment variable as integer or returns
// the default if not set or invalid
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {

	// Application startup sequence:
	// 1. Load configuration from environment variables
	// 2. Initialize logger
	// 3. Set up storage, RPC client, and parser
	// 4. Create notification service
	// 5. Start block monitor
	// 6. Subscribe to test addresses
	// 7. Start HTTP API server

	// Get configuration from environment variables or use defaults
	rpcEndpoint := getEnvOrDefault("RPC_ENDPOINT", defaultRPCEndpoint)
	serverHost := getEnvOrDefault("SERVER_HOST", defaultServerHost)
	serverPort := getEnvOrDefault("SERVER_PORT", defaultServerPort)
	logFilePath := getEnvOrDefault("LOG_FILE_PATH", defaultLogPath)
	environment := getEnvOrDefault("ENVIRONMENT", defaultEnv)
	monitorDelay := getEnvIntOrDefault("MONITOR_DELAY", defaultDelay)

	cfg := config.NewConfig(
		rpcEndpoint,
		serverHost,
		serverPort,
		logFilePath,
		environment,
		monitorDelay,
	)
	if cfg.Database.Type != config.MemoryDB {
		cfg.WithDatabase(
			config.MongoDB,
			"localhost",
			"27017",
			"blockchain",
			"user",
			"password",
		)
	}

	if err := logger.Init(logFilePath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Initialize components
	storage := storage.NewMemoryStorage()
	rpcClient := parser.NewRPCClient(cfg)
	p := parser.NewParser(storage, rpcClient)

	// Always use console notification service for simplicity
	notificationService := notification.NewConsoleNotificationService()

	monitor := monitor.NewBlockMonitor(p, rpcClient, notificationService)

	// Example addresses for testing
	testAddresses := []string{
		"0xdD93e92dc32d0B2F51430b0e6dA29BDd01AF68D6",
		"0xC22c7f8bA7dE381A299ee4EB3a11E1316525ce45",
		"0x91A9CeF0099DF1D2eA4F4B825ac06B506dfDbe07",
	}

	for _, address := range testAddresses {
		if p.Subscribe(address) {
			fmt.Printf("Subscribed to address: %s\n", address)
		}
	}

	// Start block monitoring in the existing goroutine
	go monitor.StartMonitoring()

	// Start HTTP server (this will block)
	fmt.Printf("Starting HTTP server on %s\n", cfg.GetServerAddress())
	api.StartServer(p, cfg.GetServerAddress())
}
