package journal

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestJournalEntry_ValidateBalance(t *testing.T) {
	// 1. Valid balanced entry
	entry := &JournalEntry{
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-1",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(500000),
			},
			{
				LineNumber:    2,
				AccountID:     "acc-2",
				EntryPosition: EntryPositionCredit,
				Amount:        decimal.NewFromFloat(500000),
			},
		},
	}

	if err := entry.Validate(); err != nil {
		t.Fatalf("expected valid balanced entry, got error: %v", err)
	}
	if !entry.TotalDebit.Equal(decimal.NewFromFloat(500000)) {
		t.Errorf("expected total debit 500000, got %s", entry.TotalDebit)
	}
	if !entry.TotalCredit.Equal(decimal.NewFromFloat(500000)) {
		t.Errorf("expected total credit 500000, got %s", entry.TotalCredit)
	}

	// 2. Unbalanced entry
	unbalanced := &JournalEntry{
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-1",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(500000),
			},
			{
				LineNumber:    2,
				AccountID:     "acc-2",
				EntryPosition: EntryPositionCredit,
				Amount:        decimal.NewFromFloat(450000), // Selisih 50.000
			},
		},
	}
	if err := unbalanced.Validate(); err == nil {
		t.Fatal("expected error for unbalanced entry, got nil")
	}

	// 3. Minimum lines (hanya 1 baris)
	singleLine := &JournalEntry{
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-1",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(500000),
			},
		},
	}
	if err := singleLine.Validate(); err == nil {
		t.Fatal("expected error for single line, got nil")
	}

	// 4. Invalid amount (<= 0)
	zeroAmount := &JournalEntry{
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-1",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(0),
			},
			{
				LineNumber:    2,
				AccountID:     "acc-2",
				EntryPosition: EntryPositionCredit,
				Amount:        decimal.NewFromFloat(0),
			},
		},
	}
	if err := zeroAmount.Validate(); err == nil {
		t.Fatal("expected error for zero amount, got nil")
	}
}

func TestJournalEntry_Post(t *testing.T) {
	entry := &JournalEntry{
		Status: JournalStatusDraft,
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-1",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(100000),
			},
			{
				LineNumber:    2,
				AccountID:     "acc-2",
				EntryPosition: EntryPositionCredit,
				Amount:        decimal.NewFromFloat(100000),
			},
		},
	}

	// 1. Post from DRAFT -> Success
	actor := "accountant"
	if err := entry.Post(actor); err != nil {
		t.Fatalf("unexpected error posting: %v", err)
	}
	if entry.Status != JournalStatusPosted {
		t.Errorf("expected status POSTED, got %s", entry.Status)
	}
	if entry.PostedBy == nil || *entry.PostedBy != actor {
		t.Errorf("expected posted_by accountant")
	}

	// 2. Post again -> Error ErrJournalAlreadyPosted
	if err := entry.Post(actor); err != ErrJournalAlreadyPosted {
		t.Errorf("expected ErrJournalAlreadyPosted, got %v", err)
	}
}

func TestJournalEntry_CreateReversal(t *testing.T) {
	now := time.Now().UTC()
	original := &JournalEntry{
		ID:             "entry-orig-1",
		EntryNumber:    "JV/202610/00001",
		FiscalPeriodID: "fp-1",
		FiscalPeriod:   "2026-10",
		Status:         JournalStatusPosted,
		PostedAt:       &now,
		Lines: []*JournalEntryLine{
			{
				LineNumber:    1,
				AccountID:     "acc-cash",
				AccountCode:   "1111",
				AccountName:   "Kas Kasir",
				EntryPosition: EntryPositionDebit,
				Amount:        decimal.NewFromFloat(250000),
			},
			{
				LineNumber:    2,
				AccountID:     "acc-rev",
				AccountCode:   "4110",
				AccountName:   "Pendapatan Rawat Jalan",
				EntryPosition: EntryPositionCredit,
				Amount:        decimal.NewFromFloat(250000),
			},
		},
	}

	actor := "manager"
	reason := "Koreksi salah input tarif pasien"
	revNumber := "REV/202610/00001"

	reversal, err := original.CreateReversalEntry(revNumber, reason, actor)
	if err != nil {
		t.Fatalf("unexpected error creating reversal: %v", err)
	}

	// Cek status original
	if original.Status != JournalStatusReversed {
		t.Errorf("expected original status REVERSED, got %s", original.Status)
	}
	if original.ReversalReason == nil || *original.ReversalReason != reason {
		t.Errorf("expected reversal reason stored in original")
	}

	// Cek status reversal voucher
	if reversal.Status != JournalStatusPosted {
		t.Errorf("expected reversal status POSTED, got %s", reversal.Status)
	}
	if reversal.EntryNumber != revNumber {
		t.Errorf("expected reversal number %s, got %s", revNumber, reversal.EntryNumber)
	}
	if len(reversal.Lines) != 2 {
		t.Fatalf("expected 2 lines in reversal, got %d", len(reversal.Lines))
	}

	// Cek pembalikan posisi: Baris 1 yang tadinya Debit sekarang harus Kredit
	line1 := reversal.Lines[0]
	if line1.EntryPosition != EntryPositionCredit {
		t.Errorf("expected line 1 position CREDIT, got %s", line1.EntryPosition)
	}
	if !line1.Amount.Equal(decimal.NewFromFloat(250000)) {
		t.Errorf("expected line 1 amount 250000, got %s", line1.Amount)
	}

	// Baris 2 yang tadinya Kredit sekarang harus Debit
	line2 := reversal.Lines[1]
	if line2.EntryPosition != EntryPositionDebit {
		t.Errorf("expected line 2 position DEBIT, got %s", line2.EntryPosition)
	}

	// Tidak boleh re-reverse
	_, err = original.CreateReversalEntry("REV/202610/00002", "Ulang", actor)
	if err != ErrJournalAlreadyReversed {
		t.Errorf("expected ErrJournalAlreadyReversed, got %v", err)
	}
}
