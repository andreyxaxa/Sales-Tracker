package integrationtest

import (
	"testing"

	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/category"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/product"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		prd, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, prd.ID)
		assert.Equal(t, "Football Ball", prd.Name)
		assert.Equal(t, ctg.ID, prd.CategoryID)
		assert.True(t, price.Equal(prd.BasePrice))
	})

	t.Run("category not found ", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		price, err := decimal.NewFromString("3.29")
		require.NoError(t, err)

		prd, err := productRepo.Create(ctx, wrongID, "Book", price)
		require.ErrorIs(t, err, errs.ErrCategoryNotFound)
		assert.Zero(t, prd)
	})
}

func TestGetProduct(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		createdPrd, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, createdPrd.ID)

		gottenPrd, err := productRepo.Get(ctx, createdPrd.ID)
		require.NoError(t, err)
		assert.Equal(t, createdPrd.ID, gottenPrd.ID)
		assert.Equal(t, createdPrd.Name, gottenPrd.Name)
		assert.Equal(t, createdPrd.CategoryID, gottenPrd.CategoryID)
		assert.True(t, createdPrd.BasePrice.Equal(gottenPrd.BasePrice))
	})

	t.Run("product not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		prd, err := productRepo.Get(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrProductNotFound)
		assert.Zero(t, prd.ID)
	})
}

func TestDeleteProduct(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		createdProduct, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, createdProduct.ID)

		err = productRepo.Delete(ctx, createdProduct.ID)
		require.NoError(t, err)

		deletedProduct, err := productRepo.Get(ctx, createdProduct.ID)
		require.ErrorIs(t, err, errs.ErrProductNotFound)
		assert.Zero(t, deletedProduct)
	})

	t.Run("product not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		err := productRepo.Delete(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrProductNotFound)
	})
}

func TestUpdateBasePrice(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)
	productRepo := product.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		// create
		ctg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)

		price, err := decimal.NewFromString("249")
		require.NoError(t, err)

		createdProduct, err := productRepo.Create(ctx, ctg.ID, "Football Ball", price)
		require.NoError(t, err)
		assert.NotZero(t, createdProduct.ID)

		// update
		newPrice, err := decimal.NewFromString("159.5")
		require.NoError(t, err)

		updatedProduct, err := productRepo.UpdateBasePrice(ctx, createdProduct.ID, newPrice)
		require.NoError(t, err)
		assert.NotZero(t, updatedProduct.ID)
		assert.True(t, newPrice.Equal(updatedProduct.BasePrice))
	})

	t.Run("product not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		newPrice, err := decimal.NewFromString("111.111")
		require.NoError(t, err)

		prd, err := productRepo.UpdateBasePrice(ctx, wrongID, newPrice)
		require.ErrorIs(t, err, errs.ErrProductNotFound)
		assert.Zero(t, prd)
	})
}
