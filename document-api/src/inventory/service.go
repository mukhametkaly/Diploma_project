package inventory

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
	"math/rand"
	"net/http"
	"time"
)

type inventoryService struct {
}

type Service interface {
	CreateInventory(inventory models.ShortInventory) (models.ShortInventory, error)
	UpdateInventory(inventory models.ShortInventory) (models.ShortInventory, error)
	DeleteInventory(request DeleteInventoryRequest) error
	GetInventory(request GetInventoryRequest) (models.ShortInventory, error)
	InventorysFilter(request InventorysFilterRequest) ([]models.ShortInventory, error)

	InventoryAddProduct(product models.InventoryProduct) (models.InventoryProduct, error)
	InventoryUpdateProduct(product models.InventoryProduct) (models.InventoryProduct, error)
	DeleteInventoryProduct(request DeleteInventoryProductRequest) error
	GetInventoryProducts(request GetInventoryProductsRequest) ([]models.InventoryProduct, error)
}

func NewInventoryService() Service {
	return &inventoryService{}
}

func (ws inventoryService) CreateInventory(inventory models.ShortInventory) (models.ShortInventory, error) {
	if inventory.MerchantId == "" {
		return inventory, newStringError(http.StatusBadRequest, "no merchant id")
	}

	inventory.CreatedOn = time.Now()
	inventory.UpdatedOn = inventory.CreatedOn
	inventory.Status = models.StatusDraft
	for {
		inventory.DocumentNumber = setDocNumber()
		docExist, err := IfDocNumberExist(context.TODO(), inventory.MerchantId, inventory.DocumentNumber)
		if err != nil {
			return models.ShortInventory{}, newError(http.StatusInternalServerError, err)
		}
		if !docExist {
			break
		}
	}

	inventory, err := InsertInventory(context.TODO(), inventory)
	if err != nil {
		return models.ShortInventory{}, newError(http.StatusInternalServerError, err)
	}

	return inventory, err

}

func (ws inventoryService) GetInventory(request GetInventoryRequest) (models.ShortInventory, error) {

	inventory, err := GetInventoryById(context.Background(), request.InventoryId)
	if err != nil {
		return models.ShortInventory{}, newError(http.StatusInternalServerError, err)
	}

	return inventory, nil

}

func (ws inventoryService) UpdateInventory(inventory models.ShortInventory) (models.ShortInventory, error) {

	if inventory.MerchantId == "" {
		return inventory, newStringError(http.StatusBadRequest, "no merchant id")
	}

	inventory.UpdatedOn = time.Now()

	oldInventory, err := GetInventoryById(context.Background(), inventory.ID)
	if err != nil {
		return models.ShortInventory{}, newError(http.StatusInternalServerError, err)
	}

	if oldInventory.Status == models.StatusDraft {
		if inventory.Status == models.StatusProvided {
			inventory.ProvidedTime = inventory.UpdatedOn
			err = UpdateInventoryStatus(context.TODO(), inventory)
			if err != nil {
				return models.ShortInventory{}, newError(http.StatusInternalServerError, err)
			}
		}
	}

	return inventory, err

}

func (ws inventoryService) InventoryAddProduct(product models.InventoryProduct) (models.InventoryProduct, error) {
	if product.InventoryId == 0 {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "no such inventory")
	}

	if product.Barcode == "" {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortInventory, err := GetInventoryById(context.Background(), product.InventoryId)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortInventory.Status == models.StatusProvided {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "you can't add product to provided inventory")
	}

	req := GetInventoryProductsRequest{
		InventoryId: product.InventoryId,
		Barcode:     product.Barcode,
	}
	oldProducts, err := GetInventoryProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) != 0 {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "product already added")
	}

	product.CreatedOn = time.Now()

	err = InsertInventoryProduct(context.Background(), product)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	inventorySum := product.PurchasePrice * product.ActualAmount
	err = UpdateInventorySum(context.Background(), product.InventoryId, inventorySum)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err

}

func (ws inventoryService) DeleteInventoryProduct(request DeleteInventoryProductRequest) error {
	if request.InventoryId == 0 {
		return newStringError(http.StatusBadRequest, "no such inventory")
	}

	if request.Barcode == "" {
		return newStringError(http.StatusBadRequest, "no barcode")
	}

	shortInventory, err := GetInventoryById(context.Background(), request.InventoryId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	if shortInventory.Status == models.StatusProvided {
		return newStringError(http.StatusBadRequest, "you can't remove product from provided inventory")
	}

	req := GetInventoryProductsRequest{
		InventoryId: request.InventoryId,
		Barcode:     request.Barcode,
	}
	oldProducts, err := GetInventoryProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return nil
	}

	oldProduct := oldProducts[0]

	inventorySum := oldProduct.PurchasePrice * oldProduct.ActualAmount * -1
	err = UpdateInventorySum(context.Background(), request.InventoryId, inventorySum)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	err = DeleteInventoryProduct(context.Background(), request)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws inventoryService) InventoryUpdateProduct(product models.InventoryProduct) (models.InventoryProduct, error) {
	if product.InventoryId == 0 {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "no such inventory")
	}

	if product.Barcode == "" {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortInventory, err := GetInventoryById(context.Background(), product.InventoryId)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortInventory.Status == models.StatusProvided {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "you can't update product from provided inventory")
	}

	req := GetInventoryProductsRequest{
		InventoryId: product.InventoryId,
		Barcode:     product.Barcode,
	}
	oldProducts, err := GetInventoryProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return models.InventoryProduct{}, newStringError(http.StatusBadRequest, "no such product")
	}

	oldProduct := oldProducts[0]

	inventorySum := (product.PurchasePrice * product.ActualAmount) - (oldProduct.PurchasePrice * oldProduct.ActualAmount)
	err = UpdateInventorySum(context.Background(), product.InventoryId, inventorySum)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	err = UpdateInventoryProduct(context.Background(), product)
	if err != nil {
		return models.InventoryProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err
}

func (ws inventoryService) DeleteInventory(request DeleteInventoryRequest) error {
	if request.MerchantId == "" {
		return newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	if request.InventoryId == 0 {
		return newStringError(http.StatusBadRequest, "inventory id is empty")
	}

	err := DeleteInventoryById(context.TODO(), request.InventoryId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws inventoryService) GetInventoryProducts(request GetInventoryProductsRequest) ([]models.InventoryProduct, error) {
	if request.InventoryId == 0 {
		return nil, newStringError(http.StatusBadRequest, "inventory id is empty")
	}
	products, err := GetInventoryProducts(context.TODO(), request)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, err)
	}

	return products, nil
}

func (ws inventoryService) InventorysFilter(request InventorysFilterRequest) ([]models.ShortInventory, error) {
	if request.MerchantId == "" {
		return nil, newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	inventories, err := GetInventory(context.TODO(), request)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, err)
	}
	return inventories, nil

}

func setDocNumber() string {
	docNum := fmt.Sprintf("%v", rand.Intn(1000000))
	for i := len(docNum); i < 6; i++ {
		docNum = fmt.Sprintf("%v%v", "0", docNum)
	}
	return docNum
}
