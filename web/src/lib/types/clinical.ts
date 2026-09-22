// Clinical, EHR, & Body Diagram Domain Types

import type { PatientVitals } from './master/patient';

export type FindingCategory = 'pain' | 'injury' | 'fracture' | 'edema' | 'mass';

export interface BodyFinding {
  id: number;
  x: number;
  y: number;
  view: 'front' | 'back';
  category: FindingCategory;
  severity: number;
  note: string;
  createdAt: string;
}

export interface ClinicalNotes {
  chiefComplaint: string;
  anamnesis: string;
  diagnosis: string;
  disposition: string;
}

export interface PhysicalExamPayload {
  patient_mrn: string;
  practitioner_id: number;
  encounter_type: string;
  body_diagram_findings: BodyFinding[];
  clinical_notes: ClinicalNotes;
}
