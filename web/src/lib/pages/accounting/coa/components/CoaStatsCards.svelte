<script lang="ts">
  import { Layers, Landmark } from '@lucide/svelte';
  import type { AccountType } from '$lib/types/accounting/account';

  interface Props {
    stats: {
      total: number;
      headers: number;
      postable: number;
      treasury: number;
      byType: Record<AccountType, number>;
    };
  }

  let { stats }: Props = $props();
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
  <!-- Total Akun -->
  <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col justify-between">
    <div class="flex items-center justify-between text-on-surface-variant">
      <span class="text-xs font-medium">Total Akun</span>
      <Layers class="w-4 h-4 text-primary" />
    </div>
    <div class="mt-2 flex items-baseline gap-2">
      <span class="text-2xl font-bold text-on-surface">{stats.total}</span>
      <span class="text-[11px] text-on-surface-variant">({stats.headers} Induk, {stats.postable} Pos)</span>
    </div>
  </div>

  <!-- 1. Aset (Aktiva) -->
  <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col justify-between">
    <div class="flex items-center justify-between text-on-surface-variant">
      <span class="text-xs font-medium">1. Aset (Aktiva)</span>
      <span class="w-2.5 h-2.5 rounded-full bg-blue-500"></span>
    </div>
    <div class="mt-2 flex items-baseline gap-2">
      <span class="text-2xl font-bold text-blue-700 dark:text-blue-300">{stats.byType.ASSET}</span>
      <span class="text-[11px] text-blue-600 dark:text-blue-400 font-mono">1xxx (Debit)</span>
    </div>
  </div>

  <!-- 2 & 3. Kewajiban & Modal -->
  <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col justify-between">
    <div class="flex items-center justify-between text-on-surface-variant">
      <span class="text-xs font-medium">2 & 3. Kewajiban & Modal</span>
      <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
    </div>
    <div class="mt-2 flex items-baseline gap-2">
      <span class="text-2xl font-bold text-amber-700 dark:text-amber-300">{stats.byType.LIABILITY + stats.byType.EQUITY}</span>
      <span class="text-[11px] text-amber-600 dark:text-amber-400 font-mono">2xxx / 3xxx (Kredit)</span>
    </div>
  </div>

  <!-- 4 & 5. Pendapatan & Beban -->
  <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col justify-between">
    <div class="flex items-center justify-between text-on-surface-variant">
      <span class="text-xs font-medium">4 & 5. Pendapatan & Beban</span>
      <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
    </div>
    <div class="mt-2 flex items-baseline gap-2">
      <span class="text-2xl font-bold text-emerald-700 dark:text-emerald-300">{stats.byType.REVENUE + stats.byType.EXPENSE}</span>
      <span class="text-[11px] text-emerald-600 dark:text-emerald-400 font-mono">4xxx / 5xxx</span>
    </div>
  </div>

  <!-- Kas & Bank (Treasury) -->
  <div class="p-3.5 rounded-2xl bg-white dark:bg-surface-container border border-outline-variant/60 shadow-2xs flex flex-col justify-between col-span-2 sm:col-span-1">
    <div class="flex items-center justify-between text-on-surface-variant">
      <span class="text-xs font-medium">Kas & Bank (Treasury)</span>
      <Landmark class="w-4 h-4 text-teal-600 dark:text-teal-400" />
    </div>
    <div class="mt-2 flex items-baseline gap-2">
      <span class="text-2xl font-bold text-teal-700 dark:text-teal-300">{stats.treasury}</span>
      <span class="text-[11px] text-teal-600 dark:text-teal-400">Rekening RS</span>
    </div>
  </div>
</div>
