// Lightweight Hash-based SPA Router for HOSIM EHR

export const ROUTE_MAP: Record<string, string> = {
  // Desktop Launcher (Root Hub)
  'desktop': '/desktop',

  // Clinical / EHR
  'physical': '/clinical/physical',
  'anamnesis': '/clinical/anamnesis',

  // Master Data Domains
  'master-patient': '/master/patient',
  'master-practitioner': '/master/practitioner',
  'master-departement': '/master/departement',
  'master-service-unit': '/master/service-unit',
  'master-room': '/master/room',
  'master-customer': '/master/customer',
  'master-payer': '/master/customer',
  'master-referal': '/master/referal',
  'master-tariff-class': '/master/tariff-class',
  'master-item': '/master/item',
  'master-coa': '/master/coa',

  // Workspaces Preview
  'outpatient-workspace': '/outpatient/workspace',
  'emergency-workspace': '/emergency/workspace',
  'inpatient-workspace': '/inpatient/workspace',
  'pharmacy-workspace': '/pharmacy/workspace',
  'inventory-workspace': '/inventory/workspace',
  'billing-workspace': '/billing/workspace',
  'audit-workspace': '/audit/workspace',

  // Finance & Tariff Management
  'finance-tariff': '/finance/tariff',
  'finance-price-plan': '/finance/price-plan',
  'finance-tariff-matrix': '/finance/tariff-matrix',
  'finance-tariff-lookup': '/finance/tariff-lookup',
  'finance-tariff-class': '/finance/tariff-class',
  'finance-tariff-component': '/finance/tariff-component',

  // Accounting & Financial Management
  'accounting-coa': '/accounting/coa',
  'accounting-journals': '/accounting/journals',
  'accounting-ledger': '/accounting/ledger',
  'accounting-reports': '/accounting/reports',

  // User, Role & Permission Management
  'auth-users': '/auth/users',
  'auth-roles': '/auth/roles',
  'auth-permissions': '/auth/permissions',

  // Activity & Schedule
  'history': '/activity/history',
  'schedule': '/activity/schedule'
};

// Reverse map: path -> navId
export const PATH_TO_NAV: Record<string, string> = Object.entries(ROUTE_MAP).reduce((acc, [navId, path]) => {
  acc[path] = navId;
  acc[path.replace(/^\//, '')] = navId;
  return acc;
}, {} as Record<string, string>);

/**
 * Mendapatkan navId aktif berdasarkan window.location.hash saat ini
 */
export function getNavFromCurrentHash(): string {
  if (typeof window === 'undefined') return 'desktop';

  const hash = window.location.hash.replace(/^#\/?/, '/');
  if (!hash || hash === '/' || hash === '/desktop') {
    return 'desktop';
  }

  // Exact match
  if (PATH_TO_NAV[hash]) {
    return PATH_TO_NAV[hash];
  }

  // Clean trailing slash
  const cleanPath = hash.replace(/\/+$/, '');
  if (PATH_TO_NAV[cleanPath]) {
    return PATH_TO_NAV[cleanPath];
  }

  // Fallback jika hash langsung berupa ID (misal #master-patient)
  const rawHash = window.location.hash.replace(/^#/, '');
  if (ROUTE_MAP[rawHash]) {
    return rawHash;
  }

  return 'desktop';
}

/**
 * Sinkronisasi URL browser hash dengan navId
 */
export function setHashFromNav(navId: string, replace = false): void {
  if (typeof window === 'undefined') return;

  const targetPath = ROUTE_MAP[navId] || `/${navId}`;
  const targetHash = `#${targetPath}`;

  if (window.location.hash !== targetHash) {
    if (replace) {
      window.history.replaceState(null, '', targetHash);
    } else {
      window.location.hash = targetHash;
    }
  }
}

/**
 * Menentukan modul aktif berdasarkan navId
 */
export function getModuleFromNav(navId: string): string | null {
  if (!navId || navId === 'desktop') return null;
  if (navId.startsWith('master-')) return 'master';
  if (navId.startsWith('auth-')) return 'auth';
  if (navId === 'physical' || navId === 'anamnesis' || navId === 'history' || navId === 'schedule') return 'clinical';
  if (navId.startsWith('outpatient')) return 'outpatient';
  if (navId.startsWith('emergency')) return 'emergency';
  if (navId.startsWith('inpatient')) return 'inpatient';
  if (navId.startsWith('pharmacy')) return 'pharmacy';
  if (navId.startsWith('inventory')) return 'inventory';
  if (navId.startsWith('billing')) return 'billing';
  if (navId.startsWith('audit')) return 'audit';
  if (navId.startsWith('accounting')) return 'accounting';
  if (navId.startsWith('finance-')) return 'finance';
  return null;
}
