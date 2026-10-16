<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Calculator,
    CheckCircle2,
    AlertCircle,
    Layers,
    Sparkles,
    ShieldCheck,
    Receipt,
    Percent,
    ArrowRight,
    BookOpen
  } from '@lucide/svelte';
  import type {
    SimulateTariffDTO,
    SimulateTariffResult
  } from '$lib/types/finance/price_schema';
  import type { PricePlanRecord } from '$lib/types/finance/price_plan';
  import type { DepartementRecord } from '$lib/types/master/departement';
  import type { ServiceUnitRecord } from '$lib/types/master/service_unit';
  import type { TariffClassRecord } from '$lib/types/finance/tariff';
  import type { ItemSummaryRecord } from '$lib/types/master/item';
  import { simulatePriceSchemaTariff } from '$lib/api/finance/price_schema';
  import { getPricePlans } from '$lib/api/finance/price_plan';

  interface Props {
    priceSchemaId: string;
    basePricePlanId?: string | null;
    departments: DepartementRecord[];
    serviceUnits: ServiceUnitRecord[];
    tariffClasses: TariffClassRecord[];
    items: ItemSummaryRecord[];
  }

  let {
    priceSchemaId,
    basePricePlanId,
    departments,
    serviceUnits,
    tariffClasses,
    items
  }: Props = $props();

  // Simulation Form State
  let selectedPricePlanId = $state<string>('');
  let pricePlans = $state<PricePlanRecord[]>([]);
  let selectedDeptId = $state<string>('');
  let selectedUnitId = $state<string>('');
  let selectedClassId = $state<string>('');
  let selectedItemId = $state<string>('');
  let isCito = $state<boolean>(false);
  let serviceDate = $state<string>(new Date().toISOString().split('T')[0]);

  onMount(async () => {
    try {
      const res = await getPricePlans({ limit: 100 });
      pricePlans = res.data || [];
    } catch (err) {
      console.error('Gagal mengambil daftar buku tarif:', err);
    }
  });

  $effect(() => {
    selectedPricePlanId = basePricePlanId || '';
  });

  let isCalculating = $state(false);
  let calculationError = $state<string | null>(null);
  let result = $state<SimulateTariffResult | null>(null);

  // Filter service units when department changes
  const filteredUnits = $derived.by(() => {
    if (!selectedDeptId) return serviceUnits;
    return serviceUnits.filter(u => u.departement_id === selectedDeptId);
  });

  // Set default tariff class if available
  $effect(() => {
    if (tariffClasses.length > 0 && !selectedClassId) {
      selectedClassId = tariffClasses[0].id;
    }
  });

  async function handleSimulate() {
    calculationError = null;
    result = null;

    if (!selectedClassId) {
      calculationError = 'Silakan pilih Kelas Tarif';
      return;
    }

    if (!selectedItemId) {
      calculationError = 'Silakan pilih Item Tindakan Medis';
      return;
    }

    isCalculating = true;
    try {
      const payload: SimulateTariffDTO = {
        price_schema_id: priceSchemaId || undefined,
        price_plan_id: selectedPricePlanId || undefined,
        tariff_class_id: selectedClassId,
        item_id: selectedItemId,
        departement_id: selectedDeptId || undefined,
        service_unit_id: selectedUnitId || undefined,
        service_date: serviceDate || undefined,
        is_cito: isCito
      };

      const res = await simulatePriceSchemaTariff(payload);
      result = res;
    } catch (err: any) {
      calculationError = err.message || 'Gagal menghitung simulasi tarif. Pastikan item memiliki tarif di Price Plan.';
    } finally {
      isCalculating = false;
    }
  }

  function formatRp(val: number | string | undefined): string {
    const num = typeof val === 'string' ? parseFloat(val) || 0 : val || 0;
    return 'Rp ' + Math.round(num).toLocaleString('id-ID');
  }
</script>

