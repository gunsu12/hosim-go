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
    Sliders,
  } from "@lucide/svelte";

  // Subcomponents
  import PricePlanTab from "./components/PricePlanTab.svelte";
  import TariffClassTab from "./components/TariffClassTab.svelte";
  import TariffComponentTab from "./components/TariffComponentTab.svelte";
  import TariffLookupSimulator from "./components/TariffLookupSimulator.svelte";
  import TariffMatrixEditor from "./components/TariffMatrixEditor.svelte";

  interface Props {
    initialTab?: "PLAN" | "MATRIX" | "SIMULATOR" | "CLASS" | "COMPONENT";
  }

  let { initialTab = "PLAN" }: Props = $props();

  // Active Tab State
  let activeTab = $state<
    "PLAN" | "MATRIX" | "SIMULATOR" | "CLASS" | "COMPONENT"
  >("PLAN");

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
        <div>
          <h2
            class="text-lg font-bold text-[#1f1f1f] dark:text-on-surface tracking-tight"
          >
            Buku Tariff dan Komponen
          </h2>
        </div>
      </div>
      <p
        class="text-xs text-[#444746] dark:text-on-surface-variant mt-1.5 max-w-3xl leading-relaxed"
      >
        Pusat pengelolaan buku tarif pelayanan rumah sakit (standar RS & rekanan
        PKS), validasi keseimbangan pecahan komponen tagihan, siklus persetujuan
        bertahap, serta kalkulator lookup pencarian tarif real-time.
      </p>
    </div>
  </div>



  <!-- Navigation Tabs -->
  <div
    class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-outline-variant pb-1 gap-2 overflow-x-auto"
  >
    <div class="flex items-center gap-1.5">
      <!-- Tab 1: Price Plan -->
      <button
        type="button"
        onclick={() => (activeTab = "PLAN")}
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
        onclick={() => (activeTab = "MATRIX")}
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
        onclick={() => (activeTab = "SIMULATOR")}
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
        onclick={() => (activeTab = "CLASS")}
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
        onclick={() => (activeTab = "COMPONENT")}
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
  {#if activeTab === "PLAN"}
    <PricePlanTab {showToast} onCountChange={(c) => (planCount = c)} />
  {:else if activeTab === "MATRIX"}
    <TariffMatrixEditor {showToast} />
  {:else if activeTab === "SIMULATOR"}
    <TariffLookupSimulator />
  {:else if activeTab === "CLASS"}
    <TariffClassTab {showToast} onCountChange={(c) => (classCount = c)} />
  {:else}
    <TariffComponentTab {showToast} onCountChange={(c) => (compCount = c)} />
  {/if}
</div>
