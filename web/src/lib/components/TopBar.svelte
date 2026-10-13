<script lang="ts">
  import {
    Search,
    Stethoscope,
    Wifi,
    WifiOff,
    HelpCircle,
    Settings,
    SlidersHorizontal,
    Menu,
    LogOut,
    UserCheck,
    ShieldCheck,
    LayoutGrid,
    ChevronRight,
  } from "@lucide/svelte";
  import { onMount } from "svelte";
  import { checkBackendHealth } from "$lib/api/client";
  import ThemeToggle from "./ThemeToggle.svelte";
  import type { UserProfile } from "../types";

  interface Props {
    onToggleSidebar?: () => void;
    user?: UserProfile | null;
    onLogout?: () => void;
    activeModule?: string | null;
    onNavigateToDesktop?: () => void;
  }

  let {
    onToggleSidebar,
    user,
    onLogout,
    activeModule = null,
    onNavigateToDesktop,
  }: Props = $props();

  let isProfileMenuOpen = $state(false);
  let profileContainerRef = $state<HTMLDivElement | null>(null);

  function closeProfileMenu() {
    isProfileMenuOpen = false;
  }

  function handleDocumentClick(e: MouseEvent) {
    if (
      isProfileMenuOpen &&
      profileContainerRef &&
      !profileContainerRef.contains(e.target as Node)
    ) {
      closeProfileMenu();
    }
  }

  let backendStatus = $state<{
    loading: boolean;
    connected: boolean;
    data: any;
  }>({
    loading: true,
    connected: false,
    data: null,
  });

  async function refreshHealth() {
    backendStatus.loading = true;
    const res = await checkBackendHealth();
    if (res && res.status === "success") {
      backendStatus = {
        loading: false,
        connected: true,
        data: res.data,
      };
    } else {
      backendStatus = {
        loading: false,
        connected: false,
        data: res,
      };
    }
  }

  onMount(() => {
    refreshHealth();
    const interval = setInterval(refreshHealth, 15000);
    return () => clearInterval(interval);
  });
</script>

<svelte:window onclick={handleDocumentClick} />

<header
  class="h-16 px-4 bg-[#f8fafd] dark:bg-[#18191a] border-b border-transparent dark:border-[#353738] flex items-center justify-between select-none sticky top-0 z-50 transition-colors"
