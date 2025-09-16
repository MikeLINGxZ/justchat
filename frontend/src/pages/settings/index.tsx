import React, { useState } from 'react';
import { Layout, Menu, Button } from 'antd';
import {
  SettingOutlined,
  ApiOutlined,
  UserOutlined,
  SecurityScanOutlined,
  BellOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons';
import ProviderSettingPage from './provider';
import { useViewportHeight } from '@/hooks/useViewportHeight';
import styles from './index.module.scss';

const { Sider, Content } = Layout;

interface SettingsPageProps {
  className?: string;
}

const SettingsPage: React.FC<SettingsPageProps> = ({ className }) => {
  const [selectedKey, setSelectedKey] = useState('provider');
  const [showContent, setShowContent] = useState(false); // 控制移动端内容显示
  const { isMobile } = useViewportHeight(); // 使用移动端检测

  const menuItems = [
    {
      key: 'provider',
      icon: <ApiOutlined />,
      label: '模型供应商',
    },
    {
      key: 'general',
      icon: <SettingOutlined />,
      label: '通用设置',
    },
    {
      key: 'account',
      icon: <UserOutlined />,
      label: '账户设置',
    },
    {
      key: 'security',
      icon: <SecurityScanOutlined />,
      label: '安全设置',
    },
    {
      key: 'notifications',
      icon: <BellOutlined />,
      label: '通知设置',
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    setSelectedKey(key);
    // 移动端点击菜单后显示内容
    if (isMobile) {
      setShowContent(true);
    }
  };

  const handleBackToMenu = () => {
    // 移动端返回菜单
    setShowContent(false);
  };

  const renderContent = () => {
    const content = (() => {
      switch (selectedKey) {
        case 'provider':
          return <ProviderSettingPage />;
        case 'general':
          return <div className={styles.placeholder}>通用设置功能开发中...</div>;
        case 'account':
          return <div className={styles.placeholder}>账户设置功能开发中...</div>;
        case 'security':
          return <div className={styles.placeholder}>安全设置功能开发中...</div>;
        case 'notifications':
          return <div className={styles.placeholder}>通知设置功能开发中...</div>;
        default:
          return <ProviderSettingPage />;
      }
    })();

    // 移动端在内容顶部添加返回按钮
    if (isMobile) {
      return (
        <div className={styles.mobileContent}>
          <div className={styles.mobileHeader}>
            <Button 
              type="text" 
              icon={<ArrowLeftOutlined />}
              onClick={handleBackToMenu}
              className={styles.backButton}
            >
              返回
            </Button>
            <span className={styles.mobileTitle}>
              {menuItems.find(item => item.key === selectedKey)?.label}
            </span>
          </div>
          <div className={styles.mobileContentBody}>
            {content}
          </div>
        </div>
      );
    }

    return content;
  };

  return (
    <Layout className={`${styles.settingsLayout} ${className || ''}`}>
      {/* 移动端：根据状态显示菜单或内容 */}
      {isMobile ? (
        <>
          {/* 移动端菜单 */}
          <div className={`${styles.mobileMenu} ${showContent ? styles.hidden : ''}`}>
            <div className={styles.siderHeader}>
              <h3>设置</h3>
            </div>
            <Menu
              mode="inline"
              selectedKeys={[selectedKey]}
              items={menuItems}
              onClick={handleMenuClick}
              className={styles.settingsMenu}
            />
          </div>
          
          {/* 移动端内容 */}
          <div className={`${styles.mobileContentContainer} ${showContent ? styles.visible : ''}`}>
            {renderContent()}
          </div>
        </>
      ) : (
        /* 桌面端：正常的侧边栏布局 */
        <>
          <Sider
            width={240}
            className={styles.settingsSider}
            theme="light"
          >
            <div className={styles.siderHeader}>
              <h3>设置</h3>
            </div>
            <Menu
              mode="inline"
              selectedKeys={[selectedKey]}
              items={menuItems}
              onClick={handleMenuClick}
              className={styles.settingsMenu}
            />
          </Sider>
          <Layout className={styles.settingsContent}>
            <Content className={styles.contentArea}>
              {renderContent()}
            </Content>
          </Layout>
        </>
      )}
    </Layout>
  );
};

export default SettingsPage;