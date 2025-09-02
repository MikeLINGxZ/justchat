import { useState, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { authClient } from '@/api/authClient';
import type { ServerLoginRequest } from '@/api/authClient';
import { useAuthStore } from '@/stores/authStore';
import { hashPassword } from '@/utils/crypto';
import { extractErrorMessage } from '@/utils/errorHandler';

interface LoginCredentials {
  username: string;
  password: string;
}

interface UseLoginReturn {
  isLoading: boolean;
  error: string | null;
  login: (credentials: LoginCredentials) => Promise<void>;
  clearError: () => void;
}

export const useLogin = (): UseLoginReturn => {
  const navigate = useNavigate();
  const location = useLocation();
  
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 获取重定向路径
  const from = (location.state as any)?.from?.pathname || '/home';

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const login = useCallback(async (credentials: LoginCredentials) => {
    try {
      setIsLoading(true);
      setError(null);

      // 将密码转换为MD5
      const loginData: ServerLoginRequest = {
        loginField: credentials.username,
        passwordMd5: hashPassword(credentials.password),
      };

      // 直接调用API，不通过store
      const response = await authClient.login(loginData);

      // 验证响应
      if (!response.accessToken) {
        throw new Error('登录失败：未收到访问令牌');
      }

      // 保存token到localStorage
      localStorage.setItem('token', response.accessToken);

      // 更新store状态
      const authStore = useAuthStore.getState();
      const userInfo = response.userInfo ? {
        userId: response.userInfo.userId || '',
        username: response.userInfo.username || '',
        email: response.userInfo.email || '',
        createdAt: response.userInfo.createdAt,
        updatedAt: response.userInfo.updatedAt
      } : null;

      // 使用新的setAuthState方法更新store状态
      authStore.setAuthState(userInfo, response.accessToken, true);

      // 登录成功后导航到目标页面
      navigate(from, { replace: true });

    } catch (error: any) {
      // 使用新的错误处理方法提取用户友好的错误消息
      const errorMessage = extractErrorMessage(error);
      setError(errorMessage);
      
      // 清除可能的token
      localStorage.removeItem('token');
      
      // 更新store状态为未认证
      const authStore = useAuthStore.getState();
      authStore.setAuthState(null, null, false);
      
      console.error('Login failed:', error);
    } finally {
      setIsLoading(false);
    }
  }, [navigate, from]);

  return {
    isLoading,
    error,
    login,
    clearError,
  };
};
