import React from 'react';
import { CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useIsMobile } from '@/hooks/useViewportHeight';

interface SidebarHeaderProps {
    logoText: string;
    // 仅在移动端展开遮罩时用于关闭侧边栏；桌面端不显示关闭按钮
    onCloseMobileSidebar?: () => void;
}

const SidebarHeader: React.FC<SidebarHeaderProps> = ({
    logoText,
    onCloseMobileSidebar,
}) => {
    const { t } = useTranslation();
    const isMobile = useIsMobile();

    return (
        <div className="sidebar-header">
            <div className="sidebar-logo">
                <div className="logo-icon">🍋</div>
                <span className="logo-text">{logoText}</span>
            </div>
            {isMobile && onCloseMobileSidebar && (
                <button
                    className="collapse-btn"
                    onClick={onCloseMobileSidebar}
                    title={t('home.sidebar.collapse')}
                >
                    <CloseOutlined />
                </button>
            )}
        </div>
    );
};

export default SidebarHeader;
