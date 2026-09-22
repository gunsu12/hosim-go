package tariffclass

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrTariffClassNotFound          = errors.New("kelas tarif tidak ditemukan")
	ErrTariffClassCodeAlreadyExists = errors.New("kelas tarif dengan kode tersebut sudah terdaftar")
	ErrCodeRequired                 = errors.New("kode kelas tarif wajib diisi")
	ErrNameRequired                 = errors.New("nama kelas tarif wajib diisi")
)

type CreateTariffClassRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	IsActive    *bool   `json:"is_active"`
	Description *string `json:"description"`
}

type UpdateTariffClassRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	IsActive    *bool   `json:"is_active"`
	Description *string `json:"description"`
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
	CreateTariffClass(ctx context.Context, req CreateTariffClassRequest, operatorID string) (*TariffClass, error)
	UpdateTariffClass(ctx context.Context, id string, req UpdateTariffClassRequest, operatorID string) (*TariffClass, error)
	DeleteTariffClass(ctx context.Context, id string, operatorID string) error
	GetTariffClassByID(ctx context.Context, id string) (*TariffClass, error)
	ListTariffClasses(ctx context.Context, params ListParams) ([]TariffClass, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTariffClass(ctx context.Context, req CreateTariffClassRequest, operatorID string) (*TariffClass, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrTariffClassCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	tc := &TariffClass{
		Code:        req.Code,
		Name:        req.Name,
		IsActive:    isActive,
		Description: sanitizeStringPtr(req.Description),
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}

	return s.repo.Create(ctx, tc)
}

func (s *service) UpdateTariffClass(ctx context.Context, id string, req UpdateTariffClassRequest, operatorID string) (*TariffClass, error) {
	tc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTariffClassNotFound
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

	if req.Code != tc.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrTariffClassCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	tc.Code = req.Code
	tc.Name = req.Name
	tc.Description = sanitizeStringPtr(req.Description)
	if req.IsActive != nil {
		tc.IsActive = *req.IsActive
	}
	tc.UpdatedBy = operatorID

	return s.repo.Update(ctx, tc)
}

func (s *service) DeleteTariffClass(ctx context.Context, id string, operatorID string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTariffClassNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetTariffClassByID(ctx context.Context, id string) (*TariffClass, error) {
	tc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTariffClassNotFound
		}
		return nil, err
	}
	return tc, nil
}

func (s *service) ListTariffClasses(ctx context.Context, params ListParams) ([]TariffClass, int64, error) {
	return s.repo.FindAll(ctx, params)
}
