// TypeScript Definitions for Catalog Item Domain
// Matches backend internal/catalog/item (Class Table Inheritance)

export type ItemType = 'MEDICATION' | 'GENERAL' | 'ASSET' | 'TARIFF';

export type MedicationType =
  | 'Bebas'
  | 'Keras'
  | 'Narkotika'
  | 'Psikotropika'
  | 'Prekursor'
  | 'Lainnya';

export type GeneralType =
  | 'BMHP_MEDIS'
  | 'INSTRUMEN_MEDIS'
  | 'ATK'
  | 'LINEN'
  | 'KEBERSIHAN'
  | 'DAPUR'
  | 'LAINNYA';

export type SterilizationMethod =
  | 'STEAM_AUTOCLAVE'
  | 'EO_GAS'
  | 'PLASMA'
  | 'DRY_HEAT';

export type DepreciationMethod =
  | 'STRAIGHT_LINE'
  | 'DOUBLE_DECLINING'
  | 'SUM_OF_YEARS';

export type ChargeType =
  | 'ADMINISTRASI'
  | 'AKOMODASI'
  | 'TINDAKAN'
  | 'PENUNJANG'
  | 'LAINNYA';

// ==========================================
// Unit of Measure (UOM)
// ==========================================
export interface ItemUnitRecord {
  id: string;
  unit_name: string;
  conversion_factor: number;
  is_base_unit: boolean;
  is_purchase_unit: boolean;
  is_dispense_unit: boolean;
  created_at: string;
}

export interface AddUnitDTO {
  unit_name: string;
  conversion_factor: number;
  is_purchase_unit: boolean;
  is_dispense_unit: boolean;
}

export interface UpdateUnitDTO {
  unit_name: string;
  conversion_factor: number;
  is_purchase_unit: boolean;
  is_dispense_unit: boolean;
}

// ==========================================
// Category & Product Line Summaries
// ==========================================
export interface CategorySummary {
  id: string;
  code: string;
  name: string;
  item_type: ItemType;
  income_coa_code?: string;
}

export interface ProductLineSummary {
  id: string;
  code: string;
  name: string;
  inventory_coa_code?: string;
  cogs_coa_code?: string;
}

// ==========================================
// Subtype Details
// ==========================================
export interface MedicationDetail {
  kfa_code?: string;
  bpom_nie?: string;
  dosage_form?: string;
  strength_amount?: string;
  strength_unit?: string;
  default_route?: string;
  medication_type?: MedicationType;
  is_high_alert: boolean;
  is_lasa: boolean;
  is_fornas: boolean;
  is_antibiotic: boolean;
  storage_temperature?: string;
}

export interface GeneralDetail {
  general_type: GeneralType;
  is_sterile: boolean;
  is_disposable: boolean;
  is_cssd_item: boolean;
  sterilization_method?: SterilizationMethod;
}

export interface AssetDetail {
  brand?: string;
  model_name?: string;
  is_medical_equipment: boolean;
  expected_life_years?: number;
  depreciation_method?: DepreciationMethod;
  maintenance_interval_days?: number;
}

export interface TariffDetail {
  name_alias?: string;
  charge_type: ChargeType;
  notes?: string;
}

// ==========================================
// Full Item Records
// ==========================================
export interface ItemSummaryRecord {
  id: string;
  code: string;
  name: string;
  generic_name?: string;
  item_type: ItemType;
  is_active: boolean;
  category?: CategorySummary;
  product_line?: ProductLineSummary;
  base_unit?: ItemUnitRecord;
  created_at: string;
}

export interface ItemDetailRecord {
  id: string;
  code: string;
  name: string;
  generic_name?: string;
  item_type: ItemType;
  is_active: boolean;
  category?: CategorySummary;
  product_line?: ProductLineSummary;
  medication?: MedicationDetail;
  general?: GeneralDetail;
  asset?: AssetDetail;
  tariff?: TariffDetail;
  units: ItemUnitRecord[];
  created_at: string;
  updated_at: string;
}

