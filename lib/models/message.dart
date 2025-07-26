import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/models/message_role.dart';

part 'message.g.dart';

@JsonSerializable()
class Message {
  /// 对话id
  String conversation_id;

  /// 消息id
  final String id;

  /// 消息角色
  final MessageRole role;

  /// 消息内容（对于function角色可能是空字符串）
  final String content;

  /// 思考过程内容（思维链模型返回的思考过程）
  final String? reasoningContent;

  /// 创建时间
  final DateTime createdAt;

  /// 是否删除（默认 false）
  @JsonKey(defaultValue: false)
  bool deleted;

  Message({
    required this.conversation_id,
    required this.id,
    required this.role,
    required this.content,
    this.reasoningContent,
    required this.createdAt,
    this.deleted = false,
  });

  /// 表名
  static String tableName() {
    return 'messages';
  }

  /// 表创建sql
  static String createTableSql() {
    return '''
      CREATE TABLE ${tableName()} (
        id TEXT PRIMARY KEY,
        conversation_id TEXT NOT NULL,
        role TEXT NOT NULL,
        content TEXT NOT NULL,
        reasoning_content TEXT,
        created_at INTEGER NOT NULL,
        deleted INTEGER NOT NULL DEFAULT 0,
        metadata TEXT,
        FOREIGN KEY (conversation_id) REFERENCES conversations (id) ON DELETE CASCADE
      );
      
      -- 创建FTS虚拟表用于内容搜索
      CREATE VIRTUAL TABLE ${tableName()}_fts USING fts5(
        id UNINDEXED,
        conversation_id UNINDEXED,
        content,
        reasoning_content,
        content=${tableName()},
        content_rowid=rowid
      );
      
      -- 创建触发器，自动同步数据到FTS表
      CREATE TRIGGER ${tableName()}_fts_insert AFTER INSERT ON ${tableName()} BEGIN
        INSERT INTO ${tableName()}_fts(id, conversation_id, content, reasoning_content)
        VALUES (new.id, new.conversation_id, new.content, new.reasoning_content);
      END;
      
      CREATE TRIGGER ${tableName()}_fts_delete AFTER DELETE ON ${tableName()} BEGIN
        DELETE FROM ${tableName()}_fts WHERE id = old.id;
      END;
      
      CREATE TRIGGER ${tableName()}_fts_update AFTER UPDATE ON ${tableName()} BEGIN
        UPDATE ${tableName()}_fts 
        SET content = new.content, reasoning_content = new.reasoning_content
        WHERE id = new.id;
      END;
    ''';
  }

  /// 文本搜索SQL - 使用FTS进行全文搜索
  /// 
  /// 参数:
  /// - searchQuery: 搜索关键词
  /// - conversationId: 可选，限制在特定对话中搜索
  /// - limit: 可选，限制返回结果数量，默认50
  /// 
  /// 使用示例:
  /// ```dart
  /// // 在所有消息中搜索
  /// String sql = Message.searchContentSql("关键词");
  /// 
  /// // 在特定对话中搜索
  /// String sql = Message.searchContentSql("关键词", conversationId: "conv_id");
  /// 
  /// // 限制返回数量
  /// String sql = Message.searchContentSql("关键词", limit: 10);
  /// ```
  static String searchContentSql(String searchQuery, {String? conversationId, int limit = 50}) {
    String whereClause = '';
    if (conversationId != null) {
      whereClause = 'AND m.conversation_id = ?';
    }
    
    return '''
      SELECT m.* FROM ${tableName()} m
      INNER JOIN ${tableName()}_fts fts ON m.id = fts.id
      WHERE fts MATCH ? $whereClause
      AND m.deleted = 0
      ORDER BY rank
      LIMIT $limit
    ''';
  }

  /// 获取搜索SQL的参数列表
  /// 
  /// 使用示例:
  /// ```dart
  /// List<dynamic> params = Message.getSearchParams("关键词", conversationId: "conv_id");
  /// ```
  static List<dynamic> getSearchParams(String searchQuery, {String? conversationId}) {
    List<dynamic> params = [searchQuery];
    if (conversationId != null) {
      params.add(conversationId);
    }
    return params;
  }

  /// 从 JSON 创建 LlmProvider 实例
  factory Message.fromJson(Map<String, dynamic> json) => _$MessageFromJson(json);

  /// 转换为 JSON
  Map<String, dynamic> toJson() => _$MessageToJson(this);

  /// 转换为数据库 Map
  Map<String, dynamic> toMap() {
    return {
      'conversation_id': conversation_id,
      'id': id,
      'role': role.toString().split('.').last, // MessageRole.user => 'user'
      'content': content,
      'reasoning_content': reasoningContent,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': createdAt.millisecondsSinceEpoch, // 兼容现有表结构
      'deleted': deleted ? 1 : 0,
    };
  }

  /// 从数据库 Map 构造对象
  factory Message.fromMap(Map<String, dynamic> map) {
    return Message(
      conversation_id: map['conversation_id'],
      id: map['id'],
      role: MessageRole.values.byName(map['role']),
      content: map['content'],
      reasoningContent: map['reasoning_content'],
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at']),
      deleted: (map['deleted'] ?? 0) == 1,
    );
  }
}