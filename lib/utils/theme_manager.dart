import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:lemon_tea/generated/l10n.dart';

/// 主题模式枚举
enum ThemeMode {
  light,    // 浅色模式
  dark,     // 深色模式
  system,   // 跟随系统
}

/// 字体大小设置
enum FontSizeMode {
  extraSmall,
  small,
  medium,
  large,
  extraLarge,
}

/// 获取字体大小模式的名称
String getFontSizeModeName(FontSizeMode mode) {
  switch (mode) {
    case FontSizeMode.extraSmall:
      return '超小';
    case FontSizeMode.small:
      return '小';
    case FontSizeMode.medium:
      return '中';
    case FontSizeMode.large:
      return '大';
    case FontSizeMode.extraLarge:
      return '超大';
  }
}

/// 获取字体大小模式的多语言名称
String getLocalizedFontSizeModeName(BuildContext context, FontSizeMode mode) {
  switch (mode) {
    case FontSizeMode.extraSmall:
      return S.of(context).fontSizeExtraSmall;
    case FontSizeMode.small:
      return S.of(context).fontSizeSmall;
    case FontSizeMode.medium:
      return S.of(context).fontSizeMedium;
    case FontSizeMode.large:
      return S.of(context).fontSizeLarge;
    case FontSizeMode.extraLarge:
      return S.of(context).fontSizeExtraLarge;
  }
}

/// 根据字体大小模式和基础大小计算实际字体大小
double calculateFontSize(double baseSize, FontSizeMode mode) {
  switch (mode) {
    case FontSizeMode.extraSmall:
      return baseSize - 4;
    case FontSizeMode.small:
      return baseSize - 2;
    case FontSizeMode.medium:
      return baseSize;
    case FontSizeMode.large:
      return baseSize + 2;
    case FontSizeMode.extraLarge:
      return baseSize + 4;
  }
}

/// 主题管理器Provider
final themeManagerProvider = StateNotifierProvider<ThemeManager, ThemeMode>((ref) {
  return ThemeManager();
});

/// 字体大小设置的Provider
final fontSizeModeProvider = StateNotifierProvider<FontSizeModeNotifier, FontSizeMode>((ref) {
  return FontSizeModeNotifier();
});

/// 主题管理器
class ThemeManager extends StateNotifier<ThemeMode> {
  static const String _themeKey = 'theme_mode';
  
  ThemeManager() : super(ThemeMode.system) {
    _loadThemeMode();
  }

  /// 加载保存的主题模式
  Future<void> _loadThemeMode() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final themeIndex = prefs.getInt(_themeKey);
      if (themeIndex != null && themeIndex < ThemeMode.values.length) {
        state = ThemeMode.values[themeIndex];
      }
    } catch (e) {
      // 如果加载失败，使用默认的系统模式
      state = ThemeMode.system;
    }
  }

  /// 设置主题模式
  Future<void> setThemeMode(ThemeMode mode) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setInt(_themeKey, mode.index);
      state = mode;
    } catch (e) {
      // 如果保存失败，仍然更新状态
      state = mode;
    }
  }

  /// 切换主题模式
  Future<void> toggleTheme() async {
    ThemeMode newMode;
    switch (state) {
      case ThemeMode.light:
        newMode = ThemeMode.dark;
        break;
      case ThemeMode.dark:
        newMode = ThemeMode.system;
        break;
      case ThemeMode.system:
        newMode = ThemeMode.light;
        break;
    }
    await setThemeMode(newMode);
  }

  /// 获取当前主题模式的中文名称
  String getThemeModeName() {
    switch (state) {
      case ThemeMode.light:
        return '浅色模式';
      case ThemeMode.dark:
        return '深色模式';
      case ThemeMode.system:
        return '跟随系统';
    }
  }

  /// 获取当前主题模式的多语言名称
  String getLocalizedThemeModeName(BuildContext context) {
    switch (state) {
      case ThemeMode.light:
        return S.of(context).lightMode;
      case ThemeMode.dark:
        return S.of(context).darkMode;
      case ThemeMode.system:
        return S.of(context).systemMode;
    }
  }

  /// 获取下一个主题模式的多语言名称
  String getLocalizedNextThemeModeName(BuildContext context) {
    switch (state) {
      case ThemeMode.light:
        return S.of(context).darkMode;
      case ThemeMode.dark:
        return S.of(context).systemMode;
      case ThemeMode.system:
        return S.of(context).lightMode;
    }
  }

  /// 获取下一个主题模式的中文名称
  String getNextThemeModeName() {
    switch (state) {
      case ThemeMode.light:
        return '深色模式';
      case ThemeMode.dark:
        return '跟随系统';
      case ThemeMode.system:
        return '浅色模式';
    }
  }
}

/// 字体大小管理器
class FontSizeModeNotifier extends StateNotifier<FontSizeMode> {
  static const String _fontSizeKey = 'font_size_mode';
  
  FontSizeModeNotifier() : super(FontSizeMode.medium) {
    _loadFontSizeMode();
  }

  /// 加载保存的字体大小模式
  Future<void> _loadFontSizeMode() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final fontSizeIndex = prefs.getInt(_fontSizeKey);
      if (fontSizeIndex != null && fontSizeIndex < FontSizeMode.values.length) {
        state = FontSizeMode.values[fontSizeIndex];
      }
    } catch (e) {
      // 如果加载失败，使用默认的中等大小
      state = FontSizeMode.medium;
    }
  }

  /// 设置字体大小模式
  Future<void> setFontSizeMode(FontSizeMode mode) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setInt(_fontSizeKey, mode.index);
      state = mode;
    } catch (e) {
      // 如果保存失败，仍然更新状态
      state = mode;
    }
  }

  /// 获取当前字体大小模式的名称
  String getCurrentFontSizeModeName() {
    return getFontSizeModeName(state);
  }
  
  /// 获取当前字体大小模式的多语言名称
  String getLocalizedCurrentFontSizeModeName(BuildContext context) {
    return getLocalizedFontSizeModeName(context, state);
  }
} 