// This is a generated file - do not edit.
//
// Generated from rpc/service.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import 'common.pbjson.dart' as $0;

@$core.Deprecated('Use updateLlmConfigRequestDescriptor instead')
const UpdateLlmConfigRequest$json = {
  '1': 'UpdateLlmConfigRequest',
  '2': [
    {'1': 'llm_providers', '3': 1, '4': 3, '5': 11, '6': '.lemon_tea.common.LlmProvider', '10': 'llmProviders'},
  ],
};

/// Descriptor for `UpdateLlmConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateLlmConfigRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVMbG1Db25maWdSZXF1ZXN0EkIKDWxsbV9wcm92aWRlcnMYASADKAsyHS5sZW1vbl'
    '90ZWEuY29tbW9uLkxsbVByb3ZpZGVyUgxsbG1Qcm92aWRlcnM=');

@$core.Deprecated('Use updateLlmConfigResponseDescriptor instead')
const UpdateLlmConfigResponse$json = {
  '1': 'UpdateLlmConfigResponse',
};

/// Descriptor for `UpdateLlmConfigResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateLlmConfigResponseDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVMbG1Db25maWdSZXNwb25zZQ==');

@$core.Deprecated('Use modelsRequestDescriptor instead')
const ModelsRequest$json = {
  '1': 'ModelsRequest',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'base_url', '3': 2, '4': 1, '5': 9, '10': 'baseUrl'},
    {'1': 'api_key', '3': 3, '4': 1, '5': 9, '10': 'apiKey'},
  ],
};

/// Descriptor for `ModelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelsRequestDescriptor = $convert.base64Decode(
    'Cg1Nb2RlbHNSZXF1ZXN0EhIKBG5hbWUYASABKAlSBG5hbWUSGQoIYmFzZV91cmwYAiABKAlSB2'
    'Jhc2VVcmwSFwoHYXBpX2tleRgDIAEoCVIGYXBpS2V5');

@$core.Deprecated('Use modelsResponseDescriptor instead')
const ModelsResponse$json = {
  '1': 'ModelsResponse',
  '2': [
    {'1': 'models', '3': 1, '4': 3, '5': 11, '6': '.lemon_tea.common.Model', '10': 'models'},
  ],
};

/// Descriptor for `ModelsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelsResponseDescriptor = $convert.base64Decode(
    'Cg5Nb2RlbHNSZXNwb25zZRIvCgZtb2RlbHMYASADKAsyFy5sZW1vbl90ZWEuY29tbW9uLk1vZG'
    'VsUgZtb2RlbHM=');

@$core.Deprecated('Use chatRequestDescriptor instead')
const ChatRequest$json = {
  '1': 'ChatRequest',
  '2': [
    {'1': 'llm_provider_id', '3': 1, '4': 1, '5': 9, '10': 'llmProviderId'},
    {'1': 'model_id', '3': 2, '4': 1, '5': 9, '10': 'modelId'},
    {'1': 'history_messages', '3': 3, '4': 3, '5': 11, '6': '.lemon_tea.common.Message', '10': 'historyMessages'},
    {'1': 'message', '3': 5, '4': 1, '5': 11, '6': '.lemon_tea.common.Message', '10': 'message'},
    {'1': 'prompt', '3': 6, '4': 1, '5': 9, '10': 'prompt'},
  ],
};

/// Descriptor for `ChatRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatRequestDescriptor = $convert.base64Decode(
    'CgtDaGF0UmVxdWVzdBImCg9sbG1fcHJvdmlkZXJfaWQYASABKAlSDWxsbVByb3ZpZGVySWQSGQ'
    'oIbW9kZWxfaWQYAiABKAlSB21vZGVsSWQSRAoQaGlzdG9yeV9tZXNzYWdlcxgDIAMoCzIZLmxl'
    'bW9uX3RlYS5jb21tb24uTWVzc2FnZVIPaGlzdG9yeU1lc3NhZ2VzEjMKB21lc3NhZ2UYBSABKA'
    'syGS5sZW1vbl90ZWEuY29tbW9uLk1lc3NhZ2VSB21lc3NhZ2USFgoGcHJvbXB0GAYgASgJUgZw'
    'cm9tcHQ=');

