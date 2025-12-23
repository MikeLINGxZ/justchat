import type { ModelOption } from '@/hooks/useModels';
import React, {useRef, useState, useEffect, useCallback, useLayoutEffect} from "react";
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import ChatInput from "@/components/chat/input";
import {File} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";

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
    onSendMessage?: (message: string,files:File[]) => void;
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
    const messageListRef = useRef<MessageListRef>(null);


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
                            isGenerating={isGenerating}
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
                                       onMessageListScrollToBottom={() => {
                                           messageListRef.current?.scrollToBottom();
                                       }}
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