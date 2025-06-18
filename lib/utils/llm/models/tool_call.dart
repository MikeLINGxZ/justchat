import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/utils/llm/models/tool_call_function.dart';

part 'tool_call.g.dart';

@JsonSerializable()
class ToolCall {
  final int index;
  final String id;
  final String type;
  final ToolCallFunction function;

  const ToolCall({
    required this.index,
    required this.id,
    required this.type,
    required this.function,
  });

  factory ToolCall.fromJson(Map<String, dynamic> json) =>
      _$ToolCallFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolCallToJson(this);
}

