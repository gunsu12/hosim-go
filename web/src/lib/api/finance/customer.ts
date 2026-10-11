import { authFetch } from '../client';
import type {
  CustomerRecord,
  CustomerTypeRecord,
  CreateCustomerDTO,
  UpdateCustomerDTO,
  CustomerListParams,
} from '$lib/types/finance/customer';

export interface CustomerListResponse {
  data: CustomerRecord[];
  count: number;
}

export async function getCustomers(params?: CustomerListParams): Promise<CustomerListResponse> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.customer_type_id) q.set('customer_type_id', params.customer_type_id);
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<CustomerListResponse>(`/customers${queryStr}`);
  return res.data;
}

export async function getCustomerById(id: string): Promise<CustomerRecord> {
  const res = await authFetch<CustomerRecord>(`/customers/${id}`);
  return res.data;
}

export async function createCustomer(dto: CreateCustomerDTO): Promise<CustomerRecord> {
  const res = await authFetch<CustomerRecord>('/customers', {
    method: 'POST',
    body: JSON.stringify(dto),
  });
  return res.data;
}

export async function updateCustomer(id: string, dto: UpdateCustomerDTO): Promise<CustomerRecord> {
  const res = await authFetch<CustomerRecord>(`/customers/${id}`, {
    method: 'PUT',
    body: JSON.stringify(dto),
  });
  return res.data;
}

export async function deleteCustomer(id: string): Promise<void> {
  await authFetch<void>(`/customers/${id}`, {
    method: 'DELETE',
  });
}

export async function getCustomerTypes(): Promise<CustomerTypeRecord[]> {
  const res = await authFetch<CustomerTypeRecord[]>('/customer-types');
  return res.data;
}
