<script lang="ts">
  import { Search, Stethoscope, Wifi, WifiOff, HelpCircle, Settings, SlidersHorizontal, Menu, LogOut, UserCheck, ShieldCheck } from '@lucide/svelte';
  import { onMount } from 'svelte';
  import { checkBackendHealth } from '../api';
  import { auth } from '../stores/auth.svelte';
  import type { UserProfile } from '../types';

  interface Props {
    onToggleSidebar?: () => void;
    user?: UserProfile | null;
    onLogout?: () => void;
  }

  let { onToggleSidebar, user, onLogout }: Props = $props();

  let isProfileMenuOpen = $state(false);

  let backendStatus = $state<{ loading: boolean; connected: boolean; data: any }>({
    loading: true,
    connected: false,
    data: null
  });

  async function refreshHealth() {
    backendStatus.loading = true;
    const res = await checkBackendHealth();
    if (res && res.status === 'success') {
      backendStatus = {
        loading: false,
        connected: true,
        data: res.data
      };
    } else {
      backendStatus = {
        loading: false,
        connected: false,
        data: res
      };
    }
  }

  onMount(() => {
    refreshHealth();
    const interval = setInterval(refreshHealth, 15000);
    return () => clearInterval(interval);
  });
</script>

<header class="h-16 px-4 bg-[#f8fafd] flex items-center justify-between select-none sticky top-0 z-50">
  <!-- Brand Identity ala Google Drive dengan Hamburger Toggle -->
  <div class="flex items-center gap-2 shrink-0">
    <button
      type="button"
      onclick={onToggleSidebar}
      class="p-2.5 text-[#444746] hover:text-[#1f1f1f] hover:bg-[#e9eef6] rounded-full transition-colors cursor-pointer"
      title="Buka / Tutup Menu Samping"
    >
      <Menu class="w-5 h-5" />
    </button>
    <div class="flex items-center gap-2.5 pl-1">
      <div class="flex items-center justify-center w-10 h-10 rounded-xl bg-[#0b57d0] text-white shadow-xs">
        <Stethoscope class="w-5 h-5" />
      </div>
      <div>
        <div class="flex items-center gap-1.5">
          <span class="text-lg font-semibold text-[#1f1f1f] tracking-tight">HOSIM</span>
          <span class="text-xs text-[#444746] font-normal">Klinis</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Google Drive Style Search Bar (Pencarian Pasien & No. RM) -->
  <div class="flex-1 max-w-2xl mx-4 hidden md:block">
    <div class="relative flex items-center h-12 w-full rounded-full bg-[#e9eef6] hover:bg-[#e1e5ea] focus-within:bg-white focus-within:shadow-md transition-all duration-200 px-4 border border-transparent focus-within:border-[#c4c7c5]/50">
      <Search class="w-5 h-5 text-[#444746] mr-3 shrink-0" />
      <input
        type="text"
        placeholder="Cari data pasien, No. Rekam Medis (RM), NIK, atau diagnosa..."
        class="w-full bg-transparent text-sm text-[#1f1f1f] placeholder:text-[#444746] focus:outline-none"
      />
      <button type="button" class="p-1 text-[#444746] hover:text-[#1f1f1f] hover:bg-[#d8dde3]/60 rounded-full cursor-pointer ml-1" title="Filter Lanjutan">
        <SlidersHorizontal class="w-4 h-4" />
      </button>
    </div>
  </div>

  <!-- Status Backend & Profil Pengguna -->
  <div class="flex items-center gap-3 relative">
    <!-- Backend Connection Indicator -->
    <div class="flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium {backendStatus.connected ? 'bg-[#c4eed0]/60 text-[#072711] border border-[#146c2e]/20' : 'bg-[#f0f4f9] text-[#444746] border border-[#e1e5ea]'}">
      {#if backendStatus.connected}
        <span class="w-2 h-2 rounded-full bg-[#1e8e3e]"></span>
        <Wifi class="w-3.5 h-3.5 text-[#146c2e]" />
        <span class="hidden sm:inline">Go Gin API: Terhubung</span>
        <span class="sm:hidden">Online</span>
      {:else}
        <span class="w-2 h-2 rounded-full bg-[#bdc1c6]"></span>
        <WifiOff class="w-3.5 h-3.5 text-[#747775]" />
        <span class="hidden sm:inline">Go API: Offline</span>
        <span class="sm:hidden">Offline</span>
      {/if}
    </div>

    <!-- Active Role Indicator Pill -->
    {#if user?.role === 'ADMIN'}
      <div class="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-purple-50 text-purple-700 border border-purple-200 shadow-xs" title="Superadmin: Akses Penuh">
        <ShieldCheck class="w-3.5 h-3.5 text-purple-600" />
        <span>Admin Sistem</span>
      </div>
    {:else}
      <div class="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-blue-50 text-blue-700 border border-blue-200 shadow-xs" title="Praktisi Medis: Hak Akses Klinis">
        <Stethoscope class="w-3.5 h-3.5 text-blue-600" />
        <span>Dokter / Nakes</span>
      </div>
    {/if}

    <!-- Help & Settings Icons -->
    <button type="button" class="p-2 text-[#444746] hover:text-[#1f1f1f] hover:bg-[#e9eef6] rounded-full transition-colors cursor-pointer hidden sm:block" title="Bantuan">
      <HelpCircle class="w-5 h-5" />
    </button>
    <button type="button" class="p-2 text-[#444746] hover:text-[#1f1f1f] hover:bg-[#e9eef6] rounded-full transition-colors cursor-pointer hidden sm:block" title="Pengaturan">
      <Settings class="w-5 h-5" />
    </button>

    <!-- Google Account Profile Circle with Dropdown Menu -->
    <button
      type="button"
      onclick={() => isProfileMenuOpen = !isProfileMenuOpen}
      class="w-9 h-9 rounded-full bg-[#0b57d0] text-white flex items-center justify-center text-sm font-semibold cursor-pointer shadow-xs hover:ring-4 hover:ring-[#d3e3fd]/60 transition-all select-none"
      title="{user?.name || 'dr. Hendra Wijaya, Sp.B'}"
    >
      {(user?.name || 'Hendra').charAt(0).toUpperCase()}
    </button>

    <!-- Profile Popover Google Style -->
    {#if isProfileMenuOpen}
      <div class="absolute right-0 top-12 mt-2 w-80 bg-white rounded-3xl border border-[#e1e5ea] shadow-lg p-5 z-50 animate-in fade-in duration-150">
        <div class="flex flex-col items-center text-center pb-4 border-b border-[#e1e5ea]">
          <div class="w-14 h-14 rounded-full {user?.role === 'ADMIN' ? 'bg-purple-600' : 'bg-[#0b57d0]'} text-white flex items-center justify-center text-xl font-bold mb-2 shadow-xs transition-colors">
            {(user?.name || 'Hendra').charAt(0).toUpperCase()}
          </div>
          <div class="font-semibold text-sm text-[#1f1f1f]">{user?.name || 'dr. Hendra Wijaya, Sp.B'}</div>
          <div class="text-xs text-[#444746] mt-0.5">{user?.email || 'hendra@hosim.local'}</div>
          <div class="mt-2 inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {user?.role === 'ADMIN' ? 'bg-purple-100 text-purple-800' : 'bg-[#e8f0fe] text-[#0b57d0]'}">
            <UserCheck class="w-3 h-3" />
            <span>Role: {user?.role || 'DOCTOR'}</span>
          </div>
        </div>

        <!-- Role Simulation Quick Switch -->
        <div class="py-3 border-b border-[#e1e5ea] flex flex-col gap-1.5">
          <div class="text-[10px] font-bold text-[#747775] uppercase tracking-wider px-1">
            Simulasi Role (Uji Coba Cepat):
          </div>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              onclick={async () => {
                await auth.login('admin', 'admin123');
                isProfileMenuOpen = false;
              }}
              class="px-2.5 py-2 rounded-xl text-xs font-semibold border text-center transition-all cursor-pointer {user?.role === 'ADMIN' ? 'bg-purple-100 text-purple-900 border-purple-300 ring-2 ring-purple-200' : 'bg-[#f0f4f9] text-[#444746] hover:bg-purple-50 hover:text-purple-700 border-[#e1e5ea]'}"
            >
              👑 Admin Sistem
            </button>
            <button
              type="button"
              onclick={async () => {
                await auth.login('dokter', 'dokter123');
                isProfileMenuOpen = false;
              }}
              class="px-2.5 py-2 rounded-xl text-xs font-semibold border text-center transition-all cursor-pointer {user?.role === 'DOCTOR' ? 'bg-blue-100 text-blue-900 border-blue-300 ring-2 ring-blue-200' : 'bg-[#f0f4f9] text-[#444746] hover:bg-blue-50 hover:text-blue-700 border-[#e1e5ea]'}"
            >
              🩺 Dokter / Nakes
            </button>
          </div>
        </div>

        <div class="pt-3 flex flex-col gap-1">
          <button
            type="button"
            onclick={() => { isProfileMenuOpen = false; if (onLogout) onLogout(); }}
            class="w-full flex items-center justify-center gap-2 px-3 py-2 rounded-xl text-xs font-semibold text-[#b3261e] hover:bg-[#fce8e6] transition-colors cursor-pointer"
          >
            <LogOut class="w-4 h-4" />
            <span>Keluar dari Akun (Logout)</span>
          </button>
        </div>
      </div>
    {/if}
  </div>
</header>
