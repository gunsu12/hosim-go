<script lang="ts">
  import {
    CheckCircle2,
    AlertTriangle,
    Copy,
    Calculator,
    Layers,
    Plus,
    Trash2
  } from '@lucide/svelte';
  import type { PricePlanRecord } from '$lib/types/finance/price_plan';
  import type { TariffComponentRecord } from '$lib/types/finance/tariff';

  export interface ComponentRowState {
    component_id: string;
    base_amount: number;
    cito_amount: number | null;
    coa_code: string | null;
  }

  export interface ClassMatrixRowState {
    tariff_class_id: string;
    tariff_class_code: string;
    tariff_class_name: string;
    existing_item_id: string | null;
    total_base_price: number;
    total_cito_price: number | null;
    is_active: boolean;
    is_enabled: boolean;
    components: ComponentRowState[];
  }

  interface Props {
    row: ClassMatrixRowState;
    rIdx: number;
    selectedPlan: PricePlanRecord | null;
    tariffComponents: TariffComponentRecord[];
    onAutoBalance: (row: ClassMatrixRowState) => void;
    onCopyToAllClasses: (row: ClassMatrixRowState) => void;
    onCalculateCito: (row: ClassMatrixRowState) => void;
    onAddComponent: (row: ClassMatrixRowState) => void;
    onRemoveComponent: (row: ClassMatrixRowState, cIdx: number) => void;
  }

  let {
    row,
    rIdx,
    selectedPlan,
    tariffComponents,
    onAutoBalance,
    onCopyToAllClasses,
    onCalculateCito,
    onAddComponent,
    onRemoveComponent
  }: Props = $props();

  function formatIDR(val: number): string {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(val || 0);
  }

  let sumComponents = $derived(
    row.components.reduce((acc, c) => acc + (Number(c.base_amount) || 0), 0)
  );
  let isBalanced = $derived(
    Math.abs(sumComponents - (Number(row.total_base_price) || 0)) <= 0.01
  );
  let diff = $derived(
    Math.abs(sumComponents - (Number(row.total_base_price) || 0))
  );
</script>

