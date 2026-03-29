# Lemon Tea Desktop

English | [简体中文](./README_CN.md)

Lemon Tea Desktop is a cross-platform AI desktop client built with Wails v3, Go, React, and TypeScript. It focuses on chat, tool calling, workflow-style task execution, and local desktop integration.

<p align="center"><img src="docs/app_home.png" alt="Lemon Tea Desktop Home" width="80%" height="auto" /></p>

## Current Features

- Multi-turn chat with streaming output and manual stop.
- Conversation management: create, rename, delete, and favorite chats.
- Automatic chat title generation.
- File attachments in chat input, including local image preview.
- Multiple model providers:
  - DeepSeek
  - Alibaba Cloud Bailian / Qwen compatible endpoint
  - OpenRouter
  - Ollama
  - Any OpenAI-compatible API
- Model management:
  - Load models from provider endpoints
  - Set provider default model
  - Add and remove custom models for a provider
  - Remember a local default chat model
- Tool calling based on CloudWeGo Eino / ADK.
- Built-in tools currently include:
  - current date
  - current time
  - block/wait tool for timing-style tasks
- Custom MCP tool integration:
  - import MCP servers from a local folder
  - enable/disable MCP tools
  - remove imported MCP tools
- Workflow-oriented execution for more complex requests:
  - task planning
  - worker execution
  - synthesis
  - review/retry
- Execution trace UI for plan steps, tool calls, stage transitions, and elapsed time.
- Task recovery for interrupted running tasks after restart.
- Prompt management UI:
  - view built-in prompt files
  - edit and save prompt files
  - reset prompt files to defaults
- General settings UI for application font size.
- Responsive chat/settings layout for desktop and mobile-sized windows.
- Cross-platform targets:
  - macOS
  - Windows
  - Linux
  - experimental iOS / Android build scaffolding

## Tech Stack

- Backend: Go
- Desktop shell: Wails v3
- Frontend: React 19 + TypeScript + Vite
- UI: Ant Design
- Agent / tool orchestration: CloudWeGo Eino
- Local storage: SQLite via GORM

## Quick Start

### 1. Clone

```bash
git clone <repo>
cd lemon_tea_desktop
```

### 2. Install dependencies

Install Wails v3:

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 doctor
```

Install frontend dependencies:

```bash
cd frontend
npm install
cd ..
```

### 3. Run development mode

Using Wails directly:

```bash
wails3 dev -config ./build/config.yml
```

Or with Task:

```bash
task dev
```

## Project Structure

```text
.
├── backend/     Go services, storage, provider adapters, agent workflow logic
├── frontend/    React UI, chat pages, settings pages, components
├── build/       Wails build and packaging configuration
├── docs/        README assets
└── main.go      Desktop app entry
```

## Roadmap

- richer multi-agent collaboration
- persistent memory visualization and editing
- background task automation
- internationalization

## License

MIT
