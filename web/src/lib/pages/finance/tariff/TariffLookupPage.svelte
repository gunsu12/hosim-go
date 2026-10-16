<script lang="ts">
  import TariffLookupSimulator from "./components/TariffLookupSimulator.svelte";
</script>

<div class="flex flex-col gap-5">
  <!-- Header Domain -->
  <div
    class="flex flex-wrap items-center justify-between gap-4 border-b border-[#e1e5ea] dark:border-outline-variant pb-3"
  >
    <div>
      <div class="flex items-center gap-2.5">
        <h2
          class="text-lg font-bold text-[#1f1f1f] dark:text-on-surface tracking-tight"
        >
          Simulator & Lookup Tarif Real-Time
        </h2>
        <span
          class="text-xs px-2.5 py-0.5 rounded-full bg-purple-50 dark:bg-purple-950/60 text-purple-700 dark:text-purple-300 font-semibold border border-purple-200 dark:border-purple-800"
        >
          Engine Simulasi
        </span>
      </div>
      <p
        class="text-xs text-[#444746] dark:text-on-surface-variant mt-1.5 max-w-3xl leading-relaxed"
      >
        Kalkulator lookup tarif untuk memverifikasi kalkulasi harga tindakan berdasarkan buku tarif aktif, skema penjamin PKS, kelas rawat, dan pecahan komponen biaya secara presisi.
      </p>
    </div>
  </div>

  <TariffLookupSimulator />
</div>
