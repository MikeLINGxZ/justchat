// This is a generated file - do not edit.
//
// Generated from rpc/common/common.proto.

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

/// Message 表示一条聊天消息，使用 oneof 分离不同角色的消息类型
class Message extends $pb.GeneratedMessage {
  factory Message({
    $core.String? content,
    $core.Iterable<ChatMessagePart>? multiContent,
    $core.String? name,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? extra,
    $core.String? role,
    $core.String? systemContent,
    $core.Iterable<ToolCall>? toolCalls,
    ResponseMeta? responseMeta,
    $core.String? reasoningContent,
    $core.String? toolCallId,
    $core.String? toolName,
    $core.String? chatUuid,
  }) {
    final result = create();
    if (content != null) result.content = content;
    if (multiContent != null) result.multiContent.addAll(multiContent);
    if (name != null) result.name = name;
    if (extra != null) result.extra.addEntries(extra);
    if (role != null) result.role = role;
    if (systemContent != null) result.systemContent = systemContent;
    if (toolCalls != null) result.toolCalls.addAll(toolCalls);
    if (responseMeta != null) result.responseMeta = responseMeta;
    if (reasoningContent != null) result.reasoningContent = reasoningContent;
    if (toolCallId != null) result.toolCallId = toolCallId;
    if (toolName != null) result.toolName = toolName;
    if (chatUuid != null) result.chatUuid = chatUuid;
    return result;
  }

  Message._();

