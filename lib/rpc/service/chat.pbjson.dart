// This is a generated file - do not edit.
//
// Generated from rpc/service/chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use completionsRequestDescriptor instead')
const CompletionsRequest$json = {
  '1': 'CompletionsRequest',
  '2': [
    {'1': 'model', '3': 1, '4': 1, '5': 9, '10': 'model'},
    {'1': 'messages', '3': 2, '4': 3, '5': 11, '6': '.lemon_tea.common.Message', '10': 'messages'},
    {'1': 'temperature', '3': 3, '4': 1, '5': 1, '9': 0, '10': 'temperature', '17': true},
    {'1': 'max_tokens', '3': 4, '4': 1, '5': 5, '9': 1, '10': 'maxTokens', '17': true},
    {'1': 'top_p', '3': 5, '4': 1, '5': 1, '9': 2, '10': 'topP', '17': true},
    {'1': 'n', '3': 6, '4': 1, '5': 5, '9': 3, '10': 'n', '17': true},
    {'1': 'stream', '3': 7, '4': 1, '5': 8, '9': 4, '10': 'stream', '17': true},
    {'1': 'stop', '3': 8, '4': 1, '5': 9, '9': 5, '10': 'stop', '17': true},
    {'1': 'stop_sequence', '3': 9, '4': 3, '5': 9, '10': 'stopSequence'},
    {'1': 'presence_penalty', '3': 10, '4': 1, '5': 1, '9': 6, '10': 'presencePenalty', '17': true},
    {'1': 'frequency_penalty', '3': 11, '4': 1, '5': 1, '9': 7, '10': 'frequencyPenalty', '17': true},
    {'1': 'repetition_penalty', '3': 12, '4': 1, '5': 1, '9': 8, '10': 'repetitionPenalty', '17': true},
    {'1': 'user', '3': 13, '4': 1, '5': 9, '9': 9, '10': 'user', '17': true},
    {'1': 'tools', '3': 14, '4': 3, '5': 11, '6': '.lemon_tea.server.ChatCompletionTool', '10': 'tools'},
    {'1': 'tool_choice', '3': 15, '4': 1, '5': 11, '6': '.lemon_tea.server.ToolChoice', '9': 10, '10': 'toolChoice', '17': true},
    {'1': 'response_format', '3': 16, '4': 1, '5': 11, '6': '.lemon_tea.server.ResponseFormat', '9': 11, '10': 'responseFormat', '17': true},
    {'1': 'seed', '3': 17, '4': 1, '5': 5, '9': 12, '10': 'seed', '17': true},
    {'1': 'metadata', '3': 18, '4': 3, '5': 11, '6': '.lemon_tea.server.CompletionsRequest.MetadataEntry', '10': 'metadata'},
    {'1': 'non_standard', '3': 19, '4': 1, '5': 8, '10': 'nonStandard'},
    {'1': 'chat_uuid', '3': 20, '4': 1, '5': 9, '10': 'chatUuid'},
    {'1': 'completions_custom', '3': 21, '4': 1, '5': 11, '6': '.lemon_tea.server.CompletionsCustom', '10': 'completionsCustom'},
  ],
  '3': [CompletionsRequest_MetadataEntry$json],
  '8': [
    {'1': '_temperature'},
    {'1': '_max_tokens'},
    {'1': '_top_p'},
    {'1': '_n'},
    {'1': '_stream'},
    {'1': '_stop'},
    {'1': '_presence_penalty'},
    {'1': '_frequency_penalty'},
    {'1': '_repetition_penalty'},
    {'1': '_user'},
    {'1': '_tool_choice'},
    {'1': '_response_format'},
    {'1': '_seed'},
  ],
};

