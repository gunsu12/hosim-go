package item_test

import (
	"context"
	"errors"
	"testing"

	"hosim-go/internal/catalog/item"
	"hosim-go/pkg/enums"
)

type mockCategoryProductLineRepo struct {
	item.Repository
	findCategoryByCodeFn    func(ctx context.Context, code string) (*item.ItemCategory, error)
	createCategoryFn        func(ctx context.Context, cat *item.ItemCategory) error
	findCategoryByIDFn      func(ctx context.Context, id string) (*item.ItemCategory, error)
	updateCategoryFn        func(ctx context.Context, cat *item.ItemCategory) error
	deleteCategoryFn        func(ctx context.Context, id string, deletedBy string) error
	findAllCategoriesFn     func(ctx context.Context, params item.CategoryListParams) ([]item.ItemCategory, int64, error)

	findProductLineByCodeFn func(ctx context.Context, code string) (*item.ItemProductLine, error)
	createProductLineFn     func(ctx context.Context, pl *item.ItemProductLine) error
	findProductLineByIDFn   func(ctx context.Context, id string) (*item.ItemProductLine, error)
	updateProductLineFn     func(ctx context.Context, pl *item.ItemProductLine) error
	deleteProductLineFn     func(ctx context.Context, id string, deletedBy string) error
	findAllProductLinesFn   func(ctx context.Context, params item.ProductLineListParams) ([]item.ItemProductLine, int64, error)
}

func (m *mockCategoryProductLineRepo) FindCategoryByCode(ctx context.Context, code string) (*item.ItemCategory, error) {
	if m.findCategoryByCodeFn != nil {
		return m.findCategoryByCodeFn(ctx, code)
	}
	return nil, item.ErrCategoryNotFound
}

func (m *mockCategoryProductLineRepo) CreateCategory(ctx context.Context, cat *item.ItemCategory) error {
	if m.createCategoryFn != nil {
		return m.createCategoryFn(ctx, cat)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) FindCategoryByID(ctx context.Context, id string) (*item.ItemCategory, error) {
	if m.findCategoryByIDFn != nil {
		return m.findCategoryByIDFn(ctx, id)
	}
	return nil, item.ErrCategoryNotFound
}

func (m *mockCategoryProductLineRepo) UpdateCategory(ctx context.Context, cat *item.ItemCategory) error {
	if m.updateCategoryFn != nil {
		return m.updateCategoryFn(ctx, cat)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) DeleteCategory(ctx context.Context, id string, deletedBy string) error {
	if m.deleteCategoryFn != nil {
		return m.deleteCategoryFn(ctx, id, deletedBy)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) FindAllCategories(ctx context.Context, params item.CategoryListParams) ([]item.ItemCategory, int64, error) {
	if m.findAllCategoriesFn != nil {
		return m.findAllCategoriesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *mockCategoryProductLineRepo) FindProductLineByCode(ctx context.Context, code string) (*item.ItemProductLine, error) {
	if m.findProductLineByCodeFn != nil {
		return m.findProductLineByCodeFn(ctx, code)
	}
	return nil, item.ErrProductLineNotFound
}

func (m *mockCategoryProductLineRepo) CreateProductLine(ctx context.Context, pl *item.ItemProductLine) error {
	if m.createProductLineFn != nil {
		return m.createProductLineFn(ctx, pl)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) FindProductLineByID(ctx context.Context, id string) (*item.ItemProductLine, error) {
	if m.findProductLineByIDFn != nil {
		return m.findProductLineByIDFn(ctx, id)
	}
	return nil, item.ErrProductLineNotFound
}

func (m *mockCategoryProductLineRepo) UpdateProductLine(ctx context.Context, pl *item.ItemProductLine) error {
	if m.updateProductLineFn != nil {
		return m.updateProductLineFn(ctx, pl)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) DeleteProductLine(ctx context.Context, id string, deletedBy string) error {
	if m.deleteProductLineFn != nil {
		return m.deleteProductLineFn(ctx, id, deletedBy)
	}
	return nil
}

func (m *mockCategoryProductLineRepo) FindAllProductLines(ctx context.Context, params item.ProductLineListParams) ([]item.ItemProductLine, int64, error) {
	if m.findAllProductLinesFn != nil {
		return m.findAllProductLinesFn(ctx, params)
	}
	return nil, 0, nil
}

func TestCategoryService_CreateSuccessAndDuplicate(t *testing.T) {
	mockRepo := &mockCategoryProductLineRepo{
		findCategoryByCodeFn: func(ctx context.Context, code string) (*item.ItemCategory, error) {
			if code == "DUPLICATE" {
				return &item.ItemCategory{ID: "cat-1", Code: code}, nil
			}
			return nil, item.ErrCategoryNotFound
		},
	}

	service := item.NewCategoryService(mockRepo)

	// Test Duplicate
	_, err := service.Create(context.Background(), item.CreateCategoryRequest{
		Code:     "DUPLICATE",
		Name:     "Test Cat",
		ItemType: enums.ItemTypeMedication,
	}, "tester")
	if !errors.Is(err, item.ErrCategoryCodeAlreadyExists) {
		t.Fatalf("expected ErrCategoryCodeAlreadyExists, got %v", err)
	}

	// Test Success
	res, err := service.Create(context.Background(), item.CreateCategoryRequest{
		Code:     "NEW-CAT",
		Name:     "Kategori Baru",
		ItemType: enums.ItemTypeMedication,
	}, "tester")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Code != "NEW-CAT" {
		t.Errorf("expected code NEW-CAT, got %s", res.Code)
	}
}

func TestProductLineService_CreateSuccessAndDuplicate(t *testing.T) {
	mockRepo := &mockCategoryProductLineRepo{
		findProductLineByCodeFn: func(ctx context.Context, code string) (*item.ItemProductLine, error) {
			if code == "DUPLICATE" {
				return &item.ItemProductLine{ID: "pl-1", Code: code}, nil
			}
			return nil, item.ErrProductLineNotFound
		},
	}

	service := item.NewProductLineService(mockRepo)

	// Test Duplicate
	_, err := service.Create(context.Background(), item.CreateProductLineRequest{
		Code: "DUPLICATE",
		Name: "Test PL",
	}, "tester")
	if !errors.Is(err, item.ErrProductLineCodeAlreadyExists) {
		t.Fatalf("expected ErrProductLineCodeAlreadyExists, got %v", err)
	}

	// Test Success
	res, err := service.Create(context.Background(), item.CreateProductLineRequest{
		Code: "NEW-PL",
		Name: "Lini Produk Baru",
	}, "tester")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Code != "NEW-PL" {
		t.Errorf("expected code NEW-PL, got %s", res.Code)
	}
}
