# 语音输入功能设计

## 概述

为 Lemon Tea Desktop 的聊天输入框添加语音输入能力，用户可以通过麦克风实时录音，将语音转换为文本输入。录音过程中显示声纹动画，并实时显示识别出的文字。

## 需求

### 功能需求

1. **语音按钮**：在输入框底部栏的发送按钮左侧添加麦克风按钮
2. **语音输入模式**：点击麦克风按钮后，输入框区域切换为语音输入模式
3. **声纹动画**：语音输入模式下，输入框区域显示实时声纹动画，声纹随用户说话音量变化
4. **实时语音识别**：录音的同时进行语音转文字，识别结果实时显示在声纹下方
5. **退出语音模式**：用户可以点击停止按钮结束录音，识别出的文本填入输入框

### 非功能需求

- 支持 macOS 和 Windows 平台
- 麦克风权限请求与错误处理
- 语音识别延迟 < 500ms（取决于网络和 STT 服务）
- 声纹动画流畅，帧率 ≥ 30fps

## 1. 交互流程

### 1.1 状态流转

```
正常输入模式 ──点击麦克风──► 语音输入模式 ──点击停止──► 正常输入模式
                              │                         ▲
                              │  发生错误/权限拒绝        │
                              └────────────────────────►│
```

### 1.2 正常输入模式

底部栏布局（新增麦克风按钮）：

```
┌──────────────────────────────────────────────────┐
│  [富文本编辑器]                                    │
├──────────────────────────────────────────────────┤
│  [+] [工具]              [模型选择] [🎤] [发送]    │
└──────────────────────────────────────────────────┘
```

- 麦克风按钮（🎤）位于模型选择器和发送按钮之间
- 按钮尺寸与发送按钮一致（32x32px）
- 默认状态：灰色图标，hover 时变为主题色

### 1.3 语音输入模式

点击麦克风按钮后，输入框区域切换为：

```
┌──────────────────────────────────────────────────┐
│                                                  │
│            ▁ ▃ ▅ ▇ █ ▇ ▅ ▃ ▁  （声纹动画）        │
│                                                  │
│     "你好，请帮我写一个函数"  （实时识别文本）       │
│                                                  │
├──────────────────────────────────────────────────┤
│  [+] [工具]              [模型选择]  [⏹ 停止]     │
└──────────────────────────────────────────────────┘
```

- 富文本编辑器区域被声纹动画 + 识别文本替代
- 底部栏的麦克风按钮和发送按钮合并为一个「停止」按钮（红色背景）
- 声纹动画居中显示，宽度约为输入框的 60%
- 识别文本在声纹下方，居中，单行滚动或多行显示

### 1.4 退出语音模式

点击停止按钮后：
- 停止录音和语音识别
- 将已识别的完整文本填入富文本编辑器（追加到现有内容后）
- 恢复正常输入模式
- 用户可继续编辑文本后发送

## 2. 技术方案

### 2.1 整体架构

```
┌─────────────────────────────────────────────┐
│  前端 (React)                                │
│  ┌─────────────┐    ┌───────────────────┐   │
│  │ VoiceInput   │    │ WaveformCanvas    │   │
│  │ Component    │    │ (Canvas 2D)       │   │
│  │              │    └───────────────────┘   │
│  │  ┌────────┐  │    ┌───────────────────┐   │
│  │  │Recorder│  │    │ TranscriptDisplay │   │
│  │  └────┬───┘  │    └───────────────────┘   │
│  └───────┼──────┘                            │
│          │ audio chunks                      │
│          ▼                                   │
│  ┌──────────────┐                            │
│  │ Wails IPC    │                            │
│  └──────┬───────┘                            │
├─────────┼───────────────────────────────────┤
│         ▼                                    │
│  Go 后端                                     │
│  ┌──────────────┐    ┌───────────────────┐   │
│  │ VoiceService │───►│ STT Provider      │   │
│  │              │    │ (语音识别服务)      │   │
│  └──────────────┘    └───────────────────┘   │
└─────────────────────────────────────────────┘
```

### 2.2 前端实现

#### 录音模块

使用浏览器原生 `MediaRecorder` API（Wails webview 支持）：

```typescript
// 核心录音逻辑
const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
const mediaRecorder = new MediaRecorder(stream, {
  mimeType: 'audio/webm;codecs=opus'
})

// 获取音频分析数据用于声纹动画
const audioContext = new AudioContext()
const analyser = audioContext.createAnalyser()
const source = audioContext.createMediaStreamSource(stream)
source.connect(analyser)
analyser.fftSize = 256
```

