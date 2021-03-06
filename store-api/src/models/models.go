package models

import "time"

const (
	StatusDraft    = "draft"
	StatusProvided = "provided"

	MerchantStatusActive   = "active"
	MerchantStatusDisabled = "disabled"
)

type ShoppingCart struct {
	ID           int64     `json:"id"`
	MerchantId   string    `json:"merchant_id"`
	TotalSum     float64   `json:"total_sum"`
	CreatedOn    time.Time `json:"created_on"`
	ProvidedTime time.Time `json:"provided_time"`
	Status       string    `json:"status"`
}

type ShoppingCartProduct struct {
	Barcode        string    `json:"barcode"`
	Name           string    `json:"name"`
	Amount         float64   `json:"amount"`
	ShoppingCartId int64     `json:"shopping_cart_id"`
	PurchasePrice  float64   `json:"purchase_price"`
	SellingPrice   float64   `json:"selling_price"`
	Total          float64   `json:"total"`
	CreatedOn      time.Time `json:"created_on"`
}

type Merchant struct {
	MerchantId   string `pg:",pk,unique"`
	MerchantName string `json:"merchant_name"`
	IE           string `json:"ie"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	BIN          string `json:"BIN"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CreatedOn    time.Time
	UpdatedOn    time.Time
}
