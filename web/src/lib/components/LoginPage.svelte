<script lang="ts">
  import {
    Stethoscope,
    User,
    Eye,
    EyeOff,
    ArrowRight,
    ShieldCheck,
    AlertCircle,
    Wifi,
    WifiOff,
    Sparkles
  } from '@lucide/svelte';
  import { auth } from '../stores/auth.svelte';
  import { checkBackendHealth } from '../api';
  import { onMount } from 'svelte';
  import type { UserProfile } from '../types';

  interface Props {
    onLoginSuccess?: (user: UserProfile) => void;
  }

  let { onLoginSuccess = () => {} }: Props = $props();

  let username = $state('dokter');
  let password = $state('dokter123');
  let showPassword = $state(false);
  let rememberMe = $state(true);
  let isLoading = $state(false);
  let errorMessage = $state('');

  let backendStatus = $state<{ loading: boolean; connected: boolean }>({
    loading: true,
    connected: false
  });

  async function checkHealth() {
    backendStatus.loading = true;
    const res = await checkBackendHealth();
    backendStatus = {
      loading: false,
      connected: !!(res && res.status === 'success')
    };
  }

  onMount(() => {
    checkHealth();
  });

  async function handleSubmit(e?: Event) {
    if (e) e.preventDefault();
    if (!username.trim() || !password.trim()) {
      errorMessage = 'Mohon masukkan username dan kata sandi Anda.';
      return;
    }

    isLoading = true;
    errorMessage = '';

    try {
      const res = await auth.login(username.trim(), password);
      if (res && res.success) {
        if (onLoginSuccess) onLoginSuccess(res.user);
      }
    } catch (err: any) {
      errorMessage = (err && err.message) || 'Username atau password salah. Cek koneksi backend.';
    } finally {
      isLoading = false;
    }
  }

  function setQuickAccount(user: string, pass: string): void {
    username = user;
    password = pass;
    errorMessage = '';
  }
</script>

