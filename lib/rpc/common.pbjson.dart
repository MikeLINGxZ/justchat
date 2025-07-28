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

@$core.Deprecated('Use fileTypeDescriptor instead')
const FileType$json = {
  '1': 'FileType',
  '2': [
    {'1': 'FILE_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'FILE_TYPE_IMAGE', '2': 1},
    {'1': 'FILE_TYPE_DOCUMENT', '2': 2},
    {'1': 'FILE_TYPE_AUDIO', '2': 3},
    {'1': 'FILE_TYPE_VIDEO', '2': 4},
    {'1': 'FILE_TYPE_OTHER', '2': 5},
  ],
};

/// Descriptor for `FileType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List fileTypeDescriptor = $convert.base64Decode(
    'CghGaWxlVHlwZRIZChVGSUxFX1RZUEVfVU5TUEVDSUZJRUQQABITCg9GSUxFX1RZUEVfSU1BR0'
    'UQARIWChJGSUxFX1RZUEVfRE9DVU1FTlQQAhITCg9GSUxFX1RZUEVfQVVESU8QAxITCg9GSUxF'
    'X1RZUEVfVklERU8QBBITCg9GSUxFX1RZUEVfT1RIRVIQBQ==');

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

@$core.Deprecated('Use fileContentDescriptor instead')
const FileContent$json = {
  '1': 'FileContent',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'mime_type', '3': 2, '4': 1, '5': 9, '10': 'mimeType'},
    {'1': 'type', '3': 3, '4': 1, '5': 14, '6': '.lemon_tea.common.FileType', '10': 'type'},
    {'1': 'data', '3': 4, '4': 1, '5': 12, '10': 'data'},
    {'1': 'size', '3': 5, '4': 1, '5': 3, '10': 'size'},
    {'1': 'url', '3': 6, '4': 1, '5': 9, '9': 0, '10': 'url', '17': true},
    {'1': 'description', '3': 7, '4': 1, '5': 9, '9': 1, '10': 'description', '17': true},
  ],
  '8': [
    {'1': '_url'},
    {'1': '_description'},
  ],
};

/// Descriptor for `FileContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List fileContentDescriptor = $convert.base64Decode(
    'CgtGaWxlQ29udGVudBISCgRuYW1lGAEgASgJUgRuYW1lEhsKCW1pbWVfdHlwZRgCIAEoCVIIbW'
    'ltZVR5cGUSLgoEdHlwZRgDIAEoDjIaLmxlbW9uX3RlYS5jb21tb24uRmlsZVR5cGVSBHR5cGUS'
    'EgoEZGF0YRgEIAEoDFIEZGF0YRISCgRzaXplGAUgASgDUgRzaXplEhUKA3VybBgGIAEoCUgAUg'
    'N1cmyIAQESJQoLZGVzY3JpcHRpb24YByABKAlIAVILZGVzY3JpcHRpb26IAQFCBgoEX3VybEIO'
    'CgxfZGVzY3JpcHRpb24=');

@$core.Deprecated('Use messageContentDescriptor instead')
const MessageContent$json = {
  '1': 'MessageContent',
  '2': [
    {'1': 'text', '3': 1, '4': 1, '5': 9, '9': 0, '10': 'text'},
    {'1': 'file', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.common.FileContent', '9': 0, '10': 'file'},
  ],
  '8': [
    {'1': 'content'},
  ],
};

/// Descriptor for `MessageContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageContentDescriptor = $convert.base64Decode(
    'Cg5NZXNzYWdlQ29udGVudBIUCgR0ZXh0GAEgASgJSABSBHRleHQSMwoEZmlsZRgCIAEoCzIdLm'
    'xlbW9uX3RlYS5jb21tb24uRmlsZUNvbnRlbnRIAFIEZmlsZUIJCgdjb250ZW50');

@$core.Deprecated('Use messageDescriptor instead')
const Message$json = {
  '1': 'Message',
  '2': [
    {'1': 'role', '3': 1, '4': 1, '5': 14, '6': '.lemon_tea.common.MessageRole', '10': 'role'},
    {'1': 'contents', '3': 2, '4': 3, '5': 11, '6': '.lemon_tea.common.MessageContent', '10': 'contents'},
    {
      '1': 'content',
      '3': 3,
      '4': 1,
      '5': 9,
      '8': {'3': true},
      '10': 'content',
    },
  ],
};

/// Descriptor for `Message`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDescriptor = $convert.base64Decode(
    'CgdNZXNzYWdlEjEKBHJvbGUYASABKA4yHS5sZW1vbl90ZWEuY29tbW9uLk1lc3NhZ2VSb2xlUg'
    'Ryb2xlEjwKCGNvbnRlbnRzGAIgAygLMiAubGVtb25fdGVhLmNvbW1vbi5NZXNzYWdlQ29udGVu'
    'dFIIY29udGVudHMSHAoHY29udGVudBgDIAEoCUICGAFSB2NvbnRlbnQ=');

