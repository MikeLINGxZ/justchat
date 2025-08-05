import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/chat_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_toolbar.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/controls/resizable_divider.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/utils/message_converter.dart';
import 'package:lemon_tea/rpc/service.pb.dart' as grpc_service;
import 'package:lemon_tea/rpc/common.pb.dart' as grpc_common;
import 'package:lemon_tea/storage/chat_storage.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'dart:async';
import 'dart:io' show Platform;

import 'package:lemon_tea/utils/style.dart';

class AssistantPage extends StatefulWidget {
  final ConversationManager? conversationManager;

  const AssistantPage({super.key, this.conversationManager});

  @override
  State<StatefulWidget> createState() => _AssistantPage();
}

class _AssistantPage extends State<AssistantPage> {
  late ConversationManager _conversationManager;
  List<Message> _historyMessages = [];
  bool _isLoading = false;
  String _currentTitle = 'AI 助手';
  Conversation? _currentConversation;
  final Client _grpcClient = Client();
  String? _selectedProviderId;
  String? _selectedModelId;
  bool _isStreaming = false; // 添加流式状态标记
  final GlobalKey<ChatViewState> _chatViewKey = GlobalKey<ChatViewState>(); // 添加ChatView的key
  StreamSubscription<grpc_service.ChatResponse>? _currentStreamSubscription; // 当前流订阅

  @override
  void initState() {
    super.initState();
    _conversationManager = widget.conversationManager ?? ConversationManager();
    // 监听 ConversationManager 的变化
    _conversationManager.addListener(_onConversationManagerChanged);
    _initializeConversation();
  }

