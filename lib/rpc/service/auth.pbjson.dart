// This is a generated file - do not edit.
//
// Generated from rpc/service/auth.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use fieldTypeDescriptor instead')
const FieldType$json = {
  '1': 'FieldType',
  '2': [
    {'1': 'FIELD_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'FIELD_TYPE_USERNAME', '2': 1},
    {'1': 'FIELD_TYPE_EMAIL', '2': 2},
  ],
};

/// Descriptor for `FieldType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List fieldTypeDescriptor = $convert.base64Decode(
    'CglGaWVsZFR5cGUSGgoWRklFTERfVFlQRV9VTlNQRUNJRklFRBAAEhcKE0ZJRUxEX1RZUEVfVV'
    'NFUk5BTUUQARIUChBGSUVMRF9UWVBFX0VNQUlMEAI=');

@$core.Deprecated('Use sendEmailVerificationCodeRequestDescriptor instead')
const SendEmailVerificationCodeRequest$json = {
  '1': 'SendEmailVerificationCodeRequest',
  '2': [
    {'1': 'email', '3': 1, '4': 1, '5': 9, '10': 'email'},
    {'1': 'username', '3': 2, '4': 1, '5': 9, '10': 'username'},
    {'1': 'code_type', '3': 3, '4': 1, '5': 14, '6': '.lemon_tea.common.VerificationCodeType', '10': 'codeType'},
  ],
};

/// Descriptor for `SendEmailVerificationCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendEmailVerificationCodeRequestDescriptor = $convert.base64Decode(
    'CiBTZW5kRW1haWxWZXJpZmljYXRpb25Db2RlUmVxdWVzdBIUCgVlbWFpbBgBIAEoCVIFZW1haW'
    'wSGgoIdXNlcm5hbWUYAiABKAlSCHVzZXJuYW1lEkMKCWNvZGVfdHlwZRgDIAEoDjImLmxlbW9u'
    'X3RlYS5jb21tb24uVmVyaWZpY2F0aW9uQ29kZVR5cGVSCGNvZGVUeXBl');

@$core.Deprecated('Use sendEmailVerificationCodeResponseDescriptor instead')
const SendEmailVerificationCodeResponse$json = {
  '1': 'SendEmailVerificationCodeResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SendEmailVerificationCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendEmailVerificationCodeResponseDescriptor = $convert.base64Decode(
    'CiFTZW5kRW1haWxWZXJpZmljYXRpb25Db2RlUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3'
    'VjY2VzcxIYCgdtZXNzYWdlGAIgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use checkFieldAvailabilityRequestDescriptor instead')
const CheckFieldAvailabilityRequest$json = {
  '1': 'CheckFieldAvailabilityRequest',
  '2': [
    {'1': 'field_type', '3': 1, '4': 1, '5': 14, '6': '.lemon_tea.server.FieldType', '10': 'fieldType'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `CheckFieldAvailabilityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkFieldAvailabilityRequestDescriptor = $convert.base64Decode(
    'Ch1DaGVja0ZpZWxkQXZhaWxhYmlsaXR5UmVxdWVzdBI6CgpmaWVsZF90eXBlGAEgASgOMhsubG'
    'Vtb25fdGVhLnNlcnZlci5GaWVsZFR5cGVSCWZpZWxkVHlwZRIUCgV2YWx1ZRgCIAEoCVIFdmFs'
    'dWU=');

@$core.Deprecated('Use checkFieldAvailabilityResponseDescriptor instead')
const CheckFieldAvailabilityResponse$json = {
  '1': 'CheckFieldAvailabilityResponse',
  '2': [
    {'1': 'available', '3': 1, '4': 1, '5': 8, '10': 'available'},
  ],
};

/// Descriptor for `CheckFieldAvailabilityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkFieldAvailabilityResponseDescriptor = $convert.base64Decode(
    'Ch5DaGVja0ZpZWxkQXZhaWxhYmlsaXR5UmVzcG9uc2USHAoJYXZhaWxhYmxlGAEgASgIUglhdm'
    'FpbGFibGU=');

@$core.Deprecated('Use registerRequestDescriptor instead')
const RegisterRequest$json = {
  '1': 'RegisterRequest',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {'1': 'password_md5', '3': 2, '4': 1, '5': 9, '10': 'passwordMd5'},
    {'1': 'email', '3': 3, '4': 1, '5': 9, '10': 'email'},
    {'1': 'email_verification_code', '3': 4, '4': 1, '5': 3, '10': 'emailVerificationCode'},
  ],
};

/// Descriptor for `RegisterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerRequestDescriptor = $convert.base64Decode(
    'Cg9SZWdpc3RlclJlcXVlc3QSGgoIdXNlcm5hbWUYASABKAlSCHVzZXJuYW1lEiEKDHBhc3N3b3'
    'JkX21kNRgCIAEoCVILcGFzc3dvcmRNZDUSFAoFZW1haWwYAyABKAlSBWVtYWlsEjYKF2VtYWls'
    'X3ZlcmlmaWNhdGlvbl9jb2RlGAQgASgDUhVlbWFpbFZlcmlmaWNhdGlvbkNvZGU=');

@$core.Deprecated('Use registerResponseDescriptor instead')
const RegisterResponse$json = {
  '1': 'RegisterResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {'1': 'user_id', '3': 3, '4': 1, '5': 9, '9': 0, '10': 'userId', '17': true},
  ],
  '8': [
    {'1': '_user_id'},
  ],
};

