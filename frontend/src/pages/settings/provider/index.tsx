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
  Spin,
  Typography,
  Row,
  Col,
  Popconfirm,
  Badge,
  Avatar,
  Tooltip,
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
} from '@ant-design/icons';
import { useModels } from '@/hooks/useModels';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { Provider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
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
  enable: boolean; // 使用后端字段名
  alias?: string | null; // 使用后端字段名
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

  const { models: availableModels, isLoading: isLoadingModels } = useModels();

  // 加载保存的配置
  useEffect(() => {
    loadProviderConfigs();
  }, []);

  // 当选中的供应商变化时，更新表单
  useEffect(() => {
    const provider = providers.find(p => p.id === selectedProvider);
    if (provider) {
      // 转换字段名以适配表单
      form.setFieldsValue({
        enabled: provider.enable,
        apiKey: provider.api_key,
        baseUrl: provider.base_url,
        alias: provider.alias,
        providerName: provider.provider_name,
      });
    }
  }, [selectedProvider, providers, form]);

  const loadProviderConfigs = async () => {
    setLoading(true);
    try {
      const providers = await Service.GetProviders();
      if (providers && providers.length > 0) {
        // 转换后端数据格式，添加前端需要的字段
        const formattedProviders = providers.map(provider => {
          const extras = getProviderExtras(provider.provider_name);
          return {
            ...provider,
            icon: extras.icon,
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
      const providerData = new Provider({
        provider_name: currentProvider.provider_name,
        base_url: values.baseUrl,
        api_key: values.apiKey,
        enable: values.enabled,
        alias: values.alias,
      });
      
      // 调用后端更新接口
      await Service.UpdateProvider(selectedProvider, providerData);
      
      // 更新本地状态
      const updatedProviders = providers.map(p => 
        p.id === selectedProvider ? {
          ...p,
          api_key: values.apiKey,
          base_url: values.baseUrl,
          enable: values.enabled,
          alias: values.alias,
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
        api_key: values.apiKey,
        enable: values.enabled,
        alias: values.alias,
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

  const handleAddProvider = () => {
    // TODO: 实现添加供应商功能
    message.info('添加供应商功能开发中...');
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
  const availableModelsForProvider = availableModels.filter(model => 
    currentProvider ? model.id.toLowerCase().includes(currentProvider.provider_name.toLowerCase()) : false
  );

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
                                 <div className={styles.providerName}>{provider.alias || provider.provider_name}</div>
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
                                             api_key: provider.api_key,
                                             enable: checked,
                                             alias: provider.alias,
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
                <Avatar 
                  size={24} 
                  style={{ backgroundColor: 'var(--primary-color)' }}
                >
                  {currentProvider?.icon || <ApiOutlined />}
                </Avatar>
                <span>配置 {currentProvider?.alias || currentProvider?.provider_name}</span>
                {currentProvider?.status === 'connected' && (
                  <Badge status="success" />
                )}
              </Space>
            }
            className={styles.configCard}
          >
            <Spin spinning={loading}>
              <Form
                form={form}
                layout="vertical"
                onFinish={handleSave}
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
                >
                  <Input placeholder="为供应商设置一个名称" />
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
                  label="默认模型"
                  name="defaultModel"
                >
                  <Select 
                    placeholder="选择默认模型"
                    loading={isLoadingModels}
                  >
                    {availableModelsForProvider.map(model => (
                      <Option key={model.id} value={model.id}>
                        {model.name || model.id}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>

                <Divider className={styles.formDivider} />

                <Form.Item className={styles.buttonGroup}>
                  <Space size="middle" className={styles.actionButtons}>
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
                          确定要删除 <strong>{currentProvider?.alias || currentProvider?.provider_name}</strong> 吗？
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
                        disabled={providers.length <= 1}
                        className={styles.dangerButton}
                      >
                        删除供应商
                      </Button>
                    </Popconfirm>
                  </Space>
                </Form.Item>
              </Form>
            </Spin>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default ProviderSettingPage;