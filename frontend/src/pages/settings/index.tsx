import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button } from 'antd';
import {
  SettingOutlined,
  ApiOutlined,
  UserOutlined,
  SecurityScanOutlined,
  BellOutlined,
  InfoCircleOutlined,
  ArrowLeftOutlined,
} from '@ant-design/icons';
import ProviderSettingPage from './provider';
import AboutPage from './about';
import GeneralSettingsPage from './general';
import { useViewportHeight } from '@/hooks/useViewportHeight';
import { initializeFontSize } from '@/stores/fontSizeStore';
import styles from './index.module.scss';

const { Sider, Content } = Layout;

interface SettingsPageProps {
  className?: string;
}

const SettingsPage: React.FC<SettingsPageProps> = ({ className }) => {
  const [selectedKey, setSelectedKey] = useState('general');
  const [showContent, setShowContent] = useState(false); // 控制移动端内容显示
  const { isMobile } = useViewportHeight(); // 使用移动端检测

  // 监听设备切换，处理从桌面端切换到移动端的情况
  useEffect(() => {
    if (isMobile) {
      // 切换到移动端时，如果当前是从桌面端切换过来，显示当前选中的内容
      // 这样用户不需要重新点击菜单
      setShowContent(true);
    } else {
      // 切换到桌面端时，重置移动端状态
      setShowContent(false);
    }
  }, [isMobile]);

  // 初始化字体大小设置
  useEffect(() => {
    // 初始化字体大小设置
    initializeFontSize();
  }, []);

  const menuItems = [
    {
      key: 'general',
      icon: <SettingOutlined />,
      label: '通用设置',
    },
    {
      key: 'provider',
      icon: <ApiOutlined />,
      label: '模型供应商',
    },
    {
      key: 'about',
      icon: <InfoCircleOutlined />,
      label: '关于',
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

  // 决定是否显示加载状态 - 已移除，使用App.tsx的统一loading

  const renderContent = () => {
    const content = (() => {
      switch (selectedKey) {
        case 'provider':
          return <ProviderSettingPage />;
        case 'general':
          return <GeneralSettingsPage />;
        case 'account':
          return <div className={styles.placeholder}>账户设置功能开发中...</div>;
        case 'security':
          return <div className={styles.placeholder}>安全设置功能开发中...</div>;
        case 'notifications':
          return <div className={styles.placeholder}>通知设置功能开发中...</div>;
        case 'about':
          return <AboutPage />;
        default:
          return <GeneralSettingsPage />;
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
              selectedKeys={[]} // 移动端不显示选中状态
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