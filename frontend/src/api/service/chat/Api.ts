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
 * ImageURLDetail 是图像 URL 的细节级别
 * - IMAGE_URL_DETAIL_UNSPECIFIED: 未指定
 *  - IMAGE_URL_DETAIL_HIGH: 高质量
 *  - IMAGE_URL_DETAIL_LOW: 低质量
 *  - IMAGE_URL_DETAIL_AUTO: 自动选择
 * @default "IMAGE_URL_DETAIL_UNSPECIFIED"
 */
export enum CommonImageURLDetail {
  IMAGE_URL_DETAIL_UNSPECIFIED = "IMAGE_URL_DETAIL_UNSPECIFIED",
  IMAGE_URL_DETAIL_HIGH = "IMAGE_URL_DETAIL_HIGH",
  IMAGE_URL_DETAIL_LOW = "IMAGE_URL_DETAIL_LOW",
  IMAGE_URL_DETAIL_AUTO = "IMAGE_URL_DETAIL_AUTO",
}

/**
 * ChatMessagePartType 是聊天消息中内容片段的类型
 * - CHAT_MESSAGE_PART_TYPE_UNSPECIFIED: 未指定
 *  - CHAT_MESSAGE_PART_TYPE_TEXT: 文本
 *  - CHAT_MESSAGE_PART_TYPE_IMAGE_URL: 图像链接
 *  - CHAT_MESSAGE_PART_TYPE_AUDIO_URL: 音频链接
 *  - CHAT_MESSAGE_PART_TYPE_VIDEO_URL: 视频链接
 *  - CHAT_MESSAGE_PART_TYPE_FILE_URL: 文件链接
 * @default "CHAT_MESSAGE_PART_TYPE_UNSPECIFIED"
 */
export enum CommonChatMessagePartType {
  CHAT_MESSAGE_PART_TYPE_UNSPECIFIED = "CHAT_MESSAGE_PART_TYPE_UNSPECIFIED",
  CHAT_MESSAGE_PART_TYPE_TEXT = "CHAT_MESSAGE_PART_TYPE_TEXT",
  CHAT_MESSAGE_PART_TYPE_IMAGE_URL = "CHAT_MESSAGE_PART_TYPE_IMAGE_URL",
  CHAT_MESSAGE_PART_TYPE_AUDIO_URL = "CHAT_MESSAGE_PART_TYPE_AUDIO_URL",
  CHAT_MESSAGE_PART_TYPE_VIDEO_URL = "CHAT_MESSAGE_PART_TYPE_VIDEO_URL",
  CHAT_MESSAGE_PART_TYPE_FILE_URL = "CHAT_MESSAGE_PART_TYPE_FILE_URL",
}

/** ChatTitleRequest 获取对话标题请求 */
export type ChatChatTitleBody = object;

/** ChatTitleSaveRequest 保存对话标题请求 */
export type ChatChatTitleSaveBody = object;

/** DeleteChatRequest 删除对话请求 */
export type ChatDeleteChatBody = object;

/** DeleteChatMessageRequest 删除对话消息请求 */
export type ChatDeleteChatMessageBody = object;

/** GetChatMessagesRequest 获取对话消息请求 */
export interface ChatGetChatMessagesBody {
  /** @format int64 */
  offset?: string;
  /** @format int64 */
  limit?: string;
}

/** ChatInfo 对话信息 */
export interface CommonChatInfo {
  /** 对话 ID */
  chatUuid?: string;
  /** 对话标题 */
  title?: string;
  /**
   * 使用的模型
   * @format int64
   */
  modelId?: string;
  /**
   * 创建时间（unix 秒）
   * @format int64
   */
  createdAt?: string;
  /**
   * 更新时间（unix 秒）
   * @format int64
   */
  updatedAt?: string;
  /**
   * 消息数量
   * @format int32
   */
  messageCount?: number;
  /** 最后一条消息预览 */
  lastMessagePreview?: string;
  /** 元数据 */
  metadata?: Record<string, string>;
}

