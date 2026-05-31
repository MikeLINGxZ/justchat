# 对话附件上传（多模态）设计

- 日期：2026-05-17
- 状态：草案
- 适用模块：`backend/service/agent`、`backend/pkg/agent`、`backend/service/file`、`frontend/src/components/chat`

## 目标

让用户在聊天输入区点击 Paperclip 按钮，通过 `backend/service/file.SelectFile` 选择本地文件作为附件，与文本一并发送给模型。支持图片、PDF 与文本类等任意类型，单条消息可携带多个附件，发送与历史重放都按多模态消息构造。

## 非目标

- 不做拖拽上传、不做粘贴板自动识别（后续可加）。
- 不实现剪贴板截图、屏幕录制等高阶来源。
- 不做附件托管/复制，不做云端上传。
- 不在前端做尺寸校验（前端无 `fs.stat`），统一在后端校验。
- 不做模型多模态能力检测：附件总是按多模态字段发送；若目标模型不支持，由上游 SDK / Provider 报错，由现有 `agent:stream:error` 通路呈现。

## 用户故事

1. 我在输入框写了一段问题，点击 Paperclip，系统弹出文件选择器，选中一张 PNG 后，输入框上方出现该图片的 chip。
2. 我再次点击 Paperclip，又添加一个 PDF。两个 chip 并列显示。
3. 我点击某个 chip 上的“×”，对应附件被移除。
4. 我点击发送，会话立即追加一条用户消息气泡，气泡内显示我的文字和两个只读附件 chip。助手开始流式回复。
5. 我重新进入这个会话，历史消息里看到当时的 chips；如果其中一个原始文件已经被我删除，chip 会显示 “(missing)”，模型上下文里也以 `[Missing attachment: <name>]` 占位。

## 架构总览

```
ChatInput (TS)
  ├─ attachments[]: Attachment 本地态
  ├─ Paperclip → File.SelectFile() → 推断 mime/kind → push
  └─ handleSend → chatStore.sendMessage({ content, attachments, ... })

chatStore.sendMessage
  ├─ 乐观插入 user message（带 attachments）
  └─ AgentBinding.SendMessage({ ..., attachments })

backend/service/agent.SendMessage  ──►  pkg/agent.ChatHandler.SendMessage
                                          ├─ storage.CreateMessage(content, attachmentsJSON)
                                          ├─ buildUserMessage(content, atts) → model.Message(ContentParts)
                                          └─ runner.Run(ctx, ..., msg)

下一轮 / 重新打开会话：
loadSessionHistory → toModelMessage
                       └─ 若 attachments 非空 → buildUserMessage 重建多模态 Message
                            ├─ ok → AddImageFilePath / AddFilePath
                            └─ err → Content 追加 "[Missing attachment: <name>]"
```

## 关键设计决策

### 1. 仅存路径（方案 A）

附件不复制、不入库内联字节。`Message` 上新增一个 `Attachments` 字段，保存 JSON 序列化后的元数据数组。发送与历史重放都按需 `os.ReadFile` 当前磁盘上的文件。

- **优点**：DB 体积可控；用户在系统文件管理器里仍是唯一一份；改动面最小。
- **代价**：源文件被移动 / 删除会导致历史会话失去附件上下文。
- **降级**：构造 `model.Message` 时，单个附件 `os.ReadFile` 失败，则不写入该 `ContentPart`，并把 `[Missing attachment: <name>]` 追加到 `Content` 末尾，让模型知道这次对话曾经有这个附件但已不可用。

### 2. 新增字段 `Attachments`（JSON 文本列）

`data_models.Message` 增加：

```go
Attachments string `gorm:"type:text" json:"attachments"`
```

- 用 `gorm.AutoMigrate` 自动建列；老数据该列为空串，正常兼容。
- 存储格式：JSON 数组，schema 见下。空数组或空串都视为“无附件”。
- 选择字段而非新建关联表：附件是“消息生命周期内的不可变快照”，无独立查询/更新诉求，单字段足够。

附件 JSON schema（后端 `pkgAgent.Attachment` 序列化结果）：

