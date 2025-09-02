import { useAuthStore } from '../stores/authStore';

/**
 * 全局错误处理器 - 检测登录过期并处理
 * @param error 错误对象
 * @returns 是否处理了登录过期错误
 */
export const handleGlobalError = (error: any): boolean => {
  // 检查是否是登录过期错误
  if (isLoginExpiredError(error)) {
    handleLoginExpired();
    return true;
  }
  return false;
};

/**
 * 检查是否是登录过期错误
 * @param error 错误对象
 * @returns 是否是登录过期错误
 */
export const isLoginExpiredError = (error: any): boolean => {
  // 检查错误响应中的message字段
  if (error?.error?.message === 'ErrCodeLoginExpired') {
    return true;
  }
  
  // 检查错误响应中的data字段
  if (error?.data?.message === 'ErrCodeLoginExpired') {
    return true;
  }
  
  // 检查直接的message字段
  if (error?.message === 'ErrCodeLoginExpired') {
    return true;
  }
  
  // 检查RpcStatus格式的错误
  if (error?.code === 2 && error?.message === 'ErrCodeLoginExpired') {
    return true;
  }
  
  return false;
};

/**
 * 处理登录过期 - 清除登录状态并跳转到登录页面
 */
export const handleLoginExpired = (): void => {
  console.warn('检测到登录过期，正在清除登录状态...');
  
  // 清除localStorage中的token
  localStorage.removeItem('token');
  
  // 清除认证状态
  const authStore = useAuthStore.getState();
  authStore.setAuthState(null, null, false);
  
  // 跳转到登录页面
  // 使用window.location.href确保完全重新加载页面
  const currentPath = window.location.pathname;
  if (currentPath !== '/login') {
    window.location.href = '/login';
  }
};

/**
 * 创建带有全局错误处理的fetch包装器
 * @param originalFetch 原始的fetch函数
 * @returns 包装后的fetch函数
 */
export const createGlobalErrorHandlerFetch = (originalFetch: typeof fetch) => {
  return async (...args: Parameters<typeof fetch>): Promise<Response> => {
    try {
      const response = await originalFetch(...args);
      
      // 如果响应不成功，尝试解析错误
      if (!response.ok) {
        try {
          const errorData = await response.clone().json();
          if (handleGlobalError(errorData)) {
            // 如果处理了登录过期错误，抛出一个特殊的错误
            throw new Error('LOGIN_EXPIRED_HANDLED');
          }
        } catch (parseError) {
          // 如果解析失败，继续正常流程
        }
      }
      
      return response;
    } catch (error) {
      // 检查网络错误或其他错误
      handleGlobalError(error);
      throw error;
    }
  };
};