<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    ChevronDown,
    Check,
    X,
    Loader2,
    Coins,
    Tag,
    Layers,
    Filter,
    ChevronLeft,
    ChevronRight,
    Sparkles
  } from '@lucide/svelte';
  import M3Button from './m3/M3Button.svelte';
  import { getItems, getItemById, getItemCategories } from '$lib/api/catalog/item';
  import type { ItemSummaryRecord, ItemCategoryRecord } from '$lib/types/master/item';

  interface Props {
    value?: string; // item_id yang terpilih
    label?: string;
    id?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    helperText?: string;
    onSelect?: (item: ItemSummaryRecord | null) => void;
  }

  let {
    value = $bindable(''),
    label = '',
    id = '',
    placeholder = 'Cari dan pilih tindakan medis / layanan...',
    required = false,
    disabled = false,
    helperText = '',
    onSelect
  }: Props = $props();

  // Selected item state for display
  let selectedItem = $state<ItemSummaryRecord | null>(null);
  let isFetchingSelected = $state(false);

  // Modal State
  let isModalOpen = $state(false);
  let isItemsLoading = $state(false);
  let items = $state<ItemSummaryRecord[]>([]);
  let categories = $state<ItemCategoryRecord[]>([]);

  // Search & Pagination Query State
  let searchQuery = $state('');
  let selectedCategoryFilter = $state('');
  let currentPage = $state(1);
  let pageSize = $state(10);
  let totalItems = $state(0);
  let totalPages = $state(1);

  let searchDebounceTimer: any = null;

  // Load categories once
  onMount(() => {
    getItemCategories({ item_type: 'TARIFF', limit: 100 })
      .then((res) => {
        categories = res.data || [];
      })
      .catch(() => {});
  });

  // Watch for external value change to fetch item details if needed
  $effect(() => {
    if (value && (!selectedItem || selectedItem.id !== value)) {
      fetchItemDetails(value);
    } else if (!value) {
      selectedItem = null;
    }
  });

  async function fetchItemDetails(itemId: string) {
    isFetchingSelected = true;
    try {
      const res = await getItemById(itemId);
      if (res) {
        selectedItem = {
          id: res.id,
          code: res.code,
          name: res.name,
          item_type: res.item_type,
          category_id: res.category_id,
          category_name: res.category?.name,
          product_line_id: res.product_line_id,
          product_line_name: res.product_line?.name,
          base_unit_id: res.base_unit_id,
          base_unit_name: res.base_unit?.name,
          is_active: res.is_active,
          created_at: res.created_at,
          updated_at: res.updated_at
        };
      }
    } catch (err) {
      console.warn('Gagal memuat detail item:', err);
    } finally {
      isFetchingSelected = false;
    }
  }

  async function loadItems() {
    isItemsLoading = true;
    try {
      const res = await getItems({
        page: currentPage,
        limit: pageSize,
        item_type: 'TARIFF',
        search: searchQuery.trim() || undefined,
        category_id: selectedCategoryFilter || undefined,
        is_active: true
      });

      items = res.data || [];
      totalItems = res.meta?.total_items ?? items.length;
      totalPages = res.meta?.total_pages ?? Math.max(1, Math.ceil(totalItems / pageSize));
    } catch (err) {
      console.error('Gagal memuat daftar tarif:', err);
      items = [];
    } finally {
      isItemsLoading = false;
    }
  }

  function handleOpenModal() {
    if (disabled) return;
    isModalOpen = true;
    currentPage = 1;
    loadItems();
  }

  function handleCloseModal() {
    isModalOpen = false;
  }

  function handleSearchInput(e: Event) {
    const target = e.target as HTMLInputElement;
    searchQuery = target.value;
    currentPage = 1;

    clearTimeout(searchDebounceTimer);
    searchDebounceTimer = setTimeout(() => {
      loadItems();
    }, 300);
  }

  function handleCategoryFilterChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    selectedCategoryFilter = target.value;
    currentPage = 1;
    loadItems();
  }

  function handlePageChange(newPage: number) {
    if (newPage < 1 || newPage > totalPages) return;
    currentPage = newPage;
    loadItems();
  }

  function handleSelectItem(item: ItemSummaryRecord) {
    value = item.id;
    selectedItem = item;
    isModalOpen = false;
    if (onSelect) {
      onSelect(item);
    }
  }

  function handleClearSelection(e: MouseEvent) {
    e.stopPropagation();
    value = '';
    selectedItem = null;
    if (onSelect) {
      onSelect(null);
    }
  }
