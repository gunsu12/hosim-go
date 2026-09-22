package departement

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page   int
	Limit  int
	Search string
}

type Repository interface {
	Create(ctx context.Context, dept *Departement) (*Departement, error)
	Update(ctx context.Context, dept *Departement) (*Departement, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Departement, error)
	FindByCode(ctx context.Context, code string) (*Departement, error)
	FindAll(ctx context.Context, params ListParams) ([]Departement, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dept *Departement) (*Departement, error) {
	result := r.db.WithContext(ctx).Create(dept)
	if result.Error != nil {
		return nil, result.Error
	}
	return dept, nil
}

func (r *repository) Update(ctx context.Context, dept *Departement) (*Departement, error) {
	result := r.db.WithContext(ctx).Save(dept)
	if result.Error != nil {
		return nil, result.Error
	}
	return dept, nil
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Departement{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Departement{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Departement, error) {
	var dept Departement
	err := r.db.WithContext(ctx).First(&dept, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Departement, error) {
	var dept Departement
	err := r.db.WithContext(ctx).First(&dept, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Departement, int64, error) {
	var depts []Departement
	var count int64

	query := r.db.WithContext(ctx).Model(&Departement{})

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(departement_type) LIKE ?",
			searchTerm, searchTerm, searchTerm,
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

	err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&depts).Error
	if err != nil {
		return nil, 0, err
	}

	return depts, count, nil
}
