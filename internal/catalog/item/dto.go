package item

import (
	"time"

	"hosim-go/pkg/enums"
)

// ==========================================
// Request DTOs
// ==========================================

// CreateMedicationRequest mendefinisikan payload untuk registrasi master obat
type CreateMedicationRequest struct {
	Code               string               `json:"code" binding:"required"`
	Name               string               `json:"name" binding:"required"`
	GenericName        *string              `json:"generic_name"`
	CategoryID         string               `json:"category_id" binding:"required"`
	ProductLineID      *string              `json:"product_line_id"`
	BaseUnitName       string               `json:"base_unit_name" binding:"required"`
	KFACode            *string              `json:"kfa_code"`
	BPOMNIE            *string              `json:"bpom_nie"`
	DosageForm         *string              `json:"dosage_form"`
	StrengthAmount     *string              `json:"strength_amount"`
	StrengthUnit       *string              `json:"strength_unit"`
	DefaultRoute       *string              `json:"default_route"`
	MedicationType     enums.MedicationType `json:"medication_type"`
	IsHighAlert        bool                 `json:"is_high_alert"`
	IsLASA             bool                 `json:"is_lasa"`
	IsFornas           bool                 `json:"is_fornas"`
	IsAntibiotic       bool                 `json:"is_antibiotic"`
	StorageTemperature *string              `json:"storage_temperature"`
}

// CreateGeneralRequest mendefinisikan payload untuk registrasi barang umum/BMHP/alat
type CreateGeneralRequest struct {
	Code                string                     `json:"code" binding:"required"`
	Name                string                     `json:"name" binding:"required"`
	GenericName         *string                    `json:"generic_name"`
	CategoryID          string                     `json:"category_id" binding:"required"`
	ProductLineID       *string                    `json:"product_line_id"`
	BaseUnitName        string                     `json:"base_unit_name" binding:"required"`
	GeneralType         enums.GeneralType          `json:"general_type" binding:"required"`
	IsSterile           bool                       `json:"is_sterile"`
	IsDisposable        *bool                      `json:"is_disposable"`
	IsCSSDItem          bool                       `json:"is_cssd_item"`
	SterilizationMethod *enums.SterilizationMethod `json:"sterilization_method"`
}

// CreateAssetRequest mendefinisikan payload untuk registrasi barang modal / alkes
type CreateAssetRequest struct {
	Code                    string                    `json:"code" binding:"required"`
	Name                    string                    `json:"name" binding:"required"`
	GenericName             *string                   `json:"generic_name"`
	CategoryID              string                    `json:"category_id" binding:"required"`
	ProductLineID           *string                   `json:"product_line_id"`
	BaseUnitName            string                    `json:"base_unit_name" binding:"required"`
	Brand                   *string                   `json:"brand"`
	ModelName               *string                   `json:"model_name"`
	IsMedicalEquipment      bool                      `json:"is_medical_equipment"`
	ExpectedLifeYears       *int                      `json:"expected_life_years"`
	DepreciationMethod      *enums.DepreciationMethod `json:"depreciation_method"`
	MaintenanceIntervalDays *int                      `json:"maintenance_interval_days"`
}

// CreateTariffRequest mendefinisikan payload untuk registrasi master tarif layanan / jasa pasien
type CreateTariffRequest struct {
	Code         string           `json:"code" binding:"required"`
	Name         string           `json:"name" binding:"required"`
	CategoryID   string           `json:"category_id" binding:"required"`
	BaseUnitName string           `json:"base_unit_name"` // opsional: default "Kali" jika kosong
	NameAlias    *string          `json:"name_alias"`
	ChargeType   enums.ChargeType `json:"charge_type" binding:"required"`
	Notes        *string          `json:"notes"`
}

// UpdateMedicationRequest mendefinisikan payload pembaruan data master obat
type UpdateMedicationRequest struct {
	Code               string               `json:"code" binding:"required"`
	Name               string               `json:"name" binding:"required"`
	GenericName        *string              `json:"generic_name"`
	CategoryID         string               `json:"category_id" binding:"required"`
	ProductLineID      *string              `json:"product_line_id"`
	KFACode            *string              `json:"kfa_code"`
	BPOMNIE            *string              `json:"bpom_nie"`
	DosageForm         *string              `json:"dosage_form"`
	StrengthAmount     *string              `json:"strength_amount"`
	StrengthUnit       *string              `json:"strength_unit"`
	DefaultRoute       *string              `json:"default_route"`
	MedicationType     enums.MedicationType `json:"medication_type"`
	IsHighAlert        bool                 `json:"is_high_alert"`
	IsLASA             bool                 `json:"is_lasa"`
	IsFornas           bool                 `json:"is_fornas"`
	IsAntibiotic       bool                 `json:"is_antibiotic"`
	StorageTemperature *string              `json:"storage_temperature"`
	IsActive           bool                 `json:"is_active"`
}

