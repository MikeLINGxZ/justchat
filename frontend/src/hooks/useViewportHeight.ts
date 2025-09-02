import { useEffect, useState } from 'react';

/**
 * 移动端视口高度检测 Hook
 * 解决移动端浏览器地址栏和底部菜单栏遮挡问题
 */
export const useViewportHeight = () => {
    const [viewportHeight, setViewportHeight] = useState(window.innerHeight);
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        // 检测是否为移动设备
        const checkIsMobile = () => {
            const userAgent = navigator.userAgent;
            const mobileRegex = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i;
            return mobileRegex.test(userAgent) || window.innerWidth <= 768;
        };

        // 更新视口高度
        const updateViewportHeight = () => {
            const vh = window.innerHeight;
            setViewportHeight(vh);
            
            // 设置 CSS 自定义属性
            document.documentElement.style.setProperty('--vh', `${vh * 0.01}px`);
            document.documentElement.style.setProperty('--viewport-height', `${vh}px`);
        };

        // 初始化
        setIsMobile(checkIsMobile());
        updateViewportHeight();

        // 监听窗口大小变化
        const handleResize = () => {
            setIsMobile(checkIsMobile());
            updateViewportHeight();
        };

        // 监听屏幕方向变化（移动端特有）
        const handleOrientationChange = () => {
            // 延迟执行，等待浏览器完成方向变化
            setTimeout(() => {
                updateViewportHeight();
            }, 100);
        };

        // 监听可视区域变化（主要为了检测移动端浏览器工具栏显示/隐藏）
        const handleVisualViewportChange = () => {
            if (window.visualViewport) {
                const vh = window.visualViewport.height;
                setViewportHeight(vh);
                document.documentElement.style.setProperty('--vh', `${vh * 0.01}px`);
                document.documentElement.style.setProperty('--viewport-height', `${vh}px`);
            }
        };

        // 添加事件监听器
        window.addEventListener('resize', handleResize);
        window.addEventListener('orientationchange', handleOrientationChange);
        
        // 现代浏览器支持 Visual Viewport API
        if (window.visualViewport) {
            window.visualViewport.addEventListener('resize', handleVisualViewportChange);
        }

        // 清理函数
        return () => {
            window.removeEventListener('resize', handleResize);
            window.removeEventListener('orientationchange', handleOrientationChange);
            
            if (window.visualViewport) {
                window.visualViewport.removeEventListener('resize', handleVisualViewportChange);
            }
        };
    }, []);

    return {
        viewportHeight,
        isMobile,
        // 提供一个方法手动触发更新
        updateHeight: () => {
            const vh = window.visualViewport?.height || window.innerHeight;
            setViewportHeight(vh);
            document.documentElement.style.setProperty('--vh', `${vh * 0.01}px`);
            document.documentElement.style.setProperty('--viewport-height', `${vh}px`);
        }
    };
};
