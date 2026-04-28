import React, { useState, useEffect } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Switch,
  Space,
  Alert,
  message,
  Popconfirm,
  Avatar,
  Tooltip,
  Dropdown,
  Tag,
  Empty,
  Modal,
} from 'antd';
import type { MenuProps } from 'antd';
import {
  ApiOutlined,
  SaveOutlined,
  ReloadOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  PlusOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  ArrowLeftOutlined,
  QuestionCircleOutlined,
  CloseOutlined,
  EllipsisOutlined,
  StarOutlined,
} from '@ant-design/icons';
import { Events } from '@wailsio/runtime';
import { useModelStore } from '@/stores/modelStore';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { Model, Provider, SupportProvider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import { ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models';
import { setDefaultModelConfig } from '@/utils/defaultModel';
import styles from './index.module.scss';

const toAssetUrl = (path: string) => {
  const base = import.meta.env.BASE_URL || '/';
  const normalizedBase = base.endsWith('/') ? base : `${base}/`;
  return `${normalizedBase}${path.replace(/^\/+/, '')}`;
};

const detectProviderTypeFromName = (providerName?: string): ProviderType | undefined => {
  const normalizedName = providerName?.trim().toLowerCase();
  if (!normalizedName) {
    return undefined;
  }

  if (normalizedName.includes('deepseek') || normalizedName.includes('深度求索')) {
    return ProviderType.ProviderTypeDeepseek;
  }
  if (
    normalizedName.includes('aliyun') ||
    normalizedName.includes('alibaba cloud') ||
    normalizedName.includes('bailian') ||
    normalizedName.includes('阿里云') ||
    normalizedName.includes('百炼') ||
    normalizedName.includes('qwen')
  ) {
    return ProviderType.ProviderTypeAliyuns;
  }
  if (normalizedName.includes('openrouter')) {
    return ProviderType.ProviderTypeOpenrouter;
  }
  if (normalizedName.includes('ollama')) {
    return ProviderType.ProviderTypeOllama;
  }
  if (
    normalizedName.includes('openai-compatible') ||
    normalizedName.includes('openai compatible api') ||
    normalizedName.includes('openai compatible') ||
    normalizedName.includes('openai 标准接口') ||
    normalizedName.includes('openai-compatible api')
  ) {
    return ProviderType.ProviderTypeOther;
  }

  return undefined;
};

const detectProviderTypeFromBaseUrl = (baseUrl?: string): ProviderType | undefined => {
  if (!baseUrl) {
    return undefined;
  }

  try {
    const hostname = new URL(baseUrl).hostname.toLowerCase();
    if (hostname.includes('deepseek')) {
      return ProviderType.ProviderTypeDeepseek;
    }
    if (hostname.includes('dashscope') || hostname.includes('aliyuncs')) {
      return ProviderType.ProviderTypeAliyuns;
    }
    if (hostname.includes('openrouter')) {
      return ProviderType.ProviderTypeOpenrouter;
    }
    if (hostname === 'localhost' || hostname === '127.0.0.1') {
      return ProviderType.ProviderTypeOllama;
    }
  } catch {
    return undefined;
  }

  return undefined;
};

interface ProviderSettingPageProps {
  className?: string;
}

interface ProviderConfig {
  id: number;
  provider_name: string;
  provider_type?: ProviderType;
  api_key: string;
  base_url: string;
  file_upload_base_url?: string | null;
  enable: boolean;
  is_default?: boolean;
  default_model_id: number | null;
  models: Model[];
  icon?: string;
  description?: string;
  status?: 'connected' | 'disconnected' | 'testing';
}

const findDefaultModelForProvider = (provider?: ProviderConfig | null) => {
  if (!provider?.default_model_id) {
    return null;
  }
  return provider.models.find(model => model.id === provider.default_model_id) || null;
};

const ProviderSettingPage: React.FC<ProviderSettingPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [testingConnection, setTestingConnection] = useState(false);
  const [providers, setProviders] = useState<ProviderConfig[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<number | null>(null); // 改为number类型
  const [isCreatingNew, setIsCreatingNew] = useState(false); // 新增：是否正在创建新供应商
  const [newProviderTempId, setNewProviderTempId] = useState<number | null>(null); // 新增：临时ID
  const [supportProviders, setSupportProviders] = useState<SupportProvider[]>([]); // 支持的供应商列表
  const [loadingSupportProviders, setLoadingSupportProviders] = useState(false); // 加载支持的供应商状态
  const [customModelName, setCustomModelName] = useState(''); // 自定义模型名称输入
  const [addingCustomModel, setAddingCustomModel] = useState(false); // 添加自定义模型loading
  const [showCustomModelInput, setShowCustomModelInput] = useState(false);
  const [settingDefaultProviderId, setSettingDefaultProviderId] = useState<number | null>(null);
  const [settingDefaultModelId, setSettingDefaultModelId] = useState<number | null>(null);

  const { refetch: refetchModels } = useModelStore();
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showEditorOnMobile, setShowEditorOnMobile] = useState(false);

  // 监听窗口大小变化，更新移动端状态
  useEffect(() => {
    const handleResize = () => {
      setIsMobile(isMobileDevice());
    };
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // 加载保存的配置
  useEffect(() => {
    loadSupportProviders(); // 先加载支持的供应商列表
  }, []);

  // 当 supportProviders 加载完成后，加载供应商配置
  useEffect(() => {
    if (supportProviders.length > 0) {
      loadProviderConfigs();
    }
  }, [supportProviders]);

  // 页面关闭前刷新模型数据
  useEffect(() => {
    return () => {
      // 组件卸载时刷新模型数据
      refetchModels();
    };
  }, [refetchModels]);

  // 当选中的供应商变化时，更新表单
  useEffect(() => {
    const provider = providers.find(p => p.id === selectedProvider);
    if (provider && !isCreatingNew) { // 在创建新供应商时不自动更新表单
      form.setFieldsValue({
        enabled: provider.enable,
        apiKey: provider.api_key,
        baseUrl: provider.base_url,
        fileUploadBaseUrl: provider.file_upload_base_url || '',
        providerName: provider.provider_name,
      });
      setCustomModelName('');
      setShowCustomModelInput(false);
    }
  }, [selectedProvider, providers, form, isCreatingNew]);

  const resolveIconSrc = (icon?: string) => {
    if (!icon) {
      return undefined;
    }
    if (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('data:')) {
      return icon;
    }
    if (icon.startsWith('/')) {
      return toAssetUrl(icon);
    }
    if (/\.(png|jpe?g|gif|webp|svg|ico)$/i.test(icon)) {
      return toAssetUrl(icon);
    }
    return `data:image/png;base64,${icon}`;
  };

  const findSupportProvider = (providerType?: ProviderType, providerName?: string, baseUrl?: string) => {
    const resolvedProviderType =
      providerType ||
      detectProviderTypeFromName(providerName) ||
      detectProviderTypeFromBaseUrl(baseUrl);
    if (resolvedProviderType) {
      const matchedByType = supportProviders.find(item => item.provider_type === resolvedProviderType);
      if (matchedByType) {
        return matchedByType;
      }
    }
    if (providerName) {
      return supportProviders.find(item => item.name === providerName);
    }
    return undefined;
  };

  const syncStoredDefaultModel = (providerList: ProviderConfig[]) => {
    const defaultProvider = providerList.find(provider => provider.is_default);
    const defaultModel = findDefaultModelForProvider(defaultProvider);
    if (!defaultProvider || !defaultModel) {
      return;
    }
    setDefaultModelConfig({ modelId: defaultModel.id, modelName: defaultModel.model });
  };

  const loadProviderConfigs = async () => {
    setLoading(true);
    try {
      const providers = await Service.GetProviders();
      if (providers && providers.length > 0) {
        // 转换后端数据格式，添加前端需要的字段
        const formattedProviders = providers.map(provider => {
          const matchedSupportProvider = findSupportProvider(
            provider.provider_type,
            provider.provider_name,
            provider.base_url,
          );
          const extras = getProviderExtras(provider.provider_name, provider.provider_type, provider.base_url);
          const icon = matchedSupportProvider?.icon || extras.icon;
          
          return {
            ...provider,
            icon: icon,
            description: extras.description,
            status: 'disconnected' as const
          };
        });
        setProviders(formattedProviders);
        syncStoredDefaultModel(formattedProviders);
        
        // 设置默认选中第一个供应商
        if (selectedProvider === null && formattedProviders.length > 0) {
          setSelectedProvider(formattedProviders[0].id);
        }
        return formattedProviders;
      }
      setProviders([]);
      return [];
    } catch (error) {
      console.error('加载供应商配置失败:', error);
      message.error(t('settings.provider.messages.loadProvidersFailed'));
      return [];
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (values: any) => {
    if (selectedProvider === null) return;
    
    setLoading(true);
    try {
      const currentProvider = providers.find(p => p.id === selectedProvider);
      if (!currentProvider) return;
      
      // 构造后端需要的Provider对象
      const providerData = new Provider({
        provider_name: values.providerName || currentProvider.provider_name,
        provider_type: currentProvider.provider_type,
        base_url: values.baseUrl,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey || '',
        enable: values.enabled,
        is_default: currentProvider.is_default || false,
        default_model_id: currentProvider.default_model_id || null,
      });
      
      await Service.UpdateProvider(selectedProvider, providerData);
      
      // 更新本地状态
      const updatedProviders = providers.map(p => 
        p.id === selectedProvider ? {
          ...p,
          provider_name: values.providerName || p.provider_name, // 更新供应商名称
          api_key: values.apiKey,
          base_url: values.baseUrl,
          file_upload_base_url: values.fileUploadBaseUrl || null,
          enable: values.enabled,
          is_default: p.is_default,
        } : p
      );
      setProviders(updatedProviders);

      message.success(t('settings.provider.messages.saveSuccess'));
    } catch (error) {
      console.error('保存失败:', error);
      message.error(t('settings.provider.messages.saveFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = async () => {
    if (selectedProvider === null) return;
    
    setTestingConnection(true);
    
    // 更新供应商状态为测试中
    const updatedProviders = providers.map(p => 
      p.id === selectedProvider ? { ...p, status: 'testing' as const } : p
    );
    setProviders(updatedProviders);
    
    try {
      const values = form.getFieldsValue();
      const currentProvider = providers.find(p => p.id === selectedProvider);
      if (!currentProvider) return;
      
      // 构造Provider对象进行测试
      const providerData = new Provider({
        provider_name: currentProvider.provider_name,
        base_url: values.baseUrl,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey,
        enable: values.enabled,
        is_default: currentProvider.is_default || false,
        default_model_id: currentProvider.default_model_id,
      });
      
      // TODO: 调用后端测试连接接口（如果有的话）
      // await Service.TestProviderConnection(providerData);
      
      // 模拟测试
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // 更新状态为已连接
      const connectedProviders = providers.map(p => 
        p.id === selectedProvider ? { ...p, status: 'connected' as const } : p
      );
      setProviders(connectedProviders);
      
      message.success(t('settings.provider.messages.testSuccess'));
    } catch (error) {
      // 更新状态为连接失败
      const failedProviders = providers.map(p => 
        p.id === selectedProvider ? { ...p, status: 'disconnected' as const } : p
      );
      setProviders(failedProviders);
      
      console.error('连接测试失败:', error);
      message.error(t('settings.provider.messages.testFailed'));
    } finally {
      setTestingConnection(false);
    }
  };

  // 加载支持的供应商列表
  const loadSupportProviders = async () => {
    setLoadingSupportProviders(true);
    try {
      const providers = await Service.GetSupportProviders();
      setSupportProviders(providers || []);
    } catch (error) {
      console.error('加载支持的供应商列表失败:', error);
      message.error(t('settings.provider.messages.loadSupportFailed'));
    } finally {
      setLoadingSupportProviders(false);
    }
  };

  // 打开添加供应商窗口
  const handleAddProvider = () => {
    void Service.OpenAddProviderWindow();
  };

  // 监听来自添加供应商窗口的变更事件，刷新列表
  useEffect(() => {
    const cancel = Events.On('settings:providers:changed', () => {
      void loadProviderConfigs();
    });
    return () => {
      cancel?.();
      Events.Off('settings:providers:changed');
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleDeleteProvider = async (providerId: number) => {
    try {
      // 调用后端删除接口
      await Service.DeleteProvider(providerId);
      
      const updatedProviders = providers.filter(p => p.id !== providerId);
      setProviders(updatedProviders);
      
      // 如果删除的是当前选中的供应商，切换到第一个
      if (selectedProvider === providerId && updatedProviders.length > 0) {
        setSelectedProvider(updatedProviders[0].id);
      } else if (selectedProvider === providerId) {
        setSelectedProvider(null);
      }

      message.success(t('settings.provider.messages.deleteSuccess'));
    } catch (error) {
      console.error('删除供应商失败:', error);
      message.error(t('settings.provider.messages.deleteFailed'));
    }
  };

  const handleDeleteCurrentProvider = async () => {
    if (selectedProvider !== null) {
      await handleDeleteProvider(selectedProvider);
    }
  };

  // 刷新供应商模型列表
  const handleRefreshModels = async () => {
    if (selectedProvider === null || isCreatingNew) return;
    
    setLoading(true);
    try {
      // 调用后端接口刷新模型
      await Service.UpdateProviderModels(selectedProvider);
      
      // 重新加载供应商列表以获取最新的模型数据
      const updatedProviders = await loadProviderConfigs();
      const refreshedProvider = updatedProviders.find(provider => provider.id === selectedProvider);
      const refreshedDefaultModel = findDefaultModelForProvider(refreshedProvider);
      if (refreshedProvider?.is_default && refreshedDefaultModel) {
        setDefaultModelConfig({ modelId: refreshedDefaultModel.id, modelName: refreshedDefaultModel.model });
        await refetchModels();
      }

      message.success(t('settings.provider.messages.refreshModelsSuccess'));
    } catch (error) {
      console.error('刷新模型失败:', error);
      message.error(t('settings.provider.messages.refreshModelsFailed'));
    } finally {
      setLoading(false);
    }
  };
  const handleAddCustomModel = async () => {
    if (selectedProvider === null || !customModelName.trim()) return;
    
    setAddingCustomModel(true);
    try {
      await Service.AddProviderCustomModel(selectedProvider, customModelName.trim());
      setCustomModelName('');
      await loadProviderConfigs();
      message.success(t('settings.provider.messages.addCustomModelSuccess'));
    } catch (error) {
      console.error('添加自定义模型失败:', error);
      message.error(t('settings.provider.messages.addCustomModelFailed'));
    } finally {
      setAddingCustomModel(false);
    }
  };

  const handleDeleteCustomModel = async (modelName: string) => {
    if (selectedProvider === null) return;
    
    try {
      await Service.DeleteProviderCustomModel(selectedProvider, modelName);
      await loadProviderConfigs();
      message.success(t('settings.provider.messages.deleteCustomModelSuccess'));
    } catch (error) {
      console.error('删除自定义模型失败:', error);
      message.error(t('settings.provider.messages.deleteCustomModelFailed'));
    }
  };

  const handleCancelCreate = () => {
    if (newProviderTempId !== null) {
      // 从列表中移除占位符
      setProviders(prev => prev.filter(p => p.id !== newProviderTempId));
      
      // 重置状态
      setIsCreatingNew(false);
      setNewProviderTempId(null);
      
      // 选中第一个供应商（如果存在）
      const realProviders = providers.filter(p => p.id > 0);
      if (realProviders.length > 0) {
        setSelectedProvider(realProviders[0].id);
      } else {
        setSelectedProvider(null);
        form.resetFields();
      }
    }
  };
  
  // 创建新供应商
  const handleCreateProvider = async (values: any) => {
    if (!isCreatingNew) return;
    
    setLoading(true);
    try {
      // 构造新供应商数据
      const newProviderData = new Provider({
        provider_name: values.providerName || t('settings.provider.configFallbackName'),
        base_url: values.baseUrl,
        api_key: values.apiKey,
        enable: values.enabled,
        is_default: false,
        default_model_id: values.defaultModel || 0,
      });
      
      // 调用后端创建接口
      await Service.AddProvider(newProviderData);
      
      // 重新加载供应商列表
      await loadProviderConfigs();
      
      // 重置创建状态
      setIsCreatingNew(false);
      setNewProviderTempId(null);

      message.success(t('settings.provider.messages.createSuccess'));
    } catch (error) {
      console.error('创建供应商失败:', error);
      message.error(t('settings.provider.messages.createFailed'));
    } finally {
      setLoading(false);
    }
  };

  // 获取供应商图标和描述的辅助函数
  const getProviderExtras = (providerName: string, providerType?: ProviderType, baseUrl?: string) => {
    const matchedSupportProvider = findSupportProvider(providerType, providerName, baseUrl);
    if (matchedSupportProvider) {
      return {
        icon: matchedSupportProvider.icon || '🔧',
        description: matchedSupportProvider.description || t('settings.provider.descriptions.thirdParty'),
      };
    }

    const extras: { [key: string]: { icon: string; description: string } } = {
      'openai': {
        icon: '🤖',
        description: t('settings.provider.descriptions.openai'),
      },
      'anthropic': {
        icon: '🧠', 
        description: t('settings.provider.descriptions.anthropic'),
      },
      'gemini': {
        icon: '✨',
        description: t('settings.provider.descriptions.google'),
      },
      'google': {
        icon: '✨',
        description: t('settings.provider.descriptions.google'),
      }
    };
    
    const key = providerName.toLowerCase();
    return extras[key] || { icon: '🔧', description: t('settings.provider.descriptions.thirdParty') };
  };

  const renderProviderAvatar = (
    icon?: string,
    providerName?: string,
    size: number = 28,
    enabled: boolean = true,
  ) => {
    const iconSrc = resolveIconSrc(icon);
    const fallbackText = providerName?.trim().slice(0, 1).toUpperCase() || 'P';

    if (icon) {
      if (iconSrc) {
        return (
          <Avatar
            size={size}
            src={iconSrc}
            style={{ backgroundColor: enabled ? 'var(--brand-soft-teal)' : 'var(--background-color-dark)' }}
          >
            {fallbackText}
          </Avatar>
        );
      }
      return (
        <Avatar
          size={size}
          style={{
            backgroundColor: enabled ? 'var(--brand-soft-teal)' : 'var(--background-color-dark)',
            fontSize: `${Math.round(size * 0.57)}px`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          {icon.length <= 2 ? icon : fallbackText}
        </Avatar>
      );
    }
    return (
      <Avatar
        size={size}
        icon={<ApiOutlined />}
        style={{ backgroundColor: enabled ? 'var(--brand-teal)' : 'var(--text-color-disabled)' }}
      />
    );
  };

  const currentProvider = providers.find(p => p.id === selectedProvider);
  const getModelDisplayName = (model: Model) => model.alias || model.model;
  const availableModelsForProvider = [...(currentProvider?.models || [])].sort((a, b) => {
    const nameA = getModelDisplayName(a).toLocaleLowerCase();
    const nameB = getModelDisplayName(b).toLocaleLowerCase();
    if (nameA !== nameB) {
      return nameA.localeCompare(nameB);
    }
    return a.model.localeCompare(b.model);
  });

  const handleSetDefaultProvider = async (providerId: number) => {
    setSettingDefaultProviderId(providerId);
    try {
      const defaultModel = await Service.SetDefaultProvider(providerId);
      setProviders(prev => prev.map(provider => ({
        ...provider,
        is_default: provider.id === providerId,
        default_model_id: provider.id === providerId ? (defaultModel?.id ?? null) : provider.default_model_id,
      })));

      if (defaultModel) {
        setDefaultModelConfig({ modelId: defaultModel.id, modelName: defaultModel.model });
        await refetchModels();
        message.success(t('settings.provider.messages.setDefaultProviderSuccess'));
      } else {
        message.warning(t('settings.provider.messages.setDefaultProviderNoModel'));
      }
    } catch (error) {
      console.error('设置默认供应商失败:', error);
      message.error(t('settings.provider.messages.setDefaultProviderFailed'));
    } finally {
      setSettingDefaultProviderId(null);
    }
  };

  const handleSetDefaultModel = async (model: Model) => {
    if (!currentProvider) return;

    setSettingDefaultModelId(model.id);
    try {
      const providerData = new Provider({
        provider_name: currentProvider.provider_name,
        provider_type: currentProvider.provider_type,
        base_url: currentProvider.base_url,
        file_upload_base_url: currentProvider.file_upload_base_url || null,
        api_key: currentProvider.api_key,
        enable: currentProvider.enable,
        is_default: currentProvider.is_default || false,
        default_model_id: model.id,
      });

      await Service.UpdateProvider(currentProvider.id, providerData);
      setProviders(prev => prev.map(provider => (
        provider.id === currentProvider.id
          ? { ...provider, default_model_id: model.id }
          : provider
      )));

      if (currentProvider.is_default) {
        setDefaultModelConfig({ modelId: model.id, modelName: model.model });
        await refetchModels();
      }

      message.success(t('settings.provider.messages.setDefaultModelSuccess'));
    } catch (error) {
      console.error('设置默认模型失败:', error);
      message.error(t('settings.provider.messages.setDefaultModelFailed'));
    } finally {
      setSettingDefaultModelId(null);
    }
  };

  const getProviderMenuItems = (provider: ProviderConfig): MenuProps['items'] => [
    {
      key: 'default',
      icon: <StarOutlined />,
      label: t('settings.provider.actions.setDefault'),
      disabled: provider.is_default || settingDefaultProviderId !== null,
    },
    {
      key: 'delete',
      danger: true,
      icon: <DeleteOutlined />,
      label: t('settings.provider.actions.delete'),
    },
  ];

  const handleProviderMenuClick = (provider: ProviderConfig): MenuProps['onClick'] => ({ key, domEvent }) => {
    domEvent.stopPropagation();
    if (key === 'default') {
      void handleSetDefaultProvider(provider.id);
      return;
    }
    if (key === 'delete') {
      Modal.confirm({
        title: t('settings.provider.confirm.deleteTitle'),
        content: t('settings.provider.confirm.deleteDescription'),
        okText: t('settings.provider.confirm.okDelete'),
        cancelText: t('common.cancel'),
        okButtonProps: { danger: true },
        onOk: () => handleDeleteProvider(provider.id),
      });
    }
  };

  const renderProviderList = () => (
    <Card
      title={t('settings.provider.listTitle')}
      className={styles.listCard}
      extra={
        <Button
          type="text"
          size="small"
          icon={<PlusOutlined />}
          onClick={handleAddProvider}
        >
          {t('settings.provider.add')}
        </Button>
      }
      styles={{ body: { padding: 0 } }}
    >
      <div className={styles.providerItems}>
        {providers.map(provider => (
          <div
            key={provider.id}
            className={`${styles.providerItem} ${
              selectedProvider === provider.id ? styles.selected : ''
            }`}
            onClick={() => {
              setSelectedProvider(provider.id);
              if (isMobile) setShowEditorOnMobile(true);
            }}
          >
            <div className={styles.providerItemHeader}>
              <div className={styles.providerItemTitle}>
                {renderProviderAvatar(provider.icon, provider.provider_name, 22, provider.enable)}
                <span className={styles.providerName}>{provider.provider_name}</span>
                {provider.is_default && (
                  <Tag className={styles.defaultProviderTag}>
                    {t('settings.provider.tags.default')}
                  </Tag>
                )}
              </div>
              <div className={styles.providerItemActions}>
                <Switch
                  size="small"
                  checked={provider.enable}
                  onChange={async (checked, e) => {
                    e.stopPropagation();
                    try {
                      const updatedProviders = providers.map(p =>
                        p.id === provider.id ? { ...p, enable: checked } : p
                      );
                      setProviders(updatedProviders);
                      const providerData = new Provider({
                        provider_name: provider.provider_name,
                        provider_type: provider.provider_type,
                        base_url: provider.base_url,
                        file_upload_base_url: provider.file_upload_base_url || null,
                        api_key: provider.api_key,
                        enable: checked,
                        is_default: provider.is_default || false,
                        default_model_id: provider.default_model_id,
                      });
                      await Service.UpdateProvider(provider.id, providerData);
                    } catch (error) {
                      console.error('更新供应商状态失败:', error);
                      message.error(t('settings.provider.messages.updateStatusFailed'));
                      const revertProviders = providers.map(p =>
                        p.id === provider.id ? { ...p, enable: !checked } : p
                      );
                      setProviders(revertProviders);
                    }
                  }}
                />
                <Dropdown
                  trigger={['click']}
                  menu={{
                    items: getProviderMenuItems(provider),
                    onClick: handleProviderMenuClick(provider),
                  }}
                >
                  <Button
                    type="text"
                    size="small"
                    icon={<EllipsisOutlined />}
                    className={styles.moreBtn}
                    onClick={(e) => e.stopPropagation()}
                  />
                </Dropdown>
              </div>
            </div>
            <div className={styles.providerItemMeta}>
              {provider.base_url}
            </div>
          </div>
        ))}
      </div>
    </Card>
  );

  const renderEditor = () => (
    <Card
      title={t('settings.provider.configTitle', { name: currentProvider?.provider_name || t('settings.provider.configFallbackName') })}
      className={styles.configCard}
    >
            <Form
                form={form}
                layout="vertical"
                onFinish={isCreatingNew ? handleCreateProvider : handleSave} // 根据状态决定调用哪个函数
                initialValues={currentProvider}
                className={styles.formShell}
              >
                <div className={styles.formScrollArea}>
                  <Alert
                    message={t('settings.provider.securityAlert')}
                    type="info"
                    showIcon
                    className={styles.securityAlert}
                  />

                  <Form.Item
                    label={t('settings.provider.fields.enabled')}
                    name="enabled"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>

                  <Form.Item
                    label={t('settings.provider.fields.providerName')}
                    name="providerName"
                    rules={[
                      { required: true, message: t('settings.provider.validation.providerNameRequired') },
                      { max: 50, message: t('settings.provider.validation.providerNameMax') },
                    ]}
                  >
                    <Input 
                      placeholder={t('settings.provider.placeholders.providerName')} 
                    />
                  </Form.Item>

                  <Form.Item
                    label={t('settings.provider.fields.apiKey')}
                    name="apiKey"
                    rules={[
                      { required: false, message: t('settings.provider.validation.apiKeyRequired') },
                    ]}
                  >
                    <Input.Password
                      placeholder={t('settings.provider.placeholders.apiKey')}
                      iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                    />
                  </Form.Item>

                  <Form.Item
                    label={t('settings.provider.fields.baseUrl')}
                    name="baseUrl"
                    rules={[
                      { required: true, message: t('settings.provider.validation.baseUrlRequired') },
                      { type: 'url', message: t('settings.provider.validation.invalidUrl') },
                    ]}
                  >
                    <Input placeholder={t('settings.provider.placeholders.baseUrl')} />
                  </Form.Item>

                  <Form.Item
                    label={
                      <Space>
                        <span>{t('settings.provider.fields.fileUploadUrl')}</span>
                        <Tooltip title={t('settings.provider.helper.fileUploadUrl')}>
                          <QuestionCircleOutlined style={{ color: 'var(--text-color-secondary)', cursor: 'help' }} />
                        </Tooltip>
                      </Space>
                    }
                    name="fileUploadBaseUrl"
                    rules={[
                      { type: 'url', message: t('settings.provider.validation.invalidUrl') },
                    ]}
                    style={{ display: 'none' }}
                  >
                    <Input placeholder={t('settings.provider.placeholders.fileUploadUrl')} />
                  </Form.Item>

                  {!isCreatingNew && currentProvider && (
                    <div className={styles.modelSection}>
                      <div className={styles.modelSectionHeader}>
                        <div className={styles.modelSectionTitle}>
                          <span>{t('settings.provider.fields.modelList')}</span>
                          <Tag className={styles.modelCountTag}>
                            {t('settings.provider.helper.modelCountShort', { count: availableModelsForProvider.length })}
                          </Tag>
                        </div>
                        <Space size={8}>
                          <Tooltip title={t('settings.provider.helper.refreshModels')}>
                            <Button
                              type="text"
                              icon={<ReloadOutlined />}
                              onClick={handleRefreshModels}
                              loading={loading}
                              disabled={isCreatingNew}
                            />
                          </Tooltip>
                          <Tooltip title={t('settings.provider.actions.addCustomModel')}>
                            <Button
                              type="text"
                              icon={<PlusOutlined />}
                              onClick={() => setShowCustomModelInput(value => !value)}
                            />
                          </Tooltip>
                        </Space>
                      </div>

                      {showCustomModelInput && (
                        <div className={styles.customModelInlineInput}>
                          <Input
                            placeholder={t('settings.provider.placeholders.customModel')}
                            value={customModelName}
                            onChange={(e) => setCustomModelName(e.target.value)}
                            onPressEnter={handleAddCustomModel}
                          />
                          <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={handleAddCustomModel}
                            loading={addingCustomModel}
                            disabled={!customModelName.trim()}
                          >
                            {t('common.add')}
                          </Button>
                        </div>
                      )}

                      <div className={styles.modelList}>
                        {availableModelsForProvider.length === 0 ? (
                          <Empty
                            image={Empty.PRESENTED_IMAGE_SIMPLE}
                            description={t('settings.provider.helper.noModels')}
                          />
                        ) : (
                          availableModelsForProvider.map(model => (
                            <div key={model.id} className={styles.modelItem}>
                              <div className={styles.modelInfo}>
                                <div className={styles.modelNameRow}>
                                  <span className={styles.modelName}>{getModelDisplayName(model)}</span>
                                  {currentProvider.default_model_id === model.id && (
                                    <Tag className={styles.defaultModelTag}>
                                      {t('settings.provider.tags.default')}
                                    </Tag>
                                  )}
                                  {model.is_custom && (
                                    <Tag className={styles.customModelTag}>
                                      {t('settings.provider.tags.custom')}
                                    </Tag>
                                  )}
                                </div>
                                {model.alias && model.alias !== model.model && (
                                  <span className={styles.modelId}>{model.model}</span>
                                )}
                              </div>
                              <div className={styles.modelActions}>
                                {currentProvider.default_model_id !== model.id && (
                                  <Button
                                    type="text"
                                    size="small"
                                    loading={settingDefaultModelId === model.id}
                                    className={styles.setDefaultModelBtn}
                                    onClick={() => handleSetDefaultModel(model)}
                                  >
                                    {t('settings.provider.actions.setDefaultModel')}
                                  </Button>
                                )}
                                {model.is_custom && (
                                  <Tooltip title={t('settings.provider.helper.deleteCustomModel')}>
                                    <Popconfirm
                                      title={t('settings.provider.confirm.deleteCustomModelTitle')}
                                      description={t('settings.provider.confirm.deleteCustomModelDescription', { name: model.model })}
                                      onConfirm={() => handleDeleteCustomModel(model.model)}
                                      okText={t('common.confirm')}
                                      cancelText={t('common.cancel')}
                                      okButtonProps={{ danger: true }}
                                    >
                                      <Button
                                        type="text"
                                        size="small"
                                        icon={<CloseOutlined />}
                                        className={styles.customModelDeleteBtn}
                                      />
                                    </Popconfirm>
                                  </Tooltip>
                                )}
                              </div>
                            </div>
                          ))
                        )}
                      </div>
                    </div>
                  )}
                </div>

                <div className={styles.formFooter}>
                  <Space wrap>
                    {isCreatingNew ? (
                      <>
                        <Button
                          type="primary"
                          htmlType="submit"
                          icon={<PlusOutlined />}
                          loading={loading}
                        >
                          {t('settings.provider.actions.createShort')}
                        </Button>
                        <Button
                          icon={<DeleteOutlined />}
                          onClick={handleCancelCreate}
                        >
                          {t('common.cancel')}
                        </Button>
                      </>
                    ) : (
                      <>
                        <Button
                          type="primary"
                          htmlType="submit"
                          icon={<SaveOutlined />}
                          loading={loading}
                        >
                          {t('settings.provider.actions.save')}
                        </Button>
                        <Button
                          icon={<ReloadOutlined />}
                          onClick={handleTestConnection}
                          loading={testingConnection}
                        >
                          {t('settings.provider.actions.testConnection')}
                        </Button>
                        <Popconfirm
                          title={t('settings.provider.confirm.deleteTitle')}
                          description={
                            <div className={styles.deleteConfirm}>
                              <ExclamationCircleOutlined style={{ color: 'var(--warning-color)', marginRight: 8 }} />
                              {t('settings.provider.confirm.deleteCurrentDescription', { name: currentProvider?.provider_name || t('settings.provider.configFallbackName') })}
                            </div>
                          }
                          onConfirm={handleDeleteCurrentProvider}
                          okText={t('settings.provider.confirm.okDelete')}
                          cancelText={t('common.cancel')}
                          okButtonProps={{ danger: true }}
                        >
                          <Button
                            danger
                            icon={<DeleteOutlined />}
                          >
                            {t('settings.provider.actions.delete')}
                          </Button>
                        </Popconfirm>
                      </>
                    )}
                  </Space>
                </div>
              </Form>
          </Card>
  );

  return (
    <div className={`${styles.providerSettings} ${className || ''}`}>
      {isMobile ? (
        <>
          {!showEditorOnMobile && renderProviderList()}
          {showEditorOnMobile && (
            <div className={styles.mobileEditor}>
              <Button
                type="text"
                className={styles.mobileBackButton}
                icon={<ArrowLeftOutlined />}
                onClick={() => setShowEditorOnMobile(false)}
              >
                {t('settings.provider.actions.backToList')}
              </Button>
              {renderEditor()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderProviderList()}</div>
          <div className={styles.editorColumn}>{renderEditor()}</div>
        </div>
      )}
    </div>
  );
};

export default ProviderSettingPage;
