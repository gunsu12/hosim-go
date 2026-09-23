package patient

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

var (
	ErrPatientNotFound                  = errors.New("pasien tidak ditemukan")
	ErrNIKAlreadyExists                 = errors.New("pasien dengan NIK tersebut sudah terdaftar")
	ErrMedicalRecordNoExists            = errors.New("nomor rekam medis tersebut sudah digunakan oleh pasien lain")
	ErrInvalidDateFormat                = errors.New("format tanggal lahir tidak valid (gunakan format YYYY-MM-DD)")
	ErrBirthDateInFuture                = errors.New("tanggal lahir tidak boleh di masa depan")
	ErrInvalidNIKFormat                 = errors.New("NIK harus berjumlah 16 digit angka")
	ErrInvalidFamilyCardFormat          = errors.New("nomor Kartu Keluarga (No. KK) harus berjumlah 16 digit angka")
	ErrInvalidGender                    = errors.New("jenis kelamin tidak valid (gunakan 'L' / 'P' atau 'male' / 'female')")
	ErrInvalidEmail                     = errors.New("format email tidak valid")
	ErrInvalidDeceasedDate              = errors.New("tanggal meninggal tidak valid (gunakan format YYYY-MM-DD atau YYYY-MM-DD HH:mm:ss, dan tidak boleh sebelum lahir atau di masa depan)")
	ErrMaxMedicalRecordExceeded         = errors.New("nomor rekam medis telah mencapai batas maksimal (99999999)")
	ErrInvalidMedicalRecordFormat       = errors.New("nomor rekam medis manual harus berupa 8 digit angka")
	ErrInvalidInsuranceExpiryDateFormat = errors.New("format masa berlaku asuransi tidak valid (gunakan format YYYY-MM-DD)")
)

type EmergencyContactRequest struct {
	Name     string `json:"name" binding:"required"`
	Relation string `json:"relation" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address"`
}

type RelationRequest struct {
	Name     string `json:"name" binding:"required"`
	Relation string `json:"relation" binding:"required"`
	Phone    string `json:"phone"`
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
	FamilyCardNo        string                    `json:"family_card_no"`    // Opsional: No Kartu Keluarga
	IHSPatientID        string                    `json:"ihs_patient_id"`    // Opsional: diisi saat sync SatuSehat
	Title               string                    `json:"title"`             // Opsional: Tn, Ny, Nn, An, By
	ShortName           string                    `json:"short_name" binding:"required"`
	FullName            string                    `json:"full_name" binding:"required"`
	MotherName          string                    `json:"mother_name"`               // Opsional: Nama Ibu Kandung
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
	Rhesus              string                    `json:"rhesus"`                // "+", "-", atau "tidak tahu"
	SpecialNeeds        string                    `json:"special_needs"`         // Disabilitas / Kebutuhan khusus (kursi roda, tuna rungu, dll)
	IsUnknown           bool                      `json:"is_unknown"`            // True jika Mr. X / Mrs. X (IGD)
	IsDeceased          bool                      `json:"is_deceased"`           // True jika pasien meninggal
	DeceasedAt          *string                   `json:"deceased_at"`           // Format: YYYY-MM-DD HH:mm:ss atau YYYY-MM-DD (opsional)
	CustomerID          *string                   `json:"customer_id"`           // ID Customer / Penjamin
	PayerID             *string                   `json:"payer_id"`              // Alias untuk customer_id (kompatibilitas)
	InsuranceType       *string                   `json:"insurance_type"`        // Jenis Asuransi
	InsuranceNumber     *string                   `json:"insurance_number"`      // Nomor Asuransi
	InsuranceExpiryDate *string                   `json:"insurance_expiry_date"` // Format: YYYY-MM-DD (opsional)
	EmergencyContacts   []EmergencyContactRequest `json:"emergency_contacts"`
	Relations           []RelationRequest         `json:"relations"`
	Addresses           []AddressRequest          `json:"addresses"`
}

