import React, { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Empty,
  Form,
  Input,
  Modal,
  Segmented,
  Select,
  Skeleton,
  Space,
  Spin,
  Tag,
  Typography,
  message,
} from 'antd';
import {
  DeleteOutlined,
  FileTextOutlined,
  PlusOutlined,
  ReloadOutlined,
  RollbackOutlined,
  SaveOutlined,
} from '@ant-design/icons';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import {
  AgentDetail,
  AgentSummary,
  CustomAgentInput,
  SkillSummary,
  Tool,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Paragraph, Text, Title } = Typography;
const { TextArea } = Input;

const ROLE_TAG_COLORS: Record<string, string> = {
  main: 'blue',
  workflow: 'purple',
  planner: 'cyan',
  worker: 'green',
  reviewer: 'orange',
  synthesizer: 'geekblue',
};

interface AgentSettingsPageProps {
  className?: string;
}

const AgentSettingsPage: React.FC<AgentSettingsPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [agents, setAgents] = useState<AgentSummary[]>([]);
  const [activeAgent, setActiveAgent] = useState<string>('');
  const [detail, setDetail] = useState<AgentDetail | null>(null);
  const [activePromptIdx, setActivePromptIdx] = useState(0);
  const [draft, setDraft] = useState('');
  const [loadingList, setLoadingList] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [saving, setSaving] = useState(false);
  const [resetting, setResetting] = useState(false);
  const [listError, setListError] = useState('');
  const [detailError, setDetailError] = useState('');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showEditorOnMobile, setShowEditorOnMobile] = useState(false);

  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [creating, setCreating] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [tools, setTools] = useState<Tool[]>([]);
  const [skillsList, setSkillsList] = useState<SkillSummary[]>([]);
  const [form] = Form.useForm();

  const [customDraft, setCustomDraft] = useState<{
    name: string;
    description: string;
    prompt: string;
    tools: string[];
    skills: string[];
  } | null>(null);

  useEffect(() => {
    const handleResize = () => setIsMobile(isMobileDevice());
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const currentPrompt = useMemo(() => {
    if (!detail || !detail.prompts || detail.prompts.length === 0) return null;
    return detail.prompts[activePromptIdx] || detail.prompts[0];
  }, [detail, activePromptIdx]);

  const isDirty = useMemo(() => {
    if (detail && detail.agent_type === 'custom' && customDraft) {
      return (
        customDraft.name !== detail.display_name ||
        customDraft.description !== detail.description ||
        customDraft.prompt !== (detail.prompts?.[0]?.content || '') ||
        JSON.stringify(customDraft.tools) !== JSON.stringify(detail.tools || []) ||
        JSON.stringify(customDraft.skills) !== JSON.stringify(detail.skills || [])
      );
    }
    return currentPrompt !== null && draft !== currentPrompt.content;
  }, [currentPrompt, draft, detail, customDraft]);

  const loadFormOptions = async () => {
    try {
      const [toolsResult, skillsResult] = await Promise.all([
        Service.GetTools(),
        Service.ListSkills(),
      ]);
      setTools(toolsResult);
      setSkillsList(skillsResult);
    } catch (e) {
      console.error('Failed to load form options:', e);
    }
  };

  const refreshList = async () => {
    setLoadingList(true);
    setListError('');
    try {
      const result = await Service.ListAgents();
      setAgents(result);
      if (!activeAgent && result.length > 0) {
        setActiveAgent(result[0].name);
      }
    } catch (error) {
      console.error('Failed to load agents:', error);
      setListError(t('settings.agents.listLoadFailedDesc'));
    } finally {
      setLoadingList(false);
    }
  };

  const loadDetail = async (name: string) => {
    if (!name) return;
    setLoadingDetail(true);
    setDetailError('');
    try {
      const result = await Service.GetAgent(name);
      setDetail(result);
      setActivePromptIdx(0);
      setDraft(result?.prompts?.[0]?.content || '');
      setActiveAgent(name);
      if (result && result.agent_type === 'custom') {
        setCustomDraft({
          name: result.display_name,
          description: result.description,
          prompt: result.prompts?.[0]?.content || '',
          tools: result.tools || [],
          skills: result.skills || [],
        });
        loadFormOptions();
      } else {
        setCustomDraft(null);
      }
      if (isMobile) setShowEditorOnMobile(true);
    } catch (error) {
      console.error('Failed to load agent:', error);
      setDetail(null);
      setDraft('');
      setCustomDraft(null);
      setDetailError(t('settings.agents.detailLoadFailedDesc'));
    } finally {
      setLoadingDetail(false);
    }
  };

  useEffect(() => {
    void refreshList();
  }, []);

  useEffect(() => {
    if (activeAgent) void loadDetail(activeAgent);
  }, [activeAgent]);

  const confirmDiscardIfNeeded = async (onConfirm: () => void) => {
    if (!isDirty) {
      onConfirm();
      return;
    }
    Modal.confirm({
      title: t('settings.agents.unsavedConfirmTitle'),
      content: t('settings.agents.unsavedConfirmContent'),
      okText: t('settings.agents.discardChanges'),
      cancelText: t('settings.agents.continueEditing'),
      okButtonProps: { danger: true },
      onOk: () => onConfirm(),
    });
  };

  const handleSelectAgent = async (name: string) => {
    if (name === activeAgent && detail) {
      if (isMobile) setShowEditorOnMobile(true);
      return;
    }
    await confirmDiscardIfNeeded(() => setActiveAgent(name));
  };

  const handleSelectPrompt = async (idx: number) => {
    if (idx === activePromptIdx) return;
    await confirmDiscardIfNeeded(() => {
      setActivePromptIdx(idx);
      if (detail?.prompts?.[idx]) {
        setDraft(detail.prompts[idx].content);
      }
    });
  };

  const handleDiscard = () => {
    if (!currentPrompt) return;
    setDraft(currentPrompt.content);
    message.success(t('settings.agents.discardSuccess'));
  };

  const handleSave = async () => {
    if (!detail || !currentPrompt) return;
    const content = draft.trim();
    if (!content) {
      message.error(t('settings.agents.contentRequired'));
      return;
    }
    setSaving(true);
    try {
      const result = await Service.UpdateAgentPrompt(detail.name, currentPrompt.name, content);
      setDetail(result);
      const promptContent = result?.prompts?.[activePromptIdx]?.content || '';
      setDraft(promptContent);
      message.success(t('settings.agents.saveSuccess'));
    } catch (error) {
      console.error('Failed to save agent prompt:', error);
      message.error(t('settings.agents.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleReset = async () => {
    if (!detail || !currentPrompt) return;
    setResetting(true);
    try {
      const result = await Service.ResetAgentPrompt(detail.name, currentPrompt.name);
      setDetail(result);
      const promptContent = result?.prompts?.[activePromptIdx]?.content || '';
      setDraft(promptContent);
      message.success(t('settings.agents.resetSuccess'));
    } catch (error) {
      console.error('Failed to reset agent prompt:', error);
      message.error(t('settings.agents.resetFailed'));
    } finally {
      setResetting(false);
    }
  };

  const handleCreate = async (values: any) => {
    setCreating(true);
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
      message.success(t('settings.agents.createSuccess'));
      setCreateModalVisible(false);
      form.resetFields();
      await refreshList();
      setActiveAgent(values.id);
    } catch (e) {
      console.error('Failed to create agent:', e);
      message.error(t('settings.agents.createFailed'));
    } finally {
      setCreating(false);
    }
  };

  const handleDeleteAgent = async () => {
    if (!detail || detail.agent_type !== 'custom') return;
    Modal.confirm({
      title: t('settings.agents.deleteConfirmTitle'),
      content: t('settings.agents.deleteConfirmContent'),
      okText: t('common.delete'),
      okButtonProps: { danger: true },
      cancelText: t('common.cancel'),
      onOk: async () => {
        setDeleting(true);
        try {
          await Service.DeleteCustomAgent(detail.name);
          message.success(t('settings.agents.deleteSuccess'));
          setDetail(null);
          setCustomDraft(null);
          setActiveAgent('');
          await refreshList();
        } catch (e) {
          console.error('Failed to delete agent:', e);
          message.error(t('settings.agents.deleteFailed'));
        } finally {
          setDeleting(false);
        }
      },
    });
  };

  const handleSaveCustomAgent = async () => {
    if (!detail || !customDraft) return;
    setSaving(true);
    try {
      const input = new CustomAgentInput({
        id: detail.name,
        name: customDraft.name,
        description: customDraft.description,
        prompt: customDraft.prompt,
        tools: customDraft.tools,
        skills: customDraft.skills,
      });
      const result = await Service.UpdateCustomAgent(input);
      setDetail(result);
      if (result) {
        setCustomDraft({
          name: result.display_name,
          description: result.description,
          prompt: result.prompts?.[0]?.content || '',
          tools: result.tools || [],
          skills: result.skills || [],
        });
      }
      message.success(t('settings.agents.saveSuccess'));
      await refreshList();
    } catch (e) {
      console.error('Failed to update agent:', e);
      message.error(t('settings.agents.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const getRoleLabel = (role: string) => {
    const key = `settings.agents.tags.${role}` as const;
    return t(key) || role;
  };

  const renderAgentList = () => (
    <Card
      className={styles.listCard}
      title={
        <div className={styles.listTitleRow}>
          <span>{t('settings.agents.listTitle')}</span>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            onClick={() => { loadFormOptions(); setCreateModalVisible(true); }}
            className={styles.addButton}
          />
        </div>
      }
    >
      {loadingList ? (
        <div className={styles.listLoading}>
          <Skeleton active paragraph={{ rows: 6 }} />
        </div>
      ) : listError ? (
        <Alert
          type="error"
          showIcon
          message={t('settings.agents.listLoadFailed')}
          description={listError}
          action={
            <Button size="small" onClick={() => void refreshList()}>
              {t('common.retry')}
            </Button>
          }
        />
      ) : (
        <div className={styles.agentList}>
          {agents.map(agent => {
            const selected = agent.name === activeAgent;
            return (
              <button
                key={agent.name}
                type="button"
                className={`${styles.agentItem} ${selected ? styles.selected : ''}`}
                onClick={() => void handleSelectAgent(agent.name)}
              >
                <div className={styles.agentItemHeader}>
                  <span className={styles.agentItemTitle}>
                    {agent.agent_type === 'custom' ? agent.display_name : agent.description}
                  </span>
                  <div className={styles.agentItemTags}>
                    <Tag
                      color={ROLE_TAG_COLORS[agent.agent_role] || 'default'}
                      bordered={false}
                    >
                      {getRoleLabel(agent.agent_role)}
                    </Tag>
                    <Tag
                      color={agent.agent_type === 'system' ? 'blue' : 'gold'}
                      bordered={false}
                    >
                      {t(`settings.agents.tags.${agent.agent_type}`)}
                    </Tag>
                  </div>
                </div>
                <div className={styles.agentItemName}>{agent.name}</div>
              </button>
            );
          })}
        </div>
      )}
    </Card>
  );

  const renderCustomAgentEditor = () => {
    if (!detail || !customDraft) return null;

    return (
      <div className={styles.customEditorForm}>
        <div className={styles.customEditorHeader}>
          <div className={styles.customEditorTitle}>
            <Title level={4}>{detail.display_name}</Title>
            <Tag color="gold" bordered={false}>
              {t('settings.agents.tags.custom')}
            </Tag>
            {isDirty && (
              <Tag color="orange" bordered={false}>
                {t('settings.prompt.tags.unsaved')}
              </Tag>
            )}
          </div>
          <Button
            danger
            icon={<DeleteOutlined />}
            onClick={() => void handleDeleteAgent()}
            loading={deleting}
          >
            {t('settings.agents.actions.delete')}
          </Button>
        </div>

        <div className={styles.customFormField}>
          <label>{t('settings.agents.form.name')}</label>
          <Input
            value={customDraft.name}
            onChange={e => setCustomDraft({ ...customDraft, name: e.target.value })}
            placeholder={t('settings.agents.form.namePlaceholder')}
          />
        </div>

        <div className={styles.customFormField}>
          <label>{t('settings.agents.form.description')}</label>
          <Input
            value={customDraft.description}
            onChange={e => setCustomDraft({ ...customDraft, description: e.target.value })}
            placeholder={t('settings.agents.form.descriptionPlaceholder')}
          />
        </div>

        <div className={styles.customFormField} style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
          <label>{t('settings.agents.form.prompt')}</label>
          <TextArea
            value={customDraft.prompt}
            onChange={e => setCustomDraft({ ...customDraft, prompt: e.target.value })}
            placeholder={t('settings.agents.form.promptPlaceholder')}
            className={styles.formTextArea}
            autoSize={false}
            style={{ flex: 1 }}
          />
        </div>

        <div className={styles.customFormField}>
          <label>{t('settings.agents.form.tools')}</label>
          <Select
            mode="multiple"
            value={customDraft.tools}
            onChange={val => setCustomDraft({ ...customDraft, tools: val })}
            placeholder={t('settings.agents.form.toolsPlaceholder')}
            options={tools.map(t => ({ label: t.name, value: t.id }))}
            allowClear
            style={{ width: '100%' }}
          />
        </div>

        <div className={styles.customFormField}>
          <label>{t('settings.agents.form.skills')}</label>
          <Select
            mode="multiple"
            value={customDraft.skills}
            onChange={val => setCustomDraft({ ...customDraft, skills: val })}
            placeholder={t('settings.agents.form.skillsPlaceholder')}
            options={skillsList.map(s => ({ label: `${s.name} - ${s.description}`, value: s.name }))}
            allowClear
            style={{ width: '100%' }}
          />
        </div>

        <div className={styles.customEditorActions}>
          <Button
            type="primary"
            icon={<SaveOutlined />}
            onClick={() => void handleSaveCustomAgent()}
            loading={saving}
            disabled={!isDirty}
          >
            {t('settings.agents.actions.save')}
          </Button>
        </div>
      </div>
    );
  };

  const renderCreateModal = () => (
    <Modal
      title={t('settings.agents.createTitle')}
      open={createModalVisible}
      onCancel={() => setCreateModalVisible(false)}
      onOk={() => form.submit()}
      confirmLoading={creating}
      okText={t('settings.agents.actions.create')}
      cancelText={t('common.cancel')}
      destroyOnClose
      width={600}
    >
      <Form form={form} layout="vertical" onFinish={handleCreate}>
        <Form.Item name="name" label={t('settings.agents.form.name')} rules={[{ required: true }]}>
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
        <Form.Item name="description" label={t('settings.agents.form.description')} rules={[{ required: true }]}>
          <Input placeholder={t('settings.agents.form.descriptionPlaceholder')} />
        </Form.Item>
        <Form.Item name="prompt" label={t('settings.agents.form.prompt')} rules={[{ required: true }]}>
          <TextArea rows={4} placeholder={t('settings.agents.form.promptPlaceholder')} className={styles.formTextArea} />
        </Form.Item>
        <Form.Item name="tools" label={t('settings.agents.form.tools')}>
          <Select
            mode="multiple"
            placeholder={t('settings.agents.form.toolsPlaceholder')}
            options={tools.map(t => ({ label: t.name, value: t.id }))}
            allowClear
          />
        </Form.Item>
        <Form.Item name="skills" label={t('settings.agents.form.skills')}>
          <Select
            mode="multiple"
            placeholder={t('settings.agents.form.skillsPlaceholder')}
            options={skillsList.map(s => ({ label: `${s.name} - ${s.description}`, value: s.name }))}
            allowClear
          />
        </Form.Item>
      </Form>
    </Modal>
  );

  const renderEditorBody = () => {
    if (loadingDetail) {
      return (
        <div className={styles.editorLoading}>
          <Spin />
        </div>
      );
    }

    if (detailError) {
      return (
        <Alert
          type="error"
          showIcon
          message={t('settings.agents.detailLoadFailed')}
          description={detailError}
          action={
            <Button size="small" onClick={() => activeAgent && void loadDetail(activeAgent)}>
              {t('common.retry')}
            </Button>
          }
        />
      );
    }

    if (!detail) {
      return (
        <div className={styles.emptyState}>
          <Empty description={t('settings.agents.empty')} />
        </div>
      );
    }

    if (detail.agent_type === 'custom') {
      return renderCustomAgentEditor();
    }

    const hasMultiplePrompts = detail.prompts && detail.prompts.length > 1;

    return (
      <>
        <div className={styles.editorHeader}>
          <div>
            <div className={styles.editorTitleRow}>
              <Title level={4}>{detail.description}</Title>
              <Tag color={ROLE_TAG_COLORS[detail.agent_role] || 'default'} bordered={false}>
                {getRoleLabel(detail.agent_role)}
              </Tag>
              <Tag color={detail.agent_type === 'system' ? 'blue' : 'gold'} bordered={false}>
                {t(`settings.agents.tags.${detail.agent_type}`)}
              </Tag>
              {isDirty && (
                <Tag color="orange" bordered={false}>
                  {t('settings.prompt.tags.unsaved')}
                </Tag>
              )}
            </div>
            <Text className={styles.agentName}>{detail.name}</Text>
          </div>

          {hasMultiplePrompts && (
            <div className={styles.promptTabs}>
              <Text className={styles.promptTabsLabel}>{t('settings.agents.promptSelector')}</Text>
              <Segmented
                value={activePromptIdx}
                options={detail.prompts.map((p, idx) => ({
                  label: p.title,
                  value: idx,
                }))}
                onChange={(val) => void handleSelectPrompt(val as number)}
              />
            </div>
          )}

          {currentPrompt && (
            <Paragraph className={styles.editorDescription}>
              {currentPrompt.description}
            </Paragraph>
          )}

          <Alert
            type="info"
            showIcon
            className={styles.tipAlert}
            message={t('settings.agents.tip')}
          />
        </div>

        <div className={styles.editorArea}>
          <TextArea
            value={draft}
            onChange={event => setDraft(event.target.value)}
            spellCheck={false}
            autoSize={false}
            className={styles.textArea}
          />
        </div>

        <div className={styles.editorActions}>
          <div className={styles.actionHint}>
            <FileTextOutlined />
            <span>{t('settings.agents.editorHint')}</span>
          </div>
          <Space wrap>
            <Button
              icon={<RollbackOutlined />}
              onClick={handleDiscard}
              disabled={!isDirty || saving || resetting}
            >
              {t('settings.agents.actions.discard')}
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => void handleReset()}
              loading={resetting}
              disabled={saving}
            >
              {t('settings.agents.actions.reset')}
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={() => void handleSave()}
              loading={saving}
              disabled={!isDirty || resetting}
            >
              {t('settings.agents.actions.save')}
            </Button>
          </Space>
        </div>
      </>
    );
  };

  const renderEditor = () => (
    <Card className={styles.editorCard}>
      {renderEditorBody()}
    </Card>
  );

  return (
    <div className={`${styles.agentSettings} ${className || ''}`}>
      {isMobile ? (
        <>
          {!showEditorOnMobile && renderAgentList()}
          {showEditorOnMobile && (
            <div className={styles.mobileEditor}>
              <Button
                type="text"
                className={styles.mobileBackButton}
                onClick={() => setShowEditorOnMobile(false)}
              >
                {t('settings.agents.actions.backToList')}
              </Button>
              {renderEditor()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderAgentList()}</div>
          <div className={styles.editorColumn}>{renderEditor()}</div>
        </div>
      )}
      {renderCreateModal()}
    </div>
  );
};

export default AgentSettingsPage;
