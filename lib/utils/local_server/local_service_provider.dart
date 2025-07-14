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
  
  /// 上次使用的端口
  int? _lastPort;
  
  /// 是否禁用自动重启（手动操作时设置为true）
  bool _disableAutoRestart = false;
  
  /// 禁用自动重启的计时器
  Timer? _disableAutoRestartTimer;
  
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
    
    // 保存上次使用的端口
    if (port != null) {
      _lastPort = port;
    }
    
    // 如果服务异常停止，且未禁用自动重启，则尝试重启
    if (previousState.isRunning && !isRunning && !_disableAutoRestart) {
      debugPrint('检测到CLI服务异常退出，尝试重启...');
      _restartService();
    } else if (previousState.isRunning && !isRunning) {
      debugPrint('CLI服务已停止，由于是手动操作，不会自动重启');
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
    // 如果状态显示服务应该运行，但实际上没有运行，且未禁用自动重启，则尝试重启
    if (state.isRunning && !_cliService.isRunning && !_disableAutoRestart) {
      debugPrint('定期检查：检测到CLI服务异常退出，尝试重启...');
      await _restartService();
    } else if (state.isRunning && !_cliService.isRunning) {
      debugPrint('定期检查：CLI服务已停止，由于是手动操作，不会自动重启');
      // 更新状态以反映实际情况
      state = state.copyWith(isRunning: false, port: null);
    }
  }
  
  /// 临时禁用自动重启（用于手动操作）
  void _temporarilyDisableAutoRestart() {
    _disableAutoRestart = true;
    
    // 取消现有的禁用计时器
    _disableAutoRestartTimer?.cancel();
    
    // 设置新的计时器，30秒后重新启用自动重启
    _disableAutoRestartTimer = Timer(const Duration(seconds: 30), () {
      _disableAutoRestart = false;
      debugPrint('自动重启功能已重新启用');
    });
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
      // 尝试使用上次的端口重启
      final port = await startService(requestedPort: _lastPort, isManualOperation: false);
      
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
  /// 
  /// [requestedPort] 请求使用的端口号，如果为null则自动分配
  /// [isManualOperation] 是否是手动操作，默认为true
  /// 返回服务端口号，如果启动失败则返回null
  Future<int?> startService({int? requestedPort, bool isManualOperation = true}) async {
    if (state.isRunning) {
      return state.port;
    }
    
    // 如果是手动操作，临时禁用自动重启
    if (isManualOperation) {
      _temporarilyDisableAutoRestart();
    }
    
    final port = await _cliService.startService(requestedPort: requestedPort);
    
    if (port != null) {
      state = state.copyWith(isRunning: true, port: port);
      _lastPort = port;
      _resetRestartAttempts();
    }
    
    return port;
  }
  
  /// 停止CLI服务
  /// 
  /// [isManualOperation] 是否是手动操作，默认为true
  Future<void> stopService({bool isManualOperation = true}) async {
    _restartTimer?.cancel();
    _restartTimer = null;
    _restartAttempts = 0;
    
    // 如果是手动操作，临时禁用自动重启
    if (isManualOperation) {
      _temporarilyDisableAutoRestart();
    }
    
    await _cliService.stopService();
    state = state.copyWith(isRunning: false, port: null);
  }
  
  /// 重启CLI服务（手动操作）
  Future<int?> restartService({int? requestedPort}) async {
    // 临时禁用自动重启
    _temporarilyDisableAutoRestart();
    
    // 停止服务
    await stopService(isManualOperation: false);
    
    // 启动服务
    return await startService(requestedPort: requestedPort, isManualOperation: false);
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
    _disableAutoRestartTimer?.cancel();
    _cliService.removeStatusListener(_onServiceStatusChanged);
    super.dispose();
  }
}

/// CLI服务状态Provider
final cliServiceProvider = StateNotifierProvider<CliServiceNotifier, CliServiceState>((ref) {
  return CliServiceNotifier();
}); 