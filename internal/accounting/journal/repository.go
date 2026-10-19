package journal

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveDraft(ctx context.Context, entry *JournalEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *repository) UpdateDraft(ctx context.Context, entry *JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi status harus DRAFT dengan row-level lock (FOR UPDATE)
		var existing JournalEntry
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, "id = ?", entry.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrJournalNotFound
			}
			return err
		}
		if existing.Status != JournalStatusDraft {
			return ErrJournalNotDraft
		}

		// 2. Hapus baris lama
		if err := tx.Where("journal_entry_id = ?", entry.ID).Delete(&JournalEntryLine{}).Error; err != nil {
			return err
		}

		// 3. Simpan baris baru
		if len(entry.Lines) > 0 {
			if err := tx.Create(&entry.Lines).Error; err != nil {
				return err
			}
		}

		// 4. Update header tanpa menyentuh asosiasi
		return tx.Omit(clause.Associations).Save(entry).Error
	})
}

func (r *repository) DeleteDraft(ctx context.Context, id string, actor string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi status harus DRAFT dengan row-level lock (FOR UPDATE)
		var existing JournalEntry
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrJournalNotFound
			}
			return err
		}
		if existing.Status != JournalStatusDraft {
			return ErrJournalNotDraft
		}

		// 2. Bersihkan baris rincian draf secara fisik agar tidak meninggalkan orphan rows
		if err := tx.Where("journal_entry_id = ?", id).Delete(&JournalEntryLine{}).Error; err != nil {
			return err
		}

		// 3. Tandai soft-delete secara atomik pada header
		now := time.Now().UTC()
		return tx.Model(&JournalEntry{}).
			Where("id = ? AND status = ?", id, JournalStatusDraft).
			Updates(map[string]any{
				"deleted_at": gorm.DeletedAt{Time: now, Valid: true},
				"deleted_by": actor,
				"updated_at": now,
				"updated_by": actor,
			}).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*JournalEntry, error) {
	var entry JournalEntry
	err := r.db.WithContext(ctx).
		Preload("Lines", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_number ASC")
		}).
		First(&entry, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}
	return &entry, nil
}

func (r *repository) FindByEntryNumber(ctx context.Context, entryNumber string) (*JournalEntry, error) {
	var entry JournalEntry
	err := r.db.WithContext(ctx).
		Preload("Lines", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_number ASC")
		}).
		First(&entry, "entry_number = ?", entryNumber).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}
	return &entry, nil
}

func (r *repository) List(ctx context.Context, params ListJournalParams) ([]*JournalEntry, int64, error) {
	query := r.db.WithContext(ctx).Model(&JournalEntry{})

	if params.StartDate != nil {
		dateStr := params.StartDate.Format("2006-01-02")
		query = query.Where("entry_date >= ?", dateStr)
	}
	if params.EndDate != nil {
		dateStr := params.EndDate.Format("2006-01-02")
		query = query.Where("entry_date <= ?", dateStr)
	}
	if params.FiscalPeriod != nil && *params.FiscalPeriod != "" {
		query = query.Where("fiscal_period = ?", *params.FiscalPeriod)
	}
	if params.Status != nil && *params.Status != "" {
		query = query.Where("status = ?", *params.Status)
	}
	if params.EntryType != nil && *params.EntryType != "" {
		query = query.Where("entry_type = ?", *params.EntryType)
	}
	if params.Search != "" {
		s := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(entry_number) LIKE ? OR LOWER(description) LIKE ? OR LOWER(reference_number) LIKE ?", s, s, s)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var entries []*JournalEntry
	err := query.
		Preload("Lines", func(db *gorm.DB) *gorm.DB {
			return db.Order("line_number ASC")
		}).
		Order("entry_date DESC, entry_number DESC").
		Offset(offset).
		Limit(limit).
		Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *repository) PostEntry(ctx context.Context, entry *JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi periode fiskal masih berstatus OPEN di dalam transaksi ini (mencegah TOCTOU race)
		var periodStatus string
		if err := tx.Table("fiscal_periods").Select("status").Where("id = ?", entry.FiscalPeriodID).Scan(&periodStatus).Error; err != nil {
			return err
		}
		if periodStatus != "OPEN" {
			return fmt.Errorf("transaksi dibatalkan: periode fiskal %s tidak berstatus OPEN (status saat ini: %s)", entry.FiscalPeriod, periodStatus)
		}

		// 2. Simpan perubahan status menjadi POSTED secara atomik bersyarat (mencegah TOCTOU race)
		res := tx.Model(&JournalEntry{}).
			Where("id = ? AND status = ?", entry.ID, JournalStatusDraft).
			Updates(map[string]any{
				"status":     JournalStatusPosted,
				"posted_at":  entry.PostedAt,
				"posted_by":  entry.PostedBy,
				"updated_at": entry.UpdatedAt,
				"updated_by": entry.UpdatedBy,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrJournalNotDraft
		}

		// 2. Ring 3: Sanity Check Neraca Saldo Global
		var discrepancy decimal.Decimal
		err := tx.Raw(`
			SELECT COALESCE(ROUND(SUM(CASE WHEN jel.entry_position = 'DEBIT' THEN jel.amount ELSE -jel.amount END), 4), 0)
			FROM journal_entry_lines jel
			JOIN journal_entries je ON je.id = jel.journal_entry_id
			WHERE je.status = 'POSTED' AND je.deleted_at IS NULL
		`).Scan(&discrepancy).Error
		if err != nil {
			return err
		}

		if !discrepancy.IsZero() {
			return fmt.Errorf("%w: selisih global buku besar = %s", ErrTrialBalanceCorrupted, discrepancy.StringFixed(4))
		}

		return nil
	})
}

