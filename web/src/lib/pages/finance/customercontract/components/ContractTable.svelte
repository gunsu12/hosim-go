<script lang="ts">
  import {
    FileText,
    Building2,
    Calendar,
    MoreVertical,
    Eye,
    Edit3,
    FileSignature,
    RotateCw,
    PauseCircle,
    PlayCircle,
    XCircle,
    Trash2,
    Layers,
  } from '@lucide/svelte';
  import type { CustomerContractRecord } from '$lib/types/finance/customer_contract';

  interface Props {
    contracts: CustomerContractRecord[];
    isLoading?: boolean;
    onViewDetail: (contract: CustomerContractRecord) => void;
    onEdit: (contract: CustomerContractRecord) => void;
    onAmend: (contract: CustomerContractRecord) => void;
    onManageSchemas: (contract: CustomerContractRecord) => void;
    onStatusTransition: (
      contract: CustomerContractRecord,
      action: 'ACTIVATE' | 'SUSPEND' | 'REACTIVATE' | 'TERMINATE',
    ) => void;
    onRenew: (contract: CustomerContractRecord) => void;
    onDelete: (contract: CustomerContractRecord) => void;
  }

  let {
    contracts = [],
    isLoading = false,
    onViewDetail,
    onEdit,
    onAmend,
    onManageSchemas,
    onStatusTransition,
    onRenew,
    onDelete,
  }: Props = $props();

  let activeMenuId = $state<string | null>(null);
  let menuStyle = $state<string>('');

  let activeContract = $derived(contracts.find((c) => c.id === activeMenuId));

  function toggleMenu(id: string, e: MouseEvent) {
    e.stopPropagation();
    if (activeMenuId === id) {
      activeMenuId = null;
      return;
    }

    const button = e.currentTarget as HTMLElement;
    const rect = button.getBoundingClientRect();

    const estimatedHeight = 240;
    const gap = 4;
    const spaceBelow = window.innerHeight - rect.bottom;
    const spaceAbove = rect.top;

    let right = window.innerWidth - rect.right;
    if (right < 8) right = 8;

    // Open upward if space below is limited and space above has more room
    if (spaceBelow < estimatedHeight && spaceAbove > spaceBelow) {
      const bottom = window.innerHeight - rect.top + gap;
      menuStyle = `bottom: ${Math.max(8, bottom)}px; right: ${right}px;`;
    } else {
      const top = rect.bottom + gap;
      menuStyle = `top: ${Math.max(8, top)}px; right: ${right}px;`;
    }

    activeMenuId = id;
  }

  function closeMenu() {
    activeMenuId = null;
  }

  function getExpirationBadge(contract: CustomerContractRecord) {
    if (!contract.expiration_date) return null;
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const exp = new Date(contract.expiration_date);
    exp.setHours(0, 0, 0, 0);
    const diffDays = Math.ceil((exp.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));

    if (diffDays < 0) {
      return { label: 'Kadaluarsa', color: 'bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300' };
    } else if (diffDays <= 30) {
      return { label: `Sisa ${diffDays} hr`, color: 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300' };
    }
    return null;
  }
</script>

<svelte:window
  onclick={closeMenu}
  onscroll={closeMenu}
  onresize={closeMenu}
  onkeydown={(e) => {
    if (e.key === 'Escape') closeMenu();
  }}
/>

<div
  class="w-full overflow-x-auto rounded-2xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] shadow-2xs"
  onscroll={closeMenu}
>
  <table class="w-full text-left border-collapse text-xs">
    <thead>
      <tr class="border-b border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5]">
        <th class="py-3 px-4 font-bold">KONTRAK (PKS)</th>
        <th class="py-3 px-4 font-bold">PENJAMIN / DEBITUR</th>
        <th class="py-3 px-4 font-bold">MASA BERLAKU</th>
        <th class="py-3 px-4 font-bold">TERMIN & LIMIT</th>
        <th class="py-3 px-4 font-bold text-center">STATUS</th>
        <th class="py-3 px-4 font-bold text-right">AKSI</th>
      </tr>
    </thead>
    <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
      {#if isLoading}
        {#each Array(4) as _}
          <tr class="animate-pulse">
            <td class="py-3.5 px-4"><div class="h-4 bg-gray-200 dark:bg-gray-800 rounded w-48"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-gray-200 dark:bg-gray-800 rounded w-36"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-gray-200 dark:bg-gray-800 rounded w-32"></div></td>
            <td class="py-3.5 px-4"><div class="h-4 bg-gray-200 dark:bg-gray-800 rounded w-28"></div></td>
            <td class="py-3.5 px-4 text-center"><div class="h-5 bg-gray-200 dark:bg-gray-800 rounded-full w-16 mx-auto"></div></td>
            <td class="py-3.5 px-4 text-right"><div class="h-6 bg-gray-200 dark:bg-gray-800 rounded w-8 ml-auto"></div></td>
          </tr>
        {/each}
      {:else if contracts.length === 0}
        <tr>
          <td colspan="6" class="py-12 px-4 text-center text-[#747775] dark:text-[#8e918f]">
            <FileText class="w-9 h-9 mx-auto mb-2 opacity-30 text-[#0b57d0]" />
            <p class="font-bold text-sm text-[#1f1f1f] dark:text-[#e3e3e3]">Tidak ada data kontrak kerjasama</p>
            <p class="text-[11px] mt-0.5">Silakan buat kontrak baru atau sesuaikan filter pencarian.</p>
          </td>
        </tr>
      {:else}
        {#each contracts as contract (contract.id)}
          {@const expBadge = getExpirationBadge(contract)}
          <tr
            class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors group cursor-pointer"
            onclick={() => onViewDetail(contract)}
          >
            <!-- 1. KONTRAK -->
            <td class="py-3.5 px-4 min-w-[220px]">
              <div class="flex items-center gap-1.5 flex-wrap">
                <span class="font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa] hover:underline">
                  {contract.contract_number}
                </span>
                {#if contract.renewed_from_contract_number}
                  <span class="px-1.5 py-0.2 rounded text-[9px] font-semibold bg-blue-100 text-blue-800 dark:bg-blue-950/60 dark:text-blue-300">
                    Renew
                  </span>
                {/if}
              </div>
              <p class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate max-w-sm mt-0.5">
                {contract.name}
              </p>
            </td>

            <!-- 2. PENJAMIN -->
            <td class="py-3.5 px-4 min-w-[180px]">
              <div class="flex items-center gap-1.5">
                <Building2 class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f] shrink-0" />
                <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3] truncate max-w-[220px]">
                  {contract.customer_name || contract.customer?.name || 'Penjamin'}
                </span>
              </div>
              {#if contract.customer_code}
                <span class="font-mono text-[10px] text-[#747775] dark:text-[#8e918f] block mt-0.5">
                  Kode: {contract.customer_code}
                </span>
              {/if}
            </td>

            <!-- 3. MASA BERLAKU -->
            <td class="py-3.5 px-4 whitespace-nowrap min-w-[160px]">
              <span class="font-mono text-[11px] text-[#1f1f1f] dark:text-[#e3e3e3] block">
                {contract.effective_date} <span class="text-[#747775] dark:text-[#8e918f]">s/d</span> {contract.expiration_date || 'Permanen'}
              </span>
              {#if expBadge}
                <span class="inline-block px-2 py-0.2 rounded-full text-[9px] font-bold mt-1 {expBadge.color}">
                  {expBadge.label}
                </span>
              {/if}
            </td>

            <!-- 4. TERMIN & LIMIT -->
            <td class="py-3.5 px-4 whitespace-nowrap min-w-[140px]">
              <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] block">
                Net {contract.payment_term_days} Hari
              </span>
              {#if contract.credit_limit !== null && contract.credit_limit !== undefined}
                <span class="text-[11px] font-mono text-[#0b57d0] dark:text-[#a8c7fa] block mt-0.5">
                  Rp {contract.credit_limit.toLocaleString('id-ID')}
                </span>
              {:else}
                <span class="text-[10px] text-emerald-600 dark:text-emerald-400 font-medium block mt-0.5">
                  Bebas Plafon
                </span>
              {/if}
            </td>

            <!-- 5. STATUS -->
            <td class="py-3.5 px-4 whitespace-nowrap text-center">
              <span
                class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase inline-block {contract.status ===
                'ACTIVE'
                  ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300'
                  : contract.status === 'DRAFT'
                    ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300'
                    : contract.status === 'SUSPENDED'
                      ? 'bg-orange-100 text-orange-800 dark:bg-orange-950/60 dark:text-orange-300'
                      : 'bg-rose-100 text-rose-800 dark:bg-rose-950/60 dark:text-rose-300'}"
              >
                {contract.status}
              </span>
            </td>

            <!-- 6. AKSI -->
            <td class="py-3.5 px-4 text-right whitespace-nowrap" onclick={(e) => e.stopPropagation()}>
              <div class="flex items-center justify-end gap-1">
                <button
                  type="button"
                  title="Lihat Detail"
                  onclick={() => onViewDetail(contract)}
                  class="w-7 h-7 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-[#444746] dark:text-[#c4c7c5] flex items-center justify-center cursor-pointer"
                >
                  <Eye class="w-3.5 h-3.5" />
                </button>

                <button
                  type="button"
                  title="Menu Opsi"
                  onclick={(e) => toggleMenu(contract.id, e)}
                  class="w-7 h-7 rounded-full {activeMenuId === contract.id ? 'bg-[#e1e5ea] dark:bg-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]' : 'hover:bg-gray-100 dark:hover:bg-gray-800 text-[#444746] dark:text-[#c4c7c5]'} flex items-center justify-center cursor-pointer transition-colors"
                >
                  <MoreVertical class="w-3.5 h-3.5" />
                </button>
              </div>
            </td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</div>

<!-- Floating Dropdown popup (rendered outside overflow container) -->
{#if activeContract}
  <div
    role="menu"
    tabindex="-1"
    class="fixed z-50 w-52 max-h-[calc(100vh-2rem)] overflow-y-auto rounded-2xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] shadow-xl p-1 text-left text-xs animate-in fade-in zoom-in-95 duration-100"
    style={menuStyle}
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.stopPropagation()}
  >
    <button
      type="button"
      onclick={() => { const c = activeContract; closeMenu(); if (c) onViewDetail(c); }}
      class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
    >
      <Eye class="w-3.5 h-3.5 text-[#0b57d0]" />
      <span>Detail Kontrak</span>
    </button>

    <button
      type="button"
      onclick={() => { const c = activeContract; closeMenu(); if (c) onManageSchemas(c); }}
      class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
    >
      <Layers class="w-3.5 h-3.5 text-indigo-600" />
      <span>Skema Tarif</span>
    </button>

    {#if activeContract.status === 'DRAFT'}
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onEdit(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
      >
        <Edit3 class="w-3.5 h-3.5 text-emerald-600" />
        <span>Edit Draft</span>
      </button>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onStatusTransition(c, 'ACTIVATE'); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-emerald-700 dark:text-emerald-400 font-semibold cursor-pointer"
      >
        <PlayCircle class="w-3.5 h-3.5 text-emerald-600" />
        <span>Aktivasi</span>
      </button>
      <div class="my-1 border-t border-[#e1e5ea] dark:border-[#353738]"></div>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onDelete(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-rose-50 dark:hover:bg-rose-950/40 text-rose-600 dark:text-rose-400 flex items-center gap-2 cursor-pointer"
      >
        <Trash2 class="w-3.5 h-3.5" />
        <span>Hapus Draft</span>
      </button>

    {:else if activeContract.status === 'ACTIVE'}
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onAmend(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
      >
        <FileSignature class="w-3.5 h-3.5 text-blue-600" />
        <span>Amandemen</span>
      </button>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onRenew(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
      >
        <RotateCw class="w-3.5 h-3.5 text-teal-600" />
        <span>Perpanjang (Renew)</span>
      </button>
      <div class="my-1 border-t border-[#e1e5ea] dark:border-[#353738]"></div>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onStatusTransition(c, 'SUSPEND'); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-amber-50 dark:hover:bg-amber-950/40 text-amber-700 dark:text-amber-400 flex items-center gap-2 cursor-pointer"
      >
        <PauseCircle class="w-3.5 h-3.5" />
        <span>Suspend</span>
      </button>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onStatusTransition(c, 'TERMINATE'); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-rose-50 dark:hover:bg-rose-950/40 text-rose-600 dark:text-rose-400 flex items-center gap-2 cursor-pointer"
      >
        <XCircle class="w-3.5 h-3.5" />
        <span>Terminate</span>
      </button>

    {:else if activeContract.status === 'SUSPENDED'}
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onStatusTransition(c, 'REACTIVATE'); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-emerald-50 dark:hover:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 font-semibold flex items-center gap-2 cursor-pointer"
      >
        <PlayCircle class="w-3.5 h-3.5" />
        <span>Aktifkan Kembali</span>
      </button>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onAmend(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
      >
        <FileSignature class="w-3.5 h-3.5 text-blue-600" />
        <span>Amandemen</span>
      </button>
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onStatusTransition(c, 'TERMINATE'); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-rose-50 dark:hover:bg-rose-950/40 text-rose-600 dark:text-rose-400 flex items-center gap-2 cursor-pointer"
      >
        <XCircle class="w-3.5 h-3.5" />
        <span>Terminate</span>
      </button>

    {:else if activeContract.status === 'TERMINATED'}
      <button
        type="button"
        onclick={() => { const c = activeContract; closeMenu(); if (c) onRenew(c); }}
        class="w-full px-3 py-1.5 rounded-xl text-left hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] flex items-center gap-2 text-[#1f1f1f] dark:text-[#e3e3e3] cursor-pointer"
      >
        <RotateCw class="w-3.5 h-3.5 text-teal-600" />
        <span>Perpanjang (Renew)</span>
      </button>
    {/if}
  </div>
{/if}
