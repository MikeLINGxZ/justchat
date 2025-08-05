// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'message.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FileContent _$FileContentFromJson(Map<String, dynamic> json) => FileContent(
  name: json['name'] as String,
  mimeType: json['mimeType'] as String,
  type: json['type'] as String,
  data: json['data'] as String?,
  size: (json['size'] as num).toInt(),
  url: json['url'] as String?,
  description: json['description'] as String?,
  localPath: json['localPath'] as String?,
);

Map<String, dynamic> _$FileContentToJson(FileContent instance) =>
    <String, dynamic>{
      'name': instance.name,
      'mimeType': instance.mimeType,
      'type': instance.type,
      'data': instance.data,
      'size': instance.size,
      'url': instance.url,
      'description': instance.description,
      'localPath': instance.localPath,
    };

Message _$MessageFromJson(Map<String, dynamic> json) => Message(
  role: $enumDecode(_$MessageRoleEnumMap, json['role']),
  content: json['content'] as String,
  reasoningContent: json['reasoningContent'] as String?,
  toolCalls:
      (json['toolCalls'] as List<dynamic>?)
          ?.map((e) => ToolCall.fromJson(e as Map<String, dynamic>))
          .toList(),
  toolCallId: json['toolCallId'] as String?,
  files:
      (json['files'] as List<dynamic>?)
          ?.map((e) => FileContent.fromJson(e as Map<String, dynamic>))
          .toList(),
  stoppedByUser: json['stoppedByUser'] as bool? ?? false,
);

Map<String, dynamic> _$MessageToJson(Message instance) => <String, dynamic>{
  'role': _$MessageRoleEnumMap[instance.role]!,
  'content': instance.content,
  'reasoningContent': instance.reasoningContent,
  'toolCalls': instance.toolCalls,
  'toolCallId': instance.toolCallId,
  'files': instance.files,
  'stoppedByUser': instance.stoppedByUser,
};

const _$MessageRoleEnumMap = {
  MessageRole.system: 'system',
  MessageRole.user: 'user',
  MessageRole.assistant: 'assistant',
  MessageRole.tool: 'tool',
};
