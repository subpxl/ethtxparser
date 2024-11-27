// Package config provides configuration management for the application.
package config

import (
	"strings"
	"time"
)

type DatabaseType string

const (
	MemoryDB DatabaseType = "memory"
	MongoDB  DatabaseType = "mongodb"
	MySQL    DatabaseType = "mysql"
)

type NetworkType string

const (
	Local  NetworkType = "local"
	Public NetworkType = "public"
)

type NetworkConfig struct {
	Type           NetworkType
	RequestTimeout time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
	RateLimitDelay time.Duration
}

type DatabaseConfig struct {
	Type     DatabaseType
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	RPCEndpoint  string
	ServerHost   string
	ServerPort   string
	LogFilePath  string
	Environment  string
	MonitorDelay int
	Database     DatabaseConfig
	Network      NetworkConfig
}

func NewConfig(RPCEndpoint, ServerHost, ServerPort, LogFilePath, Environment string, MonitorDelay int) *Config {
	networkType := Public
	if isLocalEndpoint(RPCEndpoint) {
		networkType = Local
	}

	return &Config{
		RPCEndpoint:  RPCEndpoint,
		ServerHost:   ServerHost,
		ServerPort:   ServerPort,
		MonitorDelay: MonitorDelay,
		LogFilePath:  LogFilePath,
		Environment:  Environment,
		Database: DatabaseConfig{
			Type: MemoryDB,
		},
		Network: getNetworkConfig(networkType),
	}
}

func isLocalEndpoint(endpoint string) bool {
	return strings.Contains(endpoint, "localhost") || strings.Contains(endpoint, "127.0.0.1")
}

func getNetworkConfig(netType NetworkType) NetworkConfig {
	switch netType {
	case Local:
		return NetworkConfig{
			Type:           Local,
			RequestTimeout: 5 * time.Second,
			RetryAttempts:  1,
			RetryDelay:     time.Second,
			RateLimitDelay: 100 * time.Millisecond,
		}
	default:
		return NetworkConfig{
			Type:           Public,
			RequestTimeout: 10 * time.Second,
			RetryAttempts:  3,
			RetryDelay:     2 * time.Second,
			RateLimitDelay: time.Second,
		}
	}
}

// GetServerAddress returns the complete server address in format "host:port"
func (cfg *Config) GetServerAddress() string {
	return cfg.ServerHost + ":" + cfg.ServerPort
}

// WithDatabase sets the database configuration and returns the updated Config
func (cfg *Config) WithDatabase(dbType DatabaseType, host, port, name, user, password string) *Config {
	cfg.Database = DatabaseConfig{
		Type:     dbType,
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: password,
	}
	return cfg
}
