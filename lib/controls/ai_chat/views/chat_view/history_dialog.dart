import 'package:flutter/material.dart';
import 'package:lemon_tea/models/conversation_v0.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';

class HistoryDialog extends StatefulWidget {
  final ConversationManager conversationManager;
  final Function(Conversation_v0) onConversationSelected;
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
  @override
  void initState() {
    super.initState();
    // 监听 ConversationManager 的变化
    widget.conversationManager.addListener(_onConversationManagerChanged);
  }

  @override
  void dispose() {
    widget.conversationManager.removeListener(_onConversationManagerChanged);
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，刷新UI
    if (mounted) {
      setState(() {});
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
                                  _getPreviewText(conversation.messages),
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
                                  '${conversation.messages.length}条消息',
                                  style: const TextStyle(
                                    fontSize: 11,
                                    color: Colors.grey,
                                  ),
                                ),
                                const SizedBox(width: 8),
                                IconButton(
                                  onPressed: () => widget.onConversationDeleted(conversation.id),
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

  String _getPreviewText(List<Message> messages) {
    if (messages.isEmpty) return '空对话';
    
    // 获取最后一条用户消息作为预览
    for (int i = messages.length - 1; i >= 0; i--) {
      if (messages[i].role == MessageRole.user) {
        final content = messages[i].content;
        return content.length > 50 ? '${content.substring(0, 50)}...' : content;
      }
    }
    
    // 如果没有用户消息，使用第一条消息
    final content = messages.first.content;
    return content.length > 50 ? '${content.substring(0, 50)}...' : content;
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