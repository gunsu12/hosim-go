package referal

import (
	"hosim-go/pkg/enums"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Referal struct {
	ID        string            `gorm:"size:36;primaryKey" json:"id"`
	Code      string            `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name      string            `gorm:"size:150;not null" json:"name"`
	Type      enums.ReferalType `gorm:"size:20;not null" json:"type"`
	Address   string            `gorm:"size:255" json:"address"`
	Phone     string            `gorm:"size:20" json:"phone"`
	Email     string            `gorm:"size:100" json:"email"`
	CreatedAt time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time         `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
	CreatedBy string            `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy string            `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy string            `gorm:"default:'SYSTEM'" json:"-"`
}

func (r *Referal) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewV7().String()
	}
	return nil
}

func (r *Referal) BeforeUpdate(tx *gorm.DB) error {
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = time.Now()
	}
	if r.UpdatedBy == "" {
		r.UpdatedBy = "SYSTEM"
	}
	return nil
}
