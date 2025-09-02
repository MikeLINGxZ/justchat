import { Api } from './service/models/Api';
import { getApiConfig } from '@/config/env';
import { createGlobalErrorHandlerFetch } from '@/utils/globalErrorHandler';
import type {
  ServerModelsResponse,
  CommonModel
} from './service/models/Api';

// 创建带有全局错误处理的fetch函数
const fetchWithErrorHandler = createGlobalErrorHandlerFetch(fetch);

// 创建API实例
const api = new Api({
  baseUrl: getApiConfig().baseUrl,
  customFetch: fetchWithErrorHandler,
  securityWorker: (securityData) => {
    const token = localStorage.getItem('token');
    return token ? { headers: { Authorization: `Bearer ${token}` } } : {};
  },
});

// 模型客户端
export const modelsClient = {
  /**
   * 获取可用模型列表
   * @param params 查询参数
   * @returns 模型列表响应
   */
  async getModels(params?: {
    llmProviderId?: string;
    baseUrl?: string;
    apiKey?: string;
  }): Promise<ServerModelsResponse> {
    const response = await api.v1.modelsModels(params, { secure: true });
    return response.data;
  },
};

export type { ServerModelsResponse, CommonModel };
export default modelsClient;