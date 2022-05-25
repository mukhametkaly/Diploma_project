package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case Conflict:
		w.WriteHeader(http.StatusConflict)
	case AccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case NoContentFound:
		w.WriteHeader(http.StatusNoContent)
	case DeserializeBug:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case InvalidCharacter:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(err)
}

type argError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var (
	InternalServerError = &argError{500, "something went wrong"}
	PostgresReadError   = &argError{409, "Ошибка считывания"}
	DeserializeBug      = &argError{415, "Ошибка сериализации"}
	NoContentFound      = &argError{204, "Ничего не найдено"}
	InvalidCharacter    = &argError{400, "Неправильные входные данные"}
	AccessDenied        = &argError{403, "Доступ к ресурсу запрещен"}
	Conflict            = &argError{409, "Конфликт обращения к ресурсу"}
	Unauthorized        = &argError{401, "Невалидный авторизационный токен"}
)

func (e *argError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
