<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Users,
    ShieldCheck,
    Key,
    Plus,
    Search,
    Edit2,
    Trash2,
    Check,
    X,
    AlertCircle,
    Loader2,
    RefreshCw,
    Lock,
    Mail,
    User,
    CheckSquare,
    Square,
    Filter,
    Shield,
    CheckCircle2
  } from '@lucide/svelte';
  import M3Button from '../../components/m3/M3Button.svelte';
  import { auth } from '../../stores/auth.svelte';
  import {
    getUsers,
    createUser,
    updateUser,
    deleteUser,
    getRoles,
    createRole,
    updateRole,
    deleteRole,
    getPermissions
  } from '$lib/api/auth';
  import type {
    UserRecord,
    CreateUserDTO,
    UpdateUserDTO,
    RoleRecord,
    CreateRoleDTO,
    UpdateRoleDTO,
    PermissionRecord
  } from '../../types/auth';

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

  // Global State
  let isLoading = $state(false);
  let isSaving = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Data Arrays
  let users = $state<UserRecord[]>([]);
  let roles = $state<RoleRecord[]>([]);
  let permissions = $state<PermissionRecord[]>([]);
  let totalUsers = $state(0);

  // Filters for Users
  let userSearch = $state('');
  let filterRole = $state('');
  let filterStatus = $state<string>('');

  // Filters for Roles & Permissions
  let roleSearch = $state('');
  let permSearch = $state('');
  let selectedPermModule = $state('');

  // User Modal State
  let showUserModal = $state(false);
  let userModalMode = $state<'create' | 'edit'>('create');
  let selectedUserId = $state<string | null>(null);
  let formUserUsername = $state('');
  let formUserEmail = $state('');
  let formUserName = $state('');
  let formUserPassword = $state('');
  let formUserRoleId = $state('');
  let formUserIsActive = $state(true);
  let formUserError = $state<string | null>(null);

  // User Delete State
  let showDeleteUserConfirm = $state(false);
  let userToDelete = $state<UserRecord | null>(null);
  let isDeletingUser = $state(false);

  // Role Modal State
  let showRoleModal = $state(false);
  let roleModalMode = $state<'create' | 'edit'>('create');
  let selectedRoleId = $state<string | null>(null);
  let formRoleCode = $state('');
  let formRoleName = $state('');
  let formRoleDesc = $state('');
  let formRoleSelectedPerms = $state<string[]>([]);
  let formRoleError = $state<string | null>(null);

  // Role Delete State
  let showDeleteRoleConfirm = $state(false);
  let roleToDelete = $state<RoleRecord | null>(null);
  let isDeletingRole = $state(false);

  // Fetch initial data
  async function loadData() {
    isLoading = true;
    errorMessage = null;
    try {
      const [usersData, rolesData, permsData] = await Promise.all([
        getUsers({
          search: userSearch,
          role_id: filterRole || undefined,
          is_active: filterStatus === '' ? undefined : filterStatus === 'true',
          limit: 100
        }),
        getRoles(),
        getPermissions()
      ]);
      users = usersData.data || [];
      totalUsers = usersData.count || 0;
      roles = rolesData || [];
      permissions = permsData || [];
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat data autentikasi';
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  function showSuccess(msg: string) {
    successMessage = msg;
    setTimeout(() => {
      successMessage = null;
    }, 4000);
  }

  // ==========================================
  // USER ACTIONS
  // ==========================================

  function openCreateUserModal() {
    userModalMode = 'create';
    selectedUserId = null;
    formUserUsername = '';
    formUserEmail = '';
    formUserName = '';
    formUserPassword = '';
    formUserRoleId = roles.length > 0 ? roles[0].id : '';
    formUserIsActive = true;
    formUserError = null;
    showUserModal = true;
  }

  function openEditUserModal(user: UserRecord) {
    userModalMode = 'edit';
    selectedUserId = user.id;
    formUserUsername = user.username;
    formUserEmail = user.email;
    formUserName = user.name;
    formUserPassword = ''; // kosongkan jika tidak ingin ganti
    formUserRoleId = user.role_id || (roles.find(r => r.code === user.role)?.id || '');
    formUserIsActive = user.is_active;
    formUserError = null;
    showUserModal = true;
  }

  async function handleSaveUser() {
    formUserError = null;
    if (userModalMode === 'create' && !formUserUsername.trim()) {
      formUserError = 'Username wajib diisi';
      return;
    }
    if (!formUserEmail.trim() || !formUserEmail.includes('@')) {
      formUserError = 'Email tidak valid';
      return;
    }
    if (!formUserName.trim()) {
      formUserError = 'Nama lengkap wajib diisi';
      return;
    }
    if (userModalMode === 'create' && (!formUserPassword || formUserPassword.length < 6)) {
      formUserError = 'Password wajib diisi minimal 6 karakter';
      return;
    }

    isSaving = true;
    try {
      if (userModalMode === 'create') {
        const payload: CreateUserDTO = {
          username: formUserUsername.trim(),
          email: formUserEmail.trim().toLowerCase(),
          password: formUserPassword,
          name: formUserName.trim(),
          role_id: formUserRoleId || undefined,
          is_active: formUserIsActive
        };
        await createUser(payload);
        showSuccess(`Pengguna '${formUserUsername}' berhasil dibuat!`);
      } else if (selectedUserId) {
        const payload: UpdateUserDTO = {
          name: formUserName.trim(),
          email: formUserEmail.trim().toLowerCase(),
          role_id: formUserRoleId || undefined,
          is_active: formUserIsActive,
          password: formUserPassword.trim() ? formUserPassword.trim() : undefined
        };
        await updateUser(selectedUserId, payload);
        showSuccess(`Pengguna '${formUserName}' berhasil diperbarui!`);
      }
      showUserModal = false;
      await loadData();
    } catch (err: any) {
      formUserError = err.message || 'Gagal menyimpan pengguna';
    } finally {
      isSaving = false;
    }
  }

  function confirmDeleteUser(user: UserRecord) {
    userToDelete = user;
    showDeleteUserConfirm = true;
  }

  async function handleDeleteUser() {
    if (!userToDelete) return;
    isDeletingUser = true;
    try {
      await deleteUser(userToDelete.id);
      showSuccess(`Pengguna '${userToDelete.username}' berhasil dinonaktifkan/dihapus!`);
      showDeleteUserConfirm = false;
      userToDelete = null;
      await loadData();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghapus pengguna';
    } finally {
      isDeletingUser = false;
    }
  }

  // ==========================================
  // ROLE ACTIONS
  // ==========================================

  function openCreateRoleModal() {
    roleModalMode = 'create';
    selectedRoleId = null;
    formRoleCode = '';
    formRoleName = '';
    formRoleDesc = '';
    formRoleSelectedPerms = [];
    formRoleError = null;
    showRoleModal = true;
  }

  function openEditRoleModal(role: RoleRecord) {
    roleModalMode = 'edit';
    selectedRoleId = role.id;
    formRoleCode = role.code;
    formRoleName = role.name;
    formRoleDesc = role.description || '';
    formRoleSelectedPerms = role.permissions ? role.permissions.map(p => p.id) : [];
    formRoleError = null;
    showRoleModal = true;
  }

  function togglePermission(permId: string) {
    if (formRoleSelectedPerms.includes(permId)) {
      formRoleSelectedPerms = formRoleSelectedPerms.filter(id => id !== permId);
    } else {
      formRoleSelectedPerms = [...formRoleSelectedPerms, permId];
    }
  }

  function toggleModulePermissions(moduleName: string, permIds: string[]) {
    const allSelected = permIds.every(id => formRoleSelectedPerms.includes(id));
    if (allSelected) {
      formRoleSelectedPerms = formRoleSelectedPerms.filter(id => !permIds.includes(id));
    } else {
      const toAdd = permIds.filter(id => !formRoleSelectedPerms.includes(id));
      formRoleSelectedPerms = [...formRoleSelectedPerms, ...toAdd];
    }
  }

  async function handleSaveRole() {
    formRoleError = null;
    if (roleModalMode === 'create' && !formRoleCode.trim()) {
      formRoleError = 'Kode Role wajib diisi';
      return;
    }
    if (!formRoleName.trim()) {
      formRoleError = 'Nama Role wajib diisi';
      return;
    }

    isSaving = true;
    try {
      if (roleModalMode === 'create') {
        const payload: CreateRoleDTO = {
          code: formRoleCode.trim().toUpperCase(),
          name: formRoleName.trim(),
          description: formRoleDesc.trim(),
          permission_ids: formRoleSelectedPerms
        };
        await createRole(payload);
        showSuccess(`Peran '${formRoleCode}' berhasil dibuat!`);
      } else if (selectedRoleId) {
        const payload: UpdateRoleDTO = {
          name: formRoleName.trim(),
          description: formRoleDesc.trim(),
          permission_ids: formRoleSelectedPerms
        };
        await updateRole(selectedRoleId, payload);
        showSuccess(`Peran '${formRoleName}' berhasil diperbarui!`);
      }
      showRoleModal = false;
      await loadData();
    } catch (err: any) {
      formRoleError = err.message || 'Gagal menyimpan peran';
    } finally {
      isSaving = false;
    }
  }

  function confirmDeleteRole(role: RoleRecord) {
    roleToDelete = role;
    showDeleteRoleConfirm = true;
  }

  async function handleDeleteRole() {
    if (!roleToDelete) return;
    isDeletingRole = true;
    try {
      await deleteRole(roleToDelete.id);
      showSuccess(`Peran '${roleToDelete.code}' berhasil dihapus!`);
      showDeleteRoleConfirm = false;
      roleToDelete = null;
      await loadData();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghapus peran';
    } finally {
      isDeletingRole = false;
    }
  }

  // Permissions grouped by module
  let permissionsByModule = $derived(
    permissions.reduce((acc, p) => {
      const mod = p.module || 'Umum';
      if (!acc[mod]) acc[mod] = [];
      acc[mod].push(p);
      return acc;
    }, {} as Record<string, PermissionRecord[]>)
  );

  let permissionModules = $derived(Object.keys(permissionsByModule).sort());

  let filteredRoles = $derived(
    roles.filter(r => {
      if (!roleSearch.trim()) return true;
      const q = roleSearch.toLowerCase();
      return r.code.toLowerCase().includes(q) || r.name.toLowerCase().includes(q) || (r.description || '').toLowerCase().includes(q);
    })
  );

  let filteredPermissions = $derived(
    permissions.filter(p => {
      if (selectedPermModule && p.module !== selectedPermModule) return false;
      if (!permSearch.trim()) return true;
      const q = permSearch.toLowerCase();
      return p.code.toLowerCase().includes(q) || p.name.toLowerCase().includes(q) || p.module.toLowerCase().includes(q);
    })
  );
</script>

<div class="flex flex-col gap-5 w-full animate-in fade-in duration-150">
  <!-- Module Header Banner -->
  <div class="flex flex-wrap items-center justify-between gap-4 border-b border-[#e1e5ea] pb-4">
    <div>
      <div class="flex items-center gap-2 text-xs font-semibold text-[#0b57d0] tracking-wide uppercase mb-1">
        <ShieldCheck class="w-4 h-4" />
        <span>Keamanan & Hak Akses (RBAC)</span>
      </div>
      <h1 class="text-xl font-bold text-[#1f1f1f] tracking-tight">Manajemen Pengguna, Peran & Hak Akses</h1>
      <p class="text-xs text-[#444746] mt-0.5">Kelola akun pengguna rumah sakit, konfigurasi peran, serta matriks perizinan modul.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="outlined" onclick={loadData} disabled={isLoading}>
        <RefreshCw class="w-4 h-4 {isLoading ? 'animate-spin' : ''}" />
        <span>Sinkronisasi</span>
      </M3Button>
      {#if activeTab === 'users'}
        <M3Button variant="filled" onclick={openCreateUserModal}>
          <Plus class="w-4 h-4" />
          <span>Tambah Pengguna</span>
        </M3Button>
      {:else if activeTab === 'roles'}
        <M3Button variant="filled" onclick={openCreateRoleModal}>
          <Plus class="w-4 h-4" />
          <span>Tambah Peran</span>
        </M3Button>
      {/if}
    </div>
  </div>

  <!-- Notification Alerts -->
  {#if successMessage}
    <div class="p-3.5 rounded-2xl bg-[#e6f4ea] border border-[#ceead6] text-[#137333] flex items-center justify-between shadow-xs animate-in fade-in duration-200">
      <div class="flex items-center gap-2 text-xs font-semibold">
        <CheckCircle2 class="w-4 h-4 text-[#137333]" />
        <span>{successMessage}</span>
      </div>
    </div>
  {/if}

  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-[#fce8e6] border border-[#f5c2c7] text-[#c5221f] flex items-center justify-between shadow-xs animate-in fade-in duration-200">
      <div class="flex items-center gap-2 text-xs font-semibold">
        <AlertCircle class="w-4 h-4 text-[#c5221f]" />
        <span>{errorMessage}</span>
      </div>
      <button type="button" onclick={() => errorMessage = null} class="text-[#c5221f] cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  <!-- Tab Switcher Navigation (Material 3 Segmented Tabs) -->
  <div class="flex border-b border-[#e1e5ea] gap-8">
    <button
      type="button"
      onclick={() => activeTab = 'users'}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'users' ? 'text-[#0b57d0]' : 'text-[#444746] hover:text-[#1f1f1f]'}"
    >
      <Users class="w-4 h-4" />
      <span>Daftar Pengguna ({totalUsers})</span>
      {#if activeTab === 'users'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] rounded-full"></div>
      {/if}
    </button>

    <button
      type="button"
      onclick={() => activeTab = 'roles'}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'roles' ? 'text-[#0b57d0]' : 'text-[#444746] hover:text-[#1f1f1f]'}"
    >
      <Shield class="w-4 h-4" />
      <span>Peran & Hak Akses ({roles.length})</span>
      {#if activeTab === 'roles'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] rounded-full"></div>
      {/if}
    </button>

    <button
      type="button"
      onclick={() => activeTab = 'permissions'}
      class="pb-3 text-xs font-bold transition-all relative flex items-center gap-2 cursor-pointer {activeTab === 'permissions' ? 'text-[#0b57d0]' : 'text-[#444746] hover:text-[#1f1f1f]'}"
    >
      <Key class="w-4 h-4" />
      <span>Katalog Izin ({permissions.length})</span>
      {#if activeTab === 'permissions'}
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-[#0b57d0] rounded-full"></div>
      {/if}
    </button>
  </div>

  <!-- TAB 1: PENGGUNA (USERS) -->
  {#if activeTab === 'users'}
    <!-- Search & Filter Bar -->
    <div class="flex flex-wrap items-center gap-3">
      <div class="relative flex-1 min-w-[240px]">
        <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775]" />
        <input
          type="text"
          bind:value={userSearch}
          oninput={loadData}
          placeholder="Cari username, nama, atau email..."
          class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] focus:bg-white focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
        />
      </div>

      <select
        bind:value={filterRole}
        onchange={loadData}
        class="h-10 px-3 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] text-[#1f1f1f] focus:outline-none cursor-pointer"
      >
        <option value="">Semua Peran</option>
        {#each roles as r}
          <option value={r.id}>{r.name} ({r.code})</option>
        {/each}
      </select>

      <select
        bind:value={filterStatus}
        onchange={loadData}
        class="h-10 px-3 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] text-[#1f1f1f] focus:outline-none cursor-pointer"
      >
        <option value="">Semua Status</option>
        <option value="true">Aktif</option>
        <option value="false">Nonaktif</option>
      </select>
    </div>

    <!-- Users Table -->
    <div class="bg-white rounded-2xl border border-[#e1e5ea] overflow-hidden shadow-xs">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f8fafd] border-b border-[#e1e5ea] text-[#444746] font-semibold">
            <th class="py-3 px-4">Pengguna</th>
            <th class="py-3 px-4">Email</th>
            <th class="py-3 px-4">Peran (Role)</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4">Terdaftar</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea]">
          {#each users as u}
            <tr class="hover:bg-[#f8fafd]/80 transition-colors">
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-[#e8f0fe] text-[#0b57d0] font-bold flex items-center justify-center text-xs">
                    {u.name ? u.name.charAt(0).toUpperCase() : u.username.charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <div class="font-bold text-[#1f1f1f]">{u.name}</div>
                    <div class="text-[11px] text-[#747775] font-mono">@{u.username}</div>
                  </div>
                </div>
              </td>
              <td class="py-3 px-4 text-[#444746] font-mono text-[11px]">{u.email}</td>
              <td class="py-3 px-4">
                <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {u.role === 'ADMIN' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'}">
                  {u.role_name || u.role || 'Staf'}
                </span>
              </td>
              <td class="py-3 px-4">
                {#if u.is_active}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#e6f4ea] text-[#137333]">
                    <span class="w-1.5 h-1.5 rounded-full bg-[#137333]"></span>
                    Aktif
                  </span>
                {:else}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#fce8e6] text-[#c5221f]">
                    <span class="w-1.5 h-1.5 rounded-full bg-[#c5221f]"></span>
                    Nonaktif
                  </span>
                {/if}
              </td>
              <td class="py-3 px-4 text-[#747775] text-[11px]">
                {new Date(u.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    onclick={() => openEditUserModal(u)}
                    class="p-1.5 text-[#444746] hover:text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
                    title="Ubah Data Pengguna"
                  >
                    <Edit2 class="w-4 h-4" />
                  </button>
                  {#if u.username !== 'admin'}
                    <button
                      type="button"
                      onclick={() => confirmDeleteUser(u)}
                      class="p-1.5 text-[#444746] hover:text-[#c5221f] hover:bg-[#fce8e6] rounded-lg transition-colors cursor-pointer"
                      title="Hapus / Nonaktifkan Pengguna"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  {/if}
                </div>
              </td>
            </tr>
          {/each}

          {#if users.length === 0}
            <tr>
              <td colspan="6" class="py-8 text-center text-[#747775]">
                {isLoading ? 'Memuat data pengguna...' : 'Tidak ada pengguna yang cocok dengan pencarian.'}
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>

  <!-- TAB 2: PERAN & HAK AKSES (ROLES) -->
  {:else if activeTab === 'roles'}
    <!-- Search Bar Roles -->
    <div class="flex items-center justify-between gap-3">
      <div class="relative w-72">
        <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775]" />
        <input
          type="text"
          bind:value={roleSearch}
          placeholder="Cari kode atau nama peran..."
          class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] focus:bg-white focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
        />
      </div>
    </div>

    <!-- Roles Grid Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each filteredRoles as r}
        <div class="bg-white rounded-2xl border border-[#e1e5ea] p-5 shadow-xs flex flex-col justify-between gap-4 hover:shadow-md transition-shadow">
          <div class="flex flex-col gap-2">
            <div class="flex items-start justify-between gap-2">
              <span class="font-mono text-xs font-bold px-2 py-0.5 rounded-lg {r.code === 'ADMIN' ? 'bg-purple-100 text-purple-900 border border-purple-200' : 'bg-blue-100 text-blue-900 border border-blue-200'}">
                {r.code}
              </span>
              <span class="text-[11px] font-semibold text-[#137333] bg-[#e6f4ea] px-2 py-0.5 rounded-full">
                {r.permissions ? r.permissions.length : 0} Izin Aktif
              </span>
            </div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">{r.name}</h3>
            <p class="text-xs text-[#444746] leading-relaxed line-clamp-2">
              {r.description || 'Tidak ada keterangan tambahan.'}
            </p>
          </div>

          <div class="pt-3 border-t border-[#e1e5ea] flex items-center justify-between">
            <span class="text-[10px] text-[#747775]">
              {r.code === 'ADMIN' ? 'Superadmin Bypass' : 'Kustom Hak Akses'}
            </span>
            <div class="flex items-center gap-1">
              <button
                type="button"
                onclick={() => openEditRoleModal(r)}
                class="px-2.5 py-1 text-xs font-semibold text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
              >
                Konfigurasi Izin
              </button>
              {#if r.code !== 'ADMIN'}
                <button
                  type="button"
                  onclick={() => confirmDeleteRole(r)}
                  class="p-1 text-[#444746] hover:text-[#c5221f] hover:bg-[#fce8e6] rounded-lg transition-colors cursor-pointer"
                  title="Hapus Peran"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>

  <!-- TAB 3: KATALOG HAK AKSES (PERMISSIONS) -->
  {:else if activeTab === 'permissions'}
    <!-- Search & Filter Modules -->
    <div class="flex flex-wrap items-center gap-3">
      <div class="relative flex-1 min-w-[240px]">
        <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775]" />
        <input
          type="text"
          bind:value={permSearch}
          placeholder="Cari kode atau nama izin hak akses..."
          class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] focus:bg-white focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
        />
      </div>

      <select
        bind:value={selectedPermModule}
        class="h-10 px-3 text-xs rounded-xl bg-[#f8fafd] border border-[#e1e5ea] text-[#1f1f1f] focus:outline-none cursor-pointer"
      >
        <option value="">Semua Modul</option>
        {#each permissionModules as m}
          <option value={m}>{m}</option>
        {/each}
      </select>
    </div>

    <!-- Permissions Table -->
    <div class="bg-white rounded-2xl border border-[#e1e5ea] overflow-hidden shadow-xs">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f8fafd] border-b border-[#e1e5ea] text-[#444746] font-semibold">
            <th class="py-3 px-4">Kode Izin (Code)</th>
            <th class="py-3 px-4">Nama Izin</th>
            <th class="py-3 px-4">Modul Domain</th>
            <th class="py-3 px-4">Keterangan</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea]">
          {#each filteredPermissions as p}
            <tr class="hover:bg-[#f8fafd]/80 transition-colors">
              <td class="py-2.5 px-4 font-mono font-bold text-[#0b57d0]">{p.code}</td>
              <td class="py-2.5 px-4 font-medium text-[#1f1f1f]">{p.name}</td>
              <td class="py-2.5 px-4">
                <span class="px-2 py-0.5 rounded-md bg-[#f0f4f9] text-[#444746] text-[11px]">
                  {p.module}
                </span>
              </td>
              <td class="py-2.5 px-4 text-[#747775]">{p.description || '-'}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<!-- ========================================================
     MODAL CREATE / EDIT PENGGUNA (USER)
     ======================================================== -->
{#if showUserModal}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-xl border border-[#e1e5ea] flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <h3 class="text-base font-bold text-[#1f1f1f]">
          {userModalMode === 'create' ? 'Tambah Pengguna Baru' : 'Ubah Data Pengguna'}
        </h3>
        <button type="button" onclick={() => showUserModal = false} class="p-1 rounded-full text-[#747775] hover:bg-[#f0f4f9]">
          <X class="w-5 h-5" />
        </button>
      </div>

      {#if formUserError}
        <div class="p-3 bg-[#fce8e6] border border-[#f5c2c7] text-[#c5221f] text-xs rounded-xl flex items-center gap-2">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>{formUserError}</span>
        </div>
      {/if}

      <div class="flex flex-col gap-3 text-xs">
        <div>
          <label for="form-username" class="block font-semibold text-[#444746] mb-1">Username</label>
          <input
            id="form-username"
            type="text"
            bind:value={formUserUsername}
            disabled={userModalMode === 'edit'}
            placeholder="misal: dokter.hendra"
            class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] disabled:bg-[#f0f4f9] disabled:text-[#747775] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
          />
        </div>

        <div>
          <label for="form-email" class="block font-semibold text-[#444746] mb-1">Alamat Email</label>
          <input
            id="form-email"
            type="email"
            bind:value={formUserEmail}
            placeholder="misal: hendra@hosim.local"
            class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
          />
        </div>

        <div>
          <label for="form-name" class="block font-semibold text-[#444746] mb-1">Nama Lengkap & Gelar</label>
          <input
            id="form-name"
            type="text"
            bind:value={formUserName}
            placeholder="misal: dr. Hendra Wijaya, Sp.B"
            class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
          />
        </div>

        <div>
          <label for="form-role" class="block font-semibold text-[#444746] mb-1">Peran Pengguna (Role)</label>
          <select
            id="form-role"
            bind:value={formUserRoleId}
            class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none cursor-pointer"
          >
            {#each roles as r}
              <option value={r.id}>{r.name} ({r.code})</option>
            {/each}
          </select>
        </div>

        <div>
          <label for="form-password" class="block font-semibold text-[#444746] mb-1">
            {userModalMode === 'create' ? 'Kata Sandi (Password)' : 'Reset Kata Sandi (Kosongkan bila tidak diubah)'}
          </label>
          <input
            id="form-password"
            type="password"
            bind:value={formUserPassword}
            placeholder={userModalMode === 'create' ? 'Minimal 6 karakter' : 'Tulis password baru jika ingin mereset'}
            class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
          />
        </div>

        <div class="flex items-center gap-2 pt-1">
          <input
            type="checkbox"
            id="user-is-active"
            bind:checked={formUserIsActive}
            class="rounded accent-[#0b57d0] w-4 h-4"
          />
          <label for="user-is-active" class="font-semibold text-[#1f1f1f] cursor-pointer">
            Akun Aktif (Dapat Login ke Sistem)
          </label>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-2 border-t border-[#e1e5ea]">
        <M3Button variant="outlined" onclick={() => showUserModal = false}>
          <span>Batal</span>
        </M3Button>
        <M3Button variant="filled" onclick={handleSaveUser} disabled={isSaving}>
          {#if isSaving}
            <Loader2 class="w-4 h-4 animate-spin" />
          {/if}
          <span>{userModalMode === 'create' ? 'Buat Akun' : 'Simpan Perubahan'}</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}

<!-- ========================================================
     MODAL CREATE / EDIT PERAN & MATRIKS HAK AKSES (ROLE)
     ======================================================== -->
{#if showRoleModal}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl max-w-2xl w-full max-h-[90vh] flex flex-col shadow-xl border border-[#e1e5ea]">
      <div class="flex items-center justify-between border-b border-[#e1e5ea] p-6 pb-4">
        <div>
          <h3 class="text-base font-bold text-[#1f1f1f]">
            {roleModalMode === 'create' ? 'Tambah Peran (Role) Baru' : `Konfigurasi Peran: ${formRoleName}`}
          </h3>
          <p class="text-xs text-[#444746]">Atur hak akses dan wewenang pengguna terhadap modul-modul sistem.</p>
        </div>
        <button type="button" onclick={() => showRoleModal = false} class="p-1 rounded-full text-[#747775] hover:bg-[#f0f4f9]">
          <X class="w-5 h-5" />
        </button>
      </div>

      <div class="p-6 overflow-y-auto flex flex-col gap-4 text-xs">
        {#if formRoleError}
          <div class="p-3 bg-[#fce8e6] border border-[#f5c2c7] text-[#c5221f] text-xs rounded-xl flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{formRoleError}</span>
          </div>
        {/if}

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="form-role-code" class="block font-semibold text-[#444746] mb-1">Kode Peran (Code)</label>
            <input
              id="form-role-code"
              type="text"
              bind:value={formRoleCode}
              disabled={roleModalMode === 'edit'}
              placeholder="misal: PHARMACIST"
              class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] disabled:bg-[#f0f4f9] uppercase font-mono focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
            />
          </div>

          <div>
            <label for="form-role-name" class="block font-semibold text-[#444746] mb-1">Nama Peran</label>
            <input
              id="form-role-name"
              type="text"
              bind:value={formRoleName}
              placeholder="misal: Apoteker / Petugas Farmasi"
              class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
            />
          </div>

          <div class="sm:col-span-2">
            <label for="form-role-desc" class="block font-semibold text-[#444746] mb-1">Deskripsi Wewenang</label>
            <input
              id="form-role-desc"
              type="text"
              bind:value={formRoleDesc}
              placeholder="Keterangan cakupan kerja pengguna..."
              class="w-full h-10 px-3 rounded-xl bg-white border border-[#e1e5ea] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
            />
          </div>
        </div>

        <!-- Matriks Hak Akses (Permissions Checkbox Grouped by Module) -->
        <div class="flex flex-col gap-2 pt-2 border-t border-[#e1e5ea]">
          <div class="flex items-center justify-between">
            <div class="font-bold text-[#1f1f1f]">Matriks Hak Akses Modul ({formRoleSelectedPerms.length} Dipilih)</div>
            <span class="text-[11px] text-[#747775]">Klik modul untuk centang semua</span>
          </div>

          <div class="flex flex-col gap-3 mt-1">
            {#each permissionModules as modName}
              {@const modPerms = permissionsByModule[modName] || []}
              {@const modPermIds = modPerms.map(p => p.id)}
              {@const isAllChecked = modPermIds.length > 0 && modPermIds.every(id => formRoleSelectedPerms.includes(id))}
              {@const isSomeChecked = modPermIds.some(id => formRoleSelectedPerms.includes(id))}

              <div class="p-3 rounded-2xl border border-[#e1e5ea] bg-[#f8fafd]">
                <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-2 mb-2">
                  <div class="flex items-center gap-2">
                    <button
                      type="button"
                      onclick={() => toggleModulePermissions(modName, modPermIds)}
                      class="flex items-center gap-1.5 font-bold text-xs text-[#1f1f1f] hover:text-[#0b57d0] cursor-pointer"
                    >
                      {#if isAllChecked}
                        <CheckSquare class="w-4 h-4 text-[#0b57d0]" />
                      {:else}
                        <Square class="w-4 h-4 text-[#747775]" />
                      {/if}
                      <span>{modName}</span>
                    </button>
                  </div>
                  <span class="text-[10px] text-[#747775]">
                    {modPermIds.filter(id => formRoleSelectedPerms.includes(id)).length}/{modPermIds.length} Izin
                  </span>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {#each modPerms as p}
                    <label class="flex items-center gap-2 p-1.5 rounded-lg hover:bg-white transition-colors cursor-pointer select-none">
                      <input
                        type="checkbox"
                        checked={formRoleSelectedPerms.includes(p.id)}
                        onchange={() => togglePermission(p.id)}
                        class="rounded accent-[#0b57d0] w-3.5 h-3.5"
                      />
                      <div class="flex flex-col">
                        <span class="text-[11px] font-semibold text-[#1f1f1f]">{p.name}</span>
                        <span class="text-[9px] font-mono text-[#747775]">{p.code}</span>
                      </div>
                    </label>
                  {/each}
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2 p-6 pt-3 border-t border-[#e1e5ea]">
        <M3Button variant="outlined" onclick={() => showRoleModal = false}>
          <span>Batal</span>
        </M3Button>
        <M3Button variant="filled" onclick={handleSaveRole} disabled={isSaving}>
          {#if isSaving}
            <Loader2 class="w-4 h-4 animate-spin" />
          {/if}
          <span>{roleModalMode === 'create' ? 'Buat Peran' : 'Simpan Peran'}</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}

<!-- ========================================================
     DIALOG KONFIRMASI HAPUS PENGGUNA
     ======================================================== -->
{#if showDeleteUserConfirm && userToDelete}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl max-w-sm w-full p-6 shadow-xl border border-[#e1e5ea] flex flex-col gap-4 text-center">
      <div class="w-12 h-12 rounded-full bg-[#fce8e6] text-[#c5221f] flex items-center justify-center mx-auto">
        <Trash2 class="w-6 h-6" />
      </div>
      <div>
        <h3 class="text-base font-bold text-[#1f1f1f]">Hapus / Nonaktifkan Pengguna?</h3>
        <p class="text-xs text-[#444746] mt-1 leading-relaxed">
          Akun <strong>@{userToDelete.username}</strong> ({userToDelete.name}) akan dinonaktifkan dan seluruh sesi aktifnya akan dicabut.
        </p>
      </div>
      <div class="flex justify-center gap-2 pt-2">
        <M3Button variant="outlined" onclick={() => { showDeleteUserConfirm = false; userToDelete = null; }}>
          <span>Batal</span>
        </M3Button>
        <M3Button variant="filled" onclick={handleDeleteUser} disabled={isDeletingUser}>
          {#if isDeletingUser}
            <Loader2 class="w-4 h-4 animate-spin" />
          {/if}
          <span>Ya, Nonaktifkan</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}

<!-- ========================================================
     DIALOG KONFIRMASI HAPUS PERAN
     ======================================================== -->
{#if showDeleteRoleConfirm && roleToDelete}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl max-w-sm w-full p-6 shadow-xl border border-[#e1e5ea] flex flex-col gap-4 text-center">
      <div class="w-12 h-12 rounded-full bg-[#fce8e6] text-[#c5221f] flex items-center justify-center mx-auto">
        <Trash2 class="w-6 h-6" />
      </div>
      <div>
        <h3 class="text-base font-bold text-[#1f1f1f]">Hapus Peran (Role)?</h3>
        <p class="text-xs text-[#444746] mt-1 leading-relaxed">
          Peran <strong>{roleToDelete.name} ({roleToDelete.code})</strong> akan dihapus permanen dari sistem.
        </p>
      </div>
      <div class="flex justify-center gap-2 pt-2">
        <M3Button variant="outlined" onclick={() => { showDeleteRoleConfirm = false; roleToDelete = null; }}>
          <span>Batal</span>
        </M3Button>
        <M3Button variant="filled" onclick={handleDeleteRole} disabled={isDeletingRole}>
          {#if isDeletingRole}
            <Loader2 class="w-4 h-4 animate-spin" />
          {/if}
          <span>Ya, Hapus</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}
