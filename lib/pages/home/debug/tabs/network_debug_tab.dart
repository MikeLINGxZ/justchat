import 'package:flutter/material.dart';

class NetworkDebugTab extends StatelessWidget {
  const NetworkDebugTab({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '网络调试',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '监控和分析网络请求',
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
                    Icon(Icons.network_check, color: Colors.green[600]),
                    const SizedBox(width: 8),
                    Text(
                      '网络状态',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildInfoRow('连接状态', '已连接'),
                _buildInfoRow('网络类型', 'WiFi'),
                _buildInfoRow('延迟', '15ms'),
                _buildInfoRow('带宽', '100 Mbps'),
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
                    Icon(Icons.history, color: Colors.blue[600]),
                    const SizedBox(width: 8),
                    Text(
                      '最近请求',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                _buildRequestItem('GET /api/chat', '200 OK', '2.3s'),
                _buildRequestItem('POST /api/settings', '201 Created', '1.1s'),
                _buildRequestItem('GET /api/models', '200 OK', '0.8s'),
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

  Widget _buildRequestItem(String url, String status, String duration) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        children: [
          Expanded(
            child: Text(url, style: const TextStyle(fontSize: 12)),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: status.contains('200') ? Colors.green[100] : Colors.orange[100],
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              status,
              style: TextStyle(
                fontSize: 10,
                color: status.contains('200') ? Colors.green[700] : Colors.orange[700],
              ),
            ),
          ),
          const SizedBox(width: 8),
          Text(duration, style: const TextStyle(fontSize: 12, color: Colors.grey)),
        ],
      ),
    );
  }
} 