// Sanitize membersihkan spasi di awal dan akhir seluruh field string pada request
func (req *CreatePatientRequest) Sanitize() {
	req.NIK = strings.TrimSpace(req.NIK)
	req.FamilyCardNo = strings.TrimSpace(req.FamilyCardNo)
	req.IHSPatientID = strings.TrimSpace(req.IHSPatientID)
	req.Title = strings.TrimSpace(req.Title)
	req.ShortName = strings.TrimSpace(req.ShortName)
	req.FullName = strings.TrimSpace(req.FullName)
	req.MotherName = strings.TrimSpace(req.MotherName)
	req.Gender = strings.TrimSpace(req.Gender)
	req.BirthPlace = strings.TrimSpace(req.BirthPlace)
	req.BirthDate = strings.TrimSpace(req.BirthDate)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.MaritalStatus = strings.TrimSpace(req.MaritalStatus)
	req.Religion = strings.TrimSpace(req.Religion)
	req.Education = strings.TrimSpace(req.Education)
	req.Occupation = strings.TrimSpace(req.Occupation)
	req.Nationality = strings.TrimSpace(req.Nationality)
	req.BloodType = strings.TrimSpace(req.BloodType)
	req.Rhesus = strings.TrimSpace(req.Rhesus)
	req.SpecialNeeds = strings.TrimSpace(req.SpecialNeeds)

	if req.CustomerID == nil && req.PayerID != nil {
		req.CustomerID = req.PayerID
	}
	if req.CustomerID != nil {
		trimmed := strings.TrimSpace(*req.CustomerID)
		if trimmed == "" {
			req.CustomerID = nil
			req.PayerID = nil
		} else {
			req.CustomerID = &trimmed
			req.PayerID = &trimmed
		}
	}

	if req.InsuranceType != nil {
		trimmed := strings.TrimSpace(*req.InsuranceType)
		if trimmed == "" {
			req.InsuranceType = nil
		} else {
			req.InsuranceType = &trimmed
		}
	}

	if req.InsuranceNumber != nil {
		trimmed := strings.TrimSpace(*req.InsuranceNumber)
		if trimmed == "" {
			req.InsuranceNumber = nil
		} else {
			req.InsuranceNumber = &trimmed
		}
	}

	req.MedicalRecordNo = strings.TrimSpace(req.MedicalRecordNo)

	if req.DeceasedAt != nil {
		trimmed := strings.TrimSpace(*req.DeceasedAt)
		req.DeceasedAt = &trimmed
	}
	if req.InsuranceExpiryDate != nil {
		trimmed := strings.TrimSpace(*req.InsuranceExpiryDate)
		req.InsuranceExpiryDate = &trimmed
	}

	for i := range req.EmergencyContacts {
		req.EmergencyContacts[i].Name = strings.TrimSpace(req.EmergencyContacts[i].Name)
		req.EmergencyContacts[i].Relation = strings.TrimSpace(req.EmergencyContacts[i].Relation)
		req.EmergencyContacts[i].Phone = strings.TrimSpace(req.EmergencyContacts[i].Phone)
		req.EmergencyContacts[i].Address = strings.TrimSpace(req.EmergencyContacts[i].Address)
	}

	for i := range req.Relations {
		req.Relations[i].Name = strings.TrimSpace(req.Relations[i].Name)
		req.Relations[i].Relation = strings.TrimSpace(req.Relations[i].Relation)
		req.Relations[i].Phone = strings.TrimSpace(req.Relations[i].Phone)
		req.Relations[i].Address = strings.TrimSpace(req.Relations[i].Address)
	}

	for i := range req.Addresses {
		req.Addresses[i].AddressType = strings.TrimSpace(req.Addresses[i].AddressType)
		req.Addresses[i].AddressLine = strings.TrimSpace(req.Addresses[i].AddressLine)
		req.Addresses[i].RT = strings.TrimSpace(req.Addresses[i].RT)
		req.Addresses[i].RW = strings.TrimSpace(req.Addresses[i].RW)
		req.Addresses[i].PostalCode = strings.TrimSpace(req.Addresses[i].PostalCode)
		req.Addresses[i].ProvinsiID = strings.TrimSpace(req.Addresses[i].ProvinsiID)
		req.Addresses[i].KabupatenID = strings.TrimSpace(req.Addresses[i].KabupatenID)
		req.Addresses[i].KecamatanID = strings.TrimSpace(req.Addresses[i].KecamatanID)
		req.Addresses[i].KelurahanID = strings.TrimSpace(req.Addresses[i].KelurahanID)
	}
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
	// 0. Sanitasi Input String
	req.Sanitize()

	// A. Validasi & Parse Tanggal Lahir (string -> time.Time)
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}
	if birthDate.After(time.Now()) {
		return nil, ErrBirthDateInFuture
	}

	// B. Validasi Format NIK & No KK
	if req.NIK != "" {
		if len(req.NIK) != 16 || !isDigits(req.NIK) {
			return nil, ErrInvalidNIKFormat
		}
		existing, err := s.repo.FindByNIK(ctx, req.NIK)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek NIK: %w", err)
		}
		if existing != nil {
			return nil, ErrNIKAlreadyExists
		}
	}
	if req.FamilyCardNo != "" {
		if len(req.FamilyCardNo) != 16 || !isDigits(req.FamilyCardNo) {
			return nil, ErrInvalidFamilyCardFormat
		}
	}

	// C. Validasi Gender
	gender, err := normalizeGender(req.Gender)
	if err != nil {
		return nil, err
	}

	// D. Validasi Email jika diisi
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	// E. Penomoran Otomatis No. Rekam Medis jika tidak diisi manual
	isAutoRM := req.MedicalRecordNo == ""
	var medicalRecordNo string
	if !isAutoRM {
		if len(req.MedicalRecordNo) != 8 || !isDigits(req.MedicalRecordNo) {
			return nil, ErrInvalidMedicalRecordFormat
		}
		medicalRecordNo = req.MedicalRecordNo
		// Jika diisi manual (misal migrasi), cek apakah sudah ada yang pakai
		existingRM, err := s.repo.FindByMedicalRecordNo(ctx, medicalRecordNo)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek nomor rekam medis: %w", err)
		}
		if existingRM != nil {
			return nil, ErrMedicalRecordNoExists
		}
	}

	// F. Mapping Kontak Darurat
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

	// G. Mapping Relasi / Penanggung Jawab Keluarga
	var relations []PatientRelation
	for _, r := range req.Relations {
		relations = append(relations, PatientRelation{
			Name:      r.Name,
			Relation:  r.Relation,
			Phone:     r.Phone,
			Address:   r.Address,
			IsActive:  true,
			CreatedBy: operatorID,
		})
	}

	// H. Mapping Alamat
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

	// I. Optional: Parse Masa Berlaku Asuransi jika diisi
	var insuranceExpiry *time.Time
	if req.InsuranceExpiryDate != nil && *req.InsuranceExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.InsuranceExpiryDate)
		if err != nil {
			return nil, ErrInvalidInsuranceExpiryDateFormat
		}
		insuranceExpiry = &t
	}

	// J. Optional: Parse Tanggal Kematian jika pasien meninggal
	var deceasedTime *time.Time
	if req.IsDeceased {
		if req.DeceasedAt != nil && *req.DeceasedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", *req.DeceasedAt); err == nil {
				deceasedTime = &t
			} else if t, err := time.Parse("2006-01-02", *req.DeceasedAt); err == nil {
				deceasedTime = &t
			} else {
				return nil, ErrInvalidDeceasedDate
			}
			if deceasedTime.Before(birthDate) || deceasedTime.After(time.Now()) {
				return nil, ErrInvalidDeceasedDate
			}
		}
	}

	// K. Buat Objek Entity Pasien
	patient := &Patient{
		MedicalRecordNo:     medicalRecordNo,
		NIK:                 req.NIK,
		FamilyCardNo:        req.FamilyCardNo,
		IHSPatientID:        req.IHSPatientID,
		Title:               req.Title,
		ShortName:           req.ShortName,
		FullName:            req.FullName,
		MotherName:          req.MotherName,
		Gender:              gender,
		BirthPlace:          req.BirthPlace,
		BirthDate:           birthDate,
		Phone:               req.Phone,
		Email:               req.Email,
		MaritalStatus:       req.MaritalStatus,
		Religion:            req.Religion,
		Education:           req.Education,
		Occupation:          req.Occupation,
		Nationality:         req.Nationality,
		BloodType:           normalizeBloodType(req.BloodType),
		Rhesus:              normalizeRhesus(req.Rhesus),
		SpecialNeeds:        req.SpecialNeeds,
		IsUnknown:           req.IsUnknown,
		IsDeceased:          req.IsDeceased,
		DeceasedAt:          deceasedTime,
		CustomerID:          req.CustomerID,
		InsuranceType:       req.InsuranceType,
		InsuranceNumber:     req.InsuranceNumber,
		InsuranceExpiryDate: insuranceExpiry,
		EmergencyContacts:   emergencyContacts,
		Relations:           relations,
		Addresses:           addresses,
		CreatedBy:           operatorID,
	}

	// I. Simpan ke database via repository (dengan auto-retry jika penomoran otomatis bertabrakan saat konkurensi)
	if !isAutoRM {
		patient.MedicalRecordNo = medicalRecordNo
		return s.repo.Create(ctx, patient)
	}

	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		generatedNo, err := s.generateNextMedicalRecordNo(ctx)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat nomor rekam medis otomatis: %w", err)
		}
		patient.MedicalRecordNo = generatedNo
		patient.ID = "" // Reset ID agar selalu menggenerate UUID baru per percobaan

		created, err := s.repo.Create(ctx, patient)
		if err == nil {
			return created, nil
		}

		if attempt == maxRetries {
			return nil, fmt.Errorf("gagal mendaftarkan pasien setelah %d percobaan: %w", maxRetries, err)
		}
	}

	return nil, errors.New("gagal membuat pasien baru")
}

