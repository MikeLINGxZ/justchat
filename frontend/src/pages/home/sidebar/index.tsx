import React, {useEffect, useRef, useState} from 'react';
import {
    BulbOutlined,
    CheckOutlined,
    LogoutOutlined,
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    MoonOutlined,
    PlusOutlined,
    SettingOutlined,
    SunOutlined,
    UserOutlined,
} from '@ant-design/icons';
import {useAuthStore} from '@/stores/authStore.ts';
import SidebarChats from '@/pages/home/sidebar/chat_lists.tsx';
import '@/pages/home/sidebar/index.scss';

interface SidebarProps {
    className?: string;
    currentChatUuid: string | null;
    onChatSelect: (chatUuid: string, chatTitle?: string) => void;
    isSidebarCollapsed: boolean;
    onToggleSidebar: () => void;
    onNewChat?: () => void;
    onRegisterRefreshCallback?: (callback: () => void) => void;
    onRegisterUpdateTitleCallback?: (
        callback: (chatUuid: string, newTitle: string) => void
    ) => void;
    onDeleteChat?: (chatUuid: string) => void;
}

const Index: React.FC<SidebarProps> = ({
                                           className,
                                           currentChatUuid,
                                           onChatSelect,
                                           isSidebarCollapsed,
                                           onToggleSidebar,
                                           onNewChat,
                                           onRegisterRefreshCallback,
                                           onRegisterUpdateTitleCallback,
                                           onDeleteChat,
                                       }) => {
    // @ts-ignore
    const {user, logout} = useAuthStore();
    const [isHeaderHovered, setIsHeaderHovered] = useState(false);
    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
    const [isThemeMenuOpen, setIsThemeMenuOpen] = useState(false);
    const [currentTheme, setCurrentTheme] = useState<'auto' | 'light' | 'dark'>(
        'auto'
    );
    const [activeTab, setActiveTab] = useState<'history' | 'favorites'>('history'); // 添加tab状态
    const userMenuRef = useRef<HTMLDivElement>(null);
    const themeMenuRef = useRef<HTMLDivElement>(null);
    const themeCloseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const isMobileRef = useRef(window.innerWidth <= 768);
    const wasAutoCollapsedRef = useRef(false); // 标记是否是因为移动端而自动折叠

    // 检测移动端
    useEffect(() => {
        const checkMobile = () => {
            const isCurrentlyMobile = window.innerWidth <= 768;
            const wasMobile = isMobileRef.current;
            
            // 检测从桌面端切换到移动端
            if (!wasMobile && isCurrentlyMobile) {
                // 从桌面端切换到移动端，自动折叠侧边栏
                if (!isSidebarCollapsed) {
                    onToggleSidebar();
                    wasAutoCollapsedRef.current = true; // 标记为自动折叠
                }
            } 
            // 检测从移动端切换到桌面端
            else if (wasMobile && !isCurrentlyMobile) {
                // 从移动端切换到桌面端，如果是因为移动端而折叠的，则自动展开
                if (isSidebarCollapsed && wasAutoCollapsedRef.current) {
                    onToggleSidebar();
                    wasAutoCollapsedRef.current = false; // 重置标记
                }
            }
            
            // 更新当前状态
            isMobileRef.current = isCurrentlyMobile;
        };

        // 初始化时检查是否为移动端
        checkMobile();
        
        window.addEventListener('resize', checkMobile);

        return () => window.removeEventListener('resize', checkMobile);
    }, [isSidebarCollapsed, onToggleSidebar]);

    // 点击外部关闭菜单
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (
                userMenuRef.current &&
                !userMenuRef.current.contains(event.target as Node)
            ) {
                setIsUserMenuOpen(false);
            }
            if (
                themeMenuRef.current &&
                !themeMenuRef.current.contains(event.target as Node)
            ) {
                setIsThemeMenuOpen(false);
            }
        };

        if (isUserMenuOpen || isThemeMenuOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isUserMenuOpen, isThemeMenuOpen]);

    // 从 localStorage 读取主题设置
    useEffect(() => {
        const savedTheme = localStorage.getItem('theme') as
            | 'auto'
            | 'light'
            | 'dark';
        if (savedTheme) {
            setCurrentTheme(savedTheme);
        }
    }, []);

    // 清理定时器
    useEffect(() => {
        return () => {
            if (themeCloseTimeoutRef.current) {
                clearTimeout(themeCloseTimeoutRef.current);
            }
        };
    }, []);

    // 处理新建对话
    const handleNewChat = () => {
        onNewChat?.();
    };

    const handleChatSelect = (chatUuid: string, chatTitle?: string) => {
        onChatSelect(chatUuid, chatTitle);
    };

    const handleUserMenuToggle = () => {
        setIsUserMenuOpen(!isUserMenuOpen);
        setIsThemeMenuOpen(false);
    };

    // 用户手动点击折叠/展开按钮时，重置自动折叠标记
    const handleManualToggle = () => {
        onToggleSidebar();
        wasAutoCollapsedRef.current = false; // 用户手动操作后重置标记
        setIsHeaderHovered(false);
    };

    const handleThemeChange = (theme: 'auto' | 'light' | 'dark') => {
        setCurrentTheme(theme);
        localStorage.setItem('theme', theme);

        // 应用主题
        const root = document.documentElement;
        if (theme === 'dark') {
            root.classList.add('dark');
            root.classList.remove('light');
        } else if (theme === 'light') {
            root.classList.add('light');
            root.classList.remove('dark');
        } else {
            // auto 模式，根据系统偏好设置
            root.classList.remove('dark', 'light');
            const prefersDark = window.matchMedia(
                '(prefers-color-scheme: dark)'
            ).matches;
            if (prefersDark) {
                root.classList.add('dark');
            } else {
                root.classList.add('light');
            }
        }

        // 清除定时器
        if (themeCloseTimeoutRef.current) {
            clearTimeout(themeCloseTimeoutRef.current);
            themeCloseTimeoutRef.current = null;
        }

        // 关闭所有菜单
        setIsThemeMenuOpen(false);
        setIsUserMenuOpen(false);
    };

    const handleLogout = async () => {
        try {
            await logout();
            setIsUserMenuOpen(false);
        } catch (error) {
            console.error('退出登录失败:', error);
        }
    };

    const getThemeIcon = (theme: 'auto' | 'light' | 'dark') => {
        switch (theme) {
            case 'light':
                return <SunOutlined/>;
            case 'dark':
                return <MoonOutlined/>;
            default:
                return <BulbOutlined/>;
        }
    };

    const getThemeText = (theme: 'auto' | 'light' | 'dark') => {
        switch (theme) {
            case 'light':
                return '浅色';
            case 'dark':
                return '深色(beta)';
            default:
                return '自动(beta)';
        }
    };

    const handleThemeMenuEnter = () => {
        if (themeCloseTimeoutRef.current) {
            clearTimeout(themeCloseTimeoutRef.current);
            themeCloseTimeoutRef.current = null;
        }
        setIsThemeMenuOpen(true);
    };

    const handleThemeMenuLeave = () => {
        themeCloseTimeoutRef.current = setTimeout(() => {
            setIsThemeMenuOpen(false);
        }, 150); // 150ms 延迟关闭
    };

    return (
        <div
            className={`sidebar ${isSidebarCollapsed ? 'collapsed' : ''} ${className || ''}`}
        >
            {/* 顶部区域 */}
            <div
                className="sidebar-header"
                onMouseOver={() => setIsHeaderHovered(true)}
                onMouseLeave={() => setIsHeaderHovered(false)}
            >
                {!isSidebarCollapsed && (
                    <div className="sidebar-logo">
                        <div className="logo-icon">🍋</div>
                        <span className="logo-text">LemonTea</span>
                    </div>
                )}
                {isSidebarCollapsed && (
                    <div className="sidebar-logo collapsed">
                        <div
                            className="logo-icon"
                            style={{opacity: isHeaderHovered ? 0 : 1}}
                        >
                            🍋
                        </div>
                        {isHeaderHovered && (
                            <div className="expand-icon collapse-btn" onClick={handleManualToggle}>
                                <MenuUnfoldOutlined/>
                            </div>
                        )}
                    </div>
                )}
                {!isSidebarCollapsed && (
                    <button
                        className="collapse-btn"
                        onClick={handleManualToggle}
                        title="折叠侧边栏"
                    >
                        <MenuFoldOutlined/>
                    </button>
                )}
            </div>

            {/* 功能按钮区域 */}
            <div className="sidebar-actions">
                <button className="action-btn" onClick={handleNewChat} title="新建对话">
                    <PlusOutlined className="action-icon"/>
                    {!isSidebarCollapsed && <span className="action-text">新建对话</span>}
                </button>
            </div>

            {/* Tab切换按钮区域 */}
            {!isSidebarCollapsed && (
                <div className="sidebar-tabs">
                    <div className="tab-switch">
                        <div 
                            className={`tab-option ${activeTab === 'history' ? 'active' : ''}`}
                            onClick={() => setActiveTab('history')}
                        >
                            对话
                        </div>
                        <div 
                            className={`tab-option ${activeTab === 'favorites' ? 'active' : ''}`}
                            onClick={() => setActiveTab('favorites')}
                        >
                            收藏
                        </div>
                        <div className="tab-slider" data-active={activeTab}></div>
                    </div>
                </div>
            )}

            {/* 主体区域 - 历史对话列表 */}
            {!isSidebarCollapsed && (
                <div className="sidebar-main">
                    <SidebarChats
                        currentChatUuid={currentChatUuid}
                        onChatSelect={handleChatSelect}
                        onRegisterRefreshCallback={onRegisterRefreshCallback}
                        onRegisterUpdateTitleCallback={onRegisterUpdateTitleCallback}
                        onDeleteChat={onDeleteChat}
                        activeTab={activeTab} // 传递tab状态
                    />
                </div>
            )}

            {/* 占位元素，确保底部用户信息固定在底部 */}
            <div className="sidebar-spacer"></div>

            {/* 底部区域 - 用户信息 */}
            <div className="sidebar-footer" ref={userMenuRef}>
                {!isSidebarCollapsed && (
                    <>
                        <div
                            className={`user-section ${isUserMenuOpen ? 'active' : ''}`}
                            onClick={handleUserMenuToggle}
                        >
                            <div className="user-avatar">
                                {user?.avatar ? (
                                    <img
                                        src={user.avatar}
                                        alt="用户头像"
                                        className="avatar-img"
                                    />
                                ) : (
                                    <UserOutlined className="avatar-icon"/>
                                )}
                            </div>
                            <div className="user-info">
                                <div className="user-name">
                                    {user?.username || '未登录用户'}
                                </div>
                                {user?.email && <div className="user-email">{user.email}</div>}
                            </div>
                            <div className="user-menu-icon">
                                <SettingOutlined/>
                            </div>
                        </div>

                        {/* 用户菜单 */}
                        {isUserMenuOpen && (
                            <div className="user-menu">
                                <div
                                    className="menu-item theme-item"
                                    onMouseEnter={handleThemeMenuEnter}
                                    onMouseLeave={handleThemeMenuLeave}
                                    ref={themeMenuRef}
                                >
                                    <div className="menu-item-content">
                                        {getThemeIcon(currentTheme)}
                                        <span>主题</span>
                                    </div>
                                    <div className="menu-arrow">›</div>

                                    {/* 主题子菜单 */}
                                    {isThemeMenuOpen && (
                                        <div
                                            className="theme-submenu"
                                            onMouseEnter={handleThemeMenuEnter}
                                            onMouseLeave={handleThemeMenuLeave}
                                        >
                                            {(['auto', 'light', 'dark'] as const).map(theme => (
                                                <div
                                                    key={theme}
                                                    className={`submenu-item ${currentTheme === theme ? 'active' : ''}`}
                                                    onClick={() => handleThemeChange(theme)}
                                                >
                                                    <div className="submenu-item-content">
                                                        {getThemeIcon(theme)}
                                                        <span>{getThemeText(theme)}</span>
                                                    </div>
                                                    {currentTheme === theme && (
                                                        <CheckOutlined className="check-icon"/>
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>

                                <div className="menu-item" onClick={handleLogout}>
                                    <div className="menu-item-content">
                                        <LogoutOutlined/>
                                        <span>退出登录</span>
                                    </div>
                                </div>
                            </div>
                        )}
                    </>
                )}
                {isSidebarCollapsed && (
                    <div className="user-avatar collapsed" onClick={handleUserMenuToggle}>
                        {user?.avatar ? (
                            <img src={user.avatar} alt="用户头像" className="avatar-img"/>
                        ) : (
                            <UserOutlined className="avatar-icon"/>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Index;