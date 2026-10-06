import { authFetch, type PaginatedData } from '../client';
import type { DepartementRecord, CreateDepartementDTO, UpdateDepartementDTO } from '../../types/master/departement';

export async function getDepartments(params?: {
  search?: string;
  page?: number;
  limit?: number;
}): Promise<PaginatedData<DepartementRecord>> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PaginatedData<DepartementRecord>>(`/departments${queryStr}`);
  return res.data;
}

export async function getDepartmentById(id: string): Promise<DepartementRecord> {
  const res = await authFetch<DepartementRecord>(`/departments/${id}`);
  return res.data;
}

export async function createDepartment(payload: CreateDepartementDTO): Promise<DepartementRecord> {
  const res = await authFetch<DepartementRecord>('/departments', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateDepartment(id: string, payload: UpdateDepartementDTO): Promise<DepartementRecord> {
  const res = await authFetch<DepartementRecord>(`/departments/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteDepartment(id: string): Promise<void> {
  await authFetch<null>(`/departments/${id}`, {
    method: 'DELETE',
  });
}
