import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/storage/chat_storage.dart';

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
              style: const TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        content: Padding(
          padding: const EdgeInsets.only(top: 8.0),
          child: Text(
            S.of(context).confirmDeleteConversation,
            style: TextStyle(
              fontSize: 14,
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
            child: Text(S.of(context).cancel),
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
            child: Text(S.of(context).delete),
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
            fontSize: 14,
          ),
          prefixIcon: Icon(
            Icons.search_rounded,
            color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
            size: 20,
          ),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        ),
        style: const TextStyle(fontSize: 14),
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
                fontSize: 18,
                fontWeight: FontWeight.w500,
                color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.8),
              ),
            ),
            if (_searchQuery.isEmpty) ...[
              const SizedBox(height: 8),
              Text(
                '开始一个新的对话吧',
                style: TextStyle(
                  fontSize: 14,
                  color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
                ),
              ),
              const SizedBox(height: 32),
              if (widget.onNewConversation != null)
                ElevatedButton.icon(
                  onPressed: widget.onNewConversation,
                  icon: const Icon(Icons.add_rounded, size: 18),
                  label: Text(S.of(context).newConversation),
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
              onTap: () => widget.onConversationSelected(conversation),
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
                              fontSize: 15,
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
                              fontSize: 13,
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
                                    fontSize: 11,
                                    fontWeight: FontWeight.w500,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                              ),
                              const Spacer(),
                              Text(
                                _formatDate(conversation.updatedAt),
                                style: TextStyle(
                                  fontSize: 11,
                                  color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.5),
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                    
                    // 删除按钮
                    const SizedBox(width: 8),
                    AnimatedOpacity(
                      opacity: isHovered ? 1.0 : 0.0,
                      duration: const Duration(milliseconds: 200),
                      child: IconButton(
                        onPressed: () => _handleDeleteConversation(conversation.id),
                        icon: const Icon(Icons.delete_outline_rounded, size: 18),
                        tooltip: S.of(context).deleteConversation,
                        style: IconButton.styleFrom(
                          foregroundColor: Colors.red[600],
                          backgroundColor: Colors.red[50],
                          padding: const EdgeInsets.all(8),
                          minimumSize: const Size(32, 32),
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
                      label: Text(S.of(context).newConversation),
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