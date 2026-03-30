import React from 'react';
import { Card, Space, Typography } from 'antd';
import {
  CodeOutlined,
  HeartFilled,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import styles from './index.module.scss';

const { Title, Text, Paragraph } = Typography;

const AboutPage: React.FC = () => {
  const { t } = useTranslation();
  const projectInfo = {
    name: 'Lemon Tea Desktop',
    version: '0.0.1-dev',
    description: t('settings.about.description'),
  };

  const techStack = [
    { name: 'Wails3', color: '#1890ff', description: t('settings.about.techDescriptions.wails3') },
    { name: 'React', color: '#61dafb', description: t('settings.about.techDescriptions.react') },
    { name: 'TypeScript', color: '#3178c6', description: t('settings.about.techDescriptions.typescript') },
    { name: 'Ant Design', color: '#1890ff', description: t('settings.about.techDescriptions.antd') },
    { name: 'Go', color: '#00add8', description: t('settings.about.techDescriptions.go') },
    { name: 'Vite', color: '#646cff', description: t('settings.about.techDescriptions.vite') },
    { name: 'SCSS', color: '#cf649a', description: t('settings.about.techDescriptions.scss') },
    { name: 'Zustand', color: '#ff6b35', description: t('settings.about.techDescriptions.zustand') },
  ];

  return (
    <div className={styles.aboutContainer}>
      <div className={styles.header}>
        <div className={styles.logo}>
          <div className={styles.logoIcon}>🍋</div>
          <div className={styles.logoText}>
            <Title level={2} className={styles.title}>
              {projectInfo.name}
            </Title>
            <Text className={styles.version}>v{projectInfo.version}</Text>
          </div>
        </div>
        <Paragraph className={styles.description}>
          {projectInfo.description}
        </Paragraph>
      </div>

      <div className={styles.content}>

        {/* 技术栈卡片 */}
        <Card 
          className={styles.techCard}
          title={
            <Space>
              <CodeOutlined />
              {t('settings.about.techStack')}
            </Space>
          }
          bordered={false}
        >
          <div className={styles.techGrid}>
            {techStack.map((tech, index) => (
              <div key={index} className={styles.techItem}>
                <div 
                  className={styles.techIcon} 
                  style={{ backgroundColor: `${tech.color}20`, color: tech.color }}
                >
                  {tech.name[0]}
                </div>
                <div className={styles.techContent}>
                  <Text strong className={styles.techName}>
                    {tech.name}
                  </Text>
                  <Text className={styles.techDesc}>
                    {tech.description}
                  </Text>
                </div>
              </div>
            ))}
          </div>
        </Card>

        {/* 感谢卡片 */}
        <Card 
          className={styles.thanksCard}
          bordered={false}
        >
          <div className={styles.thanksContent}>
            <Title level={4} className={styles.thanksTitle}>
              <HeartFilled style={{ color: '#ff4d4f' }} /> {t('settings.about.thanksTitle')}
            </Title>
            <Paragraph className={styles.thanksText}>
              {t('settings.about.thanksText')}
            </Paragraph>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default AboutPage;
