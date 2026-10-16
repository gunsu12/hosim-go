package priceschema

import (
	"errors"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AdjustmentType mendefinisikan jenis penyesuaian tarif pada Price Schema.
type AdjustmentType string

const (
	AdjustmentTypeNone        AdjustmentType = "NONE"
	AdjustmentTypePercent     AdjustmentType = "PERCENT"
	AdjustmentTypeFixedAmount AdjustmentType = "FIXED_AMOUNT"
	AdjustmentTypeFixedPrice  AdjustmentType = "FIXED_PRICE"
)

func (a AdjustmentType) IsValid() bool {
	switch a {
	case AdjustmentTypeNone, AdjustmentTypePercent, AdjustmentTypeFixedAmount, AdjustmentTypeFixedPrice:
		return true
	default:
		return false
	}
}

// MatchLevel menunjukkan level cakupan rule yang cocok (override hierarchy).
type MatchLevel string

const (
	MatchLevelItemUnit      MatchLevel = "ITEM_UNIT"
	MatchLevelItemDept      MatchLevel = "ITEM_DEPT"
	MatchLevelItem          MatchLevel = "ITEM"
	MatchLevelUnit          MatchLevel = "UNIT"
	MatchLevelDept          MatchLevel = "DEPT"
	MatchLevelSchemaDefault MatchLevel = "SCHEMA_DEFAULT"
)

var (
	ErrSchemaNotFound                = errors.New("price schema tidak ditemukan")
	ErrSchemaInactive                = errors.New("skema harga sedang dinonaktifkan")
	ErrSchemaCodeExists              = errors.New("kode price schema sudah terdaftar")
	ErrSchemaCodeRequired            = errors.New("kode price schema wajib diisi")
	ErrSchemaNameRequired            = errors.New("nama price schema wajib diisi")
	ErrRuleNotFound                  = errors.New("rule price schema tidak ditemukan")
	ErrNoMatchingRule                = errors.New("tidak ada rule price schema yang cocok")
	ErrCannotDeleteDefaultRule       = errors.New("rule SCHEMA_DEFAULT tidak boleh dihapus")
	ErrScopeExclusive                = errors.New("service_unit_id dan departement_id tidak boleh diisi bersamaan")
	ErrFixedPriceRequiresItem        = errors.New("FIXED_PRICE hanya diperbolehkan pada rule yang memiliki item_id")
	ErrFixedPriceNegative            = errors.New("nilai FIXED_PRICE tidak boleh negatif")
	ErrPercentMinimum                = errors.New("nilai PERCENT tidak boleh kurang dari -100%")
	ErrNoneValueMustBeZero           = errors.New("nilai adjustment bertipe NONE harus 0")
	ErrCostShareSumMismatch          = errors.New("total pay_by_patient_percent dan pay_by_customer_percent wajib 100%")
	ErrCostShareOutOfRange           = errors.New("persentase cost sharing harus bernilai antara 0% dan 100%")
	ErrTaxOutOfRange                 = errors.New("persentase pajak harus bernilai antara 0% dan 100%")
	ErrInvalidAdjustmentType         = errors.New("tipe penyesuaian tarif tidak valid")
	ErrInvalidMedicationAdjustType   = errors.New("tipe penyesuaian obat tidak valid")
	ErrDepartmentInactiveOrNotFound  = errors.New("departemen tidak ditemukan atau tidak aktif")
	ErrServiceUnitInactiveOrNotFound = errors.New("unit layanan tidak ditemukan atau tidak aktif")
	ErrItemInactiveOrNotFound        = errors.New("item/tindakan tidak ditemukan atau tidak aktif")
	ErrTariffNotConfigured           = errors.New("tarif dasar tindakan tidak ditemukan di buku tarif RS")
	ErrPricePlanNotFound             = errors.New("buku tarif acuan tidak ditemukan")
	ErrDuplicateRuleScope            = errors.New("terdapat aturan dengan cakupan yang sama dalam satu skema")
	ErrMultipleDefaultRules          = errors.New("skema harga hanya boleh memiliki satu aturan default (SCHEMA_DEFAULT)")
)

// PriceSchema merepresentasikan header skema penyesuaian harga PKS.
type PriceSchema struct {
	ID                string    `gorm:"primaryKey;size:36" json:"id"`
	Code              string    `gorm:"uniqueIndex:idx_price_schemas_code;size:50;not null" json:"code"`
	Name              string    `gorm:"size:150;not null" json:"name"`
	Description       *string   `gorm:"type:text" json:"description,omitempty"`
	BasePricePlanID   *string   `gorm:"size:36;index:idx_price_schemas_base_price_plan_id" json:"base_price_plan_id,omitempty"`
	BasePricePlanCode *string   `gorm:"->" json:"base_price_plan_code,omitempty"`
	BasePricePlanName *string   `gorm:"->" json:"base_price_plan_name,omitempty"`
	IsActive          bool      `gorm:"default:true;not null;index:idx_price_schemas_is_active" json:"is_active"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy         string    `gorm:"size:50;not null;default:'SYSTEM'" json:"created_by"`
	UpdatedBy         string    `gorm:"size:50;not null;default:'SYSTEM'" json:"updated_by"`

	Rules []PriceSchemaRule `gorm:"foreignKey:PriceSchemaID;constraint:OnDelete:CASCADE" json:"rules,omitempty"`
}

func (PriceSchema) TableName() string {
	return "price_schemas"
}

func (s *PriceSchema) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewV7().String()
	}
	return nil
}

// PriceSchemaRule merepresentasikan aturan bertingkat penyesuaian tarif.
type PriceSchemaRule struct {
	ID                        string          `gorm:"primaryKey;size:36" json:"id"`
	PriceSchemaID             string          `gorm:"size:36;not null;index" json:"price_schema_id"`
	DepartementID             *string         `gorm:"size:36;index:idx_psr_departement" json:"departement_id,omitempty"`
	ServiceUnitID             *string         `gorm:"size:36;index:idx_psr_service_unit" json:"service_unit_id,omitempty"`
	ItemID                    *string         `gorm:"size:36;index:idx_psr_lookup_item" json:"item_id,omitempty"`
	TariffAdjustmentType      AdjustmentType  `gorm:"size:20;not null;default:'NONE'" json:"tariff_adjustment_type"`
	TariffAdjustmentValue     decimal.Decimal `gorm:"type:numeric(15,4);not null;default:0" json:"tariff_adjustment_value"`
	MedicationAdjustmentType  AdjustmentType  `gorm:"size:20;not null;default:'NONE'" json:"medication_adjustment_type"`
	MedicationAdjustmentValue decimal.Decimal `gorm:"type:numeric(15,4);not null;default:0" json:"medication_adjustment_value"`
	PayByPatientPercent       decimal.Decimal `gorm:"type:numeric(5,2);not null;default:0" json:"pay_by_patient_percent"`
	PayByCustomerPercent      decimal.Decimal `gorm:"type:numeric(5,2);not null;default:100" json:"pay_by_customer_percent"`
	TaxPercent                decimal.Decimal `gorm:"type:numeric(5,2);not null;default:0" json:"tax_percent"`
	CreatedAt                 time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                 time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy                 string          `gorm:"size:50;not null;default:'SYSTEM'" json:"created_by"`
	UpdatedBy                 string          `gorm:"size:50;not null;default:'SYSTEM'" json:"updated_by"`
}

func (PriceSchemaRule) TableName() string {
	return "price_schema_rules"
}

func (r *PriceSchemaRule) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewV7().String()
	}
	return nil
}

// Level menentukan level hirarki aturan berdasarkan kombinasi scope.
func (r *PriceSchemaRule) Level() MatchLevel {
	hasItem := r.ItemID != nil && *r.ItemID != ""
	hasUnit := r.ServiceUnitID != nil && *r.ServiceUnitID != ""
	hasDept := r.DepartementID != nil && *r.DepartementID != ""

	if hasItem && hasUnit {
		return MatchLevelItemUnit
	}
	if hasItem && hasDept {
		return MatchLevelItemDept
	}
	if hasItem {
		return MatchLevelItem
	}
	if hasUnit {
		return MatchLevelUnit
	}
	if hasDept {
		return MatchLevelDept
	}
	return MatchLevelSchemaDefault
}

// Validate memeriksa invariant bisnis pada PriceSchemaRule.
func (r *PriceSchemaRule) Validate() error {
	hasDept := r.DepartementID != nil && *r.DepartementID != ""
	hasUnit := r.ServiceUnitID != nil && *r.ServiceUnitID != ""
	hasItem := r.ItemID != nil && *r.ItemID != ""

	// Invariant: Departement dan ServiceUnit tidak boleh diisi bersamaan
	if hasDept && hasUnit {
		return ErrScopeExclusive
	}

	// Validasi tipe penyesuaian tarif
	if !r.TariffAdjustmentType.IsValid() {
		return ErrInvalidAdjustmentType
	}

	zero := decimal.Zero
	minusHundred := decimal.NewFromInt(-100)
	hundred := decimal.NewFromInt(100)

	switch r.TariffAdjustmentType {
	case AdjustmentTypeNone:
		if !r.TariffAdjustmentValue.IsZero() {
			return ErrNoneValueMustBeZero
		}
	case AdjustmentTypePercent:
		if r.TariffAdjustmentValue.LessThan(minusHundred) {
			return ErrPercentMinimum
		}
	case AdjustmentTypeFixedAmount:
		// Nominal bebas (bisa positif atau negatif potongan)
	case AdjustmentTypeFixedPrice:
		if !hasItem {
			return ErrFixedPriceRequiresItem
		}
		if r.TariffAdjustmentValue.LessThan(zero) {
			return ErrFixedPriceNegative
		}
	}

	// Validasi tipe penyesuaian obat (fase 2 preparation)
	if r.MedicationAdjustmentType != "" {
		if !r.MedicationAdjustmentType.IsValid() {
			return ErrInvalidMedicationAdjustType
		}
		switch r.MedicationAdjustmentType {
		case AdjustmentTypeNone:
			if !r.MedicationAdjustmentValue.IsZero() {
				return ErrNoneValueMustBeZero
			}
		case AdjustmentTypePercent:
			if r.MedicationAdjustmentValue.LessThan(minusHundred) {
				return ErrPercentMinimum
			}
		case AdjustmentTypeFixedPrice:
			if !hasItem {
				return ErrFixedPriceRequiresItem
			}
			if r.MedicationAdjustmentValue.LessThan(zero) {
				return ErrFixedPriceNegative
			}
		}
	}

	// Invariant cost sharing
	if r.PayByPatientPercent.LessThan(zero) || r.PayByPatientPercent.GreaterThan(hundred) ||
		r.PayByCustomerPercent.LessThan(zero) || r.PayByCustomerPercent.GreaterThan(hundred) {
		return ErrCostShareOutOfRange
	}
	if !r.PayByPatientPercent.Add(r.PayByCustomerPercent).Equal(hundred) {
		return ErrCostShareSumMismatch
	}

	// Invariant pajak
	if r.TaxPercent.LessThan(zero) || r.TaxPercent.GreaterThan(hundred) {
		return ErrTaxOutOfRange
	}

	return nil
}
