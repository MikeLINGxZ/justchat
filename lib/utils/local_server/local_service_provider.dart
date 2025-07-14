import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/local_server/local_service.dart';
import 'dart:async';
import 'package:flutter/foundation.dart';

/// CLI服务状态
class CliServiceState {
  /// 服务是否正在运行
  final bool isRunning;
  
  /// 服务端口
  final int? port;
  
  /// 构造函数
  CliServiceState({
    required this.isRunning,
    this.port,
  });
  
  /// 创建初始状态
  factory CliServiceState.initial() {
    return CliServiceState(
      isRunning: false,
      port: null,
    );
  }
  
  /// 复制并修改状态
  CliServiceState copyWith({
    bool? isRunning,
    int? port,
  }) {
    return CliServiceState(
      isRunning: isRunning ?? this.isRunning,
      port: port ?? this.port,
    );
  }
}

/// CLI服务状态管理类
class CliServiceNotifier extends StateNotifier<CliServiceState> {
  /// CLI服务实例
  final CliService _cliService = CliService();
  
  /// 自动重启计时器
  Timer? _restartTimer;
  
  /// 重启尝试次数
  int _restartAttempts = 0;
  
  /// 最大重启尝试次数
  static const int _maxRestartAttempts = 5;
  
  /// 构造函数
  CliServiceNotifier() : super(CliServiceState.initial()) {
    // 设置CLI服务状态监听器
    _cliService.addStatusListener(_onServiceStatusChanged);
    
    // 设置定期检查服务状态的监听器（作为备份机制）
    _setupServiceMonitor();
  }
  
  /// 处理服务状态变化
  void _onServiceStatusChanged(bool isRunning, int? port) {
    debugPrint('CLI服务状态变化: isRunning=$isRunning, port=$port');
    
    // 更新状态
    final previousState = state;
    state = state.copyWith(isRunning: isRunning, port: port);
    
    // 如果服务异常停止，尝试重启
    if (previousState.isRunning && !isRunning) {
      debugPrint('检测到CLI服务异常退出，尝试重启...');
      _restartService();
    }
  }
  
  /// 设置服务监控（作为备份检测机制）
  void _setupServiceMonitor() {
    // 创建一个定时器，定期检查服务状态
    Timer.periodic(const Duration(seconds: 30), (timer) async {
      await _checkServiceStatus();
    });
  }
  
  /// 检查服务状态并在需要时重启
  Future<void> _checkServiceStatus() async {
    // 如果状态显示服务应该运行，但实际上没有运行，则尝试重启
    if (state.isRunning && !_cliService.isRunning) {
      debugPrint('定期检查：检测到CLI服务异常退出，尝试重启...');
      await _restartService();
    }
  }
  
  /// 重启服务的实现
  Future<void> _restartService() async {
    // 如果已经达到最大重启次数，则不再尝试
    if (_restartAttempts >= _maxRestartAttempts) {
      debugPrint('CLI服务重启失败次数过多，不再尝试重启');
      state = state.copyWith(isRunning: false, port: null);
      _resetRestartAttempts();
      return;
    }
    
    _restartAttempts++;
    
    // 计算退避时间（指数退避策略）
    final backoffSeconds = _calculateBackoffTime(_restartAttempts);
    debugPrint('将在 $backoffSeconds 秒后尝试重启CLI服务 (尝试 $_restartAttempts/$_maxRestartAttempts)');
    
    // 取消之前的重启计时器（如果有）
    _restartTimer?.cancel();
    
    // 设置新的重启计时器
    _restartTimer = Timer(Duration(seconds: backoffSeconds), () async {
      debugPrint('正在尝试重启CLI服务...');
      final port = await startService();
      
      if (port != null) {
        debugPrint('CLI服务重启成功，端口: $port');
        _resetRestartAttempts();
      } else {
        debugPrint('CLI服务重启失败，将再次尝试');
        // 失败后会由监控器再次触发重启
      }
    });
  }
  
  /// 计算退避时间（指数退避策略）
  int _calculateBackoffTime(int attempt) {
    // 基础等待时间为2秒，每次失败后翻倍，最大等待30秒
    return (2 << (attempt - 1)).clamp(2, 30);
  }
  
  /// 重置重启尝试计数
  void _resetRestartAttempts() {
    _restartAttempts = 0;
    _restartTimer?.cancel();
    _restartTimer = null;
  }
  
  /// 启动CLI服务
  Future<int?> startService() async {
    if (state.isRunning) {
      return state.port;
    }
    
    final port = await _cliService.startService();
    
    if (port != null) {
      state = state.copyWith(isRunning: true, port: port);
      _resetRestartAttempts();
    }
    
    return port;
  }
  
  /// 停止CLI服务
  Future<void> stopService() async {
    _restartTimer?.cancel();
    _restartTimer = null;
    _restartAttempts = 0;
    
    await _cliService.stopService();
    state = state.copyWith(isRunning: false, port: null);
  }
  
  /// 获取CLI服务状态
  Future<void> refreshState() async {
    state = state.copyWith(
      isRunning: _cliService.isRunning,
      port: _cliService.port,
    );
  }
  
  @override
  void dispose() {
    _restartTimer?.cancel();
    _cliService.removeStatusListener(_onServiceStatusChanged);
    super.dispose();
  }
}

/// CLI服务状态Provider
final cliServiceProvider = StateNotifierProvider<CliServiceNotifier, CliServiceState>((ref) {
  return CliServiceNotifier();
}); 