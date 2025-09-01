// This is a generated file - do not edit.
//
// Generated from rpc/service/auth.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// FieldType 字段类型枚举
class FieldType extends $pb.ProtobufEnum {
  static const FieldType FIELD_TYPE_UNSPECIFIED = FieldType._(0, _omitEnumNames ? '' : 'FIELD_TYPE_UNSPECIFIED');
  static const FieldType FIELD_TYPE_USERNAME = FieldType._(1, _omitEnumNames ? '' : 'FIELD_TYPE_USERNAME');
  static const FieldType FIELD_TYPE_EMAIL = FieldType._(2, _omitEnumNames ? '' : 'FIELD_TYPE_EMAIL');

  static const $core.List<FieldType> values = <FieldType> [
    FIELD_TYPE_UNSPECIFIED,
    FIELD_TYPE_USERNAME,
    FIELD_TYPE_EMAIL,
  ];

  static final $core.List<FieldType?> _byValue = $pb.ProtobufEnum.$_initByValueList(values, 2);
  static FieldType? valueOf($core.int value) =>  value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FieldType._(super.value, super.name);
}


const $core.bool _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