@$core.Deprecated('Use completionsRequestDescriptor instead')
const CompletionsRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CompletionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completionsRequestDescriptor = $convert.base64Decode(
    'ChJDb21wbGV0aW9uc1JlcXVlc3QSFAoFbW9kZWwYASABKAlSBW1vZGVsEjUKCG1lc3NhZ2VzGA'
    'IgAygLMhkubGVtb25fdGVhLmNvbW1vbi5NZXNzYWdlUghtZXNzYWdlcxIlCgt0ZW1wZXJhdHVy'
    'ZRgDIAEoAUgAUgt0ZW1wZXJhdHVyZYgBARIiCgptYXhfdG9rZW5zGAQgASgFSAFSCW1heFRva2'
    'Vuc4gBARIYCgV0b3BfcBgFIAEoAUgCUgR0b3BQiAEBEhEKAW4YBiABKAVIA1IBbogBARIbCgZz'
    'dHJlYW0YByABKAhIBFIGc3RyZWFtiAEBEhcKBHN0b3AYCCABKAlIBVIEc3RvcIgBARIjCg1zdG'
    '9wX3NlcXVlbmNlGAkgAygJUgxzdG9wU2VxdWVuY2USLgoQcHJlc2VuY2VfcGVuYWx0eRgKIAEo'
    'AUgGUg9wcmVzZW5jZVBlbmFsdHmIAQESMAoRZnJlcXVlbmN5X3BlbmFsdHkYCyABKAFIB1IQZn'
    'JlcXVlbmN5UGVuYWx0eYgBARIyChJyZXBldGl0aW9uX3BlbmFsdHkYDCABKAFICFIRcmVwZXRp'
    'dGlvblBlbmFsdHmIAQESFwoEdXNlchgNIAEoCUgJUgR1c2VyiAEBEjoKBXRvb2xzGA4gAygLMi'
    'QubGVtb25fdGVhLnNlcnZlci5DaGF0Q29tcGxldGlvblRvb2xSBXRvb2xzEkIKC3Rvb2xfY2hv'
    'aWNlGA8gASgLMhwubGVtb25fdGVhLnNlcnZlci5Ub29sQ2hvaWNlSApSCnRvb2xDaG9pY2WIAQ'
    'ESTgoPcmVzcG9uc2VfZm9ybWF0GBAgASgLMiAubGVtb25fdGVhLnNlcnZlci5SZXNwb25zZUZv'
    'cm1hdEgLUg5yZXNwb25zZUZvcm1hdIgBARIXCgRzZWVkGBEgASgFSAxSBHNlZWSIAQESTgoIbW'
    'V0YWRhdGEYEiADKAsyMi5sZW1vbl90ZWEuc2VydmVyLkNvbXBsZXRpb25zUmVxdWVzdC5NZXRh'
    'ZGF0YUVudHJ5UghtZXRhZGF0YRIhCgxub25fc3RhbmRhcmQYEyABKAhSC25vblN0YW5kYXJkEh'
    'sKCWNoYXRfdXVpZBgUIAEoCVIIY2hhdFV1aWQSUgoSY29tcGxldGlvbnNfY3VzdG9tGBUgASgL'
    'MiMubGVtb25fdGVhLnNlcnZlci5Db21wbGV0aW9uc0N1c3RvbVIRY29tcGxldGlvbnNDdXN0b2'
    '0aOwoNTWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFs'
    'dWU6AjgBQg4KDF90ZW1wZXJhdHVyZUINCgtfbWF4X3Rva2Vuc0IICgZfdG9wX3BCBAoCX25CCQ'
    'oHX3N0cmVhbUIHCgVfc3RvcEITChFfcHJlc2VuY2VfcGVuYWx0eUIUChJfZnJlcXVlbmN5X3Bl'
    'bmFsdHlCFQoTX3JlcGV0aXRpb25fcGVuYWx0eUIHCgVfdXNlckIOCgxfdG9vbF9jaG9pY2VCEg'
    'oQX3Jlc3BvbnNlX2Zvcm1hdEIHCgVfc2VlZA==');

