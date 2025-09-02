import React, { useCallback, useEffect, useState } from 'react';
import { BackTop, Layout, message } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import Index from './sidebar';
import type { Message } from '@/types';
import { useViewportHeight } from '@/hooks/useViewportHeight';
import { useModels } from '@/hooks/useModels';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import type {
  ServerCompletionsRequest,
  ServerGetChatMessagesResponse,
} from '@/api/chatClient';
import { chatClient } from '@/api/chatClient';
import type { CommonMessage } from '@/api/service/chat/Api';
import { packageErr } from '@/utils/converErr.ts';
import {
  ConvertApiMessageToMessage,
  ConvertMessageToApiMessage,
} from '@/utils/message.ts';

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
  const [currentMessages, setCurrentMessages] = useState<Message[]>([]);
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
          newMessages[lastIndex].isStreaming
        ) {
          newMessages[lastIndex] = {
            ...newMessages[lastIndex],
            isStreaming: false,
            content: newMessages[lastIndex].content + '\n\n[生成已停止]',
          };
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

    // 获取当前对话历史消息
    try {
      // 请求历史消息接口
      const response: ServerGetChatMessagesResponse =
        await chatClient.getChatMessages(chatUuid, {
          offset: '0',
          limit: '50',
        });
      // 有消息返回则转换渲染，没有则则重置消息为空
      if (response.messages) {
        const convertedMessages = response.messages.map(
          ConvertApiMessageToMessage
        );
        setCurrentMessages(convertedMessages);
      } else {
        setCurrentMessages([]);
      }
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
      // 设置对话标题
      if (chatTitle) {
        setChatTitle(chatTitle);
      }
      // 更新URL以反映当前选中的聊天
      navigate(`/home/${chatUuid}`, { replace: true });
      // H5移动端点击对话后自动隐藏侧边栏
      if (isMobile) {
        setIsSidebarCollapsed(true);
      }
    },
    [isMobile, navigate]
  );

  // 处理聊天列表刷新回调注册
  const handleRegisterRefreshCallback = useCallback((callback: () => void) => {
    setRefreshChatList(() => callback);
  }, []);

  // 处理标题更新回调注册
  const handleRegisterUpdateTitleCallback = useCallback(
    (callback: (chatUuid: string, newTitle: string) => void) => {
      setUpdateChatTitle(() => callback);
    },
    []
  );

  // 聊天功能 - 现在接收消息内容作为参数
  const handleSendMessage = useCallback(
    async (messageContent: string) => {
      if (!messageContent.trim() || isLoading || isStreaming) return;

      // 创建要发送的消息
      const userMessage: Message = {
        role: 'user',
        content: messageContent,
      };

      // 添加用户消息
      setCurrentMessages(prev => [...prev, userMessage]);
      setIsLoading(true);
      setIsStreaming(true);

      // 创建AbortController用于中断请求
      const controller = new AbortController();
      setAbortController(controller);

      // 预先初始化ai回复的消息
      const aiMessage: Message = {
        id: '',
        role: 'assistant',
        content: '',
        isStreaming: true,
      };
      setCurrentMessages(prev => [...prev, aiMessage]);

      // 用于跟踪是否已经接收到第一个响应
      let hasReceivedFirstResponse = false;
      try {
        // 构建请求消息历史
        const messagesHistory = [...currentMessages, userMessage];
        const apiMessages: CommonMessage[] = messagesHistory.map(
          ConvertMessageToApiMessage
        );

        // 构建请求参数
        const completionRequest: ServerCompletionsRequest = {
          model: selectedModel,
          messages: apiMessages,
          nonStandard: true, // 固定为 true
          chatUuid: currentChatUuid || '', // 空字符串表示新建对话
          temperature: 0.7,
          maxTokens: 2000,
          stream: true,
        };

        // 调用流式 Completions 接口
        await chatClient.createChatCompletionStream(
          completionRequest,
          // onMessage: 接收到新内容时的回调
          (content: string, reasoningContent?: string) => {
            console.log('收到流式内容:', { content, reasoningContent });

            // 第一次接收到内容时立即隐藏loading状态
            if (!hasReceivedFirstResponse) {
              hasReceivedFirstResponse = true;
              setIsLoading(false);
            }

            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  content: newMessages[lastIndex].content + (content || ''),
                  reasoningContent: reasoningContent
                    ? (newMessages[lastIndex].reasoningContent || '') +
                      reasoningContent
                    : newMessages[lastIndex].reasoningContent,
                };
              }
              return newMessages;
            });
          },
          // onError: 接收到错误时的回调
          (error: string) => {
            // 将错误信息作为AI消息显示在聊天界面中
            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  content: packageErr(`${error}`),
                  isStreaming: false,
                };
              }
              return newMessages;
            });
          },
          // onComplete: 流式响应完成时的回调
          (chatUuid?: string) => {
            // 标记消息流式输入完成
            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  isStreaming: false,
                };
              }
              return newMessages;
            });

            // 重置流式状态
            setIsStreaming(false);
            setAbortController(null);

            // 如果是新对话（currentChatUuid为空）且有chatUuid，更新chatUuid和标题
            if (!currentChatUuid && chatUuid) {
              setCurrentChatUuid(chatUuid);
              // 更新URL以反映新的聊天UUID
              navigate(`/home/${chatUuid}`, { replace: true });

              // 使用ChatTitle接口获取对话标题并设置
              (async () => {
                try {
                  const titleResponse = await chatClient.getChatTitle(chatUuid);
                  setChatTitle(titleResponse.title || '新的对话');
                } catch (error) {
                  console.error('获取对话标题失败:', error);
                  setChatTitle('新的对话');
                }
              })();

              // 更新侧边栏历史对话列表
              if (refreshChatList) {
                refreshChatList();
              }
            }
          },
          controller
        );
      } catch (error) {
        console.error('Send message error:', error);
        // 将错误信息作为AI消息显示在聊天界面中
        setCurrentMessages(prev => {
          const newMessages = [...prev];
          const lastIndex = newMessages.length - 1;
          if (lastIndex >= 0 && newMessages[lastIndex].role === 'assistant') {
            newMessages[lastIndex] = {
              ...newMessages[lastIndex],
              content: `错误: ${error instanceof Error ? error.message : '发送失败，请重试'}`,
              isStreaming: false,
            };
          }
          return newMessages;
        });
      } finally {
        setIsLoading(false);
        setIsStreaming(false);
        setAbortController(null);
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

  // 处理消息删除
  const handleDeleteMessage = useCallback(
    async (messageId: string) => {
      if (!currentChatUuid) {
        // 如果是新对话，只需要本地删除
        setCurrentMessages(prev => prev.filter(msg => msg.id !== messageId));
        message.success('消息已删除');
        return;
      }

      try {
        // 调用API删除消息
        await chatClient.deleteChatMessage(currentChatUuid, messageId);

        // 本地删除消息
        setCurrentMessages(prev => prev.filter(msg => msg.id !== messageId));
        message.success('消息已删除');
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
      const messageIndex = currentMessages.findIndex(
        msg => msg.id === messageId
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
      const aiMessage: Message = {
        id: '',
        role: 'assistant',
        content: '',
        isStreaming: true, // 标记为正在流式输入
        timestamp: Date.now(),
      };
      setCurrentMessages(prev => [...prev, aiMessage]);

      try {
        // 构建请求消息历史（不包含被删除的消息）
        const apiMessages: CommonMessage[] = newMessages.map(
          ConvertMessageToApiMessage
        );

        // 构建请求参数
        const completionRequest: ServerCompletionsRequest = {
          model: selectedModel,
          messages: apiMessages,
          nonStandard: true,
          chatUuid: currentChatUuid || '',
          temperature: 0.8, // 提高随机性以获得不同的回答
          maxTokens: 2000,
          stream: true,
        };

        // 调用流式 Completions 接口重新生成
        await chatClient.createChatCompletionStream(
          completionRequest,
          // onMessage: 接收到新内容时的回调
          (content: string, reasoningContent?: string) => {
            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  content: newMessages[lastIndex].content + (content || ''),
                  reasoningContent: reasoningContent
                    ? (newMessages[lastIndex].reasoningContent || '') +
                      reasoningContent
                    : newMessages[lastIndex].reasoningContent,
                };
              }
              return newMessages;
            });
          },
          // onError: 接收到错误时的回调
          (error: string) => {
            // 将错误信息作为AI消息显示在聊天界面中
            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  content: packageErr(`${error}`),
                  isStreaming: false,
                };
              }
              return newMessages;
            });
          },
          // onComplete: 流式响应完成时的回调
          (chatUuid?: string) => {
            // 标记消息流式输入完成
            setCurrentMessages(prev => {
              const newMessages = [...prev];
              const lastIndex = newMessages.length - 1;
              if (
                lastIndex >= 0 &&
                newMessages[lastIndex].role === 'assistant' &&
                newMessages[lastIndex].isStreaming
              ) {
                newMessages[lastIndex] = {
                  ...newMessages[lastIndex],
                  isStreaming: false,
                };
              }
              return newMessages;
            });

            // 重置流式状态
            setIsStreaming(false);
            setAbortController(null);

            // 重新生成时通常不需要更新chatUuid，因为已经是现有对话
          },
          controller
        );
      } catch (error) {
        console.error('Regenerate message error:', error);
        // 将错误信息作为AI消息显示在聊天界面中
        setCurrentMessages(prev => {
          const newMessages = [...prev];
          const lastIndex = newMessages.length - 1;
          if (lastIndex >= 0 && newMessages[lastIndex].role === 'assistant') {
            newMessages[lastIndex] = {
              ...newMessages[lastIndex],
              content: `错误: ${error instanceof Error ? error.message : '重新生成失败，请重试'}`,
              isStreaming: false,
            };
          }
          return newMessages;
        });
      } finally {
        setIsLoading(false);
        setIsStreaming(false);
        setAbortController(null);
      }
    },
    [currentMessages, selectedModel, currentChatUuid]
  );

  // 处理文件上传
  const handleFileUpload = useCallback((files: File[]) => {
    message.success(`已选择 ${files.length} 个文件`);
    // 这里实现文件上传逻辑
  }, []);

  // 处理图片上传
  const handleImageUpload = useCallback((files: File[]) => {
    message.success(`已选择 ${files.length} 张图片`);
    // 这里实现图片上传逻辑
  }, []);

  // 聊天内容区域
  const renderChatContent = () => (
    <Chat
      chatTitle={chatTitle}
      chatUuid={currentChatUuid || undefined}
      currentMessages={currentMessages}
      isLoading={isLoading || isLoadingMessages} // 包含消息加载状态
      selectedModel={selectedModel}
      availableModels={availableModels}
      isMobile={isMobile}
      isSidebarCollapsed={isSidebarCollapsed}
      isGenerating={isStreaming}
      onTitleChange={handleTitleChange}
      onSendMessage={handleSendMessage}
      onStopGeneration={handleStopGeneration}
      onModelChange={setSelectedModel}
      onToggleSidebar={handleToggleSidebar}
      onCopyMessage={handleCopyMessage}
      onDeleteMessage={handleDeleteMessage}
      onRegenerateMessage={handleRegenerateMessage}
      onFileUpload={handleFileUpload}
      onImageUpload={handleImageUpload}
    />
  );

  return (
    <div
      className={`homePage ${className || ''} ${isMobile ? 'mobile-viewport-height' : ''}`}
      key={forceRerender} // Safari兼容性：强制重新渲染
    >
      <Layout className="layout">
        {/* 侧边栏 */}
        <Sider width="auto">
          <Index
            className="sidebar"
            currentChatUuid={currentChatUuid}
            onChatSelect={handleChatSelect}
            isSidebarCollapsed={isSidebarCollapsed}
            onToggleSidebar={handleToggleSidebar}
            onNewChat={handleNewChat}
            onRegisterRefreshCallback={handleRegisterRefreshCallback}
            onRegisterUpdateTitleCallback={handleRegisterUpdateTitleCallback}
          />
        </Sider>

        {/* 移动端遮罩 */}
        {isMobile && !isSidebarCollapsed && (
          <div
            className="sidebarOverlay"
            onClick={() => setIsSidebarCollapsed(true)}
          />
        )}

        {/* 主内容区域 */}
        <Content
          className="mainContent"
          style={{
            marginLeft: 0,
            transition: isMobile
              ? 'none'
              : 'margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
          }}
        >
          {renderChatContent()}
        </Content>
      </Layout>

      {/* 回到顶部 */}
      <BackTop className="backTop" />
    </div>
  );
};

export default ChatPage;
