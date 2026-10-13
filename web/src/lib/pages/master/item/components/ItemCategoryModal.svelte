<script lang="ts">
  import M3Modal from "../../../../components/m3/M3Modal.svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import CoaLookupInput from "../../../../components/CoaLookupInput.svelte";
  import type { ItemType } from "../../../../types/master/item";

  interface Props {
    isOpen: boolean;
    mode: "CREATE" | "EDIT";
    code: string;
    name: string;
    type: ItemType;
    incomeCoa: string;
    isActive: boolean;
    isLoading: boolean;
    onClose: () => void;
    onSubmit: () => void;
  }

  let {
    isOpen,
    mode,
    code = $bindable(""),
    name = $bindable(""),
    type = $bindable("MEDICATION"),
    incomeCoa = $bindable(""),
    isActive = $bindable(true),
    isLoading,
    onClose,
    onSubmit,
  }: Props = $props();

  const typeOptions = [
    { value: "MEDICATION", label: "MEDICATION (Obat-obatan)" },
    { value: "GENERAL", label: "GENERAL (BMHP & Umum)" },
    { value: "ASSET", label: "ASSET (Barang Modal)" },
    { value: "TARIFF", label: "TARIFF (Layanan & Jasa)" },
  ];
</script>

<M3Modal
  {isOpen}
  title={mode === "CREATE" ? "Tambah Kategori Item Baru" : "Edit Kategori Item"}
  size="md"
  {onClose}
>
  <form
    onsubmit={(e) => {
      e.preventDefault();
      onSubmit();
    }}
    class="flex flex-col gap-3.5 text-xs"
  >
    <div>
      <label for="cat-code" class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1">
        Kode Kategori *
      </label>
      <input
        id="cat-code"
        type="text"
        bind:value={code}
        placeholder="Contoh: CAT-MED-PATEN"
        required
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] font-mono text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <div>
      <label for="cat-name" class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1">
        Nama Kategori *
      </label>
      <input
        id="cat-name"
        type="text"
        bind:value={name}
        placeholder="Contoh: Obat Paten & Resep"
        required
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <M3Select
      label="Subtipe Item *"
      bind:value={type}
      options={typeOptions}
      placeholder="Pilih Subtipe"
    />

    <CoaLookupInput
      id="cat-income-coa"
      label="COA Akun Pendapatan"
      bind:value={incomeCoa}
      placeholder="Cari atau pilih akun pendapatan (4xx)..."
      preferredType="REVENUE"
      helperText="Akun ini akan dikreditkan saat kasir memproses tagihan item dalam kategori ini."
    />

    {#if mode === "EDIT"}
      <div class="pt-1">
        <M3ToggleCard
          title="Status Kategori Aktif"
          description="Kategori aktif dapat dipilih untuk item baru"
          bind:checked={isActive}
        />
      </div>
    {/if}

    <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea] dark:border-[#444746]">
      <M3Button variant="text" type="button" onclick={onClose}>
        Batal
      </M3Button>
      <M3Button variant="filled" type="submit" disabled={isLoading}>
        <span>{mode === "CREATE" ? "Simpan Kategori" : "Perbarui Kategori"}</span>
      </M3Button>
    </div>
  </form>
</M3Modal>
