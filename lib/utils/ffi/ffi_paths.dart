import 'dart:io';
import 'package:path/path.dart' as path;
import 'package:flutter/foundation.dart';

/// 获取可用的动态库路径列表
List<String> getAvailableLibPaths() {
  final List<String> paths = [];
  
  // 根据不同平台添加不同的路径
  if (Platform.isMacOS) {
    // macOS平台路径
    paths.addAll([
      'example_ffi_arm64.dylib',
      'macos/Lib/example_ffi_arm64.dylib',
      'macos/Lib/example_ffi_chat_arm64.dylib',
      path.join(Directory.current.path, 'macos/Lib/example_ffi_arm64.dylib'),
      path.join(Directory.current.path, 'macos/Lib/example_ffi_chat_arm64.dylib'),
    ]);
  } else if (Platform.isWindows) {
    // Windows平台路径
    paths.addAll([
      'example_ffi.dll',
      'windows/lib/example_ffi.dll',
      path.join(Directory.current.path, 'windows/lib/example_ffi.dll'),
    ]);
  } else if (Platform.isLinux) {
    // Linux平台路径
    paths.addAll([
      'libexample_ffi.so',
      'linux/lib/libexample_ffi.so',
      path.join(Directory.current.path, 'linux/lib/libexample_ffi.so'),
    ]);
  } else if (Platform.isAndroid) {
    // Android平台路径
    paths.addAll([
      'libexample_ffi.so',
    ]);
  } else if (Platform.isIOS) {
    // iOS平台路径
    paths.addAll([
      'example_ffi.framework/example_ffi',
    ]);
  }
  
  return paths;
}

/// 获取特定类型的动态库路径列表
List<String> getLibPathsByType(String libType) {
  final List<String> paths = [];
  final String extension = _getPlatformExtension();
  
  switch (libType) {
    case 'chat':
      if (Platform.isMacOS) {
        paths.addAll([
          'example_ffi_chat_arm64.dylib',
          'macos/Lib/example_ffi_chat_arm64.dylib',
          path.join(Directory.current.path, 'macos/Lib/example_ffi_chat_arm64.dylib'),
        ]);
      } else {
        paths.add('example_ffi_chat$extension');
      }
      break;
    case 'core':
      if (Platform.isMacOS) {
        paths.addAll([
          'example_ffi_arm64.dylib',
          'macos/Lib/example_ffi_arm64.dylib',
          path.join(Directory.current.path, 'macos/Lib/example_ffi_arm64.dylib'),
        ]);
      } else {
        paths.add('example_ffi$extension');
      }
      break;
    default:
      // 默认返回所有可用路径
      return getAvailableLibPaths();
  }
  
  return paths;
}

/// 获取当前平台的动态库扩展名
String _getPlatformExtension() {
  if (Platform.isWindows) {
    return '.dll';
  } else if (Platform.isMacOS) {
    return '.dylib';
  } else if (Platform.isLinux || Platform.isAndroid) {
    return '.so';
  } else {
    return '';
  }
}

/// 获取动态库的完整路径
String getFullLibPath(String relativePath) {
  return path.join(Directory.current.path, relativePath);
} 