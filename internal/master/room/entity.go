package room

import (
	serviceunit "hosim-go/internal/master/service_unit"
	"hosim-go/pkg/enums"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Room struct {
	ID            string                  `gorm:"size:36;primaryKey" json:"id"`
	Code          string                  `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name          string                  `gorm:"size:150;not null" json:"name"`
	Capacity      int                     `gorm:"not null" json:"capacity"`
	Location      string                  `gorm:"size:150" json:"location"`
	RoomType      enums.RoomType          `gorm:"size:50;not null" json:"room_type"`
	ServiceUnitID string                  `gorm:"size:36;not null" json:"service_unit_id"`
	ServiceUnit   serviceunit.ServiceUnit `gorm:"foreignKey:ServiceUnitID;references:ID" json:"service_unit"`
	CreatedAt     time.Time               `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time               `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt          `gorm:"index" json:"-"`
	CreatedBy     string                  `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy     string                  `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy     string                  `gorm:"default:'SYSTEM'" json:"-"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewV7().String()
	}
	return nil
}
