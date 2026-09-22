package referal

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page   int
	Limit  int
	Search string
	Type   string
}

type Repository interface {
	Create(ctx context.Context, ref *Referal) (*Referal, error)
	Update(ctx context.Context, ref *Referal) (*Referal, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Referal, error)
	FindByCode(ctx context.Context, code string) (*Referal, error)
	FindAll(ctx context.Context, params ListParams) ([]Referal, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, ref *Referal) (*Referal, error) {
	result := r.db.WithContext(ctx).Create(ref)
	if result.Error != nil {
		return nil, result.Error
	}
	return ref, nil
}

func (r *repository) Update(ctx context.Context, ref *Referal) (*Referal, error) {
	result := r.db.WithContext(ctx).Save(ref)
	if result.Error != nil {
		return nil, result.Error
	}
	return ref, nil
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Referal{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Referal{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Referal, error) {
	var ref Referal
	err := r.db.WithContext(ctx).First(&ref, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Referal, error) {
	var ref Referal
	err := r.db.WithContext(ctx).First(&ref, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Referal, int64, error) {
	var refs []Referal
	var count int64

	query := r.db.WithContext(ctx).Model(&Referal{})

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(phone) LIKE ?",
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

	err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&refs).Error
	if err != nil {
		return nil, 0, err
	}

	return refs, count, nil
}
