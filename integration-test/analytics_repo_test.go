package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/analytics"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/category"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/product"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/sale"
	pg "github.com/andreyxaxa/Sales-Tracker/pkg/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type salesSeed struct {
	productID     int64
	quantity      int
	unitPrice     decimal.Decimal
	paymentMethod string
	soldAt        time.Time
}

func seedSales(t *testing.T, ctx context.Context, p *pg.Postgres, seeds []salesSeed) {
	t.Helper()

	salesRepo := sale.New(p)

	for _, s := range seeds {
		sale, err := salesRepo.Create(ctx, s.productID, s.quantity, s.unitPrice, s.paymentMethod, s.soldAt)
		require.NoError(t, err)
		assert.NotZero(t, sale.ID)
	}
}

func fillTables(t *testing.T, ctx context.Context, pg *pg.Postgres) {
	t.Helper()

	categoryRepo := category.New(pg)
	productRepo := product.New(pg)

	// Categories table
	// 1. Sports
	sportsCtg, err := categoryRepo.Create(ctx, "Sports", nil)
	require.NoError(t, err)
	assert.NotZero(t, sportsCtg.ID)
	// 2. Furniture
	furnitureCtg, err := categoryRepo.Create(ctx, "Furniture", nil)
	require.NoError(t, err)
	assert.NotZero(t, furnitureCtg.ID)
	// 3. Electronics
	electronicsCtg, err := categoryRepo.Create(ctx, "Electronics", nil)
	require.NoError(t, err)
	assert.NotZero(t, electronicsCtg.ID)

	// Products table
	// 1. Sports
	ballPrice, _ := decimal.NewFromString("100")
	ballPrd, err := productRepo.Create(ctx, sportsCtg.ID, "Ball", ballPrice)
	require.NoError(t, err)
	assert.NotZero(t, ballPrd.ID)

	glovesPrice := decimal.NewFromInt(200)
	glovePrd, err := productRepo.Create(ctx, sportsCtg.ID, "Boxing Gloves", glovesPrice)
	require.NoError(t, err)
	assert.NotZero(t, glovePrd.ID)
	// 2. Furniture
	sofaPrice, _ := decimal.NewFromString("500.49")
	sofaPrd, err := productRepo.Create(ctx, furnitureCtg.ID, "Sofa", sofaPrice)
	require.NoError(t, err)
	assert.NotZero(t, sofaPrd.ID)

	closetPrice, _ := decimal.NewFromString("399")
	closetPrd, err := productRepo.Create(ctx, furnitureCtg.ID, "Closet", closetPrice)
	require.NoError(t, err)
	assert.NotZero(t, closetPrd.ID)
	// 3. Electronics
	phonePrice, _ := decimal.NewFromString("249")
	phonePrd, err := productRepo.Create(ctx, electronicsCtg.ID, "Phone", phonePrice)
	require.NoError(t, err)
	assert.NotZero(t, phonePrd.ID)

	laptopPrice, _ := decimal.NewFromString("1230")
	laptopPrd, err := productRepo.Create(ctx, electronicsCtg.ID, "Laptop", laptopPrice)
	require.NoError(t, err)
	assert.NotZero(t, laptopPrd.ID)

	// Sales table (helper)
	seedSales(t, ctx, pg, []salesSeed{
		{ballPrd.ID, 2, decimal.NewFromInt(100), "cash", time.Now()},
		{ballPrd.ID, 3, decimal.NewFromInt(200), "card", time.Now()},
		{glovePrd.ID, 1, glovesPrice, "cash", time.Now()},
		{glovePrd.ID, 4, decimal.NewFromInt(139), "card", time.Now()},
		{sofaPrd.ID, 1, sofaPrice, "card", time.Now()},
		{sofaPrd.ID, 2, decimal.NewFromInt(600), "cash", time.Now()},
		{closetPrd.ID, 3, closetPrice, "card", time.Now()},
		{closetPrd.ID, 1, decimal.NewFromInt(359), "cash", time.Now()},
		{phonePrd.ID, 11, phonePrice, "cash", time.Now()},
		{phonePrd.ID, 1, decimal.NewFromInt(200), "card", time.Now()},
		{laptopPrd.ID, 19, decimal.NewFromInt(1100), "card", time.Now()},
		{laptopPrd.ID, 1, laptopPrice, "cash", time.Now()},
	}) // 2*100 + 3*200 + 1*200 + 4*139 + 1*500.49 + 2*600 + 3*399 + 1*359 + 11*249 + 1*200 + 19*1100 + 1*1230
}

