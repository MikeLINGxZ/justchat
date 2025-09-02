import { useState, useEffect, useCallback } from 'react';
import { modelsClient } from '@/api/modelsClient';
import type { CommonModel } from '@/api/modelsClient';

export interface ModelOption {
  id: string;
  name: string;
  ownedBy?: string;
  enabled?: boolean;
}

export interface UseModelsReturn {
  models: ModelOption[];
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export interface UseModelsParams {
  llmProviderId?: string;
  baseUrl?: string;
  apiKey?: string;
  autoFetch?: boolean; // 是否自动获取，默认为true
}

/**
 * 获取模型列表的自定义Hook
 * @param params 查询参数
 * @returns 模型列表状态和操作函数
 */
export const useModels = (params: UseModelsParams = {}): UseModelsReturn => {
  const { llmProviderId, baseUrl, apiKey, autoFetch = true } = params;
  
  const [models, setModels] = useState<ModelOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchModels = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await modelsClient.getModels({
        llmProviderId,
        baseUrl,
        apiKey,
      });
      
      // 转换API响应为组件需要的格式
      const modelOptions: ModelOption[] = (response.data || []).map((model: CommonModel) => ({
        id: model.id || '',
        name: model.id || '', // 使用id作为显示名称，如果有其他字段可以调整
        ownedBy: model.ownedBy,
        enabled: model.enabled,
      }));
      
      setModels(modelOptions);
    } catch (err: any) {
      const errorMessage = err?.message || '获取模型列表失败';
      setError(errorMessage);
      console.error('获取模型列表失败:', err);
    } finally {
      setIsLoading(false);
    }
  }, [llmProviderId, baseUrl, apiKey]);

  // 自动获取模型列表
  useEffect(() => {
    if (autoFetch) {
      fetchModels();
    }
  }, [fetchModels, autoFetch]);

  return {
    models,
    isLoading,
    error,
    refetch: fetchModels,
  };
};

export default useModels;