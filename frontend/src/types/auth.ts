export interface AuthUser {
  id: string;
  email: string;
  oauth_provider?: string;
  oauth_id?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: AuthUser;
}

export interface AuthState {
  user: AuthUser | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

export interface OAuthProvider {
  id: 'google' | 'github';
  name: string;
  icon: string;
}