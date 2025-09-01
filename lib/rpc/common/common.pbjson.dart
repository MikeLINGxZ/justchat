// This is a generated file - do not edit.
//
// Generated from rpc/common/common.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use imageURLDetailDescriptor instead')
const ImageURLDetail$json = {
  '1': 'ImageURLDetail',
  '2': [
    {'1': 'IMAGE_URL_DETAIL_UNSPECIFIED', '2': 0},
    {'1': 'IMAGE_URL_DETAIL_HIGH', '2': 1},
    {'1': 'IMAGE_URL_DETAIL_LOW', '2': 2},
    {'1': 'IMAGE_URL_DETAIL_AUTO', '2': 3},
  ],
};

/// Descriptor for `ImageURLDetail`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List imageURLDetailDescriptor = $convert.base64Decode(
    'Cg5JbWFnZVVSTERldGFpbBIgChxJTUFHRV9VUkxfREVUQUlMX1VOU1BFQ0lGSUVEEAASGQoVSU'
    '1BR0VfVVJMX0RFVEFJTF9ISUdIEAESGAoUSU1BR0VfVVJMX0RFVEFJTF9MT1cQAhIZChVJTUFH'
    'RV9VUkxfREVUQUlMX0FVVE8QAw==');

@$core.Deprecated('Use chatMessagePartTypeDescriptor instead')
const ChatMessagePartType$json = {
  '1': 'ChatMessagePartType',
  '2': [
    {'1': 'CHAT_MESSAGE_PART_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'CHAT_MESSAGE_PART_TYPE_TEXT', '2': 1},
    {'1': 'CHAT_MESSAGE_PART_TYPE_IMAGE_URL', '2': 2},
    {'1': 'CHAT_MESSAGE_PART_TYPE_AUDIO_URL', '2': 3},
    {'1': 'CHAT_MESSAGE_PART_TYPE_VIDEO_URL', '2': 4},
    {'1': 'CHAT_MESSAGE_PART_TYPE_FILE_URL', '2': 5},
  ],
};

/// Descriptor for `ChatMessagePartType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List chatMessagePartTypeDescriptor = $convert.base64Decode(
    'ChNDaGF0TWVzc2FnZVBhcnRUeXBlEiYKIkNIQVRfTUVTU0FHRV9QQVJUX1RZUEVfVU5TUEVDSU'
    'ZJRUQQABIfChtDSEFUX01FU1NBR0VfUEFSVF9UWVBFX1RFWFQQARIkCiBDSEFUX01FU1NBR0Vf'
    'UEFSVF9UWVBFX0lNQUdFX1VSTBACEiQKIENIQVRfTUVTU0FHRV9QQVJUX1RZUEVfQVVESU9fVV'
    'JMEAMSJAogQ0hBVF9NRVNTQUdFX1BBUlRfVFlQRV9WSURFT19VUkwQBBIjCh9DSEFUX01FU1NB'
    'R0VfUEFSVF9UWVBFX0ZJTEVfVVJMEAU=');

@$core.Deprecated('Use verificationCodeTypeDescriptor instead')
const VerificationCodeType$json = {
  '1': 'VerificationCodeType',
  '2': [
    {'1': 'VERIFICATION_CODE_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'VERIFICATION_CODE_TYPE_REGISTER', '2': 1},
    {'1': 'VERIFICATION_CODE_TYPE_RESET_PASSWORD', '2': 2},
  ],
};

/// Descriptor for `VerificationCodeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List verificationCodeTypeDescriptor = $convert.base64Decode(
    'ChRWZXJpZmljYXRpb25Db2RlVHlwZRImCiJWRVJJRklDQVRJT05fQ09ERV9UWVBFX1VOU1BFQ0'
    'lGSUVEEAASIwofVkVSSUZJQ0FUSU9OX0NPREVfVFlQRV9SRUdJU1RFUhABEikKJVZFUklGSUNB'
    'VElPTl9DT0RFX1RZUEVfUkVTRVRfUEFTU1dPUkQQAg==');

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

@$core.Deprecated('Use functionCallDescriptor instead')
const FunctionCall$json = {
  '1': 'FunctionCall',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'arguments', '3': 2, '4': 1, '5': 9, '10': 'arguments'},
  ],
};

