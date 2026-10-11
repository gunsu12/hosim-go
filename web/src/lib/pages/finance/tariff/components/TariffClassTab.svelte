<script lang="ts">
  import { onMount } from 'svelte';
  import {
    getTariffClasses,
    createTariffClass,
    updateTariffClass,
    deleteTariffClass
  } from '$lib/api/finance';
  import type {
    TariffClassRecord,
    CreateTariffClassDTO,
    UpdateTariffClassDTO
  } from '$lib/types/finance';

  import TariffClassTable from './TariffClassTable.svelte';
  import TariffClassModal from './TariffClassModal.svelte';
  import TariffDeleteModal from './TariffDeleteModal.svelte';

  interface Props {
    showToast: (msg: string, isError?: boolean) => void;
    onCountChange?: (count: number) => void;
  }

  let { showToast, onCountChange }: Props = $props();

  let classes = $state<TariffClassRecord[]>([]);
  let isClassesLoading = $state(false);
  let classSearchQuery = $state('');
  let classFilterStatus = $state<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL');
  let classCurrentPage = $state(1);
  let classTotalPages = $state(1);
  let classTotalCount = $state(0);

  let showClassModal = $state(false);
  let classModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedClass = $state<TariffClassRecord | null>(null);
  let isClassSaving = $state(false);

  let showDeleteModal = $state(false);
  let deleteClassId = $state('');
  let deleteClassCode = $state('');
  let deleteClassName = $state('');
  let isDeleting = $state(false);

  async function loadClasses() {
    isClassesLoading = true;
    try {
      const res = await getTariffClasses({
        page: classCurrentPage,
        limit: 10,
        search: classSearchQuery || undefined,
        is_active: classFilterStatus === 'ALL' ? undefined : classFilterStatus === 'ACTIVE'
      });
      classes = res.data;
      classTotalCount = res.meta.total;
      classTotalPages = res.meta.total_pages;
      onCountChange?.(classTotalCount);
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat daftar kelas tarif', true);
    } finally {
      isClassesLoading = false;
    }
  }

  function openCreateClass() {
    classModalMode = 'CREATE';
    selectedClass = null;
    showClassModal = true;
  }

  function openEditClass(tc: TariffClassRecord) {
    classModalMode = 'EDIT';
    selectedClass = tc;
    showClassModal = true;
  }

  function openDeleteClass(tc: TariffClassRecord) {
    deleteClassId = tc.id;
    deleteClassCode = tc.code;
    deleteClassName = tc.name;
    showDeleteModal = true;
  }

  async function handleSaveClass(payload: CreateTariffClassDTO | UpdateTariffClassDTO) {
    isClassSaving = true;
    try {
      if (classModalMode === 'CREATE') {
        await createTariffClass(payload as CreateTariffClassDTO);
        showToast(`Kelas tarif ${payload.name} berhasil ditambahkan`);
      } else if (selectedClass) {
        await updateTariffClass(selectedClass.id, payload as UpdateTariffClassDTO);
        showToast(`Kelas tarif ${payload.name} berhasil diperbarui`);
      }
      showClassModal = false;
      loadClasses();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan kelas tarif', true);
    } finally {
      isClassSaving = false;
    }
  }

  async function handleConfirmDelete() {
    if (!deleteClassId) return;
    isDeleting = true;
    try {
      await deleteTariffClass(deleteClassId);
      showToast(`Kelas tarif ${deleteClassName} berhasil dihapus`);
      showDeleteModal = false;
      loadClasses();
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus kelas tarif', true);
    } finally {
      isDeleting = false;
    }
  }

  onMount(() => {
    loadClasses();
  });
</script>

<TariffClassTable
  {classes}
  isLoading={isClassesLoading}
  searchQuery={classSearchQuery}
  filterStatus={classFilterStatus}
  currentPage={classCurrentPage}
  totalPages={classTotalPages}
  totalCount={classTotalCount}
  onSearchChange={(q) => {
    classSearchQuery = q;
    classCurrentPage = 1;
    loadClasses();
  }}
  onStatusChange={(s) => {
    classFilterStatus = s;
    classCurrentPage = 1;
    loadClasses();
  }}
  onPageChange={(p) => {
    classCurrentPage = p;
    loadClasses();
  }}
  onAdd={openCreateClass}
  onEdit={openEditClass}
  onDelete={openDeleteClass}
/>

<TariffClassModal
  isOpen={showClassModal}
  mode={classModalMode}
  tariffClass={selectedClass}
  isSaving={isClassSaving}
  onClose={() => (showClassModal = false)}
  onSave={handleSaveClass}
/>

<TariffDeleteModal
  isOpen={showDeleteModal}
  title="Hapus Kelas Tarif"
  itemCode={deleteClassCode}
  itemName={deleteClassName}
  {isDeleting}
  onClose={() => (showDeleteModal = false)}
  onConfirm={handleConfirmDelete}
/>
