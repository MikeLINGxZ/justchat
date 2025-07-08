import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/setting/manager.dart';

/// 字体大小工具类
class FontSizeUtils {
  /// 获取调整后的字体大小
  static double getAdjustedFontSize(WidgetRef ref, double baseSize) {
    final fontSizeMode = ref.watch(fontSizeModeProvider);
    return calculateFontSize(baseSize, fontSizeMode);
  }

  /// 获取调整后的标题字体大小 (20)
  static double getHeadingSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 19);
  }

  /// 获取调整后的副标题字体大小 (16)
  static double getSubheadingSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 17);
  }

  /// 获取调整后的正文字体大小 (14)
  static double getBodySize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 15);
  }

  /// 获取调整后的小字体大小 (12)
  static double getSmallSize(WidgetRef ref) {
    return getAdjustedFontSize(ref, 13);
  }
} 