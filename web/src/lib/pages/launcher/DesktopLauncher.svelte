<script lang="ts">
  import {
    Search,
    Sparkles,
    Calendar,
    Clock,
    Activity,
    X,
  } from "@lucide/svelte";
  import { auth } from "../../stores/auth.svelte";
  import {
    modules,
    type ModuleDefinition,
  } from "$lib/config/launcherModules";
  import { onMount } from "svelte";
  import LauncherModuleCard from "./components/LauncherModuleCard.svelte";

  export type { ModuleDefinition };

  interface Props {
    onSelectModule: (moduleId: string, defaultNavId: string) => void;
    highlightModuleId?: string | null;
    savedScrollTop?: number;
  }

  let {
    onSelectModule,
    highlightModuleId = null,
    savedScrollTop = 0,
  }: Props = $props();

  let activeHighlight = $state<string | null>(null);

  $effect(() => {
    activeHighlight = highlightModuleId;
  });

  // Jika modul yang disorot berada di luar kategori yang sedang dipilih, reset ke 'all' agar kartu ter-render
  $effect(() => {
    if (highlightModuleId) {
      const targetMod = modules.find((m) => m.id === highlightModuleId);
      if (targetMod && selectedCategory !== "all" && targetMod.category !== selectedCategory) {
        selectedCategory = "all";
      }
    }
  });

  onMount(() => {
    if (highlightModuleId) {
      const timer = setTimeout(() => {
        const el = document.getElementById(`launcher-module-${highlightModuleId}`);
        if (el) {
          el.scrollIntoView({ behavior: "smooth", block: "center" });
        } else if (savedScrollTop > 0) {
          const mainEl = document.querySelector("main");
          if (mainEl) mainEl.scrollTop = savedScrollTop;
        }
      }, 80);

      const fadeTimer = setTimeout(() => {
        activeHighlight = null;
      }, 2500);

      return () => {
        clearTimeout(timer);
        clearTimeout(fadeTimer);
      };
    } else if (savedScrollTop > 0) {
      const timer = setTimeout(() => {
        const mainEl = document.querySelector("main");
        if (mainEl) mainEl.scrollTop = savedScrollTop;
      }, 80);
      return () => clearTimeout(timer);
    }
  });

  let searchQuery = $state("");

  // Tanggal dan waktu saat ini
  const currentDateFormatted = new Intl.DateTimeFormat("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(new Date());

  // Filter modul berdasarkan pencarian dan role/permission
  let filteredModules = $derived(
    modules.filter((m) => {
      if (m.permission && !auth.hasPermission(m.permission)) {
        return false;
      }
      if (!searchQuery.trim()) return true;
      const q = searchQuery.toLowerCase();
      return (
        m.title.toLowerCase().includes(q) ||
        m.shortTitle.toLowerCase().includes(q) ||
        m.description.toLowerCase().includes(q) ||
        m.features.some((f) => f.toLowerCase().includes(q))
      );
    }),
  );

  // Filter kategori yang aktif
  let selectedCategory = $state<
    "all" | "clinical" | "support_and_diagnostic" | "admin_finance"
  >("all");

  // Pengelompokan modul per kategori
  let categories = $derived([
    {
      id: "clinical" as const,
      label: "Pelayanan Pasien & Asuhan Medis",
      modules: filteredModules.filter((m) => m.category === "clinical"),
    },
    {
      id: "support_and_diagnostic" as const,
      label: "Penunjang Medis & Rantai Pasok",
      modules: filteredModules.filter(
        (m) => m.category === "support_and_diagnostic",
      ),
    },
    {
      id: "admin_finance" as const,
      label: "Administrasi, Keuangan & Tata Kelola",
      modules: filteredModules.filter((m) => m.category === "admin_finance"),
    },
  ]);

  let displayedCategories = $derived(
    selectedCategory === "all"
      ? categories.filter((c) => c.modules.length > 0)
      : categories.filter(
          (c) => c.id === selectedCategory && c.modules.length > 0,
        ),
  );

  let categoryTabs = $derived([
    { id: "all" as const, label: "Semua Modul", count: filteredModules.length },
    {
      id: "clinical" as const,
      label: "Pelayanan Medis",
      count: filteredModules.filter((m) => m.category === "clinical").length,
    },
    {
      id: "support_and_diagnostic" as const,
      label: "Penunjang & Logistik",
      count: filteredModules.filter(
        (m) => m.category === "support_and_diagnostic",
      ).length,
    },
    {
      id: "admin_finance" as const,
      label: "Admin & Tata Kelola",
      count: filteredModules.filter((m) => m.category === "admin_finance")
        .length,
    },
  ]);
</script>

<div class="flex flex-col gap-4 pb-12 animate-in fade-in duration-200">
  <!-- Banner Header Launcher (Compact Hero Banner) -->
  <div
    class="relative overflow-hidden rounded-2xl bg-gradient-to-r from-[#0b57d0] via-[#1a73e8] to-[#174ea6] p-4 sm:p-5 text-white shadow-sm"
  >
    <div
      class="absolute right-0 top-0 -mt-6 -mr-6 w-56 h-56 rounded-full bg-white/10 blur-2xl pointer-events-none"
    ></div>

    <div
      class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-3"
    >
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span
            class="text-[10px] font-bold tracking-wider uppercase px-2 py-0.5 rounded-full bg-white/20 text-white backdrop-blur-xs flex items-center gap-1"
          >
            <Sparkles class="w-3 h-3 text-amber-300" />
            HOSIM-GO Modular Core
          </span>
          <div class="flex items-center gap-1 text-[11px] text-blue-100">
            <Calendar class="w-3 h-3" />
            <span>{currentDateFormatted}</span>
          </div>
        </div>

        <h1 class="text-base sm:text-xl font-bold tracking-tight text-white">
          Selamat datang, {auth.user?.name || "Tenaga Medis"}
        </h1>
        <p class="text-xs text-blue-100/90 mt-0.5 max-w-xl">
          Pilih modul di bawah ini untuk memulai pelayanan klinis, pencatatan
          RME, atau administrasi rumah sakit.
        </p>
      </div>

      <!-- Quick Search Bar (Prominent in Hero) -->
      <div class="w-full md:w-80 shrink-0">
        <div
          class="relative flex items-center bg-white/15 hover:bg-white/20 focus-within:bg-white focus-within:shadow-md focus-within:text-[#1f1f1f] rounded-xl px-3 py-1.5 transition-all backdrop-blur-xs border border-white/20 focus-within:border-white"
        >
          <Search
            class="w-4 h-4 text-blue-100 focus-within:text-[#747775] mr-2 shrink-0"
          />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari modul (cth: Poli, Resep, Billing)..."
            class="w-full bg-transparent text-xs placeholder:text-blue-200/80 focus:text-[#1f1f1f] focus:outline-none pr-5"
          />
          {#if searchQuery}
            <button
              type="button"
              onclick={() => (searchQuery = "")}
              class="absolute right-2.5 text-blue-200 hover:text-white focus-within:text-[#747775] focus-within:hover:text-[#1f1f1f] cursor-pointer"
              title="Hapus pencarian"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <!-- Quick Filter Category Tabs -->
  <div class="flex items-center justify-between gap-2 overflow-x-auto pb-0.5">
    <div class="flex items-center gap-1.5 shrink-0">
      {#each categoryTabs as tab}
        <button
          type="button"
          onclick={() => (selectedCategory = tab.id)}
          class="px-2.5 py-1 rounded-lg text-xs font-medium transition-all shrink-0 flex items-center gap-1.5 cursor-pointer {selectedCategory ===
          tab.id
            ? 'bg-[#0b57d0] text-white shadow-xs font-semibold'
            : 'bg-[#f1f3f4] dark:bg-surface-container text-[#444746] dark:text-on-surface-variant hover:bg-[#e1e5ea] dark:hover:bg-surface-container-high'}"
        >
          <span>{tab.label}</span>
          <span
            class="text-[10px] px-1.5 py-0.2 rounded-full {selectedCategory ===
            tab.id
              ? 'bg-white/20 text-white'
              : 'bg-white dark:bg-surface-container-high text-[#5f6368] dark:text-on-surface-variant'}"
          >
            {tab.count}
          </span>
        </button>
      {/each}
    </div>

    {#if searchQuery}
      <span
        class="text-[11px] text-[#747775] dark:text-on-surface-variant shrink-0 font-medium hidden sm:inline"
      >
        Menampilkan hasil untuk: "<span class="text-[#0b57d0] dark:text-primary font-semibold"
          >{searchQuery}</span
        >"
      </span>
    {/if}
  </div>

  <!-- Sections of Domains -->
  {#each displayedCategories as category}
    <div class="flex flex-col gap-2.5">
      <div
        class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-outline-variant pb-1.5"
      >
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 rounded-full bg-[#0b57d0]"></div>
          <h2
            class="text-xs sm:text-sm font-bold text-[#1f1f1f] dark:text-on-surface tracking-tight"
          >
            {category.label}
          </h2>
        </div>
        <span class="text-[10px] font-medium text-[#747775] dark:text-on-surface-variant"
          >{category.modules.length} Modul</span
        >
      </div>

      <!-- Module Grid -->
      <div
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-3"
      >
        {#each category.modules as mod}
          <LauncherModuleCard
            {mod}
            onSelect={onSelectModule}
            isHighlighted={activeHighlight === mod.id}
          />
        {/each}
      </div>
    </div>
  {/each}

  {#if filteredModules.length === 0}
    <div
      class="p-8 text-center bg-white dark:bg-surface-container rounded-2xl border border-[#e1e5ea] dark:border-outline-variant flex flex-col items-center justify-center gap-2 text-[#747775] dark:text-on-surface-variant"
    >
      <Search class="w-6 h-6 text-[#bdc1c6] dark:text-outline" />
      <div class="font-semibold text-xs text-[#1f1f1f] dark:text-on-surface">
        Modul Tidak Ditemukan
      </div>
      <p class="text-[11px] max-w-sm">
        Tidak ada modul yang cocok dengan kata kunci "{searchQuery}".
      </p>
      <button
        type="button"
        onclick={() => {
          searchQuery = "";
          selectedCategory = "all";
        }}
        class="mt-1 text-xs text-[#0b57d0] dark:text-primary hover:underline font-medium cursor-pointer"
      >
        Reset pencarian
      </button>
    </div>
  {/if}
</div>
