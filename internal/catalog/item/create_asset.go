package item

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"

	"hosim-go/pkg/enums"
)

var (
	ErrInvalidDepreciationMethod = errors.New("metode penyusutan aktiva tidak valid")
)

type CreateAssetUseCase interface {
	Execute(ctx context.Context, req CreateAssetRequest, actor string) (*ItemDetailResponse, error)
}

type createAssetUseCase struct {
	repo Repository
}

func NewCreateAssetUseCase(repo Repository) CreateAssetUseCase {
	return &createAssetUseCase{repo: repo}
}

func (uc *createAssetUseCase) Execute(ctx context.Context, req CreateAssetRequest, actor string) (*ItemDetailResponse, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	baseUnitName := strings.TrimSpace(req.BaseUnitName)

	if code == "" {
		return nil, errors.New("kode item wajib diisi")
	}
	if name == "" {
		return nil, errors.New("nama item wajib diisi")
	}
	if baseUnitName == "" {
		return nil, ErrBaseUnitRequired
	}
	if req.CategoryID == "" {
		return nil, errors.New("kategori item wajib dipilih")
	}

	// Cek keunikan kode
	existing, err := uc.repo.FindItemByCode(ctx, code)
	if err == nil && existing != nil {
		return nil, ErrCodeAlreadyExists
	}

	// Cek kategori
	category, err := uc.repo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category.ItemType != enums.ItemTypeAsset {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe ASSET")
	}

	// Cek product line jika diisi
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

	itemID := uuid.NewV7().String()
	now := time.Now()

	item := Item{
		ID:            itemID,
		Code:          code,
		Name:          name,
		GenericName:   req.GenericName,
		ItemType:      enums.ItemTypeAsset,
		CategoryID:    req.CategoryID,
		ProductLineID: req.ProductLineID,
		IsActive:      true,
		CreatedAt:     now,
		CreatedBy:     actor,
		UpdatedAt:     now,
		UpdatedBy:     actor,
	}

	asset := ItemAsset{
		ItemID:                  itemID,
		Brand:                   req.Brand,
		ModelName:               req.ModelName,
		IsMedicalEquipment:      req.IsMedicalEquipment,
		ExpectedLifeYears:       req.ExpectedLifeYears,
		DepreciationMethod:      req.DepreciationMethod,
		MaintenanceIntervalDays: req.MaintenanceIntervalDays,
		CreatedAt:               now,
		CreatedBy:               actor,
		UpdatedAt:               now,
		UpdatedBy:               actor,
	}

	baseUnit := ItemUnit{
		ID:               uuid.NewV7().String(),
		ItemID:           itemID,
		UnitName:         baseUnitName,
		ConversionFactor: 1.0,
		IsBaseUnit:       true,
		IsPurchaseUnit:   true,
		IsDispenseUnit:   false,
		CreatedAt:        now,
		CreatedBy:        actor,
		UpdatedAt:        now,
		UpdatedBy:        actor,
	}

	// Eksekusi transaksi atomik (items + item_assets + item_units)
	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.CreateItem(ctx, &item); err != nil {
			return err
		}
		if err := txRepo.CreateAsset(ctx, &asset); err != nil {
			return err
		}
		if err := txRepo.CreateUnit(ctx, &baseUnit); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	savedItem, err := uc.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	res := ToItemDetailResponse(savedItem)
	return &res, nil
}
