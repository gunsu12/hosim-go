<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { ChevronDown, Check, Search, X } from '@lucide/svelte';

  export interface SelectOption {
    value: string;
    label: string;
    description?: string;
    badge?: string;
    disabled?: boolean;
  }

  interface Props {
    value?: string;
    options: (SelectOption | string)[];
    label?: string;
    id?: string;
    name?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    searchable?: boolean;
    helperText?: string;
    size?: 'sm' | 'md';
    class?: string;
    onchange?: (value: string) => void;
  }

  let {
    value = $bindable(''),
    options = [],
    label = '',
    id = '',
    name = '',
    placeholder = 'Pilih salah satu...',
    required = false,
    disabled = false,
    searchable = false,
    helperText = '',
    size = 'md',
    class: className = '',
    onchange
  }: Props = $props();

  let isOpen = $state(false);
  let searchQuery = $state('');
  let containerEl: HTMLElement | null = $state(null);
  let dropdownEl: HTMLElement | null = $state(null);
  let searchInputEl: HTMLInputElement | null = $state(null);
  let dropdownStyle = $state('');

  function portal(node: HTMLElement) {
    document.body.appendChild(node);
    return {
      destroy() {
        if (node.parentNode) {
          node.parentNode.removeChild(node);
        }
      }
    };
  }

  function updatePosition() {
    if (!containerEl) return;
    const rect = containerEl.getBoundingClientRect();
    const spaceBelow = window.innerHeight - rect.bottom;
    const spaceAbove = rect.top;
    const minWidth = Math.max(rect.width, 240);

    let left = rect.left;
    if (left + minWidth > window.innerWidth - 12) {
      left = Math.max(12, window.innerWidth - minWidth - 12);
    }
    if (left < 12) left = 12;

    const shouldOpenUp = spaceBelow < 220 && spaceAbove > spaceBelow;

    if (shouldOpenUp) {
      const maxHeight = Math.min(260, Math.max(120, spaceAbove - 16));
      dropdownStyle = `position: fixed; left: ${left}px; bottom: ${window.innerHeight - rect.top + 6}px; width: ${minWidth}px; max-height: ${maxHeight}px; z-index: 99999;`;
    } else {
      const maxHeight = Math.min(260, Math.max(120, spaceBelow - 16));
      dropdownStyle = `position: fixed; left: ${left}px; top: ${rect.bottom + 6}px; width: ${minWidth}px; max-height: ${maxHeight}px; z-index: 99999;`;
    }
  }

  // Normalisasi options array agar mendukung format string[] maupun SelectOption[]
  let normalizedOptions = $derived.by(() => {
    return options.map((opt) => {
      if (typeof opt === 'string') {
        return { value: opt, label: opt };
      }
      return opt;
    });
  });

  // Filter options jika searchable aktif
  let filteredOptions = $derived.by(() => {
    const q = searchQuery.toLowerCase().trim();
    if (!q) return normalizedOptions;
    return normalizedOptions.filter(
      (opt) =>
        opt.label.toLowerCase().includes(q) ||
        opt.value.toLowerCase().includes(q) ||
        (opt.description || '').toLowerCase().includes(q) ||
        (opt.badge || '').toLowerCase().includes(q)
    );
  });

  // Objek opsi yang sedang dipilih
  let selectedOption = $derived(
    normalizedOptions.find((opt) => opt.value === value) || null
  );

  function toggleOpen() {
    if (disabled) return;
    isOpen = !isOpen;
    if (isOpen) {
      searchQuery = '';
      updatePosition();
      tick().then(() => {
        updatePosition();
        if (searchable) {
          searchInputEl?.focus();
        }
      });
    }
  }

  function handleSelect(opt: SelectOption) {
    if (opt.disabled) return;
    value = opt.value;
    isOpen = false;
    if (onchange) {
      onchange(opt.value);
    }
  }

  function handleScrollOrResize() {
    if (!isOpen || !containerEl) return;
    const rect = containerEl.getBoundingClientRect();
    if (rect.bottom < 0 || rect.top > window.innerHeight) {
      isOpen = false;
      return;
    }
    updatePosition();
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as Node;
    if (
      containerEl &&
      !containerEl.contains(target) &&
      dropdownEl &&
      !dropdownEl.contains(target)
    ) {
      isOpen = false;
    }
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      isOpen = false;
    }
  }

  onMount(() => {
    document.addEventListener('click', handleClickOutside);
    window.addEventListener('scroll', handleScrollOrResize, true);
    window.addEventListener('resize', handleScrollOrResize);
    return () => {
      document.removeEventListener('click', handleClickOutside);
      window.removeEventListener('scroll', handleScrollOrResize, true);
      window.removeEventListener('resize', handleScrollOrResize);
    };
  });