/// Descriptor for `FunctionCall`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List functionCallDescriptor = $convert.base64Decode(
    'CgxGdW5jdGlvbkNhbGwSEgoEbmFtZRgBIAEoCVIEbmFtZRIcCglhcmd1bWVudHMYAiABKAlSCW'
    'FyZ3VtZW50cw==');

@$core.Deprecated('Use toolCallDescriptor instead')
const ToolCall$json = {
  '1': 'ToolCall',
  '2': [
    {'1': 'index', '3': 1, '4': 1, '5': 5, '9': 0, '10': 'index', '17': true},
    {'1': 'id', '3': 2, '4': 1, '5': 9, '10': 'id'},
    {'1': 'type', '3': 3, '4': 1, '5': 9, '10': 'type'},
    {'1': 'function', '3': 4, '4': 1, '5': 11, '6': '.lemon_tea.common.FunctionCall', '10': 'function'},
    {'1': 'extra', '3': 5, '4': 3, '5': 11, '6': '.lemon_tea.common.ToolCall.ExtraEntry', '10': 'extra'},
  ],
  '3': [ToolCall_ExtraEntry$json],
  '8': [
    {'1': '_index'},
  ],
};

@$core.Deprecated('Use toolCallDescriptor instead')
const ToolCall_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ToolCall`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List toolCallDescriptor = $convert.base64Decode(
    'CghUb29sQ2FsbBIZCgVpbmRleBgBIAEoBUgAUgVpbmRleIgBARIOCgJpZBgCIAEoCVICaWQSEg'
    'oEdHlwZRgDIAEoCVIEdHlwZRI6CghmdW5jdGlvbhgEIAEoCzIeLmxlbW9uX3RlYS5jb21tb24u'
    'RnVuY3Rpb25DYWxsUghmdW5jdGlvbhI7CgVleHRyYRgFIAMoCzIlLmxlbW9uX3RlYS5jb21tb2'
    '4uVG9vbENhbGwuRXh0cmFFbnRyeVIFZXh0cmEaOAoKRXh0cmFFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBQggKBl9pbmRleA==');

@$core.Deprecated('Use chatMessageImageURLDescriptor instead')
const ChatMessageImageURL$json = {
  '1': 'ChatMessageImageURL',
  '2': [
    {'1': 'url', '3': 1, '4': 1, '5': 9, '10': 'url'},
    {'1': 'uri', '3': 2, '4': 1, '5': 9, '10': 'uri'},
    {'1': 'detail', '3': 3, '4': 1, '5': 14, '6': '.lemon_tea.common.ImageURLDetail', '10': 'detail'},
    {'1': 'mime_type', '3': 4, '4': 1, '5': 9, '10': 'mimeType'},
    {'1': 'extra', '3': 5, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatMessageImageURL.ExtraEntry', '10': 'extra'},
  ],
  '3': [ChatMessageImageURL_ExtraEntry$json],
};

@$core.Deprecated('Use chatMessageImageURLDescriptor instead')
const ChatMessageImageURL_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ChatMessageImageURL`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageImageURLDescriptor = $convert.base64Decode(
    'ChNDaGF0TWVzc2FnZUltYWdlVVJMEhAKA3VybBgBIAEoCVIDdXJsEhAKA3VyaRgCIAEoCVIDdX'
    'JpEjgKBmRldGFpbBgDIAEoDjIgLmxlbW9uX3RlYS5jb21tb24uSW1hZ2VVUkxEZXRhaWxSBmRl'
    'dGFpbBIbCgltaW1lX3R5cGUYBCABKAlSCG1pbWVUeXBlEkYKBWV4dHJhGAUgAygLMjAubGVtb2'
    '5fdGVhLmNvbW1vbi5DaGF0TWVzc2FnZUltYWdlVVJMLkV4dHJhRW50cnlSBWV4dHJhGjgKCkV4'
    'dHJhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use chatMessageAudioURLDescriptor instead')
