import React from 'react';
import { UserOutlined } from '@ant-design/icons';

interface MessageBubbleProps {
    content: string;
    isUser: boolean;
    senderName?: string;
    senderAvatar?: string;
    showSender: boolean;
    time: string;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({
    content,
    isUser,
    senderName,
    senderAvatar,
    showSender,
    time,
}) => {
    const renderAvatar = () => {
        // Agent 头像：解析 avatar 格式
        const avatarStr = senderAvatar || '';
        let display: React.ReactNode = <UserOutlined />;
        let style: React.CSSProperties = {};

        if (avatarStr.startsWith('image:')) {
            return (
                <div className="bubble-avatar">
                    <div className="avatar-circle" style={{ background: 'transparent', padding: 0, overflow: 'hidden' }}>
                        <img src={avatarStr.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                    </div>
                </div>
            );
        } else if (avatarStr.startsWith('emoji:')) {
            display = avatarStr.slice(6);
        } else if (avatarStr.startsWith('color:')) {
            style = { background: avatarStr.slice(6) };
            display = senderName?.charAt(0) || '?';
        } else if (avatarStr) {
            display = avatarStr;
        }

        return (
            <div className="bubble-avatar">
                <div className="avatar-circle" style={style}>{display}</div>
            </div>
        );
    };

    return (
        <div className={`message-bubble-row ${isUser ? 'user' : 'agent'}`}>
            {!isUser && renderAvatar()}
            <div className="bubble-content-wrapper">
                {showSender && senderName && (
                    <div className="bubble-sender-name">{senderName}</div>
                )}
                <div className={`bubble-content ${isUser ? 'user-bubble' : 'agent-bubble'}`}>
                    <div className="bubble-text">{content}</div>
                </div>
            </div>
        </div>
    );
};

export default MessageBubble;