@$core.Deprecated('Use completionsCustomDescriptor instead')
const CompletionsCustom$json = {
  '1': 'CompletionsCustom',
  '2': [
    {'1': 'use_memory', '3': 1, '4': 1, '5': 8, '10': 'useMemory'},
    {'1': 'use_multiple_agent', '3': 2, '4': 1, '5': 8, '10': 'useMultipleAgent'},
  ],
};

/// Descriptor for `CompletionsCustom`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completionsCustomDescriptor = $convert.base64Decode(
    'ChFDb21wbGV0aW9uc0N1c3RvbRIdCgp1c2VfbWVtb3J5GAEgASgIUgl1c2VNZW1vcnkSLAoSdX'
    'NlX211bHRpcGxlX2FnZW50GAIgASgIUhB1c2VNdWx0aXBsZUFnZW50');

@$core.Deprecated('Use completionsResponseDescriptor instead')
const CompletionsResponse$json = {
  '1': 'CompletionsResponse',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'object', '3': 2, '4': 1, '5': 9, '10': 'object'},
    {'1': 'created', '3': 3, '4': 1, '5': 3, '10': 'created'},
    {'1': 'model', '3': 4, '4': 1, '5': 9, '10': 'model'},
    {'1': 'choices', '3': 5, '4': 3, '5': 11, '6': '.lemon_tea.server.ChatCompletionChoice', '10': 'choices'},
    {'1': 'usage', '3': 6, '4': 1, '5': 11, '6': '.lemon_tea.common.TokenUsage', '9': 0, '10': 'usage', '17': true},
    {'1': 'system_fingerprint', '3': 7, '4': 1, '5': 9, '9': 1, '10': 'systemFingerprint', '17': true},
    {'1': 'metadata', '3': 8, '4': 3, '5': 11, '6': '.lemon_tea.server.CompletionsResponse.MetadataEntry', '10': 'metadata'},
    {'1': 'non_standard', '3': 9, '4': 1, '5': 8, '10': 'nonStandard'},
    {'1': 'chat_uuid', '3': 10, '4': 1, '5': 9, '10': 'chatUuid'},
  ],
  '3': [CompletionsResponse_MetadataEntry$json],
  '8': [
    {'1': '_usage'},
    {'1': '_system_fingerprint'},
  ],
};

@$core.Deprecated('Use completionsResponseDescriptor instead')
const CompletionsResponse_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CompletionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completionsResponseDescriptor = $convert.base64Decode(
    'ChNDb21wbGV0aW9uc1Jlc3BvbnNlEg4KAmlkGAEgASgJUgJpZBIWCgZvYmplY3QYAiABKAlSBm'
    '9iamVjdBIYCgdjcmVhdGVkGAMgASgDUgdjcmVhdGVkEhQKBW1vZGVsGAQgASgJUgVtb2RlbBJA'
    'CgdjaG9pY2VzGAUgAygLMiYubGVtb25fdGVhLnNlcnZlci5DaGF0Q29tcGxldGlvbkNob2ljZV'
    'IHY2hvaWNlcxI3CgV1c2FnZRgGIAEoCzIcLmxlbW9uX3RlYS5jb21tb24uVG9rZW5Vc2FnZUgA'
    'UgV1c2FnZYgBARIyChJzeXN0ZW1fZmluZ2VycHJpbnQYByABKAlIAVIRc3lzdGVtRmluZ2VycH'
    'JpbnSIAQESTwoIbWV0YWRhdGEYCCADKAsyMy5sZW1vbl90ZWEuc2VydmVyLkNvbXBsZXRpb25z'
    'UmVzcG9uc2UuTWV0YWRhdGFFbnRyeVIIbWV0YWRhdGESIQoMbm9uX3N0YW5kYXJkGAkgASgIUg'
    'tub25TdGFuZGFyZBIbCgljaGF0X3V1aWQYCiABKAlSCGNoYXRVdWlkGjsKDU1ldGFkYXRhRW50'
    'cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AUIICgZfdXNhZ2'
    'VCFQoTX3N5c3RlbV9maW5nZXJwcmludA==');