const ChatMessageAudioURL$json = {
  '1': 'ChatMessageAudioURL',
  '2': [
    {'1': 'url', '3': 1, '4': 1, '5': 9, '10': 'url'},
    {'1': 'uri', '3': 2, '4': 1, '5': 9, '10': 'uri'},
    {'1': 'mime_type', '3': 3, '4': 1, '5': 9, '10': 'mimeType'},
    {'1': 'extra', '3': 4, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatMessageAudioURL.ExtraEntry', '10': 'extra'},
  ],
  '3': [ChatMessageAudioURL_ExtraEntry$json],
};

@$core.Deprecated('Use chatMessageAudioURLDescriptor instead')
const ChatMessageAudioURL_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ChatMessageAudioURL`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageAudioURLDescriptor = $convert.base64Decode(
    'ChNDaGF0TWVzc2FnZUF1ZGlvVVJMEhAKA3VybBgBIAEoCVIDdXJsEhAKA3VyaRgCIAEoCVIDdX'
    'JpEhsKCW1pbWVfdHlwZRgDIAEoCVIIbWltZVR5cGUSRgoFZXh0cmEYBCADKAsyMC5sZW1vbl90'
    'ZWEuY29tbW9uLkNoYXRNZXNzYWdlQXVkaW9VUkwuRXh0cmFFbnRyeVIFZXh0cmEaOAoKRXh0cm'
    'FFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use chatMessageVideoURLDescriptor instead')
const ChatMessageVideoURL$json = {
  '1': 'ChatMessageVideoURL',
  '2': [
    {'1': 'url', '3': 1, '4': 1, '5': 9, '10': 'url'},
    {'1': 'uri', '3': 2, '4': 1, '5': 9, '10': 'uri'},
    {'1': 'mime_type', '3': 3, '4': 1, '5': 9, '10': 'mimeType'},
    {'1': 'extra', '3': 4, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatMessageVideoURL.ExtraEntry', '10': 'extra'},
  ],
  '3': [ChatMessageVideoURL_ExtraEntry$json],
};

@$core.Deprecated('Use chatMessageVideoURLDescriptor instead')
const ChatMessageVideoURL_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ChatMessageVideoURL`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageVideoURLDescriptor = $convert.base64Decode(
    'ChNDaGF0TWVzc2FnZVZpZGVvVVJMEhAKA3VybBgBIAEoCVIDdXJsEhAKA3VyaRgCIAEoCVIDdX'
    'JpEhsKCW1pbWVfdHlwZRgDIAEoCVIIbWltZVR5cGUSRgoFZXh0cmEYBCADKAsyMC5sZW1vbl90'
    'ZWEuY29tbW9uLkNoYXRNZXNzYWdlVmlkZW9VUkwuRXh0cmFFbnRyeVIFZXh0cmEaOAoKRXh0cm'
    'FFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use chatMessageFileURLDescriptor instead')
const ChatMessageFileURL$json = {
  '1': 'ChatMessageFileURL',
  '2': [
    {'1': 'url', '3': 1, '4': 1, '5': 9, '10': 'url'},
    {'1': 'uri', '3': 2, '4': 1, '5': 9, '10': 'uri'},
    {'1': 'mime_type', '3': 3, '4': 1, '5': 9, '10': 'mimeType'},
    {'1': 'name', '3': 4, '4': 1, '5': 9, '10': 'name'},
    {'1': 'extra', '3': 5, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatMessageFileURL.ExtraEntry', '10': 'extra'},
  ],
  '3': [ChatMessageFileURL_ExtraEntry$json],
};

@$core.Deprecated('Use chatMessageFileURLDescriptor instead')
const ChatMessageFileURL_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ChatMessageFileURL`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageFileURLDescriptor = $convert.base64Decode(
    'ChJDaGF0TWVzc2FnZUZpbGVVUkwSEAoDdXJsGAEgASgJUgN1cmwSEAoDdXJpGAIgASgJUgN1cm'
    'kSGwoJbWltZV90eXBlGAMgASgJUghtaW1lVHlwZRISCgRuYW1lGAQgASgJUgRuYW1lEkUKBWV4'
    'dHJhGAUgAygLMi8ubGVtb25fdGVhLmNvbW1vbi5DaGF0TWVzc2FnZUZpbGVVUkwuRXh0cmFFbn'
    'RyeVIFZXh0cmEaOAoKRXh0cmFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEo'
    'CVIFdmFsdWU6AjgB');

@$core.Deprecated('Use chatMessagePartDescriptor instead')
const ChatMessagePart$json = {
  '1': 'ChatMessagePart',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 14, '6': '.lemon_tea.common.ChatMessagePartType', '10': 'type'},
    {'1': 'text', '3': 2, '4': 1, '5': 9, '10': 'text'},
    {'1': 'image_url', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.common.ChatMessageImageURL', '9': 0, '10': 'imageUrl', '17': true},
    {'1': 'audio_url', '3': 4, '4': 1, '5': 11, '6': '.lemon_tea.common.ChatMessageAudioURL', '9': 1, '10': 'audioUrl', '17': true},
    {'1': 'video_url', '3': 5, '4': 1, '5': 11, '6': '.lemon_tea.common.ChatMessageVideoURL', '9': 2, '10': 'videoUrl', '17': true},
    {'1': 'file_url', '3': 6, '4': 1, '5': 11, '6': '.lemon_tea.common.ChatMessageFileURL', '9': 3, '10': 'fileUrl', '17': true},
  ],
  '8': [
    {'1': '_image_url'},
    {'1': '_audio_url'},
    {'1': '_video_url'},
    {'1': '_file_url'},
  ],
};

/// Descriptor for `ChatMessagePart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessagePartDescriptor = $convert.base64Decode(
    'Cg9DaGF0TWVzc2FnZVBhcnQSOQoEdHlwZRgBIAEoDjIlLmxlbW9uX3RlYS5jb21tb24uQ2hhdE'
    '1lc3NhZ2VQYXJ0VHlwZVIEdHlwZRISCgR0ZXh0GAIgASgJUgR0ZXh0EkcKCWltYWdlX3VybBgD'
    'IAEoCzIlLmxlbW9uX3RlYS5jb21tb24uQ2hhdE1lc3NhZ2VJbWFnZVVSTEgAUghpbWFnZVVybI'
    'gBARJHCglhdWRpb191cmwYBCABKAsyJS5sZW1vbl90ZWEuY29tbW9uLkNoYXRNZXNzYWdlQXVk'
    'aW9VUkxIAVIIYXVkaW9VcmyIAQESRwoJdmlkZW9fdXJsGAUgASgLMiUubGVtb25fdGVhLmNvbW'
    '1vbi5DaGF0TWVzc2FnZVZpZGVvVVJMSAJSCHZpZGVvVXJsiAEBEkQKCGZpbGVfdXJsGAYgASgL'
    'MiQubGVtb25fdGVhLmNvbW1vbi5DaGF0TWVzc2FnZUZpbGVVUkxIA1IHZmlsZVVybIgBAUIMCg'
    'pfaW1hZ2VfdXJsQgwKCl9hdWRpb191cmxCDAoKX3ZpZGVvX3VybEILCglfZmlsZV91cmw=');

@$core.Deprecated('Use topLogProbDescriptor instead')
const TopLogProb$json = {
  '1': 'TopLogProb',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {'1': 'log_prob', '3': 2, '4': 1, '5': 1, '10': 'logProb'},
    {'1': 'bytes', '3': 3, '4': 3, '5': 3, '10': 'bytes'},
  ],
};

/// Descriptor for `TopLogProb`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topLogProbDescriptor = $convert.base64Decode(
    'CgpUb3BMb2dQcm9iEhQKBXRva2VuGAEgASgJUgV0b2tlbhIZCghsb2dfcHJvYhgCIAEoAVIHbG'
    '9nUHJvYhIUCgVieXRlcxgDIAMoA1IFYnl0ZXM=');

