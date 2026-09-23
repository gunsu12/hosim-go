package customer

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page           int
	Limit          int
	Search         string
	CustomerTypeID string
}

type Repository interface {
	Create(ctx context.Context, customer *Customer) (*Customer, error)
	Update(ctx context.Context, customer *Customer) (*Customer, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindByCode(ctx context.Context, code string) (*Customer, error)
	FindAll(ctx context.Context, params ListParams) ([]Customer, int64, error)

	FindAllCustomerTypes(ctx context.Context) ([]CustomerType, error)
	FindCustomerTypeByID(ctx context.Context, id string) (*CustomerType, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, customer *Customer) (*Customer, error) {
	result := r.db.WithContext(ctx).Create(customer)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, customer.ID)
}

func (r *repository) Update(ctx context.Context, customer *Customer) (*Customer, error) {
	result := r.db.WithContext(ctx).Save(customer)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, customer.ID)
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Customer{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Customer{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Customer, error) {
	var c Customer
	err := r.db.WithContext(ctx).Preload("CustomerType").First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Customer, error) {
	var c Customer
	err := r.db.WithContext(ctx).Preload("CustomerType").First(&c, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Customer, int64, error) {
	var customers []Customer
	var count int64

	query := r.db.WithContext(ctx).Model(&Customer{})

	if params.CustomerTypeID != "" {
		query = query.Where("customer_type_id = ?", params.CustomerTypeID)
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

	err := query.Preload("CustomerType").Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, count, nil
}

func (r *repository) FindAllCustomerTypes(ctx context.Context) ([]CustomerType, error) {
	var types []CustomerType
	err := r.db.WithContext(ctx).Order("name ASC").Find(&types).Error
	return types, err
}

func (r *repository) FindCustomerTypeByID(ctx context.Context, id string) (*CustomerType, error) {
	var ct CustomerType
	err := r.db.WithContext(ctx).First(&ct, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}
