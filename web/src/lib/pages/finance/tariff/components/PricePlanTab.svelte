<script lang="ts">
  import { onMount } from 'svelte';
  import {
    getPricePlans,
    createPricePlan,
    updatePricePlan,
    submitPricePlan,
    approvePricePlan,
    activatePricePlan,
    archivePricePlan,
    clonePricePlan
  } from '$lib/api/finance';
  import { confirmDialog } from '$lib/stores/confirm.svelte';
  import type {
    PricePlanRecord,
    CreatePricePlanDTO,
    UpdatePricePlanDTO,
    ClonePricePlanDTO
  } from '$lib/types/finance';

  import PricePlanTable from './PricePlanTable.svelte';
  import PricePlanModal from './PricePlanModal.svelte';
  import PricePlanCloneModal from './PricePlanCloneModal.svelte';
  import PricePlanItemsModal from './PricePlanItemsModal.svelte';

  interface Props {
    showToast: (msg: string, isError?: boolean) => void;
    onCountChange?: (count: number) => void;
  }

  let { showToast, onCountChange }: Props = $props();

  let plans = $state<PricePlanRecord[]>([]);
  let isPlansLoading = $state(false);
  let planSearchQuery = $state('');
  let planFilterStatus = $state<string>('ALL');
  let planCurrentPage = $state(1);
  let planTotalPages = $state(1);
  let planTotalCount = $state(0);

  // Modal states
  let showPlanModal = $state(false);
  let planModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedPlan = $state<PricePlanRecord | null>(null);
  let isPlanSaving = $state(false);

  let showCloneModal = $state(false);
  let sourceClonePlan = $state<PricePlanRecord | null>(null);
  let isPlanCloning = $state(false);

  let showItemsModal = $state(false);
  let itemsTargetPlan = $state<PricePlanRecord | null>(null);

  async function loadPlans() {
    isPlansLoading = true;
    try {
      const res = await getPricePlans({
        page: planCurrentPage,
        limit: 10,
        search: planSearchQuery || undefined,
        status: planFilterStatus !== 'ALL' ? (planFilterStatus as any) : undefined
      });
      plans = Array.isArray(res.data) ? res.data : [];
      planTotalCount = res.meta?.total_items ?? plans.length;
      planTotalPages = res.meta?.total_pages ?? 1;
      onCountChange?.(planTotalCount);
    } catch (err: any) {
      plans = [];
      showToast(err.message || 'Gagal memuat daftar buku tarif', true);
    } finally {
      isPlansLoading = false;
    }
  }

  function openCreatePlan() {
    planModalMode = 'CREATE';
    selectedPlan = null;
    showPlanModal = true;
  }

  function openEditPlan(plan: PricePlanRecord) {
    planModalMode = 'EDIT';
    selectedPlan = plan;
    showPlanModal = true;
  }

  function openClonePlan(plan: PricePlanRecord) {
    sourceClonePlan = plan;
    showCloneModal = true;
  }

  function openManageItems(plan: PricePlanRecord) {
    itemsTargetPlan = plan;
    showItemsModal = true;
  }

  async function handleSavePlan(payload: CreatePricePlanDTO | UpdatePricePlanDTO) {
    isPlanSaving = true;
    try {
      if (planModalMode === 'CREATE') {
        await createPricePlan(payload as CreatePricePlanDTO);
        showToast(`Buku tarif ${payload.name} berhasil dibuat`);
      } else if (selectedPlan) {
        await updatePricePlan(selectedPlan.id, payload as UpdatePricePlanDTO);
        showToast(`Buku tarif ${payload.name} berhasil diperbarui`);
      }
      showPlanModal = false;
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan buku tarif', true);
    } finally {
      isPlanSaving = false;
    }
  }

  async function handleClonePlan(payload: ClonePricePlanDTO) {
    if (!sourceClonePlan) return;
    isPlanCloning = true;
    try {
      await clonePricePlan(sourceClonePlan.id, payload);
      showToast(`Buku tarif ${payload.new_name} berhasil diduplikasi ke draf baru`);
      showCloneModal = false;
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menggandakan buku tarif', true);
    } finally {
      isPlanCloning = false;
    }
  }

  async function handleSubmitPlan(plan: PricePlanRecord) {
    const ok = await confirmDialog({
      title: 'Pengajuan Buku Tarif',
      message: `Ajukan buku tarif "${plan.name}" untuk direview/approval?`,
      confirmText: 'Ajukan',
      variant: 'primary'
    });
    if (!ok) return;

    try {
      await submitPricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} berhasil diajukan (SUBMITTED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengajukan buku tarif', true);
    }
  }

  async function handleApprovePlan(plan: PricePlanRecord) {
    const ok = await confirmDialog({
      title: 'Persetujuan Buku Tarif',
      message: `Setujui buku tarif "${plan.name}"? Pastikan total pecahan komponen biaya sudah seimbang.`,
      confirmText: 'Setujui',
      variant: 'primary'
    });
    if (!ok) return;

    try {
      await approvePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} berhasil disetujui (APPROVED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyetujui buku tarif', true);
    }
  }

  async function handleActivatePlan(plan: PricePlanRecord) {
    const ok = await confirmDialog({
      title: 'Aktivasi Buku Tarif',
      message: `Aktifkan buku tarif "${plan.name}" sebagai acuan tarif transaksi pelayanan kasir?`,
      confirmText: 'Aktifkan',
      variant: 'primary'
    });
    if (!ok) return;

    try {
      await activatePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} sekarang AKTIF (ACTIVE)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengaktifkan buku tarif', true);
    }
  }

  async function handleArchivePlan(plan: PricePlanRecord) {
    const ok = await confirmDialog({
      title: 'Arsipkan Buku Tarif',
      message: `Arsipkan buku tarif "${plan.name}"? Buku tarif ini tidak akan digunakan lagi untuk transaksi baru.`,
      confirmText: 'Arsipkan',
      variant: 'warning'
    });
    if (!ok) return;

    try {
      await archivePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} telah diarsipkan (ARCHIVED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengarsipkan buku tarif', true);
    }
  }

  onMount(() => {
    loadPlans();
  });
