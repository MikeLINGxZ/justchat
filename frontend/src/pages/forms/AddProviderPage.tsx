import React, { useEffect, useState } from 'react';
import {
  Alert,
  Avatar,
  Button,
  Divider,
  Empty,
  Form,
  Input,
  message,
  Select,
  Space,
  Spin,
  Switch,
  Tooltip,
} from 'antd';
import {
  ApiOutlined,
  ArrowLeftOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
  PlusOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import {
  Provider,
  SupportProvider,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import { ProviderType } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models';
import styles from './formWindow.module.scss';

const WINDOW_NAME = 'window_form_provider';
const EVENT_KEY = 'settings:providers:changed';

const { Option } = Select;

const toAssetUrl = (path: string) => {
  const base = import.meta.env.BASE_URL || '/';
  const normalizedBase = base.endsWith('/') ? base : `${base}/`;
  return `${normalizedBase}${path.replace(/^\/+/, '')}`;
};

const resolveIconSrc = (icon?: string) => {
  if (!icon) return undefined;
  if (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('data:')) {
    return icon;
  }
  if (icon.startsWith('/')) return toAssetUrl(icon);
  if (/\.(png|jpe?g|gif|webp|svg|ico)$/i.test(icon)) return toAssetUrl(icon);
  return `data:image/png;base64,${icon}`;
};

const AddProviderPage: React.FC = () => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [supportProviders, setSupportProviders] = useState<SupportProvider[]>([]);
  const [loadingList, setLoadingList] = useState(false);
  const [selected, setSelected] = useState<SupportProvider | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    document.title = t('settings.provider.modal.selectProvider');
    void (async () => {
      setLoadingList(true);
      try {
        const providers = await Service.GetSupportProviders();
        setSupportProviders(providers || []);
      } catch (err) {
        console.error('加载支持的供应商列表失败:', err);
        message.error(t('settings.provider.messages.loadSupportFailed'));
      } finally {
        setLoadingList(false);
      }
    })();
  }, [t]);

  const handleSelect = (sp: SupportProvider) => {
    setSelected(sp);
    form.setFieldsValue({
      enabled: true,
      providerName: sp.name,
      apiKey: '',
      baseUrl: sp.base_url,
      fileUploadBaseUrl: sp.file_upload_base_url || '',
      defaultModel: undefined,
    });
  };

  const handleBackToList = () => {
    setSelected(null);
    form.resetFields();
  };

  const handleCancel = () => {
    void Service.CloseFormWindow(WINDOW_NAME);
  };

  const handleSubmit = async (values: any) => {
    if (!selected) return;
    setSubmitting(true);
    try {
      const data = new Provider({
        provider_name: values.providerName || selected.name,
        provider_type: selected.provider_type as ProviderType,
        base_url: values.baseUrl || selected.base_url,
        file_upload_base_url: values.fileUploadBaseUrl || null,
        api_key: values.apiKey || '',
        enable: values.enabled,
        default_model_id: values.defaultModel || 0,
      });
      await Service.AddProvider(data);
      void Events.Emit(EVENT_KEY, null);
      message.success(t('settings.provider.messages.createSuccess'));
      void Service.CloseFormWindow(WINDOW_NAME);
    } catch (err) {
      console.error('创建供应商失败:', err);
      message.error(t('settings.provider.messages.createFailed'));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className={styles.formWindow}>
      <div className={styles.header}>
        {selected ? (
          <Space align="center">
            <Button type="text" icon={<ArrowLeftOutlined />} onClick={handleBackToList} />
            <h2 style={{ display: 'inline' }}>
              {t('settings.provider.modal.addProvider', { name: selected.name })}
            </h2>
          </Space>
        ) : (
          <>
            <h2>{t('settings.provider.modal.selectProvider')}</h2>
          </>
        )}
      </div>

      <div className={styles.body}>
        {!selected ? (
          loadingList ? (
            <div style={{ textAlign: 'center', padding: 48 }}>
              <Spin />
            </div>
          ) : supportProviders.length === 0 ? (
            <Empty description={t('settings.provider.modal.emptyProviders')} />
          ) : (
            <div className={styles.supportProviderList}>
              {supportProviders.map(item => {
                const iconSrc = resolveIconSrc(item.icon);
                return (
                  <div
                    key={item.provider_type || item.name}
                    className={styles.supportProviderItem}
                    onClick={() => handleSelect(item)}
                  >
                    <Avatar
                      size={40}
                      src={iconSrc}
                      style={{
                        backgroundColor: item.icon ? 'transparent' : 'var(--primary-color-light)',
                        fontSize: 20,
                        flexShrink: 0,
                      }}
                    >
                      {!item.icon && (item.name || '?').charAt(0)}
                    </Avatar>
                    <div className={styles.supportProviderInfo}>
                      <span className={styles.supportProviderName}>{item.name}</span>
                      {item.description && (
                        <span className={styles.supportProviderDesc}>{item.description}</span>
                      )}
                      {item.base_url && (
                        <span className={styles.supportProviderUrl}>
                          <ApiOutlined style={{ marginRight: 4 }} />
                          {item.base_url}
                        </span>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )
        ) : (
          <Form form={form} layout="vertical" onFinish={handleSubmit}>
            <Alert
              message={t('settings.provider.securityAlert')}
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
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

            <Form.Item label={t('settings.provider.fields.apiKey')} name="apiKey">
              <Input.Password
                placeholder={t('settings.provider.placeholders.apiKey')}
                iconRender={visible => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
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
              rules={[{ type: 'url', message: t('settings.provider.validation.invalidUrl') }]}
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
                disabled
                notFoundContent={t('settings.provider.helper.saveProviderFirst')}
              >
                <Option value={0} disabled>
                  {t('settings.provider.helper.saveProviderFirst')}
                </Option>
              </Select>
            </Form.Item>

            <Divider />
          </Form>
        )}
      </div>

      <div className={styles.footer}>
        <Button onClick={handleCancel}>{t('common.cancel')}</Button>
        {selected && (
          <Button
            type="primary"
            icon={<PlusOutlined />}
            loading={submitting}
            onClick={() => form.submit()}
          >
            {t('settings.provider.actions.create')}
          </Button>
        )}
      </div>
    </div>
  );
};

export default AddProviderPage;
