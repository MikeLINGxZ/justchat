import React from 'react';
import { Card, Space, Tag, Typography, Divider, Button } from 'antd';
import {
  InfoCircleOutlined,
  CodeOutlined,
  TeamOutlined,
  GithubOutlined,
  HeartFilled,
  StarFilled,
  RocketOutlined,
} from '@ant-design/icons';
import styles from './index.module.scss';

const { Title, Text, Paragraph } = Typography;

const AboutPage: React.FC = () => {
  const projectInfo = {
    name: 'Lemon Tea Desktop',
    version: '0.0.1-dev',
    description: '基于 Wails3 构建的现代化AI聊天桌面应用',
    author: 'Lemon Tea Team',
    license: 'MIT',
    repository: 'https://github.com/your-repo/lemon-tea-desktop'
  };

  const techStack = [
    { name: 'Wails3', color: '#1890ff', description: '跨平台桌面应用框架' },
    { name: 'React', color: '#61dafb', description: '前端UI框架' },
    { name: 'TypeScript', color: '#3178c6', description: '类型安全的JavaScript' },
    { name: 'Ant Design', color: '#1890ff', description: 'React UI组件库' },
    { name: 'Go', color: '#00add8', description: '后端语言' },
    { name: 'Vite', color: '#646cff', description: '前端构建工具' },
    { name: 'SCSS', color: '#cf649a', description: 'CSS预处理器' },
    { name: 'Zustand', color: '#ff6b35', description: '状态管理库' }
  ];

  const features = [
    { icon: <RocketOutlined />, title: '高性能', description: '基于原生技术栈，运行流畅' },
    { icon: <CodeOutlined />, title: '现代化', description: '使用最新的前端技术栈' },
    { icon: <StarFilled />, title: '易用性', description: '简洁直观的用户界面' },
    { icon: <HeartFilled />, title: '开源', description: '完全开源，持续维护' }
  ];

  const handleOpenRepository = () => {
    // 这里可以调用 Wails 的 API 打开浏览器
    window.open(projectInfo.repository, '_blank');
  };

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
              技术栈
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
              <HeartFilled style={{ color: '#ff4d4f' }} /> 感谢使用
            </Title>
            <Paragraph className={styles.thanksText}>
              感谢您使用 Lemon Tea Desktop！这是一个用 ❤️ 构建的 AI 聊天桌面应用。
            </Paragraph>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default AboutPage;