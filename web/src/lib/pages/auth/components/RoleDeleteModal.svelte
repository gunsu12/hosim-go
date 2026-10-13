<script lang="ts">
  import { Trash2, AlertCircle, Loader2, ShieldAlert } from '@lucide/svelte';
  import M3Modal from '../../../components/m3/M3Modal.svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { deleteRole } from '$lib/api/auth';
  import type { RoleRecord } from '../../../types/auth';

  interface Props {
    isOpen: boolean;
    role?: RoleRecord | null;
    onClose: () => void;
    onSuccess: () => void;
  }

  let {
    isOpen = false,
    role = null,
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let isDeleting = $state(false);
  let errorMessage = $state<string | null>(null);

  // System protected roles
  const PROTECTED_ROLES = ['ADMIN', 'SUPERADMIN'];
  let isSystemProtected = $derived(role ? PROTECTED_ROLES.includes(role.code.toUpperCase()) : false);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      isDeleting = false;
    }
  });

  async function handleDelete() {
    if (!role || isSystemProtected) return;
    isDeleting = true;
    errorMessage = null;
    try {
      await deleteRole(role.id);
      onSuccess();
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal menghapus peran (role)';
    } finally {
      isDeleting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Hapus Peran (Role)?"
  subtitle={role ? `${role.name} (${role.code})` : '-'}
  icon={Trash2}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
  size="md"
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

  {#if isSystemProtected}
    <div class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200 flex items-start gap-2.5 text-xs">
      <ShieldAlert class="w-4 h-4 shrink-0 text-amber-600 dark:text-amber-400 mt-0.5" />
      <div>
        <span class="font-bold">Peran Terlindungi: </span>
        <span>Peran sistem utama ({role?.code}) tidak dapat dihapus karena dibutuhkan untuk operasional otorisasi dasar.</span>
      </div>
    </div>
  {:else}
    <p class="text-xs text-[#444746] dark:text-[#c4c7c5] leading-relaxed">
      Peran <strong>{role?.name} ({role?.code})</strong> akan dihapus permanen dari sistem. Seluruh wewenang yang melekat pada peran ini akan dicabut dari pengguna yang menggunakannya.
    </p>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
      <span>Batal</span>
    </M3Button>
    <M3Button
      variant="filled"
      onclick={handleDelete}
      disabled={isDeleting || isSystemProtected}
      class="bg-rose-600 hover:bg-rose-700 dark:bg-rose-600 dark:hover:bg-rose-700 text-white"
    >
      {#if isDeleting}
        <Loader2 class="w-4 h-4 mr-1.5 animate-spin" />
        <span>Menghapus...</span>
      {:else}
        <Trash2 class="w-4 h-4 mr-1.5" />
        <span>Ya, Hapus</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