@$core.Deprecated('Use logProbDescriptor instead')
const LogProb$json = {
  '1': 'LogProb',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
    {'1': 'log_prob', '3': 2, '4': 1, '5': 1, '10': 'logProb'},
    {'1': 'bytes', '3': 3, '4': 3, '5': 3, '10': 'bytes'},
    {'1': 'top_log_probs', '3': 4, '4': 3, '5': 11, '6': '.lemon_tea.common.TopLogProb', '10': 'topLogProbs'},
  ],
};

/// Descriptor for `LogProb`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logProbDescriptor = $convert.base64Decode(
    'CgdMb2dQcm9iEhQKBXRva2VuGAEgASgJUgV0b2tlbhIZCghsb2dfcHJvYhgCIAEoAVIHbG9nUH'
    'JvYhIUCgVieXRlcxgDIAMoA1IFYnl0ZXMSQAoNdG9wX2xvZ19wcm9icxgEIAMoCzIcLmxlbW9u'
    'X3RlYS5jb21tb24uVG9wTG9nUHJvYlILdG9wTG9nUHJvYnM=');

@$core.Deprecated('Use logProbsDescriptor instead')
const LogProbs$json = {
  '1': 'LogProbs',
  '2': [
    {'1': 'content', '3': 1, '4': 3, '5': 11, '6': '.lemon_tea.common.LogProb', '10': 'content'},
  ],
};

