import React, {useState} from "react";
import styles from "./index.module.scss";
import ReasoningContent from "@/components/chat/reasoning_message";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {Prism as SyntaxHighlighter} from "react-syntax-highlighter";
import {tomorrow} from "react-syntax-highlighter/dist/esm/styles/prism";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import type {Message} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import type {ToolUse} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";

interface ChatMessageProps {
    // 消息
    message: Message
    // 是否正在等待AI响应（用于显示loading状态）
    isLoading?: boolean
}

/** 工具调用区域（默认折叠，点击展开） */
const ToolUsesSection: React.FC<{ toolUses: ToolUse[] }> = ({ toolUses }) => {
    const [expanded, setExpanded] = useState(false);
    const count = toolUses.length;

    return (
        <div className={styles.toolUsesSection}>
            <div
                className={styles.toolUsesHeader}
                onClick={() => setExpanded(!expanded)}
                role="button"
            >
                <svg className={styles.toolUsesIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                </svg>
                <span>工具调用</span>
                <span className={styles.toolUsesCount}>({count})</span>
                <span className={styles.toolUsesChevron}>
                    {expanded ? '▼' : '▶'}
                </span>
            </div>
            {expanded && (
                <div className={styles.toolUsesList}>
                    {toolUses.map((toolUse, idx) => (
                        <ToolUseItem key={idx} toolUse={toolUse} />
                    ))}
                </div>
            )}
        </div>
    );
};

/** 单个工具调用的展示组件 */
const ToolUseItem: React.FC<{ toolUse: ToolUse }> = ({ toolUse }) => {
    const [expanded, setExpanded] = useState(false);
    const result = toolUse.tool_result?.trim() || '';
    const isLong = result.length > 120;
    const displayResult = isLong && !expanded ? result.slice(0, 120) + '…' : result;

    return (
        <div className={styles.toolUseItem}>
            <div
                className={`${styles.toolUseHeader} ${isLong ? styles.toolUseHeaderClickable : ''}`}
                onClick={() => isLong && setExpanded(!expanded)}
                role={isLong ? 'button' : undefined}
            >
                <span className={styles.toolUseBadge}>{toolUse.tool_name}</span>
                {isLong && (
                    <span className={styles.toolUseToggle}>
                        {expanded ? '收起' : '展开'}
                    </span>
                )}
            </div>
            {result && (
                <pre className={styles.toolUseResult}>{displayResult}</pre>
            )}
        </div>
    );
};

const ChatMessage: React.FC<ChatMessageProps> = ({
    message,
    isLoading = false,
}: ChatMessageProps) => {

    // 根据不同的角色选择不同的样式
    const isUser = message.role === 'user';
    let wrapperClass = styles.assistantMessageWrapper;
    if (isUser) {
        wrapperClass = styles.userMessageWrapper;
    }

    const isEmptyAssistant = !isUser &&
        !message.content?.trim() &&
        !message.reasoning_content?.trim() &&
        (message.assistant_message_extra?.tool_uses?.length ?? 0) === 0 &&
        (message.assistant_message_extra?.finish_error == "");

    // 如果是AI消息且内容和思考过程都为空，且不在loading状态，则不渲染
    if (isEmptyAssistant && !isLoading) {
        return null;
    }

    // 获取要渲染的内容：如果 content 为空，则使用 user_input_multi_content 的第一个 text 字段
    const getDisplayContent = () => {
        if (message.content.trim()) {
            return message.content;
        }
        if (message.content=="" && message.assistant_message_extra?.finish_error != "") {
            return message.assistant_message_extra?.finish_error
        }
        return '';
    };

    // 处理文件点击事件
    const handleFileClick = (filePath: string) => {
        if (filePath) {
            Service.OpenFile(filePath).catch((err) => {
                console.error('打开文件失败:', err);
            });
        }
    };

    return (
        <div className={styles.ChatMessage}>
            <div className={`${styles.message} ${wrapperClass}`}>
                <div className={styles.messageContainer} >
                    {isUser ? (
                        <>
                            <div className={styles.messageContent}>
                                {getDisplayContent()}
                            </div>
                            {(message.user_message_extra?.files?.length ?? 0) > 0 && (
                                <div className={styles.fileList}>
                                    {message.user_message_extra!.files!.map((file, index) => (
                                        <div
                                            key={index}
                                            className={styles.fileItem}
                                            onClick={() => handleFileClick(file.path)}
                                            title={`点击打开: ${file.name}`}
                                        >
                                            <span className={styles.fileType}>{file.mine_type}</span>
                                            <span className={styles.fileName}>{file.name}</span>
                                            {file.mine_type && (
                                                <span className={styles.fileMimeType}>{file.mine_type}</span>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            )}

                        </>
                    ):(
                        <div>
                            {/* 等待AI响应时显示loading动画 */}
                            {isLoading && isEmptyAssistant && (
                                <div className={styles.loadingIndicator}>
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                </div>
                            )}

                            {/* 渲染思考过程（如果存在） */}
                            {message.reasoning_content && (
                                <ReasoningContent
                                    content={message.reasoning_content}
                                    isStreaming={message.content == ""}
                                />
                            )}

                            {/* 渲染主要内容 */}
                            <div className={`${styles.messageContent} ${styles.markdownContent}`}>
                                <ReactMarkdown
                                    remarkPlugins={[remarkGfm]}
                                    components={{
                                        code(props: any) {
                                            const { node, inline, className, children, ...rest } = props;
                                            const match = /language-(\w+)/.exec(className || '');
                                            const language = match ? match[1] : '';

                                            return !inline && language ? (
                                                <SyntaxHighlighter
                                                    style={tomorrow}
                                                    language={language}
                                                    PreTag="div"
                                                    customStyle={{
                                                        margin: '8px 0',
                                                        borderRadius: '8px',
                                                        fontSize: '14px'
                                                    } as any}
                                                    {...rest}
                                                >
                                                    {String(children).replace(/\n$/, '')}
                                                </SyntaxHighlighter>
                                            ) : (
                                                <code className={`${className} ${styles.inlineCode}`} {...rest}>
                                                    {children}
                                                </code>
                                            );
                                        },
                                        // 自定义表格样式
                                        table: ({children}) => (
                                            <div className={styles.tableWrapper}>
                                                <table className={styles.markdownTable}>{children}</table>
                                            </div>
                                        ),
                                        // 自定义链接样式
                                        a: ({children, href}) => (
                                            <a
                                                href={href}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className={styles.markdownLink}
                                            >
                                                {children}
                                            </a>
                                        ),
                                        // 自定义引用块样式
                                        blockquote: ({children}) => (
                                            <blockquote className={styles.markdownBlockquote}>
                                                {children}
                                            </blockquote>
                                        ),
                                        // 自定义列表样式
                                        ul: ({children}) => (
                                            <ul className={styles.markdownList}>{children}</ul>
                                        ),
                                        ol: ({children}) => (
                                            <ol className={styles.markdownList}>{children}</ol>
                                        ),
                                    }}
                                >
                                    {message.content}
                                </ReactMarkdown>
                            </div>

                            {/* 工具调用信息 */}
                            {(message.assistant_message_extra?.tool_uses?.length ?? 0) > 0 && (
                                <ToolUsesSection toolUses={message.assistant_message_extra!.tool_uses!} />
                            )}
                        </div>
                    )}
                    {message.assistant_message_extra?.finish_reason === 'error' && (
                        <div className={styles.finishReasonError}>⚠ 因错误终止</div>
                    )}
                    {message.assistant_message_extra?.finish_reason === 'user stop' && (
                        <div className={styles.finishReasonUserStop}>⚠ 用户终止生成</div>
                    )}
                </div>
                
            </div>

        </div>
    )
}

ChatMessage.displayName = 'ChatMessage';
export default ChatMessage;