#### 声纹动画

使用 Canvas 2D 绘制实时声纹：

```typescript
// 从 AnalyserNode 获取频率数据
const dataArray = new Uint8Array(analyser.frequencyBinCount)
analyser.getByteFrequencyData(dataArray)

// 绘制竖条形声纹（类似音频均衡器）
// 选取中间频段的若干个频率点，映射为竖条高度
// 使用 requestAnimationFrame 持续更新
```

声纹样式：
- 竖条数量：约 20-30 根
- 竖条宽度：3-4px，间距 2px，圆角顶部
- 颜色：主题色（#1890ff）渐变
- 无声时：所有竖条保持最小高度（2-3px），轻微呼吸动画
- 说话时：竖条高度随频率数据实时变化，加入平滑过渡

#### 组件结构

```
ChatInput
├── ... (现有组件)
└── VoiceInputOverlay (语音模式覆盖层，条件渲染)
    ├── WaveformCanvas (声纹动画 Canvas)
    └── TranscriptText (识别文本显示)
```

新增文件：
```
frontend/src/components/chat/input/
├── voice_input.tsx              # 语音输入覆盖层组件
├── voice_input.module.scss      # 样式
├── waveform_canvas.tsx          # 声纹动画 Canvas 组件
└── use_voice_recorder.ts        # 录音 + 音频分析 Hook
```

### 2.3 后端实现

#### STT（语音转文字）服务

Go 后端负责将音频流发送给语音识别服务，返回识别文本。

```go
// backend/service/voice.go
type VoiceService struct {
    sttProvider STTProvider
}

// STT 提供者接口
type STTProvider interface {
    // 流式识别：接收音频数据，返回识别结果 channel
    StreamRecognize(ctx context.Context, audioStream <-chan []byte) (<-chan TranscriptResult, error)
}

type TranscriptResult struct {
    Text    string  // 识别出的文本
    IsFinal bool    // 是否为最终结果（非中间结果）
}
```

#### ASR 提供者选项

采用可插拔的 Provider 模式，支持多种外部 ASR 服务和本地方案：

| 提供者 | 类型 | 流式识别 | 优势 | 劣势 |
|--------|------|---------|------|------|
| 阿里云 ASR | 云端 | WebSocket 实时流 | 中文识别质量高，支持方言，价格低 | 需注册阿里云账号 |
| 腾讯云 ASR | 云端 | WebSocket 实时流 | 中文识别优秀，支持热词 | 需注册腾讯云账号 |
| 讯飞 ASR | 云端 | WebSocket 实时流 | 中文 ASR 行业标杆，支持多方言 | 需注册讯飞账号 |
| OpenAI Whisper API | 云端 | 非流式（分段） | 多语言优秀，可复用已有 API Key | 非实时，需分段发送 |
| Google Speech-to-Text | 云端 | gRPC/REST 流 | 多语言支持广 | 国内网络受限 |
| 浏览器 Web Speech API | 本地+云 | 实时流 | 零配置，无需后端 | 平台差异大，中文需联网 |

**推荐初期支持**：
1. 优先实现 **通用 OpenAI 兼容接口**（覆盖 OpenAI Whisper、Groq、DeepSeek 等兼容服务）
2. 实现 **WebSocket 流式 ASR 通用接口**（覆盖阿里云、腾讯云、讯飞等国内服务）
3. **浏览器 Web Speech API** 作为零配置的 fallback

#### ASR Provider 接口设计

```go
// backend/pkg/asr/provider.go

// ASR 提供者接口
type Provider interface {
    // Name 返回提供者名称
    Name() string
    // StreamRecognize 流式识别：接收音频数据，返回识别结果 channel
    StreamRecognize(ctx context.Context, config RecognizeConfig, audioStream <-chan []byte) (<-chan Result, error)
    // Recognize 非流式识别：一次性识别完整音频
    Recognize(ctx context.Context, config RecognizeConfig, audioData []byte) (*Result, error)
}

type RecognizeConfig struct {
    Language    string  // 语言代码，如 "zh-CN", "en-US"
    SampleRate  int     // 采样率，如 16000
    Format      string  // 音频格式，如 "pcm", "wav", "webm"
    EnableITN   bool    // 是否启用逆文本正则化（数字、日期格式化）
}

type Result struct {
    Text    string  // 识别出的文本
    IsFinal bool    // 是否为最终结果（非中间结果）
}
```

#### ASR Provider 配置

用户在设置页配置 ASR 服务，存储在数据库中：

