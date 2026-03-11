import React, {useEffect, useRef, useState} from "react";
import styles from "@/pages/home/chat/index.module.scss";
import MessageList, {type MessageListRef} from "@/components/chat/message_list";
import ChatTitle from "@/components/chat/title";
import ChatInput from "@/components/chat/input";
import {
    type Chat as ChatType, FileInfo,
    type Message,
    Model
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

    useEffect(() => {
        Service.GetModels(true,true)
            .then((models: Model[])=>{
                setAvailableModels(models);
            })
    }, []);

    useEffect(() => {
        const propUuid = chatUuid ?? "";

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
            setLoading(false);
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
                    tools: []
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
                    finish_error: ""
                },
                assistant_message_extra_content: "",
            }

            const newMessages = [...messages, userMessage, assistantMessage];
            setMessages(newMessages);

            const controller = new AbortController();
            abortControllerRef.current = controller;

            await CompletionsUtils(userMessage,(message:Message)=>{
                setMessages(prev => {
                    const updatedMessages = [...prev];
                    updatedMessages[updatedMessages.length - 1] = message;
                    return updatedMessages;
                });
            },(error:string)=>{

            },(newChatUuid:string)=>{
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
            {loading ? (
                <div className={styles.chatLoadingContainer}>
                    <div className={styles.loadingSpinner} />
                </div>
            ) : (
                <>
                    <ChatTitle title={chatUuid == "" ? "新建对话" : (chatInfo?.title ?? "新建对话")} isSidebarCollapsed={isSidebarCollapsed} onToggleSidebar={onToggleSidebar}/>
                    <div className={`${styles.chatMessagesContent}`}>
                        <MessageList
                            ref={messageListRef}
                            messages={messages}
                            isGenerating={isGenerating}
                        />
                    </div>
                    <div className={`${styles.chatInput}`}>
                        <ChatInput
                            key={inputResetKey}
                            selectedModelId={selectModel}
                            availableModels={availableModels}
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
                </>
            )}
        </div>
    );
};

export default Chat;
export type { ChatProps };