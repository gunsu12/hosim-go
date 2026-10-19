<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Layers,
    Plus,
    Trash2,
    ArrowUp,
    ArrowDown,
    AlertCircle,
    Check,
    Star,
    Sparkles,
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select, { type SelectOption } from '$lib/components/m3/M3Select.svelte';
  import type { CustomerContractRecord } from '$lib/types/finance/customer_contract';
  import type { PriceSchemaRecord } from '$lib/types/finance/price_schema';
  import { getPriceSchemas } from '$lib/api/finance/price_schema';

  interface Props {
    isOpen: boolean;
    contract: CustomerContractRecord | null;
    isSubmitting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onSubmit: (
      schemas: Array<{
        price_schema_id: string;
        priority: number;
        is_active: boolean;
        price_schema_name?: string;
        price_schema_code?: string;
      }>,
    ) => Promise<void> | void;
  }

  let {
    isOpen,
    contract,
    isSubmitting = false,
    errorMessage = null,
    onClose,
    onSubmit,
  }: Props = $props();

  interface SelectedItem {
    price_schema_id: string;
    code: string;
    name: string;
    priority: number;
    is_active: boolean;
  }

  let localSchemas = $state<SelectedItem[]>([]);
  let availableSchemas = $state<PriceSchemaRecord[]>([]);
  let isLoadingSchemas = $state(false);
  let selectedSchemaToAdd = $state('');

  async function loadAvailablePriceSchemas() {
    isLoadingSchemas = true;
    try {
      const res = await getPriceSchemas({ is_active: true, limit: 100 });
      availableSchemas = res.items || [];
    } catch (e) {
      console.warn('Gagal memuat daftar price schemas', e);
    } finally {
      isLoadingSchemas = false;
    }
  }

  $effect(() => {
    if (isOpen) {
      loadAvailablePriceSchemas();
      if (contract) {
        localSchemas = (contract.price_schemas || [])
          .slice()
          .sort((a, b) => a.priority - b.priority)
          .map((s, idx) => ({
            price_schema_id: s.price_schema_id,
            code: s.price_schema_code || 'SCHEMA',
            name: s.price_schema_name || 'Skema Tarif',
            priority: s.priority || idx + 1,
            is_active: s.is_active,
          }));
      } else {
        localSchemas = [];
      }
      selectedSchemaToAdd = '';
    }
  });

  const availableOptions = $derived.by(() => {
    const existingIds = new Set(localSchemas.map((s) => s.price_schema_id));
    return availableSchemas
      .filter((s) => !existingIds.has(s.id))
      .map((s) => ({
        value: s.id,
        label: `${s.code} • ${s.name}`,
        description: s.description || 'Skema tarif aktif',
      }));
  });

  function handleAddSchema() {
    if (!selectedSchemaToAdd) return;
    const found = availableSchemas.find((s) => s.id === selectedSchemaToAdd);
    if (!found) return;

    localSchemas.push({
      price_schema_id: found.id,
      code: found.code,
      name: found.name,
      priority: localSchemas.length + 1,
      is_active: true,
    });
    selectedSchemaToAdd = '';
  }

  function handleRemoveSchema(index: number) {
    localSchemas.splice(index, 1);
    // Recalculate priorities
    localSchemas.forEach((s, idx) => {
      s.priority = idx + 1;
    });
  }

  function handleMoveUp(index: number) {
    if (index <= 0) return;
    const temp = localSchemas[index];
    localSchemas[index] = localSchemas[index - 1];
    localSchemas[index - 1] = temp;
    localSchemas.forEach((s, idx) => {
      s.priority = idx + 1;
    });
  }

  function handleMoveDown(index: number) {
    if (index >= localSchemas.length - 1) return;
    const temp = localSchemas[index];
    localSchemas[index] = localSchemas[index + 1];
    localSchemas[index + 1] = temp;
    localSchemas.forEach((s, idx) => {
      s.priority = idx + 1;
    });
  }

  function handleSave() {
    const activeCount = localSchemas.filter((s) => s.is_active).length;
    if (contract?.status === 'ACTIVE' && activeCount === 0) {
      alert('Kontrak aktif wajib memiliki minimal 1 price schema yang berstatus aktif!');
      return;
    }

    const payload = localSchemas.map((s, idx) => ({
      price_schema_id: s.price_schema_id,
      priority: idx + 1,
      is_active: s.is_active,
      price_schema_name: s.name,
      price_schema_code: s.code,
    }));

    onSubmit(payload);
  }
</script>

<M3Modal
  {isOpen}
  title="Kelola Skema Tarif PKS"
  subtitle={contract ? `${contract.contract_number} • ${contract.name}` : ''}
  icon={Layers}
  iconClass="bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]"
  size="2xl"
  {onClose}
