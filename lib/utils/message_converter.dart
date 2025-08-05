
import 'dart:convert';
import 'package:lemon_tea/utils/llm/models/message.dart' as llm_models;
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/rpc/common.pb.dart' as grpc_common;
import 'package:lemon_tea/rpc/common.pbenum.dart' as grpc_enum;

/// 消息转换工具类
/// 
/// 负责在Flutter的消息模型和protobuf消息模型之间进行转换
class MessageConverter {
  /// 将Flutter的MessageRole转换为protobuf的RoleType
  static grpc_enum.RoleType convertRole(MessageRole role) {
    switch (role) {
      case MessageRole.user:
        return grpc_enum.RoleType.ROLE_TYPE_USER;
      case MessageRole.assistant:
        return grpc_enum.RoleType.ROLE_TYPE_ASSISTANT;
      case MessageRole.system:
        return grpc_enum.RoleType.ROLE_TYPE_SYSTEM;
      case MessageRole.tool:
        return grpc_enum.RoleType.ROLE_TYPE_TOOL;
    }
  }

  /// 将protobuf的RoleType转换为Flutter的MessageRole
  static MessageRole convertRoleFromProto(grpc_enum.RoleType role) {
    switch (role) {
      case grpc_enum.RoleType.ROLE_TYPE_USER:
        return MessageRole.user;
      case grpc_enum.RoleType.ROLE_TYPE_ASSISTANT:
        return MessageRole.assistant;
      case grpc_enum.RoleType.ROLE_TYPE_SYSTEM:
        return MessageRole.system;
      case grpc_enum.RoleType.ROLE_TYPE_TOOL:
        return MessageRole.tool;
      default:
        return MessageRole.user;
    }
  }

  /// 将文件类型转换为protobuf的ChatMessagePartType
  static grpc_enum.ChatMessagePartType convertFileType(String fileType) {
    switch (fileType.toLowerCase()) {
      case 'image':
        return grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_IMAGE_URL;
      case 'audio':
        return grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_AUDIO_URL;
      case 'video':
        return grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_VIDEO_URL;
      case 'text':
        return grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_TEXT;
      case 'document':
      case 'file':
      case 'other':
      default:
        return grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_FILE_URL;
    }
  }

  /// 将FileContent转换为ChatMessagePart
  static grpc_common.ChatMessagePart convertFileContent(llm_models.FileContent file) {
    final partType = convertFileType(file.type);
    
    switch (partType) {
      case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_TEXT:
        return grpc_common.ChatMessagePart(
          type: partType,
          text: _extractTextFromFile(file),
        );
      
      case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_IMAGE_URL:
        return grpc_common.ChatMessagePart(
          type: partType,
          imageUrl: grpc_common.ChatMessageImageURL(
            url: file.url ?? '',
            uri: _createDataUri(file),
            mimeType: file.mimeType,
            detail: grpc_enum.ImageURLDetail.IMAGE_URL_DETAIL_AUTO,
          ),
        );
      
      case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_AUDIO_URL:
        return grpc_common.ChatMessagePart(
          type: partType,
          audioUrl: grpc_common.ChatMessageAudioURL(
            url: file.url ?? '',
            uri: _createDataUri(file),
            mimeType: file.mimeType,
          ),
        );
      
      case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_VIDEO_URL:
        return grpc_common.ChatMessagePart(
          type: partType,
          videoUrl: grpc_common.ChatMessageVideoURL(
            url: file.url ?? '',
            uri: _createDataUri(file),
            mimeType: file.mimeType,
          ),
        );
      
      case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_FILE_URL:
      default:
        return grpc_common.ChatMessagePart(
          type: grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_FILE_URL,
          fileUrl: grpc_common.ChatMessageFileURL(
            url: file.url ?? '',
            uri: _createDataUri(file),
            mimeType: file.mimeType,
            name: file.name,
          ),
        );
    }
  }

  /// 从文件中提取文本内容
  static String _extractTextFromFile(llm_models.FileContent file) {
    if (file.data != null && file.data!.isNotEmpty) {
      try {
        // 解码 base64 数据
        final bytes = base64Decode(file.data!);
        // 尝试使用 UTF-8 解码为文本
        final text = utf8.decode(bytes, allowMalformed: true);
        
        // 添加文件名作为上下文（如果需要的话）
        if (file.name.isNotEmpty && !file.name.startsWith('temp_')) {
          return '// 文件: ${file.name}\n$text';
        }
        
        return text;
      } catch (e) {
        // 如果解码失败，返回错误信息
        return '// 无法读取文件内容: ${file.name}\n// 错误: $e';
      }
    }
    
    // 如果没有数据，返回文件名信息
    return '// 文件: ${file.name} (无内容数据)';
  }

  /// 创建data URI
  static String _createDataUri(llm_models.FileContent file) {
    if (file.data != null && file.data!.isNotEmpty) {
      return 'data:${file.mimeType};base64,${file.data}';
    }
    return file.url ?? '';
  }

