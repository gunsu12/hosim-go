package audit

import (
	"context"
	"testing"
)

func TestAuditWriterValidation(t *testing.T) {
	w := NewWriter(nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		entry   AuditEntry
		wantErr error
	}{
		{
			name: "missing actor",
			entry: AuditEntry{
				Action:       "CREATE",
				ResourceType: "price_schema",
				ResourceID:   "schema-1",
			},
			wantErr: ErrActorRequired,
		},
		{
			name: "missing action",
			entry: AuditEntry{
				ActorID:      "user-1",
				ActorRole:    "ADMIN",
				ResourceType: "price_schema",
				ResourceID:   "schema-1",
			},
			wantErr: ErrActionRequired,
		},
		{
			name: "missing resource",
			entry: AuditEntry{
				ActorID:   "user-1",
				ActorRole: "ADMIN",
				Action:    "CREATE",
			},
			wantErr: ErrResourceRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := w.Write(ctx, tt.entry)
			if err != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
