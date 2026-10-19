<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Calendar,
    CalendarClock,
    Lock,
    Unlock,
    CheckCircle2,
    RefreshCw,
    AlertCircle,
    MoreVertical,
    Clock,
    Sparkles,
    Check
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select, { type SelectOption } from '$lib/components/m3/M3Select.svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import {
    listFiscalPeriods,
    listFiscalYears,
    getCurrentFiscalPeriod
  } from '$lib/api/accounting';
  import type {
    FiscalPeriodRecord,
    FiscalYearRecord,
    PeriodStatus
  } from '$lib/types/accounting/fiscalPeriod';

  import FiscalPeriodStatusBadge from './components/FiscalPeriodStatusBadge.svelte';
  import FiscalPeriodLockModal from './components/FiscalPeriodLockModal.svelte';
  import FiscalPeriodUnlockModal from './components/FiscalPeriodUnlockModal.svelte';
  import FiscalPeriodCloseModal from './components/FiscalPeriodCloseModal.svelte';

  // State
  let periods = $state<FiscalPeriodRecord[]>([]);
  let fiscalYears = $state<FiscalYearRecord[]>([]);
  let currentActivePeriod = $state<FiscalPeriodRecord | null>(null);

  let selectedYear = $state<number>(new Date().getFullYear());
  let selectedStatus = $state<string>('ALL'); // 'ALL' | 'OPEN' | 'LOCKED' | 'CLOSED'

  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);
  let toastMessage = $state<string | null>(null);

  // Modals state
  let lockModalOpen = $state(false);
  let unlockModalOpen = $state(false);
  let closeModalOpen = $state(false);
  let targetPeriod = $state<FiscalPeriodRecord | null>(null);

  // Action Menu state
  let activeMenuId = $state<string | null>(null);
  let menuStyle = $state('');

  // Permissions
  const canManagePeriod = $derived(auth.hasPermission('accounting:period:manage'));
  const canClosePeriod = $derived(auth.hasPermission('accounting:period:close'));

  // Year options for dropdown
  const yearOptions = $derived<SelectOption[]>(
    fiscalYears.map((y) => ({
      value: String(y.year),
      label: `Tahun Buku ${y.year}${y.is_closed ? ' (Tutup Buku)' : ''}`,
      badge: y.is_closed ? 'Closed' : undefined
    }))
  );

  // Stats Counters
  const stats = $derived.by(() => {
    let open = 0;
    let locked = 0;
    let closed = 0;
    for (const p of periods) {
      if (p.status === 'OPEN') open++;
      else if (p.status === 'LOCKED') locked++;
      else if (p.status === 'CLOSED') closed++;
    }
    return { total: periods.length, open, locked, closed };
  });

  // Filtered Periods
  const filteredPeriods = $derived(
    periods.filter((p) => {
      if (selectedStatus !== 'ALL' && p.status !== selectedStatus) {
        return false;
      }
      return true;
    })
  );

  onMount(async () => {
    await loadInitialData();
  });

  async function loadInitialData() {
    isLoading = true;
    errorMessage = null;
    try {
      // 1. Ambil daftar tahun buku
      const yearsData = await listFiscalYears();
      fiscalYears = yearsData;

      // 2. Ambil current active period
      try {
        currentActivePeriod = await getCurrentFiscalPeriod();
        if (currentActivePeriod) {
          const currentYear = parseInt(currentActivePeriod.period_key.split('-')[0], 10);
          if (currentYear) {
            selectedYear = currentYear;
          }
        }
      } catch {
        // Fallback jika belum ada periode aktif saat ini (database baru / belum ada tahun aktif)
        currentActivePeriod = null;
      }

      // Jika selectedYear belum terdaftar di years, ambil yang pertama
      if (yearsData.length > 0 && !yearsData.some((y) => y.year === selectedYear)) {
        selectedYear = yearsData[0].year;
      }

      // 3. Load periods untuk tahun terpilih
      await fetchPeriods();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal memuat data periode fiskal';
    } finally {
      isLoading = false;
    }
  }

  async function fetchPeriods() {
    isLoading = true;
    errorMessage = null;
    try {
      const data = await listFiscalPeriods({ year: selectedYear });
      periods = data;
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal memuat daftar periode bulanan';
    } finally {
      isLoading = false;
    }
  }

  function handleYearChange(yearStr: string) {
    const y = parseInt(yearStr, 10);
    if (y && y !== selectedYear) {
      selectedYear = y;
      fetchPeriods();
    }
  }

  function showToast(msg: string) {
    toastMessage = msg;
    setTimeout(() => {
      toastMessage = null;
    }, 4000);
  }

  // Action Menu Helpers
  function toggleMenu(id: string, e: MouseEvent) {
    e.stopPropagation();
    if (activeMenuId === id) {
      activeMenuId = null;
      return;
    }

    const button = e.currentTarget as HTMLElement;
    const rect = button.getBoundingClientRect();
    const estimatedHeight = 160;
    const gap = 4;
    const spaceBelow = window.innerHeight - rect.bottom;
    const spaceAbove = rect.top;

    let right = window.innerWidth - rect.right;
    if (right < 8) right = 8;

    if (spaceBelow < estimatedHeight && spaceAbove > spaceBelow) {
      const bottom = window.innerHeight - rect.top + gap;
      menuStyle = `bottom: ${Math.max(8, bottom)}px; right: ${right}px;`;
    } else {
      const top = rect.bottom + gap;
      menuStyle = `top: ${Math.max(8, top)}px; right: ${right}px;`;
    }

    activeMenuId = id;
  }

  function closeMenu() {
    activeMenuId = null;
  }

  function openLockModal(period: FiscalPeriodRecord) {
    closeMenu();
    targetPeriod = period;
    lockModalOpen = true;
  }

  function openUnlockModal(period: FiscalPeriodRecord) {
    closeMenu();
    targetPeriod = period;
    unlockModalOpen = true;
  }

  function openCloseModal(period: FiscalPeriodRecord) {
    closeMenu();
    targetPeriod = period;
    closeModalOpen = true;
  }

  function handlePeriodUpdated(updated: FiscalPeriodRecord) {
    periods = periods.map((p) => (p.id === updated.id ? updated : p));
    if (currentActivePeriod && currentActivePeriod.id === updated.id) {
      currentActivePeriod = updated;
    }
    showToast(`Periode ${updated.period_key} (${updated.name}) berhasil diperbarui: Status ${updated.status}`);
  }
