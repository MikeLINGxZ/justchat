# Lemon Tea Desktop 插件系统设计

## 概述

为 Lemon Tea Desktop 设计一套类似 VSCode 的插件系统，支持第三方扩展添加工具、UI 视图、自定义 Agent 和消息 Hook。插件使用 JavaScript/TypeScript 编写，运行在 Node.js Extension Host 进程中。

## 需求

- **工具扩展**：插件可以注册供 Agent 调用的工具
- **UI 扩展**：插件可以在侧边栏、聊天消息、设置页和独立页面中注入 UI
- **Agent 扩展**：插件可以注册带有自定义 prompt 和工具集的新 Agent
- **消息 Hook**：插件可以在 Agent 执行前后拦截和修改消息
- **分发方式**：优先支持本地文件夹/zip 安装；插件市场后续实现
- **开发语言**：仅支持 JavaScript/TypeScript（完整 npm 生态）
- **运行时**：Node.js Extension Host 进程（类似 VSCode）

## 1. 插件结构与 Manifest

每个插件是一个文件夹，包含以下结构：

```
my-plugin/
├── package.json          # 插件清单（扩展 npm 标准格式）
├── dist/
│   ├── extension.js      # 后端入口（编译后的 JS）
│   └── ui/               # 前端静态资源（可选）
│       ├── sidebar.html
│       ├── chat-card.html
│       └── settings.html
└── node_modules/         # 依赖
```

### package.json — 插件清单

package.json 中的 `lemontea` 字段声明插件的所有能力：

```jsonc
{
  "name": "lemontea-web-search",
  "version": "1.0.0",
  "displayName": "Web Search",
  "description": "为 Agent 添加网页搜索能力",
  "main": "dist/extension.js",
  "lemontea": {
    "engine": ">=0.1.0",
    "activationEvents": [
      "onStartup",
      "onCommand:webSearch.search",
      "onAgent:*"
    ],
    "contributes": {
      "tools": [{
        "id": "web-search",
        "name": "Web Search",
        "description": "搜索网页"
      }],
      "agents": [{
        "id": "research-agent",
        "name": "Research Agent",
        "description": "专注于网页研究的 Agent"
      }],
      "views": {
        "sidebar": [{
          "id": "search-panel",
          "name": "搜索",
          "icon": "search",
          "entry": "dist/ui/sidebar.html"
        }],
        "chatCard": [{
          "id": "search-result",
          "name": "搜索结果",
          "entry": "dist/ui/chat-card.html"
        }],
        "settings": [{
          "id": "search-settings",
          "name": "网页搜索设置",
          "entry": "dist/ui/settings.html"
        }],
        "page": [{
          "id": "search-dashboard",
          "name": "搜索面板",
          "entry": "dist/ui/dashboard.html"
        }]
      },
      "hooks": {
        "onBeforeChat": true,
        "onAfterChat": true
      }
    }
  }
}
```

核心概念：
- `activationEvents` 控制懒加载；不需要的插件不会启动
- `contributes` 声明所有能力；主应用据此注册扩展点
- UI 资源是静态 HTML/JS 文件，通过 iframe 加载

## 2. Extension Host 架构与通信协议

### 架构

Go 主应用启动一个 Node.js 子进程作为 Extension Host，通过 stdin/stdout 进行 JSON-RPC 2.0 双向通信。

```
┌─────────────────────────────────────┐
│  Go（插件管理器）                     │
│  ┌───────────┐  ┌────────────────┐  │
│  │ 生命周期   │  │ 消息路由器      │  │
│  │ 管理器     │  │ (JSON-RPC)     │  │
│  └─────┬─────┘  └───────┬────────┘  │
│        │          stdin/stdout       │
│        │                │            │
├────────┼────────────────┼────────────┤
│        ▼                ▼            │
│  Node.js Extension Host 进程        │
│  ┌──────────────────────────────┐   │
│  │  Host 运行时                  │   │
│  │  ┌────────┐ ┌─────────────┐  │   │
│  │  │API     │ │插件          │  │   │
│  │  │桥接    │ │沙箱          │  │   │
│  │  └────┬───┘ └──────┬──────┘  │   │
│  │       │   lemontea.*│         │   │
│  │  ┌────▼─────────────▼──────┐ │   │
│  │  │  插件 A  │  插件 B      │ │   │
│  │  └────────────┴────────────┘ │   │
│  └──────────────────────────────┘   │
└─────────────────────────────────────┘
```

### JSON-RPC 协议

Go 与 Extension Host 之间的双向通信：

