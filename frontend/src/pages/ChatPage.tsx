import React, { useCallback, useEffect, useState } from 'react';
import { message } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { Message, RoleType } from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import { useViewportHeight } from '@/hooks/useViewportHeight';
import { useModels } from '@/hooks/useModels';
import Chat from '@/pages/home/chat';
import type { ChatProps } from '@/pages/home/chat';
import { Service } from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import { CompletionsUtils } from "@/utils/completions.ts";
import {
    MessageList
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/index.ts";

interface ChatPageProps {
    className?: string;
}

const ChatPage: React.FC<ChatPageProps> = ({ className }) => {
    // 获取路由参数和导航函数
    const { chatUuid: urlChatUuid } = useParams<{ chatUuid?: string }>();
    const navigate = useNavigate();

    // 本地状态管理
    const [currentChatUuid, setCurrentChatUuid] = useState<string>(urlChatUuid || '');
    const [currentMessages, setCurrentMessages] = useState<Message[]>([]);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [showLoadingMessage, setShowLoadingMessage] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [isStreaming, setIsStreaming] = useState(false);
    const [abortController, setAbortController] = useState<AbortController | null>(null);
    const [chatTitle, setChatTitle] = useState('新建对话');
    const [selectedModel, setSelectedModel] = useState('');
    const [initialLoading, setInitialLoading] = useState(true);

    // 使用视口高度检测 Hook
    const { isMobile } = useViewportHeight();

    // 使用模型获取 Hook
    const {
        models: availableModels,
        isLoading: isLoadingModels,
        error: modelsError,
    } = useModels();

    // 设置页面标题
    useEffect(() => {
        document.title = 'AI聊天 - Lemon Tea';
    }, []);

    // 同步URL参数与当前聊天UUID
    useEffect(() => {
        const newChatUuid = urlChatUuid || '';
        if (newChatUuid !== currentChatUuid) {
            setCurrentChatUuid(newChatUuid);
        }
    }, [urlChatUuid, currentChatUuid]);

    // 设置默认选中的模型
    useEffect(() => {
        if (availableModels.length > 0 && !selectedModel) {
            setSelectedModel(availableModels[0].id);
            setInitialLoading(false);
        }
    }, [availableModels, selectedModel]);

    // 显示模型加载错误
    useEffect(() => {
        if (modelsError) {
            message.error(`获取模型列表失败: ${modelsError}`);
            setInitialLoading(false);
        }
    }, [modelsError]);

    // 处理标题更改
    const handleTitleChange = useCallback(
        async (newTitle: string) => {
            setChatTitle(newTitle);

            // 如果是已有对话，调用 RenameChat API 更新标题
            if (currentChatUuid) {
                try {
                    await Service.RenameChat(currentChatUuid, newTitle);
                } catch (error) {
                    console.error('重命名聊天失败:', error);
                    message.error('重命名聊天失败');
                }
            }
        },
        [currentChatUuid]
    );

    // 处理复制消息
    const handleCopyMessage = useCallback((_content: string) => {
        // 复制功能已在MessageItem组件中实现
    }, []);

    // 处理停止生成
    const handleStopGeneration = useCallback(() => {
        if (abortController) {
            abortController.abort();
            setAbortController(null);
            setIsLoading(false);
            setIsStreaming(false);
            setShowLoadingMessage(false);
        }
    }, [abortController]);

    // 获取聊天消息
    const loadChatMessages = useCallback(async (chatUuid: string) => {
        // 当 chatUuid 为空的时候，表面此对话为新建对话
        if (!chatUuid) {
            setCurrentMessages([]);
            setInitialLoading(false);
            return;
        }
        
        // 显示加载动画
        setIsLoadingMessages(true);
        try {
            const response: MessageList | null = await Service.ChatMessages(chatUuid, 0, 50);
            console.log("response.messages:", response?.messages);
            setCurrentMessages(response?.messages || []);
        } catch (error) {
            console.error('获取聊天消息失败:', error);
            message.error('获取聊天消息失败');
            setCurrentMessages([]);
        } finally {
            setIsLoadingMessages(false);
            setInitialLoading(false);
        }
    }, []);

    // 当选择不同聊天时，加载消息
    useEffect(() => {
        loadChatMessages(currentChatUuid);
    }, [currentChatUuid, loadChatMessages]);

    // 处理模型更改
    const handleModelChange = useCallback((modelId: string) => {
        setSelectedModel(modelId);
    }, []);

    // 处理发送消息
    const handleSendMessage = useCallback(
        async (messageContent: string) => {
            if (!messageContent.trim() || isLoading) return;

            try {
                setIsLoading(true);
                setIsStreaming(true);
                setShowLoadingMessage(true);

                // 创建用户消息
                const userMessage = new Message();
                userMessage.role = RoleType.User;
                userMessage.content = messageContent.trim();

                // 创建AI消息占位符
                const assistantMessage = new Message();
                assistantMessage.role = RoleType.Assistant;
                assistantMessage.content = "";
                assistantMessage.reasoning_content = "";

                // 一次性更新消息列表（包含loading状态）
                const newMessages = [...currentMessages, userMessage, assistantMessage];
                setCurrentMessages(newMessages);

                console.log("[handleSendMessage] newMessages:", newMessages);

                // 创建新的 AbortController
                const newAbortController = new AbortController();
                setAbortController(newAbortController);

                // 用于累积增量数据的缓冲区
                let accumulatedContent = "";
                let accumulatedReasoningContent = "";

                await CompletionsUtils(currentChatUuid, selectedModel, userMessage, (message: Message) => {
                    if (message) {
                        console.log("message callback:", message);
                        // 隐藏loading消息
                        setShowLoadingMessage(false);

                        // 累积增量数据
                        accumulatedContent += message.content || '';
                        accumulatedReasoningContent += message.reasoning_content || '';

                        // 使用函数式更新确保获取最新状态，避免重复更新
                        setCurrentMessages(prev => {
                            const updatedMessages = [...prev];
                            const latestMsg = updatedMessages[updatedMessages.length - 1];
                            if (latestMsg && latestMsg.role === 'assistant') {
                                // 使用累积的内容更新，避免重复拼接
                                const newContent = accumulatedContent;
                                const newReasoningContent = accumulatedReasoningContent;
                                if (newContent !== latestMsg.content || newReasoningContent !== latestMsg.reasoning_content) {
                                    latestMsg.content = newContent;
                                    latestMsg.reasoning_content = newReasoningContent;
                                }
                            }
                            return updatedMessages;
                        });
                    }
                }, (error: string) => {
                    setIsLoading(false);
                    setIsStreaming(false);
                    setShowLoadingMessage(false);
                    setAbortController(null);
                    console.error('发送消息失败:', error);
                    message.error('发送消息失败');
                }, (chatUuid: string) => {
                    setIsLoading(false);
                    setIsStreaming(false);
                    setShowLoadingMessage(false);
                    setAbortController(null);
                    if (currentChatUuid === "") {
                        setCurrentChatUuid(chatUuid);
                        navigate(`/chat/${chatUuid}`, { replace: true });
                    }
                }, newAbortController);
            } catch (error) {
                setIsLoading(false);
                setIsStreaming(false);
                setShowLoadingMessage(false);
                setAbortController(null);
                console.error('发送消息失败:', error);
                message.error('发送消息失败');
            }
        },
        [
            isLoading,
            selectedModel,
            currentChatUuid,
            currentMessages,
            isStreaming,
            navigate,
        ]
    );

    // 处理删除消息
    const handleDeleteMessage = useCallback(
        async (messageId: string) => {
            // TODO: 实现删除消息功能
            console.log('删除消息:', messageId);
        },
        [currentChatUuid]
    );

    // 处理消息重新生成
    const handleRegenerateMessage = useCallback(
        async (messageId: string) => {
            // TODO: 实现重新生成消息功能
            console.log('重新生成消息:', messageId);
        },
        [currentMessages]
    );

    return (
        <Chat
            className={className}
            standalone={true}
            initialLoading={initialLoading || isLoadingModels}
            chatTitle={chatTitle}
            chatUuid={currentChatUuid}
            currentMessages={currentMessages}
            isLoading={isLoadingMessages}
            showLoadingMessage={showLoadingMessage}
            selectedModel={selectedModel}
            availableModels={availableModels}
            isMobile={isMobile}
            isGenerating={isStreaming}
            onTitleChange={handleTitleChange}
            onSendMessage={handleSendMessage}
            onStopGeneration={handleStopGeneration}
            onModelChange={handleModelChange}
            onCopyMessage={handleCopyMessage}
            onDeleteMessage={handleDeleteMessage}
            onRegenerateMessage={handleRegenerateMessage}
        />
    );
};

export default ChatPage;