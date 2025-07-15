// This is a generated file - do not edit.
//
// Generated from rpc/service.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'common.pb.dart' as $1;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class UpdateLlmConfigRequest extends $pb.GeneratedMessage {
  factory UpdateLlmConfigRequest({
    $core.Iterable<$1.LlmProvider>? llmProviders,
  }) {
    final result = create();
    if (llmProviders != null) result.llmProviders.addAll(llmProviders);
    return result;
  }

  UpdateLlmConfigRequest._();

  factory UpdateLlmConfigRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory UpdateLlmConfigRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateLlmConfigRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..pc<$1.LlmProvider>(1, _omitFieldNames ? '' : 'llmProviders', $pb.PbFieldType.PM, subBuilder: $1.LlmProvider.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateLlmConfigRequest clone() => UpdateLlmConfigRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateLlmConfigRequest copyWith(void Function(UpdateLlmConfigRequest) updates) => super.copyWith((message) => updates(message as UpdateLlmConfigRequest)) as UpdateLlmConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateLlmConfigRequest create() => UpdateLlmConfigRequest._();
  @$core.override
  UpdateLlmConfigRequest createEmptyInstance() => create();
  static $pb.PbList<UpdateLlmConfigRequest> createRepeated() => $pb.PbList<UpdateLlmConfigRequest>();
  @$core.pragma('dart2js:noInline')
  static UpdateLlmConfigRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateLlmConfigRequest>(create);
  static UpdateLlmConfigRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$1.LlmProvider> get llmProviders => $_getList(0);
}

class UpdateLlmConfigResponse extends $pb.GeneratedMessage {
  factory UpdateLlmConfigResponse() => create();

  UpdateLlmConfigResponse._();

  factory UpdateLlmConfigResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory UpdateLlmConfigResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateLlmConfigResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateLlmConfigResponse clone() => UpdateLlmConfigResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateLlmConfigResponse copyWith(void Function(UpdateLlmConfigResponse) updates) => super.copyWith((message) => updates(message as UpdateLlmConfigResponse)) as UpdateLlmConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateLlmConfigResponse create() => UpdateLlmConfigResponse._();
  @$core.override
  UpdateLlmConfigResponse createEmptyInstance() => create();
  static $pb.PbList<UpdateLlmConfigResponse> createRepeated() => $pb.PbList<UpdateLlmConfigResponse>();
  @$core.pragma('dart2js:noInline')
  static UpdateLlmConfigResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateLlmConfigResponse>(create);
  static UpdateLlmConfigResponse? _defaultInstance;
}

class ModelsRequest extends $pb.GeneratedMessage {
  factory ModelsRequest({
    $core.String? name,
    $core.String? baseUrl,
    $core.String? apiKey,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (baseUrl != null) result.baseUrl = baseUrl;
    if (apiKey != null) result.apiKey = apiKey;
    return result;
  }

  ModelsRequest._();

  factory ModelsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ModelsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ModelsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

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
    $core.Iterable<$1.Model>? models,
  }) {
    final result = create();
    if (models != null) result.models.addAll(models);
    return result;
  }

  ModelsResponse._();

  factory ModelsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ModelsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ModelsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..pc<$1.Model>(1, _omitFieldNames ? '' : 'models', $pb.PbFieldType.PM, subBuilder: $1.Model.create)
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
  $pb.PbList<$1.Model> get models => $_getList(0);
}

class ChatRequest extends $pb.GeneratedMessage {
  factory ChatRequest({
    $core.String? llmProviderId,
    $core.String? modelId,
    $core.Iterable<$1.Message>? historyMessages,
    $1.Message? message,
    $core.String? prompt,
  }) {
    final result = create();
    if (llmProviderId != null) result.llmProviderId = llmProviderId;
    if (modelId != null) result.modelId = modelId;
    if (historyMessages != null) result.historyMessages.addAll(historyMessages);
    if (message != null) result.message = message;
    if (prompt != null) result.prompt = prompt;
    return result;
  }

  ChatRequest._();

