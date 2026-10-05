package enums

import "testing"

func TestAccountType_IsValid(t *testing.T) {
	validTypes := []AccountType{
		AccountTypeAsset,
		AccountTypeLiability,
		AccountTypeEquity,
		AccountTypeRevenue,
		AccountTypeExpense,
	}

	for _, vt := range validTypes {
		if !vt.IsValid() {
			t.Errorf("expected %s to be valid", vt)
		}
		if vt.String() != string(vt) {
			t.Errorf("expected string %s, got %s", vt, vt.String())
		}
	}

	invalid := AccountType("INVALID")
	if invalid.IsValid() {
		t.Errorf("expected INVALID to be invalid")
	}
}

func TestAccountPosition_IsValid(t *testing.T) {
	if !AccountPositionDebit.IsValid() {
		t.Errorf("expected DEBIT to be valid")
	}
	if !AccountPositionCredit.IsValid() {
		t.Errorf("expected CREDIT to be valid")
	}

	invalid := AccountPosition("SIDE")
	if invalid.IsValid() {
		t.Errorf("expected SIDE to be invalid")
	}
}

func TestDefaultPositionForAccountType(t *testing.T) {
	if DefaultPositionForAccountType(AccountTypeAsset) != AccountPositionDebit {
		t.Errorf("expected ASSET to default to DEBIT")
	}
	if DefaultPositionForAccountType(AccountTypeExpense) != AccountPositionDebit {
		t.Errorf("expected EXPENSE to default to DEBIT")
	}
	if DefaultPositionForAccountType(AccountTypeLiability) != AccountPositionCredit {
		t.Errorf("expected LIABILITY to default to CREDIT")
	}
	if DefaultPositionForAccountType(AccountTypeEquity) != AccountPositionCredit {
		t.Errorf("expected EQUITY to default to CREDIT")
	}
	if DefaultPositionForAccountType(AccountTypeRevenue) != AccountPositionCredit {
		t.Errorf("expected REVENUE to default to CREDIT")
	}
}
