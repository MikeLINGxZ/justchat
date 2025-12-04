// utils/mobile.ts 或直接放在 useViewportHeight.ts 文件中
import { useEffect, useState } from 'react';

/**
 * 独立的工具函数：判断是否为移动设备
 * 可在非 React 环境中使用
 */
export const isMobileDevice = (widthThreshold = 768): boolean => {
    if (typeof window === 'undefined' || typeof navigator === 'undefined') return false;

    const userAgent = navigator.userAgent;
    const mobileRegex = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i;
    return mobileRegex.test(userAgent) || window.innerWidth <= widthThreshold;
};

/**
 * React Hook: 移动端视口高度检测
 * 解决移动端浏览器地址栏和底部菜单栏遮挡问题
 */
// 修改 useViewportHeight
export const useViewportHeight = () => {
    const [viewportHeight, setViewportHeight] = useState(window.innerHeight);
    const isMobile = useIsMobile(768); // 复用新 Hook

    useEffect(() => {
        const updateHeight = () => {
            const height = window.visualViewport?.height || window.innerHeight;
            setViewportHeight(height);
            document.documentElement.style.setProperty('--vh', `${height * 0.01}px`);
            document.documentElement.style.setProperty('--viewport-height', `${height}px`);
        };

        updateHeight();

        const handleResize = () => updateHeight();
        const handleOrientation = () => setTimeout(updateHeight, 100);
        const handleVisualViewport = () => window.visualViewport && updateHeight();

        window.addEventListener('resize', handleResize);
        window.addEventListener('orientationchange', handleOrientation);
        if (window.visualViewport) {
            window.visualViewport.addEventListener('resize', handleVisualViewport);
        }

        return () => {
            window.removeEventListener('resize', handleResize);
            window.removeEventListener('orientationchange', handleOrientation);
            if (window.visualViewport) {
                window.visualViewport.removeEventListener('resize', handleVisualViewport);
            }
        };
    }, []);

    return {
        viewportHeight,
        isMobile,
        updateHeight: () => {
            const height = window.visualViewport?.height || window.innerHeight;
            setViewportHeight(height);
            document.documentElement.style.setProperty('--vh', `${height * 0.01}px`);
            document.documentElement.style.setProperty('--viewport-height', `${height}px`);
        },
    };
};

/**
 * React Hook: 实时监听是否为移动设备
 * 自动响应窗口缩放、横竖屏切换等
 */
export const useIsMobile = (widthThreshold = 768): boolean => {
    const [isMobile, setIsMobile] = useState(() => isMobileDevice(widthThreshold));

    useEffect(() => {
        let timeout: NodeJS.Timeout;

        const update = () => {
            // 防抖优化，防止频繁触发
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                setIsMobile(isMobileDevice(widthThreshold));
            }, 100);
        };

        // 监听 resize
        window.addEventListener('resize', update);
        // 移动端横屏/竖屏切换
        window.addEventListener('orientationchange', update);

        // 清理
        return () => {
            window.removeEventListener('resize', update);
            window.removeEventListener('orientationchange', update);
            clearTimeout(timeout);
        };
    }, [widthThreshold]);

    return isMobile;
};