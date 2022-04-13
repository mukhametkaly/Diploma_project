package waybill

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	statusCode := err.(*argError).Status
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(err)
}

type argError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newError(statusCode int, error error) error {
	return &argError{Status: statusCode, Message: error.Error()}
}

func newStringError(statusCode int, errorString string) error {
	return &argError{Status: statusCode, Message: error}
}

func (e *argError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
