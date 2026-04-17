import React, { useCallback, useEffect, useState } from 'react';
import {
  Alert,
  Button,
  Form,
  Input,
  InputNumber,
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
  { value: 'event', labelKey: 'settings.memory.typeEvent' },
  { value: 'skill', labelKey: 'settings.memory.typeSkill' },
  { value: 'plan', labelKey: 'settings.memory.typePlan' },
];

interface MemoryEditFormValues {
  summary: string;
  content: string;
  type: string;
  time_range_start: string;
  time_range_end: string;
  location: string;
  characters: string;
  importance: number;
  emotional_valence: number;
}

const toDateInputValue = (dateValue: unknown) => {
  if (!dateValue) return '';
  const raw = String(dateValue);
  const dateOnly = raw.slice(0, 10);
  return /^\d{4}-\d{2}-\d{2}$/.test(dateOnly) ? dateOnly : '';
};

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
          type: m.type ?? '',
          time_range_start: toDateInputValue(m.time_range_start),
          time_range_end: toDateInputValue(m.time_range_end),
          location: m.location ?? '',
          characters: m.characters ?? '',
          importance: typeof m.importance === 'number' ? m.importance : 0.5,
          emotional_valence: typeof m.emotional_valence === 'number' ? m.emotional_valence : 0,
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
        time_range_start: values.time_range_start || null,
        time_range_end: values.time_range_end || null,
        location: values.location.trim() || null,
        characters: values.characters.trim() || null,
        importance: values.importance,
        emotional_valence: values.emotional_valence,
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
              <TextArea rows={6} maxLength={2000} showCount />
            </Form.Item>
            <Form.Item name="type" label={t('settings.memory.fieldType')}>
              <Select
                options={MEMORY_TYPES.map(mt => ({ value: mt.value, label: t(mt.labelKey) }))}
                allowClear
              />
            </Form.Item>
            <div className={styles.editGrid}>
              <Form.Item name="time_range_start" label={t('settings.memory.fieldStartDate')}>
                <Input type="date" />
              </Form.Item>
              <Form.Item name="time_range_end" label={t('settings.memory.fieldEndDate')}>
                <Input type="date" />
              </Form.Item>
            </div>
            <div className={styles.editGrid}>
              <Form.Item name="location" label={t('settings.memory.fieldLocation')}>
                <Input maxLength={120} />
              </Form.Item>
              <Form.Item name="characters" label={t('settings.memory.fieldCharacters')}>
                <Input maxLength={120} />
              </Form.Item>
            </div>
            <div className={styles.editGrid}>
              <Form.Item
                name="importance"
                label={t('settings.memory.fieldImportance')}
                rules={[{ required: true, message: t('settings.memory.importanceRequired') }]}
              >
                <InputNumber min={0} max={1} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item
                name="emotional_valence"
                label={t('settings.memory.fieldEmotion')}
                rules={[{ required: true, message: t('settings.memory.emotionRequired') }]}
              >
                <InputNumber min={-1} max={1} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
            </div>
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
