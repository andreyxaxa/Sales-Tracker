package product

import (
	"context"
	"fmt"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo"
	"github.com/shopspring/decimal"
)

type ProductUseCase struct {
	repo repo.ProductRepo
}

func New(r repo.ProductRepo) *ProductUseCase {
	return &ProductUseCase{
		repo: r,
	}
}

func (uc *ProductUseCase) Create(ctx context.Context, categoryID int64, name string, basePrice decimal.Decimal) (entity.Product, error) {
	prd, err := uc.repo.Create(ctx, categoryID, name, basePrice)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductUseCase - Create - uc.repo.Create: %w", err)
	}

	return prd, nil
}

func (uc *ProductUseCase) Get(ctx context.Context, productID int64) (entity.Product, error) {
	prd, err := uc.repo.Get(ctx, productID)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductUseCase - Get - uc.repo.Get: %w", err)
	}

	return prd, nil
}

func (uc *ProductUseCase) Delete(ctx context.Context, productID int64) error {
	err := uc.repo.Delete(ctx, productID)
	if err != nil {
		return fmt.Errorf("ProductUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}

func (uc *ProductUseCase) UpdateBasePrice(ctx context.Context, productID int64, newPrice decimal.Decimal) (entity.Product, error) {
	prd, err := uc.repo.UpdateBasePrice(ctx, productID, newPrice)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductUseCase - UpdateBasePrice - uc.repo.UpdateBasePrice: %w", err)
	}

	return prd, nil
}
