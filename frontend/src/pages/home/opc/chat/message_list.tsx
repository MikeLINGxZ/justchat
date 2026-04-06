import React, { useEffect, useRef } from 'react';
import { Spin } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { OPCPersonView } from '@/stores/opcStore';
import MessageBubble from './message_bubble';
import TypingIndicator from './typing_indicator';

interface ChatMessage {
    message_uuid: string;
    role: string;
    content: string;
    sender_person_uuid: string;
    created_at: string;
}

interface MessageListProps {
    messages: ChatMessage[];
    persons: OPCPersonView[];
    typingPersons: string[];
    isGroup: boolean;
    isLoading: boolean;
}

const MessageList: React.FC<MessageListProps> = ({
    messages,
    persons,
    typingPersons,
    isGroup,
    isLoading,
}) => {
    const { t } = useTranslation();
    const listRef = useRef<HTMLDivElement>(null);

    // 自动滚动到底部
    useEffect(() => {
        if (listRef.current) {
            listRef.current.scrollTop = listRef.current.scrollHeight;
        }
    }, [messages, typingPersons]);

    const getPersonInfo = (personUuid: string): OPCPersonView | undefined => {
        return persons.find(p => p.uuid === personUuid);
    };

    const formatTime = (timeStr: string) => {
        const date = new Date(timeStr);
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    };

    // 是否需要显示时间分割
    const shouldShowTimeDivider = (current: ChatMessage, prev?: ChatMessage) => {
        if (!prev) return true;
        const currentTime = new Date(current.created_at).getTime();
        const prevTime = new Date(prev.created_at).getTime();
        return currentTime - prevTime > 5 * 60 * 1000; // 5 分钟间隔
    };

    if (isLoading) {
        return (
            <div className="opc-message-list loading">
                <Spin />
            </div>
        );
    }

    return (
        <div className="opc-message-list" ref={listRef}>
            {messages.length === 0 && (
                <div className="message-list-empty">
                    <div className="empty-hint">{t('opc.chat.startConversation')}</div>
                </div>
            )}

            {messages.map((msg, index) => {
                const isUser = msg.role === 'user';
                const person = !isUser ? getPersonInfo(msg.sender_person_uuid) : undefined;
                const showTimeDivider = shouldShowTimeDivider(msg, messages[index - 1]);

                return (
                    <React.Fragment key={msg.message_uuid}>
                        {showTimeDivider && (
                            <div className="time-divider">
                                <span>{formatTime(msg.created_at)}</span>
                            </div>
                        )}
                        <MessageBubble
                            content={msg.content}
                            isUser={isUser}
                            senderName={person?.name}
                            senderAvatar={person?.avatar || person?.name?.charAt(0)}
                            showSender={isGroup && !isUser}
                            time={formatTime(msg.created_at)}
                        />
                    </React.Fragment>
                );
            })}

            {typingPersons.length > 0 && (
                <TypingIndicator names={typingPersons} />
            )}
        </div>
    );
};

export default MessageList;
