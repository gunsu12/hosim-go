// Authentication & User Profile Types (HOSIM EHR)

export interface UserProfile {
  id: string;
  username: string;
  email: string;
  name: string;
  role: string;
  permissions: string[];
}

export interface AuthSession {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: UserProfile;
  isDemo?: boolean;
}

export interface LoginPayload {
  username: string;
  password: string;
}
