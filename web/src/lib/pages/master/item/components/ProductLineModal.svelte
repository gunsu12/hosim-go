<script lang="ts">
  import { X } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import CoaLookupInput from '../../../../components/CoaLookupInput.svelte';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
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
    code = $bindable(''),
    name = $bindable(''),
    inventoryCoa = $bindable(''),
    cogsCoa = $bindable(''),
    isActive = $bindable(true),
    isLoading,
    onClose,
    onSubmit
  }: Props = $props();
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-md p-6 flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <h3 class="text-sm font-bold text-[#1f1f1f]">
          {mode === 'CREATE' ? 'Tambah Lini Produk Persediaan Baru' : 'Edit Lini Produk Persediaan'}
        </h3>
        <button onclick={onClose} type="button" class="text-[#747775] hover:text-[#1f1f1f]">
          <X class="w-4 h-4" />
        </button>
      </div>

      <form onsubmit={(e) => { e.preventDefault(); onSubmit(); }} class="flex flex-col gap-3 text-xs">
        <div>
          <label for="pl-code" class="block font-semibold text-[#444746] mb-1">Kode Lini Produk *</label>
          <input
            id="pl-code"
            type="text"
            bind:value={code}
            placeholder="Contoh: PL-OBAT-FARMASI"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0]"
          />
        </div>
        <div>
          <label for="pl-name" class="block font-semibold text-[#444746] mb-1">Nama Lini Produk *</label>
          <input
            id="pl-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Persediaan Obat Paten Farmasi"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
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

        {#if mode === 'EDIT'}
          <div class="pt-1">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="checkbox" bind:checked={isActive} class="rounded text-[#0b57d0]" />
              <span class="font-medium text-[#1f1f1f]">Lini Produk Aktif</span>
            </label>
          </div>
        {/if}

        <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea]">
          <button
            type="button"
            onclick={onClose}
            class="px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] cursor-pointer"
          >
            Batal
          </button>
          <M3Button variant="filled" type="submit" disabled={isLoading}>
            <span>{mode === 'CREATE' ? 'Simpan Lini Produk' : 'Perbarui Lini Produk'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
