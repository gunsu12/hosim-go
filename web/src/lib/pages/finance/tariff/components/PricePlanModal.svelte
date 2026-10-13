<script lang="ts">
  import { onMount } from 'svelte';
  import { BookOpen, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select from '$lib/components/m3/M3Select.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
  import M3DatePicker from '$lib/components/m3/M3DatePicker.svelte';
  import { getCustomers } from '$lib/api/finance/customer';
  import type { CustomerRecord } from '$lib/types/master/customer';
  import type {
    PricePlanRecord,
    CreatePricePlanDTO,
    UpdatePricePlanDTO
  } from '$lib/types/finance/price_plan';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
    plan: PricePlanRecord | null;
    isSaving: boolean;
    onClose: () => void;
    onSave: (payload: CreatePricePlanDTO | UpdatePricePlanDTO) => Promise<void>;
  }

  let {
    isOpen,
    mode,
    plan,
    isSaving,
    onClose,
    onSave
  }: Props = $props();

  let code = $state('');
  let name = $state('');
  let effectiveFrom = $state('');
  let effectiveTo = $state('');
  let defaultCitoPercent = $state(25);
  let isDefault = $state(false);
  let customerId = $state<string>('');
  let description = $state('');
  let formError = $state<string | null>(null);

  let customers = $state<CustomerRecord[]>([]);

  let customerOptions = $derived([
    { value: '', label: '-- Standar Umum (Berlaku untuk Pasien Reguler / Semua Penjamin) --' },
    ...customers.map((c) => ({
      value: c.id,
      label: `[${c.code}] ${c.name}`
    }))
  ]);

  onMount(() => {
    getCustomers({ is_active: true, limit: 100 }).then((res) => {
      customers = res.data || [];
    }).catch((err) => {
      console.error('Failed to load customers for price plan modal:', err);
    });
  });

  $effect(() => {
    if (isOpen) {
      if (mode === 'EDIT' && plan) {
        code = plan.code;
        name = plan.name;
        effectiveFrom = plan.effective_from || '';
        effectiveTo = plan.effective_to || '';
        defaultCitoPercent = plan.default_cito_percent || 25;
        isDefault = plan.is_default;
        customerId = plan.customer_id || '';
        description = plan.description || '';
      } else {
        const today = new Date().toISOString().split('T')[0];
        code = '';
        name = '';
        effectiveFrom = today;
        effectiveTo = '';
        defaultCitoPercent = 25;
        isDefault = false;
        customerId = '';
        description = '';
      }
      formError = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (mode === 'CREATE' && !code.trim()) {
      formError = 'Kode buku tarif wajib diisi';
      return;
    }
    if (!name.trim()) {
      formError = 'Nama buku tarif wajib diisi';
      return;
    }
    if (!effectiveFrom) {
      formError = 'Tanggal berlaku awal (effective_from) wajib diisi';
      return;
    }

    formError = null;
    if (mode === 'CREATE') {
      await onSave({
        code: code.trim().toUpperCase(),
        name: name.trim(),
        effective_from: effectiveFrom,
        effective_to: effectiveTo || undefined,
        default_cito_percent: Number(defaultCitoPercent),
        is_default: isDefault,
        customer_id: customerId || undefined,
        description: description.trim() || undefined
      } as CreatePricePlanDTO);
    } else {
      await onSave({
        name: name.trim(),
        effective_from: effectiveFrom,
        effective_to: effectiveTo || undefined,
        default_cito_percent: Number(defaultCitoPercent),
        is_default: isDefault,
        customer_id: customerId || undefined,
        description: description.trim() || undefined
      } as UpdatePricePlanDTO);
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'CREATE' ? 'Buat Buku Tarif Rumah Sakit Baru' : 'Edit Metadata Buku Tarif'}
  subtitle="Header SK Direksi & aturan tarif layanan kesehatan"
  icon={BookOpen}
  size="xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{formError}</span>
    </div>
  {/if}

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="pp-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Kode Dokumen / SK <span class="text-rose-500">*</span>
      </label>
      <input
        id="pp-code"
        type="text"
        bind:value={code}
        disabled={mode === 'EDIT'}
        placeholder="Misal: TPP-2026-REGULER"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all disabled:opacity-60"
      />
    </div>

    <div>
      <label for="pp-cito" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Standar CITO (% Tambahan) <span class="text-rose-500">*</span>
      </label>
      <div class="relative">
        <input
          id="pp-cito"
          type="number"
          step="any"
          min="0"
          max="100"
          bind:value={defaultCitoPercent}
          required
          class="w-full h-10 px-3.5 pr-8 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
        <span class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f] font-bold">%</span>
      </div>
    </div>
  </div>

  <div>
    <label for="pp-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Nama Buku Tarif <span class="text-rose-500">*</span>
    </label>
    <input
      id="pp-name"
      type="text"
      bind:value={name}
      placeholder="Contoh: Buku Tarif Pelayanan RS Tahun 2026 (Revisi SK 102)"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <M3DatePicker
      id="pp-eff-from"
      label="Berlaku Mulai *"
      bind:value={effectiveFrom}
      required
      placeholder="Pilih tanggal mulai"
    />

    <M3DatePicker
      id="pp-eff-to"
      label="Berlaku Sampai (Opsional)"
      bind:value={effectiveTo}
      min={effectiveFrom || ''}
      placeholder="Pilih tanggal selesai"
    />
  </div>

  <div>
    <M3Select
      label="Penjamin / Rekanan Khusus"
      bind:value={customerId}
      options={customerOptions}
      placeholder="Pilih Penjamin"
    />
    <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
      Kosongkan jika buku tarif ini adalah tarif umum RS. Pilih penjamin jika buku tarif ini adalah tarif PKS khusus rekanan.
    </p>
  </div>

  <div>
    <label for="pp-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Keterangan / Dasar Regulasi SK
    </label>
    <textarea
      id="pp-desc"
      bind:value={description}
      rows="2"
      placeholder="Nomor SK Direksi, catatan klausul kenaikan tarif, dsb..."
      class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
    ></textarea>
  </div>

  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <M3ToggleCard
      title="Jadikan Buku Tarif Standar Rumah Sakit (Default)"
      description="Buku tarif ini akan menjadi fallback utama jika pasien tidak memiliki PKS khusus penjamin."
      bind:checked={isDefault}
      colorVariant="primary"
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
        <span>{mode === 'CREATE' ? 'Simpan Buku Tarif' : 'Perbarui Metadata'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
