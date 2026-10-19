package fiscalperiod

import (
	"time"
)

// CreateFiscalYearRequest adalah payload untuk membuat tahun buku baru beserta 12 bulannya.
type CreateFiscalYearRequest struct {
	Year  int     `json:"year" binding:"required,min=2000,max=2100"`
	Notes *string `json:"notes"`
}

// FiscalPeriodActionRequest adalah payload opsional catatan saat melakukan aksi lock/unlock/close.
type FiscalPeriodActionRequest struct {
	Notes *string `json:"notes"`
}

// FiscalPeriodResponse adalah representasi DTO periode bulanan untuk client.
type FiscalPeriodResponse struct {
	ID           string       `json:"id"`
	FiscalYearID string       `json:"fiscal_year_id"`
	PeriodNumber int          `json:"period_number"`
	PeriodKey    string       `json:"period_key"`
	Name         string       `json:"name"`
	StartDate    string       `json:"start_date"`
	EndDate      string       `json:"end_date"`
	Status       PeriodStatus `json:"status"`
	LockedAt     *time.Time   `json:"locked_at,omitempty"`
	LockedBy     *string      `json:"locked_by,omitempty"`
	ClosedAt     *time.Time   `json:"closed_at,omitempty"`
	ClosedBy     *string      `json:"closed_by,omitempty"`
	Notes        *string      `json:"notes,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func ToFiscalPeriodResponse(p *FiscalPeriod) *FiscalPeriodResponse {
	if p == nil {
		return nil
	}
	return &FiscalPeriodResponse{
		ID:           p.ID,
		FiscalYearID: p.FiscalYearID,
		PeriodNumber: p.PeriodNumber,
		PeriodKey:    p.PeriodKey,
		Name:         p.Name,
		StartDate:    p.StartDate.Format("2006-01-02"),
		EndDate:      p.EndDate.Format("2006-01-02"),
		Status:       p.Status,
		LockedAt:     p.LockedAt,
		LockedBy:     p.LockedBy,
		ClosedAt:     p.ClosedAt,
		ClosedBy:     p.ClosedBy,
		Notes:        p.Notes,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func ToFiscalPeriodResponseList(periods []*FiscalPeriod) []*FiscalPeriodResponse {
	res := make([]*FiscalPeriodResponse, len(periods))
	for i, p := range periods {
		res[i] = ToFiscalPeriodResponse(p)
	}
	return res
}

// FiscalYearResponse adalah representasi DTO tahun fiskal untuk client.
type FiscalYearResponse struct {
	ID           string                  `json:"id"`
	Year         int                     `json:"year"`
	StartDate    string                  `json:"start_date"`
	EndDate      string                  `json:"end_date"`
	IsClosed     bool                    `json:"is_closed"`
	ClosedAt     *time.Time              `json:"closed_at,omitempty"`
	ClosedBy     *string                 `json:"closed_by,omitempty"`
	Notes        *string                 `json:"notes,omitempty"`
	PeriodsCount int                     `json:"periods_count"`
	Periods      []*FiscalPeriodResponse `json:"periods,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

func ToFiscalYearResponse(y *FiscalYear, includePeriods bool) *FiscalYearResponse {
	if y == nil {
		return nil
	}
	res := &FiscalYearResponse{
		ID:           y.ID,
		Year:         y.Year,
		StartDate:    y.StartDate.Format("2006-01-02"),
		EndDate:      y.EndDate.Format("2006-01-02"),
		IsClosed:     y.IsClosed,
		ClosedAt:     y.ClosedAt,
		ClosedBy:     y.ClosedBy,
		Notes:        y.Notes,
		PeriodsCount: len(y.Periods),
		CreatedAt:    y.CreatedAt,
		UpdatedAt:    y.UpdatedAt,
	}
	if includePeriods && len(y.Periods) > 0 {
		res.Periods = ToFiscalPeriodResponseList(y.Periods)
	}
	return res
}

func ToFiscalYearResponseList(years []*FiscalYear) []*FiscalYearResponse {
	res := make([]*FiscalYearResponse, len(years))
	for i, y := range years {
		res[i] = ToFiscalYearResponse(y, false)
	}
	return res
}
