// 通用类型定义
export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  message: string;
  code: number;
}

export interface PaginationParams {
  page: number;
  pageSize: number;
  total?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}

// 用户相关类型
export interface UserProfile {
  id: string;
  username: string;
  email: string;
  avatar?: string;
  bio?: string;
  phone?: string;
  createdAt: string;
  updatedAt: string;
}

// 路由相关类型
export interface RouteConfig {
  path: string;
  element: React.ComponentType;
  title?: string;
  requireAuth?: boolean;
  children?: RouteConfig[];
}

// 表单相关类型
export interface FormFieldError {
  field: string;
  message: string;
}

export interface ValidationResult {
  isValid: boolean;
  errors: FormFieldError[];
}

// Socket.IO 相关类型
export interface SocketEvent {
  type: string;
  payload: any;
  timestamp: number;
}

export interface ChatMessage {
  id: string;
  content: string;
  senderId: string;
  senderName: string;
  timestamp: number;
  type: 'text' | 'image' | 'file';
}

// 主题相关类型
export type ThemeMode = 'light' | 'dark';

export interface ThemeConfig {
  mode: ThemeMode;
  primaryColor: string;
  borderRadius: number;
}

// 通知相关类型
export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface NotificationConfig {
  type: NotificationType;
  title: string;
  message: string;
  duration?: number;
  placement?: 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight';
}

// 菜单相关类型
export interface MenuItem {
  key: string;
  label: string;
  icon?: React.ReactNode;
  path?: string;
  children?: MenuItem[];
  disabled?: boolean;
}

// 文件上传相关类型
export interface UploadFile {
  uid: string;
  name: string;
  status: 'uploading' | 'done' | 'error' | 'removed';
  url?: string;
  thumbUrl?: string;
  size?: number;
  type?: string;
}

export interface UploadResponse {
  url: string;
  filename: string;
  size: number;
  type: string;
}

// 环境变量类型
export interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_APP_TITLE: string;
  readonly VITE_SOCKET_URL: string;
  readonly VITE_UPLOAD_MAX_SIZE: string;
}

export interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// 工具类型
export type Nullable<T> = T | null;
export type Optional<T> = T | undefined;
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

// 状态管理相关类型
export interface StoreState {
  loading: boolean;
  error: string | null;
}

export interface AsyncAction<T = any> {
  pending: () => void;
  fulfilled: (payload: T) => void;
  rejected: (error: string) => void;
}

// 聊天信息接口
export interface ChatInfo {
  chatUuid: string;
  title: string;
  model: string;
  createdAt: string;
  updatedAt: string;
  messageCount: number;
  lastMessagePreview?: string;
  metadata?: Record<string, string>;
}

// 聊天状态接口
export interface ChatState {
  // 当前激活的聊天
  currentChatUuid: string | null;
  currentMessages: any[]; // 使用 any[] 替代 schema.Message[]
  
  // 聊天列表
  chatList: ChatInfo[];
  isLoadingChats: boolean;
  
  // 消息相关
  isLoadingMessages: boolean;
  isSendingMessage: boolean;
  
  // 搜索相关
  searchQuery: string;
  filteredChats: ChatInfo[];
  
  // 模型选择
  selectedModel: string;
  availableModels: string[];
  
  // UI状态
  isSidebarCollapsed: boolean;
  isTyping: boolean;
}