import React, { useState, useRef, useCallback } from 'react';
import {
    FolderOpenOutlined,
    ScissorOutlined,
    SendOutlined,
    SmileOutlined,
    UserAddOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface ChatInputProps {
    onSend: (content: string) => void;
    disabled: boolean;
    placeholder: string;
}

const ChatInput: React.FC<ChatInputProps> = ({ onSend, disabled, placeholder }) => {
    const { t } = useTranslation();
    const [value, setValue] = useState('');
    const textareaRef = useRef<HTMLTextAreaElement>(null);

    const handleSend = useCallback(() => {
        const content = value.trim();
        if (!content || disabled) return;
        onSend(content);
        setValue('');
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
        }
    }, [value, disabled, onSend]);

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        setValue(e.target.value);
        // 自动调整高度
        const textarea = e.target;
        textarea.style.height = 'auto';
        textarea.style.height = `${Math.min(textarea.scrollHeight, 120)}px`;
    };

    return (
        <div className="opc-chat-input">
            <div className="input-shell">
                <div className="input-wrapper">
                    <textarea
                        ref={textareaRef}
                        value={value}
                        onChange={handleInput}
                        onKeyDown={handleKeyDown}
                        placeholder={placeholder}
                        disabled={disabled}
                        rows={1}
                        className="message-textarea"
                    />
                </div>

                {/*<div className="input-toolbar" aria-hidden="true">*/}
                {/*    <button className="tool-btn" type="button" tabIndex={-1}>*/}
                {/*        <SmileOutlined />*/}
                {/*    </button>*/}
                {/*    <button className="tool-btn" type="button" tabIndex={-1}>*/}
                {/*        <FolderOpenOutlined />*/}
                {/*    </button>*/}
                {/*    <button className="tool-btn" type="button" tabIndex={-1}>*/}
                {/*        <ScissorOutlined />*/}
                {/*    </button>*/}
                {/*    <button className="tool-btn" type="button" tabIndex={-1}>*/}
                {/*        <UserAddOutlined />*/}
                {/*    </button>*/}
                {/*</div>*/}

                <div className="input-actions">
                    <button
                        className={`send-btn ${value.trim() && !disabled ? 'active' : ''}`}
                        onClick={handleSend}
                        disabled={!value.trim() || disabled}
                    >
                        <SendOutlined />
                        <span>{t('common.send', { defaultValue: '发送' })}</span>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ChatInput;
