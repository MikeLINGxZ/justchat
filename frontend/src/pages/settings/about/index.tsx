import React, { useCallback, useEffect, useRef, useState } from 'react';
import { Typography, message } from 'antd';
import {
  InfoCircleOutlined,
  AppstoreOutlined,
  MessageOutlined,
  GithubOutlined,
  MailOutlined,
  LinkOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { isMobileDevice } from '@/hooks/useViewportHeight';
import { Browser } from '@wailsio/runtime';
import styles from './index.module.scss';

const { Title, Text, Paragraph } = Typography;

const DiscordIcon = () => (
  <svg viewBox="0 0 24 24" width="1em" height="1em" fill="currentColor">
    <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z" />
  </svg>
);

const QQIcon = () => (
  <svg viewBox="0 0 24 24" width="1em" height="1em" fill="currentColor">
    <path d="M21.395 15.035a39.548 39.548 0 0 0-1.51-3.588c.142-.678.216-1.395.216-2.147 0-4.97-3.549-9-7.913-9C7.822.3 4.274 4.33 4.274 9.3c0 .752.074 1.47.216 2.147a39.548 39.548 0 0 0-1.51 3.588c-.346.919-.578 1.736-.578 2.358 0 .578.164.96.542 1.152.376.192.876.048 1.497-.35.311-.198.636-.45.96-.737.04.054.083.107.127.16a8.462 8.462 0 0 0 1.658 1.496c-.247.36-.4.733-.4 1.11 0 .374.157.727.488.973.487.362 1.277.504 2.876.504.934 0 1.717-.074 2.358-.192.641.118 1.424.192 2.358.192 1.599 0 2.389-.142 2.876-.504.331-.246.488-.599.488-.973 0-.377-.153-.75-.4-1.11a8.462 8.462 0 0 0 1.658-1.496c.044-.053.087-.106.127-.16.324.288.649.539.96.737.621.398 1.12.542 1.497.35.378-.192.542-.574.542-1.152 0-.622-.232-1.44-.578-2.358z" />
  </svg>
);

type SectionKey = 'about' | 'features' | 'contact';

const SECTIONS: { key: SectionKey; icon: React.ReactNode }[] = [
  { key: 'about', icon: <InfoCircleOutlined /> },
  { key: 'features', icon: <AppstoreOutlined /> },
  { key: 'contact', icon: <MessageOutlined /> },
];

const AboutPage: React.FC = () => {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState<SectionKey>('about');
  const [isMobile, setIsMobile] = useState(() => isMobileDevice());
  const contentRef = useRef<HTMLDivElement>(null);
  const sectionRefs = useRef<Record<SectionKey, HTMLDivElement | null>>({
    about: null,
    features: null,
    contact: null,
  });

  useEffect(() => {
    const handleResize = () => setIsMobile(isMobileDevice());
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Scroll spy
  const handleScroll = useCallback(() => {
    const container = contentRef.current;
    if (!container) return;

    // 滚动到底部时，直接激活最后一个 section
    if (container.scrollTop + container.clientHeight >= container.scrollHeight - 1) {
      setActiveSection(SECTIONS[SECTIONS.length - 1].key);
      return;
    }

    const scrollTop = container.scrollTop + 80;
    for (let i = SECTIONS.length - 1; i >= 0; i--) {
      const ref = sectionRefs.current[SECTIONS[i].key];
      if (ref && ref.offsetTop <= scrollTop) {
        setActiveSection(SECTIONS[i].key);
        return;
      }
    }
    setActiveSection('about');
  }, []);

  useEffect(() => {
    const container = contentRef.current;
    if (!container) return;
    container.addEventListener('scroll', handleScroll, { passive: true });
    return () => container.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);

  const scrollTo = (key: SectionKey) => {
    const ref = sectionRefs.current[key];
    if (ref && contentRef.current) {
      contentRef.current.scrollTo({ top: ref.offsetTop - 24, behavior: 'smooth' });
    }
  };

  const renderNav = () => (
    <div className={styles.navCard}>
      <div className={styles.navList}>
        {SECTIONS.map(section => (
          <button
            key={section.key}
            type="button"
            className={`${styles.navItem} ${activeSection === section.key ? styles.active : ''}`}
            onClick={() => scrollTo(section.key)}
          >
            <span className={styles.navIcon}>{section.icon}</span>
            <span className={styles.navLabel}>{t(`settings.about.nav.${section.key}`)}</span>
          </button>
        ))}
      </div>
    </div>
  );

  const renderContent = () => (
    <div className={styles.contentArea} ref={contentRef}>
      {/* About */}
      <div ref={el => { sectionRefs.current.about = el; }} className={styles.section}>
        <div className={styles.appHeader}>
          <img src="/appicon.png" alt="Lemon Tea" className={styles.appIcon} />
          <div className={styles.appInfo}>
            <Title level={3} className={styles.appName}>Lemon Tea Desktop</Title>
            <Text className={styles.appVersion}>v0.0.1-dev</Text>
          </div>
        </div>
        <Paragraph className={styles.appDesc}>
          {t('settings.about.description')}
        </Paragraph>
      </div>

      {/* Features */}
      <div ref={el => { sectionRefs.current.features = el; }} className={styles.section}>
        <Title level={4} className={styles.sectionTitle}>
          <AppstoreOutlined /> {t('settings.about.nav.features')}
        </Title>
        <div className={styles.featureList}>
          {(['multiModel', 'toolCall', 'workflow', 'mcp', 'trace'] as const).map(key => (
            <div key={key} className={styles.featureItem}>
              <div className={styles.featureDot} />
              <div>
                <Text strong>{t(`settings.about.features.${key}.title`)}</Text>
                <br />
                <Text type="secondary">{t(`settings.about.features.${key}.desc`)}</Text>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Contact */}
      <div ref={el => { sectionRefs.current.contact = el; }} className={styles.section}>
        <Title level={4} className={styles.sectionTitle}>
          <MessageOutlined /> {t('settings.about.nav.contact')}
        </Title>
        <div className={styles.contactList}>
          <div className={styles.contactItem} onClick={() => void Browser.OpenURL('https://discord.gg/EJ2g3csW77')}>
            <div className={styles.contactIcon} style={{ background: 'rgba(88, 101, 242, 0.1)', color: '#5865F2' }}>
              <DiscordIcon />
            </div>
            <div className={styles.contactInfo}>
              <Text strong>Discord</Text>
              <Text type="secondary">discord.gg/EJ2g3csW77</Text>
            </div>
            <LinkOutlined className={styles.contactArrow} />
          </div>
          <div className={styles.contactItem} onClick={() => void Browser.OpenURL('https://github.com/MikeLINGxZ/lemantea')}>
            <div className={styles.contactIcon} style={{ background: 'rgba(36, 41, 47, 0.1)', color: '#24292f' }}>
              <GithubOutlined />
            </div>
            <div className={styles.contactInfo}>
              <Text strong>GitHub</Text>
              <Text type="secondary">MikeLINGxZ/lemantea</Text>
            </div>
            <LinkOutlined className={styles.contactArrow} />
          </div>
          <div className={styles.contactItem} onClick={() => { navigator.clipboard.writeText('1087746402'); void message.success(t('settings.about.contact.copied')); }}>
            <div className={styles.contactIcon} style={{ background: 'rgba(18, 183, 245, 0.1)', color: '#12B7F5' }}>
              <QQIcon />
            </div>
            <div className={styles.contactInfo}>
              <Text strong>{t('settings.about.contact.qqGroup')}</Text>
              <Text type="secondary">1087746402</Text>
            </div>
            <Text type="secondary" className={styles.contactHint}>{t('settings.about.contact.clickToCopy')}</Text>
          </div>
          <div className={styles.contactItem} onClick={() => void Browser.OpenURL('mailto:lpxqu@qq.com')}>
            <div className={styles.contactIcon} style={{ background: 'rgba(24, 144, 255, 0.1)', color: '#1890ff' }}>
              <MailOutlined />
            </div>
            <div className={styles.contactInfo}>
              <Text strong>{t('settings.about.contact.email')}</Text>
              <Text type="secondary">lpxqu@qq.com</Text>
            </div>
            <LinkOutlined className={styles.contactArrow} />
          </div>
        </div>
        <Paragraph className={styles.contributeText}>
          {t('settings.about.contribute')}
        </Paragraph>
      </div>
    </div>
  );

  if (isMobile) {
    return (
      <div className={styles.aboutPage}>
        {renderContent()}
      </div>
    );
  }

  return (
    <div className={styles.aboutPage}>
      <div className={styles.desktopLayout}>
        <div className={styles.navColumn}>{renderNav()}</div>
        <div className={styles.contentColumn}>{renderContent()}</div>
      </div>
    </div>
  );
};

export default AboutPage;
