import 'dart:ffi';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:lemon_tea/ffi/nativefl.dart';
import 'package:ffi/ffi.dart';

class FfiDebugTab extends StatefulWidget {
  const FfiDebugTab({super.key});

  @override
  State<FfiDebugTab> createState() => _FfiDebugTabState();
}

class _FfiDebugTabState extends State<FfiDebugTab> {
  final TextEditingController _inputController = TextEditingController();
  final TextEditingController _outputController = TextEditingController();
  bool _isExecuting = false;

  @override
  void dispose() {
    _inputController.dispose();
    _outputController.dispose();
    super.dispose();
  }

  void _executeFFI() {
    setState(() {
      _isExecuting = true;
      _outputController.text = '正在执行...';
    });

    // 模拟执行过程
    Future.delayed(const Duration(milliseconds: 500), () {
      if (!mounted) return;
      
      try {
        // 获取动态库
        DynamicLibrary? dylib;
        try {
          if (Platform.isMacOS) {
            // 尝试不同的路径加载动态库
            final List<String> possiblePaths = [
              'example_ffi_arm64.dylib',
            ];
            
            for (final path in possiblePaths) {
              try {
                dylib = DynamicLibrary.open(path);
                _outputController.text = '成功加载动态库: $path\n';
                break;
              } catch (e) {
                // 继续尝试下一个路径
              }
            }
            
            if (dylib == null) {
              throw Exception('无法加载动态库，所有路径尝试均失败');
            }
          } else {
            dylib = DynamicLibrary.process();
          }
        } catch (e) {
          throw Exception('加载动态库失败: $e');
        }
        
        // 创建 Nativefl 实例
        final nativefl = Nativefl(dylib);
        
        // 获取输入文本
        final input = _inputController.text.isEmpty ? "默认输入" : _inputController.text;
        
        // 将 Dart 字符串转换为 C 字符串
        final inputPtr = input.toNativeUtf8().cast<Char>();
        
        // 调用本地函数
        final resultPtr = nativefl.ProcessString(inputPtr);
        
        // 将结果转换回 Dart 字符串
        String result = "无结果";
        if (resultPtr != nullptr) {
          result = resultPtr.cast<Utf8>().toDartString();
        }
        
        // 释放内存
        calloc.free(inputPtr.cast<Utf8>());
        
        setState(() {
          _isExecuting = false;
          _outputController.text += '执行结果：\n'
              '输入: $input\n'
              '时间: ${DateTime.now().toString()}\n'
              '状态: 成功\n'
              '返回值: $result';
        });
      } catch (e) {
        setState(() {
          _isExecuting = false;
          _outputController.text = '执行出错: $e';
        });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'FFI 调用',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '测试本地函数接口调用',
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
                    Icon(Icons.code, color: Colors.blue[600]),
                    const SizedBox(width: 8),
                    Text(
                      '函数调用',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                
                // 输入框
                TextField(
                  controller: _inputController,
                  decoration: InputDecoration(
                    labelText: '输入参数',
                    hintText: '请输入FFI调用参数',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    prefixIcon: const Icon(Icons.input),
                  ),
                  maxLines: 3,
                ),
                const SizedBox(height: 16),
                
                // 执行按钮
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    onPressed: _isExecuting ? null : _executeFFI,
                    icon: _isExecuting 
                        ? Container(
                            width: 24,
                            height: 24,
                            padding: const EdgeInsets.all(2.0),
                            child: const CircularProgressIndicator(
                              strokeWidth: 3,
                            ),
                          )
                        : const Icon(Icons.play_arrow),
                    label: Text(_isExecuting ? '执行中...' : '执行调用'),
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
                    maxLines: 10,
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
    );
  }
} 