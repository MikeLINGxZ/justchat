# Lemon Tea Desktop 🍋🍵

[English](./README.md) | 简体中文

基于 Wails3 构建的跨平台桌面 AI 聊天客户端，支持多种大语言模型提供商，提供流畅的对话体验。

## 功能特性

- **多轮对话**：支持流式输出、停止生成、自动生成对话标题。
- **会话管理**：创建 / 删除 / 重命名 / 收藏会话。
- **多 LLM 提供商**：
  - DeepSeek（深度求索）
  - Qwen（阿里云百炼 / 通义千问）
  - OpenRouter
  - Ollama（本地模型）
  - 任何 OpenAI 兼容的 API
- **工具调用**：基于 CloudWeGo Eino Agent 架构；支持内置工具，如获取当前日期/时间。
- **跨平台支持**：兼容 macOS、Windows 和 Linux——并实验性支持 iOS 和 Android。

## 快速开始

### 1. 克隆项目

```bash
git clone <repo>
cd lemon_tea_desktop
```

### 2. 安装依赖

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 doctor
```

### 3. 运行开发模式

```bash
wails3 dev
```

## 计划

1. 多 Agent 支持：由主 Agent 调度的多 Agent 系统，主 Agent 可以根据用户需求主动调用其他 Agent
2. 持久化记忆：支持全局持久化记忆以及记忆可视化、记忆编辑
3. 后台任务：支持用户下达任务，Agent 自动后台执行
4. 国际化支持

## 许可证

MIT License
