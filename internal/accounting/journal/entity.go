package journal

import (
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type JournalStatus string

const (
	JournalStatusDraft    JournalStatus = "DRAFT"
	JournalStatusPosted   JournalStatus = "POSTED"
	JournalStatusReversed JournalStatus = "REVERSED"
)

func (s JournalStatus) IsValid() bool {
	switch s {
	case JournalStatusDraft, JournalStatusPosted, JournalStatusReversed:
		return true
	default:
		return false
	}
}

type JournalType string

const (
	JournalTypeGeneral         JournalType = "GENERAL"
	JournalTypeAdjustment      JournalType = "ADJUSTMENT"
	JournalTypeClosing         JournalType = "CLOSING"
	JournalTypeReversal        JournalType = "REVERSAL"
	JournalTypeSystemBilling   JournalType = "SYSTEM_BILLING"
	JournalTypeSystemCashier   JournalType = "SYSTEM_CASHIER"
	JournalTypeSystemInventory JournalType = "SYSTEM_INVENTORY"
)

func (t JournalType) IsValid() bool {
	switch t {
	case JournalTypeGeneral, JournalTypeAdjustment, JournalTypeClosing, JournalTypeReversal,
		JournalTypeSystemBilling, JournalTypeSystemCashier, JournalTypeSystemInventory:
		return true
	default:
		return false
	}
}

type EntryPosition string

const (
	EntryPositionDebit  EntryPosition = "DEBIT"
	EntryPositionCredit EntryPosition = "CREDIT"
)

func (p EntryPosition) IsValid() bool {
	return p == EntryPositionDebit || p == EntryPositionCredit
}

var (
	ErrJournalNotFound              = errors.New("jurnal transaksi tidak ditemukan")
	ErrJournalUnbalanced            = errors.New("total debit dan kredit jurnal harus sama persis (seimbang)")
	ErrJournalMinimumLines          = errors.New("jurnal minimal harus memiliki 2 baris (1 Debit dan 1 Kredit)")
	ErrJournalNotDraft              = errors.New("operasi ini hanya diizinkan pada jurnal berstatus DRAFT")
	ErrJournalAlreadyPosted         = errors.New("jurnal sudah disahkan (POSTED)")
	ErrJournalAlreadyReversed       = errors.New("jurnal sudah dibatalkan sebelumnya (REVERSED)")
	ErrCannotReverseUnpostedJournal = errors.New("hanya jurnal yang sudah berstatus POSTED yang dapat di-reverse")
	ErrInvalidEntryPosition         = errors.New("posisi baris jurnal harus DEBIT atau CREDIT")
	ErrInvalidAmount                = errors.New("nominal baris jurnal harus lebih besar dari 0")
	ErrReversalReasonRequired       = errors.New("alasan pembatalan (reversal reason) wajib diisi")
	ErrAccountNotPostable           = errors.New("satu atau lebih akun yang dipilih tidak memenuhi syarat penjurnalan (harus aktif dan postable)")
	ErrTrialBalanceCorrupted        = errors.New("transaksi dibatalkan: ketidakseimbangan terdeteksi pada neraca saldo sistem")
	ErrBackdateReferenceRequired    = errors.New("nomor referensi dokumen (reference_number) wajib diisi untuk transaksi backdate")
)

// JournalEntry merepresentasikan header lembar voucher transaksi jurnal.
type JournalEntry struct {
	ID              string             `gorm:"primaryKey;type:varchar(36)" json:"id"`
	EntryNumber     string             `gorm:"type:varchar(50);not null" json:"entry_number"`
	EntryDate       time.Time          `gorm:"type:date;not null" json:"entry_date"`
	FiscalPeriodID  string             `gorm:"type:varchar(36);not null;index" json:"fiscal_period_id"`
	FiscalPeriod    string             `gorm:"type:varchar(7);not null;index" json:"fiscal_period"`
	EntryType       JournalType        `gorm:"type:varchar(30);not null;index" json:"entry_type"`
	Status          JournalStatus      `gorm:"type:varchar(20);not null;default:'DRAFT';index" json:"status"`
	Description     string             `gorm:"type:text;not null" json:"description"`
	ReferenceNumber *string            `gorm:"type:varchar(100)" json:"reference_number,omitempty"`
	TotalDebit      decimal.Decimal    `gorm:"type:numeric(18,2);not null;default:0" json:"total_debit"`
	TotalCredit     decimal.Decimal    `gorm:"type:numeric(18,2);not null;default:0" json:"total_credit"`
	ReversedEntryID *string            `gorm:"type:varchar(36)" json:"reversed_entry_id,omitempty"`
	ReversalReason  *string            `gorm:"type:text" json:"reversal_reason,omitempty"`
	PostedAt        *time.Time         `gorm:"type:timestamptz" json:"posted_at,omitempty"`
	PostedBy        *string            `gorm:"type:varchar(50)" json:"posted_by,omitempty"`
	ReversedAt      *time.Time         `gorm:"type:timestamptz" json:"reversed_at,omitempty"`
	ReversedBy      *string            `gorm:"type:varchar(50)" json:"reversed_by,omitempty"`
	CreatedAt       time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt     `gorm:"index" json:"-"`
	CreatedBy       string             `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"created_by"`
	UpdatedBy       string             `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"updated_by"`
	DeletedBy       *string            `gorm:"type:varchar(50)" json:"deleted_by,omitempty"`
	Lines           []*JournalEntryLine `gorm:"foreignKey:JournalEntryID;constraint:OnDelete:CASCADE" json:"lines,omitempty"`
}

func (JournalEntry) TableName() string {
	return "journal_entries"
}

func (e *JournalEntry) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewV7().String()
	}
	return nil
}

