package waybill

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/document-api/src/inventory"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
)

func makeCreateProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.Product)
		resp, err := s.CreateProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.Product)
		resp, err := s.UpdateProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int64)
		err := s.DeleteByIdProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeDeleteBatchProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.([]int64)
		err := s.DeleteBatchProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int64)
		resp, err := s.GetProductById(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeFilterProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(inventory.FilterProductsRequest)
		resp, err := s.FilterProducts(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

type DeleteWaybillRequest struct {
	WaybillId  int64  `json:"waybill_id"`
	MerchantId string `json:"merchant_id"`
}

type DeleteWaybillProductRequest struct {
	WaybillId  int64  `json:"waybill_id"`
	Barcode    string `json:"barcode"`
	MerchantId string `json:"merchant_id"`
}

type GetWaybillProductsRequest struct {
	WaybillId  int64 `json:"waybill_id"`
	MerchantId int64 `json:"merchant_id"`
}

type WaybillsFilterRequest struct {
	MerchantId int64  `json:"merchant_id"`
	Status     string `json:"status"`
}
