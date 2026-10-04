package priceplan

import (
	"errors"
	"math"
	"time"
	"uuid"

	"gorm.io/gorm"
)

// PricePlanStatus mendefinisikan siklus status buku tarif
type PricePlanStatus string

const (
	PricePlanStatusDraft     PricePlanStatus = "DRAFT"
	PricePlanStatusSubmitted PricePlanStatus = "SUBMITTED"
	PricePlanStatusApproved  PricePlanStatus = "APPROVED"
	PricePlanStatusActive    PricePlanStatus = "ACTIVE"
	PricePlanStatusArchived  PricePlanStatus = "ARCHIVED"
)

var (
	ErrPricePlanNotFound          = errors.New("buku tarif tidak ditemukan")
	ErrPricePlanCodeAlreadyExists = errors.New("kode buku tarif sudah terdaftar")
	ErrPricePlanItemNotFound      = errors.New("item buku tarif tidak ditemukan")
	ErrPricePlanImmutable         = errors.New("buku tarif yang telah diajukan/disetujui/aktif bersifat immutable dan tidak dapat diubah")
	ErrTariffComponentsMismatch   = errors.New("jumlah rincian komponen biaya tidak seimbang dengan total harga tarif")
	ErrTariffNotConfigured        = errors.New("tarif tidak ditemukan atau belum dikonfigurasi untuk tindakan dan kelas tersebut")
	ErrInvalidStateTransition     = errors.New("transisi status buku tarif tidak diizinkan")
	ErrDefaultPricePlanConflict   = errors.New("sudah terdapat buku tarif default aktif lain pada rentang tanggal yang sama")
	ErrEffectiveDateInvalid       = errors.New("tanggal mulai berlaku tidak boleh lebih besar dari tanggal berakhir")
	ErrItemAlreadyInPlan          = errors.New("item untuk kelas tarif tersebut sudah ada di dalam buku tarif ini")
	ErrComponentNotFound          = errors.New("komponen tarif tidak ditemukan")
)

