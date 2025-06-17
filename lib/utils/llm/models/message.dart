import 'package:json_annotation/json_annotation.dart';

part 'message.g.dart';

enum MessageRole {
  @JsonValue('system')
  system,
  @JsonValue('user')
  user,
  @JsonValue('assistant')
  assistant,
  @JsonValue('function')
  function,
}

@JsonSerializable()
class Message {
  /// 消息角色
  final MessageRole role;

  /// 消息内容（对于function角色可能是空字符串）
  final String content;

  /// 函数名称（role为function时必需）
  final String? name;

  /// 函数调用信息（role为assistant时可能存在）
  final Map<String, dynamic>? functionCall;

  Message({
    required this.role,
    required this.content,
    this.name,
    this.functionCall,
  }) : assert(content.isNotEmpty),
       assert(role != MessageRole.function || name != null);

  factory Message.fromJson(Map<String, dynamic> json) => _$MessageFromJson(json);
  Map<String, dynamic> toJson() => _$MessageToJson(this);

  Message copyWith({
    MessageRole? role,
    String? content,
    String? name,
    Map<String, dynamic>? functionCall,
  }) =>
      Message(
        role: role ?? this.role,
        content: content ?? this.content,
        name: name ?? this.name,
        functionCall: functionCall ?? this.functionCall,
      );
}