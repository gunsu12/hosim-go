<script lang="ts">
  import {
    ArrowLeft,
    Clock,
    Layers,
    Server,
    ShieldCheck,
    Terminal,
    Sparkles,
    CheckCircle2
  } from '@lucide/svelte';
  import M3Button from '../../components/m3/M3Button.svelte';

  interface Props {
    title: string;
    description: string;
    domainPackage: string;
    backendStatus?: string;
    plannedEndpoints?: string[];
    onBackToLauncher: () => void;
  }

  let {
    title,
    description,
    domainPackage,
    backendStatus = 'Domain entity & contracts tersedia di Go backend',
    plannedEndpoints = [],
    onBackToLauncher
  }: Props = $props();
</script>

<div class="max-w-4xl mx-auto w-full py-6 flex flex-col gap-6 animate-in fade-in duration-150">
  <!-- Top Navigation Back -->
  <div class="flex items-center justify-between">
    <button
      type="button"
      onclick={onBackToLauncher}
      class="inline-flex items-center gap-2 text-xs font-semibold text-[#0b57d0] hover:text-[#042f77] px-3 py-1.5 rounded-full hover:bg-[#e8f0fe] transition-colors cursor-pointer"
    >
      <ArrowLeft class="w-4 h-4" />
      <span>Kembali ke Beranda Modul (Desktop Launcher)</span>
    </button>
    <span class="text-xs px-2.5 py-1 rounded-full bg-[#fef7e0] text-[#7a4100] border border-[#fce8b2] font-medium">
      Workspace Preview Mode
    </span>
  </div>

  <!-- Hero Header -->
  <div class="p-6 bg-[#f8fafd] rounded-3xl border border-[#e1e5ea] flex flex-col gap-3">
    <div class="flex items-center gap-2 text-xs font-bold text-[#0b57d0] uppercase tracking-wider">
      <Layers class="w-4 h-4" />
      <span>Konteks Modul: {domainPackage}</span>
    </div>
    <h2 class="text-2xl font-bold text-[#1f1f1f]">{title}</h2>
    <p class="text-xs text-[#444746] leading-relaxed max-w-2xl">{description}</p>
  </div>

  <!-- Architecture & Integration Status Card -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- Backend Status -->
    <div class="p-5 bg-white rounded-2xl border border-[#e1e5ea] shadow-xs flex flex-col gap-3">
      <div class="flex items-center gap-2 text-xs font-bold text-[#1f1f1f]">
        <Server class="w-4 h-4 text-[#1e8e3e]" />
        <span>Status Integrasi Go Backend</span>
      </div>
      <p class="text-xs text-[#444746] leading-relaxed">
        {backendStatus}
      </p>
      <div class="p-3 bg-[#f8fafd] rounded-xl border border-[#e1e5ea] font-mono text-[11px] text-[#1f1f1f]">
        internal/{domainPackage}/...
      </div>
    </div>

    <!-- Security & Clean Architecture -->
    <div class="p-5 bg-white rounded-2xl border border-[#e1e5ea] shadow-xs flex flex-col gap-3">
      <div class="flex items-center gap-2 text-xs font-bold text-[#1f1f1f]">
        <ShieldCheck class="w-4 h-4 text-[#0b57d0]" />
        <span>Pragmatic Clean Architecture</span>
      </div>
      <p class="text-xs text-[#444746] leading-relaxed">
        Modul ini mengikuti panduan <code class="px-1.5 py-0.5 rounded bg-[#f0f4f9] text-[#0b57d0] font-mono text-[11px]">AGENTS.md</code>: boundary domain mandiri, otorisasi RBAC berbasis context, dan isolasi transaksi atomik.
      </p>
    </div>
  </div>

  <!-- Planned Endpoints & Workflows -->
  {#if plannedEndpoints.length > 0}
    <div class="p-5 bg-white rounded-2xl border border-[#e1e5ea] shadow-xs flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2 text-xs font-bold text-[#1f1f1f]">
          <Terminal class="w-4 h-4 text-[#0b57d0]" />
          <span>Alur Transaksi & API Rencana</span>
        </div>
        <span class="text-[11px] text-[#747775]">Jalur B (Business Workflow)</span>
      </div>

      <div class="flex flex-col gap-2">
        {#each plannedEndpoints as ep}
          <div class="flex items-center gap-2.5 text-xs text-[#444746] p-2.5 rounded-xl bg-[#f8fafd] border border-[#e1e5ea]">
            <CheckCircle2 class="w-4 h-4 text-[#1e8e3e] shrink-0" />
            <span class="font-mono text-[11px] font-semibold text-[#1f1f1f]">{ep}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <div class="flex justify-end pt-2">
    <M3Button variant="filled" onclick={onBackToLauncher}>
      <ArrowLeft class="w-4 h-4" />
      <span>Kembali ke Desktop Modul</span>
    </M3Button>
  </div>
</div>
