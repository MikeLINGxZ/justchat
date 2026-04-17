import React, { useEffect } from 'react';
import { Routes, Route, Navigate, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import Layout from '@/components/layout';
import { initializeStores } from '@/stores';
import { useViewportHeight } from '@/hooks/useViewportHeight';

const Chat = React.lazy(() => import('@/pages/home'));
const NotFound = React.lazy(() => import('@/pages/common/NotFound.tsx'));
const Settings = React.lazy(()=>import('@/pages/settings'));
const Onboarding = React.lazy(() => import('@/pages/onboarding'));
const PluginPage = React.lazy(() => import('@/pages/plugin'));

function EntryRedirect() {
  const [searchParams] = useSearchParams();
  const entry = searchParams.get('entry');

  const tab = searchParams.get('tab');

  switch (entry) {
    case 'settings':
      return <Navigate to={tab ? `/settings?tab=${tab}` : '/settings'} replace />;
    case 'onboarding':
      return <Navigate to="/onboarding" replace />;
    case 'home':
    case null:
      return <Navigate to="/home" replace />;
    default:
      return <Navigate to="/home" replace />;
  }
}

function App() {
  const { t } = useTranslation();

  // 初始化视口高度检测
  useViewportHeight();

  // 初始化所有stores
  useEffect(() => {
    void initializeStores();
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
          <span style={{ fontSize: '14px', color: '#666', whiteSpace: 'nowrap' }}>{t('common.loading')}</span>
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
        {/* 应用入口页 - 根据窗口入口参数分发到对应页面 */}
        <Route
          path="/"
          element={<EntryRedirect />}
        />

        {/* 聊天页面 */}
        <Route
          path="/home"
          element={<Chat />}
        />
        <Route
          path="/home/:chatUuid"
          element={<Chat />}
        />

        <Route
          path="/settings"
          element={<Settings />}
        />
        <Route
          path="/onboarding"
          element={<Onboarding />}
        />

        {/* 插件页面 */}
        <Route
          path="/plugin/:pluginId/:pageId"
          element={<PluginPage />}
        />

        {/* 其他路由 - 使用Layout */}
        <Route
          path="/app"
          element={<Layout />}
        >
        </Route>

        {/* 兼容旧链接 */}
        <Route
          path="/:chatUuid"
          element={<Chat />}
        />

        {/* 404页面 */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </React.Suspense>
  );
}

export default App;
