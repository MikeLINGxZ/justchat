import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/setting/manager.dart';

/// 字体大小工具类
class FontSizeUtils {
  /// 获取调整后的字体大小
  static double getAdjustedFontSize(WidgetRef ref, double baseSize) {
    final fontSizeMode = ref.watch(fontSizeModeProvider);
    return calculateFontSize(baseSize, fontSizeMode);
  }

  /// 获取调整后的标题字体大小 (24)
  static double getTitleLargeSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 24);
  }

  /// 获取调整后的标题字体大小 (22)
  static double getTitleSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 22);
  }

  /// 获取调整后的标题字体大小 (20)
  static double getHeadingSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 20);
  }

  /// 获取调整后的副标题字体大小 (18)
  static double getSubheadingSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 18);
  }

  /// 获取调整后的正文大字体大小 (16)
  static double getBodyLargeSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 16);
  }

  /// 获取调整后的正文字体大小 (14)
  static double getBodySize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 14);
  }

  /// 获取调整后的小字体大小 (13)
  static double getSmallSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 13);
  }

  /// 获取调整后的微小字体大小 (11)
  static double getXSmallSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 11);
  }

  /// 获取按钮文字大小 (14)
  static double getButtonSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 14);
  }

  /// 获取提示文字大小 (12)
  static double getCaptionSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 12);
  }

  /// 获取辅助说明文字大小 (10)
  static double getHelperTextSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 10);
  }
}