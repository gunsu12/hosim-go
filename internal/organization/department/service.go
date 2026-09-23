package department

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

var (
	ErrDepartementNotFound     = errors.New("departemen tidak ditemukan")
	ErrDepartmentNotFound      = ErrDepartementNotFound
	ErrCodeAlreadyExists       = errors.New("departemen dengan kode tersebut sudah terdaftar")
	ErrInvalidDepartementType  = errors.New("tipe departemen tidak valid (pilihan: emergency, outpatient, inpatient, diagnostic, medical_checkup, other)")
	ErrInvalidDepartmentType   = ErrInvalidDepartementType
	ErrInvalidEmail            = errors.New("format email tidak valid")
	ErrNameRequired            = errors.New("nama departemen wajib diisi")
	ErrCodeRequired            = errors.New("kode departemen wajib diisi")
)

type CreateDepartementRequest struct {
	Code              string                `json:"code" binding:"required"`
	Name              string                `json:"name" binding:"required"`
	Address           *string               `json:"address"`
	Phone             *string               `json:"phone"`
	Email             *string               `json:"email"`
	Website           *string               `json:"website"`
	Description       *string               `json:"description"`
	DepartementType   enums.DepartementType `json:"departement_type" binding:"required"`
	IhsOrganizationId *string               `json:"ihs_organization_id"`
	IsActive          *bool                 `json:"is_active"`
}

type UpdateDepartementRequest struct {
	Code              string                `json:"code" binding:"required"`
	Name              string                `json:"name" binding:"required"`
	Address           *string               `json:"address"`
	Phone             *string               `json:"phone"`
	Email             *string               `json:"email"`
	Website           *string               `json:"website"`
	Description       *string               `json:"description"`
	DepartementType   enums.DepartementType `json:"departement_type" binding:"required"`
	IhsOrganizationId *string               `json:"ihs_organization_id"`
	IsActive          *bool                 `json:"is_active"`
}

func isValidDepartementType(t enums.DepartementType) bool {
	switch t {
	case enums.DepartementTypeEmergency,
		enums.DepartementTypeOutpatient,
		enums.DepartementTypeInpatient,
		enums.DepartementTypeDiagnostic,
		enums.DepartementTypeMCU,
		enums.DepartementTypeOther:
		return true
	default:
		return false
	}
}

func sanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

type Service interface {
	CreateDepartement(ctx context.Context, req CreateDepartementRequest, operatorID string) (*Departement, error)
	UpdateDepartement(ctx context.Context, id string, req UpdateDepartementRequest, operatorID string) (*Departement, error)
	DeleteDepartement(ctx context.Context, id string, operatorID string) error
	GetDepartementByID(ctx context.Context, id string) (*Departement, error)
	ListDepartements(ctx context.Context, params ListParams) ([]Departement, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateDepartement(ctx context.Context, req CreateDepartementRequest, operatorID string) (*Departement, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if !isValidDepartementType(req.DepartementType) {
		return nil, ErrInvalidDepartementType
	}

	req.Email = sanitizeStringPtr(req.Email)
	if req.Email != nil {
		if _, err := mail.ParseAddress(*req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	dept := &Departement{
		Code:              req.Code,
		Name:              req.Name,
		Address:           sanitizeStringPtr(req.Address),
		Phone:             sanitizeStringPtr(req.Phone),
		Email:             req.Email,
		Website:           sanitizeStringPtr(req.Website),
		Description:       sanitizeStringPtr(req.Description),
		DepartementType:   req.DepartementType,
		IhsOrganizationId: sanitizeStringPtr(req.IhsOrganizationId),
		IsActive:          isActive,
		CreatedBy:         operatorID,
		UpdatedBy:         operatorID,
	}

	return s.repo.Create(ctx, dept)
}

func (s *service) UpdateDepartement(ctx context.Context, id string, req UpdateDepartementRequest, operatorID string) (*Departement, error) {
	dept, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepartementNotFound
		}
		return nil, err
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if !isValidDepartementType(req.DepartementType) {
		return nil, ErrInvalidDepartementType
	}

	req.Email = sanitizeStringPtr(req.Email)
	if req.Email != nil {
		if _, err := mail.ParseAddress(*req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	if req.Code != dept.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	dept.Code = req.Code
	dept.Name = req.Name
	dept.Address = sanitizeStringPtr(req.Address)
	dept.Phone = sanitizeStringPtr(req.Phone)
	dept.Email = req.Email
	dept.Website = sanitizeStringPtr(req.Website)
	dept.Description = sanitizeStringPtr(req.Description)
	dept.DepartementType = req.DepartementType
	dept.IhsOrganizationId = sanitizeStringPtr(req.IhsOrganizationId)
	if req.IsActive != nil {
		dept.IsActive = *req.IsActive
	}
	dept.UpdatedBy = operatorID

	return s.repo.Update(ctx, dept)
}

func (s *service) DeleteDepartement(ctx context.Context, id string, operatorID string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDepartementNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetDepartementByID(ctx context.Context, id string) (*Departement, error) {
	dept, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepartementNotFound
		}
		return nil, err
	}
	return dept, nil
}

func (s *service) ListDepartements(ctx context.Context, params ListParams) ([]Departement, int64, error) {
	return s.repo.FindAll(ctx, params)
}
