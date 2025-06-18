// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tool_call_function.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ToolCallFunction _$ToolCallFunctionFromJson(Map<String, dynamic> json) =>
    ToolCallFunction(
      name: json['name'] as String,
      arguments: json['arguments'] as String,
    );

Map<String, dynamic> _$ToolCallFunctionToJson(ToolCallFunction instance) =>
    <String, dynamic>{'name': instance.name, 'arguments': instance.arguments};