```go
// backend/storage/data_models/asr.go
type ASRProviderConfig struct {
    ID       string `gorm:"primarykey"`
    Name     string // 显示名称
    Type     string // "openai_compatible" | "websocket_stream" | "web_speech"
    BaseURL  string // API 地址
    APIKey   string // API 密钥（加密存储）
    Model    string // 模型名称（如 "whisper-1"）
    Language string // 默认语言
    Enabled  bool   // 是否启用
    IsDefault bool  // 是否为默认提供者
    Extra    string // JSON 格式的额外参数（各服务特有配置）
}
```

#### 已实现的 Provider

**1. OpenAI 兼容 Provider**（非流式，分段发送）

```go
// backend/pkg/asr/openai_compatible.go
type OpenAICompatibleProvider struct {
    baseURL string
    apiKey  string
    model   string
}

// 将音频分段（每段约 3-5 秒），逐段调用 /v1/audio/transcriptions
// 每段识别完成后立即返回结果，实现伪实时效果
func (p *OpenAICompatibleProvider) StreamRecognize(...) (<-chan Result, error) {
    // 1. 从 audioStream 累积音频数据
    // 2. 每积累约 3 秒音频，发送到 API
    // 3. 将识别结果写入 result channel
}
```

**2. WebSocket 流式 Provider**（实时流式）

```go
// backend/pkg/asr/websocket_stream.go
type WebSocketStreamProvider struct {
    wsURL       string
    apiKey      string
    authBuilder func() (url string, headers map[string]string) // 各服务鉴权方式不同
}

// 通过 WebSocket 实时发送音频帧，实时接收识别结果
func (p *WebSocketStreamProvider) StreamRecognize(...) (<-chan Result, error) {
    // 1. 建立 WebSocket 连接（含鉴权）
    // 2. 启动 goroutine 持续发送音频数据
    // 3. 启动 goroutine 持续接收识别结果
    // 4. 解析各服务的响应格式，统一为 Result
}
```

**3. Web Speech API Provider**（前端本地）

不经过后端，前端直接使用浏览器 API：

```typescript
// 前端 use_voice_recorder.ts
const recognition = new (window.SpeechRecognition || window.webkitSpeechRecognition)()
recognition.continuous = true
recognition.interimResults = true
recognition.lang = 'zh-CN'

recognition.onresult = (event) => {
  let interim = ''
  let final = ''
  for (let i = event.resultIndex; i < event.results.length; i++) {
    if (event.results[i].isFinal) {
      final += event.results[i][0].transcript
    } else {
      interim += event.results[i][0].transcript
    }
  }
  setTranscript(prev => prev + final)
  setInterimText(interim)
}
```

### 2.4 数据流

**外部 ASR 模式（通过后端）：**

```
用户说话
  → 麦克风采集音频流 (MediaStream)
  → 分流：
      ├── AnalyserNode → 频率数据 → Canvas 声纹动画 (每帧更新)
      └── MediaRecorder → 音频数据块
          → Wails IPC → Go VoiceService → ASR Provider
          → 识别结果 → Wails Event → 前端状态更新
              ├── interim (中间结果) → 灰色文字显示
              └── final (最终结果) → 累加到已确认文本
```

**Web Speech API 模式（纯前端）：**

```
用户说话
  → 麦克风采集音频流 (MediaStream)
  → 分流：
      ├── AnalyserNode → 频率数据 → Canvas 声纹动画 (每帧更新)
      └── SpeechRecognition → 识别结果回调
          ├── interim (中间结果) → 灰色文字显示
          └── final (最终结果) → 累加到已确认文本
```

**前端与后端 ASR 的统一抽象：**

```typescript
// 前端统一 ASR 接口
interface ASRAdapter {
  start(lang: string): void
  stop(): void
  onResult: (callback: (text: string, isFinal: boolean) => void) => void
  onError: (callback: (error: string) => void) => void
}

// Web Speech API 实现
class WebSpeechASR implements ASRAdapter { ... }

// 后端 ASR 实现（通过 Wails IPC）
class BackendASR implements ASRAdapter {
  start(lang: string) {
    Service.StartASR(lang)  // 调用 Go 后端
  }
  // 通过 Wails Event 接收识别结果
}
```

### 2.5 后端 ASR 服务集成

