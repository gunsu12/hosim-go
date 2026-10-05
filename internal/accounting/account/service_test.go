package account_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"hosim-go/internal/accounting/account"
	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

type mockAccountRepo struct {
	createFn                 func(ctx context.Context, acc *account.ChartOfAccount) error
	createWithParentUpdateFn func(ctx context.Context, acc *account.ChartOfAccount, parentID string) error
	updateFn                 func(ctx context.Context, acc *account.ChartOfAccount) error
	deleteFn                 func(ctx context.Context, id, deletedBy string) error
	deleteWithParentCheckFn  func(ctx context.Context, id string, parentID *string, deletedBy string) error
	findByIDFn               func(ctx context.Context, id string) (*account.ChartOfAccount, error)
	findByCodeFn             func(ctx context.Context, code string) (*account.ChartOfAccount, error)
	findAllFn                func(ctx context.Context, params account.AccountListParams) ([]account.ChartOfAccount, int64, error)
	findAllRawFn             func(ctx context.Context, activeOnly bool) ([]account.ChartOfAccount, error)
	findChildrenCountFn      func(ctx context.Context, parentID string) (int64, error)
	isCodeUsedInTariffFn     func(ctx context.Context, code string) (bool, error)
	setPostableFn            func(ctx context.Context, id string, isPostable bool) error
}

func (m *mockAccountRepo) Create(ctx context.Context, acc *account.ChartOfAccount) error {
	if m.createFn != nil {
		return m.createFn(ctx, acc)
	}
	return nil
}

func (m *mockAccountRepo) CreateWithParentUpdate(ctx context.Context, acc *account.ChartOfAccount, parentID string) error {
	if m.createWithParentUpdateFn != nil {
		return m.createWithParentUpdateFn(ctx, acc, parentID)
	}
	return nil
}

func (m *mockAccountRepo) Update(ctx context.Context, acc *account.ChartOfAccount) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, acc)
	}
	return nil
}

func (m *mockAccountRepo) Delete(ctx context.Context, id, deletedBy string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id, deletedBy)
	}
	return nil
}

func (m *mockAccountRepo) DeleteWithParentCheck(ctx context.Context, id string, parentID *string, deletedBy string) error {
	if m.deleteWithParentCheckFn != nil {
		return m.deleteWithParentCheckFn(ctx, id, parentID, deletedBy)
	}
	return nil
}

func (m *mockAccountRepo) FindByID(ctx context.Context, id string) (*account.ChartOfAccount, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAccountRepo) FindByCode(ctx context.Context, code string) (*account.ChartOfAccount, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(ctx, code)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAccountRepo) FindAll(ctx context.Context, params account.AccountListParams) ([]account.ChartOfAccount, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx, params)
	}
	return []account.ChartOfAccount{}, 0, nil
}

func (m *mockAccountRepo) FindAllRaw(ctx context.Context, activeOnly bool) ([]account.ChartOfAccount, error) {
	if m.findAllRawFn != nil {
		return m.findAllRawFn(ctx, activeOnly)
	}
	return []account.ChartOfAccount{}, nil
}

func (m *mockAccountRepo) FindChildrenCount(ctx context.Context, parentID string) (int64, error) {
	if m.findChildrenCountFn != nil {
		return m.findChildrenCountFn(ctx, parentID)
	}
	return 0, nil
}

func (m *mockAccountRepo) IsCodeUsedInTariff(ctx context.Context, code string) (bool, error) {
	if m.isCodeUsedInTariffFn != nil {
		return m.isCodeUsedInTariffFn(ctx, code)
	}
	return false, nil
}

func (m *mockAccountRepo) SetPostable(ctx context.Context, id string, isPostable bool) error {
	if m.setPostableFn != nil {
		return m.setPostableFn(ctx, id, isPostable)
	}
	return nil
}

func TestService_CreateAccount_RootAccount(t *testing.T) {
	ctx := context.Background()
	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(ctx context.Context, acc *account.ChartOfAccount) error {
			acc.ID = "uuid-root-1000"
			acc.CreatedAt = time.Now()
			acc.UpdatedAt = time.Now()
			return nil
		},
	}

	svc := account.NewService(repo)

	req := account.CreateAccountRequest{
		Code: "1000",
		Name: "ASET",
		Type: enums.AccountTypeAsset,
	}

	res, err := svc.CreateAccount(ctx, req, "admin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res.Code != "1000" {
		t.Errorf("expected code 1000, got: %s", res.Code)
	}
	if res.AccountLevel != 1 {
		t.Errorf("expected level 1 for root, got: %d", res.AccountLevel)
	}
	if res.Position != enums.AccountPositionDebit {
		t.Errorf("expected default position DEBIT for ASSET, got: %v", res.Position)
	}
	if !res.IsPostable {
		t.Errorf("expected is_postable=true by default for new leaf account")
	}
}

