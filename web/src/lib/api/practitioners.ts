import { authFetch, type PaginatedData } from './client';
import type { PractitionerRecord, ProfessionRecord, SpecialtyRecord } from '../types/master/practitioner';

export async function getPractitioners(params?: {
  search?: string;
  page?: number;
  limit?: number;
}): Promise<PaginatedData<PractitionerRecord>> {
  const q = new URLSearchParams();
  if (params?.search) q.set('search', params.search);
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  const queryStr = q.toString() ? `?${q.toString()}` : '';
  const res = await authFetch<PaginatedData<PractitionerRecord>>(`/practitioners${queryStr}`);
  return res.data;
}

export async function getPractitionerById(id: string): Promise<PractitionerRecord> {
  const res = await authFetch<PractitionerRecord>(`/practitioners/${id}`);
  return res.data;
}

export async function getProfessions(): Promise<ProfessionRecord[]> {
  const res = await authFetch<ProfessionRecord[]>('/professions');
  return res.data;
}

export async function getSpecialties(): Promise<SpecialtyRecord[]> {
  const res = await authFetch<SpecialtyRecord[]>('/specialties');
  return res.data;
}