```go
// backend/service/voice.go

func (s *Service) StartASR(lang string) error {
    // 1. 获取用户配置的默认 ASR Provider
    // 2. 初始化 Provider
    // 3. 开始接收前端发来的音频数据
    // 4. 将识别结果通过 Wails Event 推送到前端
}

func (s *Service) StopASR() error {
    // 停止识别，释放资源
}

func (s *Service) SendAudioChunk(data []byte) error {
    // 前端发送音频数据块到后端
}

// ASR 设置相关
func (s *Service) GetASRProviders() []ASRProviderViewModel
func (s *Service) CreateASRProvider(config ASRProviderConfig) error
func (s *Service) UpdateASRProvider(id string, config ASRProviderConfig) error
func (s *Service) DeleteASRProvider(id string) error
func (s *Service) SetDefaultASRProvider(id string) error
```

用户点击停止
  → 停止 MediaRecorder + ASR
  → 释放 MediaStream
  → 关闭 AudioContext
  → 将完整识别文本填入编辑器
  → 切回正常输入模式

## 3. 前端组件详细设计

### 3.1 VoiceInput 组件

```typescript
interface VoiceInputProps {
  onTranscriptComplete: (text: string) => void  // 录音结束，返回完整文本
  onCancel: () => void                           // 取消录音
}

// 状态
interface VoiceInputState {
  isRecording: boolean          // 是否正在录音
  transcript: string            // 已确认的识别文本
  interimText: string           // 中间识别结果（未确认）
  error: string | null          // 错误信息
}
```

### 3.2 useVoiceRecorder Hook

```typescript
interface UseVoiceRecorderReturn {
  isRecording: boolean
  transcript: string            // 累计已确认文本
  interimText: string           // 当前中间结果
  analyserNode: AnalyserNode | null  // 用于声纹动画
  error: string | null
  startRecording: () => Promise<void>
  stopRecording: () => void
}

function useVoiceRecorder(lang: string): UseVoiceRecorderReturn
```

### 3.3 WaveformCanvas 组件

```typescript
interface WaveformCanvasProps {
  analyser: AnalyserNode | null  // 音频分析节点
  isActive: boolean              // 是否正在录音
  width?: number
  height?: number
}
```

Canvas 绘制逻辑：
- 使用 `requestAnimationFrame` 循环
- 每帧从 `AnalyserNode.getByteFrequencyData()` 获取频率数据
- 选取均匀分布的频率点，映射为竖条高度
- 应用平滑算法（指数移动平均）避免闪烁
- 非录音状态时显示静态最小高度竖条

### 3.4 样式设计

```scss
// voice_input.module.scss

.voiceOverlay {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  gap: 16px;
  min-height: 120px;
}

.waveformContainer {
  width: 60%;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.transcriptText {
  font-size: 14px;
  color: var(--text-secondary);
  text-align: center;
  max-height: 3em;          // 最多显示约 2 行
  overflow-y: auto;
  width: 100%;
  padding: 0 16px;
  
  .interim {
    color: var(--text-tertiary);  // 中间结果用更浅的颜色
  }
}

.stopButton {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #dc3545;
  // 复用现有停止按钮样式
}
```

## 4. ChatInput 组件修改

### 4.1 新增状态

```typescript
const [isVoiceMode, setIsVoiceMode] = useState(false)
```

### 4.2 底部栏修改

在 `rightActions` 中，模型选择器和发送按钮之间添加麦克风按钮：

```tsx
<div className={styles.rightActions}>
  {/* 模型选择器 - 保持不变 */}
  <ModelSelector ... />
  
  {/* 语音按钮 - 新增 */}
  {!isVoiceMode && !isGenerating && (
    <button
      className={styles.voiceButton}
      onClick={() => setIsVoiceMode(true)}
      title={t('chat.input.voiceInput')}
    >
      <AudioOutlined />
    </button>
  )}
  
  {/* 发送/停止按钮 - 保持不变 */}
  {isGenerating ? <StopButton /> : <SendButton />}
</div>
```

### 4.3 输入区域条件渲染

```tsx
<div className={styles.inputContainer}>
  {isVoiceMode ? (
    <VoiceInputOverlay
      onTranscriptComplete={(text) => {
        // 将识别文本追加到编辑器
        appendTextToEditor(text)
        setIsVoiceMode(false)
      }}
      onCancel={() => setIsVoiceMode(false)}
    />
  ) : (
    <>
      <FilesContainer ... />
      <RichMarkdownEditor ... />
    </>
  )}
  <BottomBar ... />
</div>
```

## 5. ASR 设置页

在设置菜单中新增「语音识别」页，用户可以配置 ASR 服务：

