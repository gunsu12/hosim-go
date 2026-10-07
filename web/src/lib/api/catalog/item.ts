import { authFetch } from '../client';
import type {
  ItemSummaryRecord,
  ItemDetailRecord,
  ItemUnitRecord,
  ItemCategoryRecord,
  ItemProductLineRecord,
  ItemQueryParams,
  CreateMedicationDTO,
  UpdateMedicationDTO,
  CreateGeneralDTO,
  UpdateGeneralDTO,
  CreateAssetDTO,
  UpdateAssetDTO,
  CreateTariffDTO,
  UpdateTariffDTO,
  AddUnitDTO,
  UpdateUnitDTO,
  CreateCategoryDTO,
  UpdateCategoryDTO,
  CreateProductLineDTO,
  UpdateProductLineDTO
} from '../../types/master/item';

export interface PaginationMeta {
  current_page: number;
  per_page: number;
  total_items: number;
  total_pages: number;
}

export interface PaginatedResult<T> {
  data: T[];
  meta: PaginationMeta;
}

// ==========================================
// 1. UNIVERSAL ITEM CATALOG APIs
// ==========================================

export async function getItems(params?: ItemQueryParams): Promise<PaginatedResult<ItemSummaryRecord>> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.item_type) q.set('item_type', params.item_type);
  if (params?.category_id) q.set('category_id', params.category_id);
  if (params?.product_line_id) q.set('product_line_id', params.product_line_id);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<ItemSummaryRecord[]>(`/items${queryStr}`);
  return {
    data: res.data || [],
    meta: res.meta || { current_page: 1, per_page: 10, total_items: res.data?.length || 0, total_pages: 1 }
  };
}

export async function getItemById(id: string): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>(`/items/${id}`);
  return res.data;
}

export async function deleteItem(id: string): Promise<void> {
  await authFetch<null>(`/items/${id}`, {
    method: 'DELETE'
  });
}

// ==========================================
// 2. SUBTYPE MUTATIONS (MEDICATION, GENERAL, ASSET, TARIFF)
// ==========================================

export async function createMedicationItem(payload: CreateMedicationDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>('/items/medications', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateMedicationItem(id: string, payload: UpdateMedicationDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>(`/items/medications/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function createGeneralItem(payload: CreateGeneralDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>('/items/generals', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateGeneralItem(id: string, payload: UpdateGeneralDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>(`/items/generals/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function createAssetItem(payload: CreateAssetDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>('/items/assets', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateAssetItem(id: string, payload: UpdateAssetDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>(`/items/assets/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function createTariffItem(payload: CreateTariffDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>('/items/tariffs', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateTariffItem(id: string, payload: UpdateTariffDTO): Promise<ItemDetailRecord> {
  const res = await authFetch<ItemDetailRecord>(`/items/tariffs/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

// ==========================================
// 3. MULTI-UOM SATUAN ITEM APIs
// ==========================================

export async function getItemUnits(itemId: string): Promise<ItemUnitRecord[]> {
  const res = await authFetch<ItemUnitRecord[]>(`/items/${itemId}/units`);
  return res.data || [];
}

export async function addItemUnit(itemId: string, payload: AddUnitDTO): Promise<ItemUnitRecord> {
  const res = await authFetch<ItemUnitRecord>(`/items/${itemId}/units`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateItemUnit(itemId: string, unitId: string, payload: UpdateUnitDTO): Promise<ItemUnitRecord> {
  const res = await authFetch<ItemUnitRecord>(`/items/${itemId}/units/${unitId}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function deleteItemUnit(itemId: string, unitId: string): Promise<void> {
  await authFetch<null>(`/items/${itemId}/units/${unitId}`, {
    method: 'DELETE'
  });
}

// ==========================================
// 4. KATEGORI ITEM APIs
// ==========================================

export async function getItemCategories(params?: {
  page?: number;
  limit?: number;
  search?: string;
  item_type?: string;
  is_active?: boolean;
}): Promise<PaginatedResult<ItemCategoryRecord>> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.item_type) q.set('item_type', params.item_type);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<ItemCategoryRecord[]>(`/item-categories${queryStr}`);
  return {
    data: res.data || [],
    meta: res.meta || { current_page: 1, per_page: 50, total_items: res.data?.length || 0, total_pages: 1 }
  };
}

export async function getCategoryById(id: string): Promise<ItemCategoryRecord> {
  const res = await authFetch<ItemCategoryRecord>(`/item-categories/${id}`);
  return res.data;
}

export async function createCategory(payload: CreateCategoryDTO): Promise<ItemCategoryRecord> {
  const res = await authFetch<ItemCategoryRecord>('/item-categories', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateCategory(id: string, payload: UpdateCategoryDTO): Promise<ItemCategoryRecord> {
  const res = await authFetch<ItemCategoryRecord>(`/item-categories/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function deleteCategory(id: string): Promise<void> {
  await authFetch<null>(`/item-categories/${id}`, {
    method: 'DELETE'
  });
}

// ==========================================
// 5. LINI PRODUK (PRODUCT LINE) APIs
// ==========================================

export async function getItemProductLines(params?: {
  page?: number;
  limit?: number;
  search?: string;
  is_active?: boolean;
}): Promise<PaginatedResult<ItemProductLineRecord>> {
  const q = new URLSearchParams();
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.search) q.set('search', params.search);
  if (params?.is_active !== undefined) q.set('is_active', String(params.is_active));

  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<ItemProductLineRecord[]>(`/item-product-lines${queryStr}`);
  return {
    data: res.data || [],
    meta: res.meta || { current_page: 1, per_page: 50, total_items: res.data?.length || 0, total_pages: 1 }
  };
}

export async function getProductLineById(id: string): Promise<ItemProductLineRecord> {
  const res = await authFetch<ItemProductLineRecord>(`/item-product-lines/${id}`);
  return res.data;
}

export async function createProductLine(payload: CreateProductLineDTO): Promise<ItemProductLineRecord> {
  const res = await authFetch<ItemProductLineRecord>('/item-product-lines', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function updateProductLine(id: string, payload: UpdateProductLineDTO): Promise<ItemProductLineRecord> {
  const res = await authFetch<ItemProductLineRecord>(`/item-product-lines/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  return res.data;
}

export async function deleteProductLine(id: string): Promise<void> {
  await authFetch<null>(`/item-product-lines/${id}`, {
    method: 'DELETE'
  });
}
