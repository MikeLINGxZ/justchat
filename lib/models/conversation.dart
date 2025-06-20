import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';

part 'conversation.g.dart';

@JsonSerializable()
class Conversation {
  /// 对话唯一标识
  final String id;

  /// 对话标题
  String title;

  /// 消息列表
  final List<Message> messages;

  /// 创建时间
  final DateTime createdAt;

  /// 最后更新时间
  DateTime updatedAt;

  /// 是否已删除
  bool isDeleted;

  Conversation({
    required this.id,
    required this.title,
    required this.messages,
    required this.createdAt,
    required this.updatedAt,
    this.isDeleted = false,
  });

  factory Conversation.fromJson(Map<String, dynamic> json) =>
      _$ConversationFromJson(json);
  Map<String, dynamic> toJson() => _$ConversationToJson(this);

  /// 创建新对话
  factory Conversation.create({
    required String title,
    List<Message> messages = const [],
  }) {
    final now = DateTime.now();
    return Conversation(
      id: _generateId(),
      title: title,
      messages: messages,
      createdAt: now,
      updatedAt: now,
    );
  }

  /// 添加消息
  Conversation addMessage(Message message) {
    final newMessages = List<Message>.from(messages)..add(message);
    return copyWith(
      messages: newMessages,
      updatedAt: DateTime.now(),
    );
  }

  /// 更新标题
  Conversation updateTitle(String newTitle) {
    return copyWith(
      title: newTitle,
      updatedAt: DateTime.now(),
    );
  }

  /// 标记为删除
  Conversation markAsDeleted() {
    return copyWith(isDeleted: true);
  }

  /// 复制并修改
  Conversation copyWith({
    String? id,
    String? title,
    List<Message>? messages,
    DateTime? createdAt,
    DateTime? updatedAt,
    bool? isDeleted,
  }) {
    return Conversation(
      id: id ?? this.id,
      title: title ?? this.title,
      messages: messages ?? this.messages,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      isDeleted: isDeleted ?? this.isDeleted,
    );
  }

  /// 生成唯一ID
  static String _generateId() {
    return DateTime.now().millisecondsSinceEpoch.toString() +
        (1000 + (DateTime.now().microsecond % 9000)).toString();
  }
} 