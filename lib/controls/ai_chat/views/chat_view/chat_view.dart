import 'package:flutter/material.dart';
import 'package:flutter_ai_providers/flutter_ai_providers.dart';
import 'package:flutter_ai_toolkit/flutter_ai_toolkit.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/input_view.dart';

class ChatView extends StatefulWidget {
  const ChatView({super.key, this.onFileSelected, this.onSend});

  final Function(String)? onFileSelected;
  final Function(String)? onSend;

  @override
  State<StatefulWidget> createState() => _ChatView();
}

class _ChatView extends State<ChatView> {
  final double _defaultChatView = 400;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: _defaultChatView,
      child: Column(
        children: [
          // 顶部部件
          Container(height: 50, color: Colors.red),

          Expanded(child: Container(color: Colors.green)),

          // 底部部件
          SizedBox(child: InputView(onFileSelected: widget.onFileSelected, onSend: widget.onSend)),
        ],
      ),
    );
  }
}
