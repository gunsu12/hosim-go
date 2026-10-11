<script lang="ts">
  import { Sun, Moon, Monitor, Check } from '@lucide/svelte';
  import { theme, type ThemeMode } from '../stores/theme.svelte';

  interface Props {
    compact?: boolean;
    showLabel?: boolean;
  }

  let { compact = false, showLabel = false }: Props = $props();

  let isOpen = $state(false);
  let containerRef = $state<HTMLDivElement | null>(null);

  function toggleMenu(e: MouseEvent) {
    e.stopPropagation();
    isOpen = !isOpen;
  }

  function handleSelect(mode: ThemeMode) {
    theme.setMode(mode);
    isOpen = false;
  }

  function handleDocumentClick(e: MouseEvent) {
    if (isOpen && containerRef && !containerRef.contains(e.target as Node)) {
      isOpen = false;
    }
  }

  const options: { mode: ThemeMode; label: string; icon: any; desc: string }[] = [
    {
      mode: 'light',
      label: 'Mode Terang',
      icon: Sun,
      desc: 'Tampilan bersih & kontras standar'
    },
    {
      mode: 'dark',
      label: 'Mode Gelap',
      icon: Moon,
      desc: 'Nyaman di mata saat malam / minim cahaya'
    },
    {
      mode: 'system',
      label: 'Ikuti Sistem (Otomatis)',
      icon: Monitor,
      desc: 'Sinkron otomatis dengan setelan OS'
    }
  ];
</script>

<svelte:window onclick={handleDocumentClick} />

<div bind:this={containerRef} class="relative inline-block">
  <button
    type="button"
    onclick={toggleMenu}
    class="flex items-center gap-2 p-2 rounded-full transition-colors cursor-pointer text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-white hover:bg-[#e9eef6] dark:hover:bg-[#28292a] active:bg-[#e1e5ea] dark:active:bg-[#353738]"
    title={`Tema: ${theme.mode === 'system' ? 'Sistem (' + (theme.isDark ? 'Gelap' : 'Terang') + ')' : theme.mode === 'dark' ? 'Mode Gelap' : 'Mode Terang'}`}
    aria-label="Pengaturan Tema Tampilan"
    aria-haspopup="true"
    aria-expanded={isOpen}
  >
    <div class="relative w-5 h-5 flex items-center justify-center">
      {#if theme.mode === 'dark'}
        <Moon class="w-5 h-5 text-indigo-400 animate-in zoom-in-75 duration-150" />
      {:else if theme.mode === 'light'}
        <Sun class="w-5 h-5 text-amber-500 animate-in zoom-in-75 duration-150" />
      {:else}
        <!-- Mode Sistem: tampilkan Monitor dengan aksen sesuai tema aktif -->
        <Monitor class="w-5 h-5 text-[#0b57d0] dark:text-[#8ab4f8] animate-in zoom-in-75 duration-150" />
      {/if}
    </div>

    {#if showLabel}
      <span class="text-xs font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
        {theme.mode === 'system' ? 'Sistem' : theme.mode === 'dark' ? 'Gelap' : 'Terang'}
      </span>
    {/if}
  </button>

  {#if isOpen}
    <div
      class="absolute right-0 top-full mt-2 w-72 rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-2xl p-2 z-50 animate-in fade-in slide-in-from-top-1 duration-150"
      role="menu"
    >
      <div class="px-3 py-2 text-[11px] font-bold text-[#747775] dark:text-[#9aa0a6] uppercase tracking-wider border-b border-[#e1e5ea] dark:border-[#2e3032] mb-1.5 flex items-center justify-between">
        <span>Pilihan Tema Tampilan</span>
      </div>

      <div class="flex flex-col gap-1">
        {#each options as opt}
          {@const isActive = theme.mode === opt.mode}
          <button
            type="button"
            onclick={() => handleSelect(opt.mode)}
            class="group flex items-center justify-between w-full px-3 py-2.5 rounded-xl text-left transition-all duration-150 cursor-pointer {isActive
              ? 'bg-[#d3e3fd]/90 text-[#041e49] font-semibold dark:bg-[#0842a0] dark:text-[#d3e3fd]'
              : 'text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#f0f4f9] dark:hover:bg-[#28292a]'}"
            role="menuitem"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div
                class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-colors {isActive
                  ? 'bg-white/80 text-[#041e49] dark:bg-white/10 dark:text-[#d3e3fd]'
                  : 'bg-[#f0f4f9] text-[#444746] group-hover:bg-[#e1e5ea] dark:bg-[#28292a] dark:text-[#c4c7c5] dark:group-hover:bg-[#353738] dark:group-hover:text-white'}"
              >
                <opt.icon class="w-4.5 h-4.5" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs truncate font-medium {isActive ? 'text-[#041e49] dark:text-[#d3e3fd] font-semibold' : 'text-[#1f1f1f] dark:text-[#e3e3e3] group-hover:text-black dark:group-hover:text-white'}">
                  {opt.label}
                </span>
                <span class="text-[10px] truncate leading-tight {isActive ? 'text-[#041e49]/80 dark:text-[#d3e3fd]/80' : 'text-[#747775] dark:text-[#9aa0a6] group-hover:text-[#444746] dark:group-hover:text-[#c4c7c5]'}">
                  {opt.desc}
                </span>
              </div>
            </div>

            {#if isActive}
              <Check class="w-4 h-4 text-[#041e49] dark:text-[#d3e3fd] shrink-0 stroke-[2.5]" />
            {/if}
          </button>
        {/each}
      </div>

      <div class="mt-2 pt-2 border-t border-[#e1e5ea] dark:border-[#2e3032] px-3 py-1.5 flex items-center justify-between text-[11px] text-[#747775] dark:text-[#9aa0a6]">
        <span>Aktif saat ini:</span>
        <span class="inline-flex items-center gap-1.5 font-medium px-2.5 py-0.5 rounded-full {theme.isDark ? 'bg-indigo-950/70 text-indigo-300 border border-indigo-800/60' : 'bg-amber-100 text-amber-800 border border-amber-200'}">
          {#if theme.isDark}
            <Moon class="w-3 h-3 text-indigo-400" />
            <span>Mode Gelap</span>
          {:else}
            <Sun class="w-3 h-3 text-amber-600" />
            <span>Mode Terang</span>
          {/if}
        </span>
      </div>
    </div>
  {/if}
</div>
