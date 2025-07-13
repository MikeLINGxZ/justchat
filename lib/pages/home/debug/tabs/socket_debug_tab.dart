import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/socket/socket.dart';

class SocketDebugTab extends StatefulWidget {
  const SocketDebugTab({super.key});

  @override
  State<SocketDebugTab> createState() => _SocketDebugTabState();
}

class _SocketDebugTabState extends State<SocketDebugTab> {
  final TextEditingController _socketPathController = TextEditingController(
    text: '/var/folders/8z/vy1v761n3lbbdss2kp92b6d80000gn/T/lemontea.sock'
  );
  final TextEditingController _endpointController = TextEditingController(
    text: '/version'
  );
  final TextEditingController _outputController = TextEditingController();
  bool _isExecuting = false;

  @override
  void dispose() {
    _socketPathController.dispose();
    _endpointController.dispose();
    _outputController.dispose();
    super.dispose();
  }

  void _executeSocketRequest() async {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在执行请求...\n';
    });

    try {
      // 获取Socket路径和请求路径
      final socketPath = _socketPathController.text;
      final endpoint = _endpointController.text;
      
      // 创建Socket客户端
      final client = SocketHttpClient(socketPath);
      
      // 执行请求
      final response = await client.get(endpoint);
      
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行结果：\n'
            '请求路径: $endpoint\n'
            '时间: ${DateTime.now().toString()}\n'
            '状态: 成功\n'
            '响应内容: $response';
      });
    } catch (e) {
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行出错: $e';
      });
    }
  }

  void _testVersion() async {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在测试版本接口...\n';
    });

    try {
      final result = await testSocketConnection();
      
      setState(() {
        _isExecuting = false;
        if (result['success']) {
          _outputController.text += '测试成功：\n'
              '时间: ${DateTime.now().toString()}\n'
              '版本号: ${result['version']}';
        } else {
          _outputController.text += '测试失败：\n'
              '时间: ${DateTime.now().toString()}\n'
              '错误: ${result['error']}';
        }
      });
    } catch (e) {
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行出错: $e';
      });
    }
  }
  
  void _checkSocketFile() async {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在检查Socket文件状态...\n';
    });

    try {
      final socketPath = _socketPathController.text;
      final result = await checkSocketFile(socketPath);
      
      setState(() {
        _isExecuting = false;
        if (result['exists']) {
          final permissions = result['permissions'];
          _outputController.text += '文件存在：\n'
              '路径: $socketPath\n'
              '时间: ${DateTime.now().toString()}\n'
              '权限模式: ${permissions['mode']}\n'
              '修改时间: ${permissions['modified']}\n'
              '类型: ${permissions['type']}';
          
          if (permissions.containsKey('ls_output')) {
            _outputController.text += '\n详细权限: ${permissions['ls_output']}';
          }
        } else {
          _outputController.text += '文件不存在或无法访问：\n'
              '路径: $socketPath\n'
              '时间: ${DateTime.now().toString()}\n'
              '错误: ${result['error']}';
        }
      });
    } catch (e) {
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行出错: $e';
      });
    }
  }
  
  void _testWithCurl() async {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在使用curl测试Socket连接...\n';
    });

    try {
      final socketPath = _socketPathController.text;
      final endpoint = _endpointController.text;
      final result = await testSocketWithCurl(socketPath, endpoint);
      
      setState(() {
        _isExecuting = false;
        if (result['success']) {
          _outputController.text += '测试成功：\n'
              '路径: $socketPath\n'
              '请求: $endpoint\n'
              '时间: ${DateTime.now().toString()}\n'
              '退出码: ${result['exitCode']}\n'
              '标准输出: ${result['stdout']}\n'
              '标准错误: ${result['stderr']}';
        } else {
          _outputController.text += '测试失败：\n'
              '路径: $socketPath\n'
              '请求: $endpoint\n'
              '时间: ${DateTime.now().toString()}\n'
              '错误: ${result['error'] ?? "未知错误"}\n'
              '退出码: ${result['exitCode']}\n'
              '标准输出: ${result['stdout']}\n'
              '标准错误: ${result['stderr']}';
        }
      });
    } catch (e) {
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行出错: $e';
      });
    }
  }
  
  void _checkAvailableTools() async {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在检查可用的命令行工具...\n';
    });

    try {
      final tools = await checkAvailableTools();
      
      setState(() {
        _isExecuting = false;
        _outputController.text += '检查结果：\n'
            '时间: ${DateTime.now().toString()}\n\n'
            '可用工具：\n'
            'curl: ${tools['curl'] == true ? '✅ 可用' : '❌ 不可用'}\n'
            'socat: ${tools['socat'] == true ? '✅ 可用' : '❌ 不可用'}\n'
            'nc (netcat): ${tools['nc'] == true ? '✅ 可用' : '❌ 不可用'}\n'
            'python3: ${tools['python3'] == true ? '✅ 可用' : '❌ 不可用'}\n\n'
            '提示：如果curl可用，Socket连接问题可能是权限或沙盒限制导致的。';
      });
    } catch (e) {
      setState(() {
        _isExecuting = false;
        _outputController.text += '执行出错: $e';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Socket 调试',
            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            '测试Socket连接和HTTP请求',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Colors.grey[600],
            ),
          ),
          const SizedBox(height: 24),
          
          // 输入区域
          Card(
            elevation: 2,
            child: Padding(
              padding: const EdgeInsets.all(20.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.settings_ethernet, color: Colors.blue[600]),
                      const SizedBox(width: 8),
                      Text(
                        'Socket请求',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Socket路径输入框
                  TextField(
                    controller: _socketPathController,
                    decoration: InputDecoration(
                      labelText: 'Socket路径',
                      hintText: '请输入Unix Domain Socket路径',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                      prefixIcon: const Icon(Icons.folder),
                    ),
                  ),
                  const SizedBox(height: 16),
                  
                  // 请求路径输入框
                  TextField(
                    controller: _endpointController,
                    decoration: InputDecoration(
                      labelText: '请求路径',
                      hintText: '请输入HTTP请求路径',
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                      prefixIcon: const Icon(Icons.link),
                    ),
                  ),
                  const SizedBox(height: 16),
                  
                  // 按钮区域 - 第一行
                  Row(
                    children: [
                      // 执行请求按钮
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isExecuting ? null : _executeSocketRequest,
                          icon: _isExecuting 
                              ? Container(
                                  width: 24,
                                  height: 24,
                                  padding: const EdgeInsets.all(2.0),
                                  child: const CircularProgressIndicator(
                                    strokeWidth: 3,
                                  ),
                                )
                              : const Icon(Icons.send),
                          label: Text(_isExecuting ? '执行中...' : '发送请求'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.blue,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      // 测试版本按钮
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isExecuting ? null : _testVersion,
                          icon: _isExecuting 
                              ? Container(
                                  width: 24,
                                  height: 24,
                                  padding: const EdgeInsets.all(2.0),
                                  child: const CircularProgressIndicator(
                                    strokeWidth: 3,
                                  ),
                                )
                              : const Icon(Icons.verified),
                          label: Text(_isExecuting ? '测试中...' : '测试版本接口'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.green,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  
                  // 按钮区域 - 第二行
                  Row(
                    children: [
                      // 检查文件按钮
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isExecuting ? null : _checkSocketFile,
                          icon: _isExecuting 
                              ? Container(
                                  width: 24,
                                  height: 24,
                                  padding: const EdgeInsets.all(2.0),
                                  child: const CircularProgressIndicator(
                                    strokeWidth: 3,
                                  ),
                                )
                              : const Icon(Icons.find_in_page),
                          label: Text(_isExecuting ? '检查中...' : '检查Socket文件'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.orange,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      // 使用curl测试按钮
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isExecuting ? null : _testWithCurl,
                          icon: _isExecuting 
                              ? Container(
                                  width: 24,
                                  height: 24,
                                  padding: const EdgeInsets.all(2.0),
                                  child: const CircularProgressIndicator(
                                    strokeWidth: 3,
                                  ),
                                )
                              : const Icon(Icons.terminal),
                          label: Text(_isExecuting ? '测试中...' : '使用curl测试'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.purple,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  
                  // 按钮区域 - 第三行
                  Row(
                    children: [
                      // 检查可用工具按钮
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: _isExecuting ? null : _checkAvailableTools,
                          icon: _isExecuting 
                              ? Container(
                                  width: 24,
                                  height: 24,
                                  padding: const EdgeInsets.all(2.0),
                                  child: const CircularProgressIndicator(
                                    strokeWidth: 3,
                                  ),
                                )
                              : const Icon(Icons.build),
                          label: Text(_isExecuting ? '检查中...' : '检查可用工具'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.teal,
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 20),
          
          // 输出区域
          Card(
            elevation: 2,
            child: Padding(
              padding: const EdgeInsets.all(20.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.output, color: Colors.green[600]),
                      const SizedBox(width: 8),
                      Text(
                        '执行结果',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // 输出框
                  Container(
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey[300]!),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: TextField(
                      controller: _outputController,
                      decoration: const InputDecoration(
                        hintText: '执行结果将显示在这里',
                        contentPadding: EdgeInsets.all(16),
                        border: InputBorder.none,
                      ),
                      readOnly: true,
                      maxLines: 15,
                      style: TextStyle(
                        fontFamily: 'monospace',
                        fontSize: 13,
                        color: Colors.grey[800],
                      ),
                    ),
                  ),
                  
                  // 清除按钮
                  Align(
                    alignment: Alignment.centerRight,
                    child: Padding(
                      padding: const EdgeInsets.only(top: 8.0),
                      child: TextButton.icon(
                        onPressed: () {
                          setState(() {
                            _outputController.clear();
                          });
                        },
                        icon: const Icon(Icons.clear, size: 16),
                        label: const Text('清除输出'),
                        style: TextButton.styleFrom(
                          foregroundColor: Colors.grey[700],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
} 