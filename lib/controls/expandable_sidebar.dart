import 'package:flutter/material.dart';
import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/utils/style.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';

class ExpandableSidebar extends StatefulWidget {
  final int selectedIndex;
  final Function(int) onItemSelected;
  final ConversationManager conversationManager;
  final Function(Conversation)? onConversationSelected;
  final VoidCallback? onNewConversation;

  const ExpandableSidebar({
    super.key,
    required this.selectedIndex,
    required this.onItemSelected,
    required this.conversationManager,
    this.onConversationSelected,
    this.onNewConversation,
  });

  @override
  State<ExpandableSidebar> createState() => _ExpandableSidebarState();
}

class _ExpandableSidebarState extends State<ExpandableSidebar>
    with TickerProviderStateMixin {
  bool _isExpanded = false;
  bool _isUserMenuExpanded = false;
  String _searchQuery = '';
  List<Conversation> _filteredConversations = [];
  late AnimationController _animationController;
  late Animation<double> _widthAnimation;
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    );
    _widthAnimation = Tween<double>(
      begin: 64.0, // 收缩时的宽度
      end: 280.0,  // 展开时的宽度
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
    
    // 监听对话管理器变化
    widget.conversationManager.addListener(_updateConversations);
    _updateConversations();
  }

  @override
  void dispose() {
    _animationController.dispose();
    _searchController.dispose();
    widget.conversationManager.removeListener(_updateConversations);
    super.dispose();
  }

  void _updateConversations() {
    if (mounted) {
      setState(() {
        if (_searchQuery.isEmpty) {
          _filteredConversations = widget.conversationManager.conversations;
        } else {
          _filteredConversations = widget.conversationManager.searchConversations(_searchQuery);
        }
      });
    }
  }

  void _toggleExpanded() {
    setState(() {
      _isExpanded = !_isExpanded;
      if (_isExpanded) {
        _animationController.forward();
      } else {
        _animationController.reverse();
        _isUserMenuExpanded = false; // 收缩时关闭用户菜单
      }
    });
  }

  void _onSearchChanged(String query) {
    setState(() {
      _searchQuery = query;
      _updateConversations();
    });
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _widthAnimation,
      builder: (context, child) {
        final currentWidth = _widthAnimation.value;
        final shouldShowExpanded = currentWidth > 150; // 只有当宽度足够时才显示展开内容
        
        return Container(
          width: currentWidth,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(Style.radiusLv1),
            color: Style.sidebarBackground(context),
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(Style.radiusLv1),
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 20.0, horizontal: 12.0),
              child: Column(
                children: [
                  // 顶部柠檬图标按钮
                  _buildToggleButton(shouldShowExpanded),
                  const SizedBox(height: 16),
                  
                  // 功能按钮区域 - 无论展开还是折叠都显示
                  _buildFunctionButtons(shouldShowExpanded),
                  if (shouldShowExpanded) const SizedBox(height: 16),
                  
                  // 新建对话按钮
                  if (shouldShowExpanded) _buildNewChatButton(),
                  if (shouldShowExpanded) const SizedBox(height: 12),
                  
                  // 搜索框
                  if (shouldShowExpanded) _buildSearchBox(),
                  if (shouldShowExpanded) const SizedBox(height: 12),
                  
                  // 历史对话列表
                  if (shouldShowExpanded) Expanded(child: _buildChatHistory()),
                  
                  const Spacer(),
                  
                  // 底部用户菜单
                  _buildUserMenu(shouldShowExpanded),
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _buildToggleButton(bool shouldShowExpanded) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: Tooltip(
        message: _isExpanded ? '收起菜单' : '展开菜单',
        child: GestureDetector(
          onTap: _toggleExpanded,
          child: Container(
            width: double.infinity,
            height: 40,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(8),
              color: Colors.transparent,
            ),
            child: shouldShowExpanded
                ? Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    child: Row(
                      children: [
                        const Text('🍋', style: TextStyle(fontSize: 20)),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            'Lemtea',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: Style.primaryText(context),
                            ),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                        Icon(
                          Icons.keyboard_arrow_left,
                          size: 16,
                          color: Style.secondaryText(context),
                        ),
                      ],
                    ),
                  )
                : Center(
                    child: const Text('🍋', style: TextStyle(fontSize: 20)),
                  ),
          ),
        ),
      ),
    );
  }

  Widget _buildFunctionButtons(bool shouldShowExpanded) {
    if (shouldShowExpanded) {
      // 展开状态：显示所有按钮
      return Column(
        children: [
          _buildFunctionButton(
            icon: Icons.task_alt,
            title: '任务',
            isSelected: widget.selectedIndex == 1,
            onTap: () => widget.onItemSelected(1),
            showExpanded: shouldShowExpanded,
          ),
          const SizedBox(height: 8),
          _buildFunctionButton(
            icon: Icons.extension,
            title: '插件',
            isSelected: widget.selectedIndex == 3,
            onTap: () => widget.onItemSelected(3),
            showExpanded: shouldShowExpanded,
          ),
        ],
      );
    } else {
      // 折叠状态：只显示图标，并增加间距
      return Column(
        children: [
          _buildFunctionButton(
            icon: Icons.task_alt,
            title: '任务',
            isSelected: widget.selectedIndex == 1,
            onTap: () => widget.onItemSelected(1),
            showExpanded: shouldShowExpanded,
          ),
          const SizedBox(height: 14),
          _buildFunctionButton(
            icon: Icons.extension,
            title: '插件',
            isSelected: widget.selectedIndex == 3,
            onTap: () => widget.onItemSelected(3),
            showExpanded: shouldShowExpanded,
          ),
        ],
      );
    }
  }

  Widget _buildFunctionButton({
    required IconData icon,
    required String title,
    required bool isSelected,
    required VoidCallback onTap,
    required bool showExpanded,
  }) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: Tooltip(
        message: showExpanded ? '' : title,
        child: GestureDetector(
          onTap: onTap,
          child: Container(
            width: showExpanded ? double.infinity : 40,
            height: showExpanded ? 36 : 40,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(8),
              color: isSelected ? Style.secondaryColor(context) : Colors.transparent,
            ),
            padding: showExpanded ? const EdgeInsets.symmetric(horizontal: 12) : null,
            child: showExpanded
                ? Row(
                    children: [
                      Icon(
                        icon,
                        size: 18,
                        color: isSelected ? Style.primaryColor(context) : Style.secondaryText(context),
                      ),
                      const SizedBox(width: 12),
                      Text(
                        title,
                        style: TextStyle(
                          color: isSelected ? Style.primaryColor(context) : Style.primaryText(context),
                          fontSize: 14,
                          fontWeight: isSelected ? FontWeight.w500 : FontWeight.normal,
                        ),
                      ),
                    ],
                  )
                : Center(
                    child: Icon(
                      icon,
                      size: 20,
                      color: isSelected ? Style.primaryColor(context) : Style.secondaryText(context),
                    ),
                  ),
          ),
        ),
      ),
    );
  }

  Widget _buildNewChatButton() {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: () {
          widget.onNewConversation?.call();
          widget.onItemSelected(0); // 切换到助手页面
        },
        child: Container(
          width: double.infinity,
          height: 36,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            border: Border.all(
              color: Style.primaryBorder(context),
              width: 1,
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.add,
                size: 18,
                color: Style.primaryText(context),
              ),
              const SizedBox(width: 8),
              Text(
                '新建对话',
                style: TextStyle(
                  color: Style.primaryText(context),
                  fontSize: 14,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSearchBox() {
    return Container(
      height: 32,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(6),
        color: Style.tertiaryBackground(context),
      ),
      child: TextField(
        controller: _searchController,
        onChanged: _onSearchChanged,
        style: TextStyle(
          color: Style.primaryText(context),
          fontSize: 12,
        ),
        decoration: InputDecoration(
          hintText: '搜索对话...',
          hintStyle: TextStyle(
            color: Style.hintText(context),
            fontSize: 12,
          ),
          prefixIcon: Icon(
            Icons.search,
            size: 16,
            color: Style.hintText(context),
          ),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        ),
      ),
    );
  }

  Widget _buildChatHistory() {
    if (_filteredConversations.isEmpty) {
      return Center(
        child: Text(
          '暂无对话记录',
          style: TextStyle(
            color: Style.hintText(context),
            fontSize: 12,
          ),
        ),
      );
    }

    return ListView.builder(
      itemCount: _filteredConversations.length,
      itemBuilder: (context, index) {
        final conversation = _filteredConversations[index];
        return _buildChatHistoryItem(conversation);
      },
    );
  }

  Widget _buildChatHistoryItem(Conversation conversation) {
    final isSelected = widget.conversationManager.currentConversation?.id == conversation.id;
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        child: GestureDetector(
          onTap: () {
            widget.onConversationSelected?.call(conversation);
          },
          child: Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(6),
              color: isSelected ? Style.secondaryColor(context) : Colors.transparent,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  conversation.title.isNotEmpty ? conversation.title : '新对话',
                  style: TextStyle(
                    color: isSelected ? Style.primaryColor(context) : Style.primaryText(context),
                    fontSize: 13,
                    fontWeight: isSelected ? FontWeight.w500 : FontWeight.normal,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Text(
                  _formatDate(conversation.updatedAt),
                  style: TextStyle(
                    color: Style.hintText(context),
                    fontSize: 11,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildUserMenu(bool shouldShowExpanded) {
    if (!shouldShowExpanded) {
      return MouseRegion(
        cursor: SystemMouseCursors.click,
        child: Tooltip(
          message: '用户菜单',
          child: GestureDetector(
            onTap: _toggleExpanded,
            child: Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(20),
                color: Style.tertiaryBackground(context),
              ),
              child: Icon(
                Icons.person,
                size: 20,
                color: Style.primaryText(context),
              ),
            ),
          ),
        ),
      );
    }

    return Column(
      children: [
        // 用户信息
        MouseRegion(
          cursor: SystemMouseCursors.click,
          child: GestureDetector(
            onTap: () {
              setState(() {
                _isUserMenuExpanded = !_isUserMenuExpanded;
              });
            },
            child: Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
                color: _isUserMenuExpanded ? Style.secondaryColor(context) : Colors.transparent,
              ),
              child: Row(
                children: [
                  Container(
                    width: 32,
                    height: 32,
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(16),
                      color: Style.tertiaryBackground(context),
                    ),
                    child: Icon(
                      Icons.person,
                      size: 18,
                      color: Style.primaryText(context),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'lpxqu',
                          style: TextStyle(
                            color: Style.primaryText(context),
                            fontSize: 14,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                        Text(
                          'lpxqu@qq.com',
                          style: TextStyle(
                            color: Style.hintText(context),
                            fontSize: 11,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Icon(
                    _isUserMenuExpanded ? Icons.keyboard_arrow_up : Icons.keyboard_arrow_down,
                    size: 16,
                    color: Style.secondaryText(context),
                  ),
                ],
              ),
            ),
          ),
        ),
        
        // 用户菜单项
        if (_isUserMenuExpanded) ..._buildUserMenuItems(),
      ],
    );
  }

  List<Widget> _buildUserMenuItems() {
    return [
      const SizedBox(height: 8),
      
      // 设置
      _buildUserMenuItem(
        icon: Icons.settings,
        title: '设置',
        onTap: () {
          widget.onItemSelected(kDebugMode ? 5 : 4);
          setState(() {
            _isUserMenuExpanded = false;
          });
        },
      ),
      
      // 主题
      _buildThemeMenuItem(),
      
      // 调试（仅Debug模式）
      if (kDebugMode) _buildUserMenuItem(
        icon: Icons.bug_report,
        title: '调试',
        onTap: () {
          widget.onItemSelected(4);
          setState(() {
            _isUserMenuExpanded = false;
          });
        },
      ),
      
      // 退出登录
      _buildUserMenuItem(
        icon: Icons.logout,
        title: '退出登录',
        onTap: () {
          // TODO: 实现退出登录逻辑
          print('退出登录');
        },
      ),
    ];
  }

  Widget _buildUserMenuItem({
    required IconData icon,
    required String title,
    required VoidCallback onTap,
  }) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          width: double.infinity,
          height: 32,
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Row(
            children: [
              Icon(
                icon,
                size: 16,
                color: Style.secondaryText(context),
              ),
              const SizedBox(width: 12),
              Text(
                title,
                style: TextStyle(
                  color: Style.primaryText(context),
                  fontSize: 13,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildThemeMenuItem() {
    return PopupMenuButton<String>(
      offset: const Offset(-120, 0),
      child: Container(
        width: double.infinity,
        height: 32,
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Row(
          children: [
            Icon(
              Icons.palette,
              size: 16,
              color: Style.secondaryText(context),
            ),
            const SizedBox(width: 12),
            Text(
              '主题',
              style: TextStyle(
                color: Style.primaryText(context),
                fontSize: 13,
              ),
            ),
            const Spacer(),
            Icon(
              Icons.keyboard_arrow_right,
              size: 14,
              color: Style.secondaryText(context),
            ),
          ],
        ),
      ),
      itemBuilder: (context) => [
        PopupMenuItem(
          value: 'auto',
          child: Row(
            children: [
              Icon(Icons.brightness_auto, size: 16, color: Style.secondaryText(context)),
              const SizedBox(width: 8),
              Text('自动', style: TextStyle(fontSize: 13, color: Style.primaryText(context))),
            ],
          ),
        ),
        PopupMenuItem(
          value: 'light',
          child: Row(
            children: [
              Icon(Icons.light_mode, size: 16, color: Style.secondaryText(context)),
              const SizedBox(width: 8),
              Text('浅色', style: TextStyle(fontSize: 13, color: Style.primaryText(context))),
            ],
          ),
        ),
        PopupMenuItem(
          value: 'dark',
          child: Row(
            children: [
              Icon(Icons.dark_mode, size: 16, color: Style.secondaryText(context)),
              const SizedBox(width: 8),
              Text('深色', style: TextStyle(fontSize: 13, color: Style.primaryText(context))),
            ],
          ),
        ),
      ],
      onSelected: (value) {
        // TODO: 实现主题切换逻辑
        print('切换主题: $value');
      },
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);
    
    if (difference.inDays == 0) {
      return '今天';
    } else if (difference.inDays == 1) {
      return '昨天';
    } else if (difference.inDays < 7) {
      return '${difference.inDays}天前';
    } else {
      return '${date.month}/${date.day}';
    }
  }
}