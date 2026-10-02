// Authentication & User Profile Types (HOSIM EHR)

export interface UserProfile {
  id: string;
  username: string;
  email: string;
  name: string;
  role: string;
  role_id?: string;
  role_name?: string;
  permissions: string[];
  is_active?: boolean;
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

// ==========================================
// User Management Types
// ==========================================

export interface UserRecord {
  id: string;
  username: string;
  email: string;
  name: string;
  role_id?: string;
  role: string;
  role_name?: string;
  permissions: string[];
  is_active: boolean;
  created_at: string;
}

export interface CreateUserDTO {
  username: string;
  email: string;
  password: string;
  name: string;
  role_id?: string;
  is_active?: boolean;
}

export interface UpdateUserDTO {
  name: string;
  email: string;
  role_id?: string;
  is_active?: boolean;
  password?: string;
}

// ==========================================
// Role & Permission Management Types
// ==========================================

export interface PermissionRecord {
  id: string;
  code: string;
  name: string;
  module: string;
  description?: string;
  created_at?: string;
}

export interface RoleRecord {
  id: string;
  code: string;
  name: string;
  description?: string;
  permissions?: PermissionRecord[];
  created_at?: string;
}

export interface CreateRoleDTO {
  code: string;
  name: string;
  description?: string;
  permission_ids?: string[];
}

export interface UpdateRoleDTO {
  name: string;
  description?: string;
  permission_ids?: string[];
}
