// Master Tariff Class Domain Types (internal/master/tariff_class)

export interface TariffClassRecord {
  id: string;
  code: string;
  name: string;
  description: string;
  multiplier: string;
  coverageType: string;
  isActive: boolean;
}
