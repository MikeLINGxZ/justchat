import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/pages/home/assistant/assistant.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/storage/chat_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/style.dart';

/// 多标签页对话界面
class MultiTabAssistant extends ConsumerStatefulWidget {
  const MultiTabAssistant({super.key});

  @override
  ConsumerState<MultiTabAssistant> createState() => _MultiTabAssistantState();
}

class _MultiTabAssistantState extends ConsumerState<MultiTabAssistant>
    with TickerProviderStateMixin {
  late TabController _tabController;
  List<AssistantTab> _tabs = [];
  int _currentTabIndex = 0;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _initializeTabs();
  }

  @override
  void dispose() {
    _tabController.dispose();
    // 清理所有标签页的ConversationManager
    // ConversationManager继承自ChangeNotifier，会自动清理
    super.dispose();
  }

  /// 初始化标签页
  Future<void> _initializeTabs() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // 获取所有对话
      final conversations = await ChatStorage.getAllConversations();
      // 避免一次性打开多标签，仅打开一个（优先最近一个会话）
      _tabs.clear();
      if (conversations.isEmpty) {
        await _createNewTab();
      } else {
        // 仅取第一个会话作为初始标签
        await _createTabForConversation(conversations.first);
      }

      // 初始化TabController
      _tabController = TabController(
        length: _tabs.length,
        vsync: this,
        initialIndex: 0,
      );

      _tabController.addListener(_onTabChanged);
    } catch (e) {
      debugPrint('初始化标签页失败: $e');
      // 创建一个默认标签页作为备选
      await _createNewTab();
      _tabController = TabController(
        length: _tabs.length,
        vsync: this,
        initialIndex: 0,
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  /// 标签页切换监听
  void _onTabChanged() {
    setState(() {
      _currentTabIndex = _tabController.index;
    });
  }

  /// 为指定对话创建标签页
  Future<void> _createTabForConversation(Conversation conversation) async {
    final conversationManager = ConversationManager();
    await conversationManager.loadConversation(conversation.id);

    final tab = AssistantTab(
      id: conversation.id,
      title: conversation.title,
      conversationManager: conversationManager,
    );

    _tabs.add(tab);
  }

  /// 创建新的标签页
  Future<void> _createNewTab([String? title]) async {
    final conversationManager = ConversationManager();
    // 确保新的ConversationManager是完全空白的状态
    conversationManager.clearCurrentConversation();

    final tab = AssistantTab(
      id: 'new-${DateTime.now().millisecondsSinceEpoch}',
      title: title ?? '新对话',
      conversationManager: conversationManager,
    );

    setState(() {
      _tabs.add(tab);
    });

    // 如果TabController已经初始化，需要重新创建
    if (!_isLoading) {
      final oldController = _tabController;
      _tabController = TabController(
        length: _tabs.length,
        vsync: this,
        initialIndex: _tabs.length - 1, // 切换到新标签页
      );
      _tabController.addListener(_onTabChanged);
      oldController.dispose();

      setState(() {
        _currentTabIndex = _tabController.index;
      });
    }
  }

  /// 关闭标签页
  void _closeTab(int index) {
    if (_tabs.length <= 1) {
      // 至少保留一个标签页
      return;
    }

    setState(() {
      _tabs.removeAt(index);
    });

    // 重新创建TabController
    final oldController = _tabController;
    final newIndex = index >= _tabs.length ? _tabs.length - 1 : index;

    _tabController = TabController(
      length: _tabs.length,
      vsync: this,
      initialIndex: newIndex,
    );
    _tabController.addListener(_onTabChanged);
    oldController.dispose();

    setState(() {
      _currentTabIndex = _tabController.index;
    });
  }

  /// 构建标签页标题
  Widget _buildTabTitle(AssistantTab tab, int index) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(
          tab.title,
          overflow: TextOverflow.ellipsis,
          style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
          strutStyle: const StrutStyle(forceStrutHeight: true),
          textHeightBehavior: const TextHeightBehavior(
            applyHeightToFirstAscent: false,
            applyHeightToLastDescent: false,
          ),
        ),
        if (_tabs.length > 1)
          GestureDetector(
            onTap: () => _closeTab(index),
            child: Container(
              margin: const EdgeInsets.only(left: 8.0),
              padding: const EdgeInsets.all(2),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: Icon(
                Icons.close,
                size: 14,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
              ),
            ),
          ),
      ],
    );
  }

  /// 构建标签栏
  Widget _buildTabBar() {
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
        child: Container(
          padding: EdgeInsets.all(6),
          // padding: const EdgeInsets.only(left: 8, right: 8),
          decoration: BoxDecoration(
            color: Style.tertiaryBackground(context),
            // boxShadow: [
            //   BoxShadow(
            //     color: Colors.black.withOpacity(0.1),
            //     blurRadius: 10,
            //     offset: const Offset(0, 2),
            //   ),
            // ],
          ),
          child: Container(
            height: 38,
            child: Row(
            children: [
              Container(
                decoration: BoxDecoration(
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.circular(12),
                    topRight: Radius.circular(12),
                  ),
                ),
                child: TabBar(
                  controller: _tabController,
                  isScrollable: true,
                  tabAlignment: TabAlignment.start,
                  indicator: BoxDecoration(
                    color: Theme.of(context).colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  indicatorSize: TabBarIndicatorSize.tab,
                  indicatorPadding: const EdgeInsets.symmetric(
                    horizontal: 4,
                    vertical: 4,
                  ),
                  labelColor: Theme.of(context).colorScheme.onPrimaryContainer,
                  unselectedLabelColor: Theme.of(context).colorScheme.onSurface,
                  labelPadding: const EdgeInsets.symmetric(horizontal: 8),
                  dividerColor: Colors.transparent,
                  tabs:
                      _tabs.asMap().entries.map((entry) {
                        final index = entry.key;
                        final tab = entry.value;
                        return Tab(
                          key: ValueKey('tab-${tab.id}'),
                          child: Container(
                            constraints: const BoxConstraints(maxWidth: 200),
                            padding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 0,
                            ),
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: _buildTabTitle(tab, index),
                          ),
                        );
                      }).toList(),
                ),
              ),
              // 新建标签页按钮
              if (_tabs.length < 5)
                Container(
                  margin: const EdgeInsets.only(left: 8, right: 8),
                  child: GestureDetector(
                    onTap: () => _createNewTab(),
                    child: Container(
                      padding: const EdgeInsets.all(8),
                      decoration: BoxDecoration(
                        color: Theme.of(
                          context,
                        ).colorScheme.primary.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Icon(
                        Icons.add,
                        size: 14,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                    ),
                  ),
                ),
            ],
          ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_tabs.isEmpty) {
      return const Center(child: Text('没有可用的标签页'));
    }

    return Column(
      children: [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
          child: _buildTabBar(),
        ),
        Expanded(
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.only(
                bottomLeft: Radius.circular(12),
                bottomRight: Radius.circular(12),
              ),
              boxShadow: [
                BoxShadow(
                  color: Theme.of(context).colorScheme.shadow.withOpacity(0.1),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: ClipRRect(
              borderRadius: const BorderRadius.only(
                bottomLeft: Radius.circular(12),
                bottomRight: Radius.circular(12),
              ),
              child: IndexedStack(
                index: _currentTabIndex,
                children:
                    _tabs.map((tab) {
                      return AssistantPage(
                        key: ValueKey('assistant-${tab.id}'),
                        conversationManager: tab.conversationManager,
                      );
                    }).toList(),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

/// 标签页数据模型
class AssistantTab {
  final String id;
  final String title;
  final ConversationManager conversationManager;

  AssistantTab({
    required this.id,
    required this.title,
    required this.conversationManager,
  });
}
