// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Model _$ModelFromJson(Map<String, dynamic> json) => Model(
  id: json['id'] as String,
  object: json['object'] as String,
  ownedBy: json['OwnedBy'] as String,
);

Map<String, dynamic> _$ModelToJson(Model instance) => <String, dynamic>{
  'id': instance.id,
  'object': instance.object,
  'OwnedBy': instance.ownedBy,
};
