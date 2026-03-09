# Lemon Tea Desktop

基于 [Wails3](https://v3.wails.io/) 和 [CloudWeGo Eino](https://github.com/cloudwego/eino) 构建的跨平台 AI 智能体桌面客户端。

## 功能特性

- **多模型提供商支持**：深度求索 (DeepSeek)、阿里云百炼、OpenRouter、OpenAI 标准接口
- **流式对话**：支持实时流式输出，边生成边显示
- **对话记忆**：基于 Agent 的记忆与上下文管理
- **设置管理**：图形化配置 API 密钥与模型选择
- **跨平台**：支持 macOS、Windows、Linux、iOS、Android

## 技术栈

| 层级   | 技术 |
|--------|------|
| 框架   | Wails3、React 19 |
| 后端   | Go 1.25、GORM、SQLite、CloudWeGo Eino |
| 前端   | TypeScript、Vite、Ant Design、Zustand |
| 模型   | Eino-Ext（OpenAI、DeepSeek、Ollama、Qwen、Ark 等） |

## 环境要求

- Go 1.25+
- Node.js 18+
- [Wails3 CLI](https://v3.wails.io/docs/getting-started/installation)
- 至少配置一个 LLM 提供商的 API Key（在应用内设置中配置）

## 快速开始

### 1. 克隆并安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd lemon_tea_desktop

# 安装 Go 依赖
go mod tidy

# 安装前端依赖
cd frontend && npm install && cd ..
```

### 2. 开发模式运行

```bash
# 前端开发构建（如需先构建前端资源）
cd frontend && npm run build:dev && cd ..

# 启动应用（热重载）
wails3 dev
```

### 3. 生产构建

```bash
# 构建前端
cd frontend && npm run build && cd ..

# 打包应用
wails3 build
```

构建产物位于 `build` 目录。

### 4. 配置 LLM 提供商

首次运行后，进入 **设置 → 模型提供商**，添加并配置至少一个提供商（如 DeepSeek、阿里云百炼、OpenRouter 等）的 API Key 和 Base URL。

## 项目结构

```
lemon_tea_desktop/
├── main.go                 # 应用入口
├── backend/
│   ├── service/            # 业务服务（聊天、设置、提供商等）
│   ├── storage/            # 数据持久化（SQLite）
│   ├── models/             # 数据模型与视图模型
│   ├── pkg/
│   │   └── llm_provider/   # LLM 提供商实现
│   ├── agents/             # Agent 与记忆逻辑
│   └── examples/           # 示例代码
├── frontend/
│   ├── src/
│   │   ├── pages/          # 页面（home、settings、apps）
│   │   ├── components/     # 通用组件
│   │   ├── stores/         # 状态管理
│   │   └── utils/          # 工具函数
│   └── bindings/           # Wails 生成的 TypeScript 绑定
└── build/                  # 各平台构建配置
```

## 支持的提供商

| 提供商         | 类型         | 说明 |
|----------------|--------------|------|
| 深度求索       | DeepSeek     | 通用大模型 |
| 阿里云百炼     | Aliyuns      | DashScope 兼容接口 |
| OpenRouter     | OpenRouter   | 多模型聚合 |
| OpenAI 标准接口| Other        | 任意 OpenAI 兼容 API |

## 开发说明

- 修改前端代码后需执行 `npm run build` 或 `npm run build:dev`，再运行 `wails3 dev`
- 后端修改后 `wails3 dev` 会自动重载
- 数据库文件默认为项目根目录下的 `*.db`（如 `test.db`）

## 相关链接

- [Wails3 文档](https://v3.wails.io/)
- [CloudWeGo Eino](https://github.com/cloudwego/eino)
- [Wails Discord](https://discord.gg/JDdSxwjhGf)
- [Wails GitHub Discussions](https://github.com/wailsapp/wails/discussions)
