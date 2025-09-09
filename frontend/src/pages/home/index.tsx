import React, {useCallback, useEffect, useRef, useState} from 'react';
import {BackTop, Layout, message} from 'antd';
import {useNavigate, useParams} from 'react-router-dom';
import Index from './sidebar';
import {schema, view_models} from "../../../wailsjs/go/models"; // 修复导入路径
import {useViewportHeight} from '@/hooks/useViewportHeight';
import {useModels} from '@/hooks/useModels';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import {ChatMessages, Completions} from "../../../wailsjs/go/service/Service";
import {EventsOn} from "../../../wailsjs/runtime";
import styles from './index.module.scss'; // 添加样式导入

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
    const [currentMessages, setCurrentMessages] = useState<schema.Message[]>([]);
    const currentMessagesRef = useRef<schema.Message[]>([]);
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [refreshChatList, setRefreshChatList] = useState<(() => void) | null>(
        null
    );
    const [updateChatTitle, setUpdateChatTitle] = useState<((chatUuid: string, newTitle: string) => void) | null >(null);
    // Safari兼容性：添加强制重新渲染状态
    const [forceRerender, setForceRerender] = useState(0);

    // 使用视口高度检测 Hook
    const {isMobile} = useViewportHeight();

    // 使用模型获取 Hook
    const {
        models: availableModels,
        isLoading: isLoadingModels,
        error: modelsError,
    } = useModels();

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

    useEffect(() => {
        currentMessagesRef.current = currentMessages;
    });

    // 处理标题更改
    const handleTitleChange = useCallback(
        (newTitle: string) => {
            setChatTitle(newTitle);
            // 如果是已有对话，优先使用精确更新，否则刷新整个列表
            if (currentChatUuid && updateChatTitle) {
                updateChatTitle(currentChatUuid, newTitle);
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
            return;
        }
        // 显示加载动画
        setIsLoadingMessages(true);
        try {
            const response: view_models.MessageList = await ChatMessages(chatUuid, 0, 50);
            console.log("response.messages:", response.messages);
            setCurrentMessages(response.messages);
        } catch (error) {
            // todo 显示”加载历史消息错误“
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
        async (messageContent: string) => {
            if (!messageContent.trim() || isLoading) return;

            try {
                setIsLoading(true);
                setIsStreaming(true);

                // 创建用户消息
                const userMessage = new schema.Message();
                userMessage.role = "user";
                userMessage.content = messageContent.trim();
                
                // 创建AI消息占位符
                const assistantMessage = new schema.Message();
                assistantMessage.role = "assistant";
                assistantMessage.content = "";
                assistantMessage.reasoning_content = "";

                // 一次性更新消息列表
                const newMessages = [...currentMessages, userMessage, assistantMessage];
                setCurrentMessages(newMessages);

                console.log("[handleSendMessage] newMessages:", newMessages)

                // 用于跟踪是否已经接收到第一个响应
                let hasReceivedFirstResponse = false;

                // 调用Completions API
                Completions(currentChatUuid, selectedModel, userMessage).then((emitKey:string)=>{
                    if (currentChatUuid == "") {
                        setCurrentChatUuid(emitKey);
                        navigate(`/home/${emitKey}`, { replace: true });
                    }
                    const cancel = EventsOn(emitKey, (responseMessage?: schema.Message) => {
                        console.log("responseMessage:", responseMessage)
                        // 第一次接收到内容时立即隐藏loading状态
                        if (!hasReceivedFirstResponse) {
                            hasReceivedFirstResponse = true;
                            setIsLoading(false);
                        }
                        if (responseMessage) {
                            // 使用函数式更新确保获取最新状态
                            setCurrentMessages(prev => {
                                const updatedMessages = [...prev];
                                const latestMsg = updatedMessages[updatedMessages.length - 1];
                                if (latestMsg && latestMsg.role === 'assistant') {
                                    latestMsg.content = latestMsg.content + responseMessage.content;
                                    latestMsg.reasoning_content = (latestMsg.reasoning_content || '') + (responseMessage.reasoning_content || '');
                                }
                                return updatedMessages;
                            });
                        }
                        if (responseMessage?.response_meta?.finish_reason) {
                            cancel();
                            setIsStreaming(false);
                            setIsLoading(false);
                        }
                    });
                }).catch((err)=>{
                    console.error('发送消息失败:', err);
                    message.error('发送消息失败');
                    setIsLoading(false);
                    setIsStreaming(false);
                });

            } catch (error) {
                console.error('发送消息失败:', error);
                message.error('发送消息失败');
                setIsLoading(false);
                setIsStreaming(false);
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
                collapsedWidth={0}
                collapsed={isSidebarCollapsed}
                trigger={null}
                collapsible
                // Safari内核兼容性：添加transform样式确保正确隐藏
                style={{
                    transform: isSidebarCollapsed ? 'translateX(-100%)' : 'translateX(0)',
                    transition: 'transform 0.3s ease-in-out',
                }}
            >
                <Index
                    onNewChat={handleNewChat}
                    onChatSelect={handleChatSelect}
                    onRegisterRefreshCallback={handleSetRefreshChatList}
                    onRegisterUpdateTitleCallback={handleSetUpdateChatTitle}
                    currentChatUuid={currentChatUuid}
                    isSidebarCollapsed={isSidebarCollapsed}
                    onToggleSidebar={handleToggleSidebar}
                />
            </Sider>
            <Layout className={styles.mainLayout}>
                <Content className={styles.mainContent}>
                    <Chat
                        chatTitle={chatTitle}
                        chatUuid={currentChatUuid}
                        currentMessages={currentMessages}
                        isLoading={isLoadingMessages}
                        selectedModel={selectedModel}
                        availableModels={availableModels}
                        isMobile={isMobile}
                        isSidebarCollapsed={isSidebarCollapsed}
                        isGenerating={isStreaming}
                        onTitleChange={handleTitleChange}
                        onSendMessage={handleSendMessage}
                        onStopGeneration={handleStopGeneration}
                        onModelChange={handleModelChange}
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