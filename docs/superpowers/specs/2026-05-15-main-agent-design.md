# 主 Agent 设计文档

## 概述

基于 trpc-agent-go 框架实现 Lemontea 的对话引擎核心功能。本阶段目标：基础流式对话、会话管理、工具调用、后台运行状态反馈。

## 架构总览

```
Frontend (React + Zustand)
    ↓ Wails binding (RPC)        ↑ Wails Events (streaming)
Service Layer (service/agent)
    ↓
pkg/agent (Manager → Runner → LLMAgent)
    ↓                    ↓
trpc-agent-go        pkg/agent/tools/
(model, session)     (file_rw, shell, datetime, web_search, code_exec)
```

- trpc-agent-go 的 Runner + LLMAgent 作为对话引擎
- trpc-agent-go 的 session/sqlite 管理 LLM 对话上下文（传给模型的 messages）
- 我们的 storage/ 存储会话元数据和消息记录（前端展示 + 历史加载 + 分页）
- 流式输出：消费 Runner 返回的 Go channel，通过 Wails Events 推送到前端

## 数据模型

### Session（会话元数据）

```go
type Session struct {
    OrmModel
    Title   string `gorm:"type:varchar(255)"`
    Starred bool   `gorm:"type:bool;default:0"`
    Status  string `gorm:"type:varchar(50);default:'idle'"`
}
```

Status 取值：`idle` | `loading` | `done-unread` | `error-unread` | `waiting-unread`

### Message（消息记录）

```go
type Message struct {
    OrmModel
    SessionID   uint   `gorm:"index"`
    ParentID    *uint  `gorm:"index"`
    Role        string `gorm:"type:varchar(50)"`
    ContentType string `gorm:"type:varchar(50)"`
    Content     string `gorm:"type:text"`
    ModelName   string `gorm:"type:varchar(255)"`
    AgentName   string `gorm:"type:varchar(255)"`
    TokensIn    int
    TokensOut   int
    Extra       string `gorm:"type:text"`
}
```

字段说明：
- `Role`: user / assistant / tool / system
- `ContentType`: text / tool_call / tool_result / thinking / confirm_request / confirm_response
- `ParentID`: 预留消息树结构，支持未来多 agent 协作的上下文追溯
- `AgentName`: 预留多 agent 标识
- `Extra`: JSON 扩展字段

一轮对话产生的消息示例：

| # | Role | ContentType | Content |
|---|------|-------------|---------|
| 1 | user | text | "帮我查看 /tmp 目录" |
| 2 | assistant | thinking | "用户想看目录内容..." |
| 3 | assistant | tool_call | `{name: "shell", args: {cmd: "ls /tmp"}, purpose: "列出临时目录内容"}` |
| 4 | user | confirm_response | "approved" |
| 5 | tool | tool_result | "file1.txt\nfile2.log" |
| 6 | assistant | text | "/tmp 目录下有两个文件..." |

## Prompt 加载机制

通用的 prompt 管理模块 `pkg/prompt/`，支持多 agent 各自使用不同提示词。

### 加载优先级

1. 检查 `{dataDir}/prompt/{promptId}/index.md` 是否存在
2. 存在 → 读取文件内容作为 prompt
3. 不存在 → 使用代码中内置的默认 prompt（英文）

### 接口设计

```go
// pkg/prompt/prompt.go

// Load 根据 promptId 加载提示词。优先读取数据目录下的自定义文件，不存在则返回内置默认值。
func Load(promptId string) (string, error)

// Register 注册一个 promptId 对应的内置默认提示词。在 init 或启动时调用。
func Register(promptId string, defaultContent string)
```

### 使用方式

```go
// pkg/agent/manager.go 启动时注册默认 prompt
prompt.Register(prompt_id.MainAgent, defaultMainAgentPrompt)

// 创建 agent 时加载
instruction, _ := prompt.Load(prompt_id.MainAgent)
agent := llmagent.New("main",
    llmagent.WithInstruction(instruction),
)
```

### 文件路径

自定义 prompt 文件位置：`{dataDir}/prompt/main_agent/index.md`

promptId 由 `pkg/id/prompt_id/` 统一管理，新增 agent 时只需添加一个常量并注册默认 prompt。

## 后端模块

### pkg/agent/ — 核心引擎层

