import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/chat_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/history_dialog.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/controls/resizable_divider.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';
import 'package:lemon_tea/models/conversation.dart';

class AssistantPage extends StatefulWidget {
  const AssistantPage({super.key});

  @override
  State<StatefulWidget> createState() => _AssistantPage();
}

class _AssistantPage extends State<AssistantPage> {
  late ConversationManager _conversationManager;
  List<Message> _historyMessages = [];
  bool _isLoading = false;
  String _currentTitle = 'AI 助手';

  @override
  void initState() {
    super.initState();
    _conversationManager = ConversationManager();
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
        _historyMessages = _conversationManager.currentConversation?.messages ?? [];
        _currentTitle = _conversationManager.currentConversation?.title ?? 'AI 助手';
      });
    }
  }

  Future<void> _initializeConversation() async {
    setState(() {
      _isLoading = true;
    });

    try {
      await _conversationManager.initialize();
      
      // 如果没有对话，创建一个新的欢迎对话
      if (_conversationManager.conversations.isEmpty) {
        await _createWelcomeConversation();
      } else {
        // 加载最新的对话
        final latestConversation = _conversationManager.conversations.first;
        await _conversationManager.loadConversation(latestConversation.id);
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

  Future<void> _createWelcomeConversation() async {
    final welcomeMessage = Message(
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
    );

    await _conversationManager.createConversation(
      title: '欢迎对话',
      initialMessages: [welcomeMessage],
    );
  }

  Future<void> _handleSendMessage(String message) async {
    if (message.trim().isEmpty) return;

    final userMessage = Message(role: MessageRole.user, content: message);
    
    // 保存到存储
    await _conversationManager.addMessageToCurrent(userMessage);

    // 如果是第一条用户消息，更新对话标题
    if (_conversationManager.currentConversation?.messages.length == 2) {
      final title = _conversationManager.generateTitleFromMessage(message);
      await _conversationManager.updateConversationTitle(
        _conversationManager.currentConversation!.id,
        title,
      );
    }

    // 模拟AI回复（这里可以集成真实的AI服务）
    await Future.delayed(const Duration(milliseconds: 500));
    
    final aiMessage = Message(
      role: MessageRole.assistant,
      content: '我收到了你的消息："$message"。这是一个模拟的AI回复。',
    );

    // 保存到存储
    await _conversationManager.addMessageToCurrent(aiMessage);
  }

  Future<void> _handleNewConversation() async {
    await _conversationManager.createConversation(title: '新对话');
  }

  Future<void> _handleLoadConversation(Conversation conversation) async {
    await _conversationManager.loadConversation(conversation.id);
  }

  Future<void> _handleDeleteConversation(String conversationId) async {
    // 保存context引用
    final currentContext = context;
    
    // 显示确认对话框
    final shouldDelete = await showDialog<bool>(
      context: currentContext,
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: const Text('确定要删除这个对话吗？删除后无法恢复。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: TextButton.styleFrom(
              foregroundColor: Colors.red,
            ),
            child: const Text('删除'),
          ),
        ],
      ),
    );

    if (shouldDelete != true) return;

    await _conversationManager.deleteConversation(conversationId);
    
    // 如果删除的是当前对话，加载最新的对话或创建新对话
    if (_conversationManager.currentConversation?.id == conversationId) {
      if (_conversationManager.conversations.isNotEmpty) {
        await _handleLoadConversation(_conversationManager.conversations.first);
      } else {
        await _handleNewConversation();
      }
    }
    
    // 关闭对话框，让用户重新打开以查看更新后的列表
    if (Navigator.of(currentContext).canPop()) {
      Navigator.of(currentContext).pop();
    }
  }

  void _showHistoryDialog() {
    showDialog(
      context: context,
      builder: (context) => HistoryDialog(
        conversationManager: _conversationManager,
        onConversationSelected: _handleLoadConversation,
        onConversationDeleted: _handleDeleteConversation,
        onNewConversation: _handleNewConversation,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ResizableDivider(
      leftChild: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : ChatView(
              historyMessages: _historyMessages,
              onSend: _handleSendMessage,
              onNewConversation: _handleNewConversation,
              onHistoryTap: _showHistoryDialog,
              currentTitle: _currentTitle,
            ),
      rightChild: Text("data"),
      leftWidth: 500.0,
      minLeftWidth: 400.0,
      maxLeftWidth: 800.0,
      dividerWidth: 1.0,
    );
  }
}
