import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import i18n from '@/i18n';

// 字体大小偏移值定义（相对于默认值的偏移）
export const FONT_SIZE_OFFSETS = {
  VERY_SMALL: -4,   // 10px (默认14px - 4px)
  SMALL: -2,        // 12px (默认14px - 2px)
  NORMAL: 0,        // 14px (默认值)
  LARGE: 2,         // 16px (默认14px + 2px)
  VERY_LARGE: 4,    // 18px (默认14px + 4px)
  EXTRA_LARGE: 6,   // 20px (默认14px + 6px)
} as const;

export type FontSizeOffset = typeof FONT_SIZE_OFFSETS[keyof typeof FONT_SIZE_OFFSETS];

// 字体大小选项配置
export const FONT_SIZE_OPTIONS = [
  { value: FONT_SIZE_OFFSETS.VERY_SMALL, labelKey: 'settings.general.fontSizes.-4', description: '10px' },
  { value: FONT_SIZE_OFFSETS.SMALL, labelKey: 'settings.general.fontSizes.-2', description: '12px' },
  { value: FONT_SIZE_OFFSETS.NORMAL, labelKey: 'settings.general.fontSizes.0', description: '14px' },
  { value: FONT_SIZE_OFFSETS.LARGE, labelKey: 'settings.general.fontSizes.2', description: '16px' },
  { value: FONT_SIZE_OFFSETS.VERY_LARGE, labelKey: 'settings.general.fontSizes.4', description: '18px' },
  { value: FONT_SIZE_OFFSETS.EXTRA_LARGE, labelKey: 'settings.general.fontSizes.6', description: '20px' },
];

interface FontSizeState {
  // 当前字体大小偏移值
  fontSizeOffset: FontSizeOffset;
  
  // 设置字体大小偏移值
  setFontSizeOffset: (offset: FontSizeOffset) => void;
  
  // 获取当前实际字体大小（基础大小 + 偏移值）
  getCurrentFontSize: (baseSize?: number) => number;
  
  // 重置为默认大小
  resetFontSize: () => void;
}

export const useFontSizeStore = create<FontSizeState>()(
  persist(
    (set, get) => ({
      fontSizeOffset: FONT_SIZE_OFFSETS.NORMAL,
      
      setFontSizeOffset: (offset: FontSizeOffset) => {
        set({ fontSizeOffset: offset });
        // 更新CSS变量
        updateCSSFontSize(offset);
      },
      
      getCurrentFontSize: (baseSize = 14) => {
        const { fontSizeOffset } = get();
        return baseSize + fontSizeOffset;
      },
      
      resetFontSize: () => {
        set({ fontSizeOffset: FONT_SIZE_OFFSETS.NORMAL });
        updateCSSFontSize(FONT_SIZE_OFFSETS.NORMAL);
      },
    }),
    {
      name: 'lemon-tea-font-size',
      version: 1,
      // 恢复状态时同步到CSS变量
      onRehydrateStorage: () => (state) => {
        if (state) {
          updateCSSFontSize(state.fontSizeOffset);
        }
      },
    }
  )
);

// 监听 localStorage 变化，实现跨 tab 同步
if (typeof window !== 'undefined') {
  window.addEventListener('storage', (e) => {
    // 只处理字体大小相关的存储变化
    if (e.key === 'lemon-tea-font-size' && e.newValue) {
      try {
        const newState = JSON.parse(e.newValue);
        const currentState = useFontSizeStore.getState();
        
        // 如果字体大小偏移值发生变化，更新状态和CSS
        if (newState.state?.fontSizeOffset !== undefined && 
            newState.state.fontSizeOffset !== currentState.fontSizeOffset) {
          // 直接更新store状态（不触发set以避免循环）
          useFontSizeStore.setState({ fontSizeOffset: newState.state.fontSizeOffset });
          // 同步更新CSS变量
          updateCSSFontSize(newState.state.fontSizeOffset);
        }
      } catch (error) {
        console.error('Failed to sync font size across tabs:', error);
      }
    }
  });
}

// 更新CSS变量的函数
function updateCSSFontSize(offset: FontSizeOffset) {
  const root = document.documentElement;
  
  // 基础字体大小映射
  const baseSizes = {
    'font-size-xs': 10,     // --font-size-xs
    'font-size-sm': 12,     // --font-size-sm
    'font-size': 14,        // --font-size (标准)
    'font-size-lg': 16,     // --font-size-lg
    'font-size-xl': 18,     // --font-size-xl
    'font-size-xxl': 20,    // --font-size-xxl
  };
  
  // 更新所有字体大小相关的CSS变量
  Object.entries(baseSizes).forEach(([cssVar, baseSize]) => {
    const newSize = baseSize + offset;
    root.style.setProperty(`--${cssVar}`, `${newSize}px`);
  });
  
  // 同时更新行高以保持良好的可读性
  const baseLineHeight = 1.5715;
  // 字体越大，行高相对值可以稍微小一点
  const lineHeightAdjustment = offset > 0 ? -0.05 : offset < 0 ? 0.05 : 0;
  const newLineHeight = baseLineHeight + lineHeightAdjustment;
  
  root.style.setProperty('--line-height', `${newLineHeight}`);
  root.style.setProperty('--line-height-sm', `${Math.max(1.2, newLineHeight - 0.2)}`);
  root.style.setProperty('--line-height-lg', `${newLineHeight + 0.2}`);
}

// 初始化字体大小设置（在应用启动时调用）
export function initializeFontSize() {
  const store = useFontSizeStore.getState();
  updateCSSFontSize(store.fontSizeOffset);
}

// 获取字体大小选项标签
export function getFontSizeLabel(offset: FontSizeOffset): string {
  const option = FONT_SIZE_OPTIONS.find(opt => opt.value === offset);
  return option ? i18n.t(option.labelKey) : i18n.t('settings.general.fontSizes.0');
}

// 获取字体大小描述
export function getFontSizeDescription(offset: FontSizeOffset): string {
  const option = FONT_SIZE_OPTIONS.find(opt => opt.value === offset);
  return option ? option.description : '14px';
}