<div class="bg-white dark:bg-[#1e1f20] rounded-3xl border transition-all duration-200 shadow-2xs overflow-hidden {row.is_enabled ? 'border-blue-300 dark:border-blue-800 ring-1 ring-blue-100 dark:ring-blue-900/30' : 'border-[#e1e5ea] dark:border-[#353738] opacity-95'}">
  <!-- Class Header Bar -->
  <div class="p-3.5 bg-[#f8fafd] dark:bg-[#18191a] border-b border-[#e1e5ea] dark:border-[#353738] flex flex-wrap items-center justify-between gap-3">
    <div class="flex items-center gap-3">
      <input
        type="checkbox"
        bind:checked={row.is_enabled}
        class="w-4 h-4 rounded text-blue-600 border-[#c4c7c5] dark:border-[#353738] focus:ring-blue-500 cursor-pointer"
        title="Sertakan kelas ini dalam penetapan tarif"
      />

      <div class="flex items-center gap-2">
        <span class="font-mono font-bold text-xs px-2 py-0.5 rounded bg-blue-100 text-blue-800 border border-blue-200 dark:bg-blue-950/60 dark:text-blue-300 dark:border-blue-800/60">
          {row.tariff_class_code}
        </span>
        <span class="text-xs font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">{row.tariff_class_name}</span>
        {#if row.existing_item_id}
          <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800/60">
            SUDAH TERKONFIGURASI
          </span>
        {:else}
          <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-amber-50 text-amber-800 border border-amber-200 dark:bg-amber-950/60 dark:text-amber-300 dark:border-amber-800/60">
            BELUM DISET
          </span>
        {/if}
      </div>
    </div>

    <!-- Action Buttons on Header -->
    <div class="flex items-center gap-2">
      {#if row.is_enabled}
        {#if isBalanced && row.total_base_price > 0}
          <span class="text-[11px] font-semibold px-2.5 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800/60 flex items-center gap-1">
            <CheckCircle2 class="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
            <span>Komponen Seimbang ({formatIDR(row.total_base_price)})</span>
          </span>
        {:else if row.total_base_price > 0}
          <button
            type="button"
            onclick={() => onAutoBalance(row)}
            class="text-[11px] font-semibold px-2.5 py-1 rounded-full bg-red-50 text-red-700 border border-red-200 dark:bg-red-950/60 dark:text-red-300 dark:border-red-800/60 flex items-center gap-1.5 hover:bg-red-100 dark:hover:bg-red-900/60 cursor-pointer transition-colors"
            title="Klik untuk seimbangkan otomatis"
          >
            <AlertTriangle class="w-3 h-3 text-red-600 dark:text-red-400" />
            <span>Selisih {formatIDR(diff)} • Klik Seimbangkan</span>
          </button>
        {/if}

        <button
          type="button"
          onclick={() => onCopyToAllClasses(row)}
          class="h-7 px-2.5 rounded-lg border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#353738] text-[11px] font-medium flex items-center gap-1 cursor-pointer transition-colors shadow-2xs"
          title="Salin nominal dan struktur komponen kelas ini ke semua kelas lain"
        >
          <Copy class="w-3 h-3 text-blue-600 dark:text-blue-400" />
          <span>Salin ke Kelas Lain</span>
        </button>
      {/if}
    </div>
  </div>

  <!-- Class Body Form -->
  <div class="p-4 flex flex-col gap-4">
    <!-- Row Pricing Inputs -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
      <!-- Total Base Price -->
      <div class="flex flex-col gap-1">
        <label for={`base-price-${rIdx}`} class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">
          Total Tarif Dasar (Reguler) <span class="text-red-500">*</span>
        </label>
        <div class="relative">
          <span class="absolute left-3 top-2.5 text-xs font-semibold text-[#747775] dark:text-[#8e918f]">Rp</span>
          <input
            id={`base-price-${rIdx}`}
            type="number"
            min="0"
            step="500"
            bind:value={row.total_base_price}
            oninput={() => {
              row.is_enabled = true;
            }}
            placeholder="0"
            class="w-full h-9 pl-9 pr-3 rounded-xl border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:ring-1 focus:ring-[#0b57d0]"
          />
        </div>
      </div>

      <!-- Total Cito Price (Override) -->
      <div class="flex flex-col gap-1">
        <div class="flex items-center justify-between">
          <label for={`cito-price-${rIdx}`} class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">
            Tarif CITO (Override Nominal)
          </label>
          <button
            type="button"
            onclick={() => onCalculateCito(row)}
            class="text-[10px] text-blue-600 dark:text-blue-400 hover:underline font-semibold flex items-center gap-0.5 cursor-pointer"
          >
            <Calculator class="w-2.5 h-2.5" />
            <span>Auto +{selectedPlan?.default_cito_percent || 25}%</span>
          </button>
        </div>
        <div class="relative">
          <span class="absolute left-3 top-2.5 text-xs font-semibold text-[#747775] dark:text-[#8e918f]">Rp</span>
          <input
            id={`cito-price-${rIdx}`}
            type="number"
            min="0"
            step="500"
            bind:value={row.total_cito_price}
            placeholder="Kosongkan jika formula"
            class="w-full h-9 pl-9 pr-3 rounded-xl border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:ring-1 focus:ring-[#0b57d0]"
          />
        </div>
      </div>

      <!-- Is Active Checkbox -->
      <div class="flex flex-col justify-end">
        <label class="h-9 px-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#1e1f20] flex items-center gap-2 cursor-pointer hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] transition-colors">
          <input
            type="checkbox"
            bind:checked={row.is_active}
            class="w-3.5 h-3.5 rounded text-blue-600 border-[#c4c7c5] dark:border-[#353738] focus:ring-blue-500"
          />
          <span class="text-xs font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">Aktif Digunakan di Kasir</span>
        </label>
      </div>
    </div>

    <!-- Pecahan Komponen Biaya Matrix Section -->
    <div class="bg-[#f8fafd] dark:bg-[#18191a] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex flex-col gap-2.5">
      <div class="flex items-center justify-between text-xs">
        <div class="flex items-center gap-1.5 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
          <Layers class="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
          <span>Pecahan Komponen Biaya (Jasa Dokter, Sarana RS, dll.)</span>
        </div>
        <div class="flex items-center gap-2">
          <button
            type="button"
            onclick={() => onAddComponent(row)}
            class="px-2 py-1 rounded-lg bg-blue-50 text-blue-700 hover:bg-blue-100 border border-blue-200 dark:bg-blue-950/60 dark:text-blue-300 dark:border-blue-800/60 dark:hover:bg-blue-900/60 text-[11px] font-semibold flex items-center gap-1 cursor-pointer transition-colors"
          >
            <Plus class="w-3 h-3" />
            <span>Tambah Komponen</span>
          </button>
        </div>
      </div>

      <!-- Table Komponen Biaya -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="text-[11px] font-semibold text-[#747775] dark:text-[#8e918f] border-b border-[#e1e5ea] dark:border-[#353738]">
              <th class="py-1.5 px-2">Komponen Biaya</th>
              <th class="py-1.5 px-2 w-36">Nominal Reguler (Rp)</th>
              <th class="py-1.5 px-2 w-36">Nominal CITO (Rp)</th>
              <th class="py-1.5 px-2 w-32">Akun COA</th>
              <th class="py-1.5 px-2 w-10 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
            {#each row.components as comp, cIdx}
              {@const compDef = tariffComponents.find((tc) => tc.id === comp.component_id)}
              <tr>
                <!-- Select Component -->
                <td class="py-1.5 px-2">
                  <select
                    bind:value={comp.component_id}
                    onchange={(e: any) => {
                      const found = tariffComponents.find((tc) => tc.id === e.target.value);
                      if (found) comp.coa_code = found.default_coa_code || null;
                    }}
                    class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa]"
                  >
                    {#each tariffComponents as tc (tc.id)}
                      <option value={tc.id}>
                        [{tc.component_type}] {tc.name} ({tc.code})
                      </option>
                    {/each}
                  </select>
                </td>

                <!-- Base Amount -->
                <td class="py-1.5 px-2">
                  <input
                    type="number"
                    min="0"
                    step="500"
                    bind:value={comp.base_amount}
                    placeholder="0"
                    class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs font-medium text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa]"
                  />
                </td>

                <!-- Cito Amount -->
                <td class="py-1.5 px-2">
                  <input
                    type="number"
                    min="0"
                    step="500"
                    bind:value={comp.cito_amount}
                    placeholder="Opsional"
                    class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa]"
                  />
                </td>

                <!-- COA Code Override -->
                <td class="py-1.5 px-2">
                  <input
                    type="text"
                    bind:value={comp.coa_code}
                    placeholder={compDef?.default_coa_code || 'Akun COA'}
                    class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs font-mono text-[#444746] dark:text-[#c4c7c5] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa]"
                  />
                </td>

                <!-- Delete Button -->
                <td class="py-1.5 px-2 text-center">
                  <button
                    type="button"
                    onclick={() => onRemoveComponent(row, cIdx)}
                    class="p-1 rounded text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/40 cursor-pointer transition-colors"
                    title="Hapus baris komponen ini"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
