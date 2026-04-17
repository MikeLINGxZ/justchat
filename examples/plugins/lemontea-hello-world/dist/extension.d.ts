/**
 * Lemon Tea Desktop 示例插件 —— Hello World
 *
 * 演示插件系统的核心能力：
 * 1. 工具注册（random-joke / dice-roll）
 * 2. Agent 注册（greeting-agent）
 * 3. Hook 拦截（onBeforeChat / onAfterChat）
 * 4. 持久化存储（统计对话次数）
 */
interface Disposable {
    dispose(): void;
}
interface ToolDefinition {
    id: string;
    description: string;
    parameters?: Record<string, any>;
    execute: (params: any) => Promise<{
        content: string;
    }>;
}
interface AgentDefinition {
    id: string;
    name: string;
    description: string;
    systemPrompt: string;
    tools?: string[];
    role?: string;
}
interface LemonTeaAPI {
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
export declare function activate(api: LemonTeaAPI): Promise<void>;
export declare function deactivate(): void;
export {};
