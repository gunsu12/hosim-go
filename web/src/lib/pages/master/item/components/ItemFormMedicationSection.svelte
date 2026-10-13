<script lang="ts">
  import { Pill } from "@lucide/svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import type { MedicationType } from "../../../../types/master/item";

  interface Props {
    kfaCode: string;
    bpomNie: string;
    dosageForm: string;
    strengthAmount: string;
    strengthUnit: string;
    defaultRoute: string;
    medicationType: MedicationType;
    storageTemp: string;
    isHighAlert: boolean;
    isLasa: boolean;
    isFornas: boolean;
    isAntibiotic: boolean;
  }

  let {
    kfaCode = $bindable(""),
    bpomNie = $bindable(""),
    dosageForm = $bindable(""),
    strengthAmount = $bindable(""),
    strengthUnit = $bindable("mg"),
    defaultRoute = $bindable("Oral"),
    medicationType = $bindable("Keras"),
    storageTemp = $bindable("Suhu Ruang (15-25C)"),
    isHighAlert = $bindable(false),
    isLasa = $bindable(false),
    isFornas = $bindable(true),
    isAntibiotic = $bindable(false),
  }: Props = $props();

  const medicationTypeOptions = [
    { value: "Bebas", label: "Obat Bebas" },
    { value: "Keras", label: "Obat Keras (G)" },
    { value: "Narkotika", label: "Narkotika (O)" },
    { value: "Psikotropika", label: "Psikotropika" },
    { value: "Prekursor", label: "Prekursor Farmasi" },
    { value: "Lainnya", label: "Lainnya" },
  ];
</script>

<div
  class="p-4 rounded-2xl bg-emerald-50/40 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-800/60 flex flex-col gap-3.5 text-xs"
>
  <span
    class="font-bold text-emerald-900 dark:text-emerald-300 flex items-center gap-1.5"
  >
    <Pill class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
    <span>2. Atribut Farmasi & Regulasi Obat</span>
  </span>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div>
      <label
        for="form-med-kfa"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Kode KFA SATUSEHAT</label
      >
      <input
        id="form-med-kfa"
        type="text"
        bind:value={kfaCode}
        placeholder="93000123"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs font-mono focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <label
        for="form-med-bpom"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Nomor Izin BPOM (NIE)</label
      >
      <input
        id="form-med-bpom"
        type="text"
        bind:value={bpomNie}
        placeholder="DKL1234567809A1"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs font-mono focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <label
        for="form-med-dosage"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Bentuk Sediaan</label
      >
      <input
        id="form-med-dosage"
        type="text"
        bind:value={dosageForm}
        placeholder="Kaplet / Sirup / Injeksi"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 items-end">
    <div>
      <label
        for="form-med-strength"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Jumlah Kekuatan</label
      >
      <input
        id="form-med-strength"
        type="text"
        bind:value={strengthAmount}
        placeholder="500"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <label
        for="form-med-str-unit"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Satuan Kekuatan</label
      >
      <input
        id="form-med-str-unit"
        type="text"
        bind:value={strengthUnit}
        placeholder="mg / mcg / ml"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <M3Select
        label="Golongan Obat"
        bind:value={medicationType}
        options={medicationTypeOptions}
        placeholder="Pilih Golongan"
      />
    </div>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
    <div>
      <label
        for="form-med-route"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Rute Pemberian Default</label
      >
      <input
        id="form-med-route"
        type="text"
        bind:value={defaultRoute}
        placeholder="Oral / IV / IM / Topikal"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
    <div>
      <label
        for="form-med-storage"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Suhu Penyimpanan</label
      >
      <input
        id="form-med-storage"
        type="text"
        bind:value={storageTemp}
        placeholder="Suhu Ruang (15-25C) / Kulkas (2-8C)"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-emerald-200 dark:border-emerald-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>
  </div>

  <!-- Safety & Regulatory Toggles -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1">
    <M3ToggleCard
      title="High-Alert Medication (HAM)"
      description="Obat kewaspadaan tinggi dengan risiko bahaya fatal jika terjadi kelalaian"
      bind:checked={isHighAlert}
    />
    <M3ToggleCard
      title="LASA / NORUM"
      description="Look-Alike Sound-Alike / Nama Obat Rupa dan Ucapan Mirip"
      bind:checked={isLasa}
    />
    <M3ToggleCard
      title="Formularium Nasional (Fornas)"
      description="Tercantum dalam daftar obat acuan klaim jaminan BPJS / RS"
      bind:checked={isFornas}
    />
    <M3ToggleCard
      title="Antibiotik"
      description="Wajib monitoring evaluasi penggunaan antibiotik dan restriksi resep"
      bind:checked={isAntibiotic}
    />
  </div>
</div>
