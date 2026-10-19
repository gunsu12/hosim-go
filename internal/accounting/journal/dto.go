package journal

import (
	"time"
)

// JournalLineInput adalah payload input untuk satu baris jurnal.
type JournalLineInput struct {
	AccountID     string        `json:"account_id" binding:"required"`
	EntryPosition EntryPosition `json:"entry_position" binding:"required"`
	Amount        float64       `json:"amount" binding:"required,gt=0"`
	Memo          *string       `json:"memo,omitempty"`
	DepartmentID  *string       `json:"department_id,omitempty"`
}

// CreateJournalRequest adalah payload untuk membuat draf jurnal manual baru.
type CreateJournalRequest struct {
	EntryDate       string             `json:"entry_date" binding:"required"` // Format: YYYY-MM-DD
	EntryType       JournalType        `json:"entry_type" binding:"required"`
	Description     string             `json:"description" binding:"required"`
	ReferenceNumber *string            `json:"reference_number,omitempty"`
	Lines           []JournalLineInput `json:"lines" binding:"required,min=2"`
}

// UpdateJournalRequest adalah payload untuk memperbarui isi draf jurnal.
type UpdateJournalRequest struct {
	EntryDate       *string            `json:"entry_date,omitempty"`
	Description     *string            `json:"description,omitempty"`
	ReferenceNumber *string            `json:"reference_number,omitempty"`
	Lines           []JournalLineInput `json:"lines,omitempty"`
}

// ReverseJournalRequest adalah payload saat melakukan pembatalan jurnal berstatus POSTED.
type ReverseJournalRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// JournalLineResponse merepresentasikan rincian baris jurnal untuk client API.
type JournalLineResponse struct {
	ID            string        `json:"id"`
	LineNumber    int           `json:"line_number"`
	AccountID     string        `json:"account_id"`
	AccountCode   string        `json:"account_code"`
	AccountName   string        `json:"account_name"`
	EntryPosition EntryPosition `json:"entry_position"`
	Amount        float64       `json:"amount"`
	Memo          *string       `json:"memo,omitempty"`
	DepartmentID  *string       `json:"department_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

func ToJournalLineResponse(l *JournalEntryLine) JournalLineResponse {
	amt, _ := l.Amount.Float64()
	return JournalLineResponse{
		ID:            l.ID,
		LineNumber:    l.LineNumber,
		AccountID:     l.AccountID,
		AccountCode:   l.AccountCode,
		AccountName:   l.AccountName,
		EntryPosition: l.EntryPosition,
		Amount:        amt,
		Memo:          l.Memo,
		DepartmentID:  l.DepartmentID,
		CreatedAt:     l.CreatedAt,
	}
}

// JournalResponse merepresentasikan output voucher jurnal untuk client API.
type JournalResponse struct {
	ID              string                `json:"id"`
	EntryNumber     string                `json:"entry_number"`
	EntryDate       string                `json:"entry_date"`
	FiscalPeriodID  string                `json:"fiscal_period_id"`
	FiscalPeriod    string                `json:"fiscal_period"`
	EntryType       JournalType           `json:"entry_type"`
	Status          JournalStatus         `json:"status"`
	Description     string                `json:"description"`
	ReferenceNumber *string               `json:"reference_number,omitempty"`
	TotalDebit      float64               `json:"total_debit"`
	TotalCredit     float64               `json:"total_credit"`
	LinesCount      int                   `json:"lines_count"`
	ReversedEntryID *string               `json:"reversed_entry_id,omitempty"`
	ReversalReason  *string               `json:"reversal_reason,omitempty"`
	PostedAt        *time.Time            `json:"posted_at,omitempty"`
	PostedBy        *string               `json:"posted_by,omitempty"`
	ReversedAt      *time.Time            `json:"reversed_at,omitempty"`
	ReversedBy      *string               `json:"reversed_by,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	CreatedBy       string                `json:"created_by"`
	UpdatedAt       time.Time             `json:"updated_at"`
	UpdatedBy       string                `json:"updated_by"`
	Lines           []JournalLineResponse `json:"lines,omitempty"`
}

func ToJournalResponse(e *JournalEntry) *JournalResponse {
	if e == nil {
		return nil
	}

	debit, _ := e.TotalDebit.Float64()
	credit, _ := e.TotalCredit.Float64()

	lines := make([]JournalLineResponse, len(e.Lines))
	for i, l := range e.Lines {
		lines[i] = ToJournalLineResponse(l)
	}

	return &JournalResponse{
		ID:              e.ID,
		EntryNumber:     e.EntryNumber,
		EntryDate:       e.EntryDate.Format("2006-01-02"),
		FiscalPeriodID:  e.FiscalPeriodID,
		FiscalPeriod:    e.FiscalPeriod,
		EntryType:       e.EntryType,
		Status:          e.Status,
		Description:     e.Description,
		ReferenceNumber: e.ReferenceNumber,
		TotalDebit:      debit,
		TotalCredit:     credit,
		LinesCount:      len(e.Lines),
		ReversedEntryID: e.ReversedEntryID,
		ReversalReason:  e.ReversalReason,
		PostedAt:        e.PostedAt,
		PostedBy:        e.PostedBy,
		ReversedAt:      e.ReversedAt,
		ReversedBy:      e.ReversedBy,
		CreatedAt:       e.CreatedAt,
		CreatedBy:       e.CreatedBy,
		UpdatedAt:       e.UpdatedAt,
		UpdatedBy:       e.UpdatedBy,
		Lines:           lines,
	}
}

func ToJournalResponseList(entries []*JournalEntry) []*JournalResponse {
	res := make([]*JournalResponse, len(entries))
	for i, e := range entries {
		res[i] = ToJournalResponse(e)
	}
	return res
}
