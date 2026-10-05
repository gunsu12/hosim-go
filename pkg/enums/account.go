package enums

// AccountType mendefinisikan 5 klasifikasi baku tipe akun akuntansi
type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
)

func (t AccountType) IsValid() bool {
	switch t {
	case AccountTypeAsset, AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue, AccountTypeExpense:
		return true
	default:
		return false
	}
}

func (t AccountType) String() string {
	return string(t)
}

// AccountPosition mendefinisikan posisi saldo normal akun (Debit / Kredit)
type AccountPosition string

const (
	AccountPositionDebit  AccountPosition = "DEBIT"
	AccountPositionCredit AccountPosition = "CREDIT"
)

func (p AccountPosition) IsValid() bool {
	switch p {
	case AccountPositionDebit, AccountPositionCredit:
		return true
	default:
		return false
	}
}

func (p AccountPosition) String() string {
	return string(p)
}

// DefaultPositionForAccountType mengembalikan saldo normal default berdasarkan persamaan dasar akuntansi
func DefaultPositionForAccountType(t AccountType) AccountPosition {
	switch t {
	case AccountTypeAsset, AccountTypeExpense:
		return AccountPositionDebit
	case AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue:
		return AccountPositionCredit
	default:
		return AccountPositionDebit
	}
}
