package customer

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrCustomerNotFound          = errors.New("customer/penjamin tidak ditemukan")
	ErrCustomerCodeAlreadyExists = errors.New("customer/penjamin dengan kode tersebut sudah terdaftar")
	ErrCustomerImmutable         = errors.New("customer/penjamin ini bersifat sistem/immutable dan tidak boleh dihapus")
	ErrCustomerTypeNotFound      = errors.New("tipe customer/penjamin tidak valid / tidak ditemukan")
	ErrInvalidEmail              = errors.New("format email tidak valid")
	ErrNameRequired              = errors.New("nama customer/penjamin wajib diisi")
	ErrCodeRequired              = errors.New("kode customer/penjamin wajib diisi")
)

type CreateCustomerRequest struct {
	Code           string  `json:"code" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Address        *string `json:"address"`
	Phone          *string `json:"phone"`
	Email          *string `json:"email"`
	Website        *string `json:"website"`
	ContactPerson  *string `json:"contact_person"`
	RequireCard    *bool   `json:"require_card"`
	Description    *string `json:"description"`
	CustomerTypeID *string `json:"customer_type_id"`
	IsActive       *bool   `json:"is_active"`
	IsImmutable    *bool   `json:"is_immutable"`
}

type UpdateCustomerRequest struct {
	Code           string  `json:"code" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Address        *string `json:"address"`
	Phone          *string `json:"phone"`
	Email          *string `json:"email"`
	Website        *string `json:"website"`
	ContactPerson  *string `json:"contact_person"`
	RequireCard    *bool   `json:"require_card"`
	Description    *string `json:"description"`
	CustomerTypeID *string `json:"customer_type_id"`
	IsActive       *bool   `json:"is_active"`
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
	CreateCustomer(ctx context.Context, req CreateCustomerRequest, operatorID string) (*Customer, error)
	UpdateCustomer(ctx context.Context, id string, req UpdateCustomerRequest, operatorID string) (*Customer, error)
	DeleteCustomer(ctx context.Context, id string, operatorID string) error
	GetCustomerByID(ctx context.Context, id string) (*Customer, error)
	ListCustomers(ctx context.Context, params ListParams) ([]Customer, int64, error)
	ListCustomerTypes(ctx context.Context) ([]CustomerType, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCustomer(ctx context.Context, req CreateCustomerRequest, operatorID string) (*Customer, error) {
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

	req.CustomerTypeID = sanitizeStringPtr(req.CustomerTypeID)
	if req.CustomerTypeID != nil {
		if _, err := s.repo.FindCustomerTypeByID(ctx, *req.CustomerTypeID); err != nil {
			return nil, ErrCustomerTypeNotFound
		}
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrCustomerCodeAlreadyExists
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

	customer := &Customer{
		Code:           req.Code,
		Name:           req.Name,
		Address:        sanitizeStringPtr(req.Address),
		Phone:          sanitizeStringPtr(req.Phone),
		Email:          req.Email,
		Website:        sanitizeStringPtr(req.Website),
		ContactPerson:  sanitizeStringPtr(req.ContactPerson),
		RequireCard:    &requireCard,
		Description:    sanitizeStringPtr(req.Description),
		CustomerTypeID: req.CustomerTypeID,
		IsActive:       isActive,
		IsImmutable:    isImmutable,
		CreatedBy:      operatorID,
		UpdatedBy:      operatorID,
	}

	return s.repo.Create(ctx, customer)
}

func (s *service) UpdateCustomer(ctx context.Context, id string, req UpdateCustomerRequest, operatorID string) (*Customer, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
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

	req.CustomerTypeID = sanitizeStringPtr(req.CustomerTypeID)
	if req.CustomerTypeID != nil {
		if _, err := s.repo.FindCustomerTypeByID(ctx, *req.CustomerTypeID); err != nil {
			return nil, ErrCustomerTypeNotFound
		}
	}

	if req.Code != customer.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCustomerCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	customer.Code = req.Code
	customer.Name = req.Name
	customer.Address = sanitizeStringPtr(req.Address)
	customer.Phone = sanitizeStringPtr(req.Phone)
	customer.Email = req.Email
	customer.Website = sanitizeStringPtr(req.Website)
	customer.ContactPerson = sanitizeStringPtr(req.ContactPerson)
	if req.RequireCard != nil {
		customer.RequireCard = req.RequireCard
	}
	customer.Description = sanitizeStringPtr(req.Description)
	customer.CustomerTypeID = req.CustomerTypeID
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}
	customer.UpdatedBy = operatorID

	return s.repo.Update(ctx, customer)
}

func (s *service) DeleteCustomer(ctx context.Context, id string, operatorID string) error {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomerNotFound
		}
		return err
	}

	if customer.IsImmutable {
		return ErrCustomerImmutable
	}

	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetCustomerByID(ctx context.Context, id string) (*Customer, error) {
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return customer, nil
}

func (s *service) ListCustomers(ctx context.Context, params ListParams) ([]Customer, int64, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *service) ListCustomerTypes(ctx context.Context) ([]CustomerType, error) {
	return s.repo.FindAllCustomerTypes(ctx)
}
