package merchant

import (
	"context"
	"encoding/json"
	_ "github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
	"net/http"
)

func MakeHandler(ss Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createMerchant := kithttp.NewServer(
		makeCreateMerchantEndpoint(ss),
		decodeCreateMerchantRequest,
		encodeResponse,
		opts...,
	)

	updateMerchant := kithttp.NewServer(
		makeUpdateMerchantEndpoint(ss),
		decodeUpdateMerchantRequest,
		encodeResponse,
		opts...,
	)

	deleteMerchant := kithttp.NewServer(
		makeDeleteMerchantEndpoint(ss),
		decodeDeleteMerchantRequest,
		encodeResponse,
		opts...,
	)

	deleteMMerchant := kithttp.NewServer(
		makeDeleteBatchMerchantEndpoint(ss),
		decodeDeleteBatchMerchantRequest,
		encodeResponse,
		opts...,
	)

	getMerchant := kithttp.NewServer(
		makeGetMerchantEndpoint(ss),
		decodeGetMerchantRequest,
		encodeResponse,
		opts...,
	)

	filterMerchants := kithttp.NewServer(
		makeFilterMerchantEndpoint(ss),
		decodeFilterMerchantSync,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/merchant/create", createMerchant).Methods("POST")
	r.Handle("/v1/merchant/{id}", getMerchant).Methods("GET")
	r.Handle("/v1/merchant/{id}", updateMerchant).Methods("PUT")
	r.Handle("/v1/merchant/{id}", deleteMerchant).Methods("DELETE")
	r.Handle("/v1/merchant/delete/batch", deleteMMerchant).Methods("POST")
	r.Handle("/v1/merchant/filter", filterMerchants).Methods("POST")

	return r
}

func decodeCreateMerchantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.Merchant

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}
	return body, nil
}

func decodeUpdateMerchantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	var body models.Merchant
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}

	body.MerchantId = strId

	return body, nil
}

func decodeDeleteMerchantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	merchantId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}
	return merchantId, nil
}

func decodeGetMerchantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	return strId, nil
}

func decodeDeleteBatchMerchantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body []string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}

	return body, nil
}

func decodeFilterMerchantSync(_ context.Context, r *http.Request) (interface{}, error) {
	var body FilterMerchantsRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}
	return body, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}
