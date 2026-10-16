<script lang="ts">
  import {
    SlidersHorizontal,
    Edit2,
    Power,
    CheckCircle2,
    XCircle,
    FileSpreadsheet,
    Layers,
    Receipt
  } from '@lucide/svelte';
  import type { PriceSchemaRecord } from '$lib/types/finance/price_schema';

  interface Props {
    schemas: PriceSchemaRecord[];
    isLoading: boolean;
    onManage: (schema: PriceSchemaRecord) => void;
    onEdit: (schema: PriceSchemaRecord) => void;
    onToggleStatus: (schema: PriceSchemaRecord) => void;
  }

  let { schemas, isLoading, onManage, onEdit, onToggleStatus }: Props = $props();

  function getDefaultRuleSummary(schema: PriceSchemaRecord): string {
    const defRule = schema.rules?.find(r => r.level === 'SCHEMA_DEFAULT');
    if (!defRule) return 'Default RS (0%)';
    const type = defRule.tariff_adjustment_type;
    const val = Number(defRule.tariff_adjustment_value) || 0;
    let adj = 'Tarif 0%';
    if (type === 'PERCENT') adj = val > 0 ? `+${val}%` : `${val}%`;
    else if (type === 'FIXED_AMOUNT') adj = val > 0 ? `+Rp ${val}` : `-Rp ${Math.abs(val)}`;

    return `${adj} • Pasien ${defRule.pay_by_patient_percent}% / Penjamin ${defRule.pay_by_customer_percent}%`;
  }
</script>

<div class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-2xs overflow-hidden">
  <div class="overflow-x-auto">
    <table class="w-full text-left text-xs">
      <thead>
        <tr class="bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 text-on-surface-variant font-semibold">
          <th class="px-5 py-3.5 w-36">Kode Skema</th>
          <th class="px-5 py-3.5">Nama & Deskripsi Skema</th>
          <th class="px-5 py-3.5">Aturan Default (Fallback)</th>
          <th class="px-5 py-3.5 text-center">Aturan Khusus</th>
          <th class="px-5 py-3.5 text-center">Status</th>
          <th class="px-5 py-3.5 text-right">Aksi</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-outline-variant/40">
        {#if isLoading}
          <tr>
            <td colspan="6" class="px-5 py-12 text-center text-on-surface-variant">
              <div class="inline-flex items-center gap-2">
                <span class="loading loading-spinner loading-sm text-primary"></span>
                <span>Memuat daftar skema tarif...</span>
              </div>
            </td>
          </tr>
        {:else if schemas.length === 0}
          <tr>
            <td colspan="6" class="px-5 py-14 text-center text-on-surface-variant">
              <div class="flex flex-col items-center justify-center gap-2 max-w-sm mx-auto">
                <div class="w-12 h-12 rounded-2xl bg-surface-container-high flex items-center justify-center text-on-surface-variant mb-1">
                  <Receipt class="w-6 h-6 opacity-60" />
                </div>
                <p class="font-semibold text-sm text-on-surface">Belum ada skema tarif</p>
                <p class="text-xs text-on-surface-variant">
                  Skema tarif digunakan untuk mengatur penyesuaian harga (diskon/markup) per rekanan PKS.
                </p>
              </div>
            </td>
          </tr>
        {:else}
          {#each schemas as schema (schema.id)}
            <tr class="hover:bg-[#f8fafd]/80 dark:hover:bg-surface-container-high/60 transition-colors">
              <!-- Kode Skema -->
              <td class="px-5 py-4 font-mono font-bold text-primary">
                {schema.code}
              </td>

              <!-- Nama & Deskripsi -->
              <td class="px-5 py-4">
                <div class="flex items-center gap-2 flex-wrap mb-0.5">
                  <span class="font-bold text-on-surface text-sm">{schema.name}</span>
                  {#if schema.base_price_plan_code}
                    <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-mono font-bold bg-purple-50 text-purple-700 dark:bg-purple-900/50 dark:text-purple-300 border border-purple-200 dark:border-purple-800" title="Buku Tarif Acuan Dasar">
                      {schema.base_price_plan_code}
                    </span>
                  {/if}
                </div>
                {#if schema.description}
                  <div class="text-[11px] text-on-surface-variant line-clamp-1">{schema.description}</div>
                {:else}
                  <div class="text-[11px] text-outline italic">Tidak ada deskripsi</div>
                {/if}
              </td>

              <!-- Aturan Default -->
              <td class="px-5 py-4">
                <span class="inline-flex items-center px-2.5 py-1 rounded-xl text-[11px] font-medium bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800/50">
                  {getDefaultRuleSummary(schema)}
                </span>
              </td>

              <!-- Total Aturan Khusus -->
              <td class="px-5 py-4 text-center">
                <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl text-[11px] font-semibold bg-surface-container-high text-on-surface">
                  <Layers class="w-3.5 h-3.5 text-primary" />
                  <span>{schema.rules_count ?? (schema.rules?.length || 0)} aturan</span>
                </span>
              </td>

              <!-- Status -->
              <td class="px-5 py-4 text-center">
                {#if schema.is_active}
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800/50">
                    <CheckCircle2 class="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
                    <span>Aktif</span>
                  </span>
                {:else}
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700">
                    <XCircle class="w-3 h-3 text-zinc-500" />
                    <span>Nonaktif</span>
                  </span>
                {/if}
              </td>

              <!-- Aksi -->
              <td class="px-5 py-4 text-right">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    type="button"
                    onclick={() => onManage(schema)}
                    class="px-3 py-1.5 text-xs font-semibold rounded-xl bg-primary text-on-primary hover:bg-primary/90 transition-colors flex items-center gap-1.5 cursor-pointer shadow-2xs"
                    title="Kelola aturan hierarkis skema ini"
                  >
                    <SlidersHorizontal class="w-3.5 h-3.5" />
                    <span>Kelola Aturan</span>
                  </button>

                  <button
                    type="button"
                    onclick={() => onEdit(schema)}
                    class="p-1.5 rounded-xl text-on-surface-variant hover:text-on-surface hover:bg-surface-container-high transition-colors cursor-pointer"
                    title="Ubah info header skema"
                  >
                    <Edit2 class="w-3.5 h-3.5" />
                  </button>

                  <button
                    type="button"
                    onclick={() => onToggleStatus(schema)}
                    class="p-1.5 rounded-xl {schema.is_active ? 'text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40' : 'text-emerald-600 hover:bg-emerald-50 dark:hover:bg-emerald-950/40'} transition-colors cursor-pointer"
                    title={schema.is_active ? 'Nonaktifkan skema' : 'Aktifkan skema'}
                  >
                    <Power class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
