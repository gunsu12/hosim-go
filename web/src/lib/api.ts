// API client wrapper for HOSIM Go Backend (TypeScript)

const API_BASE = '/api/v1';

export interface HealthResponse {
  status: 'success' | 'error';
  message: string;
  data?: {
    app_name: string;
    env: string;
    database: string;
    driver: string;
    timestamp: string;
  };
  error?: string;
}

export async function checkBackendHealth(): Promise<HealthResponse> {
  try {
    const res = await fetch(`${API_BASE}/health`, {
      headers: { 'Accept': 'application/json' }
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return await res.json();
  } catch (err) {
    return {
      status: 'error',
      message: 'Tidak dapat terhubung ke Backend Go Gin',
      error: err instanceof Error ? err.message : String(err)
    };
  }
}

export async function login(username: string, password: string): Promise<any> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || 'Login gagal');
  return data;
}
