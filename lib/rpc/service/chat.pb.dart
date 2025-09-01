// This is a generated file - do not edit.
//
// Generated from rpc/service/chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import '../common/common.pb.dart' as $1;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// CreateChatCompletionRequest 与 OpenAI Chat Completions 请求体对齐
class CompletionsRequest extends $pb.GeneratedMessage {
  factory CompletionsRequest({
    $core.String? model,
    $core.Iterable<$1.Message>? messages,
    $core.double? temperature,
    $core.int? maxTokens,
    $core.double? topP,
    $core.int? n,
    $core.bool? stream,
    $core.String? stop,
    $core.Iterable<$core.String>? stopSequence,
    $core.double? presencePenalty,
    $core.double? frequencyPenalty,
    $core.double? repetitionPenalty,
    $core.String? user,
    $core.Iterable<ChatCompletionTool>? tools,
    ToolChoice? toolChoice,
    ResponseFormat? responseFormat,
    $core.int? seed,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? metadata,
    $core.bool? nonStandard,
    $core.String? chatUuid,
    CompletionsCustom? completionsCustom,
  }) {
    final result = create();
    if (model != null) result.model = model;
    if (messages != null) result.messages.addAll(messages);
    if (temperature != null) result.temperature = temperature;
    if (maxTokens != null) result.maxTokens = maxTokens;
    if (topP != null) result.topP = topP;
    if (n != null) result.n = n;
    if (stream != null) result.stream = stream;
    if (stop != null) result.stop = stop;
    if (stopSequence != null) result.stopSequence.addAll(stopSequence);
    if (presencePenalty != null) result.presencePenalty = presencePenalty;
    if (frequencyPenalty != null) result.frequencyPenalty = frequencyPenalty;
    if (repetitionPenalty != null) result.repetitionPenalty = repetitionPenalty;
    if (user != null) result.user = user;
    if (tools != null) result.tools.addAll(tools);
    if (toolChoice != null) result.toolChoice = toolChoice;
    if (responseFormat != null) result.responseFormat = responseFormat;
    if (seed != null) result.seed = seed;
    if (metadata != null) result.metadata.addEntries(metadata);
    if (nonStandard != null) result.nonStandard = nonStandard;
    if (chatUuid != null) result.chatUuid = chatUuid;
    if (completionsCustom != null) result.completionsCustom = completionsCustom;
    return result;
  }

  CompletionsRequest._();

