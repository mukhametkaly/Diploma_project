package merchant

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
	db     *pg.DB
	logger logrus.Logger
}

func GetMerchantById(ctx context.Context, id string) (merchant models.Merchant, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetMerchantById", err.Error())
		return
	}

	q := conn.ModelContext(ctx, &merchant).Where("merchant_id = ?", id)
	err = q.Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func InsertMerchant(ctx context.Context, merchant models.Merchant) (result models.Merchant, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetMerchantById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &merchant).Returning("*", merchant).Insert()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return merchant, nil
}

func UpdateMerchant(ctx context.Context, merchant models.Merchant) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetMerchantById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &merchant).
		WherePK().
		Column("merchant_name", "ie", "address", "status", "bin", "phone", "email", "updated_on").
		Update()

	if err != nil {
		Loger.Debugln("error UpdateMerchant", err.Error())
		return
	}

	return
}

func DeleteMerchantById(ctx context.Context, id string) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error DeleteMerchantById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Merchant)(nil)).
		Where("merchant_id = ?", id).
		Delete()

	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func MDeleteMerchantByIds(ctx context.Context, ids []string) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in MDeleteMerchantByIds", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, (*models.Merchant)(nil)).
		Where("merchant_id IN (?)", pg.In(ids)).
		Delete()

	if err != nil {
		Loger.Debugln("error MDeleteMerchantByIds", err.Error())
		return
	}

	return
}

func FilterMerchants(ctx context.Context, req FilterMerchantsRequest) ([]models.Merchant, error) {

	conn, err := GetPGSession()
	if err != nil {
		return nil, err
	}

	merchants := []models.Merchant{}

	query := conn.ModelContext(ctx, &merchants)
	if req.Name != "" {
		query.Where("merchant_name LIKE ?", req.Name+"%")
	}

	if req.Size == 0 {
		req.Size = 10
	}

	err = query.Limit(req.Size).Offset(req.From).Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}

	return merchants, nil

}

func CheckId(ctx context.Context, merchantId, barcode string) (bool, error) {

	conn, err := GetPGSession()
	if err != nil {
		return false, err
	}

	merchant := []models.Merchant{}

	return conn.ModelContext(ctx, &merchant).Where("merchant_id = ?", merchantId).Where("barcode = ?", barcode).Exists()

}

func GetStatistic(merchantId string) (GetStatisticResponse, error) {
	conn, err := GetPGSession()
	if err != nil {
		return GetStatisticResponse{}, err
	}

	resp := GetStatisticResponse{}

	resp.TotalSellingCount, err = conn.
		Model((*models.ShoppingCart)(nil)).
		Where("merchant_id = ?", merchantId).
		ColumnExpr("sum(total_sum)").
		SelectAndCount(&resp.TotalSellingSum)
	if err != nil {
		return GetStatisticResponse{}, err
	}

	type Product struct{}

	err = conn.
		Model((*Product)(nil)).
		Where("merchant_id = ?", merchantId).
		ColumnExpr("sum(purchase_price)").
		Select(&resp.AllProductsPurchasePrice)
	if err != nil {
		return GetStatisticResponse{}, err
	}

	err = conn.
		Model((*Product)(nil)).
		Where("merchant_id = ?", merchantId).
		ColumnExpr("sum(selling_price)").
		Select(&resp.AllProductsSellingPrice)
	if err != nil {
		return GetStatisticResponse{}, err
	}

	type Waybill struct {
		tableName struct{} `pg:"short_waybill"`
	}
	err = conn.
		Model((*Waybill)(nil)).
		Where("merchant_id = ?", merchantId).
		ColumnExpr("sum(total_sum)").
		Select(&resp.Loss)
	if err != nil {
		return GetStatisticResponse{}, err
	}

	resp.Profit = resp.TotalSellingSum

	return resp, nil

}
