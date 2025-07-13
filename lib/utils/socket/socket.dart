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
      print('尝试使用curl命令访问Socket');
      final result = await Process.run('curl', [
        '--unix-socket', 
        socketPath, 
        'http://localhost$path'
      ]);
      
      print('curl执行结果: exitCode=${result.exitCode}, stdout=${result.stdout}, stderr=${result.stderr}');
      
      if (result.exitCode == 0) {
        return result.stdout.toString().trim();
      } else {
        throw Exception('curl命令执行失败: ${result.stderr}');
      }
    } catch (e) {
      print('curl执行失败: $e');
      
      try {
        // 如果curl失败，尝试使用socat命令
        print('尝试使用socat命令访问Socket');
        
        // 准备HTTP请求内容
        final requestContent = 'GET $path HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n';
        
        // 创建临时文件存储请求内容
        final tempDir = Directory.systemTemp;
        final requestFile = File('${tempDir.path}/socket_request.txt');
        await requestFile.writeAsString(requestContent);
        
        // 使用socat命令
        // 注意：Process.run不支持直接传递stdin流，我们需要使用Process.start
        final process = await Process.start('socat', [
          '-t', '5',  // 5秒超时
          'UNIX-CONNECT:$socketPath',
          'STDIO'
        ]);
        
        // 写入请求内容
        process.stdin.write(requestContent);
        await process.stdin.flush();
        await process.stdin.close();
        
        // 读取输出
        final stdout = await process.stdout.transform(utf8.decoder).join();
        final stderr = await process.stderr.transform(utf8.decoder).join();
        final exitCode = await process.exitCode;
        
        // 删除临时文件
        await requestFile.delete();
        
        print('socat执行结果: exitCode=$exitCode, stdout=$stdout, stderr=$stderr');
        
        if (exitCode == 0) {
          final responseStr = stdout;
          final parts = responseStr.split('\r\n\r\n');
          if (parts.length >= 2) {
            return parts[1].trim();
          }
          return responseStr.trim();
        } else {
          throw Exception('socat命令执行失败: $stderr');
        }
      } catch (e) {
        print('socat执行失败: $e');
        
        // 尝试使用Python脚本
        try {
          print('尝试使用Python脚本访问Socket');
          
          // 创建临时Python脚本
          final tempDir = Directory.systemTemp;
          final pythonScript = File('${tempDir.path}/socket_client.py');
          
          await pythonScript.writeAsString('''
import socket
import sys

try:
    # 创建Unix domain socket
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    
    # 连接到socket
    sock.connect('$socketPath')
    
    # 发送HTTP请求
    request = 'GET $path HTTP/1.1\\r\\nHost: localhost\\r\\nConnection: close\\r\\n\\r\\n'
    sock.sendall(request.encode())
    
    # 接收响应
    response = b''
    while True:
        data = sock.recv(4096)
        if not data:
            break
        response += data
    
    # 关闭socket
    sock.close()
    
    # 输出响应
    print(response.decode())
    sys.exit(0)
except Exception as e:
    print(f"错误: {e}", file=sys.stderr)
    sys.exit(1)
''');
          
          // 执行Python脚本
          final result = await Process.run('python3', [pythonScript.path]);
          
          // 删除临时脚本
          await pythonScript.delete();
          
          print('Python脚本执行结果: exitCode=${result.exitCode}, stdout=${result.stdout}, stderr=${result.stderr}');
          
          if (result.exitCode == 0) {
            final responseStr = result.stdout.toString();
            final parts = responseStr.split('\r\n\r\n');
            if (parts.length >= 2) {
              return parts[1].trim();
            }
            return responseStr.trim();
          } else {
            throw Exception('Python脚本执行失败: ${result.stderr}');
          }
        } catch (e) {
          print('Python脚本执行失败: $e');
          throw Exception('所有命令行方法都失败: $e');
        }
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

/// 检查系统中可用的命令行工具
Future<Map<String, bool>> checkAvailableTools() async {
  final tools = <String, bool>{};
  
  // 检查curl
  try {
    final result = await Process.run('which', ['curl']);
    tools['curl'] = result.exitCode == 0;
  } catch (e) {
    tools['curl'] = false;
  }
  
  // 检查socat
  try {
    final result = await Process.run('which', ['socat']);
    tools['socat'] = result.exitCode == 0;
  } catch (e) {
    tools['socat'] = false;
  }
  
  // 检查nc
  try {
    final result = await Process.run('which', ['nc']);
    tools['nc'] = result.exitCode == 0;
  } catch (e) {
    tools['nc'] = false;
  }
  
  // 检查python3
  try {
    final result = await Process.run('which', ['python3']);
    tools['python3'] = result.exitCode == 0;
  } catch (e) {
    tools['python3'] = false;
  }
  
  return tools;
}
