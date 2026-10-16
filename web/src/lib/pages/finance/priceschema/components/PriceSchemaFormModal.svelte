<script lang="ts">
  import { onMount } from 'svelte';
  import { X, Check, AlertCircle, Sparkles, BookOpen } from '@lucide/svelte';
  import type {
    PriceSchemaRecord,
    CreatePriceSchemaDTO,
    UpdatePriceSchemaDTO,
    AdjustmentType
  } from '$lib/types/finance/price_schema';
  import type { PricePlanRecord } from '$lib/types/finance/price_plan';
  import { getPricePlans } from '$lib/api/finance/price_plan';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    schema: PriceSchemaRecord | null;
    isSaving: boolean;
    onClose: () => void;
    onSave: (payload: CreatePriceSchemaDTO | UpdatePriceSchemaDTO) => Promise<void>;
  }

  let { isOpen, mode, schema, isSaving, onClose, onSave }: Props = $props();

  let code = $state('');
  let name = $state('');
  let description = $state('');
  let basePricePlanId = $state('');
  let pricePlans = $state<PricePlanRecord[]>([]);

  // Initial Default Rule state (used when creating a new schema)
  let tariffAdjType = $state<AdjustmentType>('NONE');
  let tariffAdjValue = $state<number>(0);
  let patientPercent = $state<number>(0);
  let customerPercent = $state<number>(100);
  let taxPercent = $state<number>(11);

  let validationError = $state<string | null>(null);

  onMount(async () => {
    try {
      const res = await getPricePlans({ limit: 100 });
      pricePlans = res.data || [];
    } catch (err) {
      console.error('Gagal mengambil daftar buku tarif:', err);
    }
  });

  $effect(() => {
    if (isOpen) {
      validationError = null;
      if (mode === 'edit' && schema) {
        code = schema.code;
        name = schema.name;
        description = schema.description || '';
        basePricePlanId = schema.base_price_plan_id || '';
      } else {
        code = '';
        name = '';
        description = '';
        basePricePlanId = '';
        tariffAdjType = 'NONE';
        tariffAdjValue = 0;
        patientPercent = 0;
        customerPercent = 100;
        taxPercent = 11;
      }
    }
  });

  function handlePatientPctChange(val: number) {
    patientPercent = Math.max(0, Math.min(100, val || 0));
    customerPercent = 100 - patientPercent;
  }

  function handleCustomerPctChange(val: number) {
    customerPercent = Math.max(0, Math.min(100, val || 0));
    patientPercent = 100 - customerPercent;
  }

  function applyPreset(patient: number, customer: number) {
    patientPercent = patient;
    customerPercent = customer;
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    validationError = null;

    if (!name.trim()) {
      validationError = 'Nama skema wajib diisi';
      return;
    }

    if (mode === 'create') {
      if (!code.trim()) {
        validationError = 'Kode skema wajib diisi';
        return;
      }

      if (patientPercent + customerPercent !== 100) {
        validationError = 'Total pembagian tanggungan (Pasien + Penjamin) harus tepat 100%';
        return;
      }

      const createPayload: CreatePriceSchemaDTO = {
        code: code.trim().toUpperCase(),
        name: name.trim(),
        description: description.trim() || undefined,
        base_price_plan_id: basePricePlanId || undefined,
        rules: [
          {
            departement_id: null,
            service_unit_id: null,
            item_id: null,
            tariff_adjustment_type: tariffAdjType,
            tariff_adjustment_value: Number(tariffAdjValue) || 0,
            medication_adjustment_type: 'NONE',
            medication_adjustment_value: 0,
            pay_by_patient_percent: patientPercent,
            pay_by_customer_percent: customerPercent,
            tax_percent: taxPercent
          }
        ]
      };

      await onSave(createPayload);
    } else {
      const updatePayload: UpdatePriceSchemaDTO = {
        name: name.trim(),
        description: description.trim() || undefined,
        base_price_plan_id: basePricePlanId || undefined
      };
      await onSave(updatePayload);
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs animate-in fade-in duration-150">
    <div
      class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-xl max-w-lg w-full overflow-hidden flex flex-col max-h-[90vh]"
    >
      <!-- Modal Header -->
      <div class="px-6 py-4 bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-bold text-on-surface">
            {mode === 'create' ? 'Tambah Skema Tarif PKS' : 'Ubah Header Skema'}
          </h3>
          <p class="text-xs text-on-surface-variant mt-0.5">
            {mode === 'create' ? 'Tentukan kode, nama, dan aturan dasar fallback skema baru' : `Mengubah informasi skema ${code}`}
          </p>
        </div>
        <button
          type="button"
          onclick={onClose}
          class="p-1.5 rounded-full text-on-surface-variant hover:text-on-surface hover:bg-surface-container transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body -->
      <form onsubmit={handleSubmit} class="p-6 overflow-y-auto flex flex-col gap-4 text-xs">
        {#if validationError}
          <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800/60 text-rose-700 dark:text-rose-300 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{validationError}</span>
          </div>
        {/if}

        <!-- Kode Skema (Hanya saat create) -->
        {#if mode === 'create'}
          <div class="flex flex-col gap-1.5">
            <label for="schema_code" class="font-semibold text-on-surface flex items-center justify-between">
              <span>Kode Skema</span>
              <span class="text-[11px] font-normal text-on-surface-variant">Unik (Uppercase)</span>
            </label>
            <input
              id="schema_code"
              type="text"
              bind:value={code}
              placeholder="Contoh: BPJS-2026, ASURANSI-PRU"
              class="px-3.5 py-2.5 rounded-xl border border-outline-variant bg-surface text-on-surface font-mono uppercase focus:outline-hidden focus:border-primary focus:ring-1 focus:ring-primary"
              required
            />
          </div>
        {/if}

        <!-- Nama Skema -->
        <div class="flex flex-col gap-1.5">
          <label for="schema_name" class="font-semibold text-on-surface">Nama Skema Tarif</label>
          <input
            id="schema_name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Skema Penjaminan BPJS Kesehatan Kelas 1"
            class="px-3.5 py-2.5 rounded-xl border border-outline-variant bg-surface text-on-surface focus:outline-hidden focus:border-primary focus:ring-1 focus:ring-primary"
            required
          />
        </div>

        <!-- Deskripsi -->
        <div class="flex flex-col gap-1.5">
          <label for="schema_desc" class="font-semibold text-on-surface">Deskripsi (Opsional)</label>
          <textarea
            id="schema_desc"
            bind:value={description}
            rows="2"
            placeholder="Keterangan peruntukan skema, nomor PKS, atau catatan keuangan..."
            class="px-3.5 py-2.5 rounded-xl border border-outline-variant bg-surface text-on-surface focus:outline-hidden focus:border-primary focus:ring-1 focus:ring-primary resize-none"
          ></textarea>
        </div>

        <!-- Buku Tarif Acuan Dasar (Price Plan) -->
        <div class="flex flex-col gap-1.5">
          <label for="base_price_plan" class="font-semibold text-on-surface flex items-center justify-between">
            <span class="flex items-center gap-1.5">
              <BookOpen class="w-3.5 h-3.5 text-primary" />
              <span>Buku Tarif Acuan Dasar</span>
            </span>
            <span class="text-[11px] font-normal text-on-surface-variant">Opsional</span>
          </label>
          <select
            id="base_price_plan"
            bind:value={basePricePlanId}
            class="px-3.5 py-2.5 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs focus:outline-hidden focus:border-primary focus:ring-1 focus:ring-primary"
          >
            <option value="">-- Ikuti Buku Tarif Standar RS Terkini (Otomatis) --</option>
            {#each pricePlans as plan (plan.id)}
              <option value={plan.id}>
                [{plan.code}] {plan.name} {plan.is_default ? '★ (Default Standar RS)' : ''}
              </option>
            {/each}
          </select>
          <p class="text-[11px] text-on-surface-variant leading-relaxed">
            Pilih buku tarif tertentu jika penjamin ini mengacu pada versi tarif khusus/lama (misal: Buku Tarif 2024 atau BPJS TK). Kosongkan jika ingin otomatis mengikuti versi standar RS terbaru.
          </p>
        </div>

        <!-- Konfigurasi Aturan Default Awal (Hanya saat Create) -->
        {#if mode === 'create'}
          <div class="mt-2 p-4 rounded-2xl bg-surface-container-high/60 border border-outline-variant/60 flex flex-col gap-3.5">
            <div class="flex items-center gap-1.5 text-primary font-bold">
              <Sparkles class="w-4 h-4" />
              <span>Aturan Dasar Skema (SCHEMA_DEFAULT Fallback)</span>
            </div>
            <p class="text-[11px] text-on-surface-variant leading-relaxed">
              Aturan ini menjadi dasar perhitungan untuk seluruh tindakan sebelum adanya aturan spesifik pada tingkat departemen, unit, atau item.
            </p>

            <!-- Tipe Penyesuaian Tarif -->
            <div class="grid grid-cols-2 gap-2">
              <div class="flex flex-col gap-1">
                <label for="adj_type" class="font-semibold text-on-surface">Penyesuaian Tarif</label>
                <select
                  id="adj_type"
                  bind:value={tariffAdjType}
                  class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
                >
                  <option value="NONE">Harga Dasar (0%)</option>
                  <option value="PERCENT">Persentase (%)</option>
                  <option value="FIXED_AMOUNT">Nominal (+/- Rp)</option>
                </select>
              </div>

              {#if tariffAdjType !== 'NONE'}
                <div class="flex flex-col gap-1">
                  <label for="adj_val" class="font-semibold text-on-surface">
                    Nilai ({tariffAdjType === 'PERCENT' ? '%' : 'Rp'})
                  </label>
                  <input
                    id="adj_val"
                    type="number"
                    step={tariffAdjType === 'PERCENT' ? '0.01' : '1'}
                    bind:value={tariffAdjValue}
                    placeholder={tariffAdjType === 'PERCENT' ? 'Misal: -10 (Diskon)' : 'Misal: 50000'}
                    class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
                  />
                </div>
              {/if}
            </div>

            <!-- Cost Sharing -->
            <div class="flex flex-col gap-1.5">
              <div class="flex items-center justify-between">
                <span class="font-semibold text-on-surface">Porsi Tanggungan Tagihan (Cost Sharing)</span>
                <div class="flex items-center gap-1">
                  <button
                    type="button"
                    onclick={() => applyPreset(0, 100)}
                    class="px-1.5 py-0.5 rounded text-[10px] font-medium bg-surface hover:bg-surface-container border border-outline-variant/60 cursor-pointer"
                  >
                    100% Penjamin
                  </button>
                  <button
                    type="button"
                    onclick={() => applyPreset(100, 0)}
                    class="px-1.5 py-0.5 rounded text-[10px] font-medium bg-surface hover:bg-surface-container border border-outline-variant/60 cursor-pointer"
                  >
                    100% Pasien
                  </button>
                  <button
                    type="button"
                    onclick={() => applyPreset(20, 80)}
                    class="px-1.5 py-0.5 rounded text-[10px] font-medium bg-surface hover:bg-surface-container border border-outline-variant/60 cursor-pointer"
                  >
                    Co-Pay 20:80
                  </button>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2">
                <div class="flex flex-col gap-1">
                  <label for="pat_pct" class="text-[11px] text-on-surface-variant">Pasien (%)</label>
                  <input
                    id="pat_pct"
                    type="number"
                    min="0"
                    max="100"
                    value={patientPercent}
                    oninput={(e) => handlePatientPctChange(Number(e.currentTarget.value))}
                    class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs font-semibold"
                  />
                </div>

                <div class="flex flex-col gap-1">
                  <label for="cust_pct" class="text-[11px] text-on-surface-variant">Penjamin / Customer (%)</label>
                  <input
                    id="cust_pct"
                    type="number"
                    min="0"
                    max="100"
                    value={customerPercent}
                    oninput={(e) => handleCustomerPctChange(Number(e.currentTarget.value))}
                    class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs font-semibold"
                  />
                </div>
              </div>
            </div>

            <!-- Pajak PPN -->
            <div class="flex items-center justify-between pt-1 border-t border-outline-variant/40">
              <span class="text-on-surface font-medium">Tarif Pajak (PPN)</span>
              <div class="flex items-center gap-2">
                <input
                  id="tax_pct"
                  type="number"
                  min="0"
                  max="100"
                  bind:value={taxPercent}
                  class="w-16 px-2.5 py-1 text-right rounded-lg border border-outline-variant bg-surface text-xs font-semibold"
                />
                <span class="font-bold">%</span>
              </div>
            </div>
          </div>
        {/if}

        <!-- Actions -->
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-outline-variant/60 mt-2">
          <button
            type="button"
            onclick={onClose}
            class="px-4 py-2 text-xs font-semibold rounded-xl text-on-surface-variant hover:bg-surface-container transition-colors cursor-pointer"
            disabled={isSaving}
          >
            Batal
          </button>
          <button
            type="submit"
            class="px-4 py-2 text-xs font-semibold rounded-xl bg-primary text-on-primary hover:bg-primary/90 transition-colors flex items-center gap-1.5 cursor-pointer shadow-xs"
            disabled={isSaving}
          >
            {#if isSaving}
              <span class="loading loading-spinner loading-xs"></span>
              <span>Menyimpan...</span>
            {:else}
              <Check class="w-3.5 h-3.5" />
              <span>{mode === 'create' ? 'Buat Skema Baru' : 'Simpan Perubahan'}</span>
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
