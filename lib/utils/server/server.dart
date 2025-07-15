import 'dart:io';
import 'dart:convert';
import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:path/path.dart' as path;
import 'package:lemon_tea/utils/system.dart';

/// CLI服务类，用于管理CLI二进制文件的启动和停止
class Server {
  /// 单例实例
  static final Server _instance = Server._internal();

  /// 工厂构造函数
  factory Server() => _instance;

  /// 内部构造函数
  Server._internal() {
    // 注册应用退出钩子
    _registerExitHook();
  }

  /// 服务是否正在运行
  bool _isRunning = false;

  /// 服务端口
  int? _port;

  /// 进程对象
  Process? _process;

  /// 进程ID
  int? _processId;

  /// 当前应用进程ID
  final int _currentPid = pid;

  /// 监视进程
  Process? _watcherProcess;

  /// 进程输出流订阅
  StreamSubscription? _stdoutSubscription;
  StreamSubscription? _stderrSubscription;

  /// 获取服务是否正在运行
  bool get isRunning => _isRunning;

  /// 获取服务端口
  int? get port => _port;

  /// 注册应用退出钩子
  void _registerExitHook() {
    // 注册进程退出处理器
    ProcessSignal.sigint.watch().listen((_) {
      _forceKillProcess();
      exit(0);
    });

    if (!System.isWindows) {
      ProcessSignal.sigterm.watch().listen((_) {
        _forceKillProcess();
        exit(0);
      });
    }

    // 创建退出文件
    _createExitFile();
  }

  /// 创建退出文件，记录当前进程ID和子进程ID
  void _createExitFile() {
    try {
      final tempDir = Directory.systemTemp;
      final exitFile = File(
        '${tempDir.path}/lemon_tea_exit_${_currentPid}.pid',
      );

      // 确保文件存在
      if (!exitFile.existsSync()) {
        exitFile.createSync();
      }

      // 添加删除钩子
      ProcessSignal.sigint.watch().listen((_) {
        if (exitFile.existsSync()) {
          exitFile.deleteSync();
        }
      });

      if (!System.isWindows) {
        ProcessSignal.sigterm.watch().listen((_) {
          if (exitFile.existsSync()) {
            exitFile.deleteSync();
          }
        });
      }
    } catch (e) {
      debugPrint('创建退出文件失败: $e');
    }
  }

  /// 强制终止进程
  void _forceKillProcess() {
    try {
      // 先终止监视进程
      _watcherProcess?.kill();
      _watcherProcess = null;

      // 然后终止主进程
      if (_processId != null) {
        if (System.isWindows) {
          Process.runSync('taskkill', ['/F', '/PID', _processId.toString()]);
        } else {
          Process.runSync('kill', ['-9', _processId.toString()]);
        }
        _processId = null;
      }
      _cleanupProcess();
    } catch (e) {
      debugPrint('强制终止进程失败: $e');
    }
  }

