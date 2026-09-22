package tariffclass

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type TariffClass struct {
	ID          string         `gorm:"size:36;primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`
	Description *string        `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy   string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy   string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy   string         `gorm:"default:'SYSTEM'" json:"-"`
}

func (b *TariffClass) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewV7().String()
	}
	return nil
}
