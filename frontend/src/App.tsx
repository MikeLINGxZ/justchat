import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import Layout from '@/components/Layout';
import SimpleTest from '@/pages/SimpleTest';
import EnvTest from '@/pages/EnvTest';
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
        {/* 测试路由 */}
        <Route
          path="/simple-test"
          element={<SimpleTest />}
        />
        <Route
          path="/env-test"
          element={<EnvTest />}
        />

        {/* 聊天页面 - 主页面 */}
        <Route
          path="/"
          element={<Chat />}
        />
        <Route
          path="/home"
          element={<Chat />}
        />
        {/* 带chatUuid参数的聊天页面路由 */}
        <Route
          path="/home/:chatUuid"
          element={<Chat />}
        />
        <Route
          path="/:chatUuid"
          element={<Chat />}
        />

        {/* 其他路由 - 使用Layout */}
        <Route
          path="/app"
          element={<Layout />}
        >
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