/// Descriptor for `LogProbs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logProbsDescriptor = $convert.base64Decode(
    'CghMb2dQcm9icxIzCgdjb250ZW50GAEgAygLMhkubGVtb25fdGVhLmNvbW1vbi5Mb2dQcm9iUg'
    'djb250ZW50');

@$core.Deprecated('Use tokenUsageDescriptor instead')
const TokenUsage$json = {
  '1': 'TokenUsage',
  '2': [
    {'1': 'prompt_tokens', '3': 1, '4': 1, '5': 5, '10': 'promptTokens'},
    {'1': 'completion_tokens', '3': 2, '4': 1, '5': 5, '10': 'completionTokens'},
    {'1': 'total_tokens', '3': 3, '4': 1, '5': 5, '10': 'totalTokens'},
  ],
};

/// Descriptor for `TokenUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tokenUsageDescriptor = $convert.base64Decode(
    'CgpUb2tlblVzYWdlEiMKDXByb21wdF90b2tlbnMYASABKAVSDHByb21wdFRva2VucxIrChFjb2'
    '1wbGV0aW9uX3Rva2VucxgCIAEoBVIQY29tcGxldGlvblRva2VucxIhCgx0b3RhbF90b2tlbnMY'
    'AyABKAVSC3RvdGFsVG9rZW5z');