@$core.Deprecated('Use chatTitleRequestDescriptor instead')
const ChatTitleRequest$json = {
  '1': 'ChatTitleRequest',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
  ],
};

/// Descriptor for `ChatTitleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatTitleRequestDescriptor = $convert.base64Decode(
    'ChBDaGF0VGl0bGVSZXF1ZXN0EhsKCWNoYXRfdXVpZBgBIAEoCVIIY2hhdFV1aWQ=');

@$core.Deprecated('Use chatTitleResponseDescriptor instead')
const ChatTitleResponse$json = {
  '1': 'ChatTitleResponse',
  '2': [
    {'1': 'title', '3': 1, '4': 1, '5': 9, '10': 'title'},
  ],
};

/// Descriptor for `ChatTitleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatTitleResponseDescriptor = $convert.base64Decode(
    'ChFDaGF0VGl0bGVSZXNwb25zZRIUCgV0aXRsZRgBIAEoCVIFdGl0bGU=');

@$core.Deprecated('Use chatTitleSaveRequestDescriptor instead')
const ChatTitleSaveRequest$json = {
  '1': 'ChatTitleSaveRequest',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
    {'1': 'chat_title', '3': 2, '4': 1, '5': 9, '10': 'chatTitle'},
  ],
};

/// Descriptor for `ChatTitleSaveRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatTitleSaveRequestDescriptor = $convert.base64Decode(
    'ChRDaGF0VGl0bGVTYXZlUmVxdWVzdBIbCgljaGF0X3V1aWQYASABKAlSCGNoYXRVdWlkEh0KCm'
    'NoYXRfdGl0bGUYAiABKAlSCWNoYXRUaXRsZQ==');

@$core.Deprecated('Use listChatsRequestDescriptor instead')
const ListChatsRequest$json = {
  '1': 'ListChatsRequest',
  '2': [
    {'1': 'offset', '3': 1, '4': 1, '5': 3, '10': 'offset'},
    {'1': 'limit', '3': 2, '4': 1, '5': 3, '10': 'limit'},
    {'1': 'filter', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.server.ListChatsFilter', '9': 0, '10': 'filter', '17': true},
  ],
  '8': [
    {'1': '_filter'},
  ],
};

/// Descriptor for `ListChatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChatsRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0Q2hhdHNSZXF1ZXN0EhYKBm9mZnNldBgBIAEoA1IGb2Zmc2V0EhQKBWxpbWl0GAIgAS'
    'gDUgVsaW1pdBI+CgZmaWx0ZXIYAyABKAsyIS5sZW1vbl90ZWEuc2VydmVyLkxpc3RDaGF0c0Zp'
    'bHRlckgAUgZmaWx0ZXKIAQFCCQoHX2ZpbHRlcg==');

@$core.Deprecated('Use listChatsFilterDescriptor instead')
const ListChatsFilter$json = {
  '1': 'ListChatsFilter',
  '2': [
    {'1': 'tag', '3': 1, '4': 1, '5': 9, '9': 0, '10': 'tag', '17': true},
    {'1': 'keyword', '3': 2, '4': 1, '5': 9, '9': 1, '10': 'keyword', '17': true},
  ],
  '8': [
    {'1': '_tag'},
    {'1': '_keyword'},
  ],
};

/// Descriptor for `ListChatsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChatsFilterDescriptor = $convert.base64Decode(
    'Cg9MaXN0Q2hhdHNGaWx0ZXISFQoDdGFnGAEgASgJSABSA3RhZ4gBARIdCgdrZXl3b3JkGAIgAS'
    'gJSAFSB2tleXdvcmSIAQFCBgoEX3RhZ0IKCghfa2V5d29yZA==');

