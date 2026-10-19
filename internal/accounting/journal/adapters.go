package journal

import (
	"context"
	"fmt"

	"hosim-go/internal/accounting/account"
)

// AccountServiceAdapter mengadaptasi account.Service agar memenuhi interface AccountVerifier.
type AccountServiceAdapter struct {
	svc account.Service
}

func NewAccountServiceAdapter(svc account.Service) AccountVerifier {
	return &AccountServiceAdapter{svc: svc}
}

func (a *AccountServiceAdapter) ValidatePostableAccounts(ctx context.Context, accountIDs []string) (map[string]*AccountSummary, error) {
	summaries := make(map[string]*AccountSummary, len(accountIDs))

	for _, id := range accountIDs {
		// Hindari lookup ganda jika akun yang sama muncul di baris berbeda
		if _, exists := summaries[id]; exists {
			continue
		}

		acc, err := a.svc.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w: akun ID '%s' tidak ditemukan", ErrAccountNotPostable, id)
		}

		if !acc.IsActive {
			return nil, fmt.Errorf("%w: akun %s (%s) berstatus non-aktif", ErrAccountNotPostable, acc.Code, acc.Name)
		}

		if !acc.IsPostable {
			return nil, fmt.Errorf("%w: akun %s (%s) adalah akun header induk dan tidak boleh dijurnal", ErrAccountNotPostable, acc.Code, acc.Name)
		}

		summaries[id] = &AccountSummary{
			ID:         acc.ID,
			Code:       acc.Code,
			Name:       acc.Name,
			IsActive:   acc.IsActive,
			IsPostable: acc.IsPostable,
		}
	}

	return summaries, nil
}
