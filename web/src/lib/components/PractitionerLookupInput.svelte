<script lang="ts">
  import { onMount, tick } from 'svelte';
  import {
    Stethoscope,
    Search,
    ChevronDown,
    Check,
    X,
    Loader2,
    UserCheck,
    AlertCircle
  } from '@lucide/svelte';
  import { getPractitioners, getPractitionerById } from '../api/practitioners';
  import type { PractitionerRecord } from '../types/master/practitioner';

  interface Props {
    value?: string;
    label?: string;
    id?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    allowClear?: boolean;
    helperText?: string;
    professionFilter?: string; // e.g. 'DOKTER' or profession code
    onChange?: (id: string, practitioner?: PractitionerRecord | null) => void;
  }

  let {
    value = $bindable(''),
    label = '',
    id = '',
    placeholder = 'Pilih dokter atau tenaga medis...',
    required = false,
    disabled = false,
    allowClear = true,
    helperText = '',
    professionFilter = '',
    onChange
  }: Props = $props();

  const PAGE_SIZE = 20;

  let isOpen = $state(false);
  let searchQuery = $state('');
  let searchInputEl: HTMLInputElement | null = $state(null);
  let containerEl: HTMLElement | null = $state(null);

  let items = $state<PractitionerRecord[]>([]);
  let totalCount = $state<number>(0);
  let currentPage = $state<number>(1);
  let isLoading = $state(false);
  let isLoadingMore = $state(false);
  let selectedPractitioner = $state<PractitionerRecord | null>(null);

  let searchTimer: ReturnType<typeof setTimeout> | null = null;
  let hasMore = $derived(items.length < totalCount);

  // Pre-load selected item saat value berubah dari luar
  $effect(() => {
    if (value && (!selectedPractitioner || selectedPractitioner.id !== value)) {
      resolveInitialPractitioner(value);
    } else if (!value) {
      selectedPractitioner = null;
    }
  });

  async function resolveInitialPractitioner(targetId: string) {
    // Cek di items yang sudah ter-load
    const found = items.find((p) => p.id === targetId);
    if (found) {
      selectedPractitioner = found;
      return;
    }
    // Jika belum ada, fetch spesifik by ID
    try {
      const fetched = await getPractitionerById(targetId);
      selectedPractitioner = fetched;
    } catch {
      selectedPractitioner = null;
    }
  }

  async function fetchItems(page: number, query: string, append: boolean = false) {
    if (append) {
      isLoadingMore = true;
    } else {
      isLoading = true;
    }

    try {
      const res = await getPractitioners({
        search: query.trim() || undefined,
        page,
        limit: PAGE_SIZE
      });

      let loaded = res.data || [];
      if (professionFilter) {
        loaded = loaded.filter((p) => {
          const prof = typeof p.profession === 'object' && p.profession ? p.profession.code : String(p.profession || '');
          return prof.toLowerCase().includes(professionFilter.toLowerCase());
        });
      }

      totalCount = res.count || loaded.length;
      currentPage = page;

      if (append) {
        const existingIds = new Set(items.map((i) => i.id));
        const unique = loaded.filter((i) => !existingIds.has(i.id));
        items = [...items, ...unique];
      } else {
        items = loaded;
      }
    } catch (err) {
      console.error('Failed to load practitioners:', err);
    } finally {
      isLoading = false;
      isLoadingMore = false;
    }
  }

  function handleOpen() {
    if (disabled) return;
    isOpen = !isOpen;
    if (isOpen) {
      searchQuery = '';
      if (items.length === 0) {
        fetchItems(1, '', false);
      }
      tick().then(() => {
        searchInputEl?.focus();
      });
    }
  }

  function handleSearchInput(e: Event) {
    const val = (e.target as HTMLInputElement).value;
    searchQuery = val;
    if (searchTimer) clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      fetchItems(1, val, false);
    }, 300);
  }

  function handleSelect(p: PractitionerRecord) {
    value = p.id;
    selectedPractitioner = p;
    isOpen = false;
    if (onChange) {
      onChange(p.id, p);
    }
  }

  function handleClear(e: MouseEvent) {
    e.stopPropagation();
    value = '';
    selectedPractitioner = null;
    if (onChange) {
      onChange('', null);
    }
  }

  function handleLoadMore() {
    if (isLoadingMore || !hasMore) return;
    fetchItems(currentPage + 1, searchQuery, true);
  }

  // Click outside to close
  function handleClickOutside(e: MouseEvent) {
    if (containerEl && !containerEl.contains(e.target as Node)) {
      isOpen = false;
    }
  }

  onMount(() => {
    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  });