  /// 获取CLI二进制文件路径
  String _getCliBinaryPath({String? customBinaryPath}) {
    // 如果提供了自定义路径，优先使用
    if (customBinaryPath != null && customBinaryPath.isNotEmpty) {
      final customFile = File(customBinaryPath);
      if (customFile.existsSync()) {
        return customBinaryPath;
      } else {
        debugPrint('自定义二进制文件不存在: $customBinaryPath，将使用默认路径');
      }
    }

    // 获取应用程序目录
    final String appDir = path.dirname(Platform.resolvedExecutable);
    String binaryPath;

    switch (System.platform) {
      case 'windows':
        // Windows平台
        final String arch =
            Platform.version.contains('arm') ? 'arm64' : 'amd64';
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
        final String arch =
            Platform.version.contains('arm') ? 'arm64' : 'amd64';
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
        final String arch =
            Platform.version.contains('arm') ? 'arm64' : 'amd64';
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
  /// [customBinaryPath] 自定义二进制文件路径，如果为null则使用默认路径
  /// 返回服务端口号，如果启动失败则返回null
  Future<int?> startService({
    int? requestedPort,
    String? customBinaryPath,
  }) async {
    try {
      // 如果服务已经在运行，直接返回当前端口
      if (_isRunning && _port != null) {
        return _port;
      }

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
      final String binaryPath = _getCliBinaryPath(
        customBinaryPath: customBinaryPath,
      );
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

      // 直接启动进程，不使用脚本
      _process = await Process.start(binaryPath, [
        '--port',
        port.toString(),
        '--debug',
        _isDebugMode() ? 'true' : 'false',
      ]);

      // 保存进程ID
      _processId = _process!.pid;

      // 记录进程ID到退出文件
      _updateExitFile();

      // 监听进程输出
      _stdoutSubscription = _process!.stdout.transform(utf8.decoder).listen((
        data,
      ) {
        debugPrint('CLI输出: $data');
      });

      _stderrSubscription = _process!.stderr.transform(utf8.decoder).listen((
        data,
      ) {
        debugPrint('CLI错误: $data');
      });

      // 监听进程退出
      _process!.exitCode.then((exitCode) {
        debugPrint('CLI进程退出，退出码: $exitCode');
        _cleanupProcess();
      });

      // 等待一小段时间确保进程启动成功
      await Future.delayed(const Duration(seconds: 1));

      // 检查服务是否可用
      if (!await _isServiceAvailable(port)) {
        debugPrint('CLI服务启动失败，无法连接到端口: $port');
        _forceKillProcess();
        return null;
      }

      _isRunning = true;
      _port = port;

      // 仅在应用程序打包发布时注册进程清理
      // 在开发模式下，不注册监视进程，避免误杀
      if (!_isDebugMode()) {
        _registerProcessCleanup();
      }

      return port;
    } catch (e) {
      debugPrint('启动CLI服务失败: $e');
      _forceKillProcess();
      return null;
    }
  }

  /// 检查是否处于调试模式
  bool _isDebugMode() {
    bool inDebugMode = false;
    assert(() {
      inDebugMode = true;
      return true;
    }());
    return inDebugMode;
  }

  /// 注册应用退出时的进程清理
  void _registerProcessCleanup() async {
    // 创建一个独立的进程来监视主进程
    if (_processId != null) {
      try {
        if (System.isWindows) {
          // Windows平台使用批处理脚本监视父进程
          final tempDir = Directory.systemTemp;
          final watcherScript = File(
            '${tempDir.path}/lemon_tea_watcher_${_currentPid}.bat',
          );

          // 创建监视脚本，增加检测间隔和重试次数
          watcherScript.writeAsStringSync('''
@echo off
:loop
REM 等待5秒
ping -n 6 127.0.0.1 > nul
REM 检查主进程是否存在
tasklist | find "${_currentPid}" > nul
if errorlevel 1 (
    REM 主进程不存在，再次确认
    ping -n 3 127.0.0.1 > nul
    tasklist | find "${_currentPid}" > nul
    if errorlevel 1 (
        REM 确认主进程确实不存在，终止子进程
        taskkill /F /PID ${_processId} > nul 2>&1
        del "${watcherScript.path}" > nul 2>&1
        exit
    )
)
goto loop
''');

          // 启动监视脚本
          _watcherProcess = await Process.start('cmd.exe', [
            '/c',
            'start',
            '/b',
            watcherScript.path,
          ]);
        } else {
          // Unix平台使用shell脚本监视父进程
          final tempDir = Directory.systemTemp;
          final watcherScript = File(
            '${tempDir.path}/lemon_tea_watcher_${_currentPid}.sh',
          );

          // 创建监视脚本，增加检测间隔和重试次数
          watcherScript.writeAsStringSync('''
#!/bin/sh
while :; do
    # 等待5秒
    sleep 5
    # 检查主进程是否存在
    if ! ps -p ${_currentPid} > /dev/null; then
        # 主进程不存在，再次确认
        sleep 2
        if ! ps -p ${_currentPid} > /dev/null; then
            # 确认主进程确实不存在，终止子进程
            kill -9 ${_processId} > /dev/null 2>&1
            rm "${watcherScript.path}" > /dev/null 2>&1
            exit 0
        fi
    fi
done
''');

          // 设置脚本执行权限
          Process.runSync('chmod', ['+x', watcherScript.path]);

          // 启动监视脚本
          _watcherProcess = await Process.start('/bin/sh', [
            watcherScript.path,
          ]);
        }
      } catch (e) {
        debugPrint('注册进程清理失败: $e');
      }
    }
  }

  /// 更新退出文件，记录子进程ID
  void _updateExitFile() {
    try {
      if (_processId != null) {
        final tempDir = Directory.systemTemp;
        final exitFile = File(
          '${tempDir.path}/lemon_tea_exit_${_currentPid}.pid',
        );
        exitFile.writeAsStringSync(_processId.toString());
      }
    } catch (e) {
      debugPrint('更新退出文件失败: $e');
    }
  }

  /// 检查服务是否可用
  Future<bool> _isServiceAvailable(int port) async {
    try {
      // 尝试连接服务
      final socket = await Socket.connect(
        'localhost',
        port,
        timeout: const Duration(seconds: 2),
      ).catchError((e) => null);

      if (socket != null) {
        await socket.close();
        return true;
      }

      return false;
    } catch (e) {
      return false;
    }
  }

  /// 清理进程资源
  void _cleanupProcess() {
    _process?.kill();
    _process = null;
    _watcherProcess?.kill();
    _watcherProcess = null;
    _stdoutSubscription?.cancel();
    _stdoutSubscription = null;
    _stderrSubscription?.cancel();
    _stderrSubscription = null;
    _isRunning = false;
    _port = null;
  }

  /// 停止CLI服务
  Future<bool> stopService() async {
    if (!_isRunning) {
      return false;
    }

    try {
      _forceKillProcess();
      return true;
    } catch (e) {
      debugPrint('停止CLI服务失败: $e');
      return false;
    }
  }
}
