import React, { useCallback, useEffect, useState } from 'react';
import { BackTop, Layout, message } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import Index from './sidebar';
import { schema } from "../../../wailsjs/go/models"; // 修复导入路径
import { useViewportHeight } from '@/hooks/useViewportHeight';
import { useModels } from '@/hooks/useModels';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import {ChatMessages, Completions} from "../../../wailsjs/go/service/Service";
import {view_models} from "../../../wailsjs/go/models";
import {EventsOn} from "../../../wailsjs/runtime";
import styles from './index.module.scss'; // 添加样式导入

const { Content, Sider } = Layout;

interface ChatPageProps {
  className?: string;
}

const ChatPage: React.FC<ChatPageProps> = ({ className }) => {
  // 获取路由参数和导航函数
  const { chatUuid: urlChatUuid } = useParams<{ chatUuid?: string }>();
  const navigate = useNavigate();

  // 本地状态管理
  const [currentChatUuid, setCurrentChatUuid] = useState<string>(urlChatUuid || ''); // 空字符串表示新对话
  const [currentMessages, setCurrentMessages] = useState<schema.Message[]>([]); // 修改为 schema.Message[]
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [refreshChatList, setRefreshChatList] = useState<(() => void) | null>(
    null
  );
  const [updateChatTitle, setUpdateChatTitle] = useState<
    ((chatUuid: string, newTitle: string) => void) | null
  >(null);
  // Safari兼容性：添加强制重新渲染状态
  const [forceRerender, setForceRerender] = useState(0);

  // 使用视口高度检测 Hook
  const { isMobile } = useViewportHeight();

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
  const [abortController, setAbortController] =
    useState<AbortController | null>(null);

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
    if (abortController) {
      abortController.abort();
      setAbortController(null);
      setIsStreaming(false);
      setIsLoading(false);

      // 标记最后一条AI消息为完成状态
      setCurrentMessages(prev => {
        const newMessages = [...prev];
        const lastIndex = newMessages.length - 1;
        if (
          lastIndex >= 0 &&
          newMessages[lastIndex].role === 'assistant' &&
          (newMessages[lastIndex] as any).isStreaming // Message 类没有 isStreaming 属性
        ) {
          // 创建一个新的 Message 实例
          const updatedMessage = new schema.Message();
          Object.assign(updatedMessage, newMessages[lastIndex]);
          updatedMessage.content = newMessages[lastIndex].content + '\n\n[生成已停止]';
          // 移除 isStreaming 属性
          delete (updatedMessage as any).isStreaming;
          
          newMessages[lastIndex] = updatedMessage;
        }
        return newMessages;
      });
    }
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

    // 获取当前对话历史消息 (模拟实现)
    try {
      // 修复类型问题
      const response: view_models.MessageList = await ChatMessages(chatUuid,0,50);
      console.log("response.messages:",response.messages);
      setCurrentMessages(response.messages);
    } catch (error) {
      // todo 显示”加载历史消息错误“
      console.error('获取聊天消息失败:', error);
      message.error('获取聊天消息失败');
      setCurrentMessages([]);
    } finally {
      // 关闭加载动画
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
    navigate('/home', { replace: true });
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
      navigate(`/home/${chatUuid}`, { replace: true });
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
      if (!messageContent.trim() || isStreaming || isLoading) return;

      try {
        setIsLoading(true);
        setIsStreaming(true);
        
        // 创建用户消息
        const userMessage = new schema.Message(); // 修改为 schema.Message
        userMessage.role = 'user';
        userMessage.content = messageContent.trim();
        
        // 添加用户消息到聊天列表
        const updatedMessages = [...currentMessages, userMessage];
        setCurrentMessages(updatedMessages);
        
        // 创建AI消息占位符
        const aiMessage = new schema.Message(); // 修改为 schema.Message
        aiMessage.role = 'assistant';
        aiMessage.content = '';
        (aiMessage as any).isStreaming = true; // 添加 isStreaming 属性
        
        setCurrentMessages(prev => [...prev, aiMessage]);
        
        try {
          // 调用Completions API
          const emitKey: string = await Completions(currentChatUuid, selectedModel, userMessage);
          setCurrentChatUuid(emitKey);
          
          // 监听流式响应
          EventsOn(emitKey, (responseMessage?: schema.Message) => {
              console.log("responseMessage:",responseMessage)
            if (responseMessage) {
              // 更新AI消息
              setCurrentMessages(prev => {
                const newMessages = [...prev];
                const lastIndex = newMessages.length - 1;
                if (lastIndex >= 0 && newMessages[lastIndex].role === 'assistant') {
                  newMessages[lastIndex] = responseMessage;
                }
                return newMessages;
              });
            } else {
              // 流式响应结束
              setIsStreaming(false);
              setIsLoading(false);
            }
          });
        } catch (error) {
          console.error('API调用失败:', error);
          message.error('消息发送失败');
          setIsStreaming(false);
          setIsLoading(false);
          
          // 移除AI消息占位符
          setCurrentMessages(prev => prev.slice(0, -1));
        }
      } catch (error) {
        console.error('发送消息失败:', error);
        message.error('发送消息失败');
        setIsLoading(false);
        setIsStreaming(false);
      }
    },
    [currentChatUuid, currentMessages, isStreaming, isLoading, selectedModel]
  );

  // 处理删除消息
  const handleDeleteMessage = useCallback(
    async (messageId: string) => {
      try {
        // 这里应该调用API删除消息
        console.log('删除消息:', messageId);
        message.success('消息删除成功');
      } catch (error) {
        console.error('删除消息失败:', error);
        message.error('删除消息失败');
      }
    },
    [currentChatUuid]
  );

  // 处理消息重新生成
  const handleRegenerateMessage = useCallback(
    async (messageId: string) => {
      // Message 类没有 id 属性，所以我们使用索引来查找消息
      const messageIndex = currentMessages.findIndex(
        (_, index) => index.toString() === messageId // 使用索引作为 messageId
      );
      if (messageIndex === -1) return;

      // 移除该消息及之后的所有消息
      const newMessages = currentMessages.slice(0, messageIndex);
      setCurrentMessages(newMessages);
      setIsLoading(true);
      setIsStreaming(true);

      // 创建AbortController用于中断请求
      const controller = new AbortController();
      setAbortController(controller);

      // 创建AI消息占位符
      const aiMessage = new schema.Message(); // 修改为 schema.Message
      aiMessage.role = 'assistant';
      aiMessage.content = '';
      (aiMessage as any).isStreaming = true; // 标记为正在流式输入
      
      setCurrentMessages(prev => [...prev, aiMessage]);

      try {
        // 模拟重新生成延迟
        await new Promise(resolve => setTimeout(resolve, 500));
        
        // 模拟回复内容
        const mockResponse = '这是一个重新生成的回复，用于展示 UI 功能。在实际应用中，这里会调用真实的 AI 接口重新生成内容。';
        
        // 模拟打字机效果
        let currentIndex = 0;
        const interval = setInterval(() => {
          if (currentIndex <= mockResponse.length && !controller.signal.aborted) {
            const partialContent = mockResponse.slice(0, currentIndex);
            
            setCurrentMessages(prev => {
              const updatedMessages = [...prev];
              const lastIndex = updatedMessages.length - 1;
              if (lastIndex >= 0 && updatedMessages[lastIndex].role === 'assistant') {
                // 创建一个新的 Message 实例
                const updatedMessage = new schema.Message();
                Object.assign(updatedMessage, updatedMessages[lastIndex]);
                updatedMessage.content = partialContent;
                // 更新 isStreaming 属性
                (updatedMessage as any).isStreaming = currentIndex < mockResponse.length;
                
                updatedMessages[lastIndex] = updatedMessage;
              }
              return updatedMessages;
            });
            
            currentIndex++;
          } else {
            clearInterval(interval);
            if (!controller.signal.aborted) {
              setIsStreaming(false);
              setIsLoading(false);
            }
          }
        }, 30);
      } catch (error) {
        console.error('重新生成消息失败:', error);
        message.error('重新生成消息失败');
        setIsStreaming(false);
        setIsLoading(false);
        
        // 移除AI消息占位符
        setCurrentMessages(prev => prev.slice(0, -1));
      }
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
      <BackTop />
    </Layout>
  );
};

export default ChatPage;