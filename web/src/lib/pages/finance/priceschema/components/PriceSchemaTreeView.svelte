<script lang="ts">
  import {
    ChevronRight,
    ChevronDown,
    Building2,
    Hospital,
    Activity,
    PlusCircle,
    Edit2,
    RotateCcw,
    Layers,
    Plus,
    Trash2
  } from '@lucide/svelte';
  import type { PriceSchemaTreeNode } from '$lib/types/finance/price_schema';
  import { formatAdjustmentText } from '../treeUtils';

  interface Props {
    treeNodes: PriceSchemaTreeNode[];
    expandedNodeIds: Set<string>;
    onToggleNode: (id: string) => void;
    onCreateOverride: (node: PriceSchemaTreeNode) => void;
    onEditRule: (node: PriceSchemaTreeNode) => void;
    onResetRule: (node: PriceSchemaTreeNode) => void;
    onAddSubItem: (node: PriceSchemaTreeNode) => void;
  }

  let {
    treeNodes,
    expandedNodeIds,
    onToggleNode,
    onCreateOverride,
    onEditRule,
    onResetRule,
    onAddSubItem
  }: Props = $props();

  function getScopeBadge(node: PriceSchemaTreeNode) {
    if (node.isInherited) {
      return {
        label: 'INHERITED',
        classes: 'bg-zinc-100 text-zinc-600 border-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:border-zinc-700'
      };
    }

    switch (node.scopeType) {
      case 'DEPT':
        return {
          label: 'DEPT',
          classes: 'bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/50 dark:text-purple-300 dark:border-purple-800/60'
        };
      case 'UNIT':
        return {
          label: 'UNIT (OVERRIDE)',
          classes: 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/50 dark:text-blue-300 dark:border-blue-800/60'
        };
      case 'ITEM_UNIT':
        return {
          label: 'ITEM UNIT',
          classes: 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/50 dark:text-amber-300 dark:border-amber-800/60'
        };
      case 'ITEM_DEPT':
        return {
          label: 'ITEM DEPT',
          classes: 'bg-orange-50 text-orange-700 border-orange-200 dark:bg-orange-950/50 dark:text-orange-300 dark:border-orange-800/60'
        };
      case 'ITEM':
        return {
          label: 'ITEM GLOBAL',
          classes: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/50 dark:text-emerald-300 dark:border-emerald-800/60'
        };
      default:
        return {
          label: node.scopeType,
          classes: 'bg-surface-container text-on-surface border-outline-variant'
        };
    }
  }
</script>

