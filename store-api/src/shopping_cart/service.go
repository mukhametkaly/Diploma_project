package shopping_cart

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
	"math/rand"
	"net/http"
	"time"
)

type shoppingCartService struct {
}

type Service interface {
	CreateShoppingCart(shoppingCart models.ShoppingCart) (models.ShoppingCart, error)
	UpdateShoppingCart(shoppingCart models.ShoppingCart) (models.ShoppingCart, error)
	DeleteShoppingCart(request DeleteShoppingCartRequest) error
	GetShoppingCart(request GetShoppingCartRequest) (models.ShoppingCart, error)
	ShoppingCartsFilter(request ShoppingCartsFilterRequest) ([]models.ShoppingCart, error)

	ShoppingCartAddProduct(product models.ShoppingCartProduct) (models.ShoppingCartProduct, error)
	ShoppingCartUpdateProduct(product models.ShoppingCartProduct) (models.ShoppingCartProduct, error)
	DeleteShoppingCartProduct(request DeleteShoppingCartProductRequest) error
	GetShoppingCartProducts(request GetShoppingCartProductsRequest) ([]models.ShoppingCartProduct, error)
}

func NewShoppingCartService() Service {
	return &shoppingCartService{}
}

func (ws shoppingCartService) CreateShoppingCart(shoppingCart models.ShoppingCart) (models.ShoppingCart, error) {
	if shoppingCart.MerchantId == "" {
		return shoppingCart, newStringError(http.StatusBadRequest, "no merchant id")
	}

	shoppingCart.CreatedOn = time.Now()
	shoppingCart.Status = models.StatusDraft
	shoppingCart, err := InsertShoppingCart(context.TODO(), shoppingCart)
	if err != nil {
		return models.ShoppingCart{}, newError(http.StatusInternalServerError, err)
	}

	return shoppingCart, err

}

func (ws shoppingCartService) GetShoppingCart(request GetShoppingCartRequest) (models.ShoppingCart, error) {

	shoppingCart, err := GetShoppingCartById(context.Background(), request.ShoppingCartId)
	if err != nil {
		return models.ShoppingCart{}, newError(http.StatusInternalServerError, err)
	}

	return shoppingCart, nil

}

func (ws shoppingCartService) UpdateShoppingCart(shoppingCart models.ShoppingCart) (models.ShoppingCart, error) {

	if shoppingCart.MerchantId == "" {
		return shoppingCart, newStringError(http.StatusBadRequest, "no merchant id")
	}

	oldShoppingCart, err := GetShoppingCartById(context.Background(), shoppingCart.ID)
	if err != nil {
		return models.ShoppingCart{}, newError(http.StatusInternalServerError, err)
	}

	if oldShoppingCart.Status == models.StatusDraft {
		if shoppingCart.Status == models.StatusProvided {
			err = UpdateShoppingCartStatus(context.TODO(), shoppingCart)
			if err != nil {
				return models.ShoppingCart{}, newError(http.StatusInternalServerError, err)
			}
		}
	}

	return shoppingCart, err

}

func (ws shoppingCartService) ShoppingCartAddProduct(product models.ShoppingCartProduct) (models.ShoppingCartProduct, error) {
	if product.ShoppingCartId == 0 {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "no such shopping_cart")
	}

	if product.Barcode == "" {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortShoppingCart, err := GetShoppingCartById(context.Background(), product.ShoppingCartId)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortShoppingCart.Status == models.StatusProvided {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "you can't add merchant to provided shopping_cart")
	}

	req := GetShoppingCartProductsRequest{
		ShoppingCartId: product.ShoppingCartId,
		Barcode:        product.Barcode,
	}
	oldProducts, err := GetShoppingCartProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) != 0 {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "merchant already added")
	}

	product.CreatedOn = time.Now()

	err = InsertShoppingCartProduct(context.Background(), product)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	shoppingCartSum := product.PurchasePrice * product.Amount
	err = UpdateShoppingCartSum(context.Background(), product.ShoppingCartId, shoppingCartSum)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err

}

func (ws shoppingCartService) DeleteShoppingCartProduct(request DeleteShoppingCartProductRequest) error {
	if request.ShoppingCartId == 0 {
		return newStringError(http.StatusBadRequest, "no such shopping_cart")
	}

	if request.Barcode == "" {
		return newStringError(http.StatusBadRequest, "no barcode")
	}

	shortShoppingCart, err := GetShoppingCartById(context.Background(), request.ShoppingCartId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	if shortShoppingCart.Status == models.StatusProvided {
		return newStringError(http.StatusBadRequest, "you can't remove merchant from provided shopping_cart")
	}

	req := GetShoppingCartProductsRequest{
		ShoppingCartId: request.ShoppingCartId,
		Barcode:        request.Barcode,
	}
	oldProducts, err := GetShoppingCartProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return nil
	}

	oldProduct := oldProducts[0]

	shoppingCartSum := oldProduct.PurchasePrice * oldProduct.Amount * -1
	err = UpdateShoppingCartSum(context.Background(), request.ShoppingCartId, shoppingCartSum)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	err = DeleteShoppingCartProduct(context.Background(), request)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws shoppingCartService) ShoppingCartUpdateProduct(product models.ShoppingCartProduct) (models.ShoppingCartProduct, error) {
	if product.ShoppingCartId == 0 {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "no such shopping_cart")
	}

	if product.Barcode == "" {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "no barcode")
	}

	shortShoppingCart, err := GetShoppingCartById(context.Background(), product.ShoppingCartId)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	if shortShoppingCart.Status == models.StatusProvided {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "you can't update merchant from provided shopping_cart")
	}

	req := GetShoppingCartProductsRequest{
		ShoppingCartId: product.ShoppingCartId,
		Barcode:        product.Barcode,
	}
	oldProducts, err := GetShoppingCartProducts(context.Background(), req)
	if err != nil && err != pg.ErrNoRows {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	if len(oldProducts) == 0 {
		return models.ShoppingCartProduct{}, newStringError(http.StatusBadRequest, "no such merchant")
	}

	oldProduct := oldProducts[0]

	shoppingCartSum := (product.PurchasePrice * product.Amount) - (oldProduct.PurchasePrice * oldProduct.Amount)
	err = UpdateShoppingCartSum(context.Background(), product.ShoppingCartId, shoppingCartSum)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	err = UpdateShoppingCartProduct(context.Background(), product)
	if err != nil {
		return models.ShoppingCartProduct{}, newError(http.StatusInternalServerError, err)
	}

	return product, err
}

func (ws shoppingCartService) DeleteShoppingCart(request DeleteShoppingCartRequest) error {
	if request.MerchantId == "" {
		return newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	if request.ShoppingCartId == 0 {
		return newStringError(http.StatusBadRequest, "shopping_cart id is empty")
	}

	err := DeleteShoppingCartById(context.TODO(), request.ShoppingCartId)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (ws shoppingCartService) GetShoppingCartProducts(request GetShoppingCartProductsRequest) ([]models.ShoppingCartProduct, error) {
	if request.ShoppingCartId == 0 {
		return nil, newStringError(http.StatusBadRequest, "shopping_cart id is empty")
	}
	products, err := GetShoppingCartProducts(context.TODO(), request)
	if err != nil {
		return nil, newError(http.StatusInternalServerError, err)
	}

	return products, nil
}

func (ws shoppingCartService) ShoppingCartsFilter(request ShoppingCartsFilterRequest) ([]models.ShoppingCart, error) {
	if request.MerchantId == "" {
		return nil, newStringError(http.StatusBadRequest, "merchant id is empty")
	}
	inventories, err := GetShoppingCart(context.TODO(), request)
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
