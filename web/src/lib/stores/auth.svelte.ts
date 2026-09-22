// Auth Store & Session Management for HOSIM (Svelte 5 Runes + TypeScript)
import { login as apiLogin } from '../api';
import type { AuthSession, UserProfile } from '../types';

const STORAGE_KEY = 'hosim_auth_session';

function loadInitialSession(): AuthSession | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return JSON.parse(raw);
  } catch (e) {
    console.error('Failed to load session', e);
  }
  return null;
}

export function createAuthStore() {
  let session = $state<AuthSession | null>(loadInitialSession());

  return {
    get user(): UserProfile | null {
      return session?.user || null;
    },
    get token(): string | null {
      return session?.access_token || null;
    },
    get isAuthenticated(): boolean {
      return !!session?.user;
    },
    async login(username: string, password: string): Promise<{ success: boolean; user: UserProfile; isDemo?: boolean }> {
      try {
        const res = await apiLogin(username, password);
        if (res && res.status === 'success' && res.data) {
          session = res.data;
          localStorage.setItem(STORAGE_KEY, JSON.stringify(res.data));
          return { success: true, user: res.data.user };
        }
        throw new Error(res.message || 'Login gagal');
      } catch (err) {
        // Fallback demo mode jika backend offline atau login kredensial bawaan
        if (
          (username === 'admin' && password === 'admin123') ||
          (username === 'dokter' && password === 'dokter123')
        ) {
          const isDoctor = username === 'dokter';
          const mockData: AuthSession = {
            access_token: 'mock-jwt-token-demo-' + Date.now(),
            refresh_token: 'mock-refresh-token',
            token_type: 'Bearer',
            expires_in: 86400,
            user: {
              id: isDoctor ? 'USR-DOC-001' : 'USR-ADM-001',
              username: username,
              email: isDoctor ? 'hendra@hosim.local' : 'admin@hosim.local',
              name: isDoctor ? 'dr. Hendra Wijaya, Sp.B' : 'Administrator Sistem',
              role: isDoctor ? 'DOCTOR' : 'ADMIN',
              permissions: ['*']
            },
            isDemo: true
          };
          session = mockData;
          localStorage.setItem(STORAGE_KEY, JSON.stringify(mockData));
          return { success: true, user: mockData.user, isDemo: true };
        }
        throw err;
      }
    },
    logout(): void {
      session = null;
      localStorage.removeItem(STORAGE_KEY);
    }
  };
}

export const auth = createAuthStore();
