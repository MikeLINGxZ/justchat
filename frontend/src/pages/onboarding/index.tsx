import React, { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Avatar,
  Button,
  Card,
  Empty,
  Form,
  Input,
  message,
  Space,
  Spin,
  Switch,
  Typography,
} from 'antd';
import {
  ArrowLeftOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  RocketOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts';
import { Provider, SupportProvider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/models.ts';
import styles from './index.module.scss';

const { Text, Title } = Typography;
type OnboardingStep = 'intro' | 'selectProvider' | 'configProvider';

const OnboardingPage: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [step, setStep] = useState<OnboardingStep>('intro');
  const [supportProviders, setSupportProviders] = useState<SupportProvider[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<SupportProvider | null>(null);
  const [loadFailed, setLoadFailed] = useState(false);

  useEffect(() => {
    const initialize = async () => {
      try {
        const initialized = await Service.IsInitialized();
        if (initialized) {
          navigate('/home', { replace: true });
          return;
        }

        const providers = await Service.GetSupportProviders();
        setLoadFailed(false);
        setSupportProviders(providers);
      } catch (error) {
        console.error('初始化欢迎页失败:', error);
        setLoadFailed(true);
        message.error('加载欢迎页失败');
      } finally {
        setLoading(false);
      }
    };

    void initialize();
  }, [navigate]);

  const applyProviderPreset = (provider: SupportProvider) => {
    setSelectedProvider(provider);
    form.setFieldsValue({
      enabled: true,
      providerName: provider.name,
      apiKey: '',
      baseUrl: provider.base_url || '',
    });
  };

  const selectProviderHint = useMemo(() => {
    if (loadFailed) {
      return '供应商列表加载失败，请稍后重试。';
    }
    if (supportProviders.length === 0) {
      return '当前没有可用的供应商预设，暂时无法完成初始化。';
    }
    return '选择一个供应商后，再填写访问地址和 API Key。';
  }, [loadFailed, supportProviders.length]);

  const configProviderHint = useMemo(() => {
    if (!selectedProvider) {
      return '请先返回上一步选择一个供应商。';
    }
    return selectedProvider.description || '填写访问地址和 API Key 后即可开始使用。';
  }, [selectedProvider]);

  const iconSrc = (icon: string) => {
    if (!icon) {
      return undefined;
    }
    if (icon.startsWith('/') || icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('data:')) {
      return icon;
    }
    return `data:image/png;base64,${icon}`;
  };

  const handleSubmit = async (values: any) => {
    if (!selectedProvider) {
      message.warning('请先选择一个供应商');
      return;
    }

    setSaving(true);
    try {
      const provider = new Provider({
        provider_name: values.providerName || selectedProvider.name,
        provider_type: selectedProvider.provider_type,
        base_url: values.baseUrl || selectedProvider.base_url,
        file_upload_base_url: selectedProvider.file_upload_base_url || null,
        api_key: values.apiKey || '',
        enable: values.enabled,
      });

      await Service.CompleteOnboarding(provider);
    } catch (error) {
      console.error('完成初始化失败:', error);
      message.error('保存失败，请检查配置后重试');
    } finally {
      setSaving(false);
    }
  };

  const handleExit = async () => {
    try {
      await Service.ExitApp();
    } catch (error) {
      console.error('关闭应用失败:', error);
      message.error('关闭应用失败');
    }
  };

  const handleStartConfig = () => {
    setStep('selectProvider');
  };

  const handleSelectProvider = (provider: SupportProvider) => {
    applyProviderPreset(provider);
    setStep('configProvider');
  };

  if (loading) {
    return (
      <div className={styles.loadingState}>
        <Spin size="large" />
        <Text type="secondary">正在准备欢迎引导...</Text>
      </div>
    );
  }

  return (
    <div className={styles.onboardingViewport}>
      <div className={`${styles.windowShell} ${step === 'intro' ? '' : styles.windowShellConfig}`}>
        <div
          className={`${styles.onboardingPage} ${step === 'intro' ? styles.onboardingPageIntro : styles.onboardingPageConfig}`}
        >
          {step === 'intro' ? (
            <div className={styles.heroPanel}>
              <div className={styles.heroBadge}>WELCOME</div>
              <Title level={1} className={styles.heroTitle}>
                先连接一个模型供应商，
                <br />
                再开始你的第一段对话
              </Title>
            </div>
          ) : null}

          <Card
            bordered={false}
            className={`${styles.formPanel} ${step === 'intro' ? '' : styles.formPanelConfig}`}
          >
            {step === 'intro' ? (
              <div className={styles.introPanel}>
                <div className={styles.introCards}>
                  <div className={styles.introCard}>
                    <span className={styles.introCardTitle}>开始前你会完成什么</span>
                    <span className={styles.introCardText}>选择一个模型供应商，并填写最少的连接信息。</span>
                  </div>
                  <div className={styles.introCard}>
                    <span className={styles.introCardTitle}>为什么需要这一步</span>
                    <span className={styles.introCardText}>应用需要一个可用模型入口，才能在主界面发起对话。</span>
                  </div>
                  <div className={styles.introCard}>
                    <span className={styles.introCardTitle}>之后还能改吗</span>
                    <span className={styles.introCardText}>可以，进入主界面后仍可在设置页继续调整供应商和模型。</span>
                  </div>
                </div>

                <Space className={styles.actions}>
                  <Button onClick={handleExit}>关闭应用</Button>
                  <Button
                    type="primary"
                    icon={<RocketOutlined />}
                    onClick={handleStartConfig}
                  >
                    开始配置
                  </Button>
                </Space>
              </div>
            ) : step === 'selectProvider' ? (
              <>
                <div className={styles.panelTopAction}>
                  <Space>
                    <Button
                      type="text"
                      icon={<ArrowLeftOutlined />}
                      className={styles.backButton}
                      onClick={() => setStep('intro')}
                    >
                      返回介绍
                    </Button>
                  </Space>
                </div>

                <div className={styles.selectPanel}>
                  <div className={styles.sectionHeader}>
                    <Title level={3}>选择供应商</Title>
                    <Text type="secondary">{selectProviderHint}</Text>
                  </div>

                  <div className={styles.selectPanelContent}>
                    {supportProviders.length > 0 ? (
                      <div className={styles.providerGrid}>
                        {supportProviders.map((provider) => (
                          <button
                            key={provider.provider_type || provider.name}
                            type="button"
                            className={styles.providerCard}
                            onClick={() => handleSelectProvider(provider)}
                          >
                            <Avatar
                              size={44}
                              src={iconSrc(provider.icon)}
                              style={{ backgroundColor: provider.icon ? 'transparent' : '#d77a2d' }}
                            >
                              {!provider.icon ? provider.name.slice(0, 1) : null}
                            </Avatar>
                            <div className={styles.providerMeta}>
                              <span className={styles.providerName}>{provider.name}</span>
                              <span className={styles.providerUrl}>{provider.base_url || '自定义 URL'}</span>
                            </div>
                          </button>
                        ))}
                      </div>
                    ) : (
                      <div className={styles.emptyState}>
                        <Empty
                          description={loadFailed ? '供应商列表加载失败，请稍后重试。' : '暂无可用的供应商预设'}
                        />
                      </div>
                    )}
                  </div>

                  <Space className={styles.selectPanelActions}>
                    <Button onClick={handleExit}>关闭应用</Button>
                  </Space>
                </div>
              </>
            ) : (
              <>
                <div className={styles.panelTopAction}>
                  <Space>
                    <Button
                      type="text"
                      icon={<ArrowLeftOutlined />}
                      className={styles.backButton}
                      onClick={() => setStep('selectProvider')}
                    >
                      返回选择供应商
                    </Button>
                  </Space>
                </div>

                <div className={styles.configPanel}>
                  <div className={styles.sectionHeader}>
                    <Title level={3}>配置 {selectedProvider?.name || '供应商'}</Title>
                    <Text type="secondary">{configProviderHint}</Text>
                  </div>

                  {selectedProvider ? (
                    <div className={styles.currentProviderCard}>
                      <Avatar
                        size={40}
                        src={iconSrc(selectedProvider.icon)}
                        style={{ backgroundColor: selectedProvider.icon ? 'transparent' : '#d77a2d' }}
                      >
                        {!selectedProvider.icon ? selectedProvider.name.slice(0, 1) : null}
                      </Avatar>
                      <div className={styles.currentProviderMeta}>
                        <span className={styles.currentProviderName}>{selectedProvider.name}</span>
                        <span className={styles.currentProviderDesc}>
                          {selectedProvider.base_url || '自定义 URL'}
                        </span>
                      </div>
                    </div>
                  ) : null}

                  <div className={styles.configPanelContent}>

                    <Alert
                      className={styles.notice}
                      type="info"
                      showIcon
                      message="API Key 将加密保存在本地。Ollama 等本地服务可以不填写 API Key。"
                    />

                    <Form
                      form={form}
                      layout="vertical"
                      initialValues={{ enabled: true }}
                      onFinish={handleSubmit}
                      className={styles.form}
                    >
                      <Form.Item label="启用供应商" name="enabled" valuePropName="checked">
                        <Switch />
                      </Form.Item>

                      <Form.Item
                        label="供应商名称"
                        name="providerName"
                        rules={[
                          { required: true, message: '请输入供应商名称' },
                          { max: 50, message: '供应商名称不能超过 50 个字符' },
                        ]}
                      >
                        <Input placeholder="例如：我的 DeepSeek" />
                      </Form.Item>

                      <Form.Item label="API Key" name="apiKey">
                        <Input.Password
                          placeholder="请输入 API Key，本地模型可留空"
                          iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                        />
                      </Form.Item>

                      <Form.Item
                        label="API 基础 URL"
                        name="baseUrl"
                        rules={[
                          { required: true, message: '请输入 API 基础 URL' },
                          { type: 'url', message: '请输入正确的 URL' },
                        ]}
                      >
                        <Input placeholder="https://api.example.com/v1" />
                      </Form.Item>

                      <Space className={styles.actions}>
                        <Button onClick={handleExit}>关闭应用</Button>
                        <Button
                          type="primary"
                          htmlType="submit"
                          icon={<RocketOutlined />}
                          loading={saving}
                          disabled={!selectedProvider || supportProviders.length === 0 || loadFailed}
                        >
                          保存并进入
                        </Button>
                      </Space>
                    </Form>
                  </div>
                </div>
              </>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
};

export default OnboardingPage;
