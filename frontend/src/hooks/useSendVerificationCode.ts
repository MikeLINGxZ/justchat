import { useState, useCallback, useRef, useEffect } from 'react';
import { authClient } from '@/api/authClient';
import { CommonVerificationCodeType } from '@/api/service/auth/Api';
import { extractErrorMessage } from '@/utils/errorHandler';

interface SendCodeCredentials {
  email: string;
  username: string;
  codeType?: CommonVerificationCodeType;
}

interface UseSendVerificationCodeReturn {
  isSending: boolean;
  codeSent: boolean;
  countdown: number;
  error: string | null;
  sendCode: (credentials: SendCodeCredentials) => Promise<void>;
  clearError: () => void;
  resetState: () => void;
}

export const useSendVerificationCode = (
  countdownSeconds: number = 60
): UseSendVerificationCodeReturn => {
  const [isSending, setIsSending] = useState(false);
  const [codeSent, setCodeSent] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  // 清理定时器
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    };
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const resetState = useCallback(() => {
    setIsSending(false);
    setCodeSent(false);
    setCountdown(0);
    setError(null);
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const startCountdown = useCallback(() => {
    setCodeSent(true);
    setCountdown(countdownSeconds);

    // 清除之前的定时器
    if (timerRef.current) {
      clearInterval(timerRef.current);
    }

    // 开始新的倒计时
    timerRef.current = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          if (timerRef.current) {
            clearInterval(timerRef.current);
            timerRef.current = null;
          }
          setCodeSent(false);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }, [countdownSeconds]);

  const sendCode = useCallback(async (credentials: SendCodeCredentials) => {
    try {
      setIsSending(true);
      setError(null);

      await authClient.sendEmailVerificationCode({
        email: credentials.email,
        username: credentials.username,
        codeType: credentials.codeType || CommonVerificationCodeType.VERIFICATION_CODE_TYPE_REGISTER,
      });

      // 发送成功，开始倒计时
      startCountdown();

    } catch (error: any) {
      // 使用新的错误处理方法提取用户友好的错误消息
      const errorMessage = extractErrorMessage(error);
      setError(errorMessage);
      
      console.error('Send verification code failed:', error);
    } finally {
      setIsSending(false);
    }
  }, [startCountdown]);

  return {
    isSending,
    codeSent,
    countdown,
    error,
    sendCode,
    clearError,
    resetState,
  };
};
