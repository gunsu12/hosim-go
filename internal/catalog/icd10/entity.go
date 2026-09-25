package icd10

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type ICD10 struct {
	ID          string          `gorm:"size:36;primaryKey" json:"id"`
	Code        string          `gorm:"size:20;index" json:"icd10_code"`
	Text        string          `gorm:"size:255;not null" json:"icd10_text"`
	Description *string         `gorm:"text" json:"icd10_description"`
	CreatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy   *string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy   *string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy   *string         `json:"-"`
}

func (icd *ICD10) BeforeCreate(tx *gorm.DB) error {
	if icd.ID == "" {
		icd.ID = uuid.NewV7().String()
	}
	return nil
}
