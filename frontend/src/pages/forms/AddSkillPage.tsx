import React, { useEffect, useState } from 'react';
import { Button, Form, Input, message, Select } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { SkillDetail } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './formWindow.module.scss';

const WINDOW_NAME = 'window_form_skill';
const EVENT_KEY = 'settings:skills:changed';

const { TextArea } = Input;

const AddSkillPage: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    document.title = t('settings.skills.createTitle');
  }, [t]);

  const handleCancel = () => {
    void Service.CloseFormWindow(WINDOW_NAME);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      const input = new SkillDetail({
        name: values.name,
        description: values.description,
        when: values.when || '',
        version: values.version || '1.0',
        tags: values.tags || [],
        content: values.content,
      });
      const result = await Service.CreateSkill(input);
      if (result) {
        void Events.Emit(EVENT_KEY, { name: result.name });
        message.success(t('settings.skills.createSuccess'));
        void Service.CloseFormWindow(WINDOW_NAME);
      }
    } catch (error) {
      if ((error as { errorFields?: unknown }).errorFields) {
        return;
      }
      console.error('创建技能失败:', error);
      message.error(t('settings.skills.createFailed'));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className={styles.formWindow}>
      <div className={styles.header}>
        <h2>{t('settings.skills.createTitle')}</h2>
      </div>

      <div className={styles.body}>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="name"
            label={t('settings.skills.form.name')}
            rules={[
              { required: true },
              { pattern: /^[a-zA-Z0-9_-]+$/, message: t('settings.skills.form.nameInvalid') },
            ]}
          >
            <Input placeholder={t('settings.skills.form.namePlaceholder')} />
          </Form.Item>
          <Form.Item
            name="description"
            label={t('settings.skills.form.description')}
            rules={[{ required: true }]}
          >
            <Input placeholder={t('settings.skills.form.descriptionPlaceholder')} />
          </Form.Item>
          <Form.Item name="when" label={t('settings.skills.form.when')}>
            <Input placeholder={t('settings.skills.form.whenPlaceholder')} />
          </Form.Item>
          <Form.Item name="version" label={t('settings.skills.form.version')} initialValue="1.0">
            <Input />
          </Form.Item>
          <Form.Item name="tags" label={t('settings.skills.form.tags')}>
            <Select mode="tags" placeholder={t('settings.skills.form.tagsPlaceholder')} />
          </Form.Item>
          <Form.Item
            name="content"
            label={t('settings.skills.form.content')}
            rules={[{ required: true }]}
          >
            <TextArea rows={8} placeholder={t('settings.skills.form.contentPlaceholder')} />
          </Form.Item>
        </Form>
      </div>

      <div className={styles.footer}>
        <Button onClick={handleCancel}>{t('common.cancel')}</Button>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          loading={submitting}
          onClick={() => void handleSubmit()}
        >
          {t('settings.skills.actions.create')}
        </Button>
      </div>
    </div>
  );
};

export default AddSkillPage;
