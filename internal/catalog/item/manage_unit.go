package item

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"
)

type ManageUnitUseCase interface {
	ListUnits(ctx context.Context, itemID string) ([]UnitResponse, error)
	AddUnit(ctx context.Context, itemID string, req AddUnitRequest, actor string) (*UnitResponse, error)
	UpdateUnit(ctx context.Context, itemID string, unitID string, req UpdateUnitRequest, actor string) (*UnitResponse, error)
	DeleteUnit(ctx context.Context, itemID string, unitID string) error
}

type manageUnitUseCase struct {
	repo Repository
}

func NewManageUnitUseCase(repo Repository) ManageUnitUseCase {
	return &manageUnitUseCase{repo: repo}
}

func (uc *manageUnitUseCase) AddUnit(ctx context.Context, itemID string, req AddUnitRequest, actor string) (*UnitResponse, error) {
	unitName := strings.TrimSpace(req.UnitName)
	if unitName == "" {
		return nil, errors.New("nama satuan wajib diisi")
	}
	if req.ConversionFactor <= 0 {
		return nil, errors.New("faktor konversi satuan harus lebih besar dari 0")
	}

	// Pastikan item ada
	_, err := uc.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Cek apakah nama satuan sudah pernah didaftarkan untuk item ini
	existingUnits, err := uc.repo.FindUnitsByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	for _, u := range existingUnits {
		if strings.EqualFold(u.UnitName, unitName) {
			return nil, ErrDuplicateUnitName
		}
	}

	now := time.Now()
	unit := ItemUnit{
		ID:               uuid.NewV7().String(),
		ItemID:           itemID,
		UnitName:         unitName,
		ConversionFactor: req.ConversionFactor,
		IsBaseUnit:       false, // Satuan alternatif tidak boleh menggantikan base unit
		IsPurchaseUnit:   req.IsPurchaseUnit,
		IsDispenseUnit:   req.IsDispenseUnit,
		CreatedAt:        now,
		CreatedBy:        actor,
		UpdatedAt:        now,
		UpdatedBy:        actor,
	}

	if err := uc.repo.CreateUnit(ctx, &unit); err != nil {
		return nil, err
	}

	res := ToUnitResponse(&unit)
	return &res, nil
}

func (uc *manageUnitUseCase) ListUnits(ctx context.Context, itemID string) ([]UnitResponse, error) {
	_, err := uc.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	units, err := uc.repo.FindUnitsByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return ToUnitResponseList(units), nil
}

func (uc *manageUnitUseCase) UpdateUnit(ctx context.Context, itemID string, unitID string, req UpdateUnitRequest, actor string) (*UnitResponse, error) {
	unitName := strings.TrimSpace(req.UnitName)
	if unitName == "" {
		return nil, errors.New("nama satuan wajib diisi")
	}
	if req.ConversionFactor <= 0 {
		return nil, errors.New("faktor konversi satuan harus lebih besar dari 0")
	}

	unit, err := uc.repo.FindUnitByID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if unit.ItemID != itemID {
		return nil, errors.New("satuan tidak sesuai dengan item yang diminta")
	}

	if !strings.EqualFold(unit.UnitName, unitName) {
		existingUnits, err := uc.repo.FindUnitsByItemID(ctx, itemID)
		if err != nil {
			return nil, err
		}
		for _, u := range existingUnits {
			if u.ID != unitID && strings.EqualFold(u.UnitName, unitName) {
				return nil, ErrDuplicateUnitName
			}
		}
	}

	// Base unit tidak boleh diubah faktor konversinya dari 1.0
	if unit.IsBaseUnit && req.ConversionFactor != 1.0 {
		return nil, errors.New("faktor konversi base unit harus tetap 1.0")
	}

	unit.UnitName = unitName
	if !unit.IsBaseUnit {
		unit.ConversionFactor = req.ConversionFactor
		unit.IsPurchaseUnit = req.IsPurchaseUnit
		unit.IsDispenseUnit = req.IsDispenseUnit
	}
	unit.UpdatedAt = time.Now()
	unit.UpdatedBy = actor

	if err := uc.repo.UpdateUnit(ctx, unit); err != nil {
		return nil, err
	}

	res := ToUnitResponse(unit)
	return &res, nil
}

func (uc *manageUnitUseCase) DeleteUnit(ctx context.Context, itemID string, unitID string) error {
	unit, err := uc.repo.FindUnitByID(ctx, unitID)
	if err != nil {
		return err
	}

	if unit.ItemID != itemID {
		return errors.New("satuan tidak sesuai dengan item yang diminta")
	}

	if unit.IsBaseUnit {
		return ErrBaseUnitCannotBeDeleted
	}

	return uc.repo.DeleteUnit(ctx, unitID)
}
