// 环境变量配置模块

/**
 * API配置接口
 */
export interface ApiConfig {
  baseUrl: string;
  timeout: number;
  withCredentials: boolean;
}

/**
 * 获取API基础URL
 * @returns API基础URL
 */
export const getApiBaseUrl = (): string => {
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
};

/**
 * 获取API超时时间
 * @returns 超时时间（毫秒）
 */
export const getApiTimeout = (): number => {
  const timeout = import.meta.env.VITE_API_TIMEOUT;
  return timeout ? parseInt(timeout, 10) : 30000;
};

/**
 * 获取是否启用凭证传递
 * @returns 是否启用凭证传递
 */
export const getApiWithCredentials = (): boolean => {
  const withCredentials = import.meta.env.VITE_API_WITH_CREDENTIALS;
  return withCredentials === 'true' || withCredentials === true;
};

/**
 * 获取完整的API配置
 * @returns API配置对象
 */
export const getApiConfig = (): ApiConfig => {
  return {
    baseUrl: getApiBaseUrl(),
    timeout: getApiTimeout(),
    withCredentials: getApiWithCredentials(),
  };
};

/**
 * 环境变量配置
 */
export const env = {
  // API配置
  api: {
    baseUrl: getApiBaseUrl(),
    timeout: getApiTimeout(),
    withCredentials: getApiWithCredentials(),
  },
  
  // 开发模式
  isDev: import.meta.env.DEV,
  
  // 生产模式
  isProd: import.meta.env.PROD,
  
  // 模式
  mode: import.meta.env.MODE,
};

export default env;