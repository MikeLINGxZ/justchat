import 'package:flutter/material.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:markdown_widget/markdown_widget.dart';

class MessageView extends StatefulWidget {
  final List<Message> historyMessages;
  final bool isStreaming;
  
  const MessageView(
    this.historyMessages, {
    super.key,
    this.isStreaming = false,
  });

  @override
  State<StatefulWidget> createState() => _MessageViewState();
}

class _MessageViewState extends State<MessageView> {
  final ScrollController _scrollController = ScrollController();
  int _lastMessageCount = 0;
  String _lastMessageContent = '';
  final Map<int, bool> _reasoningExpanded = {}; // 跟踪每个消息的思考过程展开状态

  // 自定义 Markdown 配置，所有文字大小减少2
  late final MarkdownConfig _customLightConfig;
  late final MarkdownConfig _customDarkConfig;

  @override
  void initState() {
    super.initState();
    _lastMessageCount = widget.historyMessages.length;
    _lastMessageContent = widget.historyMessages.isNotEmpty 
        ? widget.historyMessages.last.content 
        : '';

    // 创建自定义配置
    _customLightConfig = MarkdownConfig(
      configs: [
        PConfig(textStyle: const TextStyle(fontSize: 14)), // 13 -> 14
        H1Config(
          style: const TextStyle(fontSize: 30, height: 38 / 30),
        ), // 29 -> 30
        H2Config(
          style: const TextStyle(fontSize: 22, height: 28 / 22),
        ), // 21 -> 22
        H3Config(
          style: const TextStyle(fontSize: 18, height: 23 / 18),
        ), // 17 -> 18
        H4Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
        H5Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
        H6Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
      ],
    );

    _customDarkConfig = MarkdownConfig(
      configs: [
        PConfig(textStyle: const TextStyle(fontSize: 14)), // 13 -> 14
        H1Config(
          style: const TextStyle(fontSize: 30, height: 38 / 30),
        ), // 29 -> 30
        H2Config(
          style: const TextStyle(fontSize: 22, height: 28 / 22),
        ), // 21 -> 22
        H3Config(
          style: const TextStyle(fontSize: 18, height: 23 / 18),
        ), // 17 -> 18
        H4Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
        H5Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
        H6Config(
          style: const TextStyle(fontSize: 14, height: 18 / 14),
        ), // 13 -> 14
      ],
    );

    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollToBottom();
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 100),
          curve: Curves.easeOut,
        );
      });
    }
  }

  @override
  void didUpdateWidget(MessageView oldWidget) {
    super.didUpdateWidget(oldWidget);

    final currentMessageCount = widget.historyMessages.length;
    final currentLastContent = widget.historyMessages.isNotEmpty 
        ? widget.historyMessages.last.content 
        : '';

    // 检查是否有新消息或最后一条消息内容发生变化（流式更新）
    if (currentMessageCount > _lastMessageCount || 
        (_lastMessageContent != currentLastContent && currentLastContent.isNotEmpty)) {
      
      _lastMessageCount = currentMessageCount;
      _lastMessageContent = currentLastContent;
      
      // 滚动到底部
      _scrollToBottom();
    }
  }

  Widget _buildReasoningSection(int index, String reasoningContent) {
    final isExpanded = _reasoningExpanded[index] ?? false;
    
    return Container(
      margin: const EdgeInsets.only(top: 12, bottom: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          GestureDetector(
            onTap: () {
              setState(() {
                _reasoningExpanded[index] = !isExpanded;
              });
              
              // 如果展开思考过程，延迟一下再滚动到底部
              if (!isExpanded) {
                Future.delayed(const Duration(milliseconds: 300), () {
                  _scrollToBottom();
                });
              }
            },
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Colors.orange.withOpacity(0.15),
                    Colors.orange.withOpacity(0.08),
                  ],
                  begin: Alignment.centerLeft,
                  end: Alignment.centerRight,
                ),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: Colors.orange.withOpacity(0.4),
                  width: 1.5,
                ),
                boxShadow: [
                  BoxShadow(
                    color: Colors.orange.withOpacity(0.1),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      color: Colors.orange.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Icon(
                      Icons.psychology_outlined,
                      size: 16,
                      color: Colors.orange.shade800,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    '思考过程',
                    style: TextStyle(
                      fontSize: 13,
                      color: Colors.orange.shade800,
                      fontWeight: FontWeight.w600,
                      letterSpacing: 0.5,
                    ),
                  ),
                  const SizedBox(width: 8),
                  AnimatedRotation(
                    turns: isExpanded ? 0.5 : 0,
                    duration: const Duration(milliseconds: 200),
                    child: Icon(
                      Icons.keyboard_arrow_down,
                      size: 18,
                      color: Colors.orange.shade700,
                    ),
                  ),
                ],
              ),
            ),
          ),
          AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeInOut,
            height: isExpanded ? null : 0,
            child: isExpanded ? Column(
              children: [
                const SizedBox(height: 12),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      colors: [
                        Colors.orange.withOpacity(0.03),
                        Colors.orange.withOpacity(0.08),
                      ],
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                    ),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.orange.withOpacity(0.3),
                      width: 1,
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.orange.withOpacity(0.05),
                        blurRadius: 8,
                        offset: const Offset(0, 2),
                        spreadRadius: 1,
                      ),
                    ],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // 思考过程标题栏
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.orange.withOpacity(0.15),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              Icons.lightbulb_outline,
                              size: 14,
                              color: Colors.orange.shade700,
                            ),
                            const SizedBox(width: 6),
                            Text(
                              'AI 思考过程',
                              style: TextStyle(
                                fontSize: 11,
                                color: Colors.orange.shade700,
                                fontWeight: FontWeight.w500,
                                letterSpacing: 0.3,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                      // 思考过程内容
                      MarkdownBlock(
                        data: reasoningContent,
                        config: _buildReasoningMarkdownConfig(),
                      ),
                    ],
                  ),
                ),
              ],
            ) : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }

  // 为思考过程定制的Markdown配置
  MarkdownConfig _buildReasoningMarkdownConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(
          textStyle: TextStyle(
            fontSize: 13,
            color: Colors.orange.shade800,
            height: 1.5,
          ),
        ),
        H1Config(
          style: TextStyle(
            fontSize: 24,
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.bold,
          ),
        ),
        H2Config(
          style: TextStyle(
            fontSize: 20,
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.w600,
          ),
        ),
        H3Config(
          style: TextStyle(
            fontSize: 16,
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w600,
          ),
        ),
        H4Config(
          style: TextStyle(
            fontSize: 14,
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H5Config(
          style: TextStyle(
            fontSize: 13,
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H6Config(
          style: TextStyle(
            fontSize: 13,
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
                 ),
       ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      controller: _scrollController,
      itemCount: widget.historyMessages.length,
      itemBuilder: (context, index) {
        final message = widget.historyMessages[index];
        final isLastMessage = index == widget.historyMessages.length - 1;
        final isStreamingMessage = isLastMessage && 
                                  message.role == MessageRole.assistant && 
                                  widget.isStreaming;
        
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0, horizontal: 16.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (message.role != MessageRole.user) ...[
                CircleAvatar(
                  radius: 20,
                  backgroundColor: Colors.green,
                  child: const Text(
                    'A',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Flexible(
                  child: Container(
                    constraints: const BoxConstraints(
                      minWidth: 0,
                    ),
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 8,
                    ),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // 显示思考过程（如果有的话）
                        if (message.reasoningContent != null && 
                            message.reasoningContent!.isNotEmpty)
                          AnimatedSwitcher(
                            duration: const Duration(milliseconds: 200),
                            child: _buildReasoningSection(index, message.reasoningContent!),
                          ),
                        
                        // 显示主要内容
                        Container(
                          width: double.infinity,
                          child: MarkdownBlock(
                            data: message.content.isEmpty ? ' ' : message.content,
                            config: Theme.of(context).brightness == Brightness.dark
                                ? _customDarkConfig
                                : _customLightConfig,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
              if (message.role == MessageRole.user) ...[
                const Spacer(),
                Container(
                  constraints: const BoxConstraints(
                    minWidth: 0,
                    maxWidth: 300,
                  ),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 8,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.blue.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: MarkdownBlock(
                    data: message.content,
                    config: Theme.of(context).brightness == Brightness.dark
                        ? _customDarkConfig
                        : _customLightConfig,
                  ),
                ),
                const SizedBox(width: 12),
                CircleAvatar(
                  radius: 20,
                  backgroundColor: Colors.blue,
                  child: const Text(
                    'U',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ],
          ),
        );
      },
    );
  }


}
