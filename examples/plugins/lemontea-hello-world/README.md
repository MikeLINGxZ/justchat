# Hello World 示例插件

Lemon Tea Desktop 插件系统的示例插件，演示所有核心 API。

## 功能

### 工具
- **Random Joke** (`random-joke`) - 随机讲一个编程笑话，支持中英文
- **Dice Roll** (`dice-roll`) - 掷骰子，支持自定义面数和次数

### Agent
- **Greeting Agent** - 一个幽默的问候助手，可以讲笑话和玩骰子

### Hook
- **onBeforeChat** - 统计用户对话次数
- **onAfterChat** - 记录最后对话时间

### 存储
- `joke_count` - 笑话工具使用次数
- `chat_count` - 用户对话总次数
- `last_chat_time` - 最后一次对话时间

## 构建

```bash
cnpm install
cnpm run build
```

## 安装

1. 构建插件后，打开 Lemon Tea Desktop
2. 进入 设置 → 插件
3. 点击「安装插件」，选择此文件夹
4. 插件会自动激活

## 开发

```bash
cnpm run dev   # 监听模式，自动编译
```