| 文件 | 职责 |
|------|------|
| `manager.go` | Agent 生命周期管理：创建 Runner、管理 model 实例、维护活跃会话映射 |
| `chat_handler.go` | 对话处理：接收用户消息、调用 Runner.Run()、消费 event channel、通过 Wails Events 推送流式结果、写入消息记录 |
| `stream_manager.go` | 流管理：跟踪活跃 stream、支持取消（context cancel）、管理后台运行状态 |

### pkg/agent/tools/ — 工具层

| 文件 | 类型 | 用户可见 | 需确认 |
|------|------|----------|--------|
| `registry.go` | 工具注册中心 | - | - |
| `datetime.go` | 内置 | 否 | 否 |
| `file_rw.go` | 内置 | 否 | 是 |
| `shell.go` | 内置 | 否 | 是 |
| `web_search.go` | 用户工具 | 是 | 否 |
| `code_exec.go` | 用户工具 | 是 | 是 |

#### 工具注册中心

```go
type ToolMeta struct {
    Name            string
    Description     string
    Category        string                                  // "builtin" | "user"
    RequiresConfirm bool
    FormatPurpose   func(args json.RawMessage) string       // 动态生成执行说明
}

type Registry struct { ... }
func (r *Registry) Register(meta ToolMeta, tool tool.Tool)  // 注册工具
func (r *Registry) Get(name string) (ToolMeta, tool.Tool)   // 获取工具
func (r *Registry) BuiltinTools() []tool.Tool               // 所有内置工具
func (r *Registry) UserTools() []ToolMeta                   // 用户可见工具列表
func (r *Registry) EnabledTools(enabled []string) []tool.Tool // 根据用户选择返回启用的工具
```

敏感工具的确认机制：
1. 工具调用时，chat_handler 检查 `RequiresConfirm`
2. 若需确认，暂停 stream，调用 `FormatPurpose(args)` 生成说明
3. 通过 Wails Event 发送确认请求（含工具名、参数、用途说明）
4. 前端在聊天流中渲染确认卡片（确认/拒绝/输入修改意见）
5. 用户操作后通过 `RespondToConfirm` binding 回传结果
6. 后端根据用户响应继续或中止工具执行

### service/agent/ — Wails 服务层

遵循项目 service 规范（`agent.go` + `agent_implement.go` + `agent_internal.go` + `agent_dto/`）。

公开方法：

| 方法 | 用途 |
|------|------|
| `SendMessage(input)` | 发送消息并开始流式对话 |
| `StopGeneration(input)` | 停止当前生成 |
| `RespondToConfirm(input)` | 响应工具确认（确认/拒绝/修改意见） |
| `CreateSession(input)` | 创建新会话 |
| `ListSessions(input)` | 查询会话列表（分页） |
| `LoadSessionMessages(input)` | 加载会话历史消息（分页） |
| `RenameSession(input)` | 重命名会话 |
| `DeleteSession(input)` | 删除会话 |
| `ToggleStarSession(input)` | 收藏/取消收藏 |
| `GenerateTitle(input)` | 自动生成会话标题（LLM 调用） |

### storage/ — 数据存储

| 文件 | 职责 |
|------|------|
| `session.go` | 会话 CRUD、分页查询、收藏切换、状态更新 |
| `message.go` | 消息写入、按会话分页读取 |

## 前端模块

### Wails Events（后端 → 前端）

| 事件名 | 载荷 | 用途 |
|--------|------|------|
| `agent:stream:chunk` | `{sessionId, messageId, delta, contentType}` | 流式文本/思考内容增量 |
| `agent:stream:tool_call` | `{sessionId, toolName, args, purpose}` | 工具调用信息 |
| `agent:stream:confirm_request` | `{sessionId, requestId, toolName, args, purpose}` | 敏感工具确认请求 |
| `agent:stream:tool_result` | `{sessionId, toolName, result}` | 工具执行结果 |
| `agent:stream:done` | `{sessionId, usage}` | 对话完成 |
| `agent:stream:error` | `{sessionId, error}` | 对话出错 |
| `agent:session:status` | `{sessionId, status}` | 会话状态变更 |

### chatStore.ts 改造

