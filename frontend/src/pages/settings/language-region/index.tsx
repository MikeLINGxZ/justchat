import React from 'react';
import { Card, Typography, Select, message } from 'antd';
import { GlobalOutlined, EnvironmentOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useLanguageStore } from '@/stores/languageStore';
import { translateError } from '@/utils/errorHandler';
import { LANGUAGE_OPTIONS, REGION_OPTIONS } from '@/i18n/types';
import type { AppLanguage, AppRegion } from '@/i18n/types';
import styles from './index.module.scss';

const { Title, Text } = Typography;

const LanguageRegionSettingsPage: React.FC = () => {
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
    <div className={styles.languageRegionSettings}>
      <div className={styles.pageHeader}>
        <Title level={3}>{t('settings.languageRegion.title')}</Title>
        <Text type="secondary">{t('settings.languageRegion.description')}</Text>
      </div>

      <div className={styles.content}>
        <Card title={t('settings.languageRegion.languageTitle')} className={styles.settingCard}>
          <div className={styles.settingItem}>
            <div className={styles.settingLabel}>
              <div className={styles.iconTitle}>
                <GlobalOutlined className={styles.settingIcon} />
                <Title level={5}>{t('settings.languageRegion.languageLabel')}</Title>
              </div>
              <Text type="secondary">{t('settings.languageRegion.languageDescription')}</Text>
            </div>
            <div className={styles.settingControl}>
              <Select
                value={language}
                onChange={handleLanguageChange}
                options={LANGUAGE_OPTIONS.map((option) => ({
                  value: option.value,
                  label: option.nativeLabel,
                }))}
                style={{ width: '100%', maxWidth: 320 }}
              />
              <Text type="secondary">{t('settings.languageRegion.languageHint')}</Text>
            </div>
          </div>
        </Card>

        <Card title={t('settings.languageRegion.regionTitle')} className={styles.settingCard}>
          <div className={styles.settingItem}>
            <div className={styles.settingLabel}>
              <div className={styles.iconTitle}>
                <EnvironmentOutlined className={styles.settingIcon} />
                <Title level={5}>{t('settings.languageRegion.regionLabel')}</Title>
              </div>
              <Text type="secondary">{t('settings.languageRegion.regionDescription')}</Text>
            </div>
            <div className={styles.settingControl}>
              <Select
                value={region}
                onChange={handleRegionChange}
                options={REGION_OPTIONS.map((option) => ({
                  value: option.value,
                  label: t(option.labelKey),
                }))}
                style={{ width: '100%', maxWidth: 320 }}
              />
              <Text type="secondary">{t('settings.languageRegion.regionHint')}</Text>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default LanguageRegionSettingsPage;
