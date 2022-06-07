package inventory

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/document-api/src/config"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
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

type ShortInventoryDTO struct {
	tableName      struct{} `pg:"short_inventory"`
	Id             int64    `pg:",pk,unique"`
	DocumentNumber string
	MerchantId     string
	TotalSum       float64
	CreatedOn      time.Time
	UpdatedOn      time.Time
	ProvidedTime   time.Time
	Employee       string `json:"employee"`
	Status         string
}

func (d *ShortInventoryDTO) fromDTO() models.ShortInventory {
	var inventory models.ShortInventory
	inventory.ID = d.Id
	inventory.DocumentNumber = d.DocumentNumber
	inventory.MerchantId = d.MerchantId
	inventory.TotalSum = d.TotalSum
	inventory.CreatedOn = d.CreatedOn
	inventory.UpdatedOn = d.UpdatedOn
	inventory.ProvidedTime = d.ProvidedTime
	inventory.Employee = d.Employee
	inventory.Status = d.Status
	return inventory
}

func (d *ShortInventoryDTO) toDTO(inventory models.ShortInventory) {
	d.Id = inventory.ID
	d.DocumentNumber = inventory.DocumentNumber
	d.MerchantId = inventory.MerchantId
	d.TotalSum = inventory.TotalSum
	d.CreatedOn = inventory.CreatedOn
	d.UpdatedOn = inventory.UpdatedOn
	d.ProvidedTime = inventory.ProvidedTime
	d.Employee = inventory.Employee
	d.Status = inventory.Status
}

type InventoryProductsDTO struct {
	tableName     struct{} `pg:"inventory_product"`
	Barcode       string
	Name          string    `json:"name"`
	ActualAmount  float64   `json:"actual_amount"`
	Amount        float64   `json:"amount"`
	InventoryId   int64     `json:"inventory_id"`
	PurchasePrice float64   `json:"purchase_price"`
	SellingPrice  float64   `json:"selling_price"`
	Total         float64   `json:"total"`
	CreatedOn     time.Time `json:"created_on"`
}

func (d *InventoryProductsDTO) fromDTO() models.InventoryProduct {
	var product models.InventoryProduct
	product.Barcode = d.Barcode
	product.Name = d.Name
	product.ActualAmount = d.ActualAmount
	product.Amount = d.Amount
	product.InventoryId = d.InventoryId
	product.PurchasePrice = d.PurchasePrice
	product.SellingPrice = d.SellingPrice
	product.Total = d.Total
	product.CreatedOn = d.CreatedOn
	return product
}

func (d *InventoryProductsDTO) toDTO(product models.InventoryProduct) {
	d.Barcode = product.Barcode
	d.Name = product.Name
	d.Amount = product.Amount
	d.ActualAmount = product.ActualAmount
	d.InventoryId = product.InventoryId
	d.PurchasePrice = product.PurchasePrice
	d.SellingPrice = product.SellingPrice
	d.Total = product.Total
	d.CreatedOn = product.CreatedOn
}

func InsertInventory(ctx context.Context, inventory models.ShortInventory) (models.ShortInventory, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertInventory", err.Error())
		return models.ShortInventory{}, err
	}

	var dtoInventory ShortInventoryDTO
	dtoInventory.toDTO(inventory)

	_, err = conn.ModelContext(ctx, &dtoInventory).Returning("*", &dtoInventory).Insert(&dtoInventory)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return models.ShortInventory{}, err
	}

	inventory = dtoInventory.fromDTO()

	return inventory, nil
}

func UpdateInventoryStatus(ctx context.Context, inventory models.ShortInventory) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	var dtoInventory ShortInventoryDTO
	dtoInventory.toDTO(inventory)

	_, err = conn.ModelContext(ctx, &dtoInventory).WherePK().Column("updated_on", "provided_time", "status").Update()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func UpdateInventorySum(ctx context.Context, id int64, sum float64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	var dtoInventory ShortInventoryDTO
	dtoInventory.UpdatedOn = time.Now()
	dtoInventory.Id = id

	_, err = conn.ModelContext(ctx, &dtoInventory).WherePK().Set("total_sum = total_sum + ?", sum).Column("updated_on").Update()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func DeleteInventoryById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteInventoryById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*ShortInventoryDTO)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func GetInventory(ctx context.Context, req InventorysFilterRequest) ([]models.ShortInventory, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetInventory", err.Error())
		return nil, err
	}

	inventoryDto := []ShortInventoryDTO{}

	query := conn.ModelContext(ctx, &inventoryDto).Where("merchant_id = ?", req.MerchantId)

	if req.Status != "" {
		query.Where("status = ?", req.Status)
	}

	if req.DocumentNumber != "" {
		query.Where("document_number = ?", req.DocumentNumber)
	}

	if req.Size == 0 {
		req.Size = 10
	}

	err = query.Limit(req.Size).Offset(req.From).Select()
	if err != nil {
		Loger.Debugln("error select in get list inventories", err.Error())
		return nil, err
	}

	inventories := make([]models.ShortInventory, 0, len(inventoryDto))

	for _, dto := range inventoryDto {
		inventory := dto.fromDTO()
		inventories = append(inventories, inventory)
	}

	return inventories, err
}

func IfDocNumberExist(ctx context.Context, merchantId, docNum string) (bool, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetInventory", err.Error())
		return false, err
	}

	inventoryDto := []ShortInventoryDTO{}

	return conn.ModelContext(ctx, &inventoryDto).
		Where("merchant_id = ?", merchantId).
		Where("document_number = ?", docNum).
		Exists()

}

func GetInventoryById(ctx context.Context, id int64) (models.ShortInventory, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetInventory", err.Error())
		return models.ShortInventory{}, err
	}

	inventoryDto := ShortInventoryDTO{}
	inventoryDto.Id = id

	err = conn.ModelContext(ctx, &inventoryDto).WherePK().Select()
	if err != nil {
		Loger.Debugln("error select in get list inventories", err.Error())
		return models.ShortInventory{}, err
	}

	inventory := inventoryDto.fromDTO()

	return inventory, err
}

func GetInventoryProducts(ctx context.Context, req GetInventoryProductsRequest) (product []models.InventoryProduct, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetInventoryProducts", err.Error())
		return
	}

	DtoProducts := []InventoryProductsDTO{}

	query := conn.ModelContext(ctx, &DtoProducts).Where("inventory_id = ?", req.InventoryId)
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

func DeleteInventoryProduct(ctx context.Context, req DeleteInventoryProductRequest) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteInventoryProduct", err.Error())
		return
	}

	_, err = conn.Model((*InventoryProductsDTO)(nil)).
		Where("inventory_id = ?", req.InventoryId).
		Where("barcode = ?", req.Barcode).
		Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func InsertInventoryProduct(ctx context.Context, product models.InventoryProduct) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertInventoryProduct", err.Error())
		return
	}

	var dtoProduct InventoryProductsDTO
	dtoProduct.toDTO(product)

	_, err = conn.ModelContext(ctx, &dtoProduct).Insert(&dtoProduct)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func UpdateInventoryProduct(ctx context.Context, product models.InventoryProduct) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertInventoryProduct", err.Error())
		return
	}

	var dtoProduct InventoryProductsDTO
	dtoProduct.toDTO(product)

	_, err = conn.ModelContext(ctx, &dtoProduct).
		Where("barcode = ?", dtoProduct.Barcode).
		Where("inventory_id = ?", dtoProduct.InventoryId).
		Update("received_amount", "amount", "purchase_price", "selling_price", "total")
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}