</script>

<svelte:window
  onclick={closeMenu}
  onscroll={closeMenu}
  onresize={closeMenu}
  onkeydown={(e) => {
    if (e.key === 'Escape') closeMenu();
  }}
/>

<div class="space-y-4">
  <!-- Top Banner: Info Periode Aktif Realtime -->
  <div class="p-4 sm:p-5 rounded-3xl bg-gradient-to-r from-blue-500/10 via-indigo-500/10 to-teal-500/10 border border-blue-200/60 dark:border-blue-900/40 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
    <div class="flex items-center gap-3.5">
      <div class="w-11 h-11 rounded-2xl bg-primary text-white flex items-center justify-center shadow-xs shrink-0">
        <CalendarClock class="w-5 h-5" />
      </div>
      <div>
        <div class="flex items-center gap-2">
          <h2 class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
            Periode Fiskal Bulanan (Fiscal Periods)
          </h2>
          {#if currentActivePeriod}
            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-primary/10 text-primary dark:text-[#a8c7fa] border border-primary/20">
              <Sparkles class="w-2.5 h-2.5" />
              <span>Hari Ini: {currentActivePeriod.period_key}</span>
            </span>
          {/if}
        </div>
        <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">
          Kontrol penguncian transaksi (Lock/Unlock) rekonsiliasi kas-bank dan tutup buku permanen (Closing) PSAK RS.
        </p>
      </div>
    </div>

    <!-- Active Period Summary Card -->
    {#if currentActivePeriod}
      <div class="px-3.5 py-2 rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-xs flex items-center gap-3 text-xs shrink-0">
        <div>
          <div class="text-[10px] font-semibold text-[#444746] dark:text-[#c4c7c5]">PERIODE AKTIF BERJALAN</div>
          <div class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">{currentActivePeriod.name}</div>
        </div>
        <FiscalPeriodStatusBadge status={currentActivePeriod.status} />
      </div>
    {/if}
  </div>

  <!-- Single-Row Filter Bar & Status Pills -->
  <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-3 bg-white dark:bg-[#1e1f20] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs">
    <!-- Status Filter Pills -->
    <div class="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0 scrollbar-none">
      <button
        onclick={() => (selectedStatus = 'ALL')}
        class="h-8 px-3 rounded-full text-xs font-semibold transition-all whitespace-nowrap cursor-pointer {selectedStatus === 'ALL'
          ? 'bg-primary text-white shadow-xs'
          : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea] dark:hover:bg-[#353738]'}"
      >
        Semua ({stats.total})
      </button>
      <button
        onclick={() => (selectedStatus = 'OPEN')}
        class="h-8 px-3 rounded-full text-xs font-semibold transition-all whitespace-nowrap cursor-pointer {selectedStatus === 'OPEN'
          ? 'bg-emerald-600 text-white shadow-xs'
          : 'bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 hover:bg-emerald-100 dark:hover:bg-emerald-900/60'}"
      >
        Terbuka ({stats.open})
      </button>
      <button
        onclick={() => (selectedStatus = 'LOCKED')}
        class="h-8 px-3 rounded-full text-xs font-semibold transition-all whitespace-nowrap cursor-pointer {selectedStatus === 'LOCKED'
          ? 'bg-amber-600 text-white shadow-xs'
          : 'bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-300 hover:bg-amber-100 dark:hover:bg-amber-900/60'}"
      >
        Terkunci ({stats.locked})
      </button>
      <button
        onclick={() => (selectedStatus = 'CLOSED')}
        class="h-8 px-3 rounded-full text-xs font-semibold transition-all whitespace-nowrap cursor-pointer {selectedStatus === 'CLOSED'
          ? 'bg-slate-700 text-white shadow-xs'
          : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700/80'}"
      >
        Ditutup ({stats.closed})
      </button>
    </div>

    <!-- Year Dropdown & Refresh Button -->
    <div class="flex items-center gap-2">
      <div class="w-56">
        <M3Select
          value={String(selectedYear)}
          options={yearOptions}
          size="sm"
          onchange={handleYearChange}
          placeholder="Pilih Tahun Buku..."
        />
      </div>

      <button
        onclick={fetchPeriods}
        disabled={isLoading}
        title="Segarkan Data"
        class="w-8 h-8 rounded-full flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] transition-all disabled:opacity-50 cursor-pointer"
      >
        <RefreshCw class="w-4 h-4 {isLoading ? 'animate-spin text-primary' : ''}" />
      </button>
    </div>
  </div>

  <!-- Error Banner -->
  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{errorMessage}</span>
      </div>
      <M3Button variant="tonal" onclick={fetchPeriods}>Coba Lagi</M3Button>
    </div>
  {/if}

  <!-- Data Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold w-16 text-center">No</th>
            <th class="py-3 px-4 font-semibold">Periode Key</th>
            <th class="py-3 px-4 font-semibold">Nama Bulan & Tahun</th>
            <th class="py-3 px-4 font-semibold">Rentang Kalender</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold">Catatan / Alasan</th>
            <th class="py-3 px-4 font-semibold">Jejak Audit</th>
            <th class="py-3 px-4 font-semibold text-right w-44">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#if isLoading && periods.length === 0}
            <tr>
              <td colspan="8" class="py-12 text-center text-[#444746] dark:text-[#c4c7c5]">
                <RefreshCw class="w-6 h-6 text-primary animate-spin mx-auto mb-2" />
                <p class="font-medium text-xs">Memuat 12 periode fiskal...</p>
              </td>
            </tr>
          {:else if filteredPeriods.length === 0}
            <tr>
              <td colspan="8" class="py-12 text-center text-[#444746] dark:text-[#c4c7c5]">
                <Calendar class="w-8 h-8 text-[#747775] mx-auto mb-2 opacity-50" />
                <p class="font-semibold text-sm">Tidak ada periode fiskal yang sesuai</p>
                <p class="text-xs text-[#747775] mt-1">Coba ubah filter tahun buku atau status periode.</p>
              </td>
            </tr>
          {:else}
            {#each filteredPeriods as p (p.id)}
              <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/60 transition-colors">
                <td class="py-3 px-4 text-center font-mono font-bold text-[#747775]">
                  {String(p.period_number).padStart(2, '0')}
                </td>
                <td class="py-3 px-4 font-mono font-bold text-primary dark:text-[#a8c7fa]">
                  {p.period_key}
                </td>
                <td class="py-3 px-4 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                  {p.name}
                </td>
                <td class="py-3 px-4 font-mono text-[11px] text-[#444746] dark:text-[#c4c7c5]">
                  {p.start_date} s/d {p.end_date}
                </td>
                <td class="py-3 px-4 text-center">
                  <FiscalPeriodStatusBadge status={p.status} />
                </td>
                <td class="py-3 px-4 text-[#444746] dark:text-[#c4c7c5] max-w-xs truncate text-[11px]">
                  {p.notes || '-'}
                </td>
                <td class="py-3 px-4 text-[11px] text-[#444746] dark:text-[#c4c7c5]">
                  {#if p.status === 'CLOSED'}
                    <div>
                      <span class="font-semibold text-slate-700 dark:text-slate-300">Tutup: </span>
                      <span>{p.closed_by || 'SYSTEM'}</span>
                    </div>
                  {:else if p.status === 'LOCKED'}
                    <div>
                      <span class="font-semibold text-amber-700 dark:text-amber-300">Kunci: </span>
                      <span>{p.locked_by || 'SYSTEM'}</span>
                    </div>
                  {:else}
                    <span class="text-[#747775] italic">Belum dikunci</span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    <!-- Quick Action Button -->
                    {#if p.status === 'OPEN'}
                      <button
                        onclick={() => openLockModal(p)}
                        disabled={!canManagePeriod}
                        title={!canManagePeriod ? 'Wewenang manage dibutuhkan' : 'Kunci Sementara'}
                        class="h-7 px-2.5 rounded-lg text-[11px] font-bold bg-amber-50 text-amber-700 dark:bg-amber-950/60 dark:text-amber-300 hover:bg-amber-100 border border-amber-200 dark:border-amber-800 disabled:opacity-40 transition-all cursor-pointer flex items-center gap-1"
                      >
                        <Lock class="w-3 h-3" />
                        <span>Kunci</span>
                      </button>
                    {:else if p.status === 'LOCKED'}
                      <button
                        onclick={() => openUnlockModal(p)}
                        disabled={!canManagePeriod}
                        title={!canManagePeriod ? 'Wewenang manage dibutuhkan' : 'Buka Kunci'}
                        class="h-7 px-2.5 rounded-lg text-[11px] font-bold bg-emerald-50 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300 hover:bg-emerald-100 border border-emerald-200 dark:border-emerald-800 disabled:opacity-40 transition-all cursor-pointer flex items-center gap-1"
                      >
                        <Unlock class="w-3 h-3" />
                        <span>Buka</span>
                      </button>
                    {:else}
                      <span class="h-7 px-2 rounded-lg text-[10px] font-bold bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400 border border-slate-200 dark:border-slate-700 inline-flex items-center gap-1">
                        <CheckCircle2 class="w-3 h-3 text-slate-500" />
                        <span>Tutup Buku</span>
                      </span>
                    {/if}

                    <!-- More Options Button (Smart Flipping Dropdown) -->
                    {#if p.status !== 'CLOSED'}
                      <button
                        onclick={(e) => toggleMenu(p.id, e)}
                        class="w-7 h-7 rounded-lg flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] transition-all cursor-pointer {activeMenuId === p.id ? 'bg-[#e1e5ea] dark:bg-[#353738]' : ''}"
                      >
                        <MoreVertical class="w-3.5 h-3.5" />
                      </button>
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

<!-- Smart Flipping Fixed Dropdown Menu -->
{#if activeMenuId}
  {@const target = periods.find((p) => p.id === activeMenuId)}
  {#if target}
    <div
      style={menuStyle}
      class="fixed z-50 w-52 bg-white dark:bg-[#1e1f20] rounded-2xl shadow-xl border border-[#e1e5ea] dark:border-[#353738] p-1.5 flex flex-col gap-0.5 animate-in fade-in zoom-in-95 duration-100"
    >
      <div class="px-3 py-1.5 text-[10px] font-bold font-mono text-[#747775] border-b border-[#e1e5ea] dark:border-[#353738] mb-1">
        AKSI: {target.period_key}
      </div>

      {#if target.status === 'OPEN'}
        <button
          onclick={() => openLockModal(target)}
          disabled={!canManagePeriod}
          class="w-full text-left px-3 py-2 rounded-xl text-xs font-semibold text-amber-700 dark:text-amber-400 hover:bg-amber-50 dark:hover:bg-amber-950/40 flex items-center gap-2 transition-colors disabled:opacity-40 cursor-pointer"
        >
          <Lock class="w-3.5 h-3.5" />
          <span>Kunci Sementara (Lock)</span>
        </button>
      {:else if target.status === 'LOCKED'}
        <button
          onclick={() => openUnlockModal(target)}
          disabled={!canManagePeriod}
          class="w-full text-left px-3 py-2 rounded-xl text-xs font-semibold text-emerald-700 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-950/40 flex items-center gap-2 transition-colors disabled:opacity-40 cursor-pointer"
        >
          <Unlock class="w-3.5 h-3.5" />
          <span>Buka Kunci (Unlock)</span>
        </button>

        <button
          onclick={() => openCloseModal(target)}
          disabled={!canClosePeriod}
          class="w-full text-left px-3 py-2 rounded-xl text-xs font-semibold text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 flex items-center gap-2 transition-colors disabled:opacity-40 cursor-pointer"
        >
          <CheckCircle2 class="w-3.5 h-3.5" />
          <span>Tutup Buku (Close Period)</span>
        </button>
      {/if}
    </div>
  {/if}
{/if}

<!-- Modals -->
<FiscalPeriodLockModal
  isOpen={lockModalOpen}
  period={targetPeriod}
  onClose={() => (lockModalOpen = false)}
  onSuccess={handlePeriodUpdated}
/>

<FiscalPeriodUnlockModal
  isOpen={unlockModalOpen}
  period={targetPeriod}
  onClose={() => (unlockModalOpen = false)}
  onSuccess={handlePeriodUpdated}
/>

<FiscalPeriodCloseModal
  isOpen={closeModalOpen}
  period={targetPeriod}
  onClose={() => (closeModalOpen = false)}
  onSuccess={handlePeriodUpdated}
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
