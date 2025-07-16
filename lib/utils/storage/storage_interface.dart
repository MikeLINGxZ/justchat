import 'package:lemon_tea/models/conversation_v0.dart';

/// 存储接口，定义对话存储的标准操作
abstract class StorageInterface {
  /// 保存对话
  Future<void> saveConversation(Conversation_v0 conversation);

  /// 获取所有对话（不包括已删除的）
  Future<List<Conversation_v0>> getAllConversations();

  /// 根据ID获取对话
  Future<Conversation_v0?> getConversationById(String id);

  /// 删除对话（软删除）
  Future<void> deleteConversation(String id);

  /// 永久删除对话
  Future<void> permanentlyDeleteConversation(String id);

  /// 清空所有对话
  Future<void> clearAllConversations();

  /// 获取存储统计信息
  Future<StorageStats> getStorageStats();
}

/// 存储统计信息
class StorageStats {
  final int totalConversations;
  final int activeConversations;
  final int deletedConversations;
  final int totalMessages;
  final DateTime lastUpdated;

  StorageStats({
    required this.totalConversations,
    required this.activeConversations,
    required this.deletedConversations,
    required this.totalMessages,
    required this.lastUpdated,
  });
} 