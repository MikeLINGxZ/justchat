import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/local_server/local_service.dart';

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
  
  /// 构造函数
  CliServiceNotifier() : super(CliServiceState.initial());
  
  /// 启动CLI服务
  Future<int?> startService() async {
    if (state.isRunning) {
      return state.port;
    }
    
    final port = await _cliService.startService();
    
    if (port != null) {
      state = state.copyWith(isRunning: true, port: port);
    }
    
    return port;
  }
  
  /// 停止CLI服务
  Future<void> stopService() async {
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
}

/// CLI服务状态Provider
final cliServiceProvider = StateNotifierProvider<CliServiceNotifier, CliServiceState>((ref) {
  return CliServiceNotifier();
}); 