  factory ChatRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'llmProviderId')
    ..aOS(2, _omitFieldNames ? '' : 'modelId')
    ..pc<$1.Message>(3, _omitFieldNames ? '' : 'historyMessages', $pb.PbFieldType.PM, subBuilder: $1.Message.create)
    ..aOM<$1.Message>(5, _omitFieldNames ? '' : 'message', subBuilder: $1.Message.create)
    ..aOS(6, _omitFieldNames ? '' : 'prompt')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatRequest clone() => ChatRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatRequest copyWith(void Function(ChatRequest) updates) => super.copyWith((message) => updates(message as ChatRequest)) as ChatRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatRequest create() => ChatRequest._();
  @$core.override
  ChatRequest createEmptyInstance() => create();
  static $pb.PbList<ChatRequest> createRepeated() => $pb.PbList<ChatRequest>();
  @$core.pragma('dart2js:noInline')
  static ChatRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatRequest>(create);
  static ChatRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get llmProviderId => $_getSZ(0);
  @$pb.TagNumber(1)
  set llmProviderId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLlmProviderId() => $_has(0);
  @$pb.TagNumber(1)
  void clearLlmProviderId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get modelId => $_getSZ(1);
  @$pb.TagNumber(2)
  set modelId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasModelId() => $_has(1);
  @$pb.TagNumber(2)
  void clearModelId() => $_clearField(2);

  @$pb.TagNumber(3)
  $pb.PbList<$1.Message> get historyMessages => $_getList(2);

  @$pb.TagNumber(5)
  $1.Message get message => $_getN(3);
  @$pb.TagNumber(5)
  set message($1.Message value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasMessage() => $_has(3);
  @$pb.TagNumber(5)
  void clearMessage() => $_clearField(5);
  @$pb.TagNumber(5)
  $1.Message ensureMessage() => $_ensure(3);

  @$pb.TagNumber(6)
  $core.String get prompt => $_getSZ(4);
  @$pb.TagNumber(6)
  set prompt($core.String value) => $_setString(4, value);
  @$pb.TagNumber(6)
  $core.bool hasPrompt() => $_has(4);
  @$pb.TagNumber(6)
  void clearPrompt() => $_clearField(6);
}

class ChatResponse extends $pb.GeneratedMessage {
  factory ChatResponse({
    $core.String? content,
    $core.bool? isDone,
    $core.String? requestId,
    $core.String? errorMessage,
  }) {
    final result = create();
    if (content != null) result.content = content;
    if (isDone != null) result.isDone = isDone;
    if (requestId != null) result.requestId = requestId;
    if (errorMessage != null) result.errorMessage = errorMessage;
    return result;
  }

  ChatResponse._();

  factory ChatResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'content')
    ..aOB(2, _omitFieldNames ? '' : 'isDone')
    ..aOS(3, _omitFieldNames ? '' : 'requestId')
    ..aOS(4, _omitFieldNames ? '' : 'errorMessage')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatResponse clone() => ChatResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatResponse copyWith(void Function(ChatResponse) updates) => super.copyWith((message) => updates(message as ChatResponse)) as ChatResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatResponse create() => ChatResponse._();
  @$core.override
  ChatResponse createEmptyInstance() => create();
  static $pb.PbList<ChatResponse> createRepeated() => $pb.PbList<ChatResponse>();
  @$core.pragma('dart2js:noInline')
  static ChatResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatResponse>(create);
  static ChatResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get content => $_getSZ(0);
  @$pb.TagNumber(1)
  set content($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearContent() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get isDone => $_getBF(1);
  @$pb.TagNumber(2)
  set isDone($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIsDone() => $_has(1);
  @$pb.TagNumber(2)
  void clearIsDone() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get requestId => $_getSZ(2);
  @$pb.TagNumber(3)
  set requestId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRequestId() => $_has(2);
  @$pb.TagNumber(3)
  void clearRequestId() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get errorMessage => $_getSZ(3);
  @$pb.TagNumber(4)
  set errorMessage($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasErrorMessage() => $_has(3);
  @$pb.TagNumber(4)
  void clearErrorMessage() => $_clearField(4);
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
