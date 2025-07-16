import 'package:json_annotation/json_annotation.dart';

enum MessageRole {
  @JsonValue('system')
  system,
  @JsonValue('user')
  user,
  @JsonValue('assistant')
  assistant,
  @JsonValue('tool')
  tool,
}