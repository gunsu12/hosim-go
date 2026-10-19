package journal

import (
	"context"
	"fmt"
	"strings"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
)

type Service interface {
	JournalPostingPort

	CreateDraft(ctx context.Context, req CreateJournalRequest, actor string) (*JournalResponse, error)
	UpdateDraft(ctx context.Context, id string, req UpdateJournalRequest, actor string) (*JournalResponse, error)
	DeleteDraft(ctx context.Context, id string, actor string) error
	GetJournal(ctx context.Context, id string) (*JournalResponse, error)
	ListJournals(ctx context.Context, params ListJournalParams) ([]*JournalResponse, int64, error)
	PostJournal(ctx context.Context, id string, actor string) (*JournalResponse, error)
	ReverseJournal(ctx context.Context, id string, req ReverseJournalRequest, actor string) (*JournalResponse, error)
}

type service struct {
	repo              Repository
	accountVerifier   AccountVerifier
	fiscalPeriodGuard FiscalPeriodGuard
}

func NewService(repo Repository, accountVerifier AccountVerifier, fiscalPeriodGuard FiscalPeriodGuard) Service {
	return &service{
		repo:              repo,
		accountVerifier:   accountVerifier,
		fiscalPeriodGuard: fiscalPeriodGuard,
	}
}

func (s *service) CreateDraft(ctx context.Context, req CreateJournalRequest, actor string) (*JournalResponse, error) {
	// 0. Sanitasi input string
	req.Description = strings.TrimSpace(req.Description)
	if req.ReferenceNumber != nil {
		trimmed := strings.TrimSpace(*req.ReferenceNumber)
		if trimmed == "" {
			req.ReferenceNumber = nil
		} else {
			req.ReferenceNumber = &trimmed
		}
	}

	// 1. Parse tanggal transaksi
	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (harus YYYY-MM-DD): %w", err)
	}

	// Validasi backdate: jika tanggal transaksi sebelum hari ini, nomor referensi dokumen fisik wajib diisi
	todayStr := time.Now().Format("2006-01-02")
	if entryDate.Format("2006-01-02") < todayStr && req.ReferenceNumber == nil {
		return nil, ErrBackdateReferenceRequired
	}

	// 2. Validasi periode buku via FiscalPeriodGuard
	periodInfo, err := s.fiscalPeriodGuard.ValidateDateForPosting(ctx, entryDate)
	if err != nil {
		return nil, err
	}

	// 3. Kumpulkan account IDs dan validasi kelayakan akun postable
	accountIDs := make([]string, len(req.Lines))
	for i, l := range req.Lines {
		accountIDs[i] = l.AccountID
	}
	accMap, err := s.accountVerifier.ValidatePostableAccounts(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	// 4. Generate nomor voucher urut berikutnya
	entryNumber, err := s.repo.GenerateNextEntryNumber(ctx, req.EntryType, periodInfo.PeriodKey)
	if err != nil {
		return nil, fmt.Errorf("gagal menghasilkan nomor voucher: %w", err)
	}

	now := time.Now().UTC()
	entryID := uuid.NewV7().String()

	// 5. Bentuk baris-baris rincian
	lines := make([]*JournalEntryLine, len(req.Lines))
	for i, l := range req.Lines {
		acc := accMap[l.AccountID]
		lines[i] = &JournalEntryLine{
			ID:             uuid.NewV7().String(),
			JournalEntryID: entryID,
			LineNumber:     i + 1,
			AccountID:      acc.ID,
			AccountCode:    acc.Code,
			AccountName:    acc.Name,
			EntryPosition:  l.EntryPosition,
			Amount:         decimal.NewFromFloat(l.Amount),
			Memo:           l.Memo,
			DepartmentID:   l.DepartmentID,
			CreatedAt:      now,
		}
	}

	entry := &JournalEntry{
		ID:              entryID,
		EntryNumber:     entryNumber,
		EntryDate:       entryDate,
		FiscalPeriodID:  periodInfo.PeriodID,
		FiscalPeriod:    periodInfo.PeriodKey,
		EntryType:       req.EntryType,
		Status:          JournalStatusDraft,
		Description:     req.Description,
		ReferenceNumber: req.ReferenceNumber,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       actor,
		UpdatedBy:       actor,
		Lines:           lines,
	}

	// 6. Ring 1: Validasi keseimbangan debit-kredit
	if err := entry.Validate(); err != nil {
		return nil, err
	}

	// 7. Simpan draf
	if err := s.repo.SaveDraft(ctx, entry); err != nil {
		return nil, err
	}

	return ToJournalResponse(entry), nil
}

