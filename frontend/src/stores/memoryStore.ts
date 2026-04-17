import { create } from 'zustand';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { Memory, MemoryListQuery, MemoryStats, MemoryUpdateInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';

interface MemoryState {
  memories: Memory[];
  total: number;
  stats: MemoryStats | null;
  isLoading: boolean;
  query: MemoryListQuery;

  setQuery: (query: Partial<MemoryListQuery>) => void;
  fetchMemories: () => Promise<void>;
  fetchStats: () => Promise<void>;
  deleteMemory: (id: number) => Promise<void>;
  restoreMemory: (id: number) => Promise<void>;
  updateMemory: (id: number, input: MemoryUpdateInput) => Promise<Memory | null>;
}

export const useMemoryStore = create<MemoryState>()((set, get) => ({
  memories: [],
  total: 0,
  stats: null,
  isLoading: false,
  query: {
    offset: 0,
    limit: 20,
    keyword: '',
    type: '',
    is_forgotten: false,
  },

  setQuery: (partial) => {
    set((state) => ({
      query: { ...state.query, ...partial, offset: partial.offset ?? 0 },
    }));
  },

  fetchMemories: async () => {
    set({ isLoading: true });
    try {
      const resp = await Service.GetMemories(get().query);
      set({
        memories: resp?.memories ?? [],
        total: resp?.total ?? 0,
      });
    } finally {
      set({ isLoading: false });
    }
  },

  fetchStats: async () => {
    try {
      const stats = await Service.GetMemoryStats();
      set({ stats });
    } catch {
      // ignore
    }
  },

  deleteMemory: async (id) => {
    await Service.DeleteMemory(id);
    await get().fetchMemories();
    await get().fetchStats();
  },

  restoreMemory: async (id) => {
    await Service.RestoreMemory(id);
    await get().fetchMemories();
    await get().fetchStats();
  },

  updateMemory: async (id, input) => {
    const updated = await Service.UpdateMemory(id, new MemoryUpdateInput(input));
    await get().fetchMemories();
    return updated;
  },
}));
