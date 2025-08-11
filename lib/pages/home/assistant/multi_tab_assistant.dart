import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/pages/home/assistant/assistant.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/storage/chat_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';

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

      if (conversations.isEmpty) {
        // 如果没有对话，创建一个默认标签页
        await _createNewTab();
      } else {
        // 为每个对话创建一个标签页
        for (final conversation in conversations.take(5)) {
          // 限制最多5个标签页
          await _createTabForConversation(conversation);
        }
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
    if (_tabController.indexIsChanging) {
      setState(() {
        _currentTabIndex = _tabController.index;
      });
    }
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
          style: TextStyle(
            fontSize: FontSizeUtils.getSmallSize(ref),
            fontWeight: FontWeight.w500,
            height: 1.0,
          ),
          strutStyle: const StrutStyle(
            height: 1.0,
            leading: 0.0,
            forceStrutHeight: true,
          ),
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
    return Container(
      height: 40,
      padding: const EdgeInsets.only(left: 8, right: 8, top: 2,bottom: 2),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          bottom: BorderSide(color: Theme.of(context).dividerColor, width: 0.5),
        ),
      ),
      child: Row(
        children: [
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.surface,
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
                        child: Container(
                          constraints: const BoxConstraints(maxWidth: 200),
                          padding: const EdgeInsets.symmetric(
                            horizontal: 12,
                           vertical: 0,
                          ),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: _buildTabTitle(tab, index),
                        ),
                      );
                    }).toList(),
              ),
            ),
          ),
          // 新建标签页按钮
          if (_tabs.length < 5)
            Container(
              margin: const EdgeInsets.only(left: 8, right: 16),
              child: GestureDetector(
                onTap: () => _createNewTab(),
                child: Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Theme.of(
                      context,
                    ).colorScheme.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color: Theme.of(
                        context,
                      ).colorScheme.primary.withOpacity(0.3),
                      width: 1,
                    ),
                  ),
                  child: Icon(
                    Icons.add,
                    size: 18,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                ),
              ),
            ),
        ],
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
        _buildTabBar(),
        Expanded(
          child: Container(
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.surface,
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
                children: _tabs.map((tab) {
                  return AssistantPage(
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
