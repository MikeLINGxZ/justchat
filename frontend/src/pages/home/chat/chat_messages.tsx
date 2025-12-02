import React, {forwardRef, useImperativeHandle, useRef, useCallback, useEffect, useState, useMemo} from "react";
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism';
import remarkGfm from 'remark-gfm';
import MessageAction from "@/components/MessageAction";
import ReasoningContent from "@/components/ReasoningContent";
import styles from "@/pages/home/chat/chat_messages.module.scss";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";

interface ChatMessagesProps {
    // 所有消息
    messages?: Message[]; // 修改为 Message[]
    // 是否加载中
    isLoading?: boolean;
    // 是否显示loading消息
    showLoadingMessage?: boolean;
    // 是否为移动端
    isMobile?: boolean;
    // 消息是否自动滚动到底部
    autoScrollBottom?: boolean;
    // 复制消息事件
    onCopyMessage?: (content: string) => void;
    // 删除消息事件
    onDeleteMessage?: (messageId: string) => void;
    // 重新生成消息事件
    onRegenerateMessage?: (messageId: string) => void;
    // 用户滚动事件
    onUserScroll?: (isUserScrolling: boolean) => void;
    // 类名
    className?: string;
}

export interface ChatMessagesRef {
    scrollToBottom: () => void;
    scrollToBottomSmooth: () => void;
    isAtBottom: () => boolean;
    getScrollContainer: () => HTMLDivElement | null;
    enableAutoScroll: () => void;
}

