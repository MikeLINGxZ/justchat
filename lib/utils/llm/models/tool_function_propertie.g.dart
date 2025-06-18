// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tool_function_propertie.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ToolFunctionPropertie _$ToolFunctionPropertieFromJson(
  Map<String, dynamic> json,
) => ToolFunctionPropertie(
  type: json['type'] as String,
  description: json['description'] as String,
  enums: (json['enum'] as List<dynamic>?)?.map((e) => e as String).toList(),
);

Map<String, dynamic> _$ToolFunctionPropertieToJson(
  ToolFunctionPropertie instance,
) => <String, dynamic>{
  'type': instance.type,
  'description': instance.description,
  'enum': instance.enums,
};
