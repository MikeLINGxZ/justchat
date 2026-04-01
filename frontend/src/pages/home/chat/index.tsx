import React, {useCallback, useEffect, useMemo, useRef, useState} from "react";
import {Events} from '@wailsio/runtime';
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import ChatInput from "@/components/chat/input";
import {
    type Chat as ChatType,
    FileInfo,
    type Message,
    Model,
    Task,
    Tool
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import {
    RouteType,
    TaskStatus,
    ToolApprovalDecision,
    ToolApprovalResponse
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";
import {RoleType} from "@bindings/github.com/cloudwego/eino/schema";
import {
    BuildTaskFromCompletions,
    CompletionsUtils,
    SubscribeTaskStream,
    type TaskStreamEvent
} from "@/utils/completions.ts";
import {useNavigate} from "react-router-dom";
import {notify} from "@/utils/notification.ts";
import {useTranslation} from 'react-i18next';

interface ChatProps {
    // 对话uuid
    chatUuid?: string;
    // 是否折叠菜单栏
    isSidebarCollapsed: boolean
    // 点击菜单栏事件
    onToggleSidebar:  () => void
    //
    onChatChange: (chatUuid: string) => void
    // 刷新聊天列表
    refreshChatList: (() => void) | null;
    /** 同步「正在生成」的会话 uuid 列表给侧边栏等 */
    onGeneratingUuidsChange?: (uuids: string[]) => void;
    /** 注册按 chatUuid 停止生成（与输入框停止一致） */
    onRegisterStopGenerationForChat?: (fn: (chatUuid: string) => void) => void;
}

function assistantHasSubstantiveOutput(message: Message): boolean {
    if (message.role !== RoleType.Assistant) return false;
    const content = message.content?.trim() ?? "";
    const reasoning = message.reasoning_content?.trim() ?? "";
    const prefaceContent = message.assistant_message_extra?.preface_content?.trim() ?? "";
    const prefaceReasoning = message.assistant_message_extra?.preface_reasoning_content?.trim() ?? "";
    const toolUses = message.assistant_message_extra?.tool_uses?.length ?? 0;
    const traceSteps = message.assistant_message_extra?.execution_trace?.steps?.length ?? 0;
    return content.length > 0 || reasoning.length > 0 || prefaceContent.length > 0 || prefaceReasoning.length > 0 || toolUses > 0 || traceSteps > 0;
}

function isAssistantPlaceholderMessage(message: Message): boolean {
    if (message.role !== RoleType.Assistant) return false;
    if (assistantHasSubstantiveOutput(message)) return false;
    const finishReason = message.assistant_message_extra?.finish_reason?.trim() ?? "";
    const finishError = message.assistant_message_extra?.finish_error?.trim() ?? "";
    return finishReason === "" && finishError === "";
}

function toTaskStatus(status: string): TaskStatus {
    switch (status) {
        case TaskStatus.TaskStatusPending:
            return TaskStatus.TaskStatusPending;
        case TaskStatus.TaskStatusRunning:
            return TaskStatus.TaskStatusRunning;
        case TaskStatus.TaskStatusWaitingApproval:
            return TaskStatus.TaskStatusWaitingApproval;
        case TaskStatus.TaskStatusCompleted:
            return TaskStatus.TaskStatusCompleted;
        case TaskStatus.TaskStatusFailed:
            return TaskStatus.TaskStatusFailed;
        case TaskStatus.TaskStatusStopped:
            return TaskStatus.TaskStatusStopped;
        default:
            return TaskStatus.$zero;
    }
}

function isTerminalTaskStatus(status: string | TaskStatus | undefined | null): boolean {
    return status === TaskStatus.TaskStatusCompleted ||
        status === TaskStatus.TaskStatusFailed ||
        status === TaskStatus.TaskStatusStopped;
}

function getErrorMessage(error: unknown, fallback: string): string {
    if (error instanceof Error && error.message.trim()) {
        return error.message;
    }
    return fallback;
}

function mergeAssistantMessage(current: Message, incoming: Message): Message {
    return {
        ...current,
        ...incoming,
        assistant_message_extra: incoming.assistant_message_extra
            ? {
                ...(current.assistant_message_extra || {}),
                ...incoming.assistant_message_extra,
                execution_trace: incoming.assistant_message_extra.execution_trace || current.assistant_message_extra?.execution_trace,
            }
            : current.assistant_message_extra,
    };
}

function findMessageIndexForUpsert(list: Message[], incoming: Message): number {
    if (incoming.message_uuid) {
        const exactIndex = list.findIndex(item => item.message_uuid === incoming.message_uuid);
        if (exactIndex !== -1) {
            return exactIndex;
        }
    }

    if (incoming.role !== RoleType.Assistant) {
        return -1;
    }

    for (let index = list.length - 1; index >= 0; index -= 1) {
        const message = list[index];
        if (!isAssistantPlaceholderMessage(message)) {
            continue;
        }
        const incomingChatUuid = incoming.chat_uuid ?? "";
        const currentChatUuid = message.chat_uuid ?? "";
        if (!incomingChatUuid || !currentChatUuid || incomingChatUuid === currentChatUuid) {
            return index;
        }
    }

    return -1;
}

function upsertMessage(list: Message[], incoming: Message): Message[] {
    const index = findMessageIndexForUpsert(list, incoming);
    if (index === -1) {
        return [...list, incoming];
    }
    const next = [...list];
    next[index] = mergeAssistantMessage(next[index], incoming);
    return next;
}

function mergeStreamingAssistant(propUuid: string, list: Message[], cache: Record<string, Message>): Message[] {
    const cached = cache[propUuid];
    if (!cached) return list;
    return upsertMessage(list, cached);
}


const Chat: React.FC<ChatProps> = ({
    chatUuid,
    isSidebarCollapsed,
    onToggleSidebar,
    onChatChange,
    refreshChatList,
    onGeneratingUuidsChange,
    onRegisterStopGenerationForChat,
}) => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(!!chatUuid);
    const [chatInfo, setChatInfo] = useState<ChatType | null>(null);
    // 当前对话uuid
    const [currentChatUuid,setCurrentChatUuid] = useState<string>("");
    const currentChatUuidRef = useRef<string>("");
    // 所选择的模型
    const [selectModel,setSelectModelId] = useState<number>(-1);
    const [selectModelName,setSelectModelName] = useState<string>("");
    // 可用模型
    const [availableModels,setAvailableModels] = useState<Model[]>([]);
    // 可用工具
    const [availableTools,setAvailableTools] = useState<Tool[]>([]);
    // 用户选中的工具 id 列表（持久化到 localStorage）
    const [selectedToolIds, setSelectedToolIds] = useState<string[]>(() => {
        try {
            const raw = localStorage.getItem('chat_selected_tools');
            if (!raw) return [];
            const parsed = JSON.parse(raw) as string[];
            return Array.isArray(parsed) ? parsed : [];
        } catch {
            return [];
        }
    });
    // 后端已开始流式推送的会话 uuid（支持多会话后台生成）
    const [generatingUuids, setGeneratingUuids] = useState<string[]>([]);
    const [activeTasksByChat, setActiveTasksByChat] = useState<Record<string, Task>>({});
    const activeTasksByChatRef = useRef<Record<string, Task>>({});
    const [pendingExistingChatUuids, setPendingExistingChatUuids] = useState<string[]>([]);
    const [pendingNewChatCount, setPendingNewChatCount] = useState(0);
    // 输入框文字内容
    const [inputMessage,setInputMessage] = useState<string>("");
    // 输入框选择文件
    const [inputFiles,setInputFiles] = useState<FileInfo[]>([]);
    const [displayTitle, setDisplayTitle] = useState<string>("");
    // 当前聊天消息
    const [messages,setMessages] = useState<Message[]>([]);
    // 消息列表引用
    const messageListRef = useRef<MessageListRef>(null);
    // 标记：当新对话首次获得 UUID 时，跳过因 prop 变化触发的数据重新加载
    const skipNextFetchRef = useRef<string | null>(null);
    // 用于重置 ChatInput 内部状态（切换对话/新建对话时递增）
    const [inputResetKey, setInputResetKey] = useState<number>(0);
    // 子Agent面板是否打开

    // 当前路由可见会话，用于将流式事件路由到正确会话
    const visibleChatUuidRef = useRef<string>("");
    // 非当前可见会话的进行中 assistant 快照（切回时与 DB 合并）
    const streamingAssistantByChatRef = useRef<Record<string, Message>>({});
    // 活跃任务订阅
    const taskSubscriptionsRef = useRef<Map<string, () => void>>(new Map());
    // 切换会话时的拉取序号，避免快速切换或 Strict Mode 下二次 setState 造成「闪两次 loading」
    const chatFetchSeqRef = useRef(0);

    const refreshAvailableTools = useCallback(async () => {
        const tools = await Service.GetTools();
        setAvailableTools(tools);
        return tools;
    }, []);

    const propUuid = chatUuid ?? "";
    const activeTitleChatUuid = currentChatUuid || propUuid;

    const isGenerating = useMemo(() => {
        if (propUuid !== "") {
            if (generatingUuids.includes(propUuid)) return true;
            if (pendingExistingChatUuids.includes(propUuid)) return true;
            return false;
        }
        return pendingNewChatCount > 0;
    }, [propUuid, generatingUuids, pendingExistingChatUuids, pendingNewChatCount]);

    // 侧边栏：已开始流式 + 已发请求但尚未 onStreamStarted 的已有会话
    const sidebarGeneratingUuids = useMemo(() => {
        const set = new Set(generatingUuids);
        pendingExistingChatUuids.forEach(u => {
            if (u) set.add(u);
        });
        return [...set];
    }, [generatingUuids, pendingExistingChatUuids]);

    useEffect(() => {
        onGeneratingUuidsChange?.(sidebarGeneratingUuids);
    }, [sidebarGeneratingUuids, onGeneratingUuidsChange]);

    useEffect(() => {
        if (!onRegisterStopGenerationForChat) return;
        onRegisterStopGenerationForChat((targetUuid: string) => {
            const task = activeTasksByChatRef.current[targetUuid];
            if (!task?.task_uuid) {
                return;
            }
            Service.StopTask(task.task_uuid);
        });
    }, [onRegisterStopGenerationForChat]);

    useEffect(() => {
        activeTasksByChatRef.current = activeTasksByChat;
    }, [activeTasksByChat]);

    useEffect(() => {
        currentChatUuidRef.current = currentChatUuid;
    }, [currentChatUuid]);

    useEffect(() => {
        Service.GetModels(true,true)
            .then((models: Model[])=>{
                setAvailableModels(models);
            })
        refreshAvailableTools().catch(() => {
        });
    }, [refreshAvailableTools]);

    // 当可用工具加载后，过滤掉无效内置工具，并自动纳入启用中的自定义 MCP 工具
    useEffect(() => {
        if (availableTools.length === 0) return;
        const validBuiltinIds = new Set(
            availableTools
                .filter((tool) => tool.source_type === 'builtin')
                .map((tool) => tool.id)
        );
        const enabledCustomIds = availableTools
            .filter((tool) => tool.source_type === 'mcp_custom' && tool.enabled)
            .map((tool) => tool.id);
        setSelectedToolIds(prev => {
            const next = [...new Set([
                ...prev.filter(id => validBuiltinIds.has(id)),
                ...enabledCustomIds,
            ])];
            return next.length === prev.length && next.every((id, index) => id === prev[index]) ? prev : next;
        });
    }, [availableTools]);

    // 持久化用户选择的 tools
    useEffect(() => {
        if (selectedToolIds.length === 0) {
            localStorage.removeItem('chat_selected_tools');
        } else {
            localStorage.setItem('chat_selected_tools', JSON.stringify(selectedToolIds));
        }
    }, [selectedToolIds]);

    const shouldApplyEventToMessages = useCallback((list: Message[], event: TaskStreamEvent): boolean => {
        if (visibleChatUuidRef.current === event.chat_uuid) {
            return true;
        }
        if (currentChatUuidRef.current === event.chat_uuid) {
            return true;
        }
        if (list.some(message => (message.chat_uuid ?? "") === event.chat_uuid)) {
            return true;
        }
        if (visibleChatUuidRef.current !== "" || currentChatUuidRef.current !== "") {
            return false;
        }
        return list.some(message => isAssistantPlaceholderMessage(message) && (message.chat_uuid ?? "") === "");
    }, []);

    const syncStreamingAssistantToMessages = useCallback((chatUuidToSync: string) => {
        const cached = streamingAssistantByChatRef.current[chatUuidToSync];
        if (!cached) {
            return;
        }
        setMessages(prev => upsertMessage(prev, cached));
    }, []);

    const finishTaskTracking = (event: TaskStreamEvent) => {
        const cancel = taskSubscriptionsRef.current.get(event.task_uuid);
        cancel?.();
        taskSubscriptionsRef.current.delete(event.task_uuid);
        setGeneratingUuids(prev => prev.filter(x => x !== event.chat_uuid));
        setActiveTasksByChat(prev => {
            const next = {...prev};
            delete next[event.chat_uuid];
            return next;
        });
        if (visibleChatUuidRef.current === event.chat_uuid || currentChatUuidRef.current === event.chat_uuid) {
            delete streamingAssistantByChatRef.current[event.chat_uuid];
        }
        refreshChatList?.();
    };

    const finishTaskTrackingByTask = useCallback((task: Task) => {
        const cancel = taskSubscriptionsRef.current.get(task.task_uuid);
        cancel?.();
        taskSubscriptionsRef.current.delete(task.task_uuid);
        setGeneratingUuids(prev => prev.filter(x => x !== task.chat_uuid));
        setPendingExistingChatUuids(prev => prev.filter(uuid => uuid !== task.chat_uuid));
        setActiveTasksByChat(prev => {
            const next = {...prev};
            delete next[task.chat_uuid];
            return next;
        });
        if (visibleChatUuidRef.current === task.chat_uuid || currentChatUuidRef.current === task.chat_uuid) {
            delete streamingAssistantByChatRef.current[task.chat_uuid];
        }
        refreshChatList?.();
    }, [refreshChatList]);

    const handleTaskEvent = (event: TaskStreamEvent) => {
        const assistantMessage = event.assistant_message;
        if (assistantMessage?.message_uuid) {
            streamingAssistantByChatRef.current[event.chat_uuid] = assistantMessage;
        }

        setActiveTasksByChat(prev => ({
            ...prev,
            [event.chat_uuid]: new Task({
                ...(prev[event.chat_uuid] || {}),
                task_uuid: event.task_uuid,
                chat_uuid: event.chat_uuid,
                assistant_message_uuid: assistantMessage?.message_uuid ?? prev[event.chat_uuid]?.assistant_message_uuid ?? "",
                event_key: event.event_key,
                status: toTaskStatus(event.status),
                finish_reason: event.finish_reason,
                finish_error: event.finish_error,
            }),
        }));

        if (event.status === "pending" || event.status === "running") {
            setGeneratingUuids(prev => prev.includes(event.chat_uuid) ? prev : [...prev, event.chat_uuid]);
        } else {
            setMessages(prev => {
                if (!shouldApplyEventToMessages(prev, event)) {
                    return prev;
                }
                return upsertMessage(prev, assistantMessage);
            });
            return;
        }

        if (visibleChatUuidRef.current !== event.chat_uuid) {
            return;
        }
        setMessages(prev => upsertMessage(prev, assistantMessage));
    };

    const ensureTaskSubscription = (task: Task | null | undefined) => {
        if (!task?.task_uuid || !task?.event_key) {
            return;
        }
        setActiveTasksByChat(prev => ({...prev, [task.chat_uuid]: task}));
        if (task.status === "pending" || task.status === "running") {
            setGeneratingUuids(prev => prev.includes(task.chat_uuid) ? prev : [...prev, task.chat_uuid]);
        }
        if (taskSubscriptionsRef.current.has(task.task_uuid)) {
            return;
        }
        const cancel = SubscribeTaskStream(
            task,
            handleTaskEvent,
            (error) => {
                console.error(error);
            },
            (event) => {
                finishTaskTracking(event);
            },
        );
        if (cancel) {
            taskSubscriptionsRef.current.set(task.task_uuid, cancel);
        }
    };

    const reconcileTaskAfterSubscription = useCallback(async (task: Task | null | undefined) => {
        if (!task?.task_uuid) {
            return;
        }
        try {
            const latestTask = await Service.GetTask(task.task_uuid);
            if (!latestTask || !isTerminalTaskStatus(latestTask.status)) {
                return;
            }

            console.warn("missed terminal event, reconciled from GetTask", {
                task_uuid: latestTask.task_uuid,
                chat_uuid: latestTask.chat_uuid,
                event_key: latestTask.event_key,
                status: latestTask.status,
                finish_reason: latestTask.finish_reason,
                finish_error: latestTask.finish_error,
            });

            finishTaskTrackingByTask(latestTask);

            if (visibleChatUuidRef.current !== latestTask.chat_uuid && currentChatUuidRef.current !== latestTask.chat_uuid) {
                return;
            }

            const messageList = await Service.ChatMessages(latestTask.chat_uuid, 0, 200);
            const raw = messageList?.messages ?? [];
            const merged = mergeStreamingAssistant(latestTask.chat_uuid, raw, streamingAssistantByChatRef.current);
            setMessages(merged);
            delete streamingAssistantByChatRef.current[latestTask.chat_uuid];
        } catch (error) {
            console.error("Failed to reconcile task after subscription:", error);
        }
    }, [finishTaskTrackingByTask]);

    useEffect(() => {
        Service.GetRunningTasks()
            .then((taskList) => {
                taskList?.tasks?.forEach(task => ensureTaskSubscription(task));
            })
            .catch((err) => {
                console.error("Failed to restore running tasks:", err);
            });

        return () => {
            taskSubscriptionsRef.current.forEach(cancel => cancel());
            taskSubscriptionsRef.current.clear();
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        const v = chatUuid ?? "";
        visibleChatUuidRef.current = v;

        // 如果这次 chatUuid 变化是由内部 navigate 引起的（新对话获得 UUID），只同步 currentChatUuid，不重新拉取，避免打断流式
        if (skipNextFetchRef.current && skipNextFetchRef.current === v) {
            skipNextFetchRef.current = null;
            setCurrentChatUuid(v);
            syncStreamingAssistantToMessages(v);
            return;
        }

        setCurrentChatUuid(v);

        if (!v) {
            setChatInfo(null);
            setDisplayTitle("");
            setLoading(false);
            setMessages([]);
            setInputMessage("");
            setInputFiles([]);
            setInputResetKey(prev => prev + 1);
            return;
        }
        const fetchSeq = ++chatFetchSeqRef.current;
        setLoading(true);
        setInputResetKey(prev => prev + 1);
        Promise.all([
            Service.ChatInfo(v)
                .then((info: ChatType | null) => {
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    setChatInfo(info);
                    setDisplayTitle(info?.title?.trim() || "");
                })
                .catch((err) => {
                    console.error("Failed to fetch chat info:", err);
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    setChatInfo(null);
                    setDisplayTitle("");
                }),
            Service.ChatMessages(v, 0, 200)
                .then((messageList) => {
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    const raw = messageList!.messages;
                    const merged = mergeStreamingAssistant(v, raw, streamingAssistantByChatRef.current);
                    setMessages(merged);
                    delete streamingAssistantByChatRef.current[v];
                })
                .catch((err) => {
                    console.error("Failed to fetch chat messages info:", err);
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    setMessages([]);
                }),
            Service.GetChatActiveTask(v)
                .then((task) => {
                    if (fetchSeq !== chatFetchSeqRef.current || !task) return;
                    ensureTaskSubscription(task);
                })
                .catch((err) => {
                    console.error("Failed to fetch active task:", err);
                }),
        ]).finally(() => {
            setTimeout(() => {
                if (fetchSeq !== chatFetchSeqRef.current) return;
                setLoading(false);
            }, 300);
        });
    }, [chatUuid]);

    // onModelSelectorClick 模型选择框点击事件
    const onModelSelectorClick = () => {
        Service.GetModels(true,true)
            .then((models: Model[])=>{
                setAvailableModels(models);
            })
    }

    // onSelectModelChange 所选模型变更
    const onSelectModelChange  = (modelId: number, modelName: string) => {
        setSelectModelId(modelId);
        setSelectModelName(modelName);
    }

    // onMessageChange 输入消息变更
    const onMessageChange = (message: string) => {
        setInputMessage(message);
    }

    // onSelectFileChange 输入文件变更
    const onSelectFileChange = (paths: FileInfo[]) => {
        setInputFiles(paths);
    }

    // onSendButtonClick 发送按钮点击
    const onSendButtonClick = async () => {
        if (selectModel <= 0 || !selectModelName.trim()) {
            return;
        }
        let fromEmptyChat = false;
        const prevMessages = messages;
        const pendingTitle = inputMessage.trim();
        try {
            fromEmptyChat = currentChatUuid === "";
            if (fromEmptyChat) {
                setPendingNewChatCount(n => n + 1);
            } else if (currentChatUuid) {
                setPendingExistingChatUuids(prev => [...prev, currentChatUuid]);
            }
            let userMessage:Message = {
                id: 0,
                created_at: null,
                updated_at: null,
                deleted_at: null,
                role: RoleType.User,
                chat_uuid: currentChatUuid,
                message_uuid: "",
                content: inputMessage,
                reasoning_content: "",
                user_message_extra: {
                    model_id: selectModel,
                    model_name: selectModelName,
                    files: inputFiles,
                    tools: selectedToolIds,
                    agents: [],
                },
                user_message_extra_content: "",
                assistant_message_extra: null,
                assistant_message_extra_content: "",
            }
            let assistantMessage:Message = {
                id: 0,
                created_at: null,
                updated_at: null,
                deleted_at: null,
                role: RoleType.Assistant,
                chat_uuid: currentChatUuid,
                content: "",
                reasoning_content: "",
                message_uuid: "",
                user_message_extra: null,
                user_message_extra_content: "",
                assistant_message_extra: {
                    execution_trace: {
                        steps: [],
                    },
                    route_type: RouteType.$zero,
                    retry_count: 0,
                    current_stage: "",
                    current_agent: "",
                    preface_content: "",
                    preface_reasoning_content: "",
                    finish_reason: "",
                    finish_error: "",
                    tool_uses: [],
                    pending_approvals: [],
                },
                assistant_message_extra_content: "",
            }

            setMessages(prev => [...prev, userMessage, assistantMessage]);

            const resp = await CompletionsUtils(userMessage);
            if (!resp?.chat_uuid || !resp?.task_uuid) {
                setMessages(prevMessages);
                if (fromEmptyChat) {
                    setPendingNewChatCount(n => Math.max(0, n - 1));
                } else if (currentChatUuid) {
                    setPendingExistingChatUuids(prev => prev.filter(uuid => uuid !== currentChatUuid));
                }
                return;
            }

            const streamChatUuid = resp.chat_uuid;
            assistantMessage.chat_uuid = streamChatUuid;
            assistantMessage.message_uuid = resp.message_uuid;

            if (fromEmptyChat) {
                setPendingNewChatCount(n => Math.max(0, n - 1));
                skipNextFetchRef.current = streamChatUuid;
                setCurrentChatUuid(streamChatUuid);
                if (pendingTitle) {
                    setDisplayTitle(pendingTitle);
                }
                setMessages(prev => prev.map(m => ({...m, chat_uuid: streamChatUuid})));
                navigate(`/home/${streamChatUuid}`, {replace: true});
            } else {
                setPendingExistingChatUuids(prev => prev.filter(uuid => uuid !== streamChatUuid));
            }

            setGeneratingUuids(prev =>
                prev.includes(streamChatUuid) ? prev : [...prev, streamChatUuid],
            );
            const task = BuildTaskFromCompletions(resp, assistantMessage);
            ensureTaskSubscription(task);
            console.log("chat task created:", {
                task_uuid: task.task_uuid,
                event_key: task.event_key,
                chat_uuid: task.chat_uuid,
            });
            void reconcileTaskAfterSubscription(task);
            refreshChatList?.();
        }catch (e) {
            setMessages(prevMessages);
            if (fromEmptyChat) {
                setPendingNewChatCount(n => Math.max(0, n - 1));
            } else if (currentChatUuid) {
                setPendingExistingChatUuids(prev => prev.filter(uuid => uuid !== currentChatUuid));
            }
        }
    }

    // onStopGeneration 停止生成点击
    const onStopGeneration = () => {
        const v = chatUuid ?? "";
        if (!v) {
            return;
        }
        const task = activeTasksByChatRef.current[v];
        if (!task?.task_uuid) {
            return;
        }
        Service.StopTask(task.task_uuid);
    }

    const handleTitleChange = useCallback(async (newTitle: string) => {
        const activeChatUuid = currentChatUuid || propUuid;
        if (!activeChatUuid) {
            return;
        }
        await Service.RenameChat(activeChatUuid, newTitle);
        setChatInfo(prev => (prev ? {...prev, title: newTitle} : prev));
        setDisplayTitle(newTitle);
    }, [currentChatUuid, propUuid]);

    const handleApprovalDecision = useCallback(async (approvalId: string, decision: 'allow' | 'reject') => {
        try {
            await Service.RespondToolApproval(new ToolApprovalResponse({
                approval_id: approvalId,
                decision: decision === 'allow'
                    ? ToolApprovalDecision.ToolApprovalDecisionAllow
                    : ToolApprovalDecision.ToolApprovalDecisionReject,
                comment: "",
            }));
        } catch (error: any) {
            notify.error(t('home.chat.approvalFailed'), getErrorMessage(error, t('home.chat.approvalFailedDesc')));
        }
    }, [t]);

    const handleSendApprovalComment = useCallback(async (approvalId: string, comment: string) => {
        try {
            await Service.RespondToolApproval(new ToolApprovalResponse({
                approval_id: approvalId,
                decision: ToolApprovalDecision.ToolApprovalDecisionCustom,
                comment,
            }));
        } catch (error: any) {
            notify.error(t('home.chat.approvalCommentFailed'), getErrorMessage(error, t('home.chat.approvalCommentFailedDesc')));
        }
    }, [t]);

    useEffect(() => {
        if (!chatInfo?.title?.trim()) {
            return;
        }
        setDisplayTitle(chatInfo.title);
    }, [chatInfo?.title]);

    useEffect(() => {
        if (!activeTitleChatUuid) {
            return;
        }
        const eventKey = `event:chat_title:${activeTitleChatUuid}`;
        const cancel = Events.On(eventKey, (event) => {
            const payload = event.data as { chat_uuid?: string; title?: string };
            const nextTitle = payload?.title;
            if (payload?.chat_uuid !== activeTitleChatUuid || !nextTitle) {
                return;
            }
            setDisplayTitle(nextTitle);
            setChatInfo(prev => (prev ? {...prev, title: nextTitle} : prev));
        });

        return () => {
            cancel?.();
            Events.Off(eventKey);
        };
    }, [activeTitleChatUuid]);

    return (
        <div className={`${styles.chatPage}`}>
            {/* 主内容始终渲染 */}
            <ChatTitle
                title={displayTitle}
                uuid={propUuid}
                onTitleChange={handleTitleChange}
                isSidebarCollapsed={isSidebarCollapsed}
                onToggleSidebar={onToggleSidebar}
            />
            <div className={`${styles.chatMessagesContent}`}>
                <MessageList
                    key={currentChatUuid || 'new'}
                    ref={messageListRef}
                    messages={messages}
                    isGenerating={isGenerating}
                    useInstantScrollOnFirstLoad
                    onApprovalDecision={handleApprovalDecision}
                    onSendApprovalComment={handleSendApprovalComment}
                />
            </div>
            <div className={`${styles.chatInput}`}>
                <ChatInput
                    key={inputResetKey}
                    selectedModelId={selectModel}
                    hasSelectedModel={selectModel > 0 && !!selectModelName.trim()}
                    availableModels={availableModels}
                    availableTools={availableTools}
                    selectedToolIds={selectedToolIds}
                    onSelectedToolsChange={setSelectedToolIds}
                    onRefreshTools={refreshAvailableTools}
                    isGenerating={isGenerating}
                    onMessageChange={onMessageChange}
                    onSendButtonClick={onSendButtonClick}
                    onSelectModelChange={onSelectModelChange}
                    onSelectFileChange={onSelectFileChange}
                    onStopGeneration={onStopGeneration}
                    onModelSelectorClick={onModelSelectorClick}
                    onMessageListScrollToBottom={() => {
                        messageListRef.current?.scrollToBottom();
                    }}
                />
            </div>
            {/* loading 蒙层覆盖在主内容之上 */}
            {loading && (
                <div className={styles.chatLoadingContainer}>
                    <div className={styles.loadingSpinner} />
                </div>
            )}
        </div>
    );
};

export default Chat;
export type { ChatProps };
