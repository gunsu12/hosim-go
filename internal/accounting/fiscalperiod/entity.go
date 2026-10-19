package fiscalperiod

import (
	"errors"
	"fmt"
	"time"
	"uuid"

	"gorm.io/gorm"
)

type PeriodStatus string

const (
	PeriodStatusOpen   PeriodStatus = "OPEN"
	PeriodStatusLocked PeriodStatus = "LOCKED"
	PeriodStatusClosed PeriodStatus = "CLOSED"
)

func (s PeriodStatus) IsValid() bool {
	switch s {
	case PeriodStatusOpen, PeriodStatusLocked, PeriodStatusClosed:
		return true
	default:
		return false
	}
}

var (
	ErrYearAlreadyExists              = errors.New("tahun fiskal sudah terdaftar")
	ErrYearNotFound                   = errors.New("tahun fiskal tidak ditemukan")
	ErrInvalidYear                    = errors.New("tahun fiskal harus antara 2000 dan 2100")
	ErrPeriodNotFound                 = errors.New("periode fiskal tidak ditemukan")
	ErrPeriodAlreadyClosed            = errors.New("periode fiskal sudah ditutup permanen (CLOSED)")
	ErrCannotLockClosedPeriod         = errors.New("tidak dapat mengunci periode yang sudah ditutup")
	ErrPeriodNotLocked                = errors.New("hanya periode berstatus LOCKED yang dapat di-unlock")
	ErrPeriodAlreadyOpen              = errors.New("periode fiskal sudah berstatus OPEN")
	ErrCannotCloseYearWithOpenPeriods = errors.New("seluruh 12 periode bulanan harus berstatus CLOSED sebelum menutup tahun buku")
	ErrYearAlreadyClosed              = errors.New("tahun buku sudah ditutup permanen")
	ErrPeriodLocked                   = errors.New("periode fiskal sedang terkunci sementara (LOCKED)")
	ErrPeriodClosed                   = errors.New("periode fiskal sudah ditutup permanen (CLOSED)")
)

// FiscalYear merepresentasikan tahun buku akuntansi (1 Januari - 31 Desember).
type FiscalYear struct {
	ID        string          `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Year      int             `gorm:"type:integer;uniqueIndex;not null" json:"year"`
	StartDate time.Time       `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time       `gorm:"type:date;not null" json:"end_date"`
	IsClosed  bool            `gorm:"type:boolean;not null;default:false" json:"is_closed"`
	ClosedAt  *time.Time      `gorm:"type:timestamptz" json:"closed_at"`
	ClosedBy  *string         `gorm:"type:varchar(50)" json:"closed_by"`
	Notes     *string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time       `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy string          `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"created_by"`
	UpdatedBy string          `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"updated_by"`
	Periods   []*FiscalPeriod `gorm:"foreignKey:FiscalYearID" json:"periods,omitempty"`
}

func (FiscalYear) TableName() string {
	return "fiscal_years"
}

func (y *FiscalYear) BeforeCreate(tx *gorm.DB) error {
	if y.ID == "" {
		y.ID = uuid.NewV7().String()
	}
	return nil
}

// FiscalPeriod merepresentasikan satu periode akuntansi bulanan (misal: 2026-10).
type FiscalPeriod struct {
	ID           string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	FiscalYearID string       `gorm:"type:varchar(36);not null;index" json:"fiscal_year_id"`
	PeriodNumber int          `gorm:"type:integer;not null" json:"period_number"`
	PeriodKey    string       `gorm:"type:varchar(7);uniqueIndex;not null" json:"period_key"`
	Name         string       `gorm:"type:varchar(100);not null" json:"name"`
	StartDate    time.Time    `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time    `gorm:"type:date;not null" json:"end_date"`
	Status       PeriodStatus `gorm:"type:varchar(20);not null;default:'OPEN';index" json:"status"`
	LockedAt     *time.Time   `gorm:"type:timestamptz" json:"locked_at"`
	LockedBy     *string      `gorm:"type:varchar(50)" json:"locked_by"`
	ClosedAt     *time.Time   `gorm:"type:timestamptz" json:"closed_at"`
	ClosedBy     *string      `gorm:"type:varchar(50)" json:"closed_by"`
	Notes        *string      `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time    `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy    string       `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"created_by"`
	UpdatedBy    string       `gorm:"type:varchar(50);not null;default:'SYSTEM'" json:"updated_by"`
}

func (FiscalPeriod) TableName() string {
	return "fiscal_periods"
}

func (p *FiscalPeriod) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}

