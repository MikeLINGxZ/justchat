import 'package:flutter/material.dart';

/// 主题工具类
/// 根据传入的context判断当前主题，并提供相应的自定义颜色
class Style {
  Style._(); // 私有构造函数，防止实例化

  /// 判断当前是否为暗黑主题
  static bool isDarkTheme(BuildContext context) {
    return Theme.of(context).brightness == Brightness.dark;
  }

  /// 判断当前是否为明亮主题
  static bool isLightTheme(BuildContext context) {
    return Theme.of(context).brightness == Brightness.light;
  }

  /// 获取当前主题的亮度
  static Brightness getBrightness(BuildContext context) {
    return Theme.of(context).brightness;
  }

  // ==================== 背景颜色 ====================

  /// 主背景色
  static Color primaryBackground(BuildContext context) {
    return isDarkTheme(context) 
        ? const Color(0xFF171717)  // 暗黑主题
        : const Color(0xFFFFFFFF); // 明亮主题
  }

  /// 次要背景色
  static Color secondaryBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF282828)  // 暗黑主题
        : const Color(0xFFededed); // 明亮主题
  }

  /// 卡片背景色
  static Color cardBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF2D2D2D)  // 暗黑主题
        : const Color(0xFFFFFFFF); // 明亮主题
  }

  /// 输入框背景色
  static Color inputBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF3D3D3D)  // 暗黑主题
        : const Color(0xFFF8F8F8); // 明亮主题
  }

  /// 悬浮背景色
  static Color hoverBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF404040)  // 暗黑主题
        : const Color(0xFFE0E0E0); // 明亮主题
  }

  // ==================== 文本颜色 ====================

  /// 主要文本颜色
  static Color primaryText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFFFFFFF)  // 暗黑主题
        : const Color(0xFF000000); // 明亮主题
  }

  /// 次要文本颜色
  static Color secondaryText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFB0B0B0)  // 暗黑主题
        : const Color(0xFF666666); // 明亮主题
  }

  /// 禁用文本颜色
  static Color disabledText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF666666)  // 暗黑主题
        : const Color(0xFFCCCCCC); // 明亮主题
  }

  /// 提示文本颜色
  static Color hintText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF888888)  // 暗黑主题
        : const Color(0xFF999999); // 明亮主题
  }

  /// 链接文本颜色
  static Color linkText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF4FC3F7)  // 暗黑主题
        : const Color(0xFF1976D2); // 明亮主题
  }

  // ==================== 边框颜色 ====================

  /// 主要边框颜色
  static Color primaryBorder(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF444444)  // 暗黑主题
        : const Color(0xFFE0E0E0); // 明亮主题
  }

  /// 次要边框颜色
  static Color secondaryBorder(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF2D2D2D)  // 暗黑主题
        : const Color(0xFFF0F0F0); // 明亮主题
  }

  /// 聚焦边框颜色
  static Color focusedBorder(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF6BB6FF)  // 暗黑主题
        : const Color(0xFF2196F3); // 明亮主题
  }

  /// 分割线颜色
  static Color divider(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF26303d)  // 暗黑主题
        : const Color(0xFFE8E8E8); // 明亮主题
  }

  // ==================== 状态颜色 ====================

  /// 成功颜色
  static Color success(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF4CAF50)  // 暗黑主题
        : const Color(0xFF388E3C); // 明亮主题
  }

  /// 错误颜色
  static Color error(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFF44336)  // 暗黑主题
        : const Color(0xFFD32F2F); // 明亮主题
  }

  /// 警告颜色
  static Color warning(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFFF9800)  // 暗黑主题
        : const Color(0xFFF57C00); // 明亮主题
  }

  /// 信息颜色
  static Color info(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF2196F3)  // 暗黑主题
        : const Color(0xFF1976D2); // 明亮主题
  }

  // ==================== 按钮颜色 ====================

  /// 主要按钮背景色
  static Color primaryButton(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF6200EA)  // 暗黑主题
        : const Color(0xFF3F51B5); // 明亮主题
  }

  /// 次要按钮背景色
  static Color secondaryButton(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF424242)  // 暗黑主题
        : const Color(0xFFE0E0E0); // 明亮主题
  }

  /// 按钮文本颜色
  static Color buttonText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFFFFFFF)  // 暗黑主题
        : const Color(0xFFFFFFFF); // 明亮主题
  }

  /// 次要按钮文本颜色
  static Color secondaryButtonText(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFFFFFFFF)  // 暗黑主题
        : const Color(0xFF000000); // 明亮主题
  }

  // ==================== 特殊组件颜色 ====================

  /// 侧边栏背景色
  static Color sidebarBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF0d0d0d)  // 暗黑主题
        : const Color(0xFFffffff); // 明亮主题
  }

  /// 标题栏背景色
  static Color titleBarBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF282828)  // 暗黑主题
        : const Color(0xFFededed); // 明亮主题
  }

  /// 聊天气泡背景色（用户）
  static Color userChatBubble(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF262626)  // 暗黑主题
        : const Color(0xFFf9f9f9); // 明亮主题
  }

  /// 聊天气泡背景色（助手）
  static Color assistantChatBubble(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0x00FFFFFF)  // 暗黑主题
        : const Color(0x00FFFFFF); // 明亮主题
  }

  /// 代码块背景色
  static Color codeBackground(BuildContext context) {
    return isDarkTheme(context)
        ? const Color(0xFF1E1E1E)  // 暗黑主题
        : const Color(0xFFF8F8F8); // 明亮主题
  }

  // ==================== 阴影和效果 ====================

  /// 获取卡片阴影
  static List<BoxShadow> cardShadow(BuildContext context) {
    return isDarkTheme(context)
        ? [
            BoxShadow(
              color: Colors.black.withOpacity(0.3),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ]
        : [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ];
  }

  /// 获取轻微阴影
  static List<BoxShadow> subtleShadow(BuildContext context) {
    return isDarkTheme(context)
        ? [
            BoxShadow(
              color: Colors.black.withOpacity(0.2),
              blurRadius: 4,
              offset: const Offset(0, 1),
            ),
          ]
        : [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 4,
              offset: const Offset(0, 1),
            ),
          ];
  }

  // ==================== 透明度变体 ====================

  /// 带透明度的主要背景色
  static Color primaryBackgroundWithOpacity(BuildContext context, double opacity) {
    return primaryBackground(context).withOpacity(opacity);
  }

  /// 带透明度的主要文本色
  static Color primaryTextWithOpacity(BuildContext context, double opacity) {
    return primaryText(context).withOpacity(opacity);
  }

  /// 带透明度的主要按钮色
  static Color primaryButtonWithOpacity(BuildContext context, double opacity) {
    return primaryButton(context).withOpacity(opacity);
  }

  static double radiusLv1 = 10.0;
}
