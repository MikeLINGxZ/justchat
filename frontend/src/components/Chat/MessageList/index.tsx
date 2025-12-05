import React, {
    forwardRef,
    useImperativeHandle,
    useRef,
    useCallback,
    useEffect,
    useState,
    useMemo,
    type ReactNode
} from "react";
import {Message} from "@bindings/github.com/cloudwego/eino/schema/index.ts";
import styles from "./index.module.scss";
import ChatMessage from "@/components/Chat/Message";

interface MessageListProps {
    // 类名
    className?: string;
    // 所有消息
    messages?: Message[]; // 修改为 Message[]
    // 是否加载中
    isLoading?: boolean;
    // 消息是否在生成中
    isGenerating?: boolean;
}


const MessageList:  React.FC<MessageListProps> = ({
    className,
    messages = [],
    isLoading,
    isGenerating
}) => {



    return (
        <div className={`${className} ${styles.MessageList}`}>
            {/* 消息 */}
            <div className={styles.content}>
                {
                    messages.map((message, index) => (
                        <div>
                            <ChatMessage message={message}/>
                        </div>
                    ))
                }
                {/* 滚动到底部按钮 */}
                <div className={`${styles.scrollButton}`}>
                    <div className={`${styles.scrollButtonIcon}`}>
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M7 10L12 15L17 10"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </svg>
                    </div>
                </div>
            </div>
        </div>
    )
};

MessageList.displayName = 'MessageList';

export default MessageList;