@$core.Deprecated('Use responseMetaDescriptor instead')
const ResponseMeta$json = {
  '1': 'ResponseMeta',
  '2': [
    {'1': 'finish_reason', '3': 1, '4': 1, '5': 9, '10': 'finishReason'},
    {'1': 'usage', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.common.TokenUsage', '9': 0, '10': 'usage', '17': true},
    {'1': 'log_probs', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.common.LogProbs', '9': 1, '10': 'logProbs', '17': true},
  ],
  '8': [
    {'1': '_usage'},
    {'1': '_log_probs'},
  ],
};

/// Descriptor for `ResponseMeta`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseMetaDescriptor = $convert.base64Decode(
    'CgxSZXNwb25zZU1ldGESIwoNZmluaXNoX3JlYXNvbhgBIAEoCVIMZmluaXNoUmVhc29uEjcKBX'
    'VzYWdlGAIgASgLMhwubGVtb25fdGVhLmNvbW1vbi5Ub2tlblVzYWdlSABSBXVzYWdliAEBEjwK'
    'CWxvZ19wcm9icxgDIAEoCzIaLmxlbW9uX3RlYS5jb21tb24uTG9nUHJvYnNIAVIIbG9nUHJvYn'
    'OIAQFCCAoGX3VzYWdlQgwKCl9sb2dfcHJvYnM=');

@$core.Deprecated('Use messageDescriptor instead')
const Message$json = {
  '1': 'Message',
  '2': [
    {'1': 'content', '3': 1, '4': 1, '5': 9, '10': 'content'},
    {'1': 'multi_content', '3': 2, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatMessagePart', '10': 'multiContent'},
    {'1': 'name', '3': 3, '4': 1, '5': 9, '10': 'name'},
    {'1': 'extra', '3': 4, '4': 3, '5': 11, '6': '.lemon_tea.common.Message.ExtraEntry', '10': 'extra'},
    {'1': 'role', '3': 5, '4': 1, '5': 9, '10': 'role'},
    {'1': 'system_content', '3': 10, '4': 1, '5': 9, '10': 'systemContent'},
    {'1': 'tool_calls', '3': 11, '4': 3, '5': 11, '6': '.lemon_tea.common.ToolCall', '10': 'toolCalls'},
    {'1': 'response_meta', '3': 12, '4': 1, '5': 11, '6': '.lemon_tea.common.ResponseMeta', '10': 'responseMeta'},
    {'1': 'reasoning_content', '3': 13, '4': 1, '5': 9, '10': 'reasoningContent'},
    {'1': 'tool_call_id', '3': 14, '4': 1, '5': 9, '10': 'toolCallId'},
    {'1': 'tool_name', '3': 15, '4': 1, '5': 9, '10': 'toolName'},
    {'1': 'chat_uuid', '3': 16, '4': 1, '5': 9, '10': 'chatUuid'},
  ],
  '3': [Message_ExtraEntry$json],
};

@$core.Deprecated('Use messageDescriptor instead')
const Message_ExtraEntry$json = {
  '1': 'ExtraEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Message`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDescriptor = $convert.base64Decode(
    'CgdNZXNzYWdlEhgKB2NvbnRlbnQYASABKAlSB2NvbnRlbnQSRgoNbXVsdGlfY29udGVudBgCIA'
    'MoCzIhLmxlbW9uX3RlYS5jb21tb24uQ2hhdE1lc3NhZ2VQYXJ0UgxtdWx0aUNvbnRlbnQSEgoE'
    'bmFtZRgDIAEoCVIEbmFtZRI6CgVleHRyYRgEIAMoCzIkLmxlbW9uX3RlYS5jb21tb24uTWVzc2'
    'FnZS5FeHRyYUVudHJ5UgVleHRyYRISCgRyb2xlGAUgASgJUgRyb2xlEiUKDnN5c3RlbV9jb250'
    'ZW50GAogASgJUg1zeXN0ZW1Db250ZW50EjkKCnRvb2xfY2FsbHMYCyADKAsyGi5sZW1vbl90ZW'
    'EuY29tbW9uLlRvb2xDYWxsUgl0b29sQ2FsbHMSQwoNcmVzcG9uc2VfbWV0YRgMIAEoCzIeLmxl'
    'bW9uX3RlYS5jb21tb24uUmVzcG9uc2VNZXRhUgxyZXNwb25zZU1ldGESKwoRcmVhc29uaW5nX2'
    'NvbnRlbnQYDSABKAlSEHJlYXNvbmluZ0NvbnRlbnQSIAoMdG9vbF9jYWxsX2lkGA4gASgJUgp0'
    'b29sQ2FsbElkEhsKCXRvb2xfbmFtZRgPIAEoCVIIdG9vbE5hbWUSGwoJY2hhdF91dWlkGBAgAS'
    'gJUghjaGF0VXVpZBo4CgpFeHRyYUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIg'
    'ASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use chatInfoDescriptor instead')
const ChatInfo$json = {
  '1': 'ChatInfo',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
    {'1': 'model_id', '3': 3, '4': 1, '5': 3, '10': 'modelId'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 3, '10': 'createdAt'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 3, '10': 'updatedAt'},
    {'1': 'message_count', '3': 6, '4': 1, '5': 5, '10': 'messageCount'},
    {'1': 'last_message_preview', '3': 7, '4': 1, '5': 9, '10': 'lastMessagePreview'},
    {'1': 'metadata', '3': 8, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatInfo.MetadataEntry', '10': 'metadata'},
  ],
  '3': [ChatInfo_MetadataEntry$json],
};

@$core.Deprecated('Use chatInfoDescriptor instead')
const ChatInfo_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ChatInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatInfoDescriptor = $convert.base64Decode(
    'CghDaGF0SW5mbxIbCgljaGF0X3V1aWQYASABKAlSCGNoYXRVdWlkEhQKBXRpdGxlGAIgASgJUg'
    'V0aXRsZRIZCghtb2RlbF9pZBgDIAEoA1IHbW9kZWxJZBIdCgpjcmVhdGVkX2F0GAQgASgDUglj'
    'cmVhdGVkQXQSHQoKdXBkYXRlZF9hdBgFIAEoA1IJdXBkYXRlZEF0EiMKDW1lc3NhZ2VfY291bn'
    'QYBiABKAVSDG1lc3NhZ2VDb3VudBIwChRsYXN0X21lc3NhZ2VfcHJldmlldxgHIAEoCVISbGFz'
    'dE1lc3NhZ2VQcmV2aWV3EkQKCG1ldGFkYXRhGAggAygLMigubGVtb25fdGVhLmNvbW1vbi5DaG'
    'F0SW5mby5NZXRhZGF0YUVudHJ5UghtZXRhZGF0YRo7Cg1NZXRhZGF0YUVudHJ5EhAKA2tleRgB'
    'IAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

