package priceplan

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// GetPricePlanUseCase mengorkestrasi pencarian detail dan daftar buku tarif
type GetPricePlanUseCase interface {
	GetByID(ctx context.Context, id string) (*PricePlanResponse, error)
	List(ctx context.Context, params PlanListParams) ([]PricePlanResponse, int64, error)
}

type getPricePlanUseCase struct {
	repo Repository
}

func NewGetPricePlanUseCase(repo Repository) GetPricePlanUseCase {
	return &getPricePlanUseCase{repo: repo}
}

func (uc *getPricePlanUseCase) GetByID(ctx context.Context, id string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	return toPlanResponse(plan), nil
}

func (uc *getPricePlanUseCase) List(ctx context.Context, params PlanListParams) ([]PricePlanResponse, int64, error) {
	plans, total, err := uc.repo.FindAllPlans(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]PricePlanResponse, len(plans))
	for i := range plans {
		responses[i] = *toPlanResponse(&plans[i])
	}

	return responses, total, nil
}
