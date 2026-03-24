package request

import "time"

type CreateSaleRequest struct {
	ProductID     int64     `json:"product_id" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required"`
	UnitPrice     string    `json:"unit_price" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	SoldAt        time.Time `json:"sold_at" validate:"required,lte"`
}
