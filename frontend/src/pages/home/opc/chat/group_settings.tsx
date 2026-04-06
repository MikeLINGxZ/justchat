import React, { useEffect, useState } from 'react';
import { Button, Input, Select, message } from 'antd';
import { CloseOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import type { OPCPersonView, OPCGroupView } from '@/stores/opcStore';
import './group_settings.scss';

const { TextArea } = Input;

interface GroupSettingsProps {
    group: OPCGroupView;
    allPersons: OPCPersonView[];
    open: boolean;
    onClose: () => void;
    onUpdated: () => void;
}

const GroupSettings: React.FC<GroupSettingsProps> = ({
    group,
    allPersons,
    open,
    onClose,
    onUpdated,
}) => {
    const { t } = useTranslation();
    const [name, setName] = useState(group.name);
    const [description, setDescription] = useState(group.description);
    const [memberUuids, setMemberUuids] = useState<string[]>([]);
    const [saving, setSaving] = useState(false);
    const [addMemberOpen, setAddMemberOpen] = useState(false);

    useEffect(() => {
        setName(group.name);
        setDescription(group.description);
        setMemberUuids(group.members?.map(m => m.uuid) || []);
    }, [group]);

    const handleSave = async () => {
        if (!name.trim()) return;
        setSaving(true);
        try {
            await Service.OPCUpdateGroup({
                uuid: group.uuid,
                name: name.trim(),
                description: description,
                member_uuids: memberUuids,
            });
            message.success(t('opc.group.saveSuccess'));
            onUpdated();
        } catch (err) {
            message.error(t('opc.group.saveFailed'));
        } finally {
            setSaving(false);
        }
    };

    const handleRemoveMember = (uuid: string) => {
        setMemberUuids(prev => prev.filter(u => u !== uuid));
    };

    const handleAddMember = (uuid: string) => {
        if (!memberUuids.includes(uuid)) {
            setMemberUuids(prev => [...prev, uuid]);
        }
        setAddMemberOpen(false);
    };

    const members = memberUuids
        .map(uuid => allPersons.find(p => p.uuid === uuid))
        .filter(Boolean) as OPCPersonView[];

    const availablePersons = allPersons.filter(p => !memberUuids.includes(p.uuid));

    if (!open) return null;

    return (
        <div className="group-settings-panel">
            <div className="panel-header">
                <span className="panel-title">{t('opc.group.settings')}</span>
                <CloseOutlined className="panel-close" onClick={onClose} />
            </div>

            <div className="panel-body">
                <div className="setting-section">
                    <label className="setting-label">{t('opc.group.name')}</label>
                    <Input
                        value={name}
                        onChange={e => setName(e.target.value)}
                        placeholder={t('opc.group.namePlaceholder')}
                    />
                </div>

                <div className="setting-section">
                    <label className="setting-label">{t('opc.group.description')}</label>
                    <TextArea
                        value={description}
                        onChange={e => setDescription(e.target.value)}
                        placeholder={t('opc.group.descriptionPlaceholder')}
                        rows={2}
                    />
                </div>

                <div className="setting-section">
                    <div className="section-header">
                        <label className="setting-label">{t('opc.group.members')} ({members.length})</label>
                    </div>

                    <div className="member-list">
                        {members.map(member => (
                            <div key={member.uuid} className="member-item">
                                <div className="member-avatar">
                                    {renderPersonAvatar(member)}
                                </div>
                                <div className="member-info">
                                    <div className="member-name">{member.name}</div>
                                    {member.role && <div className="member-role">{member.role}</div>}
                                </div>
                                <Button
                                    type="text"
                                    size="small"
                                    danger
                                    icon={<DeleteOutlined />}
                                    onClick={() => handleRemoveMember(member.uuid)}
                                />
                            </div>
                        ))}
                    </div>

                    {addMemberOpen ? (
                        <Select
                            autoFocus
                            open
                            placeholder={t('opc.group.membersPlaceholder')}
                            options={availablePersons.map(p => ({
                                label: `${p.name}${p.role ? ` (${p.role})` : ''}`,
                                value: p.uuid,
                            }))}
                            onSelect={handleAddMember}
                            onBlur={() => setAddMemberOpen(false)}
                            style={{ width: '100%', marginTop: 8 }}
                        />
                    ) : (
                        <Button
                            type="dashed"
                            icon={<PlusOutlined />}
                            onClick={() => setAddMemberOpen(true)}
                            style={{ width: '100%', marginTop: 8 }}
                        >
                            {t('opc.group.addMember')}
                        </Button>
                    )}
                </div>
            </div>

            <div className="panel-footer">
                <Button type="primary" onClick={handleSave} loading={saving} block>
                    {t('common.save')}
                </Button>
            </div>
        </div>
    );
};

function renderPersonAvatar(person: OPCPersonView) {
    const avatar = person.avatar || '';
    if (avatar.startsWith('image:')) {
        return (
            <div className="avatar-sm" style={{ background: 'transparent', padding: 0, overflow: 'hidden' }}>
                <img src={avatar.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: 'inherit' }} />
            </div>
        );
    }
    if (avatar.startsWith('emoji:')) {
        return <div className="avatar-sm">{avatar.slice(6)}</div>;
    }
    if (avatar.startsWith('color:')) {
        return <div className="avatar-sm" style={{ background: avatar.slice(6) }}>{person.name.charAt(0)}</div>;
    }
    return <div className="avatar-sm">{avatar || person.name.charAt(0)}</div>;
}

export default GroupSettings;
