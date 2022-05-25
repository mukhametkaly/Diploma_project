package shopping_cart

import (
	"context"
	"encoding/json"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
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

	createShoppingCart := kithttp.NewServer(
		makeCreateShoppingCartEndpoint(ss),
		decodeCreateShoppingCartRequest,
		encodeResponse,
		opts...,
	)

	getShoppingCart := kithttp.NewServer(
		makeGetShoppingCartEndpoint(ss),
		decodeGetShoppingCartRequest,
		encodeResponse,
		opts...,
	)

	updateShoppingCart := kithttp.NewServer(
		makeUpdateShoppingCartEndpoint(ss),
		decodeUpdateShoppingCartRequest,
		encodeResponse,
		opts...,
	)

	deleteShoppingCart := kithttp.NewServer(
		makeDeleteShoppingCartEndpoint(ss),
		decodeDeleteShoppingCartRequest,
		encodeResponse,
		opts...,
	)

	filterShoppingCart := kithttp.NewServer(
		makeFilterShoppingCartEndpoint(ss),
		decodeFilterShoppingCartRequest,
		encodeResponse,
		opts...,
	)

	addShoppingCartProduct := kithttp.NewServer(
		makeAddProductEndpoint(ss),
		decodeShoppingCartAddProductRequest,
		encodeResponse,
		opts...,
	)

	updateShoppingCartProduct := kithttp.NewServer(
		makeUpdateShoppingCartProductEndpoint(ss),
		decodeShoppingCartUpdateProductRequest,
		encodeResponse,
		opts...,
	)

	deleteShoppingCartProduct := kithttp.NewServer(
		makeDeleteShoppingCartProductEndpoint(ss),
		decodeDeleteShoppingCartProductRequest,
		encodeResponse,
		opts...,
	)

	getShoppingCartProduct := kithttp.NewServer(
		makeGetShoppingCartProductEndpoint(ss),
		decodeGetShoppingCartProductsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/v1/shopping_cart/create", createShoppingCart).Methods("POST")
	r.Handle("/v1/shopping_cart/{id}", getShoppingCart).Methods("GET")
	r.Handle("/v1/shopping_cart/{id}", updateShoppingCart).Methods("PUT")
	r.Handle("/v1/shopping_cart/{merchantId}/{id}", deleteShoppingCart).Methods("DELETE")
	r.Handle("/v1/shopping_cart/filter", filterShoppingCart).Methods("POST")

	r.Handle("/v1/shopping_cart/add/merchant", addShoppingCartProduct).Methods("POST")
	r.Handle("/v1/shopping_cart/merchant/{shoppingCartId}/{barcode}", updateShoppingCartProduct).Methods("PUT")
	r.Handle("/v1/shopping_cart/merchant/{shoppingCartId}/{barcode}", deleteShoppingCartProduct).Methods("DELETE")
	r.Handle("/v1/shopping_cart/merchant/get", getShoppingCartProduct).Methods("POST")

	return r
}

func decodeCreateShoppingCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.ShoppingCart

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeUpdateShoppingCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping cart id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.ID = id

	return body, nil
}

func decodeDeleteShoppingCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping_cart id")
	}

	merchantId, ok := vars["merchantId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	req := DeleteShoppingCartRequest{
		ShoppingCartId: id,
		MerchantId:     merchantId,
	}

	return req, nil
}

func decodeFilterShoppingCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body ShoppingCartsFilterRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}
	return body, nil
}

func decodeGetShoppingCartRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping_cart id")
	}

	id, err := strconv.ParseInt(strId, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body GetShoppingCartRequest = GetShoppingCartRequest{
		MerchantId:     "",
		ShoppingCartId: id,
	}

	return body, nil
}

func decodeShoppingCartAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body models.ShoppingCartProduct

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	return body, nil
}

func decodeShoppingCartUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["barcode"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping cart id")
	}

	shoppingCartIdStr, ok := vars["shoppingCartId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping cart id")
	}
	shoppingCartId, err := strconv.ParseInt(shoppingCartIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	var body models.ShoppingCartProduct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	body.Barcode = id
	body.ShoppingCartId = shoppingCartId

	return body, nil
}

func decodeDeleteShoppingCartProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	shoppingCartIdStr, ok := vars["shoppingCartId"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no shopping cart id")
	}
	shoppingCartId, err := strconv.ParseInt(shoppingCartIdStr, 0, 64)
	if err != nil {
		return nil, newError(http.StatusBadRequest, err)
	}

	barcode, ok := vars["barcode"]
	if !ok {
		return nil, newStringError(http.StatusBadRequest, "no merchant id")
	}

	req := DeleteShoppingCartProductRequest{
		ShoppingCartId: shoppingCartId,
		Barcode:        barcode,
	}

	return req, nil
}

func decodeGetShoppingCartProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body GetShoppingCartProductsRequest

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
