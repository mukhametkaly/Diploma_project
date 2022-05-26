package catalog

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
)

func makeCreateProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeUpdateProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeDeleteProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeDeleteBatchProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeGetProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeFilterProductEndpoint(s ProductService) endpoint.Endpoint {
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

func makeCreateCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.Category)
		resp, err := s.CreateCategory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.Category)
		resp, err := s.UpdateCategory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int64)
		err := s.DeleteByIdCategory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeDeleteBatchCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.([]int64)
		err := s.DeleteBatchCategory(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(int64)
		resp, err := s.GetCategoryById(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeFilterCategoryEndpoint(s CategoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		merchantId := request.(string)
		resp, err := s.FilterCategories(merchantId)
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
