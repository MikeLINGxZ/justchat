// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'llm_provider.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LlmProvider _$LlmProviderFromJson(Map<String, dynamic> json) => LlmProvider(
  name: json['name'] as String,
  baseUrl: json['baseUrl'] as String,
  apiKey: json['apiKey'] as String?,
  alias: json['alias'] as String?,
  description: json['description'] as String?,
  models:
      (json['models'] as List<dynamic>?)
          ?.map((e) => Model.fromJson(e as Map<String, dynamic>))
          .toList(),
);

Map<String, dynamic> _$LlmProviderToJson(LlmProvider instance) =>
    <String, dynamic>{
      'name': instance.name,
      'baseUrl': instance.baseUrl,
      'apiKey': instance.apiKey,
      'alias': instance.alias,
      'description': instance.description,
      'models': instance.models,
    };
