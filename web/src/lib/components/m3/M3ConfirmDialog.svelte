<script lang="ts">
  import {
    AlertTriangle,
    AlertCircle,
    HelpCircle,
    Info,
    X
  } from '@lucide/svelte';
  import { confirmStore, type DialogVariant } from '../../stores/confirm.svelte';

  interface Props {
    isOpen?: boolean;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    variant?: DialogVariant;
    isAlertOnly?: boolean;
    onConfirm?: () => void;
    onCancel?: () => void;
  }

  let {
    isOpen: propIsOpen,
    title: propTitle,
    message: propMessage,
    confirmText: propConfirmText,
    cancelText: propCancelText,
    variant: propVariant,
    isAlertOnly: propIsAlertOnly,
    onConfirm,
    onCancel
  }: Props = $props();

  // If props are passed, use props; otherwise fallback to global store
  const isControlled = $derived(propIsOpen !== undefined);
  const activeIsOpen = $derived(isControlled ? !!propIsOpen : confirmStore.isOpen);
  const activeTitle = $derived(isControlled ? (propTitle || 'Konfirmasi') : confirmStore.title);
  const activeMessage = $derived(isControlled ? (propMessage || '') : confirmStore.message);
  const activeConfirmText = $derived(isControlled ? (propConfirmText || 'Ya, Lanjutkan') : confirmStore.confirmText);
  const activeCancelText = $derived(isControlled ? (propCancelText || 'Batal') : confirmStore.cancelText);
  const activeVariant = $derived(isControlled ? (propVariant || 'primary') : confirmStore.variant);
  const activeIsAlertOnly = $derived(isControlled ? !!propIsAlertOnly : confirmStore.isAlertOnly);

  function handleConfirmAction() {
    if (isControlled) {
      onConfirm?.();
    } else {
      confirmStore.handleConfirm();
    }
  }

  function handleCancelAction() {
    if (isControlled) {
      onCancel?.();
    } else {
      confirmStore.handleCancel();
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (!activeIsOpen) return;
    if (e.key === 'Escape') {
      e.preventDefault();
      handleCancelAction();
    } else if (e.key === 'Enter') {
      e.preventDefault();
      handleConfirmAction();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if activeIsOpen}
  <!-- Backdrop Overlay -->
  <div
    class="fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150"
    role="presentation"
    onclick={handleCancelAction}
  >
    <!-- M3 Dialog Surface -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      role="alertdialog"
      aria-modal="true"
      tabindex="-1"
      aria-labelledby="m3-dialog-title"
      aria-describedby="m3-dialog-desc"
      class="w-full max-w-md bg-white dark:bg-[#1e1f20] rounded-[28px] shadow-2xl border border-[#e1e5ea] dark:border-[#353738] p-6 flex flex-col gap-4 animate-in zoom-in-95 duration-150"
      onclick={(e) => e.stopPropagation()}
    >
      <div class="flex items-start gap-4">
        <!-- Variant Icon Avatar -->
        {#if activeVariant === 'danger'}
          <div class="w-12 h-12 rounded-2xl bg-red-100 text-red-600 dark:bg-red-950/70 dark:text-red-400 border border-red-200/80 dark:border-red-800/60 flex items-center justify-center shrink-0">
            <AlertTriangle class="w-6 h-6" />
          </div>
        {:else if activeVariant === 'warning'}
          <div class="w-12 h-12 rounded-2xl bg-amber-100 text-amber-700 dark:bg-amber-950/70 dark:text-amber-300 border border-amber-200/80 dark:border-amber-800/60 flex items-center justify-center shrink-0">
            <AlertCircle class="w-6 h-6" />
          </div>
        {:else if activeVariant === 'info'}
          <div class="w-12 h-12 rounded-2xl bg-indigo-100 text-indigo-700 dark:bg-indigo-950/70 dark:text-indigo-300 border border-indigo-200/80 dark:border-indigo-800/60 flex items-center justify-center shrink-0">
            <Info class="w-6 h-6" />
          </div>
        {:else}
          <!-- Default: primary -->
          <div class="w-12 h-12 rounded-2xl bg-blue-100 text-[#0b57d0] dark:bg-blue-950/70 dark:text-[#a8c7fa] border border-blue-200/80 dark:border-blue-800/60 flex items-center justify-center shrink-0">
            <HelpCircle class="w-6 h-6" />
          </div>
        {/if}

        <!-- Dialog Content -->
        <div class="flex-1 min-w-0">
          <h3 id="m3-dialog-title" class="text-base sm:text-lg font-bold text-[#1f1f1f] dark:text-[#e3e3e3] leading-snug">
            {activeTitle}
          </h3>
          <p id="m3-dialog-desc" class="text-xs sm:text-sm text-[#444746] dark:text-[#c4c7c5] leading-relaxed mt-1.5 break-words">
            {activeMessage}
          </p>
        </div>
      </div>

      <!-- Action Buttons (M3 Align Right) -->
      <div class="flex items-center justify-end gap-2.5 pt-2">
        {#if !activeIsAlertOnly}
          <button
            type="button"
            onclick={handleCancelAction}
            class="px-5 h-10 rounded-full text-xs font-semibold text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] hover:text-[#1f1f1f] dark:hover:text-white transition-colors cursor-pointer"
          >
            {activeCancelText}
          </button>
        {/if}

        <button
          type="button"
          onclick={handleConfirmAction}
          class="px-6 h-10 rounded-full text-xs font-semibold transition-all duration-150 cursor-pointer shadow-xs {activeVariant === 'danger'
            ? 'bg-red-600 hover:bg-red-700 text-white dark:bg-red-500 dark:hover:bg-red-600'
            : activeVariant === 'warning'
              ? 'bg-amber-600 hover:bg-amber-700 text-white dark:bg-amber-500 dark:text-black dark:hover:bg-amber-400'
              : 'bg-[#0b57d0] hover:bg-[#0842a0] text-white dark:bg-[#a8c7fa] dark:text-[#062e6f] dark:hover:bg-[#8ab4f8]'}"
        >
          {activeConfirmText}
        </button>
      </div>
    </div>
  </div>
{/if}
