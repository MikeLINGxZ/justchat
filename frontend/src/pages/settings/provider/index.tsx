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
  Row,
  Col,
  Popconfirm,
  Badge,
  Avatar,
  Tooltip,
  Modal,
  List,
} from 'antd';
import {
  ApiOutlined,
  SaveOutlined,
  ReloadOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  PlusOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  SettingOutlined,
  ArrowLeftOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons';
import { useModels } from '@/hooks/useModels';
import { useModelStore } from '@/stores/modelStore';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { Provider, SupportProvider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import styles from './index.module.scss';

const { Title, Text } = Typography;
const { Option } = Select;

interface ProviderSettingPageProps {
  className?: string;
}

interface ProviderConfig {
  id: number; // 改为number类型，与后端一致
  provider_name: string; // 使用后端字段名
  api_key: string; // 使用后端字段名
  base_url: string; // 使用后端字段名
  file_upload_base_url?: string | null; // 文件上传URL
  enable: boolean; // 使用后端字段名
  default_model_id: number | null; // 默认模型ID，允许null
  models: any[]; // 供应商模型列表
  // 前端额外字段
  icon?: string;
  description?: string;
  status?: 'connected' | 'disconnected' | 'testing';
}

const ProviderSettingPage: React.FC<ProviderSettingPageProps> = ({ className }) => {
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

  const { models: availableModels, isLoading: isLoadingModels } = useModels();
  const { refetch: refetchModels } = useModelStore();

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

  const loadProviderConfigs = async () => {
    setLoading(true);
    try {
      const providers = await Service.GetProviders();
      if (providers && providers.length > 0) {
        // 转换后端数据格式，添加前端需要的字段
        const formattedProviders = providers.map(provider => {
          // 根据 provider_type 从 supportProviders 中匹配对应的 icon
          let icon: string | undefined;
          if (provider.provider_type && supportProviders.length > 0) {
            const matchedSupportProvider = supportProviders.find(
              sp => sp.provider_type === provider.provider_type
            );
            if (matchedSupportProvider && matchedSupportProvider.icon) {
              icon = matchedSupportProvider.icon;
            }
          }
          
          // 获取默认的 extras（用于 description 和 fallback icon）
          const extras = getProviderExtras(provider.provider_name);
          
          // 如果没有匹配到 icon，使用默认的 getProviderExtras
          if (!icon) {
            icon = extras.icon;
          }
          
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
      message.error('加载供应商配置失败');
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
        provider_name: values.providerName || currentProvider.provider_name, // 使用用户输入的名称
        provider_type: values.provider_type,
        base_url: values.baseUrl,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey,
        enable: values.enabled,
        default_model_id: defaultModelId,
      });
      
      // 调用后端更新接口
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

      message.success('保存成功');
    } catch (error) {
      console.error('保存失败:', error);
      message.error('保存失败');
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
      
      message.success('连接测试成功');
    } catch (error) {
      // 更新状态为连接失败
      const failedProviders = providers.map(p => 
        p.id === selectedProvider ? { ...p, status: 'disconnected' as const } : p
      );
      setProviders(failedProviders);
      
      console.error('连接测试失败:', error);
      message.error('连接测试失败');
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
      message.error('加载支持的供应商列表失败');
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
        provider_type: values.provider_type,
        base_url: values.baseUrl || selectedSupportProvider.base_url,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey,
        enable: values.enabled,
        default_model_id: values.defaultModel || 0,
      });
      
      // 调用后端创建接口
      await Service.AddProvider(newProviderData);
      
      // 重新加载供应商列表
      await loadProviderConfigs();
      
      // 关闭对话框并重置状态
      setAddProviderModalVisible(false);
      setSelectedSupportProvider(null);
      addProviderForm.resetFields();

      message.success('供应商创建成功');
    } catch (error) {
      console.error('创建供应商失败:', error);
      message.error('创建供应商失败');
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

      message.success('供应商删除成功');
    } catch (error) {
      console.error('删除供应商失败:', error);
      message.error('删除供应商失败');
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

      message.success('模型列表刷新成功');
    } catch (error) {
      console.error('刷新模型失败:', error);
      message.error('刷新模型失败');
    } finally {
      setLoading(false);
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
        provider_name: values.providerName || '新建供应商',
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

      message.success('供应商创建成功');
    } catch (error) {
      console.error('创建供应商失败:', error);
      message.error('创建供应商失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取供应商图标和描述的辅助函数
  const getProviderExtras = (providerName: string) => {
    const extras: { [key: string]: { icon: string; description: string } } = {
      'openai': {
        icon: '🤖',
        description: '强大的GPT系列模型，支持聊天和文本生成'
      },
      'anthropic': {
        icon: '🧠', 
        description: 'Claude系列模型，注重安全性和有用性'
      },
      'gemini': {
        icon: '✨',
        description: 'Google最新的多模态AI模型'
      },
      'google': {
        icon: '✨',
        description: 'Google最新的多模态AI模型'
      }
    };
    
    const key = providerName.toLowerCase();
    return extras[key] || { icon: '🔧', description: '第三方AI模型提供商' };
  };

  const getProviderIcon = (provider: ProviderConfig) => {
    if (provider.icon) {
      // 判断是否是图片路径（以 / 开头或者是 http/https 开头）
      const isImagePath = provider.icon.startsWith('/') || 
                         provider.icon.startsWith('http://') || 
                         provider.icon.startsWith('https://') ||
                         provider.icon.startsWith('data:');
      
      if (isImagePath) {
        return (
          <Avatar 
            size={28} 
            src={provider.icon}
            style={{ 
              backgroundColor: provider.enable ? 'var(--primary-color-light)' : 'var(--background-color-dark)',
            }}
          />
        );
      } else {
        // 否则作为 emoji 或文本显示
        return (
          <Avatar 
            size={28} 
            style={{ 
              backgroundColor: provider.enable ? 'var(--primary-color-light)' : 'var(--background-color-dark)',
              fontSize: '16px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}
          >
            {provider.icon}
          </Avatar>
        );
      }
    }
    return (
      <Avatar 
        size={28} 
        icon={<ApiOutlined />} 
        style={{ 
          backgroundColor: provider.enable ? 'var(--primary-color)' : 'var(--text-color-disabled)' 
        }} 
      />
    );
  };

  const currentProvider = providers.find(p => p.id === selectedProvider);
  // 使用供应商的模型列表，而不是过滤所有模型
  const availableModelsForProvider = currentProvider?.models || [];

  return (
    <div className={`${styles.providerSettings} ${className || ''}`}>


      <Row gutter={[24, 24]}>
        <Col xs={24} lg={8}>
          <Card 
            title={
              <Space>
                <SettingOutlined />
                <span>供应商列表</span>
              </Space>
            }
            className={styles.providerList}
            extra={
              <Tooltip title="添加新的供应商">
                <Button 
                  type="primary" 
                  size="small" 
                  icon={<PlusOutlined />}
                  onClick={handleAddProvider}
                  className={styles.addBtn}
                >
                  添加
                </Button>
              </Tooltip>
            }
            bodyStyle={{ padding: 0 }}
          >
            <div className={styles.providerItems}>
              {providers.map(provider => (
                <div
                  key={provider.id}
                  className={`${styles.providerItem} ${
                    selectedProvider === provider.id ? styles.selected : ''
                  }`}
                  onClick={() => setSelectedProvider(provider.id)}
                >
                 <div className={`${styles.providerItemBox}`}>
                     <div className={styles.providerContent}>
                         <div className={styles.providerLeft}>
                             {getProviderIcon(provider)}
                             <div className={styles.providerDetails}>
                                 <div className={styles.providerName}>{provider.provider_name}</div>
                             </div>
                         </div>
                         <div className={styles.providerActions}>
                             <Tooltip title={provider.enable ? '已启用' : '未启用'}>
                                 <Switch
                                     size="small"
                                     checked={provider.enable}
                                     className={styles.enableSwitch}
                                     onChange={async (checked) => {
                                         try {
                                           // 先更新UI状态
                                           const updatedProviders = providers.map(p =>
                                               p.id === provider.id ? { ...p, enable: checked } : p
                                           );
                                           setProviders(updatedProviders);
                                                           // 调用后端接口更新
                                           const providerData = new Provider({
                                             provider_name: provider.provider_name,
                                             base_url: provider.base_url,
                                             file_upload_base_url: provider.file_upload_base_url || null,
                                             api_key: provider.api_key,
                                             enable: checked,
                                             default_model_id: provider.default_model_id,
                                           });
                                           await Service.UpdateProvider(provider.id, providerData);
                                         } catch (error) {
                                           console.error('更新供应商状态失败:', error);
                                           message.error('更新供应商状态失败');
                                           // 恢复UI状态
                                           const revertProviders = providers.map(p =>
                                               p.id === provider.id ? { ...p, enable: !checked } : p
                                           );
                                           setProviders(revertProviders);
                                         }
                                     }}
                                 />
                             </Tooltip>
                             <Tooltip title="删除供应商">
                                 <Popconfirm
                                     title="删除供应商"
                                     description="确定要删除这个供应商吗？"
                                     onConfirm={() => handleDeleteProvider(provider.id)}
                                     okText="确定"
                                     cancelText="取消"
                                 >
                                     <Button
                                         type="text"
                                         size="small"
                                         icon={<DeleteOutlined />}
                                         className={styles.deleteBtn}
                                         onClick={(e) => e.stopPropagation()}
                                     />
                                 </Popconfirm>
                             </Tooltip>
                         </div>
                     </div>
                 </div>
                </div>
              ))}
            </div>
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card 
            title={
              <Space>
                {currentProvider?.icon ? (
                  (() => {
                    const isImagePath = currentProvider.icon.startsWith('/') || 
                                       currentProvider.icon.startsWith('http://') || 
                                       currentProvider.icon.startsWith('https://') ||
                                       currentProvider.icon.startsWith('data:');
                    return (
                      <Avatar 
                        size={24} 
                        src={isImagePath ? currentProvider.icon : undefined}
                        style={{ 
                          backgroundColor: isImagePath ? 'transparent' : 'var(--primary-color)',
                          fontSize: '14px',
                        }}
                      >
                        {!isImagePath && currentProvider.icon}
                      </Avatar>
                    );
                  })()
                ) : (
                  <Avatar 
                    size={24} 
                    icon={<ApiOutlined />}
                    style={{ backgroundColor: 'var(--primary-color)' }}
                  />
                )}
                <span>配置 {currentProvider?.provider_name}</span>
                {currentProvider?.status === 'connected' && (
                  <Badge status="success" />
                )}
              </Space>
            }
            className={styles.configCard}
          >
            <Form
                form={form}
                layout="vertical"
                onFinish={isCreatingNew ? handleCreateProvider : handleSave} // 根据状态决定调用哪个函数
                initialValues={currentProvider}
              >
                <Alert
                  message="API密钥将加密保存在本地，不会上传到任何服务器。"
                  type="info"
                  showIcon
                  style={{ marginBottom: 16 }}
                  className={styles.securityAlert}
                />

                <Form.Item
                  label="启用状态"
                  name="enabled"
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>

                <Form.Item
                  label="供应商名称"
                  name="providerName"
                  rules={[
                    { required: true, message: '请输入供应商名称' },
                    { max: 50, message: '供应商名称不能超过50个字符' },
                  ]}
                >
                  <Input 
                    placeholder="为供应商设置一个名称" 
                  />
                </Form.Item>

                <Form.Item
                  label="API 密钥"
                  name="apiKey"
                  rules={[
                    { required: true, message: '请输入API密钥' },
                  ]}
                >
                  <Input.Password
                    placeholder="请输入API密钥"
                    iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                  />
                </Form.Item>

                <Form.Item
                  label="API 基础URL"
                  name="baseUrl"
                  rules={[
                    { required: true, message: '请输入API基础URL' },
                    { type: 'url', message: '请输入正确的URL格式' },
                  ]}
                >
                  <Input placeholder="https://api.example.com/v1" />
                </Form.Item>

                <Form.Item
                  label={
                    <Space>
                      <span>文件上传URL</span>
                      <Tooltip title="多模态模型文件上传地址">
                        <QuestionCircleOutlined style={{ color: 'var(--text-color-secondary)', cursor: 'help' }} />
                      </Tooltip>
                    </Space>
                  }
                  name="fileUploadBaseUrl"
                  rules={[
                    { type: 'url', message: '请输入正确的URL格式' },
                  ]}
                >
                  <Input placeholder="https://api.example.com/v1/uploads" />
                </Form.Item>

                <Form.Item
                  label="默认模型"
                  name="defaultModel"
                  help={`当前供应商共有 ${availableModelsForProvider.length} 个可用模型`}
                >
                  <Input.Group compact>
                    <Form.Item 
                      name="defaultModel" 
                      noStyle
                    >
                      <Select 
                        key={`defaultModel-${selectedProvider}`} // 添加key以在供应商切换时重置组件
                        placeholder="选择默认模型"
                        allowClear
                        showSearch
                        notFoundContent="没有可用模型"
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
                        {availableModelsForProvider.map(model => (
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
                    <Tooltip title="刷新模型列表">
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

                <Divider className={styles.formDivider} />

                <Form.Item className={styles.buttonGroup}>
                  <Space size="middle" className={styles.actionButtons}>
                    {isCreatingNew ? (
                      // 新建供应商时的按钮
                      <>
                        <Button 
                          type="primary" 
                          htmlType="submit" 
                          icon={<PlusOutlined />}
                          loading={loading}
                          size="middle"
                          className={styles.primaryButton}
                        >
                          新建
                        </Button>
                        <Button 
                          icon={<DeleteOutlined />}
                          onClick={handleCancelCreate}
                          size="middle"
                          className={styles.testButton}
                        >
                          取消
                        </Button>
                      </>
                    ) : (
                      // 编辑供应商时的按钮
                      <>
                        <Button 
                          type="primary" 
                          htmlType="submit" 
                          icon={<SaveOutlined />}
                          loading={loading}
                          size="middle"
                          className={styles.primaryButton}
                        >
                          保存配置
                        </Button>
                        <Button 
                          icon={<ReloadOutlined />}
                          onClick={handleTestConnection}
                          loading={testingConnection}
                          size="middle"
                          className={styles.testButton}
                        >
                          测试连接
                        </Button>
                        <Popconfirm
                          title="删除供应商"
                          description={
                            <div className={styles.deleteConfirm}>
                              <ExclamationCircleOutlined style={{ color: 'var(--warning-color)', marginRight: 8 }} />
                              确定要删除 <strong>{currentProvider?.provider_name}</strong> 吗？
                            </div>
                          }
                          onConfirm={handleDeleteCurrentProvider}
                          okText="确定删除"
                          cancelText="取消"
                          okButtonProps={{ danger: true }}
                        >
                          <Button 
                            danger
                            icon={<DeleteOutlined />}
                            size="middle"
                            className={styles.dangerButton}
                          >
                            删除供应商
                          </Button>
                        </Popconfirm>
                      </>
                    )}
                  </Space>
                </Form.Item>
              </Form>
            </Card>
        </Col>
      </Row>

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
                title="返回选择列表"
              />
              <span>添加供应商 - {selectedSupportProvider.name}</span>
            </Space>
          ) : (
            '选择供应商'
          )
        }
        open={addProviderModalVisible}
        onCancel={() => {
          setAddProviderModalVisible(false);
          setSelectedSupportProvider(null);
          addProviderForm.resetFields();
        }}
        footer={null}
        width={selectedSupportProvider ? 700 : 600}
        centered
      >
        {!selectedSupportProvider ? (
          // 供应商列表
          <List
            loading={loadingSupportProviders}
            dataSource={supportProviders}
            renderItem={(item) => {
              const extras = getProviderExtras(item.name);
              return (
                <List.Item
                  style={{
                    cursor: 'pointer',
                    padding: '16px',
                    borderRadius: '8px',
                    marginBottom: '8px',
                    border: '1px solid var(--border-color)',
                    transition: 'all 0.3s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = 'var(--background-color-light)';
                    e.currentTarget.style.borderColor = 'var(--primary-color)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'transparent';
                    e.currentTarget.style.borderColor = 'var(--border-color)';
                  }}
                  onClick={() => handleSelectSupportProvider(item)}
                >
                  <List.Item.Meta
                    avatar={
                      <Avatar 
                        size={40} 
                        src={
                          item.icon 
                            ? (item.icon.startsWith('data:') || 
                                item.icon.startsWith('http://') || 
                                item.icon.startsWith('https://') ||
                                item.icon.startsWith('/')
                                ? item.icon 
                                : `data:image/png;base64,${item.icon}`)
                            : undefined
                        }
                        style={{ 
                          backgroundColor: item.icon ? 'transparent' : 'var(--primary-color-light)',
                          fontSize: '20px',
                        }}
                      >
                        {!item.icon && extras.icon}
                      </Avatar>
                    }
                    title={
                      <Space>
                        <span style={{ fontSize: '16px', fontWeight: 500 }}>{item.name}</span>
                      </Space>
                    }
                    description={
                      <div>
                        <div style={{ marginTop: '4px', color: 'var(--text-color-secondary)' }}>
                          {item.description || extras.description}
                        </div>
                        {item.base_url && (
                          <div style={{ marginTop: '8px', fontSize: '12px', color: 'var(--text-color-disabled)' }}>
                            <ApiOutlined style={{ marginRight: '4px' }} />
                            {item.base_url}
                          </div>
                        )}
                      </div>
                    }
                  />
                  <Button type="primary" icon={<PlusOutlined />}>
                    选择
                  </Button>
                </List.Item>
              );
            }}
            locale={{ emptyText: '暂无可用的供应商' }}
          />
        ) : (
          // 供应商配置表单
          <Form
            form={addProviderForm}
            layout="vertical"
            onFinish={handleCreateProviderInModal}
          >
            <Alert
              message="API密钥将加密保存在本地，不会上传到任何服务器。"
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />

            <Form.Item
              label="启用状态"
              name="enabled"
              valuePropName="checked"
            >
              <Switch />
            </Form.Item>

            <Form.Item
              label="供应商名称"
              name="providerName"
              rules={[
                { required: true, message: '请输入供应商名称' },
                { max: 50, message: '供应商名称不能超过50个字符' },
              ]}
            >
              <Input 
                placeholder="为供应商设置一个名称" 
              />
            </Form.Item>

            <Form.Item
              label="API 密钥"
              name="apiKey"
              rules={[
                { required: true, message: '请输入API密钥' },
              ]}
            >
              <Input.Password
                placeholder="请输入API密钥"
                iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
              />
            </Form.Item>

            <Form.Item
              label="API 基础URL"
              name="baseUrl"
              rules={[
                { required: true, message: '请输入API基础URL' },
                { type: 'url', message: '请输入正确的URL格式' },
              ]}
            >
              <Input placeholder="https://api.example.com/v1" />
            </Form.Item>

            <Form.Item
              label={
                <Space>
                  <span>文件上传URL</span>
                  <Tooltip title="多模态模型文件上传地址">
                    <QuestionCircleOutlined style={{ color: 'var(--text-color-secondary)', cursor: 'help' }} />
                  </Tooltip>
                </Space>
              }
              name="fileUploadBaseUrl"
              rules={[
                { type: 'url', message: '请输入正确的URL格式' },
              ]}
            >
              <Input placeholder="https://api.example.com/v1/uploads" />
            </Form.Item>

            <Form.Item
              label="默认模型"
              name="defaultModel"
              help="保存供应商后可以刷新模型列表"
            >
              <Select 
                placeholder="保存供应商后可选择默认模型"
                allowClear
                showSearch
                disabled
                notFoundContent="请先保存供应商"
              >
                <Option value={0} disabled>请先保存供应商</Option>
              </Select>
            </Form.Item>

            <Divider />

            <Form.Item style={{ marginBottom: 0 }}>
              <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
                <Button 
                  onClick={() => {
                    setAddProviderModalVisible(false);
                    setSelectedSupportProvider(null);
                    addProviderForm.resetFields();
                  }}
                >
                  取消
                </Button>
                <Button 
                  type="primary" 
                  htmlType="submit" 
                  icon={<PlusOutlined />}
                  loading={loading}
                >
                  创建供应商
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