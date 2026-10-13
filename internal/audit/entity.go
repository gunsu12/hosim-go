package audit

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

// AuditLog merepresentasikan catatan audit append-only di tabel audit_logs.
type AuditLog struct {
	ID            string    `gorm:"primaryKey;size:36" json:"id"`
	ActorID       string    `gorm:"size:50;not null;index" json:"actor_id"`
	ActorUsername *string   `gorm:"size:100" json:"actor_username,omitempty"`
	ActorRole     string    `gorm:"size:50;not null" json:"actor_role"`
	Action        string    `gorm:"size:50;not null" json:"action"`
	ResourceType  string    `gorm:"size:50;not null;index:idx_audit_logs_res" json:"resource_type"`
	ResourceID    string    `gorm:"size:36;not null;index:idx_audit_logs_res" json:"resource_id"`
	BeforePayload *string   `gorm:"type:jsonb" json:"before_payload,omitempty"`
	AfterPayload  *string   `gorm:"type:jsonb" json:"after_payload,omitempty"`
	CorrelationID *string   `gorm:"size:100" json:"correlation_id,omitempty"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewV7().String()
	}
	return nil
}
