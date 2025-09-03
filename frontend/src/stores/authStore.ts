import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id?: string;
  userId: string;
  username: string;
  email: string;
  phone?: string;
  bio?: string;
  avatar?: string;
  createdAt?: string;
  updatedAt?: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

interface AuthActions {
  clearError: () => void;
  setLoading: (loading: boolean) => void;
  // 设置认证状态的方法
  setAuthState: (user: User | null, token: string | null, isAuthenticated: boolean) => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, _) => ({
      // 初始状态 - 默认为已认证状态
      user: {
        userId: 'demo-user',
        username: 'Demo User',
        email: 'demo@example.com',
      },
      token: 'demo-token',
      isAuthenticated: true,
      isLoading: false,
      error: null,

      // 设置认证状态
      setAuthState: (user: User | null, token: string | null, isAuthenticated: boolean) => {
        set({
          user,
          token,
          isAuthenticated,
          isLoading: false,
          error: null,
        });
      },

      // 清除错误
      clearError: () => {
        set({ error: null });
      },

      // 设置加载状态
      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state: AuthStore) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

// 简化初始化认证状态
export const initializeAuth = () => {
  // 不需要任何验证，直接使用默认认证状态
  console.log('Auth initialized with demo state');
};