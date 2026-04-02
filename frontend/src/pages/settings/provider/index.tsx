import React, { useState, useEffect } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Select,
  Switch,
  Space,
  Divider,
  Alert,
  message,
  Typography,
  Popconfirm,
  Avatar,
  Tooltip,
  Modal,
} from 'antd';
import {
  ApiOutlined,
  SaveOutlined,
  ReloadOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  PlusOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  SettingOutlined,
  ArrowLeftOutlined,
  QuestionCircleOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import { useModels } from '@/hooks/useModels';
import { useModelStore } from '@/stores/modelStore';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { Provider, SupportProvider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import { ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models';
import styles from './index.module.scss';

const { Text } = Typography;
const { Option } = Select;

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
  default_model_id: number | null;
  models: any[];
  icon?: string;
  description?: string;
  status?: 'connected' | 'disconnected' | 'testing';
}

const ProviderSettingPage: React.FC<ProviderSettingPageProps> = ({ className }) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [testingConnection, setTestingConnection] = useState(false);
  const [providers, setProviders] = useState<ProviderConfig[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<number | null>(null); // 改为number类型
  const [isCreatingNew, setIsCreatingNew] = useState(false); // 新增：是否正在创建新供应商
  const [newProviderTempId, setNewProviderTempId] = useState<number | null>(null); // 新增：临时ID
  const [addProviderModalVisible, setAddProviderModalVisible] = useState(false); // 添加供应商对话框显示状态
  const [supportProviders, setSupportProviders] = useState<SupportProvider[]>([]); // 支持的供应商列表
  const [loadingSupportProviders, setLoadingSupportProviders] = useState(false); // 加载支持的供应商状态
  const [selectedSupportProvider, setSelectedSupportProvider] = useState<SupportProvider | null>(null); // 选中的支持供应商
  const [addProviderForm] = Form.useForm(); // 添加供应商表单
  const [customModelName, setCustomModelName] = useState(''); // 自定义模型名称输入
  const [addingCustomModel, setAddingCustomModel] = useState(false); // 添加自定义模型loading

  const { models: availableModels, isLoading: isLoadingModels } = useModels();
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
      // 转换字段名以适配表单
      // 只有当default_model_id存在且大于0时才设置，否则传递undefined以确保不被选中
      const defaultModelValue = (provider.default_model_id && provider.default_model_id > 0) ? provider.default_model_id : undefined;
      
      form.setFieldsValue({
        enabled: provider.enable,
        apiKey: provider.api_key,
        baseUrl: provider.base_url,
        fileUploadBaseUrl: provider.file_upload_base_url || '',
        providerName: provider.provider_name,
        defaultModel: defaultModelValue,
      });
      
      // 如果默认模型值为undefined，显式重置表单字段以确保清空选中状态
      if (defaultModelValue === undefined) {
        form.setFieldValue('defaultModel', undefined);
      }
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
        
        // 设置默认选中第一个供应商
        if (selectedProvider === null && formattedProviders.length > 0) {
          setSelectedProvider(formattedProviders[0].id);
        }
      }
    } catch (error) {
      console.error('加载供应商配置失败:', error);
      message.error(t('settings.provider.messages.loadProvidersFailed'));
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
      const defaultModelId = values.defaultModel || 0; // 如果没有选择默认模型，传递0
      const providerData = new Provider({
        provider_name: values.providerName || currentProvider.provider_name,
        provider_type: currentProvider.provider_type,
        base_url: values.baseUrl,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey || '',
        enable: values.enabled,
        default_model_id: defaultModelId,
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
          default_model_id: defaultModelId,
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

  // 打开添加供应商对话框
  const handleAddProvider = () => {
    setAddProviderModalVisible(true);
    setSelectedSupportProvider(null);
    addProviderForm.resetFields();
    loadSupportProviders();
  };

  // 选择支持的供应商
  const handleSelectSupportProvider = (supportProvider: SupportProvider) => {
    setSelectedSupportProvider(supportProvider);
    // 设置表单默认值
    addProviderForm.setFieldsValue({
      enabled: true,
      providerName: supportProvider.name,
      apiKey: '',
      baseUrl: supportProvider.base_url,
      fileUploadBaseUrl: supportProvider.file_upload_base_url || '',
      defaultModel: undefined,
    });
  };

  // 取消选择供应商，返回列表
  const handleCancelSelectProvider = () => {
    setSelectedSupportProvider(null);
    addProviderForm.resetFields();
  };

  // 在对话框中创建供应商
  const handleCreateProviderInModal = async (values: any) => {
    if (!selectedSupportProvider) return;
    
    setLoading(true);
    try {
      // 构造新供应商数据
      const newProviderData = new Provider({
        provider_name: values.providerName || selectedSupportProvider.name,
        provider_type: selectedSupportProvider.provider_type,
        base_url: values.baseUrl || selectedSupportProvider.base_url,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey || '',
        enable: values.enabled,
        default_model_id: values.defaultModel || 0,
      });
      
      await Service.AddProvider(newProviderData);
      
      // 重新加载供应商列表
      await loadProviderConfigs();
      
      // 关闭对话框并重置状态
      setAddProviderModalVisible(false);
      setSelectedSupportProvider(null);
      addProviderForm.resetFields();

      message.success(t('settings.provider.messages.createSuccess'));
    } catch (error) {
      console.error('创建供应商失败:', error);
      message.error(t('settings.provider.messages.createFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteProvider = async (providerId: number) => {
    try {
      // 调用后端删除接口
      await Service.DeleteProvider(providerId);
      
      const updatedProviders = providers.filter(p => p.id !== providerId);
      setProviders(updatedProviders);
      
      // 如果删除的是当前选中的供应商，切换到第一个
      if (selectedProvider === providerId && updatedProviders.length > 0) {
        setSelectedProvider(updatedProviders[0].id);
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
      await loadProviderConfigs();

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
            style={{ backgroundColor: enabled ? 'var(--primary-color-light)' : 'var(--background-color-dark)' }}
          >
            {fallbackText}
          </Avatar>
        );
      }
      return (
        <Avatar
          size={size}
          style={{
            backgroundColor: enabled ? 'var(--primary-color-light)' : 'var(--background-color-dark)',
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
        style={{ backgroundColor: enabled ? 'var(--primary-color)' : 'var(--text-color-disabled)' }}
      />
    );
  };

  const currentProvider = providers.find(p => p.id === selectedProvider);
  // 使用供应商的模型列表，而不是过滤所有模型
  const availableModelsForProvider = currentProvider?.models || [];

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
                <Popconfirm
                  title={t('settings.provider.confirm.deleteTitle')}
                  description={t('settings.provider.confirm.deleteDescription')}
                  onConfirm={() => handleDeleteProvider(provider.id)}
                  okText={t('common.confirm')}
                  cancelText={t('common.cancel')}
                >
                  <Button
                    type="text"
                    size="small"
                    icon={<DeleteOutlined />}
                    className={styles.deleteBtn}
                    onClick={(e) => e.stopPropagation()}
                  />
                </Popconfirm>
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

                  <Form.Item
                    label={t('settings.provider.fields.defaultModel')}
                    name="defaultModel"
                    help={t('settings.provider.helper.modelCount', { count: availableModelsForProvider.length })}
                  >
                    <Input.Group compact>
                      <Form.Item 
                        name="defaultModel" 
                        noStyle
                      >
                        <Select 
                          key={`defaultModel-${selectedProvider}`} // 添加key以在供应商切换时重置组件
                          placeholder={t('settings.provider.placeholders.selectDefaultModel')}
                          allowClear
                          showSearch
                          notFoundContent={t('settings.provider.helper.noModels')}
                          style={{ width: 'calc(100% - 40px)' }} // 为刷新按钮留出空间
                          filterOption={(input, option) => {
                            if (!input) return true;
                            const searchValue = input.toLowerCase();
                            // 从 option 中获取模型数据
                            const modelId = option?.value;
                            const model = availableModelsForProvider.find(m => m.id === modelId);
                            if (!model) return false;
                            
                            // 搜索模型名称、别名和模型 ID
                            const modelName = (model.model || '').toLowerCase();
                            const modelAlias = (model.alias || '').toLowerCase();
                            const modelIdStr = String(model.id || '').toLowerCase();
                            
                            return modelName.includes(searchValue) || 
                                   modelAlias.includes(searchValue) || 
                                   modelIdStr.includes(searchValue);
                          }}
                        >
                          {availableModelsForProvider.filter(m => !m.is_custom).map(model => (
                            <Option key={model.id} value={model.id}>
                              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <span>{model.alias || model.model}</span>
                                <span style={{ color: 'var(--text-color-secondary)', fontSize: '12px' }}>
                                  {model.model}
                                </span>
                              </div>
                            </Option>
                          ))}
                        </Select>
                      </Form.Item>
                      <Tooltip title={t('settings.provider.helper.refreshModels')}>
                        <Button 
                          icon={<ReloadOutlined />}
                          onClick={handleRefreshModels}
                          loading={loading}
                          disabled={isCreatingNew}
                          style={{ width: '40px' }}
                        />
                      </Tooltip>
                    </Input.Group>
                  </Form.Item>

                  {!isCreatingNew && currentProvider && (
                    <div className={styles.customModelSection}>
                      <div className={styles.customModelHeader}>
                        <div className={styles.customModelTitle}>
                          <SettingOutlined className={styles.customModelIcon} />
                          <span>{t('settings.provider.helper.customModels')}</span>
                        </div>
                        <Text type="secondary" className={styles.customModelDesc}>
                          {t('settings.provider.helper.customModelsDesc')}
                        </Text>
                      </div>
                      <div className={styles.customModelInput}>
                        <Input
                          placeholder={t('settings.provider.placeholders.customModel')}
                          value={customModelName}
                          onChange={(e) => setCustomModelName(e.target.value)}
                          onPressEnter={handleAddCustomModel}
                          className={styles.customModelInputField}
                        />
                        <Button
                          type="primary"
                          icon={<PlusOutlined />}
                          onClick={handleAddCustomModel}
                          loading={addingCustomModel}
                          disabled={!customModelName.trim()}
                          className={styles.customModelAddBtn}
                        >
                          {t('common.add')}
                        </Button>
                      </div>
                      {availableModelsForProvider.filter(m => m.is_custom).length > 0 && (
                        <div className={styles.customModelList}>
                          {availableModelsForProvider
                            .filter(m => m.is_custom)
                            .map(model => (
                              <div key={model.id} className={styles.customModelItem}>
                                <span className={styles.customModelName}>
                                  {model.alias || model.model}
                                </span>
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
                              </div>
                            ))}
                        </div>
                      )}
                    </div>
                  )}

                  <Divider className={styles.formDivider} />
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

      {/* 添加供应商对话框 */}
      <Modal
        title={
          selectedSupportProvider ? (
            <Space>
              <Button
                type="text"
                icon={<ArrowLeftOutlined />}
                onClick={handleCancelSelectProvider}
                style={{ padding: 0 }}
              />
              <span>{t('settings.provider.modal.addProvider', { name: selectedSupportProvider.name })}</span>
            </Space>
          ) : (
            t('settings.provider.modal.selectProvider')
          )
        }
        open={addProviderModalVisible}
        onCancel={() => {
          setAddProviderModalVisible(false);
          setSelectedSupportProvider(null);
          addProviderForm.resetFields();
        }}
        footer={null}
        width={isMobile ? 'calc(100vw - 32px)' : (selectedSupportProvider ? 700 : 600)}
        centered
        getContainer={() => document.body}
        zIndex={2002}
        wrapClassName={styles.addProviderModal}
      >
        {!selectedSupportProvider ? (
          <div className={styles.supportProviderList}>
            {supportProviders.map((item) => {
              const extras = getProviderExtras(item.name, item.provider_type);
              const iconSrc = resolveIconSrc(item.icon);
              return (
                <div
                  key={item.provider_type || item.name}
                  className={styles.supportProviderItem}
                  onClick={() => handleSelectSupportProvider(item)}
                >
                  <Avatar
                    size={40}
                    src={iconSrc}
                    style={{
                      backgroundColor: item.icon ? 'transparent' : 'var(--primary-color-light)',
                      fontSize: '20px',
                      flexShrink: 0,
                    }}
                  >
                    {!item.icon && extras.icon}
                  </Avatar>
                  <div className={styles.supportProviderInfo}>
                    <span className={styles.supportProviderName}>{item.name}</span>
                    <div className={styles.supportProviderDesc}>
                      {item.description || extras.description}
                    </div>
                    {item.base_url && (
                      <div className={styles.supportProviderUrl}>
                        <ApiOutlined style={{ marginRight: 4 }} />
                        {item.base_url}
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
            {supportProviders.length === 0 && !loadingSupportProviders && (
              <div className={styles.supportProviderEmpty}>{t('settings.provider.modal.emptyProviders')}</div>
            )}
          </div>
        ) : (
          <Form
            form={addProviderForm}
            layout="vertical"
            onFinish={handleCreateProviderInModal}
          >
            <Alert
              message={t('settings.provider.securityAlert')}
              type="info"
              showIcon
              className={styles.modalAlert}
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
              <Input placeholder={t('settings.provider.placeholders.providerName')} />
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

            <Form.Item
              label={t('settings.provider.fields.defaultModel')}
              name="defaultModel"
              help={t('settings.provider.helper.chooseModelAfterSave')}
            >
              <Select
                placeholder={t('settings.provider.placeholders.chooseModelAfterSave')}
                allowClear
                showSearch
                disabled
                notFoundContent={t('settings.provider.helper.saveProviderFirst')}
              >
                <Option value={0} disabled>{t('settings.provider.helper.saveProviderFirst')}</Option>
              </Select>
            </Form.Item>

            <Divider />

            <Form.Item className={styles.modalFormFooter}>
              <Space className={styles.modalFormActions}>
                <Button
                  onClick={() => {
                    setAddProviderModalVisible(false);
                    setSelectedSupportProvider(null);
                    addProviderForm.resetFields();
                  }}
                >
                  {t('common.cancel')}
                </Button>
                <Button
                  type="primary"
                  htmlType="submit"
                  icon={<PlusOutlined />}
                  loading={loading}
                >
                  {t('settings.provider.actions.create')}
                </Button>
              </Space>
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  );
};

export default ProviderSettingPage;
