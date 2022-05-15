package models

import "time"

const (
	StatusDraft    = "draft"
	StatusProvided = "provided"
)

type ShortInventory struct {
	ID             int64     `json:"id"`
	DocumentNumber string    `json:"document_number"`
	MerchantId     string    `json:"merchant_id"`
	TotalSum       float64   `json:"total_sum"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	ProvidedTime   time.Time `json:"provided_time"`
	Status         string    `json:"status"`
}

type ShortWaybill struct {
	ID             int64     `json:"id"`
	DocumentNumber string    `json:"document_number"`
	MerchantId     string    `json:"merchant_id"`
	StockId        int64     `json:"stock_id"`
	TotalSum       float64   `json:"total_sum"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	ProvidedTime   time.Time `json:"provided_time"`
	Status         string    `json:"status"`
}

type InventoryProduct struct {
	Barcode       string    `json:"barcode"`
	Name          string    `json:"name"`
	ActualAmount  float64   `json:"actual_amount"`
	Amount        float64   `json:"amount"`
	InventoryId   int64     `json:"inventory_id"`
	PurchasePrice float64   `json:"purchase_price"`
	SellingPrice  float64   `json:"selling_price"`
	Total         float64   `json:"total"`
	CreatedOn     time.Time `json:"created_on"`
}

type WaybillProduct struct {
	Barcode        string    `json:"barcode"`
	Name           string    `json:"name"`
	ReceivedAmount float64   `json:"received_amount"`
	Amount         float64   `json:"amount"`
	WaybillId      int64     `json:"waybill_id"`
	PurchasePrice  float64   `json:"purchase_price"`
	SellingPrice   float64   `json:"selling_price"`
	Total          float64   `json:"total"`
	CreatedOn      time.Time `json:"created_on"`
}
