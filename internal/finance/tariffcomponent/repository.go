package tariffcomponent

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page          int
	Limit         int
	Search        string
	IsActive      *bool
	ComponentType string
}

type Repository interface {
	Create(ctx context.Context, tc *TariffComponent) (*TariffComponent, error)
	Update(ctx context.Context, tc *TariffComponent) (*TariffComponent, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*TariffComponent, error)
	FindByCode(ctx context.Context, code string) (*TariffComponent, error)
	FindAll(ctx context.Context, params ListParams) ([]TariffComponent, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, tc *TariffComponent) (*TariffComponent, error) {
	result := r.db.WithContext(ctx).Create(tc)
	if result.Error != nil {
		return nil, result.Error
	}
	return tc, nil
}

func (r *repository) Update(ctx context.Context, tc *TariffComponent) (*TariffComponent, error) {
	result := r.db.WithContext(ctx).Save(tc)
	if result.Error != nil {
		return nil, result.Error
	}
	return tc, nil
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&TariffComponent{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&TariffComponent{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*TariffComponent, error) {
	var tc TariffComponent
	err := r.db.WithContext(ctx).First(&tc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*TariffComponent, error) {
	var tc TariffComponent
	err := r.db.WithContext(ctx).First(&tc, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]TariffComponent, int64, error) {
	components := make([]TariffComponent, 0)
	var count int64

	query := r.db.WithContext(ctx).Model(&TariffComponent{})

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.ComponentType != "" && params.ComponentType != "ALL" {
		query = query.Where("component_type = ?", params.ComponentType)
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

	err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&components).Error
	if err != nil {
		return nil, 0, err
	}

	return components, count, nil
}