func TestRevenue(t *testing.T) {
	pg := newTestPostgres()
	analyticsRepo := analytics.New(pg)

	expectedSum := decimal.NewFromInt(2*100 + 3*200 + 1*200 + 4*139 + 2*600 + 3*399 + 1*359 + 11*249 + 1*200 + 19*1100 + 1*1230).
		Add(decimal.NewFromFloat(1 * 500.49))

	params := entity.AnalyticsParams{
		From: time.Now().AddDate(0, 0, -1),
		To:   time.Now().AddDate(0, 0, 1),
	}

	t.Run("sum success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// analytics
		result, err := analyticsRepo.GetRevenue(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 1)

		assert.True(t, expectedSum.Equal(result[0].Value))
	})

	t.Run("sum with params", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		p := entity.AnalyticsParams{
			From: params.From,
			To:   params.To,
		}

		gb := entity.GroupBy("product")
		p.GroupBy = &gb

		sb := "value"
		p.SortBy = &sb

		so := entity.SortOrder("asc")
		p.SortOrder = so

		// analytics
		result, err := analyticsRepo.GetRevenue(ctx, p)
		require.NoError(t, err)
		require.Len(t, result, 6)
	})

	t.Run("avg success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// analytics
		result, err := analyticsRepo.GetAvgRevenue(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 1)

		expectedAvg := expectedSum.Div(decimal.NewFromInt(12))

		assert.True(t, expectedAvg.Equal(result[0].Value))
	})

	t.Run("avg with params", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		p := entity.AnalyticsParams{
			From: params.From,
			To:   params.To,
		}

		gb := entity.GroupBy("category")
		p.GroupBy = &gb

		sb := "value"
		p.SortBy = &sb

		so := entity.SortOrder("desc")
		p.SortOrder = so

		// analytics
		result, err := analyticsRepo.GetAvgRevenue(ctx, p)
		require.NoError(t, err)
		require.Len(t, result, 3)
	})
}

func TestSalesCount(t *testing.T) {
	pg := newTestPostgres()
	analyticsRepo := analytics.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		params := entity.AnalyticsParams{
			From: time.Now().AddDate(0, 0, -1),
			To:   time.Now().AddDate(0, 0, 1),
		}

		// analytics
		result, err := analyticsRepo.GetSalesCount(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 1)

		expectedCount := decimal.NewFromInt(12)

		assert.True(t, expectedCount.Equal(result[0].Value))
	})

	t.Run("sales count with params", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		params := entity.AnalyticsParams{
			From: time.Now().AddDate(0, 0, -1),
			To:   time.Now().AddDate(0, 0, 1),
		}

		gb := entity.GroupBy("product")
		params.GroupBy = &gb

		sb := "value"
		params.SortBy = &sb

		so := entity.SortOrder("desc")
		params.SortOrder = so

		// analytics
		result, err := analyticsRepo.GetSalesCount(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 6)
		assert.True(t, decimal.NewFromInt(2).Equal(result[0].Value))
	})
}

func TestGetMedian(t *testing.T) {
	pg := newTestPostgres()
	analyticsRepo := analytics.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		params := entity.AnalyticsParams{
			From: time.Now().AddDate(0, 0, -1),
			To:   time.Now().AddDate(0, 0, 1),
		}

		// analytics
		result, err := analyticsRepo.GetMedian(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 1)

		expectedMedian, _ := decimal.NewFromString("578")
		assert.True(t, expectedMedian.Equal(result[0].Value))
	})
}

func TestGetP90(t *testing.T) {
	pg := newTestPostgres()
	analyticsRepo := analytics.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		fillTables(t, ctx, pg)

		// params
		params := entity.AnalyticsParams{
			From: time.Now().AddDate(0, 0, -1),
			To:   time.Now().AddDate(0, 0, 1),
		}

		// analytics
		result, err := analyticsRepo.GetPercentile90(ctx, params)
		require.NoError(t, err)
		require.Len(t, result, 1)

		expectedP90, _ := decimal.NewFromString("2588.1")
		assert.True(t, expectedP90.Equal(result[0].Value.Round(1)))
	})
}
