<script lang="ts">
  import { Receipt } from "@lucide/svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import type { ChargeType } from "../../../../types/master/item";

  interface Props {
    nameAlias: string;
    chargeType: ChargeType;
    tariffNotes: string;
  }

  let {
    nameAlias = $bindable(""),
    chargeType = $bindable("TINDAKAN"),
    tariffNotes = $bindable(""),
  }: Props = $props();

  const chargeTypeOptions = [
    { value: "TINDAKAN", label: "Tindakan Medis / Operasi" },
    { value: "ADMINISTRASI", label: "Administrasi & Pendaftaran" },
    { value: "AKOMODASI", label: "Akomodasi / Kamar Rawat Inap" },
    { value: "PENUNJANG", label: "Penunjang Medis (Lab / Radiologi)" },
    { value: "LAINNYA", label: "Lainnya" },
  ];
</script>

<div
  class="p-4 rounded-2xl bg-blue-50/40 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800/60 flex flex-col gap-3.5 text-xs"
>
  <span
    class="font-bold text-blue-900 dark:text-blue-300 flex items-center gap-1.5"
  >
    <Receipt class="w-4 h-4 text-blue-600 dark:text-blue-400" />
    <span>2. Atribut Tarif Layanan & Jasa Medis</span>
  </span>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 items-end">
    <div>
      <label
        for="form-tar-alias"
        class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >Nama Alias / English Name</label
      >
      <input
        id="form-tar-alias"
        type="text"
        bind:value={nameAlias}
        placeholder="Misal: Outpatient Specialist Consultation"
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-blue-200 dark:border-blue-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <div>
      <M3Select
        label="Klasifikasi Beban Tagihan *"
        bind:value={chargeType}
        options={chargeTypeOptions}
        placeholder="Pilih Klasifikasi"
      />
    </div>
  </div>

  <div>
    <label
      for="form-tar-notes"
      class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
      >Catatan Operasional / Deskripsi Penagihan</label
    >
    <textarea
      id="form-tar-notes"
      bind:value={tariffNotes}
      rows="2"
      placeholder="Deskripsi layanan atau syarat penagihan..."
      class="w-full p-2.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-blue-200 dark:border-blue-800 text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
    ></textarea>
  </div>
</div>
