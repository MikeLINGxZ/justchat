import React, { useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface SidebarHeaderProps {
    logoText: string;
    isSidebarCollapsed: boolean;
    onToggleSidebar: () => void;
}

const SidebarHeader: React.FC<SidebarHeaderProps> = ({
    logoText,
    isSidebarCollapsed,
    onToggleSidebar,
}) => {
    const { t } = useTranslation();
    const [isHovered, setIsHovered] = useState(false);

    const handleToggle = () => {
        onToggleSidebar();
        setIsHovered(false);
    };

    return (
        <div
            className="sidebar-header"
            onMouseOver={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
        >
            {!isSidebarCollapsed && (
                <div className="sidebar-logo">
                    <div className="logo-icon">🍋</div>
                    <span className="logo-text">{logoText}</span>
                </div>
            )}
            {isSidebarCollapsed && (
                <div className="sidebar-logo collapsed">
                    <div className="logo-icon" style={{ opacity: isHovered ? 0 : 1 }}>🍋</div>
                    {isHovered && (
                        <div className="expand-icon collapse-btn" onClick={handleToggle}>
                            <MenuUnfoldOutlined />
                        </div>
                    )}
                </div>
            )}
            {!isSidebarCollapsed && (
                <button className="collapse-btn" onClick={handleToggle} title={t('home.sidebar.collapse')}>
                    <MenuFoldOutlined />
                </button>
            )}
        </div>
    );
};

export default SidebarHeader;
