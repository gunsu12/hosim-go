<script lang="ts">
  import { Trash2, AlertCircle, Loader2, ShieldAlert } from '@lucide/svelte';
  import M3Modal from '../../../components/m3/M3Modal.svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { deleteUser } from '$lib/api/auth';
  import type { UserRecord } from '../../../types/auth';

  interface Props {
    isOpen: boolean;
    user?: UserRecord | null;
    onClose: () => void;
    onSuccess: () => void;
  }

  let {
    isOpen = false,
    user = null,
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let isDeleting = $state(false);
  let errorMessage = $state<string | null>(null);

  let isAdminProtected = $derived(user?.username === 'admin');

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      isDeleting = false;
    }
  });

  async function handleDelete() {
    if (!user || isAdminProtected) return;
    isDeleting = true;
    errorMessage = null;
    try {
      await deleteUser(user.id);
      onSuccess();
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal menonaktifkan pengguna';
    } finally {
      isDeleting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Nonaktifkan Pengguna?"
  subtitle={user ? `@${user.username} (${user.name})` : '-'}
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

  {#if isAdminProtected}
    <div class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200 flex items-start gap-2.5 text-xs">
      <ShieldAlert class="w-4 h-4 shrink-0 text-amber-600 dark:text-amber-400 mt-0.5" />
      <div>
        <span class="font-bold">Akun Terlindungi: </span>
        <span>Akun superadministrator bawaan sistem (admin) tidak dapat dihapus atau dinonaktifkan.</span>
      </div>
    </div>
  {:else}
    <p class="text-xs text-[#444746] dark:text-[#c4c7c5] leading-relaxed">
      Akun <strong>@{user?.username}</strong> ({user?.name}) akan dinonaktifkan dari sistem. Seluruh sesi login aktif akan dicabut dan pengguna tidak akan dapat mengakses modul apa pun.
    </p>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
      <span>Batal</span>
    </M3Button>
    <M3Button
      variant="filled"
      onclick={handleDelete}
      disabled={isDeleting || isAdminProtected}
      class="bg-rose-600 hover:bg-rose-700 dark:bg-rose-600 dark:hover:bg-rose-700 text-white"
    >
      {#if isDeleting}
        <Loader2 class="w-4 h-4 mr-1.5 animate-spin" />
        <span>Menonaktifkan...</span>
      {:else}
        <Trash2 class="w-4 h-4 mr-1.5" />
        <span>Ya, Nonaktifkan</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
