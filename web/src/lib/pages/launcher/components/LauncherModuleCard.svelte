<script lang="ts">
  import { ChevronRight } from '@lucide/svelte';
  import type { ModuleDefinition } from '$lib/config/launcherModules';

  interface Props {
    mod: ModuleDefinition;
    onSelect: (moduleId: string, defaultNavId: string) => void;
  }

  let { mod, onSelect }: Props = $props();
</script>

<div
  role="button"
  tabindex="0"
  onclick={() => onSelect(mod.id, mod.defaultNav)}
  onkeydown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') onSelect(mod.id, mod.defaultNav);
  }}
  class="group relative text-left bg-white dark:bg-surface-container rounded-xl p-3.5 border border-[#e1e5ea] dark:border-outline-variant hover:shadow-md transition-all duration-150 cursor-pointer flex flex-col justify-between gap-2.5 {mod.color.hoverBg}"
>
  <!-- Top: Icon & Status -->
  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between gap-2">
      <div
        class="w-8 h-8 rounded-lg {mod.color.bg} {mod.color.text} flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform duration-150"
      >
        <mod.icon class="w-4 h-4" />
      </div>

      <div class="flex items-center gap-1.5 shrink-0">
        {#if mod.quickStats}
          <span class="text-[9px] font-mono text-[#747775] dark:text-on-surface-variant hidden sm:inline">
            {mod.quickStats}
          </span>
        {/if}
        <span
          class="text-[9px] font-semibold px-1.5 py-0.5 rounded-full {mod.status === 'ready'
            ? 'bg-[#e6f4ea] text-[#137333] border border-[#ceead6] dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800/60'
            : 'bg-[#f1f3f4] text-[#5f6368] border border-[#dadce0] dark:bg-surface-container-high dark:text-on-surface-variant dark:border-outline-variant'}"
        >
          {mod.statusText}
        </span>
      </div>
    </div>

    <!-- Title & Description -->
    <div>
      <h3
        class="text-xs sm:text-sm font-bold text-[#1f1f1f] dark:text-on-surface group-hover:text-[#0b57d0] dark:group-hover:text-primary transition-colors flex items-center justify-between gap-1"
      >
        <span class="truncate">{mod.title}</span>
        <ChevronRight
          class="w-3.5 h-3.5 text-[#747775] dark:text-on-surface-variant opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all shrink-0"
        />
      </h3>
      <p class="text-[11px] text-[#444746] dark:text-on-surface-variant leading-snug line-clamp-2 mt-0.5">
        {mod.description}
      </p>
    </div>
  </div>

  <!-- Features Pills -->
  <div class="pt-2 border-t border-[#f0f4f9] dark:border-outline-variant/40 flex flex-wrap items-center gap-1">
    {#each mod.features as feat}
      <span
        class="text-[9px] px-1.5 py-0.5 rounded bg-[#f8fafd] dark:bg-surface-container-high border border-[#e1e5ea] dark:border-outline-variant text-[#5f6368] dark:text-on-surface-variant leading-tight"
      >
        {feat}
      </span>
    {/each}
  </div>
</div>
