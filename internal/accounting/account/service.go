package account

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Service mengorkestrasi logika bisnis dan validasi hierarki bagan akun
type Service interface {
	CreateAccount(ctx context.Context, req CreateAccountRequest, operator string) (*AccountResponse, error)
	UpdateAccount(ctx context.Context, id string, req UpdateAccountRequest, operator string) (*AccountResponse, error)
	DeleteAccount(ctx context.Context, id string, operator string) error
	GetByID(ctx context.Context, id string) (*AccountResponse, error)
	GetByCode(ctx context.Context, code string) (*AccountResponse, error)
	List(ctx context.Context, params AccountListParams) ([]AccountResponse, int64, error)
	GetTree(ctx context.Context, activeOnly bool) ([]AccountTreeNode, error)
	ListPostable(ctx context.Context) ([]AccountResponse, error)
	ListTreasury(ctx context.Context) ([]AccountResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAccount(ctx context.Context, req CreateAccountRequest, operator string) (*AccountResponse, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	if req.Code == "" {
		return nil, ErrAccountCodeRequired
	}
	if req.Name == "" {
		return nil, ErrAccountNameRequired
	}
	if !IsValidType(req.Type) {
		return nil, ErrInvalidAccountType
	}

	// 1. Cek keunikan kode akun
	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrAccountCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. Evaluasi hierarki induk (Parent)
	var parent *ChartOfAccount
	accountLevel := 1
	var parentCode *string

	if req.ParentID != nil && strings.TrimSpace(*req.ParentID) != "" {
		pID := strings.TrimSpace(*req.ParentID)
		parent, err = s.repo.FindByID(ctx, pID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrParentAccountNotFound
			}
			return nil, err
		}

		// Invariant: Tipe akun anak wajib sama dengan tipe akun induk
		if req.Type != parent.Type {
			return nil, ErrTypeMismatchWithParent
		}

		accountLevel = parent.AccountLevel + 1
		pCode := parent.Code
		parentCode = &pCode
		req.ParentID = &pID
	} else {
		req.ParentID = nil
	}

	// 3. Tentukan posisi saldo normal (default atau override)
	position := DefaultPositionForType(req.Type)
	if req.Position != nil && *req.Position != "" {
		if !IsValidPosition(*req.Position) {
			return nil, ErrInvalidAccountPosition
		}
		position = *req.Position
	}

	// 4. Invariant Akun Perbendaharaan (Treasury)
	if req.IsTreasuryAccount {
		if req.Type != AccountTypeAsset {
			return nil, ErrTreasuryMustBeAsset
		}
	}

	currency := "IDR"
	if strings.TrimSpace(req.Currency) != "" {
		currency = strings.TrimSpace(req.Currency)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isPostable := true
	if req.IsPostable != nil {
		isPostable = *req.IsPostable
	}

	acc := &ChartOfAccount{
		ParentID:          req.ParentID,
		ParentCode:        parentCode,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		Type:              req.Type,
		Position:          position,
		AccountLevel:      accountLevel,
		IsPostable:        isPostable,
		IsTreasuryAccount: req.IsTreasuryAccount,
		BankName:          req.BankName,
		BankAccountNumber: req.BankAccountNumber,
		Currency:          currency,
		IsActive:          isActive,
		CreatedBy:         operator,
		UpdatedBy:         operator,
	}

	if err := acc.Validate(); err != nil {
		return nil, err
	}

	// 5. Simpan dan perbarui status induk secara atomik bila ada induk
	if parent != nil {
		if err := s.repo.CreateWithParentUpdate(ctx, acc, parent.ID); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Create(ctx, acc); err != nil {
			return nil, err
		}
	}

	return toAccountResponse(acc), nil
}

func (s *service) UpdateAccount(ctx context.Context, id string, req UpdateAccountRequest, operator string) (*AccountResponse, error) {
	acc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrAccountNameRequired
	}

	if req.Position != nil && *req.Position != "" {
		if !IsValidPosition(*req.Position) {
			return nil, ErrInvalidAccountPosition
		}
		acc.Position = *req.Position
	}

	// Invariant Akun Perbendaharaan
	if req.IsTreasuryAccount {
		if acc.Type != AccountTypeAsset {
			return nil, ErrTreasuryMustBeAsset
		}
		if !acc.IsPostable {
			return nil, ErrTreasuryMustBePostable
		}
	}

	acc.Name = req.Name
	acc.Description = req.Description
	acc.IsTreasuryAccount = req.IsTreasuryAccount
	acc.BankName = req.BankName
	acc.BankAccountNumber = req.BankAccountNumber
	if strings.TrimSpace(req.Currency) != "" {
		acc.Currency = strings.TrimSpace(req.Currency)
	}
	if req.IsActive != nil {
		acc.IsActive = *req.IsActive
	}
	acc.UpdatedBy = operator

	if err := acc.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, acc); err != nil {
		return nil, err
	}

	return toAccountResponse(acc), nil
}