  factory CompletionsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory CompletionsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CompletionsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'model')
    ..pc<$1.Message>(2, _omitFieldNames ? '' : 'messages', $pb.PbFieldType.PM, subBuilder: $1.Message.create)
    ..a<$core.double>(3, _omitFieldNames ? '' : 'temperature', $pb.PbFieldType.OD)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'maxTokens', $pb.PbFieldType.O3)
    ..a<$core.double>(5, _omitFieldNames ? '' : 'topP', $pb.PbFieldType.OD)
    ..a<$core.int>(6, _omitFieldNames ? '' : 'n', $pb.PbFieldType.O3)
    ..aOB(7, _omitFieldNames ? '' : 'stream')
    ..aOS(8, _omitFieldNames ? '' : 'stop')
    ..pPS(9, _omitFieldNames ? '' : 'stopSequence')
    ..a<$core.double>(10, _omitFieldNames ? '' : 'presencePenalty', $pb.PbFieldType.OD)
    ..a<$core.double>(11, _omitFieldNames ? '' : 'frequencyPenalty', $pb.PbFieldType.OD)
    ..a<$core.double>(12, _omitFieldNames ? '' : 'repetitionPenalty', $pb.PbFieldType.OD)
    ..aOS(13, _omitFieldNames ? '' : 'user')
    ..pc<ChatCompletionTool>(14, _omitFieldNames ? '' : 'tools', $pb.PbFieldType.PM, subBuilder: ChatCompletionTool.create)
    ..aOM<ToolChoice>(15, _omitFieldNames ? '' : 'toolChoice', subBuilder: ToolChoice.create)
    ..aOM<ResponseFormat>(16, _omitFieldNames ? '' : 'responseFormat', subBuilder: ResponseFormat.create)
    ..a<$core.int>(17, _omitFieldNames ? '' : 'seed', $pb.PbFieldType.O3)
    ..m<$core.String, $core.String>(18, _omitFieldNames ? '' : 'metadata', entryClassName: 'CompletionsRequest.MetadataEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.server'))
    ..aOB(19, _omitFieldNames ? '' : 'nonStandard')
    ..aOS(20, _omitFieldNames ? '' : 'chatUuid')
    ..aOM<CompletionsCustom>(21, _omitFieldNames ? '' : 'completionsCustom', subBuilder: CompletionsCustom.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsRequest clone() => CompletionsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsRequest copyWith(void Function(CompletionsRequest) updates) => super.copyWith((message) => updates(message as CompletionsRequest)) as CompletionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CompletionsRequest create() => CompletionsRequest._();
  @$core.override
  CompletionsRequest createEmptyInstance() => create();
  static $pb.PbList<CompletionsRequest> createRepeated() => $pb.PbList<CompletionsRequest>();
  @$core.pragma('dart2js:noInline')
  static CompletionsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CompletionsRequest>(create);
  static CompletionsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get model => $_getSZ(0);
  @$pb.TagNumber(1)
  set model($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasModel() => $_has(0);
  @$pb.TagNumber(1)
  void clearModel() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$1.Message> get messages => $_getList(1);

  @$pb.TagNumber(3)
  $core.double get temperature => $_getN(2);
  @$pb.TagNumber(3)
  set temperature($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTemperature() => $_has(2);
  @$pb.TagNumber(3)
  void clearTemperature() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.int get maxTokens => $_getIZ(3);
  @$pb.TagNumber(4)
  set maxTokens($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMaxTokens() => $_has(3);
  @$pb.TagNumber(4)
  void clearMaxTokens() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.double get topP => $_getN(4);
  @$pb.TagNumber(5)
  set topP($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(5)
  $core.bool hasTopP() => $_has(4);
  @$pb.TagNumber(5)
  void clearTopP() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get n => $_getIZ(5);
  @$pb.TagNumber(6)
  set n($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasN() => $_has(5);
  @$pb.TagNumber(6)
  void clearN() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.bool get stream => $_getBF(6);
  @$pb.TagNumber(7)
  set stream($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(7)
  $core.bool hasStream() => $_has(6);
  @$pb.TagNumber(7)
  void clearStream() => $_clearField(7);

  @$pb.TagNumber(8)
  $core.String get stop => $_getSZ(7);
  @$pb.TagNumber(8)
  set stop($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasStop() => $_has(7);
  @$pb.TagNumber(8)
  void clearStop() => $_clearField(8);

  @$pb.TagNumber(9)
  $pb.PbList<$core.String> get stopSequence => $_getList(8);

  @$pb.TagNumber(10)
  $core.double get presencePenalty => $_getN(9);
  @$pb.TagNumber(10)
  set presencePenalty($core.double value) => $_setDouble(9, value);
  @$pb.TagNumber(10)
  $core.bool hasPresencePenalty() => $_has(9);
  @$pb.TagNumber(10)
  void clearPresencePenalty() => $_clearField(10);

  @$pb.TagNumber(11)
  $core.double get frequencyPenalty => $_getN(10);
  @$pb.TagNumber(11)
  set frequencyPenalty($core.double value) => $_setDouble(10, value);
  @$pb.TagNumber(11)
  $core.bool hasFrequencyPenalty() => $_has(10);
  @$pb.TagNumber(11)
  void clearFrequencyPenalty() => $_clearField(11);

  @$pb.TagNumber(12)
  $core.double get repetitionPenalty => $_getN(11);
  @$pb.TagNumber(12)
  set repetitionPenalty($core.double value) => $_setDouble(11, value);
  @$pb.TagNumber(12)
  $core.bool hasRepetitionPenalty() => $_has(11);
  @$pb.TagNumber(12)
  void clearRepetitionPenalty() => $_clearField(12);

  @$pb.TagNumber(13)
  $core.String get user => $_getSZ(12);
  @$pb.TagNumber(13)
  set user($core.String value) => $_setString(12, value);
  @$pb.TagNumber(13)
  $core.bool hasUser() => $_has(12);
  @$pb.TagNumber(13)
  void clearUser() => $_clearField(13);

  @$pb.TagNumber(14)
  $pb.PbList<ChatCompletionTool> get tools => $_getList(13);

  @$pb.TagNumber(15)
  ToolChoice get toolChoice => $_getN(14);
  @$pb.TagNumber(15)
  set toolChoice(ToolChoice value) => $_setField(15, value);
  @$pb.TagNumber(15)
  $core.bool hasToolChoice() => $_has(14);
  @$pb.TagNumber(15)
  void clearToolChoice() => $_clearField(15);
  @$pb.TagNumber(15)
  ToolChoice ensureToolChoice() => $_ensure(14);

  @$pb.TagNumber(16)
  ResponseFormat get responseFormat => $_getN(15);
  @$pb.TagNumber(16)
  set responseFormat(ResponseFormat value) => $_setField(16, value);
  @$pb.TagNumber(16)
  $core.bool hasResponseFormat() => $_has(15);
  @$pb.TagNumber(16)
  void clearResponseFormat() => $_clearField(16);
  @$pb.TagNumber(16)
  ResponseFormat ensureResponseFormat() => $_ensure(15);

  @$pb.TagNumber(17)
  $core.int get seed => $_getIZ(16);
  @$pb.TagNumber(17)
  set seed($core.int value) => $_setSignedInt32(16, value);
  @$pb.TagNumber(17)
  $core.bool hasSeed() => $_has(16);
  @$pb.TagNumber(17)
  void clearSeed() => $_clearField(17);

  @$pb.TagNumber(18)
  $pb.PbMap<$core.String, $core.String> get metadata => $_getMap(17);

  /// 以下为业务字段
  @$pb.TagNumber(19)
  $core.bool get nonStandard => $_getBF(18);
  @$pb.TagNumber(19)
  set nonStandard($core.bool value) => $_setBool(18, value);
  @$pb.TagNumber(19)
  $core.bool hasNonStandard() => $_has(18);
  @$pb.TagNumber(19)
  void clearNonStandard() => $_clearField(19);

  @$pb.TagNumber(20)
  $core.String get chatUuid => $_getSZ(19);
  @$pb.TagNumber(20)
  set chatUuid($core.String value) => $_setString(19, value);
  @$pb.TagNumber(20)
  $core.bool hasChatUuid() => $_has(19);
  @$pb.TagNumber(20)
  void clearChatUuid() => $_clearField(20);

  @$pb.TagNumber(21)
  CompletionsCustom get completionsCustom => $_getN(20);
  @$pb.TagNumber(21)
  set completionsCustom(CompletionsCustom value) => $_setField(21, value);
  @$pb.TagNumber(21)
  $core.bool hasCompletionsCustom() => $_has(20);
  @$pb.TagNumber(21)
  void clearCompletionsCustom() => $_clearField(21);
  @$pb.TagNumber(21)
  CompletionsCustom ensureCompletionsCustom() => $_ensure(20);
}

/// CompletionsCustom 业务自定义
class CompletionsCustom extends $pb.GeneratedMessage {
  factory CompletionsCustom({
    $core.bool? useMemory,
    $core.bool? useMultipleAgent,
  }) {
    final result = create();
    if (useMemory != null) result.useMemory = useMemory;
    if (useMultipleAgent != null) result.useMultipleAgent = useMultipleAgent;
    return result;
  }

  CompletionsCustom._();

  factory CompletionsCustom.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory CompletionsCustom.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CompletionsCustom', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'useMemory')
    ..aOB(2, _omitFieldNames ? '' : 'useMultipleAgent')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsCustom clone() => CompletionsCustom()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsCustom copyWith(void Function(CompletionsCustom) updates) => super.copyWith((message) => updates(message as CompletionsCustom)) as CompletionsCustom;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CompletionsCustom create() => CompletionsCustom._();
  @$core.override
  CompletionsCustom createEmptyInstance() => create();
  static $pb.PbList<CompletionsCustom> createRepeated() => $pb.PbList<CompletionsCustom>();
  @$core.pragma('dart2js:noInline')
  static CompletionsCustom getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CompletionsCustom>(create);
  static CompletionsCustom? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get useMemory => $_getBF(0);
  @$pb.TagNumber(1)
  set useMemory($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUseMemory() => $_has(0);
  @$pb.TagNumber(1)
  void clearUseMemory() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.bool get useMultipleAgent => $_getBF(1);
  @$pb.TagNumber(2)
  set useMultipleAgent($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUseMultipleAgent() => $_has(1);
  @$pb.TagNumber(2)
  void clearUseMultipleAgent() => $_clearField(2);
}

/// ChatCompletionChunk 与 OpenAI Chat Completions 流式 chunk 对齐
class CompletionsResponse extends $pb.GeneratedMessage {
  factory CompletionsResponse({
    $core.String? id,
    $core.String? object,
    $fixnum.Int64? created,
    $core.String? model,
    $core.Iterable<ChatCompletionChoice>? choices,
    $1.TokenUsage? usage,
    $core.String? systemFingerprint,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? metadata,
    $core.bool? nonStandard,
    $core.String? chatUuid,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (object != null) result.object = object;
    if (created != null) result.created = created;
    if (model != null) result.model = model;
    if (choices != null) result.choices.addAll(choices);
    if (usage != null) result.usage = usage;
    if (systemFingerprint != null) result.systemFingerprint = systemFingerprint;
    if (metadata != null) result.metadata.addEntries(metadata);
    if (nonStandard != null) result.nonStandard = nonStandard;
    if (chatUuid != null) result.chatUuid = chatUuid;
    return result;
  }

  CompletionsResponse._();

  factory CompletionsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory CompletionsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CompletionsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'object')
    ..aInt64(3, _omitFieldNames ? '' : 'created')
    ..aOS(4, _omitFieldNames ? '' : 'model')
    ..pc<ChatCompletionChoice>(5, _omitFieldNames ? '' : 'choices', $pb.PbFieldType.PM, subBuilder: ChatCompletionChoice.create)
    ..aOM<$1.TokenUsage>(6, _omitFieldNames ? '' : 'usage', subBuilder: $1.TokenUsage.create)
    ..aOS(7, _omitFieldNames ? '' : 'systemFingerprint')
    ..m<$core.String, $core.String>(8, _omitFieldNames ? '' : 'metadata', entryClassName: 'CompletionsResponse.MetadataEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.server'))
    ..aOB(9, _omitFieldNames ? '' : 'nonStandard')
    ..aOS(10, _omitFieldNames ? '' : 'chatUuid')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsResponse clone() => CompletionsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompletionsResponse copyWith(void Function(CompletionsResponse) updates) => super.copyWith((message) => updates(message as CompletionsResponse)) as CompletionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CompletionsResponse create() => CompletionsResponse._();
  @$core.override
  CompletionsResponse createEmptyInstance() => create();
  static $pb.PbList<CompletionsResponse> createRepeated() => $pb.PbList<CompletionsResponse>();
  @$core.pragma('dart2js:noInline')
  static CompletionsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CompletionsResponse>(create);
  static CompletionsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get object => $_getSZ(1);
  @$pb.TagNumber(2)
  set object($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasObject() => $_has(1);
  @$pb.TagNumber(2)
  void clearObject() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get created => $_getI64(2);
  @$pb.TagNumber(3)
  set created($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCreated() => $_has(2);
  @$pb.TagNumber(3)
  void clearCreated() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get model => $_getSZ(3);
  @$pb.TagNumber(4)
  set model($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasModel() => $_has(3);
  @$pb.TagNumber(4)
  void clearModel() => $_clearField(4);

  @$pb.TagNumber(5)
  $pb.PbList<ChatCompletionChoice> get choices => $_getList(4);

  @$pb.TagNumber(6)
  $1.TokenUsage get usage => $_getN(5);
  @$pb.TagNumber(6)
  set usage($1.TokenUsage value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasUsage() => $_has(5);
  @$pb.TagNumber(6)
  void clearUsage() => $_clearField(6);
  @$pb.TagNumber(6)
  $1.TokenUsage ensureUsage() => $_ensure(5);

  @$pb.TagNumber(7)
  $core.String get systemFingerprint => $_getSZ(6);
  @$pb.TagNumber(7)
  set systemFingerprint($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasSystemFingerprint() => $_has(6);
  @$pb.TagNumber(7)
  void clearSystemFingerprint() => $_clearField(7);

  @$pb.TagNumber(8)
  $pb.PbMap<$core.String, $core.String> get metadata => $_getMap(7);

  /// 以下为业务字段
  @$pb.TagNumber(9)
  $core.bool get nonStandard => $_getBF(8);
  @$pb.TagNumber(9)
  set nonStandard($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(9)
  $core.bool hasNonStandard() => $_has(8);
  @$pb.TagNumber(9)
  void clearNonStandard() => $_clearField(9);

  @$pb.TagNumber(10)
  $core.String get chatUuid => $_getSZ(9);
  @$pb.TagNumber(10)
  set chatUuid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasChatUuid() => $_has(9);
  @$pb.TagNumber(10)
  void clearChatUuid() => $_clearField(10);
}

/// ChatTitleRequest 获取对话标题请求
class ChatTitleRequest extends $pb.GeneratedMessage {
  factory ChatTitleRequest({
    $core.String? chatUuid,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    return result;
  }

  ChatTitleRequest._();

  factory ChatTitleRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatTitleRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatTitleRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleRequest clone() => ChatTitleRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleRequest copyWith(void Function(ChatTitleRequest) updates) => super.copyWith((message) => updates(message as ChatTitleRequest)) as ChatTitleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatTitleRequest create() => ChatTitleRequest._();
  @$core.override
  ChatTitleRequest createEmptyInstance() => create();
  static $pb.PbList<ChatTitleRequest> createRepeated() => $pb.PbList<ChatTitleRequest>();
  @$core.pragma('dart2js:noInline')
  static ChatTitleRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatTitleRequest>(create);
  static ChatTitleRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);
}

/// ChatTitleResponse 获取对话标题响应
class ChatTitleResponse extends $pb.GeneratedMessage {
  factory ChatTitleResponse({
    $core.String? title,
  }) {
    final result = create();
    if (title != null) result.title = title;
    return result;
  }

  ChatTitleResponse._();

  factory ChatTitleResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatTitleResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatTitleResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'title')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleResponse clone() => ChatTitleResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleResponse copyWith(void Function(ChatTitleResponse) updates) => super.copyWith((message) => updates(message as ChatTitleResponse)) as ChatTitleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatTitleResponse create() => ChatTitleResponse._();
  @$core.override
  ChatTitleResponse createEmptyInstance() => create();
  static $pb.PbList<ChatTitleResponse> createRepeated() => $pb.PbList<ChatTitleResponse>();
  @$core.pragma('dart2js:noInline')
  static ChatTitleResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatTitleResponse>(create);
  static ChatTitleResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get title => $_getSZ(0);
  @$pb.TagNumber(1)
  set title($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTitle() => $_has(0);
  @$pb.TagNumber(1)
  void clearTitle() => $_clearField(1);
}

/// ChatTitleSaveRequest 保存对话标题请求
class ChatTitleSaveRequest extends $pb.GeneratedMessage {
  factory ChatTitleSaveRequest({
    $core.String? chatUuid,
    $core.String? chatTitle,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    if (chatTitle != null) result.chatTitle = chatTitle;
    return result;
  }

  ChatTitleSaveRequest._();

  factory ChatTitleSaveRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatTitleSaveRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatTitleSaveRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..aOS(2, _omitFieldNames ? '' : 'chatTitle')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleSaveRequest clone() => ChatTitleSaveRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatTitleSaveRequest copyWith(void Function(ChatTitleSaveRequest) updates) => super.copyWith((message) => updates(message as ChatTitleSaveRequest)) as ChatTitleSaveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatTitleSaveRequest create() => ChatTitleSaveRequest._();
  @$core.override
  ChatTitleSaveRequest createEmptyInstance() => create();
  static $pb.PbList<ChatTitleSaveRequest> createRepeated() => $pb.PbList<ChatTitleSaveRequest>();
  @$core.pragma('dart2js:noInline')
  static ChatTitleSaveRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatTitleSaveRequest>(create);
  static ChatTitleSaveRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get chatTitle => $_getSZ(1);
  @$pb.TagNumber(2)
  set chatTitle($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasChatTitle() => $_has(1);
  @$pb.TagNumber(2)
  void clearChatTitle() => $_clearField(2);
}

/// ListChatsRequest 获取对话列表请求
class ListChatsRequest extends $pb.GeneratedMessage {
  factory ListChatsRequest({
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
    ListChatsFilter? filter,
  }) {
    final result = create();
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    if (filter != null) result.filter = filter;
    return result;
  }

  ListChatsRequest._();

  factory ListChatsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ListChatsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListChatsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'offset')
    ..aInt64(2, _omitFieldNames ? '' : 'limit')
    ..aOM<ListChatsFilter>(3, _omitFieldNames ? '' : 'filter', subBuilder: ListChatsFilter.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsRequest clone() => ListChatsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsRequest copyWith(void Function(ListChatsRequest) updates) => super.copyWith((message) => updates(message as ListChatsRequest)) as ListChatsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListChatsRequest create() => ListChatsRequest._();
  @$core.override
  ListChatsRequest createEmptyInstance() => create();
  static $pb.PbList<ListChatsRequest> createRepeated() => $pb.PbList<ListChatsRequest>();
  @$core.pragma('dart2js:noInline')
  static ListChatsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListChatsRequest>(create);
  static ListChatsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get offset => $_getI64(0);
  @$pb.TagNumber(1)
  set offset($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOffset() => $_has(0);
  @$pb.TagNumber(1)
  void clearOffset() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get limit => $_getI64(1);
  @$pb.TagNumber(2)
  set limit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(2)
  void clearLimit() => $_clearField(2);

  @$pb.TagNumber(3)
  ListChatsFilter get filter => $_getN(2);
  @$pb.TagNumber(3)
  set filter(ListChatsFilter value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasFilter() => $_has(2);
  @$pb.TagNumber(3)
  void clearFilter() => $_clearField(3);
  @$pb.TagNumber(3)
  ListChatsFilter ensureFilter() => $_ensure(2);
}

class ListChatsFilter extends $pb.GeneratedMessage {
  factory ListChatsFilter({
    $core.String? tag,
    $core.String? keyword,
  }) {
    final result = create();
    if (tag != null) result.tag = tag;
    if (keyword != null) result.keyword = keyword;
    return result;
  }

  ListChatsFilter._();

  factory ListChatsFilter.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ListChatsFilter.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListChatsFilter', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'tag')
    ..aOS(2, _omitFieldNames ? '' : 'keyword')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsFilter clone() => ListChatsFilter()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsFilter copyWith(void Function(ListChatsFilter) updates) => super.copyWith((message) => updates(message as ListChatsFilter)) as ListChatsFilter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListChatsFilter create() => ListChatsFilter._();
  @$core.override
  ListChatsFilter createEmptyInstance() => create();
  static $pb.PbList<ListChatsFilter> createRepeated() => $pb.PbList<ListChatsFilter>();
  @$core.pragma('dart2js:noInline')
  static ListChatsFilter getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListChatsFilter>(create);
  static ListChatsFilter? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get tag => $_getSZ(0);
  @$pb.TagNumber(1)
  set tag($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTag() => $_has(0);
  @$pb.TagNumber(1)
  void clearTag() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get keyword => $_getSZ(1);
  @$pb.TagNumber(2)
  set keyword($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasKeyword() => $_has(1);
  @$pb.TagNumber(2)
  void clearKeyword() => $_clearField(2);
}

/// ListChatsResponse 获取对话列表响应
class ListChatsResponse extends $pb.GeneratedMessage {
  factory ListChatsResponse({
    $core.Iterable<$1.ChatInfo>? chats,
    $fixnum.Int64? totalCount,
  }) {
    final result = create();
    if (chats != null) result.chats.addAll(chats);
    if (totalCount != null) result.totalCount = totalCount;
    return result;
  }

  ListChatsResponse._();

  factory ListChatsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ListChatsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListChatsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..pc<$1.ChatInfo>(1, _omitFieldNames ? '' : 'chats', $pb.PbFieldType.PM, subBuilder: $1.ChatInfo.create)
    ..aInt64(2, _omitFieldNames ? '' : 'totalCount')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsResponse clone() => ListChatsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListChatsResponse copyWith(void Function(ListChatsResponse) updates) => super.copyWith((message) => updates(message as ListChatsResponse)) as ListChatsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListChatsResponse create() => ListChatsResponse._();
  @$core.override
  ListChatsResponse createEmptyInstance() => create();
  static $pb.PbList<ListChatsResponse> createRepeated() => $pb.PbList<ListChatsResponse>();
  @$core.pragma('dart2js:noInline')
  static ListChatsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListChatsResponse>(create);
  static ListChatsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$1.ChatInfo> get chats => $_getList(0);

  @$pb.TagNumber(2)
  $fixnum.Int64 get totalCount => $_getI64(1);
  @$pb.TagNumber(2)
  set totalCount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTotalCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalCount() => $_clearField(2);
}

/// DeleteChatRequest 删除对话请求
class DeleteChatRequest extends $pb.GeneratedMessage {
  factory DeleteChatRequest({
    $core.String? chatUuid,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    return result;
  }

  DeleteChatRequest._();

  factory DeleteChatRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory DeleteChatRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteChatRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatRequest clone() => DeleteChatRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatRequest copyWith(void Function(DeleteChatRequest) updates) => super.copyWith((message) => updates(message as DeleteChatRequest)) as DeleteChatRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteChatRequest create() => DeleteChatRequest._();
  @$core.override
  DeleteChatRequest createEmptyInstance() => create();
  static $pb.PbList<DeleteChatRequest> createRepeated() => $pb.PbList<DeleteChatRequest>();
  @$core.pragma('dart2js:noInline')
  static DeleteChatRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteChatRequest>(create);
  static DeleteChatRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);
}

/// GetChatMessagesRequest 获取对话消息请求
class GetChatMessagesRequest extends $pb.GeneratedMessage {
  factory GetChatMessagesRequest({
    $core.String? chatUuid,
    $fixnum.Int64? offset,
    $fixnum.Int64? limit,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    if (offset != null) result.offset = offset;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetChatMessagesRequest._();

  factory GetChatMessagesRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory GetChatMessagesRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetChatMessagesRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..aInt64(2, _omitFieldNames ? '' : 'offset')
    ..aInt64(3, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChatMessagesRequest clone() => GetChatMessagesRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChatMessagesRequest copyWith(void Function(GetChatMessagesRequest) updates) => super.copyWith((message) => updates(message as GetChatMessagesRequest)) as GetChatMessagesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChatMessagesRequest create() => GetChatMessagesRequest._();
  @$core.override
  GetChatMessagesRequest createEmptyInstance() => create();
  static $pb.PbList<GetChatMessagesRequest> createRepeated() => $pb.PbList<GetChatMessagesRequest>();
  @$core.pragma('dart2js:noInline')
  static GetChatMessagesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetChatMessagesRequest>(create);
  static GetChatMessagesRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get offset => $_getI64(1);
  @$pb.TagNumber(2)
  set offset($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(2)
  void clearOffset() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get limit => $_getI64(2);
  @$pb.TagNumber(3)
  set limit($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(3)
  void clearLimit() => $_clearField(3);
}

/// GetChatMessagesResponse 获取对话消息响应
class GetChatMessagesResponse extends $pb.GeneratedMessage {
  factory GetChatMessagesResponse({
    $core.Iterable<$1.Message>? messages,
    $core.int? totalCount,
  }) {
    final result = create();
    if (messages != null) result.messages.addAll(messages);
    if (totalCount != null) result.totalCount = totalCount;
    return result;
  }

  GetChatMessagesResponse._();

  factory GetChatMessagesResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory GetChatMessagesResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetChatMessagesResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..pc<$1.Message>(1, _omitFieldNames ? '' : 'messages', $pb.PbFieldType.PM, subBuilder: $1.Message.create)
    ..a<$core.int>(2, _omitFieldNames ? '' : 'totalCount', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChatMessagesResponse clone() => GetChatMessagesResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChatMessagesResponse copyWith(void Function(GetChatMessagesResponse) updates) => super.copyWith((message) => updates(message as GetChatMessagesResponse)) as GetChatMessagesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChatMessagesResponse create() => GetChatMessagesResponse._();
  @$core.override
  GetChatMessagesResponse createEmptyInstance() => create();
  static $pb.PbList<GetChatMessagesResponse> createRepeated() => $pb.PbList<GetChatMessagesResponse>();
  @$core.pragma('dart2js:noInline')
  static GetChatMessagesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetChatMessagesResponse>(create);
  static GetChatMessagesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$1.Message> get messages => $_getList(0);

  @$pb.TagNumber(2)
  $core.int get totalCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalCount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTotalCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalCount() => $_clearField(2);
}

/// DeleteChatMessageRequest 删除对话消息请求
class DeleteChatMessageRequest extends $pb.GeneratedMessage {
  factory DeleteChatMessageRequest({
    $core.String? chatUuid,
    $core.String? messageId,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    if (messageId != null) result.messageId = messageId;
    return result;
  }

  DeleteChatMessageRequest._();

  factory DeleteChatMessageRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory DeleteChatMessageRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteChatMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..aOS(2, _omitFieldNames ? '' : 'messageId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatMessageRequest clone() => DeleteChatMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatMessageRequest copyWith(void Function(DeleteChatMessageRequest) updates) => super.copyWith((message) => updates(message as DeleteChatMessageRequest)) as DeleteChatMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteChatMessageRequest create() => DeleteChatMessageRequest._();
  @$core.override
  DeleteChatMessageRequest createEmptyInstance() => create();
  static $pb.PbList<DeleteChatMessageRequest> createRepeated() => $pb.PbList<DeleteChatMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static DeleteChatMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteChatMessageRequest>(create);
  static DeleteChatMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get messageId => $_getSZ(1);
  @$pb.TagNumber(2)
  set messageId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMessageId() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessageId() => $_clearField(2);
}

/// DeleteChatMessageResponse 删除对话消息响应
class DeleteChatMessageResponse extends $pb.GeneratedMessage {
  factory DeleteChatMessageResponse({
    $core.bool? success,
    $core.String? message,
  }) {
    final result = create();
    if (success != null) result.success = success;
    if (message != null) result.message = message;
    return result;
  }

  DeleteChatMessageResponse._();

  factory DeleteChatMessageResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory DeleteChatMessageResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteChatMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'success')
    ..aOS(2, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatMessageResponse clone() => DeleteChatMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteChatMessageResponse copyWith(void Function(DeleteChatMessageResponse) updates) => super.copyWith((message) => updates(message as DeleteChatMessageResponse)) as DeleteChatMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteChatMessageResponse create() => DeleteChatMessageResponse._();
  @$core.override
  DeleteChatMessageResponse createEmptyInstance() => create();
  static $pb.PbList<DeleteChatMessageResponse> createRepeated() => $pb.PbList<DeleteChatMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static DeleteChatMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteChatMessageResponse>(create);
  static DeleteChatMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get success => $_getBF(0);
  @$pb.TagNumber(1)
  set success($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSuccess() => $_has(0);
  @$pb.TagNumber(1)
  void clearSuccess() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(2)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);
}

/// ChatCompletionTool 工具定义
class ChatCompletionTool extends $pb.GeneratedMessage {
  factory ChatCompletionTool({
    $core.String? type,
    ChatCompletionFunction? function,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (function != null) result.function = function;
    return result;
  }

  ChatCompletionTool._();

  factory ChatCompletionTool.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionTool.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionTool', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'type')
    ..aOM<ChatCompletionFunction>(2, _omitFieldNames ? '' : 'function', subBuilder: ChatCompletionFunction.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionTool clone() => ChatCompletionTool()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionTool copyWith(void Function(ChatCompletionTool) updates) => super.copyWith((message) => updates(message as ChatCompletionTool)) as ChatCompletionTool;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionTool create() => ChatCompletionTool._();
  @$core.override
  ChatCompletionTool createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionTool> createRepeated() => $pb.PbList<ChatCompletionTool>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionTool getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionTool>(create);
  static ChatCompletionTool? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get type => $_getSZ(0);
  @$pb.TagNumber(1)
  set type($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  @$pb.TagNumber(2)
  ChatCompletionFunction get function => $_getN(1);
  @$pb.TagNumber(2)
  set function(ChatCompletionFunction value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasFunction() => $_has(1);
  @$pb.TagNumber(2)
  void clearFunction() => $_clearField(2);
  @$pb.TagNumber(2)
  ChatCompletionFunction ensureFunction() => $_ensure(1);
}

/// ChatCompletionFunction 函数定义
class ChatCompletionFunction extends $pb.GeneratedMessage {
  factory ChatCompletionFunction({
    $core.String? name,
    $core.String? description,
    $core.String? parameters,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (parameters != null) result.parameters = parameters;
    return result;
  }

  ChatCompletionFunction._();

  factory ChatCompletionFunction.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionFunction.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionFunction', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'description')
    ..aOS(3, _omitFieldNames ? '' : 'parameters')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionFunction clone() => ChatCompletionFunction()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionFunction copyWith(void Function(ChatCompletionFunction) updates) => super.copyWith((message) => updates(message as ChatCompletionFunction)) as ChatCompletionFunction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionFunction create() => ChatCompletionFunction._();
  @$core.override
  ChatCompletionFunction createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionFunction> createRepeated() => $pb.PbList<ChatCompletionFunction>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionFunction getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionFunction>(create);
  static ChatCompletionFunction? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(2)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(2)
  void clearDescription() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get parameters => $_getSZ(2);
  @$pb.TagNumber(3)
  set parameters($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasParameters() => $_has(2);
  @$pb.TagNumber(3)
  void clearParameters() => $_clearField(3);
}

enum ToolChoice_Choice {
  mode, 
  named, 
  notSet
}

/// ToolChoice 工具选择策略
class ToolChoice extends $pb.GeneratedMessage {
  factory ToolChoice({
    $core.String? mode,
    ChatCompletionNamedToolChoice? named,
  }) {
    final result = create();
    if (mode != null) result.mode = mode;
    if (named != null) result.named = named;
    return result;
  }

  ToolChoice._();

  factory ToolChoice.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ToolChoice.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, ToolChoice_Choice> _ToolChoice_ChoiceByTag = {
    1 : ToolChoice_Choice.mode,
    2 : ToolChoice_Choice.named,
    0 : ToolChoice_Choice.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ToolChoice', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..oo(0, [1, 2])
    ..aOS(1, _omitFieldNames ? '' : 'mode')
    ..aOM<ChatCompletionNamedToolChoice>(2, _omitFieldNames ? '' : 'named', subBuilder: ChatCompletionNamedToolChoice.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ToolChoice clone() => ToolChoice()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ToolChoice copyWith(void Function(ToolChoice) updates) => super.copyWith((message) => updates(message as ToolChoice)) as ToolChoice;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ToolChoice create() => ToolChoice._();
  @$core.override
  ToolChoice createEmptyInstance() => create();
  static $pb.PbList<ToolChoice> createRepeated() => $pb.PbList<ToolChoice>();
  @$core.pragma('dart2js:noInline')
  static ToolChoice getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ToolChoice>(create);
  static ToolChoice? _defaultInstance;

  ToolChoice_Choice whichChoice() => _ToolChoice_ChoiceByTag[$_whichOneof(0)]!;
  void clearChoice() => $_clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  $core.String get mode => $_getSZ(0);
  @$pb.TagNumber(1)
  set mode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMode() => $_has(0);
  @$pb.TagNumber(1)
  void clearMode() => $_clearField(1);

  @$pb.TagNumber(2)
  ChatCompletionNamedToolChoice get named => $_getN(1);
  @$pb.TagNumber(2)
  set named(ChatCompletionNamedToolChoice value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasNamed() => $_has(1);
  @$pb.TagNumber(2)
  void clearNamed() => $_clearField(2);
  @$pb.TagNumber(2)
  ChatCompletionNamedToolChoice ensureNamed() => $_ensure(1);
}

/// ChatCompletionNamedToolChoice 指定特定工具
class ChatCompletionNamedToolChoice extends $pb.GeneratedMessage {
  factory ChatCompletionNamedToolChoice({
    $core.String? type,
    ChatCompletionFunction? function,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (function != null) result.function = function;
    return result;
  }

  ChatCompletionNamedToolChoice._();

  factory ChatCompletionNamedToolChoice.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionNamedToolChoice.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionNamedToolChoice', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'type')
    ..aOM<ChatCompletionFunction>(2, _omitFieldNames ? '' : 'function', subBuilder: ChatCompletionFunction.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionNamedToolChoice clone() => ChatCompletionNamedToolChoice()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionNamedToolChoice copyWith(void Function(ChatCompletionNamedToolChoice) updates) => super.copyWith((message) => updates(message as ChatCompletionNamedToolChoice)) as ChatCompletionNamedToolChoice;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionNamedToolChoice create() => ChatCompletionNamedToolChoice._();
  @$core.override
  ChatCompletionNamedToolChoice createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionNamedToolChoice> createRepeated() => $pb.PbList<ChatCompletionNamedToolChoice>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionNamedToolChoice getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionNamedToolChoice>(create);
  static ChatCompletionNamedToolChoice? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get type => $_getSZ(0);
  @$pb.TagNumber(1)
  set type($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  @$pb.TagNumber(2)
  ChatCompletionFunction get function => $_getN(1);
  @$pb.TagNumber(2)
  set function(ChatCompletionFunction value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasFunction() => $_has(1);
  @$pb.TagNumber(2)
  void clearFunction() => $_clearField(2);
  @$pb.TagNumber(2)
  ChatCompletionFunction ensureFunction() => $_ensure(1);
}

/// ResponseFormat 响应格式
class ResponseFormat extends $pb.GeneratedMessage {
  factory ResponseFormat({
    $core.String? type,
    $core.String? jsonSchema,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (jsonSchema != null) result.jsonSchema = jsonSchema;
    return result;
  }

  ResponseFormat._();

  factory ResponseFormat.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ResponseFormat.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ResponseFormat', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'type')
    ..aOS(2, _omitFieldNames ? '' : 'jsonSchema')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseFormat clone() => ResponseFormat()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseFormat copyWith(void Function(ResponseFormat) updates) => super.copyWith((message) => updates(message as ResponseFormat)) as ResponseFormat;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseFormat create() => ResponseFormat._();
  @$core.override
  ResponseFormat createEmptyInstance() => create();
  static $pb.PbList<ResponseFormat> createRepeated() => $pb.PbList<ResponseFormat>();
  @$core.pragma('dart2js:noInline')
  static ResponseFormat getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ResponseFormat>(create);
  static ResponseFormat? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get type => $_getSZ(0);
  @$pb.TagNumber(1)
  set type($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get jsonSchema => $_getSZ(1);
  @$pb.TagNumber(2)
  set jsonSchema($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasJsonSchema() => $_has(1);
  @$pb.TagNumber(2)
  void clearJsonSchema() => $_clearField(2);
}

/// ChatCompletionChoice 聊天完成选择
class ChatCompletionChoice extends $pb.GeneratedMessage {
  factory ChatCompletionChoice({
    $core.int? index,
    $1.Message? delta,
    $1.LogProbs? logprobs,
    $core.String? finishReason,
  }) {
    final result = create();
    if (index != null) result.index = index;
    if (delta != null) result.delta = delta;
    if (logprobs != null) result.logprobs = logprobs;
    if (finishReason != null) result.finishReason = finishReason;
    return result;
  }

  ChatCompletionChoice._();

  factory ChatCompletionChoice.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionChoice.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionChoice', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'index', $pb.PbFieldType.O3)
    ..aOM<$1.Message>(2, _omitFieldNames ? '' : 'delta', subBuilder: $1.Message.create)
    ..aOM<$1.LogProbs>(3, _omitFieldNames ? '' : 'logprobs', subBuilder: $1.LogProbs.create)
    ..aOS(4, _omitFieldNames ? '' : 'finishReason')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChoice clone() => ChatCompletionChoice()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChoice copyWith(void Function(ChatCompletionChoice) updates) => super.copyWith((message) => updates(message as ChatCompletionChoice)) as ChatCompletionChoice;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionChoice create() => ChatCompletionChoice._();
  @$core.override
  ChatCompletionChoice createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionChoice> createRepeated() => $pb.PbList<ChatCompletionChoice>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionChoice getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionChoice>(create);
  static ChatCompletionChoice? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get index => $_getIZ(0);
  @$pb.TagNumber(1)
  set index($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIndex() => $_has(0);
  @$pb.TagNumber(1)
  void clearIndex() => $_clearField(1);

  @$pb.TagNumber(2)
  $1.Message get delta => $_getN(1);
  @$pb.TagNumber(2)
  set delta($1.Message value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasDelta() => $_has(1);
  @$pb.TagNumber(2)
  void clearDelta() => $_clearField(2);
  @$pb.TagNumber(2)
  $1.Message ensureDelta() => $_ensure(1);

  @$pb.TagNumber(3)
  $1.LogProbs get logprobs => $_getN(2);
  @$pb.TagNumber(3)
  set logprobs($1.LogProbs value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasLogprobs() => $_has(2);
  @$pb.TagNumber(3)
  void clearLogprobs() => $_clearField(3);
  @$pb.TagNumber(3)
  $1.LogProbs ensureLogprobs() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.String get finishReason => $_getSZ(3);
  @$pb.TagNumber(4)
  set finishReason($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasFinishReason() => $_has(3);
  @$pb.TagNumber(4)
  void clearFinishReason() => $_clearField(4);
}

/// ChatCompletionChunk 流式响应块（用于流式传输）
class ChatCompletionChunk extends $pb.GeneratedMessage {
  factory ChatCompletionChunk({
    $core.String? id,
    $core.String? object,
    $fixnum.Int64? created,
    $core.String? model,
    $core.Iterable<ChatCompletionChunkChoice>? choices,
    $1.TokenUsage? usage,
    $core.String? systemFingerprint,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (object != null) result.object = object;
    if (created != null) result.created = created;
    if (model != null) result.model = model;
    if (choices != null) result.choices.addAll(choices);
    if (usage != null) result.usage = usage;
    if (systemFingerprint != null) result.systemFingerprint = systemFingerprint;
    return result;
  }

  ChatCompletionChunk._();

  factory ChatCompletionChunk.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionChunk.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionChunk', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'object')
    ..aInt64(3, _omitFieldNames ? '' : 'created')
    ..aOS(4, _omitFieldNames ? '' : 'model')
    ..pc<ChatCompletionChunkChoice>(5, _omitFieldNames ? '' : 'choices', $pb.PbFieldType.PM, subBuilder: ChatCompletionChunkChoice.create)
    ..aOM<$1.TokenUsage>(6, _omitFieldNames ? '' : 'usage', subBuilder: $1.TokenUsage.create)
    ..aOS(7, _omitFieldNames ? '' : 'systemFingerprint')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunk clone() => ChatCompletionChunk()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunk copyWith(void Function(ChatCompletionChunk) updates) => super.copyWith((message) => updates(message as ChatCompletionChunk)) as ChatCompletionChunk;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunk create() => ChatCompletionChunk._();
  @$core.override
  ChatCompletionChunk createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionChunk> createRepeated() => $pb.PbList<ChatCompletionChunk>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunk getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionChunk>(create);
  static ChatCompletionChunk? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get object => $_getSZ(1);
  @$pb.TagNumber(2)
  set object($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasObject() => $_has(1);
  @$pb.TagNumber(2)
  void clearObject() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get created => $_getI64(2);
  @$pb.TagNumber(3)
  set created($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCreated() => $_has(2);
  @$pb.TagNumber(3)
  void clearCreated() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get model => $_getSZ(3);
  @$pb.TagNumber(4)
  set model($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasModel() => $_has(3);
  @$pb.TagNumber(4)
  void clearModel() => $_clearField(4);

  @$pb.TagNumber(5)
  $pb.PbList<ChatCompletionChunkChoice> get choices => $_getList(4);

  @$pb.TagNumber(6)
  $1.TokenUsage get usage => $_getN(5);
  @$pb.TagNumber(6)
  set usage($1.TokenUsage value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasUsage() => $_has(5);
  @$pb.TagNumber(6)
  void clearUsage() => $_clearField(6);
  @$pb.TagNumber(6)
  $1.TokenUsage ensureUsage() => $_ensure(5);

  @$pb.TagNumber(7)
  $core.String get systemFingerprint => $_getSZ(6);
  @$pb.TagNumber(7)
  set systemFingerprint($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasSystemFingerprint() => $_has(6);
  @$pb.TagNumber(7)
  void clearSystemFingerprint() => $_clearField(7);
}

/// ChatCompletionChunkChoice 流式响应选择
class ChatCompletionChunkChoice extends $pb.GeneratedMessage {
  factory ChatCompletionChunkChoice({
    $core.int? index,
    ChatCompletionChunkDelta? delta,
    $1.LogProbs? logprobs,
    $core.String? finishReason,
  }) {
    final result = create();
    if (index != null) result.index = index;
    if (delta != null) result.delta = delta;
    if (logprobs != null) result.logprobs = logprobs;
    if (finishReason != null) result.finishReason = finishReason;
    return result;
  }

  ChatCompletionChunkChoice._();

  factory ChatCompletionChunkChoice.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionChunkChoice.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionChunkChoice', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'index', $pb.PbFieldType.O3)
    ..aOM<ChatCompletionChunkDelta>(2, _omitFieldNames ? '' : 'delta', subBuilder: ChatCompletionChunkDelta.create)
    ..aOM<$1.LogProbs>(3, _omitFieldNames ? '' : 'logprobs', subBuilder: $1.LogProbs.create)
    ..aOS(4, _omitFieldNames ? '' : 'finishReason')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunkChoice clone() => ChatCompletionChunkChoice()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunkChoice copyWith(void Function(ChatCompletionChunkChoice) updates) => super.copyWith((message) => updates(message as ChatCompletionChunkChoice)) as ChatCompletionChunkChoice;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunkChoice create() => ChatCompletionChunkChoice._();
  @$core.override
  ChatCompletionChunkChoice createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionChunkChoice> createRepeated() => $pb.PbList<ChatCompletionChunkChoice>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunkChoice getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionChunkChoice>(create);
  static ChatCompletionChunkChoice? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get index => $_getIZ(0);
  @$pb.TagNumber(1)
  set index($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIndex() => $_has(0);
  @$pb.TagNumber(1)
  void clearIndex() => $_clearField(1);

  @$pb.TagNumber(2)
  ChatCompletionChunkDelta get delta => $_getN(1);
  @$pb.TagNumber(2)
  set delta(ChatCompletionChunkDelta value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasDelta() => $_has(1);
  @$pb.TagNumber(2)
  void clearDelta() => $_clearField(2);
  @$pb.TagNumber(2)
  ChatCompletionChunkDelta ensureDelta() => $_ensure(1);

  @$pb.TagNumber(3)
  $1.LogProbs get logprobs => $_getN(2);
  @$pb.TagNumber(3)
  set logprobs($1.LogProbs value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasLogprobs() => $_has(2);
  @$pb.TagNumber(3)
  void clearLogprobs() => $_clearField(3);
  @$pb.TagNumber(3)
  $1.LogProbs ensureLogprobs() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.String get finishReason => $_getSZ(3);
  @$pb.TagNumber(4)
  set finishReason($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasFinishReason() => $_has(3);
  @$pb.TagNumber(4)
  void clearFinishReason() => $_clearField(4);
}

/// ChatCompletionChunkDelta 流式响应增量内容
class ChatCompletionChunkDelta extends $pb.GeneratedMessage {
  factory ChatCompletionChunkDelta({
    $core.String? role,
    $core.String? content,
    $core.Iterable<$1.ToolCall>? toolCalls,
  }) {
    final result = create();
    if (role != null) result.role = role;
    if (content != null) result.content = content;
    if (toolCalls != null) result.toolCalls.addAll(toolCalls);
    return result;
  }

  ChatCompletionChunkDelta._();

  factory ChatCompletionChunkDelta.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatCompletionChunkDelta.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatCompletionChunkDelta', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.server'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'role')
    ..aOS(2, _omitFieldNames ? '' : 'content')
    ..pc<$1.ToolCall>(3, _omitFieldNames ? '' : 'toolCalls', $pb.PbFieldType.PM, subBuilder: $1.ToolCall.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunkDelta clone() => ChatCompletionChunkDelta()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatCompletionChunkDelta copyWith(void Function(ChatCompletionChunkDelta) updates) => super.copyWith((message) => updates(message as ChatCompletionChunkDelta)) as ChatCompletionChunkDelta;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunkDelta create() => ChatCompletionChunkDelta._();
  @$core.override
  ChatCompletionChunkDelta createEmptyInstance() => create();
  static $pb.PbList<ChatCompletionChunkDelta> createRepeated() => $pb.PbList<ChatCompletionChunkDelta>();
  @$core.pragma('dart2js:noInline')
  static ChatCompletionChunkDelta getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatCompletionChunkDelta>(create);
  static ChatCompletionChunkDelta? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get role => $_getSZ(0);
  @$pb.TagNumber(1)
  set role($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRole() => $_has(0);
  @$pb.TagNumber(1)
  void clearRole() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get content => $_getSZ(1);
  @$pb.TagNumber(2)
  set content($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasContent() => $_has(1);
  @$pb.TagNumber(2)
  void clearContent() => $_clearField(2);

  @$pb.TagNumber(3)
  $pb.PbList<$1.ToolCall> get toolCalls => $_getList(2);
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
