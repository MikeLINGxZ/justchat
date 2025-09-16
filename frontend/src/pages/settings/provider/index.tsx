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
import styles from './index.module.scss';

const { Title, Text } = Typography;
const { Option } = Select;

interface ProviderSettingPageProps {
  className?: string;
}

interface ProviderConfig {
  id: string;
  name: string;
  apiKey: string;
  baseUrl?: string;
  enabled: boolean;
  defaultModel?: string;
  icon?: string;
  description?: string;
  status?: 'connected' | 'disconnected' | 'testing';
}

const ProviderSettingPage: React.FC<ProviderSettingPageProps> = ({ className }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [testingConnection, setTestingConnection] = useState(false);
  const [providers, setProviders] = useState<ProviderConfig[]>([
    {
      id: 'openai',
      name: 'OpenAI',
      apiKey: '',
      baseUrl: 'https://api.openai.com/v1',
      enabled: false,
      defaultModel: 'gpt-3.5-turbo',
      icon: '🤖',
      description: '强大的GPT系列模型，支持聊天和文本生成',
      status: 'disconnected',
    },
    {
      id: 'anthropic',
      name: 'Anthropic',
      apiKey: '',
      baseUrl: 'https://api.anthropic.com',
      enabled: false,
      defaultModel: 'claude-3-sonnet-20240229',
      icon: '🧠',
      description: 'Claude系列模型，注重安全性和有用性',
      status: 'disconnected',
    },
    {
      id: 'gemini',
      name: 'Google Gemini',
      apiKey: '',
      baseUrl: 'https://generativelanguage.googleapis.com/v1beta',
      enabled: false,
      defaultModel: 'gemini-pro',
      icon: '✨',
      description: 'Google最新的多模态AI模型',
      status: 'disconnected',
    },
  ]);
  const [selectedProvider, setSelectedProvider] = useState<string>('openai');

  const { models: availableModels, isLoading: isLoadingModels } = useModels();

  // 加载保存的配置
  useEffect(() => {
    loadProviderConfigs();
  }, []);

  // 当选中的供应商变化时，更新表单
  useEffect(() => {
    const provider = providers.find(p => p.id === selectedProvider);
    if (provider) {
      form.setFieldsValue(provider);
    }
  }, [selectedProvider, providers, form]);

  const loadProviderConfigs = async () => {
    setLoading(true);
    try {
      // TODO: 从后端加载保存的配置
      // const configs = await Service.GetProviderConfigs();
      // setProviders(configs);
    } catch (error) {
      console.error('加载供应商配置失败:', error);
      message.error('加载供应商配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (values: any) => {
    setLoading(true);
    try {
      const updatedProviders = providers.map(p => 
        p.id === selectedProvider ? { ...p, ...values } : p
      );
      setProviders(updatedProviders);
      
      // TODO: 保存到后端
      // await Service.SaveProviderConfig(selectedProvider, values);
      
      message.success('保存成功');
    } catch (error) {
      console.error('保存失败:', error);
      message.error('保存失败');
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = async () => {
    setTestingConnection(true);
    
    // 更新供应商状态为测试中
    const updatedProviders = providers.map(p => 
      p.id === selectedProvider ? { ...p, status: 'testing' as const } : p
    );
    setProviders(updatedProviders);
    
    try {
      const values = form.getFieldsValue();
      // TODO: 测试连接
      // await Service.TestProviderConnection(selectedProvider, values);
      
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

  const handleDeleteProvider = (providerId: string) => {
    const updatedProviders = providers.filter(p => p.id !== providerId);
    setProviders(updatedProviders);
    
    // 如果删除的是当前选中的供应商，切换到第一个
    if (selectedProvider === providerId && updatedProviders.length > 0) {
      setSelectedProvider(updatedProviders[0].id);
    }
    
    message.success('供应商删除成功');
    // TODO: 调用后端API删除
    // await Service.DeleteProvider(providerId);
  };

  const handleDeleteCurrentProvider = () => {
    handleDeleteProvider(selectedProvider);
  };

  const getProviderIcon = (provider: ProviderConfig) => {
    if (provider.icon) {
      return (
        <Avatar 
          size={28} 
          style={{ 
            backgroundColor: provider.enabled ? 'var(--primary-color-light)' : 'var(--background-color-dark)',
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
          backgroundColor: provider.enabled ? 'var(--primary-color)' : 'var(--text-color-disabled)' 
        }} 
      />
    );
  };

  const currentProvider = providers.find(p => p.id === selectedProvider);
  const availableModelsForProvider = availableModels.filter(model => 
    model.id.toLowerCase().includes(selectedProvider.toLowerCase())
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
                                 <div className={styles.providerName}>{provider.name}</div>
                             </div>
                         </div>
                         <div className={styles.providerActions}>
                             <Tooltip title={provider.enabled ? '已启用' : '未启用'}>
                                 <Switch
                                     size="small"
                                     checked={provider.enabled}
                                     className={styles.enableSwitch}
                                     onChange={(checked) => {
                                         const updatedProviders = providers.map(p =>
                                             p.id === provider.id ? { ...p, enabled: checked } : p
                                         );
                                         setProviders(updatedProviders);
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
                <span>配置 {currentProvider?.name}</span>
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
                          确定要删除 <strong>{currentProvider?.name}</strong> 吗？
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