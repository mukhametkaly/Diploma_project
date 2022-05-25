package shopping_cart

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
)

func makeCreateShoppingCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShoppingCart)
		resp, err := s.CreateShoppingCart(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateShoppingCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShoppingCart)
		resp, err := s.UpdateShoppingCart(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeGetShoppingCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetShoppingCartRequest)
		resp, err := s.GetShoppingCart(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteShoppingCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteShoppingCartRequest)
		err := s.DeleteShoppingCart(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeFilterShoppingCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ShoppingCartsFilterRequest)
		resp, err := s.ShoppingCartsFilter(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShoppingCartProduct)
		resp, err := s.ShoppingCartAddProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeUpdateShoppingCartProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ShoppingCartProduct)
		resp, err := s.ShoppingCartUpdateProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

func makeDeleteShoppingCartProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteShoppingCartProductRequest)
		err := s.DeleteShoppingCartProduct(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return nil, nil
	}
}

func makeGetShoppingCartProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetShoppingCartProductsRequest)
		resp, err := s.GetShoppingCartProducts(req)
		if err != nil {
			Loger.Println(err)
			return nil, err
		}
		return resp, nil
	}
}

type DeleteShoppingCartRequest struct {
	ShoppingCartId int64  `json:"shoppingCart_id"`
	MerchantId     string `json:"merchant_id"`
}

type DeleteShoppingCartProductRequest struct {
	ShoppingCartId int64  `json:"shoppingCart_id"`
	Barcode        string `json:"barcode"`
}

type GetShoppingCartProductsRequest struct {
	ShoppingCartId int64  `json:"shoppingCart_id"`
	Barcode        string `json:"barcode"`
}

type ShoppingCartsFilterRequest struct {
	MerchantId string `json:"merchant_id"`
	Status     string `json:"status"`
}

type GetShoppingCartRequest struct {
	MerchantId     string `json:"merchant_id"`
	ShoppingCartId int64  `json:"shoppingCart_id"`
}
