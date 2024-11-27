package utils

import (
    "testing"
)

func TestString2Int64(t *testing.T) {
    testCases := []struct {
        name        string
        input       string
        base        int
        expected    int64
        expectError bool
    }{
        {
            name:        "valid decimal",
            input:      "123",
            base:       10,
            expected:   123,
            expectError: false,
        },
        {
            name:        "valid hex",
            input:      "7b",
            base:       16,
            expected:   123,
            expectError: false,
        },
        {
            name:        "valid binary",
            input:      "1111011",
            base:       2,
            expected:   123,
            expectError: false,
        },
        {
            name:        "negative decimal",
            input:      "-123",
            base:       10,
            expected:   -123,
            expectError: false,
        },
        {
            name:        "invalid input",
            input:      "xyz",
            base:       10,
            expectError: true,
        },
        {
            name:        "empty string",
            input:      "",
            base:       10,
            expectError: true,
        },
        {
            name:        "invalid base",
            input:      "123",
            base:       37,
            expectError: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := String2Int64(tc.input, tc.base)
            
            if tc.expectError {
                if err == nil {
                    t.Error("Expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error: %v", err)
                }
                if result != tc.expected {
                    t.Errorf("Expected %d but got %d", tc.expected, result)
                }
            }
        })
    }
}