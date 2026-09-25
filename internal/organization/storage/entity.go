package storage

import (
	"hosim-go/internal/organization/department"
	"hosim-go/pkg/enums"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Storage struct {
	ID              string                 `gorm:"size:36;primaryKey" json:"id"`
	DepartementID   *string                `gorm:"index;size:36" json:"departement_id,omitempty"`
	Departement     *department.Department `gorm:"foreignKey:DepartementID;references:ID" json:"departement,omitempty"`
	Name            string                 `gorm:"size:150;not null" json:"name"`
	Type            enums.StorageType      `gorm:"size:20;not null" json:"storage_type"`
	Capacity        float64                `gorm:"not null" json:"capacity"`
	StorageParentID *string                `gorm:"index;size:36" json:"storage_parent_id,omitempty"`
	StorageParent   *Storage               `gorm:"foreignKey:StorageParentID;references:ID" json:"storage_parent,omitempty"`
	Location        string                 `gorm:"size:150;not null" json:"location"`
	CreatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time              `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `gorm:"index" json:"-"`
	CreatedBy       string                 `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy       string                 `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy       string                 `json:"-"`
}

func (s *Storage) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewV7().String()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.CreatedBy == "" {
		s.CreatedBy = "SYSTEM"
	}
	return nil
}

func (s *Storage) BeforeUpdate(tx *gorm.DB) error {
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now()
	}
	if s.UpdatedBy == "" {
		s.UpdatedBy = "SYSTEM"
	}
	return nil
}
