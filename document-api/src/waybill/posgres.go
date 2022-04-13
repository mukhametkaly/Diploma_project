package waybill

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

type ShortWaybillDTO struct {
	tableName      struct{} `pg:"short_waybill"`
	Id             int64    `pg:",pk,unique"`
	DocumentNumber string
	MerchantId     int64
	StockId        int64
	TotalSum       float64
	CreatedOn      time.Time
	UpdatedOn      time.Time
	ProvidedTime   time.Time
	Status         string
}

func (d *ShortWaybillDTO) fromDTO() (models.ShortWaybill, error) {
	var waybill models.ShortWaybill
	waybill.ID = d.Id
	waybill.DocumentNumber = d.DocumentNumber
	waybill.MerchantId = d.MerchantId
	waybill.StockId = d.StockId
	waybill.TotalSum = d.TotalSum
	waybill.CreatedOn = d.CreatedOn
	waybill.UpdatedOn = d.UpdatedOn
	waybill.ProvidedTime = d.ProvidedTime
	waybill.Status = d.Status
	return waybill, nil
}

func (d *ShortWaybillDTO) toDTO(waybill models.ShortWaybill) {
	d.Id = waybill.ID
	d.DocumentNumber = waybill.DocumentNumber
	d.MerchantId = waybill.MerchantId
	d.StockId = waybill.StockId
	d.TotalSum = waybill.TotalSum
	d.CreatedOn = waybill.CreatedOn
	d.UpdatedOn = waybill.UpdatedOn
	d.ProvidedTime = waybill.ProvidedTime
	d.Status = waybill.Status
}

type WaybillProductsDTO struct {
	tableName     struct{} `pg:"short_waybill"`
	Barcode       string
	Name          string
	AmountBefore  float64
	AmountAfter   float64
	WaybillId     int64
	PurchasePrice float64
	SellingPrice  float64
	Total         float64
	CreatedOn     time.Time
}

func (d *WaybillProductsDTO) fromDTO() models.WaybillProduct {
	var product models.WaybillProduct
	product.Barcode = d.Barcode
	product.Name = d.Name
	product.AmountBefore = d.AmountBefore
	product.AmountAfter = d.AmountAfter
	product.WaybillId = d.WaybillId
	product.PurchasePrice = d.PurchasePrice
	product.SellingPrice = d.SellingPrice
	product.Total = d.Total
	product.CreatedOn = d.CreatedOn
	return product
}

func (d *WaybillProductsDTO) toDTO(product models.WaybillProduct) {
	d.Barcode = product.Barcode
	d.Name = product.Name
	d.AmountBefore = product.AmountBefore
	d.AmountAfter = product.AmountAfter
	d.WaybillId = product.WaybillId
	d.PurchasePrice = product.PurchasePrice
	d.SellingPrice = product.SellingPrice
	d.Total = product.Total
	d.CreatedOn = product.CreatedOn
}

func InsertWaybill(ctx context.Context, waybill models.ShortWaybill) (wazybil models.ShortWaybill, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertWaybill", err.Error())
		return
	}

	var dtoWaybill ShortWaybillDTO
	dtoWaybill.toDTO(waybill)

	_, err = conn.ModelContext(ctx, &dtoWaybill).Insert(&dtoWaybill)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func UpdateWaybill(ctx context.Context, waybill models.ShortWaybill) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	var dtoWaybill ShortWaybillDTO
	dtoWaybill.toDTO(waybill)

	_, err = conn.ModelContext(ctx, &waybill).WherePK().Update(waybill)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func DeleteWaybillById(ctx context.Context, id int64) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteWaybillById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*ShortWaybillDTO)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func GetWaybillProducts(ctx context.Context, req GetWaybillProductsRequest) (product []models.WaybillProduct, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetWaybillProducts", err.Error())
		return
	}

	DtoProducts := []WaybillProductsDTO{}

	q := conn.ModelContext(ctx, &DtoProducts).Where("waybill_id = ?", req.WaybillId)
	err = q.Order("created_on ASC").Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func DeleteWaybillProductById(ctx context.Context, req DeleteWaybillProductRequest) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in DeleteWaybillProductById", err.Error())
		return
	}

	_, err = conn.Model((*WaybillProductsDTO)(nil)).
		Where("waybill_id = ?", req.WaybillId).
		Where("barcode = ?", req.Barcode).
		Delete()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func InsertWaybillProduct(ctx context.Context, product models.WaybillProduct) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in InsertWaybillProduct", err.Error())
		return
	}

	var dtoProduct WaybillProductsDTO
	dtoProduct.toDTO(product)

	_, err = conn.ModelContext(ctx, &dtoProduct).Insert(&dtoProduct)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func GetWaybills(ctx context.Context, req WaybillsFilterRequest) ([]models.WaybillProduct, error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetWaybill", err.Error())
		return nil, err
	}

	dtoProducts := []WaybillProductsDTO{}

	query := conn.ModelContext(ctx, &dtoProducts).Where("merchant_id = ?", req.MerchantId)

	if req.Status != "" {
		query.Where("status = ?", req.Status)
	}

	err = query.Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return nil, err
	}

	products := make([]models.WaybillProduct, 0, len(dtoProducts))

	for _, dto := range dtoProducts {
		product := dto.fromDTO()
		products = append(products, product)
	}

	return products, err
}