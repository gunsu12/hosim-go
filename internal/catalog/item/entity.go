package item

import (
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

// ItemCategory merepresentasikan pengelompokan item untuk pemetaan COA pendapatan
type ItemCategory struct {
	ID              string         `gorm:"primaryKey;size:36" json:"id"`
	Code            string         `gorm:"uniqueIndex:idx_item_categories_code,where:deleted_at IS NULL;size:50;not null" json:"code"`
	Name            string         `gorm:"size:150;not null" json:"name"`
	ItemType        enums.ItemType `gorm:"size:20;not null" json:"item_type"` // Filter subtipe item
	IncomeCOACode   *string        `gorm:"size:50" json:"income_coa_code,omitempty"`
	DiscountCOACode *string        `gorm:"size:50" json:"discount_coa_code,omitempty"`
	SalesTaxCOACode *string        `gorm:"size:50" json:"sales_tax_coa_code,omitempty"`
	IsActive        bool           `gorm:"default:true;not null" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by,omitempty"`
}

func (ItemCategory) TableName() string {
	return "item_categories"
}

func (c *ItemCategory) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewV7().String()
	}
	return nil
}

// ItemProductLine merepresentasikan lini produk persediaan dan pemetaan COA persediaan/HPP/aset
type ItemProductLine struct {
	ID                             string         `gorm:"primaryKey;size:36" json:"id"`
	Code                           string         `gorm:"uniqueIndex:idx_item_product_lines_code,where:deleted_at IS NULL;size:50;not null" json:"code"`
	Name                           string         `gorm:"size:150;not null" json:"name"`
	InventoryCOACode               *string        `gorm:"size:50" json:"inventory_coa_code,omitempty"`
	CogsCOACode                    *string        `gorm:"size:50" json:"cogs_coa_code,omitempty"`
	PurchaseDiscountCOACode        *string        `gorm:"size:50" json:"purchase_discount_coa_code,omitempty"`
	PurchaseTaxCOACode             *string        `gorm:"size:50" json:"purchase_tax_coa_code,omitempty"`
	AssetCOACode                   *string        `gorm:"size:50" json:"asset_coa_code,omitempty"`
	AssetAccumulationCOACode       *string        `gorm:"size:50" json:"asset_accumulation_coa_code,omitempty"`
	AssetDepreciationExpenseCOACode *string        `gorm:"size:50" json:"asset_depreciation_expense_coa_code,omitempty"`
	IsActive                       bool           `gorm:"default:true;not null" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by,omitempty"`
}

func (ItemProductLine) TableName() string {
	return "item_product_lines"
}

func (pl *ItemProductLine) BeforeCreate(tx *gorm.DB) error {
	if pl.ID == "" {
		pl.ID = uuid.NewV7().String()
	}
	return nil
}

// Item adalah entitas universal / aggregate root untuk katalog produk fisik dan layanan tarif rumah sakit
type Item struct {
	ID            string         `gorm:"primaryKey;size:36" json:"id"`
	Code          string         `gorm:"uniqueIndex:idx_items_code,where:deleted_at IS NULL;size:50;not null" json:"code"`
	Name          string         `gorm:"size:200;not null;index:idx_items_name" json:"name"`
	GenericName   *string        `gorm:"size:200" json:"generic_name,omitempty"`
	ItemType      enums.ItemType `gorm:"size:20;not null;index:idx_items_type" json:"item_type"`
	CategoryID    string         `gorm:"index:idx_items_category;size:36;not null" json:"category_id"`
	ProductLineID *string        `gorm:"index:idx_items_product_line;size:36" json:"product_line_id,omitempty"`
	IsActive      bool           `gorm:"default:true;not null;index:idx_items_active" json:"is_active"`

	// Relasi Asosiasi / Lookup
	Category    *ItemCategory    `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	ProductLine *ItemProductLine `gorm:"foreignKey:ProductLineID;references:ID" json:"product_line,omitempty"`

	// Relasi Subtipe (1:1 Class Table Inheritance)
	Medication *ItemMedication `gorm:"foreignKey:ItemID;references:ID" json:"medication,omitempty"`
	General    *ItemGeneral    `gorm:"foreignKey:ItemID;references:ID" json:"general,omitempty"`
	Asset      *ItemAsset      `gorm:"foreignKey:ItemID;references:ID" json:"asset,omitempty"`
	Tariff     *ItemTariff     `gorm:"foreignKey:ItemID;references:ID" json:"tariff,omitempty"`

	// Relasi Multi-Satuan (1:N)
	Units []ItemUnit `gorm:"foreignKey:ItemID;references:ID" json:"units,omitempty"`

	// Audit Trail
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by,omitempty"`
}

func (Item) TableName() string {
	return "items"
}

