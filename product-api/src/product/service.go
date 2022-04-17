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
	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Product{}, InvalidCharacter
	}
	if request.CategoryId == 0 {
		InvalidCharacter.Message = "no category id"
		return models.Product{}, InvalidCharacter
	}
	if request.Name == "" {
		InvalidCharacter.Message = "no product name"
		return models.Product{}, InvalidCharacter
	}
	if request.Barcode == "" {
		InvalidCharacter.Message = "no product barcode"
		return models.Product{}, InvalidCharacter
	}

	if request.UnitType != "piece" && request.UnitType != "weight" {
		request.UnitType = "piece"
	}

	barcodeExists, err := CheckBarcode(context.Background(), request.MerchantId, request.Barcode)
	if err != nil {
		Loger.Debugf("%v", err.Error())
		InternalServerError.Message = err.Error()
		return models.Product{}, err
	}

	if barcodeExists {
		InvalidCharacter.Message = "product with same barcode exist"
		return models.Product{}, InvalidCharacter
	}

	product, err := InsertProduct(context.Background(), request)
	if err != nil {
		return product, InternalServerError
	}
	return product, err

}

func (s *service) UpdateProduct(request models.Product) (models.Product, error) {
	var product models.Product

	if request.CategoryId == 0 {
		InvalidCharacter.Message = "no category id"
		return models.Product{}, InvalidCharacter
	}

	err := UpdateProduct(context.Background(), request)
	if err != nil {
		return product, InternalServerError
	}
	return request, nil
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
	if request.MerchantId == "" {
		return nil, InvalidCharacter
	}

	products, err := FilterProducts(context.TODO(), request)
	if err != nil {
		return nil, Conflict
	}

	return products, nil
}
