import React, { useEffect, useState } from 'react';
import { Form, Input, Modal, Select, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { useOPCStore } from '@/stores/opcStore';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';

const { TextArea } = Input;

interface EditGroupData {
    uuid: string;
    name: string;
    description: string;
    member_uuids: string[];
}

interface CreateGroupDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess: () => void;
    editData?: EditGroupData | null;
}

const CreateGroupDialog: React.FC<CreateGroupDialogProps> = ({ open, onClose, onSuccess, editData }) => {
    const { t } = useTranslation();
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const { persons } = useOPCStore();

    const isEdit = !!editData;

    useEffect(() => {
        if (open && editData) {
            form.setFieldsValue({
                name: editData.name,
                description: editData.description,
                members: editData.member_uuids,
            });
        } else if (open) {
            form.resetFields();
        }
    }, [open, editData]);

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);

            const input = {
                uuid: isEdit ? editData!.uuid : '',
                name: values.name,
                description: values.description || '',
                member_uuids: values.members || [],
            };

            if (isEdit) {
                await Service.OPCUpdateGroup(input);
                message.success(t('opc.group.editSuccess'));
            } else {
                await Service.OPCCreateGroup(input);
                message.success(t('opc.group.createSuccess'));
            }

            form.resetFields();
            onSuccess();
            onClose();
        } catch (err: any) {
            if (err?.errorFields) return;
            message.error(isEdit ? t('opc.group.editFailed') : t('opc.group.createFailed'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            title={isEdit ? t('opc.group.editTitle') : t('opc.group.createTitle')}
            open={open}
            onOk={handleSubmit}
            onCancel={() => { form.resetFields(); onClose(); }}
            confirmLoading={loading}
            okText={t('common.confirm')}
            cancelText={t('common.cancel')}
            width={480}
            destroyOnClose
        >
            <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
                <Form.Item
                    label={t('opc.group.name')}
                    name="name"
                    rules={[{ required: true, message: t('opc.group.nameRequired') }]}
                >
                    <Input placeholder={t('opc.group.namePlaceholder')} />
                </Form.Item>

                <Form.Item label={t('opc.group.description')} name="description">
                    <TextArea rows={2} placeholder={t('opc.group.descriptionPlaceholder')} />
                </Form.Item>

                <Form.Item
                    label={t('opc.group.members')}
                    name="members"
                    rules={[{ required: true, message: t('opc.group.membersRequired') }]}
                >
                    <Select
                        mode="multiple"
                        placeholder={t('opc.group.membersPlaceholder')}
                        options={persons.map(p => ({
                            label: `${p.name}${p.role ? ` (${p.role})` : ''}`,
                            value: p.uuid,
                        }))}
                        allowClear
                    />
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default CreateGroupDialog;
