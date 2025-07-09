import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { AuthState, LoginRequest, RegisterRequest } from '@/types/auth';
import { authAPI } from '@/lib/auth-api';

interface AuthStore extends AuthState {
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  getProfile: () => Promise<void>;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
  oauthLogin: (provider: string) => Promise<void>;
  handleOAuthCallback: (provider: string, code: string) => Promise<void>;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isLoading: false,
      error: null,

      login: async (data: LoginRequest) => {
        try {
          set({ isLoading: true, error: null });
          const response = await authAPI.login(data);
          
          // 保存到localStorage
          if (typeof window !== 'undefined') {
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));
          }
          
          set({ 
            user: response.user, 
            token: response.token, 
            isLoading: false 
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Login failed', 
            isLoading: false 
          });
          throw error;
        }
      },

      register: async (data: RegisterRequest) => {
        try {
          set({ isLoading: true, error: null });
          const response = await authAPI.register(data);
          
          // 保存到localStorage
          if (typeof window !== 'undefined') {
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));
          }
          
          set({ 
            user: response.user, 
            token: response.token, 
            isLoading: false 
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Registration failed', 
            isLoading: false 
          });
          throw error;
        }
      },

      logout: () => {
        authAPI.logout();
        set({ user: null, token: null, error: null });
      },

      getProfile: async () => {
        try {
          set({ isLoading: true });
          const user = await authAPI.getProfile();
          set({ user, isLoading: false });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to get profile', 
            isLoading: false 
          });
          // 如果获取用户信息失败，可能token已过期，执行登出
          get().logout();
        }
      },

      oauthLogin: async (provider: string) => {
        try {
          set({ isLoading: true, error: null });
          const { auth_url } = await authAPI.getOAuthURL(provider);
          window.location.href = auth_url;
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'OAuth login failed', 
            isLoading: false 
          });
        }
      },

      handleOAuthCallback: async (provider: string, code: string) => {
        try {
          set({ isLoading: true, error: null });
          const response = await authAPI.oauthCallback(provider, code);
          
          // 保存到localStorage
          if (typeof window !== 'undefined') {
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));
          }
          
          set({ 
            user: response.user, 
            token: response.token, 
            isLoading: false 
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'OAuth callback failed', 
            isLoading: false 
          });
          throw error;
        }
      },

      clearError: () => set({ error: null }),
      setLoading: (loading: boolean) => set({ isLoading: loading }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token 
      }),
    }
  )
);