import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/system.dart';
import 'package:lemon_tea/utils/server/server.dart';
import 'package:lemon_tea/utils/storage/local_storage.dart';
import 'dart:io';
import 'package:file_picker/file_picker.dart';

class ServerDebugTab extends ConsumerStatefulWidget {
  const ServerDebugTab({super.key});

  @override
  ConsumerState<ServerDebugTab> createState() => _ServerDebugTabState();
}

class _ServerDebugTabState extends ConsumerState<ServerDebugTab> {
  final TextEditingController _portController = TextEditingController();
  final TextEditingController _logController = TextEditingController();
  final TextEditingController _binaryPathController = TextEditingController();
  bool _isChangingPort = false;
  bool _isRestarting = false;
  bool _isStopping = false;
  bool _isRunning = false;
  int? _currentPort;
  String? _customBinaryPath;
  bool _isDebugMode = false;
  
  // 创建SERVER服务实例
  final Server _server = Server();
  
  // 本地存储实例
  final LocalStorage _localStorage = LocalStorage();
  
  // 存储键名
  static const String _portKey = 'server_debug_port';
  static const String _binaryPathKey = 'server_debug_binary_path';

  @override
  void initState() {
    super.initState();
    _isDebugMode = _checkIsDebugMode();
    _loadSavedSettings();
    _checkServiceStatus();
    _addLogMessage('SERVER调试页面已初始化');
    
    // 添加端口控制器监听器，用于实时更新应用按钮状态
    _portController.addListener(_onPortTextChanged);
  }

  @override
  void dispose() {
    _portController.removeListener(_onPortTextChanged);
    _portController.dispose();
    _logController.dispose();
    _binaryPathController.dispose();
    super.dispose();
  }
  
  // 检查是否为调试模式
  bool _checkIsDebugMode() {
    bool inDebugMode = false;
    assert(() {
      inDebugMode = true;
      return true;
    }());
    return inDebugMode;
  }
  
  // 加载保存的设置
  Future<void> _loadSavedSettings() async {
    try {
      // 加载保存的端口
      final savedPort = await _localStorage.getInt(_portKey);
      if (savedPort != null) {
        setState(() {
          _currentPort = savedPort;
        });
        _portController.text = savedPort.toString();
      }
      
      // 加载保存的二进制路径
      final savedPath = await _localStorage.getString(_binaryPathKey);
      if (savedPath != null && savedPath.isNotEmpty) {
        setState(() {
          _customBinaryPath = savedPath;
        });
        _binaryPathController.text = savedPath;
      }
    } catch (e) {
      _addLogMessage('加载保存的设置失败: $e');
    }
  }
  
  // 保存设置
  Future<void> _saveSettings() async {
    try {
      // 保存端口
      if (_currentPort != null) {
        await _localStorage.setInt(_portKey, _currentPort!);
      }
      
      // 保存二进制路径
      if (_customBinaryPath != null && _customBinaryPath!.isNotEmpty) {
        await _localStorage.setString(_binaryPathKey, _customBinaryPath!);
      }
    } catch (e) {
      _addLogMessage('保存设置失败: $e');
    }
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
      _isRunning = _server.isRunning;
      _currentPort = _server.port;
    });
    
    // 如果当前端口不为空，更新端口控制器
    if (_currentPort != null) {
      _portController.text = _currentPort.toString();
    }
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
  
  // 选择二进制文件
  Future<void> _selectBinaryFile() async {
    try {
      FilePickerResult? result = await FilePicker.platform.pickFiles(
        type: FileType.custom,
        allowedExtensions: System.isWindows ? ['exe'] : null,
      );
      
      if (result != null && result.files.single.path != null) {
        final path = result.files.single.path!;
        final file = File(path);
        
        if (await file.exists()) {
          setState(() {
            _customBinaryPath = path;
            _binaryPathController.text = path;
          });
          _addLogMessage('已选择二进制文件: $path');
          
          // 保存设置
          await _saveSettings();
        } else {
          _addLogMessage('所选文件不存在');
        }
      }
    } catch (e) {
      _addLogMessage('选择文件时发生错误: $e');
    }
  }

  Future<void> _restartServer() async {
    if (_isRestarting) return;
    
    setState(() {
      _isRestarting = true;
    });
    
    _addLogMessage('正在重启SERVER...');
    
    try {
      // 先停止服务
      await _server.stopService();
      
      // 然后启动服务，使用当前端口和自定义二进制路径（如果在调试模式下）
      final port = await _server.startService(
        requestedPort: _currentPort,
        customBinaryPath: _isDebugMode ? _customBinaryPath : null,
      );
      
      if (port != null) {
        _addLogMessage('SERVER重启成功，端口: $port');
        setState(() {
          _isRunning = true;
          _currentPort = port;
        });
        _updatePortController();
        
        // 保存设置
        await _saveSettings();
      } else {
        _addLogMessage('SERVER重启失败');
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

  Future<void> _stopServer() async {
    if (_isStopping) return;
    
    setState(() {
      _isStopping = true;
    });
    
    _addLogMessage('正在停止SERVER...');
    
    try {
      final result = await _server.stopService();
      
      if (result) {
        _addLogMessage('SERVER已停止');
        setState(() {
          _isRunning = false;
        });
      } else {
        _addLogMessage('停止SERVER失败');
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
      await _server.stopService();
      
      // 然后使用新端口启动服务
      final port = await _server.startService(
        requestedPort: newPort,
        customBinaryPath: _isDebugMode ? _customBinaryPath : null,
      );
      
      if (port != null) {
        _addLogMessage('端口更改成功，新端口: $port');
        setState(() {
          _isRunning = true;
          _currentPort = port;
        });
        _updatePortController();
        
        // 保存设置
        await _saveSettings();
      } else {
        _addLogMessage('端口更改失败');
        setState(() {
          _isRunning = false;
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
          'SERVER 调试',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '管理和监控本地SERVER',
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
                    if (_isRunning && _currentPort != null) ...[
                      const SizedBox(width: 8),
                      Text(
                        '端口: $_currentPort',
                        style: TextStyle(
                          fontSize: 12,
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                    const Spacer(),
                    if (_isDebugMode)
                      Chip(
                        label: const Text('调试模式'),
                        backgroundColor: Colors.amber[100],
                        labelStyle: TextStyle(color: Colors.amber[900]),
                        avatar: Icon(Icons.bug_report, size: 16, color: Colors.amber[900]),
                      ),
                  ],
                ),
                const SizedBox(height: 20),
                
                // 调试模式下的二进制文件路径设置
                if (_isDebugMode) ...[
                  Text(
                    '二进制文件路径',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: _binaryPathController,
                          decoration: InputDecoration(
                            labelText: '二进制文件路径',
                            hintText: '选择或输入SERVER二进制文件路径',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                            prefixIcon: const Icon(Icons.file_present),
                          ),
                          enabled: !_isChangingPort && !_isRestarting && !_isStopping,
                          onChanged: (value) {
                            _customBinaryPath = value.isNotEmpty ? value : null;
                          },
                        ),
                      ),
                      const SizedBox(width: 12),
                      ElevatedButton(
                        onPressed: !_isChangingPort && !_isRestarting && !_isStopping
                            ? _selectBinaryFile
                            : null,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.deepPurple,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                        ),
                        child: const Text('浏览...'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                ],
                
                // 端口设置
                Text(
                  '端口设置',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
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
                            ? _restartServer
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
                            ? _stopServer
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