/// Descriptor for `RegisterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerResponseDescriptor = $convert.base64Decode(
    'ChBSZWdpc3RlclJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSGAoHbWVzc2FnZR'
    'gCIAEoCVIHbWVzc2FnZRIcCgd1c2VyX2lkGAMgASgJSABSBnVzZXJJZIgBAUIKCghfdXNlcl9p'
    'ZA==');

@$core.Deprecated('Use loginRequestDescriptor instead')
const LoginRequest$json = {
  '1': 'LoginRequest',
  '2': [
    {'1': 'login_field', '3': 1, '4': 1, '5': 9, '10': 'loginField'},
    {'1': 'password_md5', '3': 2, '4': 1, '5': 9, '10': 'passwordMd5'},
  ],
};

/// Descriptor for `LoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginRequestDescriptor = $convert.base64Decode(
    'CgxMb2dpblJlcXVlc3QSHwoLbG9naW5fZmllbGQYASABKAlSCmxvZ2luRmllbGQSIQoMcGFzc3'
    'dvcmRfbWQ1GAIgASgJUgtwYXNzd29yZE1kNQ==');

@$core.Deprecated('Use loginResponseDescriptor instead')
const LoginResponse$json = {
  '1': 'LoginResponse',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '9': 0, '10': 'accessToken', '17': true},
    {'1': 'expires_in', '3': 2, '4': 1, '5': 3, '9': 1, '10': 'expiresIn', '17': true},
    {'1': 'user_info', '3': 3, '4': 1, '5': 11, '6': '.lemon_tea.server.UserInfo', '9': 2, '10': 'userInfo', '17': true},
  ],
  '8': [
    {'1': '_access_token'},
    {'1': '_expires_in'},
    {'1': '_user_info'},
  ],
};

/// Descriptor for `LoginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginResponseDescriptor = $convert.base64Decode(
    'Cg1Mb2dpblJlc3BvbnNlEiYKDGFjY2Vzc190b2tlbhgBIAEoCUgAUgthY2Nlc3NUb2tlbogBAR'
    'IiCgpleHBpcmVzX2luGAIgASgDSAFSCWV4cGlyZXNJbogBARI8Cgl1c2VyX2luZm8YAyABKAsy'
    'Gi5sZW1vbl90ZWEuc2VydmVyLlVzZXJJbmZvSAJSCHVzZXJJbmZviAEBQg8KDV9hY2Nlc3NfdG'
    '9rZW5CDQoLX2V4cGlyZXNfaW5CDAoKX3VzZXJfaW5mbw==');

@$core.Deprecated('Use resetPasswordRequestDescriptor instead')
const ResetPasswordRequest$json = {
  '1': 'ResetPasswordRequest',
  '2': [
    {'1': 'login_field', '3': 1, '4': 1, '5': 9, '10': 'loginField'},
    {'1': 'email_verification_code', '3': 2, '4': 1, '5': 9, '10': 'emailVerificationCode'},
    {'1': 'new_password_md5', '3': 3, '4': 1, '5': 9, '10': 'newPasswordMd5'},
  ],
};

/// Descriptor for `ResetPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetPasswordRequestDescriptor = $convert.base64Decode(
    'ChRSZXNldFBhc3N3b3JkUmVxdWVzdBIfCgtsb2dpbl9maWVsZBgBIAEoCVIKbG9naW5GaWVsZB'
    'I2ChdlbWFpbF92ZXJpZmljYXRpb25fY29kZRgCIAEoCVIVZW1haWxWZXJpZmljYXRpb25Db2Rl'
    'EigKEG5ld19wYXNzd29yZF9tZDUYAyABKAlSDm5ld1Bhc3N3b3JkTWQ1');

@$core.Deprecated('Use logoutRequestDescriptor instead')
const LogoutRequest$json = {
  '1': 'LogoutRequest',
};

/// Descriptor for `LogoutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutRequestDescriptor = $convert.base64Decode(
    'Cg1Mb2dvdXRSZXF1ZXN0');

@$core.Deprecated('Use userInfoDescriptor instead')
const UserInfo$json = {
  '1': 'UserInfo',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 3, '10': 'userId'},
    {'1': 'username', '3': 2, '4': 1, '5': 9, '10': 'username'},
    {'1': 'email', '3': 3, '4': 1, '5': 9, '10': 'email'},
    {'1': 'created_at', '3': 4, '4': 1, '5': 3, '10': 'createdAt'},
    {'1': 'updated_at', '3': 5, '4': 1, '5': 3, '10': 'updatedAt'},
  ],
};

/// Descriptor for `UserInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userInfoDescriptor = $convert.base64Decode(
    'CghVc2VySW5mbxIXCgd1c2VyX2lkGAEgASgDUgZ1c2VySWQSGgoIdXNlcm5hbWUYAiABKAlSCH'
    'VzZXJuYW1lEhQKBWVtYWlsGAMgASgJUgVlbWFpbBIdCgpjcmVhdGVkX2F0GAQgASgDUgljcmVh'
    'dGVkQXQSHQoKdXBkYXRlZF9hdBgFIAEoA1IJdXBkYXRlZEF0');

