<script lang="ts">
  import { Lock, Unlock, CheckCircle2 } from '@lucide/svelte';
  import type { PeriodStatus } from '$lib/types/accounting/fiscalPeriod';

  interface Props {
    status: PeriodStatus;
    showIcon?: boolean;
    class?: string;
  }

  let { status, showIcon = true, class: className = '' }: Props = $props();
</script>

{#if status === 'OPEN'}
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800/60 {className}"
  >
    {#if showIcon}
      <Unlock class="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
    {/if}
    <span>Terbuka (OPEN)</span>
  </span>
{:else if status === 'LOCKED'}
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-300 border border-amber-200 dark:border-amber-800/60 {className}"
  >
    {#if showIcon}
      <Lock class="w-3 h-3 text-amber-600 dark:text-amber-400" />
    {/if}
    <span>Terkunci Sementara (LOCKED)</span>
  </span>
{:else if status === 'CLOSED'}
  <span
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 border border-slate-200 dark:border-slate-700 {className}"
  >
    {#if showIcon}
      <CheckCircle2 class="w-3 h-3 text-slate-600 dark:text-slate-400" />
    {/if}
    <span>Tutup Buku (CLOSED)</span>
  </span>
{/if}
