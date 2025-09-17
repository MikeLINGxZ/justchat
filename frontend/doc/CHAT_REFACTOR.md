# Chat 组件重构说明

## 优化内容

### 1. 重构 Chat 组件，支持独立使用和嵌入使用

Chat 组件现在支持两种使用模式：

#### 嵌入模式（Embedded Mode）
```tsx
<Chat
    standalone={false}
    chatTitle="对话标题"
    chatUuid="chat-uuid"
    currentMessages={messages}
    // ... 其他 props
/>
```

#### 独立模式（Standalone Mode）
```tsx
<Chat
    standalone={true}
    chatTitle="对话标题"
    chatUuid="chat-uuid"
    currentMessages={messages}
    // ... 其他 props
/>
```

### 2. 添加了加载动画

- 在聊天内容加载时显示精美的自定义加载动画
- **替换了原有的圆圈旋转动画**，改为更现代的波浪条形动画
- **解决了文本换行问题**，使用 `white-space: nowrap` 确保文本不换行
- 加载动画至少持续 500ms，确保良好的用户体验
- 支持亮色和暗色主题
- 响应式设计，移动端和桌面端都有适配

**新动画特性：**
- 5个渐变色彩条以波浪形式起伏
- 文本后跟随3个渐隐渐显的省略号
- 流畅的动画效果，视觉体验更佳
- 移动端适配：较小的条形尺寸和字体

### 3. 新增独立的 ChatPage 组件

创建了 `/src/pages/ChatPage.tsx`，可以作为独立页面使用：

```tsx
import ChatPage from '@/pages/ChatPage';

// 在路由中使用
<Route path="/chat/:chatUuid?" component={ChatPage} />
```

## 使用方式

### 在现有项目中嵌入使用

```tsx
import Chat from '@/pages/home/chat';

function MyComponent() {
    return (
        <div>
            <Chat
                standalone={false}
                chatTitle="我的对话"
                currentMessages={messages}
                onSendMessage={handleSendMessage}
                onModelChange={handleModelChange}
                // ... 其他必要的 props
            />
        </div>
    );
}
```

### 作为独立页面使用

```tsx
import ChatPage from '@/pages/ChatPage';

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/chat/:chatUuid?" element={<ChatPage />} />
            </Routes>
        </Router>
    );
}
```

## 主要改进

1. **灵活性增强**：Chat 组件现在可以既作为页面组件使用，也可以嵌入到其他组件中
2. **加载体验优化**：添加了骨架屏加载动画，提升用户体验
3. **类型安全**：所有 props 都支持可选，提高了组件的容错性
4. **样式优化**：支持独立模式和嵌入模式的不同样式表现

## 注意事项

- 独立模式下会显示标题栏，嵌入模式下标题栏的显示取决于是否传入 `onTitleChange` 回调
- 加载动画使用了 Ant Design 的 Spin 组件，确保了一致的视觉体验
- 所有的事件回调都是可选的，未传入时组件会优雅降级