</script>

<div class="flex flex-col gap-1.5 w-full">
  {#if label}
    <label for={id || undefined} class="text-xs font-semibold text-[#1f1f1f] flex items-center justify-between">
      <span class="flex items-center gap-1.5">
        <Coins class="w-3.5 h-3.5 text-emerald-600" />
        <span>{label}</span>
        {#if required}
          <span class="text-red-500">*</span>
        {/if}
      </span>
      <span class="text-[10px] font-mono text-[#747775] font-normal">item_type: TARIFF</span>
    </label>
  {/if}

  <!-- Display Trigger Button / Input -->
  <div
    class="relative group w-full cursor-pointer"
    onclick={handleOpenModal}
    role="button"
    tabindex="0"
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        handleOpenModal();
      }
    }}
  >
    <div
      class="w-full min-h-10 px-3 py-1.5 rounded-xl border transition-all flex items-center justify-between gap-2 {disabled
        ? 'bg-[#f0f4f9] border-[#e1e5ea] text-[#747775] cursor-not-allowed'
        : 'bg-white border-[#c4c7c5] hover:border-[#0b57d0] focus-within:border-[#0b57d0] focus-within:ring-1 focus-within:ring-[#0b57d0]'}"
    >
      {#if isFetchingSelected}
        <div class="flex items-center gap-2 text-xs text-[#747775]">
          <Loader2 class="w-3.5 h-3.5 animate-spin text-[#0b57d0]" />
          <span>Memuat info tindakan...</span>
        </div>
      {:else if selectedItem}
        <!-- Selected Item Badge -->
        <div class="flex items-center gap-2 overflow-hidden flex-1">
          <span class="font-mono font-bold text-[11px] px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 border border-emerald-300 shrink-0">
            {selectedItem.code}
          </span>
          <div class="flex flex-col overflow-hidden">
            <span class="text-xs font-bold text-[#1f1f1f] truncate leading-tight">
              {selectedItem.name}
            </span>
            {#if selectedItem.category_name}
              <span class="text-[10px] text-[#5f6368] truncate leading-tight">
                {selectedItem.category_name}
              </span>
            {/if}
          </div>
        </div>

        {#if !disabled}
          <div class="flex items-center gap-1 shrink-0">
            <button
              type="button"
              onclick={handleClearSelection}
              class="p-1 rounded-full text-[#747775] hover:text-red-600 hover:bg-red-50 transition-colors"
              title="Hapus pilihan"
            >
              <X class="w-3.5 h-3.5" />
            </button>
            <div class="w-px h-4 bg-[#e1e5ea]"></div>
            <span class="text-[11px] font-semibold text-[#0b57d0] hover:underline px-1">
              Ganti
            </span>
          </div>
        {/if}
      {:else}
        <!-- Placeholder Empty -->
        <div class="flex items-center gap-2 text-xs text-[#747775]">
          <Search class="w-4 h-4 text-[#747775]" />
          <span>{placeholder}</span>
        </div>

        <button
          type="button"
          class="px-2.5 py-1 rounded-lg bg-blue-50 text-blue-700 hover:bg-blue-100 text-[11px] font-semibold flex items-center gap-1 transition-colors"
        >
          <span>Cari Tarif</span>
          <ChevronDown class="w-3 h-3" />
        </button>
      {/if}
    </div>
  </div>

  {#if helperText}
    <p class="text-[11px] text-[#747775] leading-normal">{helperText}</p>
  {/if}
</div>

<!-- ======================================================== -->
<!-- MODAL GENERAL LOOKUP TARIF MEDIS -->
<!-- ======================================================== -->
{#if isModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs animate-in fade-in duration-200">
    <div class="bg-white w-full max-w-3xl rounded-3xl shadow-2xl border border-[#e1e5ea] overflow-hidden flex flex-col max-h-[85vh] animate-in zoom-in-95 duration-200">
      <!-- Modal Header -->
      <div class="p-4 bg-[#f8fafd] border-b border-[#e1e5ea] flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-2xl bg-emerald-600/10 text-emerald-700 flex items-center justify-center font-bold">
            <Coins class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">
              Pencarian Katalog Tindakan Medis & Tarif RS
            </h3>
            <p class="text-[11px] text-[#747775]">
              Cari tindakan medis berdasarkan kode, nama, atau kategori dari master catalog (item_type = TARIFF)
            </p>
          </div>
        </div>

        <button
          type="button"
          onclick={handleCloseModal}
          class="p-1.5 rounded-full text-[#747775] hover:bg-gray-100 transition-colors"
          title="Tutup"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Search & Category Filters Bar -->
      <div class="p-4 border-b border-[#e1e5ea] bg-white flex flex-wrap items-center gap-3">
        <!-- Search Input -->
        <div class="relative flex-1 min-w-[240px]">
          <Search class="w-4 h-4 absolute left-3 top-2.5 text-[#747775]" />
          <input
            type="text"
            value={searchQuery}
            oninput={handleSearchInput}
            placeholder="Ketik kode atau nama tindakan (misal: Konsultasi, USG, Darah Lengkap)..."
            class="w-full h-9 pl-9 pr-8 rounded-xl border border-[#c4c7c5] text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0]"
          />
          {#if searchQuery}
            <button
              type="button"
              onclick={() => {
                searchQuery = '';
                currentPage = 1;
                loadItems();
              }}
              class="absolute right-2.5 top-2.5 text-[#747775] hover:text-[#1f1f1f]"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          {/if}
        </div>

        <!-- Category Dropdown Filter -->
        <div class="w-48 shrink-0">
          <select
            value={selectedCategoryFilter}
            onchange={handleCategoryFilterChange}
            class="w-full h-9 px-3 rounded-xl border border-[#c4c7c5] bg-white text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0]"
          >
            <option value="">Semua Kategori</option>
            {#each categories as cat (cat.id)}
              <option value={cat.id}>{cat.name}</option>
            {/each}
          </select>
        </div>
      </div>

      <!-- Items List Table -->
      <div class="flex-1 overflow-y-auto min-h-[300px]">
        {#if isItemsLoading}
          <div class="py-20 flex flex-col items-center justify-center gap-2 text-xs text-[#747775]">
            <Loader2 class="w-7 h-7 animate-spin text-[#0b57d0]" />
            <span>Mencari tindakan medis...</span>
          </div>
        {:else if items.length === 0}
          <div class="py-20 flex flex-col items-center justify-center gap-2 text-center text-[#747775] px-4">
            <Coins class="w-8 h-8 text-[#b0b8c4]" />
            <span class="text-xs font-semibold text-[#1f1f1f]">Tidak ada tindakan medis yang ditemukan</span>
            <p class="text-[11px] max-w-sm">
              Coba gunakan kata kunci pencarian lain atau pastikan tindakan medis sudah terdaftar di Master Katalog Item dengan tipe TARIFF.
            </p>
          </div>
        {:else}
          <table class="w-full text-left text-xs border-collapse">
            <thead class="sticky top-0 bg-[#f8fafd] border-b border-[#e1e5ea] z-10 text-[11px] font-semibold text-[#747775]">
              <tr>
                <th class="py-2.5 px-4 w-32">Kode Tindakan</th>
                <th class="py-2.5 px-4">Nama Layanan Medis</th>
                <th class="py-2.5 px-4 w-40">Kategori</th>
                <th class="py-2.5 px-4 w-24 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#e1e5ea]">
              {#each items as item (item.id)}
                {@const isSelected = value === item.id}
                <tr
                  class="hover:bg-blue-50/60 cursor-pointer transition-colors {isSelected ? 'bg-blue-50/90 font-medium' : ''}"
                  onclick={() => handleSelectItem(item)}
                >
                  <td class="py-2.5 px-4 font-mono font-bold text-xs text-[#0b57d0]">
                    {item.code}
                  </td>
                  <td class="py-2.5 px-4">
                    <div class="font-semibold text-[#1f1f1f]">{item.name}</div>
                    {#if item.product_line_name}
                      <div class="text-[10px] text-[#5f6368] mt-0.5">{item.product_line_name}</div>
                    {/if}
                  </td>
                  <td class="py-2.5 px-4 text-[#444746]">
                    <span class="text-[11px] px-2 py-0.5 rounded-md bg-gray-100 border border-gray-200">
                      {item.category_name || 'Tarif Medis'}
                    </span>
                  </td>
                  <td class="py-2.5 px-4 text-center">
                    {#if isSelected}
                      <span class="px-2.5 py-1 rounded-lg bg-emerald-100 text-emerald-800 text-[11px] font-bold flex items-center justify-center gap-1">
                        <Check class="w-3.5 h-3.5" />
                        <span>Terpilih</span>
                      </span>
                    {:else}
                      <button
                        type="button"
                        onclick={(e) => {
                          e.stopPropagation();
                          handleSelectItem(item);
                        }}
                        class="px-3 py-1 rounded-lg bg-blue-600 text-white hover:bg-blue-700 text-[11px] font-semibold transition-colors"
                      >
                        Pilih
                      </button>
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>

      <!-- Modal Footer & Pagination -->
      <div class="p-3 bg-[#f8fafd] border-t border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3 text-xs">
        <div class="text-[11px] text-[#5f6368]">
          Menampilkan <span class="font-semibold text-[#1f1f1f]">{items.length}</span> dari total <span class="font-semibold text-[#1f1f1f]">{totalItems}</span> tindakan
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            disabled={currentPage <= 1 || isItemsLoading}
            onclick={() => handlePageChange(currentPage - 1)}
            class="h-8 px-2.5 rounded-lg border border-[#c4c7c5] bg-white text-[#444746] hover:bg-[#f0f4f9] disabled:opacity-40 text-xs font-medium flex items-center gap-1 transition-colors cursor-pointer"
          >
            <ChevronLeft class="w-3.5 h-3.5" />
            <span>Sebelumnya</span>
          </button>

          <span class="text-xs text-[#1f1f1f] font-mono px-2">
            {currentPage} / {totalPages}
          </span>

          <button
            type="button"
            disabled={currentPage >= totalPages || isItemsLoading}
            onclick={() => handlePageChange(currentPage + 1)}
            class="h-8 px-2.5 rounded-lg border border-[#c4c7c5] bg-white text-[#444746] hover:bg-[#f0f4f9] disabled:opacity-40 text-xs font-medium flex items-center gap-1 transition-colors cursor-pointer"
          >
            <span>Selanjutnya</span>
            <ChevronRight class="w-3.5 h-3.5" />
          </button>

          <div class="w-px h-4 bg-[#c4c7c5] mx-1"></div>

          <button
            type="button"
            onclick={handleCloseModal}
            class="h-8 px-3.5 rounded-lg bg-gray-200 hover:bg-gray-300 text-xs font-semibold text-[#1f1f1f] transition-colors"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
