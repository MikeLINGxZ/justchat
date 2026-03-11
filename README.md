# Lemon Tea Desktop 🍋🍵

一款基于 **Wails3** 构建的跨平台桌面 AI 聊天客户端，支持多种大语言模型供应商，提供流畅的对话体验。

## 功能特性

- **多轮对话**：支持流式输出、停止生成、自动生成对话标题
- **对话管理**：新建 / 删除 / 重命名 / 收藏对话
- **多 LLM 供应商**：
   - Deepseek（深度求索）
   - Qwen（阿里云百炼 / 通义千问）
   - OpenRouter
   - Ollama（本地模型）
   - 任意 OpenAI 兼容接口
- **工具调用 (Tool Calling)**：基于 CloudWeGo Eino Agent 架构，支持获取当前日期/时间等工具
- **Markdown 渲染**：支持代码高亮、GFM 语法
- **个性化设置**：字体大小调节、模型供应商管理
- **文件操作**：支持选择和打开本地文件
- **跨平台**：支持 macOS、Windows、Linux（以及实验性 iOS/Android）

## 技术栈

| 层级 | 技术 |
|------|------|
| 桌面框架 | [Wails3](https://v3.wails.io/) (v3.0.0-alpha.74) |
| 后端 | Go 1.25 |
| 前端 | React 19 + TypeScript 5.8 + Vite 7 |
| UI 组件库 | Ant Design 6 |
| 样式 | Sass / SCSS Modules |
| 状态管理 | Zustand |
| 路由 | React Router 7 |
| LLM 框架 | [CloudWeGo Eino](https://github.com/cloudwego/eino) |
| 数据库 | SQLite (GORM) |
| 构建工具 | [Task](https://taskfile.dev/) |

## 项目结构

```
lemon_tea_desktop/
├── main.go                  # Go 入口，Wails3 应用启动
├── Taskfile.yml             # 任务定义（构建、运行、打包）
├── go.mod                   # Go 模块定义
│
├── frontend/                # 前端代码
│   ├── src/
│   │   ├── main.tsx         # React 入口
│   │   ├── App.tsx          # 路由与布局
│   │   ├── components/      # 组件（chat、input、message 等）
│   │   ├── pages/           # 页面（home、settings）
│   │   ├── stores/          # Zustand 状态管理
│   │   ├── hooks/           # 自定义 Hooks
│   │   └── utils/           # 工具函数
│   ├── bindings/            # Wails 自动生成的 Go → TS 绑定
│   ├── package.json
│   └── vite.config.ts
│
├── backend/                 # 后端代码
│   ├── service/             # 业务服务（chat、provider、file）
│   ├── storage/             # SQLite 数据存储
│   ├── models/              # 数据模型
│   ├── pkg/
│   │   ├── llm_provider/    # LLM 供应商适配层
│   │   └── logger/          # 日志
│   ├── agents/              # Agent 示例
│   └── utils/               # 工具函数
│
└── build/                   # 构建配置与平台资源
    ├── config.yml           # Wails3 配置
    ├── darwin/              # macOS 构建
    ├── windows/             # Windows 构建
    └── linux/               # Linux 构建
```

## 环境要求

- [Go](https://go.dev/) >= 1.25
- [Node.js](https://nodejs.org/) >= 18
- [Wails3 CLI](https://v3.wails.io/)
- [Task](https://taskfile.dev/)（任务运行器）
- 包管理器：npm / cnpm / pnpm

## 快速开始

### 1. 克隆项目

```bash
git clone <仓库地址>
cd lemon_tea_desktop
```

### 2. 安装依赖

```bash
# 安装前端依赖
cd frontend
npm install
cd ..
```

### 3. 开发模式运行

```bash
task dev
```

这将启动 Wails3 开发服务器，前端支持热重载，默认端口 `9246`。

### 4. 构建生产版本

```bash
# 构建可执行文件
task build

# 打包应用（如 macOS .app）
task package
```

## 其他命令

| 命令 | 说明 |
|------|------|
| `task dev` | 开发模式（热重载） |
| `task build` | 构建当前平台可执行文件 |
| `task package` | 打包为平台原生应用 |
| `task run` | 运行已构建的应用 |
| `task build:server` | 构建无 GUI 的 HTTP 服务模式 |
| `task run:server` | 运行服务模式 |
| `task build:docker` | 构建 Docker 镜像 |
| `task run:docker` | 构建并运行 Docker 容器 |

## 支持的 LLM 供应商

| 供应商 | API 地址 | 说明 |
|--------|---------|------|
| Deepseek | `https://api.deepseek.com/v1` | 深度求索 |
| Qwen | `https://dashscope.aliyuncs.com/compatible-mode/v1` | 阿里云百炼 |
| OpenRouter | `https://openrouter.ai/api/v1` | 多模型聚合 |
| Ollama | `http://localhost:11434/v1` | 本地模型 |
| 自定义 | 任意 OpenAI 兼容 Base URL | 自定义供应商 |

## 许可证

MIT License
