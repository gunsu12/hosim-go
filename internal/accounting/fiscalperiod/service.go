package fiscalperiod

import (
	"context"
	"time"
)

// Service mendefinisikan seluruh operasi use-case kalender akuntansi dan periode fiskal.
type Service interface {
	FiscalPeriodGuard

	CreateFiscalYear(ctx context.Context, req CreateFiscalYearRequest, actor string) (*FiscalYearResponse, error)
	GetFiscalYear(ctx context.Context, id string) (*FiscalYearResponse, error)
	ListFiscalYears(ctx context.Context) ([]*FiscalYearResponse, error)
	CloseFiscalYear(ctx context.Context, id string, actor string) (*FiscalYearResponse, error)

	ListFiscalPeriods(ctx context.Context, year *int, status *string) ([]*FiscalPeriodResponse, error)
	GetFiscalPeriod(ctx context.Context, id string) (*FiscalPeriodResponse, error)
	GetCurrentPeriod(ctx context.Context) (*FiscalPeriodResponse, error)
	LockFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error)
	UnlockFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error)
	CloseFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateFiscalYear(ctx context.Context, req CreateFiscalYearRequest, actor string) (*FiscalYearResponse, error) {
	// 1. Cek apakah tahun fiskal sudah ada
	existing, err := s.repo.FindYearByYear(ctx, req.Year)
	if err == nil && existing != nil {
		return nil, ErrYearAlreadyExists
	}

	// 2. Buat entity tahun beserta 12 periode bulanannya
	fy, periods, err := GenerateFiscalYearWithPeriods(req.Year, req.Notes, actor)
	if err != nil {
		return nil, err
	}

	// 3. Simpan secara atomik
	if err := s.repo.CreateFiscalYearWithPeriods(ctx, fy, periods); err != nil {
		return nil, err
	}

	fy.Periods = periods
	return ToFiscalYearResponse(fy, true), nil
}

func (s *service) GetFiscalYear(ctx context.Context, id string) (*FiscalYearResponse, error) {
	fy, err := s.repo.FindYearByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToFiscalYearResponse(fy, true), nil
}

func (s *service) ListFiscalYears(ctx context.Context) ([]*FiscalYearResponse, error) {
	years, err := s.repo.ListYears(ctx)
	if err != nil {
		return nil, err
	}
	return ToFiscalYearResponseList(years), nil
}

func (s *service) CloseFiscalYear(ctx context.Context, id string, actor string) (*FiscalYearResponse, error) {
	fy, err := s.repo.FindYearByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validasi seluruh 12 periode bulanan harus CLOSED
	if err := fy.CloseYear(fy.Periods, actor); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateYear(ctx, fy); err != nil {
		return nil, err
	}

	return ToFiscalYearResponse(fy, true), nil
}

func (s *service) ListFiscalPeriods(ctx context.Context, year *int, status *string) ([]*FiscalPeriodResponse, error) {
	periods, err := s.repo.ListPeriods(ctx, year, status)
	if err != nil {
		return nil, err
	}
	return ToFiscalPeriodResponseList(periods), nil
}

func (s *service) GetFiscalPeriod(ctx context.Context, id string) (*FiscalPeriodResponse, error) {
	p, err := s.repo.FindPeriodByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToFiscalPeriodResponse(p), nil
}

func (s *service) GetCurrentPeriod(ctx context.Context) (*FiscalPeriodResponse, error) {
	p, err := s.repo.FindPeriodByDate(ctx, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	return ToFiscalPeriodResponse(p), nil
}

func (s *service) LockFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error) {
	p, err := s.repo.FindPeriodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.Lock(actor, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePeriod(ctx, p); err != nil {
		return nil, err
	}

	return ToFiscalPeriodResponse(p), nil
}

func (s *service) UnlockFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error) {
	p, err := s.repo.FindPeriodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.Unlock(actor, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePeriod(ctx, p); err != nil {
		return nil, err
	}

	return ToFiscalPeriodResponse(p), nil
}

func (s *service) CloseFiscalPeriod(ctx context.Context, id string, notes *string, actor string) (*FiscalPeriodResponse, error) {
	p, err := s.repo.FindPeriodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.Close(actor, notes); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePeriod(ctx, p); err != nil {
		return nil, err
	}

	return ToFiscalPeriodResponse(p), nil
}

// ValidateDateForPosting mengimplementasikan interface FiscalPeriodGuard untuk modul jurnal.
func (s *service) ValidateDateForPosting(ctx context.Context, entryDate time.Time) (*ActivePeriodInfo, error) {
	p, err := s.repo.FindPeriodByDate(ctx, entryDate)
	if err != nil {
		return nil, err
	}

	info := &ActivePeriodInfo{
		PeriodID:  p.ID,
		PeriodKey: p.PeriodKey,
		Name:      p.Name,
		Status:    p.Status,
	}

	if p.Status == PeriodStatusLocked {
		return info, ErrPeriodLocked
	}
	if p.Status == PeriodStatusClosed {
		return info, ErrPeriodClosed
	}

	return info, nil
}
