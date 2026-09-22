// Master Referal Domain Types (internal/master/referal)

export type ReferalType = 'PUSKESMAS' | 'KLINIK_PRATAMA' | 'RS_TIPE_C' | 'RS_TIPE_B' | 'DOKTER_PRAKTEK';

export interface ReferalRecord {
  id: string;
  code: string;
  name: string;
  type: ReferalType;
  typeLabel: string;
  address: string;
  phone: string;
  email: string;
  status: 'ACTIVE' | 'INACTIVE';
}
