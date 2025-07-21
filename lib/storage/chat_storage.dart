import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/storage/sqlite_util.dart';
import 'package:flutter/foundation.dart';
import 'package:uuid/uuid.dart';

/// 聊天存储类，负责对话和消息的数据库操作
class ChatStorage {
  static const _uuid = Uuid();

  /// 获取所有对话列表
  /// 按更新时间倒序排序
  static Future<List<Conversation>> getAllConversations() async {
    try {
      final results = await SqliteUtil.instance.query(
        Conversation.tableName(),
        where: 'deleted = ?',
        whereArgs: [0],
        orderBy: 'updated_at DESC',
      );
      return results.map((map) => Conversation.fromMap(map)).toList();
    } catch (e) {
      debugPrint('获取所有对话失败: $e');
      return [];
    }
  }

  /// 通过ID获取对话
  static Future<Conversation?> getConversationById(String id) async {
    try {
      final results = await SqliteUtil.instance.query(
        Conversation.tableName(),
        where: 'id = ? AND deleted = ?',
        whereArgs: [id, 0],
      );
      
      if (results.isNotEmpty) {
        return Conversation.fromMap(results.first);
      }
      return null;
    } catch (e) {
      debugPrint('通过ID获取对话失败: $e');
      return null;
    }
  }

  /// 创建新对话
  static Future<Conversation?> createConversation({
    required String title,
    String? defaultProviderId,
    String? defaultModelId,
  }) async {
    try {
      final now = DateTime.now();
      final conversation = Conversation(
        id: _uuid.v4(),
        title: title,
        createdAt: now,
        updatedAt: now,
        defaultProviderId: defaultProviderId,
        defaultModelId: defaultModelId,
      );
      
      final result = await SqliteUtil.instance.insert(
        Conversation.tableName(),
        conversation.toMap(),
      );
      
      if (result > 0) {
        return conversation;
      }
      return null;
    } catch (e) {
      debugPrint('创建对话失败: $e');
      return null;
    }
  }

  /// 更新对话
  static Future<bool> updateConversation(Conversation conversation) async {
    try {
      conversation.updatedAt = DateTime.now();
      final result = await SqliteUtil.instance.update(
        Conversation.tableName(),
        conversation.toMap(),
        where: 'id = ?',
        whereArgs: [conversation.id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新对话失败: $e');
      return false;
    }
  }

  /// 更新对话标题
  static Future<bool> updateConversationTitle(String conversationId, String title) async {
    try {
      final result = await SqliteUtil.instance.update(
        Conversation.tableName(),
        {
          'title': title,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        where: 'id = ?',
        whereArgs: [conversationId],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新对话标题失败: $e');
      return false;
    }
  }

  /// 软删除对话
  static Future<bool> deleteConversation(String id) async {
    try {
      final result = await SqliteUtil.instance.update(
        Conversation.tableName(),
        {
          'deleted': 1,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('删除对话失败: $e');
      return false;
    }
  }

  /// 获取对话的所有消息
  /// 按创建时间正序排序
  static Future<List<Message>> getMessagesByConversationId(String conversationId) async {
    try {
      final results = await SqliteUtil.instance.query(
        Message.tableName(),
        where: 'conversation_id = ? AND deleted = ?',
        whereArgs: [conversationId, 0],
        orderBy: 'created_at ASC',
      );
      return results.map((map) => Message.fromMap(map)).toList();
    } catch (e) {
      debugPrint('获取对话消息失败: $e');
      return [];
    }
  }

  /// 添加消息到对话
  static Future<Message?> addMessage({
    required String conversationId,
    required String role,
    required String content,
  }) async {
    try {
      final now = DateTime.now();
      final message = Message(
        conversation_id: conversationId,
        id: _uuid.v4(),
        role: role == 'user' ? MessageRole.user : MessageRole.assistant,
        content: content,
        createdAt: now,
      );
      
      final result = await SqliteUtil.instance.insert(
        Message.tableName(),
        message.toMap(),
      );
      
      if (result > 0) {
        // 更新对话的最后更新时间
        await SqliteUtil.instance.update(
          Conversation.tableName(),
          {'updated_at': now.millisecondsSinceEpoch},
          where: 'id = ?',
          whereArgs: [conversationId],
        );
        return message;
      }
      return null;
    } catch (e) {
      debugPrint('添加消息失败: $e');
      return null;
    }
  }

  /// 更新消息
  static Future<bool> updateMessage(Message message) async {
    try {
      final result = await SqliteUtil.instance.update(
        Message.tableName(),
        message.toMap(),
        where: 'id = ?',
        whereArgs: [message.id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新消息失败: $e');
      return false;
    }
  }

  /// 软删除消息
  static Future<bool> deleteMessage(String id) async {
    try {
      final result = await SqliteUtil.instance.update(
        Message.tableName(),
        {'deleted': 1},
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('删除消息失败: $e');
      return false;
    }
  }

  /// 清空对话的所有消息
  static Future<bool> clearConversationMessages(String conversationId) async {
    try {
      final result = await SqliteUtil.instance.update(
        Message.tableName(),
        {'deleted': 1},
        where: 'conversation_id = ?',
        whereArgs: [conversationId],
      );
      return result >= 0;
    } catch (e) {
      debugPrint('清空对话消息失败: $e');
      return false;
    }
  }

  /// 获取对话的消息数量
  static Future<int> getMessageCountByConversationId(String conversationId) async {
    try {
      final result = await SqliteUtil.instance.count(
        Message.tableName(),
        where: 'conversation_id = ? AND deleted = ?',
        whereArgs: [conversationId, 0],
      );
      return result;
    } catch (e) {
      debugPrint('获取对话消息数量失败: $e');
      return 0;
    }
  }
} 