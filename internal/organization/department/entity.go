package department

import (
	"hosim-go/pkg/enums"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Department struct {
	ID                string                `gorm:"primaryKey;size:36" json:"id"`
	Code              string                `gorm:"uniqueIndex;size:255" json:"code"`
	Name              string                `gorm:"size:255" json:"name"`
	Address           *string               `gorm:"size:255" json:"address,omitempty"`
	Phone             *string               `gorm:"size:255" json:"phone,omitempty"`
	Email             *string               `gorm:"size:255" json:"email,omitempty"`
	Website           *string               `gorm:"size:255" json:"website,omitempty"`
	Description       *string               `gorm:"size:255" json:"description,omitempty"`
	DepartementType   enums.DepartementType `gorm:"size:20;index" json:"departement_type"`
	IhsOrganizationId *string               `gorm:"size:255" json:"ihs_organization_id,omitempty"`
	IsActive          bool                  `gorm:"default:true;index" json:"is_active"`
	CreatedAt         time.Time             `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt        `gorm:"index" json:"-"`
	CreatedBy         string                `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy         string                `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy         string                `json:"-"`
}

// TableName mengembalikan nama tabel departements agar sesuai dengan migrasi database
func (Department) TableName() string {
	return "departements"
}

// Departement adalah type alias untuk Department demi backwards compatibility
type Departement = Department

func (d *Department) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewV7().String()
	}
	return nil
}