@$core.Deprecated('Use listChatsResponseDescriptor instead')
const ListChatsResponse$json = {
  '1': 'ListChatsResponse',
  '2': [
    {'1': 'chats', '3': 1, '4': 3, '5': 11, '6': '.lemon_tea.common.ChatInfo', '10': 'chats'},
    {'1': 'total_count', '3': 2, '4': 1, '5': 3, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListChatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChatsResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0Q2hhdHNSZXNwb25zZRIwCgVjaGF0cxgBIAMoCzIaLmxlbW9uX3RlYS5jb21tb24uQ2'
    'hhdEluZm9SBWNoYXRzEh8KC3RvdGFsX2NvdW50GAIgASgDUgp0b3RhbENvdW50');

@$core.Deprecated('Use deleteChatRequestDescriptor instead')
const DeleteChatRequest$json = {
  '1': 'DeleteChatRequest',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
  ],
};

/// Descriptor for `DeleteChatRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteChatRequestDescriptor = $convert.base64Decode(
    'ChFEZWxldGVDaGF0UmVxdWVzdBIbCgljaGF0X3V1aWQYASABKAlSCGNoYXRVdWlk');

@$core.Deprecated('Use getChatMessagesRequestDescriptor instead')
const GetChatMessagesRequest$json = {
  '1': 'GetChatMessagesRequest',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
    {'1': 'offset', '3': 2, '4': 1, '5': 3, '10': 'offset'},
    {'1': 'limit', '3': 3, '4': 1, '5': 3, '10': 'limit'},
  ],
};

/// Descriptor for `GetChatMessagesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChatMessagesRequestDescriptor = $convert.base64Decode(
    'ChZHZXRDaGF0TWVzc2FnZXNSZXF1ZXN0EhsKCWNoYXRfdXVpZBgBIAEoCVIIY2hhdFV1aWQSFg'
    'oGb2Zmc2V0GAIgASgDUgZvZmZzZXQSFAoFbGltaXQYAyABKANSBWxpbWl0');

@$core.Deprecated('Use getChatMessagesResponseDescriptor instead')
const GetChatMessagesResponse$json = {
  '1': 'GetChatMessagesResponse',
  '2': [
    {'1': 'messages', '3': 1, '4': 3, '5': 11, '6': '.lemon_tea.common.Message', '10': 'messages'},
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `GetChatMessagesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChatMessagesResponseDescriptor = $convert.base64Decode(
    'ChdHZXRDaGF0TWVzc2FnZXNSZXNwb25zZRI1CghtZXNzYWdlcxgBIAMoCzIZLmxlbW9uX3RlYS'
    '5jb21tb24uTWVzc2FnZVIIbWVzc2FnZXMSHwoLdG90YWxfY291bnQYAiABKAVSCnRvdGFsQ291'
    'bnQ=');

@$core.Deprecated('Use deleteChatMessageRequestDescriptor instead')
const DeleteChatMessageRequest$json = {
  '1': 'DeleteChatMessageRequest',
  '2': [
    {'1': 'chat_uuid', '3': 1, '4': 1, '5': 9, '10': 'chatUuid'},
    {'1': 'message_id', '3': 2, '4': 1, '5': 9, '10': 'messageId'},
  ],
};

/// Descriptor for `DeleteChatMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteChatMessageRequestDescriptor = $convert.base64Decode(
    'ChhEZWxldGVDaGF0TWVzc2FnZVJlcXVlc3QSGwoJY2hhdF91dWlkGAEgASgJUghjaGF0VXVpZB'
    'IdCgptZXNzYWdlX2lkGAIgASgJUgltZXNzYWdlSWQ=');

@$core.Deprecated('Use deleteChatMessageResponseDescriptor instead')
const DeleteChatMessageResponse$json = {
  '1': 'DeleteChatMessageResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DeleteChatMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteChatMessageResponseDescriptor = $convert.base64Decode(
    'ChlEZWxldGVDaGF0TWVzc2FnZVJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSGA'
    'oHbWVzc2FnZRgCIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use chatCompletionToolDescriptor instead')
const ChatCompletionTool$json = {
  '1': 'ChatCompletionTool',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 9, '10': 'type'},
    {'1': 'function', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.server.ChatCompletionFunction', '10': 'function'},
  ],
};

/// Descriptor for `ChatCompletionTool`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionToolDescriptor = $convert.base64Decode(
    'ChJDaGF0Q29tcGxldGlvblRvb2wSEgoEdHlwZRgBIAEoCVIEdHlwZRJECghmdW5jdGlvbhgCIA'
    'EoCzIoLmxlbW9uX3RlYS5zZXJ2ZXIuQ2hhdENvbXBsZXRpb25GdW5jdGlvblIIZnVuY3Rpb24=');

