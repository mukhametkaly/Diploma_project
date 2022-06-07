package inventory

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
)

func makeCreateInventoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShortInventory)
		resp, err := s.CreateInventory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateInventoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShortInventory)
		resp, err := s.UpdateInventory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeGetInventoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetInventoryRequest)
		resp, err := s.GetInventory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteInventoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteInventoryRequest)
		err := s.DeleteInventory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeFilterInventoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InventorysFilterRequest)
		resp, err := s.InventorysFilter(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.InventoryProduct)
		resp, err := s.InventoryAddProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateInventoryProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.InventoryProduct)
		resp, err := s.InventoryUpdateProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteInventoryProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteInventoryProductRequest)
		err := s.DeleteInventoryProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetInventoryProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetInventoryProductsRequest)
		resp, err := s.GetInventoryProducts(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

type DeleteInventoryRequest struct {
	InventoryId int64  `json:"inventory_id"`
	MerchantId  string `json:"merchant_id"`
}

type DeleteInventoryProductRequest struct {
	InventoryId int64  `json:"inventory_id"`
	Barcode     string `json:"barcode"`
}

type GetInventoryProductsRequest struct {
	InventoryId int64  `json:"inventory_id"`
	Barcode     string `json:"barcode"`
}

type InventorysFilterRequest struct {
	MerchantId     string `json:"merchant_id"`
	Status         string `json:"status"`
	DocumentNumber string `json:"document_number"`
	From           int    `json:"from"`
	Size           int    `json:"size"`
}

type GetInventoryRequest struct {
	MerchantId  string `json:"merchant_id"`
	InventoryId int64  `json:"inventory_id"`
}
