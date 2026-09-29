package item_test

import (
	"context"
	"errors"
	"testing"

	"hosim-go/internal/catalog/item"
	"hosim-go/pkg/enums"
)

type mockItemRepo struct {
	item.Repository
	findItemByCodeFn     func(ctx context.Context, code string) (*item.Item, error)
	findCategoryByIDFn   func(ctx context.Context, id string) (*item.ItemCategory, error)
	productLineExistsFn  func(ctx context.Context, id string) (bool, error)
	findItemByIDFn       func(ctx context.Context, id string) (*item.Item, error)
	findUnitsByItemIDFn  func(ctx context.Context, itemID string) ([]item.ItemUnit, error)
	findUnitByIDFn       func(ctx context.Context, id string) (*item.ItemUnit, error)
	createItemFn         func(ctx context.Context, i *item.Item) error
	createMedicationFn   func(ctx context.Context, m *item.ItemMedication) error
	createTariffFn       func(ctx context.Context, t *item.ItemTariff) error
	createUnitFn         func(ctx context.Context, u *item.ItemUnit) error
	deleteUnitFn         func(ctx context.Context, id string) error
	withTransactionFn    func(ctx context.Context, fn func(txRepo item.Repository) error) error
}

func (m *mockItemRepo) FindItemByCode(ctx context.Context, code string) (*item.Item, error) {
	if m.findItemByCodeFn != nil {
		return m.findItemByCodeFn(ctx, code)
	}
	return nil, item.ErrItemNotFound
}

func (m *mockItemRepo) FindCategoryByID(ctx context.Context, id string) (*item.ItemCategory, error) {
	if m.findCategoryByIDFn != nil {
		return m.findCategoryByIDFn(ctx, id)
	}
	return nil, item.ErrCategoryNotFound
}

func (m *mockItemRepo) ProductLineExists(ctx context.Context, id string) (bool, error) {
	if m.productLineExistsFn != nil {
		return m.productLineExistsFn(ctx, id)
	}
	return true, nil
}

func (m *mockItemRepo) FindItemByID(ctx context.Context, id string) (*item.Item, error) {
	if m.findItemByIDFn != nil {
		return m.findItemByIDFn(ctx, id)
	}
	return &item.Item{ID: id, Name: "Dummy"}, nil
}

func (m *mockItemRepo) FindUnitsByItemID(ctx context.Context, itemID string) ([]item.ItemUnit, error) {
	if m.findUnitsByItemIDFn != nil {
		return m.findUnitsByItemIDFn(ctx, itemID)
	}
	return nil, nil
}

func (m *mockItemRepo) FindUnitByID(ctx context.Context, id string) (*item.ItemUnit, error) {
	if m.findUnitByIDFn != nil {
		return m.findUnitByIDFn(ctx, id)
	}
	return nil, item.ErrUnitNotFound
}

func (m *mockItemRepo) CreateItem(ctx context.Context, i *item.Item) error {
	if m.createItemFn != nil {
		return m.createItemFn(ctx, i)
	}
	return nil
}

func (m *mockItemRepo) CreateMedication(ctx context.Context, med *item.ItemMedication) error {
	if m.createMedicationFn != nil {
		return m.createMedicationFn(ctx, med)
	}
	return nil
}

func (m *mockItemRepo) CreateTariff(ctx context.Context, t *item.ItemTariff) error {
	if m.createTariffFn != nil {
		return m.createTariffFn(ctx, t)
	}
	return nil
}

func (m *mockItemRepo) CreateUnit(ctx context.Context, u *item.ItemUnit) error {
	if m.createUnitFn != nil {
		return m.createUnitFn(ctx, u)
	}
	return nil
}

func (m *mockItemRepo) DeleteUnit(ctx context.Context, id string) error {
	if m.deleteUnitFn != nil {
		return m.deleteUnitFn(ctx, id)
	}
	return nil
}

