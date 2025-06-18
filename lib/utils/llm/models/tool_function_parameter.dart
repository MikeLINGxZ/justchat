import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/utils/llm/models/tool_function_propertie.dart';

part 'tool_function_parameter.g.dart';

@JsonSerializable()
class ToolFunctionParameter {
  final String type;
  final Map<String, ToolFunctionPropertie> properties;
  final List<String>? required;

  const ToolFunctionParameter({
    required this.type,
    required this.properties,
    this.required,
  });

  factory ToolFunctionParameter.fromJson(Map<String, dynamic> json) =>
      _$ToolFunctionParameterFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolFunctionParameterToJson(this);
}