package fiscalperiod

import (
	"testing"
)

func TestGenerateFiscalYearWithPeriods_NormalYear(t *testing.T) {
	year := 2026
	notes := "Tahun Buku 2026"
	actor := "admin"

	fy, periods, err := GenerateFiscalYearWithPeriods(year, &notes, actor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fy.Year != 2026 {
		t.Errorf("expected year 2026, got %d", fy.Year)
	}
	if len(periods) != 12 {
		t.Fatalf("expected 12 periods, got %d", len(periods))
	}

	// Cek Februari 2026 (bukan tahun kabisat -> 28 hari)
	feb := periods[1]
	if feb.PeriodKey != "2026-02" {
		t.Errorf("expected period key 2026-02, got %s", feb.PeriodKey)
	}
	if feb.Name != "Februari 2026" {
		t.Errorf("expected name Februari 2026, got %s", feb.Name)
	}
	if feb.EndDate.Day() != 28 {
		t.Errorf("expected Feb 2026 end date day 28, got %d", feb.EndDate.Day())
	}

	// Cek status awal semua OPEN
	for _, p := range periods {
		if p.Status != PeriodStatusOpen {
			t.Errorf("expected period %s status OPEN, got %s", p.PeriodKey, p.Status)
		}
	}
}

func TestGenerateFiscalYearWithPeriods_LeapYear(t *testing.T) {
	year := 2024 // Tahun kabisat -> Februari 29 hari
	fy, periods, err := GenerateFiscalYearWithPeriods(year, nil, "system")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fy.Year != 2024 {
		t.Errorf("expected year 2024, got %d", fy.Year)
	}

	feb := periods[1]
	if feb.EndDate.Day() != 29 {
		t.Errorf("expected Feb 2024 end date day 29, got %d", feb.EndDate.Day())
	}
}

func TestGenerateFiscalYearWithPeriods_InvalidYear(t *testing.T) {
	_, _, err := GenerateFiscalYearWithPeriods(1999, nil, "system")
	if err != ErrInvalidYear {
		t.Errorf("expected ErrInvalidYear, got %v", err)
	}

	_, _, err = GenerateFiscalYearWithPeriods(2101, nil, "system")
	if err != ErrInvalidYear {
		t.Errorf("expected ErrInvalidYear, got %v", err)
	}
}

func TestPeriod_StateTransitions(t *testing.T) {
	p := &FiscalPeriod{
		PeriodKey: "2026-10",
		Status:    PeriodStatusOpen,
	}

	// 1. Lock from OPEN
	notes := "Proses rekonsiliasi"
	err := p.Lock("manager", &notes)
	if err != nil {
		t.Fatalf("unexpected lock error: %v", err)
	}
	if p.Status != PeriodStatusLocked {
		t.Errorf("expected status LOCKED, got %s", p.Status)
	}
	if p.LockedBy == nil || *p.LockedBy != "manager" {
		t.Errorf("expected locked_by manager")
	}

	// 2. Unlock from LOCKED
	unlockNotes := "Buka revisi"
	err = p.Unlock("manager", &unlockNotes)
	if err != nil {
		t.Fatalf("unexpected unlock error: %v", err)
	}
	if p.Status != PeriodStatusOpen {
		t.Errorf("expected status OPEN, got %s", p.Status)
	}

	// 3. Close from OPEN
	closeNotes := "Tutup permanen"
	err = p.Close("director", &closeNotes)
	if err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if p.Status != PeriodStatusClosed {
		t.Errorf("expected status CLOSED, got %s", p.Status)
	}

	// 4. Cannot Lock or Unlock a CLOSED period
	if err := p.Lock("manager", nil); err != ErrCannotLockClosedPeriod {
		t.Errorf("expected ErrCannotLockClosedPeriod, got %v", err)
	}
	if err := p.Unlock("manager", nil); err != ErrPeriodAlreadyClosed {
		t.Errorf("expected ErrPeriodAlreadyClosed, got %v", err)
	}
}

func TestFiscalYear_CloseYear(t *testing.T) {
	fy, periods, _ := GenerateFiscalYearWithPeriods(2026, nil, "system")

	// 1. Fails when periods are still OPEN
	err := fy.CloseYear(periods, "director")
	if err == nil {
		t.Fatal("expected error closing year with open periods, got nil")
	}

	// 2. Close 11 months, month 12 still open -> must still fail
	for i := 0; i < 11; i++ {
		_ = periods[i].Close("director", nil)
	}
	err = fy.CloseYear(periods, "director")
	if err == nil {
		t.Fatal("expected error closing year with 1 open period, got nil")
	}

	// 3. Close all 12 months -> must succeed
	_ = periods[11].Close("director", nil)
	err = fy.CloseYear(periods, "director")
	if err != nil {
		t.Fatalf("unexpected error closing year with all closed periods: %v", err)
	}
	if !fy.IsClosed {
		t.Errorf("expected is_closed true")
	}
	if fy.ClosedBy == nil || *fy.ClosedBy != "director" {
		t.Errorf("expected closed_by director")
	}
}
