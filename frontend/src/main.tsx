import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import '@/styles/index.scss';
import App from '@/App.tsx';
import { initializeFontSize } from '@/stores/fontSizeStore';

// 初始化字体大小设置
initializeFontSize();

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

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <ConfigProvider
        locale={zhCN}
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
  </StrictMode>,
);
