package serviceunit

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrServiceUnitNotFound          = errors.New("unit layanan tidak ditemukan")
	ErrServiceUnitCodeAlreadyExists = errors.New("unit layanan dengan kode tersebut sudah terdaftar")
	ErrInvalidEmail                 = errors.New("format email tidak valid")
	ErrCodeRequired                 = errors.New("kode unit layanan wajib diisi")
	ErrNameRequired                 = errors.New("nama unit layanan wajib diisi")
)

type CreateServiceUnitRequest struct {
	Code                 string  `json:"code" binding:"required"`
	Name                 string  `json:"name" binding:"required"`
	DepartementID        *string `json:"departement_id"`
	Address              *string `json:"address"`
	Phone                *string `json:"phone"`
	Email                *string `json:"email"`
	Website              *string `json:"website"`
	Description          *string `json:"description"`
	IhsLocationId        *string `json:"ihs_location_id"`
	IsRegistrationTarget *bool   `json:"is_registration_target"`
	IsActive             *bool   `json:"is_active"`
}

type UpdateServiceUnitRequest struct {
	Code                 string  `json:"code" binding:"required"`
	Name                 string  `json:"name" binding:"required"`
	DepartementID        *string `json:"departement_id"`
	Address              *string `json:"address"`
	Phone                *string `json:"phone"`
	Email                *string `json:"email"`
	Website              *string `json:"website"`
	Description          *string `json:"description"`
	IhsLocationId        *string `json:"ihs_location_id"`
	IsRegistrationTarget *bool   `json:"is_registration_target"`
	IsActive             *bool   `json:"is_active"`
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
	CreateServiceUnit(ctx context.Context, req CreateServiceUnitRequest, operatorID string) (*ServiceUnit, error)
	UpdateServiceUnit(ctx context.Context, id string, req UpdateServiceUnitRequest, operatorID string) (*ServiceUnit, error)
	DeleteServiceUnit(ctx context.Context, id string, operatorID string) error
	GetServiceUnitByID(ctx context.Context, id string) (*ServiceUnit, error)
	ListServiceUnits(ctx context.Context, params ListParams) ([]ServiceUnit, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateServiceUnit(ctx context.Context, req CreateServiceUnitRequest, operatorID string) (*ServiceUnit, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}

	req.Email = sanitizeStringPtr(req.Email)
	if req.Email != nil {
		if _, err := mail.ParseAddress(*req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrServiceUnitCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isRegistrationTarget := true
	if req.IsRegistrationTarget != nil {
		isRegistrationTarget = *req.IsRegistrationTarget
	}

	su := &ServiceUnit{
		Code:                 req.Code,
		Name:                 req.Name,
		DepartementID:        sanitizeStringPtr(req.DepartementID),
		Address:              sanitizeStringPtr(req.Address),
		Phone:                sanitizeStringPtr(req.Phone),
		Email:                req.Email,
		Website:              sanitizeStringPtr(req.Website),
		Description:          sanitizeStringPtr(req.Description),
		IhsLocationId:        sanitizeStringPtr(req.IhsLocationId),
		IsRegistrationTarget: isRegistrationTarget,
		IsActive:             isActive,
		CreatedBy:            operatorID,
		UpdatedBy:            operatorID,
	}

	return s.repo.Create(ctx, su)
}

func (s *service) UpdateServiceUnit(ctx context.Context, id string, req UpdateServiceUnitRequest, operatorID string) (*ServiceUnit, error) {
	su, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceUnitNotFound
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

	req.Email = sanitizeStringPtr(req.Email)
	if req.Email != nil {
		if _, err := mail.ParseAddress(*req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	if req.Code != su.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrServiceUnitCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	su.Code = req.Code
	su.Name = req.Name
	su.DepartementID = sanitizeStringPtr(req.DepartementID)
	su.Address = sanitizeStringPtr(req.Address)
	su.Phone = sanitizeStringPtr(req.Phone)
	su.Email = req.Email
	su.Website = sanitizeStringPtr(req.Website)
	su.Description = sanitizeStringPtr(req.Description)
	su.IhsLocationId = sanitizeStringPtr(req.IhsLocationId)
	if req.IsRegistrationTarget != nil {
		su.IsRegistrationTarget = *req.IsRegistrationTarget
	}
	if req.IsActive != nil {
		su.IsActive = *req.IsActive
	}
	su.UpdatedBy = operatorID

	return s.repo.Update(ctx, su)
}

func (s *service) DeleteServiceUnit(ctx context.Context, id string, operatorID string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrServiceUnitNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetServiceUnitByID(ctx context.Context, id string) (*ServiceUnit, error) {
	su, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceUnitNotFound
		}
		return nil, err
	}
	return su, nil
}

func (s *service) ListServiceUnits(ctx context.Context, params ListParams) ([]ServiceUnit, int64, error) {
	return s.repo.FindAll(ctx, params)
}
