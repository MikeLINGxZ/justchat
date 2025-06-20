import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/storage/storage_interface.dart';
import 'package:lemon_tea/utils/storage/local_storage.dart';

/// 对话管理器，负责管理对话的创建、保存、加载等操作
class ConversationManager extends ChangeNotifier {
  final StorageInterface _storage;
  
  List<Conversation> _conversations = [];
  Conversation? _currentConversation;
  bool _isLoading = false;

  ConversationManager({StorageInterface? storage}) 
      : _storage = storage ?? LocalStorage();

  /// 获取所有对话
  List<Conversation> get conversations => List.unmodifiable(_conversations);
  
  /// 获取当前对话
  Conversation? get currentConversation => _currentConversation;
  
  /// 是否正在加载
  bool get isLoading => _isLoading;

  /// 初始化，加载所有对话
  Future<void> initialize() async {
    _setLoading(true);
    try {
      _conversations = await _storage.getAllConversations();
      notifyListeners();
    } catch (e) {
      debugPrint('Failed to load conversations: $e');
    } finally {
      _setLoading(false);
    }
  }

  /// 创建新对话
  Future<Conversation> createConversation({
    String title = '新对话',
    List<Message> initialMessages = const [],
  }) async {
    final conversation = Conversation.create(
      title: title,
      messages: initialMessages,
    );
    
    await _storage.saveConversation(conversation);
    _conversations.insert(0, conversation);
    _currentConversation = conversation;
    
    notifyListeners();
    return conversation;
  }

  /// 加载对话
  Future<void> loadConversation(String id) async {
    final conversation = await _storage.getConversationById(id);
    if (conversation != null) {
      _currentConversation = conversation;
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

    final updatedConversation = _currentConversation!.addMessage(message);
    await _storage.saveConversation(updatedConversation);
    
    // 更新当前对话和列表中的对话
    _currentConversation = updatedConversation;
    final index = _conversations.indexWhere((c) => c.id == updatedConversation.id);
    if (index != -1) {
      _conversations[index] = updatedConversation;
      // 将更新的对话移到列表顶部
      _conversations.removeAt(index);
      _conversations.insert(0, updatedConversation);
    }
    
    notifyListeners();
  }

  /// 更新对话标题
  Future<void> updateConversationTitle(String id, String newTitle) async {
    final conversation = await _storage.getConversationById(id);
    if (conversation != null) {
      final updatedConversation = conversation.updateTitle(newTitle);
      await _storage.saveConversation(updatedConversation);
      
      // 更新列表中的对话
      final index = _conversations.indexWhere((c) => c.id == id);
      if (index != -1) {
        _conversations[index] = updatedConversation;
      }
      
      // 如果是当前对话，也要更新
      if (_currentConversation?.id == id) {
        _currentConversation = updatedConversation;
      }
      
      notifyListeners();
    }
  }

  /// 删除对话
  Future<void> deleteConversation(String id) async {
    await _storage.deleteConversation(id);
    
    // 从列表中移除
    _conversations.removeWhere((c) => c.id == id);
    
    // 如果删除的是当前对话，清空当前对话
    if (_currentConversation?.id == id) {
      _currentConversation = null;
    }
    
    notifyListeners();
  }

  /// 永久删除对话
  Future<void> permanentlyDeleteConversation(String id) async {
    await _storage.permanentlyDeleteConversation(id);
    
    // 从列表中移除
    _conversations.removeWhere((c) => c.id == id);
    
    // 如果删除的是当前对话，清空当前对话
    if (_currentConversation?.id == id) {
      _currentConversation = null;
    }
    
    notifyListeners();
  }

  /// 清空当前对话
  void clearCurrentConversation() {
    _currentConversation = null;
    notifyListeners();
  }

  /// 获取存储统计信息
  Future<StorageStats> getStorageStats() async {
    return await _storage.getStorageStats();
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
      
      // 搜索消息内容
      for (final message in conversation.messages) {
        if (message.content.toLowerCase().contains(lowercaseQuery)) {
          return true;
        }
      }
      
      return false;
    }).toList();
  }
} 