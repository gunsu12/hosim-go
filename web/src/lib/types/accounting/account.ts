// Chart of Accounts (COA) / Bagan Akun Types for HOSIM Accounting Domain
// Matches backend definitions in pkg/enums/account.go & internal/accounting/account/dto.go

export type AccountType = 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE';

export type AccountPosition = 'DEBIT' | 'CREDIT';

export interface AccountRecord {
  id: string;
  code: string;
  parent_id?: string | null;
  parent_code?: string | null;
  name: string;
  description?: string | null;
  type: AccountType;
  position: AccountPosition;
  account_level: number;
  is_postable: boolean;
  is_treasury_account: boolean;
  bank_name?: string | null;
  bank_account_number?: string | null;
  currency: string;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface AccountTreeNode extends AccountRecord {
  children?: AccountTreeNode[];
}

export interface CreateAccountDTO {
  code: string;
  name: string;
  parent_id?: string | null;
  type: AccountType;
  position?: AccountPosition | null;
  description?: string | null;
  is_postable?: boolean;
  is_treasury_account: boolean;
  bank_name?: string | null;
  bank_account_number?: string | null;
  currency?: string;
  is_active?: boolean;
}

export interface UpdateAccountDTO {
  name: string;
  position?: AccountPosition;
  description?: string | null;
  is_treasury_account: boolean;
  bank_name?: string | null;
  bank_account_number?: string | null;
  currency?: string;
  is_active?: boolean;
}

export interface AccountListParams {
  page?: number;
  limit?: number;
  search?: string;
  type?: AccountType;
  position?: AccountPosition;
  parent_id?: string;
  is_postable?: boolean;
  is_treasury?: boolean;
  is_active?: boolean;
  level?: number;
}

export interface AccountListResponse {
  items: AccountRecord[];
  total: number;
  page: number;
  limit: number;
}
