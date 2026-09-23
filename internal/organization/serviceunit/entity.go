package serviceunit

import (
	"hosim-go/internal/master/departement"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type ServiceUnit struct {
	ID                   string                   `gorm:"primaryKey;size:36" json:"id"`
	Code                 string                   `gorm:"uniqueIndex;size:255" json:"code"`
	Name                 string                   `gorm:"size:255" json:"name"`
	DepartementID        *string                  `gorm:"index;size:36" json:"departement_id,omitempty"`
	Departement          *departement.Departement `gorm:"foreignKey:DepartementID;references:ID" json:"departement,omitempty"`
	Address              *string                  `gorm:"size:255" json:"address,omitempty"`
	Phone                *string                  `gorm:"size:255" json:"phone,omitempty"`
	Email                *string                  `gorm:"size:255" json:"email,omitempty"`
	Website              *string                  `gorm:"size:255" json:"website,omitempty"`
	Description          *string                  `gorm:"size:255" json:"description,omitempty"`
	IhsLocationId        *string                  `gorm:"size:255" json:"ihs_location_id,omitempty"`
	IsRegistrationTarget bool                     `gorm:"default:true;index" json:"is_registration_target"`
	IsActive             bool                     `gorm:"default:true;index" json:"is_active"`
	CreatedAt            time.Time                `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time                `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt            gorm.DeletedAt           `gorm:"index" json:"-"`
	CreatedBy            string                   `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy            string                   `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy            string                   `json:"-"`
}

func (s *ServiceUnit) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewV7().String()
	}
	return nil
}
