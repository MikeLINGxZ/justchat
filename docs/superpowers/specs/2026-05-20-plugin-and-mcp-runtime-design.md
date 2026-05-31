# Lemontea 插件与 MCP 运行时设计文档

## 目标

在现有 Lemontea 桌面端中新增“插件与工具”能力，并在本期完整打通 MCP 的管理、加载和执行链路。

本期重点不是做一个“静态配置页”，而是要让用户导入的 MCP 工具能够：

- 在设置页中被添加、查看、编辑、启用、禁用、重载、删除
- 在聊天输入框中作为可选工具出现
- 在主 Agent 运行时中被真正注册为可调用工具
- 在模型调用时通过现有 `tool_call / tool_result / confirm_request` 流程执行

插件机制本期只做管理位和结构预留，不实现插件执行运行时。MCP 是本期唯一需要真正可执行的外部扩展类型。

---

## 约束与前提

- 继续沿用当前 `wails3 + React + Zustand + Go service` 架构，不做无关重构。
- 所有界面必须支持浅色、深色主题。
- 所有新增文案必须支持简体中文、英文。
- 继续沿用现有设置页视觉结构和交互模式，整体样式与“模型供应商”页保持一致。
- Go 新增函数必须补充函数注释，复杂流程必须在关键位置补充说明。
- TypeScript 必须使用显式类型，不引入 `any`。
- 前端模块引用继续使用 `@` / `@bindings`。
- 本期只支持 `stdio` 型 MCP server。
- MCP 导入后必须复制到数据目录：
  - `{data_dir}/mcp/{mcp_name}/{version}`
  - 无版本号时使用 `{data_dir}/mcp/{mcp_name}/default`
- 当前仓库已经使用 `trpc-agent-go` 作为 Agent 运行时，本期必须基于它的工具能力接入 MCP，而不是额外再造一套独立执行链。

---

## 范围边界

### 本期必须实现

- 设置页新增“插件与工具”一级菜单
- 支持从本地目录导入 `mcp工具` 或 `插件`
- MCP 元信息解析、复制、落盘、列表展示
- MCP 配置文件查看、编辑、恢复、应用
- MCP server 启动、停止、重载
- MCP 工具发现，并注册到主 Agent 的工具集合
- 输入框工具选择区域展示“已启用 MCP 工具”
- 模型可以真正调用 MCP 工具，并复用现有工具确认与结果展示链路
- 关闭、禁用、删除、改配置后，运行时工具集合同步生效

### 本期不做

- 插件运行时执行
- 插件自定义视图渲染
- 非 `stdio` 的 `http` / `sse` MCP server
- 远程插件市场
- MCP OAuth、认证向导、权限页
- 多会话共享同一个 MCP 连接池的复杂优化

---

## 用户体验目标

### 设置页

设置页新增一级菜单 `插件与工具`，整体交互与“模型供应商”页一致：

- 左侧为列表区
- 右侧为详情区

左侧列表区：

- 顶部标题：`插件与工具`
- 右侧“添加”按钮点击后弹出菜单：
  - `mcp工具`
  - `插件`
- 选择后打开文件夹选择器
- 列表项展示：
  - 名称
  - 描述
  - 作者
  - 版本
  - 类型标签：`mcp工具` / `插件`
  - 启用开关
  - 更多菜单：`重新加载`、`删除`

右侧详情区：

- 未选中项时显示空态
- 选中 `mcp工具` 时显示：
  - 名称
  - 版本
  - 作者
  - 描述
  - 启用开关
  - 运行状态
  - 配置文件编辑器
  - `恢复` / `应用` 按钮
- 选中 `插件` 时显示只读信息和预留说明，不显示执行态

### 聊天输入框

输入框的工具下拉区分成两部分：

- 内建工具
- MCP 工具

只有已启用且运行成功、工具发现成功的 MCP 工具才会出现在可选工具列表中。

工具面板底部固定显示一个按钮：

