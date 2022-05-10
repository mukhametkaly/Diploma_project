package waybill

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

	createWaybill := kithttp.NewServer(
		makeCreateWaybillEndpoint(ss),
		decodeCreateWaybillRequest,
		encodeResponse,
		opts...,
	)

	updateWaybill := kithttp.NewServer(
		makeUpdateWaybillEndpoint(ss),
		decodeUpdateWaybillRequest,
		encodeResponse,
		opts...,
	)

	deleteWaybill := kithttp.NewServer(
		makeDeleteWaybillEndpoint(ss),
		decodeDeleteWaybillRequest,
		encodeResponse,
		opts...,
	)

	filterWaybill := kithttp.NewServer(
		makeGetWaybillEndpoint(ss),
		decodeGetWaybillRequest,
		encodeResponse,
		opts...,
	)

	addWaybillProduct := kithttp.NewServer(
		makeAddProductEndpoint(ss),
		decodeWaybillAddProductRequest,
		encodeResponse,
		opts...,
	)

	updateWaybillProduct := kithttp.NewServer(
		makeUpdateWaybillProductEndpoint(ss),
		decodeWaybillUpdateProductRequest,
		encodeResponse,
		opts...,
	)

	deleteWaybillProduct := kithttp.NewServer(
		makeDeleteWaybillProductEndpoint(ss),
		decodeDeleteWaybillProductRequest,
		encodeResponse,
		opts...,
	)

	getWaybillProduct := kithttp.NewServer(
		makeGetWaybillProductEndpoint(ss),
		decodeGetWaybillProductsRequest,
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

func decodeCreateWaybillRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.ShortWaybill

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeUpdateWaybillRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no waybill id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.ShortWaybill
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.ID = id

	return body, nil
}

func decodeDeleteWaybillRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no waybill id")
	}

	merchantId, ok := vars["merchantId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	req := DeleteWaybillRequest{
		WaybillId:  id,
		MerchantId: merchantId,
	}

	return req, nil
}

func decodeGetWaybillRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body WaybillsFilterRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeWaybillAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.WaybillProduct

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	return body, nil
}

func decodeWaybillUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no waybill id")
	}

	waybillIdStr, ok := vars["waybillId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no waybill id")
	}
	waybillId, err := strconv.ParseInt(waybillIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.WaybillProduct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.Barcode = id
	body.WaybillId = waybillId

	return body, nil
}

func decodeDeleteWaybillProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	waybillIdStr, ok := vars["waybillId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no waybill id")
	}
	waybillId, err := strconv.ParseInt(waybillIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	barcode, ok := vars["barcode"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	req := DeleteWaybillProductRequest{
		WaybillId: waybillId,
		Barcode:   barcode,
	}

	return req, nil
}

func decodeGetWaybillProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body GetWaybillProductsRequest

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
