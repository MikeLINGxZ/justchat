// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'conversation_v0.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Conversation_v0 _$Conversation_v0FromJson(Map<String, dynamic> json) =>
    Conversation_v0(
      id: json['id'] as String,
      title: json['title'] as String,
      messages:
          (json['messages'] as List<dynamic>)
              .map((e) => Message.fromJson(e as Map<String, dynamic>))
              .toList(),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
      isDeleted: json['isDeleted'] as bool? ?? false,
    );

Map<String, dynamic> _$Conversation_v0ToJson(Conversation_v0 instance) =>
    <String, dynamic>{
      'id': instance.id,
      'title': instance.title,
      'messages': instance.messages,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
      'isDeleted': instance.isDeleted,
    };
