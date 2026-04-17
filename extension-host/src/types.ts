export interface Disposable {
  dispose(): void;
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
    onBeforeChat?: (ctx: any) => Promise<any>;
    onAfterChat?: (ctx: any) => Promise<any>;
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
    onBeforeChat(handler: (ctx: any) => Promise<any>): Disposable;
    onAfterChat(handler: (ctx: any) => Promise<any>): Disposable;
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
  beforeChatHooks: Array<(ctx: any) => Promise<any>>;
  afterChatHooks: Array<(ctx: any) => Promise<any>>;
  messageHandlers: Map<string, Array<(msg: any) => void>>;
}