/** ChatMessageAudioURL 表示聊天消息中的音频部分 */
export interface CommonChatMessageAudioURL {
  /** URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL */
  url?: string;
  uri?: string;
  /** MIMEType 是音频的 MIME 类型 */
  mimeType?: string;
  /** Extra 用于存储音频 URL 的额外信息 */
  extra?: Record<string, string>;
}

/** ChatMessageFileURL 表示聊天消息中的文件部分 */
export interface CommonChatMessageFileURL {
  /** 文件 URL */
  url?: string;
  uri?: string;
  /** MIMEType 是文件的 MIME 类型 */
  mimeType?: string;
  /** Name 是文件名称 */
  name?: string;
  /** Extra 用于存储文件 URL 的额外信息 */
  extra?: Record<string, string>;
}

/** ChatMessageImageURL 表示聊天消息中的图像部分 */
export interface CommonChatMessageImageURL {
  /** URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL（如 data URL） */
  url?: string;
  uri?: string;
  /**
   * Detail 是图像 URL 的质量等级
   * - IMAGE_URL_DETAIL_UNSPECIFIED: 未指定
   *  - IMAGE_URL_DETAIL_HIGH: 高质量
   *  - IMAGE_URL_DETAIL_LOW: 低质量
   *  - IMAGE_URL_DETAIL_AUTO: 自动选择
   */
  detail?: CommonImageURLDetail;
  /** MIMEType 是图像的 MIME 类型 */
  mimeType?: string;
  /** Extra 用于存储图像 URL 的额外信息 */
  extra?: Record<string, string>;
}

/** ChatMessagePart 是聊天消息中的内容片段 */
export interface CommonChatMessagePart {
  /**
   * Type 是片段的类型
   * - CHAT_MESSAGE_PART_TYPE_UNSPECIFIED: 未指定
   *  - CHAT_MESSAGE_PART_TYPE_TEXT: 文本
   *  - CHAT_MESSAGE_PART_TYPE_IMAGE_URL: 图像链接
   *  - CHAT_MESSAGE_PART_TYPE_AUDIO_URL: 音频链接
   *  - CHAT_MESSAGE_PART_TYPE_VIDEO_URL: 视频链接
   *  - CHAT_MESSAGE_PART_TYPE_FILE_URL: 文件链接
   */
  type?: CommonChatMessagePartType;
  /** Text 是文本内容，当 Type 为 "text" 时使用 */
  text?: string;
  /** ImageURL 是图像链接内容，当 Type 为 "image_url" 时使用 */
  imageUrl?: CommonChatMessageImageURL;
  /** AudioURL 是音频链接内容，当 Type 为 "audio_url" 时使用 */
  audioUrl?: CommonChatMessageAudioURL;
  /** VideoURL 是视频链接内容，当 Type 为 "video_url" 时使用 */
  videoUrl?: CommonChatMessageVideoURL;
  /** FileURL 是文件链接内容，当 Type 为 "file_url" 时使用 */
  fileUrl?: CommonChatMessageFileURL;
}

/** ChatMessageVideoURL 表示聊天消息中的视频部分 */
export interface CommonChatMessageVideoURL {
  /** URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL */
  url?: string;
  uri?: string;
  /** MIMEType 是视频的 MIME 类型 */
  mimeType?: string;
  /** Extra 用于存储视频 URL 的额外信息 */
  extra?: Record<string, string>;
}

/** Empty 空参数占位 */
export type CommonEmpty = object;

/** FunctionCall 是消息中的函数调用信息 */
export interface CommonFunctionCall {
  /** Name 是要调用的函数名称 */
  name?: string;
  /** Arguments 是调用函数所需的参数，以 JSON 格式表示 */
  arguments?: string;
}

