package priceschema

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"hosim-go/internal/audit"

	"github.com/shopspring/decimal"
)

// Mock transaction manager
type mockTxManager struct{}

func (m *mockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// Mock audit writer
type mockAuditWriter struct {
	entries []audit.AuditEntry
}

func (m *mockAuditWriter) Write(ctx context.Context, entry audit.AuditEntry) error {
	m.entries = append(m.entries, entry)
	return nil
}

// Mock in-memory repository
type mockRepository struct {
	schemas map[string]*PriceSchema
	rules   map[string]*PriceSchemaRule
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		schemas: make(map[string]*PriceSchema),
		rules:   make(map[string]*PriceSchemaRule),
	}
}

func (m *mockRepository) CreateSchema(ctx context.Context, schema *PriceSchema) error {
	if schema.ID == "" {
		schema.ID = "schema-" + schema.Code
	}
	m.schemas[schema.ID] = schema
	return nil
}

func (m *mockRepository) UpdateSchema(ctx context.Context, schema *PriceSchema) error {
	m.schemas[schema.ID] = schema
	return nil
}

func (m *mockRepository) FindSchemaByID(ctx context.Context, id string) (*PriceSchema, error) {
	s, ok := m.schemas[id]
	if !ok {
		return nil, ErrSchemaNotFound
	}
	return s, nil
}

func (m *mockRepository) FindSchemaByCode(ctx context.Context, code string) (*PriceSchema, error) {
	for _, s := range m.schemas {
		if s.Code == code {
			return s, nil
		}
	}
	return nil, ErrSchemaNotFound
}

func (m *mockRepository) ListSchemas(ctx context.Context, params ListSchemaParams) ([]PriceSchema, int64, error) {
	var list []PriceSchema
	for _, s := range m.schemas {
		list = append(list, *s)
	}
	return list, int64(len(list)), nil
}

func (m *mockRepository) CreateRule(ctx context.Context, rule *PriceSchemaRule) error {
	if rule.ID == "" {
		rule.ID = "rule-" + rule.PriceSchemaID + "-" + string(rule.Level())
	}
	m.rules[rule.ID] = rule
	return nil
}

func (m *mockRepository) UpdateRule(ctx context.Context, rule *PriceSchemaRule) error {
	m.rules[rule.ID] = rule
	return nil
}

func (m *mockRepository) DeleteRule(ctx context.Context, id string) error {
	delete(m.rules, id)
	return nil
}

func (m *mockRepository) FindRuleByID(ctx context.Context, id string) (*PriceSchemaRule, error) {
	r, ok := m.rules[id]
	if !ok {
		return nil, ErrRuleNotFound
	}
	return r, nil
}

func (m *mockRepository) FindRuleByScope(ctx context.Context, schemaID string, deptID, unitID, itemID *string) (*PriceSchemaRule, error) {
	for _, r := range m.rules {
		if r.PriceSchemaID != schemaID {
			continue
		}
		if matchPtr(r.DepartementID, deptID) && matchPtr(r.ServiceUnitID, unitID) && matchPtr(r.ItemID, itemID) {
			return r, nil
		}
	}
	return nil, nil
}

func matchPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil && *a == *b {
		return true
	}
	return false
}

func (m *mockRepository) ListRulesBySchemaID(ctx context.Context, schemaID string, params ListRuleParams) ([]PriceSchemaRule, int64, error) {
	var list []PriceSchemaRule
	for _, r := range m.rules {
		if r.PriceSchemaID == schemaID {
			list = append(list, *r)
		}
	}
	return list, int64(len(list)), nil
}

