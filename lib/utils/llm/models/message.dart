import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/tool_call.dart';

part 'message.g.dart';

/// 文件内容类
@JsonSerializable()
class FileContent {
  /// 文件名
  final String name;
  
  /// MIME类型
  final String mimeType;
  
  /// 文件类型
  final String type;
  
  /// 文件二进制数据（Base64编码）
  final String? data;
  
  /// 文件大小（字节）
  final int size;
  
  /// 文件URL（如果是远程文件）
  final String? url;
  
  /// 文件描述
  final String? description;
  
  /// 本地文件路径（用于历史记录回显和重新生成）
  final String? localPath;

  const FileContent({
    required this.name,
    required this.mimeType,
    required this.type,
    this.data,
    required this.size,
    this.url,
    this.description,
    this.localPath,
  });

  factory FileContent.fromJson(Map<String, dynamic> json) =>
      _$FileContentFromJson(json);
  Map<String, dynamic> toJson() => _$FileContentToJson(this);

  FileContent copyWith({
    String? name,
    String? mimeType,
    String? type,
    String? data,
    int? size,
    String? url,
    String? description,
    String? localPath,
  }) => FileContent(
    name: name ?? this.name,
    mimeType: mimeType ?? this.mimeType,
    type: type ?? this.type,
    data: data ?? this.data,
    size: size ?? this.size,
    url: url ?? this.url,
    description: description ?? this.description,
    localPath: localPath ?? this.localPath,
  );
}

@JsonSerializable()
class Message {
  /// 消息角色
  final MessageRole role;

  /// 消息内容（对于function角色可能是空字符串）
  final String content;

  /// 思考过程内容（思维链模型返回的思考过程）
  final String? reasoningContent;

  /// 函数调用信息（role为assistant时可能存在）
  final List<ToolCall>? toolCalls;

  // 工具调用结果id
  String? toolCallId;

  /// 附件文件列表
  final List<FileContent>? files;

  /// 是否被用户停止生成
  final bool stoppedByUser;

  // todo visible
  

  Message({
    required this.role,
    required this.content,
    this.reasoningContent,
    this.toolCalls,
    this.toolCallId,
    this.files,
    this.stoppedByUser = false,
  }) : assert(
    // 允许assistant和tool角色的content为空（用于流式输出和工具调用）
    role == MessageRole.assistant || 
    role == MessageRole.tool ||
    content.isNotEmpty ||
    (files != null && files.isNotEmpty)
  );

  /// 创建只包含文本的消息
  factory Message.text({
    required MessageRole role,
    required String content,
    String? reasoningContent,
    List<ToolCall>? toolCalls,
    String? toolCallId,
    bool stoppedByUser = false,
  }) => Message(
    role: role,
    content: content,
    reasoningContent: reasoningContent,
    toolCalls: toolCalls,
    toolCallId: toolCallId,
    stoppedByUser: stoppedByUser,
  );

  /// 创建包含文件的消息
  factory Message.withFiles({
    required MessageRole role,
    String content = '',
    required List<FileContent> files,
    String? reasoningContent,
    List<ToolCall>? toolCalls,
    String? toolCallId,
    bool stoppedByUser = false,
  }) => Message(
    role: role,
    content: content,
    files: files,
    reasoningContent: reasoningContent,
    toolCalls: toolCalls,
    toolCallId: toolCallId,
    stoppedByUser: stoppedByUser,
  );

  /// 判断消息是否包含文件
  bool get hasFiles => files != null && files!.isNotEmpty;

  /// 判断消息是否只包含文件（没有文本内容）
  bool get isFilesOnly => hasFiles && content.isEmpty;

  factory Message.fromJson(Map<String, dynamic> json) =>
      _$MessageFromJson(json);
  Map<String, dynamic> toJson() => _$MessageToJson(this);

  Message copyWith({
    MessageRole? role,
    String? content,
    String? reasoningContent,
    List<ToolCall>? toolCalls,
    String? toolCallId,
    List<FileContent>? files,
    bool? stoppedByUser,
  }) => Message(
    role: role ?? this.role,
    content: content ?? this.content,
    reasoningContent: reasoningContent ?? this.reasoningContent,
    toolCalls: toolCalls ?? this.toolCalls,
    toolCallId: toolCallId ?? this.toolCallId,
    files: files ?? this.files,
    stoppedByUser: stoppedByUser ?? this.stoppedByUser,
  );
}