</script>

<PricePlanTable
  {plans}
  isLoading={isPlansLoading}
  searchQuery={planSearchQuery}
  filterStatus={planFilterStatus}
  currentPage={planCurrentPage}
  totalPages={planTotalPages}
  totalCount={planTotalCount}
  onSearchChange={(q) => {
    planSearchQuery = q;
    planCurrentPage = 1;
    loadPlans();
  }}
  onStatusChange={(s) => {
    planFilterStatus = s;
    planCurrentPage = 1;
    loadPlans();
  }}
  onPageChange={(p) => {
    planCurrentPage = p;
    loadPlans();
  }}
  onAdd={openCreatePlan}
  onEdit={openEditPlan}
  onClone={openClonePlan}
  onManageItems={openManageItems}
  onSubmitPlan={handleSubmitPlan}
  onApprovePlan={handleApprovePlan}
  onActivatePlan={handleActivatePlan}
  onArchivePlan={handleArchivePlan}
/>

<PricePlanModal
  isOpen={showPlanModal}
  mode={planModalMode}
  plan={selectedPlan}
  isSaving={isPlanSaving}
  onClose={() => (showPlanModal = false)}
  onSave={handleSavePlan}
/>

<PricePlanCloneModal
  isOpen={showCloneModal}
  sourcePlan={sourceClonePlan}
  isCloning={isPlanCloning}
  onClose={() => (showCloneModal = false)}
  onClone={handleClonePlan}
/>

<PricePlanItemsModal
  isOpen={showItemsModal}
  plan={itemsTargetPlan}
  onClose={() => {
    showItemsModal = false;
    loadPlans();
  }}
  {showToast}
/>