// JournalEntryLine merepresentasikan satu baris debit atau kredit pada voucher jurnal.
type JournalEntryLine struct {
	ID             string          `gorm:"primaryKey;type:varchar(36)" json:"id"`
	JournalEntryID string          `gorm:"type:varchar(36);not null;index" json:"journal_entry_id"`
	LineNumber     int             `gorm:"type:integer;not null" json:"line_number"`
	AccountID      string          `gorm:"type:varchar(36);not null;index" json:"account_id"`
	AccountCode    string          `gorm:"type:varchar(50);not null;index" json:"account_code"`
	AccountName    string          `gorm:"type:varchar(150);not null" json:"account_name"`
	EntryPosition  EntryPosition   `gorm:"type:varchar(10);not null" json:"entry_position"`
	Amount         decimal.Decimal `gorm:"type:numeric(18,2);not null" json:"amount"`
	Memo           *string         `gorm:"type:text" json:"memo,omitempty"`
	DepartmentID   *string         `gorm:"type:varchar(36)" json:"department_id,omitempty"`
	CreatedAt      time.Time       `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (JournalEntryLine) TableName() string {
	return "journal_entry_lines"
}

func (l *JournalEntryLine) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.NewV7().String()
	}
	return nil
}

// RecalculateTotals menghitung ulang total debit dan kredit dari baris rincian.
func (e *JournalEntry) RecalculateTotals() {
	var totalDebit, totalCredit decimal.Decimal
	for _, line := range e.Lines {
		switch line.EntryPosition {
		case EntryPositionDebit:
			totalDebit = totalDebit.Add(line.Amount)
		case EntryPositionCredit:
			totalCredit = totalCredit.Add(line.Amount)
		}
	}
	e.TotalDebit = totalDebit
	e.TotalCredit = totalCredit
}

// Validate memeriksa keabsahan aturan bisnis dasar voucher jurnal.
func (e *JournalEntry) Validate() error {
	if len(e.Lines) < 2 {
		return ErrJournalMinimumLines
	}

	var hasDebit, hasCredit bool
	for _, l := range e.Lines {
		if !l.EntryPosition.IsValid() {
			return ErrInvalidEntryPosition
		}
		if l.Amount.LessThanOrEqual(decimal.Zero) {
			return ErrInvalidAmount
		}
		if l.EntryPosition == EntryPositionDebit {
			hasDebit = true
		}
		if l.EntryPosition == EntryPositionCredit {
			hasCredit = true
		}
	}

	if !hasDebit || !hasCredit {
		return ErrJournalMinimumLines
	}

	e.RecalculateTotals()

	if !e.TotalDebit.Equal(e.TotalCredit) {
		return fmt.Errorf("%w: Total Debit (%s) != Total Kredit (%s)", ErrJournalUnbalanced, e.TotalDebit.StringFixed(2), e.TotalCredit.StringFixed(2))
	}

	return nil
}

// Post mengesahkan draf jurnal menjadi POSTED.
func (e *JournalEntry) Post(actor string) error {
	if e.Status == JournalStatusPosted {
		return ErrJournalAlreadyPosted
	}
	if e.Status != JournalStatusDraft {
		return ErrJournalNotDraft
	}

	if err := e.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	e.Status = JournalStatusPosted
	e.PostedAt = &now
	e.PostedBy = &actor
	e.UpdatedAt = now
	e.UpdatedBy = actor
	return nil
}

// CreateReversalEntry membuat voucher jurnal pembalik (Storno) yang menolkan jurnal ini.
func (e *JournalEntry) CreateReversalEntry(reversalEntryNumber string, reason string, actor string) (*JournalEntry, error) {
	if e.Status == JournalStatusReversed {
		return nil, ErrJournalAlreadyReversed
	}
	if e.Status != JournalStatusPosted {
		return nil, ErrCannotReverseUnpostedJournal
	}
	if reason == "" {
		return nil, ErrReversalReasonRequired
	}

	now := time.Now().UTC()

	// 1. Tandai jurnal asal sebagai REVERSED
	e.Status = JournalStatusReversed
	e.ReversedAt = &now
	e.ReversedBy = &actor
	e.ReversalReason = &reason
	e.UpdatedAt = now
	e.UpdatedBy = actor

	// 2. Buat voucher pembalik baru dengan membalik posisi baris DEBIT <-> CREDIT
	reversalID := uuid.NewV7().String()
	reversedLines := make([]*JournalEntryLine, len(e.Lines))
	for i, l := range e.Lines {
		newPos := EntryPositionDebit
		if l.EntryPosition == EntryPositionDebit {
			newPos = EntryPositionCredit
		}

		lineMemo := fmt.Sprintf("Pembalik baris #%d: %s", l.LineNumber, e.EntryNumber)
		if l.Memo != nil && *l.Memo != "" {
			lineMemo = fmt.Sprintf("Pembalik: %s", *l.Memo)
		}

		reversedLines[i] = &JournalEntryLine{
			ID:             uuid.NewV7().String(),
			JournalEntryID: reversalID,
			LineNumber:     l.LineNumber,
			AccountID:      l.AccountID,
			AccountCode:    l.AccountCode,
			AccountName:    l.AccountName,
			EntryPosition:  newPos,
			Amount:         l.Amount,
			Memo:           &lineMemo,
			DepartmentID:   l.DepartmentID,
			CreatedAt:      now,
		}
	}

	refNum := e.EntryNumber
	revDesc := fmt.Sprintf("Reversal atas voucher %s: %s", e.EntryNumber, reason)

	reversalEntry := &JournalEntry{
		ID:              reversalID,
		EntryNumber:     reversalEntryNumber,
		EntryDate:       now, // Atau tanggal yang sama dengan entri awal
		FiscalPeriodID:  e.FiscalPeriodID,
		FiscalPeriod:    e.FiscalPeriod,
		EntryType:       JournalTypeReversal,
		Status:          JournalStatusPosted, // Reversal langsung POSTED
		Description:     revDesc,
		ReferenceNumber: &refNum,
		TotalDebit:      e.TotalCredit, // Saling tukar
		TotalCredit:     e.TotalDebit,
		ReversedEntryID: &e.ID,
		ReversalReason:  &reason,
		PostedAt:        &now,
		PostedBy:        &actor,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       actor,
		UpdatedBy:       actor,
		Lines:           reversedLines,
	}

	return reversalEntry, nil
}
