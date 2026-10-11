<script lang="ts">
  import {
    ChevronDown,
    ChevronRight,
    Plus,
    ArrowLeft,
    LayoutGrid,
  } from "@lucide/svelte";
  import { auth } from "../stores/auth.svelte";
  import {
    navigationStructure,
    type MenuItem,
    type SubMenuItem,
  } from "../config/navigationStructure";

  export type { MenuItem, SubMenuItem };

  interface Props {
    activeNav?: string;
    activeModule?: string | null;
    isExpanded?: boolean;
    onNewEncounter?: () => void;
    onNavigateHome?: () => void;
  }

  let {
    activeNav = $bindable("physical"),
    activeModule = null,
    isExpanded = $bindable(false),
    onNewEncounter = () => {},
    onNavigateHome = () => {},
  }: Props = $props();

  let openParentIds = $state<string[]>([]);
  let flyoutMenuId = $state<string | null>(null);

  // Filter navigasi berdasarkan konteks modul aktif (Focused Contextual Sidebar)
  const visibleNavigation = $derived(
    navigationStructure
      .filter((item) => {
        if (!activeModule) return true;
        if (activeModule === "master") return item.id === "master";
        if (activeModule === "clinical")
          return item.id === "clinical" || item.id === "activity";
        if (activeModule === "auth") return item.id === "auth";
        if (activeModule === "accounting") return item.id === "accounting";
        if (activeModule === "finance")
          return item.id === "tariff" || item.id === "customer";
        return item.id === activeModule;
      })
      .map((item) => ({
        ...item,
        children: item.children.filter(
          (sub) => !sub.permission || auth.hasPermission(sub.permission),
        ),
      }))
      .filter((item) => item.children.length > 0),
  );

  function toggleAccordion(parentId: string): void {
    if (openParentIds.includes(parentId)) {
      openParentIds = openParentIds.filter((id) => id !== parentId);
    } else {
      // Ketika klik menu yang lain, collapse accordion sebelumnya
      openParentIds = [parentId];
    }
  }

  function handleRailItemClick(itemId: string): void {
    flyoutMenuId = flyoutMenuId === itemId ? null : itemId;
  }

  function handleDocumentClick(e: MouseEvent): void {
    if (!flyoutMenuId) return;
    const target = e.target as HTMLElement | null;
    if (!target) return;

    // Cek apakah klik berada di dalam container item rail yang sedang aktif flyout-nya
    const activeItemContainer = document.querySelector(`[data-rail-item="${flyoutMenuId}"]`);
    if (activeItemContainer && activeItemContainer.contains(target)) {
      return;
    }

    flyoutMenuId = null;
  }

  function handleKeydown(e: KeyboardEvent): void {
    if (e.key === "Escape" && flyoutMenuId) {
      flyoutMenuId = null;
    }
  }

  function handleSubmenuClick(subId: string): void {
    activeNav = subId;
    flyoutMenuId = null;
  }

  function isParentActive(parent: MenuItem): boolean {
    return parent.children?.some((c) => c.id === activeNav);
  }

  // Otomatis buka accordion parent jika activeNav berada di dalamnya
  $effect(() => {
    const parent = visibleNavigation.find((p) =>
      p.children?.some((c) => c.id === activeNav),
    );
    if (parent && !openParentIds.includes(parent.id)) {
      openParentIds = [parent.id];
    }
  });

  // Otomatis tutup flyout saat sidebar diperluas
  $effect(() => {
    if (isExpanded) {
      flyoutMenuId = null;
    }
  });
</script>

<svelte:window onclick={handleDocumentClick} onkeydown={handleKeydown} />

<aside
  class="py-3 flex flex-col gap-2 select-none shrink-0 transition-all duration-200 relative {isExpanded
    ? 'w-64 px-3'
    : 'w-[72px] px-2 items-center'}"