func (m *mockRepository) FindBestMatchingRule(ctx context.Context, schemaID, itemID, serviceUnitID, departementID string) (*PriceSchemaRule, error) {
	var candidates []PriceSchemaRule
	for _, r := range m.rules {
		if r.PriceSchemaID != schemaID {
			continue
		}
		// Cek apakah rule ini memenuhi salah satu kombinasi
		lvl := r.Level()
		switch lvl {
		case MatchLevelItemUnit:
			if itemID != "" && serviceUnitID != "" && r.ItemID != nil && *r.ItemID == itemID && r.ServiceUnitID != nil && *r.ServiceUnitID == serviceUnitID {
				candidates = append(candidates, *r)
			}
		case MatchLevelItemDept:
			if itemID != "" && departementID != "" && r.ItemID != nil && *r.ItemID == itemID && r.DepartementID != nil && *r.DepartementID == departementID {
				candidates = append(candidates, *r)
			}
		case MatchLevelItem:
			if itemID != "" && r.ItemID != nil && *r.ItemID == itemID {
				candidates = append(candidates, *r)
			}
		case MatchLevelUnit:
			if serviceUnitID != "" && r.ServiceUnitID != nil && *r.ServiceUnitID == serviceUnitID {
				candidates = append(candidates, *r)
			}
		case MatchLevelDept:
			if departementID != "" && r.DepartementID != nil && *r.DepartementID == departementID {
				candidates = append(candidates, *r)
			}
		case MatchLevelSchemaDefault:
			candidates = append(candidates, *r)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoMatchingRule
	}

	priorityOf := func(r PriceSchemaRule) int {
		switch r.Level() {
		case MatchLevelItemUnit:
			return 1
		case MatchLevelItemDept:
			return 2
		case MatchLevelItem:
			return 3
		case MatchLevelUnit:
			return 4
		case MatchLevelDept:
			return 5
		default:
			return 6
		}
	}

	best := candidates[0]
	bestP := priorityOf(best)
	for _, c := range candidates[1:] {
		p := priorityOf(c)
		if p < bestP {
			best = c
			bestP = p
		}
	}

	return &best, nil
}

// Mock readers
type mockOrgReader struct{}

func (m *mockOrgReader) IsActiveDepartment(ctx context.Context, id string) (bool, error) {
	return id == "dept-valid", nil
}

func (m *mockOrgReader) GetServiceUnitInfo(ctx context.Context, id string) (*string, bool, error) {
	if id == "unit-valid" {
		dept := "dept-valid"
		return &dept, true, nil
	}
	return nil, false, nil
}

type mockItemReader struct{}

func (m *mockItemReader) IsActiveItem(ctx context.Context, id string) (bool, error) {
	return id == "item-valid", nil
}

type mockPlanReader struct{}

func (m *mockPlanReader) PlanExists(ctx context.Context, id string) (bool, error) {
	return !strings.HasPrefix(id, "plan-invalid"), nil
}

// -------------------------------------------------------------
// Tests
// -------------------------------------------------------------

func TestCreatePriceSchema_AutomaticDefaultRule(t *testing.T) {
	repo := newMockRepository()
	tx := &mockTxManager{}
	auditW := &mockAuditWriter{}
	org := &mockOrgReader{}
	itemR := &mockItemReader{}
	planR := &mockPlanReader{}

	uc := NewCreatePriceSchemaUseCase(repo, tx, auditW, org, org, itemR, planR)
	ctx := context.Background()
	actor := OperatorContext{UserID: "usr-1", Username: "admin", Role: "ADMIN"}

	req := CreatePriceSchemaRequest{
		Code: "SCHEMA_BPJS",
		Name: "Skema BPJS Kesehatan",
	}

	res, err := uc.Execute(ctx, req, actor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Code != "SCHEMA_BPJS" {
		t.Errorf("expected code SCHEMA_BPJS, got %s", res.Code)
	}
	if res.RulesCount != 1 {
		t.Fatalf("expected 1 rule (automatic SCHEMA_DEFAULT), got %d", res.RulesCount)
	}
	if res.Rules[0].Level != MatchLevelSchemaDefault {
		t.Errorf("expected level SCHEMA_DEFAULT, got %s", res.Rules[0].Level)
	}

	// Cek audit entry
	if len(auditW.entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(auditW.entries))
	}
	if auditW.entries[0].Action != "CREATE" {
		t.Errorf("expected audit action CREATE, got %s", auditW.entries[0].Action)
	}
}

func TestRuleInvariants(t *testing.T) {
	dept := "dept-1"
	unit := "unit-1"
	item := "item-1"

	// 1. Departement dan Unit tidak boleh bersamaan
	rExclusive := PriceSchemaRule{
		DepartementID:        &dept,
		ServiceUnitID:        &unit,
		TariffAdjustmentType: AdjustmentTypeNone,
		PayByPatientPercent:  decimal.Zero,
		PayByCustomerPercent: decimal.NewFromInt(100),
	}
	if err := rExclusive.Validate(); !errors.Is(err, ErrScopeExclusive) {
		t.Errorf("expected ErrScopeExclusive, got %v", err)
	}

	// 2. FIXED_PRICE tanpa item dilarang
	rFixedNoItem := PriceSchemaRule{
		DepartementID:         &dept,
		TariffAdjustmentType:  AdjustmentTypeFixedPrice,
		TariffAdjustmentValue: decimal.NewFromInt(50000),
		PayByPatientPercent:   decimal.Zero,
		PayByCustomerPercent:  decimal.NewFromInt(100),
	}
	if err := rFixedNoItem.Validate(); !errors.Is(err, ErrFixedPriceRequiresItem) {
		t.Errorf("expected ErrFixedPriceRequiresItem, got %v", err)
	}

	// 3. FIXED_PRICE negatif dilarang
	rFixedNegative := PriceSchemaRule{
		ItemID:                &item,
		TariffAdjustmentType:  AdjustmentTypeFixedPrice,
		TariffAdjustmentValue: decimal.NewFromInt(-1000),
		PayByPatientPercent:   decimal.Zero,
		PayByCustomerPercent:  decimal.NewFromInt(100),
	}
	if err := rFixedNegative.Validate(); !errors.Is(err, ErrFixedPriceNegative) {
		t.Errorf("expected ErrFixedPriceNegative, got %v", err)
	}

	// 4. Cost sharing sum != 100 dilarang
	rCostShare := PriceSchemaRule{
		TariffAdjustmentType: AdjustmentTypeNone,
		PayByPatientPercent:  decimal.NewFromInt(30),
		PayByCustomerPercent: decimal.NewFromInt(50), // 80 != 100
	}
	if err := rCostShare.Validate(); !errors.Is(err, ErrCostShareSumMismatch) {
		t.Errorf("expected ErrCostShareSumMismatch, got %v", err)
	}
}

func TestDeleteRule_CannotDeleteDefaultRule(t *testing.T) {
	repo := newMockRepository()
	tx := &mockTxManager{}
	auditW := &mockAuditWriter{}

	delUC := NewDeletePriceSchemaRuleUseCase(repo, tx, auditW)
	ctx := context.Background()
	actor := OperatorContext{UserID: "usr-1", Username: "admin", Role: "ADMIN"}

	// Buat schema default rule
	defaultRule := &PriceSchemaRule{
		ID:                   "default-rule-id",
		PriceSchemaID:        "schema-1",
		TariffAdjustmentType: AdjustmentTypeNone,
	}
	_ = repo.CreateRule(ctx, defaultRule)

	err := delUC.Execute(ctx, "schema-1", "default-rule-id", actor)
	if !errors.Is(err, ErrCannotDeleteDefaultRule) {
		t.Fatalf("expected ErrCannotDeleteDefaultRule, got %v", err)
	}
}

func TestHierarchyOverridePriority(t *testing.T) {
	repo := newMockRepository()
	org := &mockOrgReader{}
	matcher := NewRuleMatcher(repo, org)
	ctx := context.Background()

	schemaID := "schema-hierarchy-test"
	itemID := "item-darah"
	unitID := "unit-lab"
	deptID := "dept-penunjang"

	// Masukkan rule di semua level:
	// 1. ITEM_UNIT: Diskon 50%
	// 2. ITEM_DEPT: Diskon 40%
	// 3. ITEM: Diskon 30%
	// 4. UNIT: Diskon 20%
	// 5. DEPT: Diskon 10%
	// 6. SCHEMA_DEFAULT: Diskon 5%
	rules := []PriceSchemaRule{
		{
			ID:                    "r-item-unit",
			PriceSchemaID:         schemaID,
			ItemID:                &itemID,
			ServiceUnitID:         &unitID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-50),
		},
		{
			ID:                    "r-item-dept",
			PriceSchemaID:         schemaID,
			ItemID:                &itemID,
			DepartementID:         &deptID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-40),
		},
		{
			ID:                    "r-item",
			PriceSchemaID:         schemaID,
			ItemID:                &itemID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-30),
		},
		{
			ID:                    "r-unit",
			PriceSchemaID:         schemaID,
			ServiceUnitID:         &unitID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-20),
		},
		{
			ID:                    "r-dept",
			PriceSchemaID:         schemaID,
			DepartementID:         &deptID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-10),
		},
		{
			ID:                    "r-default",
			PriceSchemaID:         schemaID,
			TariffAdjustmentType:  AdjustmentTypePercent,
			TariffAdjustmentValue: decimal.NewFromInt(-5),
		},
	}

	for _, r := range rules {
		_ = repo.CreateRule(ctx, &r)
	}

	// Kasus 1: Semua ada (item + unit + dept) -> Harus menang ITEM_UNIT (1)
	m1, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m1.Level != MatchLevelItemUnit {
		t.Fatalf("expected ITEM_UNIT, got %v (err: %v)", m1.Level, err)
	}

	// Hapus ITEM_UNIT -> Sekarang harus menang ITEM_DEPT (2)
	_ = repo.DeleteRule(ctx, "r-item-unit")
	m2, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m2.Level != MatchLevelItemDept {
		t.Fatalf("expected ITEM_DEPT, got %v (err: %v)", m2.Level, err)
	}

	// Hapus ITEM_DEPT -> Sekarang harus menang ITEM (3)
	_ = repo.DeleteRule(ctx, "r-item-dept")
	m3, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m3.Level != MatchLevelItem {
		t.Fatalf("expected ITEM, got %v (err: %v)", m3.Level, err)
	}

	// Hapus ITEM -> Sekarang harus menang UNIT (4)
	_ = repo.DeleteRule(ctx, "r-item")
	m4, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m4.Level != MatchLevelUnit {
		t.Fatalf("expected UNIT, got %v (err: %v)", m4.Level, err)
	}

	// Hapus UNIT -> Sekarang harus menang DEPT (5)
	_ = repo.DeleteRule(ctx, "r-unit")
	m5, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m5.Level != MatchLevelDept {
		t.Fatalf("expected DEPT, got %v (err: %v)", m5.Level, err)
	}

	// Hapus DEPT -> Fallback ke SCHEMA_DEFAULT (6)
	_ = repo.DeleteRule(ctx, "r-dept")
	m6, err := matcher.FindMatchingRule(ctx, MatchRuleQuery{
		PriceSchemaID: schemaID,
		ItemID:        itemID,
		ServiceUnitID: unitID,
		DepartementID: deptID,
	})
	if err != nil || m6.Level != MatchLevelSchemaDefault {
		t.Fatalf("expected SCHEMA_DEFAULT, got %v (err: %v)", m6.Level, err)
	}
}

