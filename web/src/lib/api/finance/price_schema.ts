import { authFetch, type ApiResponse } from '../client';
import type {
  PriceSchemaRecord,
  PriceSchemaRuleRecord,
  CreatePriceSchemaDTO,
  UpdatePriceSchemaDTO,
  UpsertRuleDTO,
  SimulateTariffDTO,
  SimulateTariffResult
} from '../../types/finance/price_schema';

export interface PriceSchemaListParams {
  q?: string;
  is_active?: boolean;
  page?: number;
  limit?: number;
}

export interface PriceSchemaListResponse {
  items: PriceSchemaRecord[];
  meta?: {
    current_page: number;
    per_page: number;
    total_items: number;
    total_pages: number;
  };
}

export interface RuleListParams {
  level?: string;
  departement_id?: string;
  service_unit_id?: string;
  item_id?: string;
  page?: number;
  limit?: number;
}

export interface RuleListResponse {
  items: PriceSchemaRuleRecord[];
  meta?: {
    current_page: number;
    per_page: number;
    total_items: number;
    total_pages: number;
  };
}

// 1. Get List of Price Schemas
export async function getPriceSchemas(
  params?: PriceSchemaListParams
): Promise<PriceSchemaListResponse> {
  const q = new URLSearchParams();
  if (params?.q) q.set('q', params.q);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PriceSchemaRecord[]>(`/finance/price-schemas${queryStr}`);
  return {
    items: res.data || [],
    meta: res.meta
  };
}

// 2. Get Price Schema By ID
export async function getPriceSchemaById(id: string): Promise<PriceSchemaRecord> {
  const res = await authFetch<PriceSchemaRecord>(`/finance/price-schemas/${id}`);
  return res.data;
}

// 3. Create Price Schema
export async function createPriceSchema(
  payload: CreatePriceSchemaDTO
): Promise<PriceSchemaRecord> {
  const res = await authFetch<PriceSchemaRecord>('/finance/price-schemas', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

// 4. Update Price Schema Header
export async function updatePriceSchema(
  id: string,
  payload: UpdatePriceSchemaDTO
): Promise<PriceSchemaRecord> {
  const res = await authFetch<PriceSchemaRecord>(`/finance/price-schemas/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

// 5. Deactivate Price Schema
export async function deactivatePriceSchema(id: string): Promise<PriceSchemaRecord> {
  const res = await authFetch<PriceSchemaRecord>(
    `/finance/price-schemas/${id}/deactivate`,
    {
      method: 'POST'
    }
  );
  return res.data;
}

// 6. Get Rules for a Price Schema
export async function getPriceSchemaRules(
  schemaId: string,
  params?: RuleListParams
): Promise<RuleListResponse> {
  const q = new URLSearchParams();
  if (params?.level) q.set('level', params.level);
  if (params?.departement_id) q.set('departement_id', params.departement_id);
  if (params?.service_unit_id) q.set('service_unit_id', params.service_unit_id);
  if (params?.item_id) q.set('item_id', params.item_id);
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PriceSchemaRuleRecord[]>(
    `/finance/price-schemas/${schemaId}/rules${queryStr}`
  );
  return {
    items: res.data || [],
    meta: res.meta
  };
}

// 7. Upsert Rule in a Price Schema
export async function upsertPriceSchemaRule(
  schemaId: string,
  payload: UpsertRuleDTO
): Promise<PriceSchemaRuleRecord> {
  const res = await authFetch<PriceSchemaRuleRecord>(
    `/finance/price-schemas/${schemaId}/rules`,
    {
      method: 'PUT',
      body: JSON.stringify(payload)
    }
  );
  return res.data;
}

// 8. Delete / Reset Rule from a Price Schema
export async function deletePriceSchemaRule(
  schemaId: string,
  ruleId: string
): Promise<void> {
  await authFetch<null>(
    `/finance/price-schemas/${schemaId}/rules/${ruleId}`,
    {
      method: 'DELETE'
    }
  );
}

// 9. Simulate Price Schema Tariff Calculation
export async function simulatePriceSchemaTariff(
  payload: SimulateTariffDTO
): Promise<SimulateTariffResult> {
  const res = await authFetch<SimulateTariffResult>(
    '/finance/price-schemas/simulate',
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
  return res.data;
}
