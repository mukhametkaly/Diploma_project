package models

type Product struct {
	ID            int64   `json:"id"`
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	CategoryId    int64   `json:"category_id"`
	MerchantId    int64   `json:"merchant_id"`
	StockId       int64   `json:"stock_id"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	Amount        float64 `json:"amount"`
	UnitType      string  `json:"unit_type"`
}
