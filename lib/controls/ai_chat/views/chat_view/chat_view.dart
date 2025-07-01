import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/input_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/title_bar_view.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/generated/l10n.dart';

class ChatView extends StatefulWidget {
  const ChatView({
    super.key, 
    this.onFileSelected, 
    this.onSend, 
    required this.historyMessages,
    this.onNewConversation,
    this.currentTitle = '',
  });

  final Function(String)? onFileSelected;
  final Function(String)? onSend;
  final List<Message> historyMessages;
  final VoidCallback? onNewConversation;
  final String currentTitle;

  @override
  State<StatefulWidget> createState() => _ChatView();
}

class _ChatView extends State<ChatView> {
  @override
  Widget build(BuildContext context) {
    final displayTitle = widget.currentTitle.isEmpty ? S.of(context).aiAssistant : widget.currentTitle;
    
    return ConstrainedBox(
      constraints: const BoxConstraints(
        minWidth: 400.0,
      ),
      child: Column(
        children: [
          // 顶部部件
          SizedBox(
            child: TitleBarView(
              title: displayTitle,
              onAddTap: widget.onNewConversation,
            ),
          ),

          Expanded(child: MessageView(widget.historyMessages)),

          // 底部部件
          SizedBox(child: InputView(onFileSelected: widget.onFileSelected, onSend: (msg) {
            widget.onSend?.call(msg);
          })),
        ],
      ),
    );
  }
}
