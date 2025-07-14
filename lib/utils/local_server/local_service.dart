import 'dart:io';
import 'dart:convert';
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
  
  /// 状态监听器
  final List<Function(bool isRunning, int? port)> _statusListeners = [];
  
  /// 最后一次检查时间
  DateTime? _lastHeartbeat;
  
  /// 心跳检测计时器
  Timer? _heartbeatTimer;
  
  /// 进程监控计时器
  Timer? _processMonitorTimer;
  
  /// 进程输出流订阅
  StreamSubscription? _stdoutSubscription;
  
  /// 进程错误流订阅
  StreamSubscription? _stderrSubscription;
  
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
      // Windows平台
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
      // Linux平台
      final String arch = Platform.version.contains('arm') ? 'arm64' : 'amd64';
      binaryPath = path.join(appDir, 'data', 'flutter_assets', 'cli', 'lemon_tea_local_linux_$arch');
    } else {
      throw UnsupportedError('不支持的平台: ${Platform.operatingSystem}');
    }
    
    return binaryPath;
  }
  
  /// 添加状态监听器
  void addStatusListener(Function(bool isRunning, int? port) listener) {
    _statusListeners.add(listener);
  }
  
  /// 移除状态监听器
  void removeStatusListener(Function(bool isRunning, int? port) listener) {
    _statusListeners.remove(listener);
  }
  
  /// 通知所有监听器状态变化
  void _notifyStatusChange() {
    for (final listener in _statusListeners) {
      listener(_isRunning, _port);
    }
  }
  
  /// 启动心跳检测
  void _startHeartbeatMonitor() {
    _lastHeartbeat = DateTime.now();
    
    // 取消现有计时器
    _heartbeatTimer?.cancel();
    
    // 创建新计时器，每5秒检查一次心跳
    _heartbeatTimer = Timer.periodic(const Duration(seconds: 5), (timer) {
      if (_lastHeartbeat != null) {
        final difference = DateTime.now().difference(_lastHeartbeat!);
        
        // 如果超过10秒没有收到心跳，认为服务已经异常退出
        if (difference.inSeconds > 10 && _isRunning) {
          debugPrint('CLI服务心跳超时，认为服务已异常退出');
          _isRunning = false;
          _notifyStatusChange();
        }
      }
    });
  }
  
  /// 停止心跳检测
  void _stopHeartbeatMonitor() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = null;
    _lastHeartbeat = null;
  }
  
  /// 启动进程监控
  void _startProcessMonitor() {
    // 取消现有计时器
    _processMonitorTimer?.cancel();
    
    // 创建新计时器，每2秒检查一次进程状态
    _processMonitorTimer = Timer.periodic(const Duration(seconds: 2), (timer) {
      if (_cliProcess != null) {
        try {
          // 尝试获取进程PID，如果进程已经退出，这里会抛出异常
          final pid = _cliProcess!.pid;
          if (pid != 0) {
            // 更新心跳时间
            _lastHeartbeat = DateTime.now();
          }
        } catch (e) {
          // 进程已经退出
          debugPrint('CLI进程已退出');
          _cleanupResources();
        }
      }
    });
  }
  
  /// 停止进程监控
  void _stopProcessMonitor() {
    _processMonitorTimer?.cancel();
    _processMonitorTimer = null;
  }
  
  /// 启动CLI服务
  /// 
  /// [requestedPort] 请求使用的端口号，如果为null则自动分配
  /// 返回服务端口号，如果启动失败则返回null
  Future<int?> startService({int? requestedPort}) async {
    // 确保先停止任何现有的服务
    await stopService();
    
    try {
      // 获取端口
      int? port;
      if (requestedPort != null) {
        // 检查请求的端口是否可用
        final isAvailable = await System.isPortAvailable(requestedPort);
        if (isAvailable) {
          port = requestedPort;
        } else {
          debugPrint('请求的端口 $requestedPort 不可用，将自动分配端口');
        }
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
      
      // 创建Completer用于等待服务启动
      final completer = Completer<int?>();
      
      try {
        // 启动CLI进程
        debugPrint('正在启动CLI进程: $binaryPath --port $port');
        
        // 使用Process.run而不是Process.start，这样可以避免stdio连接问题
        // 但这意味着我们无法获取实时输出，而是在进程结束后获取所有输出
        if (System.isWindows) {
          // Windows平台使用不同的参数
          _cliProcess = await Process.start(
            binaryPath,
            ['--port', port.toString()],
            mode: ProcessStartMode.normal,  // 在Windows上使用normal模式
            runInShell: true,
          );
        } else {
          // macOS和Linux平台
          _cliProcess = await Process.start(
            binaryPath,
            ['--port', port.toString()],
            mode: ProcessStartMode.normal,  // 使用normal模式代替detached
            runInShell: false,  // 不使用shell
          );
        }
        
        // 设置超时
        final timeout = Timer(const Duration(seconds: 10), () {
          if (!completer.isCompleted) {
            debugPrint('启动CLI服务超时');
            stopService();
            completer.complete(null);
          }
        });
        
        // 监听进程输出
        _stdoutSubscription = _cliProcess!.stdout.transform(utf8.decoder).listen(
          (data) {
            debugPrint('CLI输出: $data');
            
            // 检查输出中是否包含服务启动成功的信息
            if (data.contains('服务已启动') || data.contains('server started') || data.contains('listening on')) {
              if (!completer.isCompleted) {
                _isRunning = true;
                _port = port;
                debugPrint('CLI服务已启动，端口: $_port');
                
                // 启动心跳监控和进程监控
                _startHeartbeatMonitor();
                _startProcessMonitor();
                
                // 通知状态变化
                _notifyStatusChange();
                
                // 取消超时
                timeout.cancel();
                
                completer.complete(_port);
              }
            }
          },
          onError: (error) {
            debugPrint('CLI输出流错误: $error');
            if (!completer.isCompleted) {
              completer.complete(null);
            }
          },
          onDone: () {
            debugPrint('CLI输出流已关闭');
          },
        );
        
        _stderrSubscription = _cliProcess!.stderr.transform(utf8.decoder).listen(
          (data) {
            debugPrint('CLI错误: $data');
          },
          onError: (error) {
            debugPrint('CLI错误流错误: $error');
          },
          onDone: () {
            debugPrint('CLI错误流已关闭');
          },
        );
        
        // 监听进程退出
        unawaited(_cliProcess!.exitCode.then((exitCode) {
          debugPrint('CLI进程退出，退出码: $exitCode');
          
          final wasRunning = _isRunning;
          _isRunning = false;
          _port = null;
          
          // 停止监控
          _stopHeartbeatMonitor();
          _stopProcessMonitor();
          
          // 清理资源
          _cleanupResources();
          
          // 只有在之前是运行状态时才通知状态变化
          if (wasRunning) {
            _notifyStatusChange();
          }
          
          // 如果completer还没有完成，则完成它
          if (!completer.isCompleted) {
            completer.complete(null);
          }
        }));
        
        // 假设进程已经启动成功，等待一段时间后检查
        Timer(const Duration(seconds: 2), () {
          if (!completer.isCompleted) {
            try {
              if (_cliProcess != null && _cliProcess!.pid != 0) {
                _isRunning = true;
                _port = port;
                debugPrint('CLI服务已启动（假定），端口: $_port');
                
                // 启动心跳监控和进程监控
                _startHeartbeatMonitor();
                _startProcessMonitor();
                
                // 通知状态变化
                _notifyStatusChange();
                
                // 取消超时
                timeout.cancel();
                
                completer.complete(_port);
              }
            } catch (e) {
              debugPrint('检查进程状态失败: $e');
            }
          }
        });
      } catch (e) {
        debugPrint('启动CLI进程失败: $e');
        if (!completer.isCompleted) {
          completer.complete(null);
        }
      }
      
      // 等待服务启动或失败
      return await completer.future;
    } catch (e) {
      debugPrint('启动CLI服务失败: $e');
      // 清理资源
      _cleanupResources();
      return null;
    }
  }
  
  /// 清理资源
  void _cleanupResources() {
    // 取消流订阅
    _stdoutSubscription?.cancel();
    _stdoutSubscription = null;
    _stderrSubscription?.cancel();
    _stderrSubscription = null;
    
    // 停止监控
    _stopHeartbeatMonitor();
    _stopProcessMonitor();
    
    // 终止进程
    if (_cliProcess != null) {
      try {
        // 尝试正常终止进程
        if (System.isWindows) {
          // Windows平台使用taskkill终止进程
          Process.runSync('taskkill', ['/F', '/PID', '${_cliProcess!.pid}']);
        } else {
          // macOS和Linux平台使用kill命令
          Process.runSync('kill', ['-9', '${_cliProcess!.pid}']);
        }
        
        // 尝试使用Dart API终止进程（作为备份）
        _cliProcess!.kill();
      } catch (e) {
        debugPrint('终止CLI进程失败: $e');
      } finally {
        _cliProcess = null;
      }
    }
    
    final wasRunning = _isRunning;
    _isRunning = false;
    _port = null;
    
    // 只有在之前是运行状态时才通知状态变化
    if (wasRunning) {
      _notifyStatusChange();
    }
  }
  
  /// 停止CLI服务
  Future<void> stopService() async {
    if (!_isRunning && _cliProcess == null) {
      debugPrint('CLI服务未运行');
      return;
    }
    
    try {
      debugPrint('正在停止CLI服务...');
      
      // 清理资源和终止进程
      _cleanupResources();
      
      debugPrint('CLI服务已停止');
    } catch (e) {
      debugPrint('停止CLI服务失败: $e');
      // 确保状态重置
      _cleanupResources();
    }
  }
} 