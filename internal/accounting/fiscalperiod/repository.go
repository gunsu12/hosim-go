package fiscalperiod

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateFiscalYearWithPeriods(ctx context.Context, year *FiscalYear, periods []*FiscalPeriod) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Simpan induk tahun buku
		if err := tx.Create(year).Error; err != nil {
			return err
		}
		// 2. Simpan 12 periode bulanan
		if len(periods) > 0 {
			if err := tx.Create(&periods).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) FindYearByID(ctx context.Context, id string) (*FiscalYear, error) {
	var fy FiscalYear
	err := r.db.WithContext(ctx).
		Preload("Periods", func(db *gorm.DB) *gorm.DB {
			return db.Order("period_number ASC")
		}).
		First(&fy, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrYearNotFound
		}
		return nil, err
	}
	return &fy, nil
}

func (r *repository) FindYearByYear(ctx context.Context, year int) (*FiscalYear, error) {
	var fy FiscalYear
	err := r.db.WithContext(ctx).
		Preload("Periods", func(db *gorm.DB) *gorm.DB {
			return db.Order("period_number ASC")
		}).
		First(&fy, "year = ?", year).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrYearNotFound
		}
		return nil, err
	}
	return &fy, nil
}

func (r *repository) ListYears(ctx context.Context) ([]*FiscalYear, error) {
	var years []*FiscalYear
	err := r.db.WithContext(ctx).
		Preload("Periods", func(db *gorm.DB) *gorm.DB {
			return db.Order("period_number ASC")
		}).
		Order("year DESC").
		Find(&years).Error
	if err != nil {
		return nil, err
	}
	return years, nil
}

func (r *repository) UpdateYear(ctx context.Context, year *FiscalYear) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current FiscalYear
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, "id = ?", year.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrYearNotFound
			}
			return err
		}
		if current.IsClosed && !year.IsClosed {
			return ErrYearAlreadyClosed
		}
		return tx.Model(&FiscalYear{}).Where("id = ?", year.ID).Updates(map[string]any{
			"is_closed":  year.IsClosed,
			"closed_at":  year.ClosedAt,
			"closed_by":  year.ClosedBy,
			"notes":      year.Notes,
			"updated_at": year.UpdatedAt,
			"updated_by": year.UpdatedBy,
		}).Error
	})
}

func (r *repository) FindPeriodByID(ctx context.Context, id string) (*FiscalPeriod, error) {
	var p FiscalPeriod
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPeriodNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindPeriodByKey(ctx context.Context, key string) (*FiscalPeriod, error) {
	var p FiscalPeriod
	err := r.db.WithContext(ctx).First(&p, "period_key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPeriodNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindPeriodByDate(ctx context.Context, date time.Time) (*FiscalPeriod, error) {
	dateStr := date.Format("2006-01-02")
	var p FiscalPeriod
	err := r.db.WithContext(ctx).
		Where("? BETWEEN start_date AND end_date", dateStr).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPeriodNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) ListPeriodsByYearID(ctx context.Context, yearID string) ([]*FiscalPeriod, error) {
	var periods []*FiscalPeriod
	err := r.db.WithContext(ctx).
		Where("fiscal_year_id = ?", yearID).
		Order("period_number ASC").
		Find(&periods).Error
	if err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *repository) ListPeriods(ctx context.Context, year *int, status *string) ([]*FiscalPeriod, error) {
	query := r.db.WithContext(ctx).Model(&FiscalPeriod{})

	if year != nil {
		query = query.Joins("JOIN fiscal_years ON fiscal_years.id = fiscal_periods.fiscal_year_id").
			Where("fiscal_years.year = ?", *year)
	}

	if status != nil && *status != "" {
		query = query.Where("fiscal_periods.status = ?", *status)
	}

	var periods []*FiscalPeriod
	err := query.Order("fiscal_periods.period_key ASC").Find(&periods).Error
	if err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *repository) UpdatePeriod(ctx context.Context, period *FiscalPeriod) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current FiscalPeriod
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, "id = ?", period.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPeriodNotFound
			}
			return err
		}
		if current.Status == PeriodStatusClosed && period.Status != PeriodStatusClosed {
			return ErrCannotLockClosedPeriod
		}
		return tx.Model(&FiscalPeriod{}).Where("id = ?", period.ID).Updates(map[string]any{
			"status":     period.Status,
			"locked_at":  period.LockedAt,
			"locked_by":  period.LockedBy,
			"closed_at":  period.ClosedAt,
			"closed_by":  period.ClosedBy,
			"notes":      period.Notes,
			"updated_at": period.UpdatedAt,
			"updated_by": period.UpdatedBy,
		}).Error
	})
}