// Mock TariffBaseLookup
type mockTariffLookup struct{}

func (m *mockTariffLookup) LookupTariffBase(ctx context.Context, pricePlanID *string, tariffClassID, itemID string, serviceDate time.Time, isCito bool) (*TariffBaseResult, error) {
	coaSarana := "410.01"
	coaDokter := "214.01"
	planID := "plan-standard"
	planCode := "STANDARD_2026"
	if pricePlanID != nil {
		planID = *pricePlanID
		planCode = "CUSTOM_" + *pricePlanID
	}
	return &TariffBaseResult{
		PricePlanID:    planID,
		PricePlanCode:  planCode,
		TotalBasePrice: decimal.NewFromInt(200000),
		TotalCitoPrice: decimal.NewFromInt(250000),
		Components: []ComponentInput{
			{
				ComponentID:       "c-1",
				ComponentName:     "Jasa RS",
				ComponentType:     "SARANA",
				COACode:           coaSarana,
				BaseAmount:        decimal.NewFromInt(120000),
				CitoAmount:        decimal.NewFromInt(150000),
				IsHospitalRevenue: true,
			},
			{
				ComponentID:       "c-2",
				ComponentName:     "Jasa Dokter",
				ComponentType:     "MEDIS",
				COACode:           coaDokter,
				BaseAmount:        decimal.NewFromInt(80000),
				CitoAmount:        decimal.NewFromInt(100000),
				IsHospitalRevenue: false,
			},
		},
	}, nil
}

