import 'package:flutter/material.dart';
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
  List<Message> _previousMessages = [];

  @override
  void initState() {
    super.initState();
    _previousMessages = List.from(widget.historyMessages);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  @override
  void didUpdateWidget(MessageView oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.historyMessages != _previousMessages) {
      _previousMessages = List.from(widget.historyMessages);
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (_scrollController.hasClients) {
          _scrollController.animateTo(
            _scrollController.position.maxScrollExtent,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOut,
          );
        }
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
            children: [
              CircleAvatar(
                radius: 14,
                backgroundColor:
                    message.role == 'user' ? Colors.blue : Colors.green,
                child: Text(
                  message.role == 'user' ? 'U' : 'A',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color:
                        message.role == 'user'
                            ? Colors.blue.withOpacity(0.1)
                            : Colors.green.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: MarkdownBlock(
                    data: message.content,
                    config:
                        Theme.of(context).brightness == Brightness.dark
                            ? MarkdownConfig.darkConfig
                            : MarkdownConfig.defaultConfig,
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