func (i *Item) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewV7().String()
	}
	return nil
}

// ItemMedication menyimpan atribut spesifik obat, cairan infus, dan vaksin (1:1 ke items)
type ItemMedication struct {
	ItemID             string               `gorm:"primaryKey;size:36" json:"item_id"`
	KFACode            *string              `gorm:"column:kfa_code;size:50;index:idx_item_meds_kfa" json:"kfa_code,omitempty"`
	BPOMNIE            *string              `gorm:"column:bpom_nie;size:50;index:idx_item_meds_nie" json:"bpom_nie,omitempty"`
	DosageForm         *string              `gorm:"size:50" json:"dosage_form,omitempty"`
	StrengthAmount     *string              `gorm:"size:50" json:"strength_amount,omitempty"`
	StrengthUnit       *string              `gorm:"size:20" json:"strength_unit,omitempty"`
	DefaultRoute       *string              `gorm:"size:50" json:"default_route,omitempty"`
	MedicationType     enums.MedicationType `gorm:"size:30" json:"medication_type,omitempty"`
	IsHighAlert        bool                 `gorm:"default:false;not null" json:"is_high_alert"`
	IsLASA             bool                 `gorm:"column:is_lasa;default:false;not null" json:"is_lasa"`
	IsFornas           bool                 `gorm:"default:false;not null" json:"is_fornas"`
	IsAntibiotic       bool                 `gorm:"default:false;not null" json:"is_antibiotic"`
	StorageTemperature *string              `gorm:"size:50" json:"storage_temperature,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (ItemMedication) TableName() string {
	return "item_medications"
}

// ItemGeneral menyimpan atribut BMHP, instrumen bedah, linen steril, dan logistik umum (1:1 ke items)
type ItemGeneral struct {
	ItemID              string                     `gorm:"primaryKey;size:36" json:"item_id"`
	GeneralType         enums.GeneralType          `gorm:"size:30;not null" json:"general_type"`
	IsSterile           bool                       `gorm:"default:false;not null" json:"is_sterile"`
	IsDisposable        bool                       `gorm:"default:true;not null" json:"is_disposable"`
	IsCSSDItem          bool                       `gorm:"column:is_cssd_item;default:false;not null" json:"is_cssd_item"`
	SterilizationMethod *enums.SterilizationMethod `gorm:"size:30" json:"sterilization_method,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (ItemGeneral) TableName() string {
	return "item_generals"
}

// ItemAsset menyimpan atribut barang modal dan alat kesehatan yang disusutkan (1:1 ke items)
type ItemAsset struct {
	ItemID                  string                    `gorm:"primaryKey;size:36" json:"item_id"`
	Brand                   *string                   `gorm:"size:100" json:"brand,omitempty"`
	ModelName               *string                   `gorm:"size:100" json:"model_name,omitempty"`
	IsMedicalEquipment      bool                      `gorm:"default:false;not null" json:"is_medical_equipment"`
	ExpectedLifeYears       *int                      `json:"expected_life_years,omitempty"`
	DepreciationMethod      *enums.DepreciationMethod `gorm:"size:30" json:"depreciation_method,omitempty"`
	MaintenanceIntervalDays *int                      `json:"maintenance_interval_days,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (ItemAsset) TableName() string {
	return "item_assets"
}

// ItemTariff menyimpan atribut master layanan dan tarif jasa rumah sakit (1:1 ke items)
type ItemTariff struct {
	ItemID     string           `gorm:"primaryKey;size:36" json:"item_id"`
	NameAlias  *string          `gorm:"size:200" json:"name_alias,omitempty"`
	ChargeType enums.ChargeType `gorm:"size:30;not null" json:"charge_type"`
	Notes      *string          `gorm:"type:text" json:"notes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (ItemTariff) TableName() string {
	return "item_tariffs"
}

// ItemUnit merepresentasikan multi-satuan ukur (UOM) dan rasio konversi terhadap base unit (1:N ke items)
type ItemUnit struct {
	ID               string    `gorm:"primaryKey;size:36" json:"id"`
	ItemID           string    `gorm:"index:idx_item_units_item;size:36;not null" json:"item_id"`
	UnitName         string    `gorm:"size:50;not null" json:"unit_name"`
	ConversionFactor float64   `gorm:"type:numeric(12,4);not null;default:1" json:"conversion_factor"`
	IsBaseUnit       bool      `gorm:"default:false;not null" json:"is_base_unit"`
	IsPurchaseUnit   bool      `gorm:"default:false;not null" json:"is_purchase_unit"`
	IsDispenseUnit   bool      `gorm:"default:false;not null" json:"is_dispense_unit"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (ItemUnit) TableName() string {
	return "item_units"
}

func (u *ItemUnit) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewV7().String()
	}
	return nil
}
