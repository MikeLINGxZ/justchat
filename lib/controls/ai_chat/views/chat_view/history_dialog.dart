import 'package:flutter/material.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/storage/llm_storage.dart';

class HistoryDialog extends StatefulWidget {
  final ConversationManager conversationManager;
  final Function(Conversation) onConversationSelected;
  final Function(String) onConversationDeleted;
  final VoidCallback? onNewConversation;

  const HistoryDialog({
    super.key,
    required this.conversationManager,
    required this.onConversationSelected,
    required this.onConversationDeleted,
    this.onNewConversation,
  });

  @override
  State<HistoryDialog> createState() => _HistoryDialogState();
}

class _HistoryDialogState extends State<HistoryDialog> {
  Map<String, int> _messageCountCache = {};
  Map<String, String> _previewCache = {};
  bool _isLoadingDetails = false;

  @override
  void initState() {
    super.initState();
    // 监听 ConversationManager 的变化
    widget.conversationManager.addListener(_onConversationManagerChanged);
    // 异步加载对话详情
    _loadConversationDetails();
  }

  @override
  void dispose() {
    widget.conversationManager.removeListener(_onConversationManagerChanged);
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，刷新UI并重新加载详情
    if (mounted) {
      setState(() {});
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
          final messageCount = await LlmStorage.getConversationMessageCount(conversation.id);
          final preview = await LlmStorage.getConversationPreview(conversation.id);
          
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

  @override
  Widget build(BuildContext context) {
    final conversations = widget.conversationManager.conversations;
    final currentConversation = widget.conversationManager.currentConversation;
    
    return Dialog(
      child: Container(
        width: 600,
        height: 500,
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text(
                  '对话历史',
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Spacer(),
                if (_isLoadingDetails)
                  const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                const SizedBox(width: 8),
                if (widget.onNewConversation != null)
                  TextButton.icon(
                    onPressed: () {
                      Navigator.of(context).pop();
                      widget.onNewConversation!();
                    },
                    icon: const Icon(Icons.add),
                    label: const Text('新对话'),
                  ),
                const SizedBox(width: 8),
                IconButton(
                  onPressed: () => Navigator.of(context).pop(),
                  icon: const Icon(Icons.close),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: conversations.isEmpty
                  ? const Center(
                      child: Text(
                        '暂无对话历史',
                        style: TextStyle(
                          fontSize: 16,
                          color: Colors.grey,
                        ),
                      ),
                    )
                  : ListView.builder(
                      itemCount: conversations.length,
                      itemBuilder: (context, index) {
                        final conversation = conversations[index];
                        final isCurrent = currentConversation?.id == conversation.id;
                        final messageCount = _messageCountCache[conversation.id] ?? 0;
                        final preview = _previewCache[conversation.id] ?? '加载中...';
                        
                        return Card(
                          margin: const EdgeInsets.only(bottom: 8),
                          color: isCurrent 
                              ? Theme.of(context).colorScheme.primary.withValues(alpha: 0.1)
                              : null,
                          child: ListTile(
                            title: Text(
                              conversation.title,
                              style: TextStyle(
                                fontWeight: isCurrent ? FontWeight.bold : FontWeight.normal,
                              ),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  preview,
                                  maxLines: 2,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(fontSize: 12),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  _formatDate(conversation.updatedAt),
                                  style: const TextStyle(
                                    fontSize: 11,
                                    color: Colors.grey,
                                  ),
                                ),
                              ],
                            ),
                            trailing: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Text(
                                  '${messageCount}条消息',
                                  style: const TextStyle(
                                    fontSize: 11,
                                    color: Colors.grey,
                                  ),
                                ),
                                const SizedBox(width: 8),
                                IconButton(
                                  onPressed: () => _showDeleteConfirmDialog(context, conversation),
                                  icon: const Icon(Icons.delete_outline, size: 18),
                                  tooltip: '删除对话',
                                  style: IconButton.styleFrom(
                                    foregroundColor: Colors.red,
                                    backgroundColor: Colors.red.withValues(alpha: 0.1),
                                  ),
                                ),
                              ],
                            ),
                            onTap: () {
                              Navigator.of(context).pop();
                              widget.onConversationSelected(conversation);
                            },
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  /// 显示删除确认对话框
  void _showDeleteConfirmDialog(BuildContext context, Conversation conversation) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('确认删除'),
          content: Text('确定要删除对话"${conversation.title}"吗？此操作不可恢复。'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('取消'),
            ),
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
                widget.onConversationDeleted(conversation.id);
                // 从缓存中移除
                _messageCountCache.remove(conversation.id);
                _previewCache.remove(conversation.id);
              },
              style: TextButton.styleFrom(foregroundColor: Colors.red),
              child: const Text('删除'),
            ),
          ],
        );
      },
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