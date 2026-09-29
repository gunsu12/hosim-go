package item

import (
	"context"
	"errors"
	"strings"
	"time"

	"hosim-go/pkg/enums"
)

type UpdateItemUseCase interface {
	UpdateMedication(ctx context.Context, id string, req UpdateMedicationRequest, actor string) (*ItemDetailResponse, error)
	UpdateGeneral(ctx context.Context, id string, req UpdateGeneralRequest, actor string) (*ItemDetailResponse, error)
	UpdateAsset(ctx context.Context, id string, req UpdateAssetRequest, actor string) (*ItemDetailResponse, error)
	UpdateTariff(ctx context.Context, id string, req UpdateTariffRequest, actor string) (*ItemDetailResponse, error)
	DeleteItem(ctx context.Context, id string, actor string) error
}

type updateItemUseCase struct {
	repo Repository
}

func NewUpdateItemUseCase(repo Repository) UpdateItemUseCase {
	return &updateItemUseCase{repo: repo}
}

func (uc *updateItemUseCase) UpdateMedication(ctx context.Context, id string, req UpdateMedicationRequest, actor string) (*ItemDetailResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.ItemType != enums.ItemTypeMedication {
		return nil, ErrSubtypeMismatch
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama item wajib diisi")
	}

	if code != item.Code {
		existing, err := uc.repo.FindItemByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCodeAlreadyExists
		}
	}

	category, err := uc.repo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category.ItemType != enums.ItemTypeMedication {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe MEDICATION")
	}

	if req.ProductLineID != nil && *req.ProductLineID != "" {
		plExists, err := uc.repo.ProductLineExists(ctx, *req.ProductLineID)
		if err != nil {
			return nil, err
		}
		if !plExists {
			return nil, ErrProductLineNotFound
		}
	}

	now := time.Now()
	item.Code = code
	item.Name = name
	item.GenericName = req.GenericName
	item.CategoryID = req.CategoryID
	item.ProductLineID = req.ProductLineID
	item.IsActive = req.IsActive
	item.UpdatedAt = now
	item.UpdatedBy = actor

	med := ItemMedication{
		ItemID:             id,
		KFACode:            req.KFACode,
		BPOMNIE:            req.BPOMNIE,
		DosageForm:         req.DosageForm,
		StrengthAmount:     req.StrengthAmount,
		StrengthUnit:       req.StrengthUnit,
		DefaultRoute:       req.DefaultRoute,
		MedicationType:     req.MedicationType,
		IsHighAlert:        req.IsHighAlert,
		IsLASA:             req.IsLASA,
		IsFornas:           req.IsFornas,
		IsAntibiotic:       req.IsAntibiotic,
		StorageTemperature: req.StorageTemperature,
		UpdatedAt:          now,
		UpdatedBy:          actor,
	}

	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.UpdateItem(ctx, item); err != nil {
			return err
		}
		if err := txRepo.UpdateMedication(ctx, &med); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToItemDetailResponse(updated)
	return &res, nil
}

func (uc *updateItemUseCase) UpdateGeneral(ctx context.Context, id string, req UpdateGeneralRequest, actor string) (*ItemDetailResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.ItemType != enums.ItemTypeGeneral {
		return nil, ErrSubtypeMismatch
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama item wajib diisi")
	}

	if code != item.Code {
		existing, err := uc.repo.FindItemByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCodeAlreadyExists
		}
	}

	category, err := uc.repo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category.ItemType != enums.ItemTypeGeneral {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe GENERAL")
	}

	if req.ProductLineID != nil && *req.ProductLineID != "" {
		plExists, err := uc.repo.ProductLineExists(ctx, *req.ProductLineID)
		if err != nil {
			return nil, err
		}
		if !plExists {
			return nil, ErrProductLineNotFound
		}
	}

	if !req.GeneralType.IsValid() {
		return nil, ErrInvalidGeneralType
	}
	if req.IsCSSDItem && (req.SterilizationMethod == nil || !req.SterilizationMethod.IsValid()) {
		return nil, ErrCSSDSterilizationRequired
	}

	isDisposable := true
	if req.IsDisposable != nil {
		isDisposable = *req.IsDisposable
	}

	now := time.Now()
	item.Code = code
	item.Name = name
	item.GenericName = req.GenericName
	item.CategoryID = req.CategoryID
	item.ProductLineID = req.ProductLineID
	item.IsActive = req.IsActive
	item.UpdatedAt = now
	item.UpdatedBy = actor

	gen := ItemGeneral{
		ItemID:              id,
		GeneralType:         req.GeneralType,
		IsSterile:           req.IsSterile,
		IsDisposable:        isDisposable,
		IsCSSDItem:          req.IsCSSDItem,
		SterilizationMethod: req.SterilizationMethod,
		UpdatedAt:           now,
		UpdatedBy:           actor,
	}

	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.UpdateItem(ctx, item); err != nil {
			return err
		}
		if err := txRepo.UpdateGeneral(ctx, &gen); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToItemDetailResponse(updated)
	return &res, nil
}

