package integrationtest

import (
	"testing"
	"time"

	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/category"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/product"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/sale"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSale(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)
	salesRepo := sale.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		// category
		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		// product
		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		prd, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, prd.ID)

		// sale
		quantity := 5
		paymentMethod := "online"
		now := time.Now()
		sale, err := salesRepo.Create(ctx, prd.ID, quantity, prd.BasePrice, paymentMethod, now)
		require.NoError(t, err)
		assert.NotZero(t, sale.ID)
		assert.NotZero(t, sale.CreatedAt)
		assert.NotZero(t, sale.SoldAt)
		assert.Equal(t, prd.ID, sale.ProductID)
		assert.Equal(t, quantity, sale.Quantity)
		assert.Equal(t, paymentMethod, sale.PaymentMethod)
		// TODO: expected: time.Local ||| actual: time.UTC
		//assert.Equal(t, now, sale.SoldAt)
		assert.True(t, prd.BasePrice.Equal(sale.UnitPrice))
	})

	t.Run("product not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		quantity := 5
		paymentMethod := "online"
		price, err := decimal.NewFromString("150.150")
		require.NoError(t, err)

		sale, err := salesRepo.Create(ctx, wrongID, quantity, price, paymentMethod, time.Now())
		require.ErrorIs(t, err, errs.ErrProductNotFound)
		assert.Zero(t, sale)
	})
}

func TestGetSale(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)
	salesRepo := sale.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		// category
		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		// product
		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		prd, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, prd.ID)

		// sale
		quantity := 5
		paymentMethod := "online"
		createdSale, err := salesRepo.Create(ctx, prd.ID, quantity, prd.BasePrice, paymentMethod, time.Now())
		require.NoError(t, err)
		assert.NotZero(t, createdSale.ID)

		// get
		gottenSale, err := salesRepo.Get(ctx, createdSale.ID)
		require.NoError(t, err)
		assert.Equal(t, createdSale.ID, gottenSale.ID)
		assert.Equal(t, createdSale.CreatedAt, gottenSale.CreatedAt)
		assert.Equal(t, createdSale.PaymentMethod, gottenSale.PaymentMethod)
		assert.Equal(t, createdSale.ProductID, gottenSale.ProductID)
		assert.Equal(t, createdSale.Quantity, gottenSale.Quantity)
		assert.Equal(t, createdSale.SoldAt, gottenSale.SoldAt)
		assert.True(t, createdSale.UnitPrice.Equal(gottenSale.UnitPrice))
	})

	t.Run("sale not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		sale, err := salesRepo.Get(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrSaleNotFound)
		assert.Zero(t, sale)
	})
}

func TestDeleteSale(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)
	salesRepo := sale.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		// category
		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		// product
		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		prd, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, prd.ID)

		// sale
		quantity := 5
		paymentMethod := "online"
		createdSale, err := salesRepo.Create(ctx, prd.ID, quantity, prd.BasePrice, paymentMethod, time.Now())
		require.NoError(t, err)
		assert.NotZero(t, createdSale.ID)

		// delete
		err = salesRepo.Delete(ctx, createdSale.ID)
		require.NoError(t, err)

		// get - err
		gottenSale, err := salesRepo.Get(ctx, createdSale.ID)
		require.ErrorIs(t, err, errs.ErrSaleNotFound)
		assert.Zero(t, gottenSale)
	})

	t.Run("sale not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		err := salesRepo.Delete(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrSaleNotFound)
	})
}