</script>

<svelte:window
  onkeydown={(e) => {
    if (isOpen && e.key === 'Escape') {
      isOpen = false;
    }
  }}
/>

<div class="relative w-full {className}" bind:this={containerEl}>
  {#if label}
    <label for={id || undefined} class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
      {label}
      {#if required}<span class="text-rose-500">*</span>{/if}
    </label>
  {/if}

  <!-- Hidden native input agar form validation HTML tetap bekerja jika form di-submit -->
  {#if name}
    <input type="hidden" {name} {value} {required} />
  {/if}

  <!-- Trigger Box (Material Design 3 Surface) -->
  <button
    type="button"
    id={id || undefined}
    onclick={toggleOpen}
    {disabled}
    class="w-full {size === 'sm' ? 'h-8 px-2.5 rounded-lg' : 'h-10 px-3.5 rounded-xl'} bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-left flex items-center justify-between gap-2 cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed transition-all {isOpen ? 'border-[#0b57d0] dark:border-[#a8c7fa] bg-white dark:bg-[#1e1f20] ring-2 ring-[#0b57d0]/20 dark:ring-[#a8c7fa]/20' : ''}"
  >
    <div class="flex items-center gap-2 min-w-0 flex-1">
      {#if selectedOption}
        <span class="text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
          {selectedOption.label}
        </span>
        {#if selectedOption.badge}
          <span class="px-2 py-0.5 rounded-md text-[10px] font-mono font-semibold bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa] shrink-0">
            {selectedOption.badge}
          </span>
        {/if}
      {:else}
        <span class="text-[#747775] dark:text-[#8e918f] font-normal truncate">
          {placeholder}
        </span>
      {/if}
    </div>

    <ChevronDown
      class="{size === 'sm' ? 'w-3.5 h-3.5' : 'w-4 h-4'} text-[#747775] dark:text-[#8e918f] shrink-0 transition-transform duration-200 {isOpen ? 'rotate-180 text-[#0b57d0] dark:text-[#a8c7fa]' : ''}"
    />
  </button>

  {#if helperText}
    <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">{helperText}</p>
  {/if}

  <!-- M3 Floating Options Surface (Soft-Curved Rounded-2xl dengan Rounded-xl Options) via Portal -->
  {#if isOpen}
    <div
      use:portal
      bind:this={dropdownEl}
      style={dropdownStyle}
      class="bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] rounded-2xl shadow-2xl overflow-hidden flex flex-col animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Optional Search Bar -->
      {#if searchable || normalizedOptions.length > 7}
        <div class="p-2 border-b border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a]/50">
          <div class="relative">
            <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
            <input
              bind:this={searchInputEl}
              type="text"
              bind:value={searchQuery}
              placeholder="Cari pilihan..."
              class="w-full h-8 pl-7 pr-2.5 text-xs rounded-lg bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:outline-hidden"
            />
          </div>
        </div>
      {/if}

      <!-- Options List -->
      <div class="overflow-y-auto p-1.5 space-y-0.5">
        {#if filteredOptions.length === 0}
          <div class="p-4 text-center text-xs text-[#747775] dark:text-[#8e918f]">
            Tidak ada pilihan yang cocok.
          </div>
        {:else}
          {#each filteredOptions as opt (opt.value)}
            {@const isSelected = opt.value === value}
            <button
              type="button"
              onclick={() => handleSelect(opt)}
              disabled={opt.disabled}
              class="w-full text-left px-3 py-2 rounded-xl text-xs font-medium flex items-center justify-between gap-2.5 transition-all cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed {isSelected
                ? 'bg-[#d3e3fd]/60 dark:bg-[#0b57d0]/30 text-[#041e49] dark:text-[#d3e3fd] font-bold'
                : 'text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a]'}"
            >
              <div class="flex flex-col min-w-0">
                <div class="flex items-center gap-2">
                  <span class="truncate">{opt.label}</span>
                  {#if opt.badge}
                    <span class="px-1.5 py-0.2 rounded text-[10px] font-mono {isSelected ? 'bg-white/80 dark:bg-black/30 text-[#041e49] dark:text-[#d3e3fd]' : 'bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]'}">
                      {opt.badge}
                    </span>
                  {/if}
                </div>
                {#if opt.description}
                  <span class="text-[10px] text-[#747775] dark:text-[#8e918f] font-normal truncate">
                    {opt.description}
                  </span>
                {/if}
              </div>

              {#if isSelected}
                <Check class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa] shrink-0" />
              {/if}
            </button>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>