// ==========================================
// Category Full Record
// ==========================================
export interface ItemCategoryRecord {
  id: string;
  code: string;
  name: string;
  item_type: ItemType;
  income_coa_code?: string;
  discount_coa_code?: string;
  sales_tax_coa_code?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCategoryDTO {
  code: string;
  name: string;
  item_type: ItemType;
  income_coa_code?: string;
  discount_coa_code?: string;
  sales_tax_coa_code?: string;
  is_active?: boolean;
}

export interface UpdateCategoryDTO {
  code: string;
  name: string;
  item_type: ItemType;
  income_coa_code?: string;
  discount_coa_code?: string;
  sales_tax_coa_code?: string;
  is_active?: boolean;
}

// ==========================================
// Product Line Full Record
// ==========================================
export interface ItemProductLineRecord {
  id: string;
  code: string;
  name: string;
  inventory_coa_code?: string;
  cogs_coa_code?: string;
  purchase_discount_coa_code?: string;
  purchase_tax_coa_code?: string;
  asset_coa_code?: string;
  asset_accumulation_coa_code?: string;
  asset_depreciation_expense_coa_code?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateProductLineDTO {
  code: string;
  name: string;
  inventory_coa_code?: string;
  cogs_coa_code?: string;
  purchase_discount_coa_code?: string;
  purchase_tax_coa_code?: string;
  asset_coa_code?: string;
  asset_accumulation_coa_code?: string;
  asset_depreciation_expense_coa_code?: string;
  is_active?: boolean;
}

export interface UpdateProductLineDTO {
  code: string;
  name: string;
  inventory_coa_code?: string;
  cogs_coa_code?: string;
  purchase_discount_coa_code?: string;
  purchase_tax_coa_code?: string;
  asset_coa_code?: string;
  asset_accumulation_coa_code?: string;
  asset_depreciation_expense_coa_code?: string;
  is_active?: boolean;
}

// ==========================================
// Create & Update Item DTOs
// ==========================================
export interface CreateMedicationDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  base_unit_name: string;
  kfa_code?: string;
  bpom_nie?: string;
  dosage_form?: string;
  strength_amount?: string;
  strength_unit?: string;
  default_route?: string;
  medication_type?: MedicationType;
  is_high_alert?: boolean;
  is_lasa?: boolean;
  is_fornas?: boolean;
  is_antibiotic?: boolean;
  storage_temperature?: string;
}

export interface UpdateMedicationDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  kfa_code?: string;
  bpom_nie?: string;
  dosage_form?: string;
  strength_amount?: string;
  strength_unit?: string;
  default_route?: string;
  medication_type?: MedicationType;
  is_high_alert?: boolean;
  is_lasa?: boolean;
  is_fornas?: boolean;
  is_antibiotic?: boolean;
  storage_temperature?: string;
  is_active?: boolean;
}

export interface CreateGeneralDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  base_unit_name: string;
  general_type: GeneralType;
  is_sterile?: boolean;
  is_disposable?: boolean;
  is_cssd_item?: boolean;
  sterilization_method?: SterilizationMethod;
}

export interface UpdateGeneralDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  general_type: GeneralType;
  is_sterile?: boolean;
  is_disposable?: boolean;
  is_cssd_item?: boolean;
  sterilization_method?: SterilizationMethod;
  is_active?: boolean;
}

export interface CreateAssetDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  base_unit_name: string;
  brand?: string;
  model_name?: string;
  is_medical_equipment?: boolean;
  expected_life_years?: number;
  depreciation_method?: DepreciationMethod;
  maintenance_interval_days?: number;
}

export interface UpdateAssetDTO {
  code: string;
  name: string;
  generic_name?: string;
  category_id: string;
  product_line_id?: string;
  brand?: string;
  model_name?: string;
  is_medical_equipment?: boolean;
  expected_life_years?: number;
  depreciation_method?: DepreciationMethod;
  maintenance_interval_days?: number;
  is_active?: boolean;
}

export interface CreateTariffDTO {
  code: string;
  name: string;
  name_alias?: string;
  category_id: string;
  base_unit_name?: string;
  charge_type: ChargeType;
  notes?: string;
}

export interface UpdateTariffDTO {
  code: string;
  name: string;
  name_alias?: string;
  category_id: string;
  charge_type: ChargeType;
  notes?: string;
  is_active?: boolean;
}

// ==========================================
// Query Parameters & Pagination
// ==========================================
export interface ItemQueryParams {
  page?: number;
  limit?: number;
  search?: string;
  item_type?: ItemType;
  category_id?: string;
  product_line_id?: string;
  is_active?: boolean;
}

export interface ItemListResponse {
  data: ItemSummaryRecord[];
  meta: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
