import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/input_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_toolbar.dart';
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
    this.selectedProviderId,
    this.selectedModelId,
    this.onModelSelected,
    this.isStreaming = false,
    this.visibleWidth = double.infinity,
    this.messageToolBar
  });

  final Function(String)? onFileSelected;
  final Function(String)? onSend;
  final List<Message> historyMessages;
  final VoidCallback? onNewConversation;
  final String currentTitle;
  final String? selectedProviderId;
  final String? selectedModelId;
  final Function(String providerId, String modelId)? onModelSelected;
  final bool isStreaming;
  final double visibleWidth;
  final MessageToolbar? messageToolBar;

  @override
  State<StatefulWidget> createState() => _ChatView();
}

class _ChatView extends State<ChatView> {
  @override
  Widget build(BuildContext context) {
    final displayTitle = widget.currentTitle.isEmpty ? S.of(context).aiAssistant : widget.currentTitle;
    
    return ConstrainedBox(
      constraints: const BoxConstraints(
        minWidth: 300.0,
      ),
      child: Column(
        children: [
          // 顶部部件
          SizedBox(
            child: TitleBarView(
              title: displayTitle,
              onAddTap: widget.onNewConversation,
              visibleWidth: widget.visibleWidth,
            ),
          ),
          Expanded(
            child: MessageView(
              widget.historyMessages,
              isStreaming: widget.isStreaming,
              visibleWidth: widget.visibleWidth,
              messageToolBar: widget.messageToolBar,
            ),
          ),

          // 底部部件
          SizedBox(
            width: widget.visibleWidth,
            child: InputView(
              onFileSelected: widget.onFileSelected,
              onSend: (msg) {
                widget.onSend?.call(msg);
              },
              selectedProviderId: widget.selectedProviderId,
              selectedModelId: widget.selectedModelId,
              onModelSelected: widget.onModelSelected,
              isStreaming: widget.isStreaming,
            ),
          ),
        ],
      ),
    );
  }
}