```
┌──────────────────────────────────────────────────┐
│  设置                                             │
│  ├── 通用                                         │
│  ├── 模型                                         │
│  ├── Agent                                        │
│  ├── 技能                                         │
│  ├── 语音识别  ◄── 新增                            │
│  └── 关于                                         │
├──────────────────────────────────────────────────┤
│                                                  │
│  ASR 服务                          [+ 添加服务]    │
│                                                  │
│  ┌────────────────────────────────────────────┐  │
│  │ ⭐ OpenAI Whisper          默认             │  │
│  │ https://api.openai.com                     │  │
│  │ 模型: whisper-1                            │  │
│  │                       [设为默认] [编辑] [删除] │ │
│  └────────────────────────────────────────────┘  │
│                                                  │
│  ┌────────────────────────────────────────────┐  │
│  │   阿里云 ASR                                │  │
│  │ wss://nls-gateway.aliyuncs.com              │  │
│  │                       [设为默认] [编辑] [删除] │ │
│  └────────────────────────────────────────────┘  │
│                                                  │
│  ┌────────────────────────────────────────────┐  │
│  │   浏览器内置 (Web Speech API)               │  │
│  │ 零配置，无需 API Key                        │  │
│  │                       [设为默认]             │  │
│  └────────────────────────────────────────────┘  │
│                                                  │
│  默认语言：[中文 (zh-CN)  ▾]                      │
│                                                  │
└──────────────────────────────────────────────────┘
```

添加/编辑 ASR 服务时弹出表单：

| 字段 | 说明 |
|------|------|
| 名称 | 服务显示名称 |
| 类型 | OpenAI 兼容 / WebSocket 流式 / 浏览器内置 |
| API 地址 | 服务端点 URL |
| API Key | 密钥（加密存储） |
| 模型 | 模型名称（OpenAI 兼容类型需要） |
| 默认语言 | 默认识别语言 |

## 6. Go 端模块结构

```
backend/
├── pkg/
│   └── asr/
│       ├── provider.go              # Provider 接口定义
│       ├── openai_compatible.go     # OpenAI 兼容实现
│       └── websocket_stream.go      # WebSocket 流式实现
├── service/
│   └── voice.go                     # 语音服务（ASR 管理 + 音频流处理）
└── storage/
    └── data_models/
        └── asr.go                   # ASR 配置数据模型
```

## 7. 错误处理

| 场景 | 处理方式 |
|------|---------|
| 麦克风权限被拒绝 | 显示提示："请在系统设置中允许麦克风权限"，退出语音模式 |
| 未配置任何 ASR 服务 | 自动使用浏览器内置 Web Speech API；若不支持则提示用户去设置页配置 ASR 服务 |
| ASR 服务连接失败 | 显示错误提示，仍然保持录音和声纹动画（录音可用但无文字） |
| ASR 服务认证失败 | 提示 "API Key 无效或已过期，请检查 ASR 设置" |
| 识别过程中网络断开 | 显示网络错误提示，已识别文本保留 |
| 长时间无语音输入 | 根据 Provider 类型处理：Web Speech 自动重启，外部 ASR 发送静音帧保持连接 |
| ASR 服务返回错误 | 显示具体错误信息，允许用户切换到其他 Provider 重试 |

## 8. 国际化

新增 i18n key：

```json
{
  "chat.input.voiceInput": "语音输入",
  "chat.input.voiceRecording": "正在录音...",
  "chat.input.voiceStop": "停止录音",
  "chat.input.voiceError.noPermission": "无法访问麦克风，请检查系统权限设置",
  "chat.input.voiceError.notSupported": "当前环境不支持语音识别，请在设置中配置 ASR 服务",
  "chat.input.voiceError.networkError": "语音识别服务连接失败",
  "chat.input.voiceError.authError": "ASR 服务认证失败，请检查 API Key",
  "settings.asr.title": "语音识别",
  "settings.asr.addProvider": "添加服务",
  "settings.asr.setDefault": "设为默认",
  "settings.asr.defaultLanguage": "默认语言",
  "settings.asr.type.openaiCompatible": "OpenAI 兼容",
  "settings.asr.type.websocketStream": "WebSocket 流式",
  "settings.asr.type.webSpeech": "浏览器内置"
}
```

## 9. 后续扩展

以下能力不在初期范围内，可后续迭代：

- **语音唤醒**：无需点击按钮，通过唤醒词触发语音输入
- **语音输出（TTS）**：Agent 回复可以朗读出来
- **多语言自动检测**：自动识别用户正在说的语言
- **语音消息**：直接发送语音消息而非转文字
- **本地 Whisper 模型**：使用 whisper.cpp 实现完全离线识别
