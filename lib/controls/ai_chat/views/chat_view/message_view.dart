import 'package:flutter/material.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:markdown_widget/markdown_widget.dart';

class MessageView extends StatefulWidget {
  final List<Message> historyMessages;
  const MessageView(this.historyMessages, {super.key});

  @override
  State<StatefulWidget> createState() => _MessageViewState();
}

class _MessageViewState extends State<MessageView> {
  final ScrollController _scrollController = ScrollController();
  int _lastMessageCount = 0;

  // 自定义 Markdown 配置，所有文字大小减少2
  late final MarkdownConfig _customLightConfig;
  late final MarkdownConfig _customDarkConfig;

  @override
  void initState() {
    super.initState();
    _lastMessageCount = widget.historyMessages.length;

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
      _scrollController.jumpTo(_scrollController.position.maxScrollExtent);
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  @override
  void didUpdateWidget(MessageView oldWidget) {
    super.didUpdateWidget(oldWidget);

    // 当有新消息到达时滚动到底部
    if (widget.historyMessages.length > _lastMessageCount) {
      _lastMessageCount = widget.historyMessages.length;
      Future.delayed(Duration(milliseconds: 300), () {
        _scrollController.jumpTo(_scrollController.position.maxScrollExtent);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      controller: _scrollController,
      itemCount: widget.historyMessages.length,
      itemBuilder: (context, index) {
        final message = widget.historyMessages[index];
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0, horizontal: 16.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            // mainAxisAlignment: message.role == MessageRole.user ? MainAxisAlignment.start : MainAxisAlignment.end,
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
                      // maxWidth: 400, // 限制最大宽度，避免过宽
                    ),
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 8,
                    ),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: MarkdownBlock(
                      data: message.content,
                      config:
                          Theme.of(context).brightness == Brightness.dark
                              ? _customDarkConfig
                              : _customLightConfig,
                    ),
                  ),
                ),
              ],
              if (message.role == MessageRole.user) ...[
                const Spacer(),
                Container(
                  constraints: const BoxConstraints(
                    minWidth: 0,
                    maxWidth: 300, // 限制最大宽度，避免过宽
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
                    config:
                        Theme.of(context).brightness == Brightness.dark
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
