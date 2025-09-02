import { Api } from './service/auth/Api';
import { getApiConfig } from '@/config/env';
import { createGlobalErrorHandlerFetch } from '@/utils/globalErrorHandler';
import type {
  ServerLoginRequest,
  ServerLoginResponse,
  ServerRegisterRequest,
  ServerResetPasswordRequest,
  ServerSendEmailVerificationCodeRequest,
  ServerCheckFieldAvailabilityRequest,
  ServerCheckFieldAvailabilityResponse,
  ServerLogoutRequest,
  CommonEmpty
} from './service/auth/Api';

// 创建API客户端实例
class AuthClient {
  private api: Api<any>;

  constructor() {
    const apiConfig = getApiConfig();
    
    this.api = new Api({
      baseUrl: apiConfig.baseUrl,
      baseApiParams: {
        credentials: apiConfig.withCredentials ? 'include' : 'same-origin',
        headers: {
          'Content-Type': 'application/json',
        },
      },
      securityWorker: (_securityData) => {
        const token = localStorage.getItem('token');
        if (token) {
          return {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          };
        }
        return {};
      },
      customFetch: createGlobalErrorHandlerFetch(fetch),
    });
  }

  // 登录
  async login(data: ServerLoginRequest): Promise<ServerLoginResponse> {
    const response = await this.api.v1.authLogin(data);
    return response.data;
  }

  // 注册
  async register(data: ServerRegisterRequest): Promise<CommonEmpty> {
    const response = await this.api.v1.authRegister(data);
    return response.data;
  }

  // 登出
  async logout(data: ServerLogoutRequest = {}): Promise<CommonEmpty> {
    const response = await this.api.v1.authLogout(data);
    return response.data;
  }

  // 发送邮件验证码
  async sendEmailVerificationCode(data: ServerSendEmailVerificationCodeRequest): Promise<CommonEmpty> {
    const response = await this.api.v1.authSendEmailVerificationCode(data);
    return response.data;
  }

  // 检查字段可用性
  async checkFieldAvailability(data: ServerCheckFieldAvailabilityRequest): Promise<ServerCheckFieldAvailabilityResponse> {
    const response = await this.api.v1.authCheckFieldAvailability(data);
    return response.data;
  }

  // 重置密码
  async resetPassword(data: ServerResetPasswordRequest): Promise<CommonEmpty> {
    const response = await this.api.v1.authResetPassword(data);
    return response.data;
  }
}

// 创建单例实例
const authClient = new AuthClient();

export { authClient };
export type {
  ServerLoginRequest,
  ServerLoginResponse,
  ServerRegisterRequest,
  ServerResetPasswordRequest,
  ServerSendEmailVerificationCodeRequest,
  ServerCheckFieldAvailabilityRequest,
  ServerCheckFieldAvailabilityResponse,
  ServerLogoutRequest,
  CommonEmpty
};