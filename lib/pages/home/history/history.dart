import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/llm/models/message.dart' as llm_message;
import 'package:lemon_tea/models/message.dart' as db_message;
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/storage/chat_storage.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_view.dart';

class HistoryPage extends ConsumerStatefulWidget {
  final ConversationManager conversationManager;
  final Function(Conversation) onConversationSelected;
  final Function(String) onConversationDeleted;
  final VoidCallback? onNewConversation;

  const HistoryPage({
    super.key,
    required this.conversationManager,
    required this.onConversationSelected,
    required this.onConversationDeleted,
    this.onNewConversation,
  });

  @override
  ConsumerState<HistoryPage> createState() => _HistoryPageState();
}

class _HistoryPageState extends ConsumerState<HistoryPage> with TickerProviderStateMixin {
  String _searchQuery = '';
  List<Conversation> _filteredConversations = [];
  Map<String, int> _messageCountCache = {};
  Map<String, String> _previewCache = {};
  bool _isLoadingDetails = false;
  String? _hoveredConversationId;
  late AnimationController _fadeController;

  @override
  void initState() {
    super.initState();
    _fadeController = AnimationController(
      duration: const Duration(milliseconds: 300),
      vsync: this,
    );
    // 监听 ConversationManager 的变化
    widget.conversationManager.addListener(_onConversationManagerChanged);
    _updateFilteredConversations();
    // 异步加载对话详情
    _loadConversationDetails();
    _fadeController.forward();
  }

  @override
  void dispose() {
    _fadeController.dispose();
    widget.conversationManager.removeListener(_onConversationManagerChanged);
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，刷新UI并重新加载详情
    if (mounted) {
      setState(() {
        _updateFilteredConversations();
      });
      _loadConversationDetails();
    }
  }

  /// 加载所有对话的详情信息（消息数量和预览）
  Future<void> _loadConversationDetails() async {
    if (_isLoadingDetails) return;
    
    setState(() {
      _isLoadingDetails = true;
    });

    try {
      final conversations = widget.conversationManager.conversations;
      
      // 并行加载所有对话的详情
      final futures = conversations.map((conversation) async {
        try {
          final messages = await ChatStorage.getMessagesByConversationId(conversation.id);
          final messageCount = messages.length;
          String preview = '暂无消息';
          
          if (messages.isNotEmpty) {
            // 获取最后一条用户或助手消息作为预览
            final lastMessage = messages.lastWhere(
              (msg) => msg.role == 'user' || msg.role == 'assistant',
              orElse: () => messages.last,
            );
            preview = lastMessage.content.length > 50 
                ? '${lastMessage.content.substring(0, 50)}...' 
                : lastMessage.content;
            // 移除换行符和多余空格
            preview = preview.replaceAll('\n', ' ').trim();
          }
          
          return {
            'id': conversation.id,
            'messageCount': messageCount,
            'preview': preview,
          };
        } catch (e) {
          debugPrint('加载对话 ${conversation.id} 详情失败: $e');
          return {
            'id': conversation.id,
            'messageCount': 0,
            'preview': '加载失败',
          };
        }
      });

      final results = await Future.wait(futures);
      
      // 更新缓存
      for (final result in results) {
        _messageCountCache[result['id'] as String] = result['messageCount'] as int;
        _previewCache[result['id'] as String] = result['preview'] as String;
      }
      
      if (mounted) {
        setState(() {});
      }
    } catch (e) {
      debugPrint('加载对话详情失败: $e');
    } finally {
      if (mounted) {
        setState(() {
          _isLoadingDetails = false;
        });
      }
    }
  }

  void _updateFilteredConversations() {
    final conversations = widget.conversationManager.conversations;
    if (_searchQuery.isEmpty) {
      _filteredConversations = conversations;
    } else {
      _filteredConversations = widget.conversationManager.searchConversations(_searchQuery);
    }
  }

  /// 转换数据库消息为 LLM 消息
  List<llm_message.Message> _convertDbMessagesToLlmMessages(List<db_message.Message> dbMessages) {
    return dbMessages.map((dbMsg) => llm_message.Message(
      role: dbMsg.role,
      content: dbMsg.content,
    )).toList();
  }

