package auth

import (
	"context"
	"encoding/json"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/models"
	"net/http"
	"strings"

	_ "github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(ss Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	registrationUser := kithttp.NewServer(
		makeRegistrateEndpoint(ss),
		decodeRegistrateUserRequest,
		encodeResponse,
		opts...,
	)

	loginUser := kithttp.NewServer(
		makeLoginEndpoint(ss),
		decodeLoginRequest,
		encodeLoginResponse,
		opts...,
	)

	auth := kithttp.NewServer(
		makeAuthEndpoint(ss),
		decodeAuthRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/auth/registration", registrationUser).Methods("POST")
	r.Handle("/v1/auth/login", loginUser).Methods("POST")
	r.Handle("/v1/auth/validation", auth).Methods("POST")

	return r
}

func decodeRegistrateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newErrorString(http.StatusBadRequest, "cannot unmarshal request")
	}
	return body, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newErrorString(http.StatusBadRequest, "cannot unmarshal request")
	}
	return body, nil
}

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	BearerToken := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(BearerToken, "Bearer ")
	if tokenString == "" {
		return nil, newErrorString(http.StatusUnauthorized, "unauthorized")
	}
	return tokenString, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Authorization", "Bearer "+response.(string))
	return nil
}

type errorer interface {
	error() error
}
