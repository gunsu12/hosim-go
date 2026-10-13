package audit

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"hosim-go/internal/database"

	"gorm.io/gorm"
)

var (
	ErrActorRequired    = errors.New("actor_id dan actor_role wajib diisi pada audit entry")
	ErrResourceRequired = errors.New("resource_type dan resource_id wajib diisi pada audit entry")
	ErrActionRequired   = errors.New("action wajib diisi pada audit entry")
)

// AuditEntry adalah payload audit yang diteruskan oleh application layer.
type AuditEntry struct {
	ActorID       string
	ActorUsername string
	ActorRole     string
	Action        string
	ResourceType  string
	ResourceID    string
	Before        any
	After         any
	CorrelationID string
	Timestamp     time.Time
}

// AuditWriter adalah port untuk mencatat audit trail secara append-only.
type AuditWriter interface {
	Write(ctx context.Context, entry AuditEntry) error
}

type writer struct {
	db *gorm.DB
}

// NewWriter membuat instance AuditWriter yang mendukung transaksi database context.
func NewWriter(db *gorm.DB) AuditWriter {
	return &writer{db: db}
}

func (w *writer) Write(ctx context.Context, entry AuditEntry) error {
	entry.ActorID = strings.TrimSpace(entry.ActorID)
	entry.ActorRole = strings.TrimSpace(entry.ActorRole)
	entry.Action = strings.TrimSpace(entry.Action)
	entry.ResourceType = strings.TrimSpace(entry.ResourceType)
	entry.ResourceID = strings.TrimSpace(entry.ResourceID)

	if entry.ActorID == "" || entry.ActorRole == "" {
		return ErrActorRequired
	}
	if entry.Action == "" {
		return ErrActionRequired
	}
	if entry.ResourceType == "" || entry.ResourceID == "" {
		return ErrResourceRequired
	}

	var beforeJSON *string
	if entry.Before != nil {
		b, err := json.Marshal(entry.Before)
		if err != nil {
			return err
		}
		s := string(b)
		beforeJSON = &s
	}

	var afterJSON *string
	if entry.After != nil {
		b, err := json.Marshal(entry.After)
		if err != nil {
			return err
		}
		s := string(b)
		afterJSON = &s
	}

	var usernamePtr *string
	if u := strings.TrimSpace(entry.ActorUsername); u != "" {
		usernamePtr = &u
	}

	var corrPtr *string
	if c := strings.TrimSpace(entry.CorrelationID); c != "" {
		corrPtr = &c
	}

	createdAt := entry.Timestamp.UTC()
	if entry.Timestamp.IsZero() {
		createdAt = time.Now().UTC()
	}

	record := AuditLog{
		ActorID:       entry.ActorID,
		ActorUsername: usernamePtr,
		ActorRole:     entry.ActorRole,
		Action:        entry.Action,
		ResourceType:  entry.ResourceType,
		ResourceID:    entry.ResourceID,
		BeforePayload: beforeJSON,
		AfterPayload:  afterJSON,
		CorrelationID: corrPtr,
		CreatedAt:     createdAt,
	}

	// Gunakan koneksi transaksi aktif jika tersedia di context
	db := database.DBFromContextOr(ctx, w.db)
	return db.Create(&record).Error
}
