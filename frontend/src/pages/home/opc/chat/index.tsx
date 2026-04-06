import React, { useCallback, useEffect, useRef, useState } from 'react';
import { message } from 'antd';
import { EllipsisOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
import { useOPCStore } from '@/stores/opcStore';
import type { OPCPersonView, OPCGroupView } from '@/stores/opcStore';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { useModelStore } from '@/stores/modelStore';
import MessageList from './message_list';
import ChatInput from './chat_input';
import GroupSettings from './group_settings';
import PersonSettings from './person_settings';
import ContactInfo from './contact_info';
import './index.scss';

interface OPCChatAreaProps {
    isSidebarCollapsed: boolean;
    onToggleSidebar: () => void;
}

interface ChatMessage {
    message_uuid: string;
    role: string;
    content: string;
    sender_person_uuid: string;
    created_at: string;
}

const OPCChatArea: React.FC<OPCChatAreaProps> = ({
    isSidebarCollapsed,
    onToggleSidebar,
}) => {
    const { t } = useTranslation();
    const { selectedType, selectedUuid, persons, groups, setSelected, setSidebarTab } = useOPCStore();
    const { models } = useModelStore();
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [typingPersons, setTypingPersons] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isSending, setIsSending] = useState(false);
    const [groupSettingsOpen, setGroupSettingsOpen] = useState(false);
    const [personInfoOpen, setPersonInfoOpen] = useState(false);
    const eventKeyRef = useRef<string | null>(null);

    // 获取当前选中的联系人信息
    const selectedPerson = selectedType === 'person'
        ? persons.find(p => p.uuid === selectedUuid)
        : null;
    const selectedGroup = selectedType === 'group'
        ? groups.find(g => g.uuid === selectedUuid)
        : null;

    const chatUuid = selectedPerson?.chat_uuid || selectedGroup?.chat_uuid || '';
    const chatTitle = selectedPerson?.name || selectedGroup?.name || '';

    // 加载聊天消息
    const loadMessages = useCallback(async () => {
        if (!chatUuid) {
            setMessages([]);
            return;
        }
        setIsLoading(true);
        try {
            const result = await Service.ChatMessages(chatUuid, 0, 100);
            if (result?.messages) {
                setMessages(result.messages.map((msg: any) => ({
                    message_uuid: msg.message_uuid,
                    role: msg.role,
                    content: msg.content,
                    sender_person_uuid: msg.sender_person_uuid || '',
                    created_at: msg.created_at,
                })));
            } else {
                setMessages([]);
            }
        } catch (err) {
            console.error('Failed to load messages:', err);
        } finally {
            setIsLoading(false);
        }
    }, [chatUuid]);

    useEffect(() => {
        loadMessages();
    }, [loadMessages]);

    // 切换联系人时清理事件监听和 typing 状态
    useEffect(() => {
        if (eventKeyRef.current) {
            Events.Off(eventKeyRef.current);
            eventKeyRef.current = null;
        }
        setTypingPersons([]);
        setIsSending(false);
        setPersonInfoOpen(false);
    }, [selectedUuid]);

    // 组件卸载时清理
    useEffect(() => {
        return () => {
            if (eventKeyRef.current) {
                Events.Off(eventKeyRef.current);
                eventKeyRef.current = null;
            }
        };
    }, []);

    // 发送消息
    const handleSend = async (content: string) => {
        if (!chatUuid || !content.trim()) return;

        const defaultModel = models.length > 0 ? models[0] : null;
        if (!defaultModel) {
            message.error(t('opc.chat.noModel'));
            return;
        }

        setIsSending(true);

        // 立即显示用户消息
        const userMsg: ChatMessage = {
            message_uuid: `temp-${Date.now()}`,
            role: 'user',
            content: content,
            sender_person_uuid: '',
            created_at: new Date().toISOString(),
        };
        setMessages(prev => [...prev, userMsg]);

        try {
            let result;

            if (selectedType === 'person' && selectedPerson) {
                result = await Service.OPCPersonChat({
                    chat_uuid: chatUuid,
                    person_uuid: selectedPerson.uuid,
                    content: content,
                    model_id: defaultModel.id,
                    model_name: defaultModel.name,
                });
            } else if (selectedType === 'group' && selectedGroup) {
                result = await Service.OPCGroupChat({
                    chat_uuid: chatUuid,
                    group_uuid: selectedGroup.uuid,
                    content: content,
                    model_id: defaultModel.id,
                    model_name: defaultModel.name,
                });
            }

            if (result?.event_key) {
                // 清理之前的监听
                if (eventKeyRef.current) {
                    Events.Off(eventKeyRef.current);
                }
                eventKeyRef.current = result.event_key;

                // 监听事件
                Events.On(result.event_key, (event: any) => {
                    const data = event?.data?.[0] || event?.data || event;
                    if (!data) return;

                    switch (data.type) {
                        case 'opc:typing': {
                            const personName = data.person_name || '';
                            setTypingPersons(prev =>
                                prev.includes(personName) ? prev : [...prev, personName]
                            );
                            break;
                        }
                        case 'opc:message': {
                            const msg = data.message;
                            if (msg) {
                                setMessages(prev => [...prev, {
                                    message_uuid: msg.message_uuid || msg.MessageUuid || `msg-${Date.now()}`,
                                    role: msg.role || msg.Role || 'assistant',
                                    content: msg.content || msg.Content || '',
                                    sender_person_uuid: msg.sender_person_uuid || msg.SenderPersonUuid || '',
                                    created_at: msg.created_at || msg.CreatedAt || new Date().toISOString(),
                                }]);
                                // 移除该人员的 typing 状态
                                const senderUuid = msg.sender_person_uuid || msg.SenderPersonUuid || '';
                                if (senderUuid) {
                                    const senderPerson = persons.find(p => p.uuid === senderUuid);
                                    if (senderPerson) {
                                        setTypingPersons(prev => prev.filter(n => n !== senderPerson.name));
                                    }
                                }
                            }
                            break;
                        }
                        case 'opc:complete': {
                            setTypingPersons([]);
                            setIsSending(false);
                            loadMessages();
                            // 刷新侧边栏数据（更新 last_message，使对话出现在对话列表中）
                            window.dispatchEvent(new CustomEvent('opc-refresh'));
                            break;
                        }
                    }
                });
            }
        } catch (err) {
            console.error('Failed to send message:', err);
            message.error(t('opc.chat.sendFailed'));
            setIsSending(false);
        }
    };

    // 空状态
    if (!selectedType || !selectedUuid) {
        return (
            <div className="opc-chat-empty">
                <div className="empty-icon">💬</div>
                <div className="empty-text">{t('opc.chat.selectContact')}</div>
            </div>
        );
    }

    // 联系人信息卡片（从联系人 tab 点击进入）
    if (selectedType === 'contact') {
        const contactPerson = persons.find(p => p.uuid === selectedUuid);
        if (!contactPerson) {
            return (
                <div className="opc-chat-empty">
                    <div className="empty-icon">💬</div>
                    <div className="empty-text">{t('opc.chat.selectContact')}</div>
                </div>
            );
        }
        return (
            <ContactInfo
                person={contactPerson}
                onSendMessage={() => {
                    setSidebarTab('conversations');
                    setSelected('person', selectedUuid);
                }}
            />
        );
    }

    return (
        <div className="opc-chat-area">
            <div className="opc-chat-main">
                <div className="opc-chat-header">
                    {isSidebarCollapsed && (
                        <button className="toggle-sidebar-btn" onClick={onToggleSidebar}>
                            ☰
                        </button>
                    )}
                    <div className="header-info">
                        <div className="header-name">{chatTitle}</div>
                        {selectedType === 'person' && selectedPerson && (
                            <div className="header-role">{selectedPerson.role}</div>
                        )}
                        {selectedType === 'group' && selectedGroup && (
                            <div className="header-role">
                                {selectedGroup.members?.length || 0} {t('opc.sidebar.members')}
                            </div>
                        )}
                    </div>
                    {selectedType === 'group' && !groupSettingsOpen && (
                        <button
                            className="header-settings-btn"
                            onClick={() => setGroupSettingsOpen(true)}
                        >
                            <EllipsisOutlined />
                        </button>
                    )}
                    {selectedType === 'person' && !personInfoOpen && (
                        <button
                            className="header-settings-btn"
                            onClick={() => setPersonInfoOpen(true)}
                        >
                            <EllipsisOutlined />
                        </button>
                    )}
                </div>

                <MessageList
                    messages={messages}
                    persons={persons}
                    typingPersons={typingPersons}
                    isGroup={selectedType === 'group'}
                    isLoading={isLoading}
                />

                <ChatInput
                    onSend={handleSend}
                    disabled={isSending}
                    placeholder={
                        selectedType === 'person'
                            ? t('opc.chat.inputPlaceholderPerson', { name: chatTitle })
                            : t('opc.chat.inputPlaceholderGroup')
                    }
                />
            </div>

            {selectedType === 'group' && selectedGroup && (
                <GroupSettings
                    group={selectedGroup}
                    allPersons={persons}
                    open={groupSettingsOpen}
                    onClose={() => setGroupSettingsOpen(false)}
                    onUpdated={() => {
                        window.dispatchEvent(new CustomEvent('opc-refresh'));
                    }}
                />
            )}

            {selectedType === 'person' && selectedPerson && (
                <PersonSettings
                    person={selectedPerson}
                    open={personInfoOpen}
                    onClose={() => setPersonInfoOpen(false)}
                    onUpdated={() => {
                        window.dispatchEvent(new CustomEvent('opc-refresh'));
                    }}
                />
            )}
        </div>
    );
};

export default OPCChatArea;