```json
[
  { "name": "diagram.png", "path": "/Users/foo/Pictures/diagram.png", "mime": "image/png", "kind": "image" },
  { "name": "spec.pdf",    "path": "/Users/foo/Docs/spec.pdf",       "mime": "application/pdf", "kind": "file"  }
]
```

字段含义：

| 字段 | 用途 |
|------|------|
| `name` | 显示名（取自路径 `filepath.Base`，UI 也用它做 chip 文案） |
| `path` | 绝对路径；由前端从 `File.SelectFile` 拿到，原样回传 |
| `mime` | mime 类型；后端推断（`mime.TypeByExtension`），失败则 `application/octet-stream` |
| `kind` | `image` 或 `file`；后端从 mime 推断，决定走 `AddImageFilePath` 还是 `AddFilePath` |

### 3. 多模态消息构造统一在 `pkg/agent`

新增 `backend/pkg/agent/attachment.go`：

```go
type Attachment struct {
    Name string `json:"name"`
    Path string `json:"path"`
    Mime string `json:"mime"`
    Kind string `json:"kind"` // "image" | "file"
}

// MarshalAttachments 序列化为存库的 JSON 字符串；空切片返回空串。
func MarshalAttachments(atts []Attachment) (string, error)

// UnmarshalAttachments 反序列化；空串返回 nil。
func UnmarshalAttachments(s string) ([]Attachment, error)

// NormalizeAttachment 用 path 推导缺失的 name/mime/kind。
func NormalizeAttachment(a Attachment) Attachment

// BuildUserMessage 把文本 + 附件构造成 model.Message。
// 单个附件读取失败：跳过该 ContentPart，且把 "[Missing attachment: <name>]" 追加到 content。
func BuildUserMessage(content string, atts []Attachment) model.Message
```

`ChatHandler.SendMessage` 与 `toModelMessage` 都通过 `BuildUserMessage` 走，保证发送态与重放态一致。

### 4. 校验在后端

前端只做提示，后端 `SendMessage` 校验：

- `len(attachments) > 10` → 返回 `error: too many attachments (max 10)`。
- 任一文件 `os.Stat` 大小 `> 20 * 1024 * 1024` → 返回 `error: attachment <name> exceeds 20MB`。
- 找不到文件 → 返回 `error: attachment <name> not found`（创建态严格校验；历史重放则降级为占位）。

前端在 `chatStore.sendMessage` 调用前可加入轻量提示，但权威检查在后端。

### 5. 顺手修复 `SelectFile` 的两个 bug

`backend/service/file/file.go:37-51`：

- `SetTitle` 使用了双重 `i18n.TCurrent`，应改为单层；
- title 的 key 写成了 `"select.folder"`，应改为 `"select.file"`；
- `CanChooseFiles(false)` 与函数名「SelectFile」相反，应为 `CanChooseFiles(true)`。

仅修这条函数；`SelectFolder` 的同款双重 t 问题不在本次范围。

## 后端改动清单

### 数据层
- `backend/models/data_models/message.go`：新增 `Attachments string` 字段（gorm `type:text`，json `attachments`）。

### `pkg/agent`
- 新文件 `attachment.go`：定义 `Attachment` 与上述四个工具函数。
- `chat_handler.go`：
  - `SendMessageParams` 增加 `Attachments []Attachment`。
  - `SendMessage`：调用 `MarshalAttachments`，把 JSON 写入 `data_models.Message.Attachments`；构造 `model.Message` 改用 `BuildUserMessage(params.Content, params.Attachments)`；保留 `Content` 字段不变以兼容旧逻辑（如 title 生成）。
  - `toModelMessage`：当 `stored.Role == "user"` 且 `stored.Attachments != ""` 时，反序列化并走 `BuildUserMessage` 重建；否则维持现有行为。
- 测试 `attachment_test.go`：覆盖
  - 单张图片正常构造（mock 真实小 PNG 文件）
  - 单个 PDF 正常构造
  - 多附件混合构造
  - 文件缺失降级（占位文本追加，且其余 part 仍正常）
  - 空 attachments 与现行 `NewUserMessage` 行为一致

