import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/chat_view.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/controls/resizable_divider.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/rpc/service.pb.dart' as grpc_service;
import 'package:lemon_tea/rpc/common.pb.dart' as grpc_common;
import 'package:lemon_tea/rpc/common.pbenum.dart' as grpc_enum;
import 'package:lemon_tea/storage/chat_storage.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'dart:async';
import 'dart:io' show Platform;

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

  @override
  void initState() {
    super.initState();
    _conversationManager = widget.conversationManager ?? ConversationManager();
    // 监听 ConversationManager 的变化
    _conversationManager.addListener(_onConversationManagerChanged);
    _initializeConversation();
  }

  @override
  void dispose() {
    _conversationManager.removeListener(_onConversationManagerChanged);
    super.dispose();
  }

  void _onConversationManagerChanged() {
    // 当 ConversationManager 发生变化时，更新UI
    if (mounted) {
      setState(() {
        _historyMessages = _conversationManager.messages;
        _currentTitle = _conversationManager.currentConversation?.title ?? 'AI 助手';
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
      
      // 获取所有对话
      final conversations = await ChatStorage.getAllConversations();
      
      if (conversations.isEmpty) {
        // 如果没有对话，创建一个新的欢迎对话
        await _createWelcomeConversation();
      } else {
        // 加载最新的对话
        await _loadConversation(conversations.first);
      }
    } catch (e) {
      debugPrint('Failed to initialize conversation: $e');
      // 如果初始化失败，创建一个本地欢迎消息
      _historyMessages = [
        Message(
          role: MessageRole.assistant,
          content: """# 欢迎使用 Markdown

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
    print("代码高亮示例")""",
        ),
      ];
    } finally {
      setState(() {
        _isLoading = false;
      });
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
    final dbMessages = await ChatStorage.getMessagesByConversationId(conversation.id);
    
    // 转换数据库Message为LLM Message
    final messages = dbMessages.map((dbMsg) => Message(
      role: dbMsg.role,
      content: dbMsg.content,
    )).toList();
    
    setState(() {
      _currentTitle = conversation.title;
      _historyMessages = messages;
      _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
      _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
    });
  }

  Future<void> _createWelcomeConversation() async {
    final conversation = await ChatStorage.createConversation(
      title: '欢迎对话',
      defaultProviderId: 'deepseek',
      defaultModelId: 'deepseek-chat',
    );
    
    if (conversation != null) {
      _currentConversation = conversation;
      
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
      
      await ChatStorage.addMessage(
        conversationId: conversation.id,
        role: 'assistant',
        content: welcomeContent,
      );
      
      setState(() {
        _currentTitle = conversation.title;
        _historyMessages = [Message(
          role: MessageRole.assistant,
          content: welcomeContent,
        )];
        _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
        _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
      });
    }
  }

  Future<void> _handleSendMessage(String message) async {
    if (message.trim().isEmpty || _currentConversation == null) return;

    // 添加用户消息到界面
    final userMessage = Message(role: MessageRole.user, content: message);
    setState(() {
      _historyMessages.add(userMessage);
    });

    // 保存用户消息到数据库
    await ChatStorage.addMessage(
      conversationId: _currentConversation!.id,
      role: 'user',
      content: message,
    );

    // 如果是第一条用户消息，更新对话标题
    final messageCount = await ChatStorage.getMessageCountByConversationId(_currentConversation!.id);
    if (messageCount == 2) { // 1条欢迎消息 + 1条用户消息
      final title = _generateTitleFromMessage(message);
      await ChatStorage.updateConversationTitle(_currentConversation!.id, title);
      setState(() {
        _currentTitle = title;
      });
    }

    try {
      // 检查gRPC客户端
      if (_grpcClient.stub == null) {
        await _grpcClient.init();
        if (_grpcClient.stub == null) {
          throw Exception("gRPC客户端初始化失败，请检查服务是否启动");
        }
      }
      
      // 准备历史消息
      List<grpc_common.Message> grpcMessages = _historyMessages.map((msg) {
        return grpc_common.Message(
          role: msg.role == MessageRole.user 
              ? grpc_enum.MessageRole.MESSAGE_ROLE_USER 
              : grpc_enum.MessageRole.MESSAGE_ROLE_ASSISTANT,
          content: msg.content,
        );
      }).toList();
      
      // 创建聊天请求
      final providerId = _currentConversation!.defaultProviderId ?? '3c64dc4d-ffa7-408f-be2b-91f1bb150e82';
      final modelId = _currentConversation!.defaultModelId ?? 'deepseek-chat';
      
      debugPrint('使用模型配置: providerId=$providerId, modelId=$modelId');
      
      final chatRequest = grpc_service.ChatRequest(
        llmProviderId: providerId,
        modelId: modelId,
        historyMessages: grpcMessages,
        message: grpc_common.Message(
          role: grpc_enum.MessageRole.MESSAGE_ROLE_USER,
          content: message,
        ),
        prompt: "你是一个有用的AI助手。",
      );
      
      // 创建请求流
      final controller = StreamController<grpc_service.ChatRequest>();
      controller.add(chatRequest);
      controller.close();
      
      // 调用gRPC聊天接口
      final responseStream = _grpcClient.stub!.chat(controller.stream);
      
      String fullResponse = '';
      await for (final response in responseStream) {
        if (response.errorMessage.isNotEmpty) {
          throw Exception(response.errorMessage);
        }
        
        fullResponse += response.content;
        
        if (response.isDone) {
          // 聊天完成，添加AI回复到界面
          final aiMessage = Message(
            role: MessageRole.assistant,
            content: fullResponse,
          );
          setState(() {
            _historyMessages.add(aiMessage);
          });
          
          // 保存AI回复到数据库
          await ChatStorage.addMessage(
            conversationId: _currentConversation!.id,
            role: 'assistant',
            content: fullResponse,
          );
          break;
        }
      }
    } catch (e) {
      debugPrint('Error during chat: $e');
      
      // 发生错误时，添加错误消息到界面
      final errorMessage = Message(
        role: MessageRole.assistant,
        content: '抱歉，处理您的请求时发生错误：${e.toString()}',
      );
      setState(() {
        _historyMessages.add(errorMessage);
      });
      
      // 保存错误消息到数据库
      await ChatStorage.addMessage(
        conversationId: _currentConversation!.id,
        role: 'assistant',
        content: '抱歉，处理您的请求时发生错误：${e.toString()}',
      );
    }
  }

  Future<void> _handleNewConversation() async {
    final conversation = await ChatStorage.createConversation(
      title: '新对话',
      defaultProviderId: _selectedProviderId ?? 'deepseek',
      defaultModelId: _selectedModelId ?? 'deepseek-chat',
    );
    
    if (conversation != null) {
      setState(() {
        _currentConversation = conversation;
        _currentTitle = conversation.title;
        _historyMessages = [];
        _selectedProviderId = conversation.defaultProviderId ?? 'deepseek';
        _selectedModelId = conversation.defaultModelId ?? 'deepseek-chat';
      });
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

  @override
  Widget build(BuildContext context) {
    return ResizableDivider(
      leftChild: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : ChatView(
              historyMessages: _historyMessages,
              onSend: _handleSendMessage,
              onNewConversation: _historyMessages.isEmpty ? null : _handleNewConversation,
              currentTitle: _currentTitle,
              selectedProviderId: _selectedProviderId,
              selectedModelId: _selectedModelId,
              onModelSelected: _handleModelSelected,
            ),
      rightChild: Text("data"),
      leftWidth: 500.0,
      minLeftWidth: 400.0,
      maxLeftWidth: 800.0,
      dividerWidth: 1.0,
    );
  }
}