import React, {useEffect, useRef, useState} from "react";
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
import {useIsMobile} from "@/hooks/useViewportHeight.ts";

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

}

const Chat: React.FC<ChatProps> = ({
    chatUuid,
    isSidebarCollapsed,
    onToggleSidebar,
    onChatChange,
    refreshChatList,
}) => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState<boolean>(!!chatUuid);
    const [chatInfo, setChatInfo] = useState<ChatType | null>(null);
    const isMobile =  useIsMobile();
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
    // 是否在生成
    const [isGenerating,setIsGenerating] = useState<boolean>(false);
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
    // 用于中止生成请求
    const abortControllerRef = useRef<AbortController | null>(null);
    // 正在生成的聊天 uuid，用于切换时忽略过期的流式消息
    const generatingForChatUuidRef = useRef<string | null>(null);

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
        const propUuid = chatUuid ?? "";

        // 切换聊天时：中止进行中的生成，并忽略其后续的流式消息
        generatingForChatUuidRef.current = null;
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
            setIsGenerating(false);
        }

        // 如果这次 chatUuid 变化是由内部 navigate 引起的（新对话获得 UUID），跳过重新加载
        if (skipNextFetchRef.current && skipNextFetchRef.current === propUuid) {
            skipNextFetchRef.current = null;
            setCurrentChatUuid(propUuid);
            return;
        }

        setCurrentChatUuid(propUuid);

        if (!propUuid) {
            setChatInfo(null);
            setLoading(false);
            setMessages([]);
            setInputMessage("");
            setInputFiles([]);
            setIsGenerating(false);
            setInputResetKey(prev => prev + 1);
            return;
        }
        setLoading(true);
        setInputResetKey(prev => prev + 1);
        Promise.all([
            Service.ChatInfo(propUuid)
                .then((info: ChatType | null) => {
                    setChatInfo(info);
                })
                .catch((err) => {
                    console.error("Failed to fetch chat info:", err);
                    setChatInfo(null);
                }),
            Service.ChatMessages(propUuid, 0, 200)
                .then((messageList) => {
                    setMessages(messageList!.messages);
                })
                .catch((err) => {
                    console.error("Failed to fetch chat messages info:", err);
                    setMessages([]);
                }),
        ]).finally(() => {
            setTimeout(() => {
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
        try {
            generatingForChatUuidRef.current = currentChatUuid;
            setIsGenerating(true);
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

            const controller = new AbortController();
            abortControllerRef.current = controller;

            await CompletionsUtils(userMessage,(message:Message)=>{
                // 若已切换到其他聊天，忽略此流式消息，避免污染当前聊天
                if (generatingForChatUuidRef.current === null) return;
                setMessages(prev => {
                    const updatedMessages = [...prev];
                    updatedMessages[updatedMessages.length - 1] = message;
                    return updatedMessages;
                });
            },(error:string)=>{

            },(newChatUuid:string)=>{
                generatingForChatUuidRef.current = null;
                abortControllerRef.current = null;
                setIsGenerating(false);
                if (currentChatUuid != "" && currentChatUuid == newChatUuid) {
                    return
                }
                setCurrentChatUuid(newChatUuid);
                setMessages(prev => prev.map(msg => ({...msg, chat_uuid: newChatUuid})));
                // 标记跳过下一次因 prop 变化触发的数据加载，避免闪烁
                skipNextFetchRef.current = newChatUuid;
                navigate(`/home/${newChatUuid}`, {replace: true});
                if (refreshChatList) {
                    setTimeout(()=>{
                        refreshChatList();
                    },1000)
                }
            }, controller)
        }catch (e) {
            generatingForChatUuidRef.current = null;
            abortControllerRef.current = null;
            setIsGenerating(false);
        }
    }

    // onStopGeneration 停止生成点击
    const onStopGeneration = () => {
        abortControllerRef.current?.abort();
        abortControllerRef.current = null;
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