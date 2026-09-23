package tarifclass

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page     int
	Limit    int
	Search   string
	IsActive *bool
}

type Repository interface {
	Create(ctx context.Context, tc *TariffClass) (*TariffClass, error)
	Update(ctx context.Context, tc *TariffClass) (*TariffClass, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*TariffClass, error)
	FindByCode(ctx context.Context, code string) (*TariffClass, error)
	FindAll(ctx context.Context, params ListParams) ([]TariffClass, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tc *TariffClass) (*TariffClass, error) {
	result := r.db.WithContext(ctx).Create(tc)
	if result.Error != nil {
		return nil, result.Error
	}
	return tc, nil
}

func (r *repository) Update(ctx context.Context, tc *TariffClass) (*TariffClass, error) {
	result := r.db.WithContext(ctx).Save(tc)
	if result.Error != nil {
		return nil, result.Error
	}
	return tc, nil
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&TariffClass{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&TariffClass{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*TariffClass, error) {
	var tc TariffClass
	err := r.db.WithContext(ctx).First(&tc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*TariffClass, error) {
	var tc TariffClass
	err := r.db.WithContext(ctx).First(&tc, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]TariffClass, int64, error) {
	var classes []TariffClass
	var count int64

	query := r.db.WithContext(ctx).Model(&TariffClass{})

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&classes).Error
	if err != nil {
		return nil, 0, err
	}

	return classes, count, nil
}
