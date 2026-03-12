# Lemon Tea Desktop 🍋🍵

English | [简体中文](./README_CN.md)

A cross-platform desktop AI chat client built based on Wails3, supporting multiple large language model providers, offers a smooth conversation experience.

## Features

- **Multi-turn Conversations**：Supports streaming output, stopping generation, and auto-generating conversation titles.
- **Conversation Management**：Create / delete / rename / favorite conversations.
- **Multiple LLM Providers**：
  - DeepSeek (DeepSeek AI)
  - Qwen (Alibaba Cloud Bailian / Tongyi Qwen)
  - OpenRouter
  - Ollama (local models)
  - Any OpenAI-compatible API
- **Tool Calling**：Built on the CloudWeGo Eino Agent architecture; supports built-in tools such as fetching the current date/time.
- **Cross-platform Support**：Compatible with macOS, Windows, and Linux—and experimental support for iOS and Android.

## Quick Start

### 1. Clone Project

```bash
git clone <repo>
cd lemon_tea_desktop
```

### 2. Install Deps

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 doctor
```

### 3. run dev

```bash
wails3 dev
```

## Plans
1. Multi-agent Support: A multi-agent system orchestrated by a main agent, which can proactively invoke other agents based on user needs
2. Persistent Memory: Supports global persistent memory with memory visualization and editing capabilities
3. Background Tasks: Supports user-assigned tasks that agents can execute automatically in the background
4. i18 support

## License

MIT License
