<script lang="ts">
  import { CheckCircle2, AlertCircle, X } from "@lucide/svelte";
  import TariffMatrixEditor from "./components/TariffMatrixEditor.svelte";

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
  <div
    class="flex flex-wrap items-center justify-between gap-4 border-b border-[#e1e5ea] dark:border-outline-variant pb-3"
  >
    <div>
      <div class="flex items-center gap-2.5">
        <h2
          class="text-lg font-bold text-[#1f1f1f] dark:text-on-surface tracking-tight"
        >
          Editor Matriks Tarif (Multi-Kelas)
        </h2>
      </div>
      <p
        class="text-xs text-[#444746] dark:text-on-surface-variant mt-1.5 max-w-3xl leading-relaxed"
      >
        Konfigurasi batch tarif per tindakan/item di seluruh tingkatan kelas
        perawatan (VVIP hingga Kelas 3) serta validasi keseimbangan pecahan
        komponen biaya.
      </p>
    </div>
  </div>

  <TariffMatrixEditor {showToast} />
</div>
