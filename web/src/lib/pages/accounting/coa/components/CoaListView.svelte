<script lang="ts">
  import {
    FolderTree,
    Landmark,
    PlusCircle,
    Edit2,
    Trash2
  } from '@lucide/svelte';
  import type { AccountRecord } from '$lib/types/accounting/account';
  import { accountTypeConfig } from '../constants';

  interface Props {
    accounts: AccountRecord[];
    totalAccounts: number;
    canCreate: boolean;
    canUpdate: boolean;
    canDelete: boolean;
    onCreateChild: (parent: AccountRecord) => void;
    onEdit: (account: AccountRecord) => void;
    onDelete: (account: AccountRecord) => void;
  }

  let {
    accounts,
    totalAccounts,
    canCreate,
    canUpdate,
    canDelete,
    onCreateChild,
    onEdit,
    onDelete
  }: Props = $props();
</script>

<div class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-2xs overflow-hidden">
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse text-xs">
      <thead>
        <tr class="bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 text-on-surface-variant select-none">
          <th class="py-3 px-4 font-semibold">Kode Akun</th>
          <th class="py-3 px-4 font-semibold">Nama Akun</th>
          <th class="py-3 px-4 font-semibold">Tipe Akun</th>
          <th class="py-3 px-4 font-semibold text-center">Posisi</th>
          <th class="py-3 px-4 font-semibold text-center">Level</th>
          <th class="py-3 px-4 font-semibold">Akun Induk</th>
          <th class="py-3 px-4 font-semibold text-center">Sifat</th>
          <th class="py-3 px-4 font-semibold">Rekening Kas/Bank</th>
          <th class="py-3 px-4 font-semibold text-center">Status</th>
          <th class="py-3 px-4 font-semibold text-right">Aksi</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-outline-variant/40 text-on-surface">
        {#if accounts.length === 0}
          <tr>
            <td colspan="10" class="py-12 text-center text-on-surface-variant">
              <FolderTree class="w-8 h-8 text-outline mx-auto mb-2 opacity-50" />
              <p class="font-medium">Tidak ada data akun ditemukan</p>
            </td>
          </tr>
        {:else}
          {#each accounts as acc (acc.id)}
            {@const typeMeta = accountTypeConfig[acc.type] || accountTypeConfig.ASSET}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-surface-container-high/80 transition-colors {acc.is_postable ? '' : 'bg-[#f0f4f9]/50 dark:bg-surface-container-high/40 font-medium'}">
              <td class="py-3 px-4 font-mono font-bold {acc.is_postable ? 'text-on-surface-variant' : 'text-primary'}">
                {acc.code}
              </td>
              <td class="py-3 px-4">
                <div class="font-semibold text-on-surface">{acc.name}</div>
                {#if acc.description}
                  <div class="text-[11px] text-on-surface-variant truncate max-w-xs">{acc.description}</div>
                {/if}
              </td>
              <td class="py-3 px-4">
                <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border {typeMeta.badgeBg} {typeMeta.badgeText}">
                  {acc.type}
                </span>
              </td>
              <td class="py-3 px-4 text-center">
                {#if acc.position === 'DEBIT'}
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-teal-50 text-teal-700 border border-teal-200 dark:bg-teal-950/60 dark:text-teal-300 dark:border-teal-800/50">
                    DEBIT
                  </span>
                {:else}
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-200 dark:bg-indigo-950/60 dark:text-indigo-300 dark:border-indigo-800/50">
                    KREDIT
                  </span>
                {/if}
              </td>
              <td class="py-3 px-4 text-center font-mono text-[11px] text-on-surface-variant">
                L{acc.account_level}
              </td>
              <td class="py-3 px-4 font-mono text-xs text-on-surface-variant">
                {acc.parent_code || '-'}
              </td>
              <td class="py-3 px-4 text-center">
                {#if acc.is_postable}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800/50">
                    Postable
                  </span>
                {:else}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-surface-container-highest text-on-surface-variant border border-outline-variant/60">
                    Induk
                  </span>
                {/if}
              </td>
              <td class="py-3 px-4">
                {#if acc.is_treasury_account}
                  <div class="flex items-center gap-1.5 text-xs text-teal-800 dark:text-teal-300 font-medium">
                    <Landmark class="w-3.5 h-3.5 text-teal-600 dark:text-teal-400 shrink-0" />
                    <div>
                      <div>{acc.bank_name || 'Kas RS'}</div>
                      {#if acc.bank_account_number}
                        <div class="text-[10px] font-mono text-on-surface-variant">{acc.bank_account_number}</div>
                      {/if}
                    </div>
                  </div>
                {:else}
                  <span class="text-outline">-</span>
                {/if}
              </td>
              <td class="py-3 px-4 text-center">
                {#if acc.is_active}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border dark:border-emerald-800/50">
                    AKTIF
                  </span>
                {:else}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-surface-container-highest text-on-surface-variant dark:border dark:border-outline-variant/60">
                    NON-AKTIF
                  </span>
                {/if}
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  {#if canCreate}
                    <button
                      type="button"
                      onclick={() => onCreateChild(acc)}
                      class="p-1.5 text-emerald-600 hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-950/50 rounded-lg transition-colors cursor-pointer"
                      title="Tambah sub-akun"
                    >
                      <PlusCircle class="w-3.5 h-3.5" />
                    </button>
                  {/if}
                  {#if canUpdate}
                    <button
                      type="button"
                      onclick={() => onEdit(acc)}
                      class="p-1.5 text-primary hover:bg-primary/10 rounded-lg transition-colors cursor-pointer"
                      title="Edit Akun"
                    >
                      <Edit2 class="w-3.5 h-3.5" />
                    </button>
                  {/if}
                  {#if canDelete}
                    <button
                      type="button"
                      onclick={() => onDelete(acc)}
                      class="p-1.5 text-rose-600 hover:bg-rose-50 dark:text-rose-400 dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
                      title="Hapus Akun"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  {/if}
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  <div class="p-3 bg-[#f8fafd] dark:bg-surface-container-high border-t border-outline-variant/60 text-xs text-on-surface-variant flex items-center justify-between">
    <span>Menampilkan {accounts.length} dari total {totalAccounts} akun</span>
    <span class="font-mono text-[11px]">HOSIM-GO Chart of Accounts Engine</span>
  </div>
</div>
