package payer

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page        int
	Limit       int
	Search      string
	PayerTypeID string
}

type Repository interface {
	Create(ctx context.Context, payer *Payer) (*Payer, error)
	Update(ctx context.Context, payer *Payer) (*Payer, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Payer, error)
	FindByCode(ctx context.Context, code string) (*Payer, error)
	FindAll(ctx context.Context, params ListParams) ([]Payer, int64, error)

	FindAllPayerTypes(ctx context.Context) ([]PayerType, error)
	FindPayerTypeByID(ctx context.Context, id string) (*PayerType, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, payer *Payer) (*Payer, error) {
	result := r.db.WithContext(ctx).Create(payer)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, payer.ID)
}

func (r *repository) Update(ctx context.Context, payer *Payer) (*Payer, error) {
	result := r.db.WithContext(ctx).Save(payer)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, payer.ID)
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Payer{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Payer{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Payer, error) {
	var p Payer
	err := r.db.WithContext(ctx).Preload("PayerType").First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Payer, error) {
	var p Payer
	err := r.db.WithContext(ctx).Preload("PayerType").First(&p, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Payer, int64, error) {
	var payers []Payer
	var count int64

	query := r.db.WithContext(ctx).Model(&Payer{})

	if params.PayerTypeID != "" {
		query = query.Where("payer_type_id = ?", params.PayerTypeID)
	}

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(contact_person) LIKE ?",
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

	err := query.Preload("PayerType").Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&payers).Error
	if err != nil {
		return nil, 0, err
	}

	return payers, count, nil
}

func (r *repository) FindAllPayerTypes(ctx context.Context) ([]PayerType, error) {
	var types []PayerType
	err := r.db.WithContext(ctx).Order("name ASC").Find(&types).Error
	return types, err
}

func (r *repository) FindPayerTypeByID(ctx context.Context, id string) (*PayerType, error) {
	var pt PayerType
	err := r.db.WithContext(ctx).First(&pt, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pt, nil
}
