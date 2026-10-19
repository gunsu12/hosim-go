<script lang="ts">
  import { onMount } from 'svelte';
  import { FileText, Check, AlertCircle, Building2, Calendar, ShieldCheck } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select, { type SelectOption } from '$lib/components/m3/M3Select.svelte';
  import M3DatePicker from '$lib/components/m3/M3DatePicker.svelte';
  import type {
    CustomerContractRecord,
    CreateContractDTO,
    UpdateContractDraftDTO,
    UpdateContractInfoDTO,
    GuaranteeRequirement,
  } from '$lib/types/finance/customer_contract';
  import type { CustomerRecord } from '$lib/types/finance/customer';
  import type { PriceSchemaRecord } from '$lib/types/finance/price_schema';
  import { getCustomers } from '$lib/api/finance/customer';
  import { getPriceSchemas } from '$lib/api/finance/price_schema';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
    contract?: CustomerContractRecord | null;
    isSubmitting?: boolean;
    onClose: () => void;
    onSubmit: (
      dto: CreateContractDTO | UpdateContractDraftDTO | UpdateContractInfoDTO,
    ) => Promise<void> | void;
  }

  let {
    isOpen,
    mode = 'CREATE',
    contract = null,
    isSubmitting = false,
    onClose,
    onSubmit,
  }: Props = $props();

  // Reference data
  let customers = $state<CustomerRecord[]>([]);
  let priceSchemas = $state<PriceSchemaRecord[]>([]);
  let isLoadingRefs = $state(false);

  // Form states
  let customerId = $state('');
  let contractNumber = $state('');
  let name = $state('');
  let effectiveDate = $state('');
  let hasExpiration = $state(true);
  let expirationDate = $state('');
  let paymentTermDays = $state<number>(30);
  let paymentTermNotes = $state('');
  let hasCreditLimit = $state(false);
  let creditLimitAmount = $state<number | string>(0);
  let guaranteeRequirement = $state<GuaranteeRequirement>('NONE');
  let guaranteeNotes = $state('');
  let notes = $state('');
  let initialPriceSchemaId = $state('');
  let localError = $state<string | null>(null);

  const guaranteeOptions: SelectOption[] = [
    { value: 'NONE', label: 'Tanpa Jaminan Khusus (NONE)', description: 'Pasien penjamin cukup kartu/identitas standar' },
    { value: 'SPPK', label: 'Wajib SPPK / SEP Online', description: 'Surat Penetapan/Perjanjian Pelayanan Kesehatan (BPJS/Asuransi)' },
    { value: 'LETTER_OF_GUARANTEE', label: 'Wajib Letter of Guarantee (GL)', description: 'Surat Jaminan Awal Resmi dari Asuransi/TPA Perusahaan' },
  ];

  async function loadReferenceData() {
    isLoadingRefs = true;
    try {
      const [custRes, schemaRes] = await Promise.all([
        getCustomers({ limit: 100, is_active: true }),
        getPriceSchemas({ limit: 100, is_active: true }),
      ]);
      customers = custRes.data || [];
      priceSchemas = schemaRes.items || [];
      if (!customerId && customers.length > 0 && mode === 'CREATE') {
        customerId = customers[0].id;
      }
      if (!initialPriceSchemaId && priceSchemas.length > 0 && mode === 'CREATE') {
        initialPriceSchemaId = priceSchemas[0].id;
      }
    } catch (e) {
      console.warn('Gagal memuat referensi customer & skema tarif', e);
    } finally {
      isLoadingRefs = false;
    }
  }

  const customerOptions = $derived<SelectOption[]>(
    customers.map((c) => ({
      value: c.id,
      label: `${c.code} • ${c.name}`,
      description: c.customer_type?.name || 'Penjamin/Debitur',
    })),
  );

  const priceSchemaOptions = $derived<SelectOption[]>(
    priceSchemas.map((s) => ({
      value: s.id,
      label: `${s.code} • ${s.name}`,
      description: s.description || 'Skema tarif aktif',
    })),
  );

  $effect(() => {
    if (isOpen) {
      localError = null;
      loadReferenceData();

      if (mode === 'EDIT' && contract) {
        customerId = contract.customer_id;
        contractNumber = contract.contract_number;
        name = contract.name;
        effectiveDate = contract.effective_date;

        if (contract.expiration_date) {
          hasExpiration = true;
          expirationDate = contract.expiration_date;
        } else {
          hasExpiration = false;
          expirationDate = '';
        }

        paymentTermDays = contract.payment_term_days;
        paymentTermNotes = contract.payment_term_notes || '';

        if (contract.credit_limit !== null && contract.credit_limit !== undefined) {
          hasCreditLimit = true;
          creditLimitAmount = contract.credit_limit;
        } else {
          hasCreditLimit = false;
          creditLimitAmount = 0;
        }

        guaranteeRequirement = contract.guarantee_requirement;
        guaranteeNotes = contract.guarantee_notes || '';
        notes = contract.notes || '';
        initialPriceSchemaId = '';
      } else {
        // CREATE default values
        customerId = customers.length > 0 ? customers[0].id : '';
        const currentYear = new Date().getFullYear();
        contractNumber = `PKS/${currentYear}/RS/`;
        name = '';
        effectiveDate = new Date().toISOString().split('T')[0];
        hasExpiration = true;
        const expDate = new Date();
        expDate.setFullYear(expDate.getFullYear() + 1);
        expirationDate = expDate.toISOString().split('T')[0];
        paymentTermDays = 30;
        paymentTermNotes = '';
        hasCreditLimit = false;
        creditLimitAmount = 0;
        guaranteeRequirement = 'SPPK';
        guaranteeNotes = '';
        notes = '';
        initialPriceSchemaId = priceSchemas.length > 0 ? priceSchemas[0].id : '';
      }
    }
  });

  let isDraftMode = $derived(mode === 'CREATE' || contract?.status === 'DRAFT');

  function handleFormSubmit(e: SubmitEvent) {
    e.preventDefault();
    localError = null;

    if (mode === 'CREATE' && !customerId) {
      localError = 'Penjamin / debitur wajib dipilih';
      return;
    }
    if (!contractNumber.trim()) {
      localError = 'Nomor kontrak (PKS) wajib diisi';
      return;
    }
    if (!name.trim()) {
      localError = 'Nama kontrak wajib diisi';
      return;
    }
    if (isDraftMode && !effectiveDate) {
      localError = 'Tanggal mulai berlaku wajib ditentukan';
      return;
    }

    if (mode === 'CREATE') {
      const dto: CreateContractDTO = {
        customer_id: customerId,
        contract_number: contractNumber.trim().toUpperCase(),
        name: name.trim(),
        effective_date: effectiveDate,
        expiration_date: hasExpiration && expirationDate ? expirationDate : null,
        payment_term_days: Number(paymentTermDays) || 30,
        payment_term_notes: paymentTermNotes.trim() || undefined,
        credit_limit: hasCreditLimit ? Number(creditLimitAmount) || 0 : null,
        guarantee_requirement: guaranteeRequirement,
        guarantee_notes: guaranteeNotes.trim() || undefined,
        notes: notes.trim() || undefined,
        price_schemas: initialPriceSchemaId
          ? [{ price_schema_id: initialPriceSchemaId, priority: 1, is_active: true }]
          : [],
      };
      onSubmit(dto);
    } else if (contract?.status === 'DRAFT') {
      const dto: UpdateContractDraftDTO = {
        contract_number: contractNumber.trim().toUpperCase(),
        name: name.trim(),
        effective_date: effectiveDate,
        expiration_date: hasExpiration && expirationDate ? expirationDate : null,
        payment_term_days: Number(paymentTermDays) || 30,
        payment_term_notes: paymentTermNotes.trim() || undefined,
        credit_limit: hasCreditLimit ? Number(creditLimitAmount) || 0 : null,
        guarantee_requirement: guaranteeRequirement,
        guarantee_notes: guaranteeNotes.trim() || undefined,
        notes: notes.trim() || undefined,
      };
      onSubmit(dto);
    } else {
      // Non-financial info update for ACTIVE/SUSPENDED
      const dto: UpdateContractInfoDTO = {
        name: name.trim(),
        notes: notes.trim() || undefined,
        guarantee_notes: guaranteeNotes.trim() || undefined,
        payment_term_notes: paymentTermNotes.trim() || undefined,
      };
      onSubmit(dto);
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'CREATE' ? 'Buat Kontrak Kerjasama (PKS) Baru' : 'Edit Kontrak Kerjasama (PKS)'}
  subtitle={mode === 'CREATE'
    ? 'Daftarkan perjanjian kerjasama penjamin baru dalam status DRAFT'
    : contract
      ? `${contract.contract_number} • Status: ${contract.status}`
      : ''}
  icon={FileText}
  iconClass="bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]"
  size="3xl"
  onsubmit={handleFormSubmit}
  {onClose}
>
  <div class="space-y-4 text-xs">
    <!-- Local error banner -->
    {#if localError}
      <div
        class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 flex items-start gap-2.5 shadow-2xs"
      >
        <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
        <span class="font-medium leading-relaxed">{localError}</span>
      </div>
    {/if}

    <!-- SECTION 1: IDENTITAS PENJAMIN & KONTRAK -->
    <div
      class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-3.5"
    >
      <div class="flex items-center gap-2 pb-1 border-b border-[#e1e5ea] dark:border-[#353738]">
        <Building2 class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] uppercase tracking-wider text-[11px]">
          1. Identitas Penjamin & Dokumen PKS
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <!-- Debitur / Penjamin -->
        <div>
          <label
            for="contract-customer-id"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Debitur / Penjamin <span class="text-rose-500">*</span>
          </label>
          {#if mode === 'CREATE'}
            <M3Select
              id="contract-customer-id"
              options={customerOptions}
              bind:value={customerId}
              searchable
              placeholder="Pilih debitur/penjamin..."
              disabled={isLoadingRefs}
            />
          {:else}
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center font-semibold truncate"
            >
              {contract?.customer_name || contract?.customer?.name || 'Penjamin Terkait'}
            </div>
            <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-0.5 block">
              Entitas penjamin tidak dapat diubah setelah kontrak dibuat.
            </span>
          {/if}
        </div>

        <!-- Nomor Kontrak -->
        <div>
          <label
            for="contract-number-input"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Nomor Kontrak (PKS) <span class="text-rose-500">*</span>
          </label>
          <input
            id="contract-number-input"
            type="text"
            required
            disabled={!isDraftMode}
            bind:value={contractNumber}
            placeholder="cth: PKS/2026/RS/BPJS/001"
            class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden disabled:bg-gray-100 dark:disabled:bg-gray-800"
          />
        </div>
      </div>

      <!-- Nama Kontrak -->
      <div>
        <label
          for="contract-name-input"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Nama / Judul Kontrak <span class="text-rose-500">*</span>
        </label>
        <input
          id="contract-name-input"
          type="text"
          required
          bind:value={name}
          placeholder="cth: PKS Rawat Inap & Rawat Jalan Eksekutif 2026"
          class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden"
        />
        <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-0.5 block">
          Label yang jelas untuk memudahkan petugas saat memilih kontrak pada form registrasi.
        </span>
      </div>
    </div>

    <!-- SECTION 2: PERIODE & TERMIN PEMBAYARAN -->
    <div
      class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-3.5"
    >
      <div class="flex items-center gap-2 pb-1 border-b border-[#e1e5ea] dark:border-[#353738]">
        <Calendar class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] uppercase tracking-wider text-[11px]">
          2. Masa Berlaku & Termin Pembayaran
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <!-- Tanggal Mulai -->
        <div>
          {#if isDraftMode}
            <M3DatePicker
              label="Tanggal Mulai Berlaku"
              required
              bind:value={effectiveDate}
              placeholder="Pilih tanggal mulai"
            />
          {:else}
            <span class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
              Tanggal Mulai Berlaku
            </span>
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center font-mono font-semibold"
            >
              {contract?.effective_date}
            </div>
            <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-0.5 block">
              Tanggal mulai tidak dapat diubah setelah kontrak berstatus aktif.
            </span>
          {/if}
        </div>

        <!-- Tanggal Berakhir -->
        <div>
          {#if isDraftMode}
            <div class="flex items-center justify-between mb-1">
              <label for="form-expiry-date-input" class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                Tanggal Berakhir
              </label>
              <label class="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={hasExpiration}
                  class="rounded text-[#0b57d0] focus:ring-0"
                />
                <span class="text-[10px] text-[#747775] dark:text-[#8e918f]">Memiliki Batas</span>
              </label>
            </div>
            {#if hasExpiration}
              <M3DatePicker
                id="form-expiry-date-input"
                bind:value={expirationDate}
                min={effectiveDate}
                placeholder="Pilih tanggal berakhir"
              />
            {:else}
              <div
                class="h-10 px-3.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 flex items-center font-medium"
              >
                Tanpa batas waktu (Otomatis Diperpanjang)
              </div>
            {/if}
          {:else}
            <span class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
              Tanggal Berakhir
            </span>
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center font-mono font-semibold"
            >
              {contract?.expiration_date || 'Tanpa batas waktu'}
            </div>
            <span class="text-[10px] text-amber-600 dark:text-amber-400 mt-0.5 block">
              Untuk mengubah tanggal berakhir kontrak aktif, gunakan menu <strong>Amandemen</strong>.
            </span>
          {/if}
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <!-- Termin Hari -->
        <div>
          <label
            for="contract-payment-term-days"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Termin Pembayaran (Hari) <span class="text-rose-500">*</span>
          </label>
          <div class="relative flex items-center">
            <input
              id="contract-payment-term-days"
              type="number"
              min="1"
              max="365"
              required
              disabled={!isDraftMode}
              bind:value={paymentTermDays}
              class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden disabled:bg-gray-100 dark:disabled:bg-gray-800"
            />
            <span class="absolute right-3 text-[#747775] dark:text-[#8e918f] font-medium text-[11px]">
              Net {paymentTermDays}
            </span>
          </div>
        </div>

        <!-- Catatan Termin -->
        <div>
          <label
            for="contract-payment-term-notes"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Catatan Termin Pembayaran
          </label>
          <input
            id="contract-payment-term-notes"
            type="text"
            bind:value={paymentTermNotes}
            placeholder="cth: Pengajuan klaim tgl 10 tiap bulan"
            class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden"
          />
        </div>
      </div>
    </div>

    <!-- SECTION 3: KREDIT & SYARAT JAMINAN -->
    <div
      class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-3.5"
    >
      <div class="flex items-center gap-2 pb-1 border-b border-[#e1e5ea] dark:border-[#353738]">
        <ShieldCheck class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] uppercase tracking-wider text-[11px]">
          3. Batas Kredit & Syarat Jaminan Kunjungan
        </span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <!-- Batas Kredit -->
        <div>
          <div class="flex items-center justify-between mb-1">
            <label for="form-limit-toggle" class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
              Batas Kredit (Plafon)
            </label>
            {#if isDraftMode}
              <label class="flex items-center gap-1.5 cursor-pointer">
                <input
                  id="form-limit-toggle"
                  type="checkbox"
                  bind:checked={hasCreditLimit}
                  class="rounded text-[#0b57d0] focus:ring-0"
                />
                <span class="text-[10px] text-[#747775] dark:text-[#8e918f]">Batasi Plafon</span>
              </label>
            {/if}
          </div>

          {#if hasCreditLimit}
            <div class="relative flex items-center">
              <span class="absolute left-3.5 text-[#747775] dark:text-[#8e918f] font-semibold text-xs">
                Rp
              </span>
              <input
                type="number"
                min="0"
                step="1000000"
                disabled={!isDraftMode}
                bind:value={creditLimitAmount}
                placeholder="cth: 500000000"
                class="w-full h-10 pl-10 pr-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden disabled:bg-gray-100 dark:disabled:bg-gray-800"
              />
            </div>
            <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-0.5 block">
              Plafon: Rp {(Number(creditLimitAmount) || 0).toLocaleString('id-ID')}
            </span>
          {:else}
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#444746] dark:text-[#c4c7c5] flex items-center font-medium"
            >
              Tanpa batas kredit (Unlimited)
            </div>
          {/if}
        </div>

        <!-- Syarat Jaminan -->
        <div>
          <label
            for="contract-guarantee-req"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Syarat Jaminan Per Kunjungan <span class="text-rose-500">*</span>
          </label>
          <M3Select
            id="contract-guarantee-req"
            options={guaranteeOptions}
            bind:value={guaranteeRequirement}
            disabled={!isDraftMode}
            placeholder="Pilih syarat jaminan..."
          />
        </div>
      </div>

      <!-- Catatan Jaminan -->
      <div>
        <label
          for="contract-guarantee-notes"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Petunjuk / Catatan Syarat Jaminan
        </label>
        <input
          id="contract-guarantee-notes"
          type="text"
          bind:value={guaranteeNotes}
          placeholder="cth: Wajib melampirkan fotokopi kartu peserta dan surat rujukan FKTP aktif"
          class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden"
        />
      </div>

      <!-- CREATE ONLY: Initial Price Schema -->
      {#if mode === 'CREATE'}
        <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
          <label
            for="contract-initial-schema"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Skema Tarif Awal (Price Schema)
          </label>
          <M3Select
            id="contract-initial-schema"
            options={priceSchemaOptions}
            bind:value={initialPriceSchemaId}
            searchable
            placeholder="Pilih skema tarif utama untuk kontrak ini..."
            disabled={isLoadingRefs}
          />
          <p class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-1">
            Skema ini akan ditautkan sebagai prioritas 1 (default registrasi). Skema tambahan dapat diatur setelah kontrak tersimpan.
          </p>
        </div>
      {/if}

      <!-- Catatan Umum -->
      <div>
        <label
          for="contract-general-notes"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Catatan Tambahan & Keterangan Khusus
        </label>
        <textarea
          id="contract-general-notes"
          rows="2"
          bind:value={notes}
          placeholder="Catatan internal pengelola kontrak..."
          class="w-full p-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden"
        ></textarea>
      </div>
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <M3Button
      variant="filled"
      type="submit"
      disabled={isSubmitting || !contractNumber.trim() || !name.trim()}
    >
      {#if isSubmitting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Menyimpan...</span>
      {:else}
        <Check class="w-4 h-4 mr-1" />
        <span>{mode === 'CREATE' ? 'Simpan Kontrak Draft' : 'Simpan Perubahan'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