func (s *service) DeleteAccount(ctx context.Context, id string, operator string) error {
	acc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAccountNotFound
		}
		return err
	}

	// 1. Cek apakah akun memiliki anak aktif
	childrenCount, err := s.repo.FindChildrenCount(ctx, id)
	if err != nil {
		return err
	}
	if childrenCount > 0 {
		return ErrCannotDeleteAccountWithChildren
	}

	// 2. Cek apakah akun digunakan dalam tarif tindakan atau buku tarif
	inUse, err := s.repo.IsCodeUsedInTariff(ctx, acc.Code)
	if err != nil {
		return err
	}
	if inUse {
		return ErrCannotDeleteAccountInUse
	}

	// 3. Lakukan soft delete dan periksa apakah induknya perlu dipulihkan jadi postable
	return s.repo.DeleteWithParentCheck(ctx, id, acc.ParentID, operator)
}

func (s *service) GetByID(ctx context.Context, id string) (*AccountResponse, error) {
	acc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return toAccountResponse(acc), nil
}

func (s *service) GetByCode(ctx context.Context, code string) (*AccountResponse, error) {
	acc, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return toAccountResponse(acc), nil
}

func (s *service) List(ctx context.Context, params AccountListParams) ([]AccountResponse, int64, error) {
	accounts, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]AccountResponse, len(accounts))
	for i := range accounts {
		responses[i] = *toAccountResponse(&accounts[i])
	}

	return responses, total, nil
}

func (s *service) GetTree(ctx context.Context, activeOnly bool) ([]AccountTreeNode, error) {
	accounts, err := s.repo.FindAllRaw(ctx, activeOnly)
	if err != nil {
		return nil, err
	}

	return buildTree(accounts), nil
}

func (s *service) ListPostable(ctx context.Context) ([]AccountResponse, error) {
	isPostable := true
	isActive := true
	params := AccountListParams{
		Limit:      500,
		IsPostable: &isPostable,
		IsActive:   &isActive,
	}

	accounts, _, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]AccountResponse, len(accounts))
	for i := range accounts {
		responses[i] = *toAccountResponse(&accounts[i])
	}
	return responses, nil
}

func (s *service) ListTreasury(ctx context.Context) ([]AccountResponse, error) {
	isTreasury := true
	isActive := true
	params := AccountListParams{
		Limit:      500,
		IsTreasury: &isTreasury,
		IsActive:   &isActive,
	}

	accounts, _, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]AccountResponse, len(accounts))
	for i := range accounts {
		responses[i] = *toAccountResponse(&accounts[i])
	}
	return responses, nil
}

// -------------------------------------------------------------
// Helper Functions: Tree Builder & DTO Mapper
// -------------------------------------------------------------

func buildTree(accounts []ChartOfAccount) []AccountTreeNode {
	childrenMap := make(map[string][]ChartOfAccount)
	var rootAccounts []ChartOfAccount

	for _, a := range accounts {
		if a.ParentID == nil || *a.ParentID == "" {
			rootAccounts = append(rootAccounts, a)
		} else {
			childrenMap[*a.ParentID] = append(childrenMap[*a.ParentID], a)
		}
	}

	var convert func(a ChartOfAccount) AccountTreeNode
	convert = func(a ChartOfAccount) AccountTreeNode {
		childAccounts := childrenMap[a.ID]
		children := make([]AccountTreeNode, len(childAccounts))
		for i, ca := range childAccounts {
			children[i] = convert(ca)
		}

		return AccountTreeNode{
			ID:                a.ID,
			Code:              a.Code,
			ParentID:          a.ParentID,
			ParentCode:        a.ParentCode,
			Name:              a.Name,
			Description:       a.Description,
			Type:              a.Type,
			Position:          a.Position,
			AccountLevel:      a.AccountLevel,
			IsPostable:        a.IsPostable,
			IsTreasuryAccount: a.IsTreasuryAccount,
			BankName:          a.BankName,
			BankAccountNumber: a.BankAccountNumber,
			Currency:          a.Currency,
			IsActive:          a.IsActive,
			Children:          children,
		}
	}

	result := make([]AccountTreeNode, len(rootAccounts))
	for i, ra := range rootAccounts {
		result[i] = convert(ra)
	}
	return result
}

func toAccountResponse(a *ChartOfAccount) *AccountResponse {
	return &AccountResponse{
		ID:                a.ID,
		Code:              a.Code,
		ParentID:          a.ParentID,
		ParentCode:        a.ParentCode,
		Name:              a.Name,
		Description:       a.Description,
		Type:              a.Type,
		Position:          a.Position,
		AccountLevel:      a.AccountLevel,
		IsPostable:        a.IsPostable,
		IsTreasuryAccount: a.IsTreasuryAccount,
		BankName:          a.BankName,
		BankAccountNumber: a.BankAccountNumber,
		Currency:          a.Currency,
		IsActive:          a.IsActive,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}
