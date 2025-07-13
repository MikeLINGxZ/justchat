import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';

import 'tabs/debug_info_tab.dart';
import 'tabs/network_debug_tab.dart';
import 'tabs/storage_debug_tab.dart';
import 'tabs/performance_debug_tab.dart';
import 'tabs/log_debug_tab.dart';
import 'tabs/ffi_debug_tab.dart';
import 'tabs/socket_debug_tab.dart';

class DebugPage extends StatefulWidget {
  const DebugPage({super.key});

  @override
  State<DebugPage> createState() => _DebugPageState();
}

class _DebugPageState extends State<DebugPage> {
  int _selectedIndex = 0;

  final List<DebugTab> _debugTabs = [
    DebugTab(
      title: '调试信息',
      icon: Icons.info_outline,
      color: Colors.blue,
      widget: const DebugInfoTab(),
    ),
    DebugTab(
      title: '网络调试',
      icon: Icons.network_check,
      color: Colors.green,
      widget: const NetworkDebugTab(),
    ),
    DebugTab(
      title: '存储调试',
      icon: Icons.storage,
      color: Colors.orange,
      widget: const StorageDebugTab(),
    ),
    DebugTab(
      title: '性能调试',
      icon: Icons.speed,
      color: Colors.purple,
      widget: const PerformanceDebugTab(),
    ),
    DebugTab(
      title: '日志调试',
      icon: Icons.article,
      color: Colors.red,
      widget: const LogDebugTab(),
    ),
    DebugTab(
      title: 'FFI 调试',
      icon: Icons.code,
      color: Colors.indigo,
      widget: const FfiDebugTab(),
    ),
    DebugTab(
      title: 'Socket 调试',
      icon: Icons.settings_ethernet,
      color: Colors.teal,
      widget: const SocketDebugTab(),
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      body: Row(
        children: [
          // 左侧菜单
          Container(
            width: 200,
            decoration: BoxDecoration(
              color: Colors.white,
              border: Border(
                right: BorderSide(color: Colors.grey[300]!),
              ),
            ),
            child: Column(
              children: [
                // 标题
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    border: Border(
                      bottom: BorderSide(color: Colors.grey[300]!),
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        Icons.bug_report,
                        size: 24,
                        color: Colors.orange[700],
                      ),
                      const SizedBox(width: 8),
                      Text(
                        '调试菜单',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: Colors.orange[700],
                        ),
                      ),
                    ],
                  ),
                ),
                
                // 菜单项
                Expanded(
                  child: ListView.builder(
                    itemCount: _debugTabs.length,
                    itemBuilder: (context, index) {
                      final tab = _debugTabs[index];
                      final isSelected = index == _selectedIndex;
                      
                      return Container(
                        margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        child: ListTile(
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                          leading: Icon(
                            tab.icon,
                            color: isSelected ? tab.color : Colors.grey[600],
                            size: 22,
                          ),
                          title: Text(
                            tab.title,
                            style: TextStyle(
                              color: isSelected ? tab.color : Colors.grey[800],
                              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                              fontSize: 14,
                            ),
                          ),
                          tileColor: isSelected ? tab.color.withOpacity(0.1) : Colors.white,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(8),
                            side: BorderSide(
                              color: isSelected ? tab.color : Colors.grey[200]!, 
                              // width: isSelected ? 2 : 1,
                            ),
                          ),
                          onTap: () {
                            setState(() {
                              _selectedIndex = index;
                            });
                          },
                        ),
                      );
                    },
                  ),
                ),
                
                // 底部提示
                Container(
                  padding: const EdgeInsets.all(12),
                  child: Row(
                    children: [
                      Icon(Icons.warning_amber_outlined, 
                           color: Colors.amber[700], size: 16),
                      const SizedBox(width: 4),
                      Expanded(
                        child: Text(
                          '仅调试模式可见',
                          style: TextStyle(
                            color: Colors.amber[700],
                            fontSize: 12,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          
          // 右侧内容区域
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(24.0),
              child: _debugTabs[_selectedIndex].widget,
            ),
          ),
        ],
      ),
    );
  }
}

class DebugTab {
  final String title;
  final IconData icon;
  final Color color;
  final Widget widget;

  DebugTab({
    required this.title,
    required this.icon,
    required this.color,
    required this.widget,
  });
} 