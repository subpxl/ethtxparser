package main

import (
    "os"
    "testing"
)

func TestGetEnvOrDefault(t *testing.T) {
    testCases := []struct {
        name         string
        key          string
        defaultValue string
        envValue     string
        expected     string
    }{
        {
            name:         "returns default when env not set",
            key:          "TEST_KEY_1",
            defaultValue: "default",
            envValue:     "",
            expected:     "default",
        },
        {
            name:         "returns env value when set",
            key:          "TEST_KEY_2",
            defaultValue: "default",
            envValue:     "custom",
            expected:     "custom",
        },
        {
            name:         "handles empty default value",
            key:          "TEST_KEY_3",
            defaultValue: "",
            envValue:     "",
            expected:     "",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.envValue != "" {
                os.Setenv(tc.key, tc.envValue)
                defer os.Unsetenv(tc.key)
            }

            result := getEnvOrDefault(tc.key, tc.defaultValue)
            if result != tc.expected {
                t.Errorf("Expected %s, got %s", tc.expected, result)
            }
        })
    }
}

func TestGetEnvIntOrDefault(t *testing.T) {
    testCases := []struct {
        name         string
        key          string
        defaultValue int
        envValue     string
        expected     int
    }{
        {
            name:         "returns default when env not set",
            key:          "TEST_INT_1",
            defaultValue: 5,
            envValue:     "",
            expected:     5,
        },
        {
            name:         "returns env value when set",
            key:          "TEST_INT_2",
            defaultValue: 5,
            envValue:     "10",
            expected:     10,
        },
        {
            name:         "returns default for invalid int",
            key:          "TEST_INT_3",
            defaultValue: 5,
            envValue:     "invalid",
            expected:     5,
        },
        {
            name:         "handles zero default value",
            key:          "TEST_INT_4",
            defaultValue: 0,
            envValue:     "",
            expected:     0,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.envValue != "" {
                os.Setenv(tc.key, tc.envValue)
                defer os.Unsetenv(tc.key)
            }

            result := getEnvIntOrDefault(tc.key, tc.defaultValue)
            if result != tc.expected {
                t.Errorf("Expected %d, got %d", tc.expected, result)
            }
        })
    }
}

func TestDefaultConstants(t *testing.T) {
    // Test that default constants are set to expected values
    testCases := []struct {
        name     string
        actual   string
        expected string
    }{
        {
            name:     "RPC Endpoint",
            actual:   defaultRPCEndpoint,
            expected: "http://127.0.0.1:7545",
        },
        {
            name:     "Server Host",
            actual:   defaultServerHost,
            expected: "127.0.0.1",
        },
        {
            name:     "Server Port",
            actual:   defaultServerPort,
            expected: "8000",
        },
        {
            name:     "Log Path",
            actual:   defaultLogPath,
            expected: "./logs/blockchain-parser.log",
        },
        {
            name:     "Environment",
            actual:   defaultEnv,
            expected: "development",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if tc.actual != tc.expected {
                t.Errorf("Expected %s to be %s, got %s", tc.name, tc.expected, tc.actual)
            }
        })
    }

    // Test default delay separately since it's an int
    if defaultDelay != 5 {
        t.Errorf("Expected default delay to be 5, got %d", defaultDelay)
    }
}

func TestEnvironmentOverrides(t *testing.T) {
    // Save original env values
    originalEndpoint := os.Getenv("RPC_ENDPOINT")
    originalPort := os.Getenv("SERVER_PORT")
    originalDelay := os.Getenv("MONITOR_DELAY")

    // Cleanup function to restore original env values
    defer func() {
        if originalEndpoint != "" {
            os.Setenv("RPC_ENDPOINT", originalEndpoint)
        } else {
            os.Unsetenv("RPC_ENDPOINT")
        }
        if originalPort != "" {
            os.Setenv("SERVER_PORT", originalPort)
        } else {
            os.Unsetenv("SERVER_PORT")
        }
        if originalDelay != "" {
            os.Setenv("MONITOR_DELAY", originalDelay)
        } else {
            os.Unsetenv("MONITOR_DELAY")
        }
    }()

    // Set test environment variables
    os.Setenv("RPC_ENDPOINT", "http://localhost:8545")
    os.Setenv("SERVER_PORT", "9000")
    os.Setenv("MONITOR_DELAY", "10")

    // Test that environment variables override defaults
    if endpoint := getEnvOrDefault("RPC_ENDPOINT", defaultRPCEndpoint); endpoint != "http://localhost:8545" {
        t.Errorf("Expected custom RPC endpoint, got %s", endpoint)
    }

    if port := getEnvOrDefault("SERVER_PORT", defaultServerPort); port != "9000" {
        t.Errorf("Expected custom port, got %s", port)
    }

    if delay := getEnvIntOrDefault("MONITOR_DELAY", defaultDelay); delay != 10 {
        t.Errorf("Expected custom delay, got %d", delay)
    }
}