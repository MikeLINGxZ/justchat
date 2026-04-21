import React, { useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
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
    const themeSubmenuPortalRef = useRef<HTMLDivElement>(null);
    const modeSubmenuPortalRef = useRef<HTMLDivElement>(null);
    const themeCloseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const modeCloseTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    const [themeSubmenuPos, setThemeSubmenuPos] = useState<{ top: number; left: number } | null>(null);
    const [modeSubmenuPos, setModeSubmenuPos] = useState<{ top: number; left: number } | null>(null);

    // 点击外部关闭菜单
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            const target = event.target as Node;
            const inTheme =
                themeMenuRef.current?.contains(target) ||
                themeSubmenuPortalRef.current?.contains(target);
            const inMode =
                modeMenuRef.current?.contains(target) ||
                modeSubmenuPortalRef.current?.contains(target);
            const inUser = userMenuRef.current?.contains(target) || inTheme || inMode;
            if (!inUser) setIsUserMenuOpen(false);
            if (!inTheme) setIsThemeMenuOpen(false);
            if (!inMode) setIsModeMenuOpen(false);
        };
        if (isUserMenuOpen || isThemeMenuOpen || isModeMenuOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [isUserMenuOpen, isThemeMenuOpen, isModeMenuOpen]);

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

    // 计算子菜单相对视口的锚点坐标：默认贴在触发项右侧，空间不够则翻到左侧。
    const computeSubmenuPosition = (trigger: HTMLElement | null) => {
        if (!trigger) return null;
        const rect = trigger.getBoundingClientRect();
        const gap = 8;
        const estimatedWidth = 160; // min-width 140 + 内边距
        const viewportWidth = window.innerWidth;
        let left = rect.right + gap;
        if (left + estimatedWidth > viewportWidth - 8) {
            left = Math.max(8, rect.left - estimatedWidth - gap);
        }
        return { top: rect.top, left };
    };

    const handleThemeMenuEnter = () => {
        if (themeCloseTimeoutRef.current) { clearTimeout(themeCloseTimeoutRef.current); themeCloseTimeoutRef.current = null; }
        setThemeSubmenuPos(computeSubmenuPosition(themeMenuRef.current));
        setIsThemeMenuOpen(true);
    };
    const handleThemeMenuLeave = () => {
        themeCloseTimeoutRef.current = setTimeout(() => setIsThemeMenuOpen(false), 150);
    };
    const handleModeMenuEnter = () => {
        if (modeCloseTimeoutRef.current) { clearTimeout(modeCloseTimeoutRef.current); modeCloseTimeoutRef.current = null; }
        setModeSubmenuPos(computeSubmenuPosition(modeMenuRef.current));
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
                                {isModeMenuOpen && modeSubmenuPos && createPortal(
                                    <div
                                        ref={modeSubmenuPortalRef}
                                        className="sidebar-user-submenu"
                                        style={{ top: modeSubmenuPos.top, left: modeSubmenuPos.left }}
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
                                    </div>,
                                    document.body
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
                                {isThemeMenuOpen && themeSubmenuPos && createPortal(
                                    <div
                                        ref={themeSubmenuPortalRef}
                                        className="sidebar-user-submenu"
                                        style={{ top: themeSubmenuPos.top, left: themeSubmenuPos.left }}
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
                                    </div>,
                                    document.body
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
