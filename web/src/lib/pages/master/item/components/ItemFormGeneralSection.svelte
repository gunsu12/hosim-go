<script lang="ts">
  import { Package } from "@lucide/svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import type {
    GeneralType,
    SterilizationMethod,
  } from "../../../../types/master/item";

  interface Props {
    generalType: GeneralType;
    sterilizationMethod: SterilizationMethod;
    isSterile: boolean;
    isDisposable: boolean;
    isCssdItem: boolean;
  }

  let {
    generalType = $bindable("BMHP_MEDIS"),
    sterilizationMethod = $bindable("STEAM_AUTOCLAVE"),
    isSterile = $bindable(false),
    isDisposable = $bindable(true),
    isCssdItem = $bindable(false),
  }: Props = $props();

  const generalTypeOptions = [
    { value: "BMHP_MEDIS", label: "BMHP Medis" },
    { value: "INSTRUMEN_MEDIS", label: "Instrumen Medis" },
    { value: "ATK", label: "Alat Tulis Kantor (ATK)" },
    { value: "LINEN", label: "Linen & Kain Medis" },
    { value: "KEBERSIHAN", label: "Bahan Kebersihan" },
    { value: "DAPUR", label: "Bahan Dapur & Gizi" },
    { value: "LAINNYA", label: "Lainnya" },
  ];

  const sterilizationOptions = [
    { value: "STEAM_AUTOCLAVE", label: "Steam Autoclave (Uap Suhu Tinggi)" },
    { value: "EO_GAS", label: "Etilen Oksida (EO Gas)" },
    { value: "PLASMA", label: "Hidrogen Peroksida Plasma" },
    { value: "DRY_HEAT", label: "Dry Heat (Panas Kering)" },
  ];
</script>

<div
  class="p-4 rounded-2xl bg-amber-50/40 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-800/60 flex flex-col gap-3.5 text-xs"
>
  <span
    class="font-bold text-amber-900 dark:text-amber-300 flex items-center gap-1.5"
  >
    <Package class="w-4 h-4 text-amber-600 dark:text-amber-400" />
    <span>2. Atribut BMHP & Siklus CSSD</span>
  </span>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 items-end">
    <M3Select
      label="Subtipe Barang Umum *"
      bind:value={generalType}
      options={generalTypeOptions}
      placeholder="Pilih Subtipe BMHP"
    />

    <M3Select
      label="Metode Sterilisasi CSSD"
      bind:value={sterilizationMethod}
      options={sterilizationOptions}
      placeholder="Pilih Metode Sterilisasi"
    />
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-2 pt-1">
    <M3ToggleCard
      title="Produk Steril"
      description="Item telah melalui proses sterilisasi sebelum pemakaian"
      bind:checked={isSterile}
    />
    <M3ToggleCard
      title="Sekali Pakai (Disposable)"
      description="Item dibuang setelah satu kali pemakaian klinis"
      bind:checked={isDisposable}
    />
    <M3ToggleCard
      title="Dikelola dalam Siklus CSSD"
      description="Item dicatat dan diproses dalam sirkulasi unit sterilisasi sentral"
      bind:checked={isCssdItem}
    />
  </div>
</div>
