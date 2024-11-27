package api

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Status  int    `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error codes
const (
	ErrCodeMissingAddress    = "MISSING_ADDRESS"
	ErrCodeInvalidAddress    = "INVALID_ADDRESS"
	ErrCodeAlreadySubscribed = "ALREADY_SUBSCRIBED"
	ErrCodeServerError       = "SERVER_ERROR"
	ErrCodeInvalidMethod     = "INVALID_METHOD"
	ErrCodeJSONParseError    = "JSON_PARSE_ERROR"
)

// Error responses
var (
	ErrMissingAddress = &APIError{
		Status:  http.StatusBadRequest,
		Message: "Missing 'address' parameter",
		Code:    ErrCodeMissingAddress,
	}

	ErrInvalidAddress = &APIError{
		Status:  http.StatusBadRequest,
		Message: "Invalid address format",
		Code:    ErrCodeInvalidAddress,
	}

	ErrAlreadySubscribed = &APIError{
		Status:  http.StatusBadRequest,
		Message: "Address already subscribed",
		Code:    ErrCodeAlreadySubscribed,
	}

	ErrMethodNotAllowed = &APIError{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method not allowed",
		Code:    ErrCodeInvalidMethod,
	}

	ErrInternalServer = &APIError{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
		Code:    ErrCodeServerError,
	}
)

// SendError sends an error response
func SendError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// ValidateMethod checks if the request method is allowed
func ValidateMethod(w http.ResponseWriter, r *http.Request, allowed string) bool {
	if r.Method != allowed {
		SendError(w, ErrMethodNotAllowed)
		return false
	}
	return true
}
