import 'dart:io';
import 'dart:convert';
import 'dart:typed_data';
import 'dart:async';

/// 用于通过Unix Domain Socket调用HTTP服务的客户端
class SocketHttpClient {
  final String socketPath;
  
  /// 构造函数
  /// [socketPath] Unix Domain Socket路径
  SocketHttpClient(this.socketPath);
  
  /// 发送GET请求到指定路径
  /// [path] 请求路径
  /// 返回响应内容字符串
  Future<String> get(String path) async {
    Socket? socket;
    try {
      // 尝试使用不同方式连接到Unix Domain Socket
      try {
        // 方法1: 直接使用Socket.connect
        socket = await Socket.connect(
          InternetAddress(socketPath, type: InternetAddressType.unix),
          0, // 端口号对Unix Domain Socket无意义
          timeout: const Duration(seconds: 5),
        );
      } catch (e) {
        print('方法1连接失败: $e');
        
        try {
          // 方法2: 使用RawSocket
          final rawSocket = await RawSocket.connect(
            InternetAddress(socketPath, type: InternetAddressType.unix),
            0,
            timeout: const Duration(seconds: 5),
          );
          
          // 创建一个新的Socket
          socket = await Socket.connect(
            InternetAddress(socketPath, type: InternetAddressType.unix),
            0,
          );
        } catch (e) {
          print('方法2连接失败: $e');
          
          // 方法3: 使用命令行工具 curl 或 nc
          return await _getViaCommandLine(path);
        }
      }
      
      if (socket == null) {
        throw Exception('无法创建Socket连接');
      }
      
      // 构建HTTP GET请求
      final request = 'GET $path HTTP/1.1\r\n'
          'Host: localhost\r\n'
          'Connection: close\r\n\r\n';
      
      // 发送请求
      socket.write(request);
      
      // 接收响应
      final response = StringBuffer();
      final completer = Completer<String>();
      
      // 设置超时
      Future.delayed(const Duration(seconds: 10), () {
        if (!completer.isCompleted) {
          completer.completeError('请求超时');
          socket?.destroy();
        }
      });
      
      // 处理数据
      socket.listen(
        (Uint8List data) {
          response.write(utf8.decode(data));
        },
        onDone: () {
          if (!completer.isCompleted) {
            // 解析HTTP响应
            final responseStr = response.toString();
            final parts = responseStr.split('\r\n\r\n');
            if (parts.length >= 2) {
              // 返回响应体
              completer.complete(parts[1].trim());
            } else {
              completer.complete(responseStr.trim());
            }
          }
          socket?.destroy();
        },
        onError: (error) {
          if (!completer.isCompleted) {
            completer.completeError('Socket错误: $error');
          }
          socket?.destroy();
        },
        cancelOnError: true,
      );
      
      return await completer.future;
    } catch (e) {
      socket?.destroy();
      throw Exception('Socket请求失败: $e');
    }
  }
  
  /// 通过命令行工具访问Socket
  Future<String> _getViaCommandLine(String path) async {
    try {
      // 尝试使用curl命令
      final result = await Process.run('curl', [
        '--unix-socket', 
        socketPath, 
        'http://localhost$path'
      ]);
      
      if (result.exitCode == 0) {
        return result.stdout.toString().trim();
      } else {
        throw Exception('命令执行失败: ${result.stderr}');
      }
    } catch (e) {
      try {
        // 如果curl失败，尝试使用nc (netcat)命令
        final tempFile = File('${Directory.systemTemp.path}/socket_request.txt');
        await tempFile.writeAsString(
          'GET $path HTTP/1.1\r\n'
          'Host: localhost\r\n'
          'Connection: close\r\n\r\n'
        );
        
        final result = await Process.run('nc', [
          '-U', 
          socketPath,
          '<',
          tempFile.path
        ]);
        
        await tempFile.delete();
        
        if (result.exitCode == 0) {
          final responseStr = result.stdout.toString();
          final parts = responseStr.split('\r\n\r\n');
          if (parts.length >= 2) {
            return parts[1].trim();
          }
          return responseStr.trim();
        } else {
          throw Exception('nc命令执行失败: ${result.stderr}');
        }
      } catch (e) {
        throw Exception('命令行访问Socket失败: $e');
      }
    }
  }
  
  /// 获取版本号
  Future<String> getVersion() async {
    return await get('/version');
  }
}

/// 使用示例
Future<String> getVersionFromSocket() async {
  final client = SocketHttpClient('/var/folders/8z/vy1v761n3lbbdss2kp92b6d80000gn/T/lemontea.sock');
  try {
    return await client.getVersion();
  } catch (e) {
    print('获取版本失败: $e');
    return '未知版本';
  }
}

/// 测试Socket连接并获取版本号
/// 返回一个包含版本号或错误信息的Future
Future<Map<String, dynamic>> testSocketConnection() async {
  try {
    final version = await getVersionFromSocket();
    return {
      'success': true,
      'version': version,
      'error': null,
    };
  } catch (e) {
    return {
      'success': false,
      'version': null,
      'error': e.toString(),
    };
  }
}

/// 检查Socket文件是否存在并可访问
Future<Map<String, dynamic>> checkSocketFile(String path) async {
  try {
    final file = File(path);
    final exists = await file.exists();
    
    if (!exists) {
      return {
        'exists': false,
        'error': '文件不存在',
        'permissions': null,
      };
    }
    
    final stat = await file.stat();
    
    // 尝试使用ls -la命令获取更详细的权限信息
    try {
      final result = await Process.run('ls', ['-la', path]);
      final lsOutput = result.stdout.toString().trim();
      
      return {
        'exists': true,
        'error': null,
        'permissions': {
          'mode': stat.mode,
          'modified': stat.modified,
          'type': stat.type,
          'ls_output': lsOutput,
        }
      };
    } catch (e) {
      return {
        'exists': true,
        'error': null,
        'permissions': {
          'mode': stat.mode,
          'modified': stat.modified,
          'type': stat.type,
        }
      };
    }
  } catch (e) {
    return {
      'exists': false,
      'error': e.toString(),
      'permissions': null,
    };
  }
}

/// 直接使用curl命令测试Socket连接
Future<Map<String, dynamic>> testSocketWithCurl(String socketPath, String path) async {
  try {
    // 执行curl命令
    final result = await Process.run('curl', [
      '--unix-socket', 
      socketPath, 
      'http://localhost$path',
      '-v'  // 添加详细输出
    ]);
    
    return {
      'success': result.exitCode == 0,
      'stdout': result.stdout.toString(),
      'stderr': result.stderr.toString(),
      'exitCode': result.exitCode,
    };
  } catch (e) {
    return {
      'success': false,
      'error': e.toString(),
      'stdout': '',
      'stderr': '',
      'exitCode': -1,
    };
  }
}
