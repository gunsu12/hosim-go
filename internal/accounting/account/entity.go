package account

import (
	"errors"
	"strings"
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

// Re-export tipe dan posisi akun dari pkg/enums untuk akses praktis di domain account
type AccountType = enums.AccountType

const (
	AccountTypeAsset     = enums.AccountTypeAsset
	AccountTypeLiability = enums.AccountTypeLiability
	AccountTypeEquity    = enums.AccountTypeEquity
	AccountTypeRevenue   = enums.AccountTypeRevenue
	AccountTypeExpense   = enums.AccountTypeExpense
)

// Re-export posisi saldo normal dari pkg/enums
type AccountPosition = enums.AccountPosition

const (
	AccountPositionDebit  = enums.AccountPositionDebit
	AccountPositionCredit = enums.AccountPositionCredit
)

var (
	ErrAccountNotFound                 = errors.New("akun tidak ditemukan")
	ErrAccountCodeRequired             = errors.New("kode akun wajib diisi")
	ErrAccountNameRequired             = errors.New("nama akun wajib diisi")
	ErrInvalidAccountLevel             = errors.New("level akun minimal 1")
	ErrAccountCodeAlreadyExists        = errors.New("kode akun sudah terdaftar")
	ErrInvalidAccountType              = errors.New("tipe akun tidak valid, harus ASSET, LIABILITY, EQUITY, REVENUE, atau EXPENSE")
	ErrInvalidAccountPosition          = errors.New("posisi saldo normal akun tidak valid, harus DEBIT atau CREDIT")
	ErrParentAccountNotFound           = errors.New("akun induk tidak ditemukan")
	ErrTypeMismatchWithParent          = errors.New("tipe akun harus sama dengan tipe akun induk")
	ErrParentCannotBeSelf              = errors.New("akun tidak dapat menjadi induk bagi dirinya sendiri")
	ErrCannotDeleteAccountWithChildren = errors.New("akun tidak dapat dihapus karena masih memiliki akun turunan")
	ErrCannotDeleteAccountInUse        = errors.New("akun tidak dapat dihapus karena sudah digunakan dalam transaksi atau konfigurasi tarif")
	ErrTreasuryMustBeAsset             = errors.New("akun perbendaharaan/kas-bank harus bertipe ASSET")
	ErrTreasuryMustBePostable          = errors.New("akun perbendaharaan/kas-bank harus berstatus postable (bukan header)")
	ErrNonPostableAccount              = errors.New("akun header/kategori tidak dapat menerima transaksi langsung (is_postable = false)")
	ErrAccountInactive                 = errors.New("akun dalam status non-aktif")
)

// ChartOfAccount merepresentasikan entitas bagan akun rumah sakit
type ChartOfAccount struct {
	ID                string          `gorm:"size:36;primaryKey" json:"id"`
	ParentID          *string         `gorm:"size:36;index" json:"parent_id,omitempty"`
	ParentCode        *string         `gorm:"size:50" json:"parent_code,omitempty"`
	Code              string          `gorm:"size:50;not null" json:"code"`
	Name              string          `gorm:"size:150;not null" json:"name"`
	Description       *string         `gorm:"type:text" json:"description,omitempty"`
	Type              AccountType     `gorm:"size:30;not null;index" json:"type"`
	Position          AccountPosition `gorm:"size:10;not null" json:"position"`
	AccountLevel      int             `gorm:"default:1;not null" json:"account_level"`
	IsPostable        bool            `gorm:"not null;index" json:"is_postable"`
	IsTreasuryAccount bool            `gorm:"default:false;not null;index" json:"is_treasury_account"`
	BankName          *string         `gorm:"size:100" json:"bank_name,omitempty"`
	BankAccountNumber *string         `gorm:"size:50" json:"bank_account_number,omitempty"`
	Currency          string          `gorm:"size:10;default:'IDR';not null" json:"currency"`
	IsActive          bool            `gorm:"default:true;not null;index" json:"is_active"`
	CreatedAt         time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	CreatedBy         string          `gorm:"default:'SYSTEM'" json:"created_by"`
	UpdatedBy         string          `gorm:"default:'SYSTEM'" json:"updated_by"`
	DeletedBy         *string         `json:"deleted_by,omitempty"`

	// Relasi hierarki pohon (Self-Referencing)
	Parent   *ChartOfAccount  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []ChartOfAccount `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (ChartOfAccount) TableName() string {
	return "chart_of_accounts"
}

func (a *ChartOfAccount) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewV7().String()
	}
	if a.Currency == "" {
		a.Currency = "IDR"
	}
	if a.AccountLevel <= 0 {
		a.AccountLevel = 1
	}
	if a.Position == "" && a.Type != "" {
		a.Position = DefaultPositionForType(a.Type)
	}
	return nil
}

// DefaultPositionForType mengembalikan saldo normal default berdasarkan persamaan akuntansi
func DefaultPositionForType(t AccountType) AccountPosition {
	return enums.DefaultPositionForAccountType(t)
}

// IsValidType memeriksa keabsahan klasifikasi tipe akun
func IsValidType(t AccountType) bool {
	return t.IsValid()
}

// IsValidPosition memeriksa keabsahan posisi saldo normal
func IsValidPosition(p AccountPosition) bool {
	return p.IsValid()
}

// Validate memvalidasi aturan dan invariant entitas akun
func (a *ChartOfAccount) Validate() error {
	a.Code = strings.TrimSpace(a.Code)
	a.Name = strings.TrimSpace(a.Name)

	if a.Code == "" {
		return ErrAccountCodeRequired
	}
	if a.Name == "" {
		return ErrAccountNameRequired
	}
	if !IsValidType(a.Type) {
		return ErrInvalidAccountType
	}
	if !IsValidPosition(a.Position) {
		return ErrInvalidAccountPosition
	}
	if a.AccountLevel < 1 {
		return ErrInvalidAccountLevel
	}
	if a.ParentID != nil && a.ID != "" && *a.ParentID == a.ID {
		return ErrParentCannotBeSelf
	}

	// Invariant Akun Perbendaharaan (Treasury)
	if a.IsTreasuryAccount {
		if a.Type != AccountTypeAsset {
			return ErrTreasuryMustBeAsset
		}
		if !a.IsPostable {
			return ErrTreasuryMustBePostable
		}
	}

	return nil
}
