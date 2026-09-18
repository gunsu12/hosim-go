package patient

import (
	"context"

	"gorm.io/gorm"
)

type ListParams struct {
	Page   int
	Limit  int
	Search string
}

type Repository interface {
	Create(ctx context.Context, patient *Patient) (*Patient, error)
	Update(ctx context.Context, patient *Patient) (*Patient, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Patient, error)
	FindByNIK(ctx context.Context, nik string) (*Patient, error)
	FindByMedicalRecordNo(ctx context.Context, medicalRecordNo string) (*Patient, error)
	FindAll(ctx context.Context, params ListParams) ([]Patient, int, error)
	GetLastMedicalRecordNo(ctx context.Context) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// create data patient
func (r *repository) Create(ctx context.Context, patient *Patient) (*Patient, error) {
	result := r.db.WithContext(ctx).Create(patient)
	if result.Error != nil {
		return nil, result.Error
	}
	return patient, nil
}

// update data patient
func (r *repository) Update(ctx context.Context, patient *Patient) (*Patient, error) {
	result := r.db.WithContext(ctx).Save(patient)
	if result.Error != nil {
		return nil, result.Error
	}
	return patient, nil
}

// delete data patient
func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Patient{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Patient{}, "id = ?", id).Error
	})
}

// find by id
func (r *repository) FindByID(ctx context.Context, id string) (*Patient, error) {
	var patient Patient
	err := r.db.WithContext(ctx).
		Preload("EmergencyContacts", "is_active = ?", true).
		Preload("Relations", "is_active = ?", true).
		Preload("Allergies", "is_active = ?", true).
		Preload("Addresses", "is_active = ?", true).
		Preload("DrugHistories", "is_active = ?", true).
		Preload("ChronicalDiseases", "is_active = ?", true).
		First(&patient, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// find by nik
func (r *repository) FindByNIK(ctx context.Context, nik string) (*Patient, error) {
	var patient Patient
	err := r.db.WithContext(ctx).
		Preload("EmergencyContacts", "is_active = ?", true).
		Preload("Relations", "is_active = ?", true).
		Preload("Allergies", "is_active = ?", true).
		Preload("Addresses", "is_active = ?", true).
		Preload("DrugHistories", "is_active = ?", true).
		Preload("ChronicalDiseases", "is_active = ?", true).
		First(&patient, "nik = ?", nik).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// find by medical record no
func (r *repository) FindByMedicalRecordNo(ctx context.Context, medicalRecordNo string) (*Patient, error) {
	var patient Patient
	err := r.db.WithContext(ctx).
		Preload("EmergencyContacts", "is_active = ?", true).
		Preload("Relations", "is_active = ?", true).
		Preload("Allergies", "is_active = ?", true).
		Preload("Addresses", "is_active = ?", true).
		Preload("DrugHistories", "is_active = ?", true).
		Preload("ChronicalDiseases", "is_active = ?", true).
		First(&patient, "medical_record_no = ?", medicalRecordNo).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// find all yang bersih dan berkinerja tinggi
func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Patient, int, error) {
	var patients []Patient
	var count int64

	// 1. Buat base query
	query := r.db.WithContext(ctx).Model(&Patient{})

	// 2. Pasang filter HANYA JIKA user mengetik pencarian
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where(
			"medical_record_no LIKE ? OR full_name LIKE ? OR nik LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// 3. Hitung total data yang cocok (untuk keperluan info pagination)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 4. Ambil datanya sesuai limit dan offset halaman
	offset := (params.Page - 1) * params.Limit
	err := query.
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Find(&patients).Error

	if err != nil {
		return nil, 0, err
	}

	return patients, int(count), nil
}

// get last medical record no
func (r *repository) GetLastMedicalRecordNo(ctx context.Context) (string, error) {
	var lastRM string
	err := r.db.WithContext(ctx).
		Model(&Patient{}).
		Unscoped().
		Select("medical_record_no").
		Order("created_at DESC").
		Limit(1).
		Scan(&lastRM).Error

	if err != nil {
		return "", err
	}

	return lastRM, nil
}
