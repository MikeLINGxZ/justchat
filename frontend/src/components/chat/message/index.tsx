import React, {forwardRef} from "react";
import {Message} from "@bindings/github.com/cloudwego/eino/schema";
import styles from "./index.module.scss";
import ReasoningContent from "@/components/chat/reasoning_message";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {Prism as SyntaxHighlighter} from "react-syntax-highlighter";
import {tomorrow} from "react-syntax-highlighter/dist/esm/styles/prism";

interface ChatMessageProps {
    // 消息
    message: Message
}

const ChatMessage: React.FC<ChatMessageProps> = ({
    message,
}: ChatMessageProps) => {

    // 根据不同的角色选择不同的样式
    const isUser = message.role === 'user';
    let wrapperClass = styles.assistantMessageWrapper;
    if (isUser) {
        wrapperClass = styles.userMessageWrapper;
    }

    // 如果是AI消息且内容和思考过程都为空，则不渲染
    if (!isUser && !message.content.trim() && !message.reasoning_content?.trim()) {
        return null;
    }

    // todo
    //  wrapperClass = styles.errorMessageWrapper;

    return (
        <div className={styles.ChatMessage}>
            <div className={`${styles.message} ${wrapperClass}`}>
                <div className={styles.messageContainer} >
                    {isUser ? (
                        <div className={styles.messageContent}>
                            {message.content}
                        </div>
                    ):(
                        <div>
                            {/* 渲染思考过程（如果存在） */}
                            {message.reasoning_content && (
                                <ReasoningContent
                                    content={message.reasoning_content}
                                    isStreaming={(message as any).isStreaming || false} // message 类没有 isStreaming 属性
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
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

ChatMessage.displayName = 'ChatMessage';
export default ChatMessage;