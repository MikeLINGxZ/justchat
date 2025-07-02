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

/// 语言设置
enum AppLanguage {
  system,   // 跟随系统
  english,  // 英文
  chinese,  // 中文
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

/// 获取语言设置的名称
String getAppLanguageName(AppLanguage language) {
  switch (language) {
    case AppLanguage.system:
      return '跟随系统';
    case AppLanguage.english:
      return 'English';
    case AppLanguage.chinese:
      return '中文';
  }
}

/// 获取语言设置的多语言名称
String getLocalizedAppLanguageName(BuildContext context, AppLanguage language) {
  switch (language) {
    case AppLanguage.system:
      return S.of(context).systemMode;
    case AppLanguage.english:
      return 'English';
    case AppLanguage.chinese:
      return '中文';
  }
}

/// 根据语言设置获取对应的Locale
Locale? getLocaleFromAppLanguage(AppLanguage language) {
  switch (language) {
    case AppLanguage.system:
      return null; // 返回null表示跟随系统
    case AppLanguage.english:
      return const Locale('en');
    case AppLanguage.chinese:
      return const Locale('zh', 'CN');
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

/// 语言设置的Provider
final appLanguageProvider = StateNotifierProvider<AppLanguageNotifier, AppLanguage>((ref) {
  return AppLanguageNotifier();
});

/// 主题管理器
class ThemeManager extends StateNotifier<ThemeMode> {
  static const String _themeKey = 'theme_mode';
  
  ThemeManager() : super(ThemeMode.system) {
    loadThemeMode();
  }

  /// 加载保存的主题模式
  Future<void> loadThemeMode() async {
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
}

/// 字体大小管理器
class FontSizeModeNotifier extends StateNotifier<FontSizeMode> {
  static const String _fontSizeKey = 'font_size_mode';
  
  FontSizeModeNotifier() : super(FontSizeMode.medium) {
    loadFontSizeMode();
  }

  /// 加载保存的字体大小模式
  Future<void> loadFontSizeMode() async {
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

  /// 获取当前字体大小模式的多语言名称
  String getLocalizedCurrentFontSizeModeName(BuildContext context) {
    return getLocalizedFontSizeModeName(context, state);
  }
}

/// 语言设置管理器
class AppLanguageNotifier extends StateNotifier<AppLanguage> {
  static const String _languageKey = 'app_language';
  
  AppLanguageNotifier() : super(AppLanguage.system) {
    _loadLanguage();
  }

  /// 加载保存的语言设置
  Future<void> _loadLanguage() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final languageIndex = prefs.getInt(_languageKey);
      if (languageIndex != null && languageIndex < AppLanguage.values.length) {
        state = AppLanguage.values[languageIndex];
      }
    } catch (e) {
      // 如果加载失败，使用默认的系统语言
      state = AppLanguage.system;
    }
  }

  /// 设置语言
  Future<void> setLanguage(AppLanguage language) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setInt(_languageKey, language.index);
      state = language;
    } catch (e) {
      // 如果保存失败，仍然更新状态
      state = language;
    }
  }

  /// 获取当前语言设置的名称
  String getCurrentLanguageName() {
    return getAppLanguageName(state);
  }
  
  /// 获取当前语言设置的多语言名称
  String getLocalizedCurrentLanguageName(BuildContext context) {
    return getLocalizedAppLanguageName(context, state);
  }
  
  /// 获取当前语言对应的Locale
  Locale? getCurrentLocale() {
    return getLocaleFromAppLanguage(state);
  }
} 