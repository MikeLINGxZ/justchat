# AI Chat - Eino 框架示例

基于 [CloudWeGo Eino](https://github.com/cloudwego/eino) 框架实现的简单 AI 聊天应用。

## 功能特性

- 支持 **Qwen**（阿里云 DashScope）、**OpenAI**（GPT 系列）、**Ollama**（本地开源模型）
- **Agent-to-Agent 架构**：Chat 本身是主 Agent，通过 AgentTool 将子任务委托给专家 Agent
  - **ChatAgent**（主）：接收用户输入，由 LLM 决定直接回答或调用子 Agent
  - **DateTimeAgent**：日期时间工具（get_current_date、get_current_time）
  - **FruitPriceAgent**：水果价格工具（get_fruit_price，内置模拟数据）
- 流式输出，边生成边显示

## 环境要求

- Go 1.18+
- 任选其一：
  - **Qwen**：阿里云 DashScope，需配置 `ALIYUN_API_KEY`
  - **OpenAI**：需配置 `OPENAI_API_KEY`
  - **Ollama**：本地运行，无需 API Key，[下载安装](https://ollama.com)

## 快速开始

### 使用 Qwen（阿里云）

```bash
export ALIYUN_API_KEY="your-dashscope-api-key"
export QWEN_MODEL="qwen-turbo"  # 可选，默认 qwen-turbo，也可用 qwen-plus 等
go run main.go
```

### 使用 Ollama（推荐本地入门）

1. 安装并启动 Ollama：

   ```bash
   # macOS
   brew install ollama
   ollama serve
   
   # 拉取模型
   ollama pull llama3.2
   ```

2. 运行聊天程序：

   ```bash
   go run main.go
   # 或
   go build -o aichat . && ./aichat
   ```

3. 可选环境变量：
   - `MODEL_NAME`：Ollama 模型名称（默认 `llama3.2`）
   - `OLLAMA_BASE_URL`：Ollama 服务地址（默认 `http://localhost:11434`）

### 环境变量说明

| 提供商 | 环境变量 | 说明 |
|--------|----------|------|
| Qwen | `ALIYUN_API_KEY` | 阿里云 DashScope API Key（必填） |
| Qwen | `QWEN_MODEL` | 模型名称（默认 `qwen-turbo`） |
| Qwen | `QWEN_BASE_URL` | API 地址（默认 DashScope 兼容模式） |
| OpenAI | `OPENAI_API_KEY` | OpenAI API Key |
| OpenAI | `OPENAI_MODEL` | 模型名称（默认 `gpt-4o-mini`） |
| Ollama | `MODEL_NAME` | 模型名称（默认 `llama3.2`） |
| Ollama | `OLLAMA_BASE_URL` | 服务地址（默认 `http://localhost:11434`） |

### 使用 OpenAI

```bash
export OPENAI_API_KEY="your-api-key"
export OPENAI_MODEL="gpt-4o-mini"  # 可选，默认 gpt-4o-mini
go run main.go
```

## 交互命令

| 命令 | 说明 |
|------|------|
| `/exit` | 退出 |
| `/clear` | 清空（提示） |
| `/help` | 显示帮助 |

### Agent-to-Agent 架构

```
用户输入 → ChatAgent（主）
              ├─ 直接回答（通用聊天）
              ├─ 调用 DateTimeAgent（日期时间问题）
              └─ 调用 FruitPriceAgent（水果价格问题）
```

由 **ChatAgent 的 LLM** 根据用户问题决定是否调用子 Agent，而非代码中的关键词路由。

- **DateTimeAgent**：工具 `get_current_date`、`get_current_time`
- **FruitPriceAgent**：工具 `get_fruit_price`（内置 20+ 种水果模拟价格）

## 项目结构

```
eino_demo/
├── main.go       # 主程序（ChatAgent Runner、聊天循环）
├── chat_agent.go # ChatAgent 主 Agent（AgentTool 编排子 Agent）
├── agent.go      # DateTimeAgent 子 Agent
├── fruit_agent.go# FruitPriceAgent 子 Agent
├── go.mod
└── README.md
```

## 依赖

- [github.com/cloudwego/eino](https://github.com/cloudwego/eino) - LLM 应用开发框架
- [github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) - Qwen、Ollama、OpenAI 等模型实现
