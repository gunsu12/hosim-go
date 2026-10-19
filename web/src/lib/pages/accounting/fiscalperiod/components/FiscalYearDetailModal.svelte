<script lang="ts">
  import { Calendar, CheckCircle2, Lock, Unlock, Clock, FileText } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import FiscalPeriodStatusBadge from './FiscalPeriodStatusBadge.svelte';
  import type { FiscalYearRecord } from '$lib/types/accounting/fiscalPeriod';

  interface Props {
    isOpen: boolean;
    fiscalYear: FiscalYearRecord | null;
    onClose: () => void;
    onNavigateToPeriod?: (periodKey: string) => void;
  }

  let { isOpen, fiscalYear, onClose, onNavigateToPeriod }: Props = $props();

  const closedCount = $derived(
    fiscalYear?.periods?.filter((p) => p.status === 'CLOSED').length || 0
  );
  const lockedCount = $derived(
    fiscalYear?.periods?.filter((p) => p.status === 'LOCKED').length || 0
  );
  const openCount = $derived(
    fiscalYear?.periods?.filter((p) => p.status === 'OPEN').length || 0
  );
</script>

<M3Modal
  {isOpen}
  title={fiscalYear ? `Detail Tahun Fiskal ${fiscalYear.year}` : 'Detail Tahun Fiskal'}
  subtitle="Daftar dan status 12 periode bulanan akuntansi"
  icon={Calendar}
  size="3xl"
  {onClose}
>
  {#if fiscalYear}
    <div class="space-y-4">
      <!-- Status Card -->
      <div class="p-4 rounded-2xl bg-[#f0f4f9] dark:bg-[#28292a] flex flex-wrap items-center justify-between gap-4">
        <div>
          <div class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">STATUS TAHUN BUKU</div>
          <div class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center gap-2 mt-0.5">
            {#if fiscalYear.is_closed}
              <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-slate-200 text-slate-800 dark:bg-slate-700 dark:text-slate-200">
                TUTUP BUKU PERMANEN (CLOSED)
              </span>
            {:else}
              <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300">
                AKTIF BERJALAN (OPEN)
              </span>
            {/if}
          </div>
          <div class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-1 font-mono">
            {fiscalYear.start_date} s/d {fiscalYear.end_date}
          </div>
        </div>

        <!-- Quick Summary Pills -->
        <div class="flex items-center gap-2">
          <div class="text-center px-3 py-1.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738]">
            <div class="text-[10px] font-semibold text-[#444746] dark:text-[#c4c7c5]">Terbuka</div>
            <div class="text-xs font-bold text-emerald-600 dark:text-emerald-400">{openCount}</div>
          </div>
          <div class="text-center px-3 py-1.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738]">
            <div class="text-[10px] font-semibold text-[#444746] dark:text-[#c4c7c5]">Terkunci</div>
            <div class="text-xs font-bold text-amber-600 dark:text-amber-400">{lockedCount}</div>
          </div>
          <div class="text-center px-3 py-1.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738]">
            <div class="text-[10px] font-semibold text-[#444746] dark:text-[#c4c7c5]">Ditutup</div>
            <div class="text-xs font-bold text-slate-600 dark:text-slate-400">{closedCount}</div>
          </div>
        </div>
      </div>

      {#if fiscalYear.notes}
        <div class="p-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-xs">
          <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] block mb-0.5">Catatan Tahun Buku:</span>
          <span class="text-[#444746] dark:text-[#c4c7c5]">{fiscalYear.notes}</span>
        </div>
      {/if}

      <!-- Table of 12 Periods -->
      <div class="rounded-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5] font-semibold">
              <th class="py-2.5 px-3">Bulan</th>
              <th class="py-2.5 px-3">Kode</th>
              <th class="py-2.5 px-3">Nama Periode</th>
              <th class="py-2.5 px-3">Rentang Tanggal</th>
              <th class="py-2.5 px-3 text-center">Status</th>
              <th class="py-2.5 px-3">Keterangan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
            {#if fiscalYear.periods && fiscalYear.periods.length > 0}
              {#each fiscalYear.periods as p}
                <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors">
                  <td class="py-2 px-3 font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
                    {String(p.period_number).padStart(2, '0')}
                  </td>
                  <td class="py-2 px-3 font-mono font-bold text-primary dark:text-[#a8c7fa]">
                    {p.period_key}
                  </td>
                  <td class="py-2 px-3 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                    {p.name}
                  </td>
                  <td class="py-2 px-3 text-[#444746] dark:text-[#c4c7c5] font-mono text-[11px]">
                    {p.start_date} s/d {p.end_date}
                  </td>
                  <td class="py-2 px-3 text-center">
                    <FiscalPeriodStatusBadge status={p.status} />
                  </td>
                  <td class="py-2 px-3 text-[#444746] dark:text-[#c4c7c5] text-[11px] truncate max-w-xs">
                    {p.notes || '-'}
                  </td>
                </tr>
              {/each}
            {:else}
              <tr>
                <td colspan="6" class="py-6 text-center text-xs text-[#444746] dark:text-[#c4c7c5]">
                  Tidak ada periode yang dimuat.
                </td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <div class="flex items-center justify-end w-full">
      <M3Button variant="filled" onclick={onClose}>
        Tutup
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