<div class="min-h-screen bg-[#f8fafd] flex flex-col justify-between p-4 sm:p-6 font-sans select-none">
  <!-- Top Bar Minimalis -->
  <div class="flex items-center justify-between max-w-5xl mx-auto w-full pt-2">
    <div class="flex items-center gap-2.5">
      <div class="w-9 h-9 rounded-xl bg-[#0b57d0] text-white flex items-center justify-center shadow-xs">
        <Stethoscope class="w-4.5 h-4.5" />
      </div>
      <span class="text-base font-semibold text-[#1f1f1f] tracking-tight">HOSIM</span>
    </div>

    <!-- Status Koneksi Backend -->
    <div class="flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium {backendStatus.connected ? 'bg-[#c4eed0]/60 text-[#072711] border border-[#146c2e]/20' : 'bg-[#f0f4f9] text-[#444746] border border-[#e1e5ea]'}">
      {#if backendStatus.connected}
        <span class="w-2 h-2 rounded-full bg-[#1e8e3e] animate-pulse"></span>
        <Wifi class="w-3.5 h-3.5 text-[#146c2e]" />
        <span>Go Gin API: Terhubung</span>
      {:else}
        <span class="w-2 h-2 rounded-full bg-[#bdc1c6]"></span>
        <WifiOff class="w-3.5 h-3.5 text-[#747775]" />
        <span>Go API Standby / Demo Mode</span>
      {/if}
    </div>
  </div>

  <!-- Card Login Utama (Google Account / Material You Style) -->
  <div class="max-w-[460px] w-full mx-auto my-8">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] p-8 sm:p-10 shadow-xs">
      <!-- Header Login -->
      <div class="flex flex-col items-center text-center mb-7">
        <div class="w-14 h-14 rounded-2xl bg-[#d3e3fd] text-[#041e49] flex items-center justify-center mb-4 shadow-xs">
          <Stethoscope class="w-7 h-7 text-[#0b57d0]" />
        </div>
        <h1 class="text-2xl font-normal text-[#1f1f1f] tracking-tight">Masuk ke Akun</h1>
        <p class="text-sm text-[#444746] mt-1.5">Sistem Informasi Manajemen Rumah Sakit</p>
      </div>

      <!-- Pesan Error -->
      {#if errorMessage}
        <div class="mb-5 p-3.5 rounded-2xl bg-[#fce8e6] border border-[#fad2cf] text-[#c5221f] text-xs flex items-start gap-2.5 animate-in fade-in duration-200">
          <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
          <span>{errorMessage}</span>
        </div>
      {/if}

      <!-- Form Login -->
      <form onsubmit={handleSubmit} class="flex flex-col gap-4">
        <!-- Input Username / Email -->
        <div>
          <label for="login-username" class="block text-xs font-semibold text-[#444746] mb-1.5">
            Nama Pengguna atau Email
          </label>
          <div class="relative flex items-center">
            <input
              id="login-username"
              type="text"
              bind:value={username}
              placeholder="Contoh: dokter atau admin"
              autocomplete="username"
              required
              class="w-full h-12 px-4 rounded-xl bg-white border border-[#747775]/50 text-[#1f1f1f] text-sm placeholder:text-[#747775] focus:outline-none focus:border-[#0b57d0] focus:ring-3 focus:ring-[#0b57d0]/15 transition-all"
            />
            <User class="w-4 h-4 text-[#747775] absolute right-4 pointer-events-none" />
          </div>
        </div>

        <!-- Input Password -->
        <div>
          <div class="flex items-center justify-between mb-1.5">
            <label for="login-password" class="text-xs font-semibold text-[#444746]">
              Kata Sandi
            </label>
            <a href="#forgot" onclick={(e) => { e.preventDefault(); alert('Hubungi administrator IT Rumah Sakit untuk reset password.'); }} class="text-xs text-[#0b57d0] hover:underline font-medium">
              Lupa sandi?
            </a>
          </div>
          <div class="relative flex items-center">
            <input
              id="login-password"
              type={showPassword ? 'text' : 'password'}
              bind:value={password}
              placeholder="Masukkan kata sandi..."
              autocomplete="current-password"
              required
              class="w-full h-12 px-4 pr-11 rounded-xl bg-white border border-[#747775]/50 text-[#1f1f1f] text-sm placeholder:text-[#747775] focus:outline-none focus:border-[#0b57d0] focus:ring-3 focus:ring-[#0b57d0]/15 transition-all"
            />
            <button
              type="button"
              onclick={() => showPassword = !showPassword}
              class="p-2 text-[#747775] hover:text-[#1f1f1f] absolute right-2.5 rounded-full cursor-pointer transition-colors"
              title={showPassword ? 'Sembunyikan Kata Sandi' : 'Tampilkan Kata Sandi'}
            >
              {#if showPassword}
                <EyeOff class="w-4 h-4" />
              {:else}
                <Eye class="w-4 h-4" />
              {/if}
            </button>
          </div>
        </div>

        <!-- Checkbox Ingat Saya -->
        <div class="flex items-center gap-2 mt-1">
          <input
            id="remember-me"
            type="checkbox"
            bind:checked={rememberMe}
            class="w-4 h-4 rounded accent-[#0b57d0] cursor-pointer"
          />
          <label for="remember-me" class="text-xs text-[#444746] cursor-pointer">
            Ingat saya di perangkat ini
          </label>
        </div>

        <!-- Tombol Masuk -->
        <div class="mt-4 flex items-center justify-between">
          <span class="text-xs text-[#444746]">Akses SIMRS v1.0</span>
          <button
            type="submit"
            disabled={isLoading}
            class="h-11 px-7 rounded-full bg-[#0b57d0] text-white font-medium text-sm hover:bg-[#0842a0] active:bg-[#073887] shadow-xs hover:shadow transition-all duration-150 flex items-center gap-2 cursor-pointer disabled:opacity-50"
          >
            {#if isLoading}
              <span class="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin"></span>
              <span>Memverifikasi...</span>
            {:else}
              <span>Masuk</span>
              <ArrowRight class="w-4 h-4" />
            {/if}
          </button>
        </div>
      </form>

      <!-- Quick Access Akun Demo Google Style -->
      <div class="mt-8 pt-6 border-t border-[#e1e5ea]">
        <div class="flex items-center gap-1.5 text-xs text-[#444746] font-medium mb-3">
          <Sparkles class="w-3.5 h-3.5 text-[#0b57d0]" />
          <span>Pilih Akun Cepat (Quick Access):</span>
        </div>

        <div class="grid grid-cols-2 gap-2.5">
          <button
            type="button"
            onclick={() => setQuickAccount('dokter', 'dokter123')}
            class="p-2.5 rounded-xl border border-[#e1e5ea] bg-[#f8fafd] hover:bg-[#e9eef6] hover:border-[#0b57d0]/40 text-left transition-all cursor-pointer group"
          >
            <div class="text-xs font-semibold text-[#1f1f1f] group-hover:text-[#0b57d0]">Dokter Spesialis</div>
            <div class="text-[11px] text-[#444746] font-mono mt-0.5">dokter / dokter123</div>
          </button>

          <button
            type="button"
            onclick={() => setQuickAccount('admin', 'admin123')}
            class="p-2.5 rounded-xl border border-[#e1e5ea] bg-[#f8fafd] hover:bg-[#e9eef6] hover:border-[#0b57d0]/40 text-left transition-all cursor-pointer group"
          >
            <div class="text-xs font-semibold text-[#1f1f1f] group-hover:text-[#0b57d0]">Administrator</div>
            <div class="text-[11px] text-[#444746] font-mono mt-0.5">admin / admin123</div>
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Footer Standar Keamanan Klinis -->
  <div class="max-w-5xl mx-auto w-full flex flex-col sm:flex-row items-center justify-between text-xs text-[#747775] gap-2 pb-2">
    <div class="flex items-center gap-2">
      <ShieldCheck class="w-4 h-4 text-[#137333]" />
      <span>Standar Kemenkes RI • Terintegrasi SATUSEHAT & BPJS Kesehatan</span>
    </div>
    <div class="flex items-center gap-4">
      <a href="#privasi" onclick={(e) => e.preventDefault()} class="hover:underline">Kebijakan Privasi RME</a>
      <span>•</span>
      <a href="#syarat" onclick={(e) => e.preventDefault()} class="hover:underline">Ketentuan Layanan</a>
    </div>
  </div>
</div>
