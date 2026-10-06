import { authFetch, type PaginatedData } from '../client';
import type { UserRecord, CreateUserDTO, UpdateUserDTO } from '../../types/auth';

export async function getUsers(params?: {
  search?: string;
  role_id?: string;
  is_active?: boolean;
  page?: number;
  limit?: number;
}): Promise<PaginatedData<UserRecord>> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.role_id) q.set('role_id', params.role_id);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PaginatedData<UserRecord>>(`/auth/users${queryStr}`);
  return res.data;
}

export async function getUserById(id: string): Promise<UserRecord> {
  const res = await authFetch<UserRecord>(`/auth/users/${id}`);
  return res.data;
}

export async function createUser(payload: CreateUserDTO): Promise<UserRecord> {
  const res = await authFetch<UserRecord>('/auth/users', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateUser(id: string, payload: UpdateUserDTO): Promise<UserRecord> {
  const res = await authFetch<UserRecord>(`/auth/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteUser(id: string): Promise<void> {
  await authFetch<null>(`/auth/users/${id}`, {
    method: 'DELETE',
  });
}