func (s *service) UpdateDraft(ctx context.Context, id string, req UpdateJournalRequest, actor string) (*JournalResponse, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if entry.Status != JournalStatusDraft {
		return nil, ErrJournalNotDraft
	}

	now := time.Now().UTC()

	// 1. Update tanggal jika disediakan
	if req.EntryDate != nil && *req.EntryDate != "" {
		newDate, err := time.Parse("2006-01-02", *req.EntryDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal tidak valid: %w", err)
		}
		periodInfo, err := s.fiscalPeriodGuard.ValidateDateForPosting(ctx, newDate)
		if err != nil {
			return nil, err
		}

		// Jika periode berubah karena tanggal baru, sesuaikan nomor urut voucher ke periode tujuan
		if entry.FiscalPeriod != periodInfo.PeriodKey {
			newEntryNumber, err := s.repo.GenerateNextEntryNumber(ctx, entry.EntryType, periodInfo.PeriodKey)
			if err != nil {
				return nil, fmt.Errorf("gagal menghasilkan nomor voucher baru untuk periode %s: %w", periodInfo.PeriodKey, err)
			}
			entry.EntryNumber = newEntryNumber
		}

		entry.EntryDate = newDate
		entry.FiscalPeriodID = periodInfo.PeriodID
		entry.FiscalPeriod = periodInfo.PeriodKey
	}

	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed != "" {
			entry.Description = trimmed
		}
	}
	if req.ReferenceNumber != nil {
		trimmed := strings.TrimSpace(*req.ReferenceNumber)
		if trimmed == "" {
			entry.ReferenceNumber = nil
		} else {
			entry.ReferenceNumber = &trimmed
		}
	}

	// Validasi backdate: jika tanggal transaksi sebelum hari ini, nomor referensi dokumen fisik wajib diisi
	todayStr := time.Now().Format("2006-01-02")
	if entry.EntryDate.Format("2006-01-02") < todayStr && entry.ReferenceNumber == nil {
		return nil, ErrBackdateReferenceRequired
	}

	// 2. Update baris jika disediakan
	if len(req.Lines) > 0 {
		accountIDs := make([]string, len(req.Lines))
		for i, l := range req.Lines {
			accountIDs[i] = l.AccountID
		}
		accMap, err := s.accountVerifier.ValidatePostableAccounts(ctx, accountIDs)
		if err != nil {
			return nil, err
		}

		newLines := make([]*JournalEntryLine, len(req.Lines))
		for i, l := range req.Lines {
			acc := accMap[l.AccountID]
			newLines[i] = &JournalEntryLine{
				ID:             uuid.NewV7().String(),
				JournalEntryID: entry.ID,
				LineNumber:     i + 1,
				AccountID:      acc.ID,
				AccountCode:    acc.Code,
				AccountName:    acc.Name,
				EntryPosition:  l.EntryPosition,
				Amount:         decimal.NewFromFloat(l.Amount),
				Memo:           l.Memo,
				DepartmentID:   l.DepartmentID,
				CreatedAt:      now,
			}
		}
		entry.Lines = newLines
	}

	entry.UpdatedAt = now
	entry.UpdatedBy = actor

	// 3. Ring 1: Validasi keseimbangan
	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateDraft(ctx, entry); err != nil {
		return nil, err
	}

	return ToJournalResponse(entry), nil
}

func (s *service) DeleteDraft(ctx context.Context, id string, actor string) error {
	return s.repo.DeleteDraft(ctx, id, actor)
}

func (s *service) GetJournal(ctx context.Context, id string) (*JournalResponse, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToJournalResponse(entry), nil
}

func (s *service) ListJournals(ctx context.Context, params ListJournalParams) ([]*JournalResponse, int64, error) {
	entries, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return ToJournalResponseList(entries), total, nil
}

func (s *service) PostJournal(ctx context.Context, id string, actor string) (*JournalResponse, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 1. Verifikasi kembali apakah periode tanggal transaksi masih OPEN saat mau posting
	if _, err := s.fiscalPeriodGuard.ValidateDateForPosting(ctx, entry.EntryDate); err != nil {
		return nil, err
	}

	// 2. Verifikasi ulang akun postable (mencegah akun dinonaktifkan setelah draft dibuat)
	accountIDs := make([]string, len(entry.Lines))
	for i, l := range entry.Lines {
		accountIDs[i] = l.AccountID
	}
	if _, err := s.accountVerifier.ValidatePostableAccounts(ctx, accountIDs); err != nil {
		return nil, err
	}

	// 3. Ubah status menjadi POSTED
	if err := entry.Post(actor); err != nil {
		return nil, err
	}

	// 4. Ring 2 & Ring 3: Simpan transaksi & cek keseimbangan global neraca saldo
	if err := s.repo.PostEntry(ctx, entry); err != nil {
		return nil, err
	}

	return ToJournalResponse(entry), nil
}

