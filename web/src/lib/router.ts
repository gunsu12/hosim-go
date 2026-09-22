// Lightweight Hash-based SPA Router for HOSIM EHR

export const ROUTE_MAP: Record<string, string> = {
  // Clinical / EHR
  'physical': '/clinical/physical',
  'anamnesis': '/clinical/anamnesis',
  'odontogram': '/clinical/odontogram',

  // Diagnostics
  'lab': '/diagnostics/lab',
  'radiology': '/diagnostics/radiology',
  'pharmacy': '/diagnostics/pharmacy',

  // Master Data Domains
  'master-patient': '/master/patient',
  'master-practitioner': '/master/practitioner',
  'master-departement': '/master/departement',
  'master-service-unit': '/master/service-unit',
  'master-room': '/master/room',
  'master-payer': '/master/payer',
  'master-referal': '/master/referal',
  'master-tariff-class': '/master/tariff-class',

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
  if (typeof window === 'undefined') return 'physical';

  const hash = window.location.hash.replace(/^#\/?/, '/');
  if (!hash || hash === '/') {
    return 'physical';
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

  return 'physical';
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
