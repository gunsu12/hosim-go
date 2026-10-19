package fiscalperiod

import (
	"context"
	"testing"
	"time"
)

type mockRepository struct {
	years   map[int]*FiscalYear
	periods map[string]*FiscalPeriod
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		years:   make(map[int]*FiscalYear),
		periods: make(map[string]*FiscalPeriod),
	}
}

func (m *mockRepository) CreateFiscalYearWithPeriods(ctx context.Context, year *FiscalYear, periods []*FiscalPeriod) error {
	m.years[year.Year] = year
	for _, p := range periods {
		m.periods[p.PeriodKey] = p
	}
	return nil
}

func (m *mockRepository) FindYearByID(ctx context.Context, id string) (*FiscalYear, error) {
	for _, y := range m.years {
		if y.ID == id {
			return y, nil
		}
	}
	return nil, ErrYearNotFound
}

func (m *mockRepository) FindYearByYear(ctx context.Context, year int) (*FiscalYear, error) {
	if y, ok := m.years[year]; ok {
		return y, nil
	}
	return nil, ErrYearNotFound
}

func (m *mockRepository) ListYears(ctx context.Context) ([]*FiscalYear, error) {
	res := make([]*FiscalYear, 0, len(m.years))
	for _, y := range m.years {
		res = append(res, y)
	}
	return res, nil
}

func (m *mockRepository) UpdateYear(ctx context.Context, year *FiscalYear) error {
	m.years[year.Year] = year
	return nil
}

func (m *mockRepository) FindPeriodByID(ctx context.Context, id string) (*FiscalPeriod, error) {
	for _, p := range m.periods {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, ErrPeriodNotFound
}

func (m *mockRepository) FindPeriodByKey(ctx context.Context, key string) (*FiscalPeriod, error) {
	if p, ok := m.periods[key]; ok {
		return p, nil
	}
	return nil, ErrPeriodNotFound
}

func (m *mockRepository) FindPeriodByDate(ctx context.Context, date time.Time) (*FiscalPeriod, error) {
	dateOnly := date.Format("2006-01-02")
	for _, p := range m.periods {
		pStart := p.StartDate.Format("2006-01-02")
		pEnd := p.EndDate.Format("2006-01-02")
		if dateOnly >= pStart && dateOnly <= pEnd {
			return p, nil
		}
	}
	return nil, ErrPeriodNotFound
}

func (m *mockRepository) ListPeriodsByYearID(ctx context.Context, yearID string) ([]*FiscalPeriod, error) {
	res := make([]*FiscalPeriod, 0)
	for _, p := range m.periods {
		if p.FiscalYearID == yearID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockRepository) ListPeriods(ctx context.Context, year *int, status *string) ([]*FiscalPeriod, error) {
	res := make([]*FiscalPeriod, 0)
	for _, p := range m.periods {
		if status != nil && string(p.Status) != *status {
			continue
		}
		res = append(res, p)
	}
	return res, nil
}

func (m *mockRepository) UpdatePeriod(ctx context.Context, period *FiscalPeriod) error {
	m.periods[period.PeriodKey] = period
	return nil
}

func TestService_CreateFiscalYear(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	svc := NewService(repo)

	// 1. Sukses buat tahun 2026
	req := CreateFiscalYearRequest{Year: 2026}
	resp, err := svc.CreateFiscalYear(ctx, req, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Year != 2026 {
		t.Errorf("expected year 2026, got %d", resp.Year)
	}
	if len(resp.Periods) != 12 {
		t.Errorf("expected 12 periods, got %d", len(resp.Periods))
	}

	// 2. Duplikasi tahun harus gagal
	_, err = svc.CreateFiscalYear(ctx, req, "admin")
	if err != ErrYearAlreadyExists {
		t.Errorf("expected ErrYearAlreadyExists, got %v", err)
	}
}

func TestService_ValidateDateForPosting_Guard(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	svc := NewService(repo)

	// Buat tahun 2026
	_, _ = svc.CreateFiscalYear(ctx, CreateFiscalYearRequest{Year: 2026}, "admin")

	// 1. Tanggal di periode OPEN (misal 15 Oktober 2026) -> Harus lolos
	octDate := time.Date(2026, time.October, 15, 10, 0, 0, 0, time.UTC)
	info, err := svc.ValidateDateForPosting(ctx, octDate)
	if err != nil {
		t.Fatalf("expected valid date, got error: %v", err)
	}
	if info.PeriodKey != "2026-10" {
		t.Errorf("expected 2026-10, got %s", info.PeriodKey)
	}

	// 2. Kunci periode Oktober (LOCKED) -> Harus mengembalikan ErrPeriodLocked
	octPeriod, _ := repo.FindPeriodByKey(ctx, "2026-10")
	_, err = svc.LockFiscalPeriod(ctx, octPeriod.ID, nil, "manager")
	if err != nil {
		t.Fatalf("lock error: %v", err)
	}

	info, err = svc.ValidateDateForPosting(ctx, octDate)
	if err != ErrPeriodLocked {
		t.Errorf("expected ErrPeriodLocked, got %v", err)
	}
	if info == nil || info.Status != PeriodStatusLocked {
		t.Errorf("expected info status LOCKED")
	}

	// 3. Tutup periode Oktober (CLOSED) -> Harus mengembalikan ErrPeriodClosed
	_, err = svc.CloseFiscalPeriod(ctx, octPeriod.ID, nil, "director")
	if err != nil {
		t.Fatalf("close error: %v", err)
	}

	info, err = svc.ValidateDateForPosting(ctx, octDate)
	if err != ErrPeriodClosed {
		t.Errorf("expected ErrPeriodClosed, got %v", err)
	}

	// 4. Tanggal di tahun 2027 (belum dibuat kalendernya) -> Harus ErrPeriodNotFound
	year2027Date := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, err = svc.ValidateDateForPosting(ctx, year2027Date)
	if err != ErrPeriodNotFound {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}
