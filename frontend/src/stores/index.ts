// 导出auth store
export { useAuthStore, initializeAuth } from './authStore';

// 导出字体大小store
export { 
  useFontSizeStore, 
  initializeFontSize, 
  FONT_SIZE_OFFSETS, 
  FONT_SIZE_OPTIONS,
  getFontSizeLabel,
  getFontSizeDescription
} from './fontSizeStore';

// 导入初始化函数
import { initializeAuth } from './authStore';
import { useAuthStore } from './authStore';
import { initializeFontSize, useFontSizeStore } from './fontSizeStore';

// 初始化函数
export const initializeStores = () => {
  // 初始化认证状态
  initializeAuth();
  // 初始化字体大小设置
  initializeFontSize();
};

// 重置所有store
export const resetAllStores = () => {
  // 重置认证状态到初始状态
  useAuthStore.getState().setAuthState(
    {
      userId: 'demo-user',
      username: 'Local User',
      email: 'storage in local disk',
    },
    'demo-token',
    true
  );
  // 重置字体大小
  useFontSizeStore.getState().resetFontSize();
};