package priceschema

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreatePriceSchemaRequest adalah payload request untuk membuat skema harga baru.
type CreatePriceSchemaRequest struct {
	Code        string              `json:"code" binding:"required"`
	Name        string              `json:"name" binding:"required"`
	Description *string             `json:"description"`
	Rules       []UpsertRuleRequest `json:"rules"`
}

// UpdatePriceSchemaRequest adalah payload request untuk memperbarui informasi skema harga.
type UpdatePriceSchemaRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// UpsertRuleRequest adalah payload request untuk menambah atau memperbarui aturan skema harga.
type UpsertRuleRequest struct {
	DepartementID             *string         `json:"departement_id"`
	ServiceUnitID             *string         `json:"service_unit_id"`
	ItemID                    *string         `json:"item_id"`
	TariffAdjustmentType      AdjustmentType  `json:"tariff_adjustment_type" binding:"required"`
	TariffAdjustmentValue     decimal.Decimal `json:"tariff_adjustment_value"`
	MedicationAdjustmentType  AdjustmentType  `json:"medication_adjustment_type"`
	MedicationAdjustmentValue decimal.Decimal `json:"medication_adjustment_value"`
	PayByPatientPercent       decimal.Decimal `json:"pay_by_patient_percent"`
	PayByCustomerPercent      decimal.Decimal `json:"pay_by_customer_percent"`
	TaxPercent                decimal.Decimal `json:"tax_percent"`
}

// SimulateTariffRequest adalah payload input untuk simulasi kalkulasi tarif (PRD §5.1, §10).
type SimulateTariffRequest struct {
	PriceSchemaID *string `json:"price_schema_id"`
	TariffClassID string  `json:"tariff_class_id" binding:"required"`
	ItemID        string  `json:"item_id" binding:"required"`
	DepartementID *string `json:"departement_id"`
	ServiceUnitID *string `json:"service_unit_id"`
	ServiceDate   string  `json:"service_date"` // YYYY-MM-DD atau RFC3339
	IsCito        bool    `json:"is_cito"`
}

// PriceSchemaResponse adalah representasi DTO output skema harga.
type PriceSchemaResponse struct {
	ID          string                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Description *string                   `json:"description,omitempty"`
	IsActive    bool                      `json:"is_active"`
	RulesCount  int                       `json:"rules_count,omitempty"`
	Rules       []PriceSchemaRuleResponse `json:"rules,omitempty"`
	CreatedAt   string                    `json:"created_at"`
	UpdatedAt   string                    `json:"updated_at"`
	CreatedBy   string                    `json:"created_by"`
	UpdatedBy   string                    `json:"updated_by"`
}

// PriceSchemaRuleResponse adalah representasi DTO output rule skema harga.
type PriceSchemaRuleResponse struct {
	ID                        string          `json:"id"`
	PriceSchemaID             string          `json:"price_schema_id"`
	Level                     MatchLevel      `json:"level"`
	DepartementID             *string         `json:"departement_id,omitempty"`
	ServiceUnitID             *string         `json:"service_unit_id,omitempty"`
	ItemID                    *string         `json:"item_id,omitempty"`
	TariffAdjustmentType      AdjustmentType  `json:"tariff_adjustment_type"`
	TariffAdjustmentValue     decimal.Decimal `json:"tariff_adjustment_value"`
	MedicationAdjustmentType  AdjustmentType  `json:"medication_adjustment_type"`
	MedicationAdjustmentValue decimal.Decimal `json:"medication_adjustment_value"`
	PayByPatientPercent       decimal.Decimal `json:"pay_by_patient_percent"`
	PayByCustomerPercent      decimal.Decimal `json:"pay_by_customer_percent"`
	TaxPercent                decimal.Decimal `json:"tax_percent"`
	CreatedAt                 string          `json:"created_at"`
	UpdatedAt                 string          `json:"updated_at"`
	CreatedBy                 string          `json:"created_by"`
	UpdatedBy                 string          `json:"updated_by"`
}

// SimulateTariffResponse adalah rincian lengkap hasil simulasi kalkulasi tarif.
type SimulateTariffResponse struct {
	PricePlanID           string            `json:"price_plan_id"`
	PricePlanCode         string            `json:"price_plan_code"`
	BasePrice             decimal.Decimal   `json:"base_price"`
	CitoSurcharge         decimal.Decimal   `json:"cito_surcharge"`
	IsCito                bool              `json:"is_cito"`
	PriceSchemaID         *string           `json:"price_schema_id,omitempty"`
	PriceSchemaRuleID     *string           `json:"price_schema_rule_id,omitempty"`
	MatchLevel            *MatchLevel       `json:"match_level,omitempty"`
	TariffAdjustmentType  AdjustmentType    `json:"tariff_adjustment_type"`
	TariffAdjustmentValue decimal.Decimal   `json:"tariff_adjustment_value"`
	FinalPrice            decimal.Decimal   `json:"final_price"`
	TaxPercent            decimal.Decimal   `json:"tax_percent"`
	TaxAmount             decimal.Decimal   `json:"tax_amount"`
	GrossAmount           decimal.Decimal   `json:"gross_amount"`
	PayByPatientPercent   decimal.Decimal   `json:"pay_by_patient_percent"`
	PayByCustomerPercent  decimal.Decimal   `json:"pay_by_customer_percent"`
	PatientAmount         decimal.Decimal   `json:"patient_amount"`
	CustomerAmount        decimal.Decimal   `json:"customer_amount"`
	Components            []ComponentOutput `json:"components"`
}

func ToPriceSchemaResponse(s *PriceSchema) PriceSchemaResponse {
	var rulesResp []PriceSchemaRuleResponse
	if len(s.Rules) > 0 {
		rulesResp = make([]PriceSchemaRuleResponse, len(s.Rules))
		for i, r := range s.Rules {
			rulesResp[i] = ToPriceSchemaRuleResponse(&r)
		}
	}

	return PriceSchemaResponse{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Description: s.Description,
		IsActive:    s.IsActive,
		RulesCount:  len(s.Rules),
		Rules:       rulesResp,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
		CreatedBy:   s.CreatedBy,
		UpdatedBy:   s.UpdatedBy,
	}
}

func ToPriceSchemaRuleResponse(r *PriceSchemaRule) PriceSchemaRuleResponse {
	return PriceSchemaRuleResponse{
		ID:                        r.ID,
		PriceSchemaID:             r.PriceSchemaID,
		Level:                     r.Level(),
		DepartementID:             r.DepartementID,
		ServiceUnitID:             r.ServiceUnitID,
		ItemID:                    r.ItemID,
		TariffAdjustmentType:      r.TariffAdjustmentType,
		TariffAdjustmentValue:     r.TariffAdjustmentValue,
		MedicationAdjustmentType:  r.MedicationAdjustmentType,
		MedicationAdjustmentValue: r.MedicationAdjustmentValue,
		PayByPatientPercent:       r.PayByPatientPercent,
		PayByCustomerPercent:      r.PayByCustomerPercent,
		TaxPercent:                r.TaxPercent,
		CreatedAt:                 r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                 r.UpdatedAt.Format(time.RFC3339),
		CreatedBy:                 r.CreatedBy,
		UpdatedBy:                 r.UpdatedBy,
	}
}