  /// 显示对话详情对话框
  Future<void> _showConversationDetailDialog(Conversation conversation) async {
    try {
      final dbMessages = await ChatStorage.getMessagesByConversationId(conversation.id);
      final llmMessages = _convertDbMessagesToLlmMessages(dbMessages);
      
      if (!mounted) return;
      
      await showDialog(
        context: context,
        barrierDismissible: true,
        builder: (context) => _ConversationDetailDialog(
          conversation: conversation,
          messages: llmMessages,
          onContinueConversation: () {
            Navigator.of(context).pop();
            widget.onConversationSelected(conversation);
          },
        ),
      );
    } catch (e) {
      debugPrint('加载对话详情失败: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              '加载对话详情失败: $e',
              style: const TextStyle(fontSize: 14), // SnackBar 默认使用固定大小
            ),
          ),
        );
      }
    }
  }

  /// 显示更多操作菜单
  Future<void> _showMoreActionsMenu(BuildContext context, WidgetRef ref, Conversation conversation) async {
    final RenderBox button = context.findRenderObject() as RenderBox;
    final RenderBox overlay = Overlay.of(context).context.findRenderObject() as RenderBox;
    final RelativeRect position = RelativeRect.fromRect(
      Rect.fromPoints(
        button.localToGlobal(Offset.zero, ancestor: overlay),
        button.localToGlobal(button.size.bottomRight(Offset.zero), ancestor: overlay),
      ),
      Offset.zero & overlay.size,
    );

    final String? action = await showMenu<String>(
      context: context,
      position: position,
      items: [
        PopupMenuItem<String>(
          value: 'view',
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.visibility_outlined, size: 18, color: Theme.of(context).colorScheme.primary),
              const SizedBox(width: 12),
              Text('查看', style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref))),
            ],
          ),
        ),
        PopupMenuItem<String>(
          value: 'continue',
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.chat_bubble_outline, size: 18, color: Theme.of(context).colorScheme.primary),
              const SizedBox(width: 12),
              Text('继续对话', style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref))),
            ],
          ),
        ),
        const PopupMenuDivider(),
        PopupMenuItem<String>(
          value: 'delete',
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.delete_outline, size: 18, color: Colors.red[600]),
              const SizedBox(width: 12),
              Text('删除', style: TextStyle(
                color: Colors.red[600],
                fontSize: FontSizeUtils.getSmallSize(ref),
              )),
            ],
          ),
        ),
      ],
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      elevation: 8,
    );

    if (action != null) {
      switch (action) {
        case 'view':
          _showConversationDetailDialog(conversation);
          break;
        case 'continue':
          widget.onConversationSelected(conversation);
          break;
        case 'delete':
          _handleDeleteConversation(conversation.id);
          break;
      }
    }
  }

  Future<void> _handleDeleteConversation(String conversationId) async {
    // 显示确认对话框
    final shouldDelete = await showDialog<bool>(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        backgroundColor: Theme.of(context).colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        title: Row(
          children: [
            Icon(
              Icons.warning_amber_rounded,
              color: Colors.orange[600],
              size: 24,
            ),
            const SizedBox(width: 12),
            Text(
              S.of(context).confirmDelete,
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        content: Padding(
          padding: const EdgeInsets.only(top: 8.0),
          child:           Text(
            S.of(context).confirmDeleteConversation,
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
              color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.8),
              height: 1.4,
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            style: TextButton.styleFrom(
              foregroundColor: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
            ),
            child: Text(S.of(context).cancel, style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
          ),
          const SizedBox(width: 8),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red[600],
              foregroundColor: Colors.white,
              elevation: 0,
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(S.of(context).delete, style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
          ),
        ],
        actionsPadding: const EdgeInsets.fromLTRB(24, 8, 24, 24),
      ),
    );

    if (shouldDelete == true) {
      await widget.onConversationDeleted(conversationId);
    }
  }

  Widget _buildSearchBox() {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Theme.of(context).colorScheme.outline.withValues(alpha: 0.2),
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.04),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: TextField(
        decoration: InputDecoration(
          hintText: S.of(context).searchConversations,
          hintStyle: TextStyle(
            color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.5),
            fontSize: FontSizeUtils.getBodySize(ref),
          ),
          prefixIcon: Icon(
            Icons.search_rounded,
            color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
            size: 20,
          ),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        ),
        style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
        onChanged: (value) {
          setState(() {
            _searchQuery = value;
            _updateFilteredConversations();
          });
        },
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: FadeTransition(
        opacity: _fadeController,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 96,
              height: 96,
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.primary.withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(
                _searchQuery.isEmpty ? Icons.chat_bubble_outline_rounded : Icons.search_off_rounded,
                size: 48,
                color: Theme.of(context).colorScheme.primary.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 24),
            Text(
              _searchQuery.isEmpty ? S.of(context).noConversationHistory : '未找到相关对话',
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.w500,
                color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.8),
              ),
            ),
            if (_searchQuery.isEmpty) ...[
              const SizedBox(height: 8),
              Text(
                '开始一个新的对话吧',
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                  color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
                ),
              ),
              const SizedBox(height: 32),
              if (widget.onNewConversation != null)
                ElevatedButton.icon(
                  onPressed: widget.onNewConversation,
                  icon: const Icon(Icons.add_rounded, size: 18),
                  label: Text(S.of(context).newConversation, style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Theme.of(context).colorScheme.primary,
                    foregroundColor: Theme.of(context).colorScheme.onPrimary,
                    elevation: 2,
                    padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                  ),
                ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildConversationCard(Conversation conversation, bool isCurrent) {
    final isHovered = _hoveredConversationId == conversation.id;
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: MouseRegion(
        onEnter: (_) => setState(() => _hoveredConversationId = conversation.id),
        onExit: (_) => setState(() => _hoveredConversationId = null),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          decoration: BoxDecoration(
            color: isCurrent 
                ? Theme.of(context).colorScheme.primary.withValues(alpha: 0.1)
                : isHovered 
                    ? Theme.of(context).colorScheme.surface.withValues(alpha: 0.8)
                    : Theme.of(context).colorScheme.surface,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isCurrent 
                  ? Theme.of(context).colorScheme.primary.withValues(alpha: 0.3)
                  : Theme.of(context).colorScheme.outline.withValues(alpha: 0.1),
              width: isCurrent ? 1.5 : 1,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: isHovered ? 0.08 : 0.04),
                blurRadius: isHovered ? 12 : 6,
                offset: Offset(0, isHovered ? 4 : 2),
              ),
            ],
          ),
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              borderRadius: BorderRadius.circular(12),
              onTap: () => _showConversationDetailDialog(conversation),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // 会话图标
                    Container(
                      width: 40,
                      height: 40,
                      decoration: BoxDecoration(
                        color: isCurrent 
                            ? Theme.of(context).colorScheme.primary.withValues(alpha: 0.2)
                            : Theme.of(context).colorScheme.primary.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Icon(
                        Icons.chat_bubble_rounded,
                        size: 20,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                    ),
                    const SizedBox(width: 12),
                    
                    // 会话信息
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            conversation.title,
                            style: TextStyle(
                              fontSize: FontSizeUtils.getBodySize(ref),
                              fontWeight: isCurrent ? FontWeight.w600 : FontWeight.w500,
                              color: Theme.of(context).colorScheme.onSurface,
                              height: 1.3,
                            ),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 6),
                          Text(
                            _previewCache[conversation.id] ?? '加载中...',
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: TextStyle(
                              fontSize: FontSizeUtils.getSmallSize(ref),
                              color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                              height: 1.4,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            children: [
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                decoration: BoxDecoration(
                                  color: Theme.of(context).colorScheme.primary.withValues(alpha: 0.1),
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  S.of(context).messagesCount(_messageCountCache[conversation.id] ?? 0),
                                  style: TextStyle(
                                    fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                                    fontWeight: FontWeight.w500,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                              ),
                              const Spacer(),
                              Text(
                                _formatDate(conversation.updatedAt),
                                style: TextStyle(
                                  fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                                  color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.5),
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                    
                    // 更多操作按钮
                    const SizedBox(width: 8),
                    AnimatedOpacity(
                      opacity: isHovered ? 1.0 : 0.0,
                      duration: const Duration(milliseconds: 200),
                      child: Builder(
                        builder: (context) => IconButton(
                          onPressed: () => _showMoreActionsMenu(context, ref, conversation),
                          icon: const Icon(Icons.more_vert_rounded, size: 18),
                          tooltip: '更多操作',
                          style: IconButton.styleFrom(
                            foregroundColor: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                            backgroundColor: Theme.of(context).colorScheme.surface.withValues(alpha: 0.8),
                            padding: const EdgeInsets.all(8),
                            minimumSize: const Size(32, 32),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final currentConversation = widget.conversationManager.currentConversation;
    
    return Container(
      color: Theme.of(context).colorScheme.background,
      child: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 标题栏
            FadeTransition(
              opacity: _fadeController,
              child: Row(
                children: [
                  Text(
                    S.of(context).conversationHistory,
                    style: TextStyle(
                      fontSize: FontSizeUtils.getHeadingSize(ref) + 2,
                      fontWeight: FontWeight.w700,
                      color: Theme.of(context).colorScheme.onBackground,
                    ),
                  ),
                  const Spacer(),
                  if (_isLoadingDetails)
                    Container(
                      margin: const EdgeInsets.only(right: 16),
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        valueColor: AlwaysStoppedAnimation<Color>(
                          Theme.of(context).colorScheme.primary.withValues(alpha: 0.7),
                        ),
                      ),
                    ),
                  if (widget.onNewConversation != null)
                    ElevatedButton.icon(
                      onPressed: widget.onNewConversation,
                      icon: const Icon(Icons.add_rounded, size: 18),
                      label: Text(S.of(context).newConversation, style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Theme.of(context).colorScheme.primary,
                        foregroundColor: Theme.of(context).colorScheme.onPrimary,
                        elevation: 2,
                        shadowColor: Theme.of(context).colorScheme.primary.withValues(alpha: 0.3),
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(10),
                        ),
                      ),
                    ),
                ],
              ),
            ),
            const SizedBox(height: 24),
            
            // 搜索框
            FadeTransition(
              opacity: _fadeController,
              child: _buildSearchBox(),
            ),
            const SizedBox(height: 24),
            
            // 对话列表
            Expanded(
              child: _filteredConversations.isEmpty
                  ? _buildEmptyState()
                  : FadeTransition(
                      opacity: _fadeController,
                      child: ListView.builder(
                        itemCount: _filteredConversations.length,
                        itemBuilder: (context, index) {
                          final conversation = _filteredConversations[index];
                          final isCurrent = currentConversation?.id == conversation.id;
                          
                          return _buildConversationCard(conversation, isCurrent);
                        },
                      ),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);
    
    if (difference.inDays == 0) {
      if (difference.inHours == 0) {
        return '${difference.inMinutes}分钟前';
      }
      return '${difference.inHours}小时前';
    } else if (difference.inDays == 1) {
      return '昨天';
    } else if (difference.inDays < 7) {
      return '${difference.inDays}天前';
    } else {
      return '${date.month}-${date.day}';
    }
  }
}

/// 对话详情对话框
class _ConversationDetailDialog extends ConsumerWidget {
  final Conversation conversation;
  final List<llm_message.Message> messages;
  final VoidCallback onContinueConversation;

  const _ConversationDetailDialog({
    required this.conversation,
    required this.messages,
    required this.onContinueConversation,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Dialog(
      backgroundColor: Colors.transparent,
      child: Container(
        width: MediaQuery.of(context).size.width * 0.8,
        height: MediaQuery.of(context).size.height * 0.8,
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.surface,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.15),
              blurRadius: 24,
              offset: const Offset(0, 8),
            ),
          ],
        ),
        child: Column(
          children: [
            // 标题栏
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.surface,
                borderRadius: const BorderRadius.only(
                  topLeft: Radius.circular(16),
                  topRight: Radius.circular(16),
                ),
                border: Border(
                  bottom: BorderSide(
                    color: Theme.of(context).colorScheme.outline.withValues(alpha: 0.2),
                  ),
                ),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.chat_bubble_rounded,
                    size: 24,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          conversation.title,
                          style: TextStyle(
                            fontSize: FontSizeUtils.getSubheadingSize(ref),
                            fontWeight: FontWeight.w600,
                            color: Theme.of(context).colorScheme.onSurface,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${messages.length} 条消息',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                            color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                          ),
                        ),
                      ],
                    ),
                  ),
                  IconButton(
                    onPressed: () => Navigator.of(context).pop(),
                    icon: const Icon(Icons.close_rounded),
                    style: IconButton.styleFrom(
                      foregroundColor: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                      backgroundColor: Theme.of(context).colorScheme.surface,
                      padding: const EdgeInsets.all(8),
                    ),
                  ),
                ],
              ),
            ),
            
            // 消息列表
            Expanded(
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 8),
                child: messages.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.chat_bubble_outline_rounded,
                              size: 64,
                              color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.3),
                            ),
                            const SizedBox(height: 16),
                            Text(
                              '暂无消息',
                              style: TextStyle(
                                fontSize: FontSizeUtils.getSubheadingSize(ref),
                                color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
                              ),
                            ),
                          ],
                        ),
                      )
                    : MessageView(messages),
              ),
            ),
            
            // 底部按钮栏
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.surface,
                borderRadius: const BorderRadius.only(
                  bottomLeft: Radius.circular(16),
                  bottomRight: Radius.circular(16),
                ),
                border: Border(
                  top: BorderSide(
                    color: Theme.of(context).colorScheme.outline.withValues(alpha: 0.2),
                  ),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                    onPressed: () => Navigator.of(context).pop(),
                    style: TextButton.styleFrom(
                      foregroundColor: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    ),
                    child: Text('关闭', style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
                  ),
                  const SizedBox(width: 12),
                  ElevatedButton.icon(
                    onPressed: onContinueConversation,
                    icon: const Icon(Icons.chat_rounded, size: 18),
                    label: Text('继续对话', style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Theme.of(context).colorScheme.primary,
                      foregroundColor: Theme.of(context).colorScheme.onPrimary,
                      elevation: 2,
                      shadowColor: Theme.of(context).colorScheme.primary.withValues(alpha: 0.3),
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
} 