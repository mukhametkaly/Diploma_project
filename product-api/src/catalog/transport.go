package catalog

import (
	"context"
	"encoding/json"
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
		makeCreateProductEndpoint(ss),
		decodeCreateProductRequest,
		encodeResponse,
		opts...,
	)

	updateProduct := kithttp.NewServer(
		makeUpdateProductEndpoint(ss),
		decodeUpdateProductRequest,
		encodeResponse,
		opts...,
	)

	deleteProduct := kithttp.NewServer(
		makeDeleteProductEndpoint(ss),
		decodeDeleteProductRequest,
		encodeResponse,
		opts...,
	)

	deleteMProduct := kithttp.NewServer(
		makeDeleteBatchProductEndpoint(ss),
		decodeDeleteBatchProductRequest,
		encodeResponse,
		opts...,
	)

	getProduct := kithttp.NewServer(
		makeGetProductEndpoint(ss),
		decodeGetProductRequest,
		encodeResponse,
		opts...,
	)

	filterProducts := kithttp.NewServer(
		makeFilterProductEndpoint(ss),
		decodeFilterProductSync,
		encodeResponse,
		opts...,
	)

	createCategory := kithttp.NewServer(
		makeCreateCategoryEndpoint(ss),
		decodeCreateCategoryRequest,
		encodeResponse,
		opts...,
	)

	updateCategory := kithttp.NewServer(
		makeUpdateCategoryEndpoint(ss),
		decodeUpdateCategoryRequest,
		encodeResponse,
		opts...,
	)

	deleteCategory := kithttp.NewServer(
		makeDeleteCategoryEndpoint(ss),
		decodeDeleteCategoryRequest,
		encodeResponse,
		opts...,
	)

	deleteMCategory := kithttp.NewServer(
		makeDeleteBatchCategoryEndpoint(ss),
		decodeDeleteBatchCategoryRequest,
		encodeResponse,
		opts...,
	)

	getCategory := kithttp.NewServer(
		makeGetCategoryEndpoint(ss),
		decodeGetCategoryRequest,
		encodeResponse,
		opts...,
	)

	filterCategories := kithttp.NewServer(
		makeFilterCategoryEndpoint(ss),
		decodeFilterCategorySync,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/product/create", createProduct).Methods("POST")
	r.Handle("/v1/product/{id}", getProduct).Methods("GET")
	r.Handle("/v1/product/{id}", updateProduct).Methods("PUT")
	r.Handle("/v1/product/{id}", deleteProduct).Methods("DELETE")
	r.Handle("/v1/product/delete/batch", deleteMProduct).Methods("POST")
	r.Handle("/v1/product/filter", filterProducts).Methods("POST")

	r.Handle("/v1/product/category/create", createCategory).Methods("POST")
	r.Handle("/v1/product/category/{id}", getCategory).Methods("GET")
	r.Handle("/v1/product/category/{id}", updateCategory).Methods("PUT")
	r.Handle("/v1/product/category/{id}", deleteCategory).Methods("DELETE")
	r.Handle("/v1/product/category/delete/batch", deleteMCategory).Methods("POST")
	r.Handle("/v1/product/category/filter", filterCategories).Methods("POST")

	return r
}

func decodeCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.Product

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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
		return nil, InvalidCharacter
	}

	return body, nil
}

func decodeFilterProductSync(_ context.Context, r *http.Request) (interface{}, error) {
	var body FilterProductsRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}
	return body, nil
}

func decodeCreateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.Category

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}
	return body, nil
}

func decodeUpdateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, InvalidCharacter
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, InvalidCharacter
	}

	var body models.Category
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}

	body.ID = id

	return body, nil
}

func decodeDeleteCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeGetCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeDeleteBatchCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body []int64
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, InvalidCharacter
	}

	return body, nil
}

func decodeFilterCategorySync(_ context.Context, r *http.Request) (interface{}, error) {
	var body FilterCategoryRequest
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
