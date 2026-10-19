package journal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"hosim-go/internal/accounting/fiscalperiod"
)

type mockRepository struct {
	entries map[string]*JournalEntry
	counter int
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		entries: make(map[string]*JournalEntry),
	}
}

func (m *mockRepository) SaveDraft(ctx context.Context, entry *JournalEntry) error {
	m.entries[entry.ID] = entry
	return nil
}

func (m *mockRepository) UpdateDraft(ctx context.Context, entry *JournalEntry) error {
	m.entries[entry.ID] = entry
	return nil
}

func (m *mockRepository) DeleteDraft(ctx context.Context, id string, actor string) error {
	delete(m.entries, id)
	return nil
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (*JournalEntry, error) {
	if e, ok := m.entries[id]; ok {
		cp := *e
		return &cp, nil
	}
	return nil, ErrJournalNotFound
}

func (m *mockRepository) FindByEntryNumber(ctx context.Context, entryNumber string) (*JournalEntry, error) {
	for _, e := range m.entries {
		if e.EntryNumber == entryNumber {
			return e, nil
		}
	}
	return nil, ErrJournalNotFound
}

func (m *mockRepository) List(ctx context.Context, params ListJournalParams) ([]*JournalEntry, int64, error) {
	res := make([]*JournalEntry, 0, len(m.entries))
	for _, e := range m.entries {
		res = append(res, e)
	}
	return res, int64(len(res)), nil
}

func (m *mockRepository) PostEntry(ctx context.Context, entry *JournalEntry) error {
	if e, ok := m.entries[entry.ID]; !ok || e.Status != JournalStatusDraft {
		return ErrJournalNotDraft
	}
	m.entries[entry.ID] = entry
	return nil
}

func (m *mockRepository) CreatePostedEntry(ctx context.Context, entry *JournalEntry) error {
	m.entries[entry.ID] = entry
	return nil
}

func (m *mockRepository) ReverseEntry(ctx context.Context, original *JournalEntry, reversal *JournalEntry) error {
	if orig, ok := m.entries[original.ID]; !ok || orig.Status != JournalStatusPosted {
		return ErrCannotReverseUnpostedJournal
	}
	m.entries[original.ID] = original
	m.entries[reversal.ID] = reversal
	return nil
}

func (m *mockRepository) GenerateNextEntryNumber(ctx context.Context, entryType JournalType, periodKey string) (string, error) {
	m.counter++
	cleanPeriod := strings.ReplaceAll(periodKey, "-", "")
	return fmt.Sprintf("JV/%s/%05d", cleanPeriod, m.counter), nil
}

func (m *mockRepository) CheckTrialBalanceDiscrepancy(ctx context.Context) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

type mockAccountVerifier struct {
	accounts map[string]*AccountSummary
}

func (m *mockAccountVerifier) ValidatePostableAccounts(ctx context.Context, accountIDs []string) (map[string]*AccountSummary, error) {
	res := make(map[string]*AccountSummary)
	for _, id := range accountIDs {
		if acc, ok := m.accounts[id]; ok {
			if !acc.IsActive || !acc.IsPostable {
				return nil, ErrAccountNotPostable
			}
			res[id] = acc
		} else {
			return nil, ErrAccountNotPostable
		}
	}
	return res, nil
}

type mockFiscalPeriodGuard struct {
	status fiscalperiod.PeriodStatus
}

func (m *mockFiscalPeriodGuard) ValidateDateForPosting(ctx context.Context, entryDate time.Time) (*fiscalperiod.ActivePeriodInfo, error) {
	if m.status == fiscalperiod.PeriodStatusLocked {
		return nil, fiscalperiod.ErrPeriodLocked
	}
	if m.status == fiscalperiod.PeriodStatusClosed {
		return nil, fiscalperiod.ErrPeriodClosed
	}
	pKey := entryDate.Format("2006-01")
	return &fiscalperiod.ActivePeriodInfo{
		PeriodID:  "fp-" + pKey,
		PeriodKey: pKey,
		Name:      "Periode " + pKey,
		Status:    m.status,
	}, nil
}

func TestService_CreateDraftAndPost(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()

	accVerifier := &mockAccountVerifier{
		accounts: map[string]*AccountSummary{
			"acc-cash": {ID: "acc-cash", Code: "1111", Name: "Kas Kasir", IsActive: true, IsPostable: true},
			"acc-rev":  {ID: "acc-rev", Code: "4110", Name: "Pendapatan Rawat Jalan", IsActive: true, IsPostable: true},
		},
	}
	periodGuard := &mockFiscalPeriodGuard{status: fiscalperiod.PeriodStatusOpen}

	svc := NewService(repo, accVerifier, periodGuard)

	// 1. Buat draf jurnal
	refNum := "REF-202610-001"
	req := CreateJournalRequest{
		EntryDate:       "2026-10-19",
		EntryType:       JournalTypeGeneral,
		Description:     "Penerimaan kasir rawat jalan",
		ReferenceNumber: &refNum,
		Lines: []JournalLineInput{
			{
				AccountID:     "acc-cash",
				EntryPosition: EntryPositionDebit,
				Amount:        300000,
			},
			{
				AccountID:     "acc-rev",
				EntryPosition: EntryPositionCredit,
				Amount:        300000,
			},
		},
	}

	created, err := svc.CreateDraft(ctx, req, "user-kasir")
	if err != nil {
		t.Fatalf("unexpected error creating draft: %v", err)
	}

	if created.Status != JournalStatusDraft {
		t.Errorf("expected status DRAFT, got %s", created.Status)
	}
	if created.TotalDebit != 300000 || created.TotalCredit != 300000 {
		t.Errorf("expected totals 300000, got debit=%f, credit=%f", created.TotalDebit, created.TotalCredit)
	}

	// 2. Post jurnal
	posted, err := svc.PostJournal(ctx, created.ID, "user-supervisor")
	if err != nil {
		t.Fatalf("unexpected error posting journal: %v", err)
	}

	if posted.Status != JournalStatusPosted {
		t.Errorf("expected status POSTED, got %s", posted.Status)
	}
	if posted.PostedBy == nil || *posted.PostedBy != "user-supervisor" {
		t.Errorf("expected posted_by user-supervisor")
	}

	// 3. Reverse jurnal
	revReq := ReverseJournalRequest{
		Reason: "Kuitansi salah ketik",
	}
	reversed, err := svc.ReverseJournal(ctx, created.ID, revReq, "user-manager")
	if err != nil {
		t.Fatalf("unexpected error reversing journal: %v", err)
	}

	if reversed.Status != JournalStatusPosted {
		t.Errorf("expected reversal status POSTED, got %s", reversed.Status)
	}
	if reversed.EntryType != JournalTypeReversal {
		t.Errorf("expected reversal type REVERSAL, got %s", reversed.EntryType)
	}

	// Pastikan original status sekarang REVERSED
	orig, _ := svc.GetJournal(ctx, created.ID)
	if orig.Status != JournalStatusReversed {
		t.Errorf("expected original status REVERSED, got %s", orig.Status)
	}
}

func TestService_PostSystemJournal_Port(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()

	accVerifier := &mockAccountVerifier{
		accounts: map[string]*AccountSummary{
			"acc-piutang": {ID: "acc-piutang", Code: "1121", Name: "Piutang Pasien", IsActive: true, IsPostable: true},
			"acc-rev-rs":  {ID: "acc-rev-rs", Code: "4101", Name: "Pendapatan RS", IsActive: true, IsPostable: true},
			"acc-rev-dr":  {ID: "acc-rev-dr", Code: "2120", Name: "Hutang Jasa Dokter", IsActive: true, IsPostable: true},
		},
	}
	periodGuard := &mockFiscalPeriodGuard{status: fiscalperiod.PeriodStatusOpen}

	svc := NewService(repo, accVerifier, periodGuard)

	refNum := "BILL-202610-0012"
	cmd := PostSystemJournalCommand{
		EntryDate:       time.Date(2026, time.October, 19, 0, 0, 0, 0, time.UTC),
		EntryType:       JournalTypeSystemBilling,
		Description:     "Charge tindakan IGD pasien Tn. Budi",
		ReferenceNumber: &refNum,
		Lines: []SystemJournalLine{
			{
				AccountID:     "acc-piutang",
				EntryPosition: EntryPositionDebit,
				Amount:        500000,
			},
			{
				AccountID:     "acc-rev-rs",
				EntryPosition: EntryPositionCredit,
				Amount:        200000,
			},
			{
				AccountID:     "acc-rev-dr",
				EntryPosition: EntryPositionCredit,
				Amount:        300000,
			},
		},
		Actor: "SYSTEM_BILLING",
	}

	postedEntry, err := svc.PostSystemJournal(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error posting system journal: %v", err)
	}

	if postedEntry.Status != JournalStatusPosted {
		t.Errorf("expected status POSTED, got %s", postedEntry.Status)
	}
	if len(postedEntry.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(postedEntry.Lines))
	}
}

