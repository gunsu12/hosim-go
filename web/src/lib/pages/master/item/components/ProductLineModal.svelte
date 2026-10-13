<script lang="ts">
  import M3Modal from "../../../../components/m3/M3Modal.svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import CoaLookupInput from "../../../../components/CoaLookupInput.svelte";

  interface Props {
    isOpen: boolean;
    mode: "CREATE" | "EDIT";
    code: string;
    name: string;
    inventoryCoa: string;
    cogsCoa: string;
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
    inventoryCoa = $bindable(""),
    cogsCoa = $bindable(""),
    isActive = $bindable(true),
    isLoading,
    onClose,
    onSubmit,
  }: Props = $props();
</script>

<M3Modal
  {isOpen}
  title={mode === "CREATE"
    ? "Tambah Lini Produk Persediaan Baru"
    : "Edit Lini Produk Persediaan"}
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
      <label for="pl-code" class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1">
        Kode Lini Produk *
      </label>
      <input
        id="pl-code"
        type="text"
        bind:value={code}
        placeholder="Contoh: PL-OBAT-FARMASI"
        required
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] font-mono text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <div>
      <label for="pl-name" class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1">
        Nama Lini Produk *
      </label>
      <input
        id="pl-name"
        type="text"
        bind:value={name}
        placeholder="Contoh: Persediaan Obat Paten Farmasi"
        required
        class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
      />
    </div>

    <CoaLookupInput
      id="pl-inventory-coa"
      label="COA Akun Persediaan (Aset)"
      bind:value={inventoryCoa}
      placeholder="Cari atau pilih akun persediaan (114.xx)..."
      preferredType="ASSET"
      helperText="Akun neraca aset untuk mutasi dan saldo stok fisik gudang/depo."
    />

    <CoaLookupInput
      id="pl-cogs-coa"
      label="COA Akun Beban Pokok (HPP)"
      bind:value={cogsCoa}
      placeholder="Cari atau pilih akun beban pokok / HPP (5xx)..."
      preferredType="EXPENSE"
      helperText="Akun laba/rugi yang didebit saat stok dikeluarkan/terjual."
    />

    {#if mode === "EDIT"}
      <div class="pt-1">
        <M3ToggleCard
          title="Status Lini Produk Aktif"
          description="Lini produk aktif dapat dipilih untuk item persediaan baru"
          bind:checked={isActive}
        />
      </div>
    {/if}

    <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea] dark:border-[#444746]">
      <M3Button variant="text" type="button" onclick={onClose}>
        Batal
      </M3Button>
      <M3Button variant="filled" type="submit" disabled={isLoading}>
        <span>{mode === "CREATE" ? "Simpan Lini Produk" : "Perbarui Lini Produk"}</span>
      </M3Button>
    </div>
  </form>
</M3Modal>
