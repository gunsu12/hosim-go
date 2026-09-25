package referal

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

var (
	ErrReferalNotFound          = errors.New("rujukan tidak ditemukan")
	ErrReferalCodeAlreadyExists = errors.New("rujukan dengan kode tersebut sudah terdaftar")
	ErrInvalidReferalType       = errors.New("tipe rujukan tidak valid (pilihan: DOCTOR, HOSPITAL, PHARMACY, LAB, RADIOLOGY, THERAPY, OTHER)")
	ErrInvalidEmail             = errors.New("format email tidak valid")
	ErrCodeRequired             = errors.New("kode rujukan wajib diisi")
	ErrNameRequired             = errors.New("nama rujukan wajib diisi")
)

type CreateReferalRequest struct {
	Code    string            `json:"code" binding:"required"`
	Name    string            `json:"name" binding:"required"`
	Type    enums.ReferalType `json:"type" binding:"required"`
	Address string            `json:"address"`
	Phone   string            `json:"phone"`
	Email   string            `json:"email"`
}

type UpdateReferalRequest struct {
	Code    string            `json:"code" binding:"required"`
	Name    string            `json:"name" binding:"required"`
	Type    enums.ReferalType `json:"type" binding:"required"`
	Address string            `json:"address"`
	Phone   string            `json:"phone"`
	Email   string            `json:"email"`
}

func isValidReferalType(t enums.ReferalType) bool {
	return t.IsValid()
}

type Service interface {
	CreateReferal(ctx context.Context, req CreateReferalRequest, operatorID string) (*Referal, error)
	UpdateReferal(ctx context.Context, id string, req UpdateReferalRequest, operatorID string) (*Referal, error)
	DeleteReferal(ctx context.Context, id string, operatorID string) error
	GetReferalByID(ctx context.Context, id string) (*Referal, error)
	ListReferals(ctx context.Context, params ListParams) ([]Referal, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateReferal(ctx context.Context, req CreateReferalRequest, operatorID string) (*Referal, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.TrimSpace(req.Address)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)

	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if !isValidReferalType(req.Type) {
		return nil, ErrInvalidReferalType
	}

	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrReferalCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	ref := &Referal{
		Code:      req.Code,
		Name:      req.Name,
		Type:      req.Type,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	return s.repo.Create(ctx, ref)
}

func (s *service) UpdateReferal(ctx context.Context, id string, req UpdateReferalRequest, operatorID string) (*Referal, error) {
	ref, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReferalNotFound
		}
		return nil, err
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.TrimSpace(req.Address)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)

	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if !isValidReferalType(req.Type) {
		return nil, ErrInvalidReferalType
	}

	if req.Email != "" {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return nil, ErrInvalidEmail
		}
	}

	if req.Code != ref.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrReferalCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	ref.Code = req.Code
	ref.Name = req.Name
	ref.Type = req.Type
	ref.Address = req.Address
	ref.Phone = req.Phone
	ref.Email = req.Email
	ref.UpdatedBy = operatorID

	return s.repo.Update(ctx, ref)
}

func (s *service) DeleteReferal(ctx context.Context, id string, operatorID string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReferalNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetReferalByID(ctx context.Context, id string) (*Referal, error) {
	ref, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReferalNotFound
		}
		return nil, err
	}
	return ref, nil
}

func (s *service) ListReferals(ctx context.Context, params ListParams) ([]Referal, int64, error) {
	return s.repo.FindAll(ctx, params)
}
