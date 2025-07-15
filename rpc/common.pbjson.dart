// This is a generated file - do not edit.
//
// Generated from rpc/common.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use messageRoleDescriptor instead')
const MessageRole$json = {
  '1': 'MessageRole',
  '2': [
    {'1': 'MESSAGE_ROLE_UNSPECIFIED', '2': 0},
    {'1': 'MESSAGE_ROLE_SYSTEM', '2': 1},
    {'1': 'MESSAGE_ROLE_USER', '2': 2},
    {'1': 'MESSAGE_ROLE_ASSISTANT', '2': 3},
    {'1': 'MESSAGE_ROLE_FUNCTION', '2': 4},
    {'1': 'MESSAGE_ROLE_TOOL', '2': 5},
  ],
};

/// Descriptor for `MessageRole`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List messageRoleDescriptor = $convert.base64Decode(
    'CgtNZXNzYWdlUm9sZRIcChhNRVNTQUdFX1JPTEVfVU5TUEVDSUZJRUQQABIXChNNRVNTQUdFX1'
    'JPTEVfU1lTVEVNEAESFQoRTUVTU0FHRV9ST0xFX1VTRVIQAhIaChZNRVNTQUdFX1JPTEVfQVNT'
    'SVNUQU5UEAMSGQoVTUVTU0FHRV9ST0xFX0ZVTkNUSU9OEAQSFQoRTUVTU0FHRV9ST0xFX1RPT0'
    'wQBQ==');

@$core.Deprecated('Use emptyDescriptor instead')
const Empty$json = {
  '1': 'Empty',
};

/// Descriptor for `Empty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emptyDescriptor = $convert.base64Decode(
    'CgVFbXB0eQ==');

@$core.Deprecated('Use llmProviderDescriptor instead')
const LlmProvider$json = {
  '1': 'LlmProvider',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'base_url', '3': 3, '4': 1, '5': 9, '10': 'baseUrl'},
    {'1': 'api_key', '3': 4, '4': 1, '5': 9, '10': 'apiKey'},
    {'1': 'alias', '3': 5, '4': 1, '5': 9, '9': 0, '10': 'alias', '17': true},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '9': 1, '10': 'description', '17': true},
    {'1': 'models', '3': 7, '4': 3, '5': 11, '6': '.lemon_tea.common.Model', '10': 'models'},
  ],
  '8': [
    {'1': '_alias'},
    {'1': '_description'},
  ],
};

/// Descriptor for `LlmProvider`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List llmProviderDescriptor = $convert.base64Decode(
    'CgtMbG1Qcm92aWRlchIOCgJpZBgBIAEoCVICaWQSEgoEbmFtZRgCIAEoCVIEbmFtZRIZCghiYX'
    'NlX3VybBgDIAEoCVIHYmFzZVVybBIXCgdhcGlfa2V5GAQgASgJUgZhcGlLZXkSGQoFYWxpYXMY'
    'BSABKAlIAFIFYWxpYXOIAQESJQoLZGVzY3JpcHRpb24YBiABKAlIAVILZGVzY3JpcHRpb26IAQ'
    'ESLwoGbW9kZWxzGAcgAygLMhcubGVtb25fdGVhLmNvbW1vbi5Nb2RlbFIGbW9kZWxzQggKBl9h'
    'bGlhc0IOCgxfZGVzY3JpcHRpb24=');

@$core.Deprecated('Use modelDescriptor instead')
const Model$json = {
  '1': 'Model',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'object', '3': 2, '4': 1, '5': 9, '10': 'object'},
    {'1': 'owned_by', '3': 3, '4': 1, '5': 9, '10': 'ownedBy'},
    {'1': 'enabled', '3': 4, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `Model`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelDescriptor = $convert.base64Decode(
    'CgVNb2RlbBIOCgJpZBgBIAEoCVICaWQSFgoGb2JqZWN0GAIgASgJUgZvYmplY3QSGQoIb3duZW'
    'RfYnkYAyABKAlSB293bmVkQnkSGAoHZW5hYmxlZBgEIAEoCFIHZW5hYmxlZA==');

@$core.Deprecated('Use messageDescriptor instead')
const Message$json = {
  '1': 'Message',
  '2': [
    {'1': 'role', '3': 1, '4': 1, '5': 14, '6': '.lemon_tea.common.MessageRole', '10': 'role'},
    {'1': 'content', '3': 2, '4': 1, '5': 9, '10': 'content'},
  ],
};

/// Descriptor for `Message`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDescriptor = $convert.base64Decode(
    'CgdNZXNzYWdlEjEKBHJvbGUYASABKA4yHS5sZW1vbl90ZWEuY29tbW9uLk1lc3NhZ2VSb2xlUg'
    'Ryb2xlEhgKB2NvbnRlbnQYAiABKAlSB2NvbnRlbnQ=');

