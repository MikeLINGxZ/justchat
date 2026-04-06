import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface OPCPersonView {
    id: number;
    uuid: string;
    name: string;
    role: string;
    agent_id: string;
    avatar: string;
    is_pinned: boolean;
    chat_uuid: string;
    last_message: string | null;
    last_message_at: string | null;
    created_at: string;
    updated_at: string;
}

export interface OPCGroupView {
    id: number;
    uuid: string;
    chat_uuid: string;
    name: string;
    description: string;
    is_pinned: boolean;
    members: OPCPersonView[];
    last_message: string | null;
    last_message_at: string | null;
    created_at: string;
    updated_at: string;
}

export type AppMode = 'chat' | 'opc';

export type OPCSidebarTab = 'conversations' | 'contacts';

export interface OPCState {
    mode: AppMode;
    persons: OPCPersonView[];
    groups: OPCGroupView[];
    selectedType: 'person' | 'group' | 'contact' | null;
    selectedUuid: string | null;
    searchQuery: string;
    isLoading: boolean;
    sidebarTab: OPCSidebarTab;

    // Actions
    setMode: (mode: AppMode) => void;
    setSelected: (type: 'person' | 'group' | 'contact', uuid: string) => void;
    clearSelected: () => void;
    setPersons: (persons: OPCPersonView[]) => void;
    setGroups: (groups: OPCGroupView[]) => void;
    setSearchQuery: (query: string) => void;
    setLoading: (loading: boolean) => void;
    setSidebarTab: (tab: OPCSidebarTab) => void;
}

export const useOPCStore = create<OPCState>()(
    persist(
        (set) => ({
            mode: 'chat',
            persons: [],
            groups: [],
            selectedType: null,
            selectedUuid: null,
            searchQuery: '',
            isLoading: false,
            sidebarTab: 'conversations',

            setMode: (mode) => set({ mode }),
            setSelected: (type, uuid) => set({ selectedType: type, selectedUuid: uuid }),
            clearSelected: () => set({ selectedType: null, selectedUuid: null }),
            setPersons: (persons) => set({ persons }),
            setGroups: (groups) => set({ groups }),
            setSearchQuery: (query) => set({ searchQuery: query }),
            setLoading: (loading) => set({ isLoading: loading }),
            setSidebarTab: (tab) => set({ sidebarTab: tab, searchQuery: '' }),
        }),
        {
            name: 'opc-store',
            partialize: (state) => ({
                mode: state.mode,
            }),
        }
    )
);

export const initializeOPC = () => {
    // OPC store 通过 persist middleware 自动从 localStorage 恢复 mode
    // 无需额外初始化
};
