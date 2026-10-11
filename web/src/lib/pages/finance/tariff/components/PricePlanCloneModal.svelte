<script lang="ts">
  import { onMount } from 'svelte';
  import { X, Copy, Calendar, AlertCircle } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { getCustomers } from '$lib/api/finance/customer';
  import type { CustomerRecord } from '$lib/types/master/customer';
  import type {
    PricePlanRecord,
    ClonePricePlanDTO
  } from '$lib/types/finance/price_plan';

  interface Props {
    isOpen: boolean;
    sourcePlan: PricePlanRecord | null;
    isCloning: boolean;
    onClose: () => void;
    onClone: (payload: ClonePricePlanDTO) => Promise<void>;
  }

  let {
    isOpen,
    sourcePlan,
    isCloning,
    onClose,
    onClone
  }: Props = $props();

  let newCode = $state('');
  let newName = $state('');
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
      console.error('Failed to load customers for clone modal:', err);
    });
  });

  $effect(() => {
    if (isOpen && sourcePlan) {
      const today = new Date().toISOString().split('T')[0];
      newCode = `${sourcePlan.code}-COPY`;
      newName = `${sourcePlan.name} (Salinan)`;
      effectiveFrom = today;
      effectiveTo = '';
      defaultCitoPercent = sourcePlan.default_cito_percent || 25;
      isDefault = false;
      customerId = sourcePlan.customer_id || '';
      description = `Kloning dari acuan ${sourcePlan.code}`;
      formError = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!newCode.trim()) {
      formError = 'Kode buku tarif baru wajib diisi';
      return;
    }
    if (!newName.trim()) {
      formError = 'Nama buku tarif baru wajib diisi';
      return;
    }
    if (!effectiveFrom) {
      formError = 'Tanggal berlaku awal wajib diisi';
      return;
    }

    formError = null;
    await onClone({
      new_code: newCode.trim().toUpperCase(),
      new_name: newName.trim(),
      effective_from: effectiveFrom,
      effective_to: effectiveTo || undefined,
      default_cito_percent: Number(defaultCitoPercent),
      is_default: isDefault,
      customer_id: customerId || undefined,
      description: description.trim() || undefined
    });
  }
</script>

{#if isOpen && sourcePlan}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-lg p-6 flex flex-col gap-4 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-purple-50 text-purple-700 flex items-center justify-center font-bold">
            <Copy class="w-4 h-4" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">
              Gandakan (Clone) Buku Tarif ke Draf Baru
            </h3>
            <p class="text-[11px] text-[#747775]">Menyalin seluruh matriks item dan pecahan komponen biaya secara atomik</p>
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

      <!-- Source Plan Banner -->
      <div class="p-3 rounded-2xl bg-purple-50/50 border border-purple-200 text-xs text-purple-900 flex flex-col gap-1">
        <span class="text-[11px] font-semibold text-purple-700">Buku Tarif Sumber (Acuan):</span>
        <div class="flex items-center gap-2 font-bold">
          <span class="font-mono bg-white px-2 py-0.5 rounded border border-purple-200 text-purple-800">{sourcePlan.code}</span>
          <span>{sourcePlan.name}</span>
        </div>
        <span class="text-[11px] text-[#555]">
          Status: <strong>{sourcePlan.status}</strong> • Total Item: <strong>{sourcePlan.total_items}</strong>
        </span>
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
            <label for="pp-new-code" class="block font-semibold text-[#444746] mb-1">Kode Buku Tarif Baru *</label>
            <input
              id="pp-new-code"
              type="text"
              bind:value={newCode}
              placeholder="Contoh: TPP-2027-REGULER"
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <div>
            <label for="pp-new-cito" class="block font-semibold text-[#444746] mb-1">Standar CITO (% Tambahan) *</label>
            <div class="relative">
              <input
                id="pp-new-cito"
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
          <label for="pp-new-name" class="block font-semibold text-[#444746] mb-1">Nama Buku Tarif Baru *</label>
          <input
            id="pp-new-name"
            type="text"
            bind:value={newName}
            placeholder="Contoh: Buku Tarif Pelayanan RS Tahun 2027"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="pp-new-from" class="block font-semibold text-[#444746] mb-1">Berlaku Mulai *</label>
            <input
              id="pp-new-from"
              type="date"
              bind:value={effectiveFrom}
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <div>
            <label for="pp-new-to" class="block font-semibold text-[#444746] mb-1">Berlaku Sampai (Opsional)</label>
            <input
              id="pp-new-to"
              type="date"
              bind:value={effectiveTo}
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>
        </div>

        <div>
          <label for="pp-new-cust" class="block font-semibold text-[#444746] mb-1">Penjamin / Rekanan</label>
          <select
            id="pp-new-cust"
            bind:value={customerId}
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          >
            <option value="">-- Standar Umum (Semua Pasien) --</option>
            {#each customers as cust (cust.id)}
              <option value={cust.id}>[{cust.code}] {cust.name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label for="pp-new-desc" class="block font-semibold text-[#444746] mb-1">Keterangan Perubahan</label>
          <textarea
            id="pp-new-desc"
            bind:value={description}
            rows="2"
            placeholder="Catatan penyesuaian tarif pada draf kloning baru..."
            class="w-full p-2.5 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden resize-none"
          ></textarea>
        </div>

        <div class="pt-1">
          <label class="flex items-center gap-2 cursor-pointer p-2 rounded-xl bg-[#f8fafd] border border-[#e1e5ea]">
            <input type="checkbox" bind:checked={isDefault} class="rounded text-[#0b57d0]" />
            <span class="font-medium text-[#1f1f1f]">Set sebagai Buku Tarif Default RS saat diaktifkan</span>
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
          <M3Button variant="filled" type="submit" disabled={isCloning}>
            <span>{isCloning ? 'Menggandakan...' : 'Mulai Salin ke Draf Baru'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