  /// 将Flutter的Message转换为protobuf的Message
  static grpc_common.Message convertMessage(llm_models.Message message) {
    List<grpc_common.ChatMessagePart> multiContent = [];
    
    // 如果有文本内容，添加文本部分
    if (message.content.isNotEmpty) {
      multiContent.add(grpc_common.ChatMessagePart(
        type: grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_TEXT,
        text: message.content,
      ));
    }
    
    // 如果有文件内容，添加文件部分
    if (message.files != null && message.files!.isNotEmpty) {
      for (final file in message.files!) {
        multiContent.add(convertFileContent(file));
      }
    }
    
    return grpc_common.Message(
      role: convertRole(message.role),
      content: message.content, // 保持向后兼容
      multiContent: multiContent,
      reasoningContent: message.reasoningContent ?? '',
      // TODO: 处理toolCalls转换
    );
  }

  /// 批量转换消息列表
  static List<grpc_common.Message> convertMessages(List<llm_models.Message> messages) {
    return messages.map((msg) => convertMessage(msg)).toList();
  }

  /// 将protobuf的Message转换为Flutter的Message（用于接收响应）
  static llm_models.Message convertMessageFromProto(grpc_common.Message protoMessage) {
    String content = protoMessage.content;
    List<llm_models.FileContent> files = [];
    
    // 解析multiContent
    for (final part in protoMessage.multiContent) {
      switch (part.type) {
        case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_TEXT:
          if (part.hasText()) {
            // 对于普通文本部分，直接追加到内容中
            if (content.isNotEmpty) {
              content += '\n${part.text}';
            } else {
              content = part.text;
            }
          }
          break;
        case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_IMAGE_URL:
          if (part.hasImageUrl()) {
            files.add(_convertImageUrlToFileContent(part.imageUrl));
          }
          break;
        case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_AUDIO_URL:
          if (part.hasAudioUrl()) {
            files.add(_convertAudioUrlToFileContent(part.audioUrl));
          }
          break;
        case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_VIDEO_URL:
          if (part.hasVideoUrl()) {
            files.add(_convertVideoUrlToFileContent(part.videoUrl));
          }
          break;
        case grpc_enum.ChatMessagePartType.CHAT_MESSAGE_PART_TYPE_FILE_URL:
          if (part.hasFileUrl()) {
            files.add(_convertFileUrlToFileContent(part.fileUrl));
          }
          break;
        default:
          break;
      }
    }
    
    return llm_models.Message(
      role: convertRoleFromProto(protoMessage.role),
      content: content,
      reasoningContent: protoMessage.hasReasoningContent() ? protoMessage.reasoningContent : null,
      files: files.isNotEmpty ? files : null,
      // TODO: 处理toolCalls转换
    );
  }

  /// 将ImageURL转换为FileContent
  static llm_models.FileContent _convertImageUrlToFileContent(grpc_common.ChatMessageImageURL imageUrl) {
    return llm_models.FileContent(
      name: _extractFilenameFromUrl(imageUrl.url),
      mimeType: imageUrl.mimeType,
      type: 'image',
      size: 0, // protobuf中没有size信息
      url: imageUrl.url.isNotEmpty ? imageUrl.url : null,
      data: _extractBase64FromDataUri(imageUrl.uri),
    );
  }

  /// 将AudioURL转换为FileContent
  static llm_models.FileContent _convertAudioUrlToFileContent(grpc_common.ChatMessageAudioURL audioUrl) {
    return llm_models.FileContent(
      name: _extractFilenameFromUrl(audioUrl.url),
      mimeType: audioUrl.mimeType,
      type: 'audio',
      size: 0,
      url: audioUrl.url.isNotEmpty ? audioUrl.url : null,
      data: _extractBase64FromDataUri(audioUrl.uri),
    );
  }

  /// 将VideoURL转换为FileContent
  static llm_models.FileContent _convertVideoUrlToFileContent(grpc_common.ChatMessageVideoURL videoUrl) {
    return llm_models.FileContent(
      name: _extractFilenameFromUrl(videoUrl.url),
      mimeType: videoUrl.mimeType,
      type: 'video',
      size: 0,
      url: videoUrl.url.isNotEmpty ? videoUrl.url : null,
      data: _extractBase64FromDataUri(videoUrl.uri),
    );
  }

  /// 将FileURL转换为FileContent
  static llm_models.FileContent _convertFileUrlToFileContent(grpc_common.ChatMessageFileURL fileUrl) {
    return llm_models.FileContent(
      name: fileUrl.name.isNotEmpty ? fileUrl.name : _extractFilenameFromUrl(fileUrl.url),
      mimeType: fileUrl.mimeType,
      type: 'file',
      size: 0,
      url: fileUrl.url.isNotEmpty ? fileUrl.url : null,
      data: _extractBase64FromDataUri(fileUrl.uri),
    );
  }

  /// 从URL中提取文件名
  static String _extractFilenameFromUrl(String url) {
    if (url.isEmpty) return 'unknown';
    final uri = Uri.tryParse(url);
    if (uri != null && uri.pathSegments.isNotEmpty) {
      return uri.pathSegments.last;
    }
    return 'unknown';
  }

  /// 从data URI中提取base64数据
  static String? _extractBase64FromDataUri(String dataUri) {
    if (dataUri.startsWith('data:')) {
      final commaIndex = dataUri.indexOf(',');
      if (commaIndex != -1) {
        return dataUri.substring(commaIndex + 1);
      }
    }
    return null;
  }
}