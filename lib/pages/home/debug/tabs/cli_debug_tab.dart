import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/local_server/local_service_provider.dart';
import 'package:lemon_tea/utils/system.dart';

class CliDebugTab extends ConsumerStatefulWidget {
  const CliDebugTab({super.key});

  @override
  ConsumerState<CliDebugTab> createState() => _CliDebugTabState();
}

class _CliDebugTabState extends ConsumerState<CliDebugTab> {
  final TextEditingController _portController = TextEditingController();
  final TextEditingController _logController = TextEditingController();
  bool _isChangingPort = false;
  bool _isRestarting = false;
  bool _isStopping = false;

  @override
  void initState() {
    super.initState();
    _updatePortController();
    _addLogMessage('CLI调试页面已初始化');
  }

  @override
  void dispose() {
    _portController.dispose();
    _logController.dispose();
    super.dispose();
  }

  void _updatePortController() {
    final cliState = ref.read(cliServiceProvider);
    _portController.text = cliState.port?.toString() ?? '未启动';
  }

  void _addLogMessage(String message) {
    final timestamp = DateTime.now().toString().split('.')[0];
    final logEntry = '[$timestamp] $message\n';
    
    setState(() {
      _logController.text = '${_logController.text}$logEntry';
    });
  }

  Future<void> _restartCliService() async {
    if (_isRestarting) return;

    setState(() {
      _isRestarting = true;
    });

    _addLogMessage('正在重启CLI服务...');

    try {
      // 使用新的restartService方法
      final port = await ref.read(cliServiceProvider.notifier).restartService();
      
      if (port != null) {
        _addLogMessage('CLI服务已重启，端口: $port');
        _updatePortController();
      } else {
        _addLogMessage('CLI服务重启失败');
      }
    } catch (e) {
      _addLogMessage('重启CLI服务时发生错误: $e');
    } finally {
      setState(() {
        _isRestarting = false;
      });
    }
  }

  Future<void> _stopCliService() async {
    if (_isStopping) return;

    setState(() {
      _isStopping = true;
    });

    _addLogMessage('正在停止CLI服务...');

    try {
      await ref.read(cliServiceProvider.notifier).stopService();
      _addLogMessage('CLI服务已停止');
      _updatePortController();
    } catch (e) {
      _addLogMessage('停止CLI服务时发生错误: $e');
    } finally {
      setState(() {
        _isStopping = false;
      });
    }
  }

  Future<void> _changePort() async {
    if (_isChangingPort) return;

    final portText = _portController.text.trim();
    if (portText.isEmpty) {
      _addLogMessage('错误：端口不能为空');
      return;
    }

    int? newPort;
    try {
      newPort = int.parse(portText);
    } catch (e) {
      _addLogMessage('错误：无效的端口号');
      return;
    }

    if (newPort <= 0 || newPort > 65535) {
      _addLogMessage('错误：端口号必须在1-65535之间');
      return;
    }

    setState(() {
      _isChangingPort = true;
    });

    _addLogMessage('正在更改CLI服务端口到: $newPort...');

    try {
      // 检查端口是否可用
      final isPortAvailable = await System.isPortAvailable(newPort);
      if (!isPortAvailable) {
        _addLogMessage('错误：端口 $newPort 已被占用');
        setState(() {
          _isChangingPort = false;
        });
        return;
      }
      
      // 使用新端口重启服务
      final port = await ref.read(cliServiceProvider.notifier).restartService(requestedPort: newPort);
      if (port != null) {
        _addLogMessage('CLI服务已启动，端口: $port');
        _updatePortController();
      } else {
        _addLogMessage('CLI服务启动失败');
      }
    } catch (e) {
      _addLogMessage('更改端口时发生错误: $e');
    } finally {
      setState(() {
        _isChangingPort = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final cliState = ref.watch(cliServiceProvider);
    final isRunning = cliState.isRunning;
    final currentPort = cliState.port;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'CLI 服务调试',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '管理和监控本地CLI服务',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 24),
        
        // 状态区域
        Card(
          elevation: 2,
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      isRunning ? Icons.check_circle : Icons.error,
                      color: isRunning ? Colors.green[600] : Colors.red[600],
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '服务状态',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                
                // 状态信息
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: isRunning ? Colors.green[50] : Colors.red[50],
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color: isRunning ? Colors.green[300]! : Colors.red[300]!,
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        isRunning ? Icons.play_arrow : Icons.stop,
                        color: isRunning ? Colors.green[700] : Colors.red[700],
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              isRunning ? '服务运行中' : '服务已停止',
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: isRunning ? Colors.green[700] : Colors.red[700],
                              ),
                            ),
                            if (isRunning && currentPort != null)
                              Text(
                                '端口: $currentPort',
                                style: TextStyle(
                                  color: Colors.grey[800],
                                ),
                              ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),
                
                // 端口设置
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _portController,
                        decoration: InputDecoration(
                          labelText: '端口',
                          hintText: '输入端口号',
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                          prefixIcon: const Icon(Icons.settings_ethernet),
                        ),
                        keyboardType: TextInputType.number,
                        enabled: !_isChangingPort && !_isRestarting && !_isStopping,
                      ),
                    ),
                    const SizedBox(width: 12),
                    ElevatedButton(
                      onPressed: isRunning && !_isChangingPort && !_isRestarting && !_isStopping
                          ? _changePort
                          : null,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.blue,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: _isChangingPort
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text('应用'),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                
                // 控制按钮
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Expanded(
                      child: ElevatedButton.icon(
                        onPressed: !_isRestarting && !_isChangingPort && !_isStopping
                            ? _restartCliService
                            : null,
                        icon: _isRestarting
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                  color: Colors.white,
                                  strokeWidth: 2,
                                ),
                              )
                            : const Icon(Icons.refresh),
                        label: const Text('重启服务'),
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
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton.icon(
                        onPressed: isRunning && !_isStopping && !_isChangingPort && !_isRestarting
                            ? _stopCliService
                            : null,
                        icon: _isStopping
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                  color: Colors.white,
                                  strokeWidth: 2,
                                ),
                              )
                            : const Icon(Icons.stop),
                        label: const Text('停止服务'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.red,
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
        
        // 日志区域
        Expanded(
          child: Card(
            elevation: 2,
            child: Padding(
              padding: const EdgeInsets.all(20.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.article, color: Colors.indigo[600]),
                      const SizedBox(width: 8),
                      Text(
                        '操作日志',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // 日志输出
                  Expanded(
                    child: Container(
                      decoration: BoxDecoration(
                        border: Border.all(color: Colors.grey[300]!),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: TextField(
                        controller: _logController,
                        decoration: const InputDecoration(
                          hintText: '日志将显示在这里',
                          contentPadding: EdgeInsets.all(16),
                          border: InputBorder.none,
                        ),
                        readOnly: true,
                        maxLines: null,
                        expands: true,
                        style: TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 13,
                          color: Colors.grey[800],
                        ),
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
                            _logController.clear();
                            _addLogMessage('日志已清除');
                          });
                        },
                        icon: const Icon(Icons.clear, size: 16),
                        label: const Text('清除日志'),
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
        ),
      ],
    );
  }
} 