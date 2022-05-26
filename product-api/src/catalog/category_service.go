package catalog

import (
	"context"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
)

type categoryService struct {
}

// CategoryService is the interface that provides methods.
type CategoryService interface {
	CreateCategory(request models.Category) (models.Category, error)
	UpdateCategory(request models.Category) (models.Category, error)
	DeleteByIdCategory(id int64) error
	DeleteBatchCategory(ids []int64) error
	GetCategoryById(id int64) (models.Category, error)
	FilterCategories(merchant string) ([]models.Category, error)
}

func NewCategoryService() CategoryService {
	return &categoryService{}
}

func (s *categoryService) CreateCategory(request models.Category) (models.Category, error) {
	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Category{}, InvalidCharacter
	}
	if request.CategoryName == "" {
		InvalidCharacter.Message = "no category"
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

func (s *categoryService) UpdateCategory(request models.Category) (models.Category, error) {
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

func (s *categoryService) DeleteByIdCategory(id int64) error {

	err := DeleteCategoryById(context.Background(), id)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *categoryService) DeleteBatchCategory(ids []int64) error {
	err := MDeleteCategoryByIds(context.Background(), ids)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *categoryService) GetCategoryById(id int64) (models.Category, error) {
	category, err := GetCategoryById(context.Background(), id)
	if err != nil {
		return category, InternalServerError
	}
	return category, nil
}

func (s *categoryService) FilterCategories(merchantId string) ([]models.Category, error) {
	if merchantId == "" {
		return nil, InvalidCharacter
	}

	categories, err := FilterCategories(context.Background(), merchantId)
	if err != nil {
		return nil, Conflict
	}

	return categories, nil
}
