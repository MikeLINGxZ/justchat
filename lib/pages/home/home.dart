import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/window_title_bar.dart';
import 'package:lemon_tea/controls/sidebar_icon_button.dart';
import 'package:lemon_tea/pages/home/assistant/assistant.dart';
import 'package:lemon_tea/pages/home/task/task.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  int _selectedIndex = 0;
  
  final List<Widget> _pages = [
    const AssistantPage(),
    const TaskPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          const WindowTitleBar(title: "Lemon Tea"),
          Expanded(
            child: Row(
              children: [
                // 侧边栏
                Container(
                  padding: const EdgeInsets.symmetric(vertical: 20.0, horizontal: 12.0),
                  child: Column(
                    children: [
                      SidebarIconButton(
                        icon: Icons.chat_bubble_outline,
                        isSelected: _selectedIndex == 0,
                        onPressed: () {
                          setState(() {
                            _selectedIndex = 0;
                          });
                        },
                      ),
                      const SizedBox(height: 14),
                      SidebarIconButton(
                        icon: Icons.task_outlined,
                        isSelected: _selectedIndex == 1,
                        onPressed: () {
                          setState(() {
                            _selectedIndex = 1;
                          });
                        },
                      ),
                    ],
                  ),
                ),
                const VerticalDivider(thickness: 1, width: 1),
                // 内容区域
                Expanded(
                  child: _pages[_selectedIndex],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
