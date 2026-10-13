// Master Practitioner Domain Types (internal/practitioner)

export interface ProfessionRecord {
  id: string;
  code: string;
  name: string;
}

export interface SpecialtyRecord {
  id: string;
  code: string;
  name: string;
}

export interface PractitionerRecord {
  id: string;
  nik?: string;
  nip?: string;
  name: string;
  gender?: string;
  sip?: string;
  sip_expiry_date?: string;
  str?: string;
  profession_id?: string;
  profession?: ProfessionRecord | string;
  specialty_id?: string;
  specialty?: SpecialtyRecord | string;
  phone?: string;
  email?: string;
  ihs_practitioner_id?: string;
  ihs_practitioner_name?: string;
  is_active: boolean;
  // UI legacy aliases
  isActive?: boolean;
  ihsId?: string;
  ihs_number?: string;
}
