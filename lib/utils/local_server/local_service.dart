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
  
  /// CLI进程
  Process? _cliProcess;
  
  /// 服务是否正在运行
  bool _isRunning = false;
  
  /// 服务端口
  int? _port;
  
  /// Isolate实例
  Isolate? _isolate;
  
  /// 发送端口
  SendPort? _sendPort;
  
  /// 接收端口
  ReceivePort? _receivePort;
  
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
  
  /// 在Isolate中运行CLI服务的入口点
  static Future<void> _isolateEntryPoint(Map<String, dynamic> params) async {
    final SendPort sendPort = params['sendPort'] as SendPort;
    final String binaryPath = params['binaryPath'] as String;
    final int port = params['port'] as int;
    
    try {
      // 打印启动信息
      debugPrint('Isolate: 正在启动CLI进程: $binaryPath --port $port');
      
      // 启动CLI进程
      final process = await Process.start(
        binaryPath,
        ['--port', port.toString()],
        mode: ProcessStartMode.normal,
        runInShell: true,
      );
      
      // 发送进程启动成功消息
      sendPort.send({'status': 'started', 'port': port});
      
      // 监听进程输出
      process.stdout.transform(utf8.decoder).listen(
        (data) {
          debugPrint('Isolate: CLI输出: $data');
          sendPort.send({'type': 'stdout', 'data': data});
        },
        onError: (error) {
          debugPrint('Isolate: CLI输出流错误: $error');
          sendPort.send({'type': 'error', 'data': 'stdout错误: $error'});
        },
        onDone: () {
          debugPrint('Isolate: CLI输出流已关闭');
        },
      );
      
      process.stderr.transform(utf8.decoder).listen(
        (data) {
          debugPrint('Isolate: CLI错误: $data');
          sendPort.send({'type': 'stderr', 'data': data});
        },
        onError: (error) {
          debugPrint('Isolate: CLI错误流错误: $error');
          sendPort.send({'type': 'error', 'data': 'stderr错误: $error'});
        },
        onDone: () {
          debugPrint('Isolate: CLI错误流已关闭');
        },
      );
      
      // 监听进程退出
      final exitCode = await process.exitCode;
      debugPrint('Isolate: CLI进程退出，退出码: $exitCode');
      sendPort.send({'status': 'exited', 'exitCode': exitCode});
      
    } catch (e) {
      debugPrint('Isolate: 启动CLI进程失败: $e');
      sendPort.send({'status': 'error', 'message': e.toString()});
    }
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
      
      // 创建接收端口
      _receivePort = ReceivePort();
      
      // 创建Completer用于等待服务启动
      final completer = Completer<int?>();
      
      // 监听接收端口
      _receivePort!.listen((message) {
        if (message is Map<String, dynamic>) {
          if (message['status'] == 'started') {
            // 服务启动成功
            if (!completer.isCompleted) {
              _isRunning = true;
              _port = message['port'] as int;
              debugPrint('CLI服务已在Isolate中启动，端口: $_port');
              completer.complete(_port);
            }
          } else if (message['status'] == 'exited') {
            // 服务退出
            debugPrint('CLI服务已退出，退出码: ${message['exitCode']}');
            _isRunning = false;
            _port = null;
          } else if (message['status'] == 'error') {
            // 服务启动失败
            if (!completer.isCompleted) {
              debugPrint('启动CLI服务失败: ${message['message']}');
              completer.complete(null);
            }
          } else if (message['type'] == 'stdout' || message['type'] == 'stderr') {
            // 处理进程输出
            final type = message['type'] as String;
            final data = message['data'] as String;
            debugPrint('CLI ${type == 'stdout' ? '输出' : '错误'}: $data');
          }
        }
      });
      
      // 启动Isolate
      debugPrint('正在启动CLI服务Isolate...');
      _isolate = await Isolate.spawn(
        _isolateEntryPoint, 
        {
          'sendPort': _receivePort!.sendPort,
          'binaryPath': binaryPath,
          'port': port,
        }
      );
      
      // 等待服务启动或失败
      return await completer.future.timeout(
        const Duration(seconds: 10),
        onTimeout: () {
          debugPrint('启动CLI服务超时');
          stopService();
          return null;
        },
      );
    } catch (e) {
      debugPrint('启动CLI服务失败: $e');
      // 清理资源
      _cleanupResources();
      return null;
    }
  }
  
  /// 清理资源
  void _cleanupResources() {
    _isolate?.kill(priority: Isolate.immediate);
    _isolate = null;
    _receivePort?.close();
    _receivePort = null;
    _sendPort = null;
    _isRunning = false;
    _port = null;
  }
  
  /// 停止CLI服务
  Future<void> stopService() async {
    if (!_isRunning) {
      debugPrint('CLI服务未运行');
      return;
    }
    
    try {
      debugPrint('正在停止CLI服务...');
      
      // 终止Isolate
      _cleanupResources();
      
      debugPrint('CLI服务已停止');
    } catch (e) {
      debugPrint('停止CLI服务失败: $e');
      // 确保状态重置
      _cleanupResources();
    }
  }
} 