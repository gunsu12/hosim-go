<script lang="ts">
  import { X } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import CoaLookupInput from '../../../../components/CoaLookupInput.svelte';
  import type { ItemType } from '../../../../types/master/item';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
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
    code = $bindable(''),
    name = $bindable(''),
    type = $bindable('MEDICATION'),
    incomeCoa = $bindable(''),
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
          {mode === 'CREATE' ? 'Tambah Kategori Item Baru' : 'Edit Kategori Item'}
        </h3>
        <button onclick={onClose} type="button" class="text-[#747775] hover:text-[#1f1f1f]">
          <X class="w-4 h-4" />
        </button>
      </div>

      <form onsubmit={(e) => { e.preventDefault(); onSubmit(); }} class="flex flex-col gap-3 text-xs">
        <div>
          <label for="cat-code" class="block font-semibold text-[#444746] mb-1">Kode Kategori *</label>
          <input
            id="cat-code"
            type="text"
            bind:value={code}
            placeholder="Contoh: CAT-MED-PATEN"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0]"
          />
        </div>
        <div>
          <label for="cat-name" class="block font-semibold text-[#444746] mb-1">Nama Kategori *</label>
          <input
            id="cat-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Obat Paten & Resep"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
          />
        </div>
        <div>
          <label for="cat-type" class="block font-semibold text-[#444746] mb-1">Subtipe Item *</label>
          <select
            id="cat-type"
            bind:value={type}
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
          >
            <option value="MEDICATION">MEDICATION (Obat-obatan)</option>
            <option value="GENERAL">GENERAL (BMHP & Umum)</option>
            <option value="ASSET">ASSET (Barang Modal)</option>
            <option value="TARIFF">TARIFF (Layanan & Jasa)</option>
          </select>
        </div>
        <CoaLookupInput
          id="cat-income-coa"
          label="COA Akun Pendapatan"
          bind:value={incomeCoa}
          placeholder="Cari atau pilih akun pendapatan (4xx)..."
          preferredType="REVENUE"
          helperText="Akun ini akan dikreditkan saat kasir memproses tagihan item dalam kategori ini."
        />

        {#if mode === 'EDIT'}
          <div class="pt-1">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="checkbox" bind:checked={isActive} class="rounded text-[#0b57d0]" />
              <span class="font-medium text-[#1f1f1f]">Kategori Aktif</span>
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
            <span>{mode === 'CREATE' ? 'Simpan Kategori' : 'Perbarui Kategori'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
