import { authFetch, type ApiResponse } from '../client';
import type {
  CustomerContractRecord,
  CustomerContractListParams,
  CustomerContractListResponse,
  CustomerContractStats,
  CreateContractDTO,
  UpdateContractDraftDTO,
  UpdateContractInfoDTO,
  AmendContractTermsDTO,
  RenewContractDTO,
  SetContractPriceSchemasDTO,
  SelectableContract,
} from '../../types/finance/customer_contract';

const LOCAL_STORAGE_KEY = 'hosim_mock_customer_contracts';

// Fallback seed data if backend endpoints are not yet deployed
function getInitialMockContracts(): CustomerContractRecord[] {
  const now = new Date().toISOString();
  return [
    {
      id: 'cc-bpjs-2026',
      customer_id: 'cust-bpjs-001',
      customer_code: 'BPJS-KES',
      customer_name: 'BPJS Kesehatan Kantor Cabang Utama',
      contract_number: 'PKS/2026/RS/BPJS/001',
      name: 'PKS Pelayanan Rujukan Tingkat Lanjutan BPJS 2026',
      effective_date: '2026-01-01',
      expiration_date: '2026-12-31',
      status: 'ACTIVE',
      payment_term_days: 30,
      payment_term_notes: 'Klaim diajukan tanggal 10 tiap bulan; jatuh tempo 30 hari verifikasi',
      credit_limit: 2000000000,
      guarantee_requirement: 'SPPK',
      guarantee_notes: 'Wajib penerbitan SEP online VClaim dan validasi rujukan primer',
      renewed_from_contract_id: null,
      notes: 'Kontrak induk pelayanan JKN-KIS rawat jalan dan rawat inap',
      activated_at: '2026-01-01T00:00:00Z',
      activated_by: 'Direktur Utama',
      status_reason: null,
      price_schemas: [
        {
          id: 'ccps-1',
          customer_contract_id: 'cc-bpjs-2026',
          price_schema_id: 'schema-bpjs-2026',
          price_schema_code: 'SCH-BPJS-2026',
          price_schema_name: 'Skema Tarif INA-CBG & Non-CBG BPJS 2026',
          priority: 1,
          is_active: true,
          created_at: now,
          updated_at: now,
        },
      ],
      created_at: '2025-12-15T08:00:00Z',
      updated_at: now,
      created_by: 'admin_finance',
      updated_by: 'admin_finance',
    },
    {
      id: 'cc-prudential-2026',
      customer_id: 'cust-pru-002',
      customer_code: 'PRU-ASR',
      customer_name: 'PT Prudential Life Assurance Indonesia',
      contract_number: 'PKS/2026/RS/PRU/042',
      name: 'PKS Direct Billing & Asuransi Rawat Inap Eksekutif',
      effective_date: '2026-02-01',
      expiration_date: '2027-01-31',
      status: 'ACTIVE',
      payment_term_days: 45,
      payment_term_notes: 'Net 45 hari dari kwitansi asli dan billing rincian diterima TPA',
      credit_limit: 750000000,
      guarantee_requirement: 'LETTER_OF_GUARANTEE',
      guarantee_notes: 'Wajib Surat Jaminan Awal (Initial GL) dari Third Party Administrator (TPA)',
      renewed_from_contract_id: null,
      notes: 'Khusus kelas rawat VIP, VVIP, dan Eksekutif',
      activated_at: '2026-02-01T08:00:00Z',
      activated_by: 'Direktur Keuangan',
      status_reason: null,
      price_schemas: [
        {
          id: 'ccps-2',
          customer_contract_id: 'cc-prudential-2026',
          price_schema_id: 'schema-corporate-tier-a',
          price_schema_code: 'SCH-CORP-A',
          price_schema_name: 'Skema Asuransi Swasta Gold Tier',
          priority: 1,
          is_active: true,
          created_at: now,
          updated_at: now,
        },
      ],
      created_at: '2026-01-20T10:00:00Z',
      updated_at: now,
      created_by: 'admin_finance',
      updated_by: 'admin_finance',
    },
    {
      id: 'cc-inhealth-2026',
      customer_id: 'cust-inhealth-003',
      customer_code: 'MNDR-INH',
      customer_name: 'PT Asuransi Jiwa Inhealth Indonesia',
      contract_number: 'PKS/2026/RS/INH/015',
      name: 'PKS Pelayanan Managed Care & Indemnity 2026',
      effective_date: '2026-04-01',
      expiration_date: '2027-03-31',
      status: 'DRAFT',
      payment_term_days: 30,
      payment_term_notes: 'Net 30 hari via transfer bank BNI RS',
      credit_limit: 500000000,
      guarantee_requirement: 'SPPK',
      guarantee_notes: 'Menunjukkan e-Card Inhealth Mobile dan verifikasi Eligibilitas Online',
      renewed_from_contract_id: null,
      notes: 'Menunggu review tanda tangan direksi kedua belah pihak',
      activated_at: null,
      activated_by: null,
      status_reason: null,
      price_schemas: [
        {
          id: 'ccps-3',
          customer_contract_id: 'cc-inhealth-2026',
          price_schema_id: 'schema-corporate-tier-a',
          price_schema_code: 'SCH-CORP-A',
          price_schema_name: 'Skema Asuransi Swasta Gold Tier',
          priority: 1,
          is_active: true,
          created_at: now,
          updated_at: now,
        },
      ],
      created_at: '2026-03-10T09:30:00Z',
      updated_at: now,
      created_by: 'staf_kerjasama',
      updated_by: 'staf_kerjasama',
    },
    {
      id: 'cc-telkom-2025',
      customer_id: 'cust-telkom-004',
      customer_code: 'CORP-TLKM',
      customer_name: 'PT Telkom Indonesia (Persero) Tbk',
      contract_number: 'PKS/2025/RS/TLKM/088',
      name: 'PKS MCU Tahunan & Rawat Jalan Karyawan',
      effective_date: '2025-01-01',
      expiration_date: '2025-12-31',
      status: 'SUSPENDED',
      payment_term_days: 14,
      payment_term_notes: 'Net 14 hari',
      credit_limit: 100000000,
      guarantee_requirement: 'NONE',
      guarantee_notes: 'Cukup kartu tanda pengenal pegawai dan NIK',
      renewed_from_contract_id: null,
      notes: 'Ditangguhkan sementara menunggu rekonsiliasi invoice triwulan 4',
      activated_at: '2025-01-01T00:00:00Z',
      activated_by: 'Direktur Keuangan',
      status_reason: 'Audit invoice klaim akhir tahun 2025 masih berlangsung; plafon outstanding pending verifikasi',
      price_schemas: [],
      created_at: '2024-12-10T14:00:00Z',
      updated_at: now,
      created_by: 'admin_finance',
      updated_by: 'admin_finance',
    },
  ];
}

