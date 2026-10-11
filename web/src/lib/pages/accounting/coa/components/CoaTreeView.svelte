<script lang="ts">
  import {
    FolderTree,
    ChevronRight,
    ChevronDown,
    Landmark,
    PlusCircle,
    Edit2,
    Trash2
  } from '@lucide/svelte';
  import type { AccountTreeNode, AccountRecord } from '$lib/types/accounting/account';
  import { accountTypeConfig } from '../constants';

  interface Props {
    treeNodes: AccountTreeNode[];
    expandedNodeIds: Set<string>;
    canCreate: boolean;
    canUpdate: boolean;
    canDelete: boolean;
    onToggleNode: (id: string) => void;
    onCreateChild: (parent: AccountRecord) => void;
    onEdit: (account: AccountRecord) => void;
    onDelete: (account: AccountRecord) => void;
  }

  let {
    treeNodes,
    expandedNodeIds,
    canCreate,
    canUpdate,
    canDelete,
    onToggleNode,
    onCreateChild,
    onEdit,
    onDelete
  }: Props = $props();
</script>

<div class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-2xs overflow-hidden">
  <!-- Tree Table Header -->
  <div class="grid grid-cols-12 gap-3 px-4 py-3 bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 text-xs font-semibold text-on-surface-variant">
    <div class="col-span-6 sm:col-span-5 flex items-center gap-2">
      <span>Struktur Akun & Nama Buku Besar</span>
    </div>
    <div class="col-span-2 text-center">
      <span>Tipe & Golongan</span>
    </div>
    <div class="col-span-1 text-center">
      <span>Posisi</span>
    </div>
    <div class="col-span-1 text-center hidden sm:block">
      <span>Level</span>
    </div>
    <div class="col-span-2 text-center">
      <span>Sifat & Bank</span>
    </div>
    <div class="col-span-1 text-right">
      <span>Aksi</span>
    </div>
  </div>

  <!-- Recursive Tree Snippet Definition -->
  {#snippet renderTreeNode(node: AccountTreeNode, depth: number)}
    {@const hasChildren = node.children && node.children.length > 0}
    {@const isExpanded = expandedNodeIds.has(node.id)}
    {@const typeMeta = accountTypeConfig[node.type] || accountTypeConfig.ASSET}
    {@const isHeader = !node.is_postable}

    <div
      class="grid grid-cols-12 gap-3 px-4 py-2.5 items-center border-b border-outline-variant/40 hover:bg-[#f8fafd] dark:hover:bg-surface-container-high/80 transition-colors {isHeader ? 'bg-[#f0f4f9]/50 dark:bg-surface-container-high/40 font-medium' : ''}"
    >
      <!-- Column 1: Tree Indent, Caret, Code & Name -->
      <div class="col-span-6 sm:col-span-5 flex items-center min-w-0" style="padding-left: {depth * 20}px;">
        {#if hasChildren}
          <button
            type="button"
            onclick={() => onToggleNode(node.id)}
            class="p-1 -ml-1 mr-1 text-on-surface-variant hover:text-primary hover:bg-primary/10 rounded-md transition-colors cursor-pointer shrink-0"
            title={isExpanded ? 'Tutup cabang' : 'Buka cabang'}
          >
            {#if isExpanded}
              <ChevronDown class="w-4 h-4 text-primary" />
            {:else}
              <ChevronRight class="w-4 h-4 text-on-surface-variant" />
            {/if}
          </button>
        {:else}
          <span class="w-6 shrink-0 inline-block text-center text-outline text-xs">•</span>
        {/if}

        <div class="flex items-baseline gap-2 min-w-0 truncate">
          <span class="font-mono text-xs font-bold {isHeader ? 'text-primary' : 'text-on-surface-variant'} shrink-0">
            {node.code}
          </span>
          <span class="text-xs truncate {isHeader ? 'font-bold text-on-surface' : 'text-on-surface'}" title={node.name}>
            {node.name}
          </span>
        </div>
      </div>

      <!-- Column 2: Account Type Badge -->
      <div class="col-span-2 text-center">
        <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border {typeMeta.badgeBg} {typeMeta.badgeText}">
          {node.type}
        </span>
      </div>

      <!-- Column 3: Normal Balance Position -->
      <div class="col-span-1 text-center">
        {#if node.position === 'DEBIT'}
          <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-teal-50 text-teal-700 border border-teal-200 dark:bg-teal-950/60 dark:text-teal-300 dark:border-teal-800/50">
            D
          </span>
        {:else}
          <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-200 dark:bg-indigo-950/60 dark:text-indigo-300 dark:border-indigo-800/50">
            K
          </span>
        {/if}
      </div>

      <!-- Column 4: Level -->
      <div class="col-span-1 text-center hidden sm:block">
        <span class="text-[11px] font-mono text-on-surface-variant">L{node.account_level}</span>
      </div>

      <!-- Column 5: Postable / Treasury / Bank Info -->
      <div class="col-span-2 text-center flex flex-wrap items-center justify-center gap-1">
        {#if node.is_postable}
          <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800/50">
            Postable
          </span>
        {:else}
          <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-surface-container-highest text-on-surface-variant border border-outline-variant/60">
            Induk
          </span>
        {/if}

        {#if node.is_treasury_account}
          <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-teal-100 text-teal-800 border border-teal-200 dark:bg-teal-950/60 dark:text-teal-300 dark:border-teal-800/50 flex items-center gap-1" title="{node.bank_name || 'Bank'}: {node.bank_account_number || ''}">
            <Landmark class="w-2.5 h-2.5" />
            <span>Kas/Bank</span>
          </span>
        {/if}

        {#if !node.is_active}
          <span class="px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-100 text-rose-700 border border-rose-200 dark:bg-rose-950/60 dark:text-rose-300 dark:border-rose-800/50">
            Non-Aktif
          </span>
        {/if}
      </div>

      <!-- Column 6: Action Buttons -->
      <div class="col-span-1 text-right">
        <div class="flex items-center justify-end gap-1">
          {#if canCreate}
            <button
              type="button"
              onclick={() => onCreateChild(node)}
              class="p-1 text-emerald-600 hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-950/50 rounded-lg transition-colors cursor-pointer"
              title="Tambah sub-akun anak di bawah {node.code}"
            >
              <PlusCircle class="w-3.5 h-3.5" />
            </button>
          {/if}

          {#if canUpdate}
            <button
              type="button"
              onclick={() => onEdit(node)}
              class="p-1 text-primary hover:bg-primary/10 rounded-lg transition-colors cursor-pointer"
              title="Edit data akun"
            >
              <Edit2 class="w-3.5 h-3.5" />
            </button>
          {/if}

          {#if canDelete}
            <button
              type="button"
              onclick={() => onDelete(node)}
              class="p-1 text-rose-600 hover:bg-rose-50 dark:text-rose-400 dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
              title="Hapus akun"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          {/if}
        </div>
      </div>
    </div>

    <!-- Render Children Recursively if Expanded -->
    {#if hasChildren && isExpanded}
      {#each node.children as child (child.id)}
        {@render renderTreeNode(child, depth + 1)}
      {/each}
    {/if}
  {/snippet}

  <!-- Tree Nodes Body -->
  {#if treeNodes.length === 0}
    <div class="p-12 text-center text-on-surface-variant">
      <FolderTree class="w-10 h-10 text-outline mx-auto mb-2 opacity-50" />
      <p class="font-medium text-sm">Tidak ada akun yang sesuai dengan filter</p>
      <p class="text-xs text-outline mt-0.5">Coba atur ulang pencarian atau pilih kategori akun lainnya.</p>
    </div>
  {:else}
    <div class="divide-y divide-outline-variant/30">
      {#each treeNodes as rootNode (rootNode.id)}
        {@render renderTreeNode(rootNode, 0)}
      {/each}
    </div>
  {/if}
</div>
