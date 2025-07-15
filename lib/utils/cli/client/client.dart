import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:grpc/grpc.dart';
import 'package:lemon_tea/rpc/service.pb.dart';
import 'package:lemon_tea/rpc/service.pbgrpc.dart';
import 'package:lemon_tea/utils/cli/server/server.dart' as cliServer;

/// gRPC客户端类，用于与本地CLI服务通信
class Client {
  /// 单例实例
  static final Client _instance = Client._internal();

  /// 工厂构造函数
  factory Client() => _instance;

  /// 内部构造函数
  Client._internal() {
    // 监听服务器端口变化
    cliServer.Server().portStream.listen((port) {
      if (port != null && port != _port) {
        _port = port;
        debugPrint('服务端口变化: $port');
        _reconnect();
      } else if (port == null) {
        _cleanupClient();
      }
    });
  }

  /// 服务端口
  int? _port;

  /// gRPC通道
  ClientChannel? _channel;

  /// gRPC客户端存根
  LemonTeaClient? _stub;

  /// 获取gRPC客户端存根
  LemonTeaClient? get stub => _stub;

  /// 初始化客户端
  ///
  /// [port] 服务端口，如果为null则使用Server中的端口
  /// 返回是否初始化成功
  Future<bool> init({int? port}) async {
    try {
      // 如果已经初始化且端口未变，直接返回成功
      if (_stub != null && _port != null && (port == null || port == _port)) {
        return true;
      }

      // 清理现有客户端
      _cleanupClient();

      // 获取端口
      _port = port ?? cliServer.Server().port;
      if (_port == null) {
        debugPrint('无法获取服务端口');
        return false;
      }

      // 创建gRPC通道
      _channel = ClientChannel(
        '127.0.0.1',
        port: _port!,
        options: const ChannelOptions(
          credentials: ChannelCredentials.insecure(),
          idleTimeout: Duration(minutes: 1),
        ),
      );

      // 创建gRPC客户端存根
      _stub = LemonTeaClient(_channel!);

      return true;
    } catch (e) {
      debugPrint('初始化gRPC客户端失败: $e');
      _cleanupClient();
      return false;
    }
  }

  /// 重新连接
  Future<bool> _reconnect() async {
    return await init(port: _port);
  }

  /// 清理客户端资源
  void _cleanupClient() {
    _channel?.shutdown();
    _channel = null;
    _stub = null;
  }

  /// 关闭客户端
  Future<void> close() async {
    _cleanupClient();
  }
}