**Go → Extension Host：**
```jsonc
// 激活插件
{"jsonrpc":"2.0","id":1,"method":"plugin/activate","params":{"pluginId":"web-search"}}

// 执行工具
{"jsonrpc":"2.0","id":2,"method":"tool/execute","params":{"pluginId":"web-search","toolId":"web-search","input":{"query":"hello"}}}

// 触发 Hook
{"jsonrpc":"2.0","id":3,"method":"hook/onBeforeChat","params":{"messages":[...],"agentId":"main"}}
```

**Extension Host → Go：**
```jsonc
// 读取应用配置
{"jsonrpc":"2.0","id":100,"method":"app/getConfig","params":{"key":"apiKey"}}

// 向前端发送事件
{"jsonrpc":"2.0","id":101,"method":"app/emitEvent","params":{"event":"searchResult","data":{...}}}

// 注册动态工具
{"jsonrpc":"2.0","id":102,"method":"app/registerTool","params":{"toolId":"web-search","schema":{...}}}
```

### Extension Host 内部结构

```typescript
class ExtensionHost {
  private plugins: Map<string, PluginInstance>
  private rpc: JsonRpcConnection  // stdin/stdout

  async activatePlugin(pluginId: string) {
    const manifest = this.loadManifest(pluginId)
    const api = this.createPluginAPI(pluginId)
    const module = require(manifest.main)
    await module.activate(api)
  }

  createPluginAPI(pluginId: string): LemonTeaAPI {
    return {
      tools: new ToolsAPI(pluginId, this.rpc),
      agents: new AgentsAPI(pluginId, this.rpc),
      hooks: new HooksAPI(pluginId, this.rpc),
      ui: new UiAPI(pluginId, this.rpc),
      storage: new StorageAPI(pluginId, this.rpc),
    }
  }
}
```

### 插件生命周期

```
安装 → 未激活 → 激活中 → 已激活 → 停用中 → 未激活 → 卸载
                  │                    ▲
                  │    错误/崩溃        │
                  └───────────────────►│
```

1. **安装**：将插件文件夹复制到 `~/.lemontea/plugins/`，执行 `npm install --production`
2. **激活**：由 `activationEvents` 触发，调用插件的 `activate(api)` 函数
3. **已激活**：插件正常运行，响应工具调用、Hook 等
4. **停用**：调用插件的 `deactivate()` 函数，清理资源
5. **卸载**：删除插件文件夹

## 3. 插件 API

插件通过 `activate(api)` 接收 API 对象：

### 3.1 工具 API

```typescript
export function activate(api: LemonTeaAPI) {
  api.tools.register({
    id: 'web-search',
    description: '搜索网页获取信息',
    parameters: {
      type: 'object',
      properties: {
        query: { type: 'string', description: '搜索关键词' }
      },
      required: ['query']
    },
    execute: async (params) => {
      const results = await searchWeb(params.query)
      return { content: JSON.stringify(results) }
    }
  })
}
```

插件注册的工具会出现在应用的工具列表中，可以分配给 Agent 使用。

### 3.2 Agent API

```typescript
api.agents.register({
  id: 'research-agent',
  name: 'Research Agent',
  description: '专注于网页研究和总结的 Agent',
  systemPrompt: 'You are a research assistant...',
  tools: ['web-search'],
  role: 'main',
  hooks: {
    onBeforeChat: async (ctx) => {
      ctx.messages.unshift({
        role: 'system',
        content: `Current date: ${new Date().toISOString()}`
      })
      return ctx
    }
  }
})
```

注册的 Agent 与内置 Agent 平级，出现在前端 Agent 列表中。

### 3.3 Hook API

```typescript
// 全局 Hook，对所有 Agent 生效
api.hooks.onBeforeChat(async (ctx) => {
  // 可用字段：ctx.messages, ctx.agentId, ctx.tools
  ctx.messages = ctx.messages.map(msg => ({
    ...msg,
    content: filterSensitiveWords(msg.content)
  }))
  return ctx
})

api.hooks.onAfterChat(async (ctx) => {
  // 可用字段：ctx.response
  console.log(`Agent ${ctx.agentId} 回复了 ${ctx.response.length} 个字符`)
  return ctx
})
```

多个插件注册同一 Hook 时，按插件激活顺序依次执行（管道模式），前一个的输出是后一个的输入。

### 3.4 UI API

```typescript
// 向前端 iframe 发送数据
api.ui.postMessage('search-panel', { type: 'searchResults', data: results })

// 接收前端 iframe 的消息
api.ui.onMessage('search-panel', (message) => {
  if (message.type === 'doSearch') {
    performSearch(message.query)
  }
})

// 在聊天流中渲染自定义卡片
api.ui.renderChatCard('search-result', { query: 'hello world', results: [...] })
```

