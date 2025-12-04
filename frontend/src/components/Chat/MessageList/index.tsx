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
            <div className={styles.Content}>
                {
                    messages.map((message, index) => (
                        <div>
                            <ChatMessage message={message}/>
                        </div>
                    ))
                }
            </div>
        </div>
    )
};

MessageList.displayName = 'MessageList';

export default MessageList;