/** LogProb 表示一个 token 的概率信息 */
export interface CommonLogProb {
  /** Token 表示 token 的文本内容 */
  token?: string;
  /**
   * LogProb 是该 token 的对数概率
   * @format double
   */
  logProb?: number;
  /** Bytes 是该 token 的 UTF-8 字节表示（整数列表） */
  bytes?: string[];
  /** TopLogProbs 是最可能的若干 token 及其对数概率列表 */
  topLogProbs?: CommonTopLogProb[];
}

/** LogProbs 是包含 token 概率信息的顶层结构 */
export interface CommonLogProbs {
  /** Content 是包含对数概率信息的消息内容 token 列表 */
  content?: CommonLogProb[];
}

/** Message 表示一条聊天消息，使用 oneof 分离不同角色的消息类型 */
export interface CommonMessage {
  /**
   * 通用字段（所有消息类型都可能包含）
   * 主要内容（纯文本）
   */
  content?: string;
  /** 多模态内容（如图片、文本、文件等） */
  multiContent?: CommonChatMessagePart[];
  /** 消息发送者名称（可选） */
  name?: string;
  /** 扩展信息（可选） */
  extra?: Record<string, string>;
  role?: string;
  /**
   * 系统消息相关字段（仅在 role == SYSTEM 时使用）
   * 系统指令内容（原 SystemMessage.content）
   */
  systemContent?: string;
  /**
   * 助手消息相关字段（仅在 role == ASSISTANT 时使用）
   * 工具调用请求
   */
  toolCalls?: CommonToolCall[];
  /** 响应元信息 */
  responseMeta?: CommonResponseMeta;
  /** 推理过程（如思维链） */
  reasoningContent?: string;
  /**
   * 工具消息相关字段（仅在 role == TOOL 时使用）
   * 工具调用 ID
   */
  toolCallId?: string;
  /** 工具名称 */
  toolName?: string;
  /** 以下为业务字段 */
  chatUuid?: string;
}

/** ResponseMeta 收集聊天响应的元信息 */
export interface CommonResponseMeta {
  /** FinishReason 是聊天响应结束的原因 */
  finishReason?: string;
  /** Usage 是聊天响应的 token 使用情况 */
  usage?: CommonTokenUsage;
  /** LogProbs 是对数概率信息 */
  logProbs?: CommonLogProbs;
}

/** TokenUsage 表示聊天模型请求的 token 使用情况 */
export interface CommonTokenUsage {
  /**
   * PromptTokens 是提示词中的 token 数量
   * @format int32
   */
  promptTokens?: number;
  /**
   * CompletionTokens 是生成内容中的 token 数量
   * @format int32
   */
  completionTokens?: number;
  /**
   * TotalTokens 是请求中总的 token 数量
   * @format int32
   */
  totalTokens?: number;
}

/** ToolCall 是消息中的工具调用信息 */
export interface CommonToolCall {
  /**
   * Index 在一条消息包含多个工具调用时使用
   * @format int32
   */
  index?: number;
  /** ID 是工具调用的唯一标识 */
  id?: string;
  /** Type 是工具调用的类型 */
  type?: string;
  /** Function 是具体的函数调用内容 */
  function?: CommonFunctionCall;
  /** Extra 用于存储工具调用的额外信息 */
  extra?: Record<string, string>;
}

/** TopLogProb 表示某个 token 的最高对数概率信息 */
export interface CommonTopLogProb {
  /** Token 表示 token 的文本内容 */
  token?: string;
  /**
   * LogProb 是该 token 的对数概率
   * @format double
   */
  logProb?: number;
  /** Bytes 是该 token 的 UTF-8 字节表示（整数列表） */
  bytes?: string[];
}

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

/** ChatCompletionChoice 聊天完成选择 */
export interface ServerChatCompletionChoice {
  /**
   * 选择索引
   * @format int32
   */
  index?: number;
  /** 生成的消息 */
  delta?: CommonMessage;
  /** 对数概率信息 */
  logprobs?: CommonLogProbs;
  /** 结束原因："stop", "length", "tool_calls", "content_filter", "function_call" */
  finishReason?: string;
}

