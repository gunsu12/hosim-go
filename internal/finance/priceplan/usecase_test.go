package priceplan

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"gorm.io/gorm"
)

// -------------------------------------------------------------
// Unit Tests: Entity Business Logic (CITO & Component Balancing)
// -------------------------------------------------------------

func TestValidateBalancing(t *testing.T) {
	itemBalanced := &TariffPricePlanItem{
		TotalBasePrice: 250000.00,
		Components: []TariffPricePlanItemComponent{
			{BaseAmount: 75000.00},
			{BaseAmount: 175000.00},
		},
	}
	if err := itemBalanced.ValidateBalancing(); err != nil {
		t.Fatalf("expected balancing to succeed, got %v", err)
	}

	itemUnbalanced := &TariffPricePlanItem{
		TotalBasePrice: 250000.00,
		Components: []TariffPricePlanItemComponent{
			{BaseAmount: 75000.00},
			{BaseAmount: 100000.00}, // Total 175000 != 250000
		},
	}
	if err := itemUnbalanced.ValidateBalancing(); !errors.Is(err, ErrTariffComponentsMismatch) {
		t.Fatalf("expected ErrTariffComponentsMismatch, got %v", err)
	}
}

func TestCalculatePrice_Regular(t *testing.T) {
	coaSarana := "410.01.001"
	coaDokter := "214.01.001"

	item := &TariffPricePlanItem{
		TotalBasePrice: 250000.00,
		Components: []TariffPricePlanItemComponent{
			{
				ComponentID:   "comp-1",
				ComponentName: "Jasa Sarana RS",
				ComponentType: "SARANA",
				BaseAmount:    75000.00,
				COACode:       &coaSarana,
			},
			{
				ComponentID:   "comp-2",
				ComponentName: "Jasa Medis Dokter",
				ComponentType: "MEDIS_DOKTER",
				BaseAmount:    175000.00,
				COACode:       &coaDokter,
			},
		},
	}

	totalPrice, comps := item.CalculatePrice(false, 25.00)
	if totalPrice != 250000.00 {
		t.Errorf("expected regular total price 250000.00, got %f", totalPrice)
	}
	if len(comps) != 2 {
		t.Fatalf("expected 2 components, got %d", len(comps))
	}
	if comps[0].Amount != 75000.00 || comps[1].Amount != 175000.00 {
		t.Errorf("component amounts do not match regular base amount: %+v", comps)
	}
}

func TestCalculatePrice_CitoFormula(t *testing.T) {
	// Formula CITO: kenaikan default 25% (total dan per komponen)
	item := &TariffPricePlanItem{
		TotalBasePrice: 200000.00,
		TotalCitoPrice: nil, // hitung via formula 25% -> 250,000.00
		Components: []TariffPricePlanItemComponent{
			{
				ComponentID:   "comp-1",
				ComponentName: "Jasa Sarana RS",
				ComponentType: "SARANA",
				BaseAmount:    50000.00, // Cito +25% = 62,500.00
				CitoAmount:    nil,
			},
			{
				ComponentID:   "comp-2",
				ComponentName: "Jasa Medis Dokter",
				ComponentType: "MEDIS_DOKTER",
				BaseAmount:    150000.00, // Cito +25% = 187,500.00
				CitoAmount:    nil,
			},
		},
	}

	totalPrice, comps := item.CalculatePrice(true, 25.00)
	expectedTotal := 250000.00
	if math.Abs(totalPrice-expectedTotal) > 0.01 {
		t.Errorf("expected cito total price %f, got %f", expectedTotal, totalPrice)
	}
	if math.Abs(comps[0].Amount-62500.00) > 0.01 {
		t.Errorf("expected component 0 amount 62500.00, got %f", comps[0].Amount)
	}
	if math.Abs(comps[1].Amount-187500.00) > 0.01 {
		t.Errorf("expected component 1 amount 187500.00, got %f", comps[1].Amount)
	}
}

