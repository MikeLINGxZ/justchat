// 导出auth store
export { useAuthStore, initializeAuth } from './authStore';

// 导入初始化函数
import { initializeAuth } from './authStore';
import { useAuthStore } from './authStore';

// 初始化函数
export const initializeStores = () => {
  // 只初始化认证状态
  initializeAuth();
};

// 重置所有store
export const resetAllStores = () => {
  useAuthStore.getState().logout();
};