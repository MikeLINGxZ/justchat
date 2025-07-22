// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Model _$ModelFromJson(Map<String, dynamic> json) => Model(
  llmProviderId: json['llmProviderId'] as String,
  id: json['id'] as String,
  object: json['object'] as String? ?? 'model',
  ownedBy: json['owned_by'] as String,
  enabled: json['enabled'] as bool? ?? true,
  isCustom: json['isCustom'] as bool? ?? false,
  seqId: (json['seqId'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$ModelToJson(Model instance) => <String, dynamic>{
  'llmProviderId': instance.llmProviderId,
  'id': instance.id,
  'object': instance.object,
  'owned_by': instance.ownedBy,
  'enabled': instance.enabled,
  'isCustom': instance.isCustom,
  'seqId': instance.seqId,
};
