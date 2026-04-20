import React, { useEffect, useRef, useState } from 'react';
import './SidebarUserMenu.scss';
import {
    BulbOutlined,
    CheckOutlined,
    InfoCircleOutlined,
    MoonOutlined,
    SettingOutlined,
    SunOutlined,
    TeamOutlined,
    UserOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useOPCStore } from '@/stores/opcStore';
import type { AppMode } from '@/stores/opcStore';
import { OpenSettingsWindow, OpenSettingsAboutWindow } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/service.ts';

interface SidebarUserMenuProps {
    currentMode: AppMode;
}

const SidebarUserMenu: React.FC<SidebarUserMenuProps> = ({
    currentMode,
}) => {
    const { t } = useTranslation();
    const { setMode } = useOPCStore();

    const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);
    const [isThemeMenuOpen, setIsThemeMenuOpen] = useState(false);
    const [isModeMenuOpen, setIsModeMenuOpen] = useState(false);
    const [currentTheme, setCurrentTheme] = useState<'auto' | 'light' | 'dark'>('auto');

    const userMenuRef = useRef<HTMLDivElement>(null);
    const themeMenuRef = useRef<HTMLDivElement>(null);
    const modeMenuRef = useRef<HTMLDivElement>(null);
    const themeCloseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const modeCloseTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // 点击外部关闭菜单
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
                setIsUserMenuOpen(false);
            }
            if (themeMenuRef.current && !themeMenuRef.current.contains(event.target as Node)) {
                setIsThemeMenuOpen(false);
            }
        };
        if (isUserMenuOpen || isThemeMenuOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [isUserMenuOpen, isThemeMenuOpen]);

    // 从 localStorage 读取主题设置
    useEffect(() => {
        const savedTheme = localStorage.getItem('theme') as 'auto' | 'light' | 'dark';
        if (savedTheme) {
            setCurrentTheme(savedTheme);
        }
    }, []);

    // 清理定时器
    useEffect(() => {
        return () => {
            if (themeCloseTimeoutRef.current) clearTimeout(themeCloseTimeoutRef.current);
            if (modeCloseTimeoutRef.current) clearTimeout(modeCloseTimeoutRef.current);
        };
    }, []);

    const handleUserMenuToggle = () => {
        setIsUserMenuOpen(!isUserMenuOpen);
        setIsThemeMenuOpen(false);
        setIsModeMenuOpen(false);
    };

    const handleThemeChange = (theme: 'auto' | 'light' | 'dark') => {
        setCurrentTheme(theme);
        localStorage.setItem('theme', theme);
        const root = document.documentElement;
        if (theme === 'dark') {
            root.classList.add('dark');
            root.classList.remove('light');
        } else if (theme === 'light') {
            root.classList.add('light');
            root.classList.remove('dark');
        } else {
            root.classList.remove('dark', 'light');
            if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                root.classList.add('dark');
            } else {
                root.classList.add('light');
            }
        }
        if (themeCloseTimeoutRef.current) {
            clearTimeout(themeCloseTimeoutRef.current);
            themeCloseTimeoutRef.current = null;
        }
        setIsThemeMenuOpen(false);
        setIsUserMenuOpen(false);
    };

    const handleModeChange = (mode: AppMode) => {
        setMode(mode);
        setIsUserMenuOpen(false);
    };

    const getThemeIcon = (theme: 'auto' | 'light' | 'dark') => {
        switch (theme) {
            case 'light': return <SunOutlined />;
            case 'dark': return <MoonOutlined />;
            default: return <BulbOutlined />;
        }
    };

    const getThemeText = (theme: 'auto' | 'light' | 'dark') => {
        switch (theme) {
            case 'light': return t('home.sidebar.themes.light');
            case 'dark': return t('home.sidebar.themes.dark');
            default: return t('home.sidebar.themes.auto');
        }
    };

    const handleThemeMenuEnter = () => {
        if (themeCloseTimeoutRef.current) { clearTimeout(themeCloseTimeoutRef.current); themeCloseTimeoutRef.current = null; }
        setIsThemeMenuOpen(true);
    };
    const handleThemeMenuLeave = () => {
        themeCloseTimeoutRef.current = setTimeout(() => setIsThemeMenuOpen(false), 150);
    };
    const handleModeMenuEnter = () => {
        if (modeCloseTimeoutRef.current) { clearTimeout(modeCloseTimeoutRef.current); modeCloseTimeoutRef.current = null; }
        setIsModeMenuOpen(true);
    };
    const handleModeMenuLeave = () => {
        modeCloseTimeoutRef.current = setTimeout(() => setIsModeMenuOpen(false), 150);
    };

    return (
        <div className="sidebar-footer" ref={userMenuRef}>
            <>
                <div
                        className={`user-section settings-trigger ${isUserMenuOpen ? 'active' : ''}`}
                        onClick={handleUserMenuToggle}
                    >
                        <SettingOutlined className="settings-trigger-icon" />
                        <span className="settings-trigger-label">{t('home.sidebar.settings')}</span>
                    </div>

                    {/* 用户菜单 */}
                    {isUserMenuOpen && (
                        <div className="user-menu">
                            {/* 模式切换 */}
                            <div
                                className="menu-item mode-item"
                                onMouseEnter={handleModeMenuEnter}
                                onMouseLeave={handleModeMenuLeave}
                                ref={modeMenuRef}
                            >
                                <div className="menu-item-content">
                                    <TeamOutlined />
                                    <span>{t('opc.sidebar.mode')}</span>
                                </div>
                                <div className="menu-arrow">›</div>
                                {isModeMenuOpen && (
                                    <div
                                        className="theme-submenu"
                                        onMouseEnter={handleModeMenuEnter}
                                        onMouseLeave={handleModeMenuLeave}
                                    >
                                        <div
                                            className={`submenu-item ${currentMode === 'chat' ? 'active' : ''}`}
                                            onClick={() => handleModeChange('chat')}
                                        >
                                            <div className="submenu-item-content">
                                                <UserOutlined />
                                                <span>{t('opc.sidebar.modeChat')}</span>
                                            </div>
                                            {currentMode === 'chat' && <CheckOutlined className="check-icon" />}
                                        </div>
                                        <div
                                            className={`submenu-item ${currentMode === 'opc' ? 'active' : ''}`}
                                            onClick={() => handleModeChange('opc')}
                                        >
                                            <div className="submenu-item-content">
                                                <TeamOutlined />
                                                <span>{t('opc.sidebar.modeOPC')}</span>
                                                <span className="beta-tag" title="功能开发中，未完全可用">beta</span>
                                            </div>
                                            {currentMode === 'opc' && <CheckOutlined className="check-icon" />}
                                        </div>
                                    </div>
                                )}
                            </div>

                            {/* 主题切换 */}
                            <div
                                className="menu-item theme-item"
                                onMouseEnter={handleThemeMenuEnter}
                                onMouseLeave={handleThemeMenuLeave}
                                ref={themeMenuRef}
                            >
                                <div className="menu-item-content">
                                    {getThemeIcon(currentTheme)}
                                    <span>{t('home.sidebar.theme')}</span>
                                </div>
                                <div className="menu-arrow">›</div>
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
                                                {currentTheme === theme && <CheckOutlined className="check-icon" />}
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div className="menu-item" onClick={() => { OpenSettingsAboutWindow(); setIsUserMenuOpen(false); }}>
                                <div className="menu-item-content">
                                    <InfoCircleOutlined />
                                    <span>{t('home.sidebar.about')}</span>
                                </div>
                            </div>

                            <div className="menu-item" onClick={() => { OpenSettingsWindow(); setIsUserMenuOpen(false); }}>
                                <div className="menu-item-content">
                                    <SettingOutlined />
                                    <span>{t('home.sidebar.settings')}</span>
                                </div>
                            </div>
                        </div>
                    )}
            </>
        </div>
    );
};

export default SidebarUserMenu;
