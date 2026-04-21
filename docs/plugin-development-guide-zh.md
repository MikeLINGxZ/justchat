# Lemon Tea Desktop 插件开发指南

## 目录

- [概述](#概述)
- [快速开始](#快速开始)
- [插件结构](#插件结构)
- [Manifest 配置](#manifest-配置)
- [插件 API](#插件-api)
  - [工具 API](#工具-api)
  - [Agent API](#agent-api)
  - [Hook API](#hook-api)
  - [UI API](#ui-api)
  - [存储 API](#存储-api)
- [插件生命周期](#插件生命周期)
- [安装与调试](#安装与调试)
- [最佳实践](#最佳实践)
- [完整示例](#完整示例)
- [API 参考](#api-参考)

---

## 概述

Lemon Tea Desktop 支持通过插件扩展应用功能。插件使用 JavaScript/TypeScript 编写，运行在独立的 Node.js Extension Host 进程中，通过 JSON-RPC 2.0 协议与主应用通信。

插件可以：

- **注册工具** —— 供 Agent 在对话中调用
- **注册 Agent** —— 带有自定义 prompt 和工具集的独立 Agent
- **拦截消息** —— 在对话前后通过 Hook 修改或记录消息
- **扩展 UI** —— 在侧边栏、聊天卡片、设置页和独立页面中注入自定义界面
- **持久化数据** —— 每个插件拥有独立的键值存储空间

---

## 快速开始

### 1. 创建插件目录

```bash
mkdir my-plugin && cd my-plugin
```

### 2. 初始化 package.json

```json
{
  "name": "lemontea-my-plugin",
  "version": "1.0.0",
  "displayName": "My Plugin",
  "description": "我的第一个 Lemon Tea 插件",
  "main": "dist/extension.js",
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.4.0"
  },
  "lemontea": {
    "engine": ">=0.1.0",
    "activationEvents": ["onStartup"],
    "contributes": {
      "tools": [],
      "agents": [],
      "views": {
        "sidebar": [],
        "chatCard": [],
        "settings": [],
        "page": []
      },
      "hooks": {
        "onBeforeChat": false,
        "onAfterChat": false
      }
    }
  }
}
```

### 3. 初始化 TypeScript

```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "outDir": "dist",
    "rootDir": "src",
    "strict": true,
    "esModuleInterop": true,
    "declaration": true,
    "sourceMap": true
  },
  "include": ["src/**/*"]
}
```

### 4. 编写插件入口

```typescript
// src/extension.ts
export async function activate(api) {
  // 在这里注册工具、Agent、Hook 等
  api.tools.register({
    id: 'hello',
    description: '打个招呼',
    execute: async (params) => {
      return { content: 'Hello from my plugin!' };
    },
  });
}

export function deactivate() {
  // 清理资源
}
```

### 5. 构建与安装

```bash
cnpm install
cnpm run build
```

打开 Lemon Tea Desktop → 设置 → 插件 → 安装插件 → 选择插件文件夹。

---

## 插件结构

```
my-plugin/
├── package.json          # 插件清单（必须包含 lemontea 字段）
├── tsconfig.json         # TypeScript 配置
├── src/
│   └── extension.ts      # 插件入口源码
├── dist/
│   └── extension.js      # 编译后的入口（main 字段指向此文件）
└── node_modules/         # 依赖（安装时自动创建）
```

插件入口文件必须导出：

| 导出       | 类型                                       | 必须 | 说明           |
| ---------- | ------------------------------------------ | ---- | -------------- |
| `activate` | `(api: LemonTeaAPI) => void \| Promise<void>` | 是   | 插件激活时调用 |
| `deactivate` | `() => void \| Promise<void>`              | 否   | 插件停用时调用 |

---

## Manifest 配置

`package.json` 中的 `lemontea` 字段是插件的核心配置：

```jsonc
{
  "name": "lemontea-my-plugin",       // 插件唯一标识（必须）
  "version": "1.0.0",                  // 语义化版本（必须）
  "displayName": "My Plugin",          // 显示名称
  "description": "插件描述",            // 插件描述
  "main": "dist/extension.js",         // 入口文件路径（必须）
  "lemontea": {
    "engine": ">=0.1.0",               // 兼容的引擎版本（必须）
    "activationEvents": [              // 激活时机
      "onStartup"                      // 应用启动时激活
    ],
    "contributes": {                   // 声明插件能力
      "tools": [...],                  // 工具声明
      "agents": [...],                 // Agent 声明
      "views": {...},                  // UI 视图声明
      "hooks": {...}                   // Hook 声明
    }
  }
}
```

### activationEvents

| 事件                      | 说明                   |
| ------------------------- | ---------------------- |
| `onStartup`               | 应用启动时立即激活     |
| `onCommand:<commandId>`   | 指定命令触发时激活     |
| `onAgent:*`               | 任意 Agent 调用时激活  |

### contributes.tools

在 manifest 中声明工具的元信息，运行时通过 `api.tools.register()` 注册执行逻辑：

```jsonc
{
  "tools": [
    {
      "id": "web-search",                // 工具 ID（插件内唯一）
      "name": "Web Search",              // 显示名称
      "description": "搜索网页",          // 描述
      "parameters": {                    // 参数 JSON Schema（可选）
        "type": "object",
        "properties": {
          "query": {
            "type": "string",
            "description": "搜索关键词"
          }
        },
        "required": ["query"]
      }
    }
  ]
}
```

> **注意**：运行时 ID 会被自动加上命名空间前缀 `plugin:<pluginName>:<toolId>`，避免与内置工具冲突。

### contributes.agents

```jsonc
{
  "agents": [
    {
      "id": "research-agent",
      "name": "Research Agent",
      "description": "专注于网页研究的 Agent"
    }
  ]
}
```

### contributes.views

```jsonc
{
  "views": {
    "sidebar": [                        // 侧边栏面板
      {
        "id": "search-panel",
        "name": "搜索",
        "icon": "search",
        "entry": "dist/ui/sidebar.html"
      }
    ],
    "chatCard": [                       // 聊天消息卡片
      {
        "id": "search-result",
        "name": "搜索结果",
        "entry": "dist/ui/chat-card.html"
      }
    ],
    "settings": [                       // 设置页面板
      {
        "id": "search-settings",
        "name": "搜索设置",
        "entry": "dist/ui/settings.html"
      }
    ],
    "page": [                           // 独立页面
      {
        "id": "dashboard",
        "name": "搜索面板",
        "entry": "dist/ui/dashboard.html"
      }
    ]
  }
}
```

### contributes.hooks

```jsonc
{
  "hooks": {
    "onBeforeChat": true,               // 是否注册 before-chat hook
    "onAfterChat": true                 // 是否注册 after-chat hook
  }
}
```

---

## 插件 API

插件通过 `activate(api)` 函数接收 `LemonTeaAPI` 对象，所有能力均通过此对象访问。

### 工具 API

注册供 Agent 调用的工具。用户在聊天输入框中选择工具后，LLM 会在合适的时机自动调用。

```typescript
api.tools.register({
  id: 'web-search',
  description: '搜索网页获取信息',
  parameters: {
    type: 'object',
    properties: {
      query: { type: 'string', description: '搜索关键词' },
      limit: { type: 'number', description: '返回结果数量' },
    },
    required: ['query'],
  },
  execute: async (params) => {
    const results = await doSearch(params.query, params.limit);
    return { content: JSON.stringify(results) };
  },
});
```

**关键点：**

- `parameters` 使用 JSON Schema 格式描述参数，LLM 据此生成正确的调用参数
- `execute` 必须返回 `{ content: string }`，content 是工具的文本输出
- 插件工具**默认需要用户确认**后才会执行（安全机制）
- 工具在聊天输入框的工具选择器中可见，用户需手动勾选启用

### Agent API

注册自定义 Agent，与内置 Agent 平级显示。

```typescript
api.agents.register({
  id: 'research-agent',
  name: 'Research Agent',
  description: '专注于网页研究和总结的 Agent',
  systemPrompt: `你是一个研究助手。你的任务是...`,
  tools: ['web-search', 'summarize'],   // 此 Agent 可用的工具 ID
  role: 'worker',                        // 角色：main | worker
});
```

**关键点：**

- 注册的 Agent 出现在聊天输入框的 Agent 选择器中
- `tools` 数组引用的是插件内的工具 ID（不需要加命名空间前缀）
- `systemPrompt` 定义 Agent 的行为和个性
- 用户选择 Agent 后，该 Agent 的 prompt 和工具集生效

### Hook API

在对话前后拦截和处理消息。Hook 自动对所有对话生效，无需用户操作。

```typescript
// Before Chat —— 在 LLM 处理前修改消息
const disposable = api.hooks.onBeforeChat(async (ctx) => {
  // ctx.messages: 消息列表
  // ctx.agentId: 当前 Agent ID
  // ctx.tools: 当前工具列表
  // 带附件或多模态输入时，消息里还会包含 reasoning_content、
  // user_input_multi_content、multi_content、tool_calls 等字段。
  // 如果你复制/改写消息对象，请保留这些字段，不要只返回 role/content。

  // 示例：在消息前注入日期信息
  ctx.messages.unshift({
    role: 'system',
    content: `当前日期: ${new Date().toISOString()}`,
  });

  return ctx; // 必须返回 ctx（可修改后返回）
});

// After Chat —— 在 LLM 回复后处理结果
api.hooks.onAfterChat(async (ctx) => {
  // ctx.response: LLM 的回复内容

  // 示例：记录日志
  console.log(`Agent 回复了 ${ctx.response.length} 个字符`);

  return ctx;
});

// 清理：停用时调用 dispose
disposable.dispose();
```

**关键点：**

- Hook 返回 `Disposable` 对象，调用 `.dispose()` 可注销
- 多个插件的同类 Hook 按激活顺序执行（管道模式），前一个的输出是后一个的输入
- Hook 执行有超时限制，超时或出错会被跳过（不会阻塞对话）
- `onBeforeChat` 可以修改消息内容；`onAfterChat` 适合做日志/统计
- 对于带文件、图片、音频等多模态消息，务必保留 `user_input_multi_content` 等原始字段

### UI API

与插件前端 UI（iframe）通信。

```typescript
// 向前端 iframe 发送数据
api.ui.postMessage('search-panel', {
  type: 'searchResults',
  data: results,
});

// 接收前端 iframe 的消息
const disposable = api.ui.onMessage('search-panel', (message) => {
  if (message.type === 'doSearch') {
    performSearch(message.query);
  }
});

// 在聊天流中渲染自定义卡片
api.ui.renderChatCard('search-result', {
  query: 'hello world',
  results: [...],
});
```

**前端 iframe 侧：**

插件 UI 是普通的 HTML 页面，通过 `postMessage` 与后端通信：

```html
<script>
  // 接收后端发来的数据
  window.addEventListener('message', (event) => {
    if (event.data.type === 'searchResults') {
      renderResults(event.data.data);
    }
  });

  // 向后端发送消息
  function doSearch(query) {
    window.parent.postMessage({ type: 'doSearch', query }, '*');
  }
</script>
```

### 存储 API

每个插件拥有独立的持久化键值存储空间，底层基于 BoltDB。

```typescript
// 写入
await api.storage.set('config', { apiKey: '...', maxResults: 10 });

// 读取
const config = await api.storage.get('config');
// → { apiKey: '...', maxResults: 10 }

// 删除
await api.storage.delete('config');

// 读取不存在的键返回 null
const missing = await api.storage.get('nonexistent');
// → null
```

**关键点：**

- 值会自动 JSON 序列化/反序列化，支持对象、数组、字符串、数字等
- 每个插件的存储空间完全隔离，互不影响
- 卸载插件时，对应的存储空间会被清除

---

## 插件生命周期

```
安装 → 未激活 → 激活中 → 已激活 → 停用中 → 未激活 → 卸载
                  │                    ▲
                  │    错误/崩溃        │
                  └───────────────────►│
```

| 阶段     | 说明                                                                     |
| -------- | ------------------------------------------------------------------------ |
| **安装** | 插件文件夹复制到 `~/.lemon_tea/plugins/`，执行 `cnpm install --production` |
| **激活** | 由 `activationEvents` 触发，调用插件的 `activate(api)` 函数              |
| **已激活** | 插件正常运行，响应工具调用、Hook 等                                      |
| **停用** | 调用插件的 `deactivate()` 函数，清理资源                                 |
| **卸载** | 删除插件文件夹和存储数据                                                 |

### 错误恢复

- Extension Host 进程崩溃时，主应用会自动重启进程并重新激活插件
- 连续崩溃 3 次（30 秒内）会停止重试，避免崩溃循环
- 单个插件的错误不会影响其他插件

---

## 安装与调试

### 安装插件

**方式一：通过 UI 安装**

设置 → 插件 → 安装插件 → 选择插件文件夹

**方式二：手动安装**

```bash
cp -r my-plugin ~/.lemon_tea/plugins/lemontea-my-plugin
cd ~/.lemon_tea/plugins/lemontea-my-plugin
cnpm install --production
```

### 调试技巧

1. **查看日志**：插件中使用 `process.stderr.write()` 输出日志（stdout 被 JSON-RPC 占用）

   ```typescript
   function log(msg: string) {
     process.stderr.write(`[my-plugin] ${msg}\n`);
   }
   ```

2. **监听模式开发**：使用 `cnpm run dev` 自动编译 TypeScript

3. **热重载**：修改代码后，在设置页禁用再启用插件即可重载

4. **存储调试**：使用 `api.storage.get/set` 检查持久化数据是否正确

---

## 最佳实践

### 命名规范

- 插件包名以 `lemontea-` 开头：`lemontea-web-search`
- 工具/Agent ID 使用 kebab-case：`web-search`、`research-agent`
- 运行时 ID 会自动添加命名空间前缀，无需手动处理

### 错误处理

```typescript
api.tools.register({
  id: 'my-tool',
  description: '...',
  execute: async (params) => {
    try {
      const result = await doSomething(params);
      return { content: result };
    } catch (error) {
      // 返回错误信息而不是抛出异常
      return { content: `执行失败: ${error.message}` };
    }
  },
});
```

### 资源清理

在 `deactivate` 中清理所有注册的 Hook 和监听器：

```typescript
const disposables: Disposable[] = [];

export function activate(api) {
  disposables.push(api.hooks.onBeforeChat(async (ctx) => { ... }));
  disposables.push(api.ui.onMessage('panel', (msg) => { ... }));
}

export function deactivate() {
  for (const d of disposables) {
    d.dispose();
  }
  disposables.length = 0;
}
```

### 性能注意事项

- Hook 处理应尽量快速，避免阻塞对话流程
- 工具执行有 30 秒超时限制
- 避免在 `activate` 中做耗时操作，影响应用启动速度
- 使用 `storage` API 缓存数据，减少重复计算

---

## 完整示例

参见 `examples/plugins/lemontea-hello-world/`，该示例演示了：

- 注册两个工具（随机笑话 + 掷骰子）
- 注册一个 Agent（问候助手）
- 注册 Before/After Chat Hook（对话统计）
- 使用存储 API（持久化计数器）

---

## API 参考

### LemonTeaAPI

```typescript
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
```

### ToolDefinition

```typescript
interface ToolDefinition {
  id: string;                                        // 工具 ID
  description: string;                               // 工具描述
  parameters?: Record<string, any>;                  // 参数 JSON Schema
  execute: (params: any) => Promise<{ content: string }>;  // 执行函数
}
```

### AgentDefinition

```typescript
interface AgentDefinition {
  id: string;           // Agent ID
  name: string;         // 显示名称
  description: string;  // 描述
  systemPrompt: string; // 系统提示词
  tools?: string[];     // 可用工具 ID 列表
  role?: string;        // 角色（main | worker）
}
```

### Disposable

```typescript
interface Disposable {
  dispose(): void;      // 注销已注册的 Hook/监听器
}
```

### ChatContext（Hook 参数）

```typescript
interface ChatContext {
  messages: Array<{     // 消息列表
    role: string;       // "system" | "user" | "assistant"
    content: string;    // 消息内容
  }>;
  agentId: string;      // 当前 Agent ID
  tools?: string[];     // 当前工具列表
  response?: string;    // LLM 回复（仅 onAfterChat 中有值）
}
```
