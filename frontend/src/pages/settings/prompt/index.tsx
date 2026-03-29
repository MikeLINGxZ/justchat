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
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { PromptFileDetail, PromptFileSummary } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Paragraph, Text, Title } = Typography;
const { TextArea } = Input;

interface PromptSettingsPageProps {
  className?: string;
}

const PromptSettingsPage: React.FC<PromptSettingsPageProps> = ({ className }) => {
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
      setListError('加载提示词列表失败，请稍后重试。');
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
      setDetailError('读取提示词失败，请重试。');
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
      title: '放弃未保存的修改？',
      content: '当前提示词还有未保存内容，切换后这些修改会丢失。',
      okText: '放弃修改',
      cancelText: '继续编辑',
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
    message.success('已放弃未保存修改');
  };

  const handleSave = async () => {
    if (!detail) {
      return;
    }
    const content = draft.trim();
    if (!content) {
      message.error('提示词内容不能为空');
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
      message.success('提示词已保存并应用');
    } catch (error) {
      console.error('保存提示词失败:', error);
      message.error('保存失败，请稍后重试');
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
      message.success('已恢复默认提示词并立即应用');
    } catch (error) {
      console.error('恢复默认提示词失败:', error);
      message.error('恢复默认失败，请稍后重试');
    } finally {
      setResetting(false);
    }
  };

  const renderPromptList = () => (
    <Card className={styles.listCard} title="提示词文件">
      {loadingList ? (
        <div className={styles.listLoading}>
          <Skeleton active paragraph={{ rows: 6 }} />
        </div>
      ) : listError ? (
        <Alert
          type="error"
          showIcon
          message="列表加载失败"
          description={listError}
          action={
            <Button size="small" onClick={() => void refreshList(activeName || undefined)}>
              重试
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
                    {item.is_system ? 'SYSTEM' : 'USER'}
                  </Tag>
                </div>
                <div className={styles.promptItemFile}>{item.name}</div>
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
          message="读取提示词失败"
          description={detailError}
          action={
            <Button size="small" onClick={() => activeName && void loadDetail(activeName)}>
              重试
            </Button>
          }
        />
      );
    }

    if (!detail) {
      return (
        <div className={styles.emptyState}>
          <Empty description="请选择一个提示词文件开始编辑" />
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
                {detail.is_system ? 'SYSTEM' : 'USER'}
              </Tag>
              {isDirty ? (
                <Tag color="orange" bordered={false}>
                  未保存
                </Tag>
              ) : (
                <Tag color="green" bordered={false}>
                  已同步
                </Tag>
              )}
            </div>
            <Text className={styles.fileName}>{detail.name}</Text>
            <Paragraph className={styles.editorDescription}>{detail.description}</Paragraph>
          </div>
          <Alert
            type="info"
            showIcon
            className={styles.tipAlert}
            message="保存后会立即刷新内存中的提示词缓存，仅影响之后新发起的请求。"
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
            <span>支持直接编辑 markdown / 模板变量文本。</span>
          </div>
          <Space wrap>
            <Button
              icon={<RollbackOutlined />}
              onClick={handleDiscard}
              disabled={!isDirty || saving || resetting}
            >
              放弃修改
            </Button>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => void handleReset()}
              loading={resetting}
              disabled={saving}
            >
              恢复默认
            </Button>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={() => void handleSave()}
              loading={saving}
              disabled={!isDirty || resetting}
            >
              保存并应用
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
                返回提示词列表
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
