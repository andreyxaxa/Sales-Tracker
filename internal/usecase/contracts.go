package usecase

import (
	"context"
	"time"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/shopspring/decimal"
)

type (
	Category interface {
		Create(ctx context.Context, name string, description *string) (entity.Category, error)
		GetByID(ctx context.Context, id int64) (entity.Category, error)
		DeleteByID(ctx context.Context, id int64) error
	}

	Product interface {
		Create(ctx context.Context, categoryID int64, name string, basePrice decimal.Decimal) (entity.Product, error)
		Get(ctx context.Context, productID int64) (entity.Product, error)
		Delete(ctx context.Context, productID int64) error
		UpdateBasePrice(ctx context.Context, productID int64, newPrice decimal.Decimal) (entity.Product, error)
	}

	Sale interface {
		Create(
			ctx context.Context,
			productID int64,
			quantity int,
			unitPrice decimal.Decimal,
			paymentMethod string,
			soldAt time.Time,
		) (entity.Sale, error)
		Get(ctx context.Context, id int64) (entity.Sale, error)
		Delete(ctx context.Context, id int64) error
	}

	Analytics interface {
		GetRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error)
		GetAvgRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error)
		GetSalesCount(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error)
		GetMedian(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error)
		GetPercentile90(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error)
	}
)
