package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Sale struct {
	ID            int64           `json:"id"`
	ProductID     int64           `json:"product_id"`
	Quantity      int             `json:"quantity"`
	UnitPrice     decimal.Decimal `json:"unit_price"`
	PaymentMethod string          `json:"payment_method"`
	SoldAt        time.Time       `json:"sold_at"`
	CreatedAt     time.Time       `json:"created_at"`
}
