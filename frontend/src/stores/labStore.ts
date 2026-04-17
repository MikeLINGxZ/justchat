import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface EmbeddingConfig {
  provider: string;
  baseUrl: string;
  apiKey: string;
  model: string;
}

interface LabState {
  pluginSystemEnabled: boolean;
  setPluginSystemEnabled: (enabled: boolean) => void;
  memorySystemEnabled: boolean;
  setMemorySystemEnabled: (enabled: boolean) => void;
  vectorSearchEnabled: boolean;
  setVectorSearchEnabled: (enabled: boolean) => void;
  embeddingConfig: EmbeddingConfig;
  setEmbeddingConfig: (config: Partial<EmbeddingConfig>) => void;
}

export const useLabStore = create<LabState>()(
  persist(
    (set) => ({
      pluginSystemEnabled: false,
      setPluginSystemEnabled: (enabled: boolean) => set({ pluginSystemEnabled: enabled }),
      memorySystemEnabled: false,
      setMemorySystemEnabled: (enabled: boolean) => set({ memorySystemEnabled: enabled }),
      vectorSearchEnabled: false,
      setVectorSearchEnabled: (enabled: boolean) => set({ vectorSearchEnabled: enabled }),
      embeddingConfig: {
        provider: 'ollama',
        baseUrl: 'http://localhost:11434',
        apiKey: '',
        model: 'bge-m3',
      },
      setEmbeddingConfig: (config) =>
        set((state) => ({
          embeddingConfig: { ...state.embeddingConfig, ...config },
        })),
    }),
    {
      name: 'lab-settings',
    }
  )
);
