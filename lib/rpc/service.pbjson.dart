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
    {'1': 'reasoning_content', '3': 2, '4': 1, '5': 9, '10': 'reasoningContent'},
    {'1': 'is_done', '3': 3, '4': 1, '5': 8, '10': 'isDone'},
    {'1': 'request_id', '3': 4, '4': 1, '5': 9, '10': 'requestId'},
    {'1': 'error_message', '3': 5, '4': 1, '5': 9, '10': 'errorMessage'},
  ],
};

/// Descriptor for `ChatResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatResponseDescriptor = $convert.base64Decode(
    'CgxDaGF0UmVzcG9uc2USGAoHY29udGVudBgBIAEoCVIHY29udGVudBIrChFyZWFzb25pbmdfY2'
    '9udGVudBgCIAEoCVIQcmVhc29uaW5nQ29udGVudBIXCgdpc19kb25lGAMgASgIUgZpc0RvbmUS'
    'HQoKcmVxdWVzdF9pZBgEIAEoCVIJcmVxdWVzdElkEiMKDWVycm9yX21lc3NhZ2UYBSABKAlSDG'
    'Vycm9yTWVzc2FnZQ==');

