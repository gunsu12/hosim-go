package icd9

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type ICD9 struct {
	ID          string          `gorm:"size:36;primaryKey" json:"id"`
	Code        string          `gorm:"size:20;index" json:"icd9_code"`
	Text        string          `gorm:"size:255;not null" json:"icd9_text"`
	Description *string         `gorm:"text" json:"icd9_description"`
	CreatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy   *string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy   *string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy   *string         `json:"-"`
}

func (icd *ICD9) BeforeCreate(tx *gorm.DB) error {
	if icd.ID == "" {
		icd.ID = uuid.NewV7().String()
	}
	return nil
}