// Nama-nama bulan baku dalam Bahasa Indonesia untuk akuntansi RS.
var monthNamesID = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// GenerateFiscalYearWithPeriods membuat entity tahun fiskal baru beserta 12 periode bulanannya.
func GenerateFiscalYearWithPeriods(year int, notes *string, actor string) (*FiscalYear, []*FiscalPeriod, error) {
	if year < 2000 || year > 2100 {
		return nil, nil, ErrInvalidYear
	}
	if actor == "" {
		actor = "SYSTEM"
	}

	yearID := uuid.NewV7().String()
	now := time.Now().UTC()

	// 01 Januari s/d 31 Desember
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)

	fy := &FiscalYear{
		ID:        yearID,
		Year:      year,
		StartDate: startDate,
		EndDate:   endDate,
		IsClosed:  false,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: actor,
		UpdatedBy: actor,
	}

	periods := make([]*FiscalPeriod, 12)
	for m := 1; m <= 12; m++ {
		// Tanggal awal bulan
		pStart := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		// Tanggal akhir bulan (hari ke-0 dari bulan berikutnya = hari terakhir bulan ini, handle kabisat otomatis)
		pEnd := time.Date(year, time.Month(m+1), 0, 0, 0, 0, 0, time.UTC)

		periodKey := fmt.Sprintf("%04d-%02d", year, m)
		periodName := fmt.Sprintf("%s %04d", monthNamesID[m-1], year)

		periods[m-1] = &FiscalPeriod{
			ID:           uuid.NewV7().String(),
			FiscalYearID: yearID,
			PeriodNumber: m,
			PeriodKey:    periodKey,
			Name:         periodName,
			StartDate:    pStart,
			EndDate:      pEnd,
			Status:       PeriodStatusOpen,
			CreatedAt:    now,
			UpdatedAt:    now,
			CreatedBy:    actor,
			UpdatedBy:    actor,
		}
	}

	return fy, periods, nil
}

// Lock mengunci sementara periode bulanan (misal saat review akhir bulan).
func (p *FiscalPeriod) Lock(actor string, notes *string) error {
	if p.Status == PeriodStatusClosed {
		return ErrCannotLockClosedPeriod
	}
	if p.Status == PeriodStatusLocked {
		return nil // Sudah locked, idempotent
	}

	now := time.Now().UTC()
	p.Status = PeriodStatusLocked
	p.LockedAt = &now
	p.LockedBy = &actor
	if notes != nil && *notes != "" {
		p.Notes = notes
	}
	p.UpdatedAt = now
	p.UpdatedBy = actor
	return nil
}

// Unlock membuka kembali periode yang berstatus LOCKED menjadi OPEN.
func (p *FiscalPeriod) Unlock(actor string, notes *string) error {
	if p.Status == PeriodStatusClosed {
		return ErrPeriodAlreadyClosed
	}
	if p.Status == PeriodStatusOpen {
		return nil // Sudah open, idempotent
	}
	if p.Status != PeriodStatusLocked {
		return ErrPeriodNotLocked
	}

	now := time.Now().UTC()
	p.Status = PeriodStatusOpen
	p.LockedAt = nil
	p.LockedBy = nil
	if notes != nil && *notes != "" {
		p.Notes = notes
	}
	p.UpdatedAt = now
	p.UpdatedBy = actor
	return nil
}

// Close menutup buku permanen sebuah periode bulanan. Tidak dapat dibuka kembali selamanya.
func (p *FiscalPeriod) Close(actor string, notes *string) error {
	if p.Status == PeriodStatusClosed {
		return nil // Idempotent
	}

	now := time.Now().UTC()
	p.Status = PeriodStatusClosed
	p.ClosedAt = &now
	p.ClosedBy = &actor
	if notes != nil && *notes != "" {
		p.Notes = notes
	}
	p.UpdatedAt = now
	p.UpdatedBy = actor
	return nil
}

// CloseYear menutup buku tahunan jika dan hanya jika seluruh 12 periode bulanan sudah berstatus CLOSED.
func (y *FiscalYear) CloseYear(periods []*FiscalPeriod, actor string) error {
	if y.IsClosed {
		return ErrYearAlreadyClosed
	}

	if len(periods) < 12 {
		return ErrCannotCloseYearWithOpenPeriods
	}

	for _, p := range periods {
		if p.Status != PeriodStatusClosed {
			return fmt.Errorf("%w: periode %s masih berstatus %s", ErrCannotCloseYearWithOpenPeriods, p.PeriodKey, p.Status)
		}
	}

	now := time.Now().UTC()
	y.IsClosed = true
	y.ClosedAt = &now
	y.ClosedBy = &actor
	y.UpdatedAt = now
	y.UpdatedBy = actor
	return nil
}
