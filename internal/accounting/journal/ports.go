package journal

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"hosim-go/internal/accounting/fiscalperiod"
)

// AccountSummary membawa data minimal akun yang sudah divalidasi.
type AccountSummary struct {
	ID         string
	Code       string
	Name       string
	IsActive   bool
	IsPostable bool
}

// AccountVerifier adalah outbound port untuk memvalidasi kelayakan akun COA.
type AccountVerifier interface {
	ValidatePostableAccounts(ctx context.Context, accountIDs []string) (map[string]*AccountSummary, error)
}

// FiscalPeriodGuard di-re-export dari domain fiscalperiod
type FiscalPeriodGuard = fiscalperiod.FiscalPeriodGuard

// ListJournalParams membawa kriteria pencarian dan paginasi jurnal.
type ListJournalParams struct {
	StartDate    *time.Time
	EndDate      *time.Time
	FiscalPeriod *string
	Status       *JournalStatus
	EntryType    *JournalType
	Search       string
	Page         int
	Limit        int
}

// Repository adalah port persistence untuk voucher jurnal dan baris rinciannya.
type Repository interface {
	SaveDraft(ctx context.Context, entry *JournalEntry) error
	UpdateDraft(ctx context.Context, entry *JournalEntry) error
	DeleteDraft(ctx context.Context, id string, actor string) error
	FindByID(ctx context.Context, id string) (*JournalEntry, error)
	FindByEntryNumber(ctx context.Context, entryNumber string) (*JournalEntry, error)
	List(ctx context.Context, params ListJournalParams) ([]*JournalEntry, int64, error)
	PostEntry(ctx context.Context, entry *JournalEntry) error
	CreatePostedEntry(ctx context.Context, entry *JournalEntry) error
	ReverseEntry(ctx context.Context, original *JournalEntry, reversal *JournalEntry) error
	GenerateNextEntryNumber(ctx context.Context, entryType JournalType, periodKey string) (string, error)
	CheckTrialBalanceDiscrepancy(ctx context.Context) (decimal.Decimal, error)
}

// SystemJournalLine adalah baris input dari modul luar (misal Billing/Kasir).
type SystemJournalLine struct {
	AccountID     string
	EntryPosition EntryPosition
	Amount        float64
	Memo          *string
	DepartmentID  *string
}

// PostSystemJournalCommand adalah perintah posting dari modul luar.
type PostSystemJournalCommand struct {
	EntryDate       time.Time
	EntryType       JournalType // misal SYSTEM_BILLING, SYSTEM_CASHIER
	Description     string
	ReferenceNumber *string
	Lines           []SystemJournalLine
	Actor           string
}

// JournalPostingPort adalah inbound port yang diekspos ke modul luar.
type JournalPostingPort interface {
	PostSystemJournal(ctx context.Context, cmd PostSystemJournalCommand) (*JournalEntry, error)
}
