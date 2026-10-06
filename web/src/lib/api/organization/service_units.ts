import { authFetch, type PaginatedData } from '../client';
import type { ServiceUnitRecord, CreateServiceUnitDTO, UpdateServiceUnitDTO } from '../../types/master/service_unit';

export async function getServiceUnits(params?: {
  search?: string;
  departement_id?: string;
  is_registration_target?: boolean;
  is_active?: boolean;
  page?: number;
  limit?: number;
}): Promise<PaginatedData<ServiceUnitRecord>> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.departement_id) q.set('departement_id', params.departement_id);
  if (params?.is_registration_target !== undefined) q.set('is_registration_target', String(params.is_registration_target));
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PaginatedData<ServiceUnitRecord>>(`/service-units${queryStr}`);
  return res.data;
}

export async function getServiceUnitById(id: string): Promise<ServiceUnitRecord> {
  const res = await authFetch<ServiceUnitRecord>(`/service-units/${id}`);
  return res.data;
}

export async function createServiceUnit(payload: CreateServiceUnitDTO): Promise<ServiceUnitRecord> {
  const res = await authFetch<ServiceUnitRecord>('/service-units', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateServiceUnit(id: string, payload: UpdateServiceUnitDTO): Promise<ServiceUnitRecord> {
  const res = await authFetch<ServiceUnitRecord>(`/service-units/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteServiceUnit(id: string): Promise<void> {
  await authFetch<null>(`/service-units/${id}`, {
    method: 'DELETE',
  });
}