### `service/agent`
- `agent_dto/send_message.go`：`SendMessageInput` 增加 `Attachments []AttachmentInput`。
- 新文件 `agent_dto/attachment_input.go`：

  ```go
  type AttachmentInput struct {
      Path string `json:"path"`
      Name string `json:"name,omitempty"`
      Mime string `json:"mime,omitempty"`
  }
  ```
- `agent.go SendMessage`：把 `input.Attachments` 经 `NormalizeAttachment` 转换为 `pkgAgent.Attachment` 后传入，并在传入前应用上面的校验（数量、大小、存在）。
- 测试 `dto_test.go`：补一组 `SendMessageInput` 序列化用例，确认 `attachments` JSON tag 正确。
- 测试 `agent_test.go`：覆盖校验路径（过多 / 过大 / 找不到 → 返回错误）与正常路径下消息正确入库。

### `service/file`
- `file.go` `SelectFile`：修复 i18n key 与 `CanChooseFiles(true)` 与双重 t（见上面 §5）。

## 前端改动清单

### 类型
`frontend/src/types/index.ts`：

```ts
export type AttachmentKind = 'image' | 'file'
export type Attachment = {
  path: string
  name: string
  mime: string
  kind: AttachmentKind
}
export type Message = { ...existing..., attachments?: Attachment[] }
```

`loadMessages` / streaming 流：把 backend 的 `message.attachments` 字符串 `JSON.parse` 为 `Attachment[]`；解析失败则忽略并打 warn。

### 工具函数
新文件 `frontend/src/lib/attachments.ts`：

- `inferAttachmentMeta(path: string): Attachment` —— 根据扩展名推断 mime + kind + name。
- 常量 `ATTACHMENT_MAX_COUNT = 10`，`ATTACHMENT_MAX_BYTES = 20 * 1024 * 1024`。
- `isImageMime(mime: string): boolean`、`isPdfMime`、`isTextMime` 用于 chip 图标分派。

### Store
`frontend/src/store/chatStore.ts`：
- `sendMessage` 参数追加 `attachments: Attachment[]`。
- 乐观插入的 user message 把 `attachments` 写到新的 `attachments` 字段；`extra` 字段不动（仍由 tool 流程使用）。
- `AgentBinding.SendMessage` 调用透传 `attachments`。
- 状态完成后 `loadMessages` 已经从后端拿回 `attachments`，前端 store 替换即可。

### ChatInput
`frontend/src/components/chat/ChatInput.tsx`：
- 增加 `useState<Attachment[]>([])`。
- Paperclip 按钮加 onClick：
  ```ts
  const res = await File.SelectFile({})
  if (!res?.file_path) return
  if (attachments.length >= ATTACHMENT_MAX_COUNT) { toastWarn(...); return }
  const next = inferAttachmentMeta(res.file_path)
  setAttachments([...attachments, next])
  ```
- 在 `<EditorContent>` 上方渲染 `<AttachmentChips items={attachments} onRemove={...} variant="input" />`，并在容器有附件时显示一条细分隔线。
- `handleSend`：传入 `attachments` 到 `chatStore.sendMessage`，发送后 `setAttachments([])`。
- isEmpty 判定：`(editor.isEmpty && attachments.length === 0)` 才禁用 send。
- `currentConversationId` 切换时重置 `attachments`。

### 消息渲染
`frontend/src/components/chat/MessageItem.tsx`：
- user / text 分支：在文字下方追加 `<AttachmentChips items={attachments} variant="message" />`（只读，无 onRemove）；若 `attachments` 为空则不渲染。
- 兼容 `attachments` 不存在的旧消息。

