package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err.(*argError).Status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(err)
}

type argError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newError(statusCode int, err error) *argError {
	return &argError{
		Status:  statusCode,
		Message: err.Error(),
	}
}

func newErrorString(statusCode int, errorString string) *argError {
	return &argError{
		Status:  statusCode,
		Message: errorString,
	}
}

func (e *argError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
