import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/utils/llm/models/tool_function.dart';

part 'tool.g.dart';

@JsonSerializable()
class Tool {
  final String type;
  final ToolFunction function;

  const Tool({
    required this.type,
    required this.function,
  });

  factory Tool.fromJson(Map<String, dynamic> json) =>
      _$ToolFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolToJson(this);
}