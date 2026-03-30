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
import { useTranslation } from 'react-i18next';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts';
import { Provider, SupportProvider } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/models.ts';
import styles from './index.module.scss';

const { Text, Title } = Typography;
type OnboardingStep = 'intro' | 'selectProvider' | 'configProvider';

const toAssetUrl = (path: string) => {
  const base = import.meta.env.BASE_URL || '/';
  const normalizedBase = base.endsWith('/') ? base : `${base}/`;
  return `${normalizedBase}${path.replace(/^\/+/, '')}`;
};

const OnboardingPage: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [step, setStep] = useState<OnboardingStep>('intro');
  const [supportProviders, setSupportProviders] = useState<SupportProvider[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<SupportProvider | null>(null);
  const [loadFailed, setLoadFailed] = useState(false);

  useEffect(() => {
    document.title = t('app.onboardingTitle');
  }, [t]);

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
        message.error(t('onboarding.loadFailed'));
      } finally {
        setLoading(false);
      }
    };

    void initialize();
  }, [navigate, t]);

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
      return t('onboarding.selectProviderLoadFailed');
    }
    if (supportProviders.length === 0) {
      return t('onboarding.noProviders');
    }
    return t('onboarding.selectHint');
  }, [loadFailed, supportProviders.length, t]);

  const configProviderHint = useMemo(() => {
    if (!selectedProvider) {
      return t('onboarding.selectProviderFirst');
    }
    return selectedProvider.description || t('onboarding.configHintFallback');
  }, [selectedProvider, t]);

  const iconSrc = (icon: string) => {
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

  const handleSubmit = async (values: any) => {
    if (!selectedProvider) {
      message.warning(t('onboarding.chooseProviderFirst'));
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
      message.error(t('onboarding.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleExit = async () => {
    try {
      await Service.ExitApp();
    } catch (error) {
      console.error('关闭应用失败:', error);
      message.error(t('onboarding.exitFailed'));
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
        <Text type="secondary">{t('onboarding.loading')}</Text>
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
              <div className={styles.heroBadge}>{t('onboarding.badge')}</div>
              <Title level={1} className={styles.heroTitle}>
                {t('onboarding.heroTitleLine1')}
                <br />
                {t('onboarding.heroTitleLine2')}
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
                    <span className={styles.introCardTitle}>{t('onboarding.intro.whatTitle')}</span>
                    <span className={styles.introCardText}>{t('onboarding.intro.whatText')}</span>
                  </div>
                  <div className={styles.introCard}>
                    <span className={styles.introCardTitle}>{t('onboarding.intro.whyTitle')}</span>
                    <span className={styles.introCardText}>{t('onboarding.intro.whyText')}</span>
                  </div>
                  <div className={styles.introCard}>
                    <span className={styles.introCardTitle}>{t('onboarding.intro.laterTitle')}</span>
                    <span className={styles.introCardText}>{t('onboarding.intro.laterText')}</span>
                  </div>
                </div>

                <Space className={styles.actions}>
                  <Button onClick={handleExit}>{t('onboarding.actions.exit')}</Button>
                  <Button
                    type="primary"
                    icon={<RocketOutlined />}
                    onClick={handleStartConfig}
                  >
                    {t('onboarding.actions.start')}
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
                      {t('onboarding.actions.backToIntro')}
                    </Button>
                  </Space>
                </div>

                <div className={styles.selectPanel}>
                  <div className={styles.sectionHeader}>
                    <Title level={3}>{t('onboarding.selectProvider')}</Title>
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
                              style={{ backgroundColor: '#d77a2d' }}
                            >
                              {provider.name.slice(0, 1)}
                            </Avatar>
                            <div className={styles.providerMeta}>
                              <span className={styles.providerName}>{provider.name}</span>
                              <span className={styles.providerUrl}>{provider.base_url || t('onboarding.customUrl')}</span>
                            </div>
                          </button>
                        ))}
                      </div>
                    ) : (
                      <div className={styles.emptyState}>
                        <Empty
                          description={loadFailed ? t('onboarding.selectProviderLoadFailed') : t('onboarding.emptyProviders')}
                        />
                      </div>
                    )}
                  </div>

                  <Space className={styles.selectPanelActions}>
                    <Button onClick={handleExit}>{t('onboarding.actions.exit')}</Button>
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
                      {t('onboarding.actions.backToProviders')}
                    </Button>
                  </Space>
                </div>

                <div className={styles.configPanel}>
                  <div className={styles.sectionHeader}>
                    <Title level={3}>{t('onboarding.configProvider', { name: selectedProvider?.name || t('onboarding.configFallbackName') })}</Title>
                    <Text type="secondary">{configProviderHint}</Text>
                  </div>

                  {selectedProvider ? (
                    <div className={styles.currentProviderCard}>
                      <Avatar
                        size={40}
                        src={iconSrc(selectedProvider.icon)}
                        style={{ backgroundColor: '#d77a2d' }}
                      >
                        {selectedProvider.name.slice(0, 1)}
                      </Avatar>
                      <div className={styles.currentProviderMeta}>
                        <span className={styles.currentProviderName}>{selectedProvider.name}</span>
                        <span className={styles.currentProviderDesc}>
                          {selectedProvider.base_url || t('onboarding.customUrl')}
                        </span>
                      </div>
                    </div>
                  ) : null}

                  <div className={styles.configPanelContent}>

                    <Alert
                      className={styles.notice}
                      type="info"
                      showIcon
                      message={t('onboarding.securityAlert')}
                    />

                    <Form
                      form={form}
                      layout="vertical"
                      initialValues={{ enabled: true }}
                      onFinish={handleSubmit}
                      className={styles.form}
                    >
                      <Form.Item label={t('onboarding.fields.enabled')} name="enabled" valuePropName="checked">
                        <Switch />
                      </Form.Item>

                      <Form.Item
                        label={t('onboarding.fields.providerName')}
                        name="providerName"
                        rules={[
                          { required: true, message: t('onboarding.validation.providerNameRequired') },
                          { max: 50, message: t('onboarding.validation.providerNameMax') },
                        ]}
                      >
                        <Input placeholder={t('onboarding.placeholders.providerName')} />
                      </Form.Item>

                      <Form.Item label={t('onboarding.fields.apiKey')} name="apiKey">
                        <Input.Password
                          placeholder={t('onboarding.placeholders.apiKey')}
                          iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                        />
                      </Form.Item>

                      <Form.Item
                        label={t('onboarding.fields.baseUrl')}
                        name="baseUrl"
                        rules={[
                          { required: true, message: t('onboarding.validation.baseUrlRequired') },
                          { type: 'url', message: t('onboarding.validation.invalidUrl') },
                        ]}
                      >
                        <Input placeholder={t('onboarding.placeholders.baseUrl')} />
                      </Form.Item>

                      <Space className={styles.actions}>
                        <Button onClick={handleExit}>{t('onboarding.actions.exit')}</Button>
                        <Button
                          type="primary"
                          htmlType="submit"
                          icon={<RocketOutlined />}
                          loading={saving}
                          disabled={!selectedProvider || supportProviders.length === 0 || loadFailed}
                        >
                          {t('onboarding.actions.saveAndEnter')}
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
