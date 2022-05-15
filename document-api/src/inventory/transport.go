package inventory

import (
	"context"
	"encoding/json"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
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

	createInventory := kithttp.NewServer(
		makeCreateInventoryEndpoint(ss),
		decodeCreateInventoryRequest,
		encodeResponse,
		opts...,
	)

	getInventory := kithttp.NewServer(
		makeGetInventoryEndpoint(ss),
		decodeGetInventoryRequest,
		encodeResponse,
		opts...,
	)

	updateInventory := kithttp.NewServer(
		makeUpdateInventoryEndpoint(ss),
		decodeUpdateInventoryRequest,
		encodeResponse,
		opts...,
	)

	deleteInventory := kithttp.NewServer(
		makeDeleteInventoryEndpoint(ss),
		decodeDeleteInventoryRequest,
		encodeResponse,
		opts...,
	)

	filterInventory := kithttp.NewServer(
		makeFilterInventoryEndpoint(ss),
		decodeFilterInventoryRequest,
		encodeResponse,
		opts...,
	)

	addInventoryProduct := kithttp.NewServer(
		makeAddProductEndpoint(ss),
		decodeInventoryAddProductRequest,
		encodeResponse,
		opts...,
	)

	updateInventoryProduct := kithttp.NewServer(
		makeUpdateInventoryProductEndpoint(ss),
		decodeInventoryUpdateProductRequest,
		encodeResponse,
		opts...,
	)

	deleteInventoryProduct := kithttp.NewServer(
		makeDeleteInventoryProductEndpoint(ss),
		decodeDeleteInventoryProductRequest,
		encodeResponse,
		opts...,
	)

	getInventoryProduct := kithttp.NewServer(
		makeGetInventoryProductEndpoint(ss),
		decodeGetInventoryProductsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/inventory/create", createInventory).Methods("POST")
	r.Handle("/v1/inventory/{id}", getInventory).Methods("GET")
	r.Handle("/v1/inventory/{id}", updateInventory).Methods("PUT")
	r.Handle("/v1/inventory/{merchantId}/{id}", deleteInventory).Methods("DELETE")
	r.Handle("/v1/inventory/filter", filterInventory).Methods("POST")

	r.Handle("/v1/inventory/add/product", addInventoryProduct).Methods("POST")
	r.Handle("/v1/inventory/product/{inventoryId}/{barcode}", updateInventoryProduct).Methods("PUT")
	r.Handle("/v1/inventory/product/{inventoryId}/{barcode}", deleteInventoryProduct).Methods("DELETE")
	r.Handle("/v1/inventory/product/get", getInventoryProduct).Methods("POST")

	return r
}

func decodeCreateInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.ShortInventory

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeUpdateInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.ShortInventory
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.ID = id

	return body, nil
}

func decodeDeleteInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}

	merchantId, ok := vars["merchantId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	req := DeleteInventoryRequest{
		InventoryId: id,
		MerchantId:  merchantId,
	}

	return req, nil
}

func decodeFilterInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body InventorysFilterRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeGetInventoryRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body GetInventoryRequest = GetInventoryRequest{
		MerchantId:  "",
		InventoryId: id,
	}

	return body, nil
}

func decodeInventoryAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.InventoryProduct

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	return body, nil
}

func decodeInventoryUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["barcode"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}

	inventoryIdStr, ok := vars["inventoryId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}
	inventoryId, err := strconv.ParseInt(inventoryIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.InventoryProduct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.Barcode = id
	body.InventoryId = inventoryId

	return body, nil
}

func decodeDeleteInventoryProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	inventoryIdStr, ok := vars["inventoryId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no inventory id")
	}
	inventoryId, err := strconv.ParseInt(inventoryIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	barcode, ok := vars["barcode"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	req := DeleteInventoryProductRequest{
		InventoryId: inventoryId,
		Barcode:     barcode,
	}

	return req, nil
}

func decodeGetInventoryProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body GetInventoryProductsRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
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