// UpdateGeneralRequest mendefinisikan payload pembaruan data barang umum / BMHP
type UpdateGeneralRequest struct {
	Code                string                     `json:"code" binding:"required"`
	Name                string                     `json:"name" binding:"required"`
	GenericName         *string                    `json:"generic_name"`
	CategoryID          string                     `json:"category_id" binding:"required"`
	ProductLineID       *string                    `json:"product_line_id"`
	GeneralType         enums.GeneralType          `json:"general_type" binding:"required"`
	IsSterile           bool                       `json:"is_sterile"`
	IsDisposable        *bool                      `json:"is_disposable"`
	IsCSSDItem          bool                       `json:"is_cssd_item"`
	SterilizationMethod *enums.SterilizationMethod `json:"sterilization_method"`
	IsActive            bool                       `json:"is_active"`
}

// UpdateAssetRequest mendefinisikan payload pembaruan data aset modal
type UpdateAssetRequest struct {
	Code                    string                    `json:"code" binding:"required"`
	Name                    string                    `json:"name" binding:"required"`
	GenericName             *string                   `json:"generic_name"`
	CategoryID              string                    `json:"category_id" binding:"required"`
	ProductLineID           *string                   `json:"product_line_id"`
	Brand                   *string                   `json:"brand"`
	ModelName               *string                   `json:"model_name"`
	IsMedicalEquipment      bool                      `json:"is_medical_equipment"`
	ExpectedLifeYears       *int                      `json:"expected_life_years"`
	DepreciationMethod      *enums.DepreciationMethod `json:"depreciation_method"`
	MaintenanceIntervalDays *int                      `json:"maintenance_interval_days"`
	IsActive                bool                      `json:"is_active"`
}

// UpdateTariffRequest mendefinisikan payload pembaruan data tarif layanan
type UpdateTariffRequest struct {
	Code       string           `json:"code" binding:"required"`
	Name       string           `json:"name" binding:"required"`
	CategoryID string           `json:"category_id" binding:"required"`
	NameAlias  *string          `json:"name_alias"`
	ChargeType enums.ChargeType `json:"charge_type" binding:"required"`
	Notes      *string          `json:"notes"`
	IsActive   bool             `json:"is_active"`
}

// UpdateUnitRequest mendefinisikan payload pembaruan data satuan
type UpdateUnitRequest struct {
	UnitName         string  `json:"unit_name" binding:"required"`
	ConversionFactor float64 `json:"conversion_factor" binding:"required,gt=0"`
	IsPurchaseUnit   bool    `json:"is_purchase_unit"`
	IsDispenseUnit   bool    `json:"is_dispense_unit"`
}

// AddUnitRequest mendefinisikan payload penambahan satuan alternatif (UOM)
type AddUnitRequest struct {
	UnitName         string  `json:"unit_name" binding:"required"`
	ConversionFactor float64 `json:"conversion_factor" binding:"required,gt=0"`
	IsPurchaseUnit   bool    `json:"is_purchase_unit"`
	IsDispenseUnit   bool    `json:"is_dispense_unit"`
}

// QueryItemsRequest mendefinisikan query string pencarian katalog
type QueryItemsRequest struct {
	Page          int             `form:"page,default=1"`
	Limit         int             `form:"limit,default=10"`
	Search        string          `form:"search"`
	ItemType      *enums.ItemType `form:"item_type"`
	CategoryID    *string         `form:"category_id"`
	ProductLineID *string         `form:"product_line_id"`
	IsActive      *bool           `form:"is_active"`
}

// ==========================================
// Response DTOs
// ==========================================

type UnitResponse struct {
	ID               string    `json:"id"`
	UnitName         string    `json:"unit_name"`
	ConversionFactor float64   `json:"conversion_factor"`
	IsBaseUnit       bool      `json:"is_base_unit"`
	IsPurchaseUnit   bool      `json:"is_purchase_unit"`
	IsDispenseUnit   bool      `json:"is_dispense_unit"`
	CreatedAt        time.Time `json:"created_at"`
}

