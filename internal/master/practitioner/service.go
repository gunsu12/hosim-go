package practitioner

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

var (
	ErrPractitionerNotFound         = errors.New("tenaga medis tidak ditemukan")
	ErrNIKAlreadyExists             = errors.New("tenaga medis dengan NIK tersebut sudah terdaftar")
	ErrNIPAlreadyExists             = errors.New("tenaga medis dengan NIP tersebut sudah terdaftar")
	ErrIhsAlreadyExists             = errors.New("tenaga medis dengan IHS ID tersebut sudah terdaftar")
	ErrInvalidNIKFormat             = errors.New("NIK harus berjumlah 16 digit angka")
	ErrInvalidGender                = errors.New("jenis kelamin tidak valid (gunakan 'L' / 'P' atau 'male' / 'female')")
	ErrInvalidEmail                 = errors.New("format email tidak valid")
	ErrInvalidSIPExpiryDateFormat   = errors.New("format tanggal masa berlaku SIP tidak valid (gunakan format YYYY-MM-DD)")
)

type CreatePractitionerRequest struct {
	NIK                 string  `json:"nik"`                   // Opsional: boleh kosong (misal WNA / belum sync)
	NIP                 string  `json:"nip"`                   // Opsional: No Induk Pegawai
	Name                string  `json:"name" binding:"required"`
	Gender              string  `json:"gender" binding:"required"` // 'L'/'P' atau 'male'/'female'
	SIP                 string  `json:"sip"`                   // Surat Izin Praktik
	SIPExpiryDate       *string `json:"sip_expiry_date"`       // Format: YYYY-MM-DD (opsional)
	STR                 string  `json:"str"`                   // Surat Tanda Registrasi
	ProfessionID        *string `json:"profession_id"`         // UUID relasi ke Profession (opsional)
	SpecialtyID         *string `json:"specialty_id"`          // UUID relasi ke Specialty (opsional)
	Phone               string  `json:"phone"`
	Email               string  `json:"email"`
	IhsPractitionerID   string  `json:"ihs_practitioner_id"`   // SatuSehat Practitioner ID (opsional)
	IhsPractitionerName string  `json:"ihs_practitioner_name"` // Nama nakes sesuai SatuSehat (opsional)
	IsActive            *bool   `json:"is_active"`             // Opsional, default true
}

func (req *CreatePractitionerRequest) Sanitize() {
	req.NIK = strings.TrimSpace(req.NIK)
	req.NIP = strings.TrimSpace(req.NIP)
	req.Name = strings.TrimSpace(req.Name)
	req.Gender = strings.TrimSpace(req.Gender)
	req.SIP = strings.TrimSpace(req.SIP)
	req.STR = strings.TrimSpace(req.STR)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.IhsPractitionerID = strings.TrimSpace(req.IhsPractitionerID)
	req.IhsPractitionerName = strings.TrimSpace(req.IhsPractitionerName)

	if req.SIPExpiryDate != nil {
		trimmed := strings.TrimSpace(*req.SIPExpiryDate)
		if trimmed == "" {
			req.SIPExpiryDate = nil
		} else {
			req.SIPExpiryDate = &trimmed
		}
	}

	// Ubah UUID string kosong menjadi nil agar foreign key PostgreSQL tidak melanggar constraint
	if req.ProfessionID != nil {
		trimmed := strings.TrimSpace(*req.ProfessionID)
		if trimmed == "" {
			req.ProfessionID = nil
		} else {
			req.ProfessionID = &trimmed
		}
	}

	if req.SpecialtyID != nil {
		trimmed := strings.TrimSpace(*req.SpecialtyID)
		if trimmed == "" {
			req.SpecialtyID = nil
		} else {
			req.SpecialtyID = &trimmed
		}
	}
}

type UpdatePractitionerRequest struct {
	CreatePractitionerRequest
}

