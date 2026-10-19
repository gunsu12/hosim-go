<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Calendar,
    CalendarPlus,
    CheckCircle2,
    Clock,
    ShieldAlert,
    RefreshCw,
    AlertCircle,
    Eye,
    Check,
    Layers,
    Lock
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import { listFiscalYears, getFiscalYearById } from '$lib/api/accounting';
  import type { FiscalYearRecord } from '$lib/types/accounting/fiscalPeriod';

  import CreateFiscalYearModal from './components/CreateFiscalYearModal.svelte';
  import CloseFiscalYearModal from './components/CloseFiscalYearModal.svelte';
  import FiscalYearDetailModal from './components/FiscalYearDetailModal.svelte';

  // State
  let fiscalYears = $state<FiscalYearRecord[]>([]);
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);
  let toastMessage = $state<string | null>(null);

  // Modals state
  let isCreateModalOpen = $state(false);
  let isCloseModalOpen = $state(false);
  let isDetailModalOpen = $state(false);
  let selectedYearForAction = $state<FiscalYearRecord | null>(null);

  // Permissions
  const canManageYear = $derived(auth.hasPermission('accounting:period:manage'));
  const canCloseYear = $derived(auth.hasPermission('accounting:period:close'));

  // Stats Counters
  const stats = $derived.by(() => {
    let activeCount = 0;
    let closedCount = 0;
    for (const y of fiscalYears) {
      if (y.is_closed) closedCount++;
      else activeCount++;
    }
    return {
      total: fiscalYears.length,
      active: activeCount,
      closed: closedCount
    };
  });

  const existingYears = $derived(fiscalYears.map((y) => y.year));

  onMount(async () => {
    await loadYears();
  });

  async function loadYears() {
    isLoading = true;
    errorMessage = null;
    try {
      const data = await listFiscalYears();
      fiscalYears = data;
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal memuat daftar tahun buku fiskal';
    } finally {
      isLoading = false;
    }
  }

  function showToast(msg: string) {
    toastMessage = msg;
    setTimeout(() => {
      toastMessage = null;
    }, 4000);
  }

  function openCreateModal() {
    isCreateModalOpen = true;
  }

  async function openDetailModal(fy: FiscalYearRecord) {
    try {
      // Ambil detail lengkap beserta 12 periodenya
      const detailed = await getFiscalYearById(fy.id);
      selectedYearForAction = detailed;
      isDetailModalOpen = true;
    } catch (err: any) {
      showToast(err?.message || 'Gagal memuat detail tahun fiskal');
    }
  }

  async function openCloseYearModal(fy: FiscalYearRecord) {
    try {
      const detailed = await getFiscalYearById(fy.id);
      selectedYearForAction = detailed;
      isCloseModalOpen = true;
    } catch (err: any) {
      showToast(err?.message || 'Gagal memuat status 12 periode tahun fiskal');
    }
  }

  function handleYearCreated(newYear: FiscalYearRecord) {
    fiscalYears = [newYear, ...fiscalYears].sort((a, b) => b.year - a.year);
    showToast(`Tahun buku ${newYear.year} beserta 12 periode bulanan berhasil dibuat.`);
  }

  function handleYearClosed(closedYear: FiscalYearRecord) {
    fiscalYears = fiscalYears.map((y) => (y.id === closedYear.id ? closedYear : y));
    showToast(`Tahun buku ${closedYear.year} berhasil ditutup permanen.`);
  }
</script>