@$core.Deprecated('Use chatResponseDescriptor instead')
const ChatResponse$json = {
  '1': 'ChatResponse',
  '2': [
    {'1': 'content', '3': 1, '4': 1, '5': 9, '10': 'content'},
    {'1': 'is_done', '3': 2, '4': 1, '5': 8, '10': 'isDone'},
    {'1': 'request_id', '3': 3, '4': 1, '5': 9, '10': 'requestId'},
    {'1': 'error_message', '3': 4, '4': 1, '5': 9, '10': 'errorMessage'},
  ],
};

/// Descriptor for `ChatResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatResponseDescriptor = $convert.base64Decode(
    'CgxDaGF0UmVzcG9uc2USGAoHY29udGVudBgBIAEoCVIHY29udGVudBIXCgdpc19kb25lGAIgAS'
    'gIUgZpc0RvbmUSHQoKcmVxdWVzdF9pZBgDIAEoCVIJcmVxdWVzdElkEiMKDWVycm9yX21lc3Nh'
    'Z2UYBCABKAlSDGVycm9yTWVzc2FnZQ==');

const $core.Map<$core.String, $core.dynamic> LemonTeaServiceBase$json = {
  '1': 'LemonTea',
  '2': [
    {'1': 'UpdateLlmConfig', '2': '.lemon_tea.server.UpdateLlmConfigRequest', '3': '.lemon_tea.server.UpdateLlmConfigResponse'},
    {'1': 'Models', '2': '.lemon_tea.server.ModelsRequest', '3': '.lemon_tea.server.ModelsResponse'},
    {'1': 'Chat', '2': '.lemon_tea.server.ChatRequest', '3': '.lemon_tea.server.ChatResponse', '5': true, '6': true},
  ],
};

@$core.Deprecated('Use lemonTeaServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> LemonTeaServiceBase$messageJson = {
  '.lemon_tea.server.UpdateLlmConfigRequest': UpdateLlmConfigRequest$json,
  '.lemon_tea.common.LlmProvider': $0.LlmProvider$json,
  '.lemon_tea.common.Model': $0.Model$json,
  '.lemon_tea.server.UpdateLlmConfigResponse': UpdateLlmConfigResponse$json,
  '.lemon_tea.server.ModelsRequest': ModelsRequest$json,
  '.lemon_tea.server.ModelsResponse': ModelsResponse$json,
  '.lemon_tea.server.ChatRequest': ChatRequest$json,
  '.lemon_tea.common.Message': $0.Message$json,
  '.lemon_tea.server.ChatResponse': ChatResponse$json,
};

/// Descriptor for `LemonTea`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List lemonTeaServiceDescriptor = $convert.base64Decode(
    'CghMZW1vblRlYRJmCg9VcGRhdGVMbG1Db25maWcSKC5sZW1vbl90ZWEuc2VydmVyLlVwZGF0ZU'
    'xsbUNvbmZpZ1JlcXVlc3QaKS5sZW1vbl90ZWEuc2VydmVyLlVwZGF0ZUxsbUNvbmZpZ1Jlc3Bv'
    'bnNlEksKBk1vZGVscxIfLmxlbW9uX3RlYS5zZXJ2ZXIuTW9kZWxzUmVxdWVzdBogLmxlbW9uX3'
    'RlYS5zZXJ2ZXIuTW9kZWxzUmVzcG9uc2USSQoEQ2hhdBIdLmxlbW9uX3RlYS5zZXJ2ZXIuQ2hh'
    'dFJlcXVlc3QaHi5sZW1vbl90ZWEuc2VydmVyLkNoYXRSZXNwb25zZSgBMAE=');

