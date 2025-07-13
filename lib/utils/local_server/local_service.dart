import 'dart:io';
import 'dart:convert';
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
  
  /// CLI进程
  Process? _cliProcess;
  
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
    
    if (System.isWindows) {
      // todo Windows平台
      final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
      binaryPath = path.join(appDir, 'data', 'flutter_assets', 'cli', 'lemon_tea_local_windows_$arch.exe');
    } else if (System.isMacOS) {
      // macOS平台
      final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
      binaryPath = path.join(appDir, '..', 'Resources',  'lemon_tea_local_darwin_$arch');
      // 检查标准路径是否存在
      if (!File(binaryPath).existsSync()) {
        throw UnsupportedError('不存在local_server: ${binaryPath}');
      }
    } else if (System.isLinux) {
      // todo Linux平台
      final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
      binaryPath = path.join(appDir, 'data', 'flutter_assets', 'cli', 'lemon_tea_local_linux_$arch');
    } else {
      throw UnsupportedError('不支持的平台: ${Platform.operatingSystem}');
    }
    
    return binaryPath;
  }
  
  /// 启动CLI服务
  /// 
  /// 返回服务端口号，如果启动失败则返回null
  Future<int?> startService() async {
    if (_isRunning) {
      debugPrint('CLI服务已经在运行中，端口: $_port');
      return _port;
    }
    
    try {
      // 获取空闲端口
      final port = await System.findFreePort();
      if (port == null) {
        debugPrint('无法获取空闲端口');
        return null;
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
      
      // 启动CLI进程
      _cliProcess = await Process.start(
        binaryPath,
        ['--port', port.toString()],
        mode: ProcessStartMode.detached,
      );
      
      // 监听进程输出
      _cliProcess!.stdout.transform(utf8.decoder).listen((data) {
        debugPrint('CLI输出: $data');
      });
      
      _cliProcess!.stderr.transform(utf8.decoder).listen((data) {
        debugPrint('CLI错误: $data');
      });
      
      // 监听进程退出
      _cliProcess!.exitCode.then((exitCode) {
        debugPrint('CLI进程退出，退出码: $exitCode');
        _isRunning = false;
        _port = null;
      });
      
      _isRunning = true;
      _port = port;
      
      debugPrint('CLI服务已启动，端口: $port');
      return port;
    } catch (e) {
      debugPrint('启动CLI服务失败: $e');
      return null;
    }
  }
  
  /// 停止CLI服务
  Future<void> stopService() async {
    if (!_isRunning || _cliProcess == null) {
      debugPrint('CLI服务未运行');
      return;
    }
    
    try {
      // 尝试正常终止进程
      _cliProcess!.kill();
      
      // 等待进程退出
      await _cliProcess!.exitCode.timeout(
        const Duration(seconds: 5),
        onTimeout: () {
          // 如果进程没有在5秒内退出，则强制终止
          _cliProcess!.kill(ProcessSignal.sigkill);
          return -1;
        },
      );
      
      _isRunning = false;
      _port = null;
      _cliProcess = null;
      
      debugPrint('CLI服务已停止');
    } catch (e) {
      debugPrint('停止CLI服务失败: $e');
    }
  }
} 