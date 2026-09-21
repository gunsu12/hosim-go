package patient

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Simpan data pokok pasien tanpa menimpa relasi otomatis
		if err := tx.Omit(clause.Associations).Save(patient).Error; err != nil {
			return err
		}

		// Jika kontak darurat diubah, hapus (soft delete) data lama lalu simpan data baru
		if patient.EmergencyContacts != nil {
			if err := tx.Where("patient_id = ?", patient.ID).Delete(&PatientEmergencyContact{}).Error; err != nil {
				return err
			}
			for i := range patient.EmergencyContacts {
				patient.EmergencyContacts[i].ID = ""
				patient.EmergencyContacts[i].PatientID = patient.ID
			}
			if len(patient.EmergencyContacts) > 0 {
				if err := tx.Create(&patient.EmergencyContacts).Error; err != nil {
					return err
				}
			}
		}

		// Jika relasi keluarga diubah, hapus (soft delete) data lama lalu simpan data baru
		if patient.Relations != nil {
			if err := tx.Where("patient_id = ?", patient.ID).Delete(&PatientRelation{}).Error; err != nil {
				return err
			}
			for i := range patient.Relations {
				patient.Relations[i].ID = ""
				patient.Relations[i].PatientID = patient.ID
			}
			if len(patient.Relations) > 0 {
				if err := tx.Create(&patient.Relations).Error; err != nil {
					return err
				}
			}
		}

		// Jika alamat diubah, hapus (soft delete) data lama lalu simpan data baru
		if patient.Addresses != nil {
			if err := tx.Where("patient_id = ?", patient.ID).Delete(&PatientAddress{}).Error; err != nil {
				return err
			}
			for i := range patient.Addresses {
				patient.Addresses[i].ID = ""
				patient.Addresses[i].PatientID = patient.ID
			}
			if len(patient.Addresses) > 0 {
				if err := tx.Create(&patient.Addresses).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return patient, nil
}

// delete data patient
func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update deleted_by dan soft delete pada pasien utama
		if err := tx.Model(&Patient{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Patient{}, "id = ?", id).Error; err != nil {
			return err
		}

		// Cascade soft delete ke tabel-tabel anak & catat deleted_by
		if err := tx.Model(&PatientEmergencyContact{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientEmergencyContact{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&PatientAddress{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientAddress{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&PatientRelation{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientRelation{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&PatientAllergy{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientAllergy{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&PatientDrugHistory{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientDrugHistory{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&PatientChronicalDisease{}).Where("patient_id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("patient_id = ?", id).Delete(&PatientChronicalDisease{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// find by id
func (r *repository) FindByID(ctx context.Context, id string) (*Patient, error) {
	var patient Patient
	err := r.db.WithContext(ctx).
		Preload("Payer").
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
		Preload("Payer").
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
		Preload("Payer").
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

	// 2. Pasang filter HANYA JIKA user mengetik pencarian (PostgreSQL ILIKE)
	if params.Search != "" {
		searchTerm := "%" + strings.TrimSpace(params.Search) + "%"
		query = query.Where(
			"medical_record_no ILIKE ? OR full_name ILIKE ? OR nik ILIKE ? OR family_card_no ILIKE ? OR phone ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 3. Hitung total data yang cocok (untuk keperluan info pagination)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 4. Ambil datanya sesuai limit dan offset halaman
	offset := (params.Page - 1) * params.Limit
	err := query.
		Preload("Payer").
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Find(&patients).Error

	if err != nil {
		return nil, 0, err
	}

	return patients, int(count), nil
}

// get last medical record no (hanya mengambil nomor rekam medis 8 digit numerik)
func (r *repository) GetLastMedicalRecordNo(ctx context.Context) (string, error) {
	var candidates []string
	err := r.db.WithContext(ctx).
		Model(&Patient{}).
		Unscoped().
		Where("length(medical_record_no) = ?", 8).
		Order("medical_record_no DESC").
		Limit(50).
		Pluck("medical_record_no", &candidates).Error

	if err != nil {
		return "", err
	}

	for _, rm := range candidates {
		if len(rm) == 8 && isDigits(rm) {
			return rm, nil
		}
	}

	return "", nil
}