function loadLocalStore(): CustomerContractRecord[] {
  if (typeof window === 'undefined') return getInitialMockContracts();
  try {
    const raw = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed) && parsed.length > 0) return parsed;
    }
  } catch (e) {
    console.warn('Failed to load customer contracts from localStorage', e);
  }
  const initial = getInitialMockContracts();
  saveLocalStore(initial);
  return initial;
}

function saveLocalStore(records: CustomerContractRecord[]): void {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(records));
  } catch (e) {
    console.error('Failed to save customer contracts to localStorage', e);
  }
}

// 1. GET LIST OF CUSTOMER CONTRACTS
export async function getCustomerContracts(
  params?: CustomerContractListParams,
): Promise<CustomerContractListResponse> {
  const q = new URLSearchParams();
  if (params?.q) q.set('q', params.q);
  if (params?.customer_id) q.set('customer_id', params.customer_id);
  if (params?.status) q.set('status', params.status);
  if (params?.valid_on) q.set('valid_on', params.valid_on);
  if (params?.page) q.set('page', String(params.page));
  if (params?.limit) q.set('limit', String(params.limit));
  if (params?.expiring_within) q.set('expiring_within', String(params.expiring_within));

  const queryStr = q.toString() ? `?${q.toString()}` : '';

  try {
    const res = await authFetch<CustomerContractRecord[]>(`/finance/customer-contracts${queryStr}`);
    if (res.data) {
      return {
        items: res.data,
        meta: res.meta,
      };
    }
  } catch (err: any) {
    // Fallback to local store if backend returns 404 or fails
    console.info('[CustomerContract API] Using local mock store:', err?.message);
  }

  // Filter in local store
  let items = loadLocalStore();

  if (params?.customer_id) {
    items = items.filter((c) => c.customer_id === params.customer_id);
  }

  if (params?.status && params.status !== 'ALL') {
    if (params.status === 'EXPIRED') {
      const today = new Date().toISOString().split('T')[0];
      items = items.filter(
        (c) => c.status === 'ACTIVE' && c.expiration_date && c.expiration_date < today,
      );
    } else {
      items = items.filter((c) => c.status === params.status);
    }
  }

  if (params?.valid_on) {
    const d = params.valid_on;
    items = items.filter((c) => {
      if (c.status !== 'ACTIVE') return false;
      if (c.effective_date > d) return false;
      if (c.expiration_date && d > c.expiration_date) return false;
      return true;
    });
  }

  if (params?.expiring_within) {
    const days = params.expiring_within;
    const now = new Date();
    const threshold = new Date(now.getTime() + days * 24 * 60 * 60 * 1000)
      .toISOString()
      .split('T')[0];
    const today = now.toISOString().split('T')[0];
    items = items.filter(
      (c) =>
        c.status === 'ACTIVE' &&
        c.expiration_date &&
        c.expiration_date >= today &&
        c.expiration_date <= threshold,
    );
  }

  if (params?.q) {
    const query = params.q.toLowerCase();
    items = items.filter(
      (c) =>
        c.contract_number.toLowerCase().includes(query) ||
        c.name.toLowerCase().includes(query) ||
        (c.customer_name && c.customer_name.toLowerCase().includes(query)) ||
        (c.customer_code && c.customer_code.toLowerCase().includes(query)) ||
        (c.notes && c.notes.toLowerCase().includes(query)),
    );
  }

  const page = params?.page || 1;
  const limit = params?.limit || 15;
  const total = items.length;
  const totalPages = Math.ceil(total / limit) || 1;
  const paginated = items.slice((page - 1) * limit, page * limit);

  return {
    items: paginated,
    meta: {
      current_page: page,
      per_page: limit,
      total_items: total,
      total_pages: totalPages,
    },
  };
}

