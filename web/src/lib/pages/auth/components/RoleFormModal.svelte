<script lang="ts">
  import { ShieldCheck, CheckCircle2, AlertCircle, Loader2, CheckSquare, Square } from '@lucide/svelte';
  import M3Modal from '../../../components/m3/M3Modal.svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { createRole, updateRole } from '$lib/api/auth';
  import type { RoleRecord, PermissionRecord, CreateRoleDTO, UpdateRoleDTO } from '../../../types/auth';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    role?: RoleRecord | null;
    permissions: PermissionRecord[];
    onClose: () => void;
    onSuccess: (saved: RoleRecord) => void;
  }

  let {
    isOpen = false,
    mode = 'create',
    role = null,
    permissions = [],
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let code = $state('');
  let name = $state('');
  let description = $state('');
  let selectedPermIds = $state<string[]>([]);
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  // Group permissions by module
  let permissionsByModule = $derived.by(() => {
    const map: Record<string, PermissionRecord[]> = {};
    for (const p of permissions) {
      const mod = p.module || 'Umum';
      if (!map[mod]) map[mod] = [];
      map[mod].push(p);
    }
    return map;
  });

  let permissionModules = $derived(Object.keys(permissionsByModule).sort());

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      if (mode === 'edit' && role) {
        code = role.code;
        name = role.name;
        description = role.description || '';
        // Map existing role permissions
        selectedPermIds = (role.permissions || []).map((p) => p.id);
      } else {
        code = '';
        name = '';
        description = '';
        selectedPermIds = [];
      }
    }
  });

  function togglePermission(permId: string) {
    if (selectedPermIds.includes(permId)) {
      selectedPermIds = selectedPermIds.filter((id) => id !== permId);
    } else {
      selectedPermIds = [...selectedPermIds, permId];
    }
  }

  function toggleModulePermissions(modName: string, modPermIds: string[]) {
    const isAllChecked = modPermIds.every((id) => selectedPermIds.includes(id));
    if (isAllChecked) {
      selectedPermIds = selectedPermIds.filter((id) => !modPermIds.includes(id));
    } else {
      const newIds = new Set([...selectedPermIds, ...modPermIds]);
      selectedPermIds = Array.from(newIds);
    }
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    errorMessage = null;

    if (mode === 'create' && !code.trim()) {
      errorMessage = 'Kode peran wajib diisi';
      return;
    }
    if (!name.trim()) {
      errorMessage = 'Nama peran wajib diisi';
      return;
    }

    isSubmitting = true;
    try {
      if (mode === 'create') {
        const payload: CreateRoleDTO = {
          code: code.trim().toUpperCase(),
          name: name.trim(),
          description: description.trim() || undefined,
          permission_ids: selectedPermIds
        };
        const created = await createRole(payload);
        onSuccess(created);
      } else if (role) {
        const payload: UpdateRoleDTO = {
          name: name.trim(),
          description: description.trim() || undefined,
          permission_ids: selectedPermIds
        };
        const updated = await updateRole(role.id, payload);
        onSuccess(updated);
      }
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Terjadi kesalahan saat menyimpan peran pengguna';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'create' ? 'Tambah Peran (Role) Baru' : `Konfigurasi Peran: ${role?.name || ''}`}
  subtitle={mode === 'create' ? 'Tentukan wewenang dan izin modul untuk kelompok pengguna' : `Kode: ${role?.code || '-'}`}
  icon={ShieldCheck}
  size="2xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <AlertCircle class="w-4 h-4 shrink-0 text-rose-600 dark:text-rose-400 mt-0.5" />
      <div>
        <span class="font-bold">Gagal: </span>
        <span>{errorMessage}</span>
      </div>
    </div>
  {/if}

  <!-- Section 1: Profil Peran -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="form-role-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Kode Peran (Code) <span class="text-rose-500">*</span>
      </label>
      <input
        id="form-role-code"
        type="text"
        bind:value={code}
        disabled={mode === 'edit'}
        placeholder="misal: PHARMACIST"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-60 disabled:cursor-not-allowed focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        Identifier peran unik huruf besar (A-Z).
      </p>
    </div>

    <div>
      <label for="form-role-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Nama Peran <span class="text-rose-500">*</span>
      </label>
      <input
        id="form-role-name"
        type="text"
        bind:value={name}
        placeholder="misal: Apoteker / Petugas Farmasi"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        Nama tampilan yang mudah dikenali.
      </p>
    </div>

    <div class="sm:col-span-2">
      <label for="form-role-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Deskripsi Wewenang
      </label>
      <input
        id="form-role-desc"
        type="text"
        bind:value={description}
        placeholder="Keterangan cakupan kerja pengguna..."
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Section 2: Matriks Hak Akses Modul -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] space-y-3">
    <div class="flex items-center justify-between">
      <h4 class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider">
        Matriks Hak Akses Modul ({selectedPermIds.length} Izin Dipilih)
      </h4>
      <span class="text-[11px] text-[#747775] dark:text-[#8e918f]">Klik judul modul untuk pilih semua</span>
    </div>

    <div class="space-y-3 max-h-72 overflow-y-auto pr-1">
      {#each permissionModules as modName}
        {@const modPerms = permissionsByModule[modName] || []}
        {@const modPermIds = modPerms.map((p) => p.id)}
        {@const isAllChecked = modPermIds.length > 0 && modPermIds.every((id) => selectedPermIds.includes(id))}

        <div class="p-3.5 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a]/50">
          <div class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-[#353738] pb-2 mb-2.5">
            <button
              type="button"
              onclick={() => toggleModulePermissions(modName, modPermIds)}
              class="flex items-center gap-1.5 font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3] hover:text-[#0b57d0] dark:hover:text-[#a8c7fa] cursor-pointer"
            >
              {#if isAllChecked}
                <CheckSquare class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa]" />
              {:else}
                <Square class="w-4 h-4 text-[#747775] dark:text-[#8e918f]" />
              {/if}
              <span>{modName}</span>
            </button>
            <span class="text-[10px] text-[#747775] dark:text-[#8e918f] font-semibold">
              {modPermIds.filter((id) => selectedPermIds.includes(id)).length}/{modPermIds.length} Izin
            </span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {#each modPerms as p}
              <label class="flex items-center gap-2.5 p-2 rounded-xl hover:bg-white dark:hover:bg-[#1e1f20] transition-colors cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={selectedPermIds.includes(p.id)}
                  onchange={() => togglePermission(p.id)}
                  class="rounded accent-[#0b57d0] dark:accent-[#a8c7fa] w-3.5 h-3.5 cursor-pointer"
                />
                <div class="flex flex-col min-w-0">
                  <span class="text-[11px] font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">{p.name}</span>
                  <span class="text-[9px] font-mono text-[#747775] dark:text-[#8e918f] truncate">{p.code}</span>
                </div>
              </label>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      <span>Batal</span>
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSubmitting}>
      {#if isSubmitting}
        <Loader2 class="w-4 h-4 mr-1.5 animate-spin" />
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1.5" />
        <span>{mode === 'create' ? 'Buat Peran' : 'Simpan Perubahan'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
