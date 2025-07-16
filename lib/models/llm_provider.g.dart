// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'llm_provider.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

LlmProvider _$LlmProviderFromJson(Map<String, dynamic> json) => LlmProvider(
  id: json['id'] as String,
  name: json['name'] as String,
  baseUrl: json['baseUrl'] as String,
  apiKey: json['apiKey'] as String?,
  alias: json['alias'] as String?,
  description: json['description'] as String?,
  enable: json['enable'] as bool? ?? true,
  checked: json['checked'] as bool? ?? false,
  defaultProviderId: json['defaultProviderId'] as String,
  defaultModelId: json['defaultModelId'] as String,
);

Map<String, dynamic> _$LlmProviderToJson(LlmProvider instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'baseUrl': instance.baseUrl,
      'apiKey': instance.apiKey,
      'alias': instance.alias,
      'description': instance.description,
      'enable': instance.enable,
      'checked': instance.checked,
      'defaultProviderId': instance.defaultProviderId,
      'defaultModelId': instance.defaultModelId,
    };
