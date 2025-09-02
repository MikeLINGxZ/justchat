/**
 * 后端错误响应的接口定义
 */
interface ApiError {
  code: number;
  message: string;
  details: any[];
}

/**
 * 错误代码到用户友好消息的映射
 */
const ERROR_MESSAGE_MAP: Record<string, string> = {
  // 认证相关错误
  'ErrCodeInvalidAccountPassword': '用户名或密码错误',
  'ErrCodeAccountNotFound': '账户不存在',
  'ErrCodeAccountDisabled': '账户已被禁用',
  'ErrCodePasswordExpired': '密码已过期，请重置密码',
  'ErrCodeAccountLocked': '账户已被锁定，请联系管理员',
  'ErrCodeInvalidToken': '登录已过期，请重新登录',
  'ErrCodeTokenExpired': '登录已过期，请重新登录',
  
  // 注册相关错误
  'ErrCodeUsernameExists': '用户名已存在',
  'ErrCodeEmailExists': '邮箱已被注册',
  'ErrCodeInvalidEmail': '邮箱格式不正确',
  'ErrCodeWeakPassword': '密码强度不够',
  'ErrCodeInvalidVerificationCode': '验证码错误或已过期',
  
  // 权限相关错误
  'ErrCodePermissionDenied': '权限不足',
  'ErrCodeUnauthorized': '未授权访问',
  'ErrCodeForbidden': '禁止访问',
  
  // 服务器相关错误
  'ErrCodeInternalServerError': '服务器内部错误，请稍后重试',
  'ErrCodeServiceUnavailable': '服务暂时不可用，请稍后重试',
  'ErrCodeTimeout': '请求超时，请检查网络连接',
  'ErrCodeTooManyRequests': '请求过于频繁，请稍后重试',
  
  // 数据相关错误
  'ErrCodeInvalidInput': '输入数据格式不正确',
  'ErrCodeDataNotFound': '数据不存在',
  'ErrCodeDataConflict': '数据冲突，请刷新后重试',
  'ErrCodeDataCorrupted': '数据已损坏',
  
  // 通用错误
  'ErrCodeUnknown': '未知错误，请联系技术支持',
  'ErrCodeBadRequest': '请求参数错误',
  'ErrCodeNotFound': '请求的资源不存在',
};

/**
 * 默认错误消息
 */
const DEFAULT_ERROR_MESSAGE = '操作失败，请稍后重试';

/**
 * 提取并转换API错误为用户友好的消息
 * @param error - 错误对象，可能是AxiosError或其他类型
 * @returns 用户友好的错误消息
 */
export const extractErrorMessage = (error: any): string => {

  try {
    // 如果是Axios错误，提取response.data
    let apiError: ApiError | null = null;
    if (error?.error?.message != null) {
      apiError = error.error.message;
    } else {
      apiError = error
    }
    
    // 如果成功提取到API错误，尝试转换消息
    if (apiError && typeof apiError === 'string') {
      const friendlyMessage = ERROR_MESSAGE_MAP[apiError];
      if (friendlyMessage) {
        return friendlyMessage;
      }
      
      // 如果message看起来是人类可读的，直接返回
      return apiError;
    }

    // 处理常见的HTTP状态码错误
    if (error?.response?.status) {
      switch (error.response.status) {
        case 400:
          return '请求参数错误';
        case 401:
          return '未授权，请重新登录';
        case 403:
          return '权限不足';
        case 404:
          return '请求的资源不存在';
        case 429:
          return '请求过于频繁，请稍后重试';
        case 500:
          return '服务器内部错误，请稍后重试';
        case 502:
        case 503:
          return '服务暂时不可用，请稍后重试';
        case 504:
          return '请求超时，请检查网络连接';
        default:
          return DEFAULT_ERROR_MESSAGE;
      }
    }

    // 处理网络错误
    if (error?.message) {
      if (error.message.includes('Network Error') || error.message.includes('timeout')) {
        return '网络连接异常，请检查网络后重试';
      }
      
      if (error.message.includes('ECONNREFUSED') || error.message.includes('ERR_CONNECTION_REFUSED')) {
        return '无法连接到服务器，请稍后重试';
      }
    }

    // 兜底返回默认错误消息
    return apiError?.message || DEFAULT_ERROR_MESSAGE;

  } catch (e) {
    // 如果错误处理本身出错，返回默认消息
    console.error('Error while extracting error message:', e);
    return DEFAULT_ERROR_MESSAGE;
  }
};

/**
 * 添加新的错误消息映射
 * @param errorCode - 错误代码
 * @param message - 用户友好的消息
 */
export const addErrorMapping = (errorCode: string, message: string): void => {
  ERROR_MESSAGE_MAP[errorCode] = message;
};

/**
 * 获取当前的错误映射（用于调试）
 */
export const getErrorMappings = (): Record<string, string> => {
  return { ...ERROR_MESSAGE_MAP };
};

/**
 * 检查是否是API错误
 * @param error - 错误对象
 * @returns 是否是API错误
 */
export const isApiError = (error: any): error is ApiError => {
  return error && 
         typeof error.code === 'number' && 
         typeof error.message === 'string' && 
         Array.isArray(error.details);
};

/**
 * 从错误中提取错误代码
 * @param error - 错误对象
 * @returns 错误代码，如果没有则返回null
 */
export const extractErrorCode = (error: any): number | null => {
  try {
    if (error?.response?.data?.code !== undefined) {
      return error.response.data.code;
    }
    if (error?.data?.code !== undefined) {
      return error.data.code;
    }
    if (error?.code !== undefined) {
      return error.code;
    }
    return null;
  } catch {
    return null;
  }
};
