// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tool_function.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ToolFunction _$ToolFunctionFromJson(Map<String, dynamic> json) => ToolFunction(
  name: json['name'] as String,
  description: json['description'] as String,
  parameters: ToolFunctionParameter.fromJson(
    json['parameters'] as Map<String, dynamic>,
  ),
);

Map<String, dynamic> _$ToolFunctionToJson(ToolFunction instance) =>
    <String, dynamic>{
      'name': instance.name,
      'description': instance.description,
      'parameters': instance.parameters,
    };
