import type { ModelOption } from '@/hooks/useModels';
import React, {useRef, useState, useEffect, useCallback, useLayoutEffect} from "react";
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import ChatInput from "@/components/chat/input";

interface ChatProps {
    // 聊天标题
    chatTitle?: string;
    // 聊天UUID
    chatUuid?: string;
    // 所有消息
    currentMessages?: Message[];
    // 是否加载中
    isLoading?: boolean;
    // 是否显示loading消息
    showLoadingMessage?: boolean;
    // 所选模型
    selectedModel?: string;
    // 可用模型
    availableModels?: ModelOption[];
    // 是否为移动端
    isMobile?: boolean;
    // 侧边栏是否收起
    isSidebarCollapsed?: boolean;
    // 是否正在生成消息
    isGenerating?: boolean;
    // 是否为独立模式（不嵌入其他页面）
    standalone?: boolean;
    // 初始化加载中
    initialLoading?: boolean;
    // 聊天标题变更事件
    onTitleChange?: (newTitle: string) => void;
    // 点击发送按钮事件
    onSendMessage?: (message: string) => void;
    // 点击停止生成按钮事件
    onStopGeneration?: () => void;
    // 模型变更事件
    onModelChange?: (model: string) => void;
    // 模型选择框点击事件（用于刷新模型数据）
    onModelSelectorClick?: () => void;
    // 切换侧边栏事件
    onToggleSidebar?: () => void;
    // 复制消息事件
    onCopyMessage?: (content: string) => void;
    // 删除消息事件
    onDeleteMessage?: (messageId: string) => void;
    // 重新生成消息事件
    onRegenerateMessage?: (messageId: string) => void;
    // 文件上传事件
    onFileUpload?: (files: File[]) => void;
    // 图片上传事件
    onImageUpload?: (files: File[]) => void;
    // 类名
    className?: string;
}

