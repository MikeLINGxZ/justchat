import 'package:json_annotation/json_annotation.dart';

part 'conversation.g.dart';

@JsonSerializable()
class Conversation {
  /// 对话id
  final String id;

  /// 对话标题
  final String title;

  /// 创建时间
  final DateTime createdAt;

  /// 最后更新时间
  DateTime updatedAt;

  /// 是否删除（默认 false）
  @JsonKey(defaultValue: false)
  bool deleted;

  Conversation({
    required this.id,
    required this.title,
    required this.createdAt,
    required this.updatedAt,
    this.deleted = false,
  });

  /// 表创建sql
  static String createTableSql() {
    return '''
      CREATE TABLE conversations (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        model_id TEXT,
        provider_id TEXT,
        system_prompt TEXT,
        metadata TEXT
      )
    ''';
  }

  /// 生成 fromJson 方法
  factory Conversation.fromJson(Map<String, dynamic> json) => _$ConversationFromJson(json);

  /// 生成 toJson 方法
  Map<String, dynamic> toJson() => _$ConversationToJson(this);

  /// 转换为数据库 Map
  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'title': title,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': updatedAt.millisecondsSinceEpoch,
      'deleted': deleted ? 1 : 0,
    };
  }

  /// 从数据库 Map 构造对象
  factory Conversation.fromMap(Map<String, dynamic> map) {
    return Conversation(
      id: map['id'],
      title: map['title'],
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at']),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(map['updated_at']),
      deleted: map['deleted'] == 1,
    );
  }
}