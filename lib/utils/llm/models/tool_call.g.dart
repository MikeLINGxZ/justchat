// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tool_call.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ToolCall _$ToolCallFromJson(Map<String, dynamic> json) => ToolCall(
  index: (json['index'] as num).toInt(),
  id: json['id'] as String,
  type: json['type'] as String,
  function: ToolCallFunction.fromJson(json['function'] as Map<String, dynamic>),
);

Map<String, dynamic> _$ToolCallToJson(ToolCall instance) => <String, dynamic>{
  'index': instance.index,
  'id': instance.id,
  'type': instance.type,
  'function': instance.function,
};
