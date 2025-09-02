// 自动生成的类型定义文件
// 此文件由 generate_rpc_ts_axios.sh 脚本生成

// 认证相关类型
export interface LoginRequest {
  username: string;
  password_md5: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    username: string;
    email: string;
  };
}

export interface RegisterRequest {
  username: string;
  password_md5: string;
  email: string;
  email_verification_code: number;
}

export interface SendEmailVerificationCodeRequest {
  email: string;
  username: string;
  code_type: number;
}

export interface CheckFieldAvailabilityRequest {
  field_type: number;
  value: string;
}

export interface CheckFieldAvailabilityResponse {
  available: boolean;
}

export interface ResetPasswordRequest {
  email: string;
  email_verification_code: number;
  new_password_md5: string;
}

// 聊天相关类型
export interface CreateChatRequest {
  title: string;
  model_id: string;
}

export interface CreateChatResponse {
  chat_id: string;
  title: string;
  model_id: string;
}

export interface SendMessageRequest {
  chat_id: string;
  content: string;
}

export interface Message {
  message_id: string;
  content: string;
  role: string;
  timestamp: string;
}

export interface ChatInfo {
  chat_id: string;
  title: string;
  model_id: string;
  created_at: string;
}

// 模型相关类型
export interface ModelInfo {
  model_id: string;
  name: string;
  description: string;
  provider: string;
}

export interface ModelDetail extends ModelInfo {
  parameters: Record<string, any>;
}

// 通用响应类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// 字段类型枚举
export enum FieldType {
  FIELD_TYPE_UNSPECIFIED = 0,
  FIELD_TYPE_USERNAME = 1,
  FIELD_TYPE_EMAIL = 2,
}

// 验证码类型枚举
export enum VerificationCodeType {
  VERIFICATION_CODE_TYPE_UNSPECIFIED = 0,
  VERIFICATION_CODE_TYPE_REGISTER = 1,
  VERIFICATION_CODE_TYPE_RESET_PASSWORD = 2,
}
