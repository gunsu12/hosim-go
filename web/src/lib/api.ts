import type { DepartementRecord, CreateDepartementDTO, UpdateDepartementDTO } from './types/master/departement';
import type { ServiceUnitRecord, CreateServiceUnitDTO, UpdateServiceUnitDTO } from './types/master/service_unit';
import type {
  AccountRecord,
  AccountTreeNode,
  CreateAccountDTO,
  UpdateAccountDTO,
  AccountListParams,
  AccountListResponse
} from './types/accounting/account';
import type {
  UserRecord,
  CreateUserDTO,
  UpdateUserDTO,
  RoleRecord,
  CreateRoleDTO,
  UpdateRoleDTO,
  PermissionRecord
} from './types/auth';

const API_BASE = '/api/v1';
const STORAGE_KEY = 'hosim_auth_session';

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  errors?: any;
}

export interface PaginatedData<T> {
  data: T[];
  count: number;
}

export interface HealthResponse {
  status: 'success' | 'error';
  message: string;
  data?: {
    app_name: string;
    env: string;
    database: string;
    driver: string;
    timestamp: string;
  };
  error?: string;
}

export function getToken(): string | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      return parsed.access_token || null;
    }
  } catch (e) {
    console.error('Failed to read token from storage', e);
  }
  return null;
}

let refreshPromise: Promise<string | null> | null = null;

export async function refreshAuthToken(): Promise<string | null> {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return null;

  let session: any = null;
  try {
    session = JSON.parse(raw);
  } catch {
    return null;
  }

  const refreshToken = session?.refresh_token;
  if (!refreshToken) {
    localStorage.removeItem(STORAGE_KEY);
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('hosim:session-expired'));
    }
    return null;
  }

  try {
    const res = await fetch(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    const json = await res.json();
    if (!res.ok || !json.success || !json.data) {
      console.warn('[AUTH] Sesi login telah berakhir. Silakan login kembali.');
      localStorage.removeItem(STORAGE_KEY);
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('hosim:session-expired'));
      }
      return null;
    }

    const newSession = json.data;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newSession));
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('hosim:session-refreshed', { detail: newSession }));
    }
    return newSession.access_token;
  } catch (err) {
    console.error('[AUTH] Gagal memperbarui token:', err);
    return null;
  }
}

export async function authFetch<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
  let token = getToken();
  const headers: Record<string, string> = {
    'Accept': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  if (options.body && !(options.body instanceof FormData) && !headers['Content-Type']) {
    headers['Content-Type'] = 'application/json';
  }

  let res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  // Jika response 401 Unauthorized (token kedaluwarsa), lakukan refresh token secara otomatis
  if (res.status === 401 && !endpoint.includes('/auth/login') && !endpoint.includes('/auth/refresh')) {
    if (!refreshPromise) {
      refreshPromise = refreshAuthToken().finally(() => {
        refreshPromise = null;
      });
    }

    const newToken = await refreshPromise;
    if (newToken) {
      // Retry request yang gagal dengan access token yang baru
      headers['Authorization'] = `Bearer ${newToken}`;
      res = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
      });
    }
  }

  let json: any;
  try {
    json = await res.json();
  } catch {
    json = { success: false, message: `Server error (${res.status})` };
  }

  if (!res.ok) {
    const errorMsg = json?.message || json?.error || `HTTP error ${res.status}`;
    throw new Error(errorMsg);
  }

  return json;
}

export async function checkBackendHealth(): Promise<HealthResponse> {
  try {
    const res = await fetch(`${API_BASE}/health`, {
      headers: { 'Accept': 'application/json' }
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const json = await res.json();
    return {
      status: 'success',
      message: json.message || 'HOSIM-GO Backend running',
      data: json.data
    };
  } catch (err) {
    return {
      status: 'error',
      message: 'Tidak dapat terhubung ke Backend Go Gin',
      error: err instanceof Error ? err.message : String(err)
    };
  }
}

export async function login(username: string, password: string): Promise<any> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Login gagal');
  return data;
}

// ==========================================
// AUTH: USERS, ROLES & PERMISSIONS API
// ==========================================

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

export async function getPermissions(): Promise<PermissionRecord[]> {
  const res = await authFetch<PermissionRecord[]>('/auth/permissions');
  return res.data;
}

// ==========================================
// MASTER DEPARTEMENT API
// ==========================================

export async function getDepartments(params?: { search?: string; page?: number; limit?: number }): Promise<PaginatedData<DepartementRecord>> {
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

// ==========================================
// MASTER SERVICE UNIT API
// ==========================================

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
