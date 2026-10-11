// Theme Store & Mode Management for HOSIM (Svelte 5 Runes + TypeScript)
// Supports: 'system' | 'light' | 'dark'

export type ThemeMode = 'system' | 'light' | 'dark';

const STORAGE_KEY = 'hosim_theme';

function getSystemPreference(): boolean {
  if (typeof window === 'undefined') return false;
  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
}

function loadInitialMode(): ThemeMode {
  if (typeof window === 'undefined') return 'system';
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved === 'light' || saved === 'dark' || saved === 'system') {
      return saved;
    }
  } catch (e) {
    console.error('Failed to read theme preference from localStorage', e);
  }
  return 'system';
}

export function createThemeStore() {
  let mode = $state<ThemeMode>(loadInitialMode());
  let systemPrefersDark = $state<boolean>(getSystemPreference());

  // Hitung apakah mode aktif saat ini adalah dark
  let isDark = $derived<boolean>(
    mode === 'dark' || (mode === 'system' && systemPrefersDark)
  );

  let resolved = $derived<'light' | 'dark'>(isDark ? 'dark' : 'light');

  function applyToDocument(dark: boolean, currentMode: ThemeMode) {
    if (typeof document === 'undefined') return;

    const root = document.documentElement;
    if (dark) {
      root.classList.add('dark');
      root.setAttribute('data-theme', 'dark');
    } else {
      root.classList.remove('dark');
      root.setAttribute('data-theme', 'light');
    }

    // Sync meta color-scheme
    const meta = document.querySelector('meta[name="color-scheme"]') as HTMLMetaElement | null;
    if (meta) {
      meta.content = currentMode === 'system' ? 'light dark' : (dark ? 'dark' : 'light');
    }
  }

  // Terapkan tema ke DOM setiap kali resolved berganti
  $effect.root(() => {
    $effect(() => {
      applyToDocument(isDark, mode);
    });
  });

  // Setup media query listener & cross-tab sync
  if (typeof window !== 'undefined') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleMediaChange = (e: MediaQueryListEvent) => {
      systemPrefersDark = e.matches;
    };
    mediaQuery.addEventListener('change', handleMediaChange);

    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === STORAGE_KEY) {
        const val = e.newValue;
        if (val === 'light' || val === 'dark' || val === 'system') {
          mode = val;
        }
      }
    };
    window.addEventListener('storage', handleStorageChange);
  }

  function setMode(newMode: ThemeMode): void {
    mode = newMode;
    try {
      if (newMode === 'system') {
        localStorage.removeItem(STORAGE_KEY);
      } else {
        localStorage.setItem(STORAGE_KEY, newMode);
      }
    } catch (e) {
      console.error('Failed to persist theme to localStorage', e);
    }
  }

  function toggle(): void {
    // Jika sedang dark -> ubah ke light; jika light -> ubah ke dark
    if (isDark) {
      setMode('light');
    } else {
      setMode('dark');
    }
  }

  return {
    get mode(): ThemeMode {
      return mode;
    },
    get isDark(): boolean {
      return isDark;
    },
    get resolved(): 'light' | 'dark' {
      return resolved;
    },
    get systemPrefersDark(): boolean {
      return systemPrefersDark;
    },
    setMode,
    toggle,
  };
}

export const theme = createThemeStore();
