import 'dart:ui';

Color hexToColor(String hexCode) {
  // 去掉开头的 '#'（如果有）
  hexCode = hexCode.replaceAll('#', '');

  // 如果没有 Alpha 通道，默认添加 FF
  if (hexCode.length == 6) {
    hexCode = 'FF$hexCode';
  }

  // 将 16 进制字符串转换为整数
  return Color(int.parse(hexCode, radix: 16));
}