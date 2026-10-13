<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Users,
    ShieldCheck,
    Key,
    Plus,
    RefreshCw,
    CheckCircle2,
    AlertCircle
  } from '@lucide/svelte';
  import M3Button from '../../components/m3/M3Button.svelte';
  import {
    getUsers,
    getRoles,
    getPermissions
  } from '$lib/api/auth';
  import { getPractitioners } from '$lib/api/practitioners';
  import UsersTab from './components/UsersTab.svelte';
  import RolesTab from './components/RolesTab.svelte';
  import PermissionsTab from './components/PermissionsTab.svelte';
  import UserFormModal from './components/UserFormModal.svelte';
  import RoleFormModal from './components/RoleFormModal.svelte';
  import UserDeleteModal from './components/UserDeleteModal.svelte';
  import RoleDeleteModal from './components/RoleDeleteModal.svelte';
  import LinkPractitionerModal from './components/LinkPractitionerModal.svelte';
  import type {
    UserRecord,
    RoleRecord,
    PermissionRecord
  } from '../../types/auth';
  import type { PractitionerRecord } from '../../types/master/practitioner';

  interface Props {
    initialTab?: 'users' | 'roles' | 'permissions';
  }

  let { initialTab = 'users' }: Props = $props();

  let activeTab = $state<'users' | 'roles' | 'permissions'>('users');

  $effect(() => {
    if (initialTab) {
      activeTab = initialTab;
    }
  });

  // Global Data State
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  let users = $state<UserRecord[]>([]);
  let roles = $state<RoleRecord[]>([]);
  let permissions = $state<PermissionRecord[]>([]);
  let practitioners = $state<PractitionerRecord[]>([]);

  // Filters State
  let userSearch = $state('');
  let filterRole = $state('');
  let filterStatus = $state('');

  let roleSearch = $state('');
  let permSearch = $state('');
  let selectedPermModule = $state('');

  // Modals State
  let showUserModal = $state(false);
  let userModalMode = $state<'create' | 'edit'>('create');
  let selectedUser = $state<UserRecord | null>(null);

  let showRoleModal = $state(false);
  let roleModalMode = $state<'create' | 'edit'>('create');
  let selectedRole = $state<RoleRecord | null>(null);

  let showDeleteUserModal = $state(false);
  let userToDelete = $state<UserRecord | null>(null);

  let showDeleteRoleModal = $state(false);
  let roleToDelete = $state<RoleRecord | null>(null);

  let showLinkPractitionerModal = $state(false);
  let userToLink = $state<UserRecord | null>(null);

  function showSuccess(msg: string) {
    successMessage = msg;
    setTimeout(() => {
      successMessage = null;
    }, 4000);
  }

  async function loadData() {
    isLoading = true;
    errorMessage = null;
    try {
      const [uRes, rRes, pRes, prRes] = await Promise.all([
        getUsers({
          search: userSearch.trim() || undefined,
          role_id: filterRole || undefined,
          is_active: filterStatus ? filterStatus === 'true' : undefined,
          limit: 100
        }),
        getRoles(),
        getPermissions(),
        getPractitioners({ limit: 100 })
      ]);

      users = uRes.data || [];
      roles = rRes || [];
      permissions = pRes || [];
      practitioners = prRes.data || [];
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat data otentikasi & wewenang';
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  // Modal Open Handlers
  function handleOpenCreateUser() {
    userModalMode = 'create';
    selectedUser = null;
    showUserModal = true;
  }

  function handleOpenEditUser(u: UserRecord) {
    userModalMode = 'edit';
    selectedUser = u;
    showUserModal = true;
  }

  function handleOpenDeleteUser(u: UserRecord) {
    userToDelete = u;
    showDeleteUserModal = true;
  }

  function handleOpenCreateRole() {
    roleModalMode = 'create';
    selectedRole = null;
    showRoleModal = true;
  }

  function handleOpenEditRole(r: RoleRecord) {
    roleModalMode = 'edit';
    selectedRole = r;
    showRoleModal = true;
  }

  function handleOpenDeleteRole(r: RoleRecord) {
    roleToDelete = r;
    showDeleteRoleModal = true;
  }

  function handleOpenLinkPractitioner(u: UserRecord) {
    userToLink = u;
    showLinkPractitionerModal = true;
  }
</script>

<div class="flex flex-col gap-6 w-full max-w-7xl mx-auto pb-12">
  <!-- Header Banner -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-2 border-b border-[#e1e5ea] dark:border-[#353738]">
    <div class="flex items-center gap-3">
      <div class="w-10 h-10 rounded-2xl bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa] flex items-center justify-center shrink-0">
        <ShieldCheck class="w-6 h-6" />
      </div>
      <div>
        <span class="text-[11px] font-bold text-[#0b57d0] dark:text-[#a8c7fa] tracking-wider uppercase">
          Keamanan & Hak Akses (RBAC)
        </span>
        <h1 class="text-xl font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">
          Manajemen Pengguna, Peran & Hak Akses
        </h1>
        <p class="text-xs text-[#747775] dark:text-[#8e918f] mt-0.5">
          Kelola akun pengguna rumah sakit, integrasi nakes/paramedis, konfigurasi peran, serta matriks perizinan modul.
        </p>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex items-center gap-2">
      <M3Button variant="outlined" onclick={loadData} disabled={isLoading}>
        <RefreshCw class="w-4 h-4 mr-1.5 {isLoading ? 'animate-spin' : ''}" />
        <span>Sinkronisasi</span>
      </M3Button>

      {#if activeTab === 'users'}
        <M3Button variant="filled" onclick={handleOpenCreateUser}>
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Tambah Pengguna</span>
        </M3Button>
      {:else if activeTab === 'roles'}
        <M3Button variant="filled" onclick={handleOpenCreateRole}>
          <Plus class="w-4 h-4 mr-1.5" />
          <span>Tambah Peran</span>
        </M3Button>
      {/if}
    </div>
  </div>

  <!-- Toast Alerts -->
  {#if successMessage}
    <div class="p-3.5 rounded-2xl bg-[#e6f4ea] dark:bg-emerald-950/40 border border-[#ceead6] dark:border-emerald-800 text-[#137333] dark:text-emerald-300 flex items-center gap-2 text-xs font-semibold animate-in fade-in duration-200">
      <CheckCircle2 class="w-4 h-4 shrink-0" />
      <span>{successMessage}</span>
    </div>
  {/if}

  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-center gap-2 text-xs font-semibold animate-in fade-in duration-200">
      <AlertCircle class="w-4 h-4 shrink-0" />
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Navigation Tabs -->
  <div class="flex items-center gap-6 border-b border-[#e1e5ea] dark:border-[#353738]">
    <button
      type="button"
      onclick={() => (activeTab = 'users')}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'users' ? 'text-[#0b57d0] dark:text-[#a8c7fa]' : 'text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3]'}"
    >
      <Users class="w-4 h-4" />
      <span>Daftar Pengguna ({users.length})</span>
      {#if activeTab === 'users'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] dark:bg-[#a8c7fa] rounded-full"></div>
      {/if}
    </button>

    <button
      type="button"
      onclick={() => (activeTab = 'roles')}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'roles' ? 'text-[#0b57d0] dark:text-[#a8c7fa]' : 'text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3]'}"
    >
      <ShieldCheck class="w-4 h-4" />
      <span>Peran & Hak Akses ({roles.length})</span>
      {#if activeTab === 'roles'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] dark:bg-[#a8c7fa] rounded-full"></div>
      {/if}
    </button>

    <button
      type="button"
      onclick={() => (activeTab = 'permissions')}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'permissions' ? 'text-[#0b57d0] dark:text-[#a8c7fa]' : 'text-[#444746] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3]'}"
    >
      <Key class="w-4 h-4" />
      <span>Katalog Izin ({permissions.length})</span>
      {#if activeTab === 'permissions'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] dark:bg-[#a8c7fa] rounded-full"></div>
      {/if}
    </button>
  </div>

  <!-- Tab Contents -->
  {#if activeTab === 'users'}
    <UsersTab
      {users}
      {roles}
      {practitioners}
      {isLoading}
      {userSearch}
      {filterRole}
      {filterStatus}
      onSearchChange={(val) => { userSearch = val; loadData(); }}
      onFilterRoleChange={(val) => { filterRole = val; loadData(); }}
      onFilterStatusChange={(val) => { filterStatus = val; loadData(); }}
      onOpenEditUser={handleOpenEditUser}
      onOpenDeleteUser={handleOpenDeleteUser}
      onOpenLinkPractitioner={handleOpenLinkPractitioner}
    />
  {:else if activeTab === 'roles'}
    <RolesTab
      {roles}
      {isLoading}
      {roleSearch}
      onSearchChange={(val) => (roleSearch = val)}
      onOpenEditRole={handleOpenEditRole}
      onOpenDeleteRole={handleOpenDeleteRole}
    />
  {:else if activeTab === 'permissions'}
    <PermissionsTab
      {permissions}
      {isLoading}
      {permSearch}
      {selectedPermModule}
      onSearchChange={(val) => (permSearch = val)}
      onModuleChange={(val) => (selectedPermModule = val)}
    />
  {/if}
</div>

<!-- ========================================================
     MODALS
     ======================================================== -->

<!-- 1. Form Create / Edit User -->
<UserFormModal
  isOpen={showUserModal}
  mode={userModalMode}
  user={selectedUser}
  {roles}
  {practitioners}
  onClose={() => { showUserModal = false; selectedUser = null; }}
  onSuccess={(saved) => {
    showSuccess(`Data pengguna '${saved.username}' berhasil disimpan!`);
    loadData();
  }}
/>

<!-- 2. Form Create / Edit Role -->
<RoleFormModal
  isOpen={showRoleModal}
  mode={roleModalMode}
  role={selectedRole}
  {permissions}
  onClose={() => { showRoleModal = false; selectedRole = null; }}
  onSuccess={(saved) => {
    showSuccess(`Peran '${saved.name}' berhasil disimpan!`);
    loadData();
  }}
/>

<!-- 3. Konfirmasi Hapus User -->
<UserDeleteModal
  isOpen={showDeleteUserModal}
  user={userToDelete}
  onClose={() => { showDeleteUserModal = false; userToDelete = null; }}
  onSuccess={() => {
    showSuccess(`Pengguna '${userToDelete?.username}' berhasil dinonaktifkan!`);
    loadData();
  }}
/>

<!-- 4. Konfirmasi Hapus Role -->
<RoleDeleteModal
  isOpen={showDeleteRoleModal}
  role={roleToDelete}
  onClose={() => { showDeleteRoleModal = false; roleToDelete = null; }}
  onSuccess={() => {
    showSuccess(`Peran '${roleToDelete?.name}' berhasil dihapus!`);
    loadData();
  }}
/>

<!-- 5. Tautkan Tenaga Medis (Practitioner) Modal Cepat -->
<LinkPractitionerModal
  isOpen={showLinkPractitionerModal}
  user={userToLink}
  {practitioners}
  onClose={() => { showLinkPractitionerModal = false; userToLink = null; }}
  onSuccess={(updated) => {
    showSuccess(`Tautan tenaga medis untuk '${updated.username}' berhasil diperbarui!`);
    loadData();
  }}
/>
