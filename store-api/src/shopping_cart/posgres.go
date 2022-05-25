package shopping_cart

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/store-api/src/config"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
	"github.com/sirupsen/logrus"
	"time"
)

var db *pg.DB

func PGConnectStart() (*pg.DB, error) {
	conn := pg.Connect(&pg.Options{
		Addr:               fmt.Sprintf("%s:%s", config.AllConfigs.Postgres.Host, config.AllConfigs.Postgres.Port),
		User:               config.AllConfigs.Postgres.User,
		Password:           config.AllConfigs.Postgres.Password,
		Database:           config.AllConfigs.Postgres.DBName,
		IdleTimeout:        59 * time.Second,
		IdleCheckFrequency: 30 * time.Second,
	})

	err := conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetPGSession() (*pg.DB, error) {

	if db == nil {
		client, err := PGConnectStart()
		if err != nil {
			return nil, err
		} else {
			db = client
			return db, nil
		}
	} else {
		return db, nil
	}
}

// Repository is postgres repository
type Repository struct {
	db *pg.DB

	logger logrus.Logger
}

type ShoppingCartDTO struct {
	tableName    struct{} `pg:"shopping_carts"`
	Id           int64    `pg:",pk,unique"`
	MerchantId   string
	TotalSum     float64
	CreatedOn    time.Time
	ProvidedTime time.Time
	Status       string
}

func (d *ShoppingCartDTO) fromDTO() models.ShoppingCart {
	var shoppingCart models.ShoppingCart
	shoppingCart.ID = d.Id
	shoppingCart.MerchantId = d.MerchantId
	shoppingCart.TotalSum = d.TotalSum
	shoppingCart.CreatedOn = d.CreatedOn
	shoppingCart.ProvidedTime = d.ProvidedTime
	shoppingCart.Status = d.Status
	return shoppingCart
}

func (d *ShoppingCartDTO) toDTO(shoppingCart models.ShoppingCart) {
	d.Id = shoppingCart.ID
	d.MerchantId = shoppingCart.MerchantId
	d.TotalSum = shoppingCart.TotalSum
	d.CreatedOn = shoppingCart.CreatedOn
	d.ProvidedTime = shoppingCart.ProvidedTime
	d.Status = shoppingCart.Status
}

type ShoppingCartProductsDTO struct {
	tableName      struct{} `pg:"shopping_cart_products"`
	Barcode        string
	Name           string    `json:"name"`
	Amount         float64   `json:"amount"`
	ShoppingCartId int64     `json:"shopping_cart_id"`
	PurchasePrice  float64   `json:"purchase_price"`
	SellingPrice   float64   `json:"selling_price"`
	Total          float64   `json:"total"`
	CreatedOn      time.Time `json:"created_on"`
}

func (d *ShoppingCartProductsDTO) fromDTO() models.ShoppingCartProduct {
	var product models.ShoppingCartProduct
	product.Barcode = d.Barcode
	product.Name = d.Name
	product.Amount = d.Amount
	product.ShoppingCartId = d.ShoppingCartId
	product.PurchasePrice = d.PurchasePrice
	product.SellingPrice = d.SellingPrice
	product.Total = d.Total
	product.CreatedOn = d.CreatedOn
	return product
}

func (d *ShoppingCartProductsDTO) toDTO(product models.ShoppingCartProduct) {
	d.Barcode = product.Barcode
	d.Name = product.Name
	d.Amount = product.Amount
	d.ShoppingCartId = product.ShoppingCartId
	d.PurchasePrice = product.PurchasePrice
	d.SellingPrice = product.SellingPrice
	d.Total = product.Total
	d.CreatedOn = product.CreatedOn
}

func InsertShoppingCart(ctx context.Context, shoppingCart models.ShoppingCart) (models.ShoppingCart, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertShoppingCart", err.Error())
		return models.ShoppingCart{}, err
	}

	var dtoShoppingCart ShoppingCartDTO
	dtoShoppingCart.toDTO(shoppingCart)

	_, err = conn.ModelContext(ctx, &dtoShoppingCart).Returning("*", &dtoShoppingCart).Insert(&dtoShoppingCart)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return models.ShoppingCart{}, err
	}

	shoppingCart = dtoShoppingCart.fromDTO()

	return shoppingCart, nil
}

