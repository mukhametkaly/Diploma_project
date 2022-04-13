package product

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
	Status           int    `json:"status"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage"`
}

var (
	InternalServerError = &argError{}
	PostgresReadError   = &argError{409, "Ошибка считывания", "Недоступена база"}
	DeserializeBug      = &argError{415, "Ошибка сериализации", "Ошибка сериализации поискового движка"}
	NoContentFound      = &argError{204, "Ничего не найдено", "Ничего не найдено"}
	InvalidCharacter    = &argError{400, "Неправильные входные данные", "Неправильный JSON"}
	AccessDenied        = &argError{403, "Доступ к ресурсу запрещен", "Доступ к ресурсу или отдельной его части запрещен"}
	Conflict            = &argError{409, "Конфликт обращения к ресурсу", "Запрос не может быть выполнен из-за конфликтного обращения к ресурсу"}
	Unauthorized        = &argError{401, "Невалидный авторизационный токен", "Запрос не может быть выполнен из-за конфликтного обращения к ресурсу"}
)

func (e *argError) Error() string {
	return fmt.Sprintf("%s %s", e.DeveloperMessage, e.Message)
}