type Service interface {
	CreatePractitioner(ctx context.Context, req CreatePractitionerRequest, operatorID string) (*Practitioner, error)
	UpdatePractitioner(ctx context.Context, id string, req UpdatePractitionerRequest, operatorID string) (*Practitioner, error)
	DeletePractitioner(ctx context.Context, id string, operatorID string) error
	GetPractitionerByID(ctx context.Context, id string) (*Practitioner, error)
	GetPractitionerByNIK(ctx context.Context, nik string) (*Practitioner, error)
	GetPractitionerByNIP(ctx context.Context, nip string) (*Practitioner, error)
	GetPractitionerByIhsID(ctx context.Context, ihsID string) (*Practitioner, error)
	ListPractitioners(ctx context.Context, params ListParams) ([]Practitioner, int, error)
	ListProfessions(ctx context.Context) ([]Profession, error)
	ListSpecialties(ctx context.Context) ([]Specialty, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePractitioner(ctx context.Context, req CreatePractitionerRequest, operatorID string) (*Practitioner, error) {
	// 0. Sanitasi Input
	req.Sanitize()

	// 1. Validasi Gender
	gender, err := normalizeGender(req.Gender)
	if err != nil {
		return nil, err
	}

	// 2. Validasi NIK jika diisi
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

	// 3. Validasi NIP jika diisi
	if req.NIP != "" {
		existing, err := s.repo.FindByNIP(ctx, req.NIP)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek NIP: %w", err)
		}
		if existing != nil {
			return nil, ErrNIPAlreadyExists
		}
	}

	// 4. Validasi IHS Practitioner ID jika diisi
	if req.IhsPractitionerID != "" {
		existing, err := s.repo.FindByIhsPractitionerID(ctx, req.IhsPractitionerID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek IHS Practitioner ID: %w", err)
		}
		if existing != nil {
			return nil, ErrIhsAlreadyExists
		}
	}

	// 5. Validasi Email jika diisi
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	// 6. Validasi Tanggal Expired SIP jika diisi
	var parsedSIPExpiryDate *time.Time
	if req.SIPExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.SIPExpiryDate)
		if err != nil {
			return nil, ErrInvalidSIPExpiryDateFormat
		}
		parsedSIPExpiryDate = &t
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	practitioner := &Practitioner{
		NIK:                 req.NIK,
		NIP:                 req.NIP,
		Name:                req.Name,
		Gender:              gender,
		SIP:                 req.SIP,
		SIPExpiryDate:       parsedSIPExpiryDate,
		STR:                 req.STR,
		ProfessionID:        req.ProfessionID,
		SpecialtyID:         req.SpecialtyID,
		Phone:               req.Phone,
		Email:               req.Email,
		IhsPractitionerID:   req.IhsPractitionerID,
		IhsPractitionerName: req.IhsPractitionerName,
		IsActive:            isActive,
		CreatedBy:           operatorID,
		UpdatedBy:           operatorID,
	}

	created, err := s.repo.Create(ctx, practitioner)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan tenaga medis: %w", err)
	}

	// Ambil kembali data lengkap beserta relasi Preload
	return s.repo.FindByID(ctx, created.ID)
}

