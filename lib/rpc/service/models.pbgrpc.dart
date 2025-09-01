// This is a generated file - do not edit.
//
// Generated from rpc/service/models.proto.

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

import 'models.pb.dart' as $0;

export 'models.pb.dart';

/// Models 服务，提供与 OpenAI Models API 对齐的能力
@$pb.GrpcServiceName('lemon_tea.server.Models')
class ModelsClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ModelsClient(super.channel, {super.options, super.interceptors});

  /// 列出可用模型（等价 OpenAI GET /v1/models）
  $grpc.ResponseFuture<$0.ModelsResponse> models($0.ModelsRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$models, request, options: options);
  }

    // method descriptors

  static final _$models = $grpc.ClientMethod<$0.ModelsRequest, $0.ModelsResponse>(
      '/lemon_tea.server.Models/models',
      ($0.ModelsRequest value) => value.writeToBuffer(),
      $0.ModelsResponse.fromBuffer);
}

@$pb.GrpcServiceName('lemon_tea.server.Models')
abstract class ModelsServiceBase extends $grpc.Service {
  $core.String get $name => 'lemon_tea.server.Models';

  ModelsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ModelsRequest, $0.ModelsResponse>(
        'models',
        models_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ModelsRequest.fromBuffer(value),
        ($0.ModelsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ModelsResponse> models_Pre($grpc.ServiceCall $call, $async.Future<$0.ModelsRequest> $request) async {
    return models($call, await $request);
  }

  $async.Future<$0.ModelsResponse> models($grpc.ServiceCall call, $0.ModelsRequest request);

}