- `管理工具与插件`

点击后直接打开设置页的 `插件与工具` 标签。

---

## 核心架构方案

本期采用“四层结构”：

1. `plugin service`
   - 面向前端提供管理接口
2. `mcp runtime manager`
   - 负责 MCP server 生命周期与工具发现
3. `agent tool registry`
   - 统一维护内建工具和动态 MCP 工具
4. `agent manager`
   - 基于当前启用工具集合构建 `trpc-agent-go` runner

这样可以保证：

- 设置页管理的是“扩展项”
- Agent 使用的是“工具实例”
- 两者之间通过运行时管理器和统一注册表桥接

---

## 数据模型设计

### 扩展项视图模型

新增统一视图模型，供前端设置页和输入框使用：

- `id`
- `name`
- `description`
- `author`
- `version`
- `kind`
  - `mcp`
  - `plugin`
- `enabled`
- `runtime_status`
  - `stopped`
  - `starting`
  - `running`
  - `error`
- `runtime_message`
- `root_dir`
- `config_file_path`
- `tool_count`
- `tools`

其中：

- `plugin` 可以没有 `tools`
- `mcp` 的 `tools` 需要包含已发现的工具元信息

### MCP 工具元信息

新增 MCP 工具视图结构：

- `tool_id`
- `server_id`
- `name`
- `description`
- `enabled`
- `requires_confirm`

`tool_id` 必须全局唯一，避免不同 server 暴露同名工具冲突。

本期命名规则固定为：

- `mcp:{server_name}:{tool_name}`

示例：

- `mcp:filesystem:read_file`

### 本地持久化结构

配置中新增插件与工具配置段，保存用户导入结果和启用状态。

每个扩展项至少保存：

- 扩展类型
- 名称
- 描述
- 作者
- 版本
- 原始导入目录
- 复制后的数据目录
- 是否启用
- MCP 配置文件相对路径
- 上次成功发现出的工具列表快照

工具列表快照不是执行真相，只是为了让设置页和输入框在下次启动前能先显示缓存数据；真正的执行工具集合仍以运行时发现结果为准。

---

## 目录与文件结构设计

### 后端

新增目录：

```text
backend/service/plugin/
├── plugin.go
├── plugin_implement.go
├── plugin_internal.go
└── plugin_dto/
    ├── list_extensions.go
    ├── import_extension.go
    ├── toggle_extension.go
    ├── reload_extension.go
    ├── delete_extension.go
    ├── get_extension_detail.go
    └── save_extension_config.go

backend/pkg/mcp/
├── manager.go
├── server.go
├── importer.go
├── manifest.go
├── config.go
└── manager_test.go
```

说明：

- `service/plugin` 负责 Wails 绑定方法和前端数据编排
- `pkg/mcp` 负责真正的 MCP 目录解析、进程管理、工具发现和实例构建

### 前端

新增目录：

```text
frontend/src/components/settings/plugins/
├── PluginToolList.tsx
├── PluginToolListItem.tsx
├── PluginToolDetailView.tsx
├── PluginToolConfigEditor.tsx
└── AddPluginToolMenu.tsx
```

同时扩展：

- `SettingsApp.tsx`
- `SettingsPrimaryMenu.tsx`
- `settingsStore.ts`
- `types/settings.ts`
- `ChatInput.tsx`

---

## MCP 导入与目录规范

### 导入流程

用户从设置页点击“添加 -> mcp工具”后：

1. 前端调用文件夹选择器
2. 后端接收所选目录
3. 后端解析 MCP 元信息与配置文件
4. 计算目标目录：
   - 有版本：`{data_dir}/mcp/{name}/{version}`
   - 无版本：`{data_dir}/mcp/{name}/default`
5. 复制目录到目标目录
6. 写入扩展项配置
7. 若默认启用，则立即尝试启动并发现工具
8. 返回更新后的扩展列表与详情

### 元信息来源

本期采用“目录内 manifest + 配置文件联合解析”的策略：