>
  <!-- Tombol Navigasi Kembali ke Desktop Launcher -->
  {#if isExpanded}
    <button
      type="button"
      onclick={onNavigateHome}
      class="flex items-center gap-2.5 px-3.5 py-2.5 rounded-2xl bg-[#e8f0fe] text-[#0b57d0] hover:bg-[#d3e3fd] dark:bg-[#0842a0]/40 dark:text-[#a8c7fa] dark:hover:bg-[#0842a0]/70 text-xs font-semibold transition-all cursor-pointer w-full border border-[#d3e3fd]/60 dark:border-[#0842a0]/60 mb-1"
      title="Kembali ke Beranda Modul (Desktop Launcher)"
    >
      <ArrowLeft class="w-4 h-4" />
      <span>Beranda Modul</span>
    </button>
  {:else}
    <button
      type="button"
      onclick={onNavigateHome}
      class="flex items-center justify-center w-11 h-9 rounded-full bg-[#e8f0fe] text-[#0b57d0] hover:bg-[#d3e3fd] dark:bg-[#0842a0]/40 dark:text-[#a8c7fa] dark:hover:bg-[#0842a0]/70 border border-[#d3e3fd]/60 dark:border-[#0842a0]/60 transition-all cursor-pointer mb-2 shadow-xs"
      title="Kembali ke Beranda Modul (Desktop Launcher)"
    >
      <ArrowLeft class="w-4.5 h-4.5" />
    </button>
  {/if}

  <!-- Tombol Aksi Simpan / Tambah Cepat jika berada di modul Klinis -->
  {#if activeModule === "clinical"}
    {#if isExpanded}
      <button
        type="button"
        onclick={onNewEncounter}
        class="flex items-center gap-3 h-11 px-4 rounded-2xl bg-white text-[#1f1f1f] dark:bg-[#1e1f20] dark:text-[#e3e3e3] shadow-xs hover:shadow-md hover:bg-[#f8fafd] dark:hover:bg-[#28292a] active:bg-[#e9eef6] border border-[#e1e5ea] dark:border-[#353738] transition-all duration-150 cursor-pointer w-full font-medium text-xs tracking-wide mb-2"
      >
        <Plus class="w-5 h-5 text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span>Simpan RME</span>
      </button>
    {:else}
      <button
        type="button"
        onclick={onNewEncounter}
        class="flex items-center justify-center w-11 h-11 rounded-2xl bg-white text-[#0b57d0] dark:bg-[#1e1f20] dark:text-[#a8c7fa] shadow-xs hover:shadow-md hover:bg-[#f8fafd] dark:hover:bg-[#28292a] active:bg-[#e9eef6] border border-[#e1e5ea] dark:border-[#353738] transition-all duration-150 cursor-pointer group mb-2"
        title="Simpan RME Pasien"
      >
        <Plus class="w-5 h-5" />
      </button>
    {/if}
  {/if}

  <!-- List Menu Item dengan Sub-Menu Terfokus -->
  <nav class="flex flex-col gap-1.5 w-full">
    {#each visibleNavigation as item}
      {@const parentActive = isParentActive(item)}
      {@const isOpen = openParentIds.includes(item.id)}
      {@const isFlyoutOpen = flyoutMenuId === item.id}

      {#if isExpanded}
        <!-- MODE DIPERLUAS: ACCORDION TREE -->
        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => toggleAccordion(item.id)}
            class="flex items-center justify-between px-3 py-2.5 rounded-full transition-colors cursor-pointer {parentActive
              ? 'bg-[#d3e3fd]/80 text-[#041e49] dark:bg-[#0842a0]/60 dark:text-[#d3e3fd] font-semibold'
              : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] hover:text-[#1f1f1f] dark:hover:text-white'}"
          >
            <div class="flex items-center gap-3">
              <item.icon
                class="w-4.5 h-4.5 {parentActive
                  ? 'text-[#0b57d0] dark:text-[#a8c7fa]'
                  : 'text-[#444746] dark:text-[#c4c7c5]'}"
              />
              <span class="text-xs">{item.label}</span>
            </div>
            {#if isOpen}
              <ChevronDown class="w-4 h-4 text-[#747775] dark:text-[#8e918f]" />
            {:else}
              <ChevronRight class="w-4 h-4 text-[#747775] dark:text-[#8e918f]" />
            {/if}
          </button>

          {#if isOpen}
            <div
              class="flex flex-col gap-0.5 pl-6 pr-1 mt-1 border-l-2 border-[#e1e5ea] dark:border-[#353738] ml-5 py-1 animate-in fade-in duration-150"
            >
              {#each item.children as sub}
                {@const isSubActive = activeNav === sub.id}
                <button
                  type="button"
                  onclick={() => handleSubmenuClick(sub.id)}
                  class="flex items-center justify-between px-3 py-2 text-xs rounded-full transition-colors cursor-pointer text-left {isSubActive
                    ? 'bg-[#0b57d0] text-white font-semibold shadow-xs'
                    : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] hover:text-[#1f1f1f] dark:hover:text-white'}"
                >
                  <span class="truncate">{sub.label}</span>
                  {#if sub.badge}
                    <span
                      class="text-[10px] px-1.5 py-0.2 rounded-md font-semibold {isSubActive
                        ? 'bg-white/20 text-white'
                        : 'bg-[#e9eef6] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5]'}"
                    >
                      {sub.badge}
                    </span>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {:else}
        <!-- MODE MINIMAL: COMPACT ICON RAIL DENGAN FLYOUT POPOVER (GOOGLE M3 PILL) -->
        <div class="relative w-full flex flex-col items-center" data-rail-item={item.id}>
          <button
            type="button"
            onclick={() => handleRailItemClick(item.id)}
            class="flex flex-col items-center justify-center w-full py-1 group cursor-pointer"
            title={item.label}
          >
            <!-- M3 Active Pill Indicator -->
            <div
              class="flex items-center justify-center w-12 h-8 rounded-full transition-all duration-200 {parentActive
                ? 'bg-[#d3e3fd] text-[#041e49] dark:bg-[#0842a0] dark:text-[#d3e3fd] shadow-xs'
                : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] hover:text-[#1f1f1f] dark:hover:text-white'}"
            >
              <item.icon
                class="w-5 h-5 transition-transform group-hover:scale-105 {parentActive
                  ? 'text-[#0b57d0] dark:text-[#d3e3fd]'
                  : 'text-[#444746] dark:text-[#c4c7c5]'}"
              />
            </div>
            <span
              class="text-[10px] tracking-tight mt-1 truncate max-w-[56px] text-center transition-colors {parentActive
                ? 'font-bold text-[#0b57d0] dark:text-[#a8c7fa]'
                : 'font-medium text-[#444746] dark:text-[#c4c7c5] group-hover:text-[#1f1f1f] dark:group-hover:text-white'}"
            >
              {item.shortLabel}
            </span>
          </button>

          <!-- Flyout Menu Popover saat hover/klik icon rail -->
          {#if isFlyoutOpen}
            <div
              class="absolute left-14 top-0 w-56 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-2xl p-2 z-50 flex flex-col gap-1 animate-in fade-in zoom-in-95 duration-150"
            >
              <div
                class="px-3 py-1.5 text-xs font-bold text-[#1f1f1f] dark:text-[#e3e3e3] border-b border-[#e1e5ea] dark:border-[#353738] mb-1 flex items-center justify-between"
              >
                <span>{item.label}</span>
                <span class="text-[10px] text-[#747775] dark:text-[#8e918f] font-normal"
                  >{item.children.length} Menu</span
                >
              </div>
              {#each item.children as sub}
                {@const isSubActive = activeNav === sub.id}
                <button
                  type="button"
                  onclick={() => handleSubmenuClick(sub.id)}
                  class="flex items-center justify-between px-3 py-2 text-xs rounded-xl transition-colors cursor-pointer text-left {isSubActive
                    ? 'bg-[#0b57d0] text-white font-semibold shadow-xs'
                    : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] hover:text-[#1f1f1f] dark:hover:text-white'}"
                >
                  <span class="truncate">{sub.label}</span>
                  {#if sub.badge}
                    <span
                      class="text-[10px] px-1.5 py-0.2 rounded-md font-semibold {isSubActive
                        ? 'bg-white/20 text-white'
                        : 'bg-[#e9eef6] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5]'}"
                    >
                      {sub.badge}
                    </span>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    {/each}
  </nav>
</aside>