func TestCalculatePrice_CitoOverride(t *testing.T) {
	// CITO Override: Total dan rincian komponen memiliki nominal eksplisit
	customCitoTotal := 300000.00
	customCitoSarana := 100000.00
	customCitoDokter := 200000.00

	item := &TariffPricePlanItem{
		TotalBasePrice: 200000.00,
		TotalCitoPrice: &customCitoTotal,
		Components: []TariffPricePlanItemComponent{
			{
				ComponentID:   "comp-1",
				ComponentName: "Jasa Sarana RS",
				BaseAmount:    50000.00,
				CitoAmount:    &customCitoSarana,
			},
			{
				ComponentID:   "comp-2",
				ComponentName: "Jasa Medis Dokter",
				BaseAmount:    150000.00,
				CitoAmount:    &customCitoDokter,
			},
		},
	}

	totalPrice, comps := item.CalculatePrice(true, 25.00)
	if totalPrice != 300000.00 {
		t.Errorf("expected override cito price 300000.00, got %f", totalPrice)
	}
	if comps[0].Amount != 100000.00 || comps[1].Amount != 200000.00 {
		t.Errorf("expected override component amounts, got %+v", comps)
	}
}

func TestStateTransitions(t *testing.T) {
	plan := &TariffPricePlan{Status: PricePlanStatusDraft}

	// Valid transitions
	if err := plan.CanTransitionTo(PricePlanStatusSubmitted); err != nil {
		t.Errorf("expected DRAFT -> SUBMITTED to be valid, got %v", err)
	}
	if err := plan.CanTransitionTo(PricePlanStatusActive); !errors.Is(err, ErrInvalidStateTransition) {
		t.Errorf("expected DRAFT -> ACTIVE to be invalid, got %v", err)
	}

	plan.Status = PricePlanStatusSubmitted
	if err := plan.CanTransitionTo(PricePlanStatusApproved); err != nil {
		t.Errorf("expected SUBMITTED -> APPROVED to be valid, got %v", err)
	}
	if err := plan.CanTransitionTo(PricePlanStatusDraft); err != nil {
		t.Errorf("expected SUBMITTED -> DRAFT (revisi) to be valid, got %v", err)
	}

	plan.Status = PricePlanStatusApproved
	if err := plan.CanTransitionTo(PricePlanStatusActive); err != nil {
		t.Errorf("expected APPROVED -> ACTIVE to be valid, got %v", err)
	}

	plan.Status = PricePlanStatusActive
	if err := plan.CanTransitionTo(PricePlanStatusArchived); err != nil {
		t.Errorf("expected ACTIVE -> ARCHIVED to be valid, got %v", err)
	}

	plan.Status = PricePlanStatusArchived
	if err := plan.CanTransitionTo(PricePlanStatusActive); !errors.Is(err, ErrInvalidStateTransition) {
		t.Errorf("expected ARCHIVED -> ACTIVE to be invalid, got %v", err)
	}
}

// -------------------------------------------------------------
// Mock Repository untuk Pengujian Use Case Lookup Bertingkat
// -------------------------------------------------------------

type mockRepository struct {
	customPlan  *TariffPricePlan
	defaultPlan *TariffPricePlan
	itemsMap    map[string]*TariffPricePlanItem // key: planID:itemID:classID
	planByID    *TariffPricePlan
	allItems    []TariffPricePlanItem
}

