import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message.dart' as db_message;
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/storage/chat_storage.dart';

/// 对话管理器，负责管理对话的创建、保存、加载等操作
/// 现在使用新的SQLite存储系统
class ConversationManager extends ChangeNotifier {
  
  List<Conversation> _conversations = [];
  Conversation? _currentConversation;
  List<Message> _currentMessages = [];
  bool _isLoading = false;

  ConversationManager();

  /// 获取所有对话
  List<Conversation> get conversations => List.unmodifiable(_conversations);
  
  /// 获取当前对话
  Conversation? get currentConversation => _currentConversation;
  
  /// 获取当前对话的消息（兼容旧接口）
  List<Message> get messages => List.unmodifiable(_currentMessages);
  
  /// 是否正在加载
  bool get isLoading => _isLoading;

  /// 初始化，加载所有对话
  Future<void> initialize() async {
    _setLoading(true);
    try {
      _conversations = await ChatStorage.getAllConversations();
      notifyListeners();
    } catch (e) {
      debugPrint('Failed to load conversations: $e');
    } finally {
      _setLoading(false);
    }
  }

  /// 创建新对话
  Future<Conversation?> createConversation({
    String title = '新对话',
    List<Message> initialMessages = const [],
    String? defaultProviderId,
    String? defaultModelId,
  }) async {
    final conversation = await ChatStorage.createConversation(
      title: title,
      defaultProviderId: defaultProviderId ?? 'deepseek',
      defaultModelId: defaultModelId ?? 'deepseek-chat',
    );
    
    if (conversation != null) {
      // 如果有初始消息，保存它们
      for (final message in initialMessages) {
        await ChatStorage.addMessage(
          conversationId: conversation.id,
          role: message.role.toString().split('.').last,
          content: message.content,
          reasoningContent: message.reasoningContent,
        );
      }
      
      _conversations.insert(0, conversation);
      _currentConversation = conversation;
      _currentMessages = initialMessages;
      notifyListeners();
    }
    
    return conversation;
  }

  /// 加载对话
  Future<void> loadConversation(String id) async {
    final conversation = await ChatStorage.getConversationById(id);
    if (conversation != null) {
      _currentConversation = conversation;
      
      // 加载对话的消息并转换为LLM Message
      final dbMessages = await ChatStorage.getMessagesByConversationId(id);
      _currentMessages = dbMessages.map((dbMsg) => Message(
        role: dbMsg.role,
        content: dbMsg.content,
        reasoningContent: dbMsg.reasoningContent,
      )).toList();
      
      notifyListeners();
    }
  }

  /// 添加消息到当前对话
  Future<void> addMessageToCurrent(Message message) async {
    if (_currentConversation == null) {
      // 如果没有当前对话，创建一个新的
      await createConversation(initialMessages: [message]);
      return;
    }

    // 保存到数据库
    await ChatStorage.addMessage(
      conversationId: _currentConversation!.id,
      role: message.role.toString().split('.').last,
      content: message.content,
      reasoningContent: message.reasoningContent,
    );
    
    // 更新内存中的消息列表
    _currentMessages.add(message);
    
    // 更新对话的最后更新时间
    _currentConversation!.updatedAt = DateTime.now();
    
    notifyListeners();
  }

  /// 更新对话标题
  Future<void> updateConversationTitle(String id, String newTitle) async {
    final success = await ChatStorage.updateConversationTitle(id, newTitle);
    if (success) {
      // 更新列表中的对话
      final index = _conversations.indexWhere((c) => c.id == id);
      if (index != -1) {
        _conversations[index] = Conversation(
          id: _conversations[index].id,
          title: newTitle,
          createdAt: _conversations[index].createdAt,
          updatedAt: DateTime.now(),
          defaultProviderId: _conversations[index].defaultProviderId,
          defaultModelId: _conversations[index].defaultModelId,
        );
      }
      
      // 如果是当前对话，也要更新
      if (_currentConversation?.id == id) {
        _currentConversation = Conversation(
          id: _currentConversation!.id,
          title: newTitle,
          createdAt: _currentConversation!.createdAt,
          updatedAt: DateTime.now(),
          defaultProviderId: _currentConversation!.defaultProviderId,
          defaultModelId: _currentConversation!.defaultModelId,
        );
      }
      
      notifyListeners();
    }
  }

  /// 删除对话
  Future<void> deleteConversation(String id) async {
    final success = await ChatStorage.deleteConversation(id);
    if (success) {
      // 从列表中移除
      _conversations.removeWhere((c) => c.id == id);
      
      // 如果删除的是当前对话，清空当前对话
      if (_currentConversation?.id == id) {
        _currentConversation = null;
        _currentMessages.clear();
      }
      
      notifyListeners();
    }
  }

  /// 清空当前对话
  void clearCurrentConversation() {
    _currentConversation = null;
    _currentMessages.clear();
    notifyListeners();
  }

  /// 设置加载状态
  void _setLoading(bool loading) {
    _isLoading = loading;
    notifyListeners();
  }

  /// 根据消息内容生成对话标题
  String generateTitleFromMessage(String message) {
    // 先处理多行消息，取第一行
    final firstLine = message.split('\n').first;
    
    // 如果第一行长度不超过20，直接返回
    if (firstLine.length <= 20) {
      return firstLine;
    }
    
    // 如果第一行超过20个字符，截断并添加省略号
    return '${firstLine.substring(0, 20)}...';
  }

  /// 搜索对话
  List<Conversation> searchConversations(String query) {
    if (query.isEmpty) {
      return _conversations;
    }
    
    final lowercaseQuery = query.toLowerCase();
    return _conversations.where((conversation) {
      // 搜索标题
      if (conversation.title.toLowerCase().contains(lowercaseQuery)) {
        return true;
      }
      
      return false;
    }).toList();
  }
} 