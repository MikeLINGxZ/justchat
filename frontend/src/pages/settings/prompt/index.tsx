import React, { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Empty,
  Input,
  Modal,
  Skeleton,
  Space,
  Spin,
  Tag,
  Typography,
  message,
} from 'antd';
import {
  CheckOutlined,
  FileTextOutlined,
  ReloadOutlined,
  RollbackOutlined,
  SaveOutlined,
} from '@ant-design/icons';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { PromptFileDetail, PromptFileSummary } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Paragraph, Text, Title } = Typography;
const { TextArea } = Input;

interface PromptSettingsPageProps {
  className?: string;
}

const PromptSettingsPage: React.FC<PromptSettingsPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<PromptFileSummary[]>([]);
  const [activeName, setActiveName] = useState<string>('');
  const [detail, setDetail] = useState<PromptFileDetail | null>(null);
  const [draft, setDraft] = useState('');
  const [loadingList, setLoadingList] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [saving, setSaving] = useState(false);
  const [resetting, setResetting] = useState(false);
  const [listError, setListError] = useState('');
  const [detailError, setDetailError] = useState('');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showEditorOnMobile, setShowEditorOnMobile] = useState(false);

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
      const result = await Service.ListPromptFiles();
      setItems(result);

      const nextActiveName = preferredName && result.some(item => item.name === preferredName)
        ? preferredName
        : result[0]?.name || '';
      if (!activeName && nextActiveName) {
        setActiveName(nextActiveName);
      }
      return result;
    } catch (error) {
      console.error('加载提示词列表失败:', error);
      setListError(t('settings.prompt.listLoadFailedDesc'));
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
      const result = await Service.GetPromptFile(name);
      setDetail(result);
      setDraft(result?.content || '');
      setActiveName(name);
      if (isMobile) {
        setShowEditorOnMobile(true);
      }
    } catch (error) {
      console.error('加载提示词失败:', error);
      setDetail(null);
      setDraft('');
      setDetailError(t('settings.prompt.detailLoadFailedDesc'));
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

  const syncItem = (nextDetail: PromptFileDetail) => {
    setItems(prev =>
      prev.map(item =>
        item.name === nextDetail.name
          ? {
              ...item,
              title: nextDetail.title,
              description: nextDetail.description,
              is_system: nextDetail.is_system,
              updated_at: nextDetail.updated_at,
            }
          : item,
      ),
    );
  };

  const confirmDiscardIfNeeded = async (onConfirm: () => void) => {
    if (!isDirty) {
      onConfirm();
      return;
    }

    Modal.confirm({
      title: t('settings.prompt.unsavedConfirmTitle'),
      content: t('settings.prompt.unsavedConfirmContent'),
      okText: t('settings.prompt.discardChanges'),
      cancelText: t('settings.prompt.continueEditing'),
      okButtonProps: { danger: true },
      onOk: () => {
        onConfirm();
      },
    });
  };

  const handleSelectPrompt = async (name: string) => {
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
    message.success(t('settings.prompt.discardSuccess'));
  };

  const handleSave = async () => {
    if (!detail) {
      return;
    }
    const content = draft.trim();
    if (!content) {
      message.error(t('settings.prompt.contentRequired'));
      return;
    }

    setSaving(true);
    try {
      const result = await Service.UpdatePromptFile(detail.name, content);
      setDetail(result);
      setDraft(result?.content || '');
      if (result) {
        syncItem(result);
      }
      message.success(t('settings.prompt.saveSuccess'));
    } catch (error) {
      console.error('保存提示词失败:', error);
      message.error(t('settings.prompt.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleReset = async () => {
    if (!detail) {
      return;
    }
    setResetting(true);
    try {
      const result = await Service.ResetPromptFile(detail.name);
      setDetail(result);
      setDraft(result?.content || '');
      if (result) {
        syncItem(result);
      }
      message.success(t('settings.prompt.resetSuccess'));
    } catch (error) {
      console.error('恢复默认提示词失败:', error);
      message.error(t('settings.prompt.resetFailed'));
    } finally {
      setResetting(false);
    }
  };

  const renderPromptList = () => (
    <Card className={styles.listCard} title={t('settings.prompt.listTitle')}>
      {loadingList ? (
        <div className={styles.listLoading}>
          <Skeleton active paragraph={{ rows: 6 }} />
        </div>
      ) : listError ? (
        <Alert
          type="error"
          showIcon
          message={t('settings.prompt.listLoadFailed')}
          description={listError}
          action={
            <Button size="small" onClick={() => void refreshList(activeName || undefined)}>
              {t('common.retry')}
            </Button>
          }
        />
      ) : (
        <div className={styles.promptList}>
          {items.map(item => {
            const selected = item.name === activeName;
            return (
              <button
                key={item.name}
                type="button"
                className={`${styles.promptItem} ${selected ? styles.selected : ''}`}
                onClick={() => void handleSelectPrompt(item.name)}
              >
                <div className={styles.promptItemHeader}>
                  <span className={styles.promptItemTitle}>{item.title}</span>
                  <Tag color={item.is_system ? 'blue' : 'gold'} bordered={false}>
                    {item.is_system ? t('settings.prompt.tags.system') : t('settings.prompt.tags.user')}
                  </Tag>
                </div>
                <div className={styles.promptItemFile} title={item.name}>{item.name}</div>
                <div className={styles.promptItemDesc}>{item.description}</div>
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
          message={t('settings.prompt.detailLoadFailed')}
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
          <Empty description={t('settings.prompt.empty')} />
        </div>
      );
    }

    return (
      <>
        <div className={styles.editorHeader}>
          <div>
            <div className={styles.editorTitleRow}>
              <Title level={4}>{detail.title}</Title>
              <Tag color={detail.is_system ? 'blue' : 'gold'} bordered={false}>
                {detail.is_system ? t('settings.prompt.tags.system') : t('settings.prompt.tags.user')}
              </Tag>
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
            <Text className={styles.fileName} title={detail.name}>{detail.name}</Text>
            <Paragraph className={styles.editorDescription}>{detail.description}</Paragraph>
          </div>
          <Alert
            type="info"
            showIcon
            className={styles.tipAlert}
            message={t('settings.prompt.tip')}
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
            <span>{t('settings.prompt.editorHint')}</span>
          </div>
          <Space wrap>
            <Button
              icon={<RollbackOutlined />}
              onClick={handleDiscard}
              disabled={!isDirty || saving || resetting}
            >
              {t('settings.prompt.actions.discard')}
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => void handleReset()}
              loading={resetting}
              disabled={saving}
            >
              {t('settings.prompt.actions.reset')}
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={() => void handleSave()}
              loading={saving}
              disabled={!isDirty || resetting}
            >
              {t('settings.prompt.actions.save')}
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
    <div className={`${styles.promptSettings} ${className || ''}`}>
      {isMobile ? (
        <>
          {!showEditorOnMobile && renderPromptList()}
          {showEditorOnMobile && (
            <div className={styles.mobileEditor}>
              <Button
                type="text"
              className={styles.mobileBackButton}
              onClick={() => setShowEditorOnMobile(false)}
            >
              {t('settings.prompt.actions.backToList')}
            </Button>
              {renderEditor()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderPromptList()}</div>
          <div className={styles.editorColumn}>{renderEditor()}</div>
        </div>
      )}
    </div>
  );
};

export default PromptSettingsPage;
