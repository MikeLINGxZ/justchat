import React, {
    forwardRef,
    useImperativeHandle,
    useRef,
    useCallback,
    useEffect,
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
    // 所有消息
    messages?: Message[]; // 修改为 message[]
    // 是否加载中
    isLoading?: boolean;
    // 消息是否在生成中
    isGenerating?: boolean;
}

export interface MessageListRef {
    scrollToBottom: () => void;
    scrollToBottomSmooth: () => void;
    isAtBottom: () => boolean;
    enableAutoScroll: () => void;
    disableAutoScroll: () => void;
}

const MessageList: React.ForwardRefRenderFunction<MessageListRef, MessageListProps> = ({
    className,
    messages = [],
    isLoading,
    isGenerating
}, ref) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const [isAtBottom, setIsAtBottom] = useState(true);
    const [showScrollButton, setShowScrollButton] = useState(false);
    const isInitialLoadRef = useRef(true);

    // 检查是否在底部
    const checkIsAtBottom = useCallback(() => {
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

    // 滚动到底部（立即）
    const scrollToBottom = useCallback(() => {
        if (containerRef.current) {
            // 获取可滚动的父容器
            const scrollContainer = containerRef.current.closest('[class*="chatMessagesContent"]') || 
                                    containerRef.current.parentElement;
            
            if (scrollContainer && scrollContainer instanceof HTMLElement) {
                scrollContainer.scrollTop = scrollContainer.scrollHeight;
                setIsAtBottom(true);
                setShowScrollButton(false);
                return;
            }
            
            // 回退到当前容器
            containerRef.current.scrollTop = containerRef.current.scrollHeight;
            setIsAtBottom(true);
            setShowScrollButton(false);
        }
    }, []);

    // 滚动到底部（平滑）
    const scrollToBottomSmooth = useCallback(() => {
        if (containerRef.current) {
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
                return;
            }
            
            // 回退到当前容器
            containerRef.current.scrollTo({
                top: containerRef.current.scrollHeight,
                behavior: 'smooth'
            });
            setIsAtBottom(true);
            setShowScrollButton(false);
        }
    }, []);

    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        scrollToBottom,
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
        enableAutoScroll: () => {
            // 可以在这里添加自动滚动的逻辑
        },
        disableAutoScroll: () => {
            // 可以在这里添加禁用自动滚动的逻辑
        }
    }));

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

    // 初次加载时自动滚动到底部
    useEffect(() => {
        if (isInitialLoadRef.current && messages.length > 0 && !isLoading) {
            isInitialLoadRef.current = false;
            // 延迟一下确保DOM渲染完成
            setTimeout(() => {
                scrollToBottomSmooth();
            }, 100);
        }
    }, [messages.length, isLoading, scrollToBottomSmooth]);

    // 消息变化时，如果在底部则自动滚动
    useEffect(() => {
        if (!isInitialLoadRef.current && isAtBottom && messages.length > 0) {
            scrollToBottom();
        }
    }, [messages.length, isAtBottom, scrollToBottom]);

    return (
        <>
            <div ref={containerRef} className={`${className} ${styles.MessageList}`}>
                {/* 消息 */}
                <div className={styles.content}>
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
                <div className={styles.scrollToBottomButton} onClick={scrollToBottomSmooth}>
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