<script lang="ts">
  import {
    Pill,
    Package,
    Stethoscope,
    Receipt,
    Thermometer,
    Flame,
    Info,
  } from "@lucide/svelte";
  import type { ItemDetailRecord } from "../../../../types/master/item";

  interface Props {
    detailItem: ItemDetailRecord;
  }

  let { detailItem }: Props = $props();
</script>

<!-- Identitas Umum -->
<div
  class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] grid grid-cols-2 gap-3 text-xs"
>
  <div>
    <span class="text-[#747775] dark:text-[#c4c7c5] block text-[11px]"
      >Nama Generik / Zat Aktif</span
    >
    <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
      >{detailItem.generic_name || "-"}</span
    >
  </div>
  <div>
    <span class="text-[#747775] dark:text-[#c4c7c5] block text-[11px]"
      >Kategori Item</span
    >
    <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
      >{detailItem.category?.name || "-"}</span
    >
  </div>
  <div>
    <span class="text-[#747775] dark:text-[#c4c7c5] block text-[11px]"
      >Lini Produk Persediaan</span
    >
    <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
      >{detailItem.product_line?.name || "Non-Persediaan"}</span
    >
  </div>
  <div>
    <span class="text-[#747775] dark:text-[#c4c7c5] block text-[11px]"
      >Status Katalog</span
    >
    <span
      class="font-medium {detailItem.is_active
        ? 'text-emerald-700 dark:text-emerald-400 font-semibold'
        : 'text-rose-700 dark:text-rose-400'}"
    >
      {detailItem.is_active ? "Aktif Digunakan" : "Nonaktif / Diarsipkan"}
    </span>
  </div>
</div>

