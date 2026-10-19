// Fiscal Period & Fiscal Year Types for HOSIM Accounting Domain
// Matches backend definitions in internal/accounting/fiscalperiod/dto.go & entity.go

export type PeriodStatus = 'OPEN' | 'LOCKED' | 'CLOSED';

export interface FiscalPeriodRecord {
  id: string;
  fiscal_year_id: string;
  period_number: number;
  period_key: string; // e.g. "2026-10"
  name: string; // e.g. "Oktober 2026"
  start_date: string; // "YYYY-MM-DD"
  end_date: string; // "YYYY-MM-DD"
  status: PeriodStatus;
  locked_at?: string | null;
  locked_by?: string | null;
  closed_at?: string | null;
  closed_by?: string | null;
  notes?: string | null;
  created_at: string;
  updated_at: string;
}

export interface FiscalYearRecord {
  id: string;
  year: number;
  start_date: string; // "YYYY-MM-DD"
  end_date: string; // "YYYY-MM-DD"
  is_closed: boolean;
  closed_at?: string | null;
  closed_by?: string | null;
  notes?: string | null;
  periods_count: number;
  periods?: FiscalPeriodRecord[];
  created_at: string;
  updated_at: string;
}

export interface CreateFiscalYearDTO {
  year: number;
  notes?: string | null;
}

export interface FiscalPeriodActionDTO {
  notes?: string | null;
}

export interface FiscalPeriodFilterParams {
  year?: number;
  status?: string;
}