- 优先读取目录内清单文件中的名称、版本、作者、描述
- 若部分字段缺失，则从 MCP 配置中补齐
- 若仍缺失，则使用目录名和默认值兜底

这样设计的原因是：

- 不同 MCP 项目的目录结构可能不完全统一
- 需要允许导入“只有配置，没有完整 manifest”的本地 MCP 工具目录

### 配置文件编辑

右侧详情区展示的是 MCP 复制后目录里的配置文件内容，而不是原目录内容。

用户修改并点击“应用”后：

1. 前端提交新配置文本
2. 后端校验格式合法性
3. 写回复制后的配置文件
4. 若当前扩展已启用，则自动重载 MCP server
5. 更新运行状态与工具列表

点击“恢复”时：

- 恢复到当前已落盘的文件内容
- 不恢复到导入前原目录内容

这样可以确保“恢复”的语义和现有设置页一致，都是恢复到最后一次已应用状态。

---

## MCP 运行时设计

### 运行时管理器职责

`backend/pkg/mcp/Manager` 负责：

- 加载所有已配置扩展项
- 启动已启用的 MCP server
- 停止、重启指定 MCP server
- 维护 server 运行状态
- 拉取并缓存可用工具
- 将 MCP tool 转换为 `trpc-agent-go` 可注册的工具实例

### 生命周期

每个 MCP server 都维护独立状态：

1. `stopped`
2. `starting`
3. `running`
4. `error`

状态变化触发点：

- 启用时：`stopped -> starting -> running/error`
- 禁用时：`running/error -> stopped`
- 配置应用或手动重载时：`running/error -> starting -> running/error`
- 删除时：先 `stop`，再清理持久化信息

### 连接方式

本期仅支持 `stdio`：

- 根据配置文件中的 `command`、`args`、`env` 启动本地进程
- 通过 `trpc-agent-go` / `trpc-mcp-go` 推荐方式建立 MCP server 连接
- 读取工具列表并生成 tool 实例

### 错误处理

MCP server 启动或发现失败时：

- 扩展项仍保留在列表中
- 运行状态记为 `error`
- 保存错误信息到 `runtime_message`
- 不将其工具暴露给输入框和 Agent

这样用户可以直接在详情页里修配置后再应用，不需要重新导入。

---

## Agent 工具注册方案

### 当前问题

现有代码中：

- `tools.Registry` 只存 `ToolMeta`
- `Manager.buildAgentTools()` 用 `switch` 手动构造内建工具实例
- runner 缓存 key 没有包含工具集合

这会导致：

- 动态 MCP 工具无法挂入
- 工具配置变化后 runner 仍可能复用旧实例

### 新方案

将 `tools.Registry` 扩展为“工具定义注册表”：

每个工具定义至少包含：

- `meta`
- `factory`

其中：

- 内建工具：启动时静态注册
- MCP 工具：server 运行成功后动态注册

`factory` 负责在需要时创建 `trpc-agent-go` 工具实例。

### 工具唯一性

工具注册表按 `tool_id` 注册。

内建工具沿用原名：

- `datetime`
- `file_read`
- `file_write`
- `shell`
- `web_search`
- `code_exec`

MCP 工具使用全局唯一名：

- `mcp:{server_name}:{tool_name}`

前端展示时可以显示原始 `tool_name`，但向后端提交和在模型层注册时必须使用全局唯一 `tool_id`。

### 启用工具构建

`buildAgentTools(enabledUserTools []string)` 改为：

1. 从注册表读取所有内建工具
2. 按 `enabledUserTools` 读取被启用的用户工具
3. 通过各自 `factory` 生成 `toolpkg.Tool`

这样：

- 输入框选择逻辑不需要区分内建和 MCP
- 只要注册表里有动态工具，Agent 就能自然接入

---

## Runner 缓存策略

### 当前问题

当前 `GetOrCreateRunner()` 的 key 只包含：

