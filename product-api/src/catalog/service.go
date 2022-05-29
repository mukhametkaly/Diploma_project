package catalog

import (
	"context"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
	"time"
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

	CreateCategory(request models.Category) (models.Category, error)
	UpdateCategory(request models.Category) (models.Category, error)
	DeleteByIdCategory(id int64) error
	DeleteBatchCategory(ids []int64) error
	GetCategoryById(id int64) (models.Category, error)
	FilterCategories(request FilterCategoryRequest) ([]models.Category, error)
}

func NewService() Service {
	return &service{}
}

func (s *service) CreateProduct(request models.Product) (models.Product, error) {
	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Product{}, InvalidCharacter
	}
	if request.CategoryName == "" {
		InvalidCharacter.Message = "no category"
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

	exist, err := CheckCategoryExists(context.Background(), request.MerchantId, request.CategoryName)
	if err != nil {
		return models.Product{}, InternalServerError
	}

	if !exist {
		category := models.Category{
			MerchantId:   request.MerchantId,
			CategoryName: request.CategoryName,
			CreatedOn:    time.Now(),
			UpdatedOn:    time.Now(),
		}
		_, err := InsertCategory(context.Background(), category)
		if err != nil {
			Loger.Debugf("%v", err.Error())
			InternalServerError.Message = err.Error()
			return models.Product{}, err
		}
	}

	product, err := InsertProduct(context.Background(), request)
	if err != nil {
		return product, InternalServerError
	}
	return product, err

}

func (s *service) UpdateProduct(request models.Product) (models.Product, error) {
	var product models.Product

	if request.CategoryName == "" {
		InvalidCharacter.Message = "no category"
		return models.Product{}, InvalidCharacter
	}

	exist, err := CheckCategoryExists(context.Background(), request.MerchantId, request.CategoryName)
	if err != nil {
		return models.Product{}, InternalServerError
	}

	if !exist {
		category := models.Category{
			MerchantId:   request.MerchantId,
			CategoryName: request.CategoryName,
			CreatedOn:    time.Now(),
			UpdatedOn:    time.Now(),
		}
		_, err := InsertCategory(context.Background(), category)
		if err != nil {
			Loger.Debugf("%v", err.Error())
			InternalServerError.Message = err.Error()
			return models.Product{}, err
		}
	}

	err = UpdateProduct(context.Background(), request)
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

func (s *service) CreateCategory(request models.Category) (models.Category, error) {
	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Category{}, InvalidCharacter
	}
	if request.CategoryName == "" {
		InvalidCharacter.Message = "no name"
		return models.Category{}, InvalidCharacter
	}

	barcodeExists, err := CheckCategoryExists(context.Background(), request.MerchantId, request.CategoryName)
	if err != nil {
		Loger.Debugf("%v", err.Error())
		InternalServerError.Message = err.Error()
		return models.Category{}, err
	}

	if barcodeExists {
		InvalidCharacter.Message = "catalog with same barcode exist"
		return models.Category{}, InvalidCharacter
	}

	category, err := InsertCategory(context.Background(), request)
	if err != nil {
		return category, InternalServerError
	}
	return category, err

}

func (s *service) UpdateCategory(request models.Category) (models.Category, error) {
	var category models.Category

	if request.CategoryName == "" {
		InvalidCharacter.Message = "no category"
		return models.Category{}, InvalidCharacter
	}

	if request.ID == 0 {
		InvalidCharacter.Message = "no category id"
		return models.Category{}, InvalidCharacter
	}

	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Category{}, InvalidCharacter
	}

	err := UpdateCategory(context.Background(), request)
	if err != nil {
		return category, InternalServerError
	}
	return request, nil
}

func (s *service) DeleteByIdCategory(id int64) error {

	err := DeleteCategoryById(context.Background(), id)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) DeleteBatchCategory(ids []int64) error {
	err := MDeleteCategoryByIds(context.Background(), ids)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) GetCategoryById(id int64) (models.Category, error) {
	category, err := GetCategoryById(context.Background(), id)
	if err != nil {
		return category, InternalServerError
	}
	return category, nil
}

func (s *service) FilterCategories(request FilterCategoryRequest) ([]models.Category, error) {
	if request.MerchantId == "" {
		return nil, InvalidCharacter
	}

	categories, err := FilterCategories(context.Background(), request)
	if err != nil {
		return nil, Conflict
	}

	return categories, nil
}
