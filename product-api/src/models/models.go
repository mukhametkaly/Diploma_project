package models

import "time"

type Product struct {
	ID            int64   `json:"id"`
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	CategoryName  string  `json:"category_name"`
	CategoryId    int64   `json:"category_id"`
	MerchantId    string  `json:"merchant_id"`
	StockId       int64   `json:"stock_id"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	Amount        float64 `json:"amount"`
	Margin        float64 `json:"margin"`
	UnitType      string  `json:"unit_type"`
}

type Category struct {
	ID            int64  `json:"id"`
	MerchantId    string `json:"merchant_id"`
	ProductsCount int64  `json:"products_count"`
	CategoryName  string `json:"category_name"`
	Description   string
	CreatedOn     time.Time
	UpdatedOn     time.Time
}