  factory Message.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory Message.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Message', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'content')
    ..pc<ChatMessagePart>(2, _omitFieldNames ? '' : 'multiContent', $pb.PbFieldType.PM, subBuilder: ChatMessagePart.create)
    ..aOS(3, _omitFieldNames ? '' : 'name')
    ..m<$core.String, $core.String>(4, _omitFieldNames ? '' : 'extra', entryClassName: 'Message.ExtraEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..aOS(5, _omitFieldNames ? '' : 'role')
    ..aOS(10, _omitFieldNames ? '' : 'systemContent')
    ..pc<ToolCall>(11, _omitFieldNames ? '' : 'toolCalls', $pb.PbFieldType.PM, subBuilder: ToolCall.create)
    ..aOM<ResponseMeta>(12, _omitFieldNames ? '' : 'responseMeta', subBuilder: ResponseMeta.create)
    ..aOS(13, _omitFieldNames ? '' : 'reasoningContent')
    ..aOS(14, _omitFieldNames ? '' : 'toolCallId')
    ..aOS(15, _omitFieldNames ? '' : 'toolName')
    ..aOS(16, _omitFieldNames ? '' : 'chatUuid')
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

  /// 通用字段（所有消息类型都可能包含）
  @$pb.TagNumber(1)
  $core.String get content => $_getSZ(0);
  @$pb.TagNumber(1)
  set content($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasContent() => $_has(0);
  @$pb.TagNumber(1)
  void clearContent() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<ChatMessagePart> get multiContent => $_getList(1);

  @$pb.TagNumber(3)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(3)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(3)
  void clearName() => $_clearField(3);

  @$pb.TagNumber(4)
  $pb.PbMap<$core.String, $core.String> get extra => $_getMap(3);

  @$pb.TagNumber(5)
  $core.String get role => $_getSZ(4);
  @$pb.TagNumber(5)
  set role($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasRole() => $_has(4);
  @$pb.TagNumber(5)
  void clearRole() => $_clearField(5);

  /// 系统消息相关字段（仅在 role == SYSTEM 时使用）
  @$pb.TagNumber(10)
  $core.String get systemContent => $_getSZ(5);
  @$pb.TagNumber(10)
  set systemContent($core.String value) => $_setString(5, value);
  @$pb.TagNumber(10)
  $core.bool hasSystemContent() => $_has(5);
  @$pb.TagNumber(10)
  void clearSystemContent() => $_clearField(10);

  /// 助手消息相关字段（仅在 role == ASSISTANT 时使用）
  @$pb.TagNumber(11)
  $pb.PbList<ToolCall> get toolCalls => $_getList(6);

  @$pb.TagNumber(12)
  ResponseMeta get responseMeta => $_getN(7);
  @$pb.TagNumber(12)
  set responseMeta(ResponseMeta value) => $_setField(12, value);
  @$pb.TagNumber(12)
  $core.bool hasResponseMeta() => $_has(7);
  @$pb.TagNumber(12)
  void clearResponseMeta() => $_clearField(12);
  @$pb.TagNumber(12)
  ResponseMeta ensureResponseMeta() => $_ensure(7);

  @$pb.TagNumber(13)
  $core.String get reasoningContent => $_getSZ(8);
  @$pb.TagNumber(13)
  set reasoningContent($core.String value) => $_setString(8, value);
  @$pb.TagNumber(13)
  $core.bool hasReasoningContent() => $_has(8);
  @$pb.TagNumber(13)
  void clearReasoningContent() => $_clearField(13);

  /// 工具消息相关字段（仅在 role == TOOL 时使用）
  @$pb.TagNumber(14)
  $core.String get toolCallId => $_getSZ(9);
  @$pb.TagNumber(14)
  set toolCallId($core.String value) => $_setString(9, value);
  @$pb.TagNumber(14)
  $core.bool hasToolCallId() => $_has(9);
  @$pb.TagNumber(14)
  void clearToolCallId() => $_clearField(14);

  @$pb.TagNumber(15)
  $core.String get toolName => $_getSZ(10);
  @$pb.TagNumber(15)
  set toolName($core.String value) => $_setString(10, value);
  @$pb.TagNumber(15)
  $core.bool hasToolName() => $_has(10);
  @$pb.TagNumber(15)
  void clearToolName() => $_clearField(15);

  /// 以下为业务字段
  @$pb.TagNumber(16)
  $core.String get chatUuid => $_getSZ(11);
  @$pb.TagNumber(16)
  set chatUuid($core.String value) => $_setString(11, value);
  @$pb.TagNumber(16)
  $core.bool hasChatUuid() => $_has(11);
  @$pb.TagNumber(16)
  void clearChatUuid() => $_clearField(16);
}

/// ChatInfo 对话信息
class ChatInfo extends $pb.GeneratedMessage {
  factory ChatInfo({
    $core.String? chatUuid,
    $core.String? title,
    $fixnum.Int64? modelId,
    $fixnum.Int64? createdAt,
    $fixnum.Int64? updatedAt,
    $core.int? messageCount,
    $core.String? lastMessagePreview,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? metadata,
  }) {
    final result = create();
    if (chatUuid != null) result.chatUuid = chatUuid;
    if (title != null) result.title = title;
    if (modelId != null) result.modelId = modelId;
    if (createdAt != null) result.createdAt = createdAt;
    if (updatedAt != null) result.updatedAt = updatedAt;
    if (messageCount != null) result.messageCount = messageCount;
    if (lastMessagePreview != null) result.lastMessagePreview = lastMessagePreview;
    if (metadata != null) result.metadata.addEntries(metadata);
    return result;
  }

  ChatInfo._();

  factory ChatInfo.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ChatInfo.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ChatInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'lemon_tea.common'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'chatUuid')
    ..aOS(2, _omitFieldNames ? '' : 'title')
    ..aInt64(3, _omitFieldNames ? '' : 'modelId')
    ..aInt64(4, _omitFieldNames ? '' : 'createdAt')
    ..aInt64(5, _omitFieldNames ? '' : 'updatedAt')
    ..a<$core.int>(6, _omitFieldNames ? '' : 'messageCount', $pb.PbFieldType.O3)
    ..aOS(7, _omitFieldNames ? '' : 'lastMessagePreview')
    ..m<$core.String, $core.String>(8, _omitFieldNames ? '' : 'metadata', entryClassName: 'ChatInfo.MetadataEntry', keyFieldType: $pb.PbFieldType.OS, valueFieldType: $pb.PbFieldType.OS, packageName: const $pb.PackageName('lemon_tea.common'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatInfo clone() => ChatInfo()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChatInfo copyWith(void Function(ChatInfo) updates) => super.copyWith((message) => updates(message as ChatInfo)) as ChatInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChatInfo create() => ChatInfo._();
  @$core.override
  ChatInfo createEmptyInstance() => create();
  static $pb.PbList<ChatInfo> createRepeated() => $pb.PbList<ChatInfo>();
  @$core.pragma('dart2js:noInline')
  static ChatInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ChatInfo>(create);
  static ChatInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get chatUuid => $_getSZ(0);
  @$pb.TagNumber(1)
  set chatUuid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasChatUuid() => $_has(0);
  @$pb.TagNumber(1)
  void clearChatUuid() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get title => $_getSZ(1);
  @$pb.TagNumber(2)
  set title($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTitle() => $_has(1);
  @$pb.TagNumber(2)
  void clearTitle() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get modelId => $_getI64(2);
  @$pb.TagNumber(3)
  set modelId($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasModelId() => $_has(2);
  @$pb.TagNumber(3)
  void clearModelId() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get createdAt => $_getI64(3);
  @$pb.TagNumber(4)
  set createdAt($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasCreatedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearCreatedAt() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get updatedAt => $_getI64(4);
  @$pb.TagNumber(5)
  set updatedAt($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasUpdatedAt() => $_has(4);
  @$pb.TagNumber(5)
  void clearUpdatedAt() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.int get messageCount => $_getIZ(5);
  @$pb.TagNumber(6)
  set messageCount($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasMessageCount() => $_has(5);
  @$pb.TagNumber(6)
  void clearMessageCount() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get lastMessagePreview => $_getSZ(6);
  @$pb.TagNumber(7)
  set lastMessagePreview($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasLastMessagePreview() => $_has(6);
  @$pb.TagNumber(7)
  void clearLastMessagePreview() => $_clearField(7);

  @$pb.TagNumber(8)
  $pb.PbMap<$core.String, $core.String> get metadata => $_getMap(7);
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
