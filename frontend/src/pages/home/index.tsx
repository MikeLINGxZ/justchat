import React, {useCallback, useEffect, useState} from 'react';
import {BackTop, Layout, message} from 'antd';
import {useNavigate, useParams} from 'react-router-dom';
import Index from './sidebar';
import {
    ChatMessagePartType,
    Message,
    MessageInputPart,
    MessageInputImage,
    MessageInputAudio,
    MessageInputVideo,
    MessageInputFile,
    RoleType
} from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import {useViewportHeight} from '@/hooks/useViewportHeight';
import {useModelStore} from '@/stores/modelStore';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import styles from './index.module.scss';
import {CompletionsUtils} from "@/utils/completions.ts";
import {
    File,
    MessageList,
    MessagePkg
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/index.ts";

const {Content, Sider} = Layout;

interface ChatPageProps {
    className?: string;
}

const ChatPage: React.FC<ChatPageProps> = ({className}) => {
    // 获取路由参数和导航函数
    const {chatUuid: urlChatUuid} = useParams<{ chatUuid?: string }>();
    const navigate = useNavigate();

    // 本地状态管理
    const [currentChatUuid, setCurrentChatUuid] = useState<string>(urlChatUuid || ''); // 空字符串表示新对话
    const [currentMessages, setCurrentMessages] = useState<Message[]>([]);
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [refreshChatList, setRefreshChatList] = useState<(() => void) | null>(
        null
    );
    const [updateChatTitle, setUpdateChatTitle] = useState<((chatUuid: string, newTitle: string) => void) | null>(null);
    // Safari兼容性：添加强制重新渲染状态
    const [forceRerender, setForceRerender] = useState(0);
    // 添加loading消息状态
    const [showLoadingMessage, setShowLoadingMessage] = useState(false);

    // 使用视口高度检测 Hook
    const {isMobile} = useViewportHeight();

    // 使用模型 Store
    const {
        models: availableModels,
        isLoading: isLoadingModels,
        error: modelsError,
        refetch: refetchModels,
    } = useModelStore();

    // 聊天相关状态
    const [chatTitle, setChatTitle] = useState('新建对话');
    const [isLoading, setIsLoading] = useState(false);
    const [selectedModel, setSelectedModel] = useState('');
    const [isStreaming, setIsStreaming] = useState(false);
    const [abortController, setAbortController] = useState<AbortController | null>(null);

    // 移动端默认隐藏侧边栏
    useEffect(() => {
        if (isMobile) {
            setIsSidebarCollapsed(true);
        } else {
            // Safari内核兼容性：从移动端切换回桌面端时，需要强制重置transform属性
            // 添加延迟重新渲染机制，确保Safari正确应用新的CSS规则
            const timer = setTimeout(() => {
                // 强制触发组件重新渲染
                setIsSidebarCollapsed(prev => prev);
                setForceRerender(prev => prev + 1);
            }, 100);

            return () => clearTimeout(timer);
        }
    }, [isMobile]);

    // 设置页面标题
    useEffect(() => {
        document.title = 'AI聊天 - Lemon Tea';
    }, []);

    // 初始化时获取模型列表
    useEffect(() => {
        refetchModels();
    }, [refetchModels]);

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
        }
    }, [availableModels, selectedModel]);

    // 显示模型加载错误
    useEffect(() => {
        if (modelsError) {
            message.error(`获取模型列表失败: ${modelsError}`);
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
                    // 如果是已有对话，优先使用精确更新，否则刷新整个列表
                    if (updateChatTitle) {
                        updateChatTitle(currentChatUuid, newTitle);
                    } else if (refreshChatList) {
                        refreshChatList();
                    }
                } catch (error) {
                    console.error('重命名聊天失败:', error);
                    message.error('重命名聊天失败');
                }
            } else if (refreshChatList) {
                refreshChatList();
            }
        },
        [currentChatUuid, updateChatTitle, refreshChatList]
    );

    // 处理消息复制
    const handleCopyMessage = useCallback((_content: string) => {
        // 复制功能已在MessageItem组件中实现
    }, []);

    // 处理停止生成
    const handleStopGeneration = useCallback(() => {
        // todo
    }, [abortController]);

    // 获取聊天消息
    const loadChatMessages = useCallback(async (chatUuid: string) => {
        // 当 chatUuid 为空的时候，表面此对话为新建对话
        if (!chatUuid) {
            setCurrentMessages([]);
            setChatTitle('新建对话');
            return;
        }
        // 显示加载动画
        setIsLoadingMessages(true);
        try {
            const response: MessageList | null = await Service.ChatMessages(chatUuid, 0, 50);
            console.log("response.messages:", response?.messages);
            setCurrentMessages(response?.messages!);
            
            // 加载消息成功后，获取聊天信息并设置标题
            try {
                const chatListResponse = await Service.ChatList(0, 100, null, false);
                if (chatListResponse?.lists) {
                    const chat = chatListResponse.lists.find(c => c.uuid === chatUuid);
                    if (chat && chat.title) {
                        setChatTitle(chat.title);
                    }
                }
            } catch (chatError) {
                console.error('获取聊天信息失败:', chatError);
                // 如果获取聊天信息失败，不影响消息加载，只记录错误
            }
        } catch (error) {
            // todo 显示"加载历史消息错误"
            console.error('获取聊天消息失败:', error);
            message.error('获取聊天消息失败');
            setCurrentMessages([]);
        } finally {
            setIsLoadingMessages(false);
        }
    }, []);

    // 当选择不同聊天时，加载消息
    useEffect(() => {
        loadChatMessages(currentChatUuid);
    }, [currentChatUuid, loadChatMessages]);

    // handleToggleSidebar 展示/隐藏侧边菜单
    const handleToggleSidebar = () => {
        setIsSidebarCollapsed(!isSidebarCollapsed);
    };

    // 处理新建对话
    const handleNewChat = useCallback(() => {
        setCurrentChatUuid(''); // 设置为空字符串表示新对话
        setCurrentMessages([]);
        setChatTitle('新对话');
        // 更新URL为新对话状态
        navigate('/home', {replace: true});
        // 移动端新建对话后自动隐藏侧边栏
        if (isMobile) {
            setIsSidebarCollapsed(true);
        }
    }, [isMobile, navigate]);

    // 处理对话选择
    const handleChatSelect = useCallback(
        (chatUuid: string, chatTitle?: string) => {
            setCurrentChatUuid(chatUuid);
            setChatTitle(chatTitle || '新建对话');
            // 更新URL但不刷新页面
            navigate(`/home/${chatUuid}`, {replace: true});
            // 移动端选择对话后自动隐藏侧边栏
            if (isMobile) {
                setIsSidebarCollapsed(true);
            }
        },
        [isMobile, navigate]
    );

    // 处理模型更改
    const handleModelChange = useCallback((modelId: string) => {
        setSelectedModel(modelId);
    }, []);

    // 处理发送消息
    const handleSendMessage = useCallback(
        async (messageContent: string,files: File[]) => {
            if (!messageContent.trim() || isLoading) return;

            try {
                setIsLoading(true);
                setIsStreaming(true);
                setShowLoadingMessage(true); // 显示loading消息

                // 创建用户消息
                const userMessage = new Message();
                userMessage.role = RoleType.User;
                userMessage.content = messageContent.trim();
                
                // 构建 user_input_multi_content，包含文本和文件
                const userInputMultiContent: MessageInputPart[] = [
                    {
                        type: ChatMessagePartType.ChatMessagePartTypeText,
                        text: messageContent.trim(),
                    }
                ];
                
                // 添加文件到 user_input_multi_content
                for (const file of files) {
                    const extra = {
                        name: file.name,
                        path: file.file_path,
                        mime_type: file.mine_type,
                    };
                    
                    let part: MessageInputPart | null = null;
                    
                    switch (file.chat_message_part_type) {
                        case ChatMessagePartType.ChatMessagePartTypeImageURL:
                            part = new MessageInputPart({
                                type: ChatMessagePartType.ChatMessagePartTypeImageURL,
                                image: new MessageInputImage({
                                    extra: extra,
                                    mime_type: file.mine_type,
                                })
                            });
                            break;
                        case ChatMessagePartType.ChatMessagePartTypeAudioURL:
                            part = new MessageInputPart({
                                type: ChatMessagePartType.ChatMessagePartTypeAudioURL,
                                audio: new MessageInputAudio({
                                    extra: extra,
                                    mime_type: file.mine_type,
                                })
                            });
                            break;
                        case ChatMessagePartType.ChatMessagePartTypeVideoURL:
                            part = new MessageInputPart({
                                type: ChatMessagePartType.ChatMessagePartTypeVideoURL,
                                video: new MessageInputVideo({
                                    extra: extra,
                                    mime_type: file.mine_type,
                                })
                            });
                            break;
                        case ChatMessagePartType.ChatMessagePartTypeFileURL:
                            part = new MessageInputPart({
                                type: ChatMessagePartType.ChatMessagePartTypeFileURL,
                                file: new MessageInputFile({
                                    extra: extra,
                                    mime_type: file.mine_type,
                                    name: file.name,
                                })
                            });
                            break;
                    }
                    
                    if (part) {
                        userInputMultiContent.push(part);
                    }
                }
                
                userMessage.user_input_multi_content = userInputMultiContent;

                // 创建AI消息占位符
                const assistantMessage = new Message();
                assistantMessage.role = RoleType.Assistant;
                assistantMessage.content = "";
                assistantMessage.reasoning_content = "";

                // 一次性更新消息列表（包含loading状态）
                const newMessages = [...currentMessages, userMessage, assistantMessage];
                setCurrentMessages(newMessages);

                console.log("[handleSendMessage] newMessages:", newMessages)

                // 创建新的 AbortController
                const newAbortController = new AbortController();
                setAbortController(newAbortController);

                // 用于累积增量数据的缓冲区
                let accumulatedContent = "";
                let accumulatedReasoningContent = "";
                let messagePkg:MessagePkg = {
                    chatUuid: currentChatUuid,
                    model: selectedModel,
                    content: messageContent.trim(),
                    files: files,
                }
                await CompletionsUtils(messagePkg, (message: Message) => {
                    if (message) {
                        console.log("message callback:", message)
                        // 隐藏loading消息
                        setShowLoadingMessage(false);

                        // 累积增量数据
                        accumulatedContent += message.content || '';
                        accumulatedReasoningContent += message.reasoning_content || '';
                        console.log("message.content:",message.content)
                        console.log("message.reasoning_content:",message.reasoning_content)
                        console.log("accumulatedContent:",accumulatedContent)
                        console.log("accumulatedReasoningContent:",accumulatedReasoningContent)

                        // 使用函数式更新确保获取最新状态，避免重复更新
                        setCurrentMessages(prev => {
                            const updatedMessages = [...prev];
                            const latestMsg = updatedMessages[updatedMessages.length - 1];
                            if (latestMsg && latestMsg.role === 'assistant') {
                                // 使用累积的内容更新，避免重复拼接
                                const newContent = accumulatedContent;
                                const newReasoningContent = accumulatedReasoningContent;
                                console.log("newContent:", newContent);
                                if (newContent !== latestMsg.content || newReasoningContent !== latestMsg.reasoning_content) {
                                    latestMsg.content = newContent;
                                    latestMsg.reasoning_content = newReasoningContent;
                                }
                                // 设置 isStreaming 为 true，表示正在生成中
                                (latestMsg as any).isStreaming = true;
                            }
                            return updatedMessages;
                        });
                    }
                }, (error: string) => {
                    setIsLoading(false);
                    setIsStreaming(false);
                    setShowLoadingMessage(false); // 隐藏loading消息
                    // 更新消息的 isStreaming 状态为 false
                    setCurrentMessages(prev => {
                        const updatedMessages = [...prev];
                        const latestMsg = updatedMessages[updatedMessages.length - 1];
                        if (latestMsg && latestMsg.role === 'assistant') {
                            (latestMsg as any).isStreaming = false;
                        }
                        return updatedMessages;
                    });
                    setAbortController(null); // 清理 AbortController
                    console.error('发送消息失败:', error);
                    message.error('发送消息失败');
                }, (chatUuid: string) => {
                    setIsLoading(false);
                    setIsStreaming(false);
                    setShowLoadingMessage(false); // 隐藏loading消息
                    // 更新消息的 isStreaming 状态为 false，表示生成完成
                    setCurrentMessages(prev => {
                        const updatedMessages = [...prev];
                        const latestMsg = updatedMessages[updatedMessages.length - 1];
                        if (latestMsg && latestMsg.role === 'assistant') {
                            (latestMsg as any).isStreaming = false;
                        }
                        return updatedMessages;
                    });
                    setAbortController(null); // 清理 AbortController
                    if (currentChatUuid == "") {
                        setCurrentChatUuid(chatUuid);
                        // 刷新聊天列表
                        if (refreshChatList) {
                            refreshChatList();
                        }
                        navigate(`/home/${chatUuid}`, {replace: true});
                    }
                }, newAbortController)
            } catch (error) {
                setIsLoading(false);
                setIsStreaming(false);
                setShowLoadingMessage(false); // 隐藏loading消息
                // 更新消息的 isStreaming 状态为 false
                setCurrentMessages(prev => {
                    const updatedMessages = [...prev];
                    const latestMsg = updatedMessages[updatedMessages.length - 1];
                    if (latestMsg && latestMsg.role === 'assistant') {
                        (latestMsg as any).isStreaming = false;
                    }
                    return updatedMessages;
                });
                setAbortController(null); // 清理 AbortController
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
            refreshChatList,
            navigate,
        ]
    );

    // 处理删除消息
    const handleDeleteMessage = useCallback(
        async (messageId: string) => {
            // todo
        },
        [currentChatUuid]
    );

    // 处理删除聊天
    const handleDeleteChat = useCallback(
        async (chatUuid: string) => {
            try {
                await Service.DeleteChat(chatUuid);
                // 如果删除的是当前聊天，导航到新聊天页面
                if (chatUuid === currentChatUuid) {
                    handleNewChat();
                }
                // 刷新聊天列表
                if (refreshChatList) {
                    refreshChatList();
                }
                message.success('聊天删除成功');
            } catch (error) {
                console.error('删除聊天失败:', error);
                message.error('删除聊天失败');
            }
        },
        [currentChatUuid, handleNewChat, refreshChatList]
    );

    // 处理消息重新生成
    const handleRegenerateMessage = useCallback(
        async (messageId: string) => {
            // todo
        },
        [currentMessages]
    );

    // 设置刷新聊天列表的回调
    const handleSetRefreshChatList = useCallback((refreshFn: () => void) => {
        setRefreshChatList(() => refreshFn);
    }, []);

    // 设置更新聊天标题的回调
    const handleSetUpdateChatTitle = useCallback(
        (updateFn: (chatUuid: string, newTitle: string) => void) => {
            setUpdateChatTitle(() => updateFn);
        },
        []
    );

    return (
        <Layout className={`${className || ''} ${styles.chatLayout}`}>
            <Sider
                className={`${styles.sidebar} ${
                    isSidebarCollapsed ? styles.collapsed : ''
                }`}
                width={280}
                collapsedWidth={isMobile ? 0 : 50}
                collapsed={isSidebarCollapsed}
                trigger={null}
                collapsible
            >
                <Index
                    onNewChat={handleNewChat}
                    onChatSelect={handleChatSelect}
                    onRegisterRefreshCallback={handleSetRefreshChatList}
                    onRegisterUpdateTitleCallback={handleSetUpdateChatTitle}
                    onDeleteChat={handleDeleteChat}
                    currentChatUuid={currentChatUuid}
                    isSidebarCollapsed={isSidebarCollapsed}
                    onToggleSidebar={handleToggleSidebar}
                />
            </Sider>
            <Layout className={styles.mainLayout}>
                <Content className={styles.mainContent} hidden={isMobile && !isSidebarCollapsed}>
                    <Chat
                        standalone={false}
                        initialLoading={isLoadingMessages}
                        chatTitle={chatTitle}
                        chatUuid={currentChatUuid}
                        currentMessages={currentMessages}
                        isLoading={isLoadingMessages}
                        showLoadingMessage={showLoadingMessage}
                        selectedModel={selectedModel}
                        availableModels={availableModels}
                        isMobile={isMobile}
                        isSidebarCollapsed={isSidebarCollapsed}
                        isGenerating={isStreaming}
                        onTitleChange={handleTitleChange}
                        onSendMessage={handleSendMessage}
                        onStopGeneration={handleStopGeneration}
                        onModelChange={handleModelChange}
                        onModelSelectorClick={refetchModels}
                        onToggleSidebar={handleToggleSidebar}
                        onCopyMessage={handleCopyMessage}
                        onDeleteMessage={handleDeleteMessage}
                        onRegenerateMessage={handleRegenerateMessage}
                    />
                </Content>
            </Layout>
            <BackTop/>
        </Layout>
    );
};

export default ChatPage;