</script>

<div class="relative w-full" bind:this={containerEl}>
  {#if label}
    <label for={id || undefined} class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
      {label}
      {#if required}<span class="text-rose-500">*</span>{/if}
    </label>
  {/if}

  <!-- Trigger Box -->
  <button
    type="button"
    onclick={handleOpen}
    {disabled}
    class="w-full min-h-10 px-3 py-2 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-left flex items-center justify-between gap-2 cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed transition-all {isOpen ? 'border-[#0b57d0] dark:border-[#a8c7fa] bg-white dark:bg-[#1e1f20] ring-2 ring-[#0b57d0]/20 dark:ring-[#a8c7fa]/20' : ''}"
  >
    {#if selectedPractitioner}
      {@const prof = typeof selectedPractitioner.profession === 'object' && selectedPractitioner.profession ? selectedPractitioner.profession.name : (selectedPractitioner.profession || 'Nakes')}
      {@const ihs = selectedPractitioner.ihs_practitioner_id || selectedPractitioner.ihs_number || selectedPractitioner.ihsId}
      <div class="flex items-center gap-2.5 min-w-0 flex-1">
        <div class="w-6 h-6 rounded-lg bg-teal-100 dark:bg-teal-900/60 text-teal-700 dark:text-teal-300 flex items-center justify-center shrink-0">
          <Stethoscope class="w-3.5 h-3.5" />
        </div>
        <div class="flex flex-col min-w-0">
          <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
            {selectedPractitioner.name}
          </span>
          <span class="text-[10px] text-[#747775] dark:text-[#8e918f] truncate font-normal">
            {prof} {selectedPractitioner.sip ? `• SIP: ${selectedPractitioner.sip}` : ''} {ihs ? `• IHS: ${ihs}` : ''}
          </span>
        </div>
      </div>
    {:else}
      <div class="flex items-center gap-2 text-[#747775] dark:text-[#8e918f]">
        <Stethoscope class="w-4 h-4 opacity-50 shrink-0" />
        <span class="font-normal truncate">{placeholder}</span>
      </div>
    {/if}

    <div class="flex items-center gap-1 shrink-0 text-[#747775] dark:text-[#8e918f]">
      {#if allowClear && selectedPractitioner && !disabled}
        <span
          role="button"
          tabindex="0"
          onclick={handleClear}
          onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); handleClear(e as any); } }}
          class="p-1 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 rounded-full transition-colors"
          title="Hapus pilihan"
        >
          <X class="w-3.5 h-3.5" />
        </span>
      {/if}
      <ChevronDown class="w-4 h-4 transition-transform {isOpen ? 'rotate-180' : ''}" />
    </div>
  </button>

  {#if helperText}
    <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">{helperText}</p>
  {/if}

  <!-- Popover Dropdown Panel -->
  {#if isOpen}
    <div class="absolute left-0 right-0 z-50 mt-1.5 bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] rounded-2xl shadow-xl overflow-hidden flex flex-col max-h-80 animate-in fade-in zoom-in-95 duration-150">
      <!-- Search Input Header -->
      <div class="p-2.5 border-b border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a]/50">
        <div class="relative">
          <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
          <input
            bind:this={searchInputEl}
            type="text"
            value={searchQuery}
            oninput={handleSearchInput}
            placeholder="Cari nama dokter, profesi, SIP, atau STR..."
            class="w-full h-8 pl-8 pr-3 rounded-lg bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:outline-hidden"
          />
        </div>
      </div>

      <!-- Items List -->
      <div class="overflow-y-auto flex-1 divide-y divide-[#e1e5ea] dark:divide-[#353738]">
        {#if allowClear}
          <button
            type="button"
            onclick={() => {
              value = '';
              selectedPractitioner = null;
              isOpen = false;
              if (onChange) onChange('', null);
            }}
            class="w-full p-2.5 text-left text-xs text-[#747775] dark:text-[#8e918f] hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/60 flex items-center gap-2 cursor-pointer"
          >
            <span class="w-4 h-4 rounded-full border border-dashed border-[#c4c7c5] dark:border-[#5e5e5e] flex items-center justify-center shrink-0">
              <X class="w-2.5 h-2.5" />
            </span>
            <span>-- Bukan Tenaga Medis (Kosongkan) --</span>
          </button>
        {/if}

        {#if isLoading}
          <div class="p-6 text-center text-[#747775] dark:text-[#8e918f] flex items-center justify-center gap-2 text-xs">
            <Loader2 class="w-4 h-4 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
            <span>Mencari data tenaga medis...</span>
          </div>
        {:else if items.length === 0}
          <div class="p-6 text-center text-[#747775] dark:text-[#8e918f] text-xs">
            {searchQuery
              ? `Tidak ditemukan tenaga medis dengan kata kunci "${searchQuery}".`
              : 'Belum ada data tenaga medis.'}
          </div>
        {:else}
          {#each items as p (p.id)}
            {@const isSelected = value === p.id}
            {@const prof = typeof p.profession === 'object' && p.profession ? p.profession.name : (p.profession || 'Nakes')}
            {@const ihs = p.ihs_practitioner_id || p.ihs_number || p.ihsId}
            <button
              type="button"
              onclick={() => handleSelect(p)}
              class="w-full p-2.5 text-left text-xs flex items-center justify-between gap-3 transition-colors cursor-pointer {isSelected ? 'bg-teal-50 dark:bg-teal-950/50 text-teal-900 dark:text-teal-100' : 'hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/60 text-[#1f1f1f] dark:text-[#e3e3e3]'}"
            >
              <div class="flex items-center gap-2.5 min-w-0">
                <div class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 {isSelected ? 'bg-teal-600 text-white' : 'bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]'}">
                  <Stethoscope class="w-3.5 h-3.5" />
                </div>
                <div class="flex flex-col min-w-0">
                  <div class="font-bold truncate">{p.name}</div>
                  <div class="text-[10px] text-[#747775] dark:text-[#8e918f] truncate">
                    {prof} {p.sip ? `• SIP: ${p.sip}` : ''} {ihs ? `• IHS: ${ihs}` : ''}
                  </div>
                </div>
              </div>
              {#if isSelected}
                <Check class="w-4 h-4 text-teal-600 dark:text-teal-400 shrink-0" />
              {/if}
            </button>
          {/each}

          {#if hasMore}
            <div class="p-2 text-center bg-[#f8fafd] dark:bg-[#28292a]/50">
              <button
                type="button"
                onclick={handleLoadMore}
                disabled={isLoadingMore}
                class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs font-semibold text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] transition-colors cursor-pointer disabled:opacity-50"
              >
                {#if isLoadingMore}
                  <Loader2 class="w-3.5 h-3.5 animate-spin" />
                  <span>Memuat...</span>
                {:else}
                  <ChevronDown class="w-3.5 h-3.5" />
                  <span>Muat {PAGE_SIZE} lainnya (Sisa {totalCount - items.length})</span>
                {/if}
              </button>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  {/if}
</div>
