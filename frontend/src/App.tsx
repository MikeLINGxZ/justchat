import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import Layout from '@/components/Layout';
import PrivateRoute from '@/components/PrivateRoute';
import PublicRoute from '@/components/PublicRoute';
import Login from '@/pages/auth/Login';
import Register from '@/pages/auth/Register';
import ForgotPassword from '@/pages/auth/ForgotPassword';
import TestAuth from '@/pages/TestAuth';
import SimpleTest from '@/pages/SimpleTest';
import EnvTest from '@/pages/EnvTest';
import TestLoginExpiry from '@/pages/TestLoginExpiry';
import { initializeStores } from '@/stores';
import { useViewportHeight } from '@/hooks/useViewportHeight';

// 懒加载页面组件
const Profile = React.lazy(() => import('@/pages/Profile'));
const Settings = React.lazy(() => import('@/pages/Settings'));
const Chat = React.lazy(() => import('@/pages/home'));
const NotFound = React.lazy(() => import('@/pages/NotFound'));

function App() {
  // 初始化视口高度检测
  useViewportHeight();
  
  // 初始化所有stores
  useEffect(() => {
    initializeStores();
  }, []);

  return (
    <React.Suspense
      fallback={
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
          }}
        >
          <Spin size="large" tip="页面加载中..." />
        </div>
      }
    >
      <Routes>
        {/* 公共路由 */}
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/register"
          element={
            <PublicRoute>
              <Register />
            </PublicRoute>
          }
        />
        <Route
          path="/forgot-password"
          element={
            <PublicRoute>
              <ForgotPassword />
            </PublicRoute>
          }
        />
        <Route
          path="/test-auth"
          element={
            <PublicRoute>
              <TestAuth />
            </PublicRoute>
          }
        />
        <Route
          path="/simple-test"
          element={<SimpleTest />}
        />
        <Route
          path="/env-test"
          element={<EnvTest />}
        />
        <Route
          path="/test-login-expiry"
          element={<TestLoginExpiry />}
        />

        {/* 聊天页面 - 独立路由 */}
        <Route
          path="/home"
          element={
            <PrivateRoute>
              <Chat />
            </PrivateRoute>
          }
        />
        {/* 带chatUuid参数的聊天页面路由 */}
        <Route
          path="/home/:chatUuid"
          element={
            <PrivateRoute>
              <Chat />
            </PrivateRoute>
          }
        />

        {/* 私有路由 - 使用Layout */}
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }
        >
          <Route index element={<Navigate to="/home" replace />} />
          <Route path="profile" element={<Profile />} />
          <Route path="settings" element={<Settings />} />
        </Route>

        {/* 404页面 */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </React.Suspense>
  );
}

export default App;