<div class="space-y-4">
  <!-- Top Banner & Action Button -->
  <div class="p-4 sm:p-5 rounded-3xl bg-gradient-to-r from-indigo-500/10 via-purple-500/10 to-blue-500/10 border border-indigo-200/60 dark:border-indigo-900/40 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
    <div class="flex items-center gap-3.5">
      <div class="w-11 h-11 rounded-2xl bg-indigo-600 text-white flex items-center justify-center shadow-xs shrink-0">
        <Calendar class="w-5 h-5" />
      </div>
      <div>
        <h2 class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
          Tahun Buku Fiskal (Fiscal Years)
        </h2>
        <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">
          Master tahun kalender akuntansi rumah sakit, otomatisasi 12 bulan periode pembukuan, dan pengesahan tutup buku tahunan.
        </p>
      </div>
    </div>

    <!-- Action Button -->
    <M3Button
      variant="filled"
      disabled={!canManageYear}
      onclick={openCreateModal}
      class="bg-indigo-600 hover:bg-indigo-700 active:bg-indigo-800 text-white shrink-0"
    >
      <CalendarPlus class="w-4 h-4 mr-1.5" />
      <span>Buka Tahun Buku Baru</span>
    </M3Button>
  </div>

  <!-- Summary Stats Cards -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div class="p-4 rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-xs flex items-center justify-between">
      <div>
        <div class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">TOTAL TAHUN BUKU</div>
        <div class="text-xl font-bold text-[#1f1f1f] dark:text-[#e3e3e3] mt-0.5">{stats.total}</div>
        <div class="text-[10px] text-[#747775] mt-0.5">Terdaftar di sistem</div>
      </div>
      <div class="w-9 h-9 rounded-xl bg-blue-50 text-blue-600 dark:bg-blue-950/60 dark:text-blue-400 flex items-center justify-center">
        <Layers class="w-4 h-4" />
      </div>
    </div>

    <div class="p-4 rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-xs flex items-center justify-between">
      <div>
        <div class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">TAHUN BUKU AKTIF</div>
        <div class="text-xl font-bold text-emerald-600 dark:text-emerald-400 mt-0.5">{stats.active}</div>
        <div class="text-[10px] text-[#747775] mt-0.5">Transaksi diperbolehkan</div>
      </div>
      <div class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400 flex items-center justify-center">
        <Clock class="w-4 h-4" />
      </div>
    </div>

    <div class="p-4 rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-xs flex items-center justify-between">
      <div>
        <div class="text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5]">SUDAH TUTUP BUKU</div>
        <div class="text-xl font-bold text-slate-700 dark:text-slate-300 mt-0.5">{stats.closed}</div>
        <div class="text-[10px] text-[#747775] mt-0.5">Tutup buku permanen</div>
      </div>
      <div class="w-9 h-9 rounded-xl bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 flex items-center justify-center">
        <CheckCircle2 class="w-4 h-4" />
      </div>
    </div>
  </div>

  <!-- Error Banner -->
  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{errorMessage}</span>
      </div>
      <M3Button variant="tonal" onclick={loadYears}>Coba Lagi</M3Button>
    </div>
  {/if}

  <!-- Data Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold w-24">Tahun</th>
            <th class="py-3 px-4 font-semibold">Rentang Periode Kalender</th>
            <th class="py-3 px-4 font-semibold text-center">Status Pembukuan</th>
            <th class="py-3 px-4 font-semibold">Catatan / Keterangan</th>
            <th class="py-3 px-4 font-semibold">Jejak Tutup Buku</th>
            <th class="py-3 px-4 font-semibold text-right w-64">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#if isLoading && fiscalYears.length === 0}
            <tr>
              <td colspan="6" class="py-12 text-center text-[#444746] dark:text-[#c4c7c5]">
                <RefreshCw class="w-6 h-6 text-primary animate-spin mx-auto mb-2" />
                <p class="font-medium text-xs">Memuat daftar tahun buku...</p>
              </td>
            </tr>
          {:else if fiscalYears.length === 0}
            <tr>
              <td colspan="6" class="py-12 text-center text-[#444746] dark:text-[#c4c7c5]">
                <Calendar class="w-8 h-8 text-[#747775] mx-auto mb-2 opacity-50" />
                <p class="font-semibold text-sm">Belum ada tahun buku yang terdaftar</p>
                <p class="text-xs text-[#747775] mt-1 mb-3">Klik tombol di atas untuk membuka tahun buku fiskal baru.</p>
                {#if canManageYear}
                  <M3Button variant="tonal" onclick={openCreateModal}>
                    <CalendarPlus class="w-4 h-4 mr-1.5" />
                    <span>Buka Tahun Buku Baru</span>
                  </M3Button>
                {/if}
              </td>
            </tr>
          {:else}
            {#each fiscalYears as fy (fy.id)}
              <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/60 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-sm text-[#1f1f1f] dark:text-[#e3e3e3]">
                  {fy.year}
                </td>
                <td class="py-3 px-4 font-mono text-[11px] text-[#444746] dark:text-[#c4c7c5]">
                  {fy.start_date} s/d {fy.end_date}
                </td>
                <td class="py-3 px-4 text-center">
                  {#if fy.is_closed}
                    <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-bold bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 border border-slate-200 dark:border-slate-700">
                      <CheckCircle2 class="w-3 h-3 text-slate-500" />
                      <span>TUTUP BUKU (CLOSED)</span>
                    </span>
                  {:else}
                    <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800">
                      <Clock class="w-3 h-3 text-emerald-600 dark:text-emerald-400" />
                      <span>AKTIF (OPEN)</span>
                    </span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-[#444746] dark:text-[#c4c7c5] max-w-xs truncate text-[11px]">
                  {fy.notes || '-'}
                </td>
                <td class="py-3 px-4 text-[11px] text-[#444746] dark:text-[#c4c7c5]">
                  {#if fy.is_closed}
                    <div>
                      <span class="font-semibold text-slate-700 dark:text-slate-300">Oleh: </span>
                      <span>{fy.closed_by || 'SYSTEM'}</span>
                    </div>
                  {:else}
                    <span class="text-[#747775] italic">Tahun aktif berjalan</span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-2">
                    <!-- Tombol Detail 12 Bulan -->
                    <button
                      onclick={() => openDetailModal(fy)}
                      class="h-7 px-2.5 rounded-lg text-[11px] font-semibold bg-[#f0f4f9] dark:bg-[#28292a] text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#e1e5ea] dark:hover:bg-[#353738] transition-all cursor-pointer flex items-center gap-1 border border-[#e1e5ea] dark:border-[#353738]"
                    >
                      <Eye class="w-3.5 h-3.5 text-primary" />
                      <span>12 Periode</span>
                    </button>

                    <!-- Tombol Tutup Buku Tahunan -->
                    {#if !fy.is_closed}
                      <button
                        onclick={() => openCloseYearModal(fy)}
                        disabled={!canCloseYear}
                        title={!canCloseYear ? 'Wewenang tutup buku dibutuhkan' : 'Tutup Buku Akhir Tahun'}
                        class="h-7 px-2.5 rounded-lg text-[11px] font-bold bg-rose-50 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300 hover:bg-rose-100 border border-rose-200 dark:border-rose-800 disabled:opacity-40 transition-all cursor-pointer flex items-center gap-1"
                      >
                        <ShieldAlert class="w-3.5 h-3.5" />
                        <span>Tutup Tahun</span>
                      </button>
                    {:else}
                      <span class="h-7 px-2.5 rounded-lg text-[10px] font-bold bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-400 border border-slate-200 dark:border-slate-700 inline-flex items-center gap-1">
                        <CheckCircle2 class="w-3 h-3 text-slate-400" />
                        <span>Selesai</span>
                      </span>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Modals -->
<CreateFiscalYearModal
  isOpen={isCreateModalOpen}
  {existingYears}
  onClose={() => (isCreateModalOpen = false)}
  onSuccess={handleYearCreated}
/>

<CloseFiscalYearModal
  isOpen={isCloseModalOpen}
  fiscalYear={selectedYearForAction}
  onClose={() => (isCloseModalOpen = false)}
  onSuccess={handleYearClosed}
/>

<FiscalYearDetailModal
  isOpen={isDetailModalOpen}
  fiscalYear={selectedYearForAction}
  onClose={() => (isDetailModalOpen = false)}
/>

<!-- Floating Toast Feedback -->
{#if toastMessage}
  <div
    class="fixed bottom-6 right-6 z-50 px-4 py-3 rounded-2xl shadow-xl flex items-center gap-3 text-xs font-semibold animate-in slide-in-from-bottom-5 duration-200 bg-[#1e1f20] text-white dark:bg-white dark:text-[#1f1f1f]"
  >
    <Check class="w-4 h-4 text-emerald-500 shrink-0" />
    <span>{toastMessage}</span>
  </div>
{/if}
