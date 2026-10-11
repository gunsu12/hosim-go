<script lang="ts">
  import {
    Receipt,
    Layers,
    Tag,
    CheckCircle2,
    AlertCircle,
    X,
    BookOpen,
    Calculator,
    Sliders
  } from '@lucide/svelte';

  // Subcomponents
  import PricePlanTab from './components/PricePlanTab.svelte';
  import TariffClassTab from './components/TariffClassTab.svelte';
  import TariffComponentTab from './components/TariffComponentTab.svelte';
  import TariffLookupSimulator from './components/TariffLookupSimulator.svelte';
  import TariffMatrixEditor from './components/TariffMatrixEditor.svelte';

  interface Props {
    initialTab?: 'PLAN' | 'MATRIX' | 'SIMULATOR' | 'CLASS' | 'COMPONENT';
  }

  let { initialTab = 'PLAN' }: Props = $props();

  // Active Tab State
  let activeTab = $state<'PLAN' | 'MATRIX' | 'SIMULATOR' | 'CLASS' | 'COMPONENT'>('PLAN');

  $effect(() => {
    if (initialTab) {
      activeTab = initialTab;
    }
  });

  // Toast / Feedback State
  let successMessage = $state<string | null>(null);
  let errorMessage = $state<string | null>(null);

  function showToast(msg: string, isError = false) {
    if (isError) {
      errorMessage = msg;
      setTimeout(() => (errorMessage = null), 6000);
    } else {
      successMessage = msg;
      setTimeout(() => (successMessage = null), 4000);
    }
  }

  // Counts for KPI Badges
  let planCount = $state(0);
  let classCount = $state(0);
  let compCount = $state(0);
</script>