func TestService_CreateAccount_ChildAccount(t *testing.T) {
	ctx := context.Background()
	parentID := "uuid-root-1000"
	parentCode := "1000"

	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findByIDFn: func(ctx context.Context, id string) (*account.ChartOfAccount, error) {
			if id == parentID {
				return &account.ChartOfAccount{
					ID:           parentID,
					Code:         parentCode,
					Name:         "ASET",
					Type:         enums.AccountTypeAsset,
					AccountLevel: 1,
					IsPostable:   true,
				}, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		createWithParentUpdateFn: func(ctx context.Context, acc *account.ChartOfAccount, pID string) error {
			if pID != parentID {
				t.Errorf("expected parentID %s, got: %s", parentID, pID)
			}
			acc.ID = "uuid-child-1100"
			acc.CreatedAt = time.Now()
			acc.UpdatedAt = time.Now()
			return nil
		},
	}

	svc := account.NewService(repo)

	req := account.CreateAccountRequest{
		Code:     "1100",
		Name:     "ASET LANCAR",
		ParentID: &parentID,
		Type:     enums.AccountTypeAsset,
	}

	res, err := svc.CreateAccount(ctx, req, "admin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res.AccountLevel != 2 {
		t.Errorf("expected level 2 for child of level 1, got: %d", res.AccountLevel)
	}
	if res.ParentCode == nil || *res.ParentCode != "1000" {
		t.Errorf("expected parent_code 1000, got: %v", res.ParentCode)
	}
}

func TestService_CreateAccount_TypeMismatchWithParent(t *testing.T) {
	ctx := context.Background()
	parentID := "uuid-root-1000"

	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return nil, gorm.ErrRecordNotFound
		},
		findByIDFn: func(ctx context.Context, id string) (*account.ChartOfAccount, error) {
			return &account.ChartOfAccount{
				ID:           parentID,
				Code:         "1000",
				Name:         "ASET",
				Type:         enums.AccountTypeAsset,
				AccountLevel: 1,
			}, nil
		},
	}

	svc := account.NewService(repo)

	req := account.CreateAccountRequest{
		Code:     "2100",
		Name:     "LIABILITAS",
		ParentID: &parentID,
		Type:     enums.AccountTypeLiability, // Mismatch with parent's ASSET
	}

	_, err := svc.CreateAccount(ctx, req, "admin")
	if !errors.Is(err, account.ErrTypeMismatchWithParent) {
		t.Fatalf("expected ErrTypeMismatchWithParent, got: %v", err)
	}
}

func TestService_CreateAccount_DuplicateCode(t *testing.T) {
	ctx := context.Background()
	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return &account.ChartOfAccount{Code: code}, nil
		},
	}

	svc := account.NewService(repo)
	req := account.CreateAccountRequest{
		Code: "1000",
		Name: "ASET",
		Type: enums.AccountTypeAsset,
	}

	_, err := svc.CreateAccount(ctx, req, "admin")
	if !errors.Is(err, account.ErrAccountCodeAlreadyExists) {
		t.Fatalf("expected ErrAccountCodeAlreadyExists, got: %v", err)
	}
}

func TestService_CreateAccount_TreasuryRules(t *testing.T) {
	ctx := context.Background()
	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	svc := account.NewService(repo)

	// Treasury with non-ASSET type should fail
	req := account.CreateAccountRequest{
		Code:              "4100",
		Name:              "Pendapatan",
		Type:              enums.AccountTypeRevenue,
		IsTreasuryAccount: true,
	}

	_, err := svc.CreateAccount(ctx, req, "admin")
	if !errors.Is(err, account.ErrTreasuryMustBeAsset) {
		t.Fatalf("expected ErrTreasuryMustBeAsset, got: %v", err)
	}
}

