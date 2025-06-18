import 'package:json_annotation/json_annotation.dart';

part 'tool_function_propertie.g.dart';

@JsonSerializable()
class ToolFunctionPropertie {
  final String type;
  final String description;
  @JsonKey(name: 'enum')
  final List<String>? enums;

  const ToolFunctionPropertie({
    required this.type,
    required this.description,
    this.enums,
  });

  factory ToolFunctionPropertie.fromJson(Map<String, dynamic> json) =>
      _$ToolFunctionPropertieFromJson(json);
  
  Map<String, dynamic> toJson() => _$ToolFunctionPropertieToJson(this);
}