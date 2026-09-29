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
	ErrInvalidChargeType = errors.New("klasifikasi beban tarif / charge_type tidak valid")
)

type CreateTariffUseCase interface {
	Execute(ctx context.Context, req CreateTariffRequest, actor string) (*ItemDetailResponse, error)
}

type createTariffUseCase struct {
	repo Repository
}

func NewCreateTariffUseCase(repo Repository) CreateTariffUseCase {
	return &createTariffUseCase{repo: repo}
}

func (uc *createTariffUseCase) Execute(ctx context.Context, req CreateTariffRequest, actor string) (*ItemDetailResponse, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	baseUnitName := strings.TrimSpace(req.BaseUnitName)
	if baseUnitName == "" {
		baseUnitName = "Kali"
	}

	if code == "" {
		return nil, errors.New("kode item wajib diisi")
	}
	if name == "" {
		return nil, errors.New("nama item wajib diisi")
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
	if category.ItemType != enums.ItemTypeTariff {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe TARIFF")
	}

	if !req.ChargeType.IsValid() {
		return nil, ErrInvalidChargeType
	}

	itemID := uuid.NewV7().String()
	now := time.Now()

	// Invariant PRD: Layanan tarif non-inventory, ProductLineID selalu nil
	item := Item{
		ID:            itemID,
		Code:          code,
		Name:          name,
		ItemType:      enums.ItemTypeTariff,
		CategoryID:    req.CategoryID,
		ProductLineID: nil,
		IsActive:      true,
		CreatedAt:     now,
		CreatedBy:     actor,
		UpdatedAt:     now,
		UpdatedBy:     actor,
	}

	tariff := ItemTariff{
		ItemID:     itemID,
		NameAlias:  req.NameAlias,
		ChargeType: req.ChargeType,
		Notes:      req.Notes,
		CreatedAt:  now,
		CreatedBy:  actor,
		UpdatedAt:  now,
		UpdatedBy:  actor,
	}

	// Invariant PRD 4.1: Standarisasi UOM untuk tarif (is_base=true, is_dispense=true, is_purchase=false)
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

	// Eksekusi transaksi atomik (items + item_tariffs + item_units)
	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.CreateItem(ctx, &item); err != nil {
			return err
		}
		if err := txRepo.CreateTariff(ctx, &tariff); err != nil {
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
