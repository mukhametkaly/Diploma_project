package waybill

import (
	"context"
	"fmt"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
	"math/rand"
	"net/http"
	"time"
)

type waybillService struct {
}

type WaybillService interface {
	CreateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error)
	UpdateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error)
	DeleteWaybill(request DeleteWaybillRequest) error
	DeleteWaybillProduct(request DeleteWaybillProductRequest) error
	GetWaybillProducts(request GetWaybillProductsRequest) ([]models.WaybillProduct, error)
	WaybillsFilter(request WaybillsFilterRequest) ([]models.ShortWaybill, error)
	WaybillAddProduct(product models.WaybillProduct) ([]models.WaybillProduct, error)
	WaybillRemoveProduct(product models.WaybillProduct) ([]models.WaybillProduct, error)
}

func NewWaybillService() WaybillService {
	return &waybillService{}
}

func (ws waybillService) CreateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error) {

	if waybill.MerchantId == "" {
		return waybill, newStringError(http.StatusBadRequest, "no merchant id")
	}

	waybill.CreatedOn = time.Now()
	waybill.UpdatedOn = waybill.CreatedOn
	waybill.Status = "draft"
	waybill.DocumentNumber = setDocNumber()

	waybill, err := InsertWaybill(context.TODO(), waybill)

}

func (ws waybillService) UpdateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error) {

}

func (ws waybillService) DeleteWaybillProduct(request DeleteWaybillProductRequest) error {

}

func (ws waybillService) WaybillAddProduct(product models.WaybillProduct) ([]models.WaybillProduct, error) {

}

func (ws waybillService) WaybillRemoveProduct(product models.WaybillProduct) ([]models.WaybillProduct, error) {

}

func (ws waybillService) DeleteWaybill(request DeleteWaybillRequest) error {

}

func (ws waybillService) GetWaybillProducts(request GetWaybillProductsRequest) ([]models.WaybillProduct, error) {

}

func (ws waybillService) WaybillsFilter(request WaybillsFilterRequest) ([]models.ShortWaybill, error) {

}

func setDocNumber() string {
	docNum := fmt.Sprintf("%v", rand.Intn(1000000))
	for i := len(docNum); i < 6; i++ {
		docNum := fmt.Sprintf("%v%v", "0", docNum)
	}
	return docNum
}
