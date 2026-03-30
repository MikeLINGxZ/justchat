import React, { useState } from 'react';
import { Card, Typography, Slider, Button, message } from 'antd';
import { ReloadOutlined, CheckOutlined, FontSizeOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useFontSizeStore, FONT_SIZE_OPTIONS, FONT_SIZE_OFFSETS } from '@/stores/fontSizeStore';
import styles from './index.module.scss';

const { Title, Text } = Typography;

const GeneralSettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const { fontSizeOffset, setFontSizeOffset } = useFontSizeStore();
  
  const [previewOffset, setPreviewOffset] = useState(fontSizeOffset);
  const [hasChanges, setHasChanges] = useState(false);

  const handlePreviewChange = (value: number) => {
    setPreviewOffset(value as any);
    setHasChanges(value !== fontSizeOffset);
  };

  const handleApplySettings = () => {
    setFontSizeOffset(previewOffset as any);
    setHasChanges(false);
    message.success(t('settings.general.applied'));
  };

  const handleResetSettings = () => {
    setPreviewOffset(FONT_SIZE_OFFSETS.NORMAL);
    setHasChanges(FONT_SIZE_OFFSETS.NORMAL !== fontSizeOffset);
  };

  const getFontSizeLabel = (offset: number) => t(`settings.general.fontSizes.${offset}`);

  const sliderMarks = FONT_SIZE_OPTIONS.reduce((marks, option) => {
    marks[option.value] = {
      style: { fontSize: '11px', color: 'var(--text-color-secondary)' },
      label: getFontSizeLabel(option.value),
    };
    return marks;
  }, {} as any);

  return (
    <div className={styles.generalSettings}>
      <div className={styles.content}>
        <Card title={t('settings.general.displayTitle')} className={styles.settingCard}>
          <div className={styles.settingItem}>
            <div className={styles.settingLabel}>
              <Title level={5}>{t('settings.general.fontSizeTitle')}</Title>
              <Text type="secondary">{t('settings.general.fontSizeDescription')}</Text>
            </div>
            
            <div className={styles.fontSizeControl}>
              <div className={styles.sliderContainer}>
                <div className={styles.sliderHeader}>
                  <Text strong>{t('settings.general.currentSize')}</Text>
                  <div className={styles.currentSize}>
                    <Text strong className={styles.sizeLabel}>{getFontSizeLabel(previewOffset)}</Text>
                    <Text type="secondary" className={styles.sizeDesc}>
                      {14 + previewOffset}px
                    </Text>
                  </div>
                </div>
                <Slider
                  min={FONT_SIZE_OFFSETS.VERY_SMALL}
                  max={FONT_SIZE_OFFSETS.EXTRA_LARGE}
                  step={2}
                  value={previewOffset}
                  onChange={handlePreviewChange}
                  marks={sliderMarks}
                />
              </div>

              <div className={styles.presetOptions}>
                <Text strong>{t('settings.general.quickSelect')}</Text>
                <div className={styles.presetButtons}>
                  {FONT_SIZE_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      className={`${styles.presetButton} ${previewOffset === option.value ? styles.active : ''}`}
                      onClick={() => handlePreviewChange(option.value)}
                    >
                      <span className={styles.buttonLabel}>{getFontSizeLabel(option.value)}</span>
                      <span className={styles.buttonSize}>{option.description}</span>
                    </button>
                  ))}
                </div>
              </div>

              <div className={styles.previewArea}>
                <div className={styles.previewHeader}>
                  <FontSizeOutlined className={styles.previewIcon} />
                  <Title level={5}>{t('settings.general.previewTitle')}</Title>
                </div>
                <div 
                  className={styles.previewContent}
                  style={{
                    fontSize: `${14 + previewOffset}px`,
                    lineHeight: 1.5715 + (previewOffset > 0 ? -0.05 : previewOffset < 0 ? 0.05 : 0)
                  }}
                >
                  <div className={styles.previewText}>
                    {t('settings.general.previewText')}
                  </div>
                  <div className={styles.previewSmall} style={{ fontSize: `${12 + previewOffset}px` }}>
                    {t('settings.general.previewSmall')}
                  </div>
                  <div className={styles.previewBold} style={{ fontSize: `${14 + previewOffset}px`, fontWeight: 600 }}>
                    {t('settings.general.previewBold')}
                  </div>
                  <div className={styles.previewSecondary} style={{ fontSize: `${14 + previewOffset}px`, opacity: 0.65 }}>
                    {t('settings.general.previewSecondary')}
                  </div>
                </div>
              </div>

              <div className={styles.actions}>
                <Button 
                  icon={<ReloadOutlined />}
                  onClick={handleResetSettings}
                  disabled={previewOffset === FONT_SIZE_OFFSETS.NORMAL}
                  className={styles.resetButton}
                >
                  {t('settings.general.reset')}
                </Button>
                <Button 
                  type="primary" 
                  icon={<CheckOutlined />}
                  onClick={handleApplySettings}
                  disabled={!hasChanges}
                  className={styles.applyButton}
                >
                  {t('settings.general.apply')}
                </Button>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default GeneralSettingsPage;
