package response

import "time"

type SaleResponse struct {
	ID            int64     `json:"id"`
	ProductID     int64     `json:"product_id"`
	Quantity      int       `json:"quantity"`
	UnitPrice     string    `json:"unit_price"`
	PaymentMethod string    `json:"payment_method"`
	SoldAt        time.Time `json:"sold_at"`
}

type CreateSaleResponse SaleResponse
type GetSaleResponse SaleResponse
