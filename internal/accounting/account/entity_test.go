package account

import (
	"errors"
	"testing"
)

func TestDefaultPositionForType(t *testing.T) {
	tests := []struct {
		accountType AccountType
		expected    AccountPosition
	}{
		{AccountTypeAsset, AccountPositionDebit},
		{AccountTypeExpense, AccountPositionDebit},
		{AccountTypeLiability, AccountPositionCredit},
		{AccountTypeEquity, AccountPositionCredit},
		{AccountTypeRevenue, AccountPositionCredit},
	}

	for _, tt := range tests {
		got := DefaultPositionForType(tt.accountType)
		if got != tt.expected {
			t.Errorf("DefaultPositionForType(%s) = %s, want %s", tt.accountType, got, tt.expected)
		}
	}
}

func TestChartOfAccount_Validate_Success(t *testing.T) {
	acc := &ChartOfAccount{
		Code:         "1111",
		Name:         "Kas Kasir Utama",
		Type:         AccountTypeAsset,
		Position:     AccountPositionDebit,
		AccountLevel: 4,
		IsPostable:   true,
	}

	if err := acc.Validate(); err != nil {
		t.Fatalf("expected valid account, got error: %v", err)
	}
}

func TestChartOfAccount_Validate_ContraAccount(t *testing.T) {
	// Akun Kontra: Akumulasi Penyusutan (Tipe ASSET, Posisi CREDIT)
	acc := &ChartOfAccount{
		Code:         "1219",
		Name:         "Akumulasi Penyusutan Alat Medis",
		Type:         AccountTypeAsset,
		Position:     AccountPositionCredit, // Valid contra position
		AccountLevel: 3,
		IsPostable:   true,
	}

	if err := acc.Validate(); err != nil {
		t.Fatalf("expected valid contra account, got error: %v", err)
	}
}

func TestChartOfAccount_Validate_TreasuryInvariants(t *testing.T) {
	// Kasus 1: Treasury tapi bukan ASSET -> Gagal
	invalidTypeTreasury := &ChartOfAccount{
		Code:              "2111",
		Name:              "Hutang Bank Operasional",
		Type:              AccountTypeLiability,
		Position:          AccountPositionCredit,
		AccountLevel:      3,
		IsPostable:        true,
		IsTreasuryAccount: true,
	}
	if err := invalidTypeTreasury.Validate(); !errors.Is(err, ErrTreasuryMustBeAsset) {
		t.Errorf("expected ErrTreasuryMustBeAsset, got %v", err)
	}

	// Kasus 2: Treasury tapi bukan Postable (Header) -> Gagal
	nonPostableTreasury := &ChartOfAccount{
		Code:              "1110",
		Name:              "Kas dan Bank Header",
		Type:              AccountTypeAsset,
		Position:          AccountPositionDebit,
		AccountLevel:      2,
		IsPostable:        false,
		IsTreasuryAccount: true,
	}
	if err := nonPostableTreasury.Validate(); !errors.Is(err, ErrTreasuryMustBePostable) {
		t.Errorf("expected ErrTreasuryMustBePostable, got %v", err)
	}

	// Kasus 3: Treasury valid (ASSET dan Postable) -> Berhasil
	validTreasury := &ChartOfAccount{
		Code:              "1113",
		Name:              "Bank Mandiri Operasional",
		Type:              AccountTypeAsset,
		Position:          AccountPositionDebit,
		AccountLevel:      4,
		IsPostable:        true,
		IsTreasuryAccount: true,
	}
	if err := validTreasury.Validate(); err != nil {
		t.Errorf("expected valid treasury account, got error: %v", err)
	}
}

func TestChartOfAccount_Validate_ParentCannotBeSelf(t *testing.T) {
	selfID := "acc-001"
	acc := &ChartOfAccount{
		ID:           selfID,
		ParentID:     &selfID,
		Code:         "1000",
		Name:         "Aset",
		Type:         AccountTypeAsset,
		Position:     AccountPositionDebit,
		AccountLevel: 1,
		IsPostable:   false,
	}

	if err := acc.Validate(); !errors.Is(err, ErrParentCannotBeSelf) {
		t.Errorf("expected ErrParentCannotBeSelf, got %v", err)
	}
}

func TestChartOfAccount_BeforeCreate(t *testing.T) {
	acc := &ChartOfAccount{
		Code: "4101",
		Name: "Pendapatan Tindakan",
		Type: AccountTypeRevenue,
	}

	if err := acc.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}

	if acc.ID == "" {
		t.Errorf("expected ID to be generated via UUID v7")
	}
	if acc.Currency != "IDR" {
		t.Errorf("expected default currency IDR, got %s", acc.Currency)
	}
	if acc.Position != AccountPositionCredit {
		t.Errorf("expected default position CREDIT for REVENUE, got %s", acc.Position)
	}
	if acc.AccountLevel != 1 {
		t.Errorf("expected default account level 1, got %d", acc.AccountLevel)
	}
}
