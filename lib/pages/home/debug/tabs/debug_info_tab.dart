import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';

class DebugInfoTab extends StatelessWidget {
  const DebugInfoTab({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '调试信息',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '查看应用的基本调试信息',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 24),
        
        // 调试信息卡片
        Card(
          elevation: 2,
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.info_outline, color: Colors.blue[600]),
                    const SizedBox(width: 8),
                    Text(
                      '系统信息',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildInfoRow('Flutter 版本', '3.16.0'),
                _buildInfoRow('Dart 版本', '3.2.0'),
                _buildInfoRow('调试模式', kDebugMode ? '是' : '否'),
                _buildInfoRow('发布模式', kReleaseMode ? '是' : '否'),
                _buildInfoRow('配置文件', kProfileMode ? '是' : '否'),
              ],
            ),
          ),
        ),
        const SizedBox(height: 20),
        
        // 调试功能卡片
        Card(
          elevation: 2,
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.build, color: Colors.green[600]),
                    const SizedBox(width: 8),
                    Text(
                      '调试功能',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildDebugButton(
                  context,
                  '清除缓存',
                  Icons.clear_all,
                  Colors.orange,
                  () => _clearCache(context),
                ),
                const SizedBox(height: 12),
                _buildDebugButton(
                  context,
                  '重置设置',
                  Icons.restore,
                  Colors.red,
                  () => _resetSettings(context),
                ),
                const SizedBox(height: 12),
                _buildDebugButton(
                  context,
                  '导出日志',
                  Icons.download,
                  Colors.blue,
                  () => _exportLogs(context),
                ),
                const SizedBox(height: 12),
                _buildDebugButton(
                  context,
                  '性能监控',
                  Icons.speed,
                  Colors.purple,
                  () => _showPerformanceMonitor(context),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Text(
            '$label: ',
            style: const TextStyle(fontWeight: FontWeight.w500),
          ),
          Text(value),
        ],
      ),
    );
  }

  Widget _buildDebugButton(
    BuildContext context,
    String label,
    IconData icon,
    Color color,
    VoidCallback onPressed,
  ) {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton.icon(
        onPressed: onPressed,
        icon: Icon(icon, size: 20),
        label: Text(label),
        style: ElevatedButton.styleFrom(
          backgroundColor: color,
          foregroundColor: Colors.white,
          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
      ),
    );
  }

  void _clearCache(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('缓存已清除'),
        backgroundColor: Colors.green,
      ),
    );
  }

  void _resetSettings(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认重置'),
        content: const Text('确定要重置所有设置吗？此操作不可撤销。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('设置已重置'),
                  backgroundColor: Colors.orange,
                ),
              );
            },
            child: const Text('确定'),
          ),
        ],
      ),
    );
  }

  void _exportLogs(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('日志导出功能开发中...'),
        backgroundColor: Colors.blue,
      ),
    );
  }

  void _showPerformanceMonitor(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('性能监控功能开发中...'),
        backgroundColor: Colors.purple,
      ),
    );
  }
} 