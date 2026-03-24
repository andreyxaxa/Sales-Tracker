package entity

import "github.com/shopspring/decimal"

type Product struct {
	ID         int64           `json:"id"`
	CategoryID int64           `json:"category_id"`
	Name       string          `json:"name"`
	BasePrice  decimal.Decimal `json:"base_price"`
}