  // 添加自动聚焦到输入框的方法
  void _focusInputField() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _chatViewKey.currentState?.requestInputFocus();
    });
  }

  @override
  void dispose() {
    _conversationManager.removeListener(_onConversationManagerChanged);
    _currentStreamSubscription?.cancel(); // 取消当前流订阅
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，更新UI
    if (mounted) {
      setState(() {
        _historyMessages = _conversationManager.messages;
        _currentTitle =
            _conversationManager.currentConversation?.title ?? 'AI 助手';
        _currentConversation = _conversationManager.currentConversation;

        // 同步模型配置
        if (_currentConversation != null) {
          _selectedProviderId =
              _currentConversation!.defaultProviderId ?? 'deepseek';
          _selectedModelId =
              _currentConversation!.defaultModelId ?? 'deepseek-chat';
        }
      });
    }
  }

  Future<void> _initializeConversation() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // 初始化gRPC客户端
      await _grpcClient.init();

      // 确保有基本的LLM配置
      await _ensureLlmConfiguration();

      // 检查ConversationManager是否已经有当前对话
      if (_conversationManager.currentConversation != null) {
        // 如果已经有当前对话，直接使用它
        await _loadConversation(_conversationManager.currentConversation!);
        debugPrint(
          '使用已有的当前对话: ${_conversationManager.currentConversation!.title}',
        );
      } else {
        // 获取所有对话
        final conversations = await ChatStorage.getAllConversations();

        if (conversations.isEmpty) {
          // 如果没有对话，创建一个新的欢迎对话
          await _createWelcomeConversation();
        } else {
          // 加载最新的对话
          await _loadConversation(conversations.first);
        }
      }
    } catch (e) {
      debugPrint('Failed to initialize conversation: $e');

      // 如果初始化失败，尝试创建一个离线模式的对话记录
      try {
        debugPrint('尝试创建离线模式的对话记录...');
        await _createWelcomeConversation();
        debugPrint('离线模式对话创建成功');
      } catch (offlineError) {
        debugPrint('离线模式对话创建也失败: $offlineError');

        // 最后的备选方案：创建一个临时的本地对话记录和消息
        // 注意：这种情况下消息不会被保存到数据库，但至少应用不会崩溃
        debugPrint('使用临时本地模式');

        _currentConversation = Conversation(
          id: 'temp-conversation-${DateTime.now().millisecondsSinceEpoch}',
          title: 'AI 助手 (离线模式)',
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
          defaultProviderId: 'deepseek',
          defaultModelId: 'deepseek-chat',
        );

        final welcomeContent = """# AI 助手 (离线模式)

数据库连接失败，当前运行在离线模式下。
您的对话记录可能无法正常保存，请检查应用权限或重启应用。

## 可用功能
- 基础聊天功能
- Markdown 显示

## 注意事项
- 消息可能无法保存
- 重启应用后对话记录可能丢失""";

        _historyMessages = [
          Message(role: MessageRole.assistant, content: welcomeContent),
        ];

        setState(() {
          _currentTitle = _currentConversation!.title;
          _selectedProviderId = 'deepseek';
          _selectedModelId = 'deepseek-chat';
        });
      }
    } finally {
      setState(() {
        _isLoading = false;
      });
      
      // 初始化完成后自动聚焦到输入框
      _focusInputField();
    }
  }

  /// 确保有基本的LLM配置
  Future<void> _ensureLlmConfiguration() async {
    try {
      // 检查是否已有LLM提供商配置
      final providers = await LlmStorage.getAllProviders();

      if (providers.isEmpty) {
        debugPrint('未找到LLM提供商配置，创建默认配置');

        // 创建默认的DeepSeek提供商
        final defaultProvider = LlmProvider(
          id: 'deepseek',
          name: 'DeepSeek',
          baseUrl: 'https://api.deepseek.com/v1',
          apiKey: Platform.environment['DEEPSEEK_API_KEY'] ?? '',
          alias: 'DeepSeek AI',
          description: 'DeepSeek AI 服务',
          seqId: 1,
        );

        await LlmStorage.addProvider(defaultProvider);

        // 创建默认模型
        final defaultModel = Model(
          id: 'deepseek-chat',
          object: 'model',
          ownedBy: 'deepseek',
          enabled: true,
          llmProviderId: 'deepseek',
          seqId: 1,
        );

        await LlmStorage.addModel(defaultModel);

        debugPrint('默认LLM配置创建完成');
      } else {
        debugPrint('找到${providers.length}个LLM提供商配置');
      }
    } catch (e) {
      debugPrint('配置LLM提供商失败: $e');
    }
  }

  Future<void> _loadConversation(Conversation conversation) async {
    _currentConversation = conversation;
    final messages = await ChatStorage.getLlmMessagesByConversationId(
      conversation.id,
    );

    setState(() {
      _currentTitle = conversation.title;
      _historyMessages = messages;
      _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
      _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
    });
  }

  Future<void> _createWelcomeConversation() async {
    try {
      const welcomeContent = """# 欢迎使用 Markdown

这是一个简单的 Markdown 示例文档，展示常用语法：

## 标题层级
二级标题 (`##`) 到六级标题 (`######`)

## 文字样式
- **加粗文本** (`**加粗**`)
- *斜体文本* (`*斜体*`)
- ~~删除线~~ (`~~删除线~~`)
- `行内代码` (`` `行内代码` ``)

## 列表
### 无序列表
- 项目一
- 项目二
  - 子项目 (缩进两个空格)

### 有序列表
1. 第一项
2. 第二项
   1. 子项 (缩进三个空格)

## 链接与图片
[百度链接](https://www.baidu.com)  
![示例图片](https://via.placeholder.com/150 "悬浮提示")

## 代码块
```python
def hello():
    print("代码高亮示例")""";

      // 通过 ConversationManager 创建欢迎对话，这样会自动更新对话列表并通知监听者
      final conversation = await _conversationManager.createConversation(
        title: '欢迎对话',
        initialMessages: [Message(role: MessageRole.assistant, content: welcomeContent)],
        defaultProviderId: 'deepseek',
        defaultModelId: 'deepseek-chat',
      );

      if (conversation != null) {
        _currentConversation = conversation;
        debugPrint('欢迎对话创建成功: ${conversation.id}');

        setState(() {
          _currentTitle = conversation.title;
          _historyMessages = [
            Message(role: MessageRole.assistant, content: welcomeContent),
          ];
          _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
          _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
        });
      } else {
        debugPrint('错误：创建欢迎对话失败');
        throw Exception('创建欢迎对话失败');
      }
    } catch (e) {
      debugPrint('创建欢迎对话时发生异常: $e');
      rethrow; // 重新抛出异常，让调用方处理
    }
  }

  Future<void> _handleSendMessageWithFiles(String message, List<FileContent> files) async {
    await _handleSendMessage(message, files);
  }

  // 停止生成方法
  void _handleStopGeneration() {
    debugPrint('用户请求停止生成');
    
    // 取消当前流订阅
    _currentStreamSubscription?.cancel();
    _currentStreamSubscription = null;
    
    // 标记最后一条消息为用户停止
    if (_historyMessages.isNotEmpty && 
        _historyMessages.last.role == MessageRole.assistant) {
      setState(() {
        final lastMessage = _historyMessages.last;
        _historyMessages[_historyMessages.length - 1] = lastMessage.copyWith(
          stoppedByUser: true,
        );
        _isStreaming = false;
      });
    } else {
      // 更新UI状态
      if (mounted) {
        setState(() {
          _isStreaming = false;
        });
      }
    }
    
    debugPrint('生成已停止');
  }



  Future<void> _handleSendMessage(String message, [List<FileContent>? files]) async {
    files ??= [];
    
    if (message.trim().isEmpty && files.isEmpty) return;
    if (_currentConversation == null) return;

    // 检查是否为临时对话模式
    final isTemporaryConversation = _currentConversation!.id.startsWith(
      'temp-',
    );
    if (isTemporaryConversation) {
      debugPrint('警告：当前为临时对话模式，消息将不会被保存到数据库');
    }

    // 添加用户消息到界面（包含文件）
    final userMessage = files.isNotEmpty 
        ? Message.withFiles(
            role: MessageRole.user, 
            content: message,
            files: files,
          )
        : Message(role: MessageRole.user, content: message);
    setState(() {
      _historyMessages.add(userMessage);
    });

    // 保存用户消息到数据库，增加错误处理
    if (!isTemporaryConversation) {
      try {
        final savedUserMessage = await ChatStorage.addMessage(
          conversationId: _currentConversation!.id,
          role: 'user',
          content: message,
          files: files,
        );

        if (savedUserMessage == null) {
          debugPrint('警告：用户消息保存失败');
          // 可以在这里添加UI提示，比如显示一个警告图标
        } else {
          debugPrint('用户消息保存成功: ${savedUserMessage.id}');
        }
      } catch (e) {
        debugPrint('保存用户消息时发生异常: $e');
        // 保存失败不应该阻止继续发送请求
      }
    } else {
      debugPrint('跳过用户消息保存（临时对话模式）');
    }

    // 如果是第一条用户消息，更新对话标题
    if (!isTemporaryConversation) {
      try {
        final messageCount = await ChatStorage.getMessageCountByConversationId(
          _currentConversation!.id,
        );
        if (messageCount >= 2) {
          // 可能包含欢迎消息 + 用户消息
          final title = _generateTitleFromMessage(message);
          final titleUpdated = await ChatStorage.updateConversationTitle(
            _currentConversation!.id,
            title,
          );
          if (titleUpdated) {
            setState(() {
              _currentTitle = title;
            });
            debugPrint('对话标题更新成功: $title');
          } else {
            debugPrint('对话标题更新失败');
          }
        }
      } catch (e) {
        debugPrint('更新对话标题时发生异常: $e');
      }
    } else {
      // 临时对话模式下也可以更新标题，只是不保存到数据库
      if (_historyMessages.length == 1) {
        // 只有欢迎消息时
        final title = _generateTitleFromMessage(message);
        setState(() {
          _currentTitle = title;
        });
        debugPrint('更新临时对话标题: $title');
      }
    }

    // 立即添加一个空的AI消息用于流式显示
    final aiMessage = Message(role: MessageRole.assistant, content: '');
    setState(() {
      _historyMessages.add(aiMessage);
      _isStreaming = true; // 开始流式显示
    });

    String fullResponse = '';
    String fullReasoningContent = ''; // 添加思考过程内容
    bool hasError = false;
    String? aiMessageId; // 用于记录AI消息的ID

    try {
      // 检查gRPC客户端
      if (_grpcClient.stub == null) {
        await _grpcClient.init();
        if (_grpcClient.stub == null) {
          throw Exception("gRPC客户端初始化失败，请检查服务是否启动");
        }
      }

      // 准备历史消息（排除刚添加的空AI消息和当前用户消息）
      final historyMessagesForRequest = _historyMessages
          .where((msg) => 
              (msg.content.isNotEmpty || (msg.files?.isNotEmpty ?? false)) &&
              msg != _historyMessages.last && // 排除空的AI消息
              msg != _historyMessages[_historyMessages.length - 2] // 排除刚添加的用户消息
          )
          .toList();
      
      List<grpc_common.Message> grpcMessages = MessageConverter.convertMessages(historyMessagesForRequest);

      // 创建聊天请求
      final providerId = _currentConversation!.defaultProviderId ?? 'deepseek';
      final modelId = _currentConversation!.defaultModelId ?? 'deepseek-chat';

      debugPrint('使用模型配置: providerId=$providerId, modelId=$modelId');

      // 创建当前用户消息
      final currentMessage = files.isNotEmpty 
          ? Message.withFiles(
              role: MessageRole.user, 
              content: message,
              files: files,
            )
          : Message(role: MessageRole.user, content: message);

      final chatRequest = grpc_service.ChatRequest(
        llmProviderId: providerId,
        modelId: modelId,
        historyMessages: grpcMessages,
        messages: [MessageConverter.convertMessage(currentMessage)],
        prompt: "你是一个有用的AI助手。",
      );

      // 创建请求流
      final controller = StreamController<grpc_service.ChatRequest>();
      controller.add(chatRequest);
      controller.close();

      // 调用gRPC聊天接口
      final responseStream = _grpcClient.stub!.chat(controller.stream);

      // 使用StreamSubscription来处理响应，以便可以取消
      final completer = Completer<void>();
      _currentStreamSubscription = responseStream.listen(
        (response) async {
          if (response.errorMessage.isNotEmpty) {
            completer.completeError(Exception(response.errorMessage));
            return;
          }

          fullResponse += response.content;
          fullReasoningContent += response.reasoningContent; // 累积思考过程内容

          // 实时更新AI消息内容
          if (mounted) {
            setState(() {
              _historyMessages.last = Message(
                role: MessageRole.assistant,
                content: fullResponse,
                reasoningContent:
                    fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
              );
            });
          }

          if (response.isDone) {
            // 聊天完成，标记完成状态
            completer.complete();
          }
        },
        onError: (error) {
          completer.completeError(error);
        },
        onDone: () {
          _currentStreamSubscription = null;
          if (mounted) {
            setState(() {
              _isStreaming = false;
            });
          }
          if (!completer.isCompleted) {
            completer.complete();
          }
        },
      );

      // 等待流完成或被取消
      await completer.future;
    } catch (e) {
      hasError = true;
      debugPrint('Error during chat: $e');

      // 发生错误时，更新AI消息为错误内容
      final errorContent = '抱歉，处理您的请求时发生错误：${e.toString()}';

      if (mounted) {
        setState(() {
          _historyMessages.last = Message(
            role: MessageRole.assistant,
            content: errorContent,
          );
        });
      }

      // 保存错误消息到数据库
      if (!isTemporaryConversation) {
        try {
          final savedErrorMessage = await ChatStorage.addMessage(
            conversationId: _currentConversation!.id,
            role: 'assistant',
            content: errorContent,
            reasoningContent: null,
          );

          if (savedErrorMessage != null) {
            debugPrint('错误消息保存成功: ${savedErrorMessage.id}');
          } else {
            debugPrint('警告：错误消息保存失败');
          }
        } catch (dbError) {
          debugPrint('保存错误消息到数据库失败: $dbError');
        }
      } else {
        debugPrint('跳过错误消息保存（临时对话模式）');
      }
    } finally {
      // 结束流式显示
      if (mounted) {
        setState(() {
          _isStreaming = false;
        });
      }

      // 保存AI回复消息到数据库（统一保存逻辑）
      if (!isTemporaryConversation &&
          fullResponse.isNotEmpty &&
          aiMessageId == null &&
          !hasError) {
        try {
          // 检查最后一条消息是否被用户停止
          bool wasStoppedByUser = false;
          if (_historyMessages.isNotEmpty && 
              _historyMessages.last.role == MessageRole.assistant) {
            wasStoppedByUser = _historyMessages.last.stoppedByUser;
          }
          
          final savedAiMessage = await ChatStorage.addMessage(
            conversationId: _currentConversation!.id,
            role: 'assistant',
            content: fullResponse,
            reasoningContent:
                fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
            stoppedByUser: wasStoppedByUser,
          );

          if (savedAiMessage != null) {
            aiMessageId = savedAiMessage.id;
            final stopInfo = wasStoppedByUser ? '（用户停止）' : '';
            debugPrint('AI回复保存成功$stopInfo: ${savedAiMessage.id}');
          } else {
            debugPrint('警告：AI回复保存失败');
          }
        } catch (e) {
          debugPrint('保存AI回复时发生异常: $e');
        }
      } else if (!isTemporaryConversation && fullResponse.isEmpty && !hasError) {
        debugPrint('跳过AI回复保存（无内容）');
      } else if (isTemporaryConversation) {
        debugPrint('跳过AI回复保存（临时对话模式）');
      }

      // 确保在异常情况下，如果AI消息仍然为空，则移除它
      if (mounted &&
          _historyMessages.isNotEmpty &&
          _historyMessages.last.content.isEmpty &&
          !hasError) {
        setState(() {
          _historyMessages.removeLast();
        });
        debugPrint('移除了空的AI消息');
      }

      // 如果成功保存了消息，记录一下
      if (aiMessageId != null) {
        debugPrint('对话记录保存完成 - 用户消息和AI回复(ID: $aiMessageId)');
      } else if (fullResponse.isNotEmpty && !isTemporaryConversation) {
        debugPrint('警告：AI回复有内容但未能保存到数据库');
      }
    }
  }

  Future<void> _handleNewConversation() async {
    // 通过 ConversationManager 创建新对话，这样会自动更新对话列表并通知监听者
    final conversation = await _conversationManager.createConversation(
      title: '新对话',
      defaultProviderId: _selectedProviderId ?? 'deepseek',
      defaultModelId: _selectedModelId ?? 'deepseek-chat',
    );

    if (conversation != null) {
      // ConversationManager 已经设置了当前对话，这里只需要更新UI状态
      setState(() {
        _currentConversation = conversation;
        _currentTitle = conversation.title;
        _historyMessages = [];
        _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
        _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
      });
      
      // 新对话创建完成后自动聚焦到输入框
      _focusInputField();
    }
  }

  Future<void> _handleModelSelected(String providerId, String modelId) async {
    setState(() {
      _selectedProviderId = providerId;
      _selectedModelId = modelId;
    });

    // 如果有当前对话，更新对话的默认模型配置
    if (_currentConversation != null) {
      try {
        await ChatStorage.updateConversationModel(
          _currentConversation!.id,
          providerId,
          modelId,
        );

        // 更新本地对话对象
        _currentConversation = _currentConversation!.copyWith(
          defaultProviderId: providerId,
          defaultModelId: modelId,
        );

        debugPrint('已更新对话模型配置: providerId=$providerId, modelId=$modelId');
      } catch (e) {
        debugPrint('更新对话模型配置失败: $e');
      }
    }
  }

  String _generateTitleFromMessage(String message) {
    // 简单的标题生成逻辑
    if (message.length <= 20) {
      return message;
    }
    return '${message.substring(0, 17)}...';
  }



  // MessageToolbar 回调函数实现
  void _handleCopyMessage(Message message) {
    // 复制消息内容（包含Markdown格式）
    Clipboard.setData(ClipboardData(text: message.content));
    
    // 显示复制成功提示
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('已复制到剪贴板'),
        duration: Duration(seconds: 2),
      ),
    );
  }

  void _handleCopyPlainText(Message message) {
    // 复制纯文本内容（去除Markdown格式）
    String plainText = _stripMarkdown(message.content);
    Clipboard.setData(ClipboardData(text: plainText));
    
    // 显示复制成功提示
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('已复制纯文本到剪贴板'),
        duration: Duration(seconds: 2),
      ),
    );
  }

  Future<void> _handleRegenerateMessage(Message message) async {
    // 检查消息角色，只有AI消息才能重新生成
    if (message.role != MessageRole.assistant) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('只能重新生成AI回复消息'),
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    // 找到消息在历史记录中的位置
    final messageIndex = _historyMessages.indexOf(message);
    if (messageIndex == -1) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('未找到要重新生成的消息'),
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    // 找到对应的用户消息（通常是前一条消息）
    String? userMessageContent;
    List<FileContent>? userMessageFiles;
    for (int i = messageIndex - 1; i >= 0; i--) {
      if (_historyMessages[i].role == MessageRole.user) {
        userMessageContent = _historyMessages[i].content;
        userMessageFiles = _historyMessages[i].files;
        break;
      }
    }

    if (userMessageContent == null && (userMessageFiles == null || userMessageFiles.isEmpty)) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('未找到对应的用户消息'),
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    // 显示确认对话框
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('重新生成'),
        content: const Text('确定要重新生成这条AI回复吗？原回复将被替换。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('重新生成'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    // 记录原始消息内容，用于后续的数据库更新
    final originalContent = message.content;

    // 将当前位置的消息设置为空的AI消息，准备重新生成
    setState(() {
      _historyMessages[messageIndex] = Message(
        role: MessageRole.assistant, 
        content: '',
      );
    });

    // 调用重新生成方法，传入用户消息内容、要替换的索引和原始内容
    await _regenerateAIResponse(userMessageContent ?? '', messageIndex, originalContent, userMessageFiles);
  }

  Future<void> _handleDeleteMessage(Message message) async {
    // 显示确认对话框
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除消息'),
        content: const Text('确定要删除这条消息吗？此操作无法撤销。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('删除'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    // 找到消息在历史记录中的位置
    final messageIndex = _historyMessages.indexOf(message);
    if (messageIndex == -1) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('未找到要删除的消息'),
          duration: Duration(seconds: 2),
        ),
      );
      return;
    }

    // 从界面中移除消息
    setState(() {
      _historyMessages.removeAt(messageIndex);
    });

    // 从数据库中删除消息
    if (_currentConversation != null && !_currentConversation!.id.startsWith('temp-')) {
      try {
        await ChatStorage.deleteMessagesByContent(
          _currentConversation!.id,
          message.content,
          message.role.name,
        );
        
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('消息已删除'),
            duration: Duration(seconds: 2),
          ),
        );
      } catch (e) {
        debugPrint('删除消息时出错: $e');
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('删除消息失败'),
            duration: Duration(seconds: 2),
          ),
        );
      }
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('消息已删除（临时对话模式）'),
          duration: Duration(seconds: 2),
        ),
      );
    }
  }

  // 重新生成AI回复的方法（不添加用户消息）
  Future<void> _regenerateAIResponse(String userMessageContent, [int? replaceIndex, String? originalContent, List<FileContent>? files]) async {
    if (_currentConversation == null) return;

    // 检查是否为临时对话模式
    final isTemporaryConversation = _currentConversation!.id.startsWith('temp-');

    // 开始流式显示（目标位置的空消息已经在调用前设置好了）
    setState(() {
      _isStreaming = true;
    });

    String fullResponse = '';
    String fullReasoningContent = '';
    bool hasError = false;
    String? aiMessageId;

    try {
      // 检查gRPC客户端
      if (_grpcClient.stub == null) {
        await _grpcClient.init();
        if (_grpcClient.stub == null) {
          throw Exception("gRPC客户端初始化失败，请检查服务是否启动");
        }
      }

      // 准备历史消息（如果指定了replaceIndex，则只使用该索引之前的消息）
      List<Message> historyToUse;
      if (replaceIndex != null) {
        // 使用重新生成位置之前的历史消息（包含对应的用户消息，不包含要重新生成的AI消息）
        historyToUse = _historyMessages.sublist(0, replaceIndex);
        
        // 如果有文件信息，需要确保最后一条用户消息包含文件
        if (files != null && files.isNotEmpty && historyToUse.isNotEmpty) {
          final lastUserMessageIndex = historyToUse.lastIndexWhere((msg) => msg.role == MessageRole.user);
          if (lastUserMessageIndex != -1) {
            final lastUserMessage = historyToUse[lastUserMessageIndex];
            // 创建包含文件的新用户消息
            final updatedUserMessage = Message.withFiles(
              role: MessageRole.user,
              content: lastUserMessage.content,
              files: files,
            );
            historyToUse[lastUserMessageIndex] = updatedUserMessage;
          }
        }
      } else {
        // 使用所有非空消息（排除最后一个空AI消息）
        historyToUse = _historyMessages.where((msg) => msg.content.isNotEmpty || (msg.files?.isNotEmpty ?? false)).toList();
      }
      
      List<grpc_common.Message> grpcMessages = MessageConverter.convertMessages(historyToUse);

      // 创建聊天请求
      final providerId = _currentConversation!.defaultProviderId ?? 'deepseek';
      final modelId = _currentConversation!.defaultModelId ?? 'deepseek-chat';

      debugPrint('重新生成使用模型配置: providerId=$providerId, modelId=$modelId');

      final chatRequest = grpc_service.ChatRequest(
        llmProviderId: providerId,
        modelId: modelId,
        historyMessages: grpcMessages,
        messages: [], // 重新生成时不需要新的用户消息
        prompt: "你是一个有用的AI助手。",
      );

      // 创建请求流
      final controller = StreamController<grpc_service.ChatRequest>();
      controller.add(chatRequest);
      controller.close();

      // 调用gRPC聊天接口
      final responseStream = _grpcClient.stub!.chat(controller.stream);

      // 使用StreamSubscription来处理响应，以便可以取消
      final completer = Completer<void>();
      _currentStreamSubscription = responseStream.listen(
        (response) async {
          if (response.errorMessage.isNotEmpty) {
            completer.completeError(Exception(response.errorMessage));
            return;
          }

          fullResponse += response.content;
          fullReasoningContent += response.reasoningContent;

          // 实时更新AI消息内容
          if (mounted) {
            setState(() {
              final newMessage = Message(
                role: MessageRole.assistant,
                content: fullResponse,
                reasoningContent:
                    fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
              );
              
              if (replaceIndex != null) {
                _historyMessages[replaceIndex] = newMessage;
              } else {
                _historyMessages.last = newMessage;
              }
            });
          }

          if (response.isDone) {
            // 聊天完成，更新AI回复到数据库（如果有原始内容）或保存新的AI回复
            if (!isTemporaryConversation) {
              try {
                if (originalContent != null) {
                  // 更新现有消息而不是创建新消息
                  final updateResult = await ChatStorage.updateMessageByContent(
                    _currentConversation!.id,
                    originalContent,
                    'assistant',
                    fullResponse,
                    fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
                  );

                  if (updateResult) {
                    debugPrint('重新生成的AI回复更新成功');
                  } else {
                    debugPrint('警告：重新生成的AI回复更新失败');
                  }
                } else {
                  // 原始逻辑：添加新消息
                  final savedAiMessage = await ChatStorage.addMessage(
                    conversationId: _currentConversation!.id,
                    role: 'assistant',
                    content: fullResponse,
                    reasoningContent:
                        fullReasoningContent.isNotEmpty
                            ? fullReasoningContent
                            : null,
                  );

                  if (savedAiMessage != null) {
                    aiMessageId = savedAiMessage.id;
                    debugPrint('重新生成的AI回复保存成功: ${savedAiMessage.id}');
                  } else {
                    debugPrint('警告：重新生成的AI回复保存失败');
                  }
                }
              } catch (e) {
                debugPrint('保存/更新重新生成的AI回复时发生异常: $e');
              }
            } else {
              debugPrint('跳过重新生成的AI回复保存（临时对话模式）');
            }
            completer.complete();
          }
        },
        onError: (error) {
          completer.completeError(error);
        },
        onDone: () {
          _currentStreamSubscription = null;
          if (mounted) {
            setState(() {
              _isStreaming = false;
            });
          }
          if (!completer.isCompleted) {
            completer.complete();
          }
        },
      );

      // 等待流完成或被取消
      await completer.future;
    } catch (e) {
      hasError = true;
      debugPrint('重新生成时发生错误: $e');

      // 发生错误时，更新AI消息为错误内容
      final errorContent = '抱歉，重新生成时发生错误：${e.toString()}';

      if (mounted) {
        setState(() {
          final errorMessage = Message(
            role: MessageRole.assistant,
            content: errorContent,
          );
          
          if (replaceIndex != null) {
            _historyMessages[replaceIndex] = errorMessage;
          } else {
            _historyMessages.last = errorMessage;
          }
        });
      }

      // 保存错误消息到数据库
      if (!isTemporaryConversation) {
        try {
          if (originalContent != null) {
            // 更新现有消息而不是创建新消息
            final updateResult = await ChatStorage.updateMessageByContent(
              _currentConversation!.id,
              originalContent,
              'assistant',
              errorContent,
              null,
            );

            if (updateResult) {
              debugPrint('重新生成的错误消息更新成功');
            } else {
              debugPrint('警告：重新生成的错误消息更新失败');
            }
          } else {
            // 原始逻辑：添加新消息
            final savedErrorMessage = await ChatStorage.addMessage(
              conversationId: _currentConversation!.id,
              role: 'assistant',
              content: errorContent,
              reasoningContent: null,
            );

            if (savedErrorMessage != null) {
              debugPrint('重新生成的错误消息保存成功: ${savedErrorMessage.id}');
            } else {
              debugPrint('警告：重新生成的错误消息保存失败');
            }
          }
        } catch (dbError) {
          debugPrint('保存/更新重新生成的错误消息到数据库失败: $dbError');
        }
      } else {
        debugPrint('跳过重新生成的错误消息保存（临时对话模式）');
      }
    } finally {
      // 结束流式显示
      if (mounted) {
        setState(() {
          _isStreaming = false;
        });
      }

      // 确保AI回复消息被保存
      if (!isTemporaryConversation &&
          fullResponse.isNotEmpty &&
          aiMessageId == null &&
          !hasError) {
        try {
          if (originalContent != null) {
            // 更新现有消息而不是创建新消息
            final updateResult = await ChatStorage.updateMessageByContent(
              _currentConversation!.id,
              originalContent,
              'assistant',
              fullResponse,
              fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
            );

            if (updateResult) {
              debugPrint('重新生成的AI回复补充更新成功');
            } else {
              debugPrint('警告：重新生成的AI回复补充更新失败');
            }
          } else {
            // 原始逻辑：添加新消息
            final savedAiMessage = await ChatStorage.addMessage(
              conversationId: _currentConversation!.id,
              role: 'assistant',
              content: fullResponse,
              reasoningContent:
                  fullReasoningContent.isNotEmpty ? fullReasoningContent : null,
            );

            if (savedAiMessage != null) {
              aiMessageId = savedAiMessage.id;
              debugPrint('重新生成的AI回复补充保存成功: ${savedAiMessage.id}');
            } else {
              debugPrint('警告：重新生成的AI回复补充保存失败');
            }
          }
        } catch (e) {
          debugPrint('补充保存/更新重新生成的AI回复时发生异常: $e');
        }
      }

      // 确保在异常情况下，如果AI消息仍然为空，则移除它
      if (mounted &&
          _historyMessages.isNotEmpty &&
          !hasError) {
        final targetIndex = replaceIndex ?? (_historyMessages.length - 1);
        if (targetIndex < _historyMessages.length && 
            _historyMessages[targetIndex].content.isEmpty) {
          setState(() {
            _historyMessages.removeAt(targetIndex);
          });
          debugPrint('移除了空的重新生成AI消息');
        }
      }

      if (originalContent != null) {
        debugPrint('重新生成完成 - AI回复已更新');
      } else if (aiMessageId != null) {
        debugPrint('重新生成完成 - AI回复(ID: $aiMessageId)');
      } else if (fullResponse.isNotEmpty && !isTemporaryConversation) {
        debugPrint('警告：重新生成的AI回复有内容但未能保存到数据库');
      }
    }
  }

  // 辅助函数：去除Markdown格式
  String _stripMarkdown(String markdown) {
    String text = markdown;
    
    // 移除代码块
    text = text.replaceAll(RegExp(r'```[\s\S]*?```'), '');
    text = text.replaceAll(RegExp(r'`[^`]*`'), '');
    
    // 移除标题标记
    text = text.replaceAll(RegExp(r'^#{1,6}\s+', multiLine: true), '');
    
    // 移除粗体和斜体
    text = text.replaceAll(RegExp(r'\*\*([^*]+)\*\*'), r'$1');
    text = text.replaceAll(RegExp(r'\*([^*]+)\*'), r'$1');
    text = text.replaceAll(RegExp(r'__([^_]+)__'), r'$1');
    text = text.replaceAll(RegExp(r'_([^_]+)_'), r'$1');
    
    // 移除删除线
    text = text.replaceAll(RegExp(r'~~([^~]+)~~'), r'$1');
    
    // 移除链接
    text = text.replaceAll(RegExp(r'\[([^\]]+)\]\([^)]+\)'), r'$1');
    
    // 移除图片
    text = text.replaceAll(RegExp(r'!\[([^\]]*)\]\([^)]+\)'), r'$1');
    
    // 移除列表标记
    text = text.replaceAll(RegExp(r'^\s*[-*+]\s+', multiLine: true), '');
    text = text.replaceAll(RegExp(r'^\s*\d+\.\s+', multiLine: true), '');
    
    // 清理多余的空行
    text = text.replaceAll(RegExp(r'\n\s*\n'), '\n\n');
    text = text.trim();
    
    return text;
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    final sideWidget = _buildSide();

    // 如果侧边栏为空，只显示聊天界面
    if (sideWidget == null) {
      return Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: double.infinity),
          child:         ChatView(
          key: _chatViewKey, // 添加key
          historyMessages: _historyMessages,
          onSendWithFiles: _handleSendMessageWithFiles,
          onNewConversation:
              _historyMessages.isEmpty ? null : _handleNewConversation,
          currentTitle: _currentTitle,
          selectedProviderId: _selectedProviderId,
          selectedModelId: _selectedModelId,
          onModelSelected: _handleModelSelected,
          isStreaming: _isStreaming,
          visibleWidth: Style.messageViewWidth,
          onStopGeneration: _handleStopGeneration,
          messageToolBar: MessageToolbar(
            onCopy: _handleCopyMessage,
            onCopyPlainText: _handleCopyPlainText,
            onRegenerate: _handleRegenerateMessage,
            onDelete: _handleDeleteMessage,
            isVisible: !_isStreaming, // 流式生成过程中隐藏工具栏
          ),
        ),
        ),
      );
    }

    // 如果侧边栏不为空，显示分割布局
    return ResizableDivider(
      leftChild: ChatView(
        // 移除重复的key，避免GlobalKey冲突
        historyMessages: _historyMessages,
        onSendWithFiles: _handleSendMessageWithFiles,
        onNewConversation:
            _historyMessages.isEmpty ? null : _handleNewConversation,
        currentTitle: _currentTitle,
        selectedProviderId: _selectedProviderId,
        selectedModelId: _selectedModelId,
        onModelSelected: _handleModelSelected,
        isStreaming: _isStreaming,
        visibleWidth: Style.messageViewWidth,
        onStopGeneration: _handleStopGeneration,
        messageToolBar: MessageToolbar(
          onCopy: _handleCopyMessage,
          onCopyPlainText: _handleCopyPlainText,
          onRegenerate: _handleRegenerateMessage,
          onDelete: _handleDeleteMessage,
          isVisible: !_isStreaming, // 流式生成过程中隐藏工具栏
        ),
      ),
      rightChild: sideWidget,
      leftWidth: 500.0,
      minLeftWidth: 400.0,
      dividerWidth: 1.0,
    );
  }

  Widget? _buildSide() {
    // 当需要显示侧边栏时，返回具体的 widget
    // 当不需要显示侧边栏时，返回 null
    return null; // 目前返回 null，只显示对话界面
  }
}