type MedicationResponse struct {
	KFACode            *string              `json:"kfa_code,omitempty"`
	BPOMNIE            *string              `json:"bpom_nie,omitempty"`
	DosageForm         *string              `json:"dosage_form,omitempty"`
	StrengthAmount     *string              `json:"strength_amount,omitempty"`
	StrengthUnit       *string              `json:"strength_unit,omitempty"`
	DefaultRoute       *string              `json:"default_route,omitempty"`
	MedicationType     enums.MedicationType `json:"medication_type,omitempty"`
	IsHighAlert        bool                 `json:"is_high_alert"`
	IsLASA             bool                 `json:"is_lasa"`
	IsFornas           bool                 `json:"is_fornas"`
	IsAntibiotic       bool                 `json:"is_antibiotic"`
	StorageTemperature *string              `json:"storage_temperature,omitempty"`
}

type GeneralResponse struct {
	GeneralType         enums.GeneralType          `json:"general_type"`
	IsSterile           bool                       `json:"is_sterile"`
	IsDisposable        bool                       `json:"is_disposable"`
	IsCSSDItem          bool                       `json:"is_cssd_item"`
	SterilizationMethod *enums.SterilizationMethod `json:"sterilization_method,omitempty"`
}

type AssetResponse struct {
	Brand                   *string                   `json:"brand,omitempty"`
	ModelName               *string                   `json:"model_name,omitempty"`
	IsMedicalEquipment      bool                      `json:"is_medical_equipment"`
	ExpectedLifeYears       *int                      `json:"expected_life_years,omitempty"`
	DepreciationMethod      *enums.DepreciationMethod `json:"depreciation_method,omitempty"`
	MaintenanceIntervalDays *int                      `json:"maintenance_interval_days,omitempty"`
}

type TariffResponse struct {
	NameAlias  *string          `json:"name_alias,omitempty"`
	ChargeType enums.ChargeType `json:"charge_type"`
	Notes      *string          `json:"notes,omitempty"`
}

type CategorySummaryResponse struct {
	ID            string         `json:"id"`
	Code          string         `json:"code"`
	Name          string         `json:"name"`
	ItemType      enums.ItemType `json:"item_type"`
	IncomeCOACode *string        `json:"income_coa_code,omitempty"`
}