const ChatMessages:  React.ForwardRefRenderFunction<ChatMessagesRef,ChatMessagesProps> = ({
    messages = [],
    isLoading,
    showLoadingMessage,
    isMobile = false,
    autoScrollBottom,
    onCopyMessage,
    onDeleteMessage,
    onRegenerateMessage,
    onUserScroll,
    className
},ref) => {

    const chatMessagesPageRef = useRef<HTMLDivElement>(null);
    const [isUserScrolling, setIsUserScrolling] = useState(false);
    const scrollTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const lastScrollTopRef = useRef<number>(0); // 记录上次滚动位置
    const userScrollDetectionRef = useRef<NodeJS.Timeout | null>(null);
    const scrollStartTimeRef = useRef<number>(0); // 记录滚动开始时间
    const isScrollingByUserRef = useRef<boolean>(false); // 用户正在滚动的标记
    const userInteractionTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    
    // 检查是否滚动到底部
    const isAtBottom = useCallback(() => {
        if (!chatMessagesPageRef.current) return false;
        const { scrollTop, scrollHeight, clientHeight } = chatMessagesPageRef.current;
        // 允许20px的误差范围，更宽松的判断
        return scrollHeight - scrollTop - clientHeight <= 20;
    }, []);

    // 立即滚动到底部（无动画）- 添加节流优化
    const scrollToBottomInstant = useCallback(() => {
        if (chatMessagesPageRef.current) {
            const element = chatMessagesPageRef.current;
            const targetScrollTop = element.scrollHeight;
            
            // 避免不必要的滚动操作
            if (Math.abs(element.scrollTop - targetScrollTop) > 1) {
                element.scrollTop = targetScrollTop;
                lastScrollTopRef.current = targetScrollTop;
            }
        }
    }, []);

    // 平滑滚动到底部 - 添加节流优化
    const scrollToBottomSmooth = useCallback(() => {
        if (chatMessagesPageRef.current) {
            const element = chatMessagesPageRef.current;
            const targetScrollTop = element.scrollHeight;
            
            // 避免不必要的滚动操作
            if (Math.abs(element.scrollTop - targetScrollTop) > 1) {
                element.scrollTo({
                    top: targetScrollTop,
                    behavior: 'smooth'
                });
                lastScrollTopRef.current = targetScrollTop;
            }
        }
    }, []);

    // 启用自动滚动
    const enableAutoScroll = useCallback(() => {
        // 清除所有定时器
        if (userInteractionTimeoutRef.current) {
            clearTimeout(userInteractionTimeoutRef.current);
        }
        if (userScrollDetectionRef.current) {
            clearTimeout(userScrollDetectionRef.current);
        }
        
        // 重置所有滚动状态
        isScrollingByUserRef.current = false;
        setIsUserScrolling(false);
        onUserScroll?.(false);
        
        // 滚动到底部
        scrollToBottomInstant();
    }, [scrollToBottomInstant, onUserScroll]);

    // 处理用户滚动事件
    const handleScroll = useCallback((e: Event) => {
        if (!chatMessagesPageRef.current) return;
        
        const element = chatMessagesPageRef.current;
        const currentScrollTop = element.scrollTop;
        const scrollHeight = element.scrollHeight;
        const clientHeight = element.clientHeight;
        
        // 检查是否在底部
        const atBottom = scrollHeight - currentScrollTop - clientHeight <= 20;
        
        // 检测任何滚动位置变化（零容忍检测）
        const scrollDiff = Math.abs(currentScrollTop - lastScrollTopRef.current);
        
        // 如果滚动位置有任何变化且用户已被标记为正在滚动，立即触发状态
        if (scrollDiff > 0 && isScrollingByUserRef.current && !isUserScrolling) {
            // 立即停止自动滚动
            setIsUserScrolling(true);
            onUserScroll?.(true);
        }
        
        // 特性：当用户手动滚动到底部时，自动重新启用自动滚动
        if (atBottom && isUserScrolling) {
            // 清除所有定时器
            if (userInteractionTimeoutRef.current) {
                clearTimeout(userInteractionTimeoutRef.current);
                userInteractionTimeoutRef.current = null;
            }
            if (userScrollDetectionRef.current) {
                clearTimeout(userScrollDetectionRef.current);
                userScrollDetectionRef.current = null;
            }
            
            // 重置所有滚动状态，重新启用自动滚动
            isScrollingByUserRef.current = false;
            setIsUserScrolling(false);
            onUserScroll?.(false);
            
            // 更新记录的滚动位置
            lastScrollTopRef.current = currentScrollTop;
            return; // 早期返回，避免后续逻辑干扰
        }
        
        // 清除之前的定时器
        if (userScrollDetectionRef.current) {
            clearTimeout(userScrollDetectionRef.current);
        }
        
        // 只有在用户滚动状态下且不在底部时才设置结束检测
        if ((isUserScrolling || isScrollingByUserRef.current) && !atBottom) {
            userScrollDetectionRef.current = setTimeout(() => {
                // 检查是否回到底部
                if (isAtBottom()) {
                    // 用户停止滚动后自动回到底部，重新启用自动滚动
                    isScrollingByUserRef.current = false;
                    setIsUserScrolling(false);
                    onUserScroll?.(false);
                }
            }, 150);
        }
        
        // 更新记录的滚动位置
        lastScrollTopRef.current = currentScrollTop;
    }, [isUserScrolling, onUserScroll, isAtBottom]);

    // 处理用户主动滚动开始（立即响应）
    const handleUserScrollStart = useCallback(() => {
        // 立即标记用户正在滚动
        isScrollingByUserRef.current = true;
        
        // 清除之前的定时器
        if (userInteractionTimeoutRef.current) {
            clearTimeout(userInteractionTimeoutRef.current);
        }
        if (userScrollDetectionRef.current) {
            clearTimeout(userScrollDetectionRef.current);
        }
        
        // 立即触发用户滚动状态，无条件停止自动滚动
        setIsUserScrolling(true);
        onUserScroll?.(true);
        
        // 设置一个延迟来重置用户滚动标记
        userInteractionTimeoutRef.current = setTimeout(() => {
            isScrollingByUserRef.current = false;
        }, 200); // 缩短延迟时间
    }, [onUserScroll]);

    useImperativeHandle(ref, () => ({
        scrollToBottom: scrollToBottomInstant,
        scrollToBottomSmooth,
        isAtBottom,
        getScrollContainer: () => chatMessagesPageRef.current,
        enableAutoScroll
    }));

    // 添加滚动事件监听
    useEffect(() => {
        const element = chatMessagesPageRef.current;
        if (element) {
            // 添加滚动事件监听
            element.addEventListener('scroll', handleScroll, { passive: true });
            
            // 添加用户输入检测事件（立即响应）
            element.addEventListener('wheel', handleUserScrollStart, { passive: true });
            element.addEventListener('touchstart', handleUserScrollStart, { passive: true });
            element.addEventListener('touchmove', handleUserScrollStart, { passive: true });
            
            // 添加键盘事件监听（方向键滚动）
            element.addEventListener('keydown', (e) => {
                if (['ArrowUp', 'ArrowDown', 'PageUp', 'PageDown', 'Home', 'End'].includes(e.key)) {
                    handleUserScrollStart();
                }
            }, { passive: true });
            
            return () => {
                element.removeEventListener('scroll', handleScroll);
                element.removeEventListener('wheel', handleUserScrollStart);
                element.removeEventListener('touchstart', handleUserScrollStart);
                element.removeEventListener('touchmove', handleUserScrollStart);
                element.removeEventListener('keydown', handleUserScrollStart as any);
            };
        }
    }, [handleScroll, handleUserScrollStart]);

    // 优化的自动滚动逻辑 - 减少不必要的触发
    const lastMessageRef = useRef<Message | null>(null); // 修改为 Message
    const lastMessageCountRef = useRef<number>(0);
    
    useEffect(() => {
        if (!autoScrollBottom || isUserScrolling || !chatMessagesPageRef.current) return;
        
        const currentMessageCount = messages.length;
        const lastMessage = messages[messages.length - 1];
        
        // 只有在以下情况才触发滚动：
        // 1. 消息数量发生变化（新消息）
        // 2. 最后一条消息内容发生变化（流式更新）
        // 3. 加载状态变化
        const shouldScroll = 
            currentMessageCount !== lastMessageCountRef.current ||
            (lastMessage && (
                !lastMessageRef.current ||
                lastMessage.content !== lastMessageRef.current.content ||
                lastMessage.reasoning_content !== lastMessageRef.current.reasoning_content ||
                (lastMessage as any).isStreaming !== (lastMessageRef.current as any).isStreaming // Message 类没有 isStreaming 属性
            )) ||
            isLoading;
            
        if (shouldScroll) {
            // 使用 requestAnimationFrame 优化滚动性能
            requestAnimationFrame(() => {
                if (!chatMessagesPageRef.current) return;
                
                // 对于流式消息，使用立即滚动保证实时跟进
                const hasStreamingMessage = messages.some(msg => (msg as any).isStreaming); // Message 类没有 isStreaming 属性
                if (hasStreamingMessage || isLoading) {
                    scrollToBottomInstant();
                } else {
                    scrollToBottomSmooth();
                }
            });
        }
        
        // 更新引用
        lastMessageCountRef.current = currentMessageCount;
        lastMessageRef.current = lastMessage ? { ...lastMessage } as any : null; // 修复类型问题
        
    }, [messages, isLoading, autoScrollBottom, isUserScrolling, scrollToBottomInstant, scrollToBottomSmooth]);

    // 清理定时器
    useEffect(() => {
        return () => {
            if (scrollTimeoutRef.current) {
                clearTimeout(scrollTimeoutRef.current);
            }
            if (userScrollDetectionRef.current) {
                clearTimeout(userScrollDetectionRef.current);
            }
            if (userInteractionTimeoutRef.current) {
                clearTimeout(userInteractionTimeoutRef.current);
            }
        };
    }, []);

    // 初始化滚动位置
    useEffect(() => {
        if (chatMessagesPageRef.current) {
            lastScrollTopRef.current = chatMessagesPageRef.current.scrollTop || 0;
        }
    }, []);

    // 处理复制消息
    const handleCopyMessage = (content: string) => {
        if (onCopyMessage) {
            onCopyMessage(content);
        } else {
            // 默认复制到剪贴板
            navigator.clipboard.writeText(content).catch(console.error);
        }
    };

    // 处理删除消息
    const handleDeleteMessage = (messageId: string) => {
        if (onDeleteMessage) {
            onDeleteMessage(messageId);
        }
    };

    // 处理重新生成消息
    const handleRegenerateMessage = (messageId: string) => {
        if (onRegenerateMessage) {
            onRegenerateMessage(messageId);
        }
    };

    // 渲染消息操作按钮
    const renderMessageActions = (message: Message, messageIndex: number) => { // 修改为 Message
        const isUser = message.role === 'user';
        
        // 判断是否为最后一条AI消息（只有最后一条AI消息才能重新生成）
        const isLastAssistantMessage = !isUser && messageIndex === messages.length - 1;
        
        return (
            <MessageAction
                message={message}
                onCopyMessage={handleCopyMessage}
                onDeleteMessage={handleDeleteMessage}
                onRegenerateMessage={handleRegenerateMessage}
                isLastAssistantMessage={isLastAssistantMessage}
                hideByDefault={true}
                alignRight={isUser}
                isMobile={isMobile}
            />
        );
    };

    // 渲染消息内容
    const renderMessageContent = useCallback((message: Message) => { // 修改为 Message
        const isUser = message.role === 'user';
        
        // 用户消息直接显示文本，AI消息使用Markdown渲染
        if (isUser) {
            return (
                <div className={styles.messageContent}>
                    {message.content}
                </div>
            );
        }
        
        // AI消息使用Markdown渲染
        return (
            <>
                {/* 渲染思考过程（如果存在） */}
                {message.reasoning_content && (
                    <ReasoningContent 
                        content={message.reasoning_content} 
                        isStreaming={(message as any).isStreaming || false} // Message 类没有 isStreaming 属性
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
            </>
        );
    }, []);

    // 检测是否为错误消息
    const isErrorMessage = useCallback((message: Message) => {
        // todo
        // return message.role === 'assistant' && message.content.includes('错误');
        return false;
    }, []);

    // 渲染单条消息
    const renderMessage = useCallback((message: Message, index: number) => { // 修改为 Message
        const isUser = message.role === 'user';
        const isErrorMsg = isErrorMessage(message);
        
        // 如果是AI消息且内容和思考过程都为空，则不渲染
        if (!isUser && !message.content.trim() && !message.reasoning_content?.trim()) {
            return null;
        }
        
        let wrapperClass = styles.assistantMessageWrapper;
        if (isUser) {
            wrapperClass = styles.userMessageWrapper;
        } else if (isErrorMsg) {
            wrapperClass = styles.errorMessageWrapper;
        }

        return (
            <div key={index} className={`${styles.messageWrapper} ${wrapperClass}`}> {/* 使用 index 作为 key */}
                <div className={styles.messageContainer}>
                    {renderMessageContent(message)}
                </div>
                {renderMessageActions(message, index)}
            </div>
        );
    }, [isErrorMessage, renderMessageContent, renderMessageActions]);

    // 优化消息列表渲染
    const renderedMessages = useMemo(() => {
        return messages.map((message, index) => renderMessage(message, index));
    }, [messages, renderMessage]);

    return (
        <div className={`${styles.chatMessagesPage} ${className || ''}`} ref={chatMessagesPageRef}>
            <div className={styles.messagesList}>
                {renderedMessages}
                
                {/* 加载状态 */}
                {(isLoading || showLoadingMessage) && (
                    <div className={styles.loadingIndicator}>
                        <div className={styles.loadingDots}>
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                        <span>AI 正在思考中...</span>
                    </div>
                )}
            </div>
        </div>
    );
};

ChatMessages.displayName = 'ChatMessages';

export default forwardRef(ChatMessages);