import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:lemon_tea/generated/l10n.dart';

/// 语言设置
enum AppLanguage {
  system,   // 跟随系统
  english,  // 英文
  chinese,  // 中文
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

/// 语言设置的Provider
final appLanguageProvider = StateNotifierProvider<AppLanguageNotifier, AppLanguage>((ref) {
  return AppLanguageNotifier();
});

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