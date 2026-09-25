package bed

import (
	"hosim-go/internal/organization/room"
	"hosim-go/pkg/enums"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type Bed struct {
	ID        string          `gorm:"size:36;primaryKey" json:"id"`
	Code      string          `gorm:"size:36;not null" json:"code"`
	Name      string          `gorm:"size:255;not null" json:"name"`
	RoomID    *string         `gorm:"index;size:36" json:"room_id,omitempty"`
	Room      *room.Room      `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	Status    enums.BedStatus `gorm:"size:36" json:"status"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
	CreatedBy string          `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy string          `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy string          `json:"-"`
}

func (b *Bed) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewV7().String()
	}
	if b.Status == "" {
		b.Status = enums.BedStatusAvailable
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.CreatedBy == "" {
		b.CreatedBy = "SYSTEM"
	}
	return nil
}

func (b *Bed) BeforeUpdate(tx *gorm.DB) error {
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	if b.UpdatedBy == "" {
		b.UpdatedBy = "SYSTEM"
	}
	return nil
}
