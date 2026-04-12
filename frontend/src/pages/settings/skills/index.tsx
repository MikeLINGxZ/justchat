import React, { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Dropdown,
  Empty,
  Form,
  Input,
  Modal,
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
  FolderOpenOutlined,
  PlusOutlined,
  RollbackOutlined,
  SaveOutlined,
} from '@ant-design/icons';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import {
  SkillDetail,
  SkillSummary,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Paragraph, Text, Title } = Typography;
const { TextArea } = Input;

interface SkillSettingsPageProps {
  className?: string;
}

const SkillSettingsPage: React.FC<SkillSettingsPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<SkillSummary[]>([]);
  const [activeName, setActiveName] = useState<string>('');
  const [detail, setDetail] = useState<SkillDetail | null>(null);
  const [draft, setDraft] = useState('');
  const [loadingList, setLoadingList] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [listError, setListError] = useState('');
  const [detailError, setDetailError] = useState('');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showEditorOnMobile, setShowEditorOnMobile] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [createForm] = Form.useForm();

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(isMobileDevice());
    };
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const isDirty = useMemo(() => {
    return detail !== null && draft !== detail.content;
  }, [detail, draft]);

  const refreshList = async (preferredName?: string) => {
    setLoadingList(true);
    setListError('');
    try {
      const result = await Service.ListSkills();
      setItems(result);

      const nextActiveName =
        preferredName && result.some(item => item.name === preferredName)
          ? preferredName
          : result[0]?.name || '';
      if (!activeName && nextActiveName) {
        setActiveName(nextActiveName);
      }
      return result;
    } catch (error) {
      console.error('加载技能列表失败:', error);
      setListError(t('settings.skills.listLoadFailedDesc'));
      return [];
    } finally {
      setLoadingList(false);
    }
  };

  const loadDetail = async (name: string) => {
    if (!name) {
      return;
    }
    setLoadingDetail(true);
    setDetailError('');
    try {
      const result = await Service.GetSkill(name);
      setDetail(result);
      setDraft(result?.content || '');
      setActiveName(name);
      if (isMobile) {
        setShowEditorOnMobile(true);
      }
    } catch (error) {
      console.error('加载技能失败:', error);
      setDetail(null);
      setDraft('');
      setDetailError(t('settings.skills.detailLoadFailedDesc'));
    } finally {
      setLoadingDetail(false);
    }
  };

  useEffect(() => {
    void refreshList();
  }, []);

  useEffect(() => {
    if (activeName) {
      void loadDetail(activeName);
    }
  }, [activeName]);

  const confirmDiscardIfNeeded = async (onConfirm: () => void) => {
    if (!isDirty) {
      onConfirm();
      return;
    }

    Modal.confirm({
      title: t('settings.skills.unsavedConfirmTitle'),
      content: t('settings.skills.unsavedConfirmContent'),
      okText: t('settings.skills.discardChanges'),
      cancelText: t('settings.skills.continueEditing'),
      okButtonProps: { danger: true },
      onOk: () => {
        onConfirm();
      },
    });
  };

  const handleSelectSkill = async (name: string) => {
    if (name === activeName && detail) {
      if (isMobile) {
        setShowEditorOnMobile(true);
      }
      return;
    }
    await confirmDiscardIfNeeded(() => {
      setActiveName(name);
    });
  };

  const handleDiscard = () => {
    if (!detail) {
      return;
    }
    setDraft(detail.content);
    message.success(t('settings.skills.discardSuccess'));
  };

  const handleSave = async () => {
    if (!detail) {
      return;
    }
    const content = draft.trim();
    if (!content) {
      message.error(t('settings.skills.contentRequired'));
      return;
    }

    setSaving(true);
    try {
      const input = new SkillDetail({
        name: detail.name,
        description: detail.description,
        when: detail.when || '',
        version: detail.version,
        tags: detail.tags,
        content,
      });
      const result = await Service.UpdateSkill(detail.name, input);
      setDetail(result);
      setDraft(result?.content || '');
      if (result) {
        setItems(prev =>
          prev.map(item =>
            item.name === result.name
              ? {
                  ...item,
                  description: result.description,
                  version: result.version,
                  tags: result.tags,
                }
              : item,
          ),
        );
      }
      message.success(t('settings.skills.saveSuccess'));
    } catch (error) {
      console.error('保存技能失败:', error);
      message.error(t('settings.skills.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = () => {
    if (!detail) {
      return;
    }
    Modal.confirm({
      title: t('settings.skills.deleteConfirmTitle'),
      content: t('settings.skills.deleteConfirmContent'),
      okText: t('settings.skills.actions.delete'),
      okButtonProps: { danger: true },
      onOk: async () => {
        setDeleting(true);
        try {
          await Service.DeleteSkill(detail.name);
          message.success(t('settings.skills.deleteSuccess'));
          setDetail(null);
          setDraft('');
          setActiveName('');
          await refreshList();
        } catch (error) {
          console.error('删除技能失败:', error);
          message.error(t('settings.skills.deleteFailed'));
        } finally {
          setDeleting(false);
        }
      },
    });
  };

  const handleCreate = async () => {
    try {
      const values = await createForm.validateFields();
      setCreating(true);
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
        message.success(t('settings.skills.createSuccess'));
        setCreateModalOpen(false);
        createForm.resetFields();
        setActiveName(result.name);
        await refreshList(result.name);
      }
    } catch (error) {
      if ((error as { errorFields?: unknown }).errorFields) {
        return;
      }
      console.error('创建技能失败:', error);
      message.error(t('settings.skills.createFailed'));
    } finally {
      setCreating(false);
    }
  };

  const handleImportFromFolder = async () => {
    try {
      const folderPath = await Service.SelectSkillFolder();
      if (!folderPath) return;

      const imported = await Service.ImportSkillsFromFolder(folderPath);
      if (imported && imported.length > 0) {
        message.success(t('settings.skills.importSuccess', { count: imported.length }));
        await refreshList(imported[0].name);
        setActiveName(imported[0].name);
      } else {
        message.info(t('settings.skills.importEmpty'));
      }
    } catch (error) {
      console.error('导入技能失败:', error);
      message.error(t('settings.skills.createFailed'));
    }
  };

  const renderSkillList = () => (
    <Card
      className={styles.listCard}
      title={
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <span>{t('settings.skills.listTitle')}</span>
          <Dropdown
            menu={{
              items: [
                {
                  key: 'create',
                  icon: <PlusOutlined />,
                  label: t('settings.skills.actions.create'),
                  onClick: () => setCreateModalOpen(true),
                },
                {
                  key: 'import',
                  icon: <FolderOpenOutlined />,
                  label: t('settings.skills.importFromFolder'),
                  onClick: () => void handleImportFromFolder(),
                },
              ],
            }}
            trigger={['click']}
          >
            <Button type="text" size="small" icon={<PlusOutlined />} />
          </Dropdown>
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
          message={t('settings.skills.listLoadFailed')}
          description={listError}
          action={
            <Button size="small" onClick={() => void refreshList(activeName || undefined)}>
              {t('common.retry')}
            </Button>
          }
        />
      ) : (
        <div className={styles.skillList}>
          {items.map(item => {
            const selected = item.name === activeName;
            return (
              <button
                key={item.name}
                type="button"
                className={`${styles.skillItem} ${selected ? styles.selected : ''}`}
                onClick={() => void handleSelectSkill(item.name)}
              >
                <div className={styles.skillItemHeader}>
                  <span className={styles.skillItemTitle}>{item.name}</span>
                  <Tag color="default" bordered={false}>
                    {item.version}
                  </Tag>
                </div>
                <div className={styles.skillItemDesc}>{item.description}</div>
                {item.tags && item.tags.length > 0 && (
                  <div className={styles.skillItemTags}>
                    {item.tags.map(tag => (
                      <Tag key={tag} bordered={false}>
                        {tag}
                      </Tag>
                    ))}
                  </div>
                )}
              </button>
            );
          })}
        </div>
      )}
    </Card>
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
          message={t('settings.skills.detailLoadFailed')}
          description={detailError}
          action={
            <Button size="small" onClick={() => activeName && void loadDetail(activeName)}>
              {t('common.retry')}
            </Button>
          }
        />
      );
    }

    if (!detail) {
      return (
        <div className={styles.emptyState}>
          <Empty description={t('settings.skills.empty')} />
        </div>
      );
    }

    return (
      <>
        <div className={styles.editorHeader}>
          <div>
            <div className={styles.editorTitleRow}>
              <Title level={4}>{detail.name}</Title>
              {isDirty ? (
                <Tag color="orange" bordered={false}>
                  {t('settings.prompt.tags.unsaved')}
                </Tag>
              ) : (
                <Tag color="green" bordered={false}>
                  {t('settings.prompt.tags.synced')}
                </Tag>
              )}
            </div>
            <Paragraph className={styles.editorDescription}>{detail.description}</Paragraph>
            {detail.when && (
              <Paragraph type="secondary" style={{ marginBottom: 4 }}>
                {t('settings.skills.form.when')}: {detail.when}
              </Paragraph>
            )}
            <Text type="secondary">{detail.version}</Text>
            {detail.tags && detail.tags.length > 0 && (
              <div className={styles.skillItemTags} style={{ marginTop: 8 }}>
                {detail.tags.map(tag => (
                  <Tag key={tag} bordered={false}>
                    {tag}
                  </Tag>
                ))}
              </div>
            )}
          </div>
          <Alert
            type="info"
            showIcon
            className={styles.tipAlert}
            message={t('settings.skills.tip')}
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
            <span>{t('settings.skills.editorHint')}</span>
          </div>
          <Space wrap>
            <Button
              icon={<RollbackOutlined />}
              onClick={handleDiscard}
              disabled={!isDirty || saving || deleting}
            >
              {t('settings.skills.actions.discard')}
            </Button>
            <Button
              danger
              icon={<DeleteOutlined />}
              onClick={handleDelete}
              loading={deleting}
              disabled={saving}
            >
              {t('settings.skills.actions.delete')}
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={() => void handleSave()}
              loading={saving}
              disabled={!isDirty || deleting}
            >
              {t('settings.skills.actions.save')}
            </Button>
          </Space>
        </div>
      </>
    );
  };

  const renderEditor = () => (
    <Card className={styles.editorCard}>{renderEditorBody()}</Card>
  );

  const renderCreateModal = () => (
    <Modal
      title={t('settings.skills.createTitle')}
      open={createModalOpen}
      onCancel={() => {
        setCreateModalOpen(false);
        createForm.resetFields();
      }}
      onOk={() => void handleCreate()}
      confirmLoading={creating}
      okText={t('settings.skills.actions.create')}
      destroyOnClose
    >
      <Form form={createForm} layout="vertical" preserve={false}>
        <Form.Item
          name="name"
          label={t('settings.skills.form.name')}
          rules={[
            { required: true },
            {
              pattern: /^[a-zA-Z0-9_-]+$/,
              message: t('settings.skills.form.nameInvalid'),
            },
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
        <Form.Item
          name="when"
          label={t('settings.skills.form.when')}
        >
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
          <TextArea rows={6} placeholder={t('settings.skills.form.contentPlaceholder')} />
        </Form.Item>
      </Form>
    </Modal>
  );

  return (
    <div className={`${styles.skillSettings} ${className || ''}`}>
      {isMobile ? (
        <>
          {!showEditorOnMobile && renderSkillList()}
          {showEditorOnMobile && (
            <div className={styles.mobileEditor}>
              <Button
                type="text"
                className={styles.mobileBackButton}
                onClick={() => setShowEditorOnMobile(false)}
              >
                {t('settings.skills.actions.backToList')}
              </Button>
              {renderEditor()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderSkillList()}</div>
          <div className={styles.editorColumn}>{renderEditor()}</div>
        </div>
      )}
      {renderCreateModal()}
    </div>
  );
};

export default SkillSettingsPage;
