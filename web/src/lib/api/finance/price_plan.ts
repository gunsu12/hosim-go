import { authFetch } from '../client';
import type {
  PricePlanRecord,
  CreatePricePlanDTO,
  UpdatePricePlanDTO,
  ClonePricePlanDTO,
  PricePlanItemRecord,
  AddPricePlanItemDTO,
  UpdatePricePlanItemDTO,
  PricePlanListParams,
  PricePlanItemListParams,
  PricePlanListResponse,
  PricePlanItemListResponse,
  LookupTariffDTO,
  LookupTariffResponse
} from '../../types/finance/price_plan';

// ==========================================
// FINANCE: TARIFF PRICE PLAN (BUKU TARIF) API
// ==========================================

export async function getPricePlans(params?: PricePlanListParams): Promise<PricePlanListResponse> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.status) q.set('status', params.status);
  if (params?.customer_id) q.set('customer_id', params.customer_id);
  if (params?.is_default !== undefined) q.set('is_default', String(params.is_default));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<any>(`/price-plans${queryStr}`);
  
  const rawList = Array.isArray(res.data) ? res.data : (res.data?.data || []);
  const rawMeta = res.meta || res.data?.meta || {
    current_page: params?.page || 1,
    per_page: params?.limit || 10,
    total_items: rawList.length,
    total_pages: 1
  };

  return {
    data: rawList,
    meta: rawMeta
  };
}

export async function getPricePlanById(id: string): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}`);
  return res.data;
}

export async function createPricePlan(payload: CreatePricePlanDTO): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>('/price-plans', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updatePricePlan(id: string, payload: UpdatePricePlanDTO): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function clonePricePlan(id: string, payload: ClonePricePlanDTO): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}/clone`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

// ------------------------------------------
// Lifecycle State Machine Transitions
// ------------------------------------------

export async function submitPricePlan(id: string): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}/submit`, {
    method: 'POST',
  });
  return res.data;
}

export async function approvePricePlan(id: string): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}/approve`, {
    method: 'POST',
  });
  return res.data;
}

export async function activatePricePlan(id: string): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}/activate`, {
    method: 'POST',
  });
  return res.data;
}

export async function archivePricePlan(id: string): Promise<PricePlanRecord> {
  const res = await authFetch<PricePlanRecord>(`/price-plans/${id}/archive`, {
    method: 'POST',
  });
  return res.data;
}

// ==========================================
// PRICE PLAN ITEMS & COMPONENTS API
// ==========================================

export async function getPricePlanItems(
  planId: string,
  params?: PricePlanItemListParams
): Promise<PricePlanItemListResponse> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.tariff_class_id) q.set('tariff_class_id', params.tariff_class_id);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<any>(`/price-plans/${planId}/items${queryStr}`);
  
  const rawList = Array.isArray(res.data) ? res.data : (res.data?.data || []);
  const rawMeta = res.meta || res.data?.meta || {
    current_page: params?.page || 1,
    per_page: params?.limit || 10,
    total_items: rawList.length,
    total_pages: 1
  };

  return {
    data: rawList,
    meta: rawMeta
  };
}

export async function getPricePlanItemById(planId: string, itemId: string): Promise<PricePlanItemRecord> {
  const res = await authFetch<PricePlanItemRecord>(`/price-plans/${planId}/items/${itemId}`);
  return res.data;
}

export async function addPricePlanItem(planId: string, payload: AddPricePlanItemDTO): Promise<PricePlanItemRecord> {
  const res = await authFetch<PricePlanItemRecord>(`/price-plans/${planId}/items`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function updatePricePlanItem(
  planId: string,
  itemId: string,
  payload: UpdatePricePlanItemDTO
): Promise<PricePlanItemRecord> {
  const res = await authFetch<PricePlanItemRecord>(`/price-plans/${planId}/items/${itemId}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
  return res.data;
}

export async function deletePricePlanItem(planId: string, itemId: string): Promise<void> {
  await authFetch<null>(`/price-plans/${planId}/items/${itemId}`, {
    method: 'DELETE',
  });
}

// ==========================================
// CORE TARIFF LOOKUP ENGINE (PRD § 6.4)
// ==========================================

export async function lookupTariff(payload: LookupTariffDTO): Promise<LookupTariffResponse> {
  const res = await authFetch<LookupTariffResponse>('/price-plans/lookup', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
  return res.data;
}
