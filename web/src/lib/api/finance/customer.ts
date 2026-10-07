import { authFetch } from '../client';
import type { CustomerRecord } from '../../types/master/customer';

export interface CustomerListResponse {
  data: CustomerRecord[];
  meta?: any;
}

export async function getCustomers(params?: { search?: string; is_active?: boolean; limit?: number }): Promise<CustomerListResponse> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.limit) q.set('limit', String(params.limit));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<CustomerListResponse>(`/customers${queryStr}`);
  return res.data;
}
