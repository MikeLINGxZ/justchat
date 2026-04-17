import React, { useEffect, useState } from 'react';
import {
  Alert,
  Button,
  Card,
  Empty,
  Popconfirm,
  Skeleton,
  Switch,
  Tag,
  Typography,
  message,
} from 'antd';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { PluginSummary } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Paragraph, Text, Title } = Typography;

interface PluginSettingsPageProps {
  className?: string;
}

const PluginSettingsPage: React.FC<PluginSettingsPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [plugins, setPlugins] = useState<PluginSummary[]>([]);
  const [selectedPluginId, setSelectedPluginId] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [listError, setListError] = useState('');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showDetailOnMobile, setShowDetailOnMobile] = useState(false);

  useEffect(() => {
    const handleResize = () => setIsMobile(isMobileDevice());
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const selectedPlugin = plugins.find(p => p.id === selectedPluginId) || null;

  const loadPlugins = async () => {
    setLoading(true);
    setListError('');
    try {
      const result = await Service.GetInstalledPlugins();
      const list = result || [];
      setPlugins(list);
      if (!selectedPluginId && list.length > 0) {
        setSelectedPluginId(list[0].id);
      }
    } catch (error) {
      console.error('Failed to load plugins:', error);
      setListError(t('settings.plugins.loadFailed', 'Failed to load plugins'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadPlugins();
  }, []);

  const handleInstall = async () => {
    try {
      await Service.InstallPlugin();
      message.success(t('settings.plugins.installSuccess'));
      await loadPlugins();
    } catch (error) {
      console.error('Failed to install plugin:', error);
      message.error(t('settings.plugins.installFailed'));
    }
  };

  const handleToggle = async (plugin: PluginSummary, enabled: boolean) => {
    try {
      if (enabled) {
        await Service.EnablePlugin(plugin.id);
      } else {
        await Service.DisablePlugin(plugin.id);
      }
      await loadPlugins();
    } catch (error) {
      console.error('Failed to toggle plugin:', error);
    }
  };

  const handleDelete = async (pluginId: string) => {
    try {
      await Service.UninstallPlugin(pluginId);
      message.success(t('settings.plugins.deleteSuccess'));
      if (selectedPluginId === pluginId) {
        setSelectedPluginId('');
      }
      await loadPlugins();
    } catch (error) {
      console.error('Failed to delete plugin:', error);
    }
  };

  const handleSelectPlugin = (id: string) => {
    setSelectedPluginId(id);
    if (isMobile) setShowDetailOnMobile(true);
  };

  const getStateTag = (state: string) => {
    switch (state) {
      case 'active':
        return <Tag color="success" bordered={false}>{t('settings.plugins.state.active')}</Tag>;
      case 'error':
        return <Tag color="error" bordered={false}>{t('settings.plugins.state.error')}</Tag>;
      case 'installed':
      default:
        return <Tag color="default" bordered={false}>{t('settings.plugins.state.installed')}</Tag>;
    }
  };

  const renderPluginList = () => (
    <Card
      className={styles.listCard}
      title={
        <div className={styles.listTitleRow}>
          <span>{t('settings.plugins.title')}</span>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            onClick={() => void handleInstall()}
            className={styles.addButton}
          />
        </div>
      }
    >
      {loading ? (
        <div className={styles.listLoading}>
          <Skeleton active paragraph={{ rows: 6 }} />
        </div>
      ) : listError ? (
        <Alert
          type="error"
          showIcon
          message={t('settings.plugins.loadFailed', 'Load failed')}
          description={listError}
          action={
            <Button size="small" onClick={() => void loadPlugins()}>
              {t('common.retry', 'Retry')}
            </Button>
          }
        />
      ) : plugins.length === 0 ? (
        <div className={styles.emptyState}>
          <Empty description={t('settings.plugins.empty')}>
            <Text type="secondary">{t('settings.plugins.emptyDesc')}</Text>
          </Empty>
        </div>
      ) : (
        <div className={styles.pluginList}>
          {plugins.map(plugin => {
            const selected = plugin.id === selectedPluginId;
            return (
              <button
                key={plugin.id}
                type="button"
                className={`${styles.pluginItem} ${selected ? styles.selected : ''}`}
                onClick={() => handleSelectPlugin(plugin.id)}
              >
                <div className={styles.pluginItemHeader}>
                  <span className={styles.pluginItemTitle}>
                    {plugin.display_name}
                  </span>
                  <div className={styles.pluginItemTags}>
                    <Tag bordered={false}>{plugin.version}</Tag>
                    {getStateTag(plugin.state)}
                  </div>
                </div>
                <div className={styles.pluginItemName}>{plugin.id}</div>
              </button>
            );
          })}
        </div>
      )}
    </Card>
  );

  const renderDetailBody = () => {
    if (!selectedPlugin) {
      return (
        <div className={styles.emptyState}>
          <Empty description={t('settings.plugins.selectPlugin', 'Select a plugin to view details')} />
        </div>
      );
    }

    return (
      <>
        <div className={styles.detailHeader}>
          <div className={styles.detailTitleRow}>
            <Title level={4}>{selectedPlugin.display_name}</Title>
            <Tag bordered={false}>{selectedPlugin.version}</Tag>
            {getStateTag(selectedPlugin.state)}
          </div>
          <Text className={styles.pluginId}>{selectedPlugin.id}</Text>
          {selectedPlugin.description && (
            <Paragraph className={styles.detailDescription}>
              {selectedPlugin.description}
            </Paragraph>
          )}
        </div>

        <div className={styles.contribSections}>
          {selectedPlugin.tools && selectedPlugin.tools.length > 0 && (
            <div className={styles.detailSection}>
              <Text className={styles.sectionTitle}>
                {t('settings.plugins.tools')} ({selectedPlugin.tools.length})
              </Text>
              <div className={styles.contribList}>
                {selectedPlugin.tools.map(tool => (
                  <Tag key={tool.id} bordered={false} color="blue">{tool.name || tool.id}</Tag>
                ))}
              </div>
            </div>
          )}

          {selectedPlugin.agents && selectedPlugin.agents.length > 0 && (
            <div className={styles.detailSection}>
              <Text className={styles.sectionTitle}>
                {t('settings.plugins.agents')} ({selectedPlugin.agents.length})
              </Text>
              <div className={styles.contribList}>
                {selectedPlugin.agents.map(agent => (
                  <Tag key={agent.id} bordered={false} color="green">{agent.name || agent.id}</Tag>
                ))}
              </div>
            </div>
          )}

          {selectedPlugin.hooks && selectedPlugin.hooks.length > 0 && (
            <div className={styles.detailSection}>
              <Text className={styles.sectionTitle}>
                {t('settings.plugins.hooks')} ({selectedPlugin.hooks.length})
              </Text>
              <div className={styles.contribList}>
                {selectedPlugin.hooks.map(hook => (
                  <Tag key={hook} bordered={false} color="orange">{hook}</Tag>
                ))}
              </div>
            </div>
          )}

          {selectedPlugin.view_count > 0 && (
            <div className={styles.detailSection}>
              <Text className={styles.sectionTitle}>
                {t('settings.plugins.views')} ({selectedPlugin.view_count})
              </Text>
            </div>
          )}
        </div>

        <div className={styles.detailActions}>
          <div className={styles.actionRow}>
            <Text>{t('settings.plugins.enabled', 'Enabled')}</Text>
            <Switch
              checked={selectedPlugin.enabled}
              onChange={(checked) => void handleToggle(selectedPlugin, checked)}
            />
          </div>
          <Popconfirm
            title={t('settings.plugins.confirmDelete')}
            onConfirm={() => void handleDelete(selectedPlugin.id)}
            okButtonProps={{ danger: true }}
          >
            <Button
              danger
              icon={<DeleteOutlined />}
            >
              {t('common.delete', 'Delete')}
            </Button>
          </Popconfirm>
        </div>
      </>
    );
  };

  const renderDetail = () => (
    <Card className={styles.detailCard}>
      {renderDetailBody()}
    </Card>
  );

  return (
    <div className={`${styles.pluginSettings} ${className || ''}`}>
      {isMobile ? (
        <>
          {!showDetailOnMobile && renderPluginList()}
          {showDetailOnMobile && (
            <div className={styles.mobileDetail}>
              <Button
                type="text"
                className={styles.mobileBackButton}
                onClick={() => setShowDetailOnMobile(false)}
              >
                {t('settings.plugins.backToList', 'Back to list')}
              </Button>
              {renderDetail()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderPluginList()}</div>
          <div className={styles.detailColumn}>{renderDetail()}</div>
        </div>
      )}
    </div>
  );
};

export default PluginSettingsPage;
