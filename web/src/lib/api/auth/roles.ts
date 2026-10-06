import { authFetch } from '../client';
import type { RoleRecord, CreateRoleDTO, UpdateRoleDTO } from '../../types/auth';

export async function getRoles(): Promise<RoleRecord[]> {
  const res = await authFetch<RoleRecord[]>('/auth/roles');
  return res.data;
}

export async function createRole(payload: CreateRoleDTO): Promise<RoleRecord> {
  const res = await authFetch<RoleRecord>('/auth/roles', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateRole(id: string, payload: UpdateRoleDTO): Promise<RoleRecord> {
  const res = await authFetch<RoleRecord>(`/auth/roles/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteRole(id: string): Promise<void> {
  await authFetch<null>(`/auth/roles/${id}`, {
    method: 'DELETE',
  });
}
