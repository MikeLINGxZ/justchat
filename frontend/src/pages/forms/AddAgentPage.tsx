import React, { useEffect, useState } from 'react';
import { Button, Form, Input, message, Select, Spin } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import {
  CustomAgentInput,
  SkillSummary,
  Tool,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './formWindow.module.scss';

const WINDOW_NAME = 'window_form_agent';
const EVENT_KEY = 'settings:agents:changed';

const { TextArea } = Input;

const AddAgentPage: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [tools, setTools] = useState<Tool[]>([]);
  const [skillsList, setSkillsList] = useState<SkillSummary[]>([]);
  const [loadingOptions, setLoadingOptions] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    document.title = t('settings.agents.createTitle');
    void (async () => {
      setLoadingOptions(true);
      try {
        const [toolsResult, skillsResult] = await Promise.all([
          Service.GetTools(),
          Service.ListSkills(),
        ]);
        setTools(toolsResult || []);
        setSkillsList(skillsResult || []);
      } catch (err) {
        console.error('Failed to load form options:', err);
      } finally {
        setLoadingOptions(false);
      }
    })();
  }, [t]);

  const handleCancel = () => {
    void Service.CloseFormWindow(WINDOW_NAME);
  };

  const handleSubmit = async (values: any) => {
    setSubmitting(true);
    try {
      const input = new CustomAgentInput({
        id: values.id,
        name: values.name,
        description: values.description,
        prompt: values.prompt,
        tools: values.tools || [],
        skills: values.skills || [],
      });
      await Service.CreateCustomAgent(input);
      void Events.Emit(EVENT_KEY, { id: values.id });
      message.success(t('settings.agents.createSuccess'));
      void Service.CloseFormWindow(WINDOW_NAME);
    } catch (err) {
      console.error('Failed to create agent:', err);
      message.error(t('settings.agents.createFailed'));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className={styles.formWindow}>
      <div className={styles.header}>
        <h2>{t('settings.agents.createTitle')}</h2>
      </div>

      <div className={styles.body}>
        {loadingOptions ? (
          <div style={{ textAlign: 'center', padding: 48 }}>
            <Spin />
          </div>
        ) : (
          <Form form={form} layout="vertical" onFinish={handleSubmit}>
            <Form.Item
              name="name"
              label={t('settings.agents.form.name')}
              rules={[{ required: true }]}
            >
              <Input placeholder={t('settings.agents.form.namePlaceholder')} />
            </Form.Item>
            <Form.Item
              name="id"
              label={t('settings.agents.form.id')}
              rules={[
                { required: true },
                { pattern: /^[a-zA-Z0-9_-]+$/, message: t('settings.agents.form.idInvalid') },
              ]}
              extra={t('settings.agents.form.idHint')}
            >
              <Input placeholder={t('settings.agents.form.idPlaceholder')} />
            </Form.Item>
            <Form.Item
              name="description"
              label={t('settings.agents.form.description')}
              rules={[{ required: true }]}
            >
              <Input placeholder={t('settings.agents.form.descriptionPlaceholder')} />
            </Form.Item>
            <Form.Item
              name="prompt"
              label={t('settings.agents.form.prompt')}
              rules={[{ required: true }]}
            >
              <TextArea rows={6} placeholder={t('settings.agents.form.promptPlaceholder')} />
            </Form.Item>
            <Form.Item name="tools" label={t('settings.agents.form.tools')}>
              <Select
                mode="multiple"
                placeholder={t('settings.agents.form.toolsPlaceholder')}
                options={tools.map(tool => ({ label: tool.name, value: tool.id }))}
                allowClear
              />
            </Form.Item>
            <Form.Item name="skills" label={t('settings.agents.form.skills')}>
              <Select
                mode="multiple"
                placeholder={t('settings.agents.form.skillsPlaceholder')}
                options={skillsList.map(s => ({
                  label: `${s.name} - ${s.description}`,
                  value: s.name,
                }))}
                allowClear
              />
            </Form.Item>
          </Form>
        )}
      </div>

      <div className={styles.footer}>
        <Button onClick={handleCancel}>{t('common.cancel')}</Button>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          loading={submitting}
          onClick={() => form.submit()}
        >
          {t('settings.agents.actions.create')}
        </Button>
      </div>
    </div>
  );
};

export default AddAgentPage;
