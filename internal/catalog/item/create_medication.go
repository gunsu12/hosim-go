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
	ErrInvalidMedicationType = errors.New("golongan obat tidak valid")
)

type CreateMedicationUseCase interface {
	Execute(ctx context.Context, req CreateMedicationRequest, actor string) (*ItemDetailResponse, error)
}

type createMedicationUseCase struct {
	repo Repository
}

func NewCreateMedicationUseCase(repo Repository) CreateMedicationUseCase {
	return &createMedicationUseCase{repo: repo}
}

func (uc *createMedicationUseCase) Execute(ctx context.Context, req CreateMedicationRequest, actor string) (*ItemDetailResponse, error) {
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
	if category.ItemType != enums.ItemTypeMedication {
		return nil, errors.New("kategori yang dipilih bukan untuk tipe MEDICATION")
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

	if req.MedicationType != "" && !req.MedicationType.IsValid() {
		return nil, ErrInvalidMedicationType
	}

	itemID := uuid.NewV7().String()
	now := time.Now()

	item := Item{
		ID:            itemID,
		Code:          code,
		Name:          name,
		GenericName:   req.GenericName,
		ItemType:      enums.ItemTypeMedication,
		CategoryID:    req.CategoryID,
		ProductLineID: req.ProductLineID,
		IsActive:      true,
		CreatedAt:     now,
		CreatedBy:     actor,
		UpdatedAt:     now,
		UpdatedBy:     actor,
	}

	medication := ItemMedication{
		ItemID:             itemID,
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
		CreatedAt:          now,
		CreatedBy:          actor,
		UpdatedAt:          now,
		UpdatedBy:          actor,
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

	// Eksekusi transaksi atomik (items + item_medications + item_units)
	err = uc.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.CreateItem(ctx, &item); err != nil {
			return err
		}
		if err := txRepo.CreateMedication(ctx, &medication); err != nil {
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

	// Ambil kembali aggregate lengkap
	savedItem, err := uc.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	res := ToItemDetailResponse(savedItem)
	return &res, nil
}
