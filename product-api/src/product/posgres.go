package product

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/product-api/src/config"
	"github.com/mukhametkaly/Diploma/product-api/src/models"
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

type ProductDTO struct {
	tableName     struct{} `pg:"products"`
	Id            int64    `pg:",pk,unique"`
	Barcode       string
	Name          string
	CategoryId    int64
	MerchantId    int64
	StockId       int64
	PurchasePrice float64
	SellingPrice  float64
	Amount        float64
	UnitType      string
}

func (d *ProductDTO) fromDTO() models.Product {
	var product models.Product
	product.ID = d.Id
	product.Barcode = d.Barcode
	product.Name = d.Name
	product.CategoryId = d.CategoryId
	product.MerchantId = d.MerchantId
	product.StockId = d.StockId
	product.PurchasePrice = d.PurchasePrice
	product.SellingPrice = d.SellingPrice
	product.Amount = d.Amount
	product.UnitType = d.UnitType
	return product
}

func (d *ProductDTO) toDTO(product models.Product) {
	d.Id = product.ID
	d.Barcode = product.Barcode
	d.Name = product.Name
	d.CategoryId = product.CategoryId
	d.MerchantId = product.MerchantId
	d.StockId = product.StockId
	d.PurchasePrice = product.PurchasePrice
	d.SellingPrice = product.SellingPrice
	d.Amount = product.Amount
	d.UnitType = product.UnitType
}

func GetProductById(ctx context.Context, id int64) (product models.Product, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	q := conn.ModelContext(ctx, &product).Where("id = ?", id)
	err = q.Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func InsertProduct(ctx context.Context, product models.Product) (result models.Product, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	productDto := ProductDTO{}
	productDto.toDTO(product)

	_, err = conn.ModelContext(ctx, &productDto).Returning("*", productDto).Insert()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	result = productDto.fromDTO()

	return
}

func UpdateProduct(ctx context.Context, product models.Product) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &product).WherePK().Update(product)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func DeleteProductById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	_, err = conn.Model((*models.Product)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func MDeleteProductByIds(ctx context.Context, ids []int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	_, err = conn.Model((*models.Product)(nil)).Where("id IN (?)", pg.In(ids)).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}
