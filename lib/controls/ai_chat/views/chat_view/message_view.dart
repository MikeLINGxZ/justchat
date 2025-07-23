import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_toolbar.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/style.dart';
import 'package:markdown_widget/markdown_widget.dart';
import 'package:flutter/services.dart';

class MessageView extends ConsumerStatefulWidget {
  final List<Message> historyMessages;
  final bool isStreaming;
  final double visibleWidth;
  final MessageToolbar? messageToolBar;

  const MessageView(
    this.historyMessages, {
    super.key,
    this.isStreaming = false,
    this.visibleWidth = double.infinity,
    this.messageToolBar
  });

  @override
  ConsumerState<MessageView> createState() => _MessageViewState();
}

class _MessageViewState extends ConsumerState<MessageView> {
  final ScrollController _scrollController = ScrollController();
  int _lastMessageCount = 0;
  String _lastMessageContent = '';
  final Map<int, bool> _reasoningExpanded = {}; // 跟踪每个消息的思考过程展开状态
  final Map<int, bool> _messageHovered = {}; // 跟踪每个消息的悬停状态

  @override
  void initState() {
    super.initState();
    _lastMessageCount = widget.historyMessages.length;
    _lastMessageContent = widget.historyMessages.isNotEmpty 
        ? widget.historyMessages.last.content 
        : '';

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
                      fontSize: FontSizeUtils.getSmallSize(ref),
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
                                fontSize: FontSizeUtils.getXSmallSize(ref),
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
            fontSize: FontSizeUtils.getXSmallSize(ref),
            color: Colors.orange.shade800,
            height: 1.5,
          ),
        ),
        H1Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getTitleSize(ref),
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.bold,
          ),
        ),
        H2Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.w600,
          ),
        ),
        H3Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w600,
          ),
        ),
        H4Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H5Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H6Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
       ],
    );
  }

  // 动态创建 Markdown 配置
  MarkdownConfig _buildLightConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(textStyle: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
        H1Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleLargeSize(ref) + 6, height: 38 / 30),
        ),
        H2Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleSize(ref), height: 28 / 22),
        ),
        H3Config(
          style: TextStyle(fontSize: FontSizeUtils.getSubheadingSize(ref), height: 23 / 18),
        ),
        H4Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H5Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H6Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
      ],
    );
  }

  MarkdownConfig _buildDarkConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(textStyle: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
        H1Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleLargeSize(ref) + 6, height: 38 / 30),
        ),
        H2Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleSize(ref), height: 28 / 22),
        ),
        H3Config(
          style: TextStyle(fontSize: FontSizeUtils.getSubheadingSize(ref), height: 23 / 18),
        ),
        H4Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H5Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H6Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: ListView.builder(
        controller: _scrollController,
        itemCount: widget.historyMessages.length,
        itemBuilder: (context, index) {
          final message = widget.historyMessages[index];
          final isLastMessage = index == widget.historyMessages.length - 1;
          final isStreamingMessage = isLastMessage &&
              message.role == MessageRole.assistant &&
              widget.isStreaming;

          return Center(
            child: SizedBox(
              width: widget.visibleWidth,
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 20.0),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (message.role != MessageRole.user) ...[
                      CircleAvatar(
                        radius: 20,
                        backgroundColor: Colors.green,
                        child: Text(
                          'A',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: FontSizeUtils.getSubheadingSize(ref),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Flexible(
                        child: MouseRegion(
                          onEnter: (_) => setState(() => _messageHovered[index] = true),
                          onExit: (_) => setState(() => _messageHovered[index] = false),
                          child: Container(
                            constraints: const BoxConstraints(
                              minWidth: 0,
                            ),
                            padding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 2,
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
                                  color: Style.assistantChatBubble(context),
                                  child: MarkdownBlock(
                                    data: message.content.isEmpty ? ' ' : message.content,
                                    config: Theme.of(context).brightness == Brightness.dark
                                        ? _buildDarkConfig()
                                        : _buildLightConfig(),
                                  ),
                                ),
                                if (message.role == MessageRole.assistant && widget.messageToolBar != null)
                                  widget.messageToolBar!.setMessage(message, visible: _messageHovered[index] ?? false),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                    if (message.role == MessageRole.user) ...[
                      const Spacer(),
                      Container(
                        constraints: const BoxConstraints(
                          minWidth: 0,
                        ),
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color:  Style.userChatBubble(context),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: MarkdownBlock(
                          data: message.content,
                          config: Theme.of(context).brightness == Brightness.dark
                              ? _buildDarkConfig()
                              : _buildLightConfig(),
                        ),
                      ),
                      const SizedBox(width: 12),
                      CircleAvatar(
                        radius: 20,
                        backgroundColor: Colors.blue,
                        child: Text(
                          'U',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: FontSizeUtils.getSubheadingSize(ref),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          );
        },
      ),
    );
  }


}

