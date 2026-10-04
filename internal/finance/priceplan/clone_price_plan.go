package priceplan

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ClonePricePlanUseCase mengorkestrasi penggandaan buku tarif acuan ke buku tarif baru (status DRAFT)
type ClonePricePlanUseCase interface {
	Execute(ctx context.Context, sourcePlanID string, req ClonePricePlanRequest, operator string) (*PricePlanResponse, error)
}

type clonePricePlanUseCase struct {
	repo Repository
}

func NewClonePricePlanUseCase(repo Repository) ClonePricePlanUseCase {
	return &clonePricePlanUseCase{repo: repo}
}

func (uc *clonePricePlanUseCase) Execute(ctx context.Context, sourcePlanID string, req ClonePricePlanRequest, operator string) (*PricePlanResponse, error) {
	// 1. Validasi buku tarif sumber
	sourcePlan, err := uc.repo.FindPlanByID(ctx, sourcePlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	req.NewCode = strings.TrimSpace(req.NewCode)
	req.NewName = strings.TrimSpace(req.NewName)

	if req.NewCode == "" {
		return nil, errors.New("kode buku tarif baru wajib diisi")
	}
	if req.NewName == "" {
		return nil, errors.New("nama buku tarif baru wajib diisi")
	}

	// 2. Cek keunikan kode baru
	existing, err := uc.repo.FindPlanByCode(ctx, req.NewCode)
	if err == nil && existing != nil {
		return nil, ErrPricePlanCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	from, err := time.Parse("2006-01-02", strings.TrimSpace(req.EffectiveFrom))
	if err != nil {
		return nil, errors.New("format effective_from tidak valid, gunakan YYYY-MM-DD")
	}

	var to *time.Time
	if req.EffectiveTo != nil && strings.TrimSpace(*req.EffectiveTo) != "" {
		parsedTo, err := time.Parse("2006-01-02", strings.TrimSpace(*req.EffectiveTo))
		if err != nil {
			return nil, errors.New("format effective_to tidak valid, gunakan YYYY-MM-DD")
		}
		if parsedTo.Before(from) {
			return nil, ErrEffectiveDateInvalid
		}
		to = &parsedTo
	}

	citoPercent := sourcePlan.DefaultCitoPercent
	if req.DefaultCitoPercent != nil && *req.DefaultCitoPercent > 0 {
		citoPercent = *req.DefaultCitoPercent
	}

	newPlan := &TariffPricePlan{
		Code:               req.NewCode,
		Name:               req.NewName,
		Description:        req.Description,
		EffectiveFrom:      from,
		EffectiveTo:        to,
		Status:             PricePlanStatusDraft,
		IsDefault:          req.IsDefault,
		CustomerID:         req.CustomerID,
		DefaultCitoPercent: citoPercent,
		CreatedBy:          operator,
		UpdatedBy:          operator,
	}

	if err := uc.repo.ClonePlan(ctx, sourcePlanID, newPlan); err != nil {
		return nil, err
	}

	return toPlanResponse(newPlan), nil
}