func TestSimulateTariffUseCase(t *testing.T) {
	repo := newMockRepository()
	org := &mockOrgReader{}
	matcher := NewRuleMatcher(repo, org)
	lookup := &mockTariffLookup{}

	uc := NewSimulateTariffUseCase(repo, lookup, matcher)
	ctx := context.Background()

	schemaID := "schema-sim"
	_ = repo.CreateSchema(ctx, &PriceSchema{
		ID:       schemaID,
		Code:     "SCHEMA_SIM",
		Name:     "Schema Sim",
		IsActive: true,
	})
	// Rule: Diskon 10% (Rp 180.000), Pajak 10% (Rp 18.000), Gross 198.000. Patient 20%, Customer 80%
	rule := &PriceSchemaRule{
		ID:                    "rule-sim",
		PriceSchemaID:         schemaID,
		TariffAdjustmentType:  AdjustmentTypePercent,
		TariffAdjustmentValue: decimal.NewFromInt(-10),
		TaxPercent:            decimal.NewFromInt(10),
		PayByPatientPercent:   decimal.NewFromInt(20),
		PayByCustomerPercent:  decimal.NewFromInt(80),
	}
	_ = repo.CreateRule(ctx, rule)

	req := SimulateTariffRequest{
		PriceSchemaID: &schemaID,
		TariffClassID: "class-1",
		ItemID:        "item-konsul",
		IsCito:        false,
	}

	res, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.BasePrice.Equal(decimal.NewFromInt(200000)) {
		t.Errorf("expected base 200000, got %v", res.BasePrice)
	}
	if !res.FinalPrice.Equal(decimal.NewFromInt(180000)) {
		t.Errorf("expected final 180000, got %v", res.FinalPrice)
	}
	if !res.TaxAmount.Equal(decimal.NewFromInt(18000)) {
		t.Errorf("expected tax 18000, got %v", res.TaxAmount)
	}
	if !res.GrossAmount.Equal(decimal.NewFromInt(198000)) {
		t.Errorf("expected gross 198000, got %v", res.GrossAmount)
	}
	// Patient = 198.000 * 0.2 = 39.600
	if !res.PatientAmount.Equal(decimal.NewFromInt(39600)) {
		t.Errorf("expected patient 39600, got %v", res.PatientAmount)
	}
	// Customer = 198.000 - 39.600 = 158.400
	if !res.CustomerAmount.Equal(decimal.NewFromInt(158400)) {
		t.Errorf("expected customer 158400, got %v", res.CustomerAmount)
	}

	// Uji simulasi gagal jika schema dinonaktifkan
	simSchema, _ := repo.FindSchemaByID(ctx, schemaID)
	simSchema.IsActive = false
	_ = repo.UpdateSchema(ctx, simSchema)
	_, err = uc.Execute(ctx, req)
	if !errors.Is(err, ErrSchemaInactive) {
		t.Fatalf("expected ErrSchemaInactive when schema is deactivated, got %v", err)
	}
}

