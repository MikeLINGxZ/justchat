// This is a generated file - do not edit.
//
// Generated from rpc/service/models.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import '../common/common.pb.dart' as $1;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// ListModelsRequest 允许通过 providerId 或直接 base_url/api_key 指定目标提供方
class ModelsRequest extends $pb.GeneratedMessage {
  factory ModelsRequest({
    $core.String? llmProviderId,
    $core.String? baseUrl,
    $core.String? apiKey,
  }) {
    final result = create();
    if (llmProviderId != null) result.llmProviderId = llmProviderId;
    if (baseUrl != null) result.baseUrl = baseUrl;
    if (apiKey != null) result.apiKey = apiKey;
    return result;
  }

  ModelsRequest._();

  factory ModelsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ModelsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ModelsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'llmProviderId')
    ..aOS(2, _omitFieldNames ? '' : 'baseUrl')
    ..aOS(3, _omitFieldNames ? '' : 'apiKey')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ModelsRequest clone() => ModelsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ModelsRequest copyWith(void Function(ModelsRequest) updates) => super.copyWith((message) => updates(message as ModelsRequest)) as ModelsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ModelsRequest create() => ModelsRequest._();
  @$core.override
  ModelsRequest createEmptyInstance() => create();
  static $pb.PbList<ModelsRequest> createRepeated() => $pb.PbList<ModelsRequest>();
  @$core.pragma('dart2js:noInline')
  static ModelsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ModelsRequest>(create);
  static ModelsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get llmProviderId => $_getSZ(0);
  @$pb.TagNumber(1)
  set llmProviderId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLlmProviderId() => $_has(0);
  @$pb.TagNumber(1)
  void clearLlmProviderId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get baseUrl => $_getSZ(1);
  @$pb.TagNumber(2)
  set baseUrl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBaseUrl() => $_has(1);
  @$pb.TagNumber(2)
  void clearBaseUrl() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get apiKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set apiKey($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasApiKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearApiKey() => $_clearField(3);
}

class ModelsResponse extends $pb.GeneratedMessage {
  factory ModelsResponse({
    $core.String? object,
    $core.Iterable<$1.Model>? data,
  }) {
    final result = create();
    if (object != null) result.object = object;
    if (data != null) result.data.addAll(data);
    return result;
  }

  ModelsResponse._();

  factory ModelsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ModelsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ModelsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'object')
    ..pc<$1.Model>(2, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: $1.Model.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ModelsResponse clone() => ModelsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ModelsResponse copyWith(void Function(ModelsResponse) updates) => super.copyWith((message) => updates(message as ModelsResponse)) as ModelsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ModelsResponse create() => ModelsResponse._();
  @$core.override
  ModelsResponse createEmptyInstance() => create();
  static $pb.PbList<ModelsResponse> createRepeated() => $pb.PbList<ModelsResponse>();
  @$core.pragma('dart2js:noInline')
  static ModelsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ModelsResponse>(create);
  static ModelsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get object => $_getSZ(0);
  @$pb.TagNumber(1)
  set object($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasObject() => $_has(0);
  @$pb.TagNumber(1)
  void clearObject() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$1.Model> get data => $_getList(1);
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
