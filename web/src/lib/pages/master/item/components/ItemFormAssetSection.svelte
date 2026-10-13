<script lang="ts">
  import { Stethoscope } from "@lucide/svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import type { DepreciationMethod } from "../../../../types/master/item";

  interface Props {
    brand: string;
    modelName: string;
    expectedLifeYears: number;
    depreciationMethod: DepreciationMethod;
    maintenanceIntervalDays: number;
    isMedicalEquipment: boolean;
  }

  let {
    brand = $bindable(""),
    modelName = $bindable(""),
    expectedLifeYears = $bindable(5),
    depreciationMethod = $bindable("STRAIGHT_LINE"),
    maintenanceIntervalDays = $bindable(180),
    isMedicalEquipment = $bindable(true),
  }: Props = $props();

  const depreciationOptions = [
    { value: "STRAIGHT_LINE", label: "Garis Lurus (Straight Line)" },
    { value: "DOUBLE_DECLINING", label: "Saldo Menurun Ganda (Double Declining)" },
    { value: "SUM_OF_YEARS", label: "Jumlah Angka Tahun (Sum of Years)" },
  ];
</script>

<div
  class="p-4 rounded-2xl bg-indigo-50/40 dark:bg-indigo-950/20 border border-indigo-200 dark:border-indigo-800/60 flex flex-col gap-3.5 text-xs"
>
  <span
    class="font-bold text-indigo-900 dark:text-indigo-300 flex items-center gap-1.5"
  >
    <Stethoscope class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
    <span>2. Atribut Barang Modal, Depresiasi & Kalibrasi</span>
  </span>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
    <div>
      <label
        for="form-ast-brand"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Merk / Pabrikan</label
      >
      <input
        id="form-ast-brand"
        type="text"
        bind:value={brand}
        placeholder="Misal: GE Healthcare, Philips, Mindray"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-indigo-200 dark:border-indigo-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <label
        for="form-ast-model"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Tipe / Model Mesin</label
      >
      <input
        id="form-ast-model"
        type="text"
        bind:value={modelName}
        placeholder="Misal: Voluson E8 / Defibrillator BeneHeart D3"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-indigo-200 dark:border-indigo-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
    <div>
      <label
        for="form-ast-life"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Estimasi Umur (Tahun)</label
      >
      <input
        id="form-ast-life"
        type="number"
        bind:value={expectedLifeYears}
        min="1"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-indigo-200 dark:border-indigo-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs font-mono focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <div>
      <M3Select
        label="Metode Depresiasi"
        bind:value={depreciationMethod}
        options={depreciationOptions}
        placeholder="Pilih Metode"
      />
    </div>

    <div>
      <label
        for="form-ast-interval"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Interval Servis / Kalibrasi (Hari)</label
      >
      <input
        id="form-ast-interval"
        type="number"
        bind:value={maintenanceIntervalDays}
        min="1"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-indigo-200 dark:border-indigo-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs font-mono focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
  </div>

  <div class="pt-1">
    <M3ToggleCard
      title="Alat Kesehatan Elektromedik / Medis"
      description="Termasuk alat elektromedik yang memerlukan sertifikat uji fungsi berkala BPFK"
      bind:checked={isMedicalEquipment}
    />
  </div>
</div>
