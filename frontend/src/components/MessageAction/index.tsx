import React from "react";
import { Message } from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import styles from "./index.module.scss";

interface MessageActionProps {
    // 消息对象
    message: Message; // 修改为 schema.Message
    // 复制消息事件
    onCopyMessage?: (content: string) => void;
    // 删除消息事件
    onDeleteMessage?: (messageId: string) => void;
    // 重新生成消息事件
    onRegenerateMessage?: (messageId: string) => void;
    // 是否为最后一条AI消息（只有最后一条AI消息才能重新生成）
    isLastAssistantMessage?: boolean;
    // 是否默认隐藏（hover时显示）
    hideByDefault?: boolean;
    // 是否靠右对齐（用户消息）
    alignRight?: boolean;
    // 是否为移动端
    isMobile?: boolean;
    // 类名
    className?: string;
}

const MessageAction: React.FC<MessageActionProps> = ({
    message,
    onCopyMessage,
    onDeleteMessage,
    onRegenerateMessage,
    isLastAssistantMessage = false,
    hideByDefault = true,
    alignRight,
    isMobile = false,
    className
}) => {
    const isUser = message.role === 'user';

    // 处理复制消息
    const handleCopyMessage = () => {
        // 构建完整的消息内容，包含思考过程（如果存在）
        let fullContent = '';
        
        // 如果有思考过程，先添加思考过程
        if (message.reasoning_content && message.reasoning_content.trim()) {
            fullContent += `## 思考过程\n\n${message.reasoning_content}\n\n## 回答\n\n`;
        }
        
        // 添加主要内容
        fullContent += message.content;
        
        if (onCopyMessage) {
            onCopyMessage(fullContent);
        } else {
            // 默认复制到剪贴板
            navigator.clipboard.writeText(fullContent).catch(console.error);
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

    // 默认显示按钮，不再需要 hover 才显示
    const shouldHide = false; // 始终显示操作按钮
    
    return (
        <div 
            className={`${styles.messageActions} ${shouldHide ? styles.hideByDefault : ''} ${alignRight ? styles.alignRight : ''} ${className || ''}`}
            data-message-actions
        >
            {/* 复制按钮 */}
            <button
                className={styles.actionButton}
                onClick={handleCopyMessage}
                title="复制消息内容"
            >
                <svg 
                    width="14" 
                    height="14" 
                    viewBox="0 0 24 24" 
                    fill="none" 
                    stroke="currentColor" 
                    strokeWidth="2" 
                    strokeLinecap="round" 
                    strokeLinejoin="round"
                    className={styles.buttonIcon}
                >
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                </svg>
                <span className={styles.tooltipOverlay}></span>
            </button>

            {/* AI 消息且为最后一条时才显示重新生成按钮 */}
            {/*{!isUser && isLastAssistantMessage && (*/}
            {/*    <button*/}
            {/*        className={styles.actionButton}*/}
            {/*        onClick={() => handleRegenerateMessage('')} // Message 类没有 id 属性*/}
            {/*        title="重新生成回答"*/}
            {/*    >*/}
            {/*        <svg */}
            {/*            width="14" */}
            {/*            height="14" */}
            {/*            viewBox="0 0 24 24" */}
            {/*            fill="none" */}
            {/*            stroke="currentColor" */}
            {/*            strokeWidth="2" */}
            {/*            strokeLinecap="round" */}
            {/*            strokeLinejoin="round"*/}
            {/*            className={styles.buttonIcon}*/}
            {/*        >*/}
            {/*            <polyline points="23 4 23 10 17 10"></polyline>*/}
            {/*            <polyline points="1 20 1 14 7 14"></polyline>*/}
            {/*            <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15"></path>*/}
            {/*        </svg>*/}
            {/*        <span className={styles.tooltipOverlay}></span>*/}
            {/*    </button>*/}
            {/*)}*/}

            {/* 删除按钮 */}
            {/*<button*/}
            {/*    className={styles.actionButton}*/}
            {/*    onClick={() => handleDeleteMessage('')} // Message 类没有 id 属性*/}
            {/*    title="删除这条消息"*/}
            {/*>*/}
            {/*    <svg */}
            {/*        width="14" */}
            {/*        height="14" */}
            {/*        viewBox="0 0 24 24" */}
            {/*        fill="none" */}
            {/*        stroke="currentColor" */}
            {/*        strokeWidth="2" */}
            {/*        strokeLinecap="round" */}
            {/*        strokeLinejoin="round"*/}
            {/*        className={styles.buttonIcon}*/}
            {/*    >*/}
            {/*        <polyline points="3,6 5,6 21,6"></polyline>*/}
            {/*        <path d="M19,6V20a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6M8,6V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2V6"></path>*/}
            {/*    </svg>*/}
            {/*    <span className={styles.tooltipOverlay}></span>*/}
            {/*</button>*/}
        </div>
    );
};

export default MessageAction;