// TariffPricePlan merepresentasikan entitas header buku tarif rumah sakit
type TariffPricePlan struct {
	ID                 string          `gorm:"size:36;primaryKey" json:"id"`
	Code               string          `gorm:"size:50;not null" json:"code"`
	Name               string          `gorm:"size:150;not null" json:"name"`
	Description        *string         `gorm:"type:text" json:"description,omitempty"`
	EffectiveFrom      time.Time       `gorm:"type:date;not null" json:"effective_from"`
	EffectiveTo        *time.Time      `gorm:"type:date" json:"effective_to,omitempty"`
	Status             PricePlanStatus `gorm:"size:20;default:'DRAFT';not null" json:"status"`
	IsDefault          bool            `gorm:"default:false;not null" json:"is_default"`
	CustomerID         *string         `gorm:"size:36" json:"customer_id,omitempty"`
	CustomerName       *string         `gorm:"->" json:"customer_name,omitempty"`
	DefaultCitoPercent float64         `gorm:"type:numeric(5,2);default:25.00;not null" json:"default_cito_percent"`
	ApprovedAt         *time.Time      `json:"approved_at,omitempty"`
	ApprovedBy         *string         `gorm:"size:50" json:"approved_by,omitempty"`
	CreatedAt          time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`
	CreatedBy          string          `gorm:"default:'SYSTEM'" json:"created_by"`
	UpdatedBy          string          `gorm:"default:'SYSTEM'" json:"updated_by"`
	DeletedBy          *string         `json:"deleted_by,omitempty"`

	Items []TariffPricePlanItem `gorm:"foreignKey:PricePlanID" json:"items,omitempty"`
}

func (TariffPricePlan) TableName() string {
	return "tariff_price_plans"
}

func (p *TariffPricePlan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	if p.Status == "" {
		p.Status = PricePlanStatusDraft
	}
	if p.DefaultCitoPercent <= 0 {
		p.DefaultCitoPercent = 25.00
	}
	return nil
}

// CanModify memeriksa apakah isi buku tarif masih boleh dimodifikasi (hanya jika DRAFT)
func (p *TariffPricePlan) CanModify() bool {
	return p.Status == PricePlanStatusDraft
}

// CanTransitionTo memvalidasi apakah transisi ke status target diizinkan
func (p *TariffPricePlan) CanTransitionTo(target PricePlanStatus) error {
	switch p.Status {
	case PricePlanStatusDraft:
		if target == PricePlanStatusSubmitted {
			return nil
		}
	case PricePlanStatusSubmitted:
		if target == PricePlanStatusDraft || target == PricePlanStatusApproved {
			return nil
		}
	case PricePlanStatusApproved:
		if target == PricePlanStatusActive {
			return nil
		}
	case PricePlanStatusActive:
		if target == PricePlanStatusArchived {
			return nil
		}
	}
	return ErrInvalidStateTransition
}

// IsEffectiveOn memeriksa apakah tanggal transaksi berada dalam rentang berlaku buku tarif ini
func (p *TariffPricePlan) IsEffectiveOn(targetDate time.Time) bool {
	// Normalisasi ke perbandingan tanggal saja (abaikan jam)
	tYear, tMonth, tDay := targetDate.Date()
	tDate := time.Date(tYear, tMonth, tDay, 0, 0, 0, 0, time.UTC)

	fromYear, fromMonth, fromDay := p.EffectiveFrom.Date()
	fromDate := time.Date(fromYear, fromMonth, fromDay, 0, 0, 0, 0, time.UTC)

	if tDate.Before(fromDate) {
		return false
	}

	if p.EffectiveTo != nil {
		toYear, toMonth, toDay := p.EffectiveTo.Date()
		toDate := time.Date(toYear, toMonth, toDay, 23, 59, 59, 999999999, time.UTC)
		if tDate.After(toDate) {
			return false
		}
	}

	return true
}

// TariffPricePlanItem merepresentasikan tarif per item tindakan/layanan pada kelas tarif tertentu
type TariffPricePlanItem struct {
	ID             string         `gorm:"size:36;primaryKey" json:"id"`
	PricePlanID    string         `gorm:"size:36;not null;index" json:"price_plan_id"`
	ItemID         string         `gorm:"size:36;not null;index" json:"item_id"`
	ItemName       string         `gorm:"->" json:"item_name,omitempty"`
	ItemCode       string         `gorm:"->" json:"item_code,omitempty"`
	TariffClassID  string         `gorm:"size:36;not null;index" json:"tariff_class_id"`
	TariffClassCode string        `gorm:"->" json:"tariff_class_code,omitempty"`
	TariffClassName string        `gorm:"->" json:"tariff_class_name,omitempty"`
	TotalBasePrice float64        `gorm:"type:numeric(15,2);default:0.00;not null" json:"total_base_price"`
	TotalCitoPrice *float64       `gorm:"type:numeric(15,2)" json:"total_cito_price,omitempty"`
	IsActive       bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy      string         `gorm:"default:'SYSTEM'" json:"created_by"`
	UpdatedBy      string         `gorm:"default:'SYSTEM'" json:"updated_by"`
	DeletedBy      *string        `json:"deleted_by,omitempty"`

	Components []TariffPricePlanItemComponent `gorm:"foreignKey:PlanItemID" json:"components,omitempty"`
}

func (TariffPricePlanItem) TableName() string {
	return "tariff_price_plan_items"
}

func (i *TariffPricePlanItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewV7().String()
	}
	return nil
}

// ValidateBalancing memvalidasi keselarasan total komponen dengan total_base_price
func (i *TariffPricePlanItem) ValidateBalancing() error {
	var sum float64
	for _, comp := range i.Components {
		sum += comp.BaseAmount
	}

	// Gunakan toleransi desimal 0.01
	if math.Abs(sum-i.TotalBasePrice) > 0.01 {
		return ErrTariffComponentsMismatch
	}
	return nil
}

// ResolvedTariffComponent adalah representasi rincian komponen hasil perhitungan tarif
type ResolvedTariffComponent struct {
	ComponentID   string  `json:"component_id"`
	ComponentName string  `json:"component_name"`
	ComponentType string  `json:"component_type"`
	Amount        float64 `json:"amount"`
	COACode       string  `json:"coa_code"`
}

// CalculatePrice menghitung nominal final dan rincian komponen (reguler atau CITO)
func (i *TariffPricePlanItem) CalculatePrice(isCito bool, defaultCitoPercent float64) (float64, []ResolvedTariffComponent) {
	if !isCito {
		resolved := make([]ResolvedTariffComponent, len(i.Components))
		for idx, c := range i.Components {
			coa := ""
			if c.COACode != nil && *c.COACode != "" {
				coa = *c.COACode
			} else if c.DefaultCOACode != nil {
				coa = *c.DefaultCOACode
			}

			resolved[idx] = ResolvedTariffComponent{
				ComponentID:   c.ComponentID,
				ComponentName: c.ComponentName,
				ComponentType: c.ComponentType,
				Amount:        math.Round(c.BaseAmount*100) / 100,
				COACode:       coa,
			}
		}
		return math.Round(i.TotalBasePrice*100) / 100, resolved
	}

	// Kasus CITO:
	var finalTotalPrice float64
	if i.TotalCitoPrice != nil && *i.TotalCitoPrice > 0 {
		finalTotalPrice = *i.TotalCitoPrice
	} else {
		finalTotalPrice = i.TotalBasePrice * (1.0 + (defaultCitoPercent / 100.0))
	}
	finalTotalPrice = math.Round(finalTotalPrice*100) / 100

	resolved := make([]ResolvedTariffComponent, len(i.Components))
	for idx, c := range i.Components {
		var compAmount float64
		if c.CitoAmount != nil && *c.CitoAmount > 0 {
			compAmount = *c.CitoAmount
		} else {
			compAmount = c.BaseAmount * (1.0 + (defaultCitoPercent / 100.0))
		}

		coa := ""
		if c.COACode != nil && *c.COACode != "" {
			coa = *c.COACode
		} else if c.DefaultCOACode != nil {
			coa = *c.DefaultCOACode
		}

		resolved[idx] = ResolvedTariffComponent{
			ComponentID:   c.ComponentID,
			ComponentName: c.ComponentName,
			ComponentType: c.ComponentType,
			Amount:        math.Round(compAmount*100) / 100,
			COACode:       coa,
		}
	}

	return finalTotalPrice, resolved
}

// TariffPricePlanItemComponent merepresentasikan pemecahan komponen biaya suatu tarif item
type TariffPricePlanItemComponent struct {
	ID             string    `gorm:"size:36;primaryKey" json:"id"`
	PlanItemID     string    `gorm:"size:36;not null;index" json:"plan_item_id"`
	ComponentID    string    `gorm:"size:36;not null;index" json:"component_id"`
	ComponentName  string    `gorm:"->" json:"component_name,omitempty"`
	ComponentCode  string    `gorm:"->" json:"component_code,omitempty"`
	ComponentType  string    `gorm:"->" json:"component_type,omitempty"`
	DefaultCOACode *string   `gorm:"->" json:"default_coa_code,omitempty"`
	BaseAmount     float64   `gorm:"type:numeric(15,2);default:0.00;not null" json:"base_amount"`
	CitoAmount     *float64  `gorm:"type:numeric(15,2)" json:"cito_amount,omitempty"`
	COACode        *string   `gorm:"size:50" json:"coa_code,omitempty"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy      string    `gorm:"default:'SYSTEM'" json:"created_by"`
	UpdatedBy      string    `gorm:"default:'SYSTEM'" json:"updated_by"`
}

func (TariffPricePlanItemComponent) TableName() string {
	return "tariff_price_plan_item_components"
}

func (c *TariffPricePlanItemComponent) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewV7().String()
	}
	return nil
}
