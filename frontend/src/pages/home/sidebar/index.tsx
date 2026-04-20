import React, {useState} from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useOPCStore } from '@/stores/opcStore';
import SidebarHeader from '@/components/sidebar/SidebarHeader';
import SidebarUserMenu from '@/components/sidebar/SidebarUserMenu';
import SidebarChats from '@/pages/home/sidebar/chat_lists.tsx';
import '@/pages/home/sidebar/index.scss';

interface SidebarProps {
    className?: string;
    currentChatUuid: string | null;
    onChatSelect: (chatUuid: string) => void;
    isSidebarCollapsed: boolean;
    onToggleSidebar: () => void;
    onResizeStart?: (event: React.MouseEvent<HTMLDivElement>) => void;
    onNewChat?: () => void;
    onRegisterRefreshCallback?: (callback: () => void) => void;
    onDeleteChat?: (chatUuid: string) => void;
    generatingChatUuids?: string[];
    onStopGenerationForChat?: (chatUuid: string) => void;
}

const Index: React.FC<SidebarProps> = ({
                                           className,
                                           currentChatUuid,
                                           onChatSelect,
                                           isSidebarCollapsed,
                                           onToggleSidebar,
                                           onResizeStart,
                                           onNewChat,
                                           onRegisterRefreshCallback,
                                           onDeleteChat,
                                           generatingChatUuids,
                                       onStopGenerationForChat,
                                   }) => {
    const { t } = useTranslation();
    const { mode } = useOPCStore();
    const [activeTab, setActiveTab] = useState<'history' | 'favorites'>('history');

    const handleNewChat = () => {
        onNewChat?.();
    };

    const handleChatSelect = (chatUuid: string) => {
        onChatSelect(chatUuid);
    };

    return (
        <div
            className={`sidebar ${isSidebarCollapsed ? 'collapsed' : ''} ${className || ''}`}
        >
            <SidebarHeader
                logoText="LemonTea"
                onCloseMobileSidebar={onToggleSidebar}
            />

            {/* 功能按钮区域 */}
            <div className="sidebar-actions">
                <button className="action-btn" onClick={handleNewChat} title={t('home.sidebar.newChat')}>
                    <PlusOutlined className="action-icon"/>
                    <span className="action-text">{t('home.sidebar.newChat')}</span>
                </button>
            </div>

            {/* Tab切换按钮区域 */}
            <div className="sidebar-tabs">
                <div className="tab-switch">
                    <div
                        className={`tab-option ${activeTab === 'history' ? 'active' : ''}`}
                        onClick={() => setActiveTab('history')}
                    >
                        {t('home.sidebar.tabs.history')}
                    </div>
                    <div
                        className={`tab-option ${activeTab === 'favorites' ? 'active' : ''}`}
                        onClick={() => setActiveTab('favorites')}
                    >
                        {t('home.sidebar.tabs.favorites')}
                    </div>
                    <div className="tab-slider" data-active={activeTab}></div>
                </div>
            </div>

            {/* 主体区域 - 历史对话列表 */}
            <div className="sidebar-main">
                <SidebarChats
                    currentChatUuid={currentChatUuid}
                    onChatSelect={handleChatSelect}
                    onRegisterRefreshCallback={onRegisterRefreshCallback}
                    onDeleteChat={onDeleteChat}
                    activeTab={activeTab}
                    generatingChatUuids={generatingChatUuids}
                    onStopGenerationForChat={onStopGenerationForChat}
                />
            </div>

            {/* 占位元素，确保底部用户信息固定在底部 */}
            <div className="sidebar-spacer"></div>

            <SidebarUserMenu currentMode={mode} />

            {onResizeStart && (
                <div
                    className="sidebar-resize-handle"
                    onMouseDown={onResizeStart}
                    title={t('home.sidebar.resizeHandle', 'Drag to resize')}
                    role="separator"
                    aria-orientation="vertical"
                />
            )}
        </div>
    );
};

export default Index;