func (s *service) UpdatePractitioner(ctx context.Context, id string, req UpdatePractitionerRequest, operatorID string) (*Practitioner, error) {
	// 0. Sanitasi Input
	req.Sanitize()

	// 1. Cek keberadaan data
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPractitionerNotFound
		}
		return nil, fmt.Errorf("gagal mencari tenaga medis: %w", err)
	}

	// 2. Validasi Gender
	gender, err := normalizeGender(req.Gender)
	if err != nil {
		return nil, err
	}

	// 3. Validasi NIK jika diisi dan berubah
	if req.NIK != "" {
		if len(req.NIK) != 16 || !isDigits(req.NIK) {
			return nil, ErrInvalidNIKFormat
		}
		if req.NIK != existing.NIK {
			found, err := s.repo.FindByNIK(ctx, req.NIK)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("gagal mengecek NIK: %w", err)
			}
			if found != nil && found.ID != id {
				return nil, ErrNIKAlreadyExists
			}
		}
	}

	// 4. Validasi NIP jika diisi dan berubah
	if req.NIP != "" && req.NIP != existing.NIP {
		found, err := s.repo.FindByNIP(ctx, req.NIP)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek NIP: %w", err)
		}
		if found != nil && found.ID != id {
			return nil, ErrNIPAlreadyExists
		}
	}

	// 5. Validasi IHS Practitioner ID jika diisi dan berubah
	if req.IhsPractitionerID != "" && req.IhsPractitionerID != existing.IhsPractitionerID {
		found, err := s.repo.FindByIhsPractitionerID(ctx, req.IhsPractitionerID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("gagal mengecek IHS Practitioner ID: %w", err)
		}
		if found != nil && found.ID != id {
			return nil, ErrIhsAlreadyExists
		}
	}

	// 6. Validasi Email jika diisi
	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	// 7. Validasi Tanggal Expired SIP jika diisi
	var parsedSIPExpiryDate *time.Time
	if req.SIPExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.SIPExpiryDate)
		if err != nil {
			return nil, ErrInvalidSIPExpiryDateFormat
		}
		parsedSIPExpiryDate = &t
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	existing.NIK = req.NIK
	existing.NIP = req.NIP
	existing.Name = req.Name
	existing.Gender = gender
	existing.SIP = req.SIP
	existing.SIPExpiryDate = parsedSIPExpiryDate
	existing.STR = req.STR
	existing.ProfessionID = req.ProfessionID
	existing.SpecialtyID = req.SpecialtyID
	existing.Phone = req.Phone
	existing.Email = req.Email
	existing.IhsPractitionerID = req.IhsPractitionerID
	existing.IhsPractitionerName = req.IhsPractitionerName
	existing.IsActive = isActive
	existing.UpdatedBy = operatorID

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("gagal mengupdate tenaga medis: %w", err)
	}

	return s.repo.FindByID(ctx, updated.ID)
}

func (s *service) DeletePractitioner(ctx context.Context, id string, operatorID string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPractitionerNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, existing.ID, operatorID)
}

func (s *service) GetPractitionerByID(ctx context.Context, id string) (*Practitioner, error) {
	practitioner, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPractitionerNotFound
		}
		return nil, err
	}
	return practitioner, nil
}

func (s *service) GetPractitionerByNIK(ctx context.Context, nik string) (*Practitioner, error) {
	trimmed := strings.TrimSpace(nik)
	if trimmed == "" {
		return nil, ErrPractitionerNotFound
	}

	practitioner, err := s.repo.FindByNIK(ctx, trimmed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPractitionerNotFound
		}
		return nil, err
	}
	return practitioner, nil
}

func (s *service) GetPractitionerByNIP(ctx context.Context, nip string) (*Practitioner, error) {
	trimmed := strings.TrimSpace(nip)
	if trimmed == "" {
		return nil, ErrPractitionerNotFound
	}

	practitioner, err := s.repo.FindByNIP(ctx, trimmed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPractitionerNotFound
		}
		return nil, err
	}
	return practitioner, nil
}

func (s *service) GetPractitionerByIhsID(ctx context.Context, ihsID string) (*Practitioner, error) {
	trimmed := strings.TrimSpace(ihsID)
	if trimmed == "" {
		return nil, ErrPractitionerNotFound
	}

	practitioner, err := s.repo.FindByIhsPractitionerID(ctx, trimmed)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPractitionerNotFound
		}
		return nil, err
	}
	return practitioner, nil
}

func (s *service) ListPractitioners(ctx context.Context, params ListParams) ([]Practitioner, int, error) {
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

func (s *service) ListProfessions(ctx context.Context) ([]Profession, error) {
	return s.repo.FindAllProfessions(ctx)
}

func (s *service) ListSpecialties(ctx context.Context) ([]Specialty, error) {
	return s.repo.FindAllSpecialties(ctx)
}

func normalizeGender(raw string) (enums.Gender, error) {
	g, ok := enums.NormalizeGender(raw)
	if !ok || !g.IsValid() {
		return "", ErrInvalidGender
	}
	return g, nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
