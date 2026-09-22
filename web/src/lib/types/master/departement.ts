// Master Departement Domain Types (internal/master/departement)

export interface DepartementRecord {
  id: string;
  code: string;
  name: string;
  type: string;
  phone: string;
  ihsOrgId: string;
  isActive: boolean;
}
