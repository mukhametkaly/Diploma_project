package inventory

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/document-api/src/waybill"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
)

func makeCreateProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.Product)
		resp, err := s.CreateProduct(req)
		if err != nil {
			waybill.Loger.Println(err)
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
			waybill.Loger.Println(err)
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
			waybill.Loger.Println(err)
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
			waybill.Loger.Println(err)
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
			waybill.Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeFilterProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FilterProductsRequest)
		resp, err := s.FilterProducts(req)
		if err != nil {
			waybill.Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

type FilterProductsRequest struct {
}

type DeleteWaybillRequest struct {
	WaybillId  int64  `json:"waybill_id"`
	MerchantId string `json:"merchant_id"`
}

type DeleteInventoryRequest struct {
	InventoryId int64  `json:"inventory_id"`
	MerchantId  string `json:"merchant_id"`
}