func TestCreatePriceSchema_ValidateCrossModuleReferences(t *testing.T) {
	repo := newMockRepository()
	tx := &mockTxManager{}
	auditW := &mockAuditWriter{}
	org := &mockOrgReader{}
	itemR := &mockItemReader{}
	planR := &mockPlanReader{}

	uc := NewCreatePriceSchemaUseCase(repo, tx, auditW, org, org, itemR, planR)
	ctx := context.Background()
	actor := OperatorContext{UserID: "usr-1", Username: "admin", Role: "ADMIN"}

	invalidDept := "dept-invalid"
	req := CreatePriceSchemaRequest{
		Code: "SCHEMA_INVALID_DEPT",
		Name: "Test Invalid Dept",
		Rules: []UpsertRuleRequest{
			{
				DepartementID:         &invalidDept,
				TariffAdjustmentType:  AdjustmentTypePercent,
				TariffAdjustmentValue: decimal.NewFromInt(10),
				PayByPatientPercent:   decimal.Zero,
				PayByCustomerPercent:  decimal.NewFromInt(100),
			},
		},
	}

	_, err := uc.Execute(ctx, req, actor)
	if !errors.Is(err, ErrDepartmentInactiveOrNotFound) {
		t.Fatalf("expected ErrDepartmentInactiveOrNotFound, got %v", err)
	}
}