func (m *mockItemRepo) WithTransaction(ctx context.Context, fn func(txRepo item.Repository) error) error {
	if m.withTransactionFn != nil {
		return m.withTransactionFn(ctx, fn)
	}
	return fn(m)
}

func TestCreateMedication_DuplicateCode(t *testing.T) {
	mockRepo := &mockItemRepo{
		findItemByCodeFn: func(ctx context.Context, code string) (*item.Item, error) {
			return &item.Item{ID: "existing-id", Code: code}, nil
		},
	}

	uc := item.NewCreateMedicationUseCase(mockRepo)
	req := item.CreateMedicationRequest{
		Code:         "MED-001",
		Name:         "Paracetamol 500mg",
		CategoryID:   "cat-1",
		BaseUnitName: "Tablet",
	}

	_, err := uc.Execute(context.Background(), req, "tester")
	if !errors.Is(err, item.ErrCodeAlreadyExists) {
		t.Fatalf("expected ErrCodeAlreadyExists, got %v", err)
	}
}

func TestCreateTariff_NonInventoryInvariant(t *testing.T) {
	var capturedItem *item.Item
	var capturedUnit *item.ItemUnit

	mockRepo := &mockItemRepo{
		findCategoryByIDFn: func(ctx context.Context, id string) (*item.ItemCategory, error) {
			return &item.ItemCategory{ID: id, ItemType: enums.ItemTypeTariff}, nil
		},
		createItemFn: func(ctx context.Context, i *item.Item) error {
			capturedItem = i
			return nil
		},
		createUnitFn: func(ctx context.Context, u *item.ItemUnit) error {
			capturedUnit = u
			return nil
		},
		findItemByIDFn: func(ctx context.Context, id string) (*item.Item, error) {
			return &item.Item{
				ID:       id,
				Code:     "TAR-001",
				Name:     "Konsul Dokter",
				ItemType: enums.ItemTypeTariff,
				Tariff: &item.ItemTariff{
					ChargeType: enums.ChargeTypeTindakan,
				},
				Units: []item.ItemUnit{
					{
						UnitName:         "Kali",
						ConversionFactor: 1.0,
						IsBaseUnit:       true,
						IsDispenseUnit:   true,
						IsPurchaseUnit:   false,
					},
				},
			}, nil
		},
	}

	uc := item.NewCreateTariffUseCase(mockRepo)
	req := item.CreateTariffRequest{
		Code:       "TAR-001",
		Name:       "Konsul Dokter",
		CategoryID: "cat-tariff",
		ChargeType: enums.ChargeTypeTindakan,
	}

	res, err := uc.Execute(context.Background(), req, "tester")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedItem == nil {
		t.Fatal("expected item to be created")
	}
	if capturedItem.ProductLineID != nil {
		t.Errorf("expected product_line_id to be nil for tariff, got %v", *capturedItem.ProductLineID)
	}
	if capturedUnit == nil || !capturedUnit.IsBaseUnit || !capturedUnit.IsDispenseUnit || capturedUnit.IsPurchaseUnit {
		t.Errorf("invalid tariff base unit configuration: %+v", capturedUnit)
	}
	if res.Code != "TAR-001" {
		t.Errorf("expected code TAR-001, got %s", res.Code)
	}
}

func TestManageUnit_CannotDeleteBaseUnit(t *testing.T) {
	mockRepo := &mockItemRepo{
		findUnitByIDFn: func(ctx context.Context, id string) (*item.ItemUnit, error) {
			return &item.ItemUnit{
				ID:         id,
				ItemID:     "item-123",
				UnitName:   "Tablet",
				IsBaseUnit: true, // Unit dasar!
			}, nil
		},
	}

	uc := item.NewManageUnitUseCase(mockRepo)
	err := uc.DeleteUnit(context.Background(), "item-123", "unit-base")
	if !errors.Is(err, item.ErrBaseUnitCannotBeDeleted) {
		t.Fatalf("expected ErrBaseUnitCannotBeDeleted, got %v", err)
	}
}
