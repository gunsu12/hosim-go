import type { AccountType, AccountPosition } from '$lib/types/accounting/account';

export interface AccountTypeMeta {
  label: string;
  prefix: string;
  color: string;
  badgeBg: string;
  badgeText: string;
  defaultPos: AccountPosition;
}

export const accountTypeConfig: Record<AccountType, AccountTypeMeta> = {
  ASSET: {
    label: 'Aset / Aktiva',
    prefix: '1',
    color: 'blue',
    badgeBg: 'bg-blue-100 border-blue-200 dark:bg-blue-950/60 dark:border-blue-800/50',
    badgeText: 'text-blue-800 dark:text-blue-300',
    defaultPos: 'DEBIT'
  },
  LIABILITY: {
    label: 'Kewajiban / Liabilitas',
    prefix: '2',
    color: 'amber',
    badgeBg: 'bg-amber-100 border-amber-200 dark:bg-amber-950/60 dark:border-amber-800/50',
    badgeText: 'text-amber-800 dark:text-amber-300',
    defaultPos: 'CREDIT'
  },
  EQUITY: {
    label: 'Ekuitas / Modal',
    prefix: '3',
    color: 'purple',
    badgeBg: 'bg-purple-100 border-purple-200 dark:bg-purple-950/60 dark:border-purple-800/50',
    badgeText: 'text-purple-800 dark:text-purple-300',
    defaultPos: 'CREDIT'
  },
  REVENUE: {
    label: 'Pendapatan / Revenue',
    prefix: '4',
    color: 'emerald',
    badgeBg: 'bg-emerald-100 border-emerald-200 dark:bg-emerald-950/60 dark:border-emerald-800/50',
    badgeText: 'text-emerald-800 dark:text-emerald-300',
    defaultPos: 'CREDIT'
  },
  EXPENSE: {
    label: 'Beban / Biaya',
    prefix: '5',
    color: 'rose',
    badgeBg: 'bg-rose-100 border-rose-200 dark:bg-rose-950/60 dark:border-rose-800/50',
    badgeText: 'text-rose-800 dark:text-rose-300',
    defaultPos: 'DEBIT'
  }
};
