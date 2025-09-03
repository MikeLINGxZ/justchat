import { useState, useEffect, useCallback } from 'react';

// 定义模型选项接口

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
  autoFetch?: boolean; // 是否自动获取，默认为true
}

// 默认模拟数据
const getDefaultModels = (): ModelOption[] => {
  return [
    {
      id: 'gpt-4o',
      name: 'GPT-4o',
      ownedBy: 'openai',
      enabled: true,
    },
    {
      id: 'gpt-4o-mini',
      name: 'GPT-4o Mini',
      ownedBy: 'openai', 
      enabled: true,
    },
    {
      id: 'gpt-4-turbo',
      name: 'GPT-4 Turbo',
      ownedBy: 'openai',
      enabled: true,
    },
    {
      id: 'gpt-3.5-turbo',
      name: 'GPT-3.5 Turbo',
      ownedBy: 'openai',
      enabled: true,
    },
    {
      id: 'claude-3-5-sonnet-20241022',
      name: 'Claude 3.5 Sonnet',
      ownedBy: 'anthropic',
      enabled: true,
    },
    {
      id: 'claude-3-5-haiku-20241022',
      name: 'Claude 3.5 Haiku',
      ownedBy: 'anthropic',
      enabled: true,
    },
    {
      id: 'claude-3-opus-20240229',
      name: 'Claude 3 Opus',
      ownedBy: 'anthropic',
      enabled: true,
    },
    {
      id: 'gemini-1.5-pro',
      name: 'Gemini 1.5 Pro',
      ownedBy: 'google',
      enabled: true,
    },
    {
      id: 'gemini-1.5-flash',
      name: 'Gemini 1.5 Flash',
      ownedBy: 'google',
      enabled: true,
    },
  ];
};

/**
 * 获取模型列表的自定义Hook
 * @param params 查询参数
 * @returns 模型列表状态和操作函数
 */
export const useModels = (params: UseModelsParams = {}): UseModelsReturn => {
  const { autoFetch = true } = params;
  
  const [models, setModels] = useState<ModelOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchModels = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // 模拟异步加载
      await new Promise(resolve => setTimeout(resolve, 300));
      
      const models = getDefaultModels();
      setModels(models);
    } catch (err: any) {
      const errorMessage = err?.message || '获取模型列表失败';
      setError(errorMessage);
      console.error('获取模型列表失败:', err);
    } finally {
      setIsLoading(false);
    }
  }, []);

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