<script lang="ts">
  import { onMount } from 'svelte';
  import { Search, ChevronDown, Check, X, Loader2, BookOpen } from '@lucide/svelte';
  import { getPostableAccounts } from '../api/accounting';
  import type { AccountRecord, AccountType } from '../types/accounting/account';

  // Module-level shared cache for postable accounts
  let accountsCache: AccountRecord[] | null = null;
  let fetchPromise: Promise<AccountRecord[]> | null = null;

  async function loadPostableAccounts(): Promise<AccountRecord[]> {
    if (accountsCache) return accountsCache;
    if (!fetchPromise) {
      fetchPromise = getPostableAccounts()
        .then((data) => {
          accountsCache = data || [];
          return accountsCache;
        })
        .catch((err) => {
          console.error('Failed to load postable accounts:', err);
          return [];
        })
        .finally(() => {
          fetchPromise = null;
        });
    }
    return fetchPromise;
  }

  interface Props {
    value?: string;
    label?: string;
    id?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    preferredType?: AccountType | 'ALL';
    helperText?: string;
    onChange?: (code: string, account?: AccountRecord) => void;
  }

  let {
    value = $bindable(''),
    label = '',
    id = '',
    placeholder = 'Pilih atau cari kode/nama akun COA...',
    required = false,
    disabled = false,
    preferredType = 'ALL',
    helperText = '',
    onChange
  }: Props = $props();

  let accounts = $state<AccountRecord[]>([]);
  let isLoading = $state(false);
  let isOpen = $state(false);
  let searchQuery = $state('');
  let activeTypeFilter = $state<AccountType | 'ALL'>('ALL');
  let containerEl: HTMLElement | null = $state(null);
  let searchInputEl: HTMLInputElement | null = $state(null);

  const accountTypeConfig: Record<AccountType, { label: string; badgeBg: string; badgeText: string }> = {
    ASSET: {
      label: 'Aset',
      badgeBg: 'bg-blue-50 border-blue-200 text-blue-700',
      badgeText: 'text-blue-700'
    },
    LIABILITY: {
      label: 'Liabilitas',
      badgeBg: 'bg-amber-50 border-amber-200 text-amber-700',
      badgeText: 'text-amber-700'
    },
    EQUITY: {
      label: 'Ekuitas',
      badgeBg: 'bg-purple-50 border-purple-200 text-purple-700',
      badgeText: 'text-purple-700'
    },
    REVENUE: {
      label: 'Pendapatan',
      badgeBg: 'bg-emerald-50 border-emerald-200 text-emerald-700',
      badgeText: 'text-emerald-700'
    },
    EXPENSE: {
      label: 'Beban/HPP',
      badgeBg: 'bg-rose-50 border-rose-200 text-rose-700',
      badgeText: 'text-rose-700'
    }
  };

  onMount(() => {
    if (preferredType && preferredType !== 'ALL') {
      activeTypeFilter = preferredType;
    }

    isLoading = true;
    loadPostableAccounts().then((data) => {
      accounts = data;
      isLoading = false;
    });

    function handleDocumentClick(event: MouseEvent) {
      if (isOpen && containerEl && !containerEl.contains(event.target as Node)) {
        isOpen = false;
      }
    }

    document.addEventListener('click', handleDocumentClick);
    return () => {
      document.removeEventListener('click', handleDocumentClick);
    };
  });

  // Find currently selected account object
  let selectedAccount = $derived(
    value ? accounts.find((a) => a.code.toLowerCase() === value.trim().toLowerCase()) : null
  );

  // Filter accounts based on query and type filter
  let filteredAccounts = $derived.by(() => {
    const q = searchQuery.trim().toLowerCase();
    return accounts.filter((acc) => {
      // Type filter
      if (activeTypeFilter !== 'ALL' && acc.type !== activeTypeFilter) {
        return false;
      }
      // Query filter (matches code, name, or description)
      if (q) {
        const matchCode = acc.code.toLowerCase().includes(q);
        const matchName = acc.name.toLowerCase().includes(q);
        return matchCode || matchName;
      }
      return true;
    });
  });

  function toggleOpen() {
    if (disabled) return;
    isOpen = !isOpen;
    if (isOpen) {
      searchQuery = '';
      setTimeout(() => searchInputEl?.focus(), 50);
    }
  }

  function handleSelect(acc: AccountRecord) {
    value = acc.code;
    isOpen = false;
    searchQuery = '';
    onChange?.(acc.code, acc);
  }

  function handleClear(e: MouseEvent) {
    e.stopPropagation();
    value = '';
    searchQuery = '';
    onChange?.('', undefined);
  }

  function handleManualApply() {
    const trimmed = searchQuery.trim();
    if (trimmed) {
      value = trimmed;
      isOpen = false;
      onChange?.(trimmed, undefined);
    }
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      isOpen = false;
    } else if (e.key === 'Enter' && isOpen) {
      e.preventDefault();
      if (filteredAccounts.length > 0) {
        handleSelect(filteredAccounts[0]);
      } else if (searchQuery.trim()) {
        handleManualApply();
      }
    }
  }
