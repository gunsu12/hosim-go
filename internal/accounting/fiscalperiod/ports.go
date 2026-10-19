package fiscalperiod

import (
	"context"
	"time"
)

// ActivePeriodInfo membawa informasi periode aktif hasil resolusi tanggal.
type ActivePeriodInfo struct {
	PeriodID  string       `json:"period_id"`
	PeriodKey string       `json:"period_key"`
	Name      string       `json:"name"`
	Status    PeriodStatus `json:"status"`
}

// FiscalPeriodGuard adalah port read-only yang diekspos ke modul lain (khususnya journal)
// untuk memvalidasi tanggal transaksi apakah diizinkan untuk posting.
type FiscalPeriodGuard interface {
	// ValidateDateForPosting memeriksa apakah entryDate berada pada periode berstatus OPEN.
	// Jika periode LOCKED atau CLOSED, atau belum terdaftar, mengembalikan error yang sesuai.
	ValidateDateForPosting(ctx context.Context, entryDate time.Time) (*ActivePeriodInfo, error)
}

// Repository adalah abstraksi penyimpanan data tahun dan periode fiskal.
type Repository interface {
	CreateFiscalYearWithPeriods(ctx context.Context, year *FiscalYear, periods []*FiscalPeriod) error
	FindYearByID(ctx context.Context, id string) (*FiscalYear, error)
	FindYearByYear(ctx context.Context, year int) (*FiscalYear, error)
	ListYears(ctx context.Context) ([]*FiscalYear, error)
	UpdateYear(ctx context.Context, year *FiscalYear) error

	FindPeriodByID(ctx context.Context, id string) (*FiscalPeriod, error)
	FindPeriodByKey(ctx context.Context, key string) (*FiscalPeriod, error)
	FindPeriodByDate(ctx context.Context, date time.Time) (*FiscalPeriod, error)
	ListPeriodsByYearID(ctx context.Context, yearID string) ([]*FiscalPeriod, error)
	ListPeriods(ctx context.Context, year *int, status *string) ([]*FiscalPeriod, error)
	UpdatePeriod(ctx context.Context, period *FiscalPeriod) error
}
