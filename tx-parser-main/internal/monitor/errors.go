package monitor

import "fmt"

type MonitorError struct {
	Code    string
	Message string
	Err     error
}

func (e *MonitorError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

const (
	ErrBlockNumberFetch   = "BLOCK_NUMBER_FETCH_ERROR"
	ErrBlockNumberParse   = "BLOCK_NUMBER_PARSE_ERROR"
	ErrBlockFetch         = "BLOCK_FETCH_ERROR"
	ErrTimestampParse     = "TIMESTAMP_PARSE_ERROR"
	ErrTransactionProcess = "TRANSACTION_PROCESS_ERROR"
)

func NewMonitorError(code string, message string, err error) *MonitorError {
	return &MonitorError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
