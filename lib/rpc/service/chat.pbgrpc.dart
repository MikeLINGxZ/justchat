// This is a generated file - do not edit.
//
// Generated from rpc/service/chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import '../common/common.pb.dart' as $1;
import 'chat.pb.dart' as $0;

export 'chat.pb.dart';

@$pb.GrpcServiceName('lemon_tea.server.Chat')
class ChatClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ChatClient(super.channel, {super.options, super.interceptors});

  /// 流式 Chat Completions（等价于 OpenAI /v1/chat/completions + stream）
  $grpc.ResponseStream<$0.CompletionsResponse> completions($0.CompletionsRequest request, {$grpc.CallOptions? options,}) {
    return $createStreamingCall(_$completions, $async.Stream.fromIterable([request]), options: options);
  }

  /// 获取对话标题
  $grpc.ResponseFuture<$0.ChatTitleResponse> chatTitle($0.ChatTitleRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$chatTitle, request, options: options);
  }

  /// 保存对话标题
  $grpc.ResponseFuture<$1.Empty> chatTitleSave($0.ChatTitleSaveRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$chatTitleSave, request, options: options);
  }

  /// 获取历史对话记录列表
  $grpc.ResponseFuture<$0.ListChatsResponse> listChats($0.ListChatsRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$listChats, request, options: options);
  }

  /// 删除对话
  $grpc.ResponseFuture<$1.Empty> deleteChat($0.DeleteChatRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$deleteChat, request, options: options);
  }

  /// 获取对话消息
  $grpc.ResponseFuture<$0.GetChatMessagesResponse> getChatMessages($0.GetChatMessagesRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$getChatMessages, request, options: options);
  }

  /// 删除对话消息
  $grpc.ResponseFuture<$0.DeleteChatMessageResponse> deleteChatMessage($0.DeleteChatMessageRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$deleteChatMessage, request, options: options);
  }

    // method descriptors

  static final _$completions = $grpc.ClientMethod<$0.CompletionsRequest, $0.CompletionsResponse>(
      '/lemon_tea.server.Chat/Completions',
      ($0.CompletionsRequest value) => value.writeToBuffer(),
      $0.CompletionsResponse.fromBuffer);
  static final _$chatTitle = $grpc.ClientMethod<$0.ChatTitleRequest, $0.ChatTitleResponse>(
      '/lemon_tea.server.Chat/ChatTitle',
      ($0.ChatTitleRequest value) => value.writeToBuffer(),
      $0.ChatTitleResponse.fromBuffer);
  static final _$chatTitleSave = $grpc.ClientMethod<$0.ChatTitleSaveRequest, $1.Empty>(
      '/lemon_tea.server.Chat/ChatTitleSave',
      ($0.ChatTitleSaveRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$listChats = $grpc.ClientMethod<$0.ListChatsRequest, $0.ListChatsResponse>(
      '/lemon_tea.server.Chat/ListChats',
      ($0.ListChatsRequest value) => value.writeToBuffer(),
      $0.ListChatsResponse.fromBuffer);
  static final _$deleteChat = $grpc.ClientMethod<$0.DeleteChatRequest, $1.Empty>(
      '/lemon_tea.server.Chat/DeleteChat',
      ($0.DeleteChatRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$getChatMessages = $grpc.ClientMethod<$0.GetChatMessagesRequest, $0.GetChatMessagesResponse>(
      '/lemon_tea.server.Chat/GetChatMessages',
      ($0.GetChatMessagesRequest value) => value.writeToBuffer(),
      $0.GetChatMessagesResponse.fromBuffer);
  static final _$deleteChatMessage = $grpc.ClientMethod<$0.DeleteChatMessageRequest, $0.DeleteChatMessageResponse>(
      '/lemon_tea.server.Chat/DeleteChatMessage',
      ($0.DeleteChatMessageRequest value) => value.writeToBuffer(),
      $0.DeleteChatMessageResponse.fromBuffer);
}

@$pb.GrpcServiceName('lemon_tea.server.Chat')
abstract class ChatServiceBase extends $grpc.Service {
  $core.String get $name => 'lemon_tea.server.Chat';

  ChatServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CompletionsRequest, $0.CompletionsResponse>(
        'Completions',
        completions_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.CompletionsRequest.fromBuffer(value),
        ($0.CompletionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChatTitleRequest, $0.ChatTitleResponse>(
        'ChatTitle',
        chatTitle_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ChatTitleRequest.fromBuffer(value),
        ($0.ChatTitleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChatTitleSaveRequest, $1.Empty>(
        'ChatTitleSave',
        chatTitleSave_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ChatTitleSaveRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListChatsRequest, $0.ListChatsResponse>(
        'ListChats',
        listChats_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListChatsRequest.fromBuffer(value),
        ($0.ListChatsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteChatRequest, $1.Empty>(
        'DeleteChat',
        deleteChat_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteChatRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChatMessagesRequest, $0.GetChatMessagesResponse>(
        'GetChatMessages',
        getChatMessages_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetChatMessagesRequest.fromBuffer(value),
        ($0.GetChatMessagesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteChatMessageRequest, $0.DeleteChatMessageResponse>(
        'DeleteChatMessage',
        deleteChatMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteChatMessageRequest.fromBuffer(value),
        ($0.DeleteChatMessageResponse value) => value.writeToBuffer()));
  }

  $async.Stream<$0.CompletionsResponse> completions_Pre($grpc.ServiceCall $call, $async.Future<$0.CompletionsRequest> $request) async* {
    yield* completions($call, await $request);
  }

  $async.Stream<$0.CompletionsResponse> completions($grpc.ServiceCall call, $0.CompletionsRequest request);

  $async.Future<$0.ChatTitleResponse> chatTitle_Pre($grpc.ServiceCall $call, $async.Future<$0.ChatTitleRequest> $request) async {
    return chatTitle($call, await $request);
  }

  $async.Future<$0.ChatTitleResponse> chatTitle($grpc.ServiceCall call, $0.ChatTitleRequest request);

  $async.Future<$1.Empty> chatTitleSave_Pre($grpc.ServiceCall $call, $async.Future<$0.ChatTitleSaveRequest> $request) async {
    return chatTitleSave($call, await $request);
  }

  $async.Future<$1.Empty> chatTitleSave($grpc.ServiceCall call, $0.ChatTitleSaveRequest request);

  $async.Future<$0.ListChatsResponse> listChats_Pre($grpc.ServiceCall $call, $async.Future<$0.ListChatsRequest> $request) async {
    return listChats($call, await $request);
  }

  $async.Future<$0.ListChatsResponse> listChats($grpc.ServiceCall call, $0.ListChatsRequest request);

  $async.Future<$1.Empty> deleteChat_Pre($grpc.ServiceCall $call, $async.Future<$0.DeleteChatRequest> $request) async {
    return deleteChat($call, await $request);
  }

  $async.Future<$1.Empty> deleteChat($grpc.ServiceCall call, $0.DeleteChatRequest request);

  $async.Future<$0.GetChatMessagesResponse> getChatMessages_Pre($grpc.ServiceCall $call, $async.Future<$0.GetChatMessagesRequest> $request) async {
    return getChatMessages($call, await $request);
  }

  $async.Future<$0.GetChatMessagesResponse> getChatMessages($grpc.ServiceCall call, $0.GetChatMessagesRequest request);

  $async.Future<$0.DeleteChatMessageResponse> deleteChatMessage_Pre($grpc.ServiceCall $call, $async.Future<$0.DeleteChatMessageRequest> $request) async {
    return deleteChatMessage($call, await $request);
  }

  $async.Future<$0.DeleteChatMessageResponse> deleteChatMessage($grpc.ServiceCall call, $0.DeleteChatMessageRequest request);

}
