package waybill

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
)

func makeCreateWaybillEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShortWaybill)
		resp, err := s.CreateWaybill(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateWaybillEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShortWaybill)
		resp, err := s.UpdateWaybill(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteWaybillEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteWaybillRequest)
		err := s.DeleteWaybill(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetWaybillEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WaybillsFilterRequest)
		resp, err := s.WaybillsFilter(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.WaybillProduct)
		resp, err := s.WaybillAddProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateWaybillProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.WaybillProduct)
		resp, err := s.WaybillUpdateProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteWaybillProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteWaybillProductRequest)
		err := s.DeleteWaybillProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetWaybillProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetWaybillProductsRequest)
		resp, err := s.GetWaybillProducts(req)
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
	WaybillId int64  `json:"waybill_id"`
	Barcode   string `json:"barcode"`
}

type GetWaybillProductsRequest struct {
	WaybillId int64  `json:"waybill_id"`
	Barcode   string `json:"barcode"`
}

type WaybillsFilterRequest struct {
	MerchantId     string `json:"merchant_id"`
	Status         string `json:"status"`
	DocumentNumber string `json:"document_number"`
}