func (s *service) ReverseJournal(ctx context.Context, id string, req ReverseJournalRequest, actor string) (*JournalResponse, error) {
	req.Reason = strings.TrimSpace(req.Reason)

	original, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	// 1. Pastikan periode saat ini OPEN untuk menampung jurnal reversal
	periodInfo, err := s.fiscalPeriodGuard.ValidateDateForPosting(ctx, now)
	if err != nil {
		return nil, err
	}

	// 2. Generate nomor voucher pembalik
	reversalNumber, err := s.repo.GenerateNextEntryNumber(ctx, JournalTypeReversal, periodInfo.PeriodKey)
	if err != nil {
		return nil, fmt.Errorf("gagal menghasilkan nomor voucher reversal: %w", err)
	}

	// 3. Buat entity reversal
	reversalEntry, err := original.CreateReversalEntry(reversalNumber, req.Reason, actor)
	if err != nil {
		return nil, err
	}
	reversalEntry.FiscalPeriodID = periodInfo.PeriodID
	reversalEntry.FiscalPeriod = periodInfo.PeriodKey

	// 4. Simpan pembalik & update original (Ring 2 & Ring 3 check)
	if err := s.repo.ReverseEntry(ctx, original, reversalEntry); err != nil {
		return nil, err
	}

	return ToJournalResponse(reversalEntry), nil
}

// PostSystemJournal mengimplementasikan JournalPostingPort untuk modul lain (Billing, Kasir, Farmasi).
func (s *service) PostSystemJournal(ctx context.Context, cmd PostSystemJournalCommand) (*JournalEntry, error) {
	// 0. Sanitasi input
	cmd.Description = strings.TrimSpace(cmd.Description)
	if cmd.ReferenceNumber != nil {
		trimmed := strings.TrimSpace(*cmd.ReferenceNumber)
		if trimmed == "" {
			cmd.ReferenceNumber = nil
		} else {
			cmd.ReferenceNumber = &trimmed
		}
	}

	// 1. Validasi periode buku
	periodInfo, err := s.fiscalPeriodGuard.ValidateDateForPosting(ctx, cmd.EntryDate)
	if err != nil {
		return nil, err
	}

	// Validasi backdate pada system journal
	todayStr := time.Now().Format("2006-01-02")
	if cmd.EntryDate.Format("2006-01-02") < todayStr && cmd.ReferenceNumber == nil {
		return nil, ErrBackdateReferenceRequired
	}

	// 2. Validasi akun
	accountIDs := make([]string, len(cmd.Lines))
	for i, l := range cmd.Lines {
		accountIDs[i] = l.AccountID
	}
	accMap, err := s.accountVerifier.ValidatePostableAccounts(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	// 3. Generate nomor voucher
	entryNumber, err := s.repo.GenerateNextEntryNumber(ctx, cmd.EntryType, periodInfo.PeriodKey)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	entryID := uuid.NewV7().String()

	lines := make([]*JournalEntryLine, len(cmd.Lines))
	for i, l := range cmd.Lines {
		acc := accMap[l.AccountID]
		lines[i] = &JournalEntryLine{
			ID:             uuid.NewV7().String(),
			JournalEntryID: entryID,
			LineNumber:     i + 1,
			AccountID:      acc.ID,
			AccountCode:    acc.Code,
			AccountName:    acc.Name,
			EntryPosition:  l.EntryPosition,
			Amount:         decimal.NewFromFloat(l.Amount),
			Memo:           l.Memo,
			DepartmentID:   l.DepartmentID,
			CreatedAt:      now,
		}
	}

	entry := &JournalEntry{
		ID:              entryID,
		EntryNumber:     entryNumber,
		EntryDate:       cmd.EntryDate,
		FiscalPeriodID:  periodInfo.PeriodID,
		FiscalPeriod:    periodInfo.PeriodKey,
		EntryType:       cmd.EntryType,
		Status:          JournalStatusPosted, // Langsung disahkan
		Description:     cmd.Description,
		ReferenceNumber: cmd.ReferenceNumber,
		PostedAt:        &now,
		PostedBy:        &cmd.Actor,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       cmd.Actor,
		UpdatedBy:       cmd.Actor,
		Lines:           lines,
	}

	// Ring 1 check
	if err := entry.Validate(); err != nil {
		return nil, err
	}

	// Simpan & jalankan Ring 3 check secara atomik
	if err := s.repo.CreatePostedEntry(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}
