import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import '@/styles/index.scss';
import App from '@/App.tsx';
import { initializeFontSize } from '@/stores/fontSizeStore';
import { hydrateLanguagePreferences, initializeLanguage } from '@/stores/languageStore';
import { hydrateLabPreferences } from '@/stores/labStore';
import { useAppLocale } from '@/hooks/useAppLocale';
import '@/i18n';

// 初始化字体大小设置
initializeFontSize();
initializeLanguage();

// 全局禁用 input 和 textarea 的自动填充提示
function disableAutocomplete() {
  // 设置所有现有的 input 和 textarea
  const setAutocomplete = (element: HTMLElement) => {
    if (element instanceof HTMLInputElement || element instanceof HTMLTextAreaElement) {
      element.setAttribute('autocomplete', 'off');
      element.setAttribute('autocorrect', 'off');
      element.setAttribute('autocapitalize', 'off');
      element.setAttribute('spellcheck', 'false');
    }
  };

  // 设置现有的所有 input 和 textarea
  // @ts-ignore
  document.querySelectorAll('input, textarea').forEach(setAutocomplete);

  // 监听 DOM 变化，确保新创建的 input 和 textarea 也禁用 autocomplete
  const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      mutation.addedNodes.forEach((node) => {
        if (node.nodeType === Node.ELEMENT_NODE) {
          const element = node as HTMLElement;
          // 检查新添加的节点本身
          setAutocomplete(element);
          // 检查新添加节点内的所有 input 和 textarea
          // @ts-ignore
          element.querySelectorAll?.('input, textarea').forEach(setAutocomplete);
        }
      });
    });
  });

  // 开始观察整个文档的变化
  observer.observe(document.body, {
    childList: true,
    subtree: true,
  });
}

// 在 DOM 加载完成后执行
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', disableAutocomplete);
} else {
  disableAutocomplete();
}

function AppProviders() {
  const { antdLocale } = useAppLocale();

  return (
    <BrowserRouter>
      <ConfigProvider
        locale={antdLocale}
        theme={{
          token: {
            colorPrimary: '#1890ff',
            borderRadius: 6,
          },
        }}
      >
        <App />
      </ConfigProvider>
    </BrowserRouter>
  );
}

async function bootstrap() {
  await hydrateLanguagePreferences();
  await hydrateLabPreferences();
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <AppProviders />
    </StrictMode>,
  );
}

void bootstrap();
