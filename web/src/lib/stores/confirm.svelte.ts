/**
 * Confirmation & Alert Dialog Store for Hosim-GO
 * Provides Promise-based dialogs replacing window.confirm() and window.alert()
 */

export type DialogVariant = 'primary' | 'danger' | 'warning' | 'info';

export interface ConfirmDialogOptions {
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: DialogVariant;
}

export interface AlertDialogOptions {
  title?: string;
  message: string;
  buttonText?: string;
  confirmText?: string;
  variant?: DialogVariant;
}

class ConfirmStore {
  isOpen = $state(false);
  title = $state('');
  message = $state('');
  confirmText = $state('OK');
  cancelText = $state('Batal');
  variant = $state<DialogVariant>('primary');
  isAlertOnly = $state(false);

  private resolver: ((value: boolean) => void) | null = null;

  confirm(options: ConfirmDialogOptions | string): Promise<boolean> {
    const opts: ConfirmDialogOptions = typeof options === 'string' ? { message: options } : options;
    this.title = opts.title || 'Konfirmasi Tindakan';
    this.message = opts.message;
    this.confirmText = opts.confirmText || 'Ya, Lanjutkan';
    this.cancelText = opts.cancelText || 'Batal';
    this.variant = opts.variant || 'primary';
    this.isAlertOnly = false;
    this.isOpen = true;

    return new Promise<boolean>((resolve) => {
      this.resolver = resolve;
    });
  }

  alert(options: AlertDialogOptions | string): Promise<void> {
    const opts: AlertDialogOptions = typeof options === 'string' ? { message: options } : options;
    this.title = opts.title || 'Informasi';
    this.message = opts.message;
    this.confirmText = opts.buttonText || opts.confirmText || 'Mengerti';
    this.cancelText = '';
    this.variant = opts.variant || 'info';
    this.isAlertOnly = true;
    this.isOpen = true;

    return new Promise<void>((resolve) => {
      this.resolver = () => resolve();
    });
  }

  handleConfirm() {
    this.isOpen = false;
    if (this.resolver) {
      this.resolver(true);
      this.resolver = null;
    }
  }

  handleCancel() {
    this.isOpen = false;
    if (this.resolver) {
      this.resolver(false);
      this.resolver = null;
    }
  }
}

export const confirmStore = new ConfirmStore();

/**
 * Trigger an M3 confirmation dialog with a Promise resolving to boolean (true: confirmed, false: cancelled)
 * Example:
 * const ok = await confirmDialog({
 *   title: 'Hapus Item',
 *   message: 'Yakin ingin menghapus item ini?',
 *   variant: 'danger',
 *   confirmText: 'Hapus'
 * });
 * if (!ok) return;
 */
export const confirmDialog = (options: ConfirmDialogOptions | string) => confirmStore.confirm(options);

/**
 * Trigger an M3 alert / info dialog with a Promise resolving when dismissed
 */
export const alertDialog = (options: AlertDialogOptions | string) => confirmStore.alert(options);
