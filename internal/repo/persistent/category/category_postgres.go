package category

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
)

const (
	categoriesTable = "categories"

	idColumn          = "id"
	nameColumn        = "name"
	descriptionColumn = "description"
)

type CategoryRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *CategoryRepo {
	return &CategoryRepo{pg}
}

func (r *CategoryRepo) Create(ctx context.Context, name string, description *string) (entity.Category, error) {
	sql, args, err := r.Builder.
		Insert(categoriesTable).
		Columns(nameColumn, descriptionColumn).
		Values(name, description).
		Suffix("RETURNING id, name, description").
		ToSql()
	if err != nil {
		return entity.Category{}, fmt.Errorf("CategoryRepo - Create - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var category entity.Category

	err = executor.QueryRow(ctx, sql, args...).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		// если категория с таким именем уже существует
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entity.Category{}, fmt.Errorf("CategoryRepo - Create: %w", errs.ErrCategoryAlreadyExists)
		}
		return entity.Category{}, fmt.Errorf("CategoryRepo - Create - executor.QueryRow.Scan: %w", err)
	}

	return category, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (entity.Category, error) {
	sql, args, err := r.Builder.
		Select(idColumn, nameColumn, descriptionColumn).
		From(categoriesTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return entity.Category{}, fmt.Errorf("CategoryRepo - GetByID - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	var category entity.Category

	err = executor.QueryRow(ctx, sql, args...).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Category{}, fmt.Errorf("CategoryRepo - GetByID: %w", errs.ErrCategoryNotFound)
		}
		return entity.Category{}, fmt.Errorf("CategoryRepo - GetByID - executor.QueryRow.Scan: %w", err)
	}

	return category, nil
}

func (r *CategoryRepo) DeleteByID(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Delete(categoriesTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("CategoryRepo - DeleteByID - r.Builder.ToSql: %w", err)
	}

	executor := r.GetExecutor(ctx)

	tag, err := executor.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		// если в products есть продукты из этой категории
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return fmt.Errorf("CategoryRepo - DeleteByID: %w", errs.ErrCategoryHasProducts)
		}
		return fmt.Errorf("CategoryRepo - DeleteByID - executor.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("CategoryRepo - DeleteByID: %w", errs.ErrCategoryNotFound)
	}

	return nil
}
