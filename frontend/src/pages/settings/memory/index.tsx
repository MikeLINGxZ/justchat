import React, { useCallback, useEffect, useState } from 'react';
import {
  Input,
  Button,
  Tag,
  Empty,
  Spin,
  Popconfirm,
  Typography,
  Select,
  Pagination,
  Tooltip,
} from 'antd';
import {
  SearchOutlined,
  DeleteOutlined,
  UndoOutlined,
  DeploymentUnitOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { Events } from '@wailsio/runtime';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { useMemoryStore } from '@/stores/memoryStore';
import { useLabStore } from '@/stores/labStore';
import styles from './index.module.scss';

const { Text, Title } = Typography;

const MEMORY_TYPES = [
  { value: '', labelKey: 'settings.memory.filterAll' },
  { value: 'event', labelKey: 'settings.memory.typeEvent' },
  { value: 'skill', labelKey: 'settings.memory.typeSkill' },
  { value: 'plan', labelKey: 'settings.memory.typePlan' },
];

const MemorySettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const {
    memories, total, stats, isLoading, query,
    setQuery, fetchMemories, fetchStats, deleteMemory, restoreMemory,
  } = useMemoryStore();
  const { vectorSearchEnabled } = useLabStore();

  const [searchText, setSearchText] = useState('');

  useEffect(() => {
    fetchMemories();
    fetchStats();
  }, [fetchMemories, fetchStats]);

  useEffect(() => {
    const cancel = Events.On('settings:memories:changed', () => {
      void fetchMemories();
      void fetchStats();
    });
    return () => {
      cancel?.();
      Events.Off('settings:memories:changed');
    };
  }, [fetchMemories, fetchStats]);

  const handleSearch = useCallback(() => {
    setQuery({ keyword: searchText });
    fetchMemories();
  }, [searchText, setQuery, fetchMemories]);

  const handleTypeChange = useCallback((value: string) => {
    setQuery({ type: value });
    fetchMemories();
  }, [setQuery, fetchMemories]);

  const handleToggleForgotten = useCallback((showForgotten: boolean) => {
    setQuery({ is_forgotten: showForgotten });
    fetchMemories();
  }, [setQuery, fetchMemories]);

  const handlePageChange = useCallback((page: number) => {
    setQuery({ offset: (page - 1) * query.limit });
    fetchMemories();
  }, [query.limit, setQuery, fetchMemories]);

  const formatDate = (dateStr: string | null | undefined) => {
    if (!dateStr) return '';
    try {
      return new Date(dateStr).toLocaleDateString('zh-CN');
    } catch {
      return '';
    }
  };

  const handleEdit = useCallback((id: number) => {
    void Service.OpenEditMemoryWindow(id);
  }, []);

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title level={4}>{t('settings.memory.title')}</Title>
        <Text type="secondary">{t('settings.memory.description')}</Text>
      </div>

      {stats && (
        <div className={styles.statsBar}>
          <Tag>{t('settings.memory.statsTotal', { count: stats.total })}</Tag>
          <Tag color="blue">{t('settings.memory.statsWeekNew', { count: stats.week_new })}</Tag>
          <Tag color="orange">{t('settings.memory.statsForgotten', { count: stats.forgotten })}</Tag>
        </div>
      )}

      <div className={styles.filterBar}>
        <Input
          placeholder={t('settings.memory.searchPlaceholder')}
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          onPressEnter={handleSearch}
          allowClear
          style={{ width: 240 }}
        />
        <Select
          value={query.type}
          onChange={handleTypeChange}
          style={{ width: 120 }}
          options={MEMORY_TYPES.map(mt => ({ value: mt.value, label: t(mt.labelKey) }))}
        />
        <Button
          type={query.is_forgotten ? 'primary' : 'default'}
          onClick={() => handleToggleForgotten(!query.is_forgotten)}
          ghost={query.is_forgotten}
        >
          {t('settings.memory.showForgotten')}
        </Button>
      </div>

      <Spin spinning={isLoading}>
        {memories.length === 0 ? (
          <Empty description={t('settings.memory.empty')} />
        ) : (
          <div className={styles.memoryList}>
            {memories.map((m) => (
              <div key={m.id} className={styles.memoryCard}>
                <div className={styles.memoryHeader}>
                  <Text strong className={styles.memoryTitle}>
                    {m.summary || t('settings.memory.noTitle')}
                    {vectorSearchEnabled && m.has_embedding && (
                      <Tooltip title={t('settings.memory.embedded')}>
                        <DeploymentUnitOutlined className={styles.embeddingIcon} />
                      </Tooltip>
                    )}
                  </Text>
                  <Text type="secondary" className={styles.memoryDate}>
                    {formatDate(m.time_range_start as unknown as string) || formatDate(m.created_at as unknown as string)}
                  </Text>
                </div>
                <Text className={styles.memoryContent}>
                  {m.content && m.content.length > 150
                    ? m.content.slice(0, 150) + '...'
                    : m.content}
                </Text>
                <div className={styles.memoryMeta}>
                  {m.type && <Tag>{m.type}</Tag>}
                  {m.location && <Tag color="cyan">{m.location}</Tag>}
                  {m.characters && <Tag color="purple">{m.characters}</Tag>}
                  <Text type="secondary" className={styles.memoryScore}>
                    {t('settings.memory.trust')}: {(m.trust_score * 100).toFixed(0)}
                    %
                    {' · '}
                    {t('settings.memory.importance')}: {(m.importance * 100).toFixed(0)}
                    %
                  </Text>
                  <div className={styles.memoryActions}>
                    <Button
                      type="link"
                      size="small"
                      icon={<EditOutlined />}
                      onClick={() => handleEdit(m.id)}
                    >
                      {t('settings.memory.edit')}
                    </Button>
                    {m.is_forgotten ? (
                      <Button
                        type="link"
                        size="small"
                        icon={<UndoOutlined />}
                        onClick={() => restoreMemory(m.id)}
                      >
                        {t('settings.memory.restore')}
                      </Button>
                    ) : (
                      <Popconfirm
                        title={t('settings.memory.deleteConfirm')}
                        onConfirm={() => deleteMemory(m.id)}
                      >
                        <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                          {t('settings.memory.delete')}
                        </Button>
                      </Popconfirm>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </Spin>

      {total > query.limit && (
        <div className={styles.pagination}>
          <Pagination
            current={Math.floor(query.offset / query.limit) + 1}
            total={total}
            pageSize={query.limit}
            onChange={handlePageChange}
            showSizeChanger={false}
            size="small"
          />
        </div>
      )}

    </div>
  );
};

export default MemorySettingsPage;