/** ChatCompletionFunction 函数定义 */
export interface ServerChatCompletionFunction {
  /** 函数名称 */
  name?: string;
  /** 函数描述 */
  description?: string;
  /** JSON Schema格式的参数定义 */
  parameters?: string;
}

/** ChatCompletionNamedToolChoice 指定特定工具 */
export interface ServerChatCompletionNamedToolChoice {
  /** 工具类型 */
  type?: string;
  /** 指定的函数 */
  function?: ServerChatCompletionFunction;
}

/** ChatCompletionTool 工具定义 */
export interface ServerChatCompletionTool {
  /** 工具类型，通常为"function" */
  type?: string;
  /** 函数定义 */
  function?: ServerChatCompletionFunction;
}

/** ChatTitleResponse 获取对话标题响应 */
export interface ServerChatTitleResponse {
  title?: string;
}

/** CreateChatCompletionRequest 与 OpenAI Chat Completions 请求体对齐 */
export interface ServerCompletionsRequest {
  /** 模型名称（必需） */
  model?: string;
  /** 对话消息列表（必需） */
  messages?: CommonMessage[];
  /**
   * 随机性控制，0-2，默认1
   * @format double
   */
  temperature?: number;
  /**
   * 最大生成token数
   * @format int32
   */
  maxTokens?: number;
  /**
   * 核采样，0-1，默认1
   * @format double
   */
  topP?: number;
  /**
   * 生成多少个选择，默认1
   * @format int32
   */
  n?: number;
  /** 是否流式返回，默认false */
  stream?: boolean;
  /** 停止序列 */
  stop?: string;
  /** 停止序列列表 */
  stopSequence?: string[];
  /**
   * 存在惩罚，-2.0到2.0，默认0
   * @format double
   */
  presencePenalty?: number;
  /**
   * 频率惩罚，-2.0到2.0，默认0
   * @format double
   */
  frequencyPenalty?: number;
  /**
   * 重复惩罚
   * @format double
   */
  repetitionPenalty?: number;
  /** 用户标识符 */
  user?: string;
  /** 可用工具列表 */
  tools?: ServerChatCompletionTool[];
  /** 工具选择策略 */
  toolChoice?: ServerToolChoice;
  /** 响应格式 */
  responseFormat?: ServerResponseFormat;
  /**
   * 随机种子
   * @format int32
   */
  seed?: number;
  /** 额外元数据 */
  metadata?: Record<string, string>;
  /**
   * 以下为业务字段
   * 是否为非标准对话
   */
  nonStandard?: boolean;
  /** 对话id */
  chatUuid?: string;
}

/** ChatCompletionChunk 与 OpenAI Chat Completions 流式 chunk 对齐 */
export interface ServerCompletionsResponse {
  /** 响应唯一标识符 */
  id?: string;
  /** 对象类型，通常为"chat.completion" */
  object?: string;
  /**
   * 创建时间戳
   * @format int64
   */
  created?: string;
  /** 使用的模型名称 */
  model?: string;
  /** 生成的选择列表 */
  choices?: ServerChatCompletionChoice[];
  /** token使用情况 */
  usage?: CommonTokenUsage;
  /** 系统指纹 */
  systemFingerprint?: string;
  /** 额外元数据 */
  metadata?: Record<string, string>;
  /**
   * 以下为业务字段
   * 是否为非标准对话
   */
  nonStandard?: boolean;
  /** 对话id */
  chatUuid?: string;
}

/** DeleteChatMessageResponse 删除对话消息响应 */
export interface ServerDeleteChatMessageResponse {
  /** 删除是否成功 */
  success?: boolean;
  /** 响应消息 */
  message?: string;
}

/** GetChatMessagesResponse 获取对话消息响应 */
export interface ServerGetChatMessagesResponse {
  /** 消息列表 */
  messages?: CommonMessage[];
  /**
   * 总消息数量
   * @format int32
   */
  totalCount?: number;
}

