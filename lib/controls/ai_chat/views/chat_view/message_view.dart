import 'package:flutter/material.dart';
import 'package:flutter_ai_providers/flutter_ai_providers.dart';
import 'package:flutter_ai_toolkit/flutter_ai_toolkit.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';

class MessageView extends StatefulWidget {
  
  const MessageView(List<Message> historyMessages, {super.key});

  @override
  State<StatefulWidget> createState() => _MessageView();
}

class _MessageView extends State<MessageView> {
  @override
  Widget build(BuildContext context) {
    return Text("data");
  }

}