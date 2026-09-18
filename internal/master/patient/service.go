package patient

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type EmergencyContactRequest struct {
	Name     string `json:"name" binding:"required"`
	Relation string `json:"relation" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address"`
}

type AddressRequest struct {
	AddressType string `json:"address_type" binding:"required"` // "ktp", "domicile", dll
	AddressLine string `json:"address_line" binding:"required"`
	RT          string `json:"rt"`
	RW          string `json:"rw"`
	PostalCode  string `json:"postal_code"`
	ProvinsiID  string `json:"provinsi_id"`
	KabupatenID string `json:"kabupaten_id"`
	KecamatanID string `json:"kecamatan_id"`
	KelurahanID string `json:"kelurahan_id"`
	IsActive    bool   `json:"is_active"`
}

type CreatePatientRequest struct {
	MedicalRecordNo     string                    `json:"medical_record_no"` // Opsional: jika kosong akan di-generate otomatis
	NIK                 string                    `json:"nik"`               // Opsional: boleh kosong untuk Bayi Baru Lahir (BBL) / Mr. X
	IHSPatientID        string                    `json:"ihs_patient_id"`    // Opsional: diisi saat sync SatuSehat
	ShortName           string                    `json:"short_name" binding:"required"`
	FullName            string                    `json:"full_name" binding:"required"`
	Gender              string                    `json:"gender" binding:"required"` // "L" / "P" atau "male" / "female"
	BirthPlace          string                    `json:"birth_place" binding:"required"`
	BirthDate           string                    `json:"birth_date" binding:"required"` // Format: YYYY-MM-DD
	Phone               string                    `json:"phone" binding:"required"`
	Email               string                    `json:"email"`
	MaritalStatus       string                    `json:"marital_status"`
	Religion            string                    `json:"religion"`
	Education           string                    `json:"education"`
	Occupation          string                    `json:"occupation"`
	Nationality         string                    `json:"nationality"`
	BloodType           string                    `json:"blood_type"`
	InsuranceType       string                    `json:"insurance_type"`
	InsuranceNumber     string                    `json:"insurance_number"`
	InsuranceExpiryDate *string                   `json:"insurance_expiry_date"` // Format: YYYY-MM-DD (opsional)
	EmergencyContacts   []EmergencyContactRequest `json:"emergency_contacts"`
	Addresses           []AddressRequest          `json:"addresses"`
}

type UpdatePatientRequest struct {
	CreatePatientRequest
}

type Service interface {
	RegisterPatient(ctx context.Context, req CreatePatientRequest, operatorID string) (*Patient, error)
	UpdatePatient(ctx context.Context, id string, req UpdatePatientRequest, operatorID string) (*Patient, error)
	DeletePatient(ctx context.Context, id string, operatorID string) error
	GetPatientByID(ctx context.Context, id string) (*Patient, error)
	GetPatientByNIK(ctx context.Context, nik string) (*Patient, error)
	GetPatientByMedicalRecordNo(ctx context.Context, rmNo string) (*Patient, error)
	ListPatients(ctx context.Context, params ListParams) ([]Patient, int, error)
}

type service struct {
	repo Repository
}

// NewService membuat instance baru dari service pasien
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// RegisterPatient mendaftarkan pasien baru di loket
func (s *service) RegisterPatient(ctx context.Context, req CreatePatientRequest, operatorID string) (*Patient, error) {
	// A. Validasi NIK: jika diisi, pastikan belum terdaftar
	if req.NIK != "" {
		existing, err := s.repo.FindByNIK(ctx, req.NIK)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek NIK: %w", err)
		}
		if existing != nil {
			return nil, errors.New("pasien dengan NIK tersebut sudah terdaftar")
		}
	}

	// B. Parse Tanggal Lahir (string -> time.Time)
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("format tanggal lahir tidak valid (gunakan format YYYY-MM-DD)")
	}

	// C. Penomoran Otomatis No. Rekam Medis jika tidak diisi manual
	medicalRecordNo := req.MedicalRecordNo
	if medicalRecordNo == "" {
		generatedNo, err := s.generateNextMedicalRecordNo(ctx)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat nomor rekam medis otomatis: %w", err)
		}
		medicalRecordNo = generatedNo
	} else {
		// Jika diisi manual (misal migrasi), cek apakah sudah ada yang pakai
		existingRM, err := s.repo.FindByMedicalRecordNo(ctx, medicalRecordNo)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek nomor rekam medis: %w", err)
		}
		if existingRM != nil {
			return nil, errors.New("nomor rekam medis tersebut sudah digunakan oleh pasien lain")
		}
	}

	// D. Mapping Kontak Darurat
	var emergencyContacts []PatientEmergencyContact
	for _, c := range req.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, PatientEmergencyContact{
			Name:      c.Name,
			Relation:  c.Relation,
			Phone:     c.Phone,
			Address:   c.Address,
			IsActive:  true,
			CreatedBy: operatorID,
		})
	}

	// E. Mapping Alamat
	var addresses []PatientAddress
	for _, a := range req.Addresses {
		addresses = append(addresses, PatientAddress{
			AddressType: a.AddressType,
			AddressLine: a.AddressLine,
			RT:          a.RT,
			RW:          a.RW,
			PostalCode:  a.PostalCode,
			ProvinsiID:  a.ProvinsiID,
			KabupatenID: a.KabupatenID,
			KecamatanID: a.KecamatanID,
			KelurahanID: a.KelurahanID,
			IsActive:    true,
			CreatedBy:   operatorID,
		})
	}

	// F. Optional: Parse Masa Berlaku Asuransi jika diisi
	var insuranceExpiry *time.Time
	if req.InsuranceExpiryDate != nil && *req.InsuranceExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.InsuranceExpiryDate)
		if err == nil {
			insuranceExpiry = &t
		}
	}

	// G. Buat Objek Entity Pasien
	patient := &Patient{
		MedicalRecordNo:     medicalRecordNo,
		NIK:                 req.NIK,
		IHSPatientID:        req.IHSPatientID,
		ShortName:           req.ShortName,
		FullName:            req.FullName,
		Gender:              req.Gender,
		BirthPlace:          req.BirthPlace,
		BirthDate:           birthDate,
		Phone:               req.Phone,
		Email:               req.Email,
		MaritalStatus:       req.MaritalStatus,
		Religion:            req.Religion,
		Education:           req.Education,
		Occupation:          req.Occupation,
		Nationality:         req.Nationality,
		BloodType:           req.BloodType,
		InsuranceType:       req.InsuranceType,
		InsuranceNumber:     req.InsuranceNumber,
		InsuranceExpiryDate: insuranceExpiry,
		EmergencyContacts:   emergencyContacts,
		Addresses:           addresses,
		CreatedBy:           operatorID,
	}

	// H. Simpan ke database via repository
	return s.repo.Create(ctx, patient)
}

