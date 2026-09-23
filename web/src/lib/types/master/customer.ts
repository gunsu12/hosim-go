// Master Customer / Penjamin Domain Types (internal/finance/customer)

export interface CustomerRecord {
  id: string;
  code: string;
  name: string;
  typeName: string;
  phone: string;
  requireCard: boolean;
  isActive: boolean;
}
