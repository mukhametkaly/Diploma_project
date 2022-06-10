package catalog

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
	CategoryName  string
	MerchantId    string
	StockId       int64
	PurchasePrice float64
	SellingPrice  float64
	Amount        float64
	Margin        float64
	UnitType      string
}

func (d *ProductDTO) fromDTO() models.Product {
	var product models.Product
	product.ID = d.Id
	product.Barcode = d.Barcode
	product.Name = d.Name
	product.CategoryId = d.CategoryId
	product.CategoryName = d.CategoryName
	product.MerchantId = d.MerchantId
	product.StockId = d.StockId
	product.PurchasePrice = d.PurchasePrice
	product.SellingPrice = d.SellingPrice
	product.Margin = d.Margin
	product.Amount = d.Amount
	product.UnitType = d.UnitType
	return product
}

func (d *ProductDTO) toDTO(product models.Product) {
	d.Id = product.ID
	d.Barcode = product.Barcode
	d.Name = product.Name
	d.CategoryId = product.CategoryId
	d.CategoryName = product.CategoryName
	d.MerchantId = product.MerchantId
	d.StockId = product.StockId
	d.PurchasePrice = product.PurchasePrice
	d.SellingPrice = product.SellingPrice
	d.Margin = product.Margin
	d.Amount = product.Amount
	d.UnitType = product.UnitType
}

func GetProductById(ctx context.Context, id int64) (product models.Product, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	productDto := ProductDTO{}

	q := conn.ModelContext(ctx, &productDto).Where("id = ?", id)
	err = q.Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	product = productDto.fromDTO()

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

	productDto := ProductDTO{}
	productDto.toDTO(product)

	//_, err = conn.ModelContext(ctx, &catalog).WherePK().Update(catalog)

	_, err = conn.ModelContext(ctx, &product).WherePK().Column("purchase_price", "selling_price", "amount", "category_id", "name").Update()
	if err != nil {
		Loger.Debugln("error UpdateProduct", err.Error())
		return
	}

	return
}

func DeleteProductById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error DeleteProductById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Product)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func MDeleteProductByIds(ctx context.Context, ids []int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in MDeleteProductByIds", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Product)(nil)).Where("id IN (?)", pg.In(ids)).Delete()
	if err != nil {
		Loger.Debugln("error MDeleteProductByIds", err.Error())
		return
	}

	return
}

func FilterProducts(ctx context.Context, req FilterProductsRequest) ([]models.Product, error) {

	conn, err := GetPGSession()
	if err != nil {
		return nil, err
	}

	productDtos := []ProductDTO{}

	query := conn.ModelContext(ctx, &productDtos).Where("merchant_id = ?", req.MerchantId)
	if req.Barcode != "" {
		query.Where("barcode LIKE ?", req.Barcode+"%")
	}

	if req.Name != "" {
		query.Where("name LIKE ?", req.Name+"%")
	}

	if req.Size == 0 {
		req.Size = 10
	}

	err = query.Limit(req.Size).Offset(req.From).Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}

	products := make([]models.Product, 0, len(productDtos))

	for _, productDto := range productDtos {
		products = append(products, productDto.fromDTO())
	}

	return products, nil

}

func CheckBarcode(ctx context.Context, merchantId, barcode string) (bool, error) {

	conn, err := GetPGSession()
	if err != nil {
		return false, err
	}

	productDtos := []ProductDTO{}

	return conn.ModelContext(ctx, &productDtos).Where("merchant_id = ?", merchantId).Where("barcode = ?", barcode).Exists()

}

func GetCategoryById(ctx context.Context, id int64) (category models.Category, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetCategoryById", err.Error())
		return
	}

	categoryDto := models.Category{}

	q := conn.ModelContext(ctx, &categoryDto).Where("id = ?", id)
	err = q.Select()
	if err != nil && err != pg.ErrNoRows {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return categoryDto, nil
}

func InsertCategory(ctx context.Context, category models.Category) (result models.Category, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetCategoryById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &category).Returning("*", category).Insert()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return category, nil
}

func UpdateCategory(ctx context.Context, category models.Category) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetCategoryById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &category).
		WherePK().
		Column("category_name", "updated_on", "description").
		Update()

	if err != nil {
		Loger.Debugln("error UpdateCategory", err.Error())
		return
	}

	return
}

func UpdateCategoryProductsCount(ctx context.Context, count int) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetCategoryById", err.Error())
		return
	}

	category := models.Category{}

	_, err = conn.ModelContext(ctx, &category).
		WherePK().
		Set("products_count = products_count + ?", count).
		Update()

	if err != nil {
		Loger.Debugln("error UpdateCategory", err.Error())
		return
	}

	return
}

func DeleteCategoryById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error DeleteCategoryById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Category)(nil)).
		Where("id = ?", id).
		Delete()

	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func MDeleteCategoryByIds(ctx context.Context, ids []int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in MDeleteCategoryByIds", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Category)(nil)).
		Where("id IN (?)", pg.In(ids)).
		Delete()

	if err != nil {
		Loger.Debugln("error MDeleteCategoryByIds", err.Error())
		return
	}

	return
}

func GetCategoryByName(ctx context.Context, categoryName, merchantId string) (models.Category, error) {

	conn, err := GetPGSession()
	if err != nil {
		return models.Category{}, err
	}

	categoryDtos := models.Category{}

	err = conn.ModelContext(ctx, &categoryDtos).
		Where("merchant_id = ?", merchantId).
		Where("category_name = ?", categoryName).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return models.Category{}, err
	}

	return categoryDtos, nil

}

func FilterCategories(ctx context.Context, request FilterCategoryRequest) ([]models.Category, error) {

	conn, err := GetPGSession()
	if err != nil {
		return nil, err
	}

	categoryDtos := []models.Category{}

	err = conn.ModelContext(ctx, &categoryDtos).
		Where("merchant_id = ?", request.MerchantId).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}

	return categoryDtos, nil

}

func CheckCategoryExists(ctx context.Context, merchantId, categoryName string) (bool, error) {

	conn, err := GetPGSession()
	if err != nil {
		return false, err
	}

	categoryDtos := []models.Category{}

	return conn.ModelContext(ctx, &categoryDtos).Where("merchant_id = ?", merchantId).Where("category_name = ?", categoryName).Exists()

}

func UpdateProductsCount(ctx context.Context, req UpdateProductsCountRequest) error {
	conn, err := GetPGSession()
	if err != nil {
		return err
	}

	for i := range req.Nomens {
		var query *pg.Query

		product := ProductDTO{
			Barcode:    req.Nomens[i].Barcode,
			MerchantId: req.MerchantId,
			Amount:     req.Nomens[i].Amount,
		}

		switch req.Action {
		case "INC":
			query = conn.ModelContext(ctx, &product).Set("amount = amount + ?", req.Nomens[i].Amount)
		case "DEC":
			query = conn.ModelContext(ctx, &product).Set("amount = amount - ?", req.Nomens[i].Amount)
		case "UPD":
			query = conn.ModelContext(ctx, &product).Set("amount = ?", req.Nomens[i].Amount)
		}

		_, err = query.Where("merchant_id = ?", req.MerchantId).Where("barcode = ?", product.Barcode).Update()
		if err != nil {
			return err
		}
	}

	return nil
}
