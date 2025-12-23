import React, {
    forwardRef,
    useImperativeHandle,
    useRef,
    useCallback,
    useEffect,
    useLayoutEffect,
    useState,
    useMemo,
    type ReactNode
} from "react";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import styles from "./index.module.scss";
import ChatMessage from "@/components/chat/message";

interface MessageListProps {
    // 类名
    className?: string;
    // 消息列表
    messages?: Message[];
}

export interface MessageListRef {
    scrollToBottomSmooth: () => void;
    setIsGenerating: (isGenerating:boolean) => void;
    setIsAutoScroll: (isAutoScroll:boolean) => void;
    isAtBottom: () => boolean;
}

const MessageList: React.ForwardRefRenderFunction<MessageListRef, MessageListProps> = ({
    className,
    messages = [],
}, ref) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const contentRef = useRef<HTMLDivElement>(null);
    const buttonRef = useRef<HTMLDivElement>(null);
    const [isAtBottom, setIsAtBottom] = useState(true);
    const [showScrollButton, setShowScrollButton] = useState(false);
    // 是否在生成消息中
    const [isGenerating,setIsGenerating] = useState(false);
    // 是否启用自动滚动
    const [isAutoScroll,setIsAutoScroll] = useState(true);
    // 是否手动接管滚动
    const [isManualScroll,setIsManualScroll] = useState(false);
    // 是否初始化滚动
    const isInitialLoadRef = useRef(true);
    // 上次滚动时间
    const lastScrollTimeRef = useRef(0);
    // 滚动到底部定时器
    const scrollTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    // 滚动锁
    const scrollingToBottomLockRef = useRef(false);

    // 检查是否在底部
    const checkIsAtBottom = useCallback(() => {
        // 如果正在滚动到底部，跳过检查，避免在滚动过程中重新显示按钮
        if (scrollingToBottomLockRef.current) {
            return true;
        }
        if (containerRef.current) {
            // 获取可滚动的父容器
            const scrollContainer = containerRef.current.closest('[class*="chatMessagesContent"]') || 
                                    containerRef.current.parentElement;
            
            if (scrollContainer && scrollContainer instanceof HTMLElement) {
                const { scrollTop, scrollHeight, clientHeight } = scrollContainer;
                const threshold = 50; // 距离底部的阈值
                const atBottom = scrollHeight - scrollTop - clientHeight < threshold;
                setIsAtBottom(atBottom);
                setShowScrollButton(!atBottom);
                return atBottom;
            }
            
            // 回退到检查当前容器
            const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
            const threshold = 50;
            const atBottom = scrollHeight - scrollTop - clientHeight < threshold;
            setIsAtBottom(atBottom);
            setShowScrollButton(!atBottom);
            return atBottom;
        }
        return true;
    }, []);

    // 滚动到底部（平滑）
    const scrollToBottomSmooth = useCallback(() => {
        if (containerRef.current) {
            // 设置滚动标记，防止在滚动过程中重新显示按钮
            scrollingToBottomLockRef.current = true;
            
            // 获取可滚动的父容器
            const scrollContainer = containerRef.current.closest('[class*="chatMessagesContent"]') || 
                                    containerRef.current.parentElement;
            
            if (scrollContainer && scrollContainer instanceof HTMLElement) {
                scrollContainer.scrollTo({
                    top: scrollContainer.scrollHeight,
                    behavior: 'smooth'
                });
                setIsAtBottom(true);
                setShowScrollButton(false);
                
                // 平滑滚动通常需要 300-500ms，我们等待滚动完成后再清除标记
                setTimeout(() => {
                    scrollingToBottomLockRef.current = false;
                    // 滚动完成后再次检查底部状态
                    checkIsAtBottom();
                }, 600);
                return;
            }
            
            // 回退到当前容器
            containerRef.current.scrollTo({
                top: containerRef.current.scrollHeight,
                behavior: 'smooth'
            });
            setIsAtBottom(true);
            setShowScrollButton(false);
            
            // 平滑滚动通常需要 300-500ms，我们等待滚动完成后再清除标记
            setTimeout(() => {
                scrollingToBottomLockRef.current = false;
                // 滚动完成后再次检查底部状态
                checkIsAtBottom();
            }, 600);
        }
    }, [checkIsAtBottom]);



    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        scrollToBottomSmooth,
        isAtBottom: () => {
            if (containerRef.current) {
                // 获取可滚动的父容器
                const scrollContainer = containerRef.current.closest('[class*="chatMessagesContent"]') || 
                                        containerRef.current.parentElement;
                
                if (scrollContainer && scrollContainer instanceof HTMLElement) {
                    const { scrollTop, scrollHeight, clientHeight } = scrollContainer;
                    return scrollHeight - scrollTop - clientHeight < 10;
                }
                
                // 回退到当前容器
                const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
                return scrollHeight - scrollTop - clientHeight < 10;
            }
            return true;
        },
        setIsGenerating: (isGenerating:boolean) => {
            setIsGenerating(isGenerating);
        },
        setIsAutoScroll: (isAutoScroll:boolean) => {
            setIsAutoScroll(isAutoScroll);
        }
    }));

    // 更新滚动到底部按钮位置，使其相对于内容区域居中，并避免与输入框重叠
    const updateScrollToBottomButtonPosition = useCallback(() => {
        if (!contentRef.current || !buttonRef.current) return;
        
        const button = buttonRef.current;
        
        // 临时禁用过渡动画，避免位置改变时的动画效果
        const originalTransition = button.style.transition;
        button.style.transition = 'none';
        
        const contentRect = contentRef.current.getBoundingClientRect();
        
        // 计算内容区域的中心位置
        const centerX = contentRect.left + contentRect.width / 2;
        button.style.left = `${centerX}px`;
        // 确保 transform 保持，用于居中
        button.style.transform = 'translateX(-50%)';
        
        // 计算输入框高度，确保按钮不重叠
        const chatInput = document.querySelector('[class*="chatInput"]') as HTMLElement;
        if (chatInput) {
            const inputRect = chatInput.getBoundingClientRect();
            const inputHeight = inputRect.height;
            // 按钮距离输入框顶部至少 20px
            const minBottom = inputHeight + 20;
            // 移动端间距较小（只多10px）
            const isMobile = window.innerWidth <= 768;
            const bottom = isMobile ? Math.max(minBottom + 10, 100) : Math.max(minBottom, 120);
            button.style.bottom = `${bottom}px`;
        }
        
        // 使用 flushSync 或强制重排，确保样式立即应用
        void button.offsetHeight; // 强制重排
        
        // 恢复过渡动画（延迟恢复，确保位置更新完成）
        requestAnimationFrame(() => {
            button.style.transition = originalTransition || '';
        });
    }, []);

    // 监听滚动事件（监听父容器的滚动）
    useEffect(() => {
        const container = containerRef.current;
        if (!container) return;

        // 获取可滚动的父容器
        const scrollContainer = container.closest('[class*="chatMessagesContent"]') || 
                                container.parentElement;
        
        if (!scrollContainer) return;

        const handleScroll = () => {
            checkIsAtBottom();
        };

        scrollContainer.addEventListener('scroll', handleScroll);
        // 初始化检查一次
        checkIsAtBottom();
        
        return () => {
            scrollContainer.removeEventListener('scroll', handleScroll);
        };
    }, [checkIsAtBottom]);

    // 使用 useLayoutEffect 确保按钮在显示时立即设置正确位置，避免闪烁
    useLayoutEffect(() => {
        if (!showScrollButton || !buttonRef.current) return;
        // 立即同步设置位置
        updateScrollToBottomButtonPosition();
    }, [showScrollButton, updateScrollToBottomButtonPosition]);

    // 监听窗口大小变化和内容区域变化，更新按钮位置
    useEffect(() => {
        if (!showScrollButton || !buttonRef.current) return;

        const handleResize = () => {
            updateScrollToBottomButtonPosition();
        };

        window.addEventListener('resize', handleResize);
        // 使用 ResizeObserver 监听内容区域大小变化
        let contentResizeObserver: ResizeObserver | null = null;
        if (contentRef.current) {
            contentResizeObserver = new ResizeObserver(() => {
                updateScrollToBottomButtonPosition();
            });
            contentResizeObserver.observe(contentRef.current);
        }

        // 监听输入框大小变化
        const chatInput = document.querySelector('[class*="chatInput"]') as HTMLElement;
        let inputResizeObserver: ResizeObserver | null = null;
        if (chatInput) {
            inputResizeObserver = new ResizeObserver(() => {
                updateScrollToBottomButtonPosition();
            });
            inputResizeObserver.observe(chatInput);
        }

        return () => {
            window.removeEventListener('resize', handleResize);
            if (contentResizeObserver && contentRef.current) {
                contentResizeObserver.unobserve(contentRef.current);
            }
            if (inputResizeObserver && chatInput) {
                inputResizeObserver.unobserve(chatInput);
            }
        };
    }, [showScrollButton, updateScrollToBottomButtonPosition]);

    // 初次加载时自动滚动到底部
    useEffect(() => {
        if (isInitialLoadRef.current && messages.length > 0) {
            isInitialLoadRef.current = false;
            // 延迟一下确保DOM渲染完成
            setTimeout(() => {
                scrollToBottomSmooth();
            }, 100);
        }
    }, [messages.length, scrollToBottomSmooth]);

    // 消息变化时，如果在底部则自动滚动
    // 使用 useLayoutEffect 确保在 DOM 更新后立即同步滚动，避免视觉闪烁
    useLayoutEffect(() => {
        // 清理之前的滚动定时器
        if (scrollTimeoutRef.current) {
            clearTimeout(scrollTimeoutRef.current);
            scrollTimeoutRef.current = null;
        }

        if (!isInitialLoadRef.current && isAtBottom && messages.length > 0) {
            // 使用防抖机制，避免短时间内多次滚动导致闪烁
            const now = Date.now();
            const timeSinceLastScroll = now - lastScrollTimeRef.current;
            const scrollDelay = timeSinceLastScroll < 50 ? 50 - timeSinceLastScroll : 0;
            
            scrollTimeoutRef.current = setTimeout(() => {
                scrollToBottomSmooth();
                lastScrollTimeRef.current = Date.now();
            }, scrollDelay);
        }
        
        // 清理函数
        return () => {
            if (scrollTimeoutRef.current) {
                clearTimeout(scrollTimeoutRef.current);
                scrollTimeoutRef.current = null;
            }
        };
    }, [messages.length, isAtBottom]);

    return (
        <>
            <div ref={containerRef} className={`${className} ${styles.MessageList}`}>
                {/* 消息 */}
                <div ref={contentRef} className={styles.content}>
                    {
                        messages.map((message: Message, index: number) => (
                            <div key={index}>
                                <ChatMessage message={message}/>
                            </div>
                        ))
                    }
                </div>
            </div>
            {/* 滚动到底部按钮 - 使用固定定位，不跟随滚动 */}
            {showScrollButton && (
                <div 
                    ref={buttonRef}
                    className={styles.scrollToBottomButton} 
                    onClick={scrollToBottomSmooth}
                >
                    <svg
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            d="M7 13L12 18L17 13"
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                        />
                        <path
                            d="M7 6L12 11L17 6"
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                        />
                    </svg>
                </div>
            )}
        </>
    )
};

MessageList.displayName = 'MessageList';

export default forwardRef(MessageList);