func (r *repository) CreatePostedEntry(ctx context.Context, entry *JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi periode fiskal berstatus OPEN di dalam transaksi ini
		var periodStatus string
		if err := tx.Table("fiscal_periods").Select("status").Where("id = ?", entry.FiscalPeriodID).Scan(&periodStatus).Error; err != nil {
			return err
		}
		if periodStatus != "OPEN" {
			return fmt.Errorf("transaksi dibatalkan: periode fiskal %s tidak berstatus OPEN (status saat ini: %s)", entry.FiscalPeriod, periodStatus)
		}

		// 2. Simpan voucher beserta baris rinciannya
		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		// 3. Ring 3: Sanity Check Neraca Saldo Global
		var discrepancy decimal.Decimal
		err := tx.Raw(`
			SELECT COALESCE(ROUND(SUM(CASE WHEN jel.entry_position = 'DEBIT' THEN jel.amount ELSE -jel.amount END), 4), 0)
			FROM journal_entry_lines jel
			JOIN journal_entries je ON je.id = jel.journal_entry_id
			WHERE je.status = 'POSTED' AND je.deleted_at IS NULL
		`).Scan(&discrepancy).Error
		if err != nil {
			return err
		}

		if !discrepancy.IsZero() {
			return fmt.Errorf("%w: selisih global buku besar = %s", ErrTrialBalanceCorrupted, discrepancy.StringFixed(4))
		}

		return nil
	})
}

func (r *repository) ReverseEntry(ctx context.Context, original *JournalEntry, reversal *JournalEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi periode reversal berstatus OPEN di dalam transaksi ini
		var periodStatus string
		if err := tx.Table("fiscal_periods").Select("status").Where("id = ?", reversal.FiscalPeriodID).Scan(&periodStatus).Error; err != nil {
			return err
		}
		if periodStatus != "OPEN" {
			return fmt.Errorf("transaksi dibatalkan: periode fiskal %s tidak berstatus OPEN (status saat ini: %s)", reversal.FiscalPeriod, periodStatus)
		}

		// 2. Update status dan audit metadata jurnal asal menjadi REVERSED secara atomik bersyarat
		res := tx.Model(&JournalEntry{}).
			Where("id = ? AND status = ?", original.ID, JournalStatusPosted).
			Updates(map[string]any{
				"status":          JournalStatusReversed,
				"reversed_at":     original.ReversedAt,
				"reversed_by":     original.ReversedBy,
				"reversal_reason": original.ReversalReason,
				"updated_at":      original.UpdatedAt,
				"updated_by":      original.UpdatedBy,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrCannotReverseUnpostedJournal
		}

		// 3. Simpan voucher reversal beserta baris rinciannya
		if err := tx.Create(reversal).Error; err != nil {
			return err
		}

		// 3. Ring 3: Sanity Check Neraca Saldo Global
		var discrepancy decimal.Decimal
		err := tx.Raw(`
			SELECT COALESCE(ROUND(SUM(CASE WHEN jel.entry_position = 'DEBIT' THEN jel.amount ELSE -jel.amount END), 4), 0)
			FROM journal_entry_lines jel
			JOIN journal_entries je ON je.id = jel.journal_entry_id
			WHERE je.status = 'POSTED' AND je.deleted_at IS NULL
		`).Scan(&discrepancy).Error
		if err != nil {
			return err
		}

		if !discrepancy.IsZero() {
			return fmt.Errorf("%w: selisih global buku besar setelah reversal = %s", ErrTrialBalanceCorrupted, discrepancy.StringFixed(4))
		}

		return nil
	})
}

func (r *repository) GenerateNextEntryNumber(ctx context.Context, entryType JournalType, periodKey string) (string, error) {
	// Tentukan prefix berdasarkan tipe jurnal
	var prefix string
	switch entryType {
	case JournalTypeGeneral:
		prefix = "JV"
	case JournalTypeAdjustment:
		prefix = "ADJ"
	case JournalTypeClosing:
		prefix = "CLS"
	case JournalTypeReversal:
		prefix = "REV"
	case JournalTypeSystemBilling:
		prefix = "BIL"
	case JournalTypeSystemCashier:
		prefix = "CSH"
	case JournalTypeSystemInventory:
		prefix = "INV"
	default:
		prefix = "JV"
	}

	cleanPeriod := strings.ReplaceAll(periodKey, "-", "")
	pattern := fmt.Sprintf("%s/%s/%%", prefix, cleanPeriod)

	var lastEntry JournalEntry
	err := r.db.WithContext(ctx).Model(&JournalEntry{}).
		Unscoped().
		Where("entry_number LIKE ?", pattern).
		Order("entry_number DESC").
		Limit(1).
		Find(&lastEntry).Error
	if err != nil {
		return "", err
	}

	var nextSeq int64 = 1
	if lastEntry.EntryNumber != "" {
		parts := strings.Split(lastEntry.EntryNumber, "/")
		if len(parts) == 3 {
			if seq, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
				nextSeq = seq + 1
			}
		}
	}
	return fmt.Sprintf("%s/%s/%05d", prefix, cleanPeriod, nextSeq), nil
}

func (r *repository) CheckTrialBalanceDiscrepancy(ctx context.Context) (decimal.Decimal, error) {
	var discrepancy decimal.Decimal
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(ROUND(SUM(CASE WHEN jel.entry_position = 'DEBIT' THEN jel.amount ELSE -jel.amount END), 4), 0)
		FROM journal_entry_lines jel
		JOIN journal_entries je ON je.id = jel.journal_entry_id
		WHERE je.status = 'POSTED' AND je.deleted_at IS NULL
	`).Scan(&discrepancy).Error
	if err != nil {
		return decimal.Zero, err
	}
	return discrepancy, nil
}