### 3.5 存储 API

```typescript
// 每个插件有独立的键值存储空间（底层由 BoltDB bucket 实现）
await api.storage.set('lastQuery', 'hello world')
const query = await api.storage.get('lastQuery')
await api.storage.delete('lastQuery')
```

### 3.6 完整类型定义

```typescript
interface LemonTeaAPI {
  tools: {
    register(tool: ToolDefinition): void
  }
  agents: {
    register(agent: AgentDefinition): void
  }
  hooks: {
    onBeforeChat(handler: HookHandler): Disposable
    onAfterChat(handler: HookHandler): Disposable
  }
  ui: {
    postMessage(viewId: string, data: any): void
    onMessage(viewId: string, handler: (msg: any) => void): Disposable
    renderChatCard(cardId: string, data: any): void
  }
  storage: {
    get(key: string): Promise<any>
    set(key: string, value: any): Promise<void>
    delete(key: string): Promise<void>
  }
}

type Disposable = { dispose(): void }
```

## 4. 前端 UI 扩展机制

### 4.1 iframe 沙箱

所有插件 UI 通过沙箱化的 iframe 加载：

```
插件 UI 文件：~/.lemontea/plugins/web-search/dist/ui/sidebar.html
加载 URL：    lemontea-plugin://web-search/ui/sidebar.html
```

```html
<iframe
  src="lemontea-plugin://web-search/ui/sidebar.html"
  sandbox="allow-scripts allow-forms"
  style="width:100%;height:100%;border:none"
/>
```

不设置 `allow-same-origin`，插件无法访问主应用的 DOM 和存储。

### 4.2 主应用与 iframe 通信

插件 UI 引入主应用提供的轻量 SDK：

```html
<script src="lemontea-plugin://sdk/lemontea-ui.js"></script>
<script>
  const api = window.LemonTeaUI.connect()

  api.onMessage((data) => {
    if (data.type === 'searchResults') {
      renderResults(data.data)
    }
  })

  document.getElementById('btn').onclick = () => {
    api.postMessage({ type: 'doSearch', query: input.value })
  }
</script>
```

通信链路：
```
iframe (postMessage) → 主应用 React → Wails IPC → Go → JSON-RPC → Extension Host → 插件
```

### 4.3 扩展点

**侧边栏**：左侧可扩展的 tab 栏。每个插件的 sidebar view 是一个 tab，点击切换 iframe 内容。

**聊天消息卡片**：插件后端调用 `api.ui.renderChatCard()` 时，Go 端向前端发出 Wails 事件，前端在聊天流中插入 iframe 容器。卡片高度由 iframe 内容动态决定（通过 postMessage 上报高度）。

**设置页**：设置页面新增「插件」分区，列出所有声明了 settings view 的插件。点击某个插件的设置项，在右侧面板加载对应 iframe。插件通过 Storage API 读写配置。

**独立页面**：通过 React Router 注册动态路由 `/plugin/:pluginId/:pageId`。页面整体由 iframe 填充。插件可以从侧边栏或聊天卡片中跳转到独立页面。

### 4.4 前端插件注册流程

```
应用启动
  → Go 扫描 ~/.lemontea/plugins/，解析所有 manifest
  → 将插件列表 + contributes 信息通过 Wails binding 传给前端
  → 前端 React 动态：
      - 渲染侧边栏 tab
      - 注册聊天卡片渲染器
      - 添加设置菜单项
      - 注册页面路由
```

主应用前端不需要知道插件的具体逻辑——只需根据 manifest 声明挂载 iframe。

## 5. Go 端插件管理架构

### 5.1 模块划分

```
backend/
├── plugin/
│   ├── manager.go          # 插件生命周期管理（安装/激活/停用/卸载）
│   ├── manifest.go         # 解析 package.json 中的 lemontea 字段
│   ├── host.go             # Extension Host 进程管理（启动/重启/通信）
│   ├── rpc.go              # JSON-RPC 编解码与路由
│   ├── hook_chain.go       # Hook 管道管理与执行
│   ├── tool_bridge.go      # 将插件工具桥接为 Eino BaseTool
│   └── agent_bridge.go     # 将插件 Agent 桥接为 IAgent 实现
```

### 5.2 核心组件

**Manager**（管理器）：
```go
type Manager struct {
    pluginsDir  string                    // ~/.lemontea/plugins/
    manifests   map[string]*Manifest      // 已安装插件
    host        *ExtensionHost            // Node.js 进程
    hookChain   *HookChain               // Hook 管道
    toolBridge  *ToolBridge              // 工具桥接
}

func (m *Manager) Install(folderPath string) error
func (m *Manager) Uninstall(pluginId string) error
func (m *Manager) Activate(pluginId string) error
func (m *Manager) Deactivate(pluginId string) error
func (m *Manager) ListPlugins() []PluginInfo
```

