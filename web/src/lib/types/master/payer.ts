// Master Payer Domain Types (internal/master/payer)

export interface PayerRecord {
  id: string;
  code: string;
  name: string;
  typeName: string;
  phone: string;
  requireCard: boolean;
  isActive: boolean;
}
