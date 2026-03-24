package category

import (
	"context"
	"fmt"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo"
)

type CategoryUseCase struct {
	repo repo.CategoryRepo
}

func New(r repo.CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{
		repo: r,
	}
}

func (uc *CategoryUseCase) Create(ctx context.Context, name string, description *string) (entity.Category, error) {
	ctg, err := uc.repo.Create(ctx, name, description)
	if err != nil {
		return entity.Category{}, fmt.Errorf("CategoryUseCase - Create - uc.repo.Create: %w", err)
	}

	return ctg, nil
}

func (uc *CategoryUseCase) GetByID(ctx context.Context, id int64) (entity.Category, error) {
	ctg, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Category{}, fmt.Errorf("CategoryUseCase - GetByID - uc.repo.GetByID: %w", err)
	}

	return ctg, nil
}

func (uc *CategoryUseCase) DeleteByID(ctx context.Context, id int64) error {
	err := uc.repo.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CategoryUseCase - DeleteByID - uc.repo.DeleteByID: %w", err)
	}

	return nil
}
