<script lang="ts">
  import { onMount } from 'svelte';
  import {
    getTariffComponents,
    createTariffComponent,
    updateTariffComponent,
    deleteTariffComponent
  } from '$lib/api/finance';
  import type {
    TariffComponentRecord,
    CreateTariffComponentDTO,
    UpdateTariffComponentDTO
  } from '$lib/types/finance';

  import TariffComponentTable from './TariffComponentTable.svelte';
  import TariffComponentModal from './TariffComponentModal.svelte';
  import TariffDeleteModal from './TariffDeleteModal.svelte';

  interface Props {
    showToast: (msg: string, isError?: boolean) => void;
    onCountChange?: (count: number) => void;
  }

  let { showToast, onCountChange }: Props = $props();

  let components = $state<TariffComponentRecord[]>([]);
  let isComponentsLoading = $state(false);
  let compSearchQuery = $state('');
  let compFilterType = $state('ALL');
  let compFilterStatus = $state<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL');
  let compCurrentPage = $state(1);
  let compTotalPages = $state(1);
  let compTotalCount = $state(0);

  let showCompModal = $state(false);
  let compModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedComp = $state<TariffComponentRecord | null>(null);
  let isCompSaving = $state(false);

  let showDeleteModal = $state(false);
  let deleteCompId = $state('');
  let deleteCompCode = $state('');
  let deleteCompName = $state('');
  let isDeleting = $state(false);

  async function loadComponents() {
    isComponentsLoading = true;
    try {
      const res = await getTariffComponents({
        page: compCurrentPage,
        limit: 10,
        search: compSearchQuery || undefined,
        component_type: compFilterType !== 'ALL' ? compFilterType : undefined,
        is_active: compFilterStatus === 'ALL' ? undefined : compFilterStatus === 'ACTIVE'
      });
      components = res.data;
      compTotalCount = res.meta.total;
      compTotalPages = res.meta.total_pages;
      onCountChange?.(compTotalCount);
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat daftar komponen tarif', true);
    } finally {
      isComponentsLoading = false;
    }
  }

  function openCreateComponent() {
    compModalMode = 'CREATE';
    selectedComp = null;
    showCompModal = true;
  }

  function openEditComponent(tc: TariffComponentRecord) {
    compModalMode = 'EDIT';
    selectedComp = tc;
    showCompModal = true;
  }

  function openDeleteComponent(tc: TariffComponentRecord) {
    deleteCompId = tc.id;
    deleteCompCode = tc.code;
    deleteCompName = tc.name;
    showDeleteModal = true;
  }

  async function handleSaveComponent(payload: CreateTariffComponentDTO | UpdateTariffComponentDTO) {
    isCompSaving = true;
    try {
      if (compModalMode === 'CREATE') {
        await createTariffComponent(payload as CreateTariffComponentDTO);
        showToast(`Komponen tarif ${payload.name} berhasil ditambahkan`);
      } else if (selectedComp) {
        await updateTariffComponent(selectedComp.id, payload as UpdateTariffComponentDTO);
        showToast(`Komponen tarif ${payload.name} berhasil diperbarui`);
      }
      showCompModal = false;
      loadComponents();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan komponen tarif', true);
    } finally {
      isCompSaving = false;
    }
  }

  async function handleConfirmDelete() {
    if (!deleteCompId) return;
    isDeleting = true;
    try {
      await deleteTariffComponent(deleteCompId);
      showToast(`Komponen tarif ${deleteCompName} berhasil dihapus`);
      showDeleteModal = false;
      loadComponents();
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus komponen tarif', true);
    } finally {
      isDeleting = false;
    }
  }

  onMount(() => {
    loadComponents();
  });
</script>

<TariffComponentTable
  {components}
  isLoading={isComponentsLoading}
  searchQuery={compSearchQuery}
  filterType={compFilterType}
  filterStatus={compFilterStatus}
  currentPage={compCurrentPage}
  totalPages={compTotalPages}
  totalCount={compTotalCount}
  onSearchChange={(q) => {
    compSearchQuery = q;
    compCurrentPage = 1;
    loadComponents();
  }}
  onTypeChange={(t) => {
    compFilterType = t;
    compCurrentPage = 1;
    loadComponents();
  }}
  onStatusChange={(s) => {
    compFilterStatus = s;
    compCurrentPage = 1;
    loadComponents();
  }}
  onPageChange={(p) => {
    compCurrentPage = p;
    loadComponents();
  }}
  onAdd={openCreateComponent}
  onEdit={openEditComponent}
  onDelete={openDeleteComponent}
/>

<TariffComponentModal
  isOpen={showCompModal}
  mode={compModalMode}
  tariffComponent={selectedComp}
  isSaving={isCompSaving}
  onClose={() => (showCompModal = false)}
  onSave={handleSaveComponent}
/>

<TariffDeleteModal
  isOpen={showDeleteModal}
  title="Hapus Komponen Tarif"
  itemCode={deleteCompCode}
  itemName={deleteCompName}
  {isDeleting}
  onClose={() => (showDeleteModal = false)}
  onConfirm={handleConfirmDelete}
/>
