package customer

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Customer struct {
	ID             string         `gorm:"primaryKey;size:36" json:"id"`
	Code           string         `gorm:"uniqueIndex;size:255" json:"code"`
	Name           string         `gorm:"size:255" json:"name"`
	Address        *string        `gorm:"size:255" json:"address,omitempty"`
	Phone          *string        `gorm:"size:255" json:"phone,omitempty"`
	Email          *string        `gorm:"size:255" json:"email,omitempty"`
	Website        *string        `gorm:"size:255" json:"website,omitempty"`
	ContactPerson  *string        `gorm:"size:255" json:"contact_person,omitempty"`
	RequireCard    *bool          `gorm:"default:true" json:"require_card"`
	Description    *string        `gorm:"size:255" json:"description,omitempty"`
	CustomerTypeID *string        `gorm:"index;size:36" json:"customer_type_id,omitempty"`
	CustomerType   *CustomerType  `gorm:"foreignKey:CustomerTypeID;references:ID" json:"customer_type,omitempty"`
	IsActive       bool           `gorm:"default:true;index" json:"is_active"`
	IsImmutable    bool           `gorm:"default:false" json:"is_immutable"` // penjamin seperti Pribadi, Asing, KITAS bersifat immutable (wajib ada di sistem tidak boleh dihapus)
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy      string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy      string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy      string         `json:"-"`
}

/*
CustomerType menyimpan tipe customer seperti Pribadi, Perusahaan, Asuransi, Pemerintah, BPJS
*/
type CustomerType struct {
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

func (Customer) TableName() string {
	return "customers"
}

func (CustomerType) TableName() string {
	return "customer_types"
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewV7().String()
	}
	return nil
}

func (ct *CustomerType) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == "" {
		ct.ID = uuid.NewV7().String()
	}
	return nil
}
