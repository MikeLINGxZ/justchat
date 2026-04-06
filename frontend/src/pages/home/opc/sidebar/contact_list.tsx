import React from 'react';
import { Dropdown, Empty } from 'antd';
import type { MenuProps } from 'antd';
import { ClearOutlined, DeleteOutlined, EditOutlined, PushpinOutlined, TeamOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { OPCPersonView, OPCGroupView, OPCSidebarTab } from '@/stores/opcStore';

interface ContactListProps {
    persons: OPCPersonView[];
    groups: OPCGroupView[];
    searchQuery: string;
    selectedUuid: string | null;
    activeTab: OPCSidebarTab;
    onSelect: (type: 'person' | 'group' | 'contact', uuid: string) => void;
    onDeletePerson: (uuid: string) => void;
    onDeleteGroup: (uuid: string) => void;
    onTogglePin: (type: 'person' | 'group', uuid: string, pinned: boolean) => void;
    onEditPerson?: (uuid: string) => void;
    onEditGroup?: (uuid: string) => void;
    onClearConversation?: (chatUuid: string) => void;
}

type ContactItem = {
    type: 'person';
    data: OPCPersonView;
} | {
    type: 'group';
    data: OPCGroupView;
};

const ContactList: React.FC<ContactListProps> = ({
    persons,
    groups,
    searchQuery,
    selectedUuid,
    activeTab,
    onSelect,
    onDeletePerson,
    onDeleteGroup,
    onTogglePin,
    onEditPerson,
    onEditGroup,
    onClearConversation,
}) => {
    const { t } = useTranslation();

    const filteredPersons = searchQuery
        ? persons.filter(p => p.name.toLowerCase().includes(searchQuery.toLowerCase()))
        : persons;
    const filteredGroups = searchQuery
        ? groups.filter(g => g.name.toLowerCase().includes(searchQuery.toLowerCase()))
        : groups;

    let contacts: ContactItem[];
    if (activeTab === 'conversations') {
        // 对话 tab：显示有消息的人员，以及所有群聊（群聊是主动创建的，始终显示）
        contacts = [
            ...filteredPersons.filter(p => p.last_message).map(p => ({ type: 'person' as const, data: p })),
            ...filteredGroups.map(g => ({ type: 'group' as const, data: g })),
        ].sort((a, b) => {
            if (a.data.is_pinned && !b.data.is_pinned) return -1;
            if (!a.data.is_pinned && b.data.is_pinned) return 1;
            const aTime = a.data.last_message_at || a.data.updated_at;
            const bTime = b.data.last_message_at || b.data.updated_at;
            return new Date(bTime).getTime() - new Date(aTime).getTime();
        });
    } else {
        // 联系人 tab：只显示人员，不显示群聊
        contacts = filteredPersons
            .map(p => ({ type: 'person' as const, data: p }))
            .sort((a, b) => {
                if (a.data.is_pinned && !b.data.is_pinned) return -1;
                if (!a.data.is_pinned && b.data.is_pinned) return 1;
                return a.data.name.localeCompare(b.data.name);
            });
    }

    if (contacts.length === 0) {
        return (
            <div className="opc-contact-empty">
                <Empty
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                    description={
                        searchQuery
                            ? t('opc.sidebar.noSearchResults')
                            : activeTab === 'conversations'
                                ? t('opc.sidebar.noConversations')
                                : t('opc.sidebar.noContacts')
                    }
                />
            </div>
        );
    }

    const formatTime = (timeStr: string | null) => {
        if (!timeStr) return '';
        const date = new Date(timeStr);
        const now = new Date();
        if (date.toDateString() === now.toDateString()) {
            return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
        return date.toLocaleDateString([], { month: '2-digit', day: '2-digit' });
    };

    const getContextMenuItems = (item: ContactItem): MenuProps['items'] => {
        const uuid = item.data.uuid;
        const isPinned = item.data.is_pinned;
        const chatUuid = item.type === 'person'
            ? (item.data as OPCPersonView).chat_uuid
            : (item.data as OPCGroupView).chat_uuid;

        if (activeTab === 'conversations') {
            // 对话 tab：清除对话（不删除联系人）、置顶
            return [
                {
                    key: 'pin',
                    icon: <PushpinOutlined />,
                    label: isPinned ? t('opc.sidebar.unpin') : t('opc.sidebar.pin'),
                    onClick: () => onTogglePin(item.type, uuid, isPinned),
                },
                {
                    key: 'clear',
                    icon: <ClearOutlined />,
                    label: t('opc.sidebar.clearConversation'),
                    danger: true,
                    onClick: () => onClearConversation?.(chatUuid),
                },
            ];
        }

        // 联系人 tab：编辑、置顶、删除联系人
        return [
            {
                key: 'edit',
                icon: <EditOutlined />,
                label: t('opc.sidebar.edit'),
                onClick: () => {
                    if (item.type === 'person') onEditPerson?.(uuid);
                    else onEditGroup?.(uuid);
                },
            },
            {
                key: 'pin',
                icon: <PushpinOutlined />,
                label: isPinned ? t('opc.sidebar.unpin') : t('opc.sidebar.pin'),
                onClick: () => onTogglePin(item.type, uuid, isPinned),
            },
            {
                key: 'delete',
                icon: <DeleteOutlined />,
                label: t('common.delete'),
                danger: true,
                onClick: () => {
                    if (item.type === 'person') onDeletePerson(uuid);
                    else onDeleteGroup(uuid);
                },
            },
        ];
    };

    const renderAvatar = (item: ContactItem) => {
        if (item.type === 'group') {
            return <div className="avatar-circle group"><TeamOutlined /></div>;
        }
        const person = item.data as OPCPersonView;
        const avatar = person.avatar || '';
        if (avatar.startsWith('image:')) {
            return (
                <div className="avatar-circle person" style={{ background: 'transparent', padding: 0, overflow: 'hidden' }}>
                    <img src={avatar.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                </div>
            );
        }
        if (avatar.startsWith('emoji:')) {
            return <div className="avatar-circle person">{avatar.slice(6)}</div>;
        }
        if (avatar.startsWith('color:')) {
            return (
                <div className="avatar-circle person" style={{ background: avatar.slice(6) }}>
                    {person.name.charAt(0)}
                </div>
            );
        }
        return <div className="avatar-circle person">{avatar || person.name.charAt(0)}</div>;
    };

    const handleClick = (item: ContactItem) => {
        if (activeTab === 'contacts') {
            // 联系人 tab：点击显示联系人信息卡片
            onSelect('contact', item.data.uuid);
        } else {
            // 对话 tab：点击进入聊天
            onSelect(item.type, item.data.uuid);
        }
    };

    return (
        <div className="opc-contact-list">
            {contacts.map((item) => {
                const uuid = item.data.uuid;
                const isSelected = uuid === selectedUuid;
                const isPinned = item.data.is_pinned;

                return (
                    <Dropdown
                        key={`${item.type}-${uuid}`}
                        menu={{ items: getContextMenuItems(item) }}
                        trigger={['contextMenu']}
                    >
                        <div
                            className={`contact-item ${isSelected ? 'active' : ''} ${isPinned ? 'pinned' : ''}`}
                            onClick={() => handleClick(item)}
                        >
                            <div className="contact-avatar">
                                {renderAvatar(item)}
                            </div>
                            <div className="contact-info">
                                <div className="contact-name-row">
                                    <span className="contact-name">{item.data.name}</span>
                                    {activeTab === 'conversations' && (
                                        <span className="contact-time">
                                            {formatTime(item.data.last_message_at)}
                                        </span>
                                    )}
                                </div>
                                <div className="contact-preview">
                                    {activeTab === 'conversations'
                                        ? (item.type === 'person'
                                            ? (item.data as OPCPersonView).last_message || (item.data as OPCPersonView).role
                                            : (item.data as OPCGroupView).last_message || `${(item.data as OPCGroupView).members?.length || 0} ${t('opc.sidebar.members')}`)
                                        : (item.data as OPCPersonView).role
                                    }
                                </div>
                            </div>
                            {isPinned && <PushpinOutlined className="pin-icon" />}
                        </div>
                    </Dropdown>
                );
            })}
        </div>
    );
};

export default ContactList;