// UpdatePatient memperbarui data demografi pasien di loket
func (s *service) UpdatePatient(ctx context.Context, id string, req UpdatePatientRequest, operatorID string) (*Patient, error) {
	// 1. Pastikan pasien yang ingin diedit ada di database
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("pasien tidak ditemukan")
	}

	// 2. Parse Tanggal Lahir
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("format tanggal lahir tidak valid (gunakan format YYYY-MM-DD)")
	}

	// 3. Update field data pokok
	patient.ShortName = req.ShortName
	patient.FullName = req.FullName
	patient.Gender = req.Gender
	patient.BirthPlace = req.BirthPlace
	patient.BirthDate = birthDate
	patient.Phone = req.Phone
	patient.Email = req.Email
	patient.MaritalStatus = req.MaritalStatus
	patient.Religion = req.Religion
	patient.Education = req.Education
	patient.Occupation = req.Occupation
	patient.Nationality = req.Nationality
	patient.BloodType = req.BloodType
	patient.InsuranceType = req.InsuranceType
	patient.InsuranceNumber = req.InsuranceNumber
	patient.UpdatedBy = operatorID

	// 4. Update Kontak Darurat jika dikirimkan
	if req.EmergencyContacts != nil {
		var newContacts []PatientEmergencyContact
		for _, c := range req.EmergencyContacts {
			newContacts = append(newContacts, PatientEmergencyContact{
				PatientID: id,
				Name:      c.Name,
				Relation:  c.Relation,
				Phone:     c.Phone,
				Address:   c.Address,
				IsActive:  true,
				CreatedBy: operatorID,
			})
		}
		patient.EmergencyContacts = newContacts
	}

	// 5. Update Alamat jika dikirimkan
	if req.Addresses != nil {
		var newAddresses []PatientAddress
		for _, a := range req.Addresses {
			newAddresses = append(newAddresses, PatientAddress{
				PatientID:   id,
				AddressType: a.AddressType,
				AddressLine: a.AddressLine,
				RT:          a.RT,
				RW:          a.RW,
				PostalCode:  a.PostalCode,
				ProvinsiID:  a.ProvinsiID,
				KabupatenID: a.KabupatenID,
				KecamatanID: a.KecamatanID,
				KelurahanID: a.KelurahanID,
				IsActive:    true,
				CreatedBy:   operatorID,
			})
		}
		patient.Addresses = newAddresses
	}

	// 6. Simpan perubahan ke repository
	return s.repo.Update(ctx, patient)
}

// DeletePatient menghapus pasien (soft delete)
func (s *service) DeletePatient(ctx context.Context, id string, operatorID string) error {
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if patient == nil {
		return errors.New("pasien tidak ditemukan")
	}

	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetPatientByID(ctx context.Context, id string) (*Patient, error) {
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("pasien tidak ditemukan")
	}
	return patient, nil
}

func (s *service) GetPatientByNIK(ctx context.Context, nik string) (*Patient, error) {
	patient, err := s.repo.FindByNIK(ctx, nik)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("pasien tidak ditemukan")
	}
	return patient, nil
}

func (s *service) GetPatientByMedicalRecordNo(ctx context.Context, rmNo string) (*Patient, error) {
	patient, err := s.repo.FindByMedicalRecordNo(ctx, rmNo)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("pasien tidak ditemukan")
	}
	return patient, nil
}

func (s *service) ListPatients(ctx context.Context, params ListParams) ([]Patient, int, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	return s.repo.FindAll(ctx, params)
}

// generateNextMedicalRecordNo menghasilkan format 8 digit angka berurutan: 00000001, 00000002, dst.
func (s *service) generateNextMedicalRecordNo(ctx context.Context) (string, error) {
	lastRM, err := s.repo.GetLastMedicalRecordNo(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Jika belum ada data sama sekali di database, mulai dari nomor 00000001
	if lastRM == "" {
		return "00000001", nil
	}

	// Konversi string ke angka
	lastNum, err := strconv.Atoi(lastRM)
	if err != nil {
		return "", fmt.Errorf("format no rekam medis terakhir (%s) bukan angka: %w", lastRM, err)
	}

	nextNum := lastNum + 1
	// Format kembali menjadi 8 digit dengan awalan 0 (leading zeros)
	return fmt.Sprintf("%08d", nextNum), nil
}