func TestSimulateTariffUseCase_InvalidSchemaReturnsError(t *testing.T) {
	repo := newMockRepository()
	org := &mockOrgReader{}
	matcher := NewRuleMatcher(repo, org)
	lookup := &mockTariffLookup{}

	uc := NewSimulateTariffUseCase(repo, lookup, matcher)
	ctx := context.Background()

	invalidSchemaID := "schema-non-existent"
	req := SimulateTariffRequest{
		PriceSchemaID: &invalidSchemaID,
		TariffClassID: "class-1",
		ItemID:        "item-konsul",
	}

	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Fatalf("expected error for non-existent schema, got nil")
	}
}

func TestCreatePriceSchema_DuplicateScopeAndMultipleDefaults(t *testing.T) {
	repo := newMockRepository()
	tx := &mockTxManager{}
	auditW := &mockAuditWriter{}
	org := &mockOrgReader{}
	itemR := &mockItemReader{}
	planR := &mockPlanReader{}

	uc := NewCreatePriceSchemaUseCase(repo, tx, auditW, org, org, itemR, planR)
	ctx := context.Background()
	actor := OperatorContext{UserID: "usr-1", Username: "admin", Role: "ADMIN"}

	// 1. Multiple default rules
	reqMultipleDefaults := CreatePriceSchemaRequest{
		Code: "SCHEMA_MULTI_DEFAULT",
		Name: "Test Multi Defaults",
		Rules: []UpsertRuleRequest{
			{
				TariffAdjustmentType: AdjustmentTypeNone,
				PayByPatientPercent:  decimal.Zero,
				PayByCustomerPercent: decimal.NewFromInt(100),
			},
			{
				TariffAdjustmentType: AdjustmentTypeNone,
				PayByPatientPercent:  decimal.Zero,
				PayByCustomerPercent: decimal.NewFromInt(100),
			},
		},
	}
	_, err := uc.Execute(ctx, reqMultipleDefaults, actor)
	if !errors.Is(err, ErrMultipleDefaultRules) {
		t.Fatalf("expected ErrMultipleDefaultRules, got %v", err)
	}

	// 2. Duplicate scope
	itemValid := "item-valid"
	reqDuplicateScope := CreatePriceSchemaRequest{
		Code: "SCHEMA_DUP_SCOPE",
		Name: "Test Dup Scope",
		Rules: []UpsertRuleRequest{
			{
				ItemID:                &itemValid,
				TariffAdjustmentType:  AdjustmentTypePercent,
				TariffAdjustmentValue: decimal.NewFromInt(10),
				PayByPatientPercent:   decimal.Zero,
				PayByCustomerPercent:  decimal.NewFromInt(100),
			},
			{
				ItemID:                &itemValid,
				TariffAdjustmentType:  AdjustmentTypePercent,
				TariffAdjustmentValue: decimal.NewFromInt(20),
				PayByPatientPercent:   decimal.Zero,
				PayByCustomerPercent:  decimal.NewFromInt(100),
			},
		},
	}
	_, err = uc.Execute(ctx, reqDuplicateScope, actor)
	if !errors.Is(err, ErrDuplicateRuleScope) {
		t.Fatalf("expected ErrDuplicateRuleScope, got %v", err)
	}
}

