import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/controls/window_title_bar.dart';
import 'package:lemon_tea/controls/sidebar_icon_button.dart';
import 'package:lemon_tea/pages/home/assistant/assistant.dart';
import 'package:lemon_tea/pages/home/task/task.dart';
import 'package:lemon_tea/pages/home/history/history.dart';
import 'package:lemon_tea/pages/home/settings/settings.dart';
import 'package:lemon_tea/pages/home/plugins/plugins.dart';
import 'package:lemon_tea/pages/home/debug/debug.dart';
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
    // 初始化ConversationManager，加载对话历史
    _conversationManager.initialize();
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
      const PluginsPage(),
    ]);

    // 在debug模式下添加debug页面
    if (kDebugMode) {
      _pages.add(const DebugPage());
    }

    // 最后添加settings页面
    _pages.add(const SettingsPage());
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
      backgroundColor: Theme.of(context).appBarTheme.backgroundColor,
      body: Column(
        children: [
          const WindowTitleBar(title: "Lemon Tea"),
          Expanded(
            child: Container(
              padding: EdgeInsets.fromLTRB(4, 0, 4, 4),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(10),
                  color: Theme.of(context).scaffoldBackgroundColor,
                ),
                child: Row(
                  children: [
                    // 侧边栏
                    Container(
                      padding: const EdgeInsets.symmetric(
                        vertical: 20.0,
                        horizontal: 12.0,
                      ),
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
                          const SizedBox(height: 14),
                          SidebarIconButton(
                            icon: Icons.extension,
                            isSelected: _selectedIndex == 3,
                            onPressed: () {
                              setState(() {
                                _selectedIndex = 3;
                              });
                            },
                          ),
                          const Spacer(),
                          // 在debug模式下显示debug按钮
                          if (kDebugMode) ...[
                            SidebarIconButton(
                              icon: Icons.bug_report,
                              isSelected: _selectedIndex == 4,
                              onPressed: () {
                                setState(() {
                                  _selectedIndex = 4;
                                });
                              },
                            ),
                            const SizedBox(height: 14),
                          ],
                          SidebarIconButton(
                            icon: Icons.settings,
                            isSelected: _selectedIndex == (kDebugMode ? 5 : 4),
                            onPressed: () {
                              setState(() {
                                _selectedIndex = kDebugMode ? 5 : 4;
                              });
                            },
                          ),
                        ],
                      ),
                    ),
                    VerticalDivider(thickness: 1, width: 1,color: Theme.of(context).dividerColor),
                    // 内容区域
                    Expanded(child: _pages[_selectedIndex]),
                  ],
                ),
              )
            ),
          ),
        ],
      ),
    );
  }
}
