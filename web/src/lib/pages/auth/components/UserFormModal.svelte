<script lang="ts">
  import { User, CheckCircle2, AlertCircle, Loader2 } from '@lucide/svelte';
  import M3Modal from '../../../components/m3/M3Modal.svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import M3ToggleCard from '../../../components/m3/M3ToggleCard.svelte';
  import M3Select from '../../../components/m3/M3Select.svelte';
  import PractitionerLookupInput from '../../../components/PractitionerLookupInput.svelte';
  import { createUser, updateUser } from '$lib/api/auth';
  import type { UserRecord, RoleRecord, CreateUserDTO, UpdateUserDTO } from '../../../types/auth';
  import type { PractitionerRecord } from '../../../types/master/practitioner';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    user?: UserRecord | null;
    roles: RoleRecord[];
    practitioners?: PractitionerRecord[];
    onClose: () => void;
    onSuccess: (saved: UserRecord) => void;
  }

  let {
    isOpen = false,
    mode = 'create',
    user = null,
    roles = [],
    practitioners = [],
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let username = $state('');
  let email = $state('');
  let name = $state('');
  let password = $state('');
  let roleId = $state('');
  let practitionerId = $state('');
  let isActive = $state(true);
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      if (mode === 'edit' && user) {
        username = user.username;
        email = user.email;
        name = user.name;
        password = '';
        roleId = user.role_id || (roles.find((r) => r.code === user.role)?.id || '');
        practitionerId = user.practitioner_id || '';
        isActive = user.is_active;
      } else {
        username = '';
        email = '';
        name = '';
        password = '';
        roleId = roles.length > 0 ? roles[0].id : '';
        practitionerId = '';
        isActive = true;
      }
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    errorMessage = null;

    if (mode === 'create' && !username.trim()) {
      errorMessage = 'Username akun wajib diisi';
      return;
    }
    if (!email.trim() || !email.includes('@')) {
      errorMessage = 'Alamat email tidak valid';
      return;
    }
    if (!name.trim()) {
      errorMessage = 'Nama lengkap pengguna wajib diisi';
      return;
    }
    if (mode === 'create' && (!password || password.length < 6)) {
      errorMessage = 'Password wajib diisi minimal 6 karakter';
      return;
    }

    isSubmitting = true;
    try {
      if (mode === 'create') {
        const payload: CreateUserDTO = {
          username: username.trim(),
          email: email.trim().toLowerCase(),
          password: password,
          name: name.trim(),
          role_id: roleId || undefined,
          practitioner_id: practitionerId.trim() ? practitionerId.trim() : undefined,
          is_active: isActive
        };
        const created = await createUser(payload);
        onSuccess(created);
      } else if (user) {
        const payload: UpdateUserDTO = {
          name: name.trim(),
          email: email.trim().toLowerCase(),
          role_id: roleId || undefined,
          practitioner_id: practitionerId.trim(),
          is_active: isActive,
          password: password.trim() ? password.trim() : undefined
        };
        const updated = await updateUser(user.id, payload);
        onSuccess(updated);
      }
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Terjadi kesalahan saat menyimpan data pengguna';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'create' ? 'Tambah Pengguna Baru' : 'Ubah Data Pengguna'}
  subtitle={mode === 'create' ? 'Buat akun pengguna baru rumah sakit' : `@${user?.username || '-'}`}
  icon={User}
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

  <!-- Section 1: Akun & Identitas (2 Kolom) -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="form-user-username" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Username <span class="text-rose-500">*</span>
      </label>
      <input
        id="form-user-username"
        type="text"
        bind:value={username}
        disabled={mode === 'edit'}
        placeholder="misal: dokter.hendra"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-60 disabled:cursor-not-allowed focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        ID login unik di dalam sistem.
      </p>
    </div>

    <div>
      <label for="form-user-email" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Alamat Email <span class="text-rose-500">*</span>
      </label>
      <input
        id="form-user-email"
        type="email"
        bind:value={email}
        placeholder="misal: hendra@hosim.local"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        Digunakan untuk notifikasi & audit.
      </p>
    </div>

    <div class="sm:col-span-2">
      <label for="form-user-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Nama Lengkap & Gelar <span class="text-rose-500">*</span>
      </label>
      <input
        id="form-user-name"
        type="text"
        bind:value={name}
        placeholder="misal: dr. Hendra Wijaya, Sp.B"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Section 2: Wewenang & Integrasi Tenaga Medis -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] space-y-4">
    <h4 class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider">
      Wewenang & Integrasi Pelayanan Klinis
    </h4>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <M3Select
          id="form-user-role"
          label="Peran Pengguna (Role)"
          required={true}
          bind:value={roleId}
          options={roles.map((r) => ({
            value: r.id,
            label: r.name,
            badge: r.code
          }))}
          helperText="Menentukan matriks hak akses modul."
        />
      </div>

      <div>
        <PractitionerLookupInput
          id="form-user-practitioner"
          bind:value={practitionerId}
          label="Tautkan Tenaga Medis / Paramedis"
          placeholder="Cari atau pilih tenaga medis (opsional)..."
          allowClear={true}
          helperText="Dokter bertanggung jawab atas SOAP rekam medis & SATUSEHAT."
        />
      </div>
    </div>
  </div>

  <!-- Section 3: Keamanan & Status Akun -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] space-y-4">
    <h4 class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider">
      Keamanan & Status Akun
    </h4>

    <div>
      <label for="form-user-password" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        {mode === 'create' ? 'Kata Sandi (Password)' : 'Reset Kata Sandi (Kosongkan bila tidak diubah)'}
        {#if mode === 'create'}<span class="text-rose-500">*</span>{/if}
      </label>
      <input
        id="form-user-password"
        type="password"
        bind:value={password}
        placeholder={mode === 'create' ? 'Minimal 6 karakter' : 'Tulis password baru jika ingin mereset'}
        required={mode === 'create'}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <M3ToggleCard
      title="Status Akun Aktif"
      description="Pengguna dapat login dan mengakses modul sesuai wewenang peran yang diberikan."
      bind:checked={isActive}
      colorVariant="emerald"
    />
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
        <span>{mode === 'create' ? 'Buat Akun' : 'Simpan Perubahan'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
