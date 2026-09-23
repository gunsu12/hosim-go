package customer_test

import (
	"context"
	"errors"
	"testing"

	"hosim-go/internal/finance/customer"

	"gorm.io/gorm"
)

type mockCustomerRepo struct {
	customer.Repository
	createFn               func(ctx context.Context, c *customer.Customer) (*customer.Customer, error)
	updateFn               func(ctx context.Context, c *customer.Customer) (*customer.Customer, error)
	deleteFn               func(ctx context.Context, id string, deletedBy string) error
	findByIDFn             func(ctx context.Context, id string) (*customer.Customer, error)
	findByCodeFn           func(ctx context.Context, code string) (*customer.Customer, error)
	findCustomerTypeByIDFn func(ctx context.Context, id string) (*customer.CustomerType, error)
	findAllFn              func(ctx context.Context, params customer.ListParams) ([]customer.Customer, int64, error)
	findAllCustomerTypesFn func(ctx context.Context) ([]customer.CustomerType, error)
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *customer.Customer) (*customer.Customer, error) {
	if m.createFn != nil {
		return m.createFn(ctx, c)
	}
	return c, nil
}

func (m *mockCustomerRepo) Update(ctx context.Context, c *customer.Customer) (*customer.Customer, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, c)
	}
	return c, nil
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string, deletedBy string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id, deletedBy)
	}
	return nil
}

func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCustomerRepo) FindByCode(ctx context.Context, code string) (*customer.Customer, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(ctx, code)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCustomerRepo) FindCustomerTypeByID(ctx context.Context, id string) (*customer.CustomerType, error) {
	if m.findCustomerTypeByIDFn != nil {
		return m.findCustomerTypeByIDFn(ctx, id)
	}
	return &customer.CustomerType{ID: id, Name: "Asuransi"}, nil
}

func (m *mockCustomerRepo) FindAll(ctx context.Context, params customer.ListParams) ([]customer.Customer, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *mockCustomerRepo) FindAllCustomerTypes(ctx context.Context) ([]customer.CustomerType, error) {
	if m.findAllCustomerTypesFn != nil {
		return m.findAllCustomerTypesFn(ctx)
	}
	return nil, nil
}

func TestCreateCustomer_Validation(t *testing.T) {
	mockRepo := &mockCustomerRepo{}
	svc := customer.NewService(mockRepo)

	// Kode kosong
	_, err := svc.CreateCustomer(context.Background(), customer.CreateCustomerRequest{
		Code: "",
		Name: "BPJS Kesehatan",
	}, "user-1")
	if !errors.Is(err, customer.ErrCodeRequired) {
		t.Fatalf("expected ErrCodeRequired, got %v", err)
	}

	// Nama kosong
	_, err = svc.CreateCustomer(context.Background(), customer.CreateCustomerRequest{
		Code: "BPJS",
		Name: "",
	}, "user-1")
	if !errors.Is(err, customer.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}

	// Email tidak valid
	badEmail := "invalid-email"
	_, err = svc.CreateCustomer(context.Background(), customer.CreateCustomerRequest{
		Code:  "BPJS",
		Name:  "BPJS Kesehatan",
		Email: &badEmail,
	}, "user-1")
	if !errors.Is(err, customer.ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestCreateCustomer_DuplicateCode(t *testing.T) {
	mockRepo := &mockCustomerRepo{
		findByCodeFn: func(ctx context.Context, code string) (*customer.Customer, error) {
			return &customer.Customer{ID: "c-1", Code: code}, nil
		},
	}
	svc := customer.NewService(mockRepo)

	_, err := svc.CreateCustomer(context.Background(), customer.CreateCustomerRequest{
		Code: "BPJS",
		Name: "BPJS Kesehatan",
	}, "user-1")
	if !errors.Is(err, customer.ErrCustomerCodeAlreadyExists) {
		t.Fatalf("expected ErrCustomerCodeAlreadyExists, got %v", err)
	}
}

func TestCreateCustomer_Success(t *testing.T) {
	var saved *customer.Customer
	mockRepo := &mockCustomerRepo{
		createFn: func(ctx context.Context, c *customer.Customer) (*customer.Customer, error) {
			saved = c
			return c, nil
		},
	}
	svc := customer.NewService(mockRepo)

	res, err := svc.CreateCustomer(context.Background(), customer.CreateCustomerRequest{
		Code: "BPJS",
		Name: "BPJS Kesehatan",
	}, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Code != "BPJS" || saved.Name != "BPJS Kesehatan" {
		t.Fatalf("customer data mismatch: %v", res)
	}
}

func TestDeleteCustomer_Immutable(t *testing.T) {
	mockRepo := &mockCustomerRepo{
		findByIDFn: func(ctx context.Context, id string) (*customer.Customer, error) {
			return &customer.Customer{
				ID:          id,
				Code:        "PRIBADI",
				Name:        "Pasien Umum / Tunai",
				IsImmutable: true,
			}, nil
		},
	}
	svc := customer.NewService(mockRepo)

	err := svc.DeleteCustomer(context.Background(), "cust-pribadi", "user-1")
	if !errors.Is(err, customer.ErrCustomerImmutable) {
		t.Fatalf("expected ErrCustomerImmutable, got %v", err)
	}
}