</script>

<div class="relative w-full text-xs" bind:this={containerEl}>
  {#if label}
    <label for={id} class="block font-semibold text-[#444746] mb-1">
      {label}
      {#if required}<span class="text-[#c5221f]">*</span>{/if}
    </label>
  {/if}

  <!-- Trigger Box -->
  <div
    role="button"
    tabindex={disabled ? -1 : 0}
    onclick={toggleOpen}
    onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggleOpen(); } }}
    class="w-full min-h-[38px] px-3 py-1.5 rounded-xl bg-white border border-[#e1e5ea] flex items-center justify-between gap-2 cursor-pointer transition-colors hover:border-[#b0b8c4] focus:outline-hidden focus:border-[#0b57d0] {disabled ? 'opacity-60 cursor-not-allowed bg-slate-50' : ''} {isOpen ? 'border-[#0b57d0] ring-1 ring-[#0b57d0]/20' : ''}"
  >
    <div class="flex items-center gap-2 flex-1 min-w-0">
      {#if selectedAccount}
        <span class="font-mono font-bold text-xs text-[#0b57d0] shrink-0">
          {selectedAccount.code}
        </span>
        <span class="text-xs text-[#1f1f1f] truncate font-medium">
          {selectedAccount.name}
        </span>
        {@const cfg = accountTypeConfig[selectedAccount.type]}
        {#if cfg}
          <span class="px-1.5 py-0.2 rounded text-[10px] font-semibold border {cfg.badgeBg} shrink-0">
            {cfg.label}
          </span>
        {/if}
      {:else if value}
        <span class="font-mono font-bold text-xs text-[#1f1f1f]">
          {value}
        </span>
        <span class="text-[11px] text-[#747775] italic shrink-0">
          (Kode manual)
        </span>
      {:else}
        <span class="text-[#747775] text-xs">
          {placeholder}
        </span>
      {/if}
    </div>

    <div class="flex items-center gap-1 shrink-0 text-[#747775]">
      {#if isLoading}
        <Loader2 class="w-3.5 h-3.5 animate-spin text-[#0b57d0]" />
      {:else if value && !disabled}
        <button
          type="button"
          onclick={handleClear}
          title="Hapus pilihan"
          class="p-0.5 rounded-full hover:bg-[#e1e5ea] text-[#747775] hover:text-[#c5221f] transition-colors"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      {/if}
      <ChevronDown class="w-4 h-4 transition-transform duration-150 {isOpen ? 'rotate-180 text-[#0b57d0]' : ''}" />
    </div>
  </div>

  {#if helperText}
    <p class="text-[11px] text-[#747775] mt-1">{helperText}</p>
  {/if}

  <!-- Dropdown Popover -->
  {#if isOpen}
    <div
      class="absolute left-0 right-0 top-full mt-1.5 z-50 bg-white rounded-2xl border border-[#e1e5ea] shadow-xl p-2.5 flex flex-col gap-2 max-h-[340px] animate-in fade-in zoom-in-95 duration-100"
    >
      <!-- Search Input -->
      <div class="relative flex items-center">
        <Search class="w-3.5 h-3.5 text-[#747775] absolute left-2.5 pointer-events-none" />
        <input
          bind:this={searchInputEl}
          type="text"
          bind:value={searchQuery}
          onkeydown={handleKeyDown}
          placeholder="Cari nomor akun atau nama akun..."
          class="w-full h-8 pl-8 pr-7 rounded-lg bg-[#f8fafd] border border-[#e1e5ea] text-xs focus:bg-white focus:border-[#0b57d0] focus:outline-hidden"
        />
        {#if searchQuery}
          <button
            type="button"
            onclick={() => searchQuery = ''}
            class="absolute right-2 text-[#747775] hover:text-[#1f1f1f]"
          >
            <X class="w-3 h-3" />
          </button>
        {/if}
      </div>

      <!-- Quick Type Filter Chips -->
      <div class="flex items-center gap-1 overflow-x-auto pb-1 text-[11px] no-scrollbar">
        <button
          type="button"
          onclick={() => activeTypeFilter = 'ALL'}
          class="px-2 py-0.5 rounded-full font-semibold transition-colors shrink-0 {activeTypeFilter === 'ALL' ? 'bg-[#0b57d0] text-white' : 'bg-[#f0f4f9] text-[#444746] hover:bg-[#e1e5ea]'}"
        >
          Semua ({accounts.length})
        </button>
        {#each Object.entries(accountTypeConfig) as [typeKey, typeCfg]}
          {@const count = accounts.filter(a => a.type === typeKey).length}
          <button
            type="button"
            onclick={() => activeTypeFilter = typeKey as AccountType}
            class="px-2 py-0.5 rounded-full font-medium transition-colors shrink-0 {activeTypeFilter === typeKey ? 'bg-[#0b57d0] text-white' : 'bg-[#f0f4f9] text-[#444746] hover:bg-[#e1e5ea]'}"
          >
            {typeCfg.label} ({count})
          </button>
        {/each}
      </div>

      <!-- Accounts List -->
      <div class="overflow-y-auto max-h-[200px] flex flex-col gap-0.5 divide-y divide-[#f0f4f9]">
        {#if filteredAccounts.length === 0}
          <div class="p-3 text-center text-[#747775] flex flex-col items-center gap-1.5">
            <BookOpen class="w-5 h-5 text-[#b0b8c4]" />
            <span class="text-xs font-medium">Tidak ada akun COA yang cocok</span>
            {#if searchQuery.trim()}
              <button
                type="button"
                onclick={handleManualApply}
                class="mt-1 px-2.5 py-1 rounded-lg bg-[#e8f0fe] text-[#0b57d0] font-semibold text-xs hover:bg-[#d3e3fd]"
              >
                Gunakan "{searchQuery.trim()}" sebagai kode akun
              </button>
            {/if}
          </div>
        {:else}
          {#each filteredAccounts as acc (acc.id)}
            {@const isSelected = value.toLowerCase() === acc.code.toLowerCase()}
            {@const typeCfg = accountTypeConfig[acc.type]}
            <button
              type="button"
              onclick={() => handleSelect(acc)}
              class="w-full text-left p-2 rounded-xl flex items-center justify-between gap-2 transition-colors hover:bg-[#f0f4f9] cursor-pointer {isSelected ? 'bg-[#e8f0fe]/70 text-[#0b57d0]' : 'text-[#1f1f1f]'}"
            >
              <div class="flex flex-col min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-mono font-bold text-xs text-[#0b57d0]">{acc.code}</span>
                  {#if typeCfg}
                    <span class="px-1.5 py-0.2 rounded text-[10px] font-semibold border {typeCfg.badgeBg}">
                      {typeCfg.label}
                    </span>
                  {/if}
                  <span class="text-[10px] text-[#747775]">
                    {acc.position === 'DEBIT' ? 'D' : 'K'}
                  </span>
                </div>
                <span class="text-xs font-medium text-[#1f1f1f] truncate mt-0.5">{acc.name}</span>
              </div>

              {#if isSelected}
                <Check class="w-4 h-4 text-[#0b57d0] shrink-0" />
              {/if}
            </button>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>
