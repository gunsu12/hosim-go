<script lang="ts">
  import { onMount } from "svelte";
  import { Calculator, Coins, CheckCircle2, AlertCircle } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Select from "$lib/components/m3/M3Select.svelte";
  import M3ToggleCard from "$lib/components/m3/M3ToggleCard.svelte";
  import TariffLookupInput from "$lib/components/TariffLookupInput.svelte";
  import { lookupTariff } from "$lib/api/finance/price_plan";
  import { getTariffClasses } from "$lib/api/finance/tariff";
  import { getCustomers } from "$lib/api/finance/customer";
  import type { LookupTariffResponse } from "$lib/types/finance/price_plan";
  import type { TariffClassRecord } from "$lib/types/finance/tariff";
  import type { CustomerRecord } from "$lib/types/master/customer";

  let tariffClasses = $state<TariffClassRecord[]>([]);
  let customers = $state<CustomerRecord[]>([]);

  // Simulator Form Inputs
  let selectedItemId = $state("");
  let selectedClassId = $state("");
  let selectedCustomerId = $state("");
  let transactionDate = $state(new Date().toISOString().split("T")[0]);
  let isCito = $state(false);

  let isCalculating = $state(false);
  let lookupResult = $state<LookupTariffResponse | null>(null);
  let lookupError = $state<string | null>(null);

  let classOptions = $derived(
    tariffClasses.map((tc) => ({
      value: tc.id,
      label: `[${tc.code}] ${tc.name}`,
    })),
  );

  let customerOptions = $derived([
    { value: "", label: "-- Pasien Umum (Tanpa PKS Penjamin Khusus) --" },
    ...customers.map((cust) => ({
      value: cust.id,
      label: `[${cust.code}] ${cust.name}`,
    })),
  ]);

  onMount(() => {
    Promise.all([
      getTariffClasses({ limit: 100 }),
      getCustomers({ is_active: true, limit: 100 }),
    ])
      .then(([classRes, custRes]) => {
        tariffClasses = classRes.data || [];
        customers = custRes.data || [];
        if (tariffClasses.length > 0) selectedClassId = tariffClasses[0].id;
      })
      .catch((err) => {
        console.error("Failed to load simulator data:", err);
      });
  });

  async function handleCalculate() {
    if (!selectedItemId || !selectedClassId) {
      lookupError = "Silakan pilih tindakan medis dan kelas tarif layanan.";
      return;
    }

    isCalculating = true;
    lookupError = null;
    lookupResult = null;

    try {
      const res = await lookupTariff({
        item_id: selectedItemId,
        tariff_class_id: selectedClassId,
        customer_id: selectedCustomerId || undefined,
        transaction_date: transactionDate,
        is_cito: isCito,
      });
      lookupResult = res;
    } catch (err: any) {
      lookupError =
        err.message ||
        "Tarif belum dikonfigurasi untuk kombinasi tindakan, kelas, atau tanggal ini.";
    } finally {
      isCalculating = false;
    }
  }
</script>

