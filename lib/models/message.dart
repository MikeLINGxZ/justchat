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
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        metadata TEXT,
        FOREIGN KEY (conversation_id) REFERENCES conversations (id) ON DELETE CASCADE
      )
    ''';
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
      'created_at': createdAt.millisecondsSinceEpoch,
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
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at']),
      deleted: map['deleted'] == 1,
    );
  }
}