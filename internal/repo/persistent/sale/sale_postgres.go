package sale

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	salesTable = "sales"

	idColumn            = "id"
	productIDColumn     = "product_id"
	quantityColumn      = "quantity"
	unitPriceColumn     = "unit_price"
	paymentMethodColumn = "payment_method"
	soldAtColumn        = "sold_at"
	createdAtColumn     = "created_at"
)

type SalesRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *SalesRepo {
	return &SalesRepo{pg}
}

func (r *SalesRepo) Create(
	ctx context.Context,
	productID int64,
	quantity int,
	unitPrice decimal.Decimal,
	paymentMethod string,
	soldAt time.Time,
) (entity.Sale, error) {
	sql, args, err := r.Builder.
		Insert(salesTable).
		Columns(productIDColumn, quantityColumn, unitPriceColumn, paymentMethodColumn, soldAtColumn).
		Values(productID, quantity, unitPrice, paymentMethod, soldAt).
		Suffix("RETURNING id, product_id, quantity, unit_price, payment_method, sold_at, created_at").
		ToSql()
	if err != nil {
		return entity.Sale{}, fmt.Errorf("SalesRepo - Create - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var sale entity.Sale

	err = executor.QueryRow(ctx, sql, args...).Scan(
		&sale.ID,
		&sale.ProductID,
		&sale.Quantity,
		&sale.UnitPrice,
		&sale.PaymentMethod,
		&sale.SoldAt,
		&sale.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return entity.Sale{}, fmt.Errorf("SalesRepo - Create: %w", errs.ErrProductNotFound)
		}
		return entity.Sale{}, fmt.Errorf("SalesRepo - Create - executor.QueryRow.Scan: %w", err)
	}

	return sale, nil
}

func (r *SalesRepo) Get(ctx context.Context, id int64) (entity.Sale, error) {
	sql, args, err := r.Builder.
		Select(
			idColumn,
			productIDColumn,
			quantityColumn,
			unitPriceColumn,
			paymentMethodColumn,
			soldAtColumn,
			createdAtColumn,
		).
		From(salesTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return entity.Sale{}, fmt.Errorf("SalesRepo - Get - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var sale entity.Sale

	err = executor.QueryRow(ctx, sql, args...).Scan(
		&sale.ID,
		&sale.ProductID,
		&sale.Quantity,
		&sale.UnitPrice,
		&sale.PaymentMethod,
		&sale.SoldAt,
		&sale.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Sale{}, fmt.Errorf("SalesRepo - Get: %w", errs.ErrSaleNotFound)
		}
		return entity.Sale{}, fmt.Errorf("SalesRepo - Get - executor.QueryRow.Scan: %w", err)
	}

	return sale, nil
}

func (r *SalesRepo) Delete(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Delete(salesTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepo - Delete - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	tag, err := executor.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepo - Delete - executor.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SalesRepo - Delete: %w", errs.ErrSaleNotFound)
	}

	return nil
}
