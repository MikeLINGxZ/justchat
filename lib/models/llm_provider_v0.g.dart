// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'llm_provider_v0.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LlmProvider_v0 _$LlmProvider_v0FromJson(Map<String, dynamic> json) =>
    LlmProvider_v0(
      name: json['name'] as String,
      baseUrl: json['baseUrl'] as String,
      apiKey: json['apiKey'] as String?,
      alias: json['alias'] as String?,
      description: json['description'] as String?,
      models:
          (json['models'] as List<dynamic>?)
              ?.map((e) => Model_v0.fromJson(e as Map<String, dynamic>))
              .toList(),
    );

Map<String, dynamic> _$LlmProvider_v0ToJson(LlmProvider_v0 instance) =>
    <String, dynamic>{
      'name': instance.name,
      'baseUrl': instance.baseUrl,
      'apiKey': instance.apiKey,
      'alias': instance.alias,
      'description': instance.description,
      'models': instance.models,
    };
