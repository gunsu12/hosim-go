// Master Service Unit Domain Types (internal/master/service_unit)

export interface ServiceUnitRecord {
  id: string;
  code: string;
  name: string;
  departmentName: string;
  phone: string;
  ihsLocationId: string;
  isRegistrationTarget: boolean;
  isActive: boolean;
}
