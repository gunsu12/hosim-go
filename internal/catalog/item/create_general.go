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
	ErrInvalidGeneralType        = errors.New("subtipe barang umum / general_type tidak valid")
	ErrCSSDSterilizationRequired = errors.New("metode sterilisasi wajib diisi untuk barang siklus CSSD")
)

type CreateGeneralUseCase interface {
	Execute(ctx context.Context, req CreateGeneralRequest, actor string) (*ItemDetailResponse, error)
}

type createGeneralUseCase struct {
	repo Repository
}

func NewCreateGeneralUseCase(repo Repository) CreateGeneralUseCase {
	return &createGeneralUseCase{repo: repo}
}

func (uc *createGeneralUseCase) Execute(ctx context.Context, req CreateGeneralRequest, actor string) (*ItemDetailResponse, error) {
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
	if category.ItemType != enums.ItemTypeGeneral {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe GENERAL")
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

	itemID := uuid.NewV7().String()
	now := time.Now()

	item := Item{
		ID:            itemID,
		Code:          code,
		Name:          name,
		GenericName:   req.GenericName,
		ItemType:      enums.ItemTypeGeneral,
		CategoryID:    req.CategoryID,
		ProductLineID: req.ProductLineID,
		IsActive:      true,
		CreatedAt:     now,
		CreatedBy:     actor,
		UpdatedAt:     now,
		UpdatedBy:     actor,
	}

	general := ItemGeneral{
		ItemID:              itemID,
		GeneralType:         req.GeneralType,
		IsSterile:           req.IsSterile,
		IsDisposable:        isDisposable,
		IsCSSDItem:          req.IsCSSDItem,
		SterilizationMethod: req.SterilizationMethod,
		CreatedAt:           now,
		CreatedBy:           actor,
		UpdatedAt:           now,
		UpdatedBy:           actor,
	}

	baseUnit := ItemUnit{
		ID:               uuid.NewV7().String(),
		ItemID:           itemID,
		UnitName:         baseUnitName,
		ConversionFactor: 1.0,
		IsBaseUnit:       true,
		IsPurchaseUnit:   false,
		IsDispenseUnit:   true,
		CreatedAt:        now,
		CreatedBy:        actor,
		UpdatedAt:        now,
		UpdatedBy:        actor,
	}

	// Eksekusi transaksi atomik (items + item_generals + item_units)
	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.CreateItem(ctx, &item); err != nil {
			return err
		}
		if err := txRepo.CreateGeneral(ctx, &general); err != nil {
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