替换 mock 数据，对接后端：
- 会话列表：调用 `ListSessions`，支持分页
- 发送消息：调用 `SendMessage`，监听 stream events 实时更新
- 停止生成：调用 `StopGeneration`
- 会话操作：对接 `CreateSession` / `RenameSession` / `DeleteSession` / `ToggleStarSession`
- 消息加载：调用 `LoadSessionMessages`，支持分页

### 组件更新

| 组件 | 改动 |
|------|------|
| `ChatInput.tsx` | 对接真实发送/停止，移除 mock 逻辑 |
| `ChatMessages.tsx` | 渲染流式内容、工具调用卡片、确认卡片 |
| `MessageItem.tsx` | 支持 tool_call / tool_result / thinking / confirm 类型渲染 |
| `ConversationList.tsx` | 对接后端会话列表，分页加载 |
| 新增 `ToolConfirmCard.tsx` | 敏感工具确认交互组件（确认/拒绝/输入修改意见） |
| 新增 `ToolCallBlock.tsx` | 工具调用展示组件（工具名、参数、用途说明、执行结果） |

## 流式对话完整流程

```
1. 用户在 ChatInput 输入内容，点击发送
2. 前端调用 SendMessage(sessionId, content, modelId)
3. 后端 service 层转发到 chat_handler
4. chat_handler:
   a. 将用户消息写入 storage/message
   b. 构建 model.NewUserMessage，调用 runner.Run()
   c. 启动 goroutine 消费 event channel
5. 每收到一个 event:
   a. 文本 chunk → emit agent:stream:chunk
   b. thinking chunk → emit agent:stream:chunk (contentType=thinking)
   c. tool_call → 检查 RequiresConfirm
      - 不需确认 → 直接执行，emit tool_call + tool_result
      - 需要确认 → emit confirm_request，暂停等待用户响应
   d. done → 将完整消息写入 storage/message，emit agent:stream:done
   e. error → 更新会话状态，emit agent:stream:error
6. 前端监听事件，实时更新 chatStore 中的消息和会话状态
7. 用户可随时点击停止 → StopGeneration → context cancel
```

## 后台运行与状态反馈

- 用户切换到其他会话或最小化窗口时，对话在后台继续运行
- stream_manager 维护所有活跃 stream 的状态
- 会话列表标题左侧的状态指示器：
  - `loading` → loading 动画（正在生成）
  - `error-unread` → 红点（错误停止）
  - `waiting-unread` → 蓝点（等待用户确认）
  - `done-unread` → 绿点（正常完成）
  - `idle` → 无指示器
- 用户点击进入会话时，自动标记为已读（`idle`）

## 目录结构

```
backend/
  models/data_models/
    session.go                    # Session 数据模型
    message.go                    # Message 数据模型（扩展，复用现有文件）
  pkg/
    prompt/
      prompt.go                   # 通用 prompt 加载（Register + Load）
    agent/
      manager.go                  # Agent 生命周期管理
      chat_handler.go             # 聊天处理逻辑
      stream_manager.go           # 流式输出管理
      tools/
        registry.go               # 工具注册中心
        datetime.go               # 内置：日期时间
        file_rw.go                # 内置：文件读写
        shell.go                  # 内置：Shell 命令
        web_search.go             # 用户工具：网页搜索
        code_exec.go              # 用户工具：代码执行
  service/
    agent/
      agent.go                    # Wails 服务（公开方法）
      agent_implement.go          # ServiceStartup
      agent_internal.go           # 私有辅助方法
      agent_dto/
        send_message.go
        stop_generation.go
        respond_to_confirm.go
        create_session.go
        list_sessions.go
        load_session_messages.go
        rename_session.go
        delete_session.go
        toggle_star_session.go
        generate_title.go
  storage/
    session.go                    # 会话存储操作
    message.go                    # 消息存储操作

frontend/src/
  store/
    chatStore.ts                  # 改造：对接后端 + 监听 events
  components/chat/
    ChatInput.tsx                 # 改造：对接真实发送/停止
    ChatMessages.tsx              # 改造：渲染流式内容
    MessageItem.tsx               # 改造：支持多种消息类型
    ToolConfirmCard.tsx           # 新增：工具确认卡片
    ToolCallBlock.tsx             # 新增：工具调用展示
  components/sidebar/
    ConversationList.tsx          # 改造：对接后端分页
```
