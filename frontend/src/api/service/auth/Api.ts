/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
 * FieldType 字段类型枚举
 * @default "FIELD_TYPE_UNSPECIFIED"
 */
export enum ServerFieldType {
  FIELD_TYPE_UNSPECIFIED = "FIELD_TYPE_UNSPECIFIED",
  FIELD_TYPE_USERNAME = "FIELD_TYPE_USERNAME",
  FIELD_TYPE_EMAIL = "FIELD_TYPE_EMAIL",
}

/** @default "VERIFICATION_CODE_TYPE_UNSPECIFIED" */
export enum CommonVerificationCodeType {
  VERIFICATION_CODE_TYPE_UNSPECIFIED = "VERIFICATION_CODE_TYPE_UNSPECIFIED",
  VERIFICATION_CODE_TYPE_REGISTER = "VERIFICATION_CODE_TYPE_REGISTER",
  VERIFICATION_CODE_TYPE_RESET_PASSWORD = "VERIFICATION_CODE_TYPE_RESET_PASSWORD",
}

/** Empty 空参数占位 */
export type CommonEmpty = object;

export interface ProtobufAny {
  "@type"?: string;
  [key: string]: any;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/** CheckFieldAvailabilityRequest 检查字段可用性请求 */
export interface ServerCheckFieldAvailabilityRequest {
  /** 字段类型 */
  fieldType?: ServerFieldType;
  /** 字段值 */
  value?: string;
}

/** CheckFieldAvailabilityResponse 检查字段可用性返回 */
export interface ServerCheckFieldAvailabilityResponse {
  /** 是否可用 */
  available?: boolean;
}

/** LoginRequest 登录请求 */
export interface ServerLoginRequest {
  /** 登录字段（用户名或邮箱） */
  loginField?: string;
  /** 密码MD5 */
  passwordMd5?: string;
}

/** LoginResponse 登录返回 */
export interface ServerLoginResponse {
  /** 访问令牌 */
  accessToken?: string;
  /**
   * 令牌过期时间（秒）
   * @format int64
   */
  expiresIn?: string;
  /** 用户信息 */
  userInfo?: ServerUserInfo;
}

/**
 * LogoutRequest 登出请求
 * 可以为空，通过认证头获取用户信息
 */
export type ServerLogoutRequest = object;

/** RegisterRequest 注册请求 */
export interface ServerRegisterRequest {
  /** 用户名 */
  username?: string;
  /** 密码MD5 */
  passwordMd5?: string;
  /** 邮箱 */
  email?: string;
  /**
   * 邮箱验证码
   * @format int64
   */
  emailVerificationCode?: string;
}

/** ResetPasswordRequest 重置密码请求 */
export interface ServerResetPasswordRequest {
  /** 登录字段（用户名或邮箱） */
  loginField?: string;
  /** 邮箱验证码 */
  emailVerificationCode?: string;
  /** 新密码MD5 */
  newPasswordMd5?: string;
}

/** SendEmailVerificationCodeRequest 发送邮件验证码请求 */
export interface ServerSendEmailVerificationCodeRequest {
  /** 邮箱地址 */
  email?: string;
  username?: string;
  codeType?: CommonVerificationCodeType;
}

/** UserInfo 用户信息 */
export interface ServerUserInfo {
  /**
   * 用户ID
   * @format int64
   */
  userId?: string;
  /** 用户名 */
  username?: string;
  /** 邮箱 */
  email?: string;
  /**
   * 创建时间（Unix时间戳）
   * @format int64
   */
  createdAt?: string;
  /**
   * 更新时间（Unix时间戳）
   * @format int64
   */
  updatedAt?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title rpc/service/auth.proto
 * @version version not set
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  v1 = {
    /**
     * No description
     *
     * @tags Auth
     * @name AuthCheckFieldAvailability
     * @summary CheckFieldAvailability 检查用户名或邮箱是否被使用
     * @request POST:/v1/auth/check-availability
     */
    authCheckFieldAvailability: (
      body: ServerCheckFieldAvailabilityRequest,
      params: RequestParams = {},
    ) =>
      this.request<ServerCheckFieldAvailabilityResponse, RpcStatus>({
        path: `/v1/auth/check-availability`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Auth
     * @name AuthLogin
     * @summary Login 登录
     * @request POST:/v1/auth/login
     */
    authLogin: (body: ServerLoginRequest, params: RequestParams = {}) =>
      this.request<ServerLoginResponse, RpcStatus>({
        path: `/v1/auth/login`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Auth
     * @name AuthLogout
     * @summary Logout 登出
     * @request POST:/v1/auth/logout
     */
    authLogout: (body: ServerLogoutRequest, params: RequestParams = {}) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/auth/logout`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Auth
     * @name AuthRegister
     * @summary Register 注册
     * @request POST:/v1/auth/register
     */
    authRegister: (body: ServerRegisterRequest, params: RequestParams = {}) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/auth/register`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Auth
     * @name AuthResetPassword
     * @summary ResetPassword 忘记密码
     * @request POST:/v1/auth/reset-password
     */
    authResetPassword: (
      body: ServerResetPasswordRequest,
      params: RequestParams = {},
    ) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/auth/reset-password`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Auth
     * @name AuthSendEmailVerificationCode
     * @summary SendEmailVerificationCode 发送邮件验证码
     * @request POST:/v1/auth/send-email-code
     */
    authSendEmailVerificationCode: (
      body: ServerSendEmailVerificationCodeRequest,
      params: RequestParams = {},
    ) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/auth/send-email-code`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