<!-- Rincian Khusus OBAT (MEDICATION) -->
{#if detailItem.item_type === "MEDICATION" && detailItem.medication}
  {@const med = detailItem.medication}
  <div class="flex flex-col gap-3">
    <h4
      class="text-xs font-bold text-emerald-800 dark:text-emerald-400 flex items-center gap-1.5"
    >
      <Pill class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
      <span>Rincian Spesifikasi Farmasi & Regulasi Obat</span>
    </h4>

    <div
      class="grid grid-cols-2 sm:grid-cols-3 gap-3 p-4 rounded-2xl bg-emerald-50/40 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-800/60 text-xs"
    >
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Kode KFA SATUSEHAT</span
        >
        <span class="font-mono font-bold text-emerald-700 dark:text-emerald-400"
          >{med.kfa_code || "-"}</span
        >
      </div>
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Nomor Izin BPOM (NIE)</span
        >
        <span class="font-mono font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{med.bpom_nie || "-"}</span
        >
      </div>
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Bentuk Sediaan</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{med.dosage_form || "-"}</span
        >
      </div>
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Kekuatan Dosis</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{med.strength_amount || ""} {med.strength_unit || ""}</span
        >
      </div>
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Rute Pemberian</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{med.default_route || "-"}</span
        >
      </div>
      <div>
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Golongan Obat</span
        >
        <span class="font-semibold text-emerald-800 dark:text-emerald-300"
          >{med.medication_type || "Bebas"}</span
        >
      </div>
      <div class="col-span-2 sm:col-span-3">
        <span class="text-emerald-900 dark:text-emerald-300 block text-[11px]"
          >Suhu Penyimpanan</span
        >
        <span
          class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center gap-1.5 mt-0.5"
        >
          <Thermometer class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" />
          <span>{med.storage_temperature || "Suhu Ruang (15-25C)"}</span>
        </span>
      </div>
    </div>

    <!-- Regulatory & Safety Badges -->
    <div class="flex items-center gap-2 flex-wrap">
      {#if med.is_high_alert}
        <span
          class="px-2.5 py-1 rounded-lg bg-rose-100 dark:bg-rose-950/40 text-rose-800 dark:text-rose-300 text-xs font-semibold flex items-center gap-1 border border-rose-300 dark:border-rose-800"
        >
          <Flame class="w-3.5 h-3.5 text-rose-600 dark:text-rose-400" />
          High Alert Medication (HAM)
        </span>
      {/if}
      {#if med.is_lasa}
        <span
          class="px-2.5 py-1 rounded-lg bg-amber-100 dark:bg-amber-950/40 text-amber-800 dark:text-amber-300 text-xs font-semibold border border-amber-300 dark:border-amber-800"
        >
          LASA / NORUM
        </span>
      {/if}
      {#if med.is_fornas}
        <span
          class="px-2.5 py-1 rounded-lg bg-blue-100 dark:bg-blue-950/40 text-blue-800 dark:text-blue-300 text-xs font-semibold border border-blue-300 dark:border-blue-800"
        >
          Formularium Nasional (Fornas)
        </span>
      {/if}
      {#if med.is_antibiotic}
        <span
          class="px-2.5 py-1 rounded-lg bg-purple-100 dark:bg-purple-950/40 text-purple-800 dark:text-purple-300 text-xs font-semibold border border-purple-300 dark:border-purple-800"
        >
          Antibiotik
        </span>
      {/if}
    </div>
  </div>

<!-- Rincian Khusus BMHP / UMUM (GENERAL) -->
{:else if detailItem.item_type === "GENERAL" && detailItem.general}
  {@const gen = detailItem.general}
  <div class="flex flex-col gap-3">
    <h4
      class="text-xs font-bold text-amber-800 dark:text-amber-400 flex items-center gap-1.5"
    >
      <Package class="w-4 h-4 text-amber-600 dark:text-amber-400" />
      <span>Rincian Spesifikasi BMHP & Siklus CSSD</span>
    </h4>

    <div
      class="grid grid-cols-2 gap-3 p-4 rounded-2xl bg-amber-50/40 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-800/60 text-xs"
    >
      <div>
        <span class="text-amber-900 dark:text-amber-300 block text-[11px]"
          >Subtipe Barang Umum</span
        >
        <span class="font-bold text-amber-800 dark:text-amber-400">{gen.general_type}</span>
      </div>
      <div>
        <span class="text-amber-900 dark:text-amber-300 block text-[11px]"
          >Metode Sterilisasi CSSD</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{gen.sterilization_method || "N/A"}</span
        >
      </div>
      <div>
        <span class="text-amber-900 dark:text-amber-300 block text-[11px]"
          >Tingkat Sterilitas</span
        >
        <span
          class="font-medium {gen.is_sterile
            ? 'text-emerald-700 dark:text-emerald-400 font-bold'
            : 'text-[#747775] dark:text-[#c4c7c5]'}"
        >
          {gen.is_sterile ? "Produk Steril" : "Non-Steril"}
        </span>
      </div>
      <div>
        <span class="text-amber-900 dark:text-amber-300 block text-[11px]"
          >Karakteristik Pemakaian</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
          {gen.is_disposable
            ? "Single-Use (Sekali Pakai)"
            : "Reusable (Bisa Dipakai Ulang)"}
        </span>
      </div>
      <div>
        <span class="text-amber-900 dark:text-amber-300 block text-[11px]"
          >Pengelolaan CSSD</span
        >
        <span
          class="font-medium {gen.is_cssd_item
            ? 'text-purple-700 dark:text-purple-400 font-semibold'
            : 'text-[#747775] dark:text-[#c4c7c5]'}"
        >
          {gen.is_cssd_item ? "Dikelola dalam Siklus CSSD" : "Bukan Siklus CSSD"}
        </span>
      </div>
    </div>
  </div>

<!-- Rincian Khusus ASET MODAL (ASSET) -->
{:else if detailItem.item_type === "ASSET" && detailItem.asset}
  {@const ast = detailItem.asset}
  <div class="flex flex-col gap-3">
    <h4
      class="text-xs font-bold text-indigo-800 dark:text-indigo-400 flex items-center gap-1.5"
    >
      <Stethoscope class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
      <span>Rincian Barang Modal, Depresiasi & Kalibrasi</span>
    </h4>

    <div
      class="grid grid-cols-2 sm:grid-cols-3 gap-3 p-4 rounded-2xl bg-indigo-50/40 dark:bg-indigo-950/20 border border-indigo-200 dark:border-indigo-800/60 text-xs"
    >
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Merk / Pabrikan (Brand)</span
        >
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{ast.brand || "-"}</span
        >
      </div>
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Model / Seri Mesin</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{ast.model_name || "-"}</span
        >
      </div>
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Klasifikasi Alat</span
        >
        <span
          class="font-semibold {ast.is_medical_equipment
            ? 'text-indigo-700 dark:text-indigo-400'
            : 'text-[#747775] dark:text-[#c4c7c5]'}"
        >
          {ast.is_medical_equipment
            ? "Alat Medis Elektromedik"
            : "Aset Non-Medis"}
        </span>
      </div>
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Masa Manfaat (Umur Ekonomis)</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{ast.expected_life_years || 5} Tahun</span
        >
      </div>
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Metode Depresiasi</span
        >
        <span class="font-mono font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{ast.depreciation_method || "STRAIGHT_LINE"}</span
        >
      </div>
      <div>
        <span class="text-indigo-900 dark:text-indigo-300 block text-[11px]"
          >Interval Servis & Kalibrasi</span
        >
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
          >{ast.maintenance_interval_days || 180} Hari Sekali</span
        >
      </div>
    </div>
  </div>

<!-- Rincian Khusus TARIF LAYANAN (TARIFF) -->
{:else if detailItem.item_type === "TARIFF" && detailItem.tariff}
  {@const tar = detailItem.tariff}
  <div class="flex flex-col gap-3">
    <h4
      class="text-xs font-bold text-blue-800 dark:text-blue-400 flex items-center gap-1.5"
    >
      <Receipt class="w-4 h-4 text-blue-600 dark:text-blue-400" />
      <span>Rincian Master Tarif Layanan & Jasa Medis</span>
    </h4>

    <div
      class="p-4 rounded-2xl bg-blue-50/40 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800/60 text-xs flex flex-col gap-3"
    >
      <div class="grid grid-cols-2 gap-3">
        <div>
          <span class="text-blue-900 dark:text-blue-300 block text-[11px]"
            >Nama Alias / Bahasa Internasional</span
          >
          <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]"
            >{tar.name_alias || "-"}</span
          >
        </div>
        <div>
          <span class="text-blue-900 dark:text-blue-300 block text-[11px]"
            >Klasifikasi Beban Tagihan</span
          >
          <span class="font-bold text-blue-700 dark:text-blue-400"
            >{tar.charge_type}</span
          >
        </div>
      </div>
      <div>
        <span class="text-blue-900 dark:text-blue-300 block text-[11px]"
          >Catatan / Deskripsi Penagihan</span
        >
        <p class="text-[#444746] dark:text-[#c4c7c5] mt-0.5 leading-relaxed">
          {tar.notes || "Tidak ada catatan tambahan."}
        </p>
      </div>
      <div
        class="p-2.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-blue-200 dark:border-blue-800 text-[11px] text-blue-900 dark:text-blue-300 flex items-center gap-2"
      >
        <Info class="w-4 h-4 text-blue-600 dark:text-blue-400 shrink-0" />
        <span>
          Nominal harga, pemetaan kelas tarif (VVIP/VIP/Kls 1-3), surcharge
          cito, dan jasa sarana/medis diatur pada <strong
            >Modul Finance (Buku Tarif / Price Plan)</strong
          >.
        </span>
      </div>
    </div>
  </div>
{/if}
