// TypeScript Definitions for Finance Customer Contract (PKS Penjamin)
// Sesuai dengan internal/finance/customercontract/PRD.md

import type { CustomerRecord } from './customer';
import type { PriceSchemaRecord } from './price_schema';

export type ContractStatus = 'DRAFT' | 'ACTIVE' | 'SUSPENDED' | 'TERMINATED';

export type GuaranteeRequirement = 'NONE' | 'SPPK' | 'LETTER_OF_GUARANTEE';

export interface CustomerContractPriceSchemaRecord {
  id: string;
  customer_contract_id: string;
  price_schema_id: string;
  price_schema_name?: string;
  price_schema_code?: string;
  priority: number;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
  price_schema?: PriceSchemaRecord;
}

export interface CustomerContractRecord {
  id: string;
  customer_id: string;
  customer?: CustomerRecord | null;
  customer_name?: string;
  customer_code?: string;
  contract_number: string;
  name: string;
  effective_date: string; // YYYY-MM-DD
  expiration_date?: string | null; // YYYY-MM-DD (null = tanpa batas)
  status: ContractStatus;
  payment_term_days: number;
  payment_term_notes?: string | null;
  credit_limit?: number | null; // null = tanpa batas, 0 = tidak boleh kredit, >0 = limit
  guarantee_requirement: GuaranteeRequirement;
  guarantee_notes?: string | null;
  renewed_from_contract_id?: string | null;
  renewed_from_contract_number?: string | null;
  notes?: string | null;
  activated_at?: string | null;
  activated_by?: string | null;
  terminated_at?: string | null;
  terminated_by?: string | null;
  status_reason?: string | null;
  price_schemas?: CustomerContractPriceSchemaRecord[];
  created_at: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

export interface CreateContractDTO {
  customer_id: string;
  contract_number: string;
  name: string;
  effective_date: string;
  expiration_date?: string | null;
  payment_term_days: number;
  payment_term_notes?: string;
  credit_limit?: number | null;
  guarantee_requirement: GuaranteeRequirement;
  guarantee_notes?: string;
  notes?: string;
  price_schemas?: Array<{
    price_schema_id: string;
    priority: number;
    is_active?: boolean;
  }>;
}

export interface UpdateContractDraftDTO {
  contract_number: string;
  name: string;
  effective_date: string;
  expiration_date?: string | null;
  payment_term_days: number;
  payment_term_notes?: string;
  credit_limit?: number | null;
  guarantee_requirement: GuaranteeRequirement;
  guarantee_notes?: string;
  notes?: string;
}

export interface UpdateContractInfoDTO {
  name: string;
  notes?: string;
  guarantee_notes?: string;
  payment_term_notes?: string;
}

export interface AmendContractTermsDTO {
  payment_term_days: number;
  credit_limit?: number | null;
  expiration_date?: string | null;
  guarantee_requirement: GuaranteeRequirement;
  reason: string; // Wajib menurut PRD D11 & D12
}

export interface RenewContractDTO {
  contract_number: string;
  name: string;
  effective_date: string;
  expiration_date?: string | null;
  notes?: string;
}

export interface SetContractPriceSchemasDTO {
  schemas: Array<{
    price_schema_id: string;
    priority: number;
    is_active?: boolean;
  }>;
}

export interface SelectableContract {
  contract_id: string;
  contract_number: string;
  name: string;
  guarantee_requirement: GuaranteeRequirement;
  price_schemas: Array<{
    price_schema_id: string;
    price_schema_code: string;
    price_schema_name: string;
    priority: number;
  }>;
}

export interface CustomerContractListParams {
  customer_id?: string;
  status?: string;
  valid_on?: string;
  q?: string;
  page?: number;
  limit?: number;
  expiring_within?: number; // e.g. 30 days
}

export interface CustomerContractListResponse {
  items: CustomerContractRecord[];
  meta?: {
    current_page: number;
    per_page: number;
    total_items: number;
    total_pages: number;
  };
}

export interface CustomerContractStats {
  total: number;
  draft: number;
  active: number;
  suspended: number;
  terminated: number;
  expiring_soon: number;
}
