package tariffcomponent

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrTariffComponentNotFound          = errors.New("komponen tarif tidak ditemukan")
	ErrTariffComponentCodeAlreadyExists = errors.New("komponen tarif dengan kode tersebut sudah terdaftar")
	ErrCodeRequired                     = errors.New("kode komponen tarif wajib diisi")
	ErrNameRequired                     = errors.New("nama komponen tarif wajib diisi")
)

type CreateTariffComponentRequest struct {
	Code               string  `json:"code" binding:"required"`
	Name               string  `json:"name" binding:"required"`
	Description        *string `json:"description"`
	IsHospitalRevenue  bool    `json:"is_hospital_revenue"`
	IsOperatorRevenue  bool    `json:"is_operator_revenue"`
	IsParamedicRevenue bool    `json:"is_paramedic_revenue"`
	ComponentType      *string `json:"component_type"`
	DefaultCOACode     *string `json:"default_coa_code"`
	IsActive           *bool   `json:"is_active"`
}

type UpdateTariffComponentRequest struct {
	Code               string  `json:"code" binding:"required"`
	Name               string  `json:"name" binding:"required"`
	Description        *string `json:"description"`
	IsHospitalRevenue  bool    `json:"is_hospital_revenue"`
	IsOperatorRevenue  bool    `json:"is_operator_revenue"`
	IsParamedicRevenue bool    `json:"is_paramedic_revenue"`
	ComponentType      *string `json:"component_type"`
	DefaultCOACode     *string `json:"default_coa_code"`
	IsActive           *bool   `json:"is_active"`
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
	CreateTariffComponent(ctx context.Context, req CreateTariffComponentRequest, operatorID string) (*TariffComponent, error)
	UpdateTariffComponent(ctx context.Context, id string, req UpdateTariffComponentRequest, operatorID string) (*TariffComponent, error)
	DeleteTariffComponent(ctx context.Context, id string, operatorID string) error
	GetTariffComponentByID(ctx context.Context, id string) (*TariffComponent, error)
	ListTariffComponents(ctx context.Context, params ListParams) ([]TariffComponent, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (comp *service) CreateTariffComponent(ctx context.Context, req CreateTariffComponentRequest, operatorID string) (*TariffComponent, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}

	existing, err := comp.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrTariffComponentCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	compType := "LAINNYA"
	if req.ComponentType != nil && strings.TrimSpace(*req.ComponentType) != "" {
		compType = strings.TrimSpace(*req.ComponentType)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	tc := &TariffComponent{
		Code:               req.Code,
		Name:               req.Name,
		IsHospitalRevenue:  req.IsHospitalRevenue,
		IsOperatorRevenue:  req.IsOperatorRevenue,
		IsParamedicRevenue: req.IsParamedicRevenue,
		ComponentType:      compType,
		DefaultCOACode:     sanitizeStringPtr(req.DefaultCOACode),
		IsActive:           isActive,
		Description:        sanitizeStringPtr(req.Description),
		CreatedBy:          operatorID,
		UpdatedBy:          operatorID,
	}

	return comp.repo.Create(ctx, tc)
}

func (comp *service) UpdateTariffComponent(ctx context.Context, id string, req UpdateTariffComponentRequest, operatorID string) (*TariffComponent, error) {
	tc, err := comp.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTariffComponentNotFound
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
		existing, err := comp.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrTariffComponentCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	tc.Code = req.Code
	tc.Name = req.Name
	tc.Description = sanitizeStringPtr(req.Description)

	tc.IsHospitalRevenue = req.IsHospitalRevenue
	tc.IsOperatorRevenue = req.IsOperatorRevenue
	tc.IsParamedicRevenue = req.IsParamedicRevenue

	if req.ComponentType != nil && strings.TrimSpace(*req.ComponentType) != "" {
		tc.ComponentType = strings.TrimSpace(*req.ComponentType)
	}
	if req.DefaultCOACode != nil {
		tc.DefaultCOACode = sanitizeStringPtr(req.DefaultCOACode)
	}
	if req.IsActive != nil {
		tc.IsActive = *req.IsActive
	}

	tc.UpdatedBy = operatorID

	return comp.repo.Update(ctx, tc)
}

func (comp *service) DeleteTariffComponent(ctx context.Context, id string, operatorID string) error {
	_, err := comp.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTariffComponentNotFound
		}
		return err
	}

	return comp.repo.Delete(ctx, id, operatorID)
}

func (comp *service) GetTariffComponentByID(ctx context.Context, id string) (*TariffComponent, error) {
	comps, err := comp.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTariffComponentNotFound
		}
		return nil, err
	}

	return comps, nil
}

func (comp *service) ListTariffComponents(ctx context.Context, params ListParams) ([]TariffComponent, int64, error) {
	return comp.repo.FindAll(ctx, params)
}
