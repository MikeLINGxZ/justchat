import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/system.dart';
import 'package:lemon_tea/utils/local_server/local_server.dart';

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
  bool _isRunning = false;
  int? _currentPort;
  
  // 创建CLI服务实例
  final CliService _cliService = CliService();

  @override
  void initState() {
    super.initState();
    _checkServiceStatus();
    _updatePortController();
    _addLogMessage('CLI调试页面已初始化');
    
    // 添加端口控制器监听器，用于实时更新应用按钮状态
    _portController.addListener(_onPortTextChanged);
  }

  @override
  void dispose() {
    _portController.removeListener(_onPortTextChanged);
    _portController.dispose();
    _logController.dispose();
    super.dispose();
  }
  
  // 端口文本变化监听
  void _onPortTextChanged() {
    // 触发重建以更新应用按钮状态
    setState(() {});
  }
  
  // 检查当前输入的端口是否与服务端口相同
  bool get _isPortUnchanged {
    if (_currentPort == null) return false;
    
    try {
      final inputPort = int.parse(_portController.text.trim());
      return inputPort == _currentPort;
    } catch (e) {
      return false;
    }
  }
  
  // 检查服务状态
  Future<void> _checkServiceStatus() async {
    setState(() {
      _isRunning = _cliService.isRunning;
      _currentPort = _cliService.port;
    });
  }

  void _updatePortController() {
    if (_currentPort != null) {
      _portController.text = _currentPort.toString();
    } else {
      _portController.text = '';
    }
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
      // 先停止服务
      await _cliService.stopService();
      
      // 然后启动服务，使用当前端口
      final port = await _cliService.startService(
        requestedPort: _currentPort,
      );
      
      if (port != null) {
        _addLogMessage('CLI服务重启成功，端口: $port');
        setState(() {
          _isRunning = true;
          _currentPort = port;
        });
        _updatePortController();
      } else {
        _addLogMessage('CLI服务重启失败');
        setState(() {
          _isRunning = false;
        });
      }
    } catch (e) {
      _addLogMessage('重启过程中发生错误: $e');
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
      final result = await _cliService.stopService();
      
      if (result) {
        _addLogMessage('CLI服务已停止');
        setState(() {
          _isRunning = false;
        });
      } else {
        _addLogMessage('停止CLI服务失败');
      }
    } catch (e) {
      _addLogMessage('停止过程中发生错误: $e');
    } finally {
      setState(() {
        _isStopping = false;
      });
    }
  }

  Future<void> _changePort() async {
    if (_isChangingPort) return;
    
    final String portText = _portController.text.trim();
    if (portText.isEmpty) {
      _addLogMessage('请输入有效的端口号');
      return;
    }
    
    int? newPort;
    try {
      newPort = int.parse(portText);
    } catch (e) {
      _addLogMessage('端口号格式无效');
      return;
    }
    
    if (newPort <= 0 || newPort > 65535) {
      _addLogMessage('端口号必须在1-65535之间');
      return;
    }
    
    // 如果端口未变，则不需要操作
    if (_currentPort != null && newPort == _currentPort) {
      _addLogMessage('端口未变更，无需操作');
      return;
    }
    
    setState(() {
      _isChangingPort = true;
    });
    
    _addLogMessage('正在更改端口到 $newPort...');
    
    try {
      // 检查端口是否可用
      final isAvailable = await System.isPortAvailable(newPort);
      if (!isAvailable) {
        _addLogMessage('端口 $newPort 已被占用，请选择其他端口');
        setState(() {
          _isChangingPort = false;
        });
        return;
      }
      
      // 先停止服务
      await _cliService.stopService();
      
      // 然后使用新端口启动服务
      final port = await _cliService.startService(
        requestedPort: newPort,
      );
      
      if (port != null) {
        _addLogMessage('端口更改成功，新端口: $port');
        setState(() {
          _isRunning = true;
          _currentPort = port;
        });
        _updatePortController();
      } else {
        _addLogMessage('端口更改失败');
        setState(() {
          _isRunning = false;
          _currentPort = null;
        });
      }
    } catch (e) {
      _addLogMessage('更改端口过程中发生错误: $e');
    } finally {
      setState(() {
        _isChangingPort = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
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
                      _isRunning ? Icons.check_circle : Icons.error,
                      color: _isRunning ? Colors.green[600] : Colors.red[600],
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
                    color: _isRunning ? Colors.green[50] : Colors.red[50],
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color: _isRunning ? Colors.green[300]! : Colors.red[300]!,
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        _isRunning ? Icons.play_arrow : Icons.stop,
                        color: _isRunning ? Colors.green[700] : Colors.red[700],
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              _isRunning ? '服务运行中' : '服务已停止',
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: _isRunning ? Colors.green[700] : Colors.red[700],
                              ),
                            ),
                            if (_isRunning && _currentPort != null)
                              Text(
                                '端口: $_currentPort',
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
                      onPressed: _isRunning && !_isChangingPort && !_isRestarting && !_isStopping && !_isPortUnchanged
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
                        onPressed: _isRunning && !_isStopping && !_isChangingPort && !_isRestarting
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