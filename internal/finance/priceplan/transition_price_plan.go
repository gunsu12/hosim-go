package priceplan

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// SubmitPricePlanUseCase mengajukan buku tarif dari DRAFT ke SUBMITTED
type SubmitPricePlanUseCase interface {
	Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error)
}

type submitPricePlanUseCase struct {
	repo Repository
}

func NewSubmitPricePlanUseCase(repo Repository) SubmitPricePlanUseCase {
	return &submitPricePlanUseCase{repo: repo}
}

func (uc *submitPricePlanUseCase) Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if err := plan.CanTransitionTo(PricePlanStatusSubmitted); err != nil {
		return nil, err
	}

	plan.Status = PricePlanStatusSubmitted
	plan.UpdatedBy = operator

	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

// ApprovePricePlanUseCase menyetujui buku tarif (SUBMITTED -> APPROVED)
// Sesuai PRD § 4.2: Memvalidasi keseimbangan komponen biaya pada SELURUH item tarif
type ApprovePricePlanUseCase interface {
	Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error)
}

type approvePricePlanUseCase struct {
	repo Repository
}

func NewApprovePricePlanUseCase(repo Repository) ApprovePricePlanUseCase {
	return &approvePricePlanUseCase{repo: repo}
}

func (uc *approvePricePlanUseCase) Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if err := plan.CanTransitionTo(PricePlanStatusApproved); err != nil {
		return nil, err
	}

	// 1. Ambil seluruh item dan komponen dalam buku tarif
	items, err := uc.repo.FindAllItemsByPlanID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Pastikan buku tarif memiliki minimal satu item
	if len(items) == 0 {
		return nil, errors.New("buku tarif tidak memiliki item tarif dan tidak dapat disetujui")
	}

	// 3. Validasi keseimbangan komponen untuk setiap item (PRD § 4.2)
	for _, item := range items {
		if err := item.ValidateBalancing(); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	plan.Status = PricePlanStatusApproved
	plan.ApprovedAt = &now
	plan.ApprovedBy = &operator
	plan.UpdatedBy = operator

	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

// ActivatePricePlanUseCase mengaktifkan buku tarif resmi (APPROVED -> ACTIVE)
type ActivatePricePlanUseCase interface {
	Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error)
}

type activatePricePlanUseCase struct {
	repo Repository
}

func NewActivatePricePlanUseCase(repo Repository) ActivatePricePlanUseCase {
	return &activatePricePlanUseCase{repo: repo}
}

func (uc *activatePricePlanUseCase) Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if err := plan.CanTransitionTo(PricePlanStatusActive); err != nil {
		return nil, err
	}

	// Jika plan default, pastikan tidak ada konflik buku tarif default aktif lain pada periode yang sama
	if plan.IsDefault {
		conflict, err := uc.repo.CheckDefaultConflict(ctx, plan.ID, plan.EffectiveFrom, plan.EffectiveTo)
		if err != nil {
			return nil, err
		}
		if conflict {
			return nil, ErrDefaultPricePlanConflict
		}
	}

	plan.Status = PricePlanStatusActive
	plan.UpdatedBy = operator

	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

// ArchivePricePlanUseCase mengarsipkan buku tarif lampau (ACTIVE -> ARCHIVED)
type ArchivePricePlanUseCase interface {
	Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error)
}

type archivePricePlanUseCase struct {
	repo Repository
}

func NewArchivePricePlanUseCase(repo Repository) ArchivePricePlanUseCase {
	return &archivePricePlanUseCase{repo: repo}
}

func (uc *archivePricePlanUseCase) Execute(ctx context.Context, id string, operator string) (*PricePlanResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if err := plan.CanTransitionTo(PricePlanStatusArchived); err != nil {
		return nil, err
	}

	plan.Status = PricePlanStatusArchived
	plan.UpdatedBy = operator

	if err := uc.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}
