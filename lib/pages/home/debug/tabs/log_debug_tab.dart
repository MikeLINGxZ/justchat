import 'package:flutter/material.dart';

class LogDebugTab extends StatelessWidget {
  const LogDebugTab({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '日志调试',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '查看应用日志和错误信息',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: Colors.grey[600],
          ),
        ),
        const SizedBox(height: 24),
        
        Card(
          elevation: 2,
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.article, color: Colors.red[600]),
                    const SizedBox(width: 8),
                    Text(
                      '日志级别',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildLogLevel('ERROR', 5, Colors.red),
                _buildLogLevel('WARNING', 12, Colors.orange),
                _buildLogLevel('INFO', 45, Colors.blue),
                _buildLogLevel('DEBUG', 128, Colors.grey),
              ],
            ),
          ),
        ),
        const SizedBox(height: 20),
        
        Card(
          elevation: 2,
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.list, color: Colors.green[600]),
                    const SizedBox(width: 8),
                    Text(
                      '最近日志',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildLogEntry('INFO', '应用启动成功', '2024-01-15 10:30:15'),
                _buildLogEntry('DEBUG', '加载配置文件', '2024-01-15 10:30:16'),
                _buildLogEntry('WARNING', '网络连接较慢', '2024-01-15 10:30:17'),
                _buildLogEntry('ERROR', 'API请求失败', '2024-01-15 10:30:18'),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildLogLevel(String level, int count, Color color) {
    Color backgroundColor;
    Color textColor;
    
    if (color == Colors.red) {
      backgroundColor = Colors.red[100]!;
      textColor = Colors.red[700]!;
    } else if (color == Colors.orange) {
      backgroundColor = Colors.orange[100]!;
      textColor = Colors.orange[700]!;
    } else if (color == Colors.blue) {
      backgroundColor = Colors.blue[100]!;
      textColor = Colors.blue[700]!;
    } else {
      backgroundColor = Colors.grey[100]!;
      textColor = Colors.grey[700]!;
    }
    
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: backgroundColor,
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              level,
              style: TextStyle(
                fontSize: 10,
                color: textColor,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(width: 8),
          Text('$count 条'),
        ],
      ),
    );
  }

  Widget _buildLogEntry(String level, String message, String timestamp) {
    Color backgroundColor;
    Color textColor;
    
    switch (level) {
      case 'ERROR':
        backgroundColor = Colors.red[100]!;
        textColor = Colors.red[700]!;
        break;
      case 'WARNING':
        backgroundColor = Colors.orange[100]!;
        textColor = Colors.orange[700]!;
        break;
      case 'INFO':
        backgroundColor = Colors.blue[100]!;
        textColor = Colors.blue[700]!;
        break;
      default:
        backgroundColor = Colors.grey[100]!;
        textColor = Colors.grey[700]!;
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: backgroundColor,
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              level,
              style: TextStyle(
                fontSize: 10,
                color: textColor,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(message, style: const TextStyle(fontSize: 12)),
          ),
          Text(timestamp, style: const TextStyle(fontSize: 10, color: Colors.grey)),
        ],
      ),
    );
  }
} 