### 新组件
`frontend/src/components/chat/AttachmentChips.tsx`：
- Props：`items: Attachment[]`、`variant: 'input' | 'message'`、`onRemove?: (index: number) => void`。
- 渲染列表：每个 chip 包含 icon + 截断 name + （input 模式）小 × 按钮。
- icon 选择：`image` → `<ImageIcon />`；mime 是 PDF → `<FileText />`；其他 → `<Paperclip />`。
- 图片缩略：variant=`message` 且 `kind==='image'` 时，使用 `<img src={file://${path}}/>` 渲染 32×32 缩略图；onError 回退到 icon。Wails Webview 允许 `file://` 协议；如不可用，保留 icon 即可，不阻塞主功能。
- 主题：使用现有 token (`bg-muted`、`border-border`、`text-muted-foreground`)；字号继承外部 `text-xs` / `text-sm`，跟随极小/小/标准/大/超大五档字号缩放。

### i18n
`frontend/src/i18n/locales/zh-CN.ts` & `en.ts`：

| key | zh-CN | en |
|---|---|---|
| `input.attach` | 附加文件 | Attach file |
| `input.attachLimitCount` | 最多附加 {{max}} 个文件 | Up to {{max}} files per message |
| `input.attachLimitSize` | 单个文件不超过 {{mb}}MB | File must be at most {{mb}}MB |
| `chat.attachmentMissing` | 附件已丢失 | Attachment missing |
| `chat.attachmentRemove` | 移除附件 | Remove attachment |

后端 `backend/pkg/i18n` locales 也增加 `select.file`：选择文件 / Select file。

## 数据流细节

### 发送流（首次）
1. 前端 `handleSend` → `chatStore.sendMessage({ content, attachments, ... })`。
2. Store 乐观写入两条临时消息（user + assistant placeholder），user 带 `attachments`。
3. `AgentBinding.SendMessage` 发往 Wails；DTO `Attachments: AttachmentInput[]`。
4. `service/agent.SendMessage` 校验数量/大小/存在；通过则转换为 `pkgAgent.Attachment` 数组（`NormalizeAttachment` 补齐 name/mime/kind）。
5. `pkgAgent.ChatHandler.SendMessage`：
   - `MarshalAttachments` → 写入 `Message.Attachments`；
   - `BuildUserMessage(content, atts)` → `model.Message`（含 `ContentParts`）；
   - `runner.Run(ctx, ..., msg)`。
6. 错误（校验失败 / runner.Run 失败）：现有错误链上抛 → 前端通过 `agent:stream:error` 或 SendMessage 返回值感知。

### 历史重放
1. `loadSessionHistory` → `ListMessagesForSession` 读出 `data_models.Message`。
2. `toModelMessage`：user + text 分支扩展为
   ```go
   atts, _ := UnmarshalAttachments(stored.Attachments)
   if len(atts) == 0 {
       return model.NewUserMessage(stored.Content), true
   }
   return BuildUserMessage(stored.Content, atts), true
   ```
3. `BuildUserMessage` 内对每个附件 `os.ReadFile`：
   - 图片：`AddImageFilePath(path, "auto")`；
   - 其他：`AddFilePath(path)`；
   - 任一失败：跳过该 ContentPart，向 `Content` 追加 `\n[Missing attachment: <name>]`，并最终用 `Content + ContentParts` 一起返回。
4. title 生成 (`agent.go GenerateTitle`) 只看 `Content`，不受影响。

### 前端渲染
1. `loadMessages` 把后端 `MessageItem.attachments` 字符串解析为 `Attachment[]`，写入 store。
2. `MessageItem` 在 user/text 分支末尾渲染 `<AttachmentChips>`，与文本同气泡。

## 边界与错误处理

| 场景 | 处理 |
|---|---|
| 用户取消文件选择 | `File.SelectFile` 返回空路径 → 不改动 attachments。 |
| 用户多次选择同一文件 | 允许（不去重）；产品决策可后续收紧。 |
| `inferMimeType` 失败 | `mime` 落到 `application/octet-stream`，`kind` 为 `file`。 |
| 模型不支持多模态 | Provider/SDK 报错 → `agent:stream:error` 走现有错误通路。 |
| 文件超过 20MB | 后端 SendMessage 校验失败；前端把错误信息提示在输入区下方（沿用 `agent:stream:error` 渲染或 `sendMessage` 返回错误处理）。 |
| 历史文件缺失 | `BuildUserMessage` 注入占位文本；UI chip 显示 missing 状态。 |
| Wails Webview 不支持 `file://` | 缩略图回退到 icon；功能不阻塞。 |
| 主题 / 字号 | chip 用现有 design token；继承父容器字号；浅深主题均使用 `bg-muted / text-foreground`。 |

