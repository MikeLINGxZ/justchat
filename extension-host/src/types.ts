export interface Disposable {
  dispose(): void;
}

export interface HookChatMessage {
  role: string;
  content?: string;
  reasoning_content?: string;
  user_input_multi_content?: any[];
  multi_content?: any[];
  tool_calls?: any[];
  [key: string]: any;
}

export interface HookChatContext {
  messages: HookChatMessage[];
  agentId?: string;
  tools?: string[];
  response?: string;
}

export interface ToolDefinition {
  id: string;
  description: string;
  parameters?: Record<string, any>;  // JSON Schema
  execute: (params: any) => Promise<{ content: string }>;
}

export interface AgentDefinition {
  id: string;
  name: string;
  description: string;
  systemPrompt: string;
  tools?: string[];
  role?: string;
  hooks?: {
    onBeforeChat?: (ctx: HookChatContext) => Promise<HookChatContext>;
    onAfterChat?: (ctx: HookChatContext) => Promise<HookChatContext>;
  };
}

export interface LemonTeaAPI {
  tools: {
    register(tool: ToolDefinition): void;
  };
  agents: {
    register(agent: AgentDefinition): void;
  };
  hooks: {
    onBeforeChat(handler: (ctx: HookChatContext) => Promise<HookChatContext>): Disposable;
    onAfterChat(handler: (ctx: HookChatContext) => Promise<HookChatContext>): Disposable;
  };
  ui: {
    postMessage(viewId: string, data: any): void;
    onMessage(viewId: string, handler: (msg: any) => void): Disposable;
    renderChatCard(cardId: string, data: any): void;
  };
  storage: {
    get(key: string): Promise<any>;
    set(key: string, value: any): Promise<void>;
    delete(key: string): Promise<void>;
  };
}

export interface PluginModule {
  activate(api: LemonTeaAPI): void | Promise<void>;
  deactivate?(): void | Promise<void>;
}

export interface PluginInstance {
  id: string;
  dir: string;
  module: PluginModule;
  tools: Map<string, ToolDefinition>;
  beforeChatHooks: Array<(ctx: HookChatContext) => Promise<HookChatContext>>;
  afterChatHooks: Array<(ctx: HookChatContext) => Promise<HookChatContext>>;
  messageHandlers: Map<string, Array<(msg: any) => void>>;
}