func UpdateShoppingCartStatus(ctx context.Context, shoppingCart models.ShoppingCart) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	var dtoShoppingCart ShoppingCartDTO
	dtoShoppingCart.toDTO(shoppingCart)

	_, err = conn.ModelContext(ctx, &dtoShoppingCart).WherePK().Column("updated_on", "provided_time", "status").Update()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func UpdateShoppingCartSum(ctx context.Context, id int64, sum float64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	var dtoShoppingCart ShoppingCartDTO
	dtoShoppingCart.Id = id

	_, err = conn.ModelContext(ctx, &dtoShoppingCart).WherePK().Set("total_sum = total_sum + ?", sum).Column("updated_on").Update()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func DeleteShoppingCartById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteShoppingCartById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*ShoppingCartDTO)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func GetShoppingCart(ctx context.Context, req ShoppingCartsFilterRequest) ([]models.ShoppingCart, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetShoppingCart", err.Error())
		return nil, err
	}

	shoppingCartDto := []ShoppingCartDTO{}

	query := conn.ModelContext(ctx, &shoppingCartDto).Where("merchant_id = ?", req.MerchantId)

	if req.Status != "" {
		query.Where("status = ?", req.Status)
	}

	err = query.Select()
	if err != nil {
		Loger.Debugln("error select in get list inventories", err.Error())
		return nil, err
	}

	inventories := make([]models.ShoppingCart, 0, len(shoppingCartDto))

	for _, dto := range shoppingCartDto {
		shoppingCart := dto.fromDTO()
		inventories = append(inventories, shoppingCart)
	}

	return inventories, err
}

func IfDocNumberExist(ctx context.Context, merchantId, docNum string) (bool, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetShoppingCart", err.Error())
		return false, err
	}

	shoppingCartDto := []ShoppingCartDTO{}

	return conn.ModelContext(ctx, &shoppingCartDto).
		Where("merchant_id = ?", merchantId).
		Where("document_number = ?", docNum).
		Exists()

}

func GetShoppingCartById(ctx context.Context, id int64) (models.ShoppingCart, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetShoppingCart", err.Error())
		return models.ShoppingCart{}, err
	}

	shoppingCartDto := ShoppingCartDTO{}
	shoppingCartDto.Id = id

	err = conn.ModelContext(ctx, &shoppingCartDto).WherePK().Select()
	if err != nil {
		Loger.Debugln("error select in get list inventories", err.Error())
		return models.ShoppingCart{}, err
	}

	shoppingCart := shoppingCartDto.fromDTO()

	return shoppingCart, err
}

func GetShoppingCartProducts(ctx context.Context, req GetShoppingCartProductsRequest) (product []models.ShoppingCartProduct, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetShoppingCartProducts", err.Error())
		return
	}

	DtoProducts := []ShoppingCartProductsDTO{}

	query := conn.ModelContext(ctx, &DtoProducts).Where("shoppingCart_id = ?", req.ShoppingCartId)
	if req.Barcode != "" {
		query.Where("barcode = ?", req.Barcode)
	}

	err = query.Order("created_on ASC").Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	for _, productDTO := range DtoProducts {
		product = append(product, productDTO.fromDTO())
	}

	return
}

func DeleteShoppingCartProduct(ctx context.Context, req DeleteShoppingCartProductRequest) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteShoppingCartProduct", err.Error())
		return
	}

	_, err = conn.Model((*ShoppingCartProductsDTO)(nil)).
		Where("shoppingCart_id = ?", req.ShoppingCartId).
		Where("barcode = ?", req.Barcode).
		Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func InsertShoppingCartProduct(ctx context.Context, product models.ShoppingCartProduct) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertShoppingCartProduct", err.Error())
		return
	}

	var dtoProduct ShoppingCartProductsDTO
	dtoProduct.toDTO(product)

	_, err = conn.ModelContext(ctx, &dtoProduct).Insert(&dtoProduct)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func UpdateShoppingCartProduct(ctx context.Context, product models.ShoppingCartProduct) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertShoppingCartProduct", err.Error())
		return
	}

	var dtoProduct ShoppingCartProductsDTO
	dtoProduct.toDTO(product)

	_, err = conn.ModelContext(ctx, &dtoProduct).
		Where("barcode = ?", dtoProduct.Barcode).
		Where("shoppingCart_id = ?", dtoProduct.ShoppingCartId).
		Update("received_amount", "amount", "purchase_price", "selling_price", "total")
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}
