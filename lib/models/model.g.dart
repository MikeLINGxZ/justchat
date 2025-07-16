// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Model _$ModelFromJson(Map<String, dynamic> json) => Model(
  llmProviderId: json['llm_provider_id'] as String,
  id: json['id'] as String,
  object: json['object'] as String? ?? 'model',
  ownedBy: json['owned_by'] as String,
  enabled: json['enabled'] as bool? ?? true,
);

Map<String, dynamic> _$ModelToJson(Model instance) => <String, dynamic>{
  'llm_provider_id': instance.llmProviderId,
  'id': instance.id,
  'object': instance.object,
  'owned_by': instance.ownedBy,
  'enabled': instance.enabled,
};