func TestPriceSchema_BasePricePlanIntegration(t *testing.T) {
	repo := newMockRepository()
	tx := &mockTxManager{}
	auditW := &mockAuditWriter{}
	org := &mockOrgReader{}
	itemR := &mockItemReader{}
	planR := &mockPlanReader{}

	createUC := NewCreatePriceSchemaUseCase(repo, tx, auditW, org, org, itemR, planR)
	ctx := context.Background()
	actor := OperatorContext{UserID: "usr-1", Username: "admin", Role: "ADMIN"}

	// 1. Gagal jika base_price_plan_id tidak valid
	invalidPlanID := "plan-invalid"
	reqInvalid := CreatePriceSchemaRequest{
		Code:            "SCHEMA_BPJS_INVALID",
		Name:            "Skema BPJS Plan Invalid",
		BasePricePlanID: &invalidPlanID,
	}
	_, err := createUC.Execute(ctx, reqInvalid, actor)
	if !errors.Is(err, ErrPricePlanNotFound) {
		t.Fatalf("expected ErrPricePlanNotFound, got %v", err)
	}

	// 2. Sukses jika base_price_plan_id valid
	validPlanID := "plan-bpjs-2024"
	reqValid := CreatePriceSchemaRequest{
		Code:            "SCHEMA_BPJS_2024",
		Name:            "Skema BPJS TK 2024",
		BasePricePlanID: &validPlanID,
	}
	res, err := createUC.Execute(ctx, reqValid, actor)
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}
	if res.BasePricePlanID == nil || *res.BasePricePlanID != validPlanID {
		t.Fatalf("expected base_price_plan_id %s, got %v", validPlanID, res.BasePricePlanID)
	}

	// 3. Simulasi otomatis menggunakan base_price_plan_id dari skema
	matcher := NewRuleMatcher(repo, org)
	lookup := &mockTariffLookup{}
	simUC := NewSimulateTariffUseCase(repo, lookup, matcher)

	simReq := SimulateTariffRequest{
		PriceSchemaID: &res.ID,
		TariffClassID: "class-1",
		ItemID:        "item-konsul",
	}
	simRes, err := simUC.Execute(ctx, simReq)
	if err != nil {
		t.Fatalf("unexpected error simulating tariff: %v", err)
	}

	if simRes.PricePlanID != validPlanID {
		t.Fatalf("expected simulated price_plan_id %s, got %s", validPlanID, simRes.PricePlanID)
	}

	// 4. Update Price Schema dengan base_price_plan_id invalid
	updateUC := NewUpdatePriceSchemaUseCase(repo, tx, auditW, planR)
	invalidUpdatePlanID := "plan-invalid-update"
	updateReqInvalid := UpdatePriceSchemaRequest{
		Name:            "Skema Update Invalid",
		BasePricePlanID: &invalidUpdatePlanID,
	}
	_, err = updateUC.Execute(ctx, res.ID, updateReqInvalid, actor)
	if !errors.Is(err, ErrPricePlanNotFound) {
		t.Fatalf("expected ErrPricePlanNotFound on update with invalid plan, got %v", err)
	}

	// 5. Update Price Schema dengan base_price_plan_id valid
	anotherValidPlanID := "plan-bpjs-2024"
	updateReqValid := UpdatePriceSchemaRequest{
		Name:            "Skema Update Valid",
		BasePricePlanID: &anotherValidPlanID,
	}
	updateRes, err := updateUC.Execute(ctx, res.ID, updateReqValid, actor)
	if err != nil {
		t.Fatalf("unexpected error updating schema: %v", err)
	}
	if updateRes.BasePricePlanID == nil || *updateRes.BasePricePlanID != anotherValidPlanID {
		t.Fatalf("expected base_price_plan_id %s, got %v", anotherValidPlanID, updateRes.BasePricePlanID)
	}

	// 6. Update Price Schema dengan string kosong (menghapus base_price_plan_id)
	emptyPlanID := "   "
	updateReqClear := UpdatePriceSchemaRequest{
		Name:            "Skema Update Clear",
		BasePricePlanID: &emptyPlanID,
	}
	clearRes, err := updateUC.Execute(ctx, res.ID, updateReqClear, actor)
	if err != nil {
		t.Fatalf("unexpected error clearing schema base price plan: %v", err)
	}
	if clearRes.BasePricePlanID != nil {
		t.Fatalf("expected base_price_plan_id to be cleared, got %v", clearRes.BasePricePlanID)
	}
}

