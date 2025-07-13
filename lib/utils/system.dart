import 'dart:io';
import 'dart:math';

class System {
  static bool get isDesktop {
    return Platform.isWindows || Platform.isMacOS || Platform.isLinux;
  }
  static bool get isMobile {
    return Platform.isAndroid || Platform.isIOS;
  }

  static bool get isAndroid {
    return Platform.isAndroid;
  }

  static bool get isIOS {
    return Platform.isIOS;
  }

  static bool get isWindows {
    return Platform.isWindows;
  }

  static bool get isMacOS {
    return Platform.isMacOS;
  }

  static bool get isLinux {
    return Platform.isLinux;
  }

  /// 获取一个空闲的端口号
  ///
  /// [startPort] 开始搜索的端口号，默认为8000
  /// [endPort] 结束搜索的端口号，默认为9000
  /// 返回一个空闲的端口号，如果没有找到则返回null
  static Future<int?> findFreePort({int startPort = 8000, int endPort = 9000}) async {
    if (startPort < 0 || endPort > 65535 || startPort > endPort) {
      throw ArgumentError('端口范围无效: $startPort-$endPort');
    }

    // 随机化起始端口，避免每次都从同一个端口开始
    final random = Random();
    int port = startPort + random.nextInt(endPort - startPort);

    // 尝试在端口范围内找到一个空闲端口
    final int maxAttempts = endPort - startPort;
    int attempts = 0;

    while (attempts < maxAttempts) {
      try {
        // 尝试绑定端口来检查它是否可用
        final serverSocket = await ServerSocket.bind(InternetAddress.loopbackIPv4, port, shared: true);
        await serverSocket.close();
        return port;
      } catch (e) {
        // 端口不可用，尝试下一个
        port = startPort + ((port - startPort + 1) % (endPort - startPort));
        attempts++;
      }
    }

    // 如果所有端口都被占用，则返回null
    return null;
  }
}