<div class="flex flex-col gap-5">
  <!-- Toast Feedback Alerts -->
  {#if successMessage}
    <div
      class="p-3.5 rounded-2xl bg-[#e6f4ea] dark:bg-emerald-950/40 border border-[#ceead6] dark:border-emerald-800/60 text-[#137333] dark:text-emerald-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2 text-xs font-semibold">
        <CheckCircle2 class="w-4 h-4 text-[#137333] dark:text-emerald-400" />
        <span>{successMessage}</span>
      </div>
      <button
        onclick={() => (successMessage = null)}
        type="button"
        class="text-[#137333] dark:text-emerald-400 hover:opacity-75 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div
      class="p-3.5 rounded-2xl bg-[#fce8e6] dark:bg-rose-950/40 border border-[#fad2cf] dark:border-rose-800/60 text-[#c5221f] dark:text-rose-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2 text-xs font-semibold">
        <AlertCircle class="w-4 h-4 text-[#c5221f] dark:text-rose-400" />
        <span>{errorMessage}</span>
      </div>
      <button
        onclick={() => (errorMessage = null)}
        type="button"
        class="text-[#c5221f] dark:text-rose-400 hover:opacity-75 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <div class="flex items-center gap-2.5">
        <div
          class="w-9 h-9 rounded-2xl bg-emerald-600/10 text-emerald-700 dark:text-emerald-400 flex items-center justify-center font-bold"
        >
          <Receipt class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-lg font-bold text-[#1f1f1f] dark:text-on-surface tracking-tight">Tarif</h2>
          <div class="flex items-center gap-2 mt-0.5">
            <span class="text-xs text-[#747775] dark:text-on-surface-variant"
              >• Buku Tarif (Price Plan), Kelas, Pecahan Komponen & Simulator</span
            >
          </div>
        </div>
      </div>
      <p class="text-xs text-[#444746] dark:text-on-surface-variant mt-1.5 max-w-3xl leading-relaxed">
        Pusat pengelolaan buku tarif pelayanan rumah sakit (standar RS & rekanan
        PKS), validasi keseimbangan pecahan komponen tagihan, siklus persetujuan
        bertahap, serta kalkulator lookup pencarian tarif real-time.
      </p>
    </div>
  </div>

  <!-- KPI / Stats Summary Cards -->
  <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
    <!-- Stat 1: Buku Tarif (Price Plans) -->
    <button
      type="button"
      onclick={() => (activeTab = 'PLAN')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab ===
      'PLAN'
        ? 'bg-indigo-50/90 border-indigo-300 ring-2 ring-indigo-400/20 dark:bg-indigo-950/50 dark:border-indigo-500/80 dark:ring-2 dark:ring-indigo-500/20 shadow-xs'
        : 'bg-white border-[#e1e5ea] hover:border-indigo-300 hover:bg-indigo-50/20 dark:bg-surface-container-low dark:border-outline-variant dark:hover:border-indigo-500/40 dark:hover:bg-indigo-950/20'}"
    >
      <div class="flex items-center justify-between text-xs font-semibold mb-1">
        <span
          class="flex items-center gap-1.5 {activeTab === 'PLAN'
            ? 'text-indigo-900 dark:text-indigo-200'
            : 'text-[#1f1f1f] dark:text-on-surface'}"
        >
          <BookOpen
            class="w-4 h-4 {activeTab === 'PLAN'
              ? 'text-indigo-600 dark:text-indigo-400'
              : 'text-indigo-500 dark:text-indigo-400'}"
          />
          Buku Tarif
        </span>
        <span
          class="text-[10px] font-mono px-1.5 py-0.5 rounded font-semibold {activeTab ===
          'PLAN'
            ? 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/80 dark:text-indigo-200 dark:border dark:border-indigo-700/60'
            : 'bg-[#f0f4f9] text-[#444746] dark:bg-surface-container dark:text-on-surface-variant'}"
          >PLAN</span
        >
      </div>
      <div
        class="text-xl font-bold {activeTab === 'PLAN'
          ? 'text-[#1f1f1f] dark:text-white'
          : 'text-[#1f1f1f] dark:text-on-surface'}"
      >
        {planCount}
      </div>
      <div
        class="text-[11px] {activeTab === 'PLAN'
          ? 'text-[#444746] dark:text-indigo-200/70'
          : 'text-[#747775] dark:text-on-surface-variant'} mt-0.5"
      >
        Daftar Buku Tarif
      </div>
    </button>

    <!-- Stat 2: Simulator & Lookup -->
    <button
      type="button"
      onclick={() => (activeTab = 'SIMULATOR')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab ===
      'SIMULATOR'
        ? 'bg-purple-50/90 border-purple-300 ring-2 ring-purple-400/20 dark:bg-purple-950/50 dark:border-purple-500/80 dark:ring-2 dark:ring-purple-500/20 shadow-xs'
        : 'bg-white border-[#e1e5ea] hover:border-purple-300 hover:bg-purple-50/20 dark:bg-surface-container-low dark:border-outline-variant dark:hover:border-purple-500/40 dark:hover:bg-purple-950/20'}"
    >
      <div class="flex items-center justify-between text-xs font-semibold mb-1">
        <span
          class="flex items-center gap-1.5 {activeTab === 'SIMULATOR'
            ? 'text-purple-900 dark:text-purple-200'
            : 'text-[#1f1f1f] dark:text-on-surface'}"
        >
          <Calculator
            class="w-4 h-4 {activeTab === 'SIMULATOR'
              ? 'text-purple-600 dark:text-purple-400'
              : 'text-purple-500 dark:text-purple-400'}"
          />
          Lookup Engine
        </span>
        <span
          class="text-[10px] font-mono px-1.5 py-0.5 rounded font-semibold {activeTab ===
          'SIMULATOR'
            ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/80 dark:text-purple-200 dark:border dark:border-purple-700/60'
            : 'bg-[#f0f4f9] text-[#444746] dark:bg-surface-container dark:text-on-surface-variant'}"
          >SIM</span
        >
      </div>
      <div
        class="text-xl font-bold {activeTab === 'SIMULATOR'
          ? 'text-[#1f1f1f] dark:text-white'
          : 'text-[#1f1f1f] dark:text-on-surface'}"
      >
        3-Step
      </div>
      <div
        class="text-[11px] {activeTab === 'SIMULATOR'
          ? 'text-[#444746] dark:text-purple-200/70'
          : 'text-[#747775] dark:text-on-surface-variant'} mt-0.5"
      >
        Hierarki Rekanan & RS
      </div>
    </button>

    <!-- Stat 3: Kelas Tarif -->
    <button
      type="button"
      onclick={() => (activeTab = 'CLASS')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab ===
      'CLASS'
        ? 'bg-blue-50/90 border-blue-300 ring-2 ring-blue-400/20 dark:bg-blue-950/50 dark:border-blue-500/80 dark:ring-2 dark:ring-blue-500/20 shadow-xs'
        : 'bg-white border-[#e1e5ea] hover:border-blue-300 hover:bg-blue-50/20 dark:bg-surface-container-low dark:border-outline-variant dark:hover:border-blue-500/40 dark:hover:bg-blue-950/20'}"
    >
      <div class="flex items-center justify-between text-xs font-semibold mb-1">
        <span
          class="flex items-center gap-1.5 {activeTab === 'CLASS'
            ? 'text-blue-900 dark:text-blue-200'
            : 'text-[#1f1f1f] dark:text-on-surface'}"
        >
          <Tag
            class="w-4 h-4 {activeTab === 'CLASS'
              ? 'text-blue-600 dark:text-blue-400'
              : 'text-blue-500 dark:text-blue-400'}"
          />
          Kelas Tarif
        </span>
        <span
          class="text-[10px] font-mono px-1.5 py-0.5 rounded font-semibold {activeTab ===
          'CLASS'
            ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/80 dark:text-blue-200 dark:border dark:border-blue-700/60'
            : 'bg-[#f0f4f9] text-[#444746] dark:bg-surface-container dark:text-on-surface-variant'}"
          >CLASS</span
        >
      </div>
      <div
        class="text-xl font-bold {activeTab === 'CLASS'
          ? 'text-[#1f1f1f] dark:text-white'
          : 'text-[#1f1f1f] dark:text-on-surface'}"
      >
        {classCount}
      </div>
      <div
        class="text-[11px] {activeTab === 'CLASS'
          ? 'text-[#444746] dark:text-blue-200/70'
          : 'text-[#747775] dark:text-on-surface-variant'} mt-0.5"
      >
        Kelas Perawatan RS
      </div>
    </button>

    <!-- Stat 4: Total Komponen -->
    <button
      type="button"
      onclick={() => (activeTab = 'COMPONENT')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab ===
      'COMPONENT'
        ? 'bg-emerald-50/90 border-emerald-300 ring-2 ring-emerald-400/20 dark:bg-emerald-950/50 dark:border-emerald-500/80 dark:ring-2 dark:ring-emerald-500/20 shadow-xs'
        : 'bg-white border-[#e1e5ea] hover:border-emerald-300 hover:bg-emerald-50/20 dark:bg-surface-container-low dark:border-outline-variant dark:hover:border-emerald-500/40 dark:hover:bg-emerald-950/20'}"
    >
      <div class="flex items-center justify-between text-xs font-semibold mb-1">
        <span
          class="flex items-center gap-1.5 {activeTab === 'COMPONENT'
            ? 'text-emerald-900 dark:text-emerald-200'
            : 'text-[#1f1f1f] dark:text-on-surface'}"
        >
          <Layers
            class="w-4 h-4 {activeTab === 'COMPONENT'
              ? 'text-emerald-600 dark:text-emerald-400'
              : 'text-emerald-500 dark:text-emerald-400'}"
          />
          Komponen Tarif
        </span>
        <span
          class="text-[10px] font-mono px-1.5 py-0.5 rounded font-semibold {activeTab ===
          'COMPONENT'
            ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/80 dark:text-emerald-200 dark:border dark:border-emerald-700/60'
            : 'bg-[#f0f4f9] text-[#444746] dark:bg-surface-container dark:text-on-surface-variant'}"
          >COMP</span
        >
      </div>
      <div
        class="text-xl font-bold {activeTab === 'COMPONENT'
          ? 'text-[#1f1f1f] dark:text-white'
          : 'text-[#1f1f1f] dark:text-on-surface'}"
      >
        {compCount}
      </div>
      <div
        class="text-[11px] {activeTab === 'COMPONENT'
          ? 'text-[#444746] dark:text-emerald-200/70'
          : 'text-[#747775] dark:text-on-surface-variant'} mt-0.5"
      >
        Pecahan Jasa RS & Dokter
      </div>
    </button>
  </div>

  <!-- Navigation Tabs -->
  <div
    class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-outline-variant pb-1 gap-2 overflow-x-auto"
  >
    <div class="flex items-center gap-1.5">
      <!-- Tab 1: Price Plan -->
      <button
        type="button"
        onclick={() => (activeTab = 'PLAN')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
        'PLAN'
          ? 'bg-indigo-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-on-surface-variant hover:bg-[#f0f4f9] dark:hover:bg-surface-container-high'}"
      >
        <BookOpen class="w-3.5 h-3.5" />
        <span>Buku Tarif ({planCount})</span>
      </button>

      <!-- Tab 2: Editor Matriks Multi-Kelas -->
      <button
        type="button"
        onclick={() => (activeTab = 'MATRIX')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
        'MATRIX'
          ? 'bg-[#0b57d0] text-white shadow-xs'
          : 'text-[#444746] dark:text-on-surface-variant hover:bg-[#f0f4f9] dark:hover:bg-surface-container-high'}"
      >
        <Sliders class="w-3.5 h-3.5" />
        <span>Editor Matriks Tarif (Multi-Kelas)</span>
      </button>

      <!-- Tab 3: Lookup Simulator -->
      <button
        type="button"
        onclick={() => (activeTab = 'SIMULATOR')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
        'SIMULATOR'
          ? 'bg-purple-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-on-surface-variant hover:bg-[#f0f4f9] dark:hover:bg-surface-container-high'}"
      >
        <Calculator class="w-3.5 h-3.5" />
        <span>Simulator & Lookup Tarif</span>
      </button>

      <!-- Tab 4: Kelas Tarif -->
      <button
        type="button"
        onclick={() => (activeTab = 'CLASS')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
        'CLASS'
          ? 'bg-blue-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-on-surface-variant hover:bg-[#f0f4f9] dark:hover:bg-surface-container-high'}"
      >
        <Tag class="w-3.5 h-3.5" />
        <span>Kelas Perawatan ({classCount})</span>
      </button>

      <!-- Tab 5: Komponen Tarif -->
      <button
        type="button"
        onclick={() => (activeTab = 'COMPONENT')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
        'COMPONENT'
          ? 'bg-emerald-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-on-surface-variant hover:bg-[#f0f4f9] dark:hover:bg-surface-container-high'}"
      >
        <Layers class="w-3.5 h-3.5" />
        <span>Komponen Tarif ({compCount})</span>
      </button>
    </div>
  </div>

  <!-- Tab Content Rendering -->
  {#if activeTab === 'PLAN'}
    <PricePlanTab {showToast} onCountChange={(c) => (planCount = c)} />
  {:else if activeTab === 'MATRIX'}
    <TariffMatrixEditor {showToast} />
  {:else if activeTab === 'SIMULATOR'}
    <TariffLookupSimulator />
  {:else if activeTab === 'CLASS'}
    <TariffClassTab {showToast} onCountChange={(c) => (classCount = c)} />
  {:else}
    <TariffComponentTab {showToast} onCountChange={(c) => (compCount = c)} />
  {/if}
</div>
