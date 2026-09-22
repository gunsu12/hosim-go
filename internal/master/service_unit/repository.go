package serviceunit

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page                 int
	Limit                int
	Search               string
	DepartementID        string
	IsRegistrationTarget *bool
	IsActive             *bool
}

type Repository interface {
	Create(ctx context.Context, su *ServiceUnit) (*ServiceUnit, error)
	Update(ctx context.Context, su *ServiceUnit) (*ServiceUnit, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*ServiceUnit, error)
	FindByCode(ctx context.Context, code string) (*ServiceUnit, error)
	FindAll(ctx context.Context, params ListParams) ([]ServiceUnit, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, su *ServiceUnit) (*ServiceUnit, error) {
	result := r.db.WithContext(ctx).Create(su)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, su.ID)
}

func (r *repository) Update(ctx context.Context, su *ServiceUnit) (*ServiceUnit, error) {
	result := r.db.WithContext(ctx).Save(su)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, su.ID)
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ServiceUnit{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&ServiceUnit{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*ServiceUnit, error) {
	var su ServiceUnit
	err := r.db.WithContext(ctx).Preload("Departement").First(&su, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &su, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*ServiceUnit, error) {
	var su ServiceUnit
	err := r.db.WithContext(ctx).Preload("Departement").First(&su, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &su, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]ServiceUnit, int64, error) {
	var units []ServiceUnit
	var count int64

	query := r.db.WithContext(ctx).Model(&ServiceUnit{})

	if params.DepartementID != "" {
		query = query.Where("departement_id = ?", params.DepartementID)
	}

	if params.IsRegistrationTarget != nil {
		query = query.Where("is_registration_target = ?", *params.IsRegistrationTarget)
	}

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

	err := query.Preload("Departement").Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&units).Error
	if err != nil {
		return nil, 0, err
	}

	return units, count, nil
}
