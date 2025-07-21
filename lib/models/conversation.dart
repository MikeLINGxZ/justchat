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

  /// 对话默认使用的供应商
  final String? defaultProviderId;

  /// 对话默认使用的模型
  final String? defaultModelId;

  Conversation({
    required this.id,
    required this.title,
    required this.createdAt,
    required this.updatedAt,
    this.deleted = false,
    this.defaultModelId,
    this.defaultProviderId,
  });

  /// 表名
  static String tableName() {
    return 'conversations';
  }

  /// 表创建sql
  static String createTableSql() {
    return '''
      CREATE TABLE ${tableName()} (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        deleted INTEGER NOT NULL DEFAULT 0,
        model_id TEXT,
        provider_id TEXT,
        system_prompt TEXT,
        default_provider_id TEXT,
        default_model_id TEXT,
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
      'default_provider_id': defaultProviderId,
      'default_model_id': defaultModelId
    };
  }

  /// 从数据库 Map 构造对象
  factory Conversation.fromMap(Map<String, dynamic> map) {
    return Conversation(
      id: map['id'],
      title: map['title'],
      createdAt: DateTime.fromMillisecondsSinceEpoch(map['created_at']),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(map['updated_at']),
      deleted: (map['deleted'] ?? 0) == 1,
      defaultProviderId: map['default_provider_id'],
      defaultModelId: map['default_model_id'],
    );
  }

  /// 创建当前对象的副本，可选择性地替换指定字段
  Conversation copyWith({
    String? id,
    String? title,
    DateTime? createdAt,
    DateTime? updatedAt,
    bool? deleted,
    String? defaultProviderId,
    String? defaultModelId,
  }) {
    return Conversation(
      id: id ?? this.id,
      title: title ?? this.title,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      deleted: deleted ?? this.deleted,
      defaultProviderId: defaultProviderId ?? this.defaultProviderId,
      defaultModelId: defaultModelId ?? this.defaultModelId,
    );
  }
}