// UpdatePatient memperbarui data demografi pasien di loket
func (s *service) UpdatePatient(ctx context.Context, id string, req UpdatePatientRequest, operatorID string) (*Patient, error) {
	// 1. Pastikan pasien yang ingin diedit ada di database
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// 0. Sanitasi Input String
	req.Sanitize()

	// 2. Parse Tanggal Lahir
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}
	if birthDate.After(time.Now()) {
		return nil, ErrBirthDateInFuture
	}

	// 3. Validasi NIK jika diubah/diisi
	if req.NIK != "" && req.NIK != patient.NIK {
		if len(req.NIK) != 16 || !isDigits(req.NIK) {
			return nil, ErrInvalidNIKFormat
		}
		existing, err := s.repo.FindByNIK(ctx, req.NIK)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek NIK: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, ErrNIKAlreadyExists
		}
		patient.NIK = req.NIK
	}

	// 3b. Validasi No. KK jika diisi
	if req.FamilyCardNo != "" {
		if len(req.FamilyCardNo) != 16 || !isDigits(req.FamilyCardNo) {
			return nil, ErrInvalidFamilyCardFormat
		}
	}

	// 3c. Validasi Gender
	gender, err := normalizeGender(req.Gender)
	if err != nil {
		return nil, err
	}

	// 3d. Validasi Email jika diisi
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	// 4. Update field data pokok
	patient.Title = req.Title
	patient.FamilyCardNo = req.FamilyCardNo
	patient.MotherName = req.MotherName
	patient.ShortName = req.ShortName
	patient.FullName = req.FullName
	patient.Gender = gender
	patient.BirthPlace = req.BirthPlace
	patient.BirthDate = birthDate
	patient.Phone = req.Phone
	patient.Email = req.Email
	patient.MaritalStatus = req.MaritalStatus
	patient.Religion = req.Religion
	patient.Education = req.Education
	patient.Occupation = req.Occupation
	patient.Nationality = req.Nationality
	patient.BloodType = normalizeBloodType(req.BloodType)
	patient.Rhesus = normalizeRhesus(req.Rhesus)
	patient.SpecialNeeds = req.SpecialNeeds
	patient.IsUnknown = req.IsUnknown
	patient.IsDeceased = req.IsDeceased
	if req.IHSPatientID != "" {
		patient.IHSPatientID = req.IHSPatientID
	}
	if !req.IsDeceased {
		patient.DeceasedAt = nil
	} else if req.DeceasedAt != nil && *req.DeceasedAt != "" {
		var dTime *time.Time
		if t, err := time.Parse("2006-01-02 15:04:05", *req.DeceasedAt); err == nil {
			dTime = &t
		} else if t, err := time.Parse("2006-01-02", *req.DeceasedAt); err == nil {
			dTime = &t
		} else {
			return nil, ErrInvalidDeceasedDate
		}
		if dTime.Before(birthDate) || dTime.After(time.Now()) {
			return nil, ErrInvalidDeceasedDate
		}
		patient.DeceasedAt = dTime
	} else if req.DeceasedAt != nil && *req.DeceasedAt == "" {
		patient.DeceasedAt = nil
	}

	// Update data customer / penjamin & asuransi
	patient.CustomerID = req.CustomerID
	patient.InsuranceType = req.InsuranceType
	patient.InsuranceNumber = req.InsuranceNumber
	if req.InsuranceExpiryDate != nil && *req.InsuranceExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.InsuranceExpiryDate)
		if err != nil {
			return nil, ErrInvalidInsuranceExpiryDateFormat
		}
		patient.InsuranceExpiryDate = &t
	} else if req.InsuranceExpiryDate != nil && *req.InsuranceExpiryDate == "" {
		patient.InsuranceExpiryDate = nil
	}
	patient.UpdatedBy = operatorID

	// 5. Update Kontak Darurat jika dikirimkan
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

	// 5b. Update Relasi Keluarga jika dikirimkan
	if req.Relations != nil {
		var newRelations []PatientRelation
		for _, r := range req.Relations {
			newRelations = append(newRelations, PatientRelation{
				PatientID: id,
				Name:      r.Name,
				Relation:  r.Relation,
				Phone:     r.Phone,
				Address:   r.Address,
				IsActive:  true,
				CreatedBy: operatorID,
			})
		}
		patient.Relations = newRelations
	}

	// 6. Update Alamat jika dikirimkan
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

	// 7. Simpan perubahan ke repository
	return s.repo.Update(ctx, patient)
}

