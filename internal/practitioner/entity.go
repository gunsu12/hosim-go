package practitioner

import (
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

type Practitioner struct {
	ID                  string       `gorm:"primaryKey;size:36" json:"id"`
	NIK                 string       `gorm:"column:nik;uniqueIndex:idx_practitioners_nik,where:nik != '' AND nik IS NOT NULL AND deleted_at IS NULL;size:16" json:"nik"`
	NIP                 string       `gorm:"column:nip;uniqueIndex:idx_practitioners_nip,where:nip != '' AND nip IS NOT NULL AND deleted_at IS NULL;size:50" json:"nip,omitempty"`
	Name                string       `gorm:"size:150;not null" json:"name"`
	Gender              enums.Gender `gorm:"size:10;not null" json:"gender"`
	SIP                 string       `gorm:"column:sip;size:50" json:"sip,omitempty"`
	SIPExpiryDate       *time.Time   `gorm:"column:sip_expiry_date" json:"sip_expiry_date,omitempty"`
	STR                 string       `gorm:"column:str;size:50" json:"str,omitempty"` // str biasanya seumur hidup
	ProfessionID        *string     `gorm:"index;size:36" json:"profession_id,omitempty"`
	Profession          *Profession `gorm:"foreignKey:ProfessionID;references:ID" json:"profession,omitempty"`
	SpecialtyID         *string     `gorm:"index;size:36" json:"specialty_id,omitempty"`
	Specialty           *Specialty  `gorm:"foreignKey:SpecialtyID;references:ID" json:"specialty,omitempty"`
	Phone               string      `gorm:"size:20" json:"phone"`
	Email               string      `gorm:"size:100" json:"email"`
	IhsPractitionerID   string      `gorm:"uniqueIndex:idx_practitioners_ihs,where:ihs_practitioner_id != '' AND ihs_practitioner_id IS NOT NULL AND deleted_at IS NULL;size:50" json:"ihs_practitioner_id,omitempty"`
	IhsPractitionerName string      `gorm:"size:150" json:"ihs_practitioner_name,omitempty"`
	IsActive            bool        `gorm:"default:true;index" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

type Specialty struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	Code      string         `gorm:"uniqueIndex;size:20;not null" json:"code"` // Misal: SP-A, SP-PD
	Name      string         `gorm:"size:150;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

type Profession struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	Code      string         `gorm:"uniqueIndex;size:20;not null" json:"code"` // Misal: DOKTER, PERAWAT, BIDAN, FARMASI
	Name      string         `gorm:"size:150;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy string         `json:"deleted_by"`
}

func (p *Practitioner) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}

func (s *Specialty) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewV7().String()
	}
	return nil
}

func (p *Profession) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}