func (m *mockRepository) CreatePlan(ctx context.Context, plan *TariffPricePlan) error { return nil }
func (m *mockRepository) UpdatePlan(ctx context.Context, plan *TariffPricePlan) error { return nil }
func (m *mockRepository) FindPlanByID(ctx context.Context, id string) (*TariffPricePlan, error) {
	if m.planByID != nil {
		return m.planByID, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockRepository) FindPlanByCode(ctx context.Context, code string) (*TariffPricePlan, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockRepository) FindAllPlans(ctx context.Context, params PlanListParams) ([]TariffPricePlan, int64, error) {
	return nil, 0, nil
}
func (m *mockRepository) CheckDefaultConflict(ctx context.Context, planID string, effectiveFrom time.Time, effectiveTo *time.Time) (bool, error) {
	return false, nil
}
func (m *mockRepository) AddItem(ctx context.Context, item *TariffPricePlanItem) error { return nil }
func (m *mockRepository) UpdateItem(ctx context.Context, item *TariffPricePlanItem) error {
	return nil
}
func (m *mockRepository) DeleteItem(ctx context.Context, planID, itemID string, deletedBy string) error {
	return nil
}
func (m *mockRepository) FindItemByID(ctx context.Context, planID, itemID string) (*TariffPricePlanItem, error) {
	return nil, nil
}
func (m *mockRepository) FindItemsByPlanID(ctx context.Context, planID string, params ItemListParams) ([]TariffPricePlanItem, int64, error) {
	return nil, 0, nil
}
func (m *mockRepository) FindAllItemsByPlanID(ctx context.Context, planID string) ([]TariffPricePlanItem, error) {
	return m.allItems, nil
}
func (m *mockRepository) BatchUpsertItems(ctx context.Context, planID string, items []TariffPricePlanItem, operator string) error {
	return nil
}

func (m *mockRepository) FindActiveCustomPlan(ctx context.Context, customerID string, txDate time.Time) (*TariffPricePlan, error) {
	if m.customPlan != nil && m.customPlan.CustomerID != nil && *m.customPlan.CustomerID == customerID {
		return m.customPlan, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) FindActiveDefaultPlan(ctx context.Context, txDate time.Time) (*TariffPricePlan, error) {
	if m.defaultPlan != nil {
		return m.defaultPlan, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) FindPlanItemForLookup(ctx context.Context, planID, itemID, tariffClassID string) (*TariffPricePlanItem, error) {
	key := planID + ":" + itemID + ":" + tariffClassID
	if it, ok := m.itemsMap[key]; ok {
		return it, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) ClonePlan(ctx context.Context, sourcePlanID string, newPlan *TariffPricePlan) error {
	return nil
}

// -------------------------------------------------------------
// Test: LookupTariff Hierarchical Fallback Resolution (PRD § 4.4)
// -------------------------------------------------------------

func TestLookupTariffUseCase_HierarchicalResolution(t *testing.T) {
	ctx := context.Background()
	custID := "cust-bpjs-001"
	otherCustID := "cust-inhealth-002"
	itemID := "item-konsultasi-spesialis"
	classID := "class-vip"

	customPlan := &TariffPricePlan{
		ID:                 "plan-bpjs",
		Code:               "TPP-2026-BPJS",
		Name:               "Tarif Khusus BPJS 2026",
		Status:             PricePlanStatusActive,
		CustomerID:         &custID,
		DefaultCitoPercent: 25.00,
	}

	defaultPlan := &TariffPricePlan{
		ID:                 "plan-default",
		Code:               "TPP-2026-REGULER",
		Name:               "Tarif Standar RS 2026",
		Status:             PricePlanStatusActive,
		IsDefault:          true,
		DefaultCitoPercent: 30.00,
	}

	customItem := &TariffPricePlanItem{
		ID:              "item-custom-1",
		PricePlanID:     "plan-bpjs",
		ItemID:          itemID,
		ItemName:        "Konsultasi Dokter Spesialis",
		TariffClassID:   classID,
		TariffClassName: "Kelas VIP",
		TotalBasePrice:  180000.00,
		IsActive:        true,
		Components: []TariffPricePlanItemComponent{
			{ComponentID: "c1", ComponentName: "Jasa Sarana", BaseAmount: 50000.00},
			{ComponentID: "c2", ComponentName: "Jasa Medis", BaseAmount: 130000.00},
		},
	}

	defaultItem := &TariffPricePlanItem{
		ID:              "item-default-1",
		PricePlanID:     "plan-default",
		ItemID:          itemID,
		ItemName:        "Konsultasi Dokter Spesialis",
		TariffClassID:   classID,
		TariffClassName: "Kelas VIP",
		TotalBasePrice:  250000.00,
		IsActive:        true,
		Components: []TariffPricePlanItemComponent{
			{ComponentID: "c1", ComponentName: "Jasa Sarana", BaseAmount: 75000.00},
			{ComponentID: "c2", ComponentName: "Jasa Medis", BaseAmount: 175000.00},
		},
	}

	mockRepo := &mockRepository{
		customPlan:  customPlan,
		defaultPlan: defaultPlan,
		itemsMap: map[string]*TariffPricePlanItem{
			"plan-bpjs:item-konsultasi-spesialis:class-vip":    customItem,
			"plan-default:item-konsultasi-spesialis:class-vip": defaultItem,
		},
	}

	uc := NewLookupTariffUseCase(mockRepo)

	// Kasus 1: Langkah 1 - Custom Plan Hit
	t.Run("Step 1: Custom Plan Hit", func(t *testing.T) {
		res, err := uc.Execute(ctx, LookupTariffRequest{
			CustomerID:      &custID,
			TariffClassID:   classID,
			ItemID:          itemID,
			TransactionDate: "2026-09-27T10:00:00Z",
			IsCito:          false,
		})
		if err != nil {
			t.Fatalf("expected lookup success, got %v", err)
		}
		if !res.IsCustomPlan {
			t.Errorf("expected is_custom_plan to be true")
		}
		if res.PricePlanCode != "TPP-2026-BPJS" {
			t.Errorf("expected plan code TPP-2026-BPJS, got %s", res.PricePlanCode)
		}
		if res.TotalPrice != 180000.00 {
			t.Errorf("expected price 180000.00, got %f", res.TotalPrice)
		}
	})

	// Kasus 2: Langkah 2 - Fallback ke Default/Standard Plan
	t.Run("Step 2: Fallback to General Plan when customer has no custom plan", func(t *testing.T) {
		res, err := uc.Execute(ctx, LookupTariffRequest{
			CustomerID:      &otherCustID, // Penjamin lain yang tidak punya custom plan
			TariffClassID:   classID,
			ItemID:          itemID,
			TransactionDate: "2026-09-27T10:00:00Z",
			IsCito:          false,
		})
		if err != nil {
			t.Fatalf("expected fallback success, got %v", err)
		}
		if res.IsCustomPlan {
			t.Errorf("expected is_custom_plan to be false for fallback")
		}
		if res.PricePlanCode != "TPP-2026-REGULER" {
			t.Errorf("expected plan code TPP-2026-REGULER, got %s", res.PricePlanCode)
		}
		if res.TotalPrice != 250000.00 {
			t.Errorf("expected price 250000.00, got %f", res.TotalPrice)
		}
	})

	// Kasus 3: Pasien Non-Penjamin (Pribadi / Umum tanpa customer_id)
	t.Run("Step 2: Fallback to General Plan for direct patient without customer_id", func(t *testing.T) {
		res, err := uc.Execute(ctx, LookupTariffRequest{
			CustomerID:      nil,
			TariffClassID:   classID,
			ItemID:          itemID,
			TransactionDate: "2026-09-27T10:00:00Z",
			IsCito:          false,
		})
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if res.IsCustomPlan {
			t.Errorf("expected is_custom_plan to be false")
		}
		if res.TotalPrice != 250000.00 {
			t.Errorf("expected price 250000.00, got %f", res.TotalPrice)
		}
	})

	// Kasus 4: Langkah 3 - Tindakan tidak terdaftar di custom maupun default plan
	t.Run("Step 3: Not Found returns ErrTariffNotConfigured", func(t *testing.T) {
		_, err := uc.Execute(ctx, LookupTariffRequest{
			CustomerID:      &custID,
			TariffClassID:   classID,
			ItemID:          "item-tidak-ada",
			TransactionDate: "2026-09-27T10:00:00Z",
			IsCito:          false,
		})
		if !errors.Is(err, ErrTariffNotConfigured) {
			t.Fatalf("expected ErrTariffNotConfigured, got %v", err)
		}
	})

	// Kasus 5: Resolusi CITO pada Lookup
	t.Run("Lookup with is_cito = true", func(t *testing.T) {
		res, err := uc.Execute(ctx, LookupTariffRequest{
			CustomerID:      nil, // General Plan (DefaultCitoPercent = 30%)
			TariffClassID:   classID,
			ItemID:          itemID,
			TransactionDate: "2026-09-27T10:00:00Z",
			IsCito:          true,
		})
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		// Base: 250,000 + 30% = 325,000
		expectedCitoTotal := 325000.00
		if math.Abs(res.TotalPrice-expectedCitoTotal) > 0.01 {
			t.Errorf("expected cito total %f, got %f", expectedCitoTotal, res.TotalPrice)
		}
		if !res.IsCito {
			t.Errorf("expected is_cito to be true in response")
		}
	})
}

func TestCalculatePrice_ZeroCitoOverride(t *testing.T) {
	zeroCito := 0.00
	item := &TariffPricePlanItem{
		TotalBasePrice: 200000.00,
		TotalCitoPrice: &zeroCito,
		Components: []TariffPricePlanItemComponent{
			{
				ComponentID: "c1",
				BaseAmount:  200000.00,
				CitoAmount:  &zeroCito,
			},
		},
	}

	totalPrice, comps := item.CalculatePrice(true, 25.00)
	if totalPrice != 0.00 {
		t.Errorf("expected 0.00 for explicit zero CITO override, got %f", totalPrice)
	}
	if len(comps) != 1 || comps[0].Amount != 0.00 {
		t.Errorf("expected component amount 0.00, got %+v", comps)
	}
}

func TestApprovePricePlan_EmptyItems(t *testing.T) {
	plan := &TariffPricePlan{
		ID:     "plan-empty",
		Status: PricePlanStatusSubmitted,
	}
	mockRepo := &mockRepository{
		planByID: plan,
		allItems: []TariffPricePlanItem{}, // 0 items
	}
	uc := NewApprovePricePlanUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), "plan-empty", "ADMIN")
	if err == nil {
		t.Fatalf("expected error when approving plan with 0 items, got nil")
	}
}

func TestCreatePricePlan_DefaultWithCustomerDisallowed(t *testing.T) {
	mockRepo := &mockRepository{}
	uc := NewCreatePricePlanUseCase(mockRepo)

	custID := "cust-123"
	_, err := uc.Execute(context.Background(), CreatePricePlanRequest{
		Code:          "TPP-TEST",
		Name:          "Test Plan",
		EffectiveFrom: "2026-01-01",
		IsDefault:     true,
		CustomerID:    &custID,
	}, "ADMIN")

	if err == nil {
		t.Fatalf("expected error when creating default plan associated with customer, got nil")
	}
}

func TestManageItems_DuplicateComponentDisallowed(t *testing.T) {
	plan := &TariffPricePlan{
		ID:     "plan-draft",
		Status: PricePlanStatusDraft,
	}
	mockRepo := &mockRepository{
		planByID: plan,
	}
	uc := NewManageItemsUseCase(mockRepo)

	_, err := uc.AddItem(context.Background(), "plan-draft", AddItemRequest{
		ItemID:         "item-1",
		TariffClassID:  "class-1",
		TotalBasePrice: 100000,
		Components: []ItemComponentRequest{
			{ComponentID: "c1", BaseAmount: 50000},
			{ComponentID: "c1", BaseAmount: 50000}, // Duplicate component_id
		},
	}, "ADMIN")

	if !errors.Is(err, ErrDuplicateComponent) {
		t.Fatalf("expected ErrDuplicateComponent, got %v", err)
	}
}