<div class="p-6 rounded-3xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col gap-5">
  <div class="flex items-center gap-2.5 pb-3 border-b border-outline-variant/40">
    <div class="w-8 h-8 rounded-xl bg-purple-100 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300 flex items-center justify-center">
      <Calculator class="w-4 h-4" />
    </div>
    <div>
      <h3 class="text-sm font-bold text-on-surface">Simulator & Uji Lookup Penyesuaian Tarif</h3>
      <p class="text-xs text-on-surface-variant">
        Uji coba perhitungan tarif secara instan untuk membuktikan aturan penyesuaian PKS bekerja akurat sebelum digunakan pada pelayanan pasien.
      </p>
    </div>
  </div>

  {#if calculationError}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800/60 text-rose-700 dark:text-rose-300 flex items-center gap-2 text-xs">
      <AlertCircle class="w-4 h-4 shrink-0" />
      <span>{calculationError}</span>
    </div>
  {/if}

  <!-- Form Grid -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3 text-xs">
    <!-- Buku Tarif Acuan -->
    <div class="flex flex-col gap-1">
      <label for="sim_plan" class="font-semibold text-on-surface flex items-center justify-between">
        <span class="flex items-center gap-1">
          <BookOpen class="w-3.5 h-3.5 text-primary" />
          <span>Buku Tarif Acuan</span>
        </span>
        {#if basePricePlanId && !selectedPricePlanId}
          <span class="text-[10px] text-purple-600 font-bold">(Terkunci)</span>
        {/if}
      </label>
      <select
        id="sim_plan"
        bind:value={selectedPricePlanId}
        class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs font-medium"
      >
        <option value="">
          {basePricePlanId ? '-- Gunakan Buku Tarif Skema --' : '-- Ikuti Standar RS Otomatis --'}
        </option>
        {#each pricePlans as plan (plan.id)}
          <option value={plan.id}>
            [{plan.code}] {plan.name} {plan.is_default ? '★' : ''}
          </option>
        {/each}
      </select>
    </div>

    <!-- Departemen -->
    <div class="flex flex-col gap-1">
      <label for="sim_dept" class="font-semibold text-on-surface">Departemen</label>
      <select
        id="sim_dept"
        bind:value={selectedDeptId}
        onchange={() => (selectedUnitId = '')}
        class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
      >
        <option value="">Semua Departemen (Opsional)</option>
        {#each departments as d (d.id)}
          <option value={d.id}>{d.name}</option>
        {/each}
      </select>
    </div>

    <!-- Unit Layanan -->
    <div class="flex flex-col gap-1">
      <label for="sim_unit" class="font-semibold text-on-surface">Unit Layanan</label>
      <select
        id="sim_unit"
        bind:value={selectedUnitId}
        class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
      >
        <option value="">Semua Unit (Opsional)</option>
        {#each filteredUnits as u (u.id)}
          <option value={u.id}>{u.name}</option>
        {/each}
      </select>
    </div>

    <!-- Kelas Tarif -->
    <div class="flex flex-col gap-1">
      <label for="sim_class" class="font-semibold text-on-surface">Kelas Rawat / Tarif *</label>
      <select
        id="sim_class"
        bind:value={selectedClassId}
        class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs font-semibold"
        required
      >
        <option value="" disabled>Pilih Kelas Tarif</option>
        {#each tariffClasses as tc (tc.id)}
          <option value={tc.id}>{tc.name} ({tc.code})</option>
        {/each}
      </select>
    </div>

    <!-- Item Tindakan -->
    <div class="flex flex-col gap-1">
      <label for="sim_item" class="font-semibold text-on-surface">Item Tindakan Medis *</label>
      <select
        id="sim_item"
        bind:value={selectedItemId}
        class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs font-semibold"
        required
      >
        <option value="" disabled>Pilih Tindakan Medis</option>
        {#each items as it (it.id)}
          <option value={it.id}>{it.name} ({it.code})</option>
        {/each}
      </select>
    </div>
  </div>

  <div class="flex items-center justify-between pt-2">
    <div class="flex items-center gap-4 text-xs">
      <label class="flex items-center gap-2 cursor-pointer font-medium text-on-surface">
        <input
          type="checkbox"
          bind:checked={isCito}
          class="rounded text-primary focus:ring-primary w-4 h-4 cursor-pointer"
        />
        <span>Tindakan Darurat (Cito)</span>
      </label>

      <div class="flex items-center gap-1.5 text-on-surface-variant">
        <span>Tanggal:</span>
        <input
          type="date"
          bind:value={serviceDate}
          class="px-2 py-1 rounded-lg border border-outline-variant bg-surface text-xs font-mono"
        />
      </div>
    </div>

    <button
      type="button"
      onclick={handleSimulate}
      class="px-5 py-2 text-xs font-bold rounded-xl bg-purple-600 hover:bg-purple-700 text-white flex items-center gap-1.5 cursor-pointer shadow-xs transition-colors"
      disabled={isCalculating}
    >
      {#if isCalculating}
        <span class="loading loading-spinner loading-xs"></span>
        <span>Menghitung...</span>
      {:else}
        <Calculator class="w-3.5 h-3.5" />
        <span>Hitung Simulasi Tarif</span>
      {/if}
    </button>
  </div>

  <!-- Simulation Results Breakdown -->
  {#if result}
    <div class="mt-2 p-5 rounded-3xl bg-linear-to-br from-[#f8fafd] to-surface border border-outline-variant/60 shadow-xs flex flex-col gap-4 text-xs animate-in fade-in duration-200">
      <!-- Top Result Header: Winning Rule Level -->
      <div class="flex flex-wrap items-center justify-between gap-2 pb-3 border-b border-outline-variant/40">
        <div class="flex items-center gap-2">
          <span class="font-bold text-on-surface">Aturan Terpilih (Winning Rule):</span>
          <span class="px-2.5 py-0.5 rounded-full text-xs font-bold bg-purple-100 text-purple-800 dark:bg-purple-900/60 dark:text-purple-200 border border-purple-200 dark:border-purple-800">
            {result.match_level || 'SCHEMA_DEFAULT'}
          </span>
          <span class="text-[11px] text-on-surface-variant">
            (Penyesuaian: {result.tariff_adjustment_type} {result.tariff_adjustment_value})
          </span>
        </div>

        <div class="font-mono text-xs text-outline">
          Buku Tarif Dasar: {result.price_plan_code}
        </div>
      </div>

      <!-- Financial Metrics Grid -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <!-- Tarif Dasar RS -->
        <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 flex flex-col gap-1">
          <span class="text-[11px] text-on-surface-variant">Harga Dasar RS</span>
          <span class="text-sm font-bold text-on-surface font-mono">{formatRp(result.base_price)}</span>
          {#if result.is_cito}
            <span class="text-[10px] text-amber-600 font-semibold">+ Cito: {formatRp(result.cito_surcharge)}</span>
          {/if}
        </div>

        <!-- Penyesuaian Skema -->
        <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 flex flex-col gap-1">
          <span class="text-[11px] text-on-surface-variant">Tarif Setelah Penyesuaian</span>
          <span class="text-sm font-bold text-primary font-mono">{formatRp(result.final_price)}</span>
          <span class="text-[10px] text-on-surface-variant">PPN ({result.tax_percent}%): {formatRp(result.tax_amount)}</span>
        </div>

        <!-- Beban Pasien -->
        <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 flex flex-col gap-1">
          <span class="text-[11px] text-on-surface-variant">Beban Pasien ({result.pay_by_patient_percent}%)</span>
          <span class="text-sm font-bold text-amber-600 dark:text-amber-400 font-mono">{formatRp(result.patient_amount)}</span>
          <span class="text-[10px] text-on-surface-variant">Tagihan Pribadi Pasien</span>
        </div>

        <!-- Beban Penjamin / Customer -->
        <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 flex flex-col gap-1">
          <span class="text-[11px] text-on-surface-variant">Beban Penjamin ({result.pay_by_customer_percent}%)</span>
          <span class="text-sm font-bold text-emerald-600 dark:text-emerald-400 font-mono">{formatRp(result.customer_amount)}</span>
          <span class="text-[10px] text-on-surface-variant">Klaim PKS Customer</span>
        </div>
      </div>

      <!-- Component Distribution Breakdown Table -->
      {#if result.components && result.components.length > 0}
        <div class="mt-2 pt-3 border-t border-outline-variant/40">
          <h4 class="font-bold text-on-surface mb-2">Alokasi Pecahan Komponen Biaya (Proporsional)</h4>
          <div class="rounded-xl border border-outline-variant/60 overflow-hidden">
            <table class="w-full text-left text-xs">
              <thead class="bg-surface-container-high text-on-surface-variant font-semibold">
                <tr>
                  <th class="px-3.5 py-2">Komponen Tarif</th>
                  <th class="px-3.5 py-2">Akun COA</th>
                  <th class="px-3.5 py-2 text-right">Tarif Dasar</th>
                  <th class="px-3.5 py-2 text-right">Hasil Penyesuaian</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-outline-variant/40">
                {#each result.components as comp (comp.component_id)}
                  <tr>
                    <td class="px-3.5 py-2 font-medium text-on-surface">
                      <span>{comp.component_name}</span>
                      {#if comp.component_code || comp.component_type}
                        <span class="text-[10px] text-outline font-normal ml-1">
                          ({comp.component_code || comp.component_type})
                        </span>
                      {/if}
                    </td>
                    <td class="px-3.5 py-2 font-mono text-[11px] text-outline">
                      {comp.coa_code || '-'}
                    </td>
                    <td class="px-3.5 py-2 text-right font-mono text-on-surface-variant">
                      {formatRp(comp.base_amount ?? comp.amount)}
                    </td>
                    <td class="px-3.5 py-2 text-right font-mono font-bold text-primary">
                      {formatRp(comp.final_amount ?? comp.amount)}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