@$core.Deprecated('Use chatCompletionFunctionDescriptor instead')
const ChatCompletionFunction$json = {
  '1': 'ChatCompletionFunction',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 2, '4': 1, '5': 9, '9': 0, '10': 'description', '17': true},
    {'1': 'parameters', '3': 3, '4': 1, '5': 9, '9': 1, '10': 'parameters', '17': true},
  ],
  '8': [
    {'1': '_description'},
    {'1': '_parameters'},
  ],
};

/// Descriptor for `ChatCompletionFunction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionFunctionDescriptor = $convert.base64Decode(
    'ChZDaGF0Q29tcGxldGlvbkZ1bmN0aW9uEhIKBG5hbWUYASABKAlSBG5hbWUSJQoLZGVzY3JpcH'
    'Rpb24YAiABKAlIAFILZGVzY3JpcHRpb26IAQESIwoKcGFyYW1ldGVycxgDIAEoCUgBUgpwYXJh'
    'bWV0ZXJziAEBQg4KDF9kZXNjcmlwdGlvbkINCgtfcGFyYW1ldGVycw==');

@$core.Deprecated('Use toolChoiceDescriptor instead')
const ToolChoice$json = {
  '1': 'ToolChoice',
  '2': [
    {'1': 'mode', '3': 1, '4': 1, '5': 9, '9': 0, '10': 'mode'},
    {'1': 'named', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.server.ChatCompletionNamedToolChoice', '9': 0, '10': 'named'},
  ],
  '8': [
    {'1': 'choice'},
  ],
};

/// Descriptor for `ToolChoice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List toolChoiceDescriptor = $convert.base64Decode(
    'CgpUb29sQ2hvaWNlEhQKBG1vZGUYASABKAlIAFIEbW9kZRJHCgVuYW1lZBgCIAEoCzIvLmxlbW'
    '9uX3RlYS5zZXJ2ZXIuQ2hhdENvbXBsZXRpb25OYW1lZFRvb2xDaG9pY2VIAFIFbmFtZWRCCAoG'
    'Y2hvaWNl');

@$core.Deprecated('Use chatCompletionNamedToolChoiceDescriptor instead')
const ChatCompletionNamedToolChoice$json = {
  '1': 'ChatCompletionNamedToolChoice',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 9, '10': 'type'},
    {'1': 'function', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.server.ChatCompletionFunction', '10': 'function'},
  ],
};

/// Descriptor for `ChatCompletionNamedToolChoice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionNamedToolChoiceDescriptor = $convert.base64Decode(
    'Ch1DaGF0Q29tcGxldGlvbk5hbWVkVG9vbENob2ljZRISCgR0eXBlGAEgASgJUgR0eXBlEkQKCG'
    'Z1bmN0aW9uGAIgASgLMigubGVtb25fdGVhLnNlcnZlci5DaGF0Q29tcGxldGlvbkZ1bmN0aW9u'
    'UghmdW5jdGlvbg==');

