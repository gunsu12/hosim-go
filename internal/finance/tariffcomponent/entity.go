package tariffcomponent

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type TariffComponent struct {
	ID                 string         `gorm:"size:36;primaryKey" json:"id"`
	Code               string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name               string         `gorm:"size:150;not null" json:"name"`
	Description        *string        `gorm:"size:255" json:"description,omitempty"`
	IsHospitalRevenue  bool           `gorm:"default:false;not null" json:"is_hospital_revenue"`
	IsOperatorRevenue  bool           `gorm:"default:false;not null" json:"is_operator_revenue"`
	IsParamedicRevenue bool           `gorm:"default:false;not null" json:"is_paramedic_revenue"`
	CreatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy          string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy          string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy          string         `gorm:"default:'SYSTEM'" json:"-"`
}

func (TariffComponent) TableName() string {
	return "tariff_components"
}

func (tc *TariffComponent) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == "" {
		tc.ID = uuid.NewV7().String()
	}
	return nil
}
