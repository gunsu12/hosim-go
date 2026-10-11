<script lang="ts">
  import { Tag, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
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

<M3Modal
  {isOpen}
  title={mode === 'CREATE' ? 'Tambah Kelas Tarif Baru' : 'Edit Kelas Tarif'}
  subtitle="Klasifikasi kamar & strata tarif rawat inap/jalan"
  icon={Tag}
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{formError}</span>
    </div>
  {/if}

  <div>
    <label for="tc-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Kode Kelas <span class="text-rose-500">*</span>
    </label>
    <input
      id="tc-code"
      type="text"
      bind:value={code}
      placeholder="Contoh: VVIP, VIP, KLS-1, KLS-2"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
    <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
      Harus unik sebagai pembeda strata kamar & kelas pelayanan.
    </p>
  </div>

  <div>
    <label for="tc-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Nama Kelas Tarif <span class="text-rose-500">*</span>
    </label>
    <input
      id="tc-name"
      type="text"
      bind:value={name}
      placeholder="Contoh: Kelas 1 (KRIS Utama)"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
  </div>

  <div>
    <label for="tc-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Deskripsi / Fasilitas Kamar
    </label>
    <textarea
      id="tc-desc"
      bind:value={description}
      rows="3"
      placeholder="Keterangan fasilitas tempat tidur, penunjang medis, atau rincian kelas..."
      class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
    ></textarea>
  </div>

  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <M3ToggleCard
      title="Kelas Tarif Aktif Digunakan"
      description="Nonaktifkan bila kelas ini sudah tidak tersedia atau tidak menerima admission baru."
      bind:checked={isActive}
      colorVariant="emerald"
    />
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSaving}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSaving}>
      {#if isSaving}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>{mode === 'CREATE' ? 'Simpan Kelas Tarif' : 'Perbarui Kelas Tarif'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
