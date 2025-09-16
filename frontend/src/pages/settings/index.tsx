import React, { useState } from 'react';
import { Layout, Menu } from 'antd';
import {
  SettingOutlined,
  ApiOutlined,
  UserOutlined,
  SecurityScanOutlined,
  BellOutlined,
} from '@ant-design/icons';
import ProviderSettingPage from './provider';
import styles from './index.module.scss';

const { Sider, Content } = Layout;

interface SettingsPageProps {
  className?: string;
}

const SettingsPage: React.FC<SettingsPageProps> = ({ className }) => {
  const [selectedKey, setSelectedKey] = useState('provider');

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
  };

  const renderContent = () => {
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
  };

  return (
    <Layout className={`${styles.settingsLayout} ${className || ''}`}>
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
    </Layout>
  );
};

export default SettingsPage;