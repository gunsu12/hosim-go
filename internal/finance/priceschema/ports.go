package priceschema

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// MatchedRule membawa data rule yang terpilih oleh RuleMatcher.
type MatchedRule struct {
	RuleID                string          `json:"rule_id"`
	PriceSchemaID         string          `json:"price_schema_id"`
	Level                 MatchLevel      `json:"level"`
	TariffAdjustmentType  AdjustmentType  `json:"tariff_adjustment_type"`
	TariffAdjustmentValue decimal.Decimal `json:"tariff_adjustment_value"`
	PayByPatientPercent   decimal.Decimal `json:"pay_by_patient_percent"`
	PayByCustomerPercent  decimal.Decimal `json:"pay_by_customer_percent"`
	TaxPercent            decimal.Decimal `json:"tax_percent"`
}

// MatchRuleQuery adalah kriteria input pencarian rule tarif.
type MatchRuleQuery struct {
	PriceSchemaID string
	ItemID        string
	DepartementID string
	ServiceUnitID string
}

// RuleMatcher adalah read contract yang diekspos priceschema ke modul lain (billing, registrasi, dsb).
type RuleMatcher interface {
	// FindMatchingRule mengembalikan rule paling spesifik sesuai urutan override (PRD §6.1).
	// Mengembalikan ErrNoMatchingRule bila schema tidak memiliki rule yang cocok.
	FindMatchingRule(ctx context.Context, q MatchRuleQuery) (MatchedRule, error)
}

// DepartmentReader adalah port verifikasi departemen ke modul organization.
type DepartmentReader interface {
	IsActiveDepartment(ctx context.Context, id string) (bool, error)
}

// ServiceUnitReader adalah port verifikasi unit layanan dan perolehan departemen unit.
type ServiceUnitReader interface {
	GetServiceUnitInfo(ctx context.Context, id string) (deptID *string, isActive bool, err error)
}

// ItemReader adalah port verifikasi keberadaan dan status aktif item di modul catalog.
type ItemReader interface {
	IsActiveItem(ctx context.Context, id string) (bool, error)
}

// TariffBaseLookup adalah port lookup tarif dasar ke modul priceplan.
type TariffBaseLookup interface {
	LookupTariffBase(ctx context.Context, tariffClassID, itemID string, serviceDate time.Time, isCito bool) (*TariffBaseResult, error)
}

// TariffBaseResult adalah data tarif dasar hasil lookup dari Price Plan.
type TariffBaseResult struct {
	PricePlanID    string
	PricePlanCode  string
	TotalBasePrice decimal.Decimal
	TotalCitoPrice decimal.Decimal
	Components     []ComponentInput
}

// ListSchemaParams memuat parameter pencarian dan paginasi price schemas.
type ListSchemaParams struct {
	Search   string
	IsActive *bool
	Page     int
	Limit    int
}

// ListRuleParams memuat parameter pencarian dan paginasi price schema rules.
type ListRuleParams struct {
	Level         string
	DepartementID string
	ServiceUnitID string
	ItemID        string
	Page          int
	Limit         int
}

// Repository adalah port persistence untuk entitas PriceSchema dan PriceSchemaRule.
type Repository interface {
	CreateSchema(ctx context.Context, schema *PriceSchema) error
	UpdateSchema(ctx context.Context, schema *PriceSchema) error
	FindSchemaByID(ctx context.Context, id string) (*PriceSchema, error)
	FindSchemaByCode(ctx context.Context, code string) (*PriceSchema, error)
	ListSchemas(ctx context.Context, params ListSchemaParams) ([]PriceSchema, int64, error)

	CreateRule(ctx context.Context, rule *PriceSchemaRule) error
	UpdateRule(ctx context.Context, rule *PriceSchemaRule) error
	DeleteRule(ctx context.Context, id string) error
	FindRuleByID(ctx context.Context, id string) (*PriceSchemaRule, error)
	FindRuleByScope(ctx context.Context, schemaID string, deptID, unitID, itemID *string) (*PriceSchemaRule, error)
	ListRulesBySchemaID(ctx context.Context, schemaID string, params ListRuleParams) ([]PriceSchemaRule, int64, error)
	FindBestMatchingRule(ctx context.Context, schemaID, itemID, serviceUnitID, departementID string) (*PriceSchemaRule, error)
}