export interface ServerListChatsFilter {
  tag?: string;
  /** 搜索关键字，为空则不搜索 */
  keyword?: string;
}

/** ListChatsRequest 获取对话列表请求 */
export interface ServerListChatsRequest {
  /** @format int64 */
  offset?: string;
  /** @format int64 */
  limit?: string;
  filter?: ServerListChatsFilter;
}

/** ListChatsResponse 获取对话列表响应 */
export interface ServerListChatsResponse {
  /** 对话列表 */
  chats?: CommonChatInfo[];
  /**
   * 总数量
   * @format int64
   */
  totalCount?: string;
}

/** ResponseFormat 响应格式 */
export interface ServerResponseFormat {
  /** "text" 或 "json_object" */
  type?: string;
  /** JSON Schema定义（当type为json_object时） */
  jsonSchema?: string;
}

/** ToolChoice 工具选择策略 */
export interface ServerToolChoice {
  /** "none", "auto", "required" */
  mode?: string;
  /** 指定特定工具 */
  named?: ServerChatCompletionNamedToolChoice;
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
 * @title rpc/service/chat.proto
 * @version version not set
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  v1 = {
    /**
     * No description
     *
     * @tags Chat
     * @name ChatCompletions
     * @summary 流式 Chat Completions（等价于 OpenAI /v1/chat/completions + stream）
     * @request POST:/v1/chat/completions
     */
    chatCompletions: (
      body: ServerCompletionsRequest,
      params: RequestParams = {},
    ) =>
      this.request<
        {
          result?: ServerCompletionsResponse;
          error?: RpcStatus;
        },
        RpcStatus
      >({
        path: `/v1/chat/completions`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatDeleteChat
     * @summary 删除对话
     * @request POST:/v1/chats/delete/{chatUuid}
     */
    chatDeleteChat: (
      chatUuid: string,
      body: ChatDeleteChatBody,
      params: RequestParams = {},
    ) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/chats/delete/${chatUuid}`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatListChats
     * @summary 获取历史对话记录列表
     * @request POST:/v1/chats/get
     */
    chatListChats: (body: ServerListChatsRequest, params: RequestParams = {}) =>
      this.request<ServerListChatsResponse, RpcStatus>({
        path: `/v1/chats/get`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatDeleteChatMessage
     * @summary 删除对话消息
     * @request POST:/v1/chats/messages/delete/{chatUuid}/{messageId}
     */
    chatDeleteChatMessage: (
      chatUuid: string,
      messageId: string,
      body: ChatDeleteChatMessageBody,
      params: RequestParams = {},
    ) =>
      this.request<ServerDeleteChatMessageResponse, RpcStatus>({
        path: `/v1/chats/messages/delete/${chatUuid}/${messageId}`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatGetChatMessages
     * @summary 获取对话消息
     * @request POST:/v1/chats/messages/get/{chatUuid}
     */
    chatGetChatMessages: (
      chatUuid: string,
      body: ChatGetChatMessagesBody,
      params: RequestParams = {},
    ) =>
      this.request<ServerGetChatMessagesResponse, RpcStatus>({
        path: `/v1/chats/messages/get/${chatUuid}`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatChatTitleSave
     * @summary 保存对话标题
     * @request POST:/v1/chats/title/save/{chatUuid}/{chatTitle}
     */
    chatChatTitleSave: (
      chatUuid: string,
      chatTitle: string,
      body: ChatChatTitleSaveBody,
      params: RequestParams = {},
    ) =>
      this.request<CommonEmpty, RpcStatus>({
        path: `/v1/chats/title/save/${chatUuid}/${chatTitle}`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Chat
     * @name ChatChatTitle
     * @summary 获取对话标题
     * @request POST:/v1/chats/title/{chatUuid}
     */
    chatChatTitle: (
      chatUuid: string,
      body: ChatChatTitleBody,
      params: RequestParams = {},
    ) =>
      this.request<ServerChatTitleResponse, RpcStatus>({
        path: `/v1/chats/title/${chatUuid}`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
