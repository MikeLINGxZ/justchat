import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/storage/sqlite_util.dart';
import 'package:flutter/foundation.dart';

class LlmStorage {
  // 1. 获取所有llm_provider
  static Future<List<LlmProvider>> getAllProviders() async {
    try {
      final results = await SqliteUtil.instance.query(
        LlmProvider.tableName(),
        orderBy: 'seq_id ASC, name ASC', // 按seq_id升序排序，相同时按名称排序
      );
      return results.map((map) => LlmProvider.fromMap(map)).toList();
    } catch (e) {
      debugPrint('获取所有LLM提供商失败: $e');
      return [];
    }
  }

  // 2. 通过id获取llm_provider
  static Future<LlmProvider?> getProviderById(String id) async {
    try {
      final results = await SqliteUtil.instance.query(
        LlmProvider.tableName(),
        where: 'id = ?',
        whereArgs: [id],
      );
      
      if (results.isNotEmpty) {
        return LlmProvider.fromMap(results.first);
      }
      return null;
    } catch (e) {
      debugPrint('通过ID获取LLM提供商失败: $e');
      return null;
    }
  }

  // 3. 通过id更新llm_provider
  static Future<bool> updateProvider(LlmProvider provider) async {
    try {
      final result = await SqliteUtil.instance.update(
        LlmProvider.tableName(),
        provider.toMap(),
        where: 'id = ?',
        whereArgs: [provider.id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新LLM提供商失败: $e');
      return false;
    }
  }

  // 4. 通过id删除llm_provider
  static Future<bool> deleteProvider(String id) async {
    try {
      final result = await SqliteUtil.instance.delete(
        LlmProvider.tableName(),
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('删除LLM提供商失败: $e');
      return false;
    }
  }

  // 5. 添加llm_provider
  static Future<bool> addProvider(LlmProvider provider) async {
    try {
      final result = await SqliteUtil.instance.insert(LlmProvider.tableName(), provider.toMap());
      return result > 0;
    } catch (e) {
      debugPrint('添加LLM提供商失败: $e');
      return false;
    }
  }

  // 6. 通过llm_provider_id获取所有模型
  static Future<List<Model>> getModelsByProviderId(String providerId) async {
    try {
      final results = await SqliteUtil.instance.query(
        Model.tableName(),
        where: 'llm_provider_id = ?',
        whereArgs: [providerId],
        orderBy: 'seq_id ASC, id ASC', // 按seq_id升序排序，相同时按id排序
      );
      return results.map((map) => Model.fromMap(map)).toList();
    } catch (e) {
      debugPrint('获取提供商的所有模型失败: $e');
      return [];
    }
  }

  // 7. 通过id获取模型
  static Future<Model?> getModelById(String id) async {
    try {
      final results = await SqliteUtil.instance.query(
        Model.tableName(),
        where: 'id = ?',
        whereArgs: [id],
      );
      
      if (results.isNotEmpty) {
        return Model.fromMap(results.first);
      }
      return null;
    } catch (e) {
      debugPrint('通过ID获取模型失败: $e');
      return null;
    }
  }

  // 8. 通过id更新模型
  static Future<bool> updateModel(Model model) async {
    try {
      final result = await SqliteUtil.instance.update(
        Model.tableName(),
        model.toMap(),
        where: 'id = ?',
        whereArgs: [model.id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新模型失败: $e');
      return false;
    }
  }

  // 9. 通过id删除模型
  static Future<bool> deleteModel(String id) async {
    try {
      final result = await SqliteUtil.instance.delete(
        Model.tableName(),
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('删除模型失败: $e');
      return false;
    }
  }

  // 10. 添加模型
  static Future<bool> addModel(Model model) async {
    try {
      final result = await SqliteUtil.instance.insert(Model.tableName(), model.toMap());
      return result > 0;
    } catch (e) {
      debugPrint('添加模型失败: $e');
      return false;
    }
  }
  
  // 11. 添加带有自定义字段的模型
  static Future<bool> addModelWithCustomFields(Map<String, dynamic> modelMap) async {
    try {
      final result = await SqliteUtil.instance.insert(Model.tableName(), modelMap);
      return result > 0;
    } catch (e) {
      debugPrint('添加自定义模型失败: $e');
      return false;
    }
  }

  // 12. 更新模型序号
  static Future<bool> updateModelSeqId(String id, int seqId) async {
    try {
      final result = await SqliteUtil.instance.update(
        Model.tableName(),
        {'seq_id': seqId},
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新模型序号失败: $e');
      return false;
    }
  }

  // 13. 更新提供商序号
  static Future<bool> updateProviderSeqId(String id, int seqId) async {
    try {
      final result = await SqliteUtil.instance.update(
        LlmProvider.tableName(),
        {'seq_id': seqId},
        where: 'id = ?',
        whereArgs: [id],
      );
      return result > 0;
    } catch (e) {
      debugPrint('更新提供商序号失败: $e');
      return false;
    }
  }

  // 14. 获取最大模型序号
  static Future<int> getMaxModelSeqId(String providerId) async {
    try {
      final result = await SqliteUtil.instance.rawQuery(
        'SELECT MAX(seq_id) as max_seq_id FROM ${Model.tableName()} WHERE llm_provider_id = ?',
        [providerId],
      );
      return result.first['max_seq_id'] as int? ?? 0;
    } catch (e) {
      debugPrint('获取最大模型序号失败: $e');
      return 0;
    }
  }

  // 15. 获取最大提供商序号
  static Future<int> getMaxProviderSeqId() async {
    try {
      final result = await SqliteUtil.instance.rawQuery(
        'SELECT MAX(seq_id) as max_seq_id FROM ${LlmProvider.tableName()}',
      );
      return result.first['max_seq_id'] as int? ?? 0;
    } catch (e) {
      debugPrint('获取最大提供商序号失败: $e');
      return 0;
    }
  }

  // ==================== 对话历史相关功能 ====================

  // 16. 获取所有历史对话（按更新时间倒序）
  static Future<List<Map<String, dynamic>>> getAllConversations() async {
    try {
      final results = await SqliteUtil.instance.query(
        'conversations',
        where: 'deleted = ?',
        whereArgs: [0],
        orderBy: 'updated_at DESC',
      );
      return results;
    } catch (e) {
      debugPrint('获取所有对话失败: $e');
      return [];
    }
  }

  // 17. 获取对话的消息数量
  static Future<int> getConversationMessageCount(String conversationId) async {
    try {
      final result = await SqliteUtil.instance.count(
        'messages',
        where: 'conversation_id = ? AND deleted = ?',
        whereArgs: [conversationId, 0],
      );
      return result;
    } catch (e) {
      debugPrint('获取对话消息数量失败: $e');
      return 0;
    }
  }

  // 18. 获取对话的最新消息预览
  static Future<String> getConversationPreview(String conversationId) async {
    try {
      final results = await SqliteUtil.instance.query(
        'messages',
        where: 'conversation_id = ? AND deleted = ? AND role = ?',
        whereArgs: [conversationId, 0, 'user'],
        orderBy: 'created_at DESC',
        limit: 1,
      );
      
      if (results.isNotEmpty) {
        final content = results.first['content'] as String;
        // 限制预览长度，去除换行符
        final preview = content.replaceAll('\n', ' ').trim();
        if (preview.length > 50) {
          return '${preview.substring(0, 50)}...';
        }
        return preview.isEmpty ? '暂无内容' : preview;
      }
      return '暂无消息';
    } catch (e) {
      debugPrint('获取对话预览失败: $e');
      return '加载失败';
    }
  }

  // 新增：获取对话的详细预览信息（包含思考过程统计）
  static Future<Map<String, dynamic>> getConversationDetailedPreview(String conversationId) async {
    try {
      // 获取最新用户消息
      final userResults = await SqliteUtil.instance.query(
        'messages',
        where: 'conversation_id = ? AND deleted = ? AND role = ?',
        whereArgs: [conversationId, 0, 'user'],
        orderBy: 'created_at DESC',
        limit: 1,
      );
      
      String preview = '暂无消息';
      if (userResults.isNotEmpty) {
        final content = userResults.first['content'] as String;
        final cleanContent = content.replaceAll('\n', ' ').trim();
        if (cleanContent.length > 50) {
          preview = '${cleanContent.substring(0, 50)}...';
        } else {
          preview = cleanContent.isEmpty ? '暂无内容' : cleanContent;
        }
      }
      
      // 统计包含思考过程的消息数量
      final reasoningCountResult = await SqliteUtil.instance.rawQuery(
        'SELECT COUNT(*) as count FROM messages WHERE conversation_id = ? AND deleted = ? AND reasoning_content IS NOT NULL AND reasoning_content != ""',
        [conversationId, 0],
      );
      
      final reasoningCount = reasoningCountResult.first['count'] as int? ?? 0;
      
      // 调试输出
      debugPrint('getConversationDetailedPreview: conversationId=$conversationId, reasoningCount=$reasoningCount');
      
      return {
        'preview': preview,
        'hasReasoning': reasoningCount > 0,
        'reasoningCount': reasoningCount,
      };
    } catch (e) {
      debugPrint('获取对话详细预览失败: $e');
      return {
        'preview': '加载失败',
        'hasReasoning': false,
        'reasoningCount': 0,
      };
    }
  }

  // 19. 获取对话详细信息（包含消息数量和预览）
  static Future<Map<String, dynamic>?> getConversationWithDetails(String conversationId) async {
    try {
      final conversationResults = await SqliteUtil.instance.query(
        'conversations',
        where: 'id = ? AND deleted = ?',
        whereArgs: [conversationId, 0],
      );
      
      if (conversationResults.isEmpty) {
        return null;
      }
      
      final conversation = conversationResults.first;
      final messageCount = await getConversationMessageCount(conversationId);
      final preview = await getConversationPreview(conversationId);
      
      return {
        ...conversation,
        'message_count': messageCount,
        'preview': preview,
      };
    } catch (e) {
      debugPrint('获取对话详细信息失败: $e');
      return null;
    }
  }

  // 20. 搜索对话（根据标题）
  static Future<List<Map<String, dynamic>>> searchConversations(String query) async {
    try {
      final results = await SqliteUtil.instance.query(
        'conversations',
        where: 'deleted = ? AND title LIKE ?',
        whereArgs: [0, '%$query%'],
        orderBy: 'updated_at DESC',
      );
      return results;
    } catch (e) {
      debugPrint('搜索对话失败: $e');
      return [];
    }
  }
}