// DeletePatient menghapus pasien (soft delete)
func (s *service) DeletePatient(ctx context.Context, id string, operatorID string) error {
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPatientNotFound
		}
		return err
	}
	if patient == nil {
		return ErrPatientNotFound
	}

	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetPatientByID(ctx context.Context, id string) (*Patient, error) {
	patient, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}
	return patient, nil
}

func (s *service) GetPatientByNIK(ctx context.Context, nik string) (*Patient, error) {
	patient, err := s.repo.FindByNIK(ctx, nik)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}
	return patient, nil
}

func (s *service) GetPatientByMedicalRecordNo(ctx context.Context, rmNo string) (*Patient, error) {
	patient, err := s.repo.FindByMedicalRecordNo(ctx, rmNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	if patient == nil {
		return nil, ErrPatientNotFound
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
	if params.Limit > 100 {
		params.Limit = 100
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
	if nextNum > 99999999 {
		return "", ErrMaxMedicalRecordExceeded
	}
	// Format kembali menjadi 8 digit dengan awalan 0 (leading zeros)
	return fmt.Sprintf("%08d", nextNum), nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func normalizeGender(g string) (enums.Gender, error) {
	if gender, ok := enums.NormalizeGender(g); ok {
		return gender, nil
	}
	return "", ErrInvalidGender
}

func normalizeBloodType(bt string) string {
	bt = strings.ToUpper(strings.TrimSpace(bt))
	switch bt {
	case "A", "B", "AB", "O":
		return bt
	default:
		return ""
	}
}

func normalizeRhesus(rh string) string {
	switch strings.ToLower(strings.TrimSpace(rh)) {
	case "+", "positif", "positive", "pos":
		return "+"
	case "-", "negatif", "negative", "neg":
		return "-"
	default:
		return ""
	}
}
