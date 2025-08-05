// This is a generated file - do not edit.
//
// Generated from rpc/common.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'common.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'common.pbenum.dart';

/// Empty 空参数占位
class Empty extends $pb.GeneratedMessage {
  factory Empty() => create();

  Empty._();

  factory Empty.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory Empty.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Empty', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Empty clone() => Empty()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Empty copyWith(void Function(Empty) updates) => super.copyWith((message) => updates(message as Empty)) as Empty;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Empty create() => Empty._();
  @$core.override
  Empty createEmptyInstance() => create();
  static $pb.PbList<Empty> createRepeated() => $pb.PbList<Empty>();
  @$core.pragma('dart2js:noInline')
  static Empty getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Empty>(create);
  static Empty? _defaultInstance;
}

/// LlmProvider llm API配置信息
class LlmProvider extends $pb.GeneratedMessage {
  factory LlmProvider({
    $core.String? id,
    $core.String? name,
    $core.String? baseUrl,
    $core.String? apiKey,
    $core.String? alias,
    $core.String? description,
    $core.Iterable<Model>? models,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (name != null) result.name = name;
    if (baseUrl != null) result.baseUrl = baseUrl;
    if (apiKey != null) result.apiKey = apiKey;
    if (alias != null) result.alias = alias;
    if (description != null) result.description = description;
    if (models != null) result.models.addAll(models);
    return result;
  }

  LlmProvider._();

  factory LlmProvider.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory LlmProvider.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LlmProvider', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'name')
    ..aOS(3, _omitFieldNames ? '' : 'baseUrl')
    ..aOS(4, _omitFieldNames ? '' : 'apiKey')
    ..aOS(5, _omitFieldNames ? '' : 'alias')
    ..aOS(6, _omitFieldNames ? '' : 'description')
    ..pc<Model>(7, _omitFieldNames ? '' : 'models', $pb.PbFieldType.PM, subBuilder: Model.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LlmProvider clone() => LlmProvider()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LlmProvider copyWith(void Function(LlmProvider) updates) => super.copyWith((message) => updates(message as LlmProvider)) as LlmProvider;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LlmProvider create() => LlmProvider._();
  @$core.override
  LlmProvider createEmptyInstance() => create();
  static $pb.PbList<LlmProvider> createRepeated() => $pb.PbList<LlmProvider>();
  @$core.pragma('dart2js:noInline')
  static LlmProvider getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LlmProvider>(create);
  static LlmProvider? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(2)
  void clearName() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get baseUrl => $_getSZ(2);
  @$pb.TagNumber(3)
  set baseUrl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasBaseUrl() => $_has(2);
  @$pb.TagNumber(3)
  void clearBaseUrl() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get apiKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set apiKey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasApiKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearApiKey() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get alias => $_getSZ(4);
  @$pb.TagNumber(5)
  set alias($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasAlias() => $_has(4);
  @$pb.TagNumber(5)
  void clearAlias() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get description => $_getSZ(5);
  @$pb.TagNumber(6)
  set description($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasDescription() => $_has(5);
  @$pb.TagNumber(6)
  void clearDescription() => $_clearField(6);

  @$pb.TagNumber(7)
  $pb.PbList<Model> get models => $_getList(6);
}

/// Model llm模型信息
class Model extends $pb.GeneratedMessage {
  factory Model({
    $core.String? id,
    $core.String? object,
    $core.String? ownedBy,
    $core.bool? enabled,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (object != null) result.object = object;
    if (ownedBy != null) result.ownedBy = ownedBy;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  Model._();

  factory Model.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory Model.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Model', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'object')
    ..aOS(3, _omitFieldNames ? '' : 'ownedBy')
    ..aOB(4, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Model clone() => Model()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Model copyWith(void Function(Model) updates) => super.copyWith((message) => updates(message as Model)) as Model;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Model create() => Model._();
  @$core.override
  Model createEmptyInstance() => create();
  static $pb.PbList<Model> createRepeated() => $pb.PbList<Model>();
  @$core.pragma('dart2js:noInline')
  static Model getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Model>(create);
  static Model? _defaultInstance;

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
  $core.String get ownedBy => $_getSZ(2);
  @$pb.TagNumber(3)
  set ownedBy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasOwnedBy() => $_has(2);
  @$pb.TagNumber(3)
  void clearOwnedBy() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get enabled => $_getBF(3);
  @$pb.TagNumber(4)
  set enabled($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasEnabled() => $_has(3);
  @$pb.TagNumber(4)
  void clearEnabled() => $_clearField(4);
}

/// FunctionCall 是消息中的函数调用信息
class FunctionCall extends $pb.GeneratedMessage {
  factory FunctionCall({
    $core.String? name,
    $core.String? arguments,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arguments != null) result.arguments = arguments;
    return result;
  }

  FunctionCall._();

  factory FunctionCall.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory FunctionCall.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'FunctionCall', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'arguments')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FunctionCall clone() => FunctionCall()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FunctionCall copyWith(void Function(FunctionCall) updates) => super.copyWith((message) => updates(message as FunctionCall)) as FunctionCall;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FunctionCall create() => FunctionCall._();
  @$core.override
  FunctionCall createEmptyInstance() => create();
  static $pb.PbList<FunctionCall> createRepeated() => $pb.PbList<FunctionCall>();
  @$core.pragma('dart2js:noInline')
  static FunctionCall getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<FunctionCall>(create);
  static FunctionCall? _defaultInstance;

  /// Name 是要调用的函数名称
  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  /// Arguments 是调用函数所需的参数，以 JSON 格式表示
  @$pb.TagNumber(2)
  $core.String get arguments => $_getSZ(1);
  @$pb.TagNumber(2)
  set arguments($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasArguments() => $_has(1);
  @$pb.TagNumber(2)
  void clearArguments() => $_clearField(2);
}

/// ToolCall 是消息中的工具调用信息
class ToolCall extends $pb.GeneratedMessage {
  factory ToolCall({
    $core.int? index,
    $core.String? id,
    $core.String? type,
    FunctionCall? function,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (index != null) result.index = index;
    if (id != null) result.id = id;
    if (type != null) result.type = type;
    if (function != null) result.function = function;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  ToolCall._();

  factory ToolCall.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ToolCall.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ToolCall', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'index', $pb.PbFieldType.O3)
    ..aOS(2, _omitFieldNames ? '' : 'id')
    ..aOS(3, _omitFieldNames ? '' : 'type')
    ..aOM<FunctionCall>(4, _omitFieldNames ? '' : 'function', subBuilder: FunctionCall.create)
    ..m<$core.String, $core.String>(5, _omitFieldNames ? '' : 'extra', entryClassName: 'ToolCall.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ToolCall clone() => ToolCall()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ToolCall copyWith(void Function(ToolCall) updates) => super.copyWith((message) => updates(message as ToolCall)) as ToolCall;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ToolCall create() => ToolCall._();
  @$core.override
  ToolCall createEmptyInstance() => create();
  static $pb.PbList<ToolCall> createRepeated() => $pb.PbList<ToolCall>();
  @$core.pragma('dart2js:noInline')
  static ToolCall getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ToolCall>(create);
  static ToolCall? _defaultInstance;

  /// Index 在一条消息包含多个工具调用时使用
  @$pb.TagNumber(1)
  $core.int get index => $_getIZ(0);
  @$pb.TagNumber(1)
  set index($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIndex() => $_has(0);
  @$pb.TagNumber(1)
  void clearIndex() => $_clearField(1);

  /// ID 是工具调用的唯一标识
  @$pb.TagNumber(2)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(2)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(2)
  void clearId() => $_clearField(2);

  /// Type 是工具调用的类型
  @$pb.TagNumber(3)
  $core.String get type => $_getSZ(2);
  @$pb.TagNumber(3)
  set type($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(3)
  void clearType() => $_clearField(3);

  /// Function 是具体的函数调用内容
  @$pb.TagNumber(4)
  FunctionCall get function => $_getN(3);
  @$pb.TagNumber(4)
  set function(FunctionCall value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasFunction() => $_has(3);
  @$pb.TagNumber(4)
  void clearFunction() => $_clearField(4);
  @$pb.TagNumber(4)
  FunctionCall ensureFunction() => $_ensure(3);

  /// Extra 用于存储工具调用的额外信息
  @$pb.TagNumber(5)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(4);
}

/// ChatMessageImageURL 表示聊天消息中的图像部分
class ChatMessageImageURL extends $pb.GeneratedMessage {
  factory ChatMessageImageURL({
    $core.String? url,
    $core.String? uri,
    ImageURLDetail? detail,
    $core.String? mimeType,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (url != null) result.url = url;
    if (uri != null) result.uri = uri;
    if (detail != null) result.detail = detail;
    if (mimeType != null) result.mimeType = mimeType;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  ChatMessageImageURL._();

  factory ChatMessageImageURL.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatMessageImageURL.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatMessageImageURL', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'url')
    ..aOS(2, _omitFieldNames ? '' : 'uri')
    ..e<ImageURLDetail>(3, _omitFieldNames ? '' : 'detail', $pb.PbFieldType.OE, defaultOrMaker: ImageURLDetail.IMAGE_URL_DETAIL_UNSPECIFIED, valueOf: ImageURLDetail.valueOf, enumValues: ImageURLDetail.values)
    ..aOS(4, _omitFieldNames ? '' : 'mimeType')
    ..m<$core.String, $core.String>(5, _omitFieldNames ? '' : 'extra', entryClassName: 'ChatMessageImageURL.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageImageURL clone() => ChatMessageImageURL()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageImageURL copyWith(void Function(ChatMessageImageURL) updates) => super.copyWith((message) => updates(message as ChatMessageImageURL)) as ChatMessageImageURL;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatMessageImageURL create() => ChatMessageImageURL._();
  @$core.override
  ChatMessageImageURL createEmptyInstance() => create();
  static $pb.PbList<ChatMessageImageURL> createRepeated() => $pb.PbList<ChatMessageImageURL>();
  @$core.pragma('dart2js:noInline')
  static ChatMessageImageURL getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatMessageImageURL>(create);
  static ChatMessageImageURL? _defaultInstance;

  /// URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL（如 data URL）
  @$pb.TagNumber(1)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(1)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearUrl() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get uri => $_getSZ(1);
  @$pb.TagNumber(2)
  set uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUri() => $_has(1);
  @$pb.TagNumber(2)
  void clearUri() => $_clearField(2);

  /// Detail 是图像 URL 的质量等级
  @$pb.TagNumber(3)
  ImageURLDetail get detail => $_getN(2);
  @$pb.TagNumber(3)
  set detail(ImageURLDetail value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasDetail() => $_has(2);
  @$pb.TagNumber(3)
  void clearDetail() => $_clearField(3);

  /// MIMEType 是图像的 MIME 类型
  @$pb.TagNumber(4)
  $core.String get mimeType => $_getSZ(3);
  @$pb.TagNumber(4)
  set mimeType($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasMimeType() => $_has(3);
  @$pb.TagNumber(4)
  void clearMimeType() => $_clearField(4);

  /// Extra 用于存储图像 URL 的额外信息
  @$pb.TagNumber(5)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(4);
}

/// ChatMessageAudioURL 表示聊天消息中的音频部分
class ChatMessageAudioURL extends $pb.GeneratedMessage {
  factory ChatMessageAudioURL({
    $core.String? url,
    $core.String? uri,
    $core.String? mimeType,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (url != null) result.url = url;
    if (uri != null) result.uri = uri;
    if (mimeType != null) result.mimeType = mimeType;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  ChatMessageAudioURL._();

  factory ChatMessageAudioURL.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatMessageAudioURL.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatMessageAudioURL', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'url')
    ..aOS(2, _omitFieldNames ? '' : 'uri')
    ..aOS(3, _omitFieldNames ? '' : 'mimeType')
    ..m<$core.String, $core.String>(4, _omitFieldNames ? '' : 'extra', entryClassName: 'ChatMessageAudioURL.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageAudioURL clone() => ChatMessageAudioURL()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageAudioURL copyWith(void Function(ChatMessageAudioURL) updates) => super.copyWith((message) => updates(message as ChatMessageAudioURL)) as ChatMessageAudioURL;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatMessageAudioURL create() => ChatMessageAudioURL._();
  @$core.override
  ChatMessageAudioURL createEmptyInstance() => create();
  static $pb.PbList<ChatMessageAudioURL> createRepeated() => $pb.PbList<ChatMessageAudioURL>();
  @$core.pragma('dart2js:noInline')
  static ChatMessageAudioURL getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatMessageAudioURL>(create);
  static ChatMessageAudioURL? _defaultInstance;

  /// URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL
  @$pb.TagNumber(1)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(1)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearUrl() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get uri => $_getSZ(1);
  @$pb.TagNumber(2)
  set uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUri() => $_has(1);
  @$pb.TagNumber(2)
  void clearUri() => $_clearField(2);

  /// MIMEType 是音频的 MIME 类型
  @$pb.TagNumber(3)
  $core.String get mimeType => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimeType($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMimeType() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimeType() => $_clearField(3);

  /// Extra 用于存储音频 URL 的额外信息
  @$pb.TagNumber(4)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(3);
}

/// ChatMessageVideoURL 表示聊天消息中的视频部分
class ChatMessageVideoURL extends $pb.GeneratedMessage {
  factory ChatMessageVideoURL({
    $core.String? url,
    $core.String? uri,
    $core.String? mimeType,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (url != null) result.url = url;
    if (uri != null) result.uri = uri;
    if (mimeType != null) result.mimeType = mimeType;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  ChatMessageVideoURL._();

  factory ChatMessageVideoURL.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatMessageVideoURL.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatMessageVideoURL', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'url')
    ..aOS(2, _omitFieldNames ? '' : 'uri')
    ..aOS(3, _omitFieldNames ? '' : 'mimeType')
    ..m<$core.String, $core.String>(4, _omitFieldNames ? '' : 'extra', entryClassName: 'ChatMessageVideoURL.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageVideoURL clone() => ChatMessageVideoURL()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageVideoURL copyWith(void Function(ChatMessageVideoURL) updates) => super.copyWith((message) => updates(message as ChatMessageVideoURL)) as ChatMessageVideoURL;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatMessageVideoURL create() => ChatMessageVideoURL._();
  @$core.override
  ChatMessageVideoURL createEmptyInstance() => create();
  static $pb.PbList<ChatMessageVideoURL> createRepeated() => $pb.PbList<ChatMessageVideoURL>();
  @$core.pragma('dart2js:noInline')
  static ChatMessageVideoURL getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatMessageVideoURL>(create);
  static ChatMessageVideoURL? _defaultInstance;

  /// URL 可以是传统 URL 或符合 RFC-2397 的特殊 URL
  @$pb.TagNumber(1)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(1)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearUrl() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get uri => $_getSZ(1);
  @$pb.TagNumber(2)
  set uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUri() => $_has(1);
  @$pb.TagNumber(2)
  void clearUri() => $_clearField(2);

  /// MIMEType 是视频的 MIME 类型
  @$pb.TagNumber(3)
  $core.String get mimeType => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimeType($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMimeType() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimeType() => $_clearField(3);

  /// Extra 用于存储视频 URL 的额外信息
  @$pb.TagNumber(4)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(3);
}

/// ChatMessageFileURL 表示聊天消息中的文件部分
class ChatMessageFileURL extends $pb.GeneratedMessage {
  factory ChatMessageFileURL({
    $core.String? url,
    $core.String? uri,
    $core.String? mimeType,
    $core.String? name,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (url != null) result.url = url;
    if (uri != null) result.uri = uri;
    if (mimeType != null) result.mimeType = mimeType;
    if (name != null) result.name = name;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  ChatMessageFileURL._();

  factory ChatMessageFileURL.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatMessageFileURL.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatMessageFileURL', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'url')
    ..aOS(2, _omitFieldNames ? '' : 'uri')
    ..aOS(3, _omitFieldNames ? '' : 'mimeType')
    ..aOS(4, _omitFieldNames ? '' : 'name')
    ..m<$core.String, $core.String>(5, _omitFieldNames ? '' : 'extra', entryClassName: 'ChatMessageFileURL.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageFileURL clone() => ChatMessageFileURL()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessageFileURL copyWith(void Function(ChatMessageFileURL) updates) => super.copyWith((message) => updates(message as ChatMessageFileURL)) as ChatMessageFileURL;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatMessageFileURL create() => ChatMessageFileURL._();
  @$core.override
  ChatMessageFileURL createEmptyInstance() => create();
  static $pb.PbList<ChatMessageFileURL> createRepeated() => $pb.PbList<ChatMessageFileURL>();
  @$core.pragma('dart2js:noInline')
  static ChatMessageFileURL getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatMessageFileURL>(create);
  static ChatMessageFileURL? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(1)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearUrl() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get uri => $_getSZ(1);
  @$pb.TagNumber(2)
  set uri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUri() => $_has(1);
  @$pb.TagNumber(2)
  void clearUri() => $_clearField(2);

  /// MIMEType 是文件的 MIME 类型
  @$pb.TagNumber(3)
  $core.String get mimeType => $_getSZ(2);
  @$pb.TagNumber(3)
  set mimeType($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMimeType() => $_has(2);
  @$pb.TagNumber(3)
  void clearMimeType() => $_clearField(3);

  /// Name 是文件名称
  @$pb.TagNumber(4)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(4)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(4)
  void clearName() => $_clearField(4);

  /// Extra 用于存储文件 URL 的额外信息
  @$pb.TagNumber(5)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(4);
}

/// ChatMessagePart 是聊天消息中的内容片段
class ChatMessagePart extends $pb.GeneratedMessage {
  factory ChatMessagePart({
    ChatMessagePartType? type,
    $core.String? text,
    ChatMessageImageURL? imageUrl,
    ChatMessageAudioURL? audioUrl,
    ChatMessageVideoURL? videoUrl,
    ChatMessageFileURL? fileUrl,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (text != null) result.text = text;
    if (imageUrl != null) result.imageUrl = imageUrl;
    if (audioUrl != null) result.audioUrl = audioUrl;
    if (videoUrl != null) result.videoUrl = videoUrl;
    if (fileUrl != null) result.fileUrl = fileUrl;
    return result;
  }

  ChatMessagePart._();

  factory ChatMessagePart.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatMessagePart.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatMessagePart', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..e<ChatMessagePartType>(1, _omitFieldNames ? '' : 'type', $pb.PbFieldType.OE, defaultOrMaker: ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_UNSPECIFIED, valueOf: ChatMessagePartType.valueOf, enumValues: ChatMessagePartType.values)
    ..aOS(2, _omitFieldNames ? '' : 'text')
    ..aOM<ChatMessageImageURL>(3, _omitFieldNames ? '' : 'imageUrl', subBuilder: ChatMessageImageURL.create)
    ..aOM<ChatMessageAudioURL>(4, _omitFieldNames ? '' : 'audioUrl', subBuilder: ChatMessageAudioURL.create)
    ..aOM<ChatMessageVideoURL>(5, _omitFieldNames ? '' : 'videoUrl', subBuilder: ChatMessageVideoURL.create)
    ..aOM<ChatMessageFileURL>(6, _omitFieldNames ? '' : 'fileUrl', subBuilder: ChatMessageFileURL.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessagePart clone() => ChatMessagePart()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatMessagePart copyWith(void Function(ChatMessagePart) updates) => super.copyWith((message) => updates(message as ChatMessagePart)) as ChatMessagePart;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatMessagePart create() => ChatMessagePart._();
  @$core.override
  ChatMessagePart createEmptyInstance() => create();
  static $pb.PbList<ChatMessagePart> createRepeated() => $pb.PbList<ChatMessagePart>();
  @$core.pragma('dart2js:noInline')
  static ChatMessagePart getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatMessagePart>(create);
  static ChatMessagePart? _defaultInstance;

  /// Type 是片段的类型
  @$pb.TagNumber(1)
  ChatMessagePartType get type => $_getN(0);
  @$pb.TagNumber(1)
  set type(ChatMessagePartType value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  /// Text 是文本内容，当 Type 为 "text" 时使用
  @$pb.TagNumber(2)
  $core.String get text => $_getSZ(1);
  @$pb.TagNumber(2)
  set text($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasText() => $_has(1);
  @$pb.TagNumber(2)
  void clearText() => $_clearField(2);

  /// ImageURL 是图像链接内容，当 Type 为 "image_url" 时使用
  @$pb.TagNumber(3)
  ChatMessageImageURL get imageUrl => $_getN(2);
  @$pb.TagNumber(3)
  set imageUrl(ChatMessageImageURL value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasImageUrl() => $_has(2);
  @$pb.TagNumber(3)
  void clearImageUrl() => $_clearField(3);
  @$pb.TagNumber(3)
  ChatMessageImageURL ensureImageUrl() => $_ensure(2);

  /// AudioURL 是音频链接内容，当 Type 为 "audio_url" 时使用
  @$pb.TagNumber(4)
  ChatMessageAudioURL get audioUrl => $_getN(3);
  @$pb.TagNumber(4)
  set audioUrl(ChatMessageAudioURL value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasAudioUrl() => $_has(3);
  @$pb.TagNumber(4)
  void clearAudioUrl() => $_clearField(4);
  @$pb.TagNumber(4)
  ChatMessageAudioURL ensureAudioUrl() => $_ensure(3);

  /// VideoURL 是视频链接内容，当 Type 为 "video_url" 时使用
  @$pb.TagNumber(5)
  ChatMessageVideoURL get videoUrl => $_getN(4);
  @$pb.TagNumber(5)
  set videoUrl(ChatMessageVideoURL value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasVideoUrl() => $_has(4);
  @$pb.TagNumber(5)
  void clearVideoUrl() => $_clearField(5);
  @$pb.TagNumber(5)
  ChatMessageVideoURL ensureVideoUrl() => $_ensure(4);

  /// FileURL 是文件链接内容，当 Type 为 "file_url" 时使用
  @$pb.TagNumber(6)
  ChatMessageFileURL get fileUrl => $_getN(5);
  @$pb.TagNumber(6)
  set fileUrl(ChatMessageFileURL value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasFileUrl() => $_has(5);
  @$pb.TagNumber(6)
  void clearFileUrl() => $_clearField(6);
  @$pb.TagNumber(6)
  ChatMessageFileURL ensureFileUrl() => $_ensure(5);
}

/// TopLogProb 表示某个 token 的最高对数概率信息
class TopLogProb extends $pb.GeneratedMessage {
  factory TopLogProb({
    $core.String? token,
    $core.double? logProb,
    $core.Iterable<$fixnum.Int64>? bytes,
  }) {
    final result = create();
    if (token != null) result.token = token;
    if (logProb != null) result.logProb = logProb;
    if (bytes != null) result.bytes.addAll(bytes);
    return result;
  }

  TopLogProb._();

  factory TopLogProb.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory TopLogProb.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'TopLogProb', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'token')
    ..a<$core.double>(2, _omitFieldNames ? '' : 'logProb', $pb.PbFieldType.OD)
    ..p<$fixnum.Int64>(3, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.K6)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TopLogProb clone() => TopLogProb()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TopLogProb copyWith(void Function(TopLogProb) updates) => super.copyWith((message) => updates(message as TopLogProb)) as TopLogProb;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TopLogProb create() => TopLogProb._();
  @$core.override
  TopLogProb createEmptyInstance() => create();
  static $pb.PbList<TopLogProb> createRepeated() => $pb.PbList<TopLogProb>();
  @$core.pragma('dart2js:noInline')
  static TopLogProb getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TopLogProb>(create);
  static TopLogProb? _defaultInstance;

  /// Token 表示 token 的文本内容
  @$pb.TagNumber(1)
  $core.String get token => $_getSZ(0);
  @$pb.TagNumber(1)
  set token($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);

  /// LogProb 是该 token 的对数概率
  @$pb.TagNumber(2)
  $core.double get logProb => $_getN(1);
  @$pb.TagNumber(2)
  set logProb($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLogProb() => $_has(1);
  @$pb.TagNumber(2)
  void clearLogProb() => $_clearField(2);

  /// Bytes 是该 token 的 UTF-8 字节表示（整数列表）
  @$pb.TagNumber(3)
  $pb.PbList<$fixnum.Int64> get bytes => $_getList(2);
}

/// LogProb 表示一个 token 的概率信息
class LogProb extends $pb.GeneratedMessage {
  factory LogProb({
    $core.String? token,
    $core.double? logProb,
    $core.Iterable<$fixnum.Int64>? bytes,
    $core.Iterable<TopLogProb>? topLogProbs,
  }) {
    final result = create();
    if (token != null) result.token = token;
    if (logProb != null) result.logProb = logProb;
    if (bytes != null) result.bytes.addAll(bytes);
    if (topLogProbs != null) result.topLogProbs.addAll(topLogProbs);
    return result;
  }

  LogProb._();

  factory LogProb.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory LogProb.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LogProb', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'token')
    ..a<$core.double>(2, _omitFieldNames ? '' : 'logProb', $pb.PbFieldType.OD)
    ..p<$fixnum.Int64>(3, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.K6)
    ..pc<TopLogProb>(4, _omitFieldNames ? '' : 'topLogProbs', $pb.PbFieldType.PM, subBuilder: TopLogProb.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogProb clone() => LogProb()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogProb copyWith(void Function(LogProb) updates) => super.copyWith((message) => updates(message as LogProb)) as LogProb;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LogProb create() => LogProb._();
  @$core.override
  LogProb createEmptyInstance() => create();
  static $pb.PbList<LogProb> createRepeated() => $pb.PbList<LogProb>();
  @$core.pragma('dart2js:noInline')
  static LogProb getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LogProb>(create);
  static LogProb? _defaultInstance;

  /// Token 表示 token 的文本内容
  @$pb.TagNumber(1)
  $core.String get token => $_getSZ(0);
  @$pb.TagNumber(1)
  set token($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasToken() => $_has(0);
  @$pb.TagNumber(1)
  void clearToken() => $_clearField(1);

  /// LogProb 是该 token 的对数概率
  @$pb.TagNumber(2)
  $core.double get logProb => $_getN(1);
  @$pb.TagNumber(2)
  set logProb($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLogProb() => $_has(1);
  @$pb.TagNumber(2)
  void clearLogProb() => $_clearField(2);

  /// Bytes 是该 token 的 UTF-8 字节表示（整数列表）
  @$pb.TagNumber(3)
  $pb.PbList<$fixnum.Int64> get bytes => $_getList(2);

  /// TopLogProbs 是最可能的若干 token 及其对数概率列表
  @$pb.TagNumber(4)
  $pb.PbList<TopLogProb> get topLogProbs => $_getList(3);
}

/// LogProbs 是包含 token 概率信息的顶层结构
class LogProbs extends $pb.GeneratedMessage {
  factory LogProbs({
    $core.Iterable<LogProb>? content,
  }) {
    final result = create();
    if (content != null) result.content.addAll(content);
    return result;
  }

  LogProbs._();

  factory LogProbs.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory LogProbs.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LogProbs', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..pc<LogProb>(1, _omitFieldNames ? '' : 'content', $pb.PbFieldType.PM, subBuilder: LogProb.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogProbs clone() => LogProbs()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogProbs copyWith(void Function(LogProbs) updates) => super.copyWith((message) => updates(message as LogProbs)) as LogProbs;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LogProbs create() => LogProbs._();
  @$core.override
  LogProbs createEmptyInstance() => create();
  static $pb.PbList<LogProbs> createRepeated() => $pb.PbList<LogProbs>();
  @$core.pragma('dart2js:noInline')
  static LogProbs getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LogProbs>(create);
  static LogProbs? _defaultInstance;

  /// Content 是包含对数概率信息的消息内容 token 列表
  @$pb.TagNumber(1)
  $pb.PbList<LogProb> get content => $_getList(0);
}

/// TokenUsage 表示聊天模型请求的 token 使用情况
class TokenUsage extends $pb.GeneratedMessage {
  factory TokenUsage({
    $core.int? promptTokens,
    $core.int? completionTokens,
    $core.int? totalTokens,
  }) {
    final result = create();
    if (promptTokens != null) result.promptTokens = promptTokens;
    if (completionTokens != null) result.completionTokens = completionTokens;
    if (totalTokens != null) result.totalTokens = totalTokens;
    return result;
  }

  TokenUsage._();

  factory TokenUsage.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory TokenUsage.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'TokenUsage', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'promptTokens', $pb.PbFieldType.O3)
    ..a<$core.int>(2, _omitFieldNames ? '' : 'completionTokens', $pb.PbFieldType.O3)
    ..a<$core.int>(3, _omitFieldNames ? '' : 'totalTokens', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TokenUsage clone() => TokenUsage()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TokenUsage copyWith(void Function(TokenUsage) updates) => super.copyWith((message) => updates(message as TokenUsage)) as TokenUsage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TokenUsage create() => TokenUsage._();
  @$core.override
  TokenUsage createEmptyInstance() => create();
  static $pb.PbList<TokenUsage> createRepeated() => $pb.PbList<TokenUsage>();
  @$core.pragma('dart2js:noInline')
  static TokenUsage getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TokenUsage>(create);
  static TokenUsage? _defaultInstance;

  /// PromptTokens 是提示词中的 token 数量
  @$pb.TagNumber(1)
  $core.int get promptTokens => $_getIZ(0);
  @$pb.TagNumber(1)
  set promptTokens($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPromptTokens() => $_has(0);
  @$pb.TagNumber(1)
  void clearPromptTokens() => $_clearField(1);

  /// CompletionTokens 是生成内容中的 token 数量
  @$pb.TagNumber(2)
  $core.int get completionTokens => $_getIZ(1);
  @$pb.TagNumber(2)
  set completionTokens($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCompletionTokens() => $_has(1);
  @$pb.TagNumber(2)
  void clearCompletionTokens() => $_clearField(2);

  /// TotalTokens 是请求中总的 token 数量
  @$pb.TagNumber(3)
  $core.int get totalTokens => $_getIZ(2);
  @$pb.TagNumber(3)
  set totalTokens($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTotalTokens() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalTokens() => $_clearField(3);
}

/// ResponseMeta 收集聊天响应的元信息
class ResponseMeta extends $pb.GeneratedMessage {
  factory ResponseMeta({
    $core.String? finishReason,
    TokenUsage? usage,
    LogProbs? logProbs,
  }) {
    final result = create();
    if (finishReason != null) result.finishReason = finishReason;
    if (usage != null) result.usage = usage;
    if (logProbs != null) result.logProbs = logProbs;
    return result;
  }

  ResponseMeta._();

  factory ResponseMeta.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ResponseMeta.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ResponseMeta', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'finishReason')
    ..aOM<TokenUsage>(2, _omitFieldNames ? '' : 'usage', subBuilder: TokenUsage.create)
    ..aOM<LogProbs>(3, _omitFieldNames ? '' : 'logProbs', subBuilder: LogProbs.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseMeta clone() => ResponseMeta()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseMeta copyWith(void Function(ResponseMeta) updates) => super.copyWith((message) => updates(message as ResponseMeta)) as ResponseMeta;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseMeta create() => ResponseMeta._();
  @$core.override
  ResponseMeta createEmptyInstance() => create();
  static $pb.PbList<ResponseMeta> createRepeated() => $pb.PbList<ResponseMeta>();
  @$core.pragma('dart2js:noInline')
  static ResponseMeta getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ResponseMeta>(create);
  static ResponseMeta? _defaultInstance;

  /// FinishReason 是聊天响应结束的原因
  @$pb.TagNumber(1)
  $core.String get finishReason => $_getSZ(0);
  @$pb.TagNumber(1)
  set finishReason($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasFinishReason() => $_has(0);
  @$pb.TagNumber(1)
  void clearFinishReason() => $_clearField(1);

  /// Usage 是聊天响应的 token 使用情况
  @$pb.TagNumber(2)
  TokenUsage get usage => $_getN(1);
  @$pb.TagNumber(2)
  set usage(TokenUsage value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasUsage() => $_has(1);
  @$pb.TagNumber(2)
  void clearUsage() => $_clearField(2);
  @$pb.TagNumber(2)
  TokenUsage ensureUsage() => $_ensure(1);

  /// LogProbs 是对数概率信息
  @$pb.TagNumber(3)
  LogProbs get logProbs => $_getN(2);
  @$pb.TagNumber(3)
  set logProbs(LogProbs value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasLogProbs() => $_has(2);
  @$pb.TagNumber(3)
  void clearLogProbs() => $_clearField(3);
  @$pb.TagNumber(3)
  LogProbs ensureLogProbs() => $_ensure(2);
}

/// Message 表示一条聊天消息
class Message extends $pb.GeneratedMessage {
  factory Message({
    RoleType? role,
    $core.String? content,
    $core.Iterable<ChatMessagePart>? multiContent,
    $core.String? name,
    $core.Iterable<ToolCall>? toolCalls,
    $core.String? toolCallId,
    $core.String? toolName,
    ResponseMeta? responseMeta,
    $core.String? reasoningContent,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
  }) {
    final result = create();
    if (role != null) result.role = role;
    if (content != null) result.content = content;
    if (multiContent != null) result.multiContent.addAll(multiContent);
    if (name != null) result.name = name;
    if (toolCalls != null) result.toolCalls.addAll(toolCalls);
    if (toolCallId != null) result.toolCallId = toolCallId;
    if (toolName != null) result.toolName = toolName;
    if (responseMeta != null) result.responseMeta = responseMeta;
    if (reasoningContent != null) result.reasoningContent = reasoningContent;
    if (extra != null) result.extra.addEntries(extra);
    return result;
  }

  Message._();

  factory Message.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory Message.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Message', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..e<RoleType>(1, _omitFieldNames ? '' : 'role', $pb.PbFieldType.OE, defaultOrMaker: RoleType.ROLE_TYPE_UNSPECIFIED, valueOf: RoleType.valueOf, enumValues: RoleType.values)
    ..aOS(2, _omitFieldNames ? '' : 'content')
    ..pc<ChatMessagePart>(3, _omitFieldNames ? '' : 'multiContent', $pb.PbFieldType.PM, subBuilder: ChatMessagePart.create)
    ..aOS(4, _omitFieldNames ? '' : 'name')
    ..pc<ToolCall>(5, _omitFieldNames ? '' : 'toolCalls', $pb.PbFieldType.PM, subBuilder: ToolCall.create)
    ..aOS(6, _omitFieldNames ? '' : 'toolCallId')
    ..aOS(7, _omitFieldNames ? '' : 'toolName')
    ..aOM<ResponseMeta>(8, _omitFieldNames ? '' : 'responseMeta', subBuilder: ResponseMeta.create)
    ..aOS(9, _omitFieldNames ? '' : 'reasoningContent')
    ..m<$core.String, $core.String>(10, _omitFieldNames ? '' : 'extra', entryClassName: 'Message.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Message clone() => Message()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Message copyWith(void Function(Message) updates) => super.copyWith((message) => updates(message as Message)) as Message;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Message create() => Message._();
  @$core.override
  Message createEmptyInstance() => create();
  static $pb.PbList<Message> createRepeated() => $pb.PbList<Message>();
  @$core.pragma('dart2js:noInline')
  static Message getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Message>(create);
  static Message? _defaultInstance;

  @$pb.TagNumber(1)
  RoleType get role => $_getN(0);
  @$pb.TagNumber(1)
  set role(RoleType value) => $_setField(1, value);
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

  /// 如果 multi_content 非空，则优先使用它而不是 content
  @$pb.TagNumber(3)
  $pb.PbList<ChatMessagePart> get multiContent => $_getList(2);

  @$pb.TagNumber(4)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(4)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(4)
  void clearName() => $_clearField(4);

  /// 仅用于 Assistant 消息
  @$pb.TagNumber(5)
  $pb.PbList<ToolCall> get toolCalls => $_getList(4);

  /// 仅用于 Tool 消息
  @$pb.TagNumber(6)
  $core.String get toolCallId => $_getSZ(5);
  @$pb.TagNumber(6)
  set toolCallId($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasToolCallId() => $_has(5);
  @$pb.TagNumber(6)
  void clearToolCallId() => $_clearField(6);

  /// 仅用于 Tool 消息
  @$pb.TagNumber(7)
  $core.String get toolName => $_getSZ(6);
  @$pb.TagNumber(7)
  set toolName($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasToolName() => $_has(6);
  @$pb.TagNumber(7)
  void clearToolName() => $_clearField(7);

  @$pb.TagNumber(8)
  ResponseMeta get responseMeta => $_getN(7);
  @$pb.TagNumber(8)
  set responseMeta(ResponseMeta value) => $_setField(8, value);
  @$pb.TagNumber(8)
  $core.bool hasResponseMeta() => $_has(7);
  @$pb.TagNumber(8)
  void clearResponseMeta() => $_clearField(8);
  @$pb.TagNumber(8)
  ResponseMeta ensureResponseMeta() => $_ensure(7);

  /// ReasoningContent 是模型的推理思考过程
  @$pb.TagNumber(9)
  $core.String get reasoningContent => $_getSZ(8);
  @$pb.TagNumber(9)
  set reasoningContent($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasReasoningContent() => $_has(8);
  @$pb.TagNumber(9)
  void clearReasoningContent() => $_clearField(9);

  /// 针对特定模型实现的自定义扩展信息
  @$pb.TagNumber(10)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(9);
}

/// MessageList 表示一组消息的列表
class MessageList extends $pb.GeneratedMessage {
  factory MessageList({
    $core.Iterable<Message>? messages,
  }) {
    final result = create();
    if (messages != null) result.messages.addAll(messages);
    return result;
  }

  MessageList._();

  factory MessageList.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MessageList.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MessageList', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..pc<Message>(1, _omitFieldNames ? '' : 'messages', $pb.PbFieldType.PM, subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageList clone() => MessageList()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageList copyWith(void Function(MessageList) updates) => super.copyWith((message) => updates(message as MessageList)) as MessageList;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MessageList create() => MessageList._();
  @$core.override
  MessageList createEmptyInstance() => create();
  static $pb.PbList<MessageList> createRepeated() => $pb.PbList<MessageList>();
  @$core.pragma('dart2js:noInline')
  static MessageList getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MessageList>(create);
  static MessageList? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Message> get messages => $_getList(0);
}

class CreateMessageRequest extends $pb.GeneratedMessage {
  factory CreateMessageRequest({
    Message? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CreateMessageRequest._();

  factory CreateMessageRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory CreateMessageRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOM<Message>(1, _omitFieldNames ? '' : 'message', subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateMessageRequest clone() => CreateMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateMessageRequest copyWith(void Function(CreateMessageRequest) updates) => super.copyWith((message) => updates(message as CreateMessageRequest)) as CreateMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateMessageRequest create() => CreateMessageRequest._();
  @$core.override
  CreateMessageRequest createEmptyInstance() => create();
  static $pb.PbList<CreateMessageRequest> createRepeated() => $pb.PbList<CreateMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateMessageRequest>(create);
  static CreateMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  Message get message => $_getN(0);
  @$pb.TagNumber(1)
  set message(Message value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => $_clearField(1);
  @$pb.TagNumber(1)
  Message ensureMessage() => $_ensure(0);
}

class CreateMessageResponse extends $pb.GeneratedMessage {
  factory CreateMessageResponse({
    $core.String? messageId,
    Message? message,
  }) {
    final result = create();
    if (messageId != null) result.messageId = messageId;
    if (message != null) result.message = message;
    return result;
  }

  CreateMessageResponse._();

  factory CreateMessageResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory CreateMessageResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'messageId')
    ..aOM<Message>(2, _omitFieldNames ? '' : 'message', subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateMessageResponse clone() => CreateMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateMessageResponse copyWith(void Function(CreateMessageResponse) updates) => super.copyWith((message) => updates(message as CreateMessageResponse)) as CreateMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateMessageResponse create() => CreateMessageResponse._();
  @$core.override
  CreateMessageResponse createEmptyInstance() => create();
  static $pb.PbList<CreateMessageResponse> createRepeated() => $pb.PbList<CreateMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateMessageResponse>(create);
  static CreateMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get messageId => $_getSZ(0);
  @$pb.TagNumber(1)
  set messageId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMessageId() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessageId() => $_clearField(1);

  @$pb.TagNumber(2)
  Message get message => $_getN(1);
  @$pb.TagNumber(2)
  set message(Message value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);
  @$pb.TagNumber(2)
  Message ensureMessage() => $_ensure(1);
}

class GetMessageRequest extends $pb.GeneratedMessage {
  factory GetMessageRequest({
    $core.String? messageId,
  }) {
    final result = create();
    if (messageId != null) result.messageId = messageId;
    return result;
  }

  GetMessageRequest._();

  factory GetMessageRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory GetMessageRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'messageId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMessageRequest clone() => GetMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMessageRequest copyWith(void Function(GetMessageRequest) updates) => super.copyWith((message) => updates(message as GetMessageRequest)) as GetMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMessageRequest create() => GetMessageRequest._();
  @$core.override
  GetMessageRequest createEmptyInstance() => create();
  static $pb.PbList<GetMessageRequest> createRepeated() => $pb.PbList<GetMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static GetMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetMessageRequest>(create);
  static GetMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get messageId => $_getSZ(0);
  @$pb.TagNumber(1)
  set messageId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMessageId() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessageId() => $_clearField(1);
}

class GetMessageResponse extends $pb.GeneratedMessage {
  factory GetMessageResponse({
    Message? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  GetMessageResponse._();

  factory GetMessageResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory GetMessageResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOM<Message>(1, _omitFieldNames ? '' : 'message', subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMessageResponse clone() => GetMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMessageResponse copyWith(void Function(GetMessageResponse) updates) => super.copyWith((message) => updates(message as GetMessageResponse)) as GetMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMessageResponse create() => GetMessageResponse._();
  @$core.override
  GetMessageResponse createEmptyInstance() => create();
  static $pb.PbList<GetMessageResponse> createRepeated() => $pb.PbList<GetMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static GetMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetMessageResponse>(create);
  static GetMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Message get message => $_getN(0);
  @$pb.TagNumber(1)
  set message(Message value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => $_clearField(1);
  @$pb.TagNumber(1)
  Message ensureMessage() => $_ensure(0);
}

class UpdateMessageRequest extends $pb.GeneratedMessage {
  factory UpdateMessageRequest({
    $core.String? messageId,
    Message? message,
  }) {
    final result = create();
    if (messageId != null) result.messageId = messageId;
    if (message != null) result.message = message;
    return result;
  }

  UpdateMessageRequest._();

  factory UpdateMessageRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory UpdateMessageRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'messageId')
    ..aOM<Message>(2, _omitFieldNames ? '' : 'message', subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMessageRequest clone() => UpdateMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMessageRequest copyWith(void Function(UpdateMessageRequest) updates) => super.copyWith((message) => updates(message as UpdateMessageRequest)) as UpdateMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMessageRequest create() => UpdateMessageRequest._();
  @$core.override
  UpdateMessageRequest createEmptyInstance() => create();
  static $pb.PbList<UpdateMessageRequest> createRepeated() => $pb.PbList<UpdateMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static UpdateMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateMessageRequest>(create);
  static UpdateMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get messageId => $_getSZ(0);
  @$pb.TagNumber(1)
  set messageId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMessageId() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessageId() => $_clearField(1);

  @$pb.TagNumber(2)
  Message get message => $_getN(1);
  @$pb.TagNumber(2)
  set message(Message value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);
  @$pb.TagNumber(2)
  Message ensureMessage() => $_ensure(1);
}

class UpdateMessageResponse extends $pb.GeneratedMessage {
  factory UpdateMessageResponse({
    Message? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UpdateMessageResponse._();

  factory UpdateMessageResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory UpdateMessageResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UpdateMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOM<Message>(1, _omitFieldNames ? '' : 'message', subBuilder: Message.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMessageResponse clone() => UpdateMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMessageResponse copyWith(void Function(UpdateMessageResponse) updates) => super.copyWith((message) => updates(message as UpdateMessageResponse)) as UpdateMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMessageResponse create() => UpdateMessageResponse._();
  @$core.override
  UpdateMessageResponse createEmptyInstance() => create();
  static $pb.PbList<UpdateMessageResponse> createRepeated() => $pb.PbList<UpdateMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static UpdateMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UpdateMessageResponse>(create);
  static UpdateMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  Message get message => $_getN(0);
  @$pb.TagNumber(1)
  set message(Message value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => $_clearField(1);
  @$pb.TagNumber(1)
  Message ensureMessage() => $_ensure(0);
}

class DeleteMessageRequest extends $pb.GeneratedMessage {
  factory DeleteMessageRequest({
    $core.String? messageId,
  }) {
    final result = create();
    if (messageId != null) result.messageId = messageId;
    return result;
  }

  DeleteMessageRequest._();

  factory DeleteMessageRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory DeleteMessageRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'messageId')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageRequest clone() => DeleteMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageRequest copyWith(void Function(DeleteMessageRequest) updates) => super.copyWith((message) => updates(message as DeleteMessageRequest)) as DeleteMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageRequest create() => DeleteMessageRequest._();
  @$core.override
  DeleteMessageRequest createEmptyInstance() => create();
  static $pb.PbList<DeleteMessageRequest> createRepeated() => $pb.PbList<DeleteMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteMessageRequest>(create);
  static DeleteMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get messageId => $_getSZ(0);
  @$pb.TagNumber(1)
  set messageId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMessageId() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessageId() => $_clearField(1);
}

class DeleteMessageResponse extends $pb.GeneratedMessage {
  factory DeleteMessageResponse({
    $core.bool? success,
  }) {
    final result = create();
    if (success != null) result.success = success;
    return result;
  }

  DeleteMessageResponse._();

  factory DeleteMessageResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory DeleteMessageResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'DeleteMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'success')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageResponse clone() => DeleteMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageResponse copyWith(void Function(DeleteMessageResponse) updates) => super.copyWith((message) => updates(message as DeleteMessageResponse)) as DeleteMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageResponse create() => DeleteMessageResponse._();
  @$core.override
  DeleteMessageResponse createEmptyInstance() => create();
  static $pb.PbList<DeleteMessageResponse> createRepeated() => $pb.PbList<DeleteMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DeleteMessageResponse>(create);
  static DeleteMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get success => $_getBF(0);
  @$pb.TagNumber(1)
  set success($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSuccess() => $_has(0);
  @$pb.TagNumber(1)
  void clearSuccess() => $_clearField(1);
}

class ListMessagesRequest extends $pb.GeneratedMessage {
  factory ListMessagesRequest({
    $core.int? pageSize,
    $core.String? pageToken,
    $core.String? filter,
  }) {
    final result = create();
    if (pageSize != null) result.pageSize = pageSize;
    if (pageToken != null) result.pageToken = pageToken;
    if (filter != null) result.filter = filter;
    return result;
  }

  ListMessagesRequest._();

  factory ListMessagesRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ListMessagesRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListMessagesRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'pageSize', $pb.PbFieldType.O3)
    ..aOS(2, _omitFieldNames ? '' : 'pageToken')
    ..aOS(3, _omitFieldNames ? '' : 'filter')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessagesRequest clone() => ListMessagesRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessagesRequest copyWith(void Function(ListMessagesRequest) updates) => super.copyWith((message) => updates(message as ListMessagesRequest)) as ListMessagesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMessagesRequest create() => ListMessagesRequest._();
  @$core.override
  ListMessagesRequest createEmptyInstance() => create();
  static $pb.PbList<ListMessagesRequest> createRepeated() => $pb.PbList<ListMessagesRequest>();
  @$core.pragma('dart2js:noInline')
  static ListMessagesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListMessagesRequest>(create);
  static ListMessagesRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get pageSize => $_getIZ(0);
  @$pb.TagNumber(1)
  set pageSize($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPageSize() => $_has(0);
  @$pb.TagNumber(1)
  void clearPageSize() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get pageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set pageToken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearPageToken() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get filter => $_getSZ(2);
  @$pb.TagNumber(3)
  set filter($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasFilter() => $_has(2);
  @$pb.TagNumber(3)
  void clearFilter() => $_clearField(3);
}

class ListMessagesResponse extends $pb.GeneratedMessage {
  factory ListMessagesResponse({
    $core.Iterable<Message>? messages,
    $core.String? nextPageToken,
    $core.int? totalCount,
  }) {
    final result = create();
    if (messages != null) result.messages.addAll(messages);
    if (nextPageToken != null) result.nextPageToken = nextPageToken;
    if (totalCount != null) result.totalCount = totalCount;
    return result;
  }

  ListMessagesResponse._();

  factory ListMessagesResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ListMessagesResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListMessagesResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..pc<Message>(1, _omitFieldNames ? '' : 'messages', $pb.PbFieldType.PM, subBuilder: Message.create)
    ..aOS(2, _omitFieldNames ? '' : 'nextPageToken')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'totalCount', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessagesResponse clone() => ListMessagesResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessagesResponse copyWith(void Function(ListMessagesResponse) updates) => super.copyWith((message) => updates(message as ListMessagesResponse)) as ListMessagesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMessagesResponse create() => ListMessagesResponse._();
  @$core.override
  ListMessagesResponse createEmptyInstance() => create();
  static $pb.PbList<ListMessagesResponse> createRepeated() => $pb.PbList<ListMessagesResponse>();
  @$core.pragma('dart2js:noInline')
  static ListMessagesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListMessagesResponse>(create);
  static ListMessagesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<Message> get messages => $_getList(0);

  @$pb.TagNumber(2)
  $core.String get nextPageToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set nextPageToken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasNextPageToken() => $_has(1);
  @$pb.TagNumber(2)
  void clearNextPageToken() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.int get totalCount => $_getIZ(2);
  @$pb.TagNumber(3)
  set totalCount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTotalCount() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalCount() => $_clearField(3);
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
