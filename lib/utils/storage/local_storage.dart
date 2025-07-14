import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/utils/storage/storage_interface.dart';

/// 本地存储实现，使用SharedPreferences
class LocalStorage implements StorageInterface {
  static const String _conversationsKey = 'conversations';
  static const String _deletedConversationsKey = 'deleted_conversations';
  
  late SharedPreferences _prefs;
  bool _initialized = false;

  /// 初始化存储
  Future<void> _ensureInitialized() async {
    if (!_initialized) {
      _prefs = await SharedPreferences.getInstance();
      _initialized = true;
    }
  }
  
  /// 保存整数值
  Future<bool> setInt(String key, int value) async {
    await _ensureInitialized();
    return await _prefs.setInt(key, value);
  }
  
  /// 获取整数值
  Future<int?> getInt(String key) async {
    await _ensureInitialized();
    return _prefs.getInt(key);
  }
  
  /// 保存字符串值
  Future<bool> setString(String key, String value) async {
    await _ensureInitialized();
    return await _prefs.setString(key, value);
  }
  
  /// 获取字符串值
  Future<String?> getString(String key) async {
    await _ensureInitialized();
    return _prefs.getString(key);
  }
  
  /// 保存布尔值
  Future<bool> setBool(String key, bool value) async {
    await _ensureInitialized();
    return await _prefs.setBool(key, value);
  }
  
  /// 获取布尔值
  Future<bool?> getBool(String key) async {
    await _ensureInitialized();
    return _prefs.getBool(key);
  }
  
  /// 删除指定键的值
  Future<bool> remove(String key) async {
    await _ensureInitialized();
    return await _prefs.remove(key);
  }

  @override
  Future<void> saveConversation(Conversation conversation) async {
    await _ensureInitialized();
    
    final conversations = await _getAllConversationsMap();
    conversations[conversation.id] = jsonEncode(conversation.toJson());
    
    await _prefs.setString(_conversationsKey, jsonEncode(conversations));
  }

  @override
  Future<List<Conversation>> getAllConversations() async {
    await _ensureInitialized();
    
    final conversations = await _getAllConversationsMap();
    final conversationList = conversations.values
        .map((json) => Conversation.fromJson(jsonDecode(json)))
        .where((conv) => !conv.isDeleted)
        .toList();
    
    // 按更新时间倒序排列
    conversationList.sort((a, b) => b.updatedAt.compareTo(a.updatedAt));
    
    return conversationList;
  }

  @override
  Future<Conversation?> getConversationById(String id) async {
    await _ensureInitialized();
    
    final conversations = await _getAllConversationsMap();
    final jsonString = conversations[id];
    
    if (jsonString != null) {
      final conversation = Conversation.fromJson(jsonDecode(jsonString));
      return conversation.isDeleted ? null : conversation;
    }
    
    return null;
  }

  @override
  Future<void> deleteConversation(String id) async {
    await _ensureInitialized();
    
    final conversation = await getConversationById(id);
    if (conversation != null) {
      final deletedConversation = conversation.markAsDeleted();
      await saveConversation(deletedConversation);
      
      // 保存到已删除列表
      final deletedConversations = await _getDeletedConversationsMap();
      deletedConversations[id] = jsonEncode(deletedConversation.toJson());
      await _prefs.setString(_deletedConversationsKey, jsonEncode(deletedConversations));
    }
  }

  @override
  Future<void> permanentlyDeleteConversation(String id) async {
    await _ensureInitialized();
    
    final conversations = await _getAllConversationsMap();
    conversations.remove(id);
    await _prefs.setString(_conversationsKey, jsonEncode(conversations));
    
    // 从已删除列表中也移除
    final deletedConversations = await _getDeletedConversationsMap();
    deletedConversations.remove(id);
    await _prefs.setString(_deletedConversationsKey, jsonEncode(deletedConversations));
  }

  @override
  Future<void> clearAllConversations() async {
    await _ensureInitialized();
    
    await _prefs.remove(_conversationsKey);
    await _prefs.remove(_deletedConversationsKey);
  }

  @override
  Future<StorageStats> getStorageStats() async {
    await _ensureInitialized();
    
    final conversations = await _getAllConversationsMap();
    final deletedConversations = await _getDeletedConversationsMap();
    
    int totalMessages = 0;
    int activeConversations = 0;
    
    for (final jsonString in conversations.values) {
      final conversation = Conversation.fromJson(jsonDecode(jsonString));
      if (!conversation.isDeleted) {
        activeConversations++;
        totalMessages += conversation.messages.length;
      }
    }
    
    return StorageStats(
      totalConversations: conversations.length,
      activeConversations: activeConversations,
      deletedConversations: deletedConversations.length,
      totalMessages: totalMessages,
      lastUpdated: DateTime.now(),
    );
  }

  /// 获取所有对话的Map（包括已删除的）
  Future<Map<String, String>> _getAllConversationsMap() async {
    final jsonString = _prefs.getString(_conversationsKey);
    if (jsonString != null) {
      final Map<String, dynamic> map = jsonDecode(jsonString);
      return map.map((key, value) => MapEntry(key, value.toString()));
    }
    return {};
  }

  /// 获取已删除对话的Map
  Future<Map<String, String>> _getDeletedConversationsMap() async {
    final jsonString = _prefs.getString(_deletedConversationsKey);
    if (jsonString != null) {
      final Map<String, dynamic> map = jsonDecode(jsonString);
      return map.map((key, value) => MapEntry(key, value.toString()));
    }
    return {};
  }

  /// 恢复已删除的对话
  Future<void> restoreConversation(String id) async {
    await _ensureInitialized();
    
    final deletedConversations = await _getDeletedConversationsMap();
    final jsonString = deletedConversations[id];
    
    if (jsonString != null) {
      final conversation = Conversation.fromJson(jsonDecode(jsonString));
      final restoredConversation = conversation.copyWith(isDeleted: false);
      await saveConversation(restoredConversation);
      
      // 从已删除列表中移除
      deletedConversations.remove(id);
      await _prefs.setString(_deletedConversationsKey, jsonEncode(deletedConversations));
    }
  }

  /// 获取已删除的对话列表
  Future<List<Conversation>> getDeletedConversations() async {
    await _ensureInitialized();
    
    final deletedConversations = await _getDeletedConversationsMap();
    final conversationList = deletedConversations.values
        .map((json) => Conversation.fromJson(jsonDecode(json)))
        .toList();
    
    // 按删除时间倒序排列
    conversationList.sort((a, b) => b.updatedAt.compareTo(a.updatedAt));
    
    return conversationList;
  }
} 