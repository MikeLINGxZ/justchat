import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/utils/llm/models/tool_function_parameter.dart';

part 'tool_function.g.dart';

@JsonSerializable()
class ToolFunction {
  final String name;
  final String description;
  final ToolFunctionParameter parameters;

  const ToolFunction({
    required this.name,
    required this.description,
    required this.parameters,
  });

  factory ToolFunction.fromJson(Map<String, dynamic> json) =>
      _$ToolFunctionFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolFunctionToJson(this);
}
