// This is a generated file - do not edit.
//
// Generated from rpc/common.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// MessageRole 对话角色
class MessageRole extends $pb.ProtobufEnum {
  static const MessageRole MESSAGE_ROLE_UNSPECIFIED = MessageRole._(0, _omitEnumNames ? '' : 'MESSAGE_ROLE_UNSPECIFIED');
  static const MessageRole MESSAGE_ROLE_SYSTEM = MessageRole._(1, _omitEnumNames ? '' : 'MESSAGE_ROLE_SYSTEM');
  static const MessageRole MESSAGE_ROLE_USER = MessageRole._(2, _omitEnumNames ? '' : 'MESSAGE_ROLE_USER');
  static const MessageRole MESSAGE_ROLE_ASSISTANT = MessageRole._(3, _omitEnumNames ? '' : 'MESSAGE_ROLE_ASSISTANT');
  static const MessageRole MESSAGE_ROLE_FUNCTION = MessageRole._(4, _omitEnumNames ? '' : 'MESSAGE_ROLE_FUNCTION');
  static const MessageRole MESSAGE_ROLE_TOOL = MessageRole._(5, _omitEnumNames ? '' : 'MESSAGE_ROLE_TOOL');

  static const $core.List<MessageRole> values = <MessageRole> [
    MESSAGE_ROLE_UNSPECIFIED,
    MESSAGE_ROLE_SYSTEM,
    MESSAGE_ROLE_USER,
    MESSAGE_ROLE_ASSISTANT,
    MESSAGE_ROLE_FUNCTION,
    MESSAGE_ROLE_TOOL,
  ];

  static final $core.List<MessageRole?> _byValue = $pb.ProtobufEnum.$_initByValueList(values, 5);
  static MessageRole? valueOf($core.int value) =>  value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MessageRole._(super.value, super.name);
}


const $core.bool _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
