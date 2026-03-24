package sale

import (
	"context"
	"fmt"
	"time"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo"
	"github.com/shopspring/decimal"
)

type SaleUseCase struct {
	repo repo.SaleRepo
}

func New(r repo.SaleRepo) *SaleUseCase {
	return &SaleUseCase{
		repo: r,
	}
}

func (uc *SaleUseCase) Create(
	ctx context.Context,
	productID int64,
	quantity int,
	unitPrice decimal.Decimal,
	paymentMethod string,
	soldAt time.Time,
) (entity.Sale, error) {
	sale, err := uc.repo.Create(ctx, productID, quantity, unitPrice, paymentMethod, soldAt)
	if err != nil {
		return entity.Sale{}, fmt.Errorf("SaleUseCase - Create - uc.repo.Create: %w", err)
	}

	return sale, nil
}

func (uc *SaleUseCase) Get(ctx context.Context, id int64) (entity.Sale, error) {
	sale, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Sale{}, fmt.Errorf("SaleUseCase - Get - uc.repo.Get: %w", err)
	}

	return sale, nil
}

func (uc *SaleUseCase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("SaleUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
