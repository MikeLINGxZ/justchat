// This is a generated file - do not edit.
//
// Generated from rpc/service/models.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use modelsRequestDescriptor instead')
const ModelsRequest$json = {
  '1': 'ModelsRequest',
  '2': [
    {'1': 'llm_provider_id', '3': 1, '4': 1, '5': 9, '10': 'llmProviderId'},
    {'1': 'base_url', '3': 2, '4': 1, '5': 9, '10': 'baseUrl'},
    {'1': 'api_key', '3': 3, '4': 1, '5': 9, '10': 'apiKey'},
  ],
};

/// Descriptor for `ModelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelsRequestDescriptor = $convert.base64Decode(
    'Cg1Nb2RlbHNSZXF1ZXN0EiYKD2xsbV9wcm92aWRlcl9pZBgBIAEoCVINbGxtUHJvdmlkZXJJZB'
    'IZCghiYXNlX3VybBgCIAEoCVIHYmFzZVVybBIXCgdhcGlfa2V5GAMgASgJUgZhcGlLZXk=');

@$core.Deprecated('Use modelsResponseDescriptor instead')
const ModelsResponse$json = {
  '1': 'ModelsResponse',
  '2': [
    {'1': 'object', '3': 1, '4': 1, '5': 9, '10': 'object'},
    {'1': 'data', '3': 2, '4': 3, '5': 11, '6': '.lemon_tea.common.Model', '10': 'data'},
  ],
};

/// Descriptor for `ModelsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelsResponseDescriptor = $convert.base64Decode(
    'Cg5Nb2RlbHNSZXNwb25zZRIWCgZvYmplY3QYASABKAlSBm9iamVjdBIrCgRkYXRhGAIgAygLMh'
    'cubGVtb25fdGVhLmNvbW1vbi5Nb2RlbFIEZGF0YQ==');