func TestService_BackdateValidation(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()

	accVerifier := &mockAccountVerifier{
		accounts: map[string]*AccountSummary{
			"acc-cash": {ID: "acc-cash", Code: "1111", Name: "Kas Kasir", IsActive: true, IsPostable: true},
			"acc-rev":  {ID: "acc-rev", Code: "4110", Name: "Pendapatan Rawat Jalan", IsActive: true, IsPostable: true},
		},
	}
	periodGuard := &mockFiscalPeriodGuard{status: fiscalperiod.PeriodStatusOpen}
	svc := NewService(repo, accVerifier, periodGuard)

	pastDate := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	// 1. Coba input backdate tanpa reference_number -> harus gagal
	reqWithoutRef := CreateJournalRequest{
		EntryDate:   pastDate,
		EntryType:   JournalTypeGeneral,
		Description: "Beli ATK 3 hari lalu tanpa nota",
		Lines: []JournalLineInput{
			{AccountID: "acc-cash", EntryPosition: EntryPositionDebit, Amount: 100000},
			{AccountID: "acc-rev", EntryPosition: EntryPositionCredit, Amount: 100000},
		},
	}
	_, err := svc.CreateDraft(ctx, reqWithoutRef, "user-tester")
	if !errors.Is(err, ErrBackdateReferenceRequired) {
		t.Fatalf("expected ErrBackdateReferenceRequired, got %v", err)
	}

	// 2. Input backdate dengan reference_number -> harus sukses
	ref := "NOTA-ATK-001"
	reqWithRef := reqWithoutRef
	reqWithRef.ReferenceNumber = &ref

	draft, err := svc.CreateDraft(ctx, reqWithRef, "user-tester")
	if err != nil {
		t.Fatalf("unexpected error creating backdate draft with reference: %v", err)
	}
	if draft.ReferenceNumber == nil || *draft.ReferenceNumber != ref {
		t.Errorf("expected reference number %s, got %v", ref, draft.ReferenceNumber)
	}

	// 3. Update draft mengubah ke tanggal backdate tetapi menghapus referensi -> harus ditolak
	emptyRef := "   "
	updateReq := UpdateJournalRequest{
		ReferenceNumber: &emptyRef,
	}
	_, err = svc.UpdateDraft(ctx, draft.ID, updateReq, "user-tester")
	if !errors.Is(err, ErrBackdateReferenceRequired) {
		t.Fatalf("expected ErrBackdateReferenceRequired on update without ref, got %v", err)
	}
}

