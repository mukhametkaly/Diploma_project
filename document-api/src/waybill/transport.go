package waybill

import (
	"context"
	"encoding/json"
	"github.com/mukhametkaly/Diploma/document-api/src/inventory"
	"github.com/mukhametkaly/Diploma/document-api/src/waybill"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
	"net/http"
	"strconv"

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

	createProduct := kithttp.NewServer(
		waybill.makeCreateProductEndpoint(ss),
		decodeCreateProductRequest,
		encodeResponse,
		opts...,
	)

	updateProduct := kithttp.NewServer(
		waybill.makeUpdateProductEndpoint(ss),
		decodeUpdateProductRequest,
		encodeResponse,
		opts...,
	)

	deleteProduct := kithttp.NewServer(
		waybill.makeDeleteProductEndpoint(ss),
		decodeDeleteProductRequest,
		encodeResponse,
		opts...,
	)

	deleteMProduct := kithttp.NewServer(
		waybill.makeDeleteBatchProductEndpoint(ss),
		decodeDeleteBatchProductRequest,
		encodeResponse,
		opts...,
	)

	getProduct := kithttp.NewServer(
		waybill.makeGetProductEndpoint(ss),
		decodeGetProductRequest,
		encodeResponse,
		opts...,
	)

	filterProducts := kithttp.NewServer(
		waybill.makeFilterProductEndpoint(ss),
		decodeFilterProductSync,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/waybill/create", createProduct).Methods("POST")
	r.Handle("/v1/waybill/{id}", getProduct).Methods("GET")
	r.Handle("/v1/waybill/{id}", updateProduct).Methods("PUT")
	r.Handle("/v1/waybill/{id}", deleteProduct).Methods("DELETE")
	r.Handle("/v1/waybill/delete/batch", deleteMProduct).Methods("POST")
	r.Handle("/v1/waybill/filter", filterProducts).Methods("POST")

	return r
}

func decodeCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.Product

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		InvalidCharacter.DeveloperMessage = err.Error()
		return nil, InvalidCharacter
	}
	return body, nil
}

func decodeUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, InvalidCharacter
	}

	var body models.Product
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		InvalidCharacter.DeveloperMessage = err.Error()
		return nil, InvalidCharacter
	}

	body.ID = id

	return body, nil
}

func decodeDeleteProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, InvalidCharacter
	}

	return id, nil
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, InvalidCharacter
	}

	return id, nil
}

func decodeDeleteBatchProductRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var body []int64
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		InvalidCharacter.DeveloperMessage = err.Error()
		return nil, InvalidCharacter
	}

	return body, nil
}

func decodeFilterProductSync(_ context.Context, r *http.Request) (interface{}, error) {
	var body inventory.FilterProductsRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		InvalidCharacter.DeveloperMessage = err.Error()
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
