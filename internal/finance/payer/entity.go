package payer

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Payer struct {
	ID            string         `gorm:"primaryKey;size:36" json:"id"`
	Code          string         `gorm:"uniqueIndex;size:255" json:"code"`
	Name          string         `gorm:"size:255" json:"name"`
	Address       *string        `gorm:"size:255" json:"address,omitempty"`
	Phone         *string        `gorm:"size:255" json:"phone,omitempty"`
	Email         *string        `gorm:"size:255" json:"email,omitempty"`
	Website       *string        `gorm:"size:255" json:"website,omitempty"`
	ContactPerson *string        `gorm:"size:255" json:"contact_person,omitempty"`
	RequireCard   *bool          `gorm:"default:true" json:"require_card"`
	Description   *string        `gorm:"size:255" json:"description,omitempty"`
	PayerTypeID   *string        `gorm:"index;size:36" json:"payer_type_id,omitempty"`
	PayerType     *PayerType     `gorm:"foreignKey:PayerTypeID;references:ID" json:"payer_type,omitempty"`
	IsActive      bool           `gorm:"default:true;index" json:"is_active"`
	IsImmutable   bool           `gorm:"default:false" json:"is_imutable"` // penjamin seperti Pribadi, Asing, KITAS harusnya immutable artinya di sistem wajib ada tidak boleh dihapus
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy     string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy     string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy     string         `json:"-"`
}

/*
store payer type like Pribadi, Perusahaan, Asuransi, Pemerintah, BPJS
*/
type PayerType struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	Code      string         `gorm:"uniqueIndex;size:255" json:"code"`
	Name      string         `gorm:"size:255" json:"name"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy string         `json:"-"`
}

func (p *Payer) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}

func (p *PayerType) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}
