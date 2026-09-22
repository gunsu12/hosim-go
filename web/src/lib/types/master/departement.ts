// Master Departement Domain Types (internal/master/departement)

export type DepartementType =
  | 'emergency'
  | 'outpatient'
  | 'inpatient'
  | 'diagnostic'
  | 'medical_checkup'
  | 'other';

export interface DepartementRecord {
  id: string;
  code: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  departement_type: DepartementType;
  ihs_organization_id?: string;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface CreateDepartementDTO {
  code: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  departement_type: DepartementType;
  ihs_organization_id?: string;
  is_active?: boolean;
}

export interface UpdateDepartementDTO {
  code: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  departement_type: DepartementType;
  ihs_organization_id?: string;
  is_active?: boolean;
}
