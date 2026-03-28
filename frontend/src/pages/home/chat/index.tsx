import React, {useCallback, useEffect, useMemo, useRef, useState} from "react";
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import ChatInput from "@/components/chat/input";
import {
    type Chat as ChatType, FileInfo,
    type Message,
    Model,
    Task,
    Tool
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import {RoleType} from "@bindings/github.com/cloudwego/eino/schema";
import {BuildTaskFromCompletions, CompletionsUtils, SubscribeTaskStream, type TaskStreamEvent} from "@/utils/completions.ts";
import {useNavigate} from "react-router-dom";

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

function mergeStreamingAssistant(propUuid: string, list: Message[], cache: Record<string, Message>): Message[] {
    const cached = cache[propUuid];
    if (!cached || list.length === 0) return list;
    const index = list.findIndex(message => message.message_uuid === cached.message_uuid);
    if (index === -1) return [...list, cached];
    const next = [...list];
    next[index] = cached;
    return next;
}

function assistantHasSubstantiveOutput(message: Message): boolean {
    if (message.role !== RoleType.Assistant) return false;
    const content = message.content?.trim() ?? "";
    const reasoning = message.reasoning_content?.trim() ?? "";
    const toolUses = message.assistant_message_extra?.tool_uses?.length ?? 0;
    return content.length > 0 || reasoning.length > 0 || toolUses > 0;
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
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(!!chatUuid);
    const [chatInfo, setChatInfo] = useState<ChatType | null>(null);
    // 当前对话uuid
    const [currentChatUuid,setCurrentChatUuid] = useState<string>("");
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
    // 当前聊天消息
    const [messages,setMessages] = useState<Message[]>([]);
    // 消息列表引用
    const messageListRef = useRef<MessageListRef>(null);
    // 标记：当新对话首次获得 UUID 时，跳过因 prop 变化触发的数据重新加载
    const skipNextFetchRef = useRef<string | null>(null);
    // 用于重置 ChatInput 内部状态（切换对话/新建对话时递增）
    const [inputResetKey, setInputResetKey] = useState<number>(0);
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

    const upsertMessage = (list: Message[], incoming: Message): Message[] => {
        const index = list.findIndex(item => item.message_uuid === incoming.message_uuid);
        if (index === -1) {
            return [...list, incoming];
        }
        const next = [...list];
        next[index] = incoming;
        return next;
    };

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
        delete streamingAssistantByChatRef.current[event.chat_uuid];
        refreshChatList?.();
    };

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
                status: event.status,
                finish_reason: event.finish_reason,
                finish_error: event.finish_error,
            }),
        }));

        if (event.status === "pending" || event.status === "running") {
            setGeneratingUuids(prev => prev.includes(event.chat_uuid) ? prev : [...prev, event.chat_uuid]);
        } else {
            setMessages(prev => {
                if (visibleChatUuidRef.current !== event.chat_uuid) {
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
            return;
        }

        setCurrentChatUuid(v);

        if (!v) {
            setChatInfo(null);
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
                })
                .catch((err) => {
                    console.error("Failed to fetch chat info:", err);
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    setChatInfo(null);
                }),
            Service.ChatMessages(v, 0, 200)
                .then((messageList) => {
                    if (fetchSeq !== chatFetchSeqRef.current) return;
                    const raw = messageList!.messages;
                    const merged = mergeStreamingAssistant(v, raw, streamingAssistantByChatRef.current);
                    setMessages(merged);
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
        let fromEmptyChat = false;
        const prevMessages = messages;
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
                    finish_reason: "",
                    finish_error: "",
                    tool_uses: [],
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

    return (
        <div className={`${styles.chatPage}`}>
            {/* 主内容始终渲染 */}
            <ChatTitle title={chatUuid == "" ? "新建对话" : (chatInfo?.title ?? "新建对话")} isSidebarCollapsed={isSidebarCollapsed} onToggleSidebar={onToggleSidebar}/>
            <div className={`${styles.chatMessagesContent}`}>
                <MessageList
                    key={currentChatUuid || 'new'}
                    ref={messageListRef}
                    messages={messages}
                    isGenerating={isGenerating}
                    useInstantScrollOnFirstLoad
                />
            </div>
            <div className={`${styles.chatInput}`}>
                <ChatInput
                    key={inputResetKey}
                    selectedModelId={selectModel}
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
