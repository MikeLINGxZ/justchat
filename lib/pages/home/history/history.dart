import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';

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

class _HistoryPageState extends ConsumerState<HistoryPage> {
  String _searchQuery = '';
  List<Conversation> _filteredConversations = [];

  @override
  void initState() {
    super.initState();
    // 监听 ConversationManager 的变化
    widget.conversationManager.addListener(_onConversationManagerChanged);
    _updateFilteredConversations();
  }

  @override
  void dispose() {
    widget.conversationManager.removeListener(_onConversationManagerChanged);
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，刷新UI
    if (mounted) {
      setState(() {
        _updateFilteredConversations();
      });
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
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: const Text('确定要删除这个对话吗？删除后无法恢复。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: TextButton.styleFrom(
              foregroundColor: Colors.red,
            ),
            child: const Text('删除'),
          ),
        ],
      ),
    );

    if (shouldDelete == true) {
      await widget.onConversationDeleted(conversationId);
    }
  }

  @override
  Widget build(BuildContext context) {
    final currentConversation = widget.conversationManager.currentConversation;
    
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 标题栏
          Row(
            children: [
              Text(
                '对话历史',
                style: TextStyle(
                  fontSize: FontSizeUtils.getHeadingSize(ref),
                  fontWeight: FontWeight.bold,
                ),
              ),
              const Spacer(),
              if (widget.onNewConversation != null)
                ElevatedButton.icon(
                  onPressed: widget.onNewConversation,
                  icon: const Icon(Icons.add),
                  label: const Text('新对话'),
                ),
            ],
          ),
          const SizedBox(height: 16),
          
          // 搜索框
          TextField(
            decoration: InputDecoration(
              hintText: '搜索对话...',
              prefixIcon: const Icon(Icons.search),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            ),
            onChanged: (value) {
              setState(() {
                _searchQuery = value;
                _updateFilteredConversations();
              });
            },
          ),
          const SizedBox(height: 16),
          
          // 对话列表
          Expanded(
            child: _filteredConversations.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          _searchQuery.isEmpty ? Icons.chat_bubble_outline : Icons.search_off,
                          size: 64,
                          color: Colors.grey[400],
                        ),
                        const SizedBox(height: 16),
                        Text(
                          _searchQuery.isEmpty ? '暂无对话历史' : '未找到相关对话',
                          style: TextStyle(
                            fontSize: 16,
                            color: Colors.grey[600],
                          ),
                        ),
                        if (_searchQuery.isEmpty) ...[
                          const SizedBox(height: 8),
                          Text(
                            '开始一个新的对话吧',
                            style: TextStyle(
                              fontSize: 14,
                              color: Colors.grey[500],
                            ),
                          ),
                        ],
                      ],
                    ),
                  )
                : ListView.builder(
                    itemCount: _filteredConversations.length,
                    itemBuilder: (context, index) {
                      final conversation = _filteredConversations[index];
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
                                onPressed: () => _handleDeleteConversation(conversation.id),
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
                            widget.onConversationSelected(conversation);
                          },
                        ),
                      );
                    },
                  ),
          ),
        ],
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