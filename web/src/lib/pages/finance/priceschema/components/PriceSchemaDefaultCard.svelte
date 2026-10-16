<script lang="ts">
  import {
    ShieldCheck,
    Edit2,
    Sparkles,
    Receipt,
    Percent,
    Landmark,
  } from "@lucide/svelte";
  import type { PriceSchemaRuleRecord } from "$lib/types/finance/price_schema";
  import { formatAdjustmentText } from "../treeUtils";

  interface Props {
    defaultRule: PriceSchemaRuleRecord | null;
    onEditDefault: () => void;
  }

  let { defaultRule, onEditDefault }: Props = $props();

  const adjType = $derived(defaultRule?.tariff_adjustment_type || "NONE");
  const adjVal = $derived(
    defaultRule?.tariff_adjustment_value
      ? Number(defaultRule.tariff_adjustment_value)
      : 0,
  );
  const patientPct = $derived(
    defaultRule?.pay_by_patient_percent
      ? Number(defaultRule.pay_by_patient_percent)
      : 0,
  );
  const customerPct = $derived(
    defaultRule?.pay_by_customer_percent
      ? Number(defaultRule.pay_by_customer_percent)
      : 100,
  );
  const taxPct = $derived(
    defaultRule?.tax_percent ? Number(defaultRule.tax_percent) : 11,
  );
</script>

<div
  class="p-5 rounded-3xl bg-linear-to-r from-blue-50/80 via-indigo-50/40 to-surface border border-blue-200/70 dark:from-blue-950/20 dark:via-indigo-950/10 dark:to-surface-container dark:border-blue-800/40 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-4"
>
  <div class="flex items-start gap-3.5">
    <div
      class="w-10 h-10 rounded-2xl bg-primary text-on-primary flex items-center justify-center shrink-0 shadow-xs"
    >
      <ShieldCheck class="w-5 h-5" />
    </div>

    <div class="flex flex-col gap-1 min-w-0">
      <div class="flex items-center gap-2">
        <h4 class="text-sm font-bold text-on-surface">Aturan Dasar Skema</h4>
        <span
          class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-bold bg-blue-100 text-blue-800 dark:bg-blue-900/60 dark:text-blue-200"
        >
          Prioritas Terakhir
        </span>
      </div>
      <p class="text-xs text-on-surface-variant max-w-2xl leading-relaxed">
        Jangkar dasar penagihan layanan jika suatu departemen, unit layanan,
        atau item tindakan tidak memiliki aturan override tersendiri.
      </p>

      <!-- Pill Metrics -->
      <div class="flex flex-wrap items-center gap-2 mt-2">
        <!-- Penyesuaian Tarif -->
        <div
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-white dark:bg-surface-container border border-outline-variant/60 text-xs font-semibold text-on-surface"
        >
          <Percent class="w-3.5 h-3.5 text-primary" />
          <span>Tarif: {formatAdjustmentText(adjType, adjVal)}</span>
        </div>

        <!-- Cost Sharing -->
        <div
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-white dark:bg-surface-container border border-outline-variant/60 text-xs font-semibold text-on-surface"
        >
          <Landmark class="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
          <span>Tanggungan: Pasien {patientPct}% • Penjamin {customerPct}%</span
          >
        </div>

        <!-- Pajak -->
        <div
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl bg-white dark:bg-surface-container border border-outline-variant/60 text-xs font-semibold text-on-surface"
        >
          <Receipt class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" />
          <span>PPN: {taxPct}%</span>
        </div>

        <!-- Farmasi (Fase 2) -->
        <div
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-xl bg-surface-container text-[11px] text-on-surface-variant"
        >
          <span>Farmasi: 0%</span>
          <span
            class="text-[9px] font-mono px-1 rounded bg-surface-container-high"
            >Fase 2</span
          >
        </div>
      </div>
    </div>
  </div>

  <div class="shrink-0 flex items-center">
    <button
      type="button"
      onclick={onEditDefault}
      class="px-4 py-2 text-xs font-semibold rounded-2xl bg-white dark:bg-surface-container hover:bg-blue-50 dark:hover:bg-surface-container-high border border-blue-200 dark:border-blue-800 text-primary flex items-center gap-2 cursor-pointer shadow-2xs transition-colors"
    >
      <Edit2 class="w-3.5 h-3.5" />
      <span>Ubah Aturan Dasar</span>
    </button>
  </div>
</div>