const Chat: React.FC<ChatProps> = ({
    chatTitle = '新建对话',
    chatUuid,
    currentMessages = [],
    isLoading = false,
    showLoadingMessage = false,
    selectedModel = '',
    availableModels = [],
    isMobile = false,
    isSidebarCollapsed = false,
    isGenerating = false,
    standalone = false,
    initialLoading = false,
    onTitleChange,
    onSendMessage,
    onStopGeneration,
    onModelChange,
    onModelSelectorClick,
    onToggleSidebar,
    onCopyMessage,
    onDeleteMessage,
    onRegenerateMessage,
    onFileUpload,
    onImageUpload,
    className
}) => {
    // 内部加载状态管理（用于保证最小加载时间）
    const [internalLoading, setInternalLoading] = useState(false);
    const [loadingStartTime, setLoadingStartTime] = useState<number | null>(null);
    const minLoadingDuration = 200; // 最小加载时间200ms
    // 滚动状态管理
    const [autoScroll, setAutoScroll] = useState(true); // 是否启用自动滚动
    const [isUserScrolling, setIsUserScrolling] = useState(false); // 用户是否在手动滚动
    const [isAtBottom, setIsAtBottom] = useState(true); // 是否在底部
    const [showScrollButton, setShowScrollButton] = useState(false); // 是否显示滚动到底部按钮
    const messageListRef = useRef<MessageListRef>(null);
    const lastMessageCountRef = useRef(0); // 记录上次消息数量
    const isGeneratingRef = useRef(false); // 记录是否正在生成
    const hasInitializedRef = useRef(false); // 记录是否已初始化
    const scrollTimeoutRef = useRef<NodeJS.Timeout | null>(null); // 滚动防抖定时器
    const lastScrollTimeRef = useRef(0); // 上次滚动时间
    const justCompletedRef = useRef(false); // 标记是否刚刚完成生成
    const isScrollingToBottomRef = useRef(false); // 标记是否正在滚动到底部

    // 更新生成状态引用，检测生成完成时刻
    useEffect(() => {
        const wasGenerating = isGeneratingRef.current;
        isGeneratingRef.current = isGenerating || false;
        
        // 如果从生成中变为完成，标记为刚刚完成
        if (wasGenerating && !isGeneratingRef.current) {
            justCompletedRef.current = true;
            // 300ms 后清除标记，避免影响后续的正常滚动
            setTimeout(() => {
                justCompletedRef.current = false;
            }, 300);
        }
    }, [isGenerating]);

    // 滚动到底部的处理函数
    const handleScrollToBottom = useCallback(() => {
        if (messageListRef.current) {
            // 设置滚动标记，防止在滚动过程中重新显示按钮
            isScrollingToBottomRef.current = true;
            setAutoScroll(true);
            setIsUserScrolling(false);
            setIsAtBottom(true);
            setShowScrollButton(false);
            
            messageListRef.current.scrollToBottomSmooth();
            
            // 平滑滚动通常需要 300-500ms，我们等待滚动完成后再清除标记
            // 使用稍长的时间确保滚动完成
            setTimeout(() => {
                isScrollingToBottomRef.current = false;
                // 滚动完成后再次检查底部状态，确保状态正确
                if (messageListRef.current) {
                    const atBottom = messageListRef.current.isAtBottom();
                    setIsAtBottom(atBottom);
                    setShowScrollButton(!atBottom);
                }
            }, 600);
        }
    }, []);

    // 定期检查底部状态（用于显示/隐藏滚动按钮）
    useEffect(() => {
        const checkBottomStatus = () => {
            // 如果正在滚动到底部，跳过检查，避免在滚动过程中重新显示按钮
            if (isScrollingToBottomRef.current) {
                return;
            }
            
            if (!autoScroll && messageListRef.current) {
                const atBottom = messageListRef.current.isAtBottom();
                setIsAtBottom(atBottom);
                setShowScrollButton(!atBottom);
            }
        };

        // 只在非自动滚动状态下进行定期检查
        if (!autoScroll) {
            const interval = setInterval(checkBottomStatus, 200);
            return () => clearInterval(interval);
        }
    }, [autoScroll]);

    // 消息变化时的自动滚动逻辑
    useEffect(() => {
        // 清理之前的滚动定时器
        if (scrollTimeoutRef.current) {
            clearTimeout(scrollTimeoutRef.current);
            scrollTimeoutRef.current = null;
        }

        const currentMessageCount = currentMessages.length;
        const hasNewMessage = currentMessageCount > lastMessageCountRef.current;
        const lastMessage = currentMessages[currentMessages.length - 1];
        const isStreamingMessage = (lastMessage as any)?.isStreaming; // message 类没有 isStreaming 属性
        
        // 更新消息数量引用
        const prevMessageCount = lastMessageCountRef.current;
        lastMessageCountRef.current = currentMessageCount;
        
        // 决定是否需要自动滚动
        // 只有在生成过程中或消息内容变化时才滚动，避免在生成完成时（isGenerating 从 true 变为 false）触发滚动
        // 如果刚刚完成生成，不触发滚动，避免闪烁
        const shouldAutoScroll = autoScroll && !justCompletedRef.current && (
            hasNewMessage || // 有新消息（消息数量增加）
            (isGenerating && currentMessageCount > 0) // 正在生成状态且有消息（包括内容更新时的滚动）
        );
        
        if (shouldAutoScroll && messageListRef.current && currentMessageCount > 0) {
            // 使用防抖机制，避免短时间内多次滚动导致闪烁
            const now = Date.now();
            const timeSinceLastScroll = now - lastScrollTimeRef.current;
            // 减少防抖延迟，提高响应速度，但仍避免频繁滚动
            const scrollDelay = timeSinceLastScroll < 16 ? 16 - timeSinceLastScroll : 0;
            
            scrollTimeoutRef.current = setTimeout(() => {
                if (messageListRef.current && autoScroll && !justCompletedRef.current) {
                    // 使用immediate滚动确保实时跟进
                    messageListRef.current.scrollToBottom();
                    setIsAtBottom(true);
                    setShowScrollButton(false);
                    lastScrollTimeRef.current = Date.now();
                }
            }, scrollDelay);
        }
        
        // 清理函数
        return () => {
            if (scrollTimeoutRef.current) {
                clearTimeout(scrollTimeoutRef.current);
                scrollTimeoutRef.current = null;
            }
        };
    }, [currentMessages, autoScroll, isGenerating]);

    // 首次进入聊天时自动滚动到底部
    useEffect(() => {
        if (!hasInitializedRef.current && messageListRef.current && currentMessages.length > 0) {
            hasInitializedRef.current = true;
            // 延迟一下确保DOM渲染完成
            setTimeout(() => {
                if (messageListRef.current) {
                    messageListRef.current.scrollToBottomSmooth();
                    setIsAtBottom(true);
                    setShowScrollButton(false);
                }
            }, 100);
        }
    }, [currentMessages.length]);

    // 用户发送新消息时确保自动滚动
    useEffect(() => {
        const messageCount = currentMessages.length;
        if (messageCount > 0) {
            const lastMessage = currentMessages[messageCount - 1];
            // 如果最后一条消息是用户消息，说明用户刚发送消息，启用自动滚动
            if (lastMessage.role === 'user') {
                setAutoScroll(true);
                setIsUserScrolling(false);
                setShowScrollButton(false);
            }
        }
    }, [currentMessages.length]);

    // 管理加载状态，确保最小加载时间
    useEffect(() => {
        if (initialLoading || isLoading) {
            if (!internalLoading) {
                setInternalLoading(true);
                setLoadingStartTime(Date.now());
            }
        } else if (internalLoading && loadingStartTime) {
            const elapsed = Date.now() - loadingStartTime;
            const remaining = Math.max(0, minLoadingDuration - elapsed);
            
            if (remaining > 0) {
                setTimeout(() => {
                    setInternalLoading(false);
                    setLoadingStartTime(null);
                }, remaining);
            } else {
                setInternalLoading(false);
                setLoadingStartTime(null);
            }
        }
    }, [initialLoading, isLoading, internalLoading, loadingStartTime]);

    // 决定是否显示加载状态
    const shouldShowLoading = internalLoading || initialLoading;

    return (
        <div className={`${styles.chatPage} ${standalone ? styles.standalone : ''} ${className || ''}`}>
            {shouldShowLoading ? (
                <div className={styles.chatLoadingContainer}>
                    <div className={styles.loadingPlaceholder}>
                        {/* 占位内容，防止布局跳动 */}
                        <div className={styles.loadingTitle}></div>
                        <div className={styles.loadingMessages}>
                            <div className={styles.loadingMessage}></div>
                            <div className={styles.loadingMessage}></div>
                            <div className={styles.loadingMessage}></div>
                        </div>
                        <div className={styles.loadingInput}></div>
                    </div>
                </div>
            ) : (
                <>
                    <ChatTitle
                        title={chatTitle}
                        uuid={chatUuid}
                        onTitleChange={onTitleChange}
                        isSidebarCollapsed={isSidebarCollapsed}
                        onToggleSidebar={onToggleSidebar}
                    />
                    <div className={`${styles.chatMessagesContent}`}>
                        <MessageList
                            ref={messageListRef}
                            messages={currentMessages}
                            isLoading={isLoading}
                        />
                    </div>
                    {(onSendMessage || onModelChange) && (
                        <div className={`${styles.chatInput}`}>
                            <ChatInput
                                       selectedModel={selectedModel}
                                       availableModels={availableModels}
                                       isGenerating={isGenerating}
                                       onSendMessage={onSendMessage || (() => {})}
                                       onStopGeneration={onStopGeneration}
                                       onModelChange={onModelChange || (() => {})}
                                       onModelSelectorClick={onModelSelectorClick}
                                       onMessageListScrollToBottom={handleScrollToBottom}
                            />
                        </div>

                    )}
                </>
            )}
        </div>
    );
};

export default Chat;
export type { ChatProps };