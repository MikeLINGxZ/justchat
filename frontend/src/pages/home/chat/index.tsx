import React, {useEffect, useMemo, useRef, useState} from "react";
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import ChatInput from "@/components/chat/input";
import {
    type Chat as ChatType, FileInfo,
    type Message,
    Model,
    Tool
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import {RoleType} from "@bindings/github.com/cloudwego/eino/schema";
import {CompletionsUtils} from "@/utils/completions.ts";
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
    const last = list[list.length - 1];
    if (last.role !== RoleType.Assistant || last.message_uuid !== cached.message_uuid) return list;
    const next = [...list];
    next[next.length - 1] = cached;
    return next;
}

const PENDING_NEW_CHAT_ABORT_KEY = "__NEW__";

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
    // 从「新建对话」路由发起生成时，后端返回的真实 chat uuid（在 navigate 完成前用于匹配事件）
    const [pendingNewChatStreamUuid, setPendingNewChatStreamUuid] = useState<string | null>(null);
    const pendingNewChatStreamUuidRef = useRef<string | null>(null);
    // 已发起 Completions、尚未收到 onStreamStarted 的会话（多任务并发时不能用单个 boolean）
    const [chatsAwaitingStreamStart, setChatsAwaitingStreamStart] = useState<string[]>([]);
    const [newChatsAwaitingStreamCount, setNewChatsAwaitingStreamCount] = useState(0);
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
    // 每个会话对应的 AbortController（仅用户点击停止时使用）
    const activeCompletionsRef = useRef<Map<string, AbortController>>(new Map());
    // 流式尚未开始时，按会话（或 __NEW__）挂起的 AbortController，支持多会话并发
    const pendingAbortStacksRef = useRef<Record<string, AbortController[]>>({});
    // 切换会话时的拉取序号，避免快速切换或 Strict Mode 下二次 setState 造成「闪两次 loading」
    const chatFetchSeqRef = useRef(0);
    // 新建对话：是否已在「首次 AI 实质输出」或 onComplete 兜底时刷新过侧边栏列表
    const newChatListRefreshDoneRef = useRef(false);

    const registerPendingAbort = (chatKey: string, controller: AbortController) => {
        const m = pendingAbortStacksRef.current;
        if (!m[chatKey]) m[chatKey] = [];
        m[chatKey].push(controller);
    };

    const unregisterPendingAbort = (chatKey: string, controller: AbortController) => {
        const arr = pendingAbortStacksRef.current[chatKey];
        if (!arr) return;
        const i = arr.indexOf(controller);
        if (i >= 0) arr.splice(i, 1);
        if (arr.length === 0) delete pendingAbortStacksRef.current[chatKey];
    };

    const tryAbortLatestPendingForKey = (chatKey: string): boolean => {
        const arr = pendingAbortStacksRef.current[chatKey];
        if (!arr?.length) return false;
        const c = arr.pop()!;
        if (arr.length === 0) delete pendingAbortStacksRef.current[chatKey];
        c.abort();
        return true;
    };

    const propUuid = chatUuid ?? "";

    const isGenerating = useMemo(() => {
        if (propUuid !== "") {
            if (generatingUuids.includes(propUuid)) return true;
            if (chatsAwaitingStreamStart.includes(propUuid)) return true;
            return false;
        }
        if (newChatsAwaitingStreamCount > 0) return true;
        return (
            pendingNewChatStreamUuid !== null &&
            generatingUuids.includes(pendingNewChatStreamUuid)
        );
    }, [
        propUuid,
        generatingUuids,
        pendingNewChatStreamUuid,
        chatsAwaitingStreamStart,
        newChatsAwaitingStreamCount,
    ]);

    // 侧边栏：已开始流式 + 已发请求但尚未 onStreamStarted 的已有会话
    const sidebarGeneratingUuids = useMemo(() => {
        const set = new Set(generatingUuids);
        chatsAwaitingStreamStart.forEach(u => {
            if (u) set.add(u);
        });
        return [...set];
    }, [generatingUuids, chatsAwaitingStreamStart]);

    useEffect(() => {
        onGeneratingUuidsChange?.(sidebarGeneratingUuids);
    }, [sidebarGeneratingUuids, onGeneratingUuidsChange]);

    useEffect(() => {
        if (!onRegisterStopGenerationForChat) return;
        onRegisterStopGenerationForChat((targetUuid: string) => {
            const mapped = activeCompletionsRef.current.get(targetUuid);
            if (mapped) {
                mapped.abort();
                return;
            }
            if (tryAbortLatestPendingForKey(targetUuid)) return;
        });
    }, [onRegisterStopGenerationForChat]);

    useEffect(() => {
        Service.GetModels(true,true)
            .then((models: Model[])=>{
                setAvailableModels(models);
            })
        Service.GetTools()
            .then((tools: Tool[])=> {
                setAvailableTools(tools);
            })
    }, []);

    // 当可用工具加载后，过滤掉已不存在的工具 ID
    useEffect(() => {
        if (availableTools.length === 0) return;
        const validIds = new Set(availableTools.map(t => t.id));
        setSelectedToolIds(prev => {
            const filtered = prev.filter(id => validIds.has(id));
            return filtered.length === prev.length ? prev : filtered;
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

    useEffect(() => {
        const v = chatUuid ?? "";
        visibleChatUuidRef.current = v;
        if (v !== "" && pendingNewChatStreamUuidRef.current !== null && v !== pendingNewChatStreamUuidRef.current) {
            setPendingNewChatStreamUuid(null);
            pendingNewChatStreamUuidRef.current = null;
        }

        // 如果这次 chatUuid 变化是由内部 navigate 引起的（新对话获得 UUID），只同步 currentChatUuid，不重新拉取，避免打断流式
        if (skipNextFetchRef.current && skipNextFetchRef.current === v) {
            skipNextFetchRef.current = null;
            setCurrentChatUuid(v);
            return;
        }

        setCurrentChatUuid(v);

        if (!v) {
            // 空白「新建对话」页与后台正在写的会话解绑，否则 pendingNewChatStreamUuid 仍指向旧 uuid 时会误显示停止按钮
            setPendingNewChatStreamUuid(null);
            pendingNewChatStreamUuidRef.current = null;
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
        ]).finally(() => {
            setTimeout(() => {
                if (fetchSeq !== chatFetchSeqRef.current) return;
                setLoading(false);
            }, 300);
        });
    }, [chatUuid]);

    useEffect(() => {
        pendingNewChatStreamUuidRef.current = pendingNewChatStreamUuid;
    }, [pendingNewChatStreamUuid]);

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

    const finishCompletionForChat = (streamChatUuid: string) => {
        activeCompletionsRef.current.delete(streamChatUuid);
        setGeneratingUuids(prev => prev.filter(x => x !== streamChatUuid));
        delete streamingAssistantByChatRef.current[streamChatUuid];
        if (pendingNewChatStreamUuidRef.current === streamChatUuid) {
            setPendingNewChatStreamUuid(null);
            pendingNewChatStreamUuidRef.current = null;
        }
    };

    // onSendButtonClick 发送按钮点击
    const onSendButtonClick = async () => {
        let streamStartedForThisRequest = false;
        let pendingAbortKey = "";
        let controller: AbortController | null = null;
        let fromEmptyChat = false;
        let sourceChatUuidForAwaiting = "";
        try {
            fromEmptyChat = currentChatUuid === "";
            sourceChatUuidForAwaiting = currentChatUuid;
            pendingAbortKey = fromEmptyChat ? PENDING_NEW_CHAT_ABORT_KEY : currentChatUuid;
            if (fromEmptyChat) {
                newChatListRefreshDoneRef.current = false;
            }
            if (fromEmptyChat) {
                setNewChatsAwaitingStreamCount(n => n + 1);
            } else if (currentChatUuid) {
                setChatsAwaitingStreamStart(prev => [...prev, currentChatUuid]);
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

            const newMessages = [...messages, userMessage, assistantMessage];
            setMessages(newMessages);

            controller = new AbortController();
            const abortCtl = controller;
            registerPendingAbort(pendingAbortKey, abortCtl);

            await CompletionsUtils(userMessage,(message:Message)=>{
                const uid = message.chat_uuid;
                streamingAssistantByChatRef.current[uid] = message;
                const visible = visibleChatUuidRef.current;
                const matchesVisible =
                    uid === visible ||
                    (visible === "" && uid === pendingNewChatStreamUuidRef.current);
                if (
                    fromEmptyChat &&
                    !newChatListRefreshDoneRef.current &&
                    matchesVisible &&
                    assistantHasSubstantiveOutput(message)
                ) {
                    newChatListRefreshDoneRef.current = true;
                    refreshChatList?.();
                }
                if (!matchesVisible) return;
                setMessages(prev => {
                    const updatedMessages = [...prev];
                    updatedMessages[updatedMessages.length - 1] = message;
                    return updatedMessages;
                });
            },(_error:string)=>{

            },(newChatUuid:string)=>{
                finishCompletionForChat(newChatUuid);
                // 已有会话：完成时不要改路由/消息树，否则用户切到别页会被强制跳回生成中的会话
                if (!fromEmptyChat) {
                    return;
                }
                const visible = visibleChatUuidRef.current;
                if (visible !== "" && visible !== newChatUuid) {
                    return;
                }
                if (visible === newChatUuid) {
                    return;
                }
                setCurrentChatUuid(newChatUuid);
                setMessages(prev => prev.map(msg => ({...msg, chat_uuid: newChatUuid})));
                skipNextFetchRef.current = newChatUuid;
                navigate(`/home/${newChatUuid}`, {replace: true});
                if (!newChatListRefreshDoneRef.current) {
                    newChatListRefreshDoneRef.current = true;
                    refreshChatList?.();
                }
            }, abortCtl, (resp) => {
                streamStartedForThisRequest = true;
                if (!resp?.chat_uuid) {
                    if (fromEmptyChat) {
                        setNewChatsAwaitingStreamCount(n => Math.max(0, n - 1));
                    } else if (sourceChatUuidForAwaiting) {
                        setChatsAwaitingStreamStart(prev => {
                            const i = prev.indexOf(sourceChatUuidForAwaiting);
                            if (i === -1) return prev;
                            const next = [...prev];
                            next.splice(i, 1);
                            return next;
                        });
                    }
                    unregisterPendingAbort(pendingAbortKey, abortCtl);
                    return;
                }
                const streamChatUuid = resp.chat_uuid;
                if (fromEmptyChat) {
                    setNewChatsAwaitingStreamCount(n => Math.max(0, n - 1));
                } else {
                    setChatsAwaitingStreamStart(prev => {
                        const i = prev.indexOf(streamChatUuid);
                        if (i === -1) return prev;
                        const next = [...prev];
                        next.splice(i, 1);
                        return next;
                    });
                }
                unregisterPendingAbort(pendingAbortKey, abortCtl);
                if (fromEmptyChat) {
                    // 尽早同步 URL 与父级 currentChatUuid，侧边栏能选中当前会话；skipNextFetch 避免打断流式去全量拉消息
                    skipNextFetchRef.current = streamChatUuid;
                    setPendingNewChatStreamUuid(streamChatUuid);
                    pendingNewChatStreamUuidRef.current = streamChatUuid;
                    setMessages(prev =>
                        prev.map(m => ({...m, chat_uuid: streamChatUuid})),
                    );
                    navigate(`/home/${streamChatUuid}`, {replace: true});
                }
                activeCompletionsRef.current.set(streamChatUuid, abortCtl);
                setGeneratingUuids(prev =>
                    prev.includes(streamChatUuid) ? prev : [...prev, streamChatUuid],
                );
            }, (abortedChatUuid) => {
                if (abortedChatUuid) {
                    finishCompletionForChat(abortedChatUuid);
                }
            })
        }catch (e) {
            if (!streamStartedForThisRequest) {
                if (fromEmptyChat) {
                    setNewChatsAwaitingStreamCount(n => Math.max(0, n - 1));
                } else if (sourceChatUuidForAwaiting) {
                    setChatsAwaitingStreamStart(prev => {
                        const i = prev.indexOf(sourceChatUuidForAwaiting);
                        if (i === -1) return prev;
                        const next = [...prev];
                        next.splice(i, 1);
                        return next;
                    });
                }
                if (controller && pendingAbortKey) {
                    unregisterPendingAbort(pendingAbortKey, controller);
                }
            }
        }
    }

    // onStopGeneration 停止生成点击
    const onStopGeneration = () => {
        const v = chatUuid ?? "";
        const targetUuid =
            v !== "" ? v : pendingNewChatStreamUuidRef.current ?? pendingNewChatStreamUuid;
        if (targetUuid) {
            const mapped = activeCompletionsRef.current.get(targetUuid);
            if (mapped) {
                mapped.abort();
                return;
            }
            if (tryAbortLatestPendingForKey(targetUuid)) return;
            return;
        }
        tryAbortLatestPendingForKey(PENDING_NEW_CHAT_ABORT_KEY);
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
