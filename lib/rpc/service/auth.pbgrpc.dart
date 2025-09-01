// This is a generated file - do not edit.
//
// Generated from rpc/service/auth.proto.

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
import 'auth.pb.dart' as $0;

export 'auth.pb.dart';

@$pb.GrpcServiceName('lemon_tea.server.Auth')
class AuthClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AuthClient(super.channel, {super.options, super.interceptors});

  /// SendEmailVerificationCode 发送邮件验证码
  $grpc.ResponseFuture<$1.Empty> sendEmailVerificationCode($0.SendEmailVerificationCodeRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$sendEmailVerificationCode, request, options: options);
  }

  /// CheckFieldAvailability 检查用户名或邮箱是否被使用
  $grpc.ResponseFuture<$0.CheckFieldAvailabilityResponse> checkFieldAvailability($0.CheckFieldAvailabilityRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$checkFieldAvailability, request, options: options);
  }

  /// Register 注册
  $grpc.ResponseFuture<$1.Empty> register($0.RegisterRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$register, request, options: options);
  }

  /// Login 登录
  $grpc.ResponseFuture<$0.LoginResponse> login($0.LoginRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$login, request, options: options);
  }

  /// ResetPassword 忘记密码
  $grpc.ResponseFuture<$1.Empty> resetPassword($0.ResetPasswordRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$resetPassword, request, options: options);
  }

  /// Logout 登出
  $grpc.ResponseFuture<$1.Empty> logout($0.LogoutRequest request, {$grpc.CallOptions? options,}) {
    return $createUnaryCall(_$logout, request, options: options);
  }

    // method descriptors

  static final _$sendEmailVerificationCode = $grpc.ClientMethod<$0.SendEmailVerificationCodeRequest, $1.Empty>(
      '/lemon_tea.server.Auth/SendEmailVerificationCode',
      ($0.SendEmailVerificationCodeRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$checkFieldAvailability = $grpc.ClientMethod<$0.CheckFieldAvailabilityRequest, $0.CheckFieldAvailabilityResponse>(
      '/lemon_tea.server.Auth/CheckFieldAvailability',
      ($0.CheckFieldAvailabilityRequest value) => value.writeToBuffer(),
      $0.CheckFieldAvailabilityResponse.fromBuffer);
  static final _$register = $grpc.ClientMethod<$0.RegisterRequest, $1.Empty>(
      '/lemon_tea.server.Auth/Register',
      ($0.RegisterRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/lemon_tea.server.Auth/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$resetPassword = $grpc.ClientMethod<$0.ResetPasswordRequest, $1.Empty>(
      '/lemon_tea.server.Auth/ResetPassword',
      ($0.ResetPasswordRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$logout = $grpc.ClientMethod<$0.LogoutRequest, $1.Empty>(
      '/lemon_tea.server.Auth/Logout',
      ($0.LogoutRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('lemon_tea.server.Auth')
abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'lemon_tea.server.Auth';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SendEmailVerificationCodeRequest, $1.Empty>(
        'SendEmailVerificationCode',
        sendEmailVerificationCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SendEmailVerificationCodeRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckFieldAvailabilityRequest, $0.CheckFieldAvailabilityResponse>(
        'CheckFieldAvailability',
        checkFieldAvailability_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CheckFieldAvailabilityRequest.fromBuffer(value),
        ($0.CheckFieldAvailabilityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterRequest, $1.Empty>(
        'Register',
        register_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.RegisterRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResetPasswordRequest, $1.Empty>(
        'ResetPassword',
        resetPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ResetPasswordRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogoutRequest, $1.Empty>(
        'Logout',
        logout_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogoutRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> sendEmailVerificationCode_Pre($grpc.ServiceCall $call, $async.Future<$0.SendEmailVerificationCodeRequest> $request) async {
    return sendEmailVerificationCode($call, await $request);
  }

  $async.Future<$1.Empty> sendEmailVerificationCode($grpc.ServiceCall call, $0.SendEmailVerificationCodeRequest request);

  $async.Future<$0.CheckFieldAvailabilityResponse> checkFieldAvailability_Pre($grpc.ServiceCall $call, $async.Future<$0.CheckFieldAvailabilityRequest> $request) async {
    return checkFieldAvailability($call, await $request);
  }

  $async.Future<$0.CheckFieldAvailabilityResponse> checkFieldAvailability($grpc.ServiceCall call, $0.CheckFieldAvailabilityRequest request);

  $async.Future<$1.Empty> register_Pre($grpc.ServiceCall $call, $async.Future<$0.RegisterRequest> $request) async {
    return register($call, await $request);
  }

  $async.Future<$1.Empty> register($grpc.ServiceCall call, $0.RegisterRequest request);

  $async.Future<$0.LoginResponse> login_Pre($grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login($grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$1.Empty> resetPassword_Pre($grpc.ServiceCall $call, $async.Future<$0.ResetPasswordRequest> $request) async {
    return resetPassword($call, await $request);
  }

  $async.Future<$1.Empty> resetPassword($grpc.ServiceCall call, $0.ResetPasswordRequest request);

  $async.Future<$1.Empty> logout_Pre($grpc.ServiceCall $call, $async.Future<$0.LogoutRequest> $request) async {
    return logout($call, await $request);
  }

  $async.Future<$1.Empty> logout($grpc.ServiceCall call, $0.LogoutRequest request);

}