@$core.Deprecated('Use responseFormatDescriptor instead')
const ResponseFormat$json = {
  '1': 'ResponseFormat',
  '2': [
    {'1': 'type', '3': 1, '4': 1, '5': 9, '10': 'type'},
    {'1': 'json_schema', '3': 2, '4': 1, '5': 9, '9': 0, '10': 'jsonSchema', '17': true},
  ],
  '8': [
    {'1': '_json_schema'},
  ],
};

/// Descriptor for `ResponseFormat`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List responseFormatDescriptor = $convert.base64Decode(
    'Cg5SZXNwb25zZUZvcm1hdBISCgR0eXBlGAEgASgJUgR0eXBlEiQKC2pzb25fc2NoZW1hGAIgAS'
    'gJSABSCmpzb25TY2hlbWGIAQFCDgoMX2pzb25fc2NoZW1h');

@$core.Deprecated('Use chatCompletionChoiceDescriptor instead')
const ChatCompletionChoice$json = {
  '1': 'ChatCompletionChoice',
  '2': [
    {'1': 'index', '3': 1, '4': 1, '5': 5, '10': 'index'},
    {'1': 'delta', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.common.Message', '10': 'delta'},
    {'1': 'logprobs', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.common.LogProbs', '9': 0, '10': 'logprobs', '17': true},
    {'1': 'finish_reason', '3': 4, '4': 1, '5': 9, '9': 1, '10': 'finishReason', '17': true},
  ],
  '8': [
    {'1': '_logprobs'},
    {'1': '_finish_reason'},
  ],
};

/// Descriptor for `ChatCompletionChoice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionChoiceDescriptor = $convert.base64Decode(
    'ChRDaGF0Q29tcGxldGlvbkNob2ljZRIUCgVpbmRleBgBIAEoBVIFaW5kZXgSLwoFZGVsdGEYAi'
    'ABKAsyGS5sZW1vbl90ZWEuY29tbW9uLk1lc3NhZ2VSBWRlbHRhEjsKCGxvZ3Byb2JzGAMgASgL'
    'MhoubGVtb25fdGVhLmNvbW1vbi5Mb2dQcm9ic0gAUghsb2dwcm9ic4gBARIoCg1maW5pc2hfcm'
    'Vhc29uGAQgASgJSAFSDGZpbmlzaFJlYXNvbogBAUILCglfbG9ncHJvYnNCEAoOX2ZpbmlzaF9y'
    'ZWFzb24=');

@$core.Deprecated('Use chatCompletionChunkDescriptor instead')
const ChatCompletionChunk$json = {
  '1': 'ChatCompletionChunk',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'object', '3': 2, '4': 1, '5': 9, '10': 'object'},
    {'1': 'created', '3': 3, '4': 1, '5': 3, '10': 'created'},
    {'1': 'model', '3': 4, '4': 1, '5': 9, '10': 'model'},
    {'1': 'choices', '3': 5, '4': 3, '5': 11, '6': '.lemon_tea.server.ChatCompletionChunkChoice', '10': 'choices'},
    {'1': 'usage', '3': 6, '4': 1, '5': 11, '6': '.lemon_tea.common.TokenUsage', '9': 0, '10': 'usage', '17': true},
    {'1': 'system_fingerprint', '3': 7, '4': 1, '5': 9, '9': 1, '10': 'systemFingerprint', '17': true},
  ],
  '8': [
    {'1': '_usage'},
    {'1': '_system_fingerprint'},
  ],
};

