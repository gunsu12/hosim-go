package priceplan

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CreatePricePlanUseCase mengorkestrasi pembuatan buku tarif baru (status DRAFT)
type CreatePricePlanUseCase interface {
	Execute(ctx context.Context, req CreatePricePlanRequest, operator string) (*PricePlanResponse, error)
}

type createPricePlanUseCase struct {
	repo Repository
}

func NewCreatePricePlanUseCase(repo Repository) CreatePricePlanUseCase {
	return &createPricePlanUseCase{repo: repo}
}

func (uc *createPricePlanUseCase) Execute(ctx context.Context, req CreatePricePlanRequest, operator string) (*PricePlanResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		return nil, errors.New("kode buku tarif wajib diisi")
	}
	if req.Name == "" {
		return nil, errors.New("nama buku tarif wajib diisi")
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

	// Cek keunikan kode
	existing, err := uc.repo.FindPlanByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrPricePlanCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	citoPercent := 25.00
	if req.DefaultCitoPercent != nil && *req.DefaultCitoPercent > 0 {
		citoPercent = *req.DefaultCitoPercent
	}

	plan := &TariffPricePlan{
		Code:               req.Code,
		Name:               req.Name,
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

	if err := uc.repo.CreatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

// UpdatePricePlanUseCase mengorkestrasi pembaruan metadata buku tarif (hanya jika DRAFT)
type UpdatePricePlanUseCase interface {
	Execute(ctx context.Context, id string, req UpdatePricePlanRequest, operator string) (*PricePlanResponse, error)
}

type updatePricePlanUseCase struct {
	repo Repository
}

func NewUpdatePricePlanUseCase(repo Repository) UpdatePricePlanUseCase {
	return &updatePricePlanUseCase{repo: repo}
}

func (uc *updatePricePlanUseCase) Execute(ctx context.Context, id string, req UpdatePricePlanRequest, operator string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if !plan.CanModify() {
		return nil, ErrPricePlanImmutable
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("nama buku tarif wajib diisi")
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

	plan.Name = req.Name
	plan.Description = req.Description
	plan.EffectiveFrom = from
	plan.EffectiveTo = to
	plan.IsDefault = req.IsDefault
	plan.CustomerID = req.CustomerID
	if req.DefaultCitoPercent != nil && *req.DefaultCitoPercent > 0 {
		plan.DefaultCitoPercent = *req.DefaultCitoPercent
	}
	plan.UpdatedBy = operator

	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

func toPlanResponse(p *TariffPricePlan) *PricePlanResponse {
	effFrom := p.EffectiveFrom.Format("2006-01-02")
	var effTo *string
	if p.EffectiveTo != nil {
		s := p.EffectiveTo.Format("2006-01-02")
		effTo = &s
	}

	return &PricePlanResponse{
		ID:                 p.ID,
		Code:               p.Code,
		Name:               p.Name,
		Description:        p.Description,
		EffectiveFrom:      effFrom,
		EffectiveTo:        effTo,
		Status:             p.Status,
		IsDefault:          p.IsDefault,
		CustomerID:         p.CustomerID,
		CustomerName:       p.CustomerName,
		DefaultCitoPercent: p.DefaultCitoPercent,
		ApprovedAt:         p.ApprovedAt,
		ApprovedBy:         p.ApprovedBy,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}
