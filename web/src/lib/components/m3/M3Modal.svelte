<script lang="ts">
  import { X } from '@lucide/svelte';
  import type { Snippet, Component } from 'svelte';

  export type ModalSize = 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl' | '4xl';

  interface Props {
    isOpen: boolean;
    title: string;
    subtitle?: string;
    size?: ModalSize;
    icon?: Component<any>;
    iconClass?: string;
    closeOnEscape?: boolean;
    closeOnBackdrop?: boolean;
    onsubmit?: (e: SubmitEvent) => void;
    onClose: () => void;
    // Snippets
    children?: Snippet;
    footer?: Snippet;
    headerExtra?: Snippet;
  }

  let {
    isOpen = false,
    title,
    subtitle = '',
    size = '2xl',
    icon,
    iconClass = 'bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]',
    closeOnEscape = true,
    closeOnBackdrop = false,
    onsubmit,
    onClose,
    children,
    footer,
    headerExtra
  }: Props = $props();

  const sizeClasses: Record<ModalSize, string> = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
    '3xl': 'max-w-3xl',
    '4xl': 'max-w-4xl'
  };

  function handleKeydown(e: KeyboardEvent) {
    if (!isOpen) return;
    if (closeOnEscape && e.key === 'Escape') {
      e.preventDefault();
      onClose();
    }
  }

  function handleBackdropClick(e: MouseEvent) {
    if (closeOnBackdrop && e.target === e.currentTarget) {
      onClose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
  <!-- Backdrop Overlay -->
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="fixed inset-0 z-50 bg-black/40 dark:bg-black/60 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto animate-in fade-in duration-150"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={handleBackdropClick}
  >
    <!-- M3 Dialog Surface -->
    <div
      class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738] shadow-2xl w-full {sizeClasses[size] || sizeClasses['2xl']} max-h-[92vh] flex flex-col overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div
        class="px-6 py-4 border-b border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between bg-[#f8fafd] dark:bg-[#28292a] shrink-0"
      >
        <div class="flex items-center gap-3 min-w-0">
          {#if icon}
            {@const IconComponent = icon}
            <div
              class="w-10 h-10 rounded-2xl flex items-center justify-center font-bold shrink-0 {iconClass}"
            >
              <IconComponent class="w-5 h-5" />
            </div>
          {/if}
          <div class="min-w-0">
            <h3 class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
              {title}
            </h3>
            {#if subtitle}
              <p class="text-xs text-[#444746] dark:text-[#c4c7c5] truncate">
                {subtitle}
              </p>
            {/if}
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          {@render headerExtra?.()}
          <button
            onclick={onClose}
            type="button"
            aria-label="Tutup dialog"
            class="text-[#747775] hover:text-[#1f1f1f] dark:hover:text-white p-1.5 rounded-full hover:bg-[#e1e5ea] dark:hover:bg-[#353738] transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Main Body & Form Wrapper -->
      {#if onsubmit}
        <form {onsubmit} class="flex-1 flex flex-col overflow-hidden min-h-0">
          <div class="flex-1 overflow-y-auto p-6 space-y-4 text-xs text-[#1f1f1f] dark:text-[#e3e3e3]">
            {@render children?.()}
          </div>

          {#if footer}
            <div
              class="px-6 py-4 border-t border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a] flex items-center justify-end gap-3 shrink-0"
            >
              {@render footer()}
            </div>
          {/if}
        </form>
      {:else}
        <div class="flex-1 overflow-y-auto p-6 space-y-4 text-xs text-[#1f1f1f] dark:text-[#e3e3e3]">
          {@render children?.()}
        </div>

        {#if footer}
          <div
            class="px-6 py-4 border-t border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a] flex items-center justify-end gap-3 shrink-0"
          >
            {@render footer()}
          </div>
        {/if}
      {/if}
    </div>
  </div>
{/if}
