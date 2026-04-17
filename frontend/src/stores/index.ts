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

// 导出语言 store
export {
  useLanguageStore,
  hydrateLanguagePreferences,
  initializeLanguage,
} from './languageStore';
export { hydrateLabPreferences, useLabStore } from './labStore';

// 导出OPC store
export { useOPCStore, initializeOPC } from './opcStore';

// 导入初始化函数
import { initializeAuth } from './authStore';
import { useAuthStore } from './authStore';
import { initializeFontSize, useFontSizeStore } from './fontSizeStore';
import { initializeLanguage, useLanguageStore } from './languageStore';
import { useLabStore } from './labStore';
import { initializeOPC } from './opcStore';

// 初始化函数
export const initializeStores = async () => {
  // 初始化认证状态
  initializeAuth();
  // 初始化字体大小设置
  initializeFontSize();
  // 初始化语言设置
  initializeLanguage();
  // 初始化OPC状态
  initializeOPC();
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
  // 重置语言
  useLanguageStore.getState().setLanguage('zh-CN');
  useLabStore.getState().setMemorySystemEnabled(false);
  useLabStore.getState().setVectorSearchEnabled(false);
};
