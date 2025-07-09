import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:lemon_tea/generated/l10n.dart';

/// 字体大小设置
enum FontSizeMode {
  extraSmall,
  small,
  medium,
  large,
  extraLarge,
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

/// 字体大小设置的Provider
final fontSizeModeProvider = StateNotifierProvider<FontSizeModeNotifier, FontSizeMode>((ref) {
  return FontSizeModeNotifier();
});

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