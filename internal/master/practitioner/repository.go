package practitioner

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
	Create(ctx context.Context, practitioner *Practitioner) (*Practitioner, error)
	Update(ctx context.Context, practitioner *Practitioner) (*Practitioner, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Practitioner, error)
	FindByNIK(ctx context.Context, nik string) (*Practitioner, error)
	FindByNIP(ctx context.Context, nip string) (*Practitioner, error)
	FindByIhsPractitionerID(ctx context.Context, ihsPractitionerID string) (*Practitioner, error)
	FindAll(ctx context.Context, params ListParams) ([]Practitioner, int, error)
	FindAllProfessions(ctx context.Context) ([]Profession, error)
	FindAllSpecialties(ctx context.Context) ([]Specialty, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, practitioner *Practitioner) (*Practitioner, error) {
	result := r.db.WithContext(ctx).Create(practitioner)
	if result.Error != nil {
		return nil, result.Error
	}
	return practitioner, nil
}

func (r *repository) Update(ctx context.Context, practitioner *Practitioner) (*Practitioner, error) {
	result := r.db.WithContext(ctx).Save(practitioner)
	if result.Error != nil {
		return nil, result.Error
	}
	return practitioner, nil
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Practitioner{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Practitioner{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Practitioner, error) {
	var practitioner Practitioner
	err := r.db.WithContext(ctx).
		Preload("Profession").
		Preload("Specialty").
		First(&practitioner, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &practitioner, nil
}

func (r *repository) FindByNIK(ctx context.Context, nik string) (*Practitioner, error) {
	var practitioner Practitioner
	err := r.db.WithContext(ctx).
		Preload("Profession").
		Preload("Specialty").
		First(&practitioner, "nik = ?", nik).Error
	if err != nil {
		return nil, err
	}
	return &practitioner, nil
}

func (r *repository) FindByNIP(ctx context.Context, nip string) (*Practitioner, error) {
	var practitioner Practitioner
	err := r.db.WithContext(ctx).
		Preload("Profession").
		Preload("Specialty").
		First(&practitioner, "nip = ?", nip).Error
	if err != nil {
		return nil, err
	}
	return &practitioner, nil
}

func (r *repository) FindByIhsPractitionerID(ctx context.Context, ihsPractitionerID string) (*Practitioner, error) {
	var practitioner Practitioner
	err := r.db.WithContext(ctx).
		Preload("Profession").
		Preload("Specialty").
		First(&practitioner, "ihs_practitioner_id = ?", ihsPractitionerID).Error
	if err != nil {
		return nil, err
	}
	return &practitioner, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Practitioner, int, error) {
	var practitioners []Practitioner
	var count int64

	query := r.db.WithContext(ctx).Model(&Practitioner{})

	if params.Search != "" {
		searchTerm := "%" + strings.TrimSpace(params.Search) + "%"
		query = query.Where(
			"name ILIKE ? OR nip ILIKE ? OR nik ILIKE ? OR phone ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := query.
		Preload("Profession").
		Preload("Specialty").
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Find(&practitioners).Error

	if err != nil {
		return nil, 0, err
	}

	return practitioners, int(count), nil
}

func (r *repository) FindAllProfessions(ctx context.Context) ([]Profession, error) {
	var professions []Profession
	err := r.db.WithContext(ctx).Order("name ASC").Find(&professions).Error
	if err != nil {
		return nil, err
	}
	return professions, nil
}

func (r *repository) FindAllSpecialties(ctx context.Context) ([]Specialty, error) {
	var specialties []Specialty
	err := r.db.WithContext(ctx).Order("name ASC").Find(&specialties).Error
	if err != nil {
		return nil, err
	}
	return specialties, nil
}
