import 'dart:io';
import 'dart:convert';
import 'dart:isolate';
import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:path/path.dart' as path;
import 'package:lemon_tea/utils/system.dart';

/// CLI服务类，用于管理CLI二进制文件的启动和停止
class CliService {
  /// 单例实例
  static final CliService _instance = CliService._internal();

  /// 工厂构造函数
  factory CliService() => _instance;

  /// 内部构造函数
  CliService._internal();

  /// 服务是否正在运行
  bool _isRunning = false;

  /// 服务端口
  int? _port;

  /// 获取服务是否正在运行
  bool get isRunning => _isRunning;

  /// 获取服务端口
  int? get port => _port;

  /// 获取CLI二进制文件路径
  String _getCliBinaryPath() {
    // 获取应用程序目录
    final String appDir = path.dirname(Platform.resolvedExecutable);
    String binaryPath;

    switch (System.platform) {
      case 'windows':
      // Windows平台
        final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
        binaryPath = path.join(
          appDir,
          'data',
          'flutter_assets',
          'cli',
          'lemon_tea_local_windows_$arch.exe',
        );
        break;

      case 'macos':
      // macOS平台
        final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
        binaryPath = path.join(
          appDir,
          '..',
          'Resources',
          'lemon_tea_local_darwin_$arch',
        );
        // 检查标准路径是否存在
        if (!File(binaryPath).existsSync()) {
          throw UnsupportedError('不存在local_server: ${binaryPath}');
        }
        break;

      case 'linux':
      // Linux平台
        final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
        binaryPath = path.join(
          appDir,
          'data',
          'flutter_assets',
          'cli',
          'lemon_tea_local_linux_$arch',
        );
        break;

      default:
        throw UnsupportedError('不支持的平台: ${Platform.operatingSystem}');
    }

    return binaryPath;
  }

  /// 启动CLI服务
  ///
  /// [requestedPort] 请求使用的端口号，如果为null则自动分配
  /// 返回服务端口号，如果启动失败则返回null
  Future<int?> startService({int? requestedPort}) async {
    try {
      // 获取端口
      int? port;
      if (requestedPort != null &&
          await System.isPortAvailable(requestedPort)) {
        port = requestedPort;
      }

      // 如果未指定端口或请求的端口不可用，则自动分配
      if (port == null) {
        port = await System.findFreePort();
        if (port == null) {
          debugPrint('无法获取空闲端口');
          return null;
        }
      }

      // 获取二进制文件路径
      final String binaryPath = _getCliBinaryPath();
      final File binaryFile = File(binaryPath);

      // 检查二进制文件是否存在
      if (!await binaryFile.exists()) {
        debugPrint('CLI二进制文件不存在: $binaryPath');
        return null;
      }

      // 确保二进制文件有执行权限（仅在非Windows平台）
      if (!System.isWindows) {
        try {
          await Process.run('chmod', ['+x', binaryPath]);
        } catch (e) {
          debugPrint('设置执行权限失败: $e');
        }
      }

      // todo 执行启动命令 {binaryPath} --port {port}


      return null;
    } catch (e) {
      debugPrint('启动CLI服务失败: $e');
      return null;
    }

  }

  stopService() {

  }
}