func (uc *updateItemUseCase) UpdateAsset(ctx context.Context, id string, req UpdateAssetRequest, actor string) (*ItemDetailResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.ItemType != enums.ItemTypeAsset {
		return nil, ErrSubtypeMismatch
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama item wajib diisi")
	}

	if code != item.Code {
		existing, err := uc.repo.FindItemByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCodeAlreadyExists
		}
	}

	category, err := uc.repo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category.ItemType != enums.ItemTypeAsset {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe ASSET")
	}

	if req.ProductLineID != nil && *req.ProductLineID != "" {
		plExists, err := uc.repo.ProductLineExists(ctx, *req.ProductLineID)
		if err != nil {
			return nil, err
		}
		if !plExists {
			return nil, ErrProductLineNotFound
		}
	}

	if req.DepreciationMethod != nil && !req.DepreciationMethod.IsValid() {
		return nil, ErrInvalidDepreciationMethod
	}

	now := time.Now()
	item.Code = code
	item.Name = name
	item.GenericName = req.GenericName
	item.CategoryID = req.CategoryID
	item.ProductLineID = req.ProductLineID
	item.IsActive = req.IsActive
	item.UpdatedAt = now
	item.UpdatedBy = actor

	asset := ItemAsset{
		ItemID:                  id,
		Brand:                   req.Brand,
		ModelName:               req.ModelName,
		IsMedicalEquipment:      req.IsMedicalEquipment,
		ExpectedLifeYears:       req.ExpectedLifeYears,
		DepreciationMethod:      req.DepreciationMethod,
		MaintenanceIntervalDays: req.MaintenanceIntervalDays,
		UpdatedAt:               now,
		UpdatedBy:               actor,
	}

	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.UpdateItem(ctx, item); err != nil {
			return err
		}
		if err := txRepo.UpdateAsset(ctx, &asset); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToItemDetailResponse(updated)
	return &res, nil
}

func (uc *updateItemUseCase) UpdateTariff(ctx context.Context, id string, req UpdateTariffRequest, actor string) (*ItemDetailResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.ItemType != enums.ItemTypeTariff {
		return nil, ErrSubtypeMismatch
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama item wajib diisi")
	}

	if code != item.Code {
		existing, err := uc.repo.FindItemByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCodeAlreadyExists
		}
	}

	category, err := uc.repo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category.ItemType != enums.ItemTypeTariff {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe TARIFF")
	}

	if !req.ChargeType.IsValid() {
		return nil, ErrInvalidChargeType
	}

	now := time.Now()
	item.Code = code
	item.Name = name
	item.CategoryID = req.CategoryID
	item.ProductLineID = nil // Non-inventory invariant
	item.IsActive = req.IsActive
	item.UpdatedAt = now
	item.UpdatedBy = actor

	tariff := ItemTariff{
		ItemID:     id,
		NameAlias:  req.NameAlias,
		ChargeType: req.ChargeType,
		Notes:      req.Notes,
		UpdatedAt:  now,
		UpdatedBy:  actor,
	}

	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.UpdateItem(ctx, item); err != nil {
			return err
		}
		if err := txRepo.UpdateTariff(ctx, &tariff); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToItemDetailResponse(updated)
	return &res, nil
}

func (uc *updateItemUseCase) DeleteItem(ctx context.Context, id string, actor string) error {
	_, err := uc.repo.FindItemByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.DeleteItem(ctx, id, actor)
}
