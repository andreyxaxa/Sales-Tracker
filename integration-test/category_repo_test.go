package integrationtest

import (
	"testing"

	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/category"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		desc := "test description"
		ctg, err := categoryRepo.Create(ctx, "Sports", &desc)

		require.NoError(t, err)
		assert.NotZero(t, ctg.ID)
		assert.Equal(t, "test description", *ctg.Description)
		assert.Equal(t, "Sports", ctg.Name)
	})

	t.Run("category already exists", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		_, err := categoryRepo.Create(ctx, "Electronics", nil)
		require.NoError(t, err)

		_, err = categoryRepo.Create(ctx, "Electronics", nil)
		require.ErrorIs(t, err, errs.ErrCategoryAlreadyExists)
	})
}

func TestGetCategory(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		desc := "some desc..."
		createdCtg, err := categoryRepo.Create(ctx, "Sports", &desc)
		require.NoError(t, err)
		assert.NotZero(t, createdCtg.ID)

		gottenCtg, err := categoryRepo.GetByID(ctx, createdCtg.ID)
		require.NoError(t, err)
		assert.Equal(t, createdCtg.ID, gottenCtg.ID)
		assert.Equal(t, createdCtg.Name, gottenCtg.Name)
		assert.Equal(t, *createdCtg.Description, *gottenCtg.Description)
	})

	t.Run("category not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		ctg, err := categoryRepo.GetByID(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrCategoryNotFound)
		assert.Zero(t, ctg)
	})
}

func TestDeleteCategory(t *testing.T) {
	pg := newTestPostgres()
	categoryRepo := category.New(pg)

	t.Run("success", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		// create
		createdCtg, err := categoryRepo.Create(ctx, "Sports", nil)
		require.NoError(t, err)
		assert.NotZero(t, createdCtg.ID)

		// delete
		err = categoryRepo.DeleteByID(ctx, createdCtg.ID)
		require.NoError(t, err)

		// get - err
		deletedCtg, err := categoryRepo.GetByID(ctx, createdCtg.ID)
		require.ErrorIs(t, err, errs.ErrCategoryNotFound)
		assert.Zero(t, deletedCtg)
	})

	t.Run("category not found", func(t *testing.T) {
		ctx, rollback := newTestContext(t)
		defer rollback()

		err := categoryRepo.DeleteByID(ctx, wrongID)
		require.ErrorIs(t, err, errs.ErrCategoryNotFound)
	})
}