func TestService_DeleteAccount_HasChildren(t *testing.T) {
	ctx := context.Background()
	accountID := "uuid-header-1100"

	repo := &mockAccountRepo{
		findByIDFn: func(ctx context.Context, id string) (*account.ChartOfAccount, error) {
			return &account.ChartOfAccount{
				ID:   accountID,
				Code: "1100",
			}, nil
		},
		findChildrenCountFn: func(ctx context.Context, parentID string) (int64, error) {
			return 1, nil // Still has children!
		},
	}

	svc := account.NewService(repo)

	err := svc.DeleteAccount(ctx, accountID, "admin")
	if !errors.Is(err, account.ErrCannotDeleteAccountWithChildren) {
		t.Fatalf("expected ErrCannotDeleteAccountWithChildren, got: %v", err)
	}
}

func TestService_DeleteAccount_UsedInTariff(t *testing.T) {
	ctx := context.Background()
	accountID := "uuid-account-4101"

	repo := &mockAccountRepo{
		findByIDFn: func(ctx context.Context, id string) (*account.ChartOfAccount, error) {
			return &account.ChartOfAccount{
				ID:   accountID,
				Code: "4101",
			}, nil
		},
		findChildrenCountFn: func(ctx context.Context, parentID string) (int64, error) {
			return 0, nil
		},
		isCodeUsedInTariffFn: func(ctx context.Context, code string) (bool, error) {
			return true, nil // Used in tariff!
		},
	}

	svc := account.NewService(repo)

	err := svc.DeleteAccount(ctx, accountID, "admin")
	if !errors.Is(err, account.ErrCannotDeleteAccountInUse) {
		t.Fatalf("expected ErrCannotDeleteAccountInUse, got: %v", err)
	}
}

func TestService_GetTree(t *testing.T) {
	ctx := context.Background()
	parentID := "id-1000"
	childID := "id-1100"
	grandChildID := "id-1110"

	accounts := []account.ChartOfAccount{
		{
			ID:           parentID,
			Code:         "1000",
			Name:         "ASET",
			Type:         enums.AccountTypeAsset,
			AccountLevel: 1,
			IsPostable:   false,
		},
		{
			ID:           childID,
			ParentID:     &parentID,
			Code:         "1100",
			Name:         "ASET LANCAR",
			Type:         enums.AccountTypeAsset,
			AccountLevel: 2,
			IsPostable:   false,
		},
		{
			ID:           grandChildID,
			ParentID:     &childID,
			Code:         "1110",
			Name:         "Kas dan Setara Kas",
			Type:         enums.AccountTypeAsset,
			AccountLevel: 3,
			IsPostable:   true,
		},
	}

	repo := &mockAccountRepo{
		findAllRawFn: func(ctx context.Context, activeOnly bool) ([]account.ChartOfAccount, error) {
			return accounts, nil
		},
	}

	svc := account.NewService(repo)

	tree, err := svc.GetTree(ctx, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tree) != 1 {
		t.Fatalf("expected 1 root node, got: %d", len(tree))
	}

	root := tree[0]
	if root.Code != "1000" {
		t.Errorf("expected root code 1000, got: %s", root.Code)
	}
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child for root, got: %d", len(root.Children))
	}

	child := root.Children[0]
	if child.Code != "1100" {
		t.Errorf("expected child code 1100, got: %s", child.Code)
	}
	if len(child.Children) != 1 {
		t.Fatalf("expected 1 child for 1100, got: %d", len(child.Children))
	}

	grandChild := child.Children[0]
	if grandChild.Code != "1110" {
		t.Errorf("expected grandchild code 1110, got: %s", grandChild.Code)
	}
	if len(grandChild.Children) != 0 {
		t.Errorf("expected leaf to have 0 children, got: %d", len(grandChild.Children))
	}
}

func TestService_CreateAccount_ExplicitHeadAccount(t *testing.T) {
	ctx := context.Background()
	repo := &mockAccountRepo{
		findByCodeFn: func(ctx context.Context, code string) (*account.ChartOfAccount, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(ctx context.Context, acc *account.ChartOfAccount) error {
			acc.ID = "uuid-head-6000"
			acc.CreatedAt = time.Now()
			acc.UpdatedAt = time.Now()
			return nil
		},
	}

	svc := account.NewService(repo)
	isPostableFalse := false

	req := account.CreateAccountRequest{
		Code:       "6000",
		Name:       "PENDAPATAN NON OPERASIONAL",
		Type:       enums.AccountTypeRevenue,
		IsPostable: &isPostableFalse, // Dibuat sebagai akun kepala/header
	}

	res, err := svc.CreateAccount(ctx, req, "admin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res.IsPostable {
		t.Errorf("expected is_postable=false for explicit head account, got: %v", res.IsPostable)
	}
}
