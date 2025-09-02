import { Api } from './service/chat/Api';
import { getApiConfig } from '@/config/env';
import { createGlobalErrorHandlerFetch } from '@/utils/globalErrorHandler';
import type {
  ServerCompletionsRequest,
  ServerCompletionsResponse,
  ServerListChatsRequest,
  ServerListChatsResponse,
  ServerGetChatMessagesResponse,
  ServerDeleteChatMessageResponse,
  ServerChatTitleResponse,
  ChatGetChatMessagesBody,
  ChatDeleteChatBody,
  ChatDeleteChatMessageBody,
  ChatChatTitleBody,
  ChatChatTitleSaveBody,
  CommonEmpty,
} from './service/chat/Api';

// 创建聊天API客户端实例
class ChatClient {
  private api: Api<any>;

  constructor() {
    const apiConfig = getApiConfig();

    this.api = new Api({
      baseUrl: apiConfig.baseUrl,
      baseApiParams: {
        credentials: apiConfig.withCredentials ? 'include' : 'same-origin',
        headers: {
          'Content-Type': 'application/json',
        },
      },
      securityWorker: _securityData => {
        const token = localStorage.getItem('token');
        if (token) {
          return {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          };
        }
        return {};
      },
      customFetch: createGlobalErrorHandlerFetch(fetch),
    });
  }

  // 创建聊天完成（流式）
  async createChatCompletion(
    data: ServerCompletionsRequest
  ): Promise<ServerCompletionsResponse> {
    const response = await this.api.v1.chatCompletions(data, { secure: true });
    return response.data.result || {};
  }

  // 获取聊天列表
  async listChats(
    params?: ServerListChatsRequest
  ): Promise<ServerListChatsResponse> {
    const response = await this.api.v1.chatListChats(params || {}, {
      secure: true,
    });
    return response.data;
  }

  // 删除聊天
  async deleteChat(chatUuid: string): Promise<CommonEmpty> {
    const response = await this.api.v1.chatDeleteChat(
      chatUuid,
      {},
      { secure: true }
    );
    return response.data;
  }

  // 获取聊天消息
  async getChatMessages(
    chatUuid: string,
    params?: ChatGetChatMessagesBody
  ): Promise<ServerGetChatMessagesResponse> {
    const response = await this.api.v1.chatGetChatMessages(
      chatUuid,
      params || {},
      { secure: true }
    );
    return response.data;
  }

  // 删除聊天消息
  async deleteChatMessage(
    chatUuid: string,
    messageId: string
  ): Promise<ServerDeleteChatMessageResponse> {
    const response = await this.api.v1.chatDeleteChatMessage(
      chatUuid,
      messageId,
      {},
      { secure: true }
    );
    return response.data;
  }

  // 获取对话标题
  async getChatTitle(chatUuid: string): Promise<ServerChatTitleResponse> {
    const response = await this.api.v1.chatChatTitle(
      chatUuid,
      {},
      { secure: true }
    );
    return response.data;
  }

  // 保存对话标题
  async saveChatTitle(
    chatUuid: string,
    chatTitle: string
  ): Promise<CommonEmpty> {
    const response = await this.api.v1.chatChatTitleSave(
      chatUuid,
      encodeURIComponent(chatTitle),
      {},
      { secure: true }
    );
    return response.data;
  }

  // 流式聊天完成
  async createChatCompletionStream(
    data: ServerCompletionsRequest,
    onMessage: (content: string, reasoningContent?: string) => void,
    onError: (error: string) => void,
    onComplete: (chatUuid?: string) => void,
    abortController?: AbortController
  ): Promise<void> {
    const apiConfig = getApiConfig();
    const token = localStorage.getItem('token');
    let chatUuid = '';
    // 设置流式请求参数
    const streamData = { ...data, stream: true };

    try {
      const response = await fetch(`${apiConfig.baseUrl}/v1/chat/completions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: token ? `Bearer ${token}` : '',
          Accept: 'text/event-stream',
        },
        credentials: apiConfig.withCredentials ? 'include' : 'same-origin',
        body: JSON.stringify(streamData),
        signal: abortController?.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error('无法获取响应流');
      }

      const decoder = new TextDecoder();
      let buffer = '';
      chatUuid = '';

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n');
          buffer = lines.pop() || '';

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const data = line.slice(6).trim();
              if (data === '[DONE]') {
                return;
              }

              try {
                const parsed = JSON.parse(data);

                // 添加调试日志
                console.log('SSE数据:', data);
                console.log('解析后的数据:', parsed);

                // 检查是否有错误
                if (parsed.error) {
                  onError(parsed.error);
                  return;
                }

                // 提取chatUuid
                if (parsed.chat_uuid && !chatUuid) {
                  chatUuid = parsed.chat_uuid;
                }

                // 提取消息内容
                if (parsed.choices && parsed.choices.length > 0) {
                  const choice = parsed.choices[0];
                  console.log('choice对象:', choice);

                  if (choice.delta) {
                    const content = choice.delta.content || '';

                    // 尝试多种可能的reasoning_content字段位置
                    let reasoningContent = '';
                    if (choice.delta.reasoning_content) {
                      reasoningContent = choice.delta.reasoning_content;
                    } else if (choice.delta.reasoningContent) {
                      reasoningContent = choice.delta.reasoningContent;
                    } else if (choice.delta.reasoning) {
                      reasoningContent = choice.delta.reasoning;
                    }

                    console.log('提取的内容:', { content, reasoningContent });
                    console.log(
                      'delta完整对象:',
                      JSON.stringify(choice.delta, null, 2)
                    );

                    // 只有当content或reasoningContent有值时才调用回调
                    if (content || reasoningContent) {
                      onMessage(content, reasoningContent);
                    }
                  }
                }
              } catch (parseError) {
                console.warn('解析SSE数据失败:', parseError, 'data:', data);
              }
            }
          }
        }
      } finally {
        reader.releaseLock();
      }
    } catch (error) {
      console.error('流式请求失败:', error);
      // 如果是AbortError，不调用onError，因为这是用户主动停止
      if (error instanceof Error && error.name === 'AbortError') {
        console.log('请求被用户中断');
      } else {
        onError(error instanceof Error ? error.message : '请求失败');
      }
    } finally {
      onComplete(chatUuid);
    }
  }
}

// 创建单例实例
const chatClient = new ChatClient();

export { chatClient };
export type {
  ServerCompletionsRequest,
  ServerCompletionsResponse,
  ServerListChatsRequest,
  ServerListChatsResponse,
  ServerGetChatMessagesResponse,
  ServerDeleteChatMessageResponse,
  ServerChatTitleResponse,
  ChatGetChatMessagesBody,
  ChatDeleteChatBody,
  ChatDeleteChatMessageBody,
  ChatChatTitleBody,
  ChatChatTitleSaveBody,
  CommonEmpty,
};
