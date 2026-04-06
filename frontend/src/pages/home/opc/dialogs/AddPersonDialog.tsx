import React, { useEffect, useState } from 'react';
import { Form, Input, Modal, Popover, Select, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';

const { TextArea } = Input;

const PRESET_EMOJIS = ['😀', '👨‍💻', '👩‍🎨', '🧑‍💼', '👨‍🔬', '👩‍⚕️', '🧑‍🏫', '👨‍🍳', '👩‍🔧', '🧑‍🚀', '👨‍⚖️', '👩‍💻', '🤖', '🦊', '🐱', '🐶'];
const PRESET_COLORS = ['#667eea', '#764ba2', '#f093fb', '#f5576c', '#4facfe', '#43e97b', '#fa709a', '#fee140'];

interface EditPersonData {
    uuid: string;
    name: string;
    role: string;
    prompt: string;
    tools: string[];
    skills: string[];
    avatar: string;
}

interface AddPersonDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess: () => void;
    editData?: EditPersonData | null;
}

const AddPersonDialog: React.FC<AddPersonDialogProps> = ({ open, onClose, onSuccess, editData }) => {
    const { t } = useTranslation();
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [tools, setTools] = useState<{ id: string; name: string }[]>([]);
    const [skills, setSkills] = useState<{ name: string; description: string }[]>([]);
    const [avatar, setAvatar] = useState('');
    const [avatarPickerOpen, setAvatarPickerOpen] = useState(false);

    const isEdit = !!editData;

    useEffect(() => {
        if (open) {
            loadToolsAndSkills();
            if (editData) {
                form.setFieldsValue({
                    name: editData.name,
                    role: editData.role,
                    prompt: editData.prompt,
                    tools: editData.tools,
                    skills: editData.skills,
                });
                setAvatar(editData.avatar || '');
            } else {
                form.resetFields();
                setAvatar('');
            }
        }
    }, [open, editData]);

    const loadToolsAndSkills = async () => {
        try {
            const [toolList, skillList] = await Promise.all([
                Service.GetTools().catch(() => []),
                Service.ListSkills().catch(() => []),
            ]);
            setTools((toolList || []).map((t: any) => ({ id: t.id, name: t.name || t.id })));
            setSkills((skillList || []).map((s: any) => ({ name: s.name, description: s.description || '' })));
        } catch (err) {
            console.error('Failed to load tools/skills:', err);
        }
    };

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            const finalAvatar = avatar || values.name.charAt(0);
            const input = {
                uuid: isEdit ? editData!.uuid : '',
                name: values.name,
                role: values.role || '',
                prompt: values.prompt || '',
                tools: values.tools || [],
                skills: values.skills || [],
                avatar: finalAvatar,
            };

            if (isEdit) {
                await Service.OPCUpdatePerson(input);
                message.success(t('opc.person.editSuccess'));
            } else {
                await Service.OPCCreatePerson(input);
                message.success(t('opc.person.createSuccess'));
            }

            form.resetFields();
            setAvatar('');
            onSuccess();
            onClose();
        } catch (err: any) {
            if (err?.errorFields) return;
            message.error(isEdit ? t('opc.person.editFailed') : t('opc.person.createFailed'));
        } finally {
            setLoading(false);
        }
    };

    const handleClose = () => {
        form.resetFields();
        setAvatar('');
        onClose();
    };

    const getAvatarDisplay = () => {
        const name = form.getFieldValue('name') || '';
        if (avatar.startsWith('image:')) {
            return <img src={avatar.slice(6)} style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: '50%' }} />;
        }
        if (avatar.startsWith('emoji:')) return avatar.slice(6);
        if (avatar.startsWith('color:')) return name.charAt(0) || '?';
        return avatar || name.charAt(0) || '?';
    };

    const getAvatarBg = () => {
        if (avatar.startsWith('image:')) return 'transparent';
        if (avatar.startsWith('color:')) return avatar.slice(6);
        return 'linear-gradient(135deg, #667eea, #764ba2)';
    };

    const handleUploadAvatar = async () => {
        try {
            const result = await Service.OPCSelectAvatar();
            if (result) {
                setAvatar(result);
                setAvatarPickerOpen(false);
            }
        } catch (err) {
            message.error(t('opc.person.avatarUploadFailed'));
        }
    };

    const avatarPicker = (
        <div style={{ width: 240 }}>
            <div
                onClick={handleUploadAvatar}
                style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'center',
                    height: 36, border: '1px dashed var(--border-color)', borderRadius: 6,
                    cursor: 'pointer', marginBottom: 12, fontSize: 13,
                    color: 'var(--text-color-secondary)', transition: 'border-color 0.2s',
                }}
                onMouseEnter={(e) => (e.currentTarget.style.borderColor = 'var(--primary-color)')}
                onMouseLeave={(e) => (e.currentTarget.style.borderColor = 'var(--border-color)')}
            >
                {t('opc.person.uploadAvatar')}
            </div>
            <div style={{ marginBottom: 8, fontSize: 12, color: 'var(--text-color-secondary)' }}>
                {t('opc.person.selectEmoji')}
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(8, 1fr)', gap: 4, marginBottom: 12 }}>
                {PRESET_EMOJIS.map(emoji => (
                    <div
                        key={emoji}
                        onClick={() => { setAvatar(`emoji:${emoji}`); setAvatarPickerOpen(false); }}
                        style={{
                            width: 28, height: 28, display: 'flex', alignItems: 'center', justifyContent: 'center',
                            cursor: 'pointer', borderRadius: 4, fontSize: 16,
                            background: avatar === `emoji:${emoji}` ? 'var(--primary-color-light)' : 'transparent',
                        }}
                    >
                        {emoji}
                    </div>
                ))}
            </div>
            <div style={{ marginBottom: 8, fontSize: 12, color: 'var(--text-color-secondary)' }}>
                {t('opc.person.selectColor')}
            </div>
            <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                {PRESET_COLORS.map(color => (
                    <div
                        key={color}
                        onClick={() => { setAvatar(`color:${color}`); setAvatarPickerOpen(false); }}
                        style={{
                            width: 24, height: 24, borderRadius: '50%', background: color, cursor: 'pointer',
                            border: avatar === `color:${color}` ? '2px solid var(--text-color)' : '2px solid transparent',
                        }}
                    />
                ))}
            </div>
        </div>
    );

    return (
        <Modal
            title={isEdit ? t('opc.person.editTitle') : t('opc.person.addTitle')}
            open={open}
            onOk={handleSubmit}
            onCancel={handleClose}
            confirmLoading={loading}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
            width={520}
            destroyOnClose
        >
            <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
                {/* 头像选择 */}
                <div style={{ display: 'flex', justifyContent: 'center', marginBottom: 20 }}>
                    <Popover
                        content={avatarPicker}
                        trigger="click"
                        open={avatarPickerOpen}
                        onOpenChange={setAvatarPickerOpen}
                        placement="bottom"
                    >
                        <div
                            style={{
                                width: 64, height: 64, borderRadius: '50%', display: 'flex',
                                alignItems: 'center', justifyContent: 'center', fontSize: 28,
                                color: '#fff', cursor: 'pointer', background: getAvatarBg(),
                                transition: 'transform 0.2s', fontWeight: 500,
                                border: '2px dashed var(--border-color)',
                            }}
                            onMouseEnter={(e) => (e.currentTarget.style.transform = 'scale(1.05)')}
                            onMouseLeave={(e) => (e.currentTarget.style.transform = 'scale(1)')}
                        >
                            {getAvatarDisplay()}
                        </div>
                    </Popover>
                </div>

                <Form.Item
                    label={t('opc.person.name')}
                    name="name"
                    rules={[{ required: true, message: t('opc.person.nameRequired') }]}
                >
                    <Input placeholder={t('opc.person.namePlaceholder')} />
                </Form.Item>

                <Form.Item label={t('opc.person.role')} name="role">
                    <Input placeholder={t('opc.person.rolePlaceholder')} />
                </Form.Item>

                <Form.Item label={t('opc.person.tools')} name="tools">
                    <Select
                        mode="multiple"
                        placeholder={t('opc.person.toolsPlaceholder')}
                        options={tools.map(t => ({ label: t.name, value: t.id }))}
                        allowClear
                    />
                </Form.Item>

                <Form.Item label={t('opc.person.skills')} name="skills">
                    <Select
                        mode="multiple"
                        placeholder={t('opc.person.skillsPlaceholder')}
                        options={skills.map(s => ({ label: s.name, value: s.name }))}
                        allowClear
                    />
                </Form.Item>

                <Form.Item label={t('opc.person.prompt')} name="prompt">
                    <TextArea rows={4} placeholder={t('opc.person.promptPlaceholder')} />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default AddPersonDialog;
