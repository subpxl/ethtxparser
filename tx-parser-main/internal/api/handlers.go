// Package api provides HTTP endpoints for interacting with the blockchain parser
package api

import (
	"blockchain-parser/internal/logger"
	"blockchain-parser/internal/parser"
	"encoding/json"
	"net/http"
)

// StartServer initializes and starts the HTTP server with all endpoints
func StartServer(p parser.Parser, address string) error {
	http.HandleFunc("/currentBlock", makeCurrentBlockHandler(p))
	http.HandleFunc("/subscribe", makeSubscribeHandler(p))
	http.HandleFunc("/transactions", makeTransactionsHandler(p))

	// IGONRE: for testing purposes
	http.HandleFunc("/subscribers", makeSubscribersList(p))

	return http.ListenAndServe(address, nil)
}

// makeCurrentBlockHandler creates a handler for /currentBlock endpoint
func makeCurrentBlockHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("current block request from %s", r.RemoteAddr)

		if !ValidateMethod(w, r, http.MethodGet) {
			logger.Warn("Invalid method %s for current block from %s", r.Method, r.RemoteAddr)
			return
		}
		currentBlock := p.GetCurrentBlock()
		respondWithJSON(w, http.StatusOK, map[string]int64{"current_block": currentBlock})
	}
}

// makeSubscribeHandler adds subscriber to list of subscribers
func makeSubscribeHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received subscription request from %s", r.RemoteAddr)

		if !ValidateMethod(w, r, http.MethodPost) {
			logger.Warn("Invalid method %s for subscription", r.Method)

			return
		}

		address := r.URL.Query().Get("address")

		if err := ValidateAddress(address); err != nil {
			logger.Error("Invalid address format: %s", address)

			SendError(w, &APIError{
				Status:  http.StatusBadRequest,
				Message: err.Message,
				Code:    ErrCodeInvalidAddress,
			})
			return
		}

		if p.IsSubscribed(address) {
			logger.Warn("Address already subscribed: %s", address)

			SendError(w, ErrAlreadySubscribed)
			return
		}

		if !p.Subscribe(address) {
			logger.Error("Failed to subscribe address: %s", address)

			SendError(w, ErrInvalidAddress)
			return
		}

		logger.Info("Successfully subscribed address: %s", address)
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "success",
			"address": address,
		})
	}
}

// makeTransactionsHandler creates a handler for /transactions endpoint
func makeTransactionsHandler(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Handling transactions request from %s", r.RemoteAddr)

		if !ValidateMethod(w, r, http.MethodGet) {
			logger.Warn("Invalid HTTP method %s for transactions endpoint", r.Method)
			return
		}

		address := r.URL.Query().Get("address")
		logger.Debug("Querying transactions for address: %s", address)

		if err := ValidateAddress(address); err != nil {
			logger.Error("Invalid address format: %s - %s", address, err.Message)
			SendError(w, &APIError{
				Status:  http.StatusBadRequest,
				Message: err.Message,
				Code:    ErrCodeInvalidAddress,
			})
			return
		}

		transactions := p.GetTransactions(address)
		if transactions == nil {
			logger.Info("No transactions found for address: %s", address)
			respondWithJSON(w, http.StatusOK, map[string]interface{}{
				"status":       "not_found",
				"address":      address,
				"transactions": []interface{}{},
			})
			return
		}

		logger.Debug("Found %d transactions for address %s", len(transactions), address)

		// Validate each transaction
		for i, tx := range transactions {
			if err := ValidateTransactionHash(tx.Hash); err != nil {
				logger.Error("Invalid transaction hash found in storage: %s for address %s", tx.Hash, address)
				SendError(w, &APIError{
					Status:  http.StatusInternalServerError,
					Message: "Invalid transaction data in storage",
					Code:    ErrCodeServerError,
				})
				return
			}
			logger.Debug("Validated transaction %d/%d: %s", i+1, len(transactions), tx.Hash)
		}

		logger.Info("Successfully returning %d transactions for address %s", len(transactions), address)
		respondWithJSON(w, http.StatusOK, transactions)
	}
}

// makeSubscribersList creates a handler for /subscribers endpoint
func makeSubscribersList(p parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Handling subscribers list request from %s", r.RemoteAddr)

		if !ValidateMethod(w, r, http.MethodGet) {
			logger.Warn("Invalid HTTP method %s for subscribers list endpoint", r.Method)
			return
		}

		subscribers := p.GetSubscribers()
		logger.Debug("Retrieved %d subscribers", len(subscribers))

		response := struct {
			Title       string   `json:"title"`
			Subscribers []string `json:"subscribers"`
		}{
			Title:       "Subscribers",
			Subscribers: subscribers,
		}

		if subscribers == nil {
			logger.Debug("No subscribers found, returning empty list")
			response.Subscribers = []string{}
		}

		logger.Info("Successfully returning subscribers list with %d entries", len(response.Subscribers))
		respondWithJSON(w, http.StatusOK, response)
	}
}

// respondWithJSON helper for formatting JSON responses
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		SendError(w, ErrInternalServer)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
