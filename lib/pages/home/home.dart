import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/window_title_bar.dart';
import 'package:lemon_tea/controls/sidebar_icon_button.dart';
import 'package:lemon_tea/pages/home/assistant/assistant.dart';
import 'package:lemon_tea/pages/home/task/task.dart';
import 'package:lemon_tea/pages/home/history/history.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  int _selectedIndex = 0;
  late ConversationManager _conversationManager;
  
  final List<Widget> _pages = [];

  @override
  void initState() {
    super.initState();
    _conversationManager = ConversationManager();
    _initializePages();
  }

  void _initializePages() {
    _pages.clear();
    _pages.addAll([
      AssistantPage(conversationManager: _conversationManager),
      const TaskPage(),
      HistoryPage(
        conversationManager: _conversationManager,
        onConversationSelected: _handleConversationSelected,
        onConversationDeleted: _handleConversationDeleted,
        onNewConversation: _handleNewConversation,
      ),
    ]);
  }

  Future<void> _handleConversationSelected(Conversation conversation) async {
    await _conversationManager.loadConversation(conversation.id);
    // 切换到助手页面
    setState(() {
      _selectedIndex = 0;
    });
  }

  Future<void> _handleConversationDeleted(String conversationId) async {
    await _conversationManager.deleteConversation(conversationId);
  }

  void _handleNewConversation() {
    // 切换到助手页面
    setState(() {
      _selectedIndex = 0;
    });
  }

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
                      const SizedBox(height: 14),
                      SidebarIconButton(
                        icon: Icons.history,
                        isSelected: _selectedIndex == 2,
                        onPressed: () {
                          setState(() {
                            _selectedIndex = 2;
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