/// Descriptor for `ChatCompletionChunk`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionChunkDescriptor = $convert.base64Decode(
    'ChNDaGF0Q29tcGxldGlvbkNodW5rEg4KAmlkGAEgASgJUgJpZBIWCgZvYmplY3QYAiABKAlSBm'
    '9iamVjdBIYCgdjcmVhdGVkGAMgASgDUgdjcmVhdGVkEhQKBW1vZGVsGAQgASgJUgVtb2RlbBJF'
    'CgdjaG9pY2VzGAUgAygLMisubGVtb25fdGVhLnNlcnZlci5DaGF0Q29tcGxldGlvbkNodW5rQ2'
    'hvaWNlUgdjaG9pY2VzEjcKBXVzYWdlGAYgASgLMhwubGVtb25fdGVhLmNvbW1vbi5Ub2tlblVz'
    'YWdlSABSBXVzYWdliAEBEjIKEnN5c3RlbV9maW5nZXJwcmludBgHIAEoCUgBUhFzeXN0ZW1GaW'
    '5nZXJwcmludIgBAUIICgZfdXNhZ2VCFQoTX3N5c3RlbV9maW5nZXJwcmludA==');

@$core.Deprecated('Use chatCompletionChunkChoiceDescriptor instead')
const ChatCompletionChunkChoice$json = {
  '1': 'ChatCompletionChunkChoice',
  '2': [
    {'1': 'index', '3': 1, '4': 1, '5': 5, '10': 'index'},
    {'1': 'delta', '3': 2, '4': 1, '5': 11, '6': '.lemon_tea.server.ChatCompletionChunkDelta', '9': 0, '10': 'delta', '17': true},
    {'1': 'logprobs', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.common.LogProbs', '9': 1, '10': 'logprobs', '17': true},
    {'1': 'finish_reason', '3': 4, '4': 1, '5': 9, '9': 2, '10': 'finishReason', '17': true},
  ],
  '8': [
    {'1': '_delta'},
    {'1': '_logprobs'},
    {'1': '_finish_reason'},
  ],
};

/// Descriptor for `ChatCompletionChunkChoice`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionChunkChoiceDescriptor = $convert.base64Decode(
    'ChlDaGF0Q29tcGxldGlvbkNodW5rQ2hvaWNlEhQKBWluZGV4GAEgASgFUgVpbmRleBJFCgVkZW'
    'x0YRgCIAEoCzIqLmxlbW9uX3RlYS5zZXJ2ZXIuQ2hhdENvbXBsZXRpb25DaHVua0RlbHRhSABS'
    'BWRlbHRhiAEBEjsKCGxvZ3Byb2JzGAMgASgLMhoubGVtb25fdGVhLmNvbW1vbi5Mb2dQcm9ic0'
    'gBUghsb2dwcm9ic4gBARIoCg1maW5pc2hfcmVhc29uGAQgASgJSAJSDGZpbmlzaFJlYXNvbogB'
    'AUIICgZfZGVsdGFCCwoJX2xvZ3Byb2JzQhAKDl9maW5pc2hfcmVhc29u');

@$core.Deprecated('Use chatCompletionChunkDeltaDescriptor instead')
const ChatCompletionChunkDelta$json = {
  '1': 'ChatCompletionChunkDelta',
  '2': [
    {'1': 'role', '3': 1, '4': 1, '5': 9, '9': 0, '10': 'role', '17': true},
    {'1': 'content', '3': 2, '4': 1, '5': 9, '9': 1, '10': 'content', '17': true},
    {'1': 'tool_calls', '3': 3, '4': 3, '5': 11, '6': '.lemon_tea.common.ToolCall', '10': 'toolCalls'},
  ],
  '8': [
    {'1': '_role'},
    {'1': '_content'},
  ],
};

/// Descriptor for `ChatCompletionChunkDelta`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatCompletionChunkDeltaDescriptor = $convert.base64Decode(
    'ChhDaGF0Q29tcGxldGlvbkNodW5rRGVsdGESFwoEcm9sZRgBIAEoCUgAUgRyb2xliAEBEh0KB2'
    'NvbnRlbnQYAiABKAlIAVIHY29udGVudIgBARI5Cgp0b29sX2NhbGxzGAMgAygLMhoubGVtb25f'
    'dGVhLmNvbW1vbi5Ub29sQ2FsbFIJdG9vbENhbGxzQgcKBV9yb2xlQgoKCF9jb250ZW50');

