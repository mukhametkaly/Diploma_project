package waybill

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
	"math/rand"
	"net/http"
	"time"
)

type waybillService struct {
}

type Service interface {
	CreateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error)
	UpdateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error)
	DeleteWaybill(request DeleteWaybillRequest) error
	GetWaybill(request GetWaybillRequest) (models.ShortWaybill, error)
	WaybillsFilter(request WaybillsFilterRequest) ([]models.ShortWaybill, error)

	WaybillAddProduct(product models.WaybillProduct) (models.WaybillProduct, error)
	WaybillUpdateProduct(product models.WaybillProduct) (models.WaybillProduct, error)
	DeleteWaybillProduct(request DeleteWaybillProductRequest) error
	GetWaybillProducts(request GetWaybillProductsRequest) ([]models.WaybillProduct, error)
}

func NewWaybillService() Service {
	return &waybillService{}
}

func (ws waybillService) CreateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error) {
	if waybill.MerchantId == "" {
		return waybill, newStringError(http.StatusBadRequest, "no merchant id")
	}

	waybill.CreatedOn = time.Now()
	waybill.UpdatedOn = waybill.CreatedOn
	waybill.Status = models.StatusDraft
	for {
		waybill.DocumentNumber = setDocNumber()
		docExist, err := IfDocNumberExist(context.TODO(), waybill.MerchantId, waybill.DocumentNumber)
		if err != nil {
			return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
		}
		if !docExist {
			break
		}
	}

	waybill, err := InsertWaybill(context.TODO(), waybill)
	if err != nil {
		return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
	}

	return waybill, err

}

func (ws waybillService) GetWaybill(request GetWaybillRequest) (models.ShortWaybill, error) {

	waybill, err := GetWaybillById(context.Background(), request.WaybillId)
	if err != nil {
		return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
	}

	return waybill, nil

}

func (ws waybillService) UpdateWaybill(waybill models.ShortWaybill) (models.ShortWaybill, error) {

	if waybill.MerchantId == "" {
		return waybill, newStringError(http.StatusBadRequest, "no merchant id")
	}

	waybill.UpdatedOn = time.Now()

	oldWaybill, err := GetWaybillById(context.Background(), waybill.ID)
	if err != nil {
		return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
	}

	if oldWaybill.Status == models.StatusDraft {
		if waybill.Status == models.StatusProvided {
			waybill.ProvidedTime = waybill.UpdatedOn
			err = UpdateWaybillStatus(context.TODO(), waybill)
			if err != nil {
				return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
			}

			products, err := GetWaybillProducts(context.Background(), GetWaybillProductsRequest{WaybillId: waybill.ID})
			if err != nil {
				return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
			}

			err = SendProductsToUpdateCount(products, oldWaybill.MerchantId)
			if err != nil {
				return models.ShortWaybill{}, newError(http.StatusInternalServerError, err)
			}
		}
	}

	return waybill, err

}

func (ws waybillService) WaybillAddProduct(product models.WaybillProduct) (models.WaybillProduct, error) {
	if product.WaybillId == 0 {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "no such waybill")
	}

	if product.Barcode == "" {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortWaybill, err := GetWaybillById(context.Background(), product.WaybillId)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortWaybill.Status == models.StatusProvided {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "you can't add product to provided waybill")
	}

	req := GetWaybillProductsRequest{
		WaybillId: product.WaybillId,
		Barcode:   product.Barcode,
	}
	oldProducts, err := GetWaybillProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) != 0 {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "product already added")
	}

	product.CreatedOn = time.Now()

	err = InsertWaybillProduct(context.Background(), product)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	waybillSum := product.PurchasePrice * product.ReceivedAmount
	err = UpdateWaybillSum(context.Background(), product.WaybillId, waybillSum)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err

}

func (ws waybillService) DeleteWaybillProduct(request DeleteWaybillProductRequest) error {
	if request.WaybillId == 0 {
		return newStringError(http.StatusBadRequest, "no such waybill")
	}

	if request.Barcode == "" {
		return newStringError(http.StatusBadRequest, "no barcode")
	}

	shortWaybill, err := GetWaybillById(context.Background(), request.WaybillId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	if shortWaybill.Status == models.StatusProvided {
		return newStringError(http.StatusBadRequest, "you can't remove product from provided waybill")
	}

	req := GetWaybillProductsRequest{
		WaybillId: request.WaybillId,
		Barcode:   request.Barcode,
	}
	oldProducts, err := GetWaybillProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return nil
	}

	oldProduct := oldProducts[0]

	waybillSum := oldProduct.PurchasePrice * oldProduct.ReceivedAmount * -1
	err = UpdateWaybillSum(context.Background(), request.WaybillId, waybillSum)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	err = DeleteWaybillProduct(context.Background(), request)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws waybillService) WaybillUpdateProduct(product models.WaybillProduct) (models.WaybillProduct, error) {
	if product.WaybillId == 0 {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "no such waybill")
	}

	if product.Barcode == "" {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortWaybill, err := GetWaybillById(context.Background(), product.WaybillId)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortWaybill.Status == models.StatusProvided {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "you can't update product from provided waybill")
	}

	req := GetWaybillProductsRequest{
		WaybillId: product.WaybillId,
		Barcode:   product.Barcode,
	}
	oldProducts, err := GetWaybillProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return models.WaybillProduct{}, newStringError(http.StatusBadRequest, "no such product")
	}

	oldProduct := oldProducts[0]

	waybillSum := (product.PurchasePrice * product.ReceivedAmount) - (oldProduct.PurchasePrice * oldProduct.ReceivedAmount)
	err = UpdateWaybillSum(context.Background(), product.WaybillId, waybillSum)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	err = UpdateWaybillProduct(context.Background(), product)
	if err != nil {
		return models.WaybillProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err
}

func (ws waybillService) DeleteWaybill(request DeleteWaybillRequest) error {
	if request.MerchantId == "" {
		return newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	if request.WaybillId == 0 {
		return newStringError(http.StatusBadRequest, "waybill id is empty")
	}

	err := DeleteWaybillById(context.TODO(), request.WaybillId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws waybillService) GetWaybillProducts(request GetWaybillProductsRequest) ([]models.WaybillProduct, error) {
	if request.WaybillId == 0 {
		return nil, newStringError(http.StatusBadRequest, "waybill id is empty")
	}
	products, err := GetWaybillProducts(context.TODO(), request)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, err)
	}

	return products, nil
}

func (ws waybillService) WaybillsFilter(request WaybillsFilterRequest) ([]models.ShortWaybill, error) {
	if request.MerchantId == "" {
		return nil, newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	waybills, err := GetWaybill(context.TODO(), request)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, err)
	}
	return waybills, nil

}

func setDocNumber() string {
	docNum := fmt.Sprintf("%v", rand.Intn(1000000))
	for i := len(docNum); i < 6; i++ {
		docNum = fmt.Sprintf("%v%v", "0", docNum)
	}
	return docNum
}