>
  {#if contract}
    <div class="space-y-4 text-xs">
      <!-- Error Message in Modal -->
      {#if errorMessage}
        <div
          class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs flex items-start gap-2.5 shadow-2xs"
        >
          <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <span class="font-medium leading-relaxed">{errorMessage}</span>
        </div>
      {/if}

      <!-- Explanation Box -->
      <div
        class="p-3.5 rounded-2xl bg-[#f0f4f9] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] leading-relaxed text-[#444746] dark:text-[#c4c7c5]"
      >
        <p class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] mb-0.5">
          Aturan Urutan Prioritas (Priority Index):
        </p>
        <p>
          Skema tarif pada <strong>Prioritas 1</strong> akan otomatis menjadi pilihan <em>default</em>
          saat pasien mendaftar di loket registrasi rawat jalan/inap. Anda dapat menaik-turunkan urutan
          menggunakan tombol panah.
        </p>
      </div>

      <!-- Add Schema Bar -->
      <div
        class="p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] flex flex-col sm:flex-row items-center gap-2.5"
      >
        <div class="flex-1 w-full">
          <M3Select
            options={availableOptions}
            bind:value={selectedSchemaToAdd}
            placeholder={isLoadingSchemas
              ? 'Memuat skema tarif...'
              : availableOptions.length > 0
                ? 'Pilih skema tarif untuk ditautkan...'
                : 'Semua skema tarif telah ditautkan'}
            disabled={isLoadingSchemas || availableOptions.length === 0}
          />
        </div>
        <M3Button
          variant="filled"
          disabled={!selectedSchemaToAdd}
          onclick={handleAddSchema}
          class="shrink-0 h-10 w-full sm:w-auto"
        >
          <Plus class="w-4 h-4 mr-1" />
          <span>Tautkan Skema</span>
        </M3Button>
      </div>

      <!-- Schemas List -->
      <div class="space-y-2">
        <div class="flex items-center justify-between px-1 text-[11px] font-semibold text-[#747775] dark:text-[#8e918f]">
          <span>DAFTAR SKEMA TARIF TERHUBUNG ({localSchemas.length})</span>
          <span>PRIORITAS & STATUS</span>
        </div>

        {#if localSchemas.length === 0}
          <div
            class="p-8 text-center rounded-2xl border border-dashed border-[#e1e5ea] dark:border-[#353738] text-[#747775] dark:text-[#8e918f]"
          >
            <Layers class="w-8 h-8 mx-auto mb-2 opacity-40" />
            <p class="font-medium">Belum ada skema tarif yang terhubung ke kontrak ini.</p>
            <p class="text-[11px] mt-0.5">Pilih skema tarif di atas untuk menambahkannya.</p>
          </div>
        {:else}
          {#each localSchemas as item, idx}
            <div
              class="p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] flex items-center justify-between gap-3 shadow-2xs hover:border-[#0b57d0]/40 transition-colors"
            >
              <!-- Left info -->
              <div class="flex items-center gap-3 min-w-0">
                <div
                  class="w-8 h-8 rounded-xl font-bold flex items-center justify-center shrink-0 {idx === 0
                    ? 'bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300'
                    : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5]'}"
                >
                  {#if idx === 0}
                    <Star class="w-4 h-4 fill-amber-500 text-amber-500" />
                  {:else}
                    <span class="text-xs font-mono">{idx + 1}</span>
                  {/if}
                </div>

                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
                      {item.name}
                    </span>
                    {#if idx === 0}
                      <span
                        class="px-2 py-0.2 rounded-full text-[9px] font-bold bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300"
                      >
                        DEFAULT REGISTRASI
                      </span>
                    {/if}
                  </div>
                  <span class="text-[11px] font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa] block mt-0.5">
                    {item.code}
                  </span>
                </div>
              </div>

              <!-- Right actions -->
              <div class="flex items-center gap-2 shrink-0">
                <!-- Toggle active -->
                <label class="flex items-center gap-1.5 cursor-pointer mr-2">
                  <input
                    type="checkbox"
                    bind:checked={item.is_active}
                    class="rounded text-[#0b57d0] focus:ring-0"
                  />
                  <span class="text-[11px] {item.is_active ? 'text-emerald-700 dark:text-emerald-400 font-semibold' : 'text-gray-400'}">
                    {item.is_active ? 'Aktif' : 'Non-aktif'}
                  </span>
                </label>

                <!-- Reorder buttons -->
                <div class="flex items-center gap-1">
                  <button
                    type="button"
                    disabled={idx === 0}
                    onclick={() => handleMoveUp(idx)}
                    class="w-7 h-7 rounded-lg border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] disabled:opacity-30 cursor-pointer"
                    title="Naikkan Prioritas"
                  >
                    <ArrowUp class="w-3.5 h-3.5" />
                  </button>
                  <button
                    type="button"
                    disabled={idx === localSchemas.length - 1}
                    onclick={() => handleMoveDown(idx)}
                    class="w-7 h-7 rounded-lg border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] disabled:opacity-30 cursor-pointer"
                    title="Turunkan Prioritas"
                  >
                    <ArrowDown class="w-3.5 h-3.5" />
                  </button>
                </div>

                <!-- Remove button -->
                <button
                  type="button"
                  onclick={() => handleRemoveSchema(idx)}
                  class="w-7 h-7 rounded-lg border border-rose-200 dark:border-rose-900/60 text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 flex items-center justify-center cursor-pointer ml-1"
                  title="Hapus Tautan"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <M3Button variant="filled" disabled={isSubmitting} onclick={handleSave}>
      {#if isSubmitting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Menyimpan...</span>
      {:else}
        <Check class="w-4 h-4 mr-1" />
        <span>Simpan Skema Tarif</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
