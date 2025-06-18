// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'tool_function_parameter.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ToolFunctionParameter _$ToolFunctionParameterFromJson(
  Map<String, dynamic> json,
) => ToolFunctionParameter(
  type: json['type'] as String,
  properties: (json['properties'] as Map<String, dynamic>).map(
    (k, e) =>
        MapEntry(k, ToolFunctionPropertie.fromJson(e as Map<String, dynamic>)),
  ),
  required:
      (json['required'] as List<dynamic>?)?.map((e) => e as String).toList(),
);

Map<String, dynamic> _$ToolFunctionParameterToJson(
  ToolFunctionParameter instance,
) => <String, dynamic>{
  'type': instance.type,
  'properties': instance.properties,
  'required': instance.required,
};
