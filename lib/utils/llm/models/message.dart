import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/tool_call.dart';

part 'message.g.dart';

@JsonSerializable()
class Message {
  /// 消息角色
  final MessageRole role;

  /// 消息内容（对于function角色可能是空字符串）
  final String content;

  /// 思考过程内容（思维链模型返回的思考过程）
  final String? reasoningContent;

  /// 函数调用信息（role为assistant时可能存在）
  final List<ToolCall>? toolCalls;

  // 工具调用结果id
  String? toolCallId;

  // todo visible
  

  Message({
    required this.role,
    required this.content,
    this.reasoningContent,
    this.toolCalls,
    this.toolCallId,
  }) : assert(
    // 允许assistant和tool角色的content为空（用于流式输出和工具调用）
    role == MessageRole.assistant || 
    role == MessageRole.tool ||
    content.isNotEmpty
  );

  factory Message.fromJson(Map<String, dynamic> json) =>
      _$MessageFromJson(json);
  Map<String, dynamic> toJson() => _$MessageToJson(this);

  Message copyWith({
    MessageRole? role,
    String? content,
    String? reasoningContent,
    List<ToolCall>? toolCalls,
    String? toolCallId,
  }) => Message(
    role: role ?? this.role,
    content: content ?? this.content,
    reasoningContent: reasoningContent ?? this.reasoningContent,
    toolCalls: toolCalls ?? this.toolCalls,
    toolCallId: this.toolCallId,
  );
}
