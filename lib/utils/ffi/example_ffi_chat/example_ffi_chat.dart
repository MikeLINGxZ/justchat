import 'dart:ffi';
import 'dart:convert';
import 'package:ffi/ffi.dart';
import 'package:lemon_tea/utils/ffi/example_ffi_chat/example_ffi_chat.g.dart';

/// 消息类，对应Go中的Message结构体
class Message {
  final String role;
  final String content;

  Message({required this.role, required this.content});

  Map<String, dynamic> toJson() => {'role': role, 'content': content};

  factory Message.fromJson(Map<String, dynamic> json) =>
      Message(role: json['role'] as String, content: json['content'] as String);
}

/// 聊天请求类，对应Go中的ChatRequest结构体
class ChatRequest {
  final String systemPrompt;
  final List<Message> messages;
  final String apiKey;
  final String baseURL;
  final String model;

  ChatRequest({
    required this.systemPrompt,
    required this.messages,
    required this.apiKey,
    required this.baseURL,
    required this.model,
  });

  Map<String, dynamic> toJson() => {
    'system_prompt': systemPrompt,
    'messages': messages.map((e) => e.toJson()).toList(),
    'api_key': apiKey,
    'base_url': baseURL,
    'model': model,
  };
}

/// 聊天响应类，对应Go中的ChatResponse结构体
class ChatResponse {
  final String content;
  final String? error;

  ChatResponse({required this.content, this.error});

  factory ChatResponse.fromJson(Map<String, dynamic> json) => ChatResponse(
    content: json['content'] as String,
    error: json['error'] as String?,
  );
}

class ExampleFfiChat {
  static ExampleFfiChat? _instance;
  static ExampleFfiChatGenerate? _generateInstance;

  static String libPath() {
    return "example_ffi_chat_arm64.dylib";
  }

  static final DynamicLibrary _dylib = DynamicLibrary.open(libPath());

  factory ExampleFfiChat() {
    _instance ??= ExampleFfiChat._internal();
    return _instance!;
  }

  ExampleFfiChat._internal();

  static ExampleFfiChatGenerate instance() {
    _generateInstance ??= ExampleFfiChatGenerate(_dylib);
    return _generateInstance!;
  }

  /// 封装Chat方法，处理JSON字符串转换
  ///
  /// 参数:
  /// - [request]: ChatRequest对象，包含聊天所需的所有参数
  ///
  /// 返回:
  /// - ChatResponse对象，包含聊天结果或错误信息
  static ChatResponse chat(ChatRequest request) {
    // 获取FFI实例
    final ffiInstance = instance();

    // 将请求对象转换为JSON字符串
    final jsonInput = jsonEncode(request.toJson());

    // 将Dart字符串转换为C字符串
    final inputPtr = jsonInput.toNativeUtf8().cast<Char>();

    // 调用本地函数
    final resultPtr = ffiInstance.Chat(inputPtr);

    // 将结果转换回Dart字符串
    String resultJson = "";
    if (resultPtr != nullptr) {
      resultJson = resultPtr.cast<Utf8>().toDartString();
    }

    // 释放内存
    calloc.free(inputPtr);

    // 解析JSON响应
    final Map<String, dynamic> responseMap = jsonDecode(resultJson);
    return ChatResponse.fromJson(responseMap);
  }
}
