// Master Customer / Penjamin Domain Types (internal/finance/customer)

export interface CustomerTypeRecord {
  id: string;
  code: string;
  name: string;
  created_at?: string;
  updated_at?: string;
}

export interface CustomerRecord {
  id: string;
  code: string;
  name: string;
  address?: string | null;
  phone?: string | null;
  email?: string | null;
  website?: string | null;
  contact_person?: string | null;
  require_card?: boolean | null;
  description?: string | null;
  customer_type_id?: string | null;
  customer_type?: CustomerTypeRecord | null;
  is_active: boolean;
  is_immutable?: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface CreateCustomerDTO {
  code: string;
  name: string;
  address?: string | null;
  phone?: string | null;
  email?: string | null;
  website?: string | null;
  contact_person?: string | null;
  require_card?: boolean;
  description?: string | null;
  customer_type_id?: string | null;
  is_active?: boolean;
}

export interface UpdateCustomerDTO {
  code: string;
  name: string;
  address?: string | null;
  phone?: string | null;
  email?: string | null;
  website?: string | null;
  contact_person?: string | null;
  require_card?: boolean;
  description?: string | null;
  customer_type_id?: string | null;
  is_active?: boolean;
}

export interface CustomerListParams {
  search?: string;
  customer_type_id?: string;
  page?: number;
  limit?: number;
  is_active?: boolean;
}
