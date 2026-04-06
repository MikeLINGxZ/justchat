import React, { useEffect, useState } from 'react';
import { Button, Input, message } from 'antd';
import { CloseOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import type { OPCPersonView } from '@/stores/opcStore';
import './group_settings.scss';

const { TextArea } = Input;

interface PersonSettingsProps {
    person: OPCPersonView;
    open: boolean;
    onClose: () => void;
    onUpdated: () => void;
}

const PersonSettings: React.FC<PersonSettingsProps> = ({
    person,
    open,
    onClose,
    onUpdated,
}) => {
    const { t } = useTranslation();
    const [name, setName] = useState(person.name);
    const [role, setRole] = useState(person.role);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        setName(person.name);
        setRole(person.role);
    }, [person]);

    const handleSave = async () => {
        if (!name.trim()) return;
        setSaving(true);
        try {
            // 先获取现有的 agent 信息以保留 prompt/tools/skills
            const agentDetail = await Service.GetAgent(person.agent_id).catch(() => null);
            await Service.OPCUpdatePerson({
                uuid: person.uuid,
                name: name.trim(),
                role: role,
                prompt: agentDetail?.prompts?.[0]?.content || '',
                tools: agentDetail?.tools || [],
                skills: agentDetail?.skills || [],
                avatar: person.avatar || '',
            });
            message.success(t('opc.person.editSuccess'));
            onUpdated();
        } catch (err) {
            message.error(t('opc.person.editFailed'));
        } finally {
            setSaving(false);
        }
    };

    if (!open) return null;

    const renderAvatar = () => {
        const avatar = person.avatar || '';
        let display: React.ReactNode = person.name.charAt(0) || '?';
        let style: React.CSSProperties = { background: 'linear-gradient(135deg, #667eea, #764ba2)', width: 64, height: 64, borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 28, fontWeight: 500, color: '#fff', margin: '0 auto 16px' };

        if (avatar.startsWith('image:')) {
            return (
                <div style={{ ...style, background: 'transparent', padding: 0, overflow: 'hidden' }}>
                    <img src={avatar.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: '50%' }} />
                </div>
            );
        } else if (avatar.startsWith('emoji:')) {
            display = avatar.slice(6);
        } else if (avatar.startsWith('color:')) {
            style = { ...style, background: avatar.slice(6) };
            display = person.name.charAt(0);
        } else if (avatar) {
            display = avatar;
        }

        return <div style={style}>{display}</div>;
    };

    return (
        <div className="group-settings-panel">
            <div className="panel-header">
                <span className="panel-title">{t('opc.person.info')}</span>
                <CloseOutlined className="panel-close" onClick={onClose} />
            </div>

            <div className="panel-body">
                <div style={{ textAlign: 'center', marginBottom: 20 }}>
                    {renderAvatar()}
                </div>

                <div className="setting-section">
                    <label className="setting-label">{t('opc.person.name')}</label>
                    <Input
                        value={name}
                        onChange={e => setName(e.target.value)}
                    />
                </div>

                <div className="setting-section">
                    <label className="setting-label">{t('opc.person.role')}</label>
                    <Input
                        value={role}
                        onChange={e => setRole(e.target.value)}
                    />
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

export default PersonSettings;
