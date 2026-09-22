package payer

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrPayerNotFound          = errors.New("penjamin tidak ditemukan")
	ErrPayerCodeAlreadyExists = errors.New("penjamin dengan kode tersebut sudah terdaftar")
	ErrPayerImmutable         = errors.New("penjamin ini bersifat sistem/immutable dan tidak boleh dihapus")
	ErrPayerTypeNotFound      = errors.New("tipe penjamin tidak valid / tidak ditemukan")
	ErrInvalidEmail           = errors.New("format email tidak valid")
	ErrNameRequired           = errors.New("nama penjamin wajib diisi")
	ErrCodeRequired           = errors.New("kode penjamin wajib diisi")
)

type CreatePayerRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Website       *string `json:"website"`
	ContactPerson *string `json:"contact_person"`
	RequireCard   *bool   `json:"require_card"`
	Description   *string `json:"description"`
	PayerTypeID   *string `json:"payer_type_id"`
	IsActive      *bool   `json:"is_active"`
	IsImmutable   *bool   `json:"is_immutable"`
}

type UpdatePayerRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Website       *string `json:"website"`
	ContactPerson *string `json:"contact_person"`
	RequireCard   *bool   `json:"require_card"`
	Description   *string `json:"description"`
	PayerTypeID   *string `json:"payer_type_id"`
	IsActive      *bool   `json:"is_active"`
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
	CreatePayer(ctx context.Context, req CreatePayerRequest, operatorID string) (*Payer, error)
	UpdatePayer(ctx context.Context, id string, req UpdatePayerRequest, operatorID string) (*Payer, error)
	DeletePayer(ctx context.Context, id string, operatorID string) error
	GetPayerByID(ctx context.Context, id string) (*Payer, error)
	ListPayers(ctx context.Context, params ListParams) ([]Payer, int64, error)
	ListPayerTypes(ctx context.Context) ([]PayerType, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePayer(ctx context.Context, req CreatePayerRequest, operatorID string) (*Payer, error) {
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

	req.PayerTypeID = sanitizeStringPtr(req.PayerTypeID)
	if req.PayerTypeID != nil {
		if _, err := s.repo.FindPayerTypeByID(ctx, *req.PayerTypeID); err != nil {
			return nil, ErrPayerTypeNotFound
		}
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrPayerCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isImmutable := false
	if req.IsImmutable != nil {
		isImmutable = *req.IsImmutable
	}

	requireCard := true
	if req.RequireCard != nil {
		requireCard = *req.RequireCard
	}

	payer := &Payer{
		Code:          req.Code,
		Name:          req.Name,
		Address:       sanitizeStringPtr(req.Address),
		Phone:         sanitizeStringPtr(req.Phone),
		Email:         req.Email,
		Website:       sanitizeStringPtr(req.Website),
		ContactPerson: sanitizeStringPtr(req.ContactPerson),
		RequireCard:   &requireCard,
		Description:   sanitizeStringPtr(req.Description),
		PayerTypeID:   req.PayerTypeID,
		IsActive:      isActive,
		IsImmutable:   isImmutable,
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}

	return s.repo.Create(ctx, payer)
}

func (s *service) UpdatePayer(ctx context.Context, id string, req UpdatePayerRequest, operatorID string) (*Payer, error) {
	payer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPayerNotFound
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

	req.PayerTypeID = sanitizeStringPtr(req.PayerTypeID)
	if req.PayerTypeID != nil {
		if _, err := s.repo.FindPayerTypeByID(ctx, *req.PayerTypeID); err != nil {
			return nil, ErrPayerTypeNotFound
		}
	}

	if req.Code != payer.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrPayerCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	payer.Code = req.Code
	payer.Name = req.Name
	payer.Address = sanitizeStringPtr(req.Address)
	payer.Phone = sanitizeStringPtr(req.Phone)
	payer.Email = req.Email
	payer.Website = sanitizeStringPtr(req.Website)
	payer.ContactPerson = sanitizeStringPtr(req.ContactPerson)
	if req.RequireCard != nil {
		payer.RequireCard = req.RequireCard
	}
	payer.Description = sanitizeStringPtr(req.Description)
	payer.PayerTypeID = req.PayerTypeID
	if req.IsActive != nil {
		payer.IsActive = *req.IsActive
	}
	payer.UpdatedBy = operatorID

	return s.repo.Update(ctx, payer)
}

func (s *service) DeletePayer(ctx context.Context, id string, operatorID string) error {
	payer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPayerNotFound
		}
		return err
	}

	if payer.IsImmutable {
		return ErrPayerImmutable
	}

	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetPayerByID(ctx context.Context, id string) (*Payer, error) {
	payer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPayerNotFound
		}
		return nil, err
	}
	return payer, nil
}

func (s *service) ListPayers(ctx context.Context, params ListParams) ([]Payer, int64, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *service) ListPayerTypes(ctx context.Context) ([]PayerType, error) {
	return s.repo.FindAllPayerTypes(ctx)
}