type ProductLineSummaryResponse struct {
	ID               string  `json:"id"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	InventoryCOACode *string `json:"inventory_coa_code,omitempty"`
	CogsCOACode      *string `json:"cogs_coa_code,omitempty"`
}

type ItemDetailResponse struct {
	ID          string                      `json:"id"`
	Code        string                      `json:"code"`
	Name        string                      `json:"name"`
	GenericName *string                     `json:"generic_name,omitempty"`
	ItemType    enums.ItemType              `json:"item_type"`
	IsActive    bool                        `json:"is_active"`
	Category    *CategorySummaryResponse    `json:"category,omitempty"`
	ProductLine *ProductLineSummaryResponse `json:"product_line,omitempty"`
	Medication  *MedicationResponse         `json:"medication,omitempty"`
	General     *GeneralResponse            `json:"general,omitempty"`
	Asset       *AssetResponse              `json:"asset,omitempty"`
	Tariff      *TariffResponse             `json:"tariff,omitempty"`
	Units       []UnitResponse              `json:"units"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type ItemSummaryResponse struct {
	ID          string                      `json:"id"`
	Code        string                      `json:"code"`
	Name        string                      `json:"name"`
	GenericName *string                     `json:"generic_name,omitempty"`
	ItemType    enums.ItemType              `json:"item_type"`
	IsActive    bool                        `json:"is_active"`
	Category    *CategorySummaryResponse    `json:"category,omitempty"`
	ProductLine *ProductLineSummaryResponse `json:"product_line,omitempty"`
	BaseUnit    *UnitResponse               `json:"base_unit,omitempty"`
	CreatedAt   time.Time                   `json:"created_at"`
}

// ==========================================
// Mapping Helpers
// ==========================================

func ToUnitResponse(u *ItemUnit) UnitResponse {
	return UnitResponse{
		ID:               u.ID,
		UnitName:         u.UnitName,
		ConversionFactor: u.ConversionFactor,
		IsBaseUnit:       u.IsBaseUnit,
		IsPurchaseUnit:   u.IsPurchaseUnit,
		IsDispenseUnit:   u.IsDispenseUnit,
		CreatedAt:        u.CreatedAt,
	}
}

func ToUnitResponseList(units []ItemUnit) []UnitResponse {
	res := make([]UnitResponse, len(units))
	for i := range units {
		res[i] = ToUnitResponse(&units[i])
	}
	return res
}

func ToItemDetailResponse(item *Item) ItemDetailResponse {
	res := ItemDetailResponse{
		ID:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		GenericName: item.GenericName,
		ItemType:    item.ItemType,
		IsActive:    item.IsActive,
		Units:       ToUnitResponseList(item.Units),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	if item.Category != nil {
		res.Category = &CategorySummaryResponse{
			ID:            item.Category.ID,
			Code:          item.Category.Code,
			Name:          item.Category.Name,
			ItemType:      item.Category.ItemType,
			IncomeCOACode: item.Category.IncomeCOACode,
		}
	}

	if item.ProductLine != nil {
		res.ProductLine = &ProductLineSummaryResponse{
			ID:               item.ProductLine.ID,
			Code:             item.ProductLine.Code,
			Name:             item.ProductLine.Name,
			InventoryCOACode: item.ProductLine.InventoryCOACode,
			CogsCOACode:      item.ProductLine.CogsCOACode,
		}
	}

	if item.Medication != nil {
		res.Medication = &MedicationResponse{
			KFACode:            item.Medication.KFACode,
			BPOMNIE:            item.Medication.BPOMNIE,
			DosageForm:         item.Medication.DosageForm,
			StrengthAmount:     item.Medication.StrengthAmount,
			StrengthUnit:       item.Medication.StrengthUnit,
			DefaultRoute:       item.Medication.DefaultRoute,
			MedicationType:     item.Medication.MedicationType,
			IsHighAlert:        item.Medication.IsHighAlert,
			IsLASA:             item.Medication.IsLASA,
			IsFornas:           item.Medication.IsFornas,
			IsAntibiotic:       item.Medication.IsAntibiotic,
			StorageTemperature: item.Medication.StorageTemperature,
		}
	}

	if item.General != nil {
		res.General = &GeneralResponse{
			GeneralType:         item.General.GeneralType,
			IsSterile:           item.General.IsSterile,
			IsDisposable:        item.General.IsDisposable,
			IsCSSDItem:          item.General.IsCSSDItem,
			SterilizationMethod: item.General.SterilizationMethod,
		}
	}

	if item.Asset != nil {
		res.Asset = &AssetResponse{
			Brand:                   item.Asset.Brand,
			ModelName:               item.Asset.ModelName,
			IsMedicalEquipment:      item.Asset.IsMedicalEquipment,
			ExpectedLifeYears:       item.Asset.ExpectedLifeYears,
			DepreciationMethod:      item.Asset.DepreciationMethod,
			MaintenanceIntervalDays: item.Asset.MaintenanceIntervalDays,
		}
	}

	if item.Tariff != nil {
		res.Tariff = &TariffResponse{
			NameAlias:  item.Tariff.NameAlias,
			ChargeType: item.Tariff.ChargeType,
			Notes:      item.Tariff.Notes,
		}
	}

	return res
}

func ToItemSummaryResponse(item *Item) ItemSummaryResponse {
	res := ItemSummaryResponse{
		ID:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		GenericName: item.GenericName,
		ItemType:    item.ItemType,
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
	}

	if item.Category != nil {
		res.Category = &CategorySummaryResponse{
			ID:            item.Category.ID,
			Code:          item.Category.Code,
			Name:          item.Category.Name,
			ItemType:      item.Category.ItemType,
			IncomeCOACode: item.Category.IncomeCOACode,
		}
	}

	if item.ProductLine != nil {
		res.ProductLine = &ProductLineSummaryResponse{
			ID:               item.ProductLine.ID,
			Code:             item.ProductLine.Code,
			Name:             item.ProductLine.Name,
			InventoryCOACode: item.ProductLine.InventoryCOACode,
			CogsCOACode:      item.ProductLine.CogsCOACode,
		}
	}

	for i := range item.Units {
		if item.Units[i].IsBaseUnit {
			base := ToUnitResponse(&item.Units[i])
			res.BaseUnit = &base
			break
		}
	}

	return res
}