<div class="flex flex-col gap-6">
  <!-- Simulator Banner & Guidance -->
  <div
    class="p-5 rounded-3xl bg-linear-to-r from-emerald-500/10 via-teal-500/10 to-blue-500/10 border border-emerald-200 dark:border-emerald-800/60 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
  >
    <div class="flex items-center gap-3.5">
      <div
        class="w-12 h-12 rounded-2xl bg-emerald-600 text-white flex items-center justify-center shadow-md shrink-0"
      >
        <Calculator class="w-6 h-6" />
      </div>
      <div>
        <h3 class="text-sm font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
          Simulator & Engine Lookup Tarif Transaksi
        </h3>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
    <!-- Form Inputs Panel -->
    <div
      class="lg:col-span-5 bg-white dark:bg-[#1e1f20] p-5 rounded-3xl border border-[#e1e5ea] dark:border-[#444746] shadow-xs flex flex-col gap-4 text-xs"
    >
      <h4
        class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-sm flex items-center gap-2 pb-2 border-b border-[#e1e5ea] dark:border-[#444746]"
      >
        <Coins class="w-4 h-4 text-[#0b57d0] dark:text-blue-400" />
        <span>Parameter Pencarian Tarif</span>
      </h4>

      <div>
        <TariffLookupInput
          id="sim-tariff-lookup"
          label="Tindakan Medis / Tarif Layanan"
          required={true}
          bind:value={selectedItemId}
          placeholder="Cari tindakan medis (misal: Konsultasi, USG)..."
        />
      </div>

      <M3Select
        label="Kelas Layanan Perawatan *"
        bind:value={selectedClassId}
        options={classOptions}
        placeholder="Pilih Kelas"
      />

      <M3Select
        label="Penjamin / Asuransi Pasien (Opsional)"
        bind:value={selectedCustomerId}
        options={customerOptions}
        placeholder="Pilih Penjamin"
      />

      <div>
        <label
          for="sim-date"
          class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
          >Tanggal Transaksi Layanan *</label
        >
        <input
          id="sim-date"
          type="date"
          bind:value={transactionDate}
          class="w-full h-10 px-3 rounded-xl border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#2b2b2b] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:border-[#0b57d0] focus:outline-hidden"
        />
      </div>

      <div class="pt-1">
        <M3ToggleCard
          title="Tindakan Darurat / CITO"
          description="Terapkan kenaikan persentase CITO pada tarif tindakan & komponen"
          bind:checked={isCito}
        />
      </div>

      <div class="pt-2">
        <M3Button
          variant="filled"
          onclick={handleCalculate}
          disabled={isCalculating}
        >
          <div class="flex items-center justify-center gap-2 py-1 w-full">
            <Calculator class="w-4 h-4" />
            <span
              >{isCalculating
                ? "Menghitung Tarif..."
                : "Cek & Hitung Tarif"}</span
            >
          </div>
        </M3Button>
      </div>
    </div>

    <!-- Results Panel -->
    <div class="lg:col-span-7 flex flex-col gap-4">
      {#if lookupError}
        <div
          class="p-5 rounded-3xl bg-[#fce8e6] dark:bg-rose-950/40 border border-[#fad2cf] dark:border-rose-800 text-[#c5221f] dark:text-rose-300 text-xs flex items-start gap-3"
        >
          <AlertCircle class="w-5 h-5 shrink-0 mt-0.5" />
          <div>
            <h5 class="font-bold text-sm">Gagal Menemukan Konfigurasi Tarif</h5>
            <p class="mt-1 leading-relaxed">{lookupError}</p>
            <p class="text-[11px] text-[#747775] dark:text-[#c4c7c5] mt-2">
              Pastikan terdapat Buku Tarif berstatus <strong>ACTIVE</strong> yang
              memuat kombinasi tindakan dan kelas ini, atau atur buku tarif default
              RS.
            </p>
          </div>
        </div>
      {:else if lookupResult}
        <!-- Success Result Card -->
        <div
          class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#444746] shadow-xs p-5 flex flex-col gap-4 text-xs animate-in fade-in duration-200"
        >
          <div
            class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-[#444746] pb-3"
          >
            <div class="flex items-center gap-2">
              <CheckCircle2
                class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
              />
              <h4 class="font-bold text-sm text-[#1f1f1f] dark:text-[#e3e3e3]">
                Hasil Resolusi Tarif Transaksi
              </h4>
            </div>
            <span
              class="px-2.5 py-0.5 rounded-full text-[11px] font-bold border {lookupResult.is_custom_plan
                ? 'bg-purple-50 dark:bg-purple-950/40 text-purple-700 dark:text-purple-300 border-purple-200 dark:border-purple-800'
                : 'bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 border-blue-200 dark:border-blue-800'}"
            >
              {lookupResult.is_custom_plan
                ? "Custom Plan Penjamin (Langkah 1)"
                : "Standar Default RS (Langkah 2)"}
            </span>
          </div>

          <!-- Total Price Highlight Banner -->
          <div
            class="p-4 rounded-2xl bg-linear-to-r from-emerald-600 to-teal-700 text-white flex items-center justify-between shadow-xs"
          >
            <div>
              <span class="text-[11px] opacity-80 block"
                >Total Tagihan Tarif Pasien</span
              >
              <div class="text-2xl font-mono font-bold mt-0.5">
                Rp {lookupResult.total_price.toLocaleString("id-ID")}
              </div>
            </div>
            <div class="text-right">
              <span
                class="px-2.5 py-1 rounded-lg bg-white/20 font-mono font-semibold text-xs"
              >
                {lookupResult.is_cito ? "STATUS: CITO" : "STATUS: REGULER"}
              </span>
            </div>
          </div>

          <!-- Plan Info Grid -->
          <div
            class="grid grid-cols-2 gap-3 p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#282a2c] border border-[#e1e5ea] dark:border-[#444746]"
          >
            <div>
              <span class="text-[#747775] dark:text-[#c4c7c5] text-[11px] block"
                >Kode Buku Tarif Acuan</span
              >
              <span
                class="font-mono font-bold text-[#0b57d0] dark:text-blue-400"
                >{lookupResult.price_plan_code}</span
              >
            </div>
            <div>
              <span class="text-[#747775] dark:text-[#c4c7c5] text-[11px] block"
                >Kelas Layanan Perawatan</span
              >
              <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]"
                >{lookupResult.tariff_class_name ||
                  lookupResult.tariff_class_code}</span
              >
            </div>
            <div class="col-span-2">
              <span class="text-[#747775] dark:text-[#c4c7c5] text-[11px] block"
                >Tindakan Medis Terpilih</span
              >
              <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]"
                >{lookupResult.item_name} [{lookupResult.item_code}]</span
              >
            </div>
          </div>

          <!-- Breakdown Component Table -->
          <div>
            <h5
              class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-2"
            >
              Pecahan Rincian Komponen Biaya & COA:
            </h5>
            <div
              class="border border-[#e1e5ea] dark:border-[#444746] rounded-2xl overflow-hidden"
            >
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr
                    class="bg-[#f0f4f9] dark:bg-[#282a2c] text-[#444746] dark:text-[#c4c7c5] font-semibold border-b border-[#e1e5ea] dark:border-[#444746]"
                  >
                    <th class="py-2.5 px-3">Nama Pos Komponen</th>
                    <th class="py-2.5 px-3">Tipe Pos</th>
                    <th class="py-2.5 px-3">Bagan Akun (COA)</th>
                    <th class="py-2.5 px-3 text-right">Nominal Akhir</th>
                  </tr>
                </thead>
                <tbody
                  class="divide-y divide-[#e1e5ea] dark:divide-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3]"
                >
                  {#each lookupResult.components as c (c.component_id)}
                    <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#2b2b2b]">
                      <td class="py-2.5 px-3 font-semibold">
                        {c.component_name}
                      </td>
                      <td class="py-2.5 px-3">
                        <span
                          class="px-2 py-0.5 rounded text-[10px] font-semibold bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300"
                        >
                          {c.component_type}
                        </span>
                      </td>
                      <td
                        class="py-2.5 px-3 font-mono text-[#0b57d0] dark:text-blue-400"
                      >
                        {c.coa_code || "-"}
                      </td>
                      <td class="py-2.5 px-3 text-right font-mono font-bold">
                        Rp {c.amount.toLocaleString("id-ID")}
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      {:else}
        <!-- Initial Blank State Placeholder -->
        <div
          class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#444746] p-12 text-center text-[#747775] dark:text-[#c4c7c5] flex flex-col items-center justify-center gap-3"
        >
          <Calculator class="w-12 h-12 text-[#b0b8c4] dark:text-[#555]" />
          <h5 class="font-bold text-sm text-[#1f1f1f] dark:text-[#e3e3e3]">
            Simulator Siap Digunakan
          </h5>
          <p class="text-xs max-w-sm">
            Pilih tindakan medis dan kelas layanan di sebelah kiri, kemudian
            klik tombol <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]"
              >"Cek & Hitung Tarif"</strong
            > untuk melihat proses resolusi hierarkis.
          </p>
        </div>
      {/if}
    </div>
  </div>
</div>
