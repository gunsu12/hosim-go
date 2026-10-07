// Domain Finance: Tariff Price Plan & Lookup Engine Types
// Corresponds to backend internal/finance/priceplan (PRD.md & dto.go)

export type PricePlanStatus = 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'ACTIVE' | 'ARCHIVED';

export interface PricePlanRecord {
  id: string;
  code: string;
  name: string;
  description?: string | null;
  effective_from: string;
  effective_to?: string | null;
  status: PricePlanStatus;
  is_default: boolean;
  customer_id?: string | null;
  customer_name?: string | null;
  default_cito_percent: number;
  approved_at?: string | null;
  approved_by?: string | null;
  total_items: number;
  created_at: string;
  updated_at: string;
}

export interface CreatePricePlanDTO {
  code: string;
  name: string;
  description?: string | null;
  effective_from: string;
  effective_to?: string | null;
  is_default?: boolean;
  customer_id?: string | null;
  default_cito_percent?: number;
}

export interface UpdatePricePlanDTO {
  name: string;
  description?: string | null;
  effective_from: string;
  effective_to?: string | null;
  is_default?: boolean;
  customer_id?: string | null;
  default_cito_percent?: number;
}

export interface ClonePricePlanDTO {
  new_code: string;
  new_name: string;
  description?: string | null;
  effective_from: string;
  effective_to?: string | null;
  is_default?: boolean;
  customer_id?: string | null;
  default_cito_percent?: number;
}

export interface ItemComponentDTO {
  component_id: string;
  base_amount: number;
  cito_amount?: number | null;
  coa_code?: string | null;
}

export interface AddPricePlanItemDTO {
  item_id: string;
  tariff_class_id: string;
  total_base_price: number;
  total_cito_price?: number | null;
  is_active?: boolean;
  components: ItemComponentDTO[];
}

export interface UpdatePricePlanItemDTO {
  total_base_price: number;
  total_cito_price?: number | null;
  is_active?: boolean;
  components: ItemComponentDTO[];
}

export interface PricePlanItemComponentRecord {
  id: string;
  plan_item_id: string;
  component_id: string;
  component_code?: string;
  component_name?: string;
  component_type?: string;
  base_amount: number;
  cito_amount?: number | null;
  coa_code?: string | null;
}

export interface PricePlanItemRecord {
  id: string;
  price_plan_id: string;
  item_id: string;
  item_code?: string;
  item_name?: string;
  tariff_class_id: string;
  tariff_class_code?: string;
  tariff_class_name?: string;
  total_base_price: number;
  total_cito_price?: number | null;
  is_active: boolean;
  components: PricePlanItemComponentRecord[];
}

export interface PricePlanListParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: PricePlanStatus;
  customer_id?: string;
  is_default?: boolean;
}

export interface PricePlanItemListParams {
  page?: number;
  limit?: number;
  search?: string;
  item_id?: string;
  tariff_class_id?: string;
  is_active?: boolean;
}

export interface BatchUpsertItemsDTO {
  items: AddPricePlanItemDTO[];
}

export interface PricePlanPaginationMeta {
  current_page: number;
  per_page: number;
  total_items: number;
  total_pages: number;
}

export interface PricePlanListResponse {
  data: PricePlanRecord[];
  meta: PricePlanPaginationMeta;
}

export interface PricePlanItemListResponse {
  data: PricePlanItemRecord[];
  meta: PricePlanPaginationMeta;
}

// Tariff Lookup Engine Types (PRD § 6.4)
export interface LookupTariffDTO {
  customer_id?: string | null;
  tariff_class_id: string;
  item_id: string;
  transaction_date: string;
  is_cito: boolean;
}

export interface ResolvedComponentRecord {
  component_id: string;
  component_name: string;
  component_type: string;
  amount: number;
  coa_code: string;
}

export interface LookupTariffResponse {
  price_plan_id: string;
  price_plan_code: string;
  is_custom_plan: boolean;
  item_id: string;
  item_code?: string;
  item_name?: string;
  tariff_class_id: string;
  tariff_class_code?: string;
  tariff_class_name?: string;
  is_cito: boolean;
  total_price: number;
  components: ResolvedComponentRecord[];
}