func TestService_UpdateDraftFiscalPeriodChange(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()

	accVerifier := &mockAccountVerifier{
		accounts: map[string]*AccountSummary{
			"acc-cash": {ID: "acc-cash", Code: "1111", Name: "Kas Kasir", IsActive: true, IsPostable: true},
			"acc-rev":  {ID: "acc-rev", Code: "4110", Name: "Pendapatan Rawat Jalan", IsActive: true, IsPostable: true},
		},
	}
	periodGuard := &mockFiscalPeriodGuard{status: fiscalperiod.PeriodStatusOpen}
	svc := NewService(repo, accVerifier, periodGuard)

	// Buat draf untuk bulan depan (misal: 2026-11-05)
	futureDate := "2026-11-05"
	req := CreateJournalRequest{
		EntryDate:   futureDate,
		EntryType:   JournalTypeGeneral,
		Description: "Draf jurnal masa depan",
		Lines: []JournalLineInput{
			{AccountID: "acc-cash", EntryPosition: EntryPositionDebit, Amount: 50000},
			{AccountID: "acc-rev", EntryPosition: EntryPositionCredit, Amount: 50000},
		},
	}
	draft, err := svc.CreateDraft(ctx, req, "user-tester")
	if err != nil {
		t.Fatalf("failed to create draft: %v", err)
	}
	if !strings.HasPrefix(draft.EntryNumber, "JV/202611/") {
		t.Errorf("expected EntryNumber prefix JV/202611/, got %s", draft.EntryNumber)
	}

	// Update draf di-backdate ke 2026-10-15 (periode 2026-10)
	backdate := "2026-10-15"
	ref := "BA-202610-REVISI"
	updateReq := UpdateJournalRequest{
		EntryDate:       &backdate,
		ReferenceNumber: &ref,
	}
	updated, err := svc.UpdateDraft(ctx, draft.ID, updateReq, "user-tester")
	if err != nil {
		t.Fatalf("failed to update draft: %v", err)
	}

	// Pastikan periode dan nomor voucher otomatis disesuaikan ke periode baru 202610
	if updated.FiscalPeriod != "2026-10" {
		t.Errorf("expected fiscal period 2026-10, got %s", updated.FiscalPeriod)
	}
	if !strings.HasPrefix(updated.EntryNumber, "JV/202610/") {
		t.Errorf("expected regenerated EntryNumber prefix JV/202610/, got %s", updated.EntryNumber)
	}
}
