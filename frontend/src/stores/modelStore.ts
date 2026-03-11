import {Model} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts';
import { create } from 'zustand';
import {persist} from "zustand/middleware";

// 定义模型选项接口
export interface ModelOption {
    id: number;
    name: string;
    ownedBy?: string;
    enabled?: boolean;
}

// 将后端模型数据转换为前端使用的格式
const convertBackendModel = (backendModel: Model): ModelOption => {
    return {
        id: backendModel.id,
        name: backendModel.model ,
        ownedBy: backendModel.owned_by,
        enabled: backendModel.enable,
    };
};

// 默认模拟数据
const getDefaultModels = (): ModelOption[] => {
    return [];
};

export interface UseModelsReturn {
    models: ModelOption[];
    isLoading: boolean;
    error: string | null;
    refetch: () => Promise<void>;
}

export const useModelStore = create<UseModelsReturn>()(
    persist(
        (set, get) => ({
            models: [],
            isLoading: false,
            error: null,
            refetch: async () => {
                set({ isLoading: true, error: null });
                
                try {
                    // 调用 Wails 后端服务获取模型列表
                    const backendModels = await Service.GetModels(true, true);
                    
                    // 转换后端数据格式为前端使用的格式
                    const convertedModels = backendModels.map(convertBackendModel);
                    
                    // 如果后端没有数据，使用默认模拟数据作为后备
                    if (convertedModels.length === 0) {
                        console.warn('后端返回空模型列表，使用默认模拟数据');
                        set({ models: getDefaultModels(), isLoading: false });
                    } else {
                        set({ models: convertedModels, isLoading: false });
                    }
                } catch (err: any) {
                    const errorMessage = err?.message || '获取模型列表失败';
                    console.error('从后端获取模型列表失败，使用默认模拟数据:', err);
                    
                    // 出错时使用默认模拟数据作为后备
                    set({ 
                        error: errorMessage, 
                        models: getDefaultModels(), 
                        isLoading: false 
                    });
                }
            },
        }),
        {
            name: 'model-store', // 持久化存储的键名
            partialize: (state) => ({ models: state.models }), // 只持久化 models 数据
        }
    )
)