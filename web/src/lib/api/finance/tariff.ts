import { authFetch } from '../client';
import type {
  TariffClassRecord,
  CreateTariffClassDTO,
  UpdateTariffClassDTO,
  TariffComponentRecord,
  CreateTariffComponentDTO,
  UpdateTariffComponentDTO,
  TariffListParams,
  TariffClassListResponse,
  TariffComponentListResponse
} from '../../types/finance/tariff';

// ==========================================
// FINANCE: TARIFF CLASSES (KELAS TARIF) API
// ==========================================

export async function getTariffClasses(params?: TariffListParams): Promise<TariffClassListResponse> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<TariffClassListResponse>(`/tariff-classes${queryStr}`);
  return res.data;
}

export async function getTariffClassById(id: string): Promise<TariffClassRecord> {
  const res = await authFetch<TariffClassRecord>(`/tariff-classes/${id}`);
  return res.data;
}

export async function createTariffClass(payload: CreateTariffClassDTO): Promise<TariffClassRecord> {
  const res = await authFetch<TariffClassRecord>('/tariff-classes', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateTariffClass(id: string, payload: UpdateTariffClassDTO): Promise<TariffClassRecord> {
  const res = await authFetch<TariffClassRecord>(`/tariff-classes/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteTariffClass(id: string): Promise<void> {
  await authFetch<null>(`/tariff-classes/${id}`, {
    method: 'DELETE',
  });
}

// ==========================================
// FINANCE: TARIFF COMPONENTS (KOMPONEN TARIF) API
// ==========================================

export async function getTariffComponents(params?: TariffListParams): Promise<TariffComponentListResponse> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.component_type) q.set('component_type', params.component_type);
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<TariffComponentListResponse>(`/tariff-components${queryStr}`);
  return res.data;
}

export async function getTariffComponentById(id: string): Promise<TariffComponentRecord> {
  const res = await authFetch<TariffComponentRecord>(`/tariff-components/${id}`);
  return res.data;
}

export async function createTariffComponent(payload: CreateTariffComponentDTO): Promise<TariffComponentRecord> {
  const res = await authFetch<TariffComponentRecord>('/tariff-components', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateTariffComponent(id: string, payload: UpdateTariffComponentDTO): Promise<TariffComponentRecord> {
  const res = await authFetch<TariffComponentRecord>(`/tariff-components/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteTariffComponent(id: string): Promise<void> {
  await authFetch<null>(`/tariff-components/${id}`, {
    method: 'DELETE',
  });
}
