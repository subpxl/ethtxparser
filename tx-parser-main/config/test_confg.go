// config/config_test.go
package config

import (
    "testing"
)

func TestNewConfig(t *testing.T) {
    tests := []struct {
        name         string
        rpcEndpoint  string
        serverHost   string
        serverPort   string
        logFilePath  string
        environment  string
        monitorDelay int
        wantAddress  string
    }{
        {
            name:         "basic config",
            rpcEndpoint:  "http://localhost:8545",
            serverHost:   "localhost",
            serverPort:   "8000",
            logFilePath:  "./logs/app.log",
            environment:  "development",
            monitorDelay: 5,
            wantAddress:  "localhost:8000",
        },
        {
            name:         "different host port",
            rpcEndpoint:  "http://127.0.0.1:8545",
            serverHost:   "127.0.0.1",
            serverPort:   "9000",
            logFilePath:  "./logs/prod.log",
            environment:  "production",
            monitorDelay: 10,
            wantAddress:  "127.0.0.1:9000",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cfg := NewConfig(
                tt.rpcEndpoint,
                tt.serverHost,
                tt.serverPort,
                tt.logFilePath,
                tt.environment,
                tt.monitorDelay,
            )

            // Verify all fields are set correctly
            if cfg.RPCEndpoint != tt.rpcEndpoint {
                t.Errorf("RPCEndpoint = %v, want %v", cfg.RPCEndpoint, tt.rpcEndpoint)
            }
            if cfg.ServerHost != tt.serverHost {
                t.Errorf("ServerHost = %v, want %v", cfg.ServerHost, tt.serverHost)
            }
            if cfg.ServerPort != tt.serverPort {
                t.Errorf("ServerPort = %v, want %v", cfg.ServerPort, tt.serverPort)
            }
            if cfg.LogFilePath != tt.logFilePath {
                t.Errorf("LogFilePath = %v, want %v", cfg.LogFilePath, tt.logFilePath)
            }
            if cfg.MonitorDelay != tt.monitorDelay {
                t.Errorf("MonitorDelay = %v, want %v", cfg.MonitorDelay, tt.monitorDelay)
            }

            // Test GetServerAddress
            if addr := cfg.GetServerAddress(); addr != tt.wantAddress {
                t.Errorf("GetServerAddress() = %v, want %v", addr, tt.wantAddress)
            }

            // Verify default database type
            if cfg.Database.Type != MemoryDB {
                t.Errorf("Default Database.Type = %v, want %v", cfg.Database.Type, MemoryDB)
            }
        })
    }
}

func TestWithDatabase(t *testing.T) {
    cfg := NewConfig(
        "http://localhost:8545",
        "localhost",
        "8000",
        "./logs/app.log",
        "development",
        5,
    )

    // Test adding MongoDB configuration
    cfg.WithDatabase(
        MongoDB,
        "localhost",
        "27017",
        "blockchain",
        "user",
        "pass",
    )

    if cfg.Database.Type != MongoDB {
        t.Errorf("Database.Type = %v, want %v", cfg.Database.Type, MongoDB)
    }
    if cfg.Database.Host != "localhost" {
        t.Errorf("Database.Host = %v, want localhost", cfg.Database.Host)
    }
    if cfg.Database.Port != "27017" {
        t.Errorf("Database.Port = %v, want 27017", cfg.Database.Port)
    }
    if cfg.Database.Name != "blockchain" {
        t.Errorf("Database.Name = %v, want blockchain", cfg.Database.Name)
    }
}