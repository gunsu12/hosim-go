// Master Patient Domain Types (internal/master/patient)

export interface PatientVitals {
  bp: string;
  hr: string;
  rr: string;
  temp: string;
  spo2: string;
}

export interface Patient {
  name: string;
  mrn: string;
  nik: string;
  birthDate: string;
  gender: string;
  bloodType: string;
  payer: string;
  allergies: string[];
  vitals: PatientVitals;
}

export interface PatientRecord {
  id: number;
  mrn: string;
  nik: string;
  name: string;
  gender: 'L' | 'P';
  birthDate: string;
  phone: string;
  bloodType: string;
  payer: string;
  status: 'ACTIVE' | 'INACTIVE';
}
