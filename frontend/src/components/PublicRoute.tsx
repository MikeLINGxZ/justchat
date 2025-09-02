import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { Spin } from 'antd';
import { useAuthStore } from '@/stores/authStore';

interface PublicRouteProps {
  children: React.ReactNode;
  redirectTo?: string;
}

const PublicRoute: React.FC<PublicRouteProps> = ({ 
  children, 
  redirectTo = '/home' 
}) => {
  const { isAuthenticated, isLoading } = useAuthStore();
  const location = useLocation();

  // 如果正在加载，显示加载状态
  if (isLoading) {
    return (
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100vh',
        }}
      >
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  // 如果已认证，重定向到指定页面（通常是仪表板）
  if (isAuthenticated) {
    // 获取原来想要访问的页面，如果没有则使用默认重定向页面
    const from = (location.state as any)?.from?.pathname || redirectTo;
    return <Navigate to={from} replace />;
  }

  // 如果未认证，渲染子组件（登录/注册页面）
  return <>{children}</>;
};

export default PublicRoute;