>
  <!-- Brand Identity & Module Context Switcher -->
  <div class="flex items-center gap-1.5 shrink-0">
    {#if activeModule}
      <button
        type="button"
        onclick={onToggleSidebar}
        class="p-2 text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] rounded-full transition-colors cursor-pointer"
        title="Buka / Tutup Menu Modul"
      >
        <Menu class="w-5 h-5" />
      </button>
    {/if}

    <button
      type="button"
      onclick={onNavigateToDesktop}
      class="flex items-center gap-2 p-1 rounded-xl hover:bg-[#e9eef6] dark:hover:bg-[#28292a] transition-colors cursor-pointer text-left group"
      title="Buka Desktop Launcher (Semua Modul)"
    >
      <div
        class="flex items-center justify-center w-9 h-9 rounded-xl bg-[#0b57d0] text-white shadow-xs group-hover:scale-105 transition-transform"
      >
        <LayoutGrid class="w-4.5 h-4.5" />
      </div>
      <div>
        <div class="flex items-center gap-1.5">
          <span class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight"
            >HOSIM</span
          >
        </div>
      </div>
    </button>

    {#if activeModule}
      <ChevronRight class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f] shrink-0" />
      <span
        class="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa] border border-[#d3e3fd] dark:border-[#0b57d0]/40 shrink-0"
      >
        {activeModule === "master"
          ? "Master Data RS"
          : activeModule === "clinical"
            ? "Pelayanan Klinis"
            : activeModule === "auth"
              ? "Pengguna & Akses"
              : activeModule === "accounting"
                ? "Akuntansi & Keuangan RS"
                : activeModule === "finance"
                  ? "Tarif & Keuangan RS"
                  : activeModule.toUpperCase()}
      </span>
    {/if}
  </div>

  <!-- Google Drive Style Search Bar (Pencarian Pasien & No. RM) -->
  <div class="flex-1 max-w-2xl mx-4 hidden md:block">
    <div
      class="relative flex items-center h-12 w-full rounded-full bg-[#e9eef6] dark:bg-[#28292a] hover:bg-[#e1e5ea] dark:hover:bg-[#353738] focus-within:bg-white dark:focus-within:bg-[#1e1f20] focus-within:shadow-md transition-all duration-200 px-4 border border-transparent focus-within:border-[#c4c7c5]/50 dark:focus-within:border-[#0b57d0]/60"
    >
      <Search class="w-5 h-5 text-[#444746] dark:text-[#8e918f] mr-3 shrink-0" />
      <input
        type="text"
        placeholder="Cari data pasien, No. Rekam Medis (RM), NIK, atau diagnosa..."
        class="w-full bg-transparent text-sm text-[#1f1f1f] dark:text-[#e3e3e3] placeholder:text-[#444746] dark:placeholder:text-[#8e918f] focus:outline-none"
      />
      <button
        type="button"
        class="p-1 text-[#444746] dark:text-[#8e918f] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] hover:bg-[#d8dde3]/60 dark:hover:bg-[#353738] rounded-full cursor-pointer ml-1"
        title="Filter Lanjutan"
      >
        <SlidersHorizontal class="w-4 h-4" />
      </button>
    </div>
  </div>

  <!-- Status Backend & Profil Pengguna -->
  <div class="flex items-center gap-2 sm:gap-3 relative">
    <!-- Theme Switcher (Terang / Gelap / Sistem) -->
    <ThemeToggle />

    <!-- Help & Settings Icons -->
    <button
      type="button"
      class="p-2 text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] rounded-full transition-colors cursor-pointer hidden sm:block"
      title="Bantuan"
    >
      <HelpCircle class="w-5 h-5" />
    </button>
    <button
      type="button"
      class="p-2 text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] rounded-full transition-colors cursor-pointer hidden sm:block"
      title="Pengaturan"
    >
      <Settings class="w-5 h-5" />
    </button>

    <!-- Google Account Profile Circle with Dropdown Menu -->
    <div bind:this={profileContainerRef} class="relative">
      <button
        type="button"
        onclick={() => (isProfileMenuOpen = !isProfileMenuOpen)}
        class="w-9 h-9 rounded-full bg-[#0b57d0] text-white flex items-center justify-center text-sm font-semibold cursor-pointer shadow-xs hover:ring-4 hover:ring-[#d3e3fd]/60 transition-all select-none"
        title={user?.name || "dr. Hendra Wijaya, Sp.B"}
      >
        {(user?.name || "Hendra").charAt(0).toUpperCase()}
      </button>

      <!-- Profile Popover Google Style -->
      {#if isProfileMenuOpen}
        <div
          class="absolute right-0 top-full pt-2 z-50 animate-in fade-in duration-150"
        >
          <div
            class="w-72 bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738] shadow-lg p-5"
          >
            <div
              class="flex flex-col items-center text-center pb-4 border-b border-[#e1e5ea] dark:border-[#353738]"
            >
              <div
                class="w-14 h-14 rounded-full {user?.role === 'ADMIN'
                  ? 'bg-purple-600'
                  : 'bg-[#0b57d0]'} text-white flex items-center justify-center text-xl font-bold mb-2 shadow-xs transition-colors"
              >
                {(user?.name || "Hendra").charAt(0).toUpperCase()}
              </div>
              <div class="font-semibold text-sm text-[#1f1f1f] dark:text-[#e3e3e3]">
                {user?.name || "dr. Hendra Wijaya, Sp.B"}
              </div>
              <div class="text-xs text-[#444746] dark:text-[#8e918f] mt-0.5">
                {user?.email || "hendra@hosim.local"}
              </div>
              <div
                class="mt-2 inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {user?.role ===
                'ADMIN'
                  ? 'bg-purple-100 dark:bg-purple-950/50 text-purple-800 dark:text-purple-300 border border-purple-200 dark:border-purple-800'
                  : 'bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa] border border-[#d3e3fd] dark:border-[#0b57d0]/40'}"
              >
                <UserCheck class="w-3 h-3" />
                <span>Role: {user?.role || "DOCTOR"}</span>
              </div>
            </div>

            <div class="pt-3 flex flex-col gap-1">
              <button
                type="button"
                onclick={() => {
                  closeProfileMenu();
                  if (onLogout) onLogout();
                }}
                class="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-xs font-semibold text-[#b3261e] dark:text-rose-400 hover:bg-[#fce8e6] dark:hover:bg-rose-950/40 transition-colors cursor-pointer"
              >
                <LogOut class="w-4 h-4" />
                <span>Keluar dari Akun (Logout)</span>
              </button>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</header>
