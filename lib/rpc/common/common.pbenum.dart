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

import 'package:protobuf/protobuf.dart' as $pb;

/// ImageURLDetail 是图像 URL 的细节级别
class ImageURLDetail extends $pb.ProtobufEnum {
  static const ImageURLDetail IMAGE_URL_DETAIL_UNSPECIFIED = ImageURLDetail._(0, _omitEnumNames ? '' : 'IMAGE_URL_DETAIL_UNSPECIFIED');
  static const ImageURLDetail IMAGE_URL_DETAIL_HIGH = ImageURLDetail._(1, _omitEnumNames ? '' : 'IMAGE_URL_DETAIL_HIGH');
  static const ImageURLDetail IMAGE_URL_DETAIL_LOW = ImageURLDetail._(2, _omitEnumNames ? '' : 'IMAGE_URL_DETAIL_LOW');
  static const ImageURLDetail IMAGE_URL_DETAIL_AUTO = ImageURLDetail._(3, _omitEnumNames ? '' : 'IMAGE_URL_DETAIL_AUTO');

  static const $core.List<ImageURLDetail> values = <ImageURLDetail> [
    IMAGE_URL_DETAIL_UNSPECIFIED,
    IMAGE_URL_DETAIL_HIGH,
    IMAGE_URL_DETAIL_LOW,
    IMAGE_URL_DETAIL_AUTO,
  ];

  static final $core.List<ImageURLDetail?> _byValue = $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ImageURLDetail? valueOf($core.int value) =>  value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImageURLDetail._(super.value, super.name);
}

/// ChatMessagePartType 是聊天消息中内容片段的类型
class ChatMessagePartType extends $pb.ProtobufEnum {
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_UNSPECIFIED = ChatMessagePartType._(0, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_UNSPECIFIED');
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_TEXT = ChatMessagePartType._(1, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_TEXT');
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_IMAGE_URL = ChatMessagePartType._(2, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_IMAGE_URL');
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_AUDIO_URL = ChatMessagePartType._(3, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_AUDIO_URL');
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_VIDEO_URL = ChatMessagePartType._(4, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_VIDEO_URL');
  static const ChatMessagePartType CHAT_MESSAGE_PART_TYPE_FILE_URL = ChatMessagePartType._(5, _omitEnumNames ? '' : 'CHAT_MESSAGE_PART_TYPE_FILE_URL');

  static const $core.List<ChatMessagePartType> values = <ChatMessagePartType> [
    CHAT_MESSAGE_PART_TYPE_UNSPECIFIED,
    CHAT_MESSAGE_PART_TYPE_TEXT,
    CHAT_MESSAGE_PART_TYPE_IMAGE_URL,
    CHAT_MESSAGE_PART_TYPE_AUDIO_URL,
    CHAT_MESSAGE_PART_TYPE_VIDEO_URL,
    CHAT_MESSAGE_PART_TYPE_FILE_URL,
  ];

  static final $core.List<ChatMessagePartType?> _byValue = $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ChatMessagePartType? valueOf($core.int value) =>  value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChatMessagePartType._(super.value, super.name);
}

class VerificationCodeType extends $pb.ProtobufEnum {
  static const VerificationCodeType VERIFICATION_CODE_TYPE_UNSPECIFIED = VerificationCodeType._(0, _omitEnumNames ? '' : 'VERIFICATION_CODE_TYPE_UNSPECIFIED');
  static const VerificationCodeType VERIFICATION_CODE_TYPE_REGISTER = VerificationCodeType._(1, _omitEnumNames ? '' : 'VERIFICATION_CODE_TYPE_REGISTER');
  static const VerificationCodeType VERIFICATION_CODE_TYPE_RESET_PASSWORD = VerificationCodeType._(2, _omitEnumNames ? '' : 'VERIFICATION_CODE_TYPE_RESET_PASSWORD');

  static const $core.List<VerificationCodeType> values = <VerificationCodeType> [
    VERIFICATION_CODE_TYPE_UNSPECIFIED,
    VERIFICATION_CODE_TYPE_REGISTER,
    VERIFICATION_CODE_TYPE_RESET_PASSWORD,
  ];

  static final $core.List<VerificationCodeType?> _byValue = $pb.ProtobufEnum.$_initByValueList(values, 2);
  static VerificationCodeType? valueOf($core.int value) =>  value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VerificationCodeType._(super.value, super.name);
}


const $core.bool _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
