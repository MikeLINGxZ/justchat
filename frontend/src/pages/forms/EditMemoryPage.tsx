import React, { useCallback, useEffect, useState } from 'react';
import {
  Alert,
  Button,
  Form,
  Input,
  message,
  Select,
  Spin,
  Typography,
} from 'antd';
import { SaveOutlined } from '@ant-design/icons';
import { useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import {
  Memory,
  MemoryUpdateInput,
} from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import { useLabStore } from '@/stores/labStore';
import { translateError } from '@/utils/errorHandler';
import styles from './formWindow.module.scss';

const WINDOW_NAME_PREFIX = 'window_form_memory';
const EVENT_KEY = 'settings:memories:changed';

const { Text } = Typography;
const { TextArea } = Input;

const MEMORY_TYPES = [
  { value: 'fact', labelKey: 'settings.memory.typeFact' },
  { value: 'information', labelKey: 'settings.memory.typeInformation' },
  { value: 'event', labelKey: 'settings.memory.typeEvent' },
];

const MEMORY_TARGETS = [
  { value: 'user', labelKey: 'settings.memory.targetUser' },
  { value: 'memory', labelKey: 'settings.memory.targetMemory' },
];

interface MemoryEditFormValues {
  summary: string;
  content: string;
  type: string;
  target: string;
}

const EditMemoryPage: React.FC = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const idParam = searchParams.get('id');
  const memoryId = idParam ? Number(idParam) : NaN;

  const [form] = Form.useForm<MemoryEditFormValues>();
  const { vectorSearchEnabled } = useLabStore();

  const [memory, setMemory] = useState<Memory | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string>('');
  const [saving, setSaving] = useState(false);

  const windowName = Number.isNaN(memoryId) ? WINDOW_NAME_PREFIX : `${WINDOW_NAME_PREFIX}_${memoryId}`;

  useEffect(() => {
    document.title = t('settings.memory.editTitle');
  }, [t]);

  useEffect(() => {
    if (Number.isNaN(memoryId)) {
      setLoading(false);
      setLoadError(t('settings.memory.invalidId', { defaultValue: 'Invalid memory id' }));
      return;
    }
    void (async () => {
      setLoading(true);
      try {
        const m = await Service.GetMemoryDetail(memoryId);
        if (!m) {
          setLoadError(t('settings.memory.notFound', { defaultValue: 'Memory not found' }));
          return;
        }
        setMemory(m);
        form.setFieldsValue({
          summary: m.summary ?? '',
          content: m.content ?? '',
          type: (m.type ?? '').trim(),
          target: (m.target ?? 'user').trim() || 'user',
        });
      } catch (err) {
        setLoadError(translateError(err, t('settings.memory.loadFailed', { defaultValue: 'Load failed' })));
      } finally {
        setLoading(false);
      }
    })();
  }, [memoryId, form, t]);

  const handleCancel = useCallback(() => {
    void Service.CloseFormWindow(windowName);
  }, [windowName]);

  const handleSubmit = useCallback(async () => {
    if (!memory) return;
    try {
      const values = await form.validateFields();
      setSaving(true);
      const payload = new MemoryUpdateInput({
        summary: values.summary.trim(),
        content: values.content.trim(),
        type: values.type ?? '',
        target: values.target || 'user',
      });
      await Service.UpdateMemory(memory.id, payload);
      void Events.Emit(EVENT_KEY, { id: memory.id });
      message.success(t('settings.memory.saveSuccess'));
      void Service.CloseFormWindow(windowName);
    } catch (error) {
      if ((error as { errorFields?: unknown }).errorFields) {
        return;
      }
      message.error(translateError(error, t('settings.memory.saveFailed')));
    } finally {
      setSaving(false);
    }
  }, [memory, form, t, windowName]);

  return (
    <div className={styles.formWindow}>
      <div className={styles.header}>
        <h2>{t('settings.memory.editTitle')}</h2>
        {memory?.summary && <div className={styles.subtitle}>{memory.summary}</div>}
      </div>

      <div className={styles.body}>
        {loading ? (
          <div style={{ textAlign: 'center', padding: 48 }}>
            <Spin />
          </div>
        ) : loadError ? (
          <Alert type="error" showIcon message={loadError} />
        ) : (
          <Form form={form} layout="vertical">
            <Form.Item name="summary" label={t('settings.memory.fieldSummary')}>
              <Input maxLength={120} />
            </Form.Item>
            <Form.Item
              name="content"
              label={t('settings.memory.fieldContent')}
              rules={[{ required: true, message: t('settings.memory.contentRequired') }]}
            >
              <TextArea rows={10} maxLength={4000} showCount />
            </Form.Item>
            <Form.Item name="type" label={t('settings.memory.fieldType')}>
              <Select
                options={MEMORY_TYPES.map(mt => ({ value: mt.value, label: t(mt.labelKey) }))}
                allowClear
              />
            </Form.Item>
            <Form.Item name="target" label={t('settings.memory.fieldTarget')}>
              <Select
                options={MEMORY_TARGETS.map(mt => ({ value: mt.value, label: t(mt.labelKey) }))}
              />
            </Form.Item>
            <Text type="secondary" className={styles.hintText}>
              {t('settings.memory.capacityHint', {
                count: memory?.char_count ?? 0,
                defaultValue: `${memory?.char_count ?? 0} chars in this entry`,
              })}
            </Text>
            {vectorSearchEnabled && (
              <Text type="secondary" className={styles.hintText}>
                {t('settings.memory.reembedHint')}
              </Text>
            )}
          </Form>
        )}
      </div>

      <div className={styles.footer}>
        <Button onClick={handleCancel}>{t('common.cancel')}</Button>
        <Button
          type="primary"
          icon={<SaveOutlined />}
          loading={saving}
          disabled={!memory}
          onClick={() => void handleSubmit()}
        >
          {t('settings.memory.actions.save', { defaultValue: t('common.confirm') })}
        </Button>
      </div>
    </div>
  );
};

export default EditMemoryPage;
