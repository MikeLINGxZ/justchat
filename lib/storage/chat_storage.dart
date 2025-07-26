import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/message.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/storage/sqlite_util.dart';
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

  /// 更新对话的默认模型配置
  static Future<bool> updateConversationModel(String conversationId, String providerId, String modelId) async {
    try {
      final result = await SqliteUtil.instance.update(
        Conversation.tableName(),
        {
          'default_provider_id': providerId,
          'default_model_id': modelId,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
        },
        where: 'id = ?',
        whereArgs: [conversationId],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新对话模型配置失败: $e');
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
    String? reasoningContent,
  }) async {
    try {
      final now = DateTime.now();
      final message = Message(
        conversation_id: conversationId,
        id: _uuid.v4(),
        role: role == 'user' ? MessageRole.user : MessageRole.assistant,
        content: content,
        reasoningContent: reasoningContent,
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

  /// 根据内容和角色删除消息（用于重新生成和删除功能）
  static Future<bool> deleteMessagesByContent(
    String conversationId,
    String content,
    String role,
  ) async {
    try {
      final result = await SqliteUtil.instance.update(
        Message.tableName(),
        {'deleted': 1},
        where: 'conversation_id = ? AND content = ? AND role = ? AND deleted = ?',
        whereArgs: [conversationId, content, role, 0],
      );
      return result > 0;
    } catch (e) {
      debugPrint('根据内容删除消息失败: $e');
      return false;
    }
  }

  /// 根据内容和角色查找消息（用于重新生成功能）
  static Future<Message?> findMessageByContent(
    String conversationId,
    String content,
    String role,
  ) async {
    try {
      final results = await SqliteUtil.instance.query(
        Message.tableName(),
        where: 'conversation_id = ? AND content = ? AND role = ? AND deleted = ?',
        whereArgs: [conversationId, content, role, 0],
        limit: 1,
      );
      
      if (results.isNotEmpty) {
        return Message.fromMap(results.first);
      }
      return null;
    } catch (e) {
      debugPrint('根据内容查找消息失败: $e');
      return null;
    }
  }

  /// 根据内容更新消息（用于重新生成功能）
  static Future<bool> updateMessageByContent(
    String conversationId,
    String oldContent,
    String role,
    String newContent,
    String? newReasoningContent,
  ) async {
    try {
      final result = await SqliteUtil.instance.update(
        Message.tableName(),
        {
          'content': newContent,
          'reasoning_content': newReasoningContent,
        },
        where: 'conversation_id = ? AND content = ? AND role = ? AND deleted = ?',
        whereArgs: [conversationId, oldContent, role, 0],
      );
      return result > 0;
    } catch (e) {
      debugPrint('根据内容更新消息失败: $e');
      return false;
    }
  }

  /// 转义FTS搜索查询，避免语法错误
  static String _escapeFtsQuery(String query) {
    // 移除可能导致FTS语法错误的特殊字符
    String escaped = query
        .replaceAll('"', '') // 移除引号
        .replaceAll("'", '') // 移除单引号
        .replaceAll('*', '') // 移除星号
        .replaceAll('(', '') // 移除括号
        .replaceAll(')', '')
        .replaceAll('[', '') // 移除方括号
        .replaceAll(']', '')
        .replaceAll('+', '') // 移除加号
        .replaceAll('-', '') // 移除减号
        .replaceAll(':', '') // 移除冒号
        .replaceAll('^', '') // 移除异或
        .replaceAll('!', '') // 移除感叹号
        .replaceAll('&', '') // 移除和号
        .replaceAll('|', '') // 移除竖线
        .replaceAll('~', '') // 移除波浪号
        .trim();
    
    // 如果转义后为空，返回原查询用于LIKE搜索
    if (escaped.isEmpty) {
      return query.trim();
    }
    
    // 将多个空格合并为单个空格
    escaped = escaped.replaceAll(RegExp(r'\s+'), ' ');
    
    return escaped;
  }

  /// 降级搜索: 当FTS搜索失败时使用LIKE搜索
  static Future<List<String>> _fallbackSearchConversationIds(String searchQuery) async {
    try {
      debugPrint('执行降级搜索: $searchQuery');
      
      final sql = '''
        SELECT DISTINCT conversation_id 
        FROM ${Message.tableName()} 
        WHERE deleted = 0 
        AND (content LIKE ? OR reasoning_content LIKE ?)
        LIMIT 50
      ''';
      
      final likeQuery = '%$searchQuery%';
      final results = await SqliteUtil.instance.rawQuery(sql, [likeQuery, likeQuery]);
      
      return results.map((result) => result['conversation_id'] as String).toList();
    } catch (e) {
      debugPrint('降级搜索也失败: $e');
      return [];
    }
  }

  /// 根据消息内容搜索对话ID列表
  /// 使用FTS全文搜索功能
  static Future<List<String>> searchConversationIdsByMessageContent(String searchQuery) async {
    try {
      if (searchQuery.trim().isEmpty) {
        return [];
      }

      // 转义查询字符串
      final escapedQuery = _escapeFtsQuery(searchQuery);
      debugPrint('原始查询: $searchQuery, 转义后: $escapedQuery');

      // 首先尝试FTS搜索
      try {
        final sql = '''
          SELECT DISTINCT m.conversation_id 
          FROM ${Message.tableName()} m
          INNER JOIN ${Message.tableName()}_fts fts ON m.id = fts.id
          WHERE fts MATCH ? 
          AND m.deleted = 0
          ORDER BY rank
          LIMIT 50
        ''';
        
        final results = await SqliteUtil.instance.rawQuery(sql, [escapedQuery]);
        
        final conversationIds = results.map((result) => result['conversation_id'] as String).toList();
        debugPrint('FTS搜索成功，找到 ${conversationIds.length} 个对话');
        
        return conversationIds;
      } catch (ftsError) {
        debugPrint('FTS搜索失败，尝试降级搜索: $ftsError');
        // FTS搜索失败，使用降级搜索
        return await _fallbackSearchConversationIds(searchQuery);
      }
    } catch (e) {
      debugPrint('搜索消息内容失败: $e');
      return [];
    }
  }

  /// 搜索包含特定内容的消息
  /// 返回搜索结果及其对话信息
  static Future<List<Map<String, dynamic>>> searchMessagesWithConversationInfo(String searchQuery) async {
    try {
      if (searchQuery.trim().isEmpty) {
        return [];
      }

      // 转义查询字符串
      final escapedQuery = _escapeFtsQuery(searchQuery);
      debugPrint('搜索消息详情 - 原始查询: $searchQuery, 转义后: $escapedQuery');

      // 首先尝试FTS搜索
      try {
        final sql = '''
          SELECT m.*, c.title as conversation_title 
          FROM ${Message.tableName()} m
          INNER JOIN ${Message.tableName()}_fts fts ON m.id = fts.id
          INNER JOIN ${Conversation.tableName()} c ON m.conversation_id = c.id
          WHERE fts MATCH ? 
          AND m.deleted = 0 
          AND c.deleted = 0
          ORDER BY rank, m.created_at DESC
          LIMIT 50
        ''';
        
        final results = await SqliteUtil.instance.rawQuery(sql, [escapedQuery]);
        debugPrint('FTS消息搜索成功，找到 ${results.length} 条消息');
        
        return results;
      } catch (ftsError) {
        debugPrint('FTS消息搜索失败，尝试降级搜索: $ftsError');
        
        // 降级到LIKE搜索
        final sql = '''
          SELECT m.*, c.title as conversation_title 
          FROM ${Message.tableName()} m
          INNER JOIN ${Conversation.tableName()} c ON m.conversation_id = c.id
          WHERE m.deleted = 0 
          AND c.deleted = 0
          AND (m.content LIKE ? OR m.reasoning_content LIKE ?)
          ORDER BY m.created_at DESC
          LIMIT 50
        ''';
        
        final likeQuery = '%$searchQuery%';
        final results = await SqliteUtil.instance.rawQuery(sql, [likeQuery, likeQuery]);
        debugPrint('降级消息搜索成功，找到 ${results.length} 条消息');
        
        return results;
      }
    } catch (e) {
      debugPrint('搜索消息和对话信息失败: $e');
      return [];
    }
  }
}