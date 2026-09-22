// Auth Store & Session Management for HOSIM (Svelte 5 Runes + TypeScript)
import { login as apiLogin } from '../api';
import type { AuthSession, UserProfile } from '../types';

const STORAGE_KEY = 'hosim_auth_session';

function loadInitialSession(): AuthSession | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      // Bersihkan jika session tersimpan masih berupa mock demo token lama
      if (parsed?.access_token?.startsWith('mock-jwt-token-demo-') || parsed?.isDemo) {
        localStorage.removeItem(STORAGE_KEY);
        return null;
      }
      return parsed;
    }
  } catch (e) {
    console.error('Failed to load session', e);
  }
  return null;
}

export function createAuthStore() {
  let session = $state<AuthSession | null>(loadInitialSession());

  if (typeof window !== 'undefined') {
    window.addEventListener('hosim:session-refreshed', (e: any) => {
      session = e.detail;
    });
    window.addEventListener('hosim:session-expired', () => {
      session = null;
    });
  }

  return {
    get user(): UserProfile | null {
      return session?.user || null;
    },
    get token(): string | null {
      return session?.access_token || null;
    },
    get isAuthenticated(): boolean {
      return !!session?.user && !!session?.access_token;
    },
    get isDemo(): boolean {
      return !!session?.isDemo;
    },
    hasRole(role: string): boolean {
      if (!session?.user?.role) return false;
      return session.user.role.toUpperCase() === role.toUpperCase();
    },
    hasPermission(permission: string): boolean {
      if (!session?.user) return false;
      if (session.user.role?.toUpperCase() === 'ADMIN') return true;
      const perms = session.user.permissions || [];
      return perms.includes('*') || perms.map(p => p.toLowerCase()).includes(permission.toLowerCase());
    },
    hasAnyPermission(...permissions: string[]): boolean {
      if (!session?.user) return false;
      if (session.user.role?.toUpperCase() === 'ADMIN') return true;
      const perms = session.user.permissions || [];
      if (perms.includes('*')) return true;
      const lowerPerms = perms.map(p => p.toLowerCase());
      return permissions.some(req => lowerPerms.includes(req.toLowerCase()));
    },
    async login(username: string, password: string): Promise<{ success: boolean; user: UserProfile; isDemo?: boolean }> {
      try {
        const res = await apiLogin(username, password);
        if (res && (res.success === true || res.status === 'success') && res.data) {
          session = res.data;
          localStorage.setItem(STORAGE_KEY, JSON.stringify(res.data));
          return { success: true, user: res.data.user };
        }
        throw new Error(res?.message || 'Login gagal');
      } catch (err: any) {
        console.error('Gagal login ke backend API:', err);
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
