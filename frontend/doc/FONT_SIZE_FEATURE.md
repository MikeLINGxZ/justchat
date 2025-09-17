# 字体大小调节功能

## 功能概述

字体大小调节功能允许用户根据个人喜好调整应用程序中所有文本的大小，提升阅读体验。该功能采用相对偏移值设计，确保字体大小调整的一致性和可维护性。

## 功能特性

### 📏 多级字体选择
- **极小** (10px): 最紧凑的文字显示
- **小** (12px): 较小的文字尺寸
- **标准** (14px): 默认推荐尺寸
- **大** (16px): 更舒适的阅读体验
- **很大** (18px): 大字体显示
- **极大** (20px): 最大的文字尺寸

### 🎛️ 多种调节方式
1. **滑块控制**: 直观的滑动条操作
2. **快速选择**: 预设选项一键切换
3. **实时预览**: 即时查看调整效果

### 💾 持久化存储
- 用户设置自动保存到本地存储
- 应用重启后自动恢复用户设置
- 跨会话保持设置状态

### 📱 响应式设计
- 移动端友好的操作界面
- 不同屏幕尺寸下的最佳显示效果

## 使用方法

### 访问字体设置
1. 打开应用设置 (通过菜单或快捷键)
2. 选择 "通用设置" 选项卡
3. 找到 "显示设置" 区域中的 "字体大小" 选项

### 调整字体大小
**方法一：使用滑块**
- 拖动滑块到所需位置
- 实时查看右侧显示的当前字体大小

**方法二：选择预设选项**
- 点击预设按钮 (极小/小/标准/大/很大/极大)
- 每个选项都显示对应的像素值

### 预览效果
- 在设置界面下方的预览区域查看文本效果
- 包含不同类型的文本样式展示
- 实时反映当前字体大小设置

### 重置设置
- 点击 "恢复默认" 按钮恢复到标准字体大小
- 点击 "应用设置" 确认当前设置

## 技术实现

### 架构设计
```typescript
// 字体大小状态管理
useFontSizeStore()
├── fontSizeOffset: 当前偏移值
├── setFontSizeOffset(): 设置偏移值
├── getCurrentFontSize(): 获取实际字体大小
└── resetFontSize(): 重置到默认值
```

### 偏移值机制
```typescript
// 基于默认14px的偏移值设计
FONT_SIZE_OFFSETS = {
  VERY_SMALL: -4,   // 10px = 14px - 4px
  SMALL: -2,        // 12px = 14px - 2px  
  NORMAL: 0,        // 14px = 14px + 0px
  LARGE: 2,         // 16px = 14px + 2px
  VERY_LARGE: 4,    // 18px = 14px + 4px
  EXTRA_LARGE: 6    // 20px = 14px + 6px
}
```

### CSS变量更新
系统自动更新以下CSS变量：
- `--font-size-xs`: 极小字体
- `--font-size-sm`: 小字体  
- `--font-size`: 标准字体
- `--font-size-lg`: 大字体
- `--font-size-xl`: 超大字体
- `--font-size-xxl`: 极大字体
- `--line-height`: 对应行高

### 组件结构
```
src/pages/settings/general/
├── index.tsx           # 通用设置页面组件
└── index.module.scss   # 样式文件

src/stores/
└── fontSizeStore.ts    # 字体大小状态管理
```

## 注意事项

### 兼容性
- 所有现有组件自动支持字体大小调节
- 通过CSS变量实现，无需修改各个组件
- 支持浅色和深色主题模式

### 性能优化
- CSS变量更新性能优异
- 状态持久化采用异步存储
- 避免不必要的重渲染

### 可扩展性
- 预留了更多字体大小选项的扩展空间
- 支持添加更多显示相关设置
- 可集成到主题系统中

## 开发指南

### 添加新的字体大小选项
```typescript
// 在 fontSizeStore.ts 中添加新选项
export const FONT_SIZE_OFFSETS = {
  // ... 现有选项
  HUGE: 8,  // 22px = 14px + 8px
};
```

### 在组件中使用字体设置
```typescript
import { useFontSizeStore } from '@/stores/fontSizeStore';

const MyComponent = () => {
  const { fontSizeOffset, getCurrentFontSize } = useFontSizeStore();
  
  // 获取当前字体大小
  const currentSize = getCurrentFontSize(); // 基于14px
  const headerSize = getCurrentFontSize(18); // 基于18px
  
  return <div>当前字体大小: {currentSize}px</div>;
};
```

### 添加新的CSS变量
```scss
// 在variables.scss中添加
:root {
  --font-size-custom: 16px; // 将被动态更新
}
```

然后在fontSizeStore.ts的updateCSSFontSize函数中添加对应更新逻辑。

## 更新日志

### v1.0.0 (当前版本)
- ✅ 基础字体大小调节功能
- ✅ 滑块和按钮双重控制方式  
- ✅ 实时预览效果
- ✅ 持久化存储
- ✅ 响应式设计
- ✅ 自动CSS变量更新
- ✅ 行高自适应调整

### 计划功能
- 🔄 更多字体相关设置 (字重、间距等)
- 🔄 主题集成
- 🔄 无障碍功能增强
- 🔄 快捷键支持