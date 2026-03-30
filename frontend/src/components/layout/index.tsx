import React from 'react';
import { Layout as AntLayout, Menu, Avatar, Dropdown, Button, Space } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  UserOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MessageOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '@/stores/authStore';
import { useViewportHeight } from '@/hooks/useViewportHeight';
import styles from './index.module.scss';

const { Header, Sider, Content } = AntLayout;

interface LayoutProps {
  children?: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuthStore();
  const [collapsed, setCollapsed] = React.useState(false);
  
  // 使用视口高度检测 Hook
  const { isMobile } = useViewportHeight();

  const userMenuItems = [
    {
      key: 'chat',
      icon: <MessageOutlined />,
      label: t('layout.chat'),
      onClick: () => navigate('/home'),
    },
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: t('layout.profile'),
      onClick: () => navigate('/app/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: t('layout.settings'),
      onClick: () => navigate('/app/settings'),
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  return (
    <AntLayout className={`${styles.layout} ${isMobile ? 'mobile-viewport-height' : ''}`}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        className={styles.sider}
      >
        <div className={styles.logo}>
          <h2>{collapsed ? 'LT' : 'Lemon Tea'}</h2>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={[
            {
              key: '/app/profile',
              icon: <UserOutlined />,
              label: t('layout.profile'),
            },
            {
              key: '/app/settings',
              icon: <SettingOutlined />,
              label: t('layout.settings'),
            },
          ]}
          onClick={handleMenuClick}
        />
      </Sider>
      <AntLayout>
        <Header className={styles.header}>
          <div className={styles.headerLeft}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              className={styles.trigger}
            />
          </div>
          <div className={styles.headerRight}>
            <Space>
              <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                <div className={styles.userInfo}>
                  <Avatar
                    size="small"
                    icon={<UserOutlined />}
                    src={user?.avatar}
                  />
                  <span className={styles.username}>{user?.username}</span>
                </div>
              </Dropdown>
            </Space>
          </div>
        </Header>
        <Content className={styles.content}>
          <div className={styles.contentInner}>
            <Outlet />
          </div>
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;
