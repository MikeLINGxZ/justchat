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
import ChatMessage from "@/components/chat/message";

interface MessageListProps {
    // 类名
    className?: string;
    // 所有消息
    messages?: Message[]; // 修改为 message[]
    // 是否加载中
    isLoading?: boolean;
    // 消息是否在生成中
    isGenerating?: boolean;
}

export interface MessageListRef {
    scrollToBottom: () => void;
    isAtBottom: () => boolean;
    enableAutoScroll: () => void;
    disableAutoScroll: () => void;
}

const MessageList: React.ForwardRefRenderFunction<MessageListRef, MessageListProps> = ({
    className,
    messages = [],
    isLoading,
    isGenerating
}, ref) => {
    const containerRef = useRef<HTMLDivElement>(null);

    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        scrollToBottom: () => {
            if (containerRef.current) {
                containerRef.current.scrollTop = containerRef.current.scrollHeight;
            }
        },
        isAtBottom: () => {
            if (containerRef.current) {
                const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
                return scrollHeight - scrollTop - clientHeight < 10;
            }
            return true;
        },
        enableAutoScroll: () => {
            // 可以在这里添加自动滚动的逻辑
        },
        disableAutoScroll: () => {
            // 可以在这里添加禁用自动滚动的逻辑
        }
    }));

    // todo 获取组件到可显示界面底部的距离，给scrollButton的bottom设置为这个距离

    return (
        <div ref={containerRef} className={`${className} ${styles.MessageList}`}>
            {/* 消息 */}
            <div className={styles.content}>
                {
                    messages.map((message: Message, index: number) => (
                        <div key={index}>
                            <ChatMessage message={message}/>
                        </div>
                    ))
                }
            </div>
        </div>
    )
};

MessageList.displayName = 'MessageList';

export default forwardRef(MessageList);