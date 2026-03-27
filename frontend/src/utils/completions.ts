import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import {
    Completions,
    type Message
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
import { Events } from '@wailsio/runtime';
import * as view_models$0
    from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/models.ts";

export async function CompletionsUtils(
    messageInput: view_models$0.Message,
    onMessage: (message: Message) => void,
    onError: (error: string) => void,
    onComplete: (chatUuid: string) => void,
    abortController?: AbortController,
    onStreamStarted?: (resp: Completions | null) => void,
    onAborted?: (chatUuid: string) => void,
): Promise<void> {
    let cancel: (() => void) | null = null;
    let isCompleted = false;
    let messageKey: string = "" ;
    let streamChatUuid = "";
    try {
        // 中止时仅通知后端停止；保留监听直至收到带 finish_reason 的最终 Emit（见后端 chat defer）
        if (abortController) {
            abortController.signal.addEventListener('abort', () => {
                if (messageKey) {
                    Service.StopCompletions(messageKey);
                }
            });
        }

        // 调用 Completions API
        const resp: Completions | null = await Service.Completions(messageInput);
        messageKey = resp?.event_key ?? ""
        streamChatUuid = resp?.chat_uuid ?? ""

        if (!messageKey) {
            onStreamStarted?.(resp);
            if (abortController?.signal.aborted && streamChatUuid) {
                onAborted?.(streamChatUuid);
            }
            return;
        }

        onStreamStarted?.(resp)
        // 设置事件监听器
        cancel = Events.On(messageKey, (event) => {
            const responseMessage: Message = event.data;
            try {
                // 处理接收到的消息
                if (responseMessage) {
                    console.log("[CompletionsUtils] responseMessage:", responseMessage)
                    onMessage(responseMessage);
                }

                // 检查是否完成
                if (responseMessage?.assistant_message_extra?.finish_reason != "") {
                    // 清理事件监听器
                    if (cancel) {
                        cancel();
                        cancel = null;
                    }
                    Events.Off(messageKey);
                    onComplete(streamChatUuid);
                    isCompleted = true;
                }
            } catch (error) {
                console.error('处理响应消息时出错:', error);
                onError(`处理响应消息时出错: ${error instanceof Error ? error.message : String(error)}`);
            }
        });

        if (abortController?.signal.aborted) {
            Service.StopCompletions(messageKey);
        }

        // 等待完成（用户停止后也需等到最终带 finish_reason 的事件）
        return new Promise<void>((resolve) => {
            const checkCompletion = () => {
                if (isCompleted) {
                    resolve();
                } else {
                    setTimeout(checkCompletion, 100);
                }
            };
            checkCompletion();
        });

    } catch (error) {
        console.error('Completions API 调用失败:', error);
        onError(`API 调用失败: ${error instanceof Error ? error.message : String(error)}`);
        if (cancel) {
            cancel();
            cancel = null;
        }
        isCompleted = true;
        throw error;
    }
}
