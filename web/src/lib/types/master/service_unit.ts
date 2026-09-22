// Master Service Unit Domain Types (internal/master/service_unit)
import type { DepartementRecord } from './departement';

export interface ServiceUnitRecord {
  id: string;
  code: string;
  name: string;
  departement_id?: string;
  departement?: DepartementRecord;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  ihs_location_id?: string;
  is_registration_target: boolean;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface CreateServiceUnitDTO {
  code: string;
  name: string;
  departement_id?: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  ihs_location_id?: string;
  is_registration_target?: boolean;
  is_active?: boolean;
}

export interface UpdateServiceUnitDTO {
  code: string;
  name: string;
  departement_id?: string;
  address?: string;
  phone?: string;
  email?: string;
  website?: string;
  description?: string;
  ihs_location_id?: string;
  is_registration_target?: boolean;
  is_active?: boolean;
}
