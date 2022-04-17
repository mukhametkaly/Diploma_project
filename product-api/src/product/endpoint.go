package product

import (
	"context"
	"github.com/go-kit/kit/endpoint"
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
		req := request.(FilterProductsRequest)
		resp, err := s.FilterProducts(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

type FilterProductsRequest struct {
	MerchantId string `json:"merchant_id"`
	Barcode    string `json:"barcode"`
	Name       string `json:"name"`
	From       int    `json:"from"`
	Size       int    `json:"size"`
}
