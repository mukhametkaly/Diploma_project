package catalog

import (
	"context"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
	"time"
)

type productService struct {
}

// ProductService is the interface that provides methods.
type ProductService interface {
	CreateProduct(request models.Product) (models.Product, error)
	UpdateProduct(request models.Product) (models.Product, error)
	DeleteByIdProduct(id int64) error
	DeleteBatchProduct(ids []int64) error
	GetProductById(id int64) (models.Product, error)
	FilterProducts(request FilterProductsRequest) ([]models.Product, error)
}

func NewService() ProductService {
	return &productService{}
}

func (s *productService) CreateProduct(request models.Product) (models.Product, error) {
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

func (s *productService) UpdateProduct(request models.Product) (models.Product, error) {
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

func (s *productService) DeleteByIdProduct(id int64) error {

	err := DeleteProductById(context.Background(), id)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *productService) DeleteBatchProduct(ids []int64) error {
	err := MDeleteProductByIds(context.Background(), ids)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *productService) GetProductById(id int64) (models.Product, error) {
	product, err := GetProductById(context.Background(), id)
	if err != nil {
		return product, InternalServerError
	}
	return product, nil
}

func (s *productService) FilterProducts(request FilterProductsRequest) ([]models.Product, error) {
	if request.MerchantId == "" {
		return nil, InvalidCharacter
	}

	products, err := FilterProducts(context.TODO(), request)
	if err != nil {
		return nil, Conflict
	}

	return products, nil
}
