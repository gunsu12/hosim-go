package bed

import (
	"hosim-go/internal/organization/room"
	"time"

	"gorm.io/gorm"
)

type Bed struct {
	ID        string         `gorm:"size:36;primaryKey" json:"id"`
	Code      string         `gorm:"size:36;not null" json:"code"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	RoomID    *string        `gorm:"index;size:36" json:"room_id,omitempty"`
	Room      *room.Room     `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy string         `json:"-"`
}
