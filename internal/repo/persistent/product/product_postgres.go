package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/andreyxaxa/Sales-Tracker/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
)

const (
	productsTable = "products"

	idColumn         = "id"
	categoryIDColumn = "category_id"
	nameColumn       = "name"
	basePriceColumn  = "base_price"
)

type ProductsRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *ProductsRepo {
	return &ProductsRepo{pg}
}

func (r *ProductsRepo) Create(ctx context.Context, categoryID int64, name string, basePrice decimal.Decimal) (entity.Product, error) {
	sql, args, err := r.Builder.
		Insert(productsTable).
		Columns(categoryIDColumn, nameColumn, basePriceColumn).
		Values(categoryID, name, basePrice).
		Suffix("RETURNING id, category_id, name, base_price").
		ToSql()
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductsRepo - Create - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var product entity.Product

	err = executor.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.BasePrice,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return entity.Product{}, fmt.Errorf("ProductsRepo - Create: %w", errs.ErrCategoryNotFound)
		}
		return entity.Product{}, fmt.Errorf("ProductsRepo - Create - executor.QueryRow.Scan: %w", err)
	}

	return product, nil
}

func (r *ProductsRepo) Get(ctx context.Context, productID int64) (entity.Product, error) {
	sql, args, err := r.Builder.
		Select(idColumn, categoryIDColumn, nameColumn, basePriceColumn).
		From(productsTable).
		Where(squirrel.Eq{idColumn: productID}).
		ToSql()
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductsRepo - Get - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var product entity.Product

	err = executor.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.BasePrice,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, fmt.Errorf("ProductsRepo - Get: %w", errs.ErrProductNotFound)
		}
		return entity.Product{}, fmt.Errorf("ProductsRepo - Get - executor.QueryRow.Scan: %w", err)
	}

	return product, nil
}

func (r *ProductsRepo) Delete(ctx context.Context, productID int64) error {
	sql, args, err := r.Builder.
		Delete(productsTable).
		Where(squirrel.Eq{idColumn: productID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ProductsRepo - Delete - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	tag, err := executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ProductsRepo - Delete - executor.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ProductsRepo - Delete: %w", errs.ErrProductNotFound)
	}

	return nil
}

func (r *ProductsRepo) UpdateBasePrice(ctx context.Context, productID int64, newPrice decimal.Decimal) (entity.Product, error) {
	sql, args, err := r.Builder.
		Update(productsTable).
		Set(basePriceColumn, newPrice).
		Where(squirrel.Eq{idColumn: productID}).
		Suffix("RETURNING id, category_id, name, base_price").
		ToSql()
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductsRepo - UpdateBasePrice - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var updatedProduct entity.Product

	err = executor.QueryRow(ctx, sql, args...).Scan(
		&updatedProduct.ID,
		&updatedProduct.CategoryID,
		&updatedProduct.Name,
		&updatedProduct.BasePrice,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, fmt.Errorf("ProductsRepo - UpdateBasePrice: %w", errs.ErrProductNotFound)
		}
		return entity.Product{}, fmt.Errorf("ProductsRepo - UpdateBasePrice - executor.QueryRow.Scan: %w", err)
	}

	return updatedProduct, nil
}