**ExtensionHost**（扩展宿主）：
```go
type ExtensionHost struct {
    cmd     *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
    rpc     *JsonRpcConn
}

func (h *ExtensionHost) Start() error
func (h *ExtensionHost) Stop() error
func (h *ExtensionHost) Restart() error        // 崩溃自动重启
func (h *ExtensionHost) Call(method string, params any) (any, error)
func (h *ExtensionHost) OnRequest(handler func(method string, params any) any)
```

**ToolBridge**（工具桥接）— 将插件工具适配为 Eino 工具接口：
```go
type PluginTool struct {
    pluginId    string
    toolId      string
    name        string
    description string
    schema      map[string]any
    host        *ExtensionHost
}

func (t *PluginTool) Tool() tool.BaseTool
func (t *PluginTool) Id() string
func (t *PluginTool) RequireConfirmation() bool { return true }  // 插件工具默认需确认
```

**HookChain**（Hook 管道）：
```go
type HookChain struct {
    beforeChat []HookEntry
    afterChat  []HookEntry
}

func (c *HookChain) RunBeforeChat(ctx *ChatContext) (*ChatContext, error)
func (c *HookChain) RunAfterChat(ctx *ChatContext) (*ChatContext, error)
```

### 5.3 与现有 Service 层的集成

```go
// service.go 新增字段
type Service struct {
    // ... 现有字段
    pluginManager *plugin.Manager
}

// 在 chat 执行流程中集成
func (s *Service) chat(...) {
    ctx := &ChatContext{Messages: messages, AgentId: agentId}

    // Before Hook
    ctx, _ = s.pluginManager.HookChain().RunBeforeChat(ctx)

    // 合并插件工具与内置 + MCP 工具
    tools := s.resolveSelectedTools(...)
    tools = append(tools, s.pluginManager.GetPluginTools()...)

    // 执行 chat
    response := s.executeChat(ctx.Messages, tools)

    // After Hook
    ctx.Response = response
    s.pluginManager.HookChain().RunAfterChat(ctx)
}
```

### 5.4 安全性

- **插件工具默认需要用户确认**，复用现有的 `tool_approval` 机制
- **崩溃隔离**：Go 端监控 Extension Host 进程状态，崩溃后自动重启，已激活的插件重新激活
- **超时控制**：工具调用和 Hook 执行都有超时限制（如 30 秒），超时自动取消
- **存储隔离**：每个插件只能访问自己的 BoltDB bucket

## 6. 插件安装与管理 UI

### 6.1 安装流程

```
用户选择插件文件夹/zip
  → Go 校验 package.json 中的 lemontea 字段是否合法
  → 检查 engine 版本兼容性
  → 复制文件到 ~/.lemontea/plugins/<plugin-id>/
  → 执行 npm install --production
  → 解析 manifest，注册 contributes
  → 前端刷新插件列表
```

### 6.2 设置页 — 插件管理界面

在设置菜单中新增「插件」页（位于 Agent 和技能之后）：

- 列出所有已安装插件，显示名称、版本、描述
- 显示能力摘要（工具数、视图数、Hook 数）
- 每个插件的操作：设置、禁用、删除
- 「安装插件」按钮打开文件夹选择对话框

### 6.3 前端 Service API

```go
func (s *Service) InstallPlugin() error
func (s *Service) UninstallPlugin(pluginId string) error
func (s *Service) EnablePlugin(pluginId string) error
func (s *Service) DisablePlugin(pluginId string) error
func (s *Service) GetInstalledPlugins() []PluginViewModel
func (s *Service) GetPluginDetail(pluginId string) PluginDetailViewModel
```

### 6.4 与现有功能的关系

| 功能 | 现有机制 | 插件系统集成 |
|------|---------|-------------|
| 工具 | 内置 + MCP | 插件工具与 MCP 工具平级，统一出现在工具列表 |
| Agent | 系统 + 自定义 | 插件 Agent 与自定义 Agent 平级，统一出现在 Agent 列表 |
| 技能 | Markdown 文件 | 保持独立；插件注册 Agent 时可引用技能 |
| Hook | 无 | 插件独有的新能力 |

MCP 作为开放标准协议继续保留，插件系统是更高层的扩展机制。两者共存：
- 简单的工具扩展 → 用 MCP
- 需要 UI + Hook + Agent 的复杂扩展 → 用插件
