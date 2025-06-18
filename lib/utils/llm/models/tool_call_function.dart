import 'package:json_annotation/json_annotation.dart';

part 'tool_call_function.g.dart';

@JsonSerializable()
class ToolCallFunction {
  final String name;
  final String arguments;

  const ToolCallFunction({
    required this.name,
    required this.arguments,
  });

  factory ToolCallFunction.fromJson(Map<String, dynamic> json) =>
      _$ToolCallFunctionFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolCallFunctionToJson(this);
}