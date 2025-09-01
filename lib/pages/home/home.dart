import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/controls/window_title_bar.dart';
import 'package:lemon_tea/controls/expandable_sidebar.dart';
import 'package:lemon_tea/pages/home/assistant/multi_tab_assistant.dart';
import 'package:lemon_tea/pages/home/task/task.dart';
import 'package:lemon_tea/pages/home/history/history.dart';
import 'package:lemon_tea/pages/home/settings/settings.dart';
import 'package:lemon_tea/pages/home/plugins/plugins.dart';
import 'package:lemon_tea/pages/home/debug/debug.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/utils/style.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  int _selectedIndex = 0;
  late ConversationManager _conversationManager;
  late List<SidebarItem> _sidebarItems;

  final List<Widget> _pages = [];

  @override
  void initState() {
    super.initState();
    _conversationManager = ConversationManager();
    _initializeSidebarItems();
    _initializePages();
    // 初始化ConversationManager，加载对话历史
    _conversationManager.initialize();
  }

  void _initializeSidebarItems() {
    _sidebarItems = [
      const SidebarItem(
        icon: Icons.chat_bubble_outline,
        title: '助手',
        index: 0,
      ),
      const SidebarItem(
        icon: Icons.task_alt,
        title: '任务',
        index: 1,
      ),
      const SidebarItem(
        icon: Icons.history,
        title: '历史',
        index: 2,
      ),
      const SidebarItem(
        icon: Icons.extension,
        title: '插件',
        index: 3,
      ),
    ];

    // 在debug模式下添加debug项
    if (kDebugMode) {
      _sidebarItems.add(const SidebarItem(
        icon: Icons.bug_report,
        title: '调试',
        index: 4,
      ));
    }

    // 最后添加settings项
    _sidebarItems.add(SidebarItem(
      icon: Icons.settings,
      title: '设置',
      index: kDebugMode ? 5 : 4,
    ));
  }

  void _initializePages() {
    _pages.clear();
    _pages.addAll([
      const MultiTabAssistant(),
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
      backgroundColor: Style.secondaryBackground(context),
      body: Column(
        children: [
          const WindowTitleBar(title: "Lemon Tea"),
          Expanded(
            child: Container(
              padding: EdgeInsets.fromLTRB(4, 0, 4, 4),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(Style.radiusLv1),
                ),
                child: Row(
                  children: [
                    // 可展开侧边栏
                    ExpandableSidebar(
                      selectedIndex: _selectedIndex,
                      onItemSelected: (index) {
                        setState(() {
                          _selectedIndex = index;
                        });
                      },
                      items: _sidebarItems,
                    ),
                    const SizedBox(width: 4),

                    // 内容区域：使用 IndexedStack 保持各页面状态，切换无销毁
                    Expanded(
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.only(
                            topLeft: Radius.circular(Style.radiusLv1),
                            bottomLeft: Radius.circular(Style.radiusLv1),
                            topRight: Radius.circular(Style.radiusLv1),
                            bottomRight: Radius.circular(Style.radiusLv1),
                          ),
                          color: Style.primaryBackground(context),
                        ),
                        child: IndexedStack(
                          index: _selectedIndex,
                          children: _pages,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
