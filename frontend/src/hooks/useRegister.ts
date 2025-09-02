import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { authClient } from '@/api/authClient';
import type { ServerRegisterRequest } from '@/api/authClient';
import { hashPassword } from '@/utils/crypto';
import { extractErrorMessage } from '@/utils/errorHandler';

interface RegisterCredentials {
  username: string;
  email: string;
  password: string;
  emailVerificationCode: string;
}

interface UseRegisterReturn {
  isLoading: boolean;
  error: string | null;
  register: (credentials: RegisterCredentials) => Promise<void>;
  clearError: () => void;
}

export const useRegister = (): UseRegisterReturn => {
  const navigate = useNavigate();
  
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const register = useCallback(async (credentials: RegisterCredentials) => {
    try {
      setIsLoading(true);
      setError(null);

      // 将密码转换为MD5
      const registerData: ServerRegisterRequest = {
        username: credentials.username,
        email: credentials.email,
        passwordMd5: hashPassword(credentials.password),
        emailVerificationCode: credentials.emailVerificationCode,
      };

      // 直接调用API，不通过store
      await authClient.register(registerData);

      // 注册成功后跳转到登录页面
      navigate('/login', { 
        replace: true,
        state: { 
          message: '注册成功，请登录',
          username: credentials.username 
        }
      });

    } catch (error: any) {
      // 使用新的错误处理方法提取用户友好的错误消息
      const errorMessage = extractErrorMessage(error);
      setError(errorMessage);
      
      console.error('Register failed:', error);
    } finally {
      setIsLoading(false);
    }
  }, [navigate]);

  return {
    isLoading,
    error,
    register,
    clearError,
  };
};
