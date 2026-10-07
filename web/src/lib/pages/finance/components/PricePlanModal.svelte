<script lang="ts">
  import { onMount } from 'svelte';
  import { X, BookOpen, Building2, Calendar, Percent } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
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

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-lg p-6 flex flex-col gap-4 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-[#0b57d0]/10 text-[#0b57d0] flex items-center justify-center font-bold">
            <BookOpen class="w-4 h-4" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">
              {mode === 'CREATE' ? 'Buat Buku Tarif Rumah Sakit Baru' : 'Edit Metadata Buku Tarif'}
            </h3>
            <p class="text-[11px] text-[#747775]">Header SK Direksi & aturan tarif layanan kesehatan</p>
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

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="pp-code" class="block font-semibold text-[#444746] mb-1">Kode Dokumen / SK *</label>
            <input
              id="pp-code"
              type="text"
              bind:value={code}
              disabled={mode === 'EDIT'}
              placeholder="Misal: TPP-2026-REGULER"
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden disabled:bg-slate-100 disabled:opacity-70"
            />
          </div>

          <div>
            <label for="pp-cito" class="block font-semibold text-[#444746] mb-1">Standar CITO (% Tambahan) *</label>
            <div class="relative">
              <input
                id="pp-cito"
                type="number"
                step="any"
                min="0"
                max="100"
                bind:value={defaultCitoPercent}
                required
                class="w-full h-9 px-3 pr-8 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden"
              />
              <span class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] font-bold">%</span>
            </div>
          </div>
        </div>

        <div>
          <label for="pp-name" class="block font-semibold text-[#444746] mb-1">Nama Buku Tarif *</label>
          <input
            id="pp-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Buku Tarif Pelayanan RS Tahun 2026 (Revisi SK 102)"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="pp-eff-from" class="block font-semibold text-[#444746] mb-1">Berlaku Mulai *</label>
            <input
              id="pp-eff-from"
              type="date"
              bind:value={effectiveFrom}
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <div>
            <label for="pp-eff-to" class="block font-semibold text-[#444746] mb-1">Berlaku Sampai (Opsional)</label>
            <input
              id="pp-eff-to"
              type="date"
              bind:value={effectiveTo}
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>
        </div>

        <div>
          <label for="pp-customer" class="block font-semibold text-[#444746] mb-1">Penjamin / Rekanan Khusus</label>
          <select
            id="pp-customer"
            bind:value={customerId}
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          >
            <option value="">-- Standar Umum (Berlaku untuk Pasien Reguler / Semua Penjamin) --</option>
            {#each customers as cust (cust.id)}
              <option value={cust.id}>[{cust.code}] {cust.name}</option>
            {/each}
          </select>
          <p class="text-[11px] text-[#747775] mt-1">
            Kosongkan jika buku tarif ini adalah tarif umum RS. Pilih penjamin jika buku tarif ini adalah tarif PKS khusus rekanan.
          </p>
        </div>

        <div>
          <label for="pp-desc" class="block font-semibold text-[#444746] mb-1">Keterangan / Dasar Regulasi SK</label>
          <textarea
            id="pp-desc"
            bind:value={description}
            rows="2"
            placeholder="Nomor SK Direksi, catatan klausul kenaikan tarif, dsb..."
            class="w-full p-2.5 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden resize-none"
          ></textarea>
        </div>

        <div class="pt-1">
          <label class="flex items-center gap-2 cursor-pointer p-2.5 rounded-xl bg-[#f8fafd] border border-[#e1e5ea]">
            <input type="checkbox" bind:checked={isDefault} class="rounded text-[#0b57d0]" />
            <div>
              <span class="font-bold text-[#1f1f1f] block text-xs">Jadikan Buku Tarif Standar Rumah Sakit (Default)</span>
              <span class="text-[11px] text-[#747775]">Buku tarif ini akan menjadi fallback utama jika pasien tidak memiliki PKS khusus penjamin.</span>
            </div>
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
            <span>{mode === 'CREATE' ? 'Simpan Buku Tarif' : 'Perbarui Metadata'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
