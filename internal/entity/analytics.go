package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type GroupBy string

const (
	GroupByDay      GroupBy = "day"
	GroupByWeek     GroupBy = "week"
	GroupByMonth    GroupBy = "month"
	GroupByCategory GroupBy = "category"
	GroupByProduct  GroupBy = "product"
)

type SortOrder string

const (
	SortAscL  SortOrder = "asc"
	SortDescL SortOrder = "desc"
	SortAscU  SortOrder = "ASC"
	SortDescU SortOrder = "DESC"
)

type AnalyticsParams struct {
	From          time.Time `json:"from"`
	To            time.Time `json:"to"`
	ProductID     *int64    `json:"product_id"`
	CategoryID    *int64    `json:"category_id"`
	GroupBy       *GroupBy  `json:"group_by"`
	SortBy        *string   `json:"sort_by"`
	SortOrder     SortOrder `json:"sord_order"`
	PaymentMethod *string   `json:"payment_method"`
}

type AnalyticsResult struct {
	GroupLabel string          `db:"group_label"`
	Value      decimal.Decimal `db:"value"`
}
