import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import {
    Completions,
    type Message
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
import { Events } from '@wailsio/runtime';
import {GenEventsKey} from "@/utils/events.ts";
import * as view_models$0
    from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/models.ts";

export async function CompletionsUtils(
    messageInput: view_models$0.Message,
    onMessage: (message: Message) => void,
    onError: (error: string) => void,
    onComplete: (chatUuid: string) => void,
    abortController?: AbortController
): Promise<void> {
    let cancel: (() => void) | null = null;
    let isCompleted = false;
    let messageKey: string = "" ;
    try {
        // 设置取消监听器：中止时停止后端并移除事件监听，避免切换聊天后仍收到流式消息
        if (abortController) {
            abortController.signal.addEventListener('abort', () => {
                if (messageKey) {
                    Service.StopCompletions(messageKey);
                    if (cancel) {
                        cancel();
                        cancel = null;
                    }
                    Events.Off(messageKey);
                }
            });
        }

        // 调用 Completions API
        const resp: Completions | null = await Service.Completions(messageInput);
        messageKey = resp?.event_key ?? ""
        // 设置事件监听器
        cancel = Events.On(resp?.event_key!, (event) => {
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
                    Events.Off(resp?.event_key!);
                    onComplete(resp?.chat_uuid!);
                    isCompleted = true;
                }
            } catch (error) {
                console.error('处理响应消息时出错:', error);
                onError(`处理响应消息时出错: ${error instanceof Error ? error.message : String(error)}`);
            }
        });

        // 等待完成或取消
        return new Promise<void>((resolve, reject) => {
            const checkCompletion = () => {
                if (isCompleted) {
                    resolve();
                } else if (abortController?.signal.aborted) {
                    reject(new Error('请求已被取消'));
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