- `baseURL`
- `modelName`

没有包含启用工具集合。工具变化后如果继续复用旧 runner，会出现：

- 输入框显示新工具，但模型实际上调用不到
- 已禁用工具仍然可能被旧 runner 调用

### 新方案

runner 缓存 key 扩展为：

- `baseURL`
- `modelName`
- `providerType`
- 启用工具集合快照

工具集合快照采用稳定排序后拼接生成。

此外，当发生以下事件时，需要清理相关缓存 runner：

- MCP 工具启用/禁用
- MCP server 重载成功
- MCP tool 名单变化
- 配置应用导致 tool schema 改变

这样可以保证：

- 输入框的工具选择
- 注册表里的真实工具实例
- runner 实际拥有的工具集合

三者始终一致。

---

## 前端设计

### 设置页一级菜单

`SettingsPrimaryTab` 新增：

- `plugins`

菜单顺序调整为：

- `general`
- `providers`
- `plugins`
- `about`

### 设置页状态

`settingsStore` 新增：

- 扩展项列表
- 当前选中扩展 id
- 配置编辑草稿
- 配置是否 dirty
- 配置应用中状态
- 运行状态刷新方法

详情页中的配置编辑草稿按“当前选中项”隔离，切换扩展时要重置到该项的已保存内容。

### 左侧列表交互

列表项点击后：

- 右侧详情随之刷新
- 不自动改变当前运行状态

启用开关变更后：

- 乐观更新 UI
- 请求后端启用/禁用
- 失败则回滚

“重新加载”点击后：

- 只作用于当前项
- 若项为 `plugin`，只刷新元信息
- 若项为 `mcp`，重启 server 并重新发现工具

### 详情区配置编辑器

本期采用简单文本编辑器，不引入复杂 IDE 能力。

需要支持：

- 查看完整文本
- 直接编辑
- dirty 检测
- 恢复
- 应用

按钮行为：

- 无修改时 `恢复` / `应用` 禁用
- 有修改时启用
- 应用成功后 dirty 清零

---

## 聊天输入框联动

### 工具列表来源

移除当前 `mockTools` 作为真实数据源的角色，改成：

- 内建工具列表来自后端静态定义
- MCP 工具列表来自后端运行时发现结果

前端工具下拉展示统一结构：

- `id`
- `name`
- `description`
- `enabled`
- `category`

### 展示分组

工具下拉分两组展示：

- 内建工具
- MCP 工具

底部固定按钮：

- `管理工具与插件`

点击后调用：

- `Window.OpenSettings({ tab: "plugins" })`

### 可选范围

只有满足以下条件的 MCP 工具才展示：

- 所属扩展项已启用
- 扩展运行状态为 `running`
- 工具发现成功

这可以避免前端勾选一个实际上不可执行的工具。

---

## 后端服务接口设计

### Plugin Service

新增公开方法：

- `ListExtensions`
- `ImportExtension`
- `GetExtensionDetail`
- `ToggleExtension`
- `ReloadExtension`
- `DeleteExtension`
- `SaveExtensionConfig`
- `ListAvailableTools`

统一签名遵循现有 Wails 绑定规范。

### 输入输出职责

- `ListExtensions`
  - 返回列表页需要的摘要字段
- `GetExtensionDetail`
  - 返回详情页完整信息与配置文本
- `ListAvailableTools`
  - 返回聊天输入框可选工具
- `ToggleExtension`
  - 负责启用/禁用并返回最新状态
- `ReloadExtension`
  - 负责重载 server 并返回最新状态与工具列表
- `SaveExtensionConfig`
  - 负责保存配置和自动重载

---

## 执行链路设计

### 初始化

应用启动后：

1. 加载配置中的扩展项
2. 对已启用的 MCP 项逐个启动
3. 成功后注册对应工具
4. 失败则记录错误状态

### 用户发送消息

用户在输入框勾选工具并发送消息时：

