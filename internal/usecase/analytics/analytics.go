package analytics

import (
	"context"
	"fmt"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo"
)

type AnalyticsUseCase struct {
	repo repo.AnalyticsRepo
}

func New(r repo.AnalyticsRepo) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		repo: r,
	}
}

func (uc *AnalyticsUseCase) GetRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	results, err := uc.repo.GetRevenue(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsUseCase - GetRevenue - uc.repo.GetRevenue: %w", err)
	}

	return results, nil
}

func (uc *AnalyticsUseCase) GetAvgRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	results, err := uc.repo.GetAvgRevenue(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsUseCase - GetAvgRevenue - uc.repo.GetAvgRevenue: %w", err)
	}

	return results, nil
}

func (uc *AnalyticsUseCase) GetSalesCount(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	results, err := uc.repo.GetSalesCount(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsUseCase - GetSalesCount - uc.repo.GetSalesCount: %w", err)
	}

	return results, nil
}

func (uc *AnalyticsUseCase) GetMedian(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	results, err := uc.repo.GetMedian(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsUseCase - GetMedian - uc.repo.GetMedian: %w", err)
	}

	return results, nil
}

func (uc *AnalyticsUseCase) GetPercentile90(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	results, err := uc.repo.GetPercentile90(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsUseCase - GetPercentile90 - uc.repo.GetPercentile90: %w", err)
	}

	return results, nil
}