## 测试计划

### 后端
- `pkgAgent.MarshalAttachments` / `UnmarshalAttachments` 往返一致。
- `BuildUserMessage`：
  - 单 image / 单 file / 混合 → ContentParts 顺序正确。
  - 缺失文件 → `Content` 含 `[Missing attachment: ...]`，且 ContentParts 不含该项。
  - 空 attachments → 等价于 `NewUserMessage`。
- `service/agent.SendMessage` 校验路径：
  - 11 个附件 → 错误。
  - 大于 20MB → 错误。
  - 不存在 → 错误。
  - 正常入库的 `Attachments` 字段是合法 JSON 并能反序列化回来。
- `ChatHandler.SendMessage` 在 history 重放时使用 `BuildUserMessage`（通过 mock storage + 现有 chat_handler_test 模式扩展）。

### 前端
- 单测 `__tests__/chatInputAttachments.test.tsx`：
  - mock `File.SelectFile` 返回路径 → 点 Paperclip → chip 出现。
  - chip 上点 × → 该附件被移除。
  - 文本为空但附件 ≥1 → Send 启用。
  - send 时 `chatStore.sendMessage` 收到含 attachments 的参数；之后输入与 attachments 都被清空。
- 单测 `__tests__/messageAttachments.test.tsx`：渲染带 attachments 的 user 消息，chip 列表存在；attachments=undefined 时不渲染。
- `inferAttachmentMeta` 单测：常见扩展名 → 正确 mime/kind/name。

### 手测清单
- 浅 / 深主题 × 极小 / 小 / 标准 / 大 / 超大 字号下，chip 在 ChatInput 与 user 气泡内均显示合理，无溢出。
- 中文 / 英文界面下，所有新增 i18n key 落地无 fallback。
- 发送：图片 + PDF + .txt 混合一条消息，模型返回后历史能复现 chip。
- 关闭并重开应用 → 历史会话仍能展示 chips。
- 删除磁盘上某个附件源文件 → 重新进入会话后 chip 显示 missing，对话能继续，不崩溃。
- 错误路径：附件 > 20MB / 数量 > 10 → 前端能看到明确错误信息。

## 影响面汇总

| 区域 | 文件 | 变更 |
|---|---|---|
| 数据模型 | `backend/models/data_models/message.go` | 加 `Attachments` 字段 |
| Storage | — | 无；走 AutoMigrate |
| 后端 pkg | `backend/pkg/agent/attachment.go`（新）、`chat_handler.go` | 多模态构造 + 序列化 |
| 后端 service | `backend/service/agent/agent.go`、`agent_dto/send_message.go`、`agent_dto/attachment_input.go`（新） | DTO 与校验 |
| 后端 file | `backend/service/file/file.go` | SelectFile bug 修复 |
| 后端 i18n | `backend/pkg/i18n/locales/*` | `select.file` |
| 前端类型 | `frontend/src/types/index.ts` | `Attachment` / `Message.attachments` |
| 前端工具 | `frontend/src/lib/attachments.ts`（新） | mime/kind 推断 |
| 前端 store | `frontend/src/store/chatStore.ts` | 透传 attachments |
| 前端组件 | `ChatInput.tsx`、`MessageItem.tsx`、`AttachmentChips.tsx`（新） | UI 渲染 |
| 前端 i18n | `frontend/src/i18n/locales/zh-CN.ts`、`en.ts` | 新增 key |

## 开放问题（实现期可定）

- chip 文件名的最大显示长度（建议 24 字符，超出 ellipsis）。
- `AttachmentChips` 在消息气泡内是放在文字上方还是下方？建议下方，使文字仍是视觉重心。
- 是否提供「点击 chip 在系统中打开」功能？本次不做，后续可加。