1. 前端提交 `enabledUserTools`
2. 后端 `Agent.Manager` 基于注册表构建对应工具集合
3. `trpc-agent-go` runner 运行
4. 模型选择并调用 MCP 工具
5. 事件通过当前已有流式消息链输出：
   - `tool_call`
   - `confirm_request`
   - `tool_result`
   - `assistant text`

### 配置变更

当用户在设置页启用、禁用、删除或重载 MCP 时：

1. 运行时管理器更新 server 状态
2. 工具注册表同步增删改
3. Agent runner 缓存失效
4. 输入框下次加载工具时拿到新结果

---

## 测试策略

### 后端

重点测试：

- MCP 目录导入路径是否正确
- 无版本号时是否落到 `default`
- 启用/禁用时运行状态是否正确切换
- 配置保存后是否触发重载
- 动态工具注册与移除是否正确
- runner 缓存 key 是否包含工具集合
- 工具变化后是否触发 runner 重建

优先为以下模块补测试：

- `backend/pkg/mcp/manager_test.go`
- `backend/pkg/agent/tools/registry_test.go`
- `backend/pkg/agent/manager_test.go`
- `backend/service/plugin/plugin_test.go`

### 前端

重点测试：

- 设置页新增一级菜单是否正常切换
- 插件与工具列表和详情联动
- 启用/禁用开关的乐观更新与回滚
- 配置编辑器 dirty 状态、恢复、应用
- 输入框工具列表是否正确显示 MCP 工具
- “管理工具与插件”按钮是否正确打开设置页

建议新增：

- `frontend/src/__tests__/pluginSettings.test.tsx`
- `frontend/src/__tests__/chatInputTools.test.tsx`

---

## 风险与取舍

### 1. MCP 项目目录结构不统一

风险：

- 不同 MCP 项目的配置和元信息文件格式可能不同

取舍：

- 本期先支持仓库约定的一种标准目录方案
- 解析器保留兜底逻辑，但不追求兼容所有社区目录结构

### 2. 动态工具名冲突

风险：

- 不同 server 可能暴露同名工具

取舍：

- 注册名强制加 `mcp:{server}:{tool}` 前缀
- 前端展示名和模型工具名允许分离

### 3. MCP server 启动慢或失败

风险：

- 启动时会影响可用工具列表

取舍：

- 列表允许展示错误状态
- 输入框只展示真正 `running` 的工具

### 4. runner 缓存与工具配置不一致

风险：

- 工具已经变化，但缓存 runner 未重建

取舍：

- 本期优先保证正确性，工具变化后主动失效缓存
- 不先做更复杂的增量复用优化

---

## 分阶段实现建议

### 阶段 1

- 建立 `plugin service` 和数据模型
- 完成 MCP 导入、列表、详情、配置读写
- 设置页 UI 打通

### 阶段 2

- 建立 `mcp runtime manager`
- 实现 `stdio` 启停、工具发现、错误状态
- 输入框改为真实工具来源

### 阶段 3

- 扩展 `tools.Registry`
- 动态注册 MCP 工具
- 更新 `Agent.Manager` 的 runner 构建与缓存失效

### 阶段 4

- 联调工具调用链
- 补测试
- 完整验证设置页、输入框、Agent 执行一致性

---

## 结论

本期最关键的设计决策有三点：

1. MCP 必须基于现有 `trpc-agent-go` 工具体系接入，而不是额外做一套执行桥。
2. `tools.Registry` 必须从“静态元信息表”升级为“可注册动态实例的工具定义表”。
3. runner 缓存必须感知工具集合变化，否则前端可见工具和实际可执行工具会长期不一致。

按照这个方案落地后，Lemontea 将具备一条完整的 MCP 工具链路：

- 设置页导入和管理
- 运行时启动和发现
- 输入框选择
- Agent 真正调用
- 前端展示结果

这也是后续继续扩展 `http/sse MCP`、插件执行运行时和更复杂工具权限体系的稳定基础。
