import { authFetch } from '../client';
import type {
  FiscalYearRecord,
  FiscalPeriodRecord,
  CreateFiscalYearDTO,
  FiscalPeriodActionDTO,
  FiscalPeriodFilterParams
} from '../../types/accounting/fiscalPeriod';

// ==========================================
// ACCOUNTING: FISCAL YEARS & PERIODS API
// ==========================================

/**
 * Mengambil daftar tahun fiskal yang terdaftar
 */
export async function listFiscalYears(): Promise<FiscalYearRecord[]> {
  const res = await authFetch<FiscalYearRecord[]>('/accounting/fiscal-years');
  return res.data;
}

/**
 * Mengambil detail satu tahun fiskal beserta 12 periode bulanannya
 */
export async function getFiscalYearById(id: string): Promise<FiscalYearRecord> {
  const res = await authFetch<FiscalYearRecord>(`/accounting/fiscal-years/${encodeURIComponent(id)}`);
  return res.data;
}

/**
 * Membuat tahun fiskal baru beserta 12 periode bulanannya
 */
export async function createFiscalYear(payload: CreateFiscalYearDTO): Promise<FiscalYearRecord> {
  const res = await authFetch<FiscalYearRecord>('/accounting/fiscal-years', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

/**
 * Menutup buku tahunan (seluruh 12 bulan harus sudah berstatus CLOSED)
 */
export async function closeFiscalYear(id: string): Promise<FiscalYearRecord> {
  const res = await authFetch<FiscalYearRecord>(`/accounting/fiscal-years/${encodeURIComponent(id)}/close`, {
    method: 'PUT',
  });
  return res.data;
}

/**
 * Mengambil daftar periode bulanan dengan filter tahun dan status
 */
export async function listFiscalPeriods(params?: FiscalPeriodFilterParams): Promise<FiscalPeriodRecord[]> {
  const q = new URLSearchParams();
  if (params?.year !== undefined && params?.year !== null) {
    q.set('year', String(params.year));
  }
  if (params?.status) {
    q.set('status', params.status);
  }
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<FiscalPeriodRecord[]>(`/accounting/fiscal-periods${queryStr}`);
  return res.data;
}

/**
 * Mengambil periode fiskal yang aktif untuk hari ini
 */
export async function getCurrentFiscalPeriod(): Promise<FiscalPeriodRecord> {
  const res = await authFetch<FiscalPeriodRecord>('/accounting/fiscal-periods/current');
  return res.data;
}

/**
 * Mengambil detail satu periode bulanan
 */
export async function getFiscalPeriodById(id: string): Promise<FiscalPeriodRecord> {
  const res = await authFetch<FiscalPeriodRecord>(`/accounting/fiscal-periods/${encodeURIComponent(id)}`);
  return res.data;
}

/**
 * Mengunci sementara periode bulanan (LOCKED)
 */
export async function lockFiscalPeriod(id: string, payload?: FiscalPeriodActionDTO): Promise<FiscalPeriodRecord> {
  const res = await authFetch<FiscalPeriodRecord>(`/accounting/fiscal-periods/${encodeURIComponent(id)}/lock`, {
    method: 'PUT',
    body: JSON.stringify(payload || {}),
  });
  return res.data;
}

/**
 * Membuka kembali periode yang terkunci (LOCKED -> OPEN)
 */
export async function unlockFiscalPeriod(id: string, payload?: FiscalPeriodActionDTO): Promise<FiscalPeriodRecord> {
  const res = await authFetch<FiscalPeriodRecord>(`/accounting/fiscal-periods/${encodeURIComponent(id)}/unlock`, {
    method: 'PUT',
    body: JSON.stringify(payload || {}),
  });
  return res.data;
}

/**
 * Menutup buku permanen sebuah periode bulanan (CLOSED)
 */
export async function closeFiscalPeriod(id: string, payload?: FiscalPeriodActionDTO): Promise<FiscalPeriodRecord> {
  const res = await authFetch<FiscalPeriodRecord>(`/accounting/fiscal-periods/${encodeURIComponent(id)}/close`, {
    method: 'PUT',
    body: JSON.stringify(payload || {}),
  });
  return res.data;
}
