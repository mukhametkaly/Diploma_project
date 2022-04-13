package product

import (
	"context"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
)

type service struct {
}

// Service is the interface that provides methods.
type Service interface {
	CreateProduct(request models.Product) (models.Product, error)
	UpdateProduct(request models.Product) (models.Product, error)
	DeleteByIdProduct(id int64) error
	DeleteBatchProduct(ids []int64) error
	GetProductById(id int64) (models.Product, error)
	FilterProducts(request FilterProductsRequest) ([]models.Product, error)
}

func NewService() Service {
	return &service{}
}

func (s *service) CreateProduct(request models.Product) (models.Product, error) {
	product, err := InsertProduct(context.Background(), request)
	if err != nil {
		return product, InternalServerError
	}
	return product, err
}

func (s *service) UpdateProduct(request models.Product) (models.Product, error) {
	var products models.Product

	err := UpdateProduct(context.Background(), request)
	if err != nil {
		return products, InternalServerError
	}
	return request, err
}

func (s *service) DeleteByIdProduct(id int64) error {

	err := DeleteProductById(context.Background(), id)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) DeleteBatchProduct(ids []int64) error {
	err := MDeleteProductByIds(context.Background(), ids)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) GetProductById(id int64) (models.Product, error) {
	product, err := GetProductById(context.Background(), id)
	if err != nil {
		return product, InternalServerError
	}
	return product, nil
}

func (s *service) FilterProducts(request FilterProductsRequest) ([]models.Product, error) {
	return []models.Product{}, nil
}