<div class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-2xs overflow-hidden">
  <!-- Table Header -->
  <div class="grid grid-cols-12 gap-3 px-5 py-3.5 bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 text-xs font-semibold text-on-surface-variant">
    <div class="col-span-5 flex items-center gap-2">
      <span>Struktur Layanan & Tindakan</span>
    </div>
    <div class="col-span-2 text-center">
      <span>Level Scope</span>
    </div>
    <div class="col-span-2 text-left">
      <span>Penyesuaian Tarif</span>
    </div>
    <div class="col-span-1 text-center">
      <span>Tanggungan</span>
    </div>
    <div class="col-span-1 text-center">
      <span>PPN</span>
    </div>
    <div class="col-span-1 text-right">
      <span>Aksi</span>
    </div>
  </div>

  <!-- Recursive Snippet Tree Renderer -->
  {#snippet renderTreeNode(node: PriceSchemaTreeNode, depth: number)}
    {@const hasChildren = node.children && node.children.length > 0}
    {@const isExpanded = expandedNodeIds.has(node.id)}
    {@const badge = getScopeBadge(node)}
    {@const isDept = node.nodeType === 'DEPT'}
    {@const isUnit = node.nodeType === 'UNIT'}
    {@const isItem = node.nodeType === 'ITEM'}

    <div
      class="grid grid-cols-12 gap-3 px-5 py-3 items-center border-b border-outline-variant/40 hover:bg-[#f8fafd] dark:hover:bg-surface-container-high/60 transition-colors {isDept ? 'bg-surface-container-lowest font-medium' : ''}"
    >
      <!-- Column 1: Indentation, Chevron, Icon, Code & Name -->
      <div class="col-span-5 flex items-center min-w-0" style="padding-left: {depth * 20}px;">
        {#if hasChildren}
          <button
            type="button"
            onclick={() => onToggleNode(node.id)}
            class="p-1 -ml-1 mr-1.5 text-on-surface-variant hover:text-primary hover:bg-primary/10 rounded-lg transition-colors cursor-pointer shrink-0"
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

        <!-- Icon according to Node Type -->
        <span class="mr-2 shrink-0">
          {#if isDept}
            <Building2 class="w-4 h-4 text-purple-600 dark:text-purple-400" />
          {:else if isUnit}
            <Hospital class="w-4 h-4 text-blue-600 dark:text-blue-400" />
          {:else}
            <Activity class="w-4 h-4 text-amber-600 dark:text-amber-400" />
          {/if}
        </span>

        <div class="flex items-baseline gap-2 min-w-0 truncate">
          {#if node.code}
            <span class="font-mono text-[11px] font-semibold text-on-surface-variant shrink-0">
              {node.code}
            </span>
          {/if}
          <span
            class="text-xs truncate {isDept ? 'font-bold text-on-surface' : node.isInherited ? 'text-on-surface-variant' : 'font-semibold text-on-surface'}"
            title={node.name}
          >
            {node.name}
          </span>

          {#if node.isInherited && node.parentRuleSummary}
            <span class="hidden xl:inline-block text-[10px] text-outline truncate" title={node.parentRuleSummary}>
              ({node.parentRuleSummary})
            </span>
          {/if}
        </div>
      </div>

      <!-- Column 2: Scope Level Badge -->
      <div class="col-span-2 text-center">
        <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-bold border {badge.classes}">
          {badge.label}
        </span>
      </div>

      <!-- Column 3: Tariff Adjustment -->
      <div class="col-span-2 text-left">
        <span class="text-xs font-semibold {node.tariffAdjustmentType === 'PERCENT' && node.tariffAdjustmentValue < 0 ? 'text-emerald-700 dark:text-emerald-400' : node.tariffAdjustmentType === 'PERCENT' && node.tariffAdjustmentValue > 0 ? 'text-blue-700 dark:text-blue-400' : 'text-on-surface'}">
          {formatAdjustmentText(node.tariffAdjustmentType, node.tariffAdjustmentValue)}
        </span>
      </div>

      <!-- Column 4: Cost Sharing (Pasien / Penjamin) -->
      <div class="col-span-1 text-center font-mono text-[11px] font-medium text-on-surface">
        {node.payByPatientPercent}% / {node.payByCustomerPercent}%
      </div>

      <!-- Column 5: PPN Tax -->
      <div class="col-span-1 text-center text-xs font-semibold text-on-surface">
        {node.taxPercent}%
      </div>

      <!-- Column 6: Actions -->
      <div class="col-span-1 text-right">
        <div class="inline-flex items-center gap-1">
          {#if node.isInherited}
            <!-- Button to create Override -->
            <button
              type="button"
              onclick={() => onCreateOverride(node)}
              class="px-2 py-1 rounded-lg text-[11px] font-semibold bg-blue-50 hover:bg-blue-100 text-blue-700 dark:bg-blue-950/40 dark:hover:bg-blue-900/60 dark:text-blue-300 border border-blue-200 dark:border-blue-800 transition-colors cursor-pointer"
              title="Tetapkan aturan khusus (override) untuk {node.name}"
            >
              + Override
            </button>
          {:else}
            <!-- If has rule, can edit or reset/delete rule -->
            <button
              type="button"
              onclick={() => onEditRule(node)}
              class="p-1.5 rounded-lg text-on-surface-variant hover:text-primary hover:bg-surface-container transition-colors cursor-pointer"
              title="Edit aturan khusus ini"
            >
              <Edit2 class="w-3.5 h-3.5" />
            </button>

            <!-- Reset / Delete override rule (revert to parent inheritance) -->
            {#if node.ruleId}
              <button
                type="button"
                onclick={() => onResetRule(node)}
                class="p-1.5 rounded-lg text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer"
                title="Hapus override (kembalikan mewarisi aturan induk)"
              >
                <RotateCcw class="w-3.5 h-3.5" />
              </button>
            {/if}
          {/if}

          <!-- If Unit, allow adding specific Item Rule -->
          {#if isUnit}
            <button
              type="button"
              onclick={() => onAddSubItem(node)}
              class="p-1.5 rounded-lg text-on-surface-variant hover:text-primary hover:bg-surface-container transition-colors cursor-pointer"
              title="Tambah pengecualian item tindakan pada {node.name}"
            >
              <PlusCircle class="w-3.5 h-3.5 text-primary" />
            </button>
          {/if}
        </div>
      </div>
    </div>

    <!-- Recursive children -->
    {#if hasChildren && isExpanded}
      {#each node.children! as child (child.id)}
        {@render renderTreeNode(child, depth + 1)}
      {/each}
    {/if}
  {/snippet}

  <!-- Render Root List -->
  {#if treeNodes.length === 0}
    <div class="px-5 py-12 text-center text-xs text-on-surface-variant">
      Tidak ada data yang sesuai dengan pencarian atau filter.
    </div>
  {:else}
    {#each treeNodes as rootNode (rootNode.id)}
      {@render renderTreeNode(rootNode, 0)}
    {/each}
  {/if}
</div>
