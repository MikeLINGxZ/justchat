import {schema} from "../../wailsjs/go/models.ts";
import {Completions} from "../../wailsjs/go/service/Service";
import {EventsOn} from "../../wailsjs/runtime";

export async function CompletionsUtils(
    chatUuid: string,
    selectedModel: string,
    userMessage: schema.Message,
    onMessage: (message: schema.Message) => void,
    onError: (error: string) => void,
    onComplete: (chatUuid: string) => void,
    abortController?: AbortController
): Promise<void> {
    let cancel: (() => void) | null = null;
    let isCompleted = false;
    let hasReceivedFirstResponse = false;

    try {
        // 设置取消监听器
        if (abortController) {
            abortController.signal.addEventListener('abort', () => {
                if (cancel) {
                    cancel();
                    cancel = null;
                }
                if (!isCompleted) {
                    onError('请求已被取消');
                }
            });
        }

        // 调用 Completions API
        const emitKey = await Completions(chatUuid, selectedModel, userMessage);

        // 设置事件监听器
        cancel = EventsOn(emitKey, (responseMessage?: schema.Message) => {
            try {
                // 第一次接收到内容时标记
                if (!hasReceivedFirstResponse && responseMessage) {
                    hasReceivedFirstResponse = true;
                }

                // 处理接收到的消息
                if (responseMessage) {
                    console.log("[CompletionsUtils] responseMessage:",responseMessage)
                    onMessage(responseMessage);
                }

                // 检查是否完成
                if (responseMessage?.response_meta?.finish_reason && 
                    responseMessage.response_meta.finish_reason !== "") {
                    isCompleted = true;
                    if (cancel) {
                        cancel();
                        cancel = null;
                    }
                    onComplete(emitKey);
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
        
        // 清理资源
        if (cancel) {
            cancel();
        }
        
        throw error;
    }
}