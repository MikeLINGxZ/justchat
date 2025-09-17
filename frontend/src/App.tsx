import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from '@/components/Layout';
import SimpleTest from '@/pages/SimpleTest';
import EnvTest from '@/pages/EnvTest';
import { initializeStores } from '@/stores';
import { useViewportHeight } from '@/hooks/useViewportHeight';

const Chat = React.lazy(() => import('@/pages/home'));
const NotFound = React.lazy(() => import('@/pages/NotFound'));
const Settings = React.lazy(()=>import('@/pages/settings'));

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
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            gap: '16px'
          }}
        >
          <div style={{ display: 'flex', gap: '4px', alignItems: 'center' }}>
            {[...Array(5)].map((_, i) => (
              <div
                key={i}
                style={{
                  width: '4px',
                  height: '30px',
                  background: 'linear-gradient(45deg, #667eea, #764ba2)',
                  borderRadius: '2px',
                  animation: `loading-wave 1.2s ease-in-out infinite ${i * 0.1}s`
                }}
              />
            ))}
          </div>
          <span style={{ fontSize: '14px', color: '#666', whiteSpace: 'nowrap' }}>加载中...</span>
          <style>{`
            @keyframes loading-wave {
              0%, 40%, 100% {
                transform: scaleY(0.4);
                opacity: 0.6;
              }
              20% {
                transform: scaleY(1);
                opacity: 1;
              }
            }
          `}</style>
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
          <Route
              path="/settings"
              element={<Settings />}
          />

        {/* 其他路由 - 使用Layout */}
        <Route
          path="/app"
          element={<Layout />}
        >
        </Route>

        {/* 404页面 */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </React.Suspense>
  );
}

export default App;
