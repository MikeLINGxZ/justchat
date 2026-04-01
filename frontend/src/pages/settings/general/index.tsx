import React, { useState } from 'react';
import { Button, Card, message, Select, Slider, Typography } from 'antd';
import {
  CheckOutlined,
  FontSizeOutlined,
  GlobalOutlined,
  EnvironmentOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { useTranslation } from 'react-i18next';
import { useFontSizeStore, FONT_SIZE_OPTIONS, FONT_SIZE_OFFSETS } from '@/stores/fontSizeStore';
import { useLanguageStore } from '@/stores/languageStore';
import { translateError } from '@/utils/errorHandler';
import { LANGUAGE_OPTIONS, REGION_OPTIONS } from '@/i18n/types';
import type { AppLanguage, AppRegion } from '@/i18n/types';
import styles from './index.module.scss';

const { Title, Text } = Typography;

type SettingSection = 'display' | 'language-region';

const GeneralSettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState<SettingSection>('display');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const [showDetailOnMobile, setShowDetailOnMobile] = useState(false);

  React.useEffect(() => {
    const handleResize = () => setIsMobile(isMobileDevice());
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const sections: { key: SettingSection; title: string; icon: React.ReactNode }[] = [
    { key: 'display', title: t('settings.general.menuDisplay'), icon: <FontSizeOutlined /> },
    { key: 'language-region', title: t('settings.general.menuLanguageRegion'), icon: <GlobalOutlined /> },
  ];

  const handleSelectSection = (key: SettingSection) => {
    setActiveSection(key);
    if (isMobile) setShowDetailOnMobile(true);
  };

  const renderSectionList = () => (
    <Card className={styles.listCard} title={t('settings.general.listTitle')}>
      <div className={styles.sectionList}>
        {sections.map(section => (
          <button
            key={section.key}
            type="button"
            className={`${styles.sectionItem} ${activeSection === section.key ? styles.selected : ''}`}
            onClick={() => handleSelectSection(section.key)}
          >
            <span className={styles.sectionIcon}>{section.icon}</span>
            <span className={styles.sectionTitle}>{section.title}</span>
          </button>
        ))}
      </div>
    </Card>
  );

  const renderDetail = () => (
    <Card className={styles.detailCard}>
      {activeSection === 'display' ? <DisplaySettings /> : <LanguageRegionSettings />}
    </Card>
  );

  return (
    <div className={styles.generalSettings}>
      {isMobile ? (
        <>
          {!showDetailOnMobile && renderSectionList()}
          {showDetailOnMobile && (
            <div className={styles.mobileDetail}>
              <Button
                type="text"
                className={styles.mobileBackButton}
                onClick={() => setShowDetailOnMobile(false)}
              >
                {t('settings.back')}
              </Button>
              {renderDetail()}
            </div>
          )}
        </>
      ) : (
        <div className={styles.desktopLayout}>
          <div className={styles.listColumn}>{renderSectionList()}</div>
          <div className={styles.detailColumn}>{renderDetail()}</div>
        </div>
      )}
    </div>
  );
};

// ---- Display Settings ----

const DisplaySettings: React.FC = () => {
  const { t } = useTranslation();
  const { fontSizeOffset, setFontSizeOffset } = useFontSizeStore();
  const [previewOffset, setPreviewOffset] = useState(fontSizeOffset);
  const [hasChanges, setHasChanges] = useState(false);

  const handlePreviewChange = (value: number) => {
    setPreviewOffset(value as any);
    setHasChanges(value !== fontSizeOffset);
  };

  const handleApply = () => {
    setFontSizeOffset(previewOffset as any);
    setHasChanges(false);
    message.success(t('settings.general.applied'));
  };

  const handleReset = () => {
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
    <div className={styles.settingContent}>
      <div className={styles.settingHeader}>
        <Title level={4}>{t('settings.general.displayTitle')}</Title>
        <Text type="secondary">{t('settings.general.fontSizeDescription')}</Text>
      </div>

      <div className={styles.sliderSection}>
        <div className={styles.sliderHeader}>
          <Text strong>{t('settings.general.currentSize')}</Text>
          <div className={styles.currentSize}>
            <Text strong>{getFontSizeLabel(previewOffset)}</Text>
            <Text type="secondary">{14 + previewOffset}px</Text>
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

      <div className={styles.presetSection}>
        <Text strong>{t('settings.general.quickSelect')}</Text>
        <div className={styles.presetButtons}>
          {FONT_SIZE_OPTIONS.map(option => (
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
          <FontSizeOutlined />
          <Text strong>{t('settings.general.previewTitle')}</Text>
        </div>
        <div
          className={styles.previewContent}
          style={{
            fontSize: `${14 + previewOffset}px`,
            lineHeight: 1.5715 + (previewOffset > 0 ? -0.05 : previewOffset < 0 ? 0.05 : 0),
          }}
        >
          <div>{t('settings.general.previewText')}</div>
          <div style={{ fontSize: `${12 + previewOffset}px` }}>{t('settings.general.previewSmall')}</div>
          <div style={{ fontWeight: 600 }}>{t('settings.general.previewBold')}</div>
          <div style={{ opacity: 0.65 }}>{t('settings.general.previewSecondary')}</div>
        </div>
      </div>

      <div className={styles.actions}>
        <Button icon={<ReloadOutlined />} onClick={handleReset} disabled={previewOffset === FONT_SIZE_OFFSETS.NORMAL}>
          {t('settings.general.reset')}
        </Button>
        <Button type="primary" icon={<CheckOutlined />} onClick={handleApply} disabled={!hasChanges}>
          {t('settings.general.apply')}
        </Button>
      </div>
    </div>
  );
};

// ---- Language & Region Settings ----

const LanguageRegionSettings: React.FC = () => {
  const { t } = useTranslation();
  const { language, region, setLanguage, setRegion } = useLanguageStore();

  const handleLanguageChange = async (nextLanguage: AppLanguage) => {
    try {
      await setLanguage(nextLanguage);
      void message.success(t('settings.languageRegion.languageChanged'));
    } catch (error) {
      void message.error(translateError(error));
    }
  };

  const handleRegionChange = async (nextRegion: AppRegion) => {
    try {
      await setRegion(nextRegion);
      void message.success(t('settings.languageRegion.regionChanged'));
    } catch (error) {
      void message.error(translateError(error));
    }
  };

  return (
    <div className={styles.settingContent}>
      <div className={styles.settingHeader}>
        <Title level={4}>{t('settings.languageRegion.title')}</Title>
        <Text type="secondary">{t('settings.languageRegion.description')}</Text>
      </div>

      <div className={styles.formSection}>
        <div className={styles.formItem}>
          <div className={styles.formLabel}>
            <div className={styles.formLabelIcon}>
              <GlobalOutlined />
              <Text strong>{t('settings.languageRegion.languageLabel')}</Text>
            </div>
            <Text type="secondary">{t('settings.languageRegion.languageDescription')}</Text>
          </div>
          <Select
            value={language}
            onChange={handleLanguageChange}
            options={LANGUAGE_OPTIONS.map(option => ({
              value: option.value,
              label: option.nativeLabel,
            }))}
            style={{ width: '100%', maxWidth: 320 }}
          />
          <Text type="secondary" className={styles.formHint}>
            {t('settings.languageRegion.languageHint')}
          </Text>
        </div>

        <div className={styles.formItem}>
          <div className={styles.formLabel}>
            <div className={styles.formLabelIcon}>
              <EnvironmentOutlined />
              <Text strong>{t('settings.languageRegion.regionLabel')}</Text>
            </div>
            <Text type="secondary">{t('settings.languageRegion.regionDescription')}</Text>
          </div>
          <Select
            value={region}
            onChange={handleRegionChange}
            options={REGION_OPTIONS.map(option => ({
              value: option.value,
              label: t(option.labelKey),
            }))}
            style={{ width: '100%', maxWidth: 320 }}
          />
          <Text type="secondary" className={styles.formHint}>
            {t('settings.languageRegion.regionHint')}
          </Text>
        </div>
      </div>
    </div>
  );
};

export default GeneralSettingsPage;
