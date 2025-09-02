// RPC 客户端工厂 (Axios 版)
// 此文件由 generate_rpc_ts_axios.sh 脚本生成

import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { getApiConfig } from '@/config/env';

// 客户端配置
export interface ClientConfig {
  baseURL?: string;
  timeout?: number;
  headers?: Record<string, string>;
  withCredentials?: boolean;
}

// 获取默认配置（从环境变量）
const getDefaultConfig = (): ClientConfig => {
  const envConfig = getApiConfig();
  return {
    baseURL: envConfig.baseUrl,
    timeout: envConfig.timeout,
    headers: {
      'Content-Type': 'application/json',
    },
    withCredentials: envConfig.withCredentials,
  };
};

// 默认配置
const DEFAULT_CONFIG: ClientConfig = getDefaultConfig();

// 响应数据接口
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 客户端工厂类
export class RpcClientFactory {
  private axiosInstance: AxiosInstance;
  private config: ClientConfig;

  constructor(config: ClientConfig = {}) {
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.axiosInstance = axios.create(this.config);
    
    // 请求拦截器
    this.axiosInstance.interceptors.request.use(
      (config) => {
        // 可以在这里添加认证 token
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse<ApiResponse>) => {
        return response;
      },
      (error) => {
        // 处理错误响应
        if (error.response?.status === 401) {
          // 未授权，清除 token 并跳转到登录页
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // 获取配置
  getConfig(): ClientConfig {
    return { ...this.config };
  }

  // 更新配置
  updateConfig(newConfig: Partial<ClientConfig>) {
    this.config = { ...this.config, ...newConfig };
    this.axiosInstance.defaults.baseURL = this.config.baseURL;
    this.axiosInstance.defaults.timeout = this.config.timeout;
    if (this.config.headers) {
      Object.assign(this.axiosInstance.defaults.headers, this.config.headers);
    }
  }

  // 通用请求方法
  async request<T>(config: AxiosRequestConfig): Promise<T> {
    try {
      const response = await this.axiosInstance.request<ApiResponse<T>>(config);
      return response.data.data;
    } catch (error) {
      throw error;
    }
  }

  // 获取认证服务客户端
  getAuthClient() {
    return {
      // 发送邮件验证码
      sendEmailVerificationCode: (data: {
        email: string;
        username: string;
        code_type: number;
      }) =>
        this.request({
          method: 'POST',
          url: '/v1/auth/send-email-code',
          data,
        }),

      // 检查字段可用性
      checkFieldAvailability: (data: {
        field_type: number;
        value: string;
      }) =>
        this.request<{ available: boolean }>({
          method: 'POST',
          url: '/v1/auth/check-availability',
          data,
        }),

      // 注册
      register: (data: {
        username: string;
        password_md5: string;
        email: string;
        email_verification_code: number;
      }) =>
        this.request({
          method: 'POST',
          url: '/v1/auth/register',
          data,
        }),

      // 登录
      login: (data: {
        username: string;
        password_md5: string;
      }) =>
        this.request<{
          token: string;
          user: {
            id: string;
            username: string;
            email: string;
          };
        }>({
          method: 'POST',
          url: '/v1/auth/login',
          data,
        }),

      // 重置密码
      resetPassword: (data: {
        email: string;
        email_verification_code: number;
        new_password_md5: string;
      }) =>
        this.request({
          method: 'POST',
          url: '/v1/auth/reset-password',
          data,
        }),

      // 登出
      logout: () =>
        this.request({
          method: 'POST',
          url: '/v1/auth/logout',
        }),
    };
  }

  // 获取聊天服务客户端
  getChatClient() {
    return {
      // 创建聊天
      createChat: (data: {
        title: string;
        model_id: string;
      }) =>
        this.request<{
          chat_id: string;
          title: string;
          model_id: string;
        }>({
          method: 'POST',
          url: '/v1/chat/create',
          data,
        }),

      // 发送消息
      sendMessage: (data: {
        chat_id: string;
        content: string;
      }) =>
        this.request<{
          message_id: string;
          content: string;
          role: string;
          timestamp: string;
        }>({
          method: 'POST',
          url: '/v1/chat/send-message',
          data,
        }),

      // 获取聊天历史
      getChatHistory: (chatId: string) =>
        this.request<Array<{
          message_id: string;
          content: string;
          role: string;
          timestamp: string;
        }>>({
          method: 'GET',
          url: `/v1/chat/${chatId}/history`,
        }),

      // 获取聊天列表
      getChatList: () =>
        this.request<Array<{
          chat_id: string;
          title: string;
          model_id: string;
          created_at: string;
        }>>({
          method: 'GET',
          url: '/v1/chat/list',
        }),

      // 删除聊天
      deleteChat: (chatId: string) =>
        this.request({
          method: 'DELETE',
          url: `/v1/chat/${chatId}`,
        }),
    };
  }

  // 获取模型服务客户端
  getModelsClient() {
    return {
      // 获取模型列表
      getModelList: () =>
        this.request<Array<{
          model_id: string;
          name: string;
          description: string;
          provider: string;
        }>>({
          method: 'GET',
          url: '/v1/models/list',
        }),

      // 获取模型详情
      getModelDetail: (modelId: string) =>
        this.request<{
          model_id: string;
          name: string;
          description: string;
          provider: string;
          parameters: Record<string, any>;
        }>({
          method: 'GET',
          url: `/v1/models/${modelId}`,
        }),
    };
  }
}

// 默认客户端实例
export const defaultRpcClient = new RpcClientFactory();

// 导出便捷方法
export const getAuthClient = () => defaultRpcClient.getAuthClient();
export const getChatClient = () => defaultRpcClient.getChatClient();
export const getModelsClient = () => defaultRpcClient.getModelsClient();

// 创建新的客户端实例
export const createRpcClient = (config: ClientConfig) => new RpcClientFactory(config);

// 导出 axios 实例（如果需要直接使用）
export const getAxiosInstance = () => defaultRpcClient['axiosInstance'];
