// Master Practitioner Domain Types (internal/master/practitioner)

export interface PractitionerRecord {
  id: string;
  name: string;
  sip: string;
  str: string;
  specialty: string;
  profession: string;
  phone: string;
  ihsId: string;
  isActive: boolean;
}
