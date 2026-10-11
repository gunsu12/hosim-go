<script lang="ts">
  import { X, Tag } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { TariffClassRecord, CreateTariffClassDTO, UpdateTariffClassDTO } from '$lib/types/finance/tariff';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
    tariffClass: TariffClassRecord | null;
    isSaving: boolean;
    onClose: () => void;
    onSave: (payload: CreateTariffClassDTO | UpdateTariffClassDTO) => Promise<void>;
  }

  let {
    isOpen,
    mode,
    tariffClass,
    isSaving,
    onClose,
    onSave
  }: Props = $props();

  let code = $state('');
  let name = $state('');
  let description = $state('');
  let isActive = $state(true);
  let formError = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      if (mode === 'EDIT' && tariffClass) {
        code = tariffClass.code;
        name = tariffClass.name;
        description = tariffClass.description || '';
        isActive = tariffClass.is_active;
      } else {
        code = '';
        name = '';
        description = '';
        isActive = true;
      }
      formError = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!code.trim()) {
      formError = 'Kode kelas tarif wajib diisi';
      return;
    }
    if (!name.trim()) {
      formError = 'Nama kelas tarif wajib diisi';
      return;
    }

    formError = null;
    await onSave({
      code: code.trim().toUpperCase(),
      name: name.trim(),
      description: description.trim() || undefined,
      is_active: isActive
    });
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-md p-6 flex flex-col gap-4">
      <!-- Modal Header -->
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-[#0b57d0]/10 text-[#0b57d0] flex items-center justify-center font-bold">
            <Tag class="w-4 h-4" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">
              {mode === 'CREATE' ? 'Tambah Kelas Tarif Baru' : 'Edit Kelas Tarif'}
            </h3>
            <p class="text-[11px] text-[#747775]">Klasifikasi kamar & strata tarif rawat inap/jalan</p>
          </div>
        </div>
        <button
          onclick={onClose}
          type="button"
          class="text-[#747775] hover:text-[#1f1f1f] p-1 rounded-full hover:bg-[#f0f4f9] transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body Form -->
      <form onsubmit={handleSubmit} class="flex flex-col gap-3 text-xs">
        {#if formError}
          <div class="p-2.5 rounded-xl bg-[#fce8e6] text-[#c5221f] text-xs font-semibold">
            {formError}
          </div>
        {/if}

        <div>
          <label for="tc-code" class="block font-semibold text-[#444746] mb-1">Kode Kelas *</label>
          <input
            id="tc-code"
            type="text"
            bind:value={code}
            placeholder="Contoh: VVIP, VIP, KLS-1, KLS-2"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden"
          />
        </div>

        <div>
          <label for="tc-name" class="block font-semibold text-[#444746] mb-1">Nama Kelas Tarif *</label>
          <input
            id="tc-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Kelas 1 (KRIS Utama)"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          />
        </div>

        <div>
          <label for="tc-desc" class="block font-semibold text-[#444746] mb-1">Deskripsi / Fasilitas Kamar</label>
          <textarea
            id="tc-desc"
            bind:value={description}
            rows="3"
            placeholder="Keterangan fasilitas tempat tidur, penunjang medis, atau rincian kelas..."
            class="w-full p-2.5 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden resize-none"
          ></textarea>
        </div>

        <div class="pt-1">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" bind:checked={isActive} class="rounded text-[#0b57d0]" />
            <span class="font-medium text-[#1f1f1f]">Kelas Tarif Aktif Digunakan</span>
          </label>
        </div>

        <!-- Footer Actions -->
        <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea] mt-1">
          <button
            type="button"
            onclick={onClose}
            class="px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9] transition-colors"
          >
            Batal
          </button>
          <M3Button variant="filled" type="submit" disabled={isSaving}>
            <span>{mode === 'CREATE' ? 'Simpan Kelas Tarif' : 'Perbarui Kelas Tarif'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
