import { authFetch } from '../client';
import type { PermissionRecord } from '../../types/auth';

export async function getPermissions(): Promise<PermissionRecord[]> {
  const res = await authFetch<PermissionRecord[]>('/auth/permissions');
  return res.data;
}
