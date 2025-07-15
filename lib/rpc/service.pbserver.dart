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

import 'package:protobuf/protobuf.dart' as $pb;

import 'service.pb.dart' as $1;
import 'service.pbjson.dart';

export 'service.pb.dart';

abstract class LemonTeaServiceBase extends $pb.GeneratedService {
  $async.Future<$1.UpdateLlmConfigResponse> updateLlmConfig($pb.ServerContext ctx, $1.UpdateLlmConfigRequest request);
  $async.Future<$1.ModelsResponse> models($pb.ServerContext ctx, $1.ModelsRequest request);
  $async.Future<$1.ChatResponse> chat($pb.ServerContext ctx, $1.ChatRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'UpdateLlmConfig': return $1.UpdateLlmConfigRequest();
      case 'Models': return $1.ModelsRequest();
      case 'Chat': return $1.ChatRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'UpdateLlmConfig': return updateLlmConfig(ctx, request as $1.UpdateLlmConfigRequest);
      case 'Models': return models(ctx, request as $1.ModelsRequest);
      case 'Chat': return chat(ctx, request as $1.ChatRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => LemonTeaServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => LemonTeaServiceBase$messageJson;
}

