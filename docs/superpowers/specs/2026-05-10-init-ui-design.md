# Lemontea 初始界面设计文档

## 目标

快速实现一个可操作的前端界面，全部使用 mock 数据写死，让产品可以被看到并交互。不接真实后端。

---

## 技术栈

| 用途 | 方案 |
|---|---|
| 样式 + 组件基础 | Tailwind CSS + shadcn/ui |
| 状态管理 | Zustand |
| i18n | i18next + react-i18next |
| Markdown 渲染 | react-markdown + remark-gfm + rehype-highlight |
| 面板拖拽分割 | react-resizable-panels |
| 图标 | lucide-react |

---

## 全局约束

- 深色 / 浅色主题：通过 Tailwind `dark:` 变体实现，主题状态存入 Zustand
- 字体大小：5 档（极小 / 小 / 标准 / 大 / 超大），在 `<html>` 上设置 `font-size` 基准值，组件使用 `rem` 单位
- 国际化：简体中文 + 英文，语言状态存入 Zustand，通过 i18next 动态切换

---

## 整体布局

```
┌──────────────────────────────────────────────┐
│  Sidebar (可拖拽调宽)  │  Chat Area           │
│  min: 200px            │  flex: 1             │
│  max: 400px            │                      │
└──────────────────────────────────────────────┘
```

使用 `react-resizable-panels` 实现左右两栏，支持拖拽边界改变宽度。

---

## 组件树

```
App
├── ThemeProvider        # 同步 Zustand 主题到 html class
├── I18nProvider         # i18next 初始化
└── MainLayout
    ├── Sidebar
    │   ├── SidebarHeader        # icon + "lemontea" 品牌区
    │   ├── ConversationList     # 新建按钮、tab、搜索、列表
    │   └── SidebarFooter        # 设置按钮 + 弹出菜单
    └── ChatArea
        ├── ChatHeader           # 对话标题（可编辑）
        ├── ChatMessages         # 消息列表
        │   ├── WelcomeScreen    # 新对话时显示
        │   └── MessageList      # 已有消息时显示
        └── ChatInput            # 输入框 + 功能栏
```

---

## Sidebar 设计

### SidebarHeader

- App icon（SVG）+ 文字 "lemontea"
- 鼠标悬停 icon 触发 CSS `bounce` 动画（Tailwind `hover:animate-bounce`）

### ConversationList

**布局（从上到下）：**
1. 新建对话按钮（`+ 新对话`）
2. Tab 切换：对话 / 收藏（toggle group，两态）
3. 搜索输入框
4. 对话分组列表（今天 / 昨天 / 过去7天 / 更久以前）
5. 底部文字："已加载全部聊天记录（x 条）"

**对话列表项状态：**
- 默认：标题 + 时间预览
- Hover：右侧出现 `...` 菜单按钮
- `...` 菜单：收藏/取消收藏、重命名、删除
- Loading 中（AI 输出）：左侧 loading spinner
- Loading 完成（用户在其他对话）：左侧绿色圆点
- Loading 出错（用户在其他对话）：左侧红色圆点
- 需要确认（用户在其他对话）：左侧蓝色圆点
- 点击进入：右侧聊天区先显示 loading，加载完显示内容

**Mock 数据结构：**
```ts
interface Conversation {
  id: string
  title: string
  createdAt: Date
  updatedAt: Date
  starred: boolean
  status: 'idle' | 'loading' | 'done' | 'error' | 'waiting'
  unread: boolean
}
```

### SidebarFooter

- 设置按钮（齿轮图标）
- 点击弹出上方菜单，含三项：
  - 主题（Hover 展开子菜单：自动 / 浅色 / 深色）
  - 关于（点击弹出新窗口）
  - 设置（点击弹出新窗口）

---

## ChatArea 设计

### ChatHeader

- 显示当前对话标题
- 新对话时标题为"新对话"
- 已有标题时，右侧显示编辑图标
- 点击编辑图标：标题变为输入框 + 保存 / 取消按钮

### ChatMessages

**新对话（WelcomeScreen）：**
1. App icon（居中大图）
2. 欢迎标题："欢迎使用 lemontea"
3. 副标题："你的智能助手，随时准备帮助你完成工作"
4. 快捷开始按钮列表（icon + 标题 + 简介，3-4 个 mock 按钮）

**已有消息（MessageList）：**
- AI 消息：靠左，无头像，支持 Markdown 渲染
- User 消息：靠右，圆角矩形背景框
- 时间戳：按跨天规则显示（参见文档规则表），居中灰色小字
- AI 消息下方：模型名称 + token 数（输入/输出，用图标区分）
- 思维链：正文未输出时默认展开，正文出现后自动折叠，可手动切换
- 对话中底部显示 loading 动画（三点跳动）
- 用户滚动离底后显示"回到底部"浮动按钮

**自动滚动逻辑：**
- AI 输出时自动滚动到底部
- 用户手动滚动 → 暂停自动滚动 + 显示"回到底部"按钮
- 点击"回到底部" → 滚动到底 + 恢复自动滚动
- 用户发送新消息 → 重新启用自动滚动

**Mock 数据结构：**
```ts
interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  thinking?: string        // 思维链内容
  createdAt: Date
  model?: string
  inputTokens?: number
  outputTokens?: number
}
```

### ChatInput

**功能栏（从左到右）：**
1. 添加按钮（`+`）：点击选择文件（mock，仅 UI）
2. 工具选择栏：点击展开工具列表，可禁用/启用（mock 数据）
3. 模型选择栏：点击展开供应商分组的模型列表（mock 数据）
4. 发送按钮（右侧，AI 输出时变为停止按钮）

**输入框：**
- 多行文本框，支持实时 Markdown 预览渲染（切换编辑/预览模式）
- Enter 发送，Shift+Enter 换行

---

## 全局状态（Zustand）

```ts
interface AppStore {
  theme: 'auto' | 'light' | 'dark'
  fontSize: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  language: 'zh-CN' | 'en'
  currentConversationId: string | null
  conversations: Conversation[]
  messages: Record<string, Message[]>
}
```

---

## 文件结构

```
frontend/src/
├── main.tsx
├── App.tsx
├── store/
│   └── appStore.ts          # Zustand store
├── i18n/
│   ├── index.ts             # i18next 初始化
│   ├── zh-CN.ts
│   └── en.ts
├── components/
│   ├── layout/
│   │   └── MainLayout.tsx
│   ├── sidebar/
│   │   ├── Sidebar.tsx
│   │   ├── SidebarHeader.tsx
│   │   ├── ConversationList.tsx
│   │   ├── ConversationItem.tsx
│   │   └── SidebarFooter.tsx
│   ├── chat/
│   │   ├── ChatArea.tsx
│   │   ├── ChatHeader.tsx
│   │   ├── ChatMessages.tsx
│   │   ├── WelcomeScreen.tsx
│   │   ├── MessageList.tsx
│   │   ├── MessageItem.tsx
│   │   ├── ThinkingBlock.tsx
│   │   └── ChatInput.tsx
│   └── ui/                  # shadcn/ui 生成的组件
├── mock/
│   └── data.ts              # 所有 mock 数据
└── lib/
    └── utils.ts             # cn() 工具函数等
```

---

## 窗口

- 关于页面：独立 Wails 窗口，简单显示版本号、logo、作者信息
- 设置页面：独立 Wails 窗口，暂时仅显示占位内容
- 主窗口：默认 1200×800，最小 800×600

---

## 不在本次范围内

- 真实后端通信
- 持久化存储
- 文件上传真实逻辑
- 工具/模型的真实数据
- 设置页面的完整实现
