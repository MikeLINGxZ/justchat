import { Events } from '@wailsio/runtime';
import { Service } from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import {
    Completions,
    Task,
    type Message,
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import type { ExecutionTrace, TraceStep } from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";
import * as view_models$0 from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/models.ts";

export interface TaskStreamEvent {
    task_uuid: string;
    chat_uuid: string;
    event_key: string;
    status: string;
    finish_reason: string;
    finish_error: string;
    execution_trace?: ExecutionTrace;
    trace_delta?: TraceStep[];
    current_stage?: string;
    current_agent?: string;
    retry_count?: number;
    assistant_message: Message;
}

function isTaskFinished(event: TaskStreamEvent): boolean {
    return event.status !== "pending" && event.status !== "running" && event.status !== "waiting_approval";
}

export function SubscribeTaskStream(
    task: Task,
    onEvent: (event: TaskStreamEvent) => void,
    onError?: (error: string) => void,
    onComplete?: (event: TaskStreamEvent) => void,
): (() => void) | null {
    if (!task?.event_key || !task?.task_uuid) {
        return null;
    }

    let cancel: (() => void) | null = null;
    cancel = Events.On(task.event_key, (event) => {
        try {
            const payload = event.data as TaskStreamEvent;
            onEvent(payload);
            if (isTaskFinished(payload)) {
                onComplete?.(payload);
                cancel?.();
                cancel = null;
                Events.Off(task.event_key);
            }
        } catch (error) {
            onError?.(`处理任务流事件时出错: ${error instanceof Error ? error.message : String(error)}`);
        }
    });

    return () => {
        cancel?.();
        cancel = null;
        Events.Off(task.event_key);
    };
}

export async function CompletionsUtils(
    messageInput: view_models$0.Message,
): Promise<Completions | null> {
    return await Service.Completions(messageInput);
}

export function BuildTaskFromCompletions(resp: Completions, assistantMessage: Message): Task {
    return new Task({
        task_uuid: resp.task_uuid,
        chat_uuid: resp.chat_uuid,
        assistant_message_uuid: assistantMessage.message_uuid || resp.message_uuid,
        event_key: resp.event_key,
        status: "pending",
    });
}
