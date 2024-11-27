package logger

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testLogFile = "./test_logs/test.log"

func cleanup() {
	Close()
	os.RemoveAll(filepath.Dir(testLogFile))
}

func readLastLogLine() (string, error) {
	file, err := os.Open(testLogFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	return lastLine, scanner.Err()
}

func TestInit(t *testing.T) {
	defer cleanup()

	tests := []struct {
		name        string
		logPath     string
		expectError bool
	}{
		{
			name:        "valid path",
			logPath:     testLogFile,
			expectError: false,
		},
		{
			name:        "invalid path",
			logPath:     "\x00invalid", // Use null character which is invalid in file paths
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.logPath)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
func TestLogging(t *testing.T) {
	defer cleanup()

	if err := Init(testLogFile); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		level    string
		message  string
		exitTest bool
	}{
		{
			name:    "debug message",
			logFunc: Debug,
			level:   "DEBUG",
			message: "test debug message",
		},
		{
			name:    "info message",
			logFunc: Info,
			level:   "INFO",
			message: "test info message",
		},
		{
			name:    "warn message",
			logFunc: Warn,
			level:   "WARN",
			message: "test warn message",
		},
		{
			name:    "error message",
			logFunc: Error,
			level:   "ERROR",
			message: "test error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(tt.message)

			// Read the last line from the log file
			lastLine, err := readLastLogLine()
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}

			// Check if the log contains the expected components
			if !strings.Contains(lastLine, tt.level) {
				t.Errorf("Expected log level %s not found in: %s", tt.level, lastLine)
			}
			if !strings.Contains(lastLine, tt.message) {
				t.Errorf("Expected message '%s' not found in: %s", tt.message, lastLine)
			}
		})
	}
}

func TestGetCallerInfo(t *testing.T) {
	info := getCallerInfo()
	if info == "unknown" {
		t.Error("Failed to get caller information")
	}
	if !strings.Contains(info, ".go:") {
		t.Errorf("Expected file information in caller info, got: %s", info)
	}
}

func TestClose(t *testing.T) {
	if err := Init(testLogFile); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	Close()

	// Try to write to the closed file
	Info("test message after close")

	// Attempting to read from the file should still work
	_, err := os.ReadFile(testLogFile)
	if err != nil {
		t.Errorf("Failed to read closed log file: %v", err)
	}
}

// Add this at the package level
var osExit = os.Exit
