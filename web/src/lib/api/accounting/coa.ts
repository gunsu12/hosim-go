import { authFetch } from '../client';
import type {
  AccountRecord,
  AccountTreeNode,
  CreateAccountDTO,
  UpdateAccountDTO,
  AccountListParams,
  AccountListResponse
} from '../../types/accounting/account';

// ==========================================
// ACCOUNTING: CHART OF ACCOUNTS (COA) API
// ==========================================

export async function getAccounts(params?: AccountListParams): Promise<AccountListResponse> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.type) q.set('type', params.type);
  if (params?.position) q.set('position', params.position);
  if (params?.parent_id) q.set('parent_id', params.parent_id);
  if (params?.is_postable !== undefined) q.set('is_postable', String(params.is_postable));
  if (params?.is_treasury !== undefined) q.set('is_treasury', String(params.is_treasury));
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.level !== undefined) q.set('level', String(params.level));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<AccountListResponse>(`/accounting/accounts${queryStr}`);
  return res.data;
}

export async function getAccountTree(activeOnly: boolean = false): Promise<AccountTreeNode[]> {
  const queryStr = activeOnly ? '?active_only=true' : '';
  const res = await authFetch<AccountTreeNode[]>(`/accounting/accounts/tree${queryStr}`);
  return res.data;
}

export async function getAccountById(id: string): Promise<AccountRecord> {
  const res = await authFetch<AccountRecord>(`/accounting/accounts/${id}`);
  return res.data;
}

export async function getPostableAccounts(): Promise<AccountRecord[]> {
  const res = await authFetch<AccountRecord[]>('/accounting/accounts/postable');
  return res.data;
}

export async function getTreasuryAccounts(): Promise<AccountRecord[]> {
  const res = await authFetch<AccountRecord[]>('/accounting/accounts/treasury');
  return res.data;
}

export async function createAccount(payload: CreateAccountDTO): Promise<AccountRecord> {
  const res = await authFetch<AccountRecord>('/accounting/accounts', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updateAccount(id: string, payload: UpdateAccountDTO): Promise<AccountRecord> {
  const res = await authFetch<AccountRecord>(`/accounting/accounts/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deleteAccount(id: string): Promise<void> {
  await authFetch<null>(`/accounting/accounts/${id}`, {
    method: 'DELETE',
  });
}
