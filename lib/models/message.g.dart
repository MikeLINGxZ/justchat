// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'message.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Message _$MessageFromJson(Map<String, dynamic> json) => Message(
  conversation_id: json['conversation_id'] as String,
  id: json['id'] as String,
  role: $enumDecode(_$MessageRoleEnumMap, json['role']),
  content: json['content'] as String,
  createdAt: DateTime.parse(json['createdAt'] as String),
  deleted: json['deleted'] as bool? ?? false,
);

Map<String, dynamic> _$MessageToJson(Message instance) => <String, dynamic>{
  'conversation_id': instance.conversation_id,
  'id': instance.id,
  'role': _$MessageRoleEnumMap[instance.role]!,
  'content': instance.content,
  'createdAt': instance.createdAt.toIso8601String(),
  'deleted': instance.deleted,
};

const _$MessageRoleEnumMap = {
  MessageRole.system: 'system',
  MessageRole.user: 'user',
  MessageRole.assistant: 'assistant',
  MessageRole.tool: 'tool',
};
