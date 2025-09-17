import type { ModelOption } from '@/hooks/useModels';
import React, {useRef, useState, useEffect, useCallback} from "react";
import ChatInput from "@/pages/home/chat/chat_input";
import styles from "@/pages/home/chat/index.module.scss";
import ChatMessages, {type ChatMessagesRef} from "@/pages/home/chat/chat_messages.tsx";
import ChatTitle from "@/pages/home/chat/chat_title.tsx";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";

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
    const minLoadingDuration = 500; // 最小加载时间500ms
    // 滚动状态管理
    const [autoScroll, setAutoScroll] = useState(true); // 是否启用自动滚动
    const [isUserScrolling, setIsUserScrolling] = useState(false); // 用户是否在手动滚动
    const [isAtBottom, setIsAtBottom] = useState(true); // 是否在底部
    const [showScrollButton, setShowScrollButton] = useState(false); // 是否显示滚动到底部按钮
    const chatMessagesRef = useRef<ChatMessagesRef>(null);
    const lastMessageCountRef = useRef(0); // 记录上次消息数量
    const isGeneratingRef = useRef(false); // 记录是否正在生成
    const hasInitializedRef = useRef(false); // 记录是否已初始化

    // 更新生成状态引用
    useEffect(() => {
        isGeneratingRef.current = isGenerating || false;
    }, [isGenerating]);

    // 处理用户滚动事件
    const handleUserScroll = useCallback((userScrolling: boolean) => {
        setIsUserScrolling(userScrolling);
        
        // 检查是否在底部
        if (chatMessagesRef.current) {
            const atBottom = chatMessagesRef.current.isAtBottom();
            setIsAtBottom(atBottom);
            
            if (userScrolling) {
                // 用户开始滚动时，立即根据情况禁用自动滚动（包括AI生成期间）
                if (!atBottom) {
                    // 用户滚动到非底部位置，立即禁用自动滚动并显示按钮
                    setAutoScroll(false);
                    setShowScrollButton(true);
                } else {
                    // 用户滚动到底部，恢复自动滚动并隐藏按钮
                    setAutoScroll(true);
                    setShowScrollButton(false);
                }
            } else {
                // 更新按钮显示状态
                setShowScrollButton(!atBottom && !autoScroll);
            }
        }
    }, [autoScroll]);

    // 滚动到底部的处理函数
    const handleScrollToBottom = useCallback(() => {
        if (chatMessagesRef.current) {
            chatMessagesRef.current.scrollToBottomSmooth();
            setAutoScroll(true);
            setIsUserScrolling(false);
            setIsAtBottom(true);
            setShowScrollButton(false);
        }
    }, []);

    // 定期检查底部状态（用于显示/隐藏滚动按钮）
    useEffect(() => {
        const checkBottomStatus = () => {
            if (!autoScroll && chatMessagesRef.current) {
                const atBottom = chatMessagesRef.current.isAtBottom();
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
        const currentMessageCount = currentMessages.length;
        const hasNewMessage = currentMessageCount > lastMessageCountRef.current;
        const lastMessage = currentMessages[currentMessages.length - 1];
        const isStreamingMessage = (lastMessage as any)?.isStreaming; // Message 类没有 isStreaming 属性
        
        // 更新消息数量引用
        lastMessageCountRef.current = currentMessageCount;
        
        // 决定是否需要自动滚动
        const shouldAutoScroll = autoScroll && (
            hasNewMessage || // 有新消息
            isStreamingMessage || // 消息正在流式生成
            isGenerating // 正在生成状态
        );
        
        if (shouldAutoScroll && chatMessagesRef.current && currentMessageCount > 0) {
            // 使用immediate滚动确保实时跟进
            chatMessagesRef.current.scrollToBottom();
            setIsAtBottom(true);
            setShowScrollButton(false);
        }
    }, [currentMessages, autoScroll, isGenerating]);

    // 首次进入聊天时自动滚动到底部
    useEffect(() => {
        if (!hasInitializedRef.current && chatMessagesRef.current && currentMessages.length > 0) {
            hasInitializedRef.current = true;
            // 延迟一下确保DOM渲染完成
            setTimeout(() => {
                if (chatMessagesRef.current) {
                    chatMessagesRef.current.scrollToBottomSmooth();
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
                    <div className={styles.customLoadingSpinner}>
                        <div className={styles.loadingAnimation}>
                            <div className={styles.loadingBar}></div>
                            <div className={styles.loadingBar}></div>
                            <div className={styles.loadingBar}></div>
                            <div className={styles.loadingBar}></div>
                            <div className={styles.loadingBar}></div>
                        </div>
                        <div className={styles.loadingText}>
                            <span>正在加载聊天内容</span>
                            <div className={styles.loadingDots}>
                                <span>.</span>
                                <span>.</span>
                                <span>.</span>
                            </div>
                        </div>
                    </div>
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
                    {standalone && onTitleChange && (
                        <ChatTitle 
                            chatTitle={chatTitle}
                            chatUuid={chatUuid}
                            onTitleChange={onTitleChange}
                            isMobile={isMobile}
                            isSidebarCollapsed={isSidebarCollapsed}
                            onToggleSidebar={onToggleSidebar}
                        />
                    )}
                    {!standalone && onTitleChange && (
                        <ChatTitle 
                            chatTitle={chatTitle}
                            chatUuid={chatUuid}
                            onTitleChange={onTitleChange}
                            isMobile={isMobile}
                            isSidebarCollapsed={isSidebarCollapsed}
                            onToggleSidebar={onToggleSidebar}
                        />
                    )}
                    <div className={`${styles.chatMessagesContent}`}>
                        <ChatMessages
                            ref={chatMessagesRef}
                            messages={currentMessages}
                            isLoading={isLoading}
                            showLoadingMessage={showLoadingMessage}
                            isMobile={isMobile}
                            onCopyMessage={onCopyMessage}
                            onDeleteMessage={onDeleteMessage}
                            onRegenerateMessage={onRegenerateMessage}
                            autoScrollBottom={autoScroll}
                            onUserScroll={handleUserScroll}
                        />
                    </div>
                    {(onSendMessage || onModelChange) && (
                        <ChatInput className={`${styles.chatInput}`}
                            selectedModel={selectedModel}
                            availableModels={availableModels}
                            isMobile={isMobile}
                            isGenerating={isGenerating}
                            onSendMessage={onSendMessage || (() => {})}
                            onStopGeneration={onStopGeneration}
                            onModelChange={onModelChange || (() => {})}
                            onFileUpload={onFileUpload}
                            onImageUpload={onImageUpload}
                            onMessageListScrollToBottom={handleScrollToBottom}
                            showScrollToBottom={showScrollButton}
                        />
                    )}
                </>
            )}
        </div>
    );
};

export default Chat;
export type { ChatProps };