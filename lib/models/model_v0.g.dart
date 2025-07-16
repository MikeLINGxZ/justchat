// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model_v0.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Model_v0 _$Model_v0FromJson(Map<String, dynamic> json) => Model_v0(
  id: json['id'] as String,
  object: json['object'] as String,
  ownedBy: json['owned_by'] as String,
  enabled: json['enabled'] as bool? ?? true,
);

Map<String, dynamic> _$Model_v0ToJson(Model_v0 instance) => <String, dynamic>{
  'id': instance.id,
  'object': instance.object,
  'owned_by': instance.ownedBy,
  'enabled': instance.enabled,
};
