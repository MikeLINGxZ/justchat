// This is a generated file - do not edit.
//
// Generated from rpc/service.proto.

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

import 'service.pb.dart' as $0;

export 'service.pb.dart';

@$pb.GrpcServiceName('lemon_tea.server.LemonTea')
class LemonTeaClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  LemonTeaClient(super.channel, {super.options, super.interceptors});

  /// UpdateLlmConfig 更新llm配置
  $grpc.ResponseFuture<$0.UpdateLlmConfigResponse> updateLlmConfig($0.UpdateLlmConfigRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$updateLlmConfig, request, options: options);
  }

  /// Models 获取llm的模型列表
  $grpc.ResponseFuture<$0.ModelsResponse> models($0.ModelsRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$models, request, options: options);
  }

  /// Chat 对话接口
  $grpc.ResponseStream<$0.ChatResponse> chat($async.Stream<$0.ChatRequest> request, {$grpc.CallOptions? options,}) {
    return $createStreamingCall(_$chat, request, options: options);
  }

  /// UploadFile 上传文件
  $grpc.ResponseFuture<$0.UploadFileResponse> uploadFile($0.UploadFileRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$uploadFile, request, options: options);
  }

    // method descriptors

  static final _$updateLlmConfig = $grpc.ClientMethod<$0.UpdateLlmConfigRequest, $0.UpdateLlmConfigResponse>(
      '/lemon_tea.server.LemonTea/UpdateLlmConfig',
      ($0.UpdateLlmConfigRequest value) => value.writeToBuffer(),
      $0.UpdateLlmConfigResponse.fromBuffer);
  static final _$models = $grpc.ClientMethod<$0.ModelsRequest, $0.ModelsResponse>(
      '/lemon_tea.server.LemonTea/Models',
      ($0.ModelsRequest value) => value.writeToBuffer(),
      $0.ModelsResponse.fromBuffer);
  static final _$chat = $grpc.ClientMethod<$0.ChatRequest, $0.ChatResponse>(
      '/lemon_tea.server.LemonTea/Chat',
      ($0.ChatRequest value) => value.writeToBuffer(),
      $0.ChatResponse.fromBuffer);
  static final _$uploadFile = $grpc.ClientMethod<$0.UploadFileRequest, $0.UploadFileResponse>(
      '/lemon_tea.server.LemonTea/UploadFile',
      ($0.UploadFileRequest value) => value.writeToBuffer(),
      $0.UploadFileResponse.fromBuffer);
}

@$pb.GrpcServiceName('lemon_tea.server.LemonTea')
abstract class LemonTeaServiceBase extends $grpc.Service {
  $core.String get $name => 'lemon_tea.server.LemonTea';

  LemonTeaServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.UpdateLlmConfigRequest, $0.UpdateLlmConfigResponse>(
        'UpdateLlmConfig',
        updateLlmConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateLlmConfigRequest.fromBuffer(value),
        ($0.UpdateLlmConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ModelsRequest, $0.ModelsResponse>(
        'Models',
        models_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ModelsRequest.fromBuffer(value),
        ($0.ModelsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChatRequest, $0.ChatResponse>(
        'Chat',
        chat,
        true,
        true,
        ($core.List<$core.int> value) => $0.ChatRequest.fromBuffer(value),
        ($0.ChatResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UploadFileRequest, $0.UploadFileResponse>(
        'UploadFile',
        uploadFile_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UploadFileRequest.fromBuffer(value),
        ($0.UploadFileResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.UpdateLlmConfigResponse> updateLlmConfig_Pre($grpc.ServiceCall $call, $async.Future<$0.UpdateLlmConfigRequest> $request) async {
    return updateLlmConfig($call, await $request);
  }

  $async.Future<$0.UpdateLlmConfigResponse> updateLlmConfig($grpc.ServiceCall call, $0.UpdateLlmConfigRequest request);

  $async.Future<$0.ModelsResponse> models_Pre($grpc.ServiceCall $call, $async.Future<$0.ModelsRequest> $request) async {
    return models($call, await $request);
  }

  $async.Future<$0.ModelsResponse> models($grpc.ServiceCall call, $0.ModelsRequest request);

  $async.Stream<$0.ChatResponse> chat($grpc.ServiceCall call, $async.Stream<$0.ChatRequest> request);

  $async.Future<$0.UploadFileResponse> uploadFile_Pre($grpc.ServiceCall $call, $async.Future<$0.UploadFileRequest> $request) async {
    return uploadFile($call, await $request);
  }

  $async.Future<$0.UploadFileResponse> uploadFile($grpc.ServiceCall call, $0.UploadFileRequest request);

}