// 2. GET CONTRACT BY ID
export async function getCustomerContractById(id: string): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}`);
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const found = items.find((c) => c.id === id);
  if (!found) {
    throw new Error(`Kontrak dengan ID '${id}' tidak ditemukan`);
  }
  return found;
}

// 3. CREATE DRAFT CONTRACT
export async function createCustomerContract(
  dto: CreateContractDTO,
): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>('/finance/customer-contracts', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  // Check unique contract_number
  if (items.some((c) => c.contract_number.toUpperCase() === dto.contract_number.toUpperCase())) {
    throw new Error(`Nomor kontrak '${dto.contract_number}' sudah terdaftar dalam sistem.`);
  }

  const now = new Date().toISOString();
  const newContract: CustomerContractRecord = {
    id: 'cc-' + Date.now(),
    customer_id: dto.customer_id,
    contract_number: dto.contract_number.trim().toUpperCase(),
    name: dto.name.trim(),
    effective_date: dto.effective_date,
    expiration_date: dto.expiration_date || null,
    status: 'DRAFT',
    payment_term_days: dto.payment_term_days,
    payment_term_notes: dto.payment_term_notes || null,
    credit_limit: dto.credit_limit !== undefined ? dto.credit_limit : null,
    guarantee_requirement: dto.guarantee_requirement,
    guarantee_notes: dto.guarantee_notes || null,
    renewed_from_contract_id: null,
    notes: dto.notes || null,
    activated_at: null,
    activated_by: null,
    terminated_at: null,
    terminated_by: null,
    status_reason: null,
    price_schemas: (dto.price_schemas || []).map((s, idx) => ({
      id: 'ccps-' + Date.now() + '-' + idx,
      customer_contract_id: 'cc-' + Date.now(),
      price_schema_id: s.price_schema_id,
      priority: s.priority || idx + 1,
      is_active: s.is_active !== undefined ? s.is_active : true,
      created_at: now,
      updated_at: now,
    })),
    created_at: now,
    updated_at: now,
    created_by: 'Petugas Keuangan',
    updated_by: 'Petugas Keuangan',
  };

  items.unshift(newContract);
  saveLocalStore(items);
  return newContract;
}

// 4. UPDATE CONTRACT (DRAFT OR INFO)
export async function updateCustomerContract(
  id: string,
  dto: UpdateContractDraftDTO | UpdateContractInfoDTO,
): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  const now = new Date().toISOString();

  if (curr.status === 'DRAFT') {
    const draftDTO = dto as UpdateContractDraftDTO;
    items[idx] = {
      ...curr,
      contract_number: draftDTO.contract_number ? draftDTO.contract_number.toUpperCase() : curr.contract_number,
      name: draftDTO.name || curr.name,
      effective_date: draftDTO.effective_date || curr.effective_date,
      expiration_date: draftDTO.expiration_date !== undefined ? draftDTO.expiration_date : curr.expiration_date,
      payment_term_days: draftDTO.payment_term_days ?? curr.payment_term_days,
      payment_term_notes: draftDTO.payment_term_notes !== undefined ? draftDTO.payment_term_notes : curr.payment_term_notes,
      credit_limit: draftDTO.credit_limit !== undefined ? draftDTO.credit_limit : curr.credit_limit,
      guarantee_requirement: draftDTO.guarantee_requirement || curr.guarantee_requirement,
      guarantee_notes: draftDTO.guarantee_notes !== undefined ? draftDTO.guarantee_notes : curr.guarantee_notes,
      notes: draftDTO.notes !== undefined ? draftDTO.notes : curr.notes,
      updated_at: now,
    };
  } else {
    // Non-financial info update for ACTIVE or SUSPENDED
    const infoDTO = dto as UpdateContractInfoDTO;
    items[idx] = {
      ...curr,
      name: infoDTO.name || curr.name,
      notes: infoDTO.notes !== undefined ? infoDTO.notes : curr.notes,
      guarantee_notes: infoDTO.guarantee_notes !== undefined ? infoDTO.guarantee_notes : curr.guarantee_notes,
      payment_term_notes: infoDTO.payment_term_notes !== undefined ? infoDTO.payment_term_notes : curr.payment_term_notes,
      updated_at: now,
    };
  }

  saveLocalStore(items);
  return items[idx];
}

// 5. AMEND CONTRACT TERMS (ACTIVE / SUSPENDED)
export async function amendCustomerContract(
  id: string,
  dto: AmendContractTermsDTO,
): Promise<CustomerContractRecord> {
  if (!dto.reason || !dto.reason.trim()) {
    throw new Error('Alasan perubahan klausul/addendum (reason) wajib diisi untuk jejak audit.');
  }

  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/amend`, {
      method: 'POST',
      body: JSON.stringify(dto),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  if (curr.status !== 'ACTIVE' && curr.status !== 'SUSPENDED') {
    throw new Error('Amandemen hanya berlaku untuk kontrak berstatus ACTIVE atau SUSPENDED.');
  }

  const now = new Date().toISOString();
  items[idx] = {
    ...curr,
    payment_term_days: dto.payment_term_days,
    credit_limit: dto.credit_limit !== undefined ? dto.credit_limit : curr.credit_limit,
    expiration_date: dto.expiration_date !== undefined ? dto.expiration_date : curr.expiration_date,
    guarantee_requirement: dto.guarantee_requirement,
    notes: curr.notes ? `${curr.notes}\n[Amandemen ${now.split('T')[0]}]: ${dto.reason}` : `[Amandemen]: ${dto.reason}`,
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 6. ACTIVATE CONTRACT
export async function activateCustomerContract(id: string): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/activate`, {
      method: 'POST',
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  if (curr.status !== 'DRAFT') {
    throw new Error('Aktivasi hanya dapat dilakukan dari status DRAFT.');
  }

  if (!curr.effective_date) {
    throw new Error('Tanggal mulai berlaku (effective_date) wajib diisi sebelum mengaktifkan kontrak.');
  }

  const activeSchemas = (curr.price_schemas || []).filter((s) => s.is_active);
  if (activeSchemas.length === 0) {
    throw new Error('Kontrak wajib memiliki minimal 1 price schema aktif sebelum dapat diaktifkan.');
  }

  const now = new Date().toISOString();
  items[idx] = {
    ...curr,
    status: 'ACTIVE',
    activated_at: now,
    activated_by: 'Administrator Keuangan',
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 7. REACTIVATE CONTRACT
export async function reactivateCustomerContract(id: string): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/reactivate`, {
      method: 'POST',
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  if (curr.status !== 'SUSPENDED') {
    throw new Error('Reaktivasi hanya dapat dilakukan dari status SUSPENDED.');
  }

  const now = new Date().toISOString();
  items[idx] = {
    ...curr,
    status: 'ACTIVE',
    status_reason: null,
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 8. SUSPEND CONTRACT
export async function suspendCustomerContract(
  id: string,
  reason: string,
): Promise<CustomerContractRecord> {
  if (!reason || !reason.trim()) {
    throw new Error('Alasan penangguhan (suspend reason) wajib diisi.');
  }

  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/suspend`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  if (curr.status !== 'ACTIVE') {
    throw new Error('Penangguhan (suspend) hanya dapat dilakukan pada kontrak yang berstatus ACTIVE.');
  }

  const now = new Date().toISOString();
  items[idx] = {
    ...curr,
    status: 'SUSPENDED',
    status_reason: reason.trim(),
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 9. TERMINATE CONTRACT
export async function terminateCustomerContract(
  id: string,
  reason: string,
): Promise<CustomerContractRecord> {
  if (!reason || !reason.trim()) {
    throw new Error('Alasan pemutusan kontrak (terminate reason) wajib diisi.');
  }

  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/terminate`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  if (curr.status !== 'ACTIVE' && curr.status !== 'SUSPENDED') {
    throw new Error('Terminasi hanya dapat dilakukan dari status ACTIVE atau SUSPENDED.');
  }

  const now = new Date().toISOString();
  items[idx] = {
    ...curr,
    status: 'TERMINATED',
    status_reason: reason.trim(),
    terminated_at: now,
    terminated_by: 'Direksi Rumah Sakit',
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 10. RENEW CONTRACT
export async function renewCustomerContract(
  id: string,
  dto: RenewContractDTO,
): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/renew`, {
      method: 'POST',
      body: JSON.stringify(dto),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const found = items.find((c) => c.id === id);
  if (!found) throw new Error('Kontrak lama tidak ditemukan');

  if (found.status !== 'ACTIVE' && found.status !== 'TERMINATED') {
    throw new Error('Perpanjangan (renew) hanya dapat dilakukan dari kontrak ACTIVE atau TERMINATED.');
  }

  const now = new Date().toISOString();
  const newContract: CustomerContractRecord = {
    id: 'cc-' + Date.now(),
    customer_id: found.customer_id,
    customer_name: found.customer_name,
    customer_code: found.customer_code,
    contract_number: dto.contract_number.trim().toUpperCase(),
    name: dto.name.trim(),
    effective_date: dto.effective_date,
    expiration_date: dto.expiration_date || null,
    status: 'DRAFT',
    payment_term_days: found.payment_term_days,
    payment_term_notes: found.payment_term_notes,
    credit_limit: found.credit_limit,
    guarantee_requirement: found.guarantee_requirement,
    guarantee_notes: found.guarantee_notes,
    renewed_from_contract_id: found.id,
    renewed_from_contract_number: found.contract_number,
    notes: dto.notes || `Diperpanjang dari kontrak PKS: ${found.contract_number}`,
    activated_at: null,
    activated_by: null,
    terminated_at: null,
    terminated_by: null,
    status_reason: null,
    price_schemas: (found.price_schemas || []).map((s) => ({
      ...s,
      id: 'ccps-' + Date.now() + '-' + Math.random().toString(36).substring(2, 6),
      customer_contract_id: 'cc-' + Date.now(),
      created_at: now,
      updated_at: now,
    })),
    created_at: now,
    updated_at: now,
    created_by: 'Petugas Keuangan',
    updated_by: 'Petugas Keuangan',
  };

  items.unshift(newContract);
  saveLocalStore(items);
  return newContract;
}

// 11. DELETE DRAFT CONTRACT
export async function deleteCustomerContract(id: string): Promise<void> {
  try {
    await authFetch<void>(`/finance/customer-contracts/${id}`, {
      method: 'DELETE',
    });
    return;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const found = items.find((c) => c.id === id);
  if (!found) throw new Error('Kontrak tidak ditemukan');

  if (found.status !== 'DRAFT') {
    throw new Error('Hanya kontrak berstatus DRAFT yang diizinkan untuk dihapus dari sistem.');
  }

  const filtered = items.filter((c) => c.id !== id);
  saveLocalStore(filtered);
}

// 12. SET PRICE SCHEMAS & PRIORITIES
export async function setContractPriceSchemas(
  id: string,
  schemas: Array<{ price_schema_id: string; priority: number; is_active?: boolean; price_schema_name?: string; price_schema_code?: string }>,
): Promise<CustomerContractRecord> {
  try {
    const res = await authFetch<CustomerContractRecord>(`/finance/customer-contracts/${id}/price-schemas`, {
      method: 'PUT',
      body: JSON.stringify({ schemas }),
    });
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  const idx = items.findIndex((c) => c.id === id);
  if (idx === -1) throw new Error('Kontrak tidak ditemukan');

  const curr = items[idx];
  const now = new Date().toISOString();

  // If active, ensure at least 1 schema is active
  const activeCount = schemas.filter((s) => s.is_active !== false).length;
  if (curr.status === 'ACTIVE' && activeCount === 0) {
    throw new Error('Kontrak aktif wajib memiliki minimal 1 price schema aktif.');
  }

  items[idx] = {
    ...curr,
    price_schemas: schemas.map((s, i) => ({
      id: 'ccps-' + Date.now() + '-' + i,
      customer_contract_id: id,
      price_schema_id: s.price_schema_id,
      price_schema_name: s.price_schema_name,
      price_schema_code: s.price_schema_code,
      priority: s.priority,
      is_active: s.is_active !== false,
      created_at: now,
      updated_at: now,
    })),
    updated_at: now,
  };

  saveLocalStore(items);
  return items[idx];
}

// 13. SELECTABLE CONTRACTS QUERY (FOR REGISTRATION/BILLING)
export async function getSelectableContracts(
  customerId: string,
  onDate?: string,
): Promise<SelectableContract[]> {
  const d = onDate || new Date().toISOString().split('T')[0];
  try {
    const res = await authFetch<SelectableContract[]>(
      `/finance/customers/${customerId}/selectable-contracts?on=${d}`,
    );
    if (res.data) return res.data;
  } catch (err) {
    // Local fallback
  }

  const items = loadLocalStore();
  return items
    .filter((c) => {
      if (c.customer_id !== customerId) return false;
      if (c.status !== 'ACTIVE') return false;
      if (c.effective_date > d) return false;
      if (c.expiration_date && d > c.expiration_date) return false;
      return true;
    })
    .map((c) => ({
      contract_id: c.id,
      contract_number: c.contract_number,
      name: c.name,
      guarantee_requirement: c.guarantee_requirement,
      price_schemas: (c.price_schemas || [])
        .filter((s) => s.is_active)
        .sort((a, b) => a.priority - b.priority)
        .map((s) => ({
          price_schema_id: s.price_schema_id,
          price_schema_code: s.price_schema_code || 'SCHEMA',
          price_schema_name: s.price_schema_name || 'Skema Tarif',
          priority: s.priority,
        })),
    }));
}
