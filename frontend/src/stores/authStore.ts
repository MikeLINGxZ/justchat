import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authClient } from '@/api/authClient';
import { extractErrorMessage } from '@/utils/errorHandler';

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
  logout: () => Promise<void>;
  getCurrentUser: () => Promise<void>;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
  // 新增：直接设置认证状态的方法
  setAuthState: (user: User | null, token: string | null, isAuthenticated: boolean) => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, _) => ({
      // 初始状态
      user: null,
      token: null,
      isAuthenticated: false,
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



      // 登出
      logout: async () => {
        try {
          await authClient.logout();
        } catch (error) {
          console.error('Logout API call failed:', error);
        } finally {
          // 清除本地存储
          localStorage.removeItem('token');
          
          set({
            user: null,
            token: null,
            isAuthenticated: false,
            error: null,
          });
        }
      },

      // 获取当前用户信息
      getCurrentUser: async () => {
        try {
          set({ isLoading: true });
          // 注意：这里需要根据实际的API接口来实现
          // 如果后端没有提供获取当前用户信息的接口，可以跳过
          set({
            isLoading: false,
          });
        } catch (error: any) {
          const errorMessage = extractErrorMessage(error);
          set({
            error: errorMessage,
            isLoading: false,
            isAuthenticated: false,
          });
          
          // 如果获取用户信息失败，清除token
          localStorage.removeItem('token');
        }
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

// 初始化认证状态
export const initializeAuth = () => {
  const token = localStorage.getItem('token');
  if (token) {
    useAuthStore.getState().getCurrentUser();
  }
};