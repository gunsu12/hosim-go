package account

import (
	"time"

	"hosim-go/pkg/enums"
)

// ==========================================
// Request DTOs
// ==========================================

// CreateAccountRequest mendefinisikan payload pembuatan akun baru
type CreateAccountRequest struct {
	Code              string                 `json:"code" binding:"required"`
	Name              string                 `json:"name" binding:"required"`
	ParentID          *string                `json:"parent_id"`
	Type              enums.AccountType      `json:"type" binding:"required"`
	Position          *enums.AccountPosition `json:"position"` // Opsional, default dihitung otomatis dari Type
	Description       *string                `json:"description"`
	IsPostable        *bool                  `json:"is_postable"` // Opsional, default true (false untuk akun kepala/header)
	IsTreasuryAccount bool                   `json:"is_treasury_account"`
	BankName          *string                `json:"bank_name"`
	BankAccountNumber *string                `json:"bank_account_number"`
	Currency          string                 `json:"currency"` // Default: IDR
	IsActive          *bool                  `json:"is_active"`
}

// UpdateAccountRequest mendefinisikan payload pembaruan metadata akun
type UpdateAccountRequest struct {
	Name              string                 `json:"name" binding:"required"`
	Position          *enums.AccountPosition `json:"position"`
	Description       *string                `json:"description"`
	IsTreasuryAccount bool                   `json:"is_treasury_account"`
	BankName          *string                `json:"bank_name"`
	BankAccountNumber *string                `json:"bank_account_number"`
	Currency          string                 `json:"currency"`
	IsActive          *bool                  `json:"is_active"`
}

// AccountListParams mendefinisikan parameter filter dan paginasi daftar akun
type AccountListParams struct {
	Page       int
	Limit      int
	Search     string
	Type       *enums.AccountType
	Position   *enums.AccountPosition
	ParentID   *string
	IsPostable *bool
	IsTreasury *bool
	IsActive   *bool
	Level      *int
}

// ==========================================
// Response DTOs
// ==========================================

// AccountResponse merepresentasikan data akun untuk output REST API
type AccountResponse struct {
	ID                string                `json:"id"`
	Code              string                `json:"code"`
	ParentID          *string               `json:"parent_id,omitempty"`
	ParentCode        *string               `json:"parent_code,omitempty"`
	Name              string                `json:"name"`
	Description       *string               `json:"description,omitempty"`
	Type              enums.AccountType     `json:"type"`
	Position          enums.AccountPosition `json:"position"`
	AccountLevel      int                   `json:"account_level"`
	IsPostable        bool                  `json:"is_postable"`
	IsTreasuryAccount bool                  `json:"is_treasury_account"`
	BankName          *string               `json:"bank_name,omitempty"`
	BankAccountNumber *string               `json:"bank_account_number,omitempty"`
	Currency          string                `json:"currency"`
	IsActive          bool                  `json:"is_active"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

// AccountTreeNode merepresentasikan simpul hierarki pohon akun
type AccountTreeNode struct {
	ID                string                `json:"id"`
	Code              string                `json:"code"`
	ParentID          *string               `json:"parent_id,omitempty"`
	ParentCode        *string               `json:"parent_code,omitempty"`
	Name              string                `json:"name"`
	Description       *string               `json:"description,omitempty"`
	Type              enums.AccountType     `json:"type"`
	Position          enums.AccountPosition `json:"position"`
	AccountLevel      int                   `json:"account_level"`
	IsPostable        bool                  `json:"is_postable"`
	IsTreasuryAccount bool                  `json:"is_treasury_account"`
	BankName          *string               `json:"bank_name,omitempty"`
	BankAccountNumber *string               `json:"bank_account_number,omitempty"`
	Currency          string                `json:"currency"`
	IsActive          bool                  `json:"is_active"`
	Children          []AccountTreeNode     `json:"children"`
}
