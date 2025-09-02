import { notification } from 'antd';
import type { NotificationPlacement } from 'antd/es/notification/interface';

export type NotificationType = 'success' | 'info' | 'warning' | 'error';

interface NotificationConfig {
  type: NotificationType;
  title: string;
  message: string;
  duration?: number;
  placement?: NotificationPlacement;
}

// 默认配置
const defaultConfig = {
  duration: 4.5,
  placement: 'topRight' as NotificationPlacement,
};

// 显示通知
export const showNotification = (config: NotificationConfig) => {
  const { type, title, message, duration = defaultConfig.duration, placement = defaultConfig.placement } = config;
  
  notification[type]({
    message: title,
    description: message,
    duration,
    placement,
  });
};

// 便捷方法
export const notify = {
  success: (title: string, message: string, duration?: number) => {
    showNotification({ type: 'success', title, message, duration });
  },
  info: (title: string, message: string, duration?: number) => {
    showNotification({ type: 'info', title, message, duration });
  },
  warning: (title: string, message: string, duration?: number) => {
    showNotification({ type: 'warning', title, message, duration });
  },
  error: (title: string, message: string, duration?: number) => {
    showNotification({ type: 'error', title, message, duration });
  },
  clear: () => {
    notification.destroy();
  },
};
