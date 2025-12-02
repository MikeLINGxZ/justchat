import { useState, useEffect, useCallback } from 'react';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts';
import { Model } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/index.ts';

// 定义模型选项接口
export interface ModelOption {
  id: string;
  name: string;
  ownedBy?: string;
  enabled?: boolean;
}

// 将后端模型数据转换为前端使用的格式
const convertBackendModel = (backendModel: Model): ModelOption => {
  return {
    id: backendModel.model,
    name: backendModel.alias || backendModel.model,
    ownedBy: backendModel.owned_by,
    enabled: backendModel.enable,
  };
};

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
  return [];
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
      // 调用 Wails 后端服务获取模型列表
      const backendModels = await Service.GetModels(true,true);
      
      // 转换后端数据格式为前端使用的格式
      const convertedModels = backendModels.map(convertBackendModel);
      
      // 如果后端没有数据，使用默认模拟数据作为后备
      if (convertedModels.length === 0) {
        console.warn('后端返回空模型列表，使用默认模拟数据');
        setModels(getDefaultModels());
      } else {
        setModels(convertedModels);
      }
    } catch (err: any) {
      const errorMessage = err?.message || '获取模型列表失败';
      setError(errorMessage);
      console.error('从后端获取模型列表失败，使用默认模拟数据:', err);
      
      // 出错时使用默认模拟数据作为后备
      setModels(getDefaultModels());
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