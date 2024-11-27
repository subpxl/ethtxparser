package api

import (
	"regexp"
	"strings"
)

const (
	addressLength = 42 // standard Ethereum address length including '0x'
)

var (
	// Regex for validating Ethereum addresses
	addressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateAddress checks if the given address is a valid Ethereum address
func ValidateAddress(address string) *ValidationError {
	if address == "" {
		return &ValidationError{
			Field:   "address",
			Message: "Address is required",
		}
	}

	// Convert to lowercase for consistent validation
	address = strings.ToLower(address)

	if len(address) != addressLength {
		return &ValidationError{
			Field:   "address",
			Message: "Address must be 42 characters long",
		}
	}

	if !strings.HasPrefix(address, "0x") {
		return &ValidationError{
			Field:   "address",
			Message: "Address must start with '0x'",
		}
	}

	if !addressRegex.MatchString(address) {
		return &ValidationError{
			Field:   "address",
			Message: "Invalid address format",
		}
	}

	return nil
}



// ValidateBlockNumber validates a block number
func ValidateBlockNumber(blockNum int64) *ValidationError {
	if blockNum < 0 {
		return &ValidationError{
			Field:   "block_number",
			Message: "Block number cannot be negative",
		}
	}
	return nil
}

// ValidateTransactionHash validates a transaction hash
func ValidateTransactionHash(hash string) *ValidationError {
	if hash == "" {
		return &ValidationError{
			Field:   "hash",
			Message: "Transaction hash is required",
		}
	}

	hash = strings.ToLower(hash)
	if !strings.HasPrefix(hash, "0x") {
		return &ValidationError{
			Field:   "hash",
			Message: "Transaction hash must start with '0x'",
		}
	}

	// Standard transaction hash length is 66 characters (0x + 64 hex characters)
	if len(hash) != 66 {
		return &ValidationError{
			Field:   "hash",
			Message: "Transaction hash must be 66 characters long",
		}
	}

	hashRegex := regexp.MustCompile("^0x[0-9a-f]{64}$")
	if !hashRegex.MatchString(hash) {
		return &ValidationError{
			Field:   "hash",
			Message: "Invalid transaction hash format",
		}
	}

	return nil
}
