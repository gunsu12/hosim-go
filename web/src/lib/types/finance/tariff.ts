// Domain Finance: Master Tarif & Komponen Tarif Types
// Corresponds to backend internal/finance/tariffclass and internal/finance/tariffcomponent

export interface TariffClassRecord {
  id: string;
  code: string;
  name: string;
  is_active: boolean;
  description?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateTariffClassDTO {
  code: string;
  name: string;
  is_active?: boolean;
  description?: string | null;
}

export interface UpdateTariffClassDTO {
  code: string;
  name: string;
  is_active?: boolean;
  description?: string | null;
}

export type TariffComponentType =
  | 'JASA_MEDIS'
  | 'JASA_RS'
  | 'SEWA_ALAT'
  | 'BAHAN_ALKES'
  | 'ADMINISTRASI'
  | 'LAINNYA';

export interface TariffComponentRecord {
  id: string;
  code: string;
  name: string;
  description?: string | null;
  is_hospital_revenue: boolean;
  is_operator_revenue: boolean;
  is_paramedic_revenue: boolean;
  component_type: TariffComponentType | string;
  default_coa_code?: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateTariffComponentDTO {
  code: string;
  name: string;
  description?: string | null;
  is_hospital_revenue: boolean;
  is_operator_revenue: boolean;
  is_paramedic_revenue: boolean;
  component_type?: string;
  default_coa_code?: string | null;
  is_active?: boolean;
}

export interface UpdateTariffComponentDTO {
  code: string;
  name: string;
  description?: string | null;
  is_hospital_revenue: boolean;
  is_operator_revenue: boolean;
  is_paramedic_revenue: boolean;
  component_type?: string;
  default_coa_code?: string | null;
  is_active?: boolean;
}

export interface TariffListParams {
  page?: number;
  limit?: number;
  search?: string;
  is_active?: boolean;
  component_type?: string;
}

export interface TariffPaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface TariffClassListResponse {
  data: TariffClassRecord[];
  meta: TariffPaginationMeta;
}

export interface TariffComponentListResponse {
  data: TariffComponentRecord[];
  meta: TariffPaginationMeta;
}
