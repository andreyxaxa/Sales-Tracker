package response

import "github.com/shopspring/decimal"

type AnalyticsResponse struct {
	GroupLabel string          `json:"group_label"`
	Value      decimal.Decimal `json:"value"`
}
