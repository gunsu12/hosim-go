<script lang="ts">
  import { onMount } from 'svelte';
  import {
    X,
    Plus,
    Search,
    Trash2,
    Edit3,
    CheckCircle2,
    AlertTriangle,
    Coins,
    Layers,
    Tag,
    Calculator,
    ChevronDown,
    Building2,
    Lock
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import {
    getPricePlanItems,
    addPricePlanItem,
    updatePricePlanItem,
    deletePricePlanItem
  } from '$lib/api/finance/price_plan';
  import { getTariffClasses, getTariffComponents } from '$lib/api/finance/tariff';
  import { getItems } from '$lib/api/catalog/item';
  import type {
    PricePlanRecord,
    PricePlanItemRecord,
    ItemComponentDTO
  } from '$lib/types/finance/price_plan';
  import type { TariffClassRecord, TariffComponentRecord } from '$lib/types/finance/tariff';
  import type { ItemSummaryRecord } from '$lib/types/master/item';

  interface Props {
    isOpen: boolean;
    plan: PricePlanRecord | null;
    onClose: () => void;
    showToast: (msg: string, isError?: boolean) => void;
  }

  let {
    isOpen,
    plan,
    onClose,
    showToast
  }: Props = $props();

  let items = $state<PricePlanItemRecord[]>([]);
  let isLoading = $state(false);
  let searchQuery = $state('');
  let filterClassId = $state('');
  let currentPage = $state(1);
  let totalPages = $state(1);
  let totalCount = $state(0);

  // Reference Options
  let tariffClasses = $state<TariffClassRecord[]>([]);
  let tariffComponents = $state<TariffComponentRecord[]>([]);
  let tariffCatalogItems = $state<ItemSummaryRecord[]>([]);

  // Item Form Modal / Inline State
  let showItemForm = $state(false);
  let itemFormMode = $state<'ADD' | 'EDIT'>('ADD');
  let editingItemId = $state<string | null>(null);
  let selectedCatalogItemId = $state('');
  let selectedClassId = $state('');
  let formBasePrice = $state<number>(0);
  let formCitoPrice = $state<number | null>(null);
  let formComponents = $state<ItemComponentDTO[]>([]);
  let isSavingItem = $state(false);
  let formError = $state<string | null>(null);

  // Balancing Calculation
  let sumComponentBase = $derived(
    formComponents.reduce((acc, c) => acc + (Number(c.base_amount) || 0), 0)
  );
  let balanceDiff = $derived(
    Math.round((Number(formBasePrice) - sumComponentBase) * 100) / 100
  );
  let isBalanced = $derived(balanceDiff === 0);

  onMount(() => {
    Promise.all([
      getTariffClasses({ limit: 100 }),
      getTariffComponents({ limit: 100 }),
      getItems({ item_type: 'TARIFF', limit: 100 })
    ]).then(([classRes, compRes, itemRes]) => {
      tariffClasses = classRes.data || [];
      tariffComponents = compRes.data || [];
      tariffCatalogItems = itemRes.data || [];
    }).catch((err) => {
      console.error('Failed to load tariff references:', err);
    });
  });

  async function loadPlanItems() {
    if (!plan) return;
    isLoading = true;
    try {
      const res = await getPricePlanItems(plan.id, {
        page: currentPage,
        limit: 20,
        search: searchQuery || undefined,
        tariff_class_id: filterClassId || undefined
      });
      items = res.data;
      totalCount = res.meta.total_items;
      totalPages = res.meta.total_pages;
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat item tarif', true);
    } finally {
      isLoading = false;
    }
  }

  $effect(() => {
    if (isOpen && plan) {
      loadPlanItems();
    }
  });

  function openAddItem() {
    itemFormMode = 'ADD';
    editingItemId = null;
    selectedCatalogItemId = '';
    selectedClassId = tariffClasses[0]?.id || '';
    formBasePrice = 0;
    formCitoPrice = null;
    // Default initial component row
    if (tariffComponents.length > 0) {
      formComponents = [{
        component_id: tariffComponents[0].id,
        base_amount: 0,
        cito_amount: null,
        coa_code: tariffComponents[0].default_coa_code || null
      }];
    } else {
      formComponents = [];
    }
    formError = null;
    showItemForm = true;
  }

  function openEditItem(item: PricePlanItemRecord) {
    itemFormMode = 'EDIT';
    editingItemId = item.id;
    selectedCatalogItemId = item.item_id;
    selectedClassId = item.tariff_class_id;
    formBasePrice = item.total_base_price;
    formCitoPrice = item.total_cito_price || null;
    formComponents = item.components.map((c) => ({
      component_id: c.component_id,
      base_amount: c.base_amount,
      cito_amount: c.cito_amount || null,
      coa_code: c.coa_code || null
    }));
    formError = null;
    showItemForm = true;
  }

  function addComponentRow() {
    const firstComp = tariffComponents[0];
    formComponents = [
      ...formComponents,
      {
        component_id: firstComp?.id || '',
        base_amount: 0,
        cito_amount: null,
        coa_code: firstComp?.default_coa_code || null
      }
    ];
  }

  function removeComponentRow(idx: number) {
    formComponents = formComponents.filter((_, i) => i !== idx);
  }

  function handleAutoBalance() {
    if (formComponents.length === 0) return;
    const lastIdx = formComponents.length - 1;
    const currentLastAmount = formComponents[lastIdx].base_amount || 0;
    const adjusted = currentLastAmount + balanceDiff;
    if (adjusted >= 0) {
      formComponents[lastIdx].base_amount = adjusted;
    }
  }

  async function handleSaveItem() {
    if (!plan) return;
    if (!selectedCatalogItemId) {
      formError = 'Pilih tindakan medis dari katalog item tarif';
      return;
    }
    if (!selectedClassId) {
      formError = 'Pilih kelas tarif layanan';
      return;
    }
    if (Number(formBasePrice) <= 0) {
      formError = 'Total tarif dasar harus lebih dari 0';
      return;
    }
    if (formComponents.length === 0) {
      formError = 'Minimal satu rincian komponen biaya wajib disertakan';
      return;
    }
    if (!isBalanced) {
      formError = `Total pecahan komponen (Rp ${sumComponentBase.toLocaleString('id-ID')}) belum seimbang dengan total tarif dasar (Rp ${Number(formBasePrice).toLocaleString('id-ID')}). Selisih: Rp ${balanceDiff.toLocaleString('id-ID')}`;
      return;
    }

    formError = null;
    isSavingItem = true;
    try {
      if (itemFormMode === 'ADD') {
        await addPricePlanItem(plan.id, {
          item_id: selectedCatalogItemId,
          tariff_class_id: selectedClassId,
          total_base_price: Number(formBasePrice),
          total_cito_price: formCitoPrice ? Number(formCitoPrice) : undefined,
          is_active: true,
          components: formComponents
        });
        showToast('Item tarif berhasil ditambahkan ke buku tarif');
      } else if (editingItemId) {
        await updatePricePlanItem(plan.id, editingItemId, {
          total_base_price: Number(formBasePrice),
          total_cito_price: formCitoPrice ? Number(formCitoPrice) : undefined,
          is_active: true,
          components: formComponents
        });
        showToast('Item tarif berhasil diperbarui');
      }

      showItemForm = false;
      loadPlanItems();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan item tarif';
    } finally {
      isSavingItem = false;
    }
  }

  async function handleDeleteItem(item: PricePlanItemRecord) {
    if (!plan) return;
    if (!confirm(`Hapus item tarif ${item.item_name} pada kelas ${item.tariff_class_name}?`)) return;
    try {
      await deletePricePlanItem(plan.id, item.id);
      showToast('Item tarif berhasil dihapus');
      loadPlanItems();
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus item tarif', true);
    }
  }
</script>

{#if isOpen && plan}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-3 sm:p-6 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-5xl h-[90vh] flex flex-col overflow-hidden">
      <!-- Modal Header -->
      <div class="p-4 sm:p-5 border-b border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between shrink-0">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-[#0b57d0]/10 text-[#0b57d0] flex items-center justify-center font-bold">
            <Coins class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono font-bold text-xs text-[#0b57d0] bg-[#e8f0fe] px-2 py-0.5 rounded">
                {plan.code}
              </span>
              <span class="text-xs px-2 py-0.5 rounded-full font-bold {plan.status === 'ACTIVE' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'}">
                {plan.status}
              </span>
              {#if plan.status !== 'DRAFT'}
                <span class="inline-flex items-center gap-1 text-[11px] text-[#747775] italic">
                  <Lock class="w-3 h-3 text-[#747775]" />
                  Terkunci (Immutable)
                </span>
              {/if}
            </div>
            <h3 class="text-sm font-bold text-[#1f1f1f] mt-0.5">{plan.name}</h3>
          </div>
        </div>
        <button
          onclick={onClose}
          type="button"
          class="text-[#747775] hover:text-[#1f1f1f] p-1.5 rounded-full hover:bg-[#e1e5ea] transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Controls & Filter Bar -->
      <div class="p-3 sm:p-4 border-b border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3 bg-white shrink-0 text-xs">
        <div class="flex items-center gap-2 flex-1 min-w-[240px]">
          <div class="relative flex-1">
            <Search class="w-3.5 h-3.5 text-[#747775] absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              type="text"
              bind:value={searchQuery}
              onkeydown={(e) => e.key === 'Enter' && loadPlanItems()}
              placeholder="Cari tindakan atau kelas (tekan Enter)..."
              class="w-full h-9 pl-9 pr-3 rounded-xl bg-[#f8fafd] border border-[#e1e5ea] text-xs focus:bg-white focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <select
            bind:value={filterClassId}
            onchange={loadPlanItems}
            class="h-9 px-3 rounded-xl border border-[#e1e5ea] bg-white text-xs focus:border-[#0b57d0] focus:outline-hidden"
          >
            <option value="">Semua Kelas Tarif</option>
            {#each tariffClasses as tc (tc.id)}
              <option value={tc.id}>[{tc.code}] {tc.name}</option>
            {/each}
          </select>
        </div>

        {#if plan.status === 'DRAFT'}
          <M3Button variant="filled" onclick={openAddItem}>
            <Plus class="w-4 h-4" />
            <span>Tambah Item Tarif</span>
          </M3Button>
        {/if}
      </div>

      <!-- Table Content Area -->
      <div class="flex-1 overflow-y-auto p-4 sm:p-5">
        {#if isLoading}
          <div class="py-16 text-center text-[#747775] flex flex-col items-center justify-center gap-2">
            <div class="w-6 h-6 border-2 border-[#0b57d0] border-t-transparent rounded-full animate-spin"></div>
            <span class="text-xs">Memuat item tarif...</span>
          </div>
        {:else if items.length === 0}
          <div class="py-16 text-center text-[#747775] flex flex-col items-center justify-center gap-2">
            <Layers class="w-10 h-10 text-[#b0b8c4]" />
            <span class="font-bold text-xs text-[#1f1f1f]">Belum ada item tarif dalam buku tarif ini</span>
            <p class="text-[11px] text-[#747775]">
              {plan.status === 'DRAFT' ? 'Klik tombol "+ Tambah Item Tarif" di atas untuk memasukkan tindakan medis dan pecahan komponen biayanya.' : 'Buku tarif ini tidak memiliki item.'}
            </p>
          </div>
        {:else}
          <div class="border border-[#e1e5ea] rounded-2xl overflow-hidden shadow-2xs">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="bg-[#f0f4f9] text-[#444746] font-semibold border-b border-[#e1e5ea]">
                  <th class="py-3 px-3 w-10 text-center">No</th>
                  <th class="py-3 px-3">Tindakan Medis</th>
                  <th class="py-3 px-3">Kelas Tarif</th>
                  <th class="py-3 px-3 text-right">Tarif Dasar</th>
                  <th class="py-3 px-3 text-right">Tarif CITO (+{plan.default_cito_percent}%)</th>
                  <th class="py-3 px-3">Pecahan Komponen Biaya</th>
                  {#if plan.status === 'DRAFT'}
                    <th class="py-3 px-3 w-20 text-right">Aksi</th>
                  {/if}
                </tr>
              </thead>
              <tbody class="divide-y divide-[#e1e5ea]">
                {#each items as item, idx (item.id)}
                  <tr class="hover:bg-[#f8fafd] transition-colors">
                    <td class="py-3 px-3 text-center font-mono text-[#747775]">
                      {idx + 1}
                    </td>
                    <td class="py-3 px-3">
                      <div class="font-semibold text-[#1f1f1f]">{item.item_name}</div>
                      <span class="font-mono text-[10px] text-[#0b57d0]">{item.item_code}</span>
                    </td>
                    <td class="py-3 px-3">
                      <span class="px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-bold border border-blue-200 text-[11px]">
                        {item.tariff_class_name || item.tariff_class_code}
                      </span>
                    </td>
                    <td class="py-3 px-3 text-right font-mono font-bold text-[#1f1f1f]">
                      Rp {item.total_base_price.toLocaleString('id-ID')}
                    </td>
                    <td class="py-3 px-3 text-right font-mono text-[#747775]">
                      {#if item.total_cito_price}
                        Rp {item.total_cito_price.toLocaleString('id-ID')}
                      {:else}
                        Rp {Math.round(item.total_base_price * (1 + plan.default_cito_percent / 100)).toLocaleString('id-ID')}
                      {/if}
                    </td>
                    <td class="py-3 px-3">
                      <div class="flex flex-wrap gap-1">
                        {#each item.components as c (c.id)}
                          <span class="px-1.5 py-0.5 rounded bg-[#f0f4f9] border border-[#e1e5ea] text-[10px] text-[#333]">
                            <strong>{c.component_name || c.component_code}:</strong> Rp {c.base_amount.toLocaleString('id-ID')}
                          </span>
                        {/each}
                      </div>
                    </td>
                    {#if plan.status === 'DRAFT'}
                      <td class="py-3 px-3 text-right">
                        <div class="flex items-center justify-end gap-1">
                          <button
                            type="button"
                            onclick={() => openEditItem(item)}
                            class="p-1 rounded text-[#0b57d0] hover:bg-[#e8f0fe] cursor-pointer"
                            title="Edit Item & Komponen"
                          >
                            <Edit3 class="w-3.5 h-3.5" />
                          </button>
                          <button
                            type="button"
                            onclick={() => handleDeleteItem(item)}
                            class="p-1 rounded text-[#c5221f] hover:bg-[#fce8e6] cursor-pointer"
                            title="Hapus Item"
                          >
                            <Trash2 class="w-3.5 h-3.5" />
                          </button>
                        </div>
                      </td>
                    {/if}
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="p-3 sm:p-4 border-t border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between text-xs shrink-0">
        <span class="text-[#747775]">
          Total <strong>{totalCount}</strong> kombinasi item tindakan per kelas
        </span>
        <button
          type="button"
          onclick={onClose}
          class="px-4 py-2 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9]"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- SUB-MODAL FORM: TAMBAH / EDIT ITEM TARIF -->
{#if showItemForm && plan}
  <div class="fixed inset-0 z-60 bg-black/50 backdrop-blur-xs flex items-center justify-center p-3 sm:p-4 animate-in fade-in duration-100">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-2xl max-h-[92vh] flex flex-col overflow-hidden">
      <!-- Sub-Modal Header -->
      <div class="p-4 border-b border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between shrink-0">
        <div>
          <h4 class="text-sm font-bold text-[#1f1f1f]">
            {itemFormMode === 'ADD' ? 'Tambah Item Tindakan Tarif' : 'Edit Item & Pecahan Komponen Biaya'}
          </h4>
          <p class="text-[11px] text-[#747775]">Wajib memenuhi invariant balancing: Total Komponen = Total Tarif Dasar</p>
        </div>
        <button onclick={() => showItemForm = false} type="button" class="text-[#747775] hover:text-[#1f1f1f]">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Sub-Modal Body -->
      <form onsubmit={(e) => { e.preventDefault(); handleSaveItem(); }} class="flex-1 overflow-y-auto p-5 flex flex-col gap-4 text-xs">
        {#if formError}
          <div class="p-3 rounded-xl bg-[#fce8e6] text-[#c5221f] text-xs font-semibold flex items-center gap-2">
            <AlertTriangle class="w-4 h-4 shrink-0" />
            <span>{formError}</span>
          </div>
        {/if}

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Pilih Tindakan Medis -->
          <div>
            <label for="f-item" class="block font-semibold text-[#444746] mb-1">Tindakan Medis (Katalog TARIFF) *</label>
            <select
              id="f-item"
              bind:value={selectedCatalogItemId}
              disabled={itemFormMode === 'EDIT'}
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden disabled:bg-slate-100"
            >
              <option value="">-- Pilih Tindakan Medis --</option>
              {#each tariffCatalogItems as catItem (catItem.id)}
                <option value={catItem.id}>[{catItem.code}] {catItem.name}</option>
              {/each}
            </select>
          </div>

          <!-- Pilih Kelas Layanan -->
          <div>
            <label for="f-class" class="block font-semibold text-[#444746] mb-1">Kelas Layanan Perawatan *</label>
            <select
              id="f-class"
              bind:value={selectedClassId}
              disabled={itemFormMode === 'EDIT'}
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden disabled:bg-slate-100"
            >
              {#each tariffClasses as tc (tc.id)}
                <option value={tc.id}>[{tc.code}] {tc.name}</option>
              {/each}
            </select>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- Total Tarif Dasar -->
          <div>
            <label for="f-base" class="block font-semibold text-[#444746] mb-1">Total Tarif Dasar (Rp) *</label>
            <input
              id="f-base"
              type="number"
              min="0"
              step="any"
              bind:value={formBasePrice}
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono font-bold text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <!-- Total Tarif CITO (Opsional) -->
          <div>
            <label for="f-cito" class="block font-semibold text-[#444746] mb-1">Total Tarif CITO (Rp, Opsional)</label>
            <input
              id="f-cito"
              type="number"
              min="0"
              step="any"
              bind:value={formCitoPrice}
              placeholder="Kosong = Otomatis +{plan.default_cito_percent}%"
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>
        </div>

        <!-- Pecahan Komponen Biaya -->
        <div class="border border-[#e1e5ea] rounded-2xl p-4 bg-[#f8fafd] flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <div>
              <span class="font-bold text-[#1f1f1f] text-xs block">Rincian Pecahan Komponen Biaya</span>
              <p class="text-[11px] text-[#747775]">Pemecahan ke pos honor dokter, sarana RS, BHP, dll.</p>
            </div>
            <button
              type="button"
              onclick={addComponentRow}
              class="px-2.5 py-1 rounded-lg bg-white border border-[#e1e5ea] hover:bg-[#e8f0fe] hover:border-[#0b57d0] text-[#0b57d0] font-semibold text-xs flex items-center gap-1 cursor-pointer"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Tambah Komponen</span>
            </button>
          </div>

          <!-- Component Rows -->
          <div class="flex flex-col gap-2">
            {#each formComponents as compRow, idx (idx)}
              <div class="flex items-center gap-2 bg-white p-2.5 rounded-xl border border-[#e1e5ea]">
                <div class="flex-1">
                  <select
                    bind:value={compRow.component_id}
                    required
                    class="w-full h-8 px-2 rounded-lg border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
                  >
                    {#each tariffComponents as tc (tc.id)}
                      <option value={tc.id}>[{tc.code}] {tc.name}</option>
                    {/each}
                  </select>
                </div>

                <div class="w-36">
                  <input
                    type="number"
                    min="0"
                    step="any"
                    bind:value={compRow.base_amount}
                    placeholder="Nominal (Rp)"
                    required
                    class="w-full h-8 px-2 rounded-lg border border-[#e1e5ea] text-xs font-mono font-bold"
                  />
                </div>

                {#if formComponents.length > 1}
                  <button
                    type="button"
                    onclick={() => removeComponentRow(idx)}
                    class="p-1 rounded text-[#c5221f] hover:bg-[#fce8e6]"
                    title="Hapus baris komponen"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                {/if}
              </div>
            {/each}
          </div>

          <!-- Component Balancing Summary Bar -->
          <div class="p-3 rounded-xl border flex items-center justify-between text-xs {isBalanced ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-rose-50 border-rose-200 text-rose-800'}">
            <div class="flex items-center gap-2">
              {#if isBalanced}
                <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
                <span class="font-bold">Komponen Seimbang (100%): Rp {sumComponentBase.toLocaleString('id-ID')}</span>
              {:else}
                <AlertTriangle class="w-4 h-4 text-rose-600 shrink-0" />
                <div>
                  <span class="font-bold">Belum Seimbang!</span>
                  <span class="text-[11px] block">
                    Total Komponen: Rp {sumComponentBase.toLocaleString('id-ID')} • Selisih: Rp {balanceDiff.toLocaleString('id-ID')}
                  </span>
                </div>
              {/if}
            </div>

            {#if !isBalanced && formComponents.length > 0}
              <button
                type="button"
                onclick={handleAutoBalance}
                class="px-2.5 py-1 rounded-lg bg-white border border-rose-300 text-rose-700 font-bold text-[11px] hover:bg-rose-100 cursor-pointer"
              >
                Seimbangkan Otomatis
              </button>
            {/if}
          </div>
        </div>

        <!-- Sub-Modal Actions -->
        <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea]">
          <button
            type="button"
            onclick={() => showItemForm = false}
            class="px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746]"
          >
            Batal
          </button>
          <M3Button variant="filled" type="submit" disabled={isSavingItem || !isBalanced}>
            <span>{isSavingItem ? 'Menyimpan...' : 'Simpan Item Tarif'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
