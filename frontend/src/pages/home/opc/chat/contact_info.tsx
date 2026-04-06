import React from 'react';
import { Button } from 'antd';
import { MessageOutlined, UserOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { OPCPersonView } from '@/stores/opcStore';

interface ContactInfoProps {
    person: OPCPersonView;
    onSendMessage: () => void;
}

const ContactInfo: React.FC<ContactInfoProps> = ({ person, onSendMessage }) => {
    const { t } = useTranslation();

    const renderAvatar = () => {
        const avatar = person.avatar || '';
        let display: React.ReactNode = person.name.charAt(0) || '?';
        let style: React.CSSProperties = { background: 'linear-gradient(135deg, #667eea, #764ba2)' };

        if (avatar.startsWith('image:')) {
            return (
                <div className="info-avatar" style={{ background: 'transparent', padding: 0, overflow: 'hidden' }}>
                    <img src={avatar.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: '50%' }} />
                </div>
            );
        } else if (avatar.startsWith('emoji:')) {
            display = avatar.slice(6);
        } else if (avatar.startsWith('color:')) {
            style = { background: avatar.slice(6) };
            display = person.name.charAt(0);
        } else if (avatar) {
            display = avatar;
        }

        return (
            <div className="info-avatar" style={style}>
                {display}
            </div>
        );
    };

    return (
        <div className="opc-contact-info">
            <div className="info-card">
                {renderAvatar()}
                <div className="info-name">{person.name}</div>
                {person.role && <div className="info-role">{person.role}</div>}

                <Button
                    type="primary"
                    icon={<MessageOutlined />}
                    onClick={onSendMessage}
                    style={{ marginTop: 24 }}
                >
                    {t('opc.contact.sendMessage')}
                </Button>
            </div>
        </div>